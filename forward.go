package cerebrum

import (
	"crypto/tls"
	"io"
	"net"
	"time"

	log "github.com/mgutz/logxi/v1"

	"github.com/blacklabeldata/grim"
	"github.com/blacklabeldata/namedtuple"
	"github.com/hashicorp/yamux"
	"golang.org/x/net/context"
	tomb "gopkg.in/tomb.v2"
)

type ConnType byte

const (
	connForward ConnType = 0x01
	connRaft             = 0x02
)

const (
	// Warn if the Raft command is larger than this.
	// If it's over 1MB something is probably being abusive.
	raftWarnSize = 1024 * 1024

	// enqueueLimit caps how long we will wait to enqueue
	// a new Raft command. Something is probably wrong if this
	// value is ever reached. However, it prevents us from blocking
	// the requesting goroutine forever.
	enqueueLimit = 30 * time.Second
)

// listen is used to listen for incoming RPC connections
func (s *cerebrum) listen() error {
	defer s.listener.Close()

	// Create tomb for connection goroutines
	var t tomb.Tomb
	t.Go(func() error {
	OUTER:
		for {

			// Accepts will only block for 1s
			s.listener.SetDeadline(time.Now().Add(s.config.ConnectionDeadline))

			select {

			// Stop server on channel receive
			case <-s.t.Dying():
				t.Kill(nil)
				break OUTER
			default:

				// Accept new connection
				tcpConn, err := s.listener.Accept()
				if err != nil {
					if neterr, ok := err.(net.Error); ok && neterr.Timeout() {
						s.logger.Trace("Connection timeout...")
					} else {
						s.logger.Warn("Connection failed", "error", err)
					}
					continue
				}

				// Handle connection
				s.logger.Info("Successful TCP connection:", tcpConn.RemoteAddr().String())
				t.Go(func() error {
					s.handleConn(t, tls.Server(tcpConn, s.config.TLSConfig))
					return nil
				})
			}
		}
		return nil
	})

	return t.Wait()
}

func (c *cerebrum) handleConn(parent tomb.Tomb, conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 1)
	_, err := conn.Read(buf)
	if err != nil {
		return
	}

	var t tomb.Tomb
	switch ConnType(buf[0]) {
	case connForward:
		c.logger.Info("Accepted forwarding connection", "addr", conn.RemoteAddr().String())

		// Return the error for the connection processing
		t.Go(func() error {
			c.handleForwardingConn(t, conn)
			return nil
		})
	case connRaft:
		c.logger.Info("Accepted raft connection", "addr", conn.RemoteAddr().String())
		c.raftLayer.Handoff(conn)
	default:
		c.logger.Warn("Closing TCP connection; unknown type", "addr", conn.RemoteAddr().String())
		return
	}
}

// handleConn is used to determine if this is a Raft or
// Consul type RPC connection and invoke the correct handler
func (s *cerebrum) handleForwardingConn(parent tomb.Tomb, conn net.Conn) {
	defer conn.Close()
	conf := yamux.DefaultConfig()
	conf.LogOutput = s.config.LogOutput
	session, _ := yamux.Server(conn, conf)
	var t tomb.Tomb
	t.Go(func() error {
		for {
			select {
			case <-parent.Dying():
			default:
				stream, err := session.Accept()
				if err != nil {
					if err != io.EOF {
						s.logger.Error("multiplex conn accept failed", "err", err)
					}
					return err
				}
				t.Go(func() error {
					return s.handleStream(t, stream)
				})
			}
		}
	})
	t.Wait()
}

func (c *cerebrum) handleStream(parent tomb.Tomb, stream net.Conn) error {
	decoder := namedtuple.NewDecoder(namedtuple.DefaultRegistry, stream)
	for {
		select {
		case <-parent.Dying():
		default:
			tuple, err := decoder.Decode()
			if err != nil {
				return err
			}

			switch {
			case tuple.Is(nodeStatus):
				if err := c.applier.Apply(tuple); err != nil {
					return err
				}
			default:
			}
		}
	}
	return nil
}

type ForwardingHandler struct {
	applier Applier
	logger  log.Logger
}

func (f *ForwardingHandler) Handle(c context.Context, conn net.Conn) {
	g := grim.ReaperWithContext(c)
	defer g.Wait()

	messages := make(chan namedtuple.Tuple, 1)
	g.Spawn(func(ctx context.Context) {
		defer conn.Close()
		for {
			select {
			case <-ctx.Done():
			case msg, ok := <-messages:
				if !ok {
					return
				}
				if err := f.applier.Apply(msg); err != nil {
					f.logger.Warn("error applying message")
				}
			}
		}
	})
	g.Spawn(func(ctx context.Context) {
		defer close(messages)
		decoder := namedtuple.NewDecoder(namedtuple.DefaultRegistry, conn)
		for {
			tuple, err := decoder.Decode()
			if err != nil {
				return
			}

			if tuple.Is(nodeStatus) {
				messages <- tuple
			}
		}
	})
	return
}
