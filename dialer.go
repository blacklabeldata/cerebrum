package cerebrum

import (
	"net"
	"time"
)

// Dialer dials an address and initializes a connection of the given type.
// If the timeout occurs before the dial completes, an error will be returned.
type Dialer interface {
	Dial(t ConnType, addr string, timeout time.Duration) (net.Conn, error)
	Shutdown() error
}

// NewDialer creates a new Dialer implementation with a connection pool.
func NewDialer(p *ConnPool) Dialer {
	return &dialer{p}
}

type dialer struct {
	pool *ConnPool
}

func (d *dialer) Dial(c ConnType, address string, timeout time.Duration) (net.Conn, error) {
	switch c {
	case connRaft:
		return d.pool.dialRaft(address, timeout)
	case connForward:
		return d.pool.dialForwarding(address, timeout)
	default:
		return nil, ErrUnknownConnType
	}
}

func (d *dialer) Shutdown() error {
	return d.pool.Shutdown()
}
