package namedtuple

import (
	"bufio"
	"bytes"
	"errors"
	"io"

	"github.com/blacklabeldata/xbinary"
)

var (
	// ErrTupleExceedsMaxSize is returned if the length of a Tuple of greater than the maximum allowable size
	// for the Decoder.
	ErrTupleExceedsMaxSize = errors.New("Tuple exceeds maximum allowable length")

	// ErrInvalidProtocolVersion is returned from Decode() if the Tuple version is unknown.
	ErrInvalidProtocolVersion = errors.New("Invalid protocol version in Tuple header")

	// ErrTupleLengthTooSmall is returned from the Decode() method if the decoded length is too small to include all the required information
	ErrTupleLengthTooSmall = errors.New("Tuple length is too short to include all the required information")

	// ErrUnknownTupleType is returned when the Tuple being decoded is of an unknown type.
	ErrUnknownTupleType = errors.New("Unknown tuple type")

	// ErrInvalidLength is returned if the byte count for the length is not 1, 2, 4 or 8.
	ErrInvalidLength = errors.New("Invalid Tuple Size: tuple length must be encoded as 1,2,4 or 8 bytes")

	// EmptyTuple is returned along with an error from the Decode() method.
	EmptyTuple = Tuple{}
)

const (

	// VersionOneTupleHeaderSize is the size of the header for version one.
	VersionOneTupleHeaderSize = 13

	// DefaultMaxSize is the default maximum size of a tuple being decoded.
	DefaultMaxSize uint64 = 4096
)

// Create a reader which reads the first byte and the content length.
// If the length exceeds the maxSize, return an error
// Create a bytes.Buffer and io.CopyN(contentLength) into the buffer
// Based on the protocol version, decode(buffer.Bytes()) into (Tuple, error)

// decoder := NewDecoder(reg, 65536)
// for _, tup, err := decoder.Decode(reader); err != nil {
// }

// Decoder decodes data into a Tuple or an error.
type Decoder interface {
	Decode() (Tuple, error)
}

// NewDecoder creates a new Decoder using a type Registry and an io.Reader.
func NewDecoder(reg Registry, r io.Reader) Decoder {
	var buf []byte
	return decoder{reg, DefaultMaxSize, bytes.NewBuffer(buf), bufio.NewReader(r)}
}

// NewDecoderSize creates a new Decoder using a type Registry, a max size and an io.Reader.
func NewDecoderSize(reg Registry, maxSize uint64, r io.Reader) Decoder {
	var buf []byte
	return decoder{reg, maxSize, bytes.NewBuffer(buf), bufio.NewReader(r)}
}

type decoder struct {
	reg     Registry
	maxSize uint64
	buffer  *bytes.Buffer
	reader  *bufio.Reader
}

func (d decoder) Decode() (Tuple, error) {

	// Reads the protocol header
	pH, err := d.reader.ReadByte()
	if err != nil {
		return EmptyTuple, err
	}

	// Parse nuber of length bytes and version
	byteCount, version := ParseProtocolHeader(pH)

	// Read bytes for content length
	b := make([]byte, byteCount)
	n, err := d.reader.Read(b)
	if err != nil {
		return EmptyTuple, err
	} else if n != int(byteCount) {
		return EmptyTuple, io.ErrUnexpectedEOF
	}

	// Parse content length based on number of bytes
	length, err := d.parseLength(byteCount, b)
	if err != nil {
		// This should not happen as the
		// Read call above also checks for length.
		return EmptyTuple, err
	}

	// Verify length against maxSize
	if length > d.maxSize {
		return EmptyTuple, ErrTupleExceedsMaxSize
	}

	// Copy Length bytes into buffer
	if _, err := io.CopyN(d.buffer, d.reader, int64(length)); err != nil {
		return EmptyTuple, err
	}

	// Depending on the protocol version, parse the tuple
	switch version {
	case 1:
		return d.parseVersionOneTuple(byteCount, version, length)
	default:
		return EmptyTuple, ErrInvalidProtocolVersion
	}
}

