package cerebrum

import (
	"errors"
	"testing"
	"time"

	log "github.com/mgutz/logxi/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestForward_NoLeader(t *testing.T) {
	applier := &MockRaftApplier{leader: ""}
	applier.On("Leader").Return(ErrNoLeader)
	fwdr := &forwarder{
		raft:   applier,
		logger: &log.NullLogger{},
	}

	var buf []byte
	err := fwdr.Forward(buf)
	applier.AssertCalled(t, "Leader")
	assert.Equal(t, ErrNoLeader, err, "Forward should return ErrNoLeader")
}

func TestForward_DialError(t *testing.T) {
	applier := &MockRaftApplier{leader: "leader"}
	applier.On("Leader").Return("leader")

	dialError := errors.New("dial error")
	dialer := &MockDialer{err: dialError}
	dialer.On("Dial", connForward, "leader", 3*time.Second).Return(nil, dialError)

	fwdr := &forwarder{
		raft:   applier,
		dialer: dialer,
		logger: &log.NullLogger{},
	}

	var buf []byte
	err := fwdr.Forward(buf)
	applier.AssertCalled(t, "Leader")
	dialer.AssertCalled(t, "Dial", connForward, "leader", 3*time.Second)
	assert.Equal(t, dialError, err, "Forward should return dial error")
}

func TestForward_WriteError(t *testing.T) {
	applier := &MockRaftApplier{leader: "leader"}
	applier.On("Leader").Return("leader")

	var buf []byte
	writeError := errors.New("write failed")
	conn := &MockConn{err: writeError}
	conn.On("Write", buf).Return(0, writeError)
	conn.On("Close").Return()

	dialer := &MockDialer{err: nil, conn: conn}
	dialer.On("Dial", connForward, "leader", 3*time.Second).Return(conn, nil)

	fwdr := &forwarder{
		raft:   applier,
		dialer: dialer,
		logger: &log.NullLogger{},
	}

	err := fwdr.Forward(buf)
	applier.AssertCalled(t, "Leader")
	dialer.AssertCalled(t, "Dial", connForward, "leader", 3*time.Second)
	conn.AssertCalled(t, "Write", buf)
	assert.Equal(t, writeError, err, "Forward should return write error")
}

func TestForward_Write(t *testing.T) {
	applier := &MockRaftApplier{leader: "leader"}
	applier.On("Leader").Return("leader")

	var buf []byte
	conn := &MockConn{}
	conn.On("Write", buf).Return(0, nil)
	conn.On("Close").Return()

	dialer := &MockDialer{err: nil, conn: conn}
	dialer.On("Dial", connForward, "leader", 3*time.Second).Return(conn, nil)

	fwdr := NewForwarder(applier, dialer, &log.NullLogger{})
	err := fwdr.Forward(buf)

	applier.AssertCalled(t, "Leader")
	dialer.AssertCalled(t, "Dial", connForward, "leader", 3*time.Second)
	conn.AssertCalled(t, "Write", buf)
	assert.Equal(t, nil, err, "Forward should return nil")
}

type MockForwarder struct {
	mock.Mock
	err error
}

func (m *MockForwarder) Forward(buf []byte) error {
	m.Called(buf)
	return m.err
}
