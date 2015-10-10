package namedtuple

import (
	"math"

	"github.com/blacklabeldata/xbinary"
)

func (b *TupleBuilder) PutUint8Array(field string, value []uint8) (wrote int, err error) {

	// field type should be
	if err = b.typeCheck(field, Uint8ArrayField); err != nil {
		return 0, err
	}

	size := len(value)
	if size < math.MaxUint8 {

		// write data
		if _, err = xbinary.LittleEndian.PutUint8Array(b.buffer, b.pos+2, value); err != nil {
			return 2, err
		}

		// write type code
		b.buffer[b.pos] = byte(UnsignedByteArray8Code.OpCode)

		// write length
		b.buffer[b.pos+1] = byte(size)

		wrote += size + 2
	} else if size < math.MaxUint16 {

		// write length
		if _, err = xbinary.LittleEndian.PutUint16(b.buffer, b.pos+1, uint16(size)); err != nil {
			return 1, err
		}

		// write value
		if _, err = xbinary.LittleEndian.PutUint8Array(b.buffer, b.pos+3, value); err != nil {
			return 3, err
		}
		// write type code
		b.buffer[b.pos] = byte(UnsignedByteArray16Code.OpCode)

		wrote += 3 + size
	} else if size < math.MaxUint32 {

		// write length
		if _, err = xbinary.LittleEndian.PutUint32(b.buffer, b.pos+1, uint32(size)); err != nil {
			return 1, err
		}

		// write value
		if _, err = xbinary.LittleEndian.PutUint8Array(b.buffer, b.pos+5, value); err != nil {
			return 5, err
		}
		// write type code
		b.buffer[b.pos] = byte(UnsignedByteArray32Code.OpCode)

		wrote += 5 + size
	} else {

		// write length
		if _, err = xbinary.LittleEndian.PutUint64(b.buffer, b.pos+1, uint64(size)); err != nil {
			return 1, err
		}

		// write value
		if _, err = xbinary.LittleEndian.PutUint8Array(b.buffer, b.pos+9, value); err != nil {
			return 9, err
		}
		// write type code
		b.buffer[b.pos] = byte(UnsignedByteArray64Code.OpCode)

		wrote += 9 + size
	}

	b.offsets[field] = b.pos
	b.pos += wrote
	return
}

func (b *TupleBuilder) PutInt8Array(field string, value []int8) (wrote int, err error) {

	// field type should be
	if err = b.typeCheck(field, Int8ArrayField); err != nil {
		return 0, err
	}

	size := len(value)
	if size < math.MaxUint8 {

		// write data
		if _, err = xbinary.LittleEndian.PutInt8Array(b.buffer, b.pos+2, value); err != nil {
			return 2, err
		}

		// write type code
		b.buffer[b.pos] = byte(ByteArray8Code.OpCode)

		// write length
		b.buffer[b.pos+1] = byte(size)

		wrote += size + 2
	} else if size < math.MaxUint16 {

		// write length
		if _, err = xbinary.LittleEndian.PutUint16(b.buffer, b.pos+1, uint16(size)); err != nil {
			return 1, err
		}

		// write value
		if _, err = xbinary.LittleEndian.PutInt8Array(b.buffer, b.pos+3, value); err != nil {
			return 3, err
		}
		// write type code
		b.buffer[b.pos] = byte(ByteArray16Code.OpCode)

		wrote += 3 + size
	} else if size < math.MaxUint32 {

		// write length
		if _, err = xbinary.LittleEndian.PutUint32(b.buffer, b.pos+1, uint32(size)); err != nil {
			return 1, err
		}

		// write value
		if _, err = xbinary.LittleEndian.PutInt8Array(b.buffer, b.pos+5, value); err != nil {
			return 5, err
		}
		// write type code
		b.buffer[b.pos] = byte(ByteArray32Code.OpCode)

		wrote += 5 + size
	} else {

		// write length
		if _, err = xbinary.LittleEndian.PutUint64(b.buffer, b.pos+1, uint64(size)); err != nil {
			return 1, err
		}

		// write value
		if _, err = xbinary.LittleEndian.PutInt8Array(b.buffer, b.pos+9, value); err != nil {
			return 9, err
		}
		// write type code
		b.buffer[b.pos] = byte(ByteArray64Code.OpCode)

		wrote += 9 + size
	}

	b.offsets[field] = b.pos
	b.pos += wrote
	return
}