func (d decoder) parseLength(byteCount uint8, buf []byte) (l uint64, err error) {
	switch byteCount {
	case 1:
		if len(buf) == 1 {
			l = uint64(buf[0])
		} else {
			err = xbinary.ErrOutOfRange
		}
	case 2:
		if size, e := xbinary.LittleEndian.Uint16(buf, 0); e == nil {
			l = uint64(size)
		} else {
			err = e
		}
	case 4:
		if size, e := xbinary.LittleEndian.Uint32(buf, 0); e == nil {
			l = uint64(size)
		} else {
			err = e
		}
	case 8:
		if size, e := xbinary.LittleEndian.Uint64(buf, 0); e == nil {
			l = uint64(size)
		} else {
			err = e
		}
	default:
		// This should never happen
		err = ErrInvalidLength
	}
	return
}

func (d decoder) parseVersionOneTuple(offsetSize uint8, protocolVersion uint8, length uint64) (t Tuple, err error) {
	buffer := d.buffer.Bytes()
	var namespaceHash, typeHash, fieldCount uint32
	var version uint8

	// The buffer needs to be at least 13 bytes. This includes the uint8 tuple version, the uint32 namespace and type hashes and the field count
	if len(buffer) < VersionOneTupleHeaderSize {
		return EmptyTuple, ErrTupleLengthTooSmall
	}

	// Read Tuple version
	version = buffer[0]

	// Read namespace hash
	namespaceHash, err = xbinary.LittleEndian.Uint32(buffer, 1)
	if err != nil {
		// Should not occur as buffer length has already been validated
		return EmptyTuple, err
	}

	// Read type hash
	typeHash, err = xbinary.LittleEndian.Uint32(buffer, 5)
	if err != nil {
		// Should not occur as buffer length has already been validated
		return EmptyTuple, err
	}

	// Check if known tuple type
	tupleType, exists := d.reg.GetWithHash(namespaceHash, typeHash)
	if !exists {
		return EmptyTuple, ErrUnknownTupleType
	}

	// Read field count
	fieldCount, err = xbinary.LittleEndian.Uint32(buffer, 9)
	if err != nil {
		// Should not occur as buffer length has already been validated
		return EmptyTuple, err
	}

	// Read field offsets
	offsets, err := readFieldOffsets(offsetSize, fieldCount, buffer)
	if err != nil {
		return EmptyTuple, err
	}

	// Slice tuple data
	pos := VersionOneTupleHeaderSize + int(fieldCount)*int(offsetSize)
	t.data = buffer[pos:]

	// Create TupleHeader
	t.Header = TupleHeader{
		ProtocolVersion: protocolVersion,
		TupleVersion:    version,
		NamespaceHash:   namespaceHash,
		Hash:            typeHash,
		FieldCount:      fieldCount,
		FieldSize:       offsetSize,
		ContentLength:   uint64(len(t.data)),
		Offsets:         offsets,
		Type:            tupleType,
	}
	return
}

func readFieldOffsets(byteCount uint8, fieldCount uint32, buffer []byte) ([]uint64, error) {
	offsets := make([]uint64, int(fieldCount))
	var err error
	switch byteCount {
	case 1:
		// Check buffer length
		if len(buffer) < int(fieldCount)+VersionOneTupleHeaderSize {
			err = ErrTupleLengthTooSmall
		} else {

			// Process offsets
			for i := VersionOneTupleHeaderSize; i < int(fieldCount)+VersionOneTupleHeaderSize; i++ {
				offsets[i-VersionOneTupleHeaderSize] = uint64(buffer[i])
			}
		}
	case 2:
		o := make([]uint16, int(fieldCount))
		err = xbinary.LittleEndian.Uint16Array(buffer, VersionOneTupleHeaderSize, &o)
		if err == nil {
			for i, offset := range o {
				offsets[i] = uint64(offset)
			}
		}
	case 4:
		o := make([]uint32, int(fieldCount))
		err = xbinary.LittleEndian.Uint32Array(buffer, VersionOneTupleHeaderSize, &o)
		if err == nil {
			for i, offset := range o {
				offsets[i] = uint64(offset)
			}
		}
	case 8:
		o := make([]uint64, int(fieldCount))
		err = xbinary.LittleEndian.Uint64Array(buffer, VersionOneTupleHeaderSize, &o)
		if err == nil {
			for i, offset := range o {
				offsets[i] = uint64(offset)
			}
		}
	default:
		err = ErrInvalidLength
	}
	return offsets, err
}
