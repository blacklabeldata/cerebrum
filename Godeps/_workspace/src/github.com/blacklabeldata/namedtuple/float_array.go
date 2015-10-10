package namedtuple

import (
	"math"

	"github.com/blacklabeldata/xbinary"
)

// PutFloat32Array writes a float array for the given field. The field type must be a 'Float32ArrayField', otherwise as error will be returned. The type code is written first followed by the array size in bytes. If the size of the array is less than `math.MaxUint8`, a byte will be used to represent the length. If the size of the array is less than `math.MaxUint16`, a 16-bit unsigned integer will be used to represent the length and so on. If the buffer is too small to store the entire array, an `xbinary.ErrOutOfRange` error will be returned. If the write is successful, the total number of bytes will be returned as well as a nil error.
func (b *TupleBuilder) PutFloat32Array(field string, value []float32) (wrote int, err error) {

	// field type should be
	if err = b.typeCheck(field, Float32ArrayField); err != nil {
		return 0, err
	}

	size := len(value)
	if size < math.MaxUint8 {

		if b.available() < size*4+2 {
			return 0, xbinary.ErrOutOfRange
		}

		// write length
		xbinary.LittleEndian.PutFloat32Array(b.buffer, b.pos+2, value)

		// write type code
		b.buffer[b.pos] = byte(FloatArray8Code.OpCode)

		// write length
		b.buffer[b.pos+1] = byte(size)

		wrote += size + 2
	} else if size < math.MaxUint16 {

		if b.available() < size*4+3 {
			return 0, xbinary.ErrOutOfRange
		}

		// write length
		xbinary.LittleEndian.PutUint16(b.buffer, b.pos+1, uint16(size))

		// write value
		xbinary.LittleEndian.PutFloat32Array(b.buffer, b.pos+3, value)

		// write type code
		b.buffer[b.pos] = byte(FloatArray16Code.OpCode)

		wrote += 3 + size
	} else if size < math.MaxUint32 {

		if b.available() < size*4+5 {
			return 0, xbinary.ErrOutOfRange
		}

		// write length
		xbinary.LittleEndian.PutUint32(b.buffer, b.pos+1, uint32(size))

		// write value
		xbinary.LittleEndian.PutFloat32Array(b.buffer, b.pos+5, value)

		// write type code
		b.buffer[b.pos] = byte(FloatArray32Code.OpCode)

		wrote += 5 + size
	} else {

		if b.available() < size*4+9 {
			return 0, xbinary.ErrOutOfRange
		}

		// write length
		xbinary.LittleEndian.PutUint64(b.buffer, b.pos+1, uint64(size))

		// write value
		xbinary.LittleEndian.PutFloat32Array(b.buffer, b.pos+9, value)

		// write type code
		b.buffer[b.pos] = byte(FloatArray64Code.OpCode)

		wrote += 9 + size
	}

	b.offsets[field] = b.pos
	b.pos += wrote
	return
}

// PutFloat64Array writes a float array for the given field. The field type must be a 'Float64ArrayField', otherwise as error will be returned. The type code is written first followed by the array size in bytes. If the size of the array is less than `math.MaxUint8`, a byte will be used to represent the length. If the size of the array is less than `math.MaxUint16`, a 16-bit unsigned integer will be used to represent the length and so on. If the buffer is too small to store the entire array, an `xbinary.ErrOutOfRange` error will be returned. If the write is successful, the total number of bytes will be returned as well as a nil error.
func (b *TupleBuilder) PutFloat64Array(field string, value []float64) (wrote int, err error) {

	// field type should be
	if err = b.typeCheck(field, Float64ArrayField); err != nil {
		return 0, err
	}

	size := len(value)
	if size < math.MaxUint8 {

		if b.available() < size*8+2 {
			return 0, xbinary.ErrOutOfRange
		}

		// write length
		xbinary.LittleEndian.PutFloat64Array(b.buffer, b.pos+2, value)

		// write type code
		b.buffer[b.pos] = byte(DoubleArray8Code.OpCode)

		// write length
		b.buffer[b.pos+1] = byte(size)

		wrote += size + 2
	} else if size < math.MaxUint16 {

		if b.available() < size*8+3 {
			return 0, xbinary.ErrOutOfRange
		}

		// write length
		xbinary.LittleEndian.PutUint16(b.buffer, b.pos+1, uint16(size))

		// write value
		xbinary.LittleEndian.PutFloat64Array(b.buffer, b.pos+3, value)

		// write type code
		b.buffer[b.pos] = byte(DoubleArray16Code.OpCode)

		wrote += 3 + size
	} else if size < math.MaxUint32 {

		if b.available() < size*8+5 {
			return 0, xbinary.ErrOutOfRange
		}

		// write length
		xbinary.LittleEndian.PutUint32(b.buffer, b.pos+1, uint32(size))

		// write value
		xbinary.LittleEndian.PutFloat64Array(b.buffer, b.pos+5, value)

		// write type code
		b.buffer[b.pos] = byte(DoubleArray32Code.OpCode)

		wrote += 5 + size
	} else {

		if b.available() < size*8+9 {
			return 0, xbinary.ErrOutOfRange
		}

		// write length
		xbinary.LittleEndian.PutUint64(b.buffer, b.pos+1, uint64(size))

		// write value
		xbinary.LittleEndian.PutFloat64Array(b.buffer, b.pos+9, value)

		// write type code
		b.buffer[b.pos] = byte(DoubleArray64Code.OpCode)

		wrote += 9 + size
	}

	b.offsets[field] = b.pos
	b.pos += wrote
	return
}
