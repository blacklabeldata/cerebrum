package cerebrum

import (
	"crypto/tls"
	"io"
	"time"

	"github.com/blacklabeldata/serfer"
	"github.com/hashicorp/raft"
)

type Config struct {

	// Bootstrap allows for a single node to become the leader for Raft.
	Bootstrap bool

	// NodeID should be unique across all nodes in the cluster.
	NodeID string

	// NodeName is the name of this node.
	NodeName string

	// DataCenter is the name of the data center for this node.
	DataCenter string

	// ServicePrefix is used to filter out unknown events.
	ServicePrefix string

	// DataPath is where all the data is stored
	DataPath string

	// LeaderElectionHandler processes leader election events.
	LeaderElectionHandler serfer.LeaderElectionHandler

	// UserEvent processes known, non-leader election events.
	UserEvent serfer.UserEventHandler

	// UnknownEventHandler processes unkown events.
	UnknownEventHandler serfer.UnknownEventHandler

	// Called when a Member joins the cluster.
	NodeJoined serfer.MemberJoinHandler

	// Called when a Member leaves the cluster by sending a leave message.
	NodeLeft serfer.MemberLeaveHandler

	// Called when a Member has been detected as failed.
	NodeFailed serfer.MemberFailureHandler

	// Called when a Member has been Readed from the cluster.
	NodeReaped serfer.MemberReapHandler

	// Called when a Member has been updated.
	NodeUpdated serfer.MemberUpdateHandler

	// Called when a serf.Query is received.
	QueryHandler serfer.QueryEventHandler

	// GossipBindAddr is the address of the Serf server.
	GossipBindAddr string

	// GossipBindPort is the port for the Serf server.
	GossipBindPort int

	// GossipAdvertiseAddr is the advertising address for the Serf server.
	GossipAdvertiseAddr string

	// GossipAdvertisePort is the advertising port for the Serf server.
	GossipAdvertisePort int

	// LogOutput is the output for all logs.
	LogOutput io.Writer

	// RaftConfig configures the Raft server.
	RaftConfig *raft.Config

	// SnapshotsRetained is the number of snapshots kept for Raft
	SnapshotsRetained int

	// LogCacheSize is the number of log entries to keep in memory.
	LogCacheSize int

	// TLSConfig is the config for Raft over TLS
	TLSConfig *tls.Config

	// RaftBindAddr is the bind address and port for the Raft TLS server.
	RaftBindAddr string

	// ReconcileInterval is the interval at which Raft makes sure the FSM has caught up.
	ReconcileInterval time.Duration

	// ConnectionDeadline is the maximum the TLS server will wait for connections.
	// This deadline also applies to the ammount of time to wait for the server to shutdown.
	ConnectionDeadline time.Duration

	// EstablishLeadership is called when a node becomes the leader of the
	// Raft cluster. This function can be called multiple times if it returns an
	// error.
	EstablishLeadership func() error

	// RevokeLeadership is called when a node loses the leadership of the
	// Raft cluster.
	RevokeLeadership func() error

	// ConsistentNodeStatus

	// Services is an array of services running on top of Cerebrum.
	Services []Service

	// ExistingNodes is an array of nodes already in the cluster.
	ExistingNodes []string
}