func (b *TupleBuilder) PutUint16Array(field string, value []uint16) (wrote int, err error) {

	// field type should be
	if err = b.typeCheck(field, Uint16ArrayField); err != nil {
		return 0, err
	}

	size := len(value)
	if size < math.MaxUint8 {

		// write length
		if _, err = xbinary.LittleEndian.PutUint16Array(b.buffer, b.pos+2, value); err != nil {
			return 2, err
		}

		// write type code
		b.buffer[b.pos] = byte(UnsignedShortArray8Code.OpCode)

		// write length
		b.buffer[b.pos+1] = byte(size)

		wrote += size + 2
	} else if size < math.MaxUint16 {

		// write length
		if _, err = xbinary.LittleEndian.PutUint16(b.buffer, b.pos+1, uint16(size)); err != nil {
			return 1, err
		}

		// write value
		if _, err = xbinary.LittleEndian.PutUint16Array(b.buffer, b.pos+3, value); err != nil {
			return 3, err
		}
		// write type code
		b.buffer[b.pos] = byte(UnsignedShortArray16Code.OpCode)

		wrote += 3 + size
	} else if size < math.MaxUint32 {

		// write length
		if _, err = xbinary.LittleEndian.PutUint32(b.buffer, b.pos+1, uint32(size)); err != nil {
			return 1, err
		}

		// write value
		if _, err = xbinary.LittleEndian.PutUint16Array(b.buffer, b.pos+5, value); err != nil {
			return 5, err
		}
		// write type code
		b.buffer[b.pos] = byte(UnsignedShortArray32Code.OpCode)

		wrote += 5 + size
	} else {

		// write length
		if _, err = xbinary.LittleEndian.PutUint64(b.buffer, b.pos+1, uint64(size)); err != nil {
			return 1, err
		}

		// write value
		if _, err = xbinary.LittleEndian.PutUint16Array(b.buffer, b.pos+9, value); err != nil {
			return 9, err
		}
		// write type code
		b.buffer[b.pos] = byte(UnsignedShortArray64Code.OpCode)

		wrote += 9 + size
	}

	b.offsets[field] = b.pos
	b.pos += wrote
	return
}

func (b *TupleBuilder) PutInt16Array(field string, value []int16) (wrote int, err error) {

	// field type should be
	if err = b.typeCheck(field, Int16ArrayField); err != nil {
		return 0, err
	}

	size := len(value)
	if size < math.MaxUint8 {

		// write length
		if _, err = xbinary.LittleEndian.PutInt16Array(b.buffer, b.pos+2, value); err != nil {
			return 2, err
		}

		// write type code
		b.buffer[b.pos] = byte(ShortArray8Code.OpCode)

		// write length
		b.buffer[b.pos+1] = byte(size)

		wrote += size + 2
	} else if size < math.MaxUint16 {

		// write length
		if _, err = xbinary.LittleEndian.PutUint16(b.buffer, b.pos+1, uint16(size)); err != nil {
			return 1, err
		}

		// write value
		if _, err = xbinary.LittleEndian.PutInt16Array(b.buffer, b.pos+3, value); err != nil {
			return 3, err
		}
		// write type code
		b.buffer[b.pos] = byte(ShortArray16Code.OpCode)

		wrote += 3 + size
	} else if size < math.MaxUint32 {

		// write length
		if _, err = xbinary.LittleEndian.PutUint32(b.buffer, b.pos+1, uint32(size)); err != nil {
			return 1, err
		}

		// write value
		if _, err = xbinary.LittleEndian.PutInt16Array(b.buffer, b.pos+5, value); err != nil {
			return 5, err
		}
		// write type code
		b.buffer[b.pos] = byte(ShortArray32Code.OpCode)

		wrote += 5 + size
	} else {

		// write length
		if _, err = xbinary.LittleEndian.PutUint64(b.buffer, b.pos+1, uint64(size)); err != nil {
			return 1, err
		}

		// write value
		if _, err = xbinary.LittleEndian.PutInt16Array(b.buffer, b.pos+9, value); err != nil {
			return 9, err
		}
		// write type code
		b.buffer[b.pos] = byte(ShortArray64Code.OpCode)

		wrote += 9 + size
	}

	b.offsets[field] = b.pos
	b.pos += wrote
	return
}

