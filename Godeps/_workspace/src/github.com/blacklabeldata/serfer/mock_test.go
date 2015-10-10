package serfer

import (
	"github.com/hashicorp/serf/serf"

	"github.com/stretchr/testify/mock"
)

// MockEventHandler mocks a basic Event handler.
type MockEventHandler struct {
	mock.Mock
}

// HandleEvent processes member events.
func (m *MockEventHandler) HandleEvent(e serf.Event) {
	m.Called(e)
	return
}

// HandleUnknownEvent processes member events.
func (m *MockEventHandler) HandleUnknownEvent(e serf.UserEvent) {
	m.Called(e)
	return
}

func (m *MockEventHandler) HandleMemberJoin(e serf.MemberEvent) {
	m.Called(e)
	return
}
func (m *MockEventHandler) HandleMemberUpdate(e serf.MemberEvent) {
	m.Called(e)
	return
}
func (m *MockEventHandler) HandleMemberLeave(e serf.MemberEvent) {
	m.Called(e)
	return
}
func (m *MockEventHandler) HandleMemberFailure(e serf.MemberEvent) {
	m.Called(e)
	return
}
func (m *MockEventHandler) HandleMemberReap(e serf.MemberEvent) {
	m.Called(e)
	return
}
func (m *MockEventHandler) HandleUserEvent(e serf.UserEvent) {
	m.Called(e)
	return
}
func (m *MockEventHandler) HandleQueryEvent(e serf.Query) {
	m.Called(e)
	return
}
func (m *MockEventHandler) HandleLeaderElection(e serf.UserEvent) {
	m.Called(e)
	return
}
func (m *MockEventHandler) Reconcile(e serf.Member) {
	m.Called(e)
	return
}

// MockEvent
type MockEvent struct {
	mock.Mock

	Type serf.EventType
	Name string
}

// EventType returns the EventType
func (m *MockEvent) EventType() serf.EventType {
	m.Called()
	return m.Type
}

// String returns the EventType name
func (m *MockEvent) String() string {
	return m.Name
}
