package cerebrum

import (
	"net"
	"time"

	"github.com/stretchr/testify/mock"
)

type MockDialer struct {
	mock.Mock
	err  error
	conn net.Conn
}

func (m *MockDialer) Dial(c ConnType, addr string, timeout time.Duration) (net.Conn, error) {
	m.Called(c, addr, timeout)
	return m.conn, m.err
}

func (m *MockDialer) Shutdown() error {
	m.Called()
	return nil
}

type MockConn struct {
	mock.Mock
	err error
}

func (m *MockConn) Read(b []byte) (n int, err error) {
	m.Called(b)
	return 0, m.err
}

func (m *MockConn) Write(b []byte) (n int, err error) {
	m.Called(b)
	return 0, m.err
}

func (m *MockConn) Close() error {
	m.Called()
	return m.err
}

func (m *MockConn) LocalAddr() net.Addr {
	m.Called()
	return nil
}

func (m *MockConn) RemoteAddr() net.Addr {
	m.Called()
	return nil
}

func (m *MockConn) SetDeadline(t time.Time) error {
	m.Called(t)
	return m.err
}

func (m *MockConn) SetReadDeadline(t time.Time) error {
	m.Called(t)
	return m.err
}

func (m *MockConn) SetWriteDeadline(t time.Time) error {
	m.Called(t)
	return m.err
}
