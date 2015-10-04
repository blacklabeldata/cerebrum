package cerebrum

import (
	"crypto/tls"
	"io"

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
	LeaderElectionHandler serfer.UserEventHandler

	// UserEvent processes known, non-leader election events.
	UserEvent serfer.UserEventHandler

	// UnknownEventHandler processes unkown events.
	UnknownEventHandler serfer.UserEventHandler

	// Called when a Member joins the cluster.
	NodeJoined serfer.MemberEventHandler

	// Called when a Member leaves the cluster by sending a leave message.
	NodeLeft serfer.MemberEventHandler

	// Called when a Member has been detected as failed.
	NodeFailed serfer.MemberEventHandler

	// Called when a Member has been Readed from the cluster.
	NodeReaped serfer.MemberEventHandler

	// Called when a Member has been updated.
	NodeUpdated serfer.MemberEventHandler

	// Called when a serf.Query is received.
	QueryHandler serfer.QueryEventHandler

	// GossipBindAddr
	GossipBindAddr string

	// GossipBindPort
	GossipBindPort int

	// GossipAdvertiseAddr
	GossipAdvertiseAddr string

	// GossipAdvertisePort
	GossipAdvertisePort int

	// LogOutput
	LogOutput io.Writer

	// RaftConfig
	RaftConfig *raft.Config

	// SnapshotsRetained is the number of snapshots kept for Raft
	SnapshotsRetained int

	// LogCacheSize is the number of log entries to keep in memory.
	LogCacheSize int

	// TLSConfig is the
	TLSConfig *tls.Config
}
