package yamuxer

import (
	"net"

	log "github.com/mgutz/logxi/v1"

	"github.com/blacklabeldata/grim"
	"golang.org/x/net/context"
)

// StreamType is the stream enum. There can only be one handler per type. 0x00
// is a reserved StreamType.
type StreamType byte

// UnknownStreamType is a cautionary type as it covered the nil byte.
var UnknownStreamType StreamType = 0x00

// Dispatcher processes stream connections and hands them off to the registered
// handler.
type Dispatcher interface {

	// Register adds a Handler for the given type.
	Register(StreamType, Handler)

	// RegisterFunc adds a HandlerFunc for the given type.
	RegisterFunc(StreamType, HandlerFunc)

	// Dispatch reads the first byte and passes the connection to the handler
	// for the given type.
	Dispatch(grim.GrimReaper, net.Conn)
}

// NewDispatcher creates a Dispatcher implementation which also allows for
// the optional handling of unknown streams. If the unknown handler is nil,
// the streams will be closed immediately.
func NewDispatcher(l log.Logger, unknownHandler Handler) Dispatcher {
	return &dispatcher{
		logger:   l,
		unknown:  unknownHandler,
		handlers: make(map[StreamType]Handler, 0),
	}
}

// dispatcher is an implementation of the Dispatcher interface.
type dispatcher struct {
	logger   log.Logger
	unknown  Handler
	handlers map[StreamType]Handler
}

// Register adds the handler based on type.
func (d *dispatcher) Register(c StreamType, h Handler) {
	d.handlers[c] = h
}

// RegisterFunc adds the HandlerFunc based on type.
func (d *dispatcher) RegisterFunc(c StreamType, h HandlerFunc) {
	d.handlers[c] = &basicHandler{h}
}

// Dispatch processes the connection.
func (d *dispatcher) Dispatch(g grim.GrimReaper, conn net.Conn) {

	// Create buffer and read connection type
	buf := make([]byte, 1)
	if _, err := conn.Read(buf); err != nil {
		d.logger.Warn("error reading StreamType",
			"remote", conn.RemoteAddr().String, "err", err)
		conn.Close()
		return
	}

	// Dispatch connection based on first byte
	ctype := StreamType(buf[0])
	if handler, ok := d.handlers[ctype]; ok {
		g.SpawnFunc(func(c context.Context) {
			handler.Handle(c, conn)
		})
		return
	}

	// Process unknown streams if there is an unknown stream handler
	if d.unknown != nil {
		g.SpawnFunc(func(c context.Context) {
			d.unknown.Handle(c, conn)
		})
		return
	}

	// If there is not an unknown stream handler, close the stream.
	d.logger.Warn("closing connection: unknown stream type",
		"type", ctype, "remote", conn.RemoteAddr().String())
	conn.Close()
}
