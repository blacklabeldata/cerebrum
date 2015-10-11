package cerebrum

import (
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/blacklabeldata/yamuxer"
	"github.com/hashicorp/yamux"
)

// Conn is a pooled connection to a Cerebrum server
type Conn struct {
	// refCount    int32
	// shouldClose int32

	addr     string
	session  *yamux.Session
	lastUsed time.Time

	pool *ConnPool
}

// markForUse does all the bookkeeping required to ready a connection for use.
func (c *Conn) markForUse() {
	c.lastUsed = time.Now()
}

func (c *Conn) Close() error {
	return c.session.Close()
}

// getClient is used to get a cached or new client
func (c *Conn) getClient() (net.Conn, error) {

	// Open a new session
	stream, err := c.session.Open()
	if err != nil {
		return nil, err
	}
	return stream, nil
}

func (c *Conn) dial(t yamuxer.StreamType) (net.Conn, error) {
	stream, err := c.getClient()
	if err != nil {
		return nil, err
	}

	if _, err = stream.Write([]byte{byte(t)}); err != nil {
		return nil, err
	}
	return stream, nil
}

func (c *Conn) dialRaft() (net.Conn, error) {
	return c.dial(connRaft)
}

func (c *Conn) dialForwarding() (net.Conn, error) {
	return c.dial(connForward)
}

// ConnPool is used to maintain a connection pool to other
// Nomad servers. This is used to reduce the latency of
// RPC requests between servers. It is only used to pool
// connections in the rpc mode. Raft connections
// are pooled separately.
type ConnPool struct {
	sync.Mutex

	// LogOutput is used to control logging
	logOutput io.Writer

	// The maximum time to keep a connection open
	maxTime time.Duration

	// Pool maps an address to a open connection
	pool map[string]*Conn

	// TLS config
	config *tls.Config

	// Used to indicate the pool is shutdown
	shutdown   bool
	shutdownCh chan struct{}
}

// NewPool is used to make a new connection pool
// Maintain at most one connection per host, for up to maxTime.
// Set maxTime to 0 to disable reaping.
// If TLS settings are provided outgoing connections use TLS.
func NewPool(logOutput io.Writer, maxTime time.Duration, config *tls.Config) *ConnPool {
	pool := &ConnPool{
		logOutput:  logOutput,
		maxTime:    maxTime,
		pool:       make(map[string]*Conn),
		config:     config,
		shutdownCh: make(chan struct{}),
	}
	if maxTime > 0 {
		go pool.reap()
	}
	return pool
}

// Shutdown is used to close the connection pool
func (p *ConnPool) Shutdown() error {
	p.Lock()
	defer p.Unlock()

	for _, conn := range p.pool {
		conn.Close()
	}
	p.pool = make(map[string]*Conn)

	if p.shutdown {
		return nil
	}
	p.shutdown = true
	close(p.shutdownCh)
	return nil
}

// Acquire is used to get a connection that is
// pooled or to return a new connection
func (p *ConnPool) acquire(addr string, timeout time.Duration) (*Conn, error) {
	// Check to see if there's a pooled connection available. This is up
	// here since it should the the vastly more common case than the rest
	// of the code here.
	p.Lock()
	defer p.Unlock()

	c := p.pool[addr]
	if c != nil {
		c.markForUse()
		p.Unlock()
		return c, nil
	}

	// Create new connection if it doesn't exist
	c, err := p.getNewConn(addr, timeout)
	if err != nil {
		return nil, err
	}
	p.pool[addr] = c

	return c, nil
}

// getNewConn is used to return a new connection
func (p *ConnPool) getNewConn(addr string, timeout time.Duration) (*Conn, error) {
	// Try to dial the conn
	conn, err := tls.DialWithDialer(&net.Dialer{Timeout: timeout}, "tcp", addr, p.config)
	if err != nil {
		return nil, err
	}

	// Setup the logger
	conf := yamux.DefaultConfig()
	conf.LogOutput = p.logOutput

	// Create a multiplexed session
	session, err := yamux.Client(conn, conf)
	if err != nil {
		conn.Close()
		return nil, err
	}

	// Wrap the connection
	c := &Conn{
		addr:     addr,
		session:  session,
		lastUsed: time.Now(),
		pool:     p,
	}
	return c, nil
}

// clearConn is used to clear any cached connection, potentially in response to an error
func (p *ConnPool) clearConn(conn *Conn) {
	// Ensure returned streams are closed
	// atomic.StoreInt32(&conn.shouldClose, 1)

	// Clear from the cache
	p.Lock()
	if c, ok := p.pool[conn.addr]; ok && c == conn {
		delete(p.pool, conn.addr)
	}

	// Close down immediately if idle
	if conn.session.NumStreams() == 0 {
		conn.Close()
	}
	p.Unlock()
}

// getClient is used to get a usable client for an address
func (p *ConnPool) getClient(addr string, timeout time.Duration) (*Conn, error) {
	conn, err := p.acquire(addr, timeout)
	if err != nil {
		return nil, fmt.Errorf("failed to get conn: %v", err)
	}
	return conn, nil
}

func (p *ConnPool) dialRaft(addr string, timeout time.Duration) (net.Conn, error) {
	// Get a usable client
	conn, err := p.getClient(addr, timeout)
	if err != nil {
		return nil, fmt.Errorf("rpc error: %v", err)
	}

	// Return raft stream
	return conn.dialRaft()
}

func (p *ConnPool) dialForwarding(addr string, timeout time.Duration) (net.Conn, error) {
	// Get a usable client
	conn, err := p.getClient(addr, timeout)
	if err != nil {
		return nil, fmt.Errorf("rpc error: %v", err)
	}

	// Return forwarding raft stream
	return conn.dialForwarding()
}

// Reap is used to close conns open over maxTime
func (p *ConnPool) reap() {
	for {
		// Sleep for a while
		select {
		case <-p.shutdownCh:
			return
		case <-time.After(time.Second):
		}

		// Reap all old conns
		p.Lock()
		var removed []string
		now := time.Now()
		for host, conn := range p.pool {
			// Skip recently used connections
			if now.Sub(conn.lastUsed) < p.maxTime {
				continue
			}

			// Skip connections with active streams
			if conn.session.NumStreams() > 0 {
				continue
			}

			// Close the conn
			conn.Close()

			// Remove from pool
			removed = append(removed, host)
		}
		for _, host := range removed {
			delete(p.pool, host)
		}
		p.Unlock()
	}
}
