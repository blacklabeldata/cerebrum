package cerebrum

import (
	"io"
	"testing"

	"github.com/hashicorp/raft"

	"github.com/stretchr/testify/mock"
)

func TestApply_Error(t *testing.T) {

}

type MockFSM struct {
	mock.Mock
	err error
}

func (m *MockFSM) Apply(log *raft.Log) interface{} {
	m.Called("Apply", log)
	return m.err
}

func (m *MockFSM) Snapshot() (raft.FSMSnapshot, error) {
	m.Called("Snapshot")
	return nil, m.err
}

func (m *MockFSM) Restore(old io.ReadCloser) error {
	m.Called("Restore", old)
	return m.err
}
