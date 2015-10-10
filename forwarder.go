package cerebrum

import (
	"time"

	log "github.com/mgutz/logxi/v1"
)

// Forwarder forwards data to the cluster leader. If the leader is unknown,
// an error will be returned. If the leader cannot be contacted, an error
// returned. If the data cannot be written to the leader, an error will be
// returned.
type Forwarder interface {
	Forward([]byte) error
}

func NewForwarder(r RaftApplier, d Dialer, l log.Logger) Forwarder {
	return &forwarder{r, d, l}
}

type forwarder struct {
	raft   RaftApplier
	dialer Dialer
	logger log.Logger
}

// Forward is used to forward an RPC call to the leader, or fail if no leader
func (f *forwarder) Forward(buf []byte) error {
	leader := f.raft.Leader()
	if leader == "" {
		f.logger.Error("No cluster leader")
		return ErrNoLeader
	}

	conn, err := f.dialer.Dial(connForward, leader, 3*time.Second)
	if err != nil {
		f.logger.Error("Failed to dial cluster leader", "leader", leader, "err", err)
		return err
	}
	defer conn.Close()
	if _, err = conn.Write(buf); err != nil {
		f.logger.Error("Failed to send data to cluster leader", "leader", leader, "err", err)
		return err
	}
	return nil
}
