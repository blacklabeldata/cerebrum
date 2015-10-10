package cerebrum

import (
	"bytes"
	"testing"
	"time"

	"github.com/blacklabeldata/namedtuple"
	"github.com/hashicorp/raft"
	log "github.com/mgutz/logxi/v1"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestApplier_FollowerState(t *testing.T) {
	raftApplier := &MockRaftApplier{state: raft.Follower}
	raftApplier.On("State").Return(raft.Follower)

	buffer := make([]byte, 512)
	builder := namedtuple.NewBuilder(nodeStatus, buffer)
	builder.PutString("ID", "id")
	builder.PutString("Name", "name")
	builder.PutString("DataCenter", "dc1")
	builder.PutUint8("Status", uint8(0))
	builder.PutString("Addr", "127.0.0.1")
	builder.PutInt32("Port", int32(9000))

	// encode data
	tuple, err := builder.Build()
	assert.Nil(t, err)

	var buf bytes.Buffer
	tuple.WriteTo(&buf)
	data := buf.Bytes()

	fwdr := &MockForwarder{}
	fwdr.On("Forward", data).Return()
	applier := &applier{
		logger:    &log.NullLogger{},
		raft:      raftApplier,
		forwarder: fwdr,
	}
	err = applier.Apply(tuple)
	assert.Nil(t, err)
}

func TestApplier_LeaderState(t *testing.T) {

	buffer := make([]byte, 512)
	builder := namedtuple.NewBuilder(nodeStatus, buffer)
	builder.PutString("ID", "id")
	builder.PutString("Name", "name")
	builder.PutString("DataCenter", "dc1")
	builder.PutUint8("Status", uint8(0))
	builder.PutString("Addr", "127.0.0.1")
	builder.PutInt32("Port", int32(9000))

	// encode data
	tuple, err := builder.Build()
	assert.Nil(t, err)

	var buf bytes.Buffer
	tuple.WriteTo(&buf)
	data := buf.Bytes()

	fwdr := &MockForwarder{}
	future := &MockApplyFuture{}
	future.On("Error").Return(nil)

	raftApplier := &MockRaftApplier{state: raft.Leader, future: future}
	raftApplier.On("State").Return(raft.Leader)
	raftApplier.On("Apply", data, enqueueLimit).Return(future)

	applier := NewApplier(raftApplier, fwdr, &log.NullLogger{})
	err = applier.Apply(tuple)
	assert.Nil(t, err)
	raftApplier.AssertCalled(t, "Apply", data, enqueueLimit)
	fwdr.AssertNotCalled(t, "Forward")
}

type MockRaftApplier struct {
	mock.Mock
	future raft.ApplyFuture
	state  raft.RaftState
	leader string
}

func (m *MockRaftApplier) Apply(cmd []byte, timeout time.Duration) raft.ApplyFuture {
	m.Called(cmd, timeout)
	return m.future
}

func (m *MockRaftApplier) State() raft.RaftState {
	m.Called()
	return m.state
}

func (m *MockRaftApplier) Leader() string {
	m.Called()
	return m.leader
}

type MockApplyFuture struct {
	mock.Mock
	err      error
	response interface{}
	index    uint64
}

func (m *MockApplyFuture) Error() error {
	m.Called()
	return m.err
}

func (m *MockApplyFuture) Response() interface{} {
	m.Called()
	return m.response
}

func (m *MockApplyFuture) Index() uint64 {
	m.Called()
	return m.index
}
