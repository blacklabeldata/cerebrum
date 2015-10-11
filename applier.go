package cerebrum

import (
	"bytes"
	"time"

	"github.com/hashicorp/raft"
	log "github.com/mgutz/logxi/v1"

	"github.com/blacklabeldata/namedtuple"
)

// Applier applies tuples to the Raft log if the node is the leader, otherwise
// the tuple will be forwarded to the leader.
type Applier interface {

	// Apply performs the Raft or forward operation depending on the node's
	// leader status.
	Apply(namedtuple.Tuple) error
}

// RaftApplier covers a few of the raft.Raft methods to make testing easier.
type RaftApplier interface {
	Apply(cmd []byte, timeout time.Duration) raft.ApplyFuture
	State() raft.RaftState
	Leader() string
}

func NewApplier(r RaftApplier, f Forwarder, l log.Logger, timeout time.Duration) Applier {
	return &applier{
		logger:       l,
		raft:         r,
		forwarder:    f,
		enqueueLimit: timeout,
	}
}

type applier struct {
	logger       log.Logger
	raft         RaftApplier
	forwarder    Forwarder
	enqueueLimit time.Duration
}

func (c *applier) Apply(tuple namedtuple.Tuple) error {
	var buf bytes.Buffer
	tuple.WriteTo(&buf)
	data := buf.Bytes()

	if c.raft.State() == raft.Leader {
		future := c.raft.Apply(data, c.enqueueLimit)
		return future.Error()
	}

	// Handle leader forwarding
	return c.forwarder.Forward(data)
}
