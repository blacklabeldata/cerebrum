package cerebrum

import (
	"net"

	log "github.com/mgutz/logxi/v1"

	"github.com/blacklabeldata/grim"
	"github.com/blacklabeldata/namedtuple"
	"github.com/blacklabeldata/yamuxer"
	"golang.org/x/net/context"
)

const (
	connForward yamuxer.StreamType = 0x01
	connRaft                       = 0x02
)

const (
// enqueueLimit caps how long we will wait to enqueue
// a new Raft command. Something is probably wrong if this
// value is ever reached. However, it prevents us from blocking
// the requesting goroutine forever.
// enqueueLimit = 30 * time.Second
)

type ForwardingHandler struct {
	applier Applier
	logger  log.Logger
}

func (f *ForwardingHandler) Handle(c context.Context, conn net.Conn) {
	f.logger.Info("Accepted raft connection", "addr", conn.RemoteAddr().String())
	g := grim.ReaperWithContext(c)
	defer g.Wait()

	messages := make(chan namedtuple.Tuple, 1)
	g.SpawnFunc(func(ctx context.Context) {
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
	g.SpawnFunc(func(ctx context.Context) {
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
