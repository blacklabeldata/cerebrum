package namedtuple

import (
	"encoding/binary"
	"errors"
	"io"
)

// ProtocolVersionMask is the lower 6 bits of the first byte of the ptotocol header (0b00111111)
const ProtocolVersionMask = 63

// ProtocolVersionMask is the upper 2 bits of the first byte of the ptotocol header (0b11000000)
const ProtocolSizeEnumMask = 192

// ParseProtocolHeader returns the number of bytes to read for the content length and the protocol version. The upper 2 bits represent the number of bytes in the content length (0-3 bits = 2**n). The lower 6 bits are the protocol version which determines how the bytes are interpreted.
func ParseProtocolHeader(header uint8) (lenBytes uint8, version uint8) {
	version = header & ProtocolVersionMask
	lenBytes = 1 << ((header & ProtocolSizeEnumMask) >> 6)
	return
}

// TupleHeader stores meta data about the tuple such as the version, the hashes and the number of fields.
type TupleHeader struct {
	ProtocolVersion uint8
	TupleVersion    uint8
	NamespaceHash   uint32
	Hash            uint32
	FieldCount      uint32
	FieldSize       uint8
	ContentLength   uint64
	Offsets         []uint64
	Type            TupleType
}

// Size returns the Version 1 header size plus the size of all the offsets
func (t *TupleHeader) Size() int {
	return VersionOneTupleHeaderSize + int(t.FieldSize)*int(t.FieldCount)
}

// WriteTo writes the TupleHeader into the given writer.
func (t *TupleHeader) WriteTo(w io.Writer) (int64, error) {

	if len(t.Offsets) != int(t.FieldCount) {
		return 0, errors.New("Invalid Header: Field count does not equal number of field offsets")
	}

	// Encode Header
	dst := make([]byte, t.Size())
	dst[0] = byte(t.TupleVersion)
	binary.LittleEndian.PutUint32(dst[1:], t.NamespaceHash)
	binary.LittleEndian.PutUint32(dst[5:], t.Hash)
	binary.LittleEndian.PutUint32(dst[9:], t.FieldCount)

	pos := int64(13)
	switch t.FieldSize {
	case 1:

		// Write field offsets
		for _, offset := range t.Offsets {
			dst[pos] = byte(offset)
			pos++
		}
	case 2:
		// Set size enum
		dst[0] |= 64

		// Write field offsets
		for _, offset := range t.Offsets {
			binary.LittleEndian.PutUint16(dst[pos:], uint16(offset))
			pos += 2
		}
	case 4:
		// Set size enum
		dst[0] |= 128

		// Write field offsets
		for _, offset := range t.Offsets {
			binary.LittleEndian.PutUint32(dst[pos:], uint32(offset))
			pos += 4
		}
	case 8:
		// Set size enum
		dst[0] |= 192

		// Write field offsets
		for _, offset := range t.Offsets {
			binary.LittleEndian.PutUint64(dst[pos:], offset)
			pos += 8
		}
	default:
		return pos, errors.New("Invalid Header: Field size must be 1,2,4 or 8 bytes")
	}

	n, err := w.Write(dst)
	return int64(n), err
}
