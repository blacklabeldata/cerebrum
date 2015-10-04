package cerebrum

import (
	"os"
	"path/filepath"
	"time"

	"github.com/blacklabeldata/serfer"
	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb"
	"github.com/hashicorp/serf/serf"
	log "github.com/mgutz/logxi/v1"
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
		// logger.Warn("Could not create data directory", "err", err.Error())
		return
	}

	// TODO: Replace with Raft IsLeader check
	isLeader := func() bool {
		return true
	}

	// Setup reconciler
	reconcilerCh := make(chan serf.Member, 32)
	reconciler := &Reconciler{reconcilerCh, isLeader}

	serfEventCh := make(chan serf.Event, 256)
	serfer := serfer.NewSerfer(serfEventCh, serfer.SerfEventHandler{
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
	})
	return &cerebrum{
		config: c,
		serfer: serfer,
	}, nil
}

type Cerebrum interface {
	Start() error
	Stop() error
}

type cerebrum struct {
	id          string
	serfEventCh chan serf.Event
	config      *Config
	logger      log.Logger
	serfer      serfer.Serfer

	// The raft instance is used among Consul nodes within the
	// DC to protect operations that require strong consistency
	raft          *raft.Raft
	raftPeers     raft.PeerStore
	raftStore     *raftboltdb.BoltStore
	raftTransport *raft.NetworkTransport
	fsm           raft.FSM
}

func (c *cerebrum) Start() error {
	return nil
}

func (c *cerebrum) Stop() error {
	return nil
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

	conf.Tags["id"] = c.id
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

	// // Create the FSM
	// var err error
	// c.fsm, err = NewFSM(statePath, c.config.LogOutput)
	// if err != nil {
	// 	return err
	// }

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

	// Create TCP transport
	trans, err := NewTLSTransport("0.0.0.0:9122", 3, 10*time.Second, c.config.LogOutput, c.config.TLSConfig)
	if err != nil {
		return err
	}

	// // Create a transport layer
	// trans := raft.NewNetworkTransport(t, 3, 10*time.Second, c.config.LogOutput)
	c.raftTransport = trans

	// Setup the peer store
	c.raftPeers = raft.NewJSONPeers(path, trans)

	// Ensure local host is always included if we are in bootstrap mode
	if c.config.Bootstrap {
		peers, err := c.raftPeers.Peers()
		if err != nil {
			store.Close()
			return err
		}
		if !raft.PeerContained(peers, trans.LocalAddr()) {
			c.raftPeers.SetPeers(raft.AddUniquePeer(peers, trans.LocalAddr()))
		}
	}

	// Make sure we set the LogOutput
	c.config.RaftConfig.LogOutput = c.config.LogOutput

	// Setup the Raft store
	c.raft, err = raft.NewRaft(c.config.RaftConfig, c.fsm, cacheStore, store,
		snapshots, c.raftPeers, trans)
	if err != nil {
		store.Close()
		trans.Close()
		return err
	}

	// // Start monitoring leadership
	// c.t.Go(func() error {
	// 	c.monitorLeadership()
	// 	return nil
	// })
	return nil
}
