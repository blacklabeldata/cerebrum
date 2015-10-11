package yamuxer

import (
	"net"

	"golang.org/x/net/context"
)

// Handler handles incoming stream connections for a stream type.
type Handler interface {
	Handle(context.Context, net.Conn)
}

// HandlerFunc handles incoming stream connections.
type HandlerFunc func(context.Context, net.Conn)

// basicHandler wraps a HandlerFunc and implements the Handler interface.
type basicHandler struct {
	handlerFunc HandlerFunc
}

// Handle forwards the connection and context to the handleFunc.
func (h *basicHandler) Handle(c context.Context, conn net.Conn) {
	h.handlerFunc(c, conn)
}
