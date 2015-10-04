package cerebrum

import (
	"crypto/tls"
	"errors"
	"io"
	"net"
	"time"

	"github.com/hashicorp/raft"
)

var (
	errNotAdvertisable = errors.New("local bind address is not advertisable")
	errNotTCP          = errors.New("local address is not a TCP address")
)

// TLSStreamLayer implements raft.StreamLayer interface over TLS.
type TLSStreamLayer struct {
	advertise net.Addr
	listener  *net.TCPListener
	config    *tls.Config
}

// NewTLSTransport returns a raft.NetworkTransport that is built on top of
// a TLS streaming transport layer.
func NewTLSTransport(
	bindAddr string,
	maxPool int,
	timeout time.Duration,
	logOutput io.Writer,
	config *tls.Config,
) (*raft.NetworkTransport, error) {

	// Try to bind
	listener, err := tls.Listen("tcp", bindAddr, config)
	if err != nil {
		return nil, err
	}

	// Create stream
	stream := &TLSStreamLayer{
		advertise: listener.Addr(),
		listener:  listener.(*net.TCPListener),
	}

	// Verify that we have a usable advertise address
	addr, ok := stream.Addr().(*net.TCPAddr)
	if !ok {
		listener.Close()
		return nil, errNotTCP
	}
	if addr.IP.IsUnspecified() {
		listener.Close()
		return nil, errNotAdvertisable
	}

	// Create the network transport
	trans := raft.NewNetworkTransport(stream, maxPool, timeout, logOutput)
	return trans, nil
}

// Dial implements the StreamLayer interface.
func (t *TLSStreamLayer) Dial(address string, timeout time.Duration) (net.Conn, error) {
	return tls.DialWithDialer(&net.Dialer{Timeout: timeout}, "tcp", address, t.config)
}

// Accept implements the net.Listener interface.
func (t *TLSStreamLayer) Accept() (c net.Conn, err error) {
	return t.listener.Accept()
}

// Close implements the net.Listener interface.
func (t *TLSStreamLayer) Close() (err error) {
	return t.listener.Close()
}

// Addr implements the net.Listener interface.
func (t *TLSStreamLayer) Addr() net.Addr {
	// Use an advertise addr if provided
	if t.advertise != nil {
		return t.advertise
	}
	return t.listener.Addr()
}