func (b *TupleBuilder) PutUint32Array(field string, value []uint32) (wrote int, err error) {

	// field type should be
	if err = b.typeCheck(field, Uint32ArrayField); err != nil {
		return 0, err
	}

	size := len(value)
	if size < math.MaxUint8 {

		// write length
		if _, err = xbinary.LittleEndian.PutUint32Array(b.buffer, b.pos+2, value); err != nil {
			return 2, err
		}

		// write type code
		b.buffer[b.pos] = byte(UnsignedIntArray8Code.OpCode)

		// write length
		b.buffer[b.pos+1] = byte(size)

		wrote += size + 2
	} else if size < math.MaxUint16 {

		// write length
		if _, err = xbinary.LittleEndian.PutUint16(b.buffer, b.pos+1, uint16(size)); err != nil {
			return 1, err
		}

		// write value
		if _, err = xbinary.LittleEndian.PutUint32Array(b.buffer, b.pos+3, value); err != nil {
			return 3, err
		}
		// write type code
		b.buffer[b.pos] = byte(UnsignedIntArray16Code.OpCode)

		wrote += 3 + size
	} else if size < math.MaxUint32 {

		// write length
		if _, err = xbinary.LittleEndian.PutUint32(b.buffer, b.pos+1, uint32(size)); err != nil {
			return 1, err
		}

		// write value
		if _, err = xbinary.LittleEndian.PutUint32Array(b.buffer, b.pos+5, value); err != nil {
			return 5, err
		}
		// write type code
		b.buffer[b.pos] = byte(UnsignedIntArray32Code.OpCode)

		wrote += 5 + size
	} else {

		// write length
		if _, err = xbinary.LittleEndian.PutUint64(b.buffer, b.pos+1, uint64(size)); err != nil {
			return 1, err
		}

		// write value
		if _, err = xbinary.LittleEndian.PutUint32Array(b.buffer, b.pos+9, value); err != nil {
			return 9, err
		}
		// write type code
		b.buffer[b.pos] = byte(UnsignedIntArray64Code.OpCode)

		wrote += 9 + size
	}

	b.offsets[field] = b.pos
	b.pos += wrote
	return
}

func (b *TupleBuilder) PutInt32Array(field string, value []int32) (wrote int, err error) {

	// field type should be
	if err = b.typeCheck(field, Int32ArrayField); err != nil {
		return 0, err
	}

	size := len(value)
	if size < math.MaxUint8 {

		// write length
		if _, err = xbinary.LittleEndian.PutInt32Array(b.buffer, b.pos+2, value); err != nil {
			return 2, err
		}

		// write type code
		b.buffer[b.pos] = byte(IntArray8Code.OpCode)

		// write length
		b.buffer[b.pos+1] = byte(size)

		wrote += size + 2
	} else if size < math.MaxUint16 {

		// write length
		if _, err = xbinary.LittleEndian.PutUint16(b.buffer, b.pos+1, uint16(size)); err != nil {
			return 1, err
		}

		// write value
		if _, err = xbinary.LittleEndian.PutInt32Array(b.buffer, b.pos+3, value); err != nil {
			return 3, err
		}
		// write type code
		b.buffer[b.pos] = byte(IntArray16Code.OpCode)

		wrote += 3 + size
	} else if size < math.MaxUint32 {

		// write length
		if _, err = xbinary.LittleEndian.PutUint32(b.buffer, b.pos+1, uint32(size)); err != nil {
			return 1, err
		}

		// write value
		if _, err = xbinary.LittleEndian.PutInt32Array(b.buffer, b.pos+5, value); err != nil {
			return 5, err
		}
		// write type code
		b.buffer[b.pos] = byte(IntArray32Code.OpCode)

		wrote += 5 + size
	} else {

		// write length
		if _, err = xbinary.LittleEndian.PutUint64(b.buffer, b.pos+1, uint64(size)); err != nil {
			return 1, err
		}

		// write value
		if _, err = xbinary.LittleEndian.PutInt32Array(b.buffer, b.pos+9, value); err != nil {
			return 9, err
		}
		// write type code
		b.buffer[b.pos] = byte(IntArray64Code.OpCode)

		wrote += 9 + size
	}

	b.offsets[field] = b.pos
	b.pos += wrote
	return
}

