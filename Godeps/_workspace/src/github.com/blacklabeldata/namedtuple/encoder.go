package namedtuple

import (
	"bytes"
	"errors"
	"io"
	"math"

	"github.com/blacklabeldata/xbinary"
)

var (
	// ErrInvalidFieldSize is retured by an encoder if the the field size is invalid
	ErrInvalidFieldSize = errors.New("Invalid Field Size: field size must be 1,2,4 or 8 bytes")
)

// Encoder encodes tuples normally into a given io.Writer.
type Encoder interface {
	Encode(Tuple) error
}

// NewEncoder creates a new encoder with the given io.Writer
func NewEncoder(w io.Writer) Encoder {
	return versionOneEncoder{w, make([]byte, 9), bytes.NewBuffer(make([]byte, 0, 4096))}
}

type versionOneEncoder struct {
	w              io.Writer
	protocolHeader []byte
	buffer         *bytes.Buffer
}

func (e versionOneEncoder) Encode(t Tuple) error {
	defer e.buffer.Reset()

	// Write tuple header and payload to buffer
	if _, err := e.writeTuple(t); err != nil {
		return err
	}

	// Write protocol header to underlying writer
	size := e.buffer.Len()
	if err := e.writeProtocolHeader(size); err != nil {
		return err
	}

	// Write buffer to writer (buffer should now contain the tuple header and the body)
	_, err := e.w.Write(e.buffer.Bytes())
	return err
}

func (e versionOneEncoder) writeProtocolHeader(size int) (err error) {

	// Set protocol version to 1
	e.protocolHeader[0] = 1

	// Write protocol version, size enum and content length
	if size < math.MaxUint8 {
		// Size enum for an 8-bit sized tuple is 0.

		// Write tuple size
		e.protocolHeader[1] = uint8(size)

		// Write buffer to writer
		_, err = e.w.Write(e.protocolHeader[:2])
	} else if size < math.MaxUint16 {

		// Set size enum (mask: 0b01000000)
		e.protocolHeader[0] |= 64

		if _, err := xbinary.LittleEndian.PutUint16(e.protocolHeader, 1, uint16(size)); err != nil {
			return err
		}

		// Write buffer to writer
		_, err = e.w.Write(e.protocolHeader[:3])
	} else if size < math.MaxUint32 {
		// Set size enum (mask: 0b01000000)
		e.protocolHeader[0] |= 128

		if _, err := xbinary.LittleEndian.PutUint32(e.protocolHeader, 1, uint32(size)); err != nil {
			return err
		}

		// Write buffer to writer
		_, err = e.w.Write(e.protocolHeader[:5])
	} else {
		// Set size enum (mask: 0b11000000)
		e.protocolHeader[0] |= 192

		if _, err := xbinary.LittleEndian.PutUint64(e.protocolHeader, 1, uint64(size)); err != nil {
			return err
		}

		// Write buffer to writer
		_, err = e.w.Write(e.protocolHeader[:9])
	}

	// Check for error when writing protocol header
	if err != nil {
		return err
	}
	return nil
}

func (e versionOneEncoder) writeTupleHeader(t TupleHeader) (int64, error) {
	return t.WriteTo(e.buffer)
}

func (e versionOneEncoder) writeTuple(t Tuple) (int64, error) {

	// write header
	var wrote int
	if wrote, err := e.writeTupleHeader(t.Header); err != nil {
		return wrote, nil
	}

	n, err := e.buffer.Write(t.data)
	if err != nil {
		return int64(n), err
	}
	return int64(wrote + n), nil
}
