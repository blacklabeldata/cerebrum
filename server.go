package cerebrum

import (
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/blacklabeldata/grim"
	"github.com/blacklabeldata/serfer"
	"github.com/blacklabeldata/yamuxer"
	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb"
	"github.com/hashicorp/serf/serf"
	log "github.com/mgutz/logxi/v1"
	"golang.org/x/net/context"
	// tomb "gopkg.in/tomb.v2"
)

// Constants
const (
	CerebrumServiceID   = "cerebrum"
	CerebrumRole        = "cerebrum-server"
	CerebrumEventPrefix = CerebrumServiceID + ":"
	CerebrumLeaderEvent = CerebrumEventPrefix + "new-leader"
)

const (
	SerfSnapshotDir   = "serf/local.snapshot"
	RaftStateDir      = "raft/"
	tmpStatePath      = "tmp/"
	SnapshotsRetained = 2

	// raftLogCacheSize is the maximum number of logs to cache in-memory.
	// This is used to reduce disk I/O for the recently committed entries.
	raftLogCacheSize = 512
)

func New(c *Config) (cer Cerebrum, err error) {

	// Create logger
	if c.LogOutput == nil {
		c.LogOutput = log.NewConcurrentWriter(os.Stderr)
	}
	logger := log.NewLogger(c.LogOutput, "kappa")

	// Create data directory
	if err = os.MkdirAll(c.DataPath, 0755); err != nil {
		logger.Warn("Could not create data directory", "err", err)
		return
	}

	// Setup reconciler
	serfEventCh := make(chan serf.Event, 256)
	reconcilerCh := make(chan serf.Member, 32)

	ctx, cancel := context.WithCancel(context.Background())
	cereb := &cerebrum{
		config:      c,
		logger:      logger,
		dialer:      NewDialer(NewPool(c.LogOutput, 5*time.Minute, c.TLSConfig)),
		serfEventCh: serfEventCh,
		reconcileCh: reconcilerCh,
		grim:        grim.ReaperWithContext(ctx),
		context:     ctx,
		cancel:      cancel,
	}

	// Create serf server
	err = cereb.setupRaft()
	if err != nil {
		err = logger.Error("Failed to start serf: %v", err)
		return nil, err
	}

	isLeader := func() bool { return cereb.raft.State() == raft.Leader }
	reconciler := &Reconciler{reconcilerCh, isLeader}
	cereb.serfer = serfer.NewSerfer(serfEventCh, serfer.SerfEventHandler{
		Logger:              log.NewLogger(c.LogOutput, CerebrumEventPrefix),
		ServicePrefix:       CerebrumEventPrefix,
		ReconcileOnJoin:     true,
		ReconcileOnLeave:    true,
		ReconcileOnFail:     true,
		ReconcileOnUpdate:   true,
		ReconcileOnReap:     true,
		NodeJoined:          c.NodeJoined,
		NodeUpdated:         c.NodeUpdated,
		NodeLeft:            c.NodeLeft,
		NodeFailed:          c.NodeFailed,
		NodeReaped:          c.NodeReaped,
		UserEvent:           c.UserEvent,
		UnknownEventHandler: c.UnknownEventHandler,
		Reconciler:          reconciler,
		IsLeader:            isLeader,
		IsLeaderEvent: func(name string) bool {
			return name == CerebrumLeaderEvent
		},
		LeaderElectionHandler: cereb,
	})

	// Create serf server
	cereb.serf, err = cereb.setupSerf()
	if err != nil {
		err = logger.Error("Failed to start serf: %v", err)
		return nil, err
	}

	cer = cereb
	return cer, nil
}

type Cerebrum interface {
	Start() error
	Stop()
}

type cerebrum struct {
	config *Config
	logger log.Logger

	// pool        *ConnPool
	dialer      Dialer
	serfEventCh chan serf.Event
	leader      string
	serf        *serf.Serf
	serfer      serfer.Serfer

	// The raft instance is used among Consul nodes within the
	// DC to protect operations that require strong consistency
	raft              *raft.Raft
	raftPeers         raft.PeerStore
	raftLayer         *RaftLayer
	raftStore         *raftboltdb.BoltStore
	raftTransport     *raft.NetworkTransport
	reconcileCh       chan serf.Member
	nodeStatusUpdater NodeStatusUpdater
	muxer             yamuxer.Yamuxer
	fsm               raft.FSM

	applier   Applier
	forwarder Forwarder

	// t       tomb.Tomb
	grim    grim.GrimReaper
	context context.Context
	cancel  context.CancelFunc
}

func (c *cerebrum) Start() error {

	// Start monitoring raft cluster
	go c.monitorLeadership()

	// Start serf handler
	c.serfer.Start()

	// Join serf cluster
	c.logger.Info("Joining cluster", "nodes", c.config.ExistingNodes)

	n, err := c.serf.Join(c.config.ExistingNodes, true)
	if err != nil && !c.config.Bootstrap {
		err = c.logger.Error("Failed to join cluster", "err", err)
		return err
	}
	c.logger.Info("Joined cluster", "nodes", n)

	// Start services
	ctx := Context{
		Context: c.context,
		Serf:    c.serf,
		Raft:    c.raft,
	}
	for _, svc := range c.config.Services {
		svc.Start(&ctx)
	}

	return nil
}

func (c *cerebrum) Stop() {
	c.cancel()
	for _, svc := range c.config.Services {
		svc.Stop()
	}

	// Shutdown serf
	c.logger.Info("Shutting down Serf server...")
	c.serf.Leave()
	c.serf.Shutdown()
	<-c.serf.ShutdownCh()

	// Stop serf event handlers
	if err := c.serfer.Stop(); err != nil {
		c.logger.Warn("error: stopping Serfer handlers", err.Error())
	}

	// c.listener.Close()
	c.muxer.Stop()
	c.dialer.Shutdown()
}