func (b *TupleBuilder) PutUint64Array(field string, value []uint64) (wrote int, err error) {

	// field type should be
	if err = b.typeCheck(field, Uint64ArrayField); err != nil {
		return 0, err
	}

	size := len(value)
	if size < math.MaxUint8 {

		// write data
		if _, err = xbinary.LittleEndian.PutUint64Array(b.buffer, b.pos+2, value); err != nil {
			return 2, err
		}

		// write type code
		b.buffer[b.pos] = byte(UnsignedLongArray8Code.OpCode)

		// write length
		b.buffer[b.pos+1] = byte(size)

		wrote += size + 2
	} else if size < math.MaxUint16 {

		// write length
		if _, err = xbinary.LittleEndian.PutUint16(b.buffer, b.pos+1, uint16(size)); err != nil {
			return 1, err
		}

		// write value
		if _, err = xbinary.LittleEndian.PutUint64Array(b.buffer, b.pos+3, value); err != nil {
			return 3, err
		}

		// write type code
		b.buffer[b.pos] = byte(UnsignedLongArray16Code.OpCode)

		wrote += 3 + size
	} else if size < math.MaxUint32 {

		// write length
		if _, err = xbinary.LittleEndian.PutUint32(b.buffer, b.pos+1, uint32(size)); err != nil {
			return 1, err
		}

		// write value
		if _, err = xbinary.LittleEndian.PutUint64Array(b.buffer, b.pos+5, value); err != nil {
			return 5, err
		}

		// write type code
		b.buffer[b.pos] = byte(UnsignedLongArray32Code.OpCode)

		wrote += 5 + size
	} else {
		// write length
		if _, err = xbinary.LittleEndian.PutUint64(b.buffer, b.pos+1, uint64(size)); err != nil {
			return 1, err
		}

		// write value
		if _, err = xbinary.LittleEndian.PutUint64Array(b.buffer, b.pos+9, value); err != nil {
			return 9, err
		}

		// write type code
		b.buffer[b.pos] = byte(UnsignedLongArray64Code.OpCode)

		wrote += 9 + size
	}

	b.offsets[field] = b.pos
	b.pos += wrote
	return
}

func (b *TupleBuilder) PutInt64Array(field string, value []int64) (wrote int, err error) {

	// field type should be
	if err = b.typeCheck(field, Int64ArrayField); err != nil {
		return 0, err
	}

	size := len(value)
	if size < math.MaxUint8 {

		// write length
		if _, err = xbinary.LittleEndian.PutInt64Array(b.buffer, b.pos+2, value); err != nil {
			return 2, err
		}

		// write type code
		b.buffer[b.pos] = byte(LongArray8Code.OpCode)

		// write length
		b.buffer[b.pos+1] = byte(size)

		wrote += size + 2
	} else if size < math.MaxUint16 {

		// write length
		if _, err = xbinary.LittleEndian.PutUint16(b.buffer, b.pos+1, uint16(size)); err != nil {
			return 1, err
		}

		// write value
		if _, err = xbinary.LittleEndian.PutInt64Array(b.buffer, b.pos+3, value); err != nil {
			return 3, err
		}

		// write type code
		b.buffer[b.pos] = byte(LongArray16Code.OpCode)

		wrote += 3 + size
	} else if size < math.MaxUint32 {

		// write length
		if _, err = xbinary.LittleEndian.PutUint32(b.buffer, b.pos+1, uint32(size)); err != nil {
			return 1, err
		}

		// write value
		if _, err = xbinary.LittleEndian.PutInt64Array(b.buffer, b.pos+5, value); err != nil {
			return 5, err
		}

		// write type code
		b.buffer[b.pos] = byte(LongArray32Code.OpCode)

		wrote += 5 + size
	} else {

		// write length
		if _, err = xbinary.LittleEndian.PutUint64(b.buffer, b.pos+1, uint64(size)); err != nil {
			return 1, err
		}

		// write value
		if _, err = xbinary.LittleEndian.PutInt64Array(b.buffer, b.pos+9, value); err != nil {
			return 9, err
		}

		// write type code
		b.buffer[b.pos] = byte(LongArray64Code.OpCode)

		wrote += 9 + size
	}

	b.offsets[field] = b.pos
	b.pos += wrote
	return
}
