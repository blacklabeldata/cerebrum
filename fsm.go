package cerebrum

import (
	"bytes"
	"errors"
	"io"
	"time"

	"github.com/blacklabeldata/namedtuple"
	"github.com/hashicorp/raft"
	log "github.com/mgutz/logxi/v1"
)

type fsm struct {
	logOutput io.Writer
	logger    log.Logger
	path      string
	userFSM   raft.FSM
}

// NewFSM is used to construct a new FSM with a blank state
func NewFSM(path string, userFSM raft.FSM, logOutput io.Writer) (raft.FSM, error) {
	fsm := &fsm{
		logOutput: logOutput,
		logger:    log.NewLogger(logOutput, ""),
		path:      path,
		userFSM:   userFSM,
	}
	return fsm, nil
}

func (c *fsm) Apply(log *raft.Log) interface{} {
	buf := log.Data

	reader := bytes.NewReader(buf)
	dec := namedtuple.NewDecoder(namedtuple.DefaultRegistry, reader)
	tup, err := dec.Decode()
	if err != nil {
		return c.userFSM.Apply(log)
	}

	switch {
	case tup.Is(nodeStatus):
		return c.applyNodeStatus(tup)
	default:
		return c.userFSM.Apply(log)
	}
}

func (f *fsm) applyNodeStatus(t namedtuple.Tuple) error {
	// TODO: Store node status in persistent storage
	f.logger.Info("applyNodeStatus")
	return nil
}

func (c *fsm) Snapshot() (raft.FSMSnapshot, error) {
	defer func(start time.Time) {
		c.logger.Info("snapshot created", "elapsed", time.Now().Sub(start))
	}(time.Now())

	// TODO: Create snapshot struct.

	return nil, errors.New("Not implemented")
	// // Create a new snapshot
	// snap, err := c.state.Snapshot()
	// if err != nil {
	// 	return nil, err
	// }
	// return &consulSnapshot{snap}, nil
}

func (c *fsm) Restore(old io.ReadCloser) error {
	defer old.Close()

	// TODO: Create file.
	// TODO: Pipe restore into new file.
	// TODO: Switch node storage.

	return errors.New("Not implemented")
}