func (c *cerebrum) setupSerf() (*serf.Serf, error) {
	// Get serf config
	conf := serf.DefaultConfig()

	// Initialize serf
	conf.Init()

	conf.NodeName = c.config.NodeName
	conf.MemberlistConfig.BindAddr = c.config.GossipBindAddr
	conf.MemberlistConfig.BindPort = c.config.GossipBindPort
	conf.MemberlistConfig.AdvertiseAddr = c.config.GossipAdvertiseAddr
	conf.MemberlistConfig.AdvertisePort = c.config.GossipAdvertisePort
	c.logger.Info("Gossip",
		"BindAddr", conf.MemberlistConfig.BindAddr,
		"BindPort", conf.MemberlistConfig.BindPort,
		"AdvertiseAddr", conf.MemberlistConfig.AdvertiseAddr,
		"AdvertisePort", conf.MemberlistConfig.AdvertisePort)

	conf.Tags["id"] = c.config.NodeID
	conf.Tags["role"] = "cerebrum-server"
	conf.Tags["dc"] = c.config.DataCenter
	// conf.Tags["port"] = fmt.Sprintf("%d", port)

	conf.MemberlistConfig.LogOutput = c.config.LogOutput
	conf.LogOutput = c.config.LogOutput
	conf.EventCh = c.serfEventCh
	conf.SnapshotPath = filepath.Join(c.config.DataPath, SerfSnapshotDir)
	conf.ProtocolVersion = conf.ProtocolVersion
	conf.RejoinAfterLeave = true
	conf.EnableNameConflictResolution = false
	conf.Merge = &mergeDelegate{c.logger}
	if err := os.MkdirAll(conf.SnapshotPath, 0755); err != nil {
		return nil, err
	}
	return serf.Create(conf)
}

// setupRaft is used to setup and initialize Raft
func (c *cerebrum) setupRaft() error {

	// If we are in bootstrap mode, enable a single node cluster
	if c.config.Bootstrap {
		c.config.RaftConfig.EnableSingleNode = true
	}

	// Create the base state path
	statePath := filepath.Join(c.config.DataPath, tmpStatePath)
	if err := os.RemoveAll(statePath); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(statePath), 0755); err != nil {
		return err
	}

	// Create the base raft path
	path := filepath.Join(c.config.DataPath, RaftStateDir)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	// Create the backend raft store for logs and stable storage
	store, err := raftboltdb.NewBoltStore(filepath.Join(path, "raft.db"))
	if err != nil {
		return err
	}
	c.raftStore = store

	// Wrap the store in a LogCache to improve performance
	cacheStore, err := raft.NewLogCache(c.config.LogCacheSize, store)
	if err != nil {
		store.Close()
		return err
	}

	// Create the snapshot store
	snapshots, err := raft.NewFileSnapshotStore(path, c.config.SnapshotsRetained, c.config.LogOutput)
	if err != nil {
		store.Close()
		return err
	}

	// Try to bind
	addr, err := net.ResolveTCPAddr("tcp", c.config.RaftBindAddr)
	if err != nil {
		return err
	}

	// Start TCP listener
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return err
	}

	// Create connection layer and transport
	layer := NewRaftLayer(c.dialer, listener.Addr(), c.config.TLSConfig)
	c.raftTransport = raft.NewNetworkTransport(layer, 3, 10*time.Second, c.config.LogOutput)

	// Create TLS connection dispatcher
	dispatcher := yamuxer.NewDispatcher(log.NewLogger(c.config.LogOutput, "dispatcher"), nil)
	dispatcher.Register(connRaft, layer)
	dispatcher.Register(connForward, &ForwardingHandler{c.applier, log.NewLogger(c.config.LogOutput, "forwarder")})

	// Create TLS connection muxer
	c.muxer = yamuxer.New(c.context, &yamuxer.Config{
		Listener:   listener,
		TLSConfig:  c.config.TLSConfig,
		Deadline:   c.config.ConnectionDeadline,
		LogOutput:  c.config.LogOutput,
		Dispatcher: dispatcher,
	})

	// Setup the peer store
	c.raftPeers = raft.NewJSONPeers(path, c.raftTransport)

	// Ensure local host is always included if we are in bootstrap mode
	if c.config.Bootstrap {
		peers, err := c.raftPeers.Peers()
		if err != nil {
			store.Close()
			return err
		}
		if !raft.PeerContained(peers, c.raftTransport.LocalAddr()) {
			c.raftPeers.SetPeers(raft.AddUniquePeer(peers, c.raftTransport.LocalAddr()))
		}
	}

	// Make sure we set the LogOutput
	c.config.RaftConfig.LogOutput = c.config.LogOutput

	// Setup the Raft store
	c.raft, err = raft.NewRaft(c.config.RaftConfig, c.fsm, cacheStore, store,
		snapshots, c.raftPeers, c.raftTransport)
	if err != nil {
		store.Close()
		c.raftTransport.Close()
		return err
	}

	// Setup forwarding and applier
	c.forwarder = NewForwarder(c.raft, c.dialer, log.NewLogger(c.config.LogOutput, "forwarder"))
	c.applier = NewApplier(c.raft, c.forwarder, log.NewLogger(c.config.LogOutput, "applier"), c.config.EnqueueTimeout)

	c.nodeStatusUpdater = NewNodeStatusUpdater(c.applier, log.NewLogger(c.config.LogOutput, "node-status-updater"))
	return nil
}
