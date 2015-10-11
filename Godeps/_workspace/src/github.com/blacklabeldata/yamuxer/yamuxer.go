package yamuxer

import (
	"crypto/tls"
	"io"
	"net"
	"time"

	"github.com/hashicorp/yamux"
	log "github.com/mgutz/logxi/v1"

	"github.com/blacklabeldata/grim"
	"golang.org/x/net/context"
)

// Yamuxer is a simple interface which starts listening for connections on
// Start and waits for all connections to be closed on Stop.
type Yamuxer interface {
	Start()
	Stop()
}

// Config is used to create a new yamuxer.
type Config struct {
	Listener   *net.TCPListener
	TLSConfig  *tls.Config
	Dispatcher Dispatcher
	Deadline   time.Duration
	LogOutput  io.Writer
}

// New creates an implementation of the Yamuxer interface using the given
// context and config.
func New(context context.Context, c *Config) Yamuxer {
	return &yamuxer{
		grim:       grim.ReaperWithContext(context),
		listener:   c.Listener,
		dispatcher: c.Dispatcher,
		deadline:   c.Deadline,
		tlsConfig:  c.TLSConfig,
		logger:     log.NewLogger(c.LogOutput, "yamuxer"),
		logOutput:  c.LogOutput,
	}
}

// yamuxer implements the Yamuxer interface.
type yamuxer struct {
	grim       grim.GrimReaper
	listener   *net.TCPListener
	dispatcher Dispatcher
	deadline   time.Duration
	tlsConfig  *tls.Config
	logger     log.Logger
	logOutput  io.Writer
}

func (y *yamuxer) Start() {
	y.grim.SpawnFunc(y.listen)
}

func (y *yamuxer) Stop() {
	y.grim.Kill()
	y.grim.Wait()
	y.listener.Close()
}

func (y *yamuxer) listen(ctx context.Context) {
	defer y.listener.Close()
OUTER:
	for {

		// Accepts will only block for 1s
		y.listener.SetDeadline(time.Now().Add(y.deadline))

		select {

		// Stop server on channel receive
		case <-ctx.Done():
			break OUTER
		default:

			// Accept new connection
			tcpConn, err := y.listener.Accept()
			if err != nil {
				if neterr, ok := err.(net.Error); ok && neterr.Timeout() {
					y.logger.Trace("Connection timeout...")
				} else {
					y.logger.Warn("Connection failed", "error", err)
				}
				continue
			}

			// Handle connection
			y.logger.Info("Successful TCP connection:", tcpConn.RemoteAddr().String())
			y.handleConn(y.grim.New(), tls.Server(tcpConn, y.tlsConfig))
		}
	}
	return
}

func (y *yamuxer) handleConn(g grim.GrimReaper, conn net.Conn) {
	defer g.Wait()

	conf := yamux.DefaultConfig()
	conf.LogOutput = y.logOutput
	session, _ := yamux.Server(conn, conf)

	streamCh := make(chan net.Conn)
	g.SpawnFunc(processStreams(g.New(), conn, streamCh, y.dispatcher))
	g.SpawnFunc(acceptStreams(y.logger, session, streamCh))
}

func processStreams(g grim.GrimReaper, conn net.Conn, streamCh chan net.Conn, dispatcher Dispatcher) grim.TaskFunc {
	return func(ctx context.Context) {
		defer conn.Close()
		for {
			select {
			case <-ctx.Done():
				return
			case stream, ok := <-streamCh:
				if !ok {
					return
				}
				dispatcher.Dispatch(g.New(), stream)
			}
		}
	}
}

func acceptStreams(logger log.Logger, session *yamux.Session, streamCh chan net.Conn) grim.TaskFunc {
	return func(ctx context.Context) {
		defer close(streamCh)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				stream, err := session.Accept()
				if err != nil {
					if err != io.EOF {
						logger.Error("multiplex conn accept failed", "err", err)
					}
					return
				}
				streamCh <- stream
			}
		}
	}
}
