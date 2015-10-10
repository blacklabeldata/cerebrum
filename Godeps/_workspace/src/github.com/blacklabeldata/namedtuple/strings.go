package namedtuple

import (
	"math"

	"github.com/blacklabeldata/xbinary"
)

// PutString writes a string value for the given field. The field type must be a `StringField` otherwise an error will be returned. The type code is written first, then a length, and finally the value. If the size of the string is `< math.MaxUint8`, a single byte will represent the length. If the size of the string is `< math.MaxUint16`, a uint16 value will represent the length and so on. If the buffer does not have enough space for the entire string and field header, an error will be returned. If successful, the number of bytes written will be returned as well as a nil error.
func (b *TupleBuilder) PutString(field string, value string) (wrote int, err error) {

	// field type should be
	if err = b.typeCheck(field, StringField); err != nil {
		return 0, err
	}

	size := len(value)
	if size < math.MaxUint8 {

		if b.available() < size+2 {
			return 0, xbinary.ErrOutOfRange
		}

		// write length
		xbinary.LittleEndian.PutString(b.buffer, b.pos+2, value)

		// write type code
		b.buffer[b.pos] = byte(String8Code.OpCode)

		// write length
		b.buffer[b.pos+1] = byte(size)

		wrote += size + 2
	} else if size < math.MaxUint16 {

		if b.available() < size+3 {
			return 0, xbinary.ErrOutOfRange
		}

		// write length
		xbinary.LittleEndian.PutUint16(b.buffer, b.pos+1, uint16(size))

		// write value
		xbinary.LittleEndian.PutString(b.buffer, b.pos+3, value)

		// write type code
		b.buffer[b.pos] = byte(String16Code.OpCode)

		wrote += 3 + size
	} else if size < math.MaxUint32 {

		if b.available() < size+5 {
			return 0, xbinary.ErrOutOfRange
		}

		// write length
		xbinary.LittleEndian.PutUint32(b.buffer, b.pos+1, uint32(size))

		// write value
		xbinary.LittleEndian.PutString(b.buffer, b.pos+5, value)

		// write type code
		b.buffer[b.pos] = byte(String32Code.OpCode)

		wrote += 5 + size
	} else {

		if b.available() < size+9 {
			return 0, xbinary.ErrOutOfRange
		}

		// write length
		xbinary.LittleEndian.PutUint64(b.buffer, b.pos+1, uint64(size))

		// write value
		xbinary.LittleEndian.PutString(b.buffer, b.pos+9, value)

		// write type code
		b.buffer[b.pos] = byte(String64Code.OpCode)

		wrote += 9 + size
	}

	b.offsets[field] = b.pos
	b.pos += wrote
	return
}
