package namedtuple

import (
	"math"

	"github.com/blacklabeldata/xbinary"
)

// PutUint8 sets an 8-bit unsigned value for the given string name. The field name must be a Uint8Field otherwise an error will be returned. If the type buffer no longer has enough space to write this value an xbinary.ErrOutOfRange error will be returned. Upon success 2 bytes should be written into the buffer and the returned error should be nil. The type code is written first then the byte value.
func (b *TupleBuilder) PutUint8(field string, value uint8) (wrote uint64, err error) {

	// field type should be a Uint8Field
	if err = b.typeCheck(field, Uint8Field); err != nil {
		return 0, err
	}

	// minimum bytes is 2 (type code + value)
	if b.available() < 2 {
		return 0, xbinary.ErrOutOfRange
	}

	// write type code
	b.buffer[b.pos] = byte(UnsignedInt8Code.OpCode)

	// write value
	b.buffer[b.pos+1] = byte(value)

	// set field offset
	b.offsets[field] = b.pos

	// incr pos
	b.pos += 2

	return 2, nil
}

// PutInt8 sets an 8-bit signed value for the given string name. The field name must be an Int8Field otherwise an error will be returned. If the type buffer no longer has enough space to write this value, an xbinary.ErrOutOfRange error will be returned. Upon success, 2 bytes should be written into the buffer and the returned error should be nil. The type code is written first then the byte value.
func (b *TupleBuilder) PutInt8(field string, value int8) (wrote uint64, err error) {

	// field type should be
	if err = b.typeCheck(field, Int8Field); err != nil {
		return 0, err
	}

	// minimum bytes is 2 (type code + value)
	if b.available() < 2 {
		return 0, xbinary.ErrOutOfRange
	}

	// write type code
	b.buffer[b.pos] = byte(Int8Code.OpCode)

	// write value
	b.buffer[b.pos+1] = byte(value)

	// set field offset
	b.offsets[field] = b.pos

	// incr pos
	b.pos += 2

	return 2, nil
}

// PutUint16 sets a 16-bit unsigned value for the given field name.The field name must be a Uint16Field otherwise an error will be returned. If the type buffer no longer has enough space to write the value, an xbinary.ErrOutOfRange error will be returned. Upon success, the number of bytes written as well as a nil error will be returned. The type code will be writtn first. If the value is `< math.MaxUint8`, only 1 byte will be written. Otherwise, the entire 16-bit value will be written.
func (b *TupleBuilder) PutUint16(field string, value uint16) (wrote uint64, err error) {

	// field type should be
	if err = b.typeCheck(field, Uint16Field); err != nil {
		return 0, err
	}

	if value < math.MaxUint8 {

		// minimum bytes is 2 (type code + value)
		if b.available() < 2 {
			return 0, xbinary.ErrOutOfRange
		}

		// write type code
		b.buffer[b.pos] = byte(UnsignedShort8Code.OpCode)

		// write value
		b.buffer[b.pos+1] = byte(value)

		// set field offset
		b.offsets[field] = b.pos

		// incr pos
		b.pos += 2

		return 2, nil
	}

	// write value
	// length check performed by xbinary
	wrote, err = xbinary.LittleEndian.PutUint16(b.buffer, b.pos+1, value)
	if err != nil {
		return 0, err
	}

	// write type code
	b.buffer[b.pos] = byte(UnsignedShort16Code.OpCode)

	// set field offset
	b.offsets[field] = b.pos

	// incr pos
	b.pos += 3

	// wrote 3 bytes
	return 3, nil

}

// PutInt16 sets a 16-bit signed value for the given field name.The field name must be an Int16Field; otherwise, an error will be returned. If the type buffer no longer has enough space to write the value, an xbinary.ErrOutOfRange error will be returned. Upon success, the number of bytes written as well as a nil error will be returned. The type code will be written first. If the value is `< math.MaxUint8`, only 1 byte will be written. Otherwise, the entire 16-bit value will be written.
func (b *TupleBuilder) PutInt16(field string, value int16) (wrote uint64, err error) {

	// field type should be
	if err = b.typeCheck(field, Int16Field); err != nil {
		return 0, err
	}

	if uint16(value) < math.MaxUint8 {

		// minimum bytes is 2 (type code + value)
		if b.available() < 2 {
			return 0, xbinary.ErrOutOfRange
		}

		// write type code
		b.buffer[b.pos] = byte(Short8Code.OpCode)

		// write value
		b.buffer[b.pos+1] = byte(value)

		// set field offset
		b.offsets[field] = b.pos

		// incr pos
		b.pos += 2

		return 2, nil
	}

	// write value
	// length check performed by xbinary
	wrote, err = xbinary.LittleEndian.PutInt16(b.buffer, b.pos+1, value)
	if err != nil {
		return 0, err
	}

	// write type code
	b.buffer[b.pos] = byte(Short16Code.OpCode)

	// set field offset
	b.offsets[field] = b.pos

	// incr pos
	b.pos += 3

	// wrote 3 bytes
	return 3, nil

}

// PutUint32 sets a 32-bit unsigned value for the given field name. The field name must be a Uint32Field, otherwise, an error will be returned. If the type buffer no longer has enough space to write the value, an `xbinary.ErrOutOfRange` error will be returned. Upon success, the number of bytes written as well as a nil error will be returned. The type code will be written first. If the value is `< math.MaxUint8`, only 1 byte will be written. If the value is `< math.MaxUint16`, only 2 bytes will be written. Otherwise, the entire 32-bit value will be written.
func (b *TupleBuilder) PutUint32(field string, value uint32) (wrote uint64, err error) {

	// field type should be
	if err = b.typeCheck(field, Uint32Field); err != nil {
		return 0, err
	}

	if value < math.MaxUint8 {

		// minimum bytes is 2 (type code + value)
		if b.available() < 2 {
			return 0, xbinary.ErrOutOfRange
		}

		// write type code
		b.buffer[b.pos] = byte(UnsignedInt8Code.OpCode)

		// write value
		b.buffer[b.pos+1] = byte(value)

		// set field offset
		b.offsets[field] = b.pos

		// incr pos
		b.pos += 2

		return 2, nil
	} else if value < math.MaxUint16 {

		// write value
		// length check performed by xbinary
		wrote, err = xbinary.LittleEndian.PutUint16(b.buffer, b.pos+1, uint16(value))
		if err != nil {
			return 0, err
		}
		// write type code
		b.buffer[b.pos] = byte(UnsignedInt16Code.OpCode)

		// set field offset
		b.offsets[field] = b.pos

		// incr pos
		b.pos += 3

		// wrote 3 bytes
		return 3, nil
	}

	// write value
	// length check performed by xbinary
	wrote, err = xbinary.LittleEndian.PutUint32(b.buffer, b.pos+1, value)
	if err != nil {
		return 0, err
	}

	// write type code
	b.buffer[b.pos] = byte(UnsignedInt32Code.OpCode)

	// set field offset
	b.offsets[field] = b.pos

	// incr pos
	b.pos += 5

	// wrote 5 bytes
	return 5, nil
}

// PutInt32 sets a 32-bit signed value for the given field name. The field name must be a Int32Field. Otherwise, an error will be returned. If the type buffer no longer has enough space to write the value, an `xbinary.ErrOutOfRange` error will be returned. Upon success, the number of bytes written as well as a nil error will be returned. The type code will be written first. If the absolute value is `< math.MaxUint8`, only 1 byte will be written. If the absolute value is `< math.MaxUint16`, only 2 bytes will be written. Otherwise, the entire 32-bit value will be written.
func (b *TupleBuilder) PutInt32(field string, value int32) (wrote uint64, err error) {

	// field type should be
	if err = b.typeCheck(field, Int32Field); err != nil {
		return 0, err
	}

	unsigned := uint32(value)
	if unsigned < math.MaxUint8 {

		// minimum bytes is 2 (type code + value)
		if b.available() < 2 {
			return 0, xbinary.ErrOutOfRange
		}

		// write type code
		b.buffer[b.pos] = byte(Int8Code.OpCode)

		// write value
		b.buffer[b.pos+1] = byte(value)

		// set field offset
		b.offsets[field] = b.pos

		// incr pos
		b.pos += 2

		return 2, nil
	} else if unsigned < math.MaxUint16 {

		// write value
		// length check performed by xbinary
		wrote, err = xbinary.LittleEndian.PutInt16(b.buffer, b.pos+1, int16(value))
		if err != nil {
			return 0, err
		}

		// write type code
		b.buffer[b.pos] = byte(Int16Code.OpCode)

		// set field offset
		b.offsets[field] = b.pos

		// incr pos
		b.pos += 3

		// wrote 3 bytes
		return 3, nil
	}

	// write value
	// length check performed by xbinary
	wrote, err = xbinary.LittleEndian.PutInt32(b.buffer, b.pos+1, value)
	if err != nil {
		return 0, err
	}

	// write type code
	b.buffer[b.pos] = byte(Int32Code.OpCode)

	// set field offset
	b.offsets[field] = b.pos

	// incr pos
	b.pos += 5

	// wrote 5 bytes
	return 5, nil
}

// PutUint64 sets a 64-bit unsigned integer for the given field name. The field name must be a Uint64Field. Otherwise, an error will be returned. If the type buffer no longer has enough space to write the value, an `xbinary.ErrOutOfRange` error will be returned. Upon success, the number of bytes written as well as a nil error will be returned. The type code will be written first. If the absolute value is `< math.MaxUint8`, only 1 byte will be written. If the absolute value is `< math.MaxUint16`, only 2 bytes will be written. If the absolute value is `< math.MaxUint32`, only 4 bytes will be written. Otherwise, the entire 64-bit value will be written.
func (b *TupleBuilder) PutUint64(field string, value uint64) (wrote uint64, err error) {

	// field type should be
	if err = b.typeCheck(field, Uint64Field); err != nil {
		return 0, err
	}

	if value < math.MaxUint8 {

		// minimum bytes is 2 (type code + value)
		if b.available() < 2 {
			return 0, xbinary.ErrOutOfRange
		}

		// write type code
		b.buffer[b.pos] = byte(UnsignedLong8Code.OpCode)

		// write value
		b.buffer[b.pos+1] = byte(value)

		// set field offset
		b.offsets[field] = b.pos

		// incr pos
		b.pos += 2

		return 2, nil
	} else if value < math.MaxUint16 {

		// write value
		// length check performed by xbinary
		wrote, err = xbinary.LittleEndian.PutUint16(b.buffer, b.pos+1, uint16(value))
		if err != nil {
			return 0, err
		}

		// write type code
		b.buffer[b.pos] = byte(UnsignedLong16Code.OpCode)

		// set field offset
		b.offsets[field] = b.pos

		// incr pos
		b.pos += 3

		// wrote 3 bytes
		return 3, nil
	} else if value < math.MaxUint32 {

		// write value
		// length check performed by xbinary
		wrote, err = xbinary.LittleEndian.PutUint32(b.buffer, b.pos+1, uint32(value))
		if err != nil {
			return 0, err
		}

		// write type code
		b.buffer[b.pos] = byte(UnsignedLong32Code.OpCode)

		// set field offset
		b.offsets[field] = b.pos

		// incr pos
		b.pos += 5

		// wrote 5 bytes
		return 5, nil
	}

	// write value
	// length check performed by xbinary
	wrote, err = xbinary.LittleEndian.PutUint64(b.buffer, b.pos+1, value)
	if err != nil {
		return 0, err
	}
	// write type code
	b.buffer[b.pos] = byte(UnsignedLong64Code.OpCode)

	// set field offset
	b.offsets[field] = b.pos

	// incr pos
	b.pos += 9

	// wrote 9 bytes
	return 9, nil

}

// PutInt64 sets a 64-bit signed integer for the given field name. The field name must be a Int64Field. Otherwise, an error will be returned. If the type buffer no longer has enough space to write the value, an `xbinary.ErrOutOfRange` error will be returned. Upon success, the number of bytes written as well as a nil error will be returned. The type code will be written first. If the absolute value is `< math.MaxUint8`, only 1 byte will be written. If the absolute value is `< math.MaxUint16`, only 2 bytes will be written. If the absolute value is `< math.MaxUint32`, only 4 bytes will be written. Otherwise, the entire 64-bit value will be written.
func (b *TupleBuilder) PutInt64(field string, value int64) (wrote uint64, err error) {

	// field type should be
	if err = b.typeCheck(field, Int64Field); err != nil {
		return 0, err
	}

	unsigned := uint64(value)
	if unsigned < math.MaxUint8 {

		// minimum bytes is 2 (type code + value)
		if b.available() < 2 {
			return 0, xbinary.ErrOutOfRange
		}

		// write type code
		b.buffer[b.pos] = byte(Long8Code.OpCode)

		// write value
		b.buffer[b.pos+1] = byte(value)

		// set field offset
		b.offsets[field] = b.pos

		// incr pos
		b.pos += 2

		return 2, nil
	} else if unsigned < math.MaxUint16 {

		// write value
		// length check performed by xbinary
		wrote, err = xbinary.LittleEndian.PutInt16(b.buffer, b.pos+1, int16(value))
		if err != nil {
			return 0, err
		}

		// write type code
		b.buffer[b.pos] = byte(Long16Code.OpCode)

		// set field offset
		b.offsets[field] = b.pos

		// incr pos
		b.pos += 3

		// wrote 3 bytes
		return 3, nil
	} else if unsigned < math.MaxUint32 {

		// write value
		// length check performed by xbinary
		wrote, err = xbinary.LittleEndian.PutInt32(b.buffer, b.pos+1, int32(value))
		if err != nil {
			return 0, err
		}

		// write type code
		b.buffer[b.pos] = byte(Long32Code.OpCode)

		// set field offset
		b.offsets[field] = b.pos

		// incr pos
		b.pos += 5

		// wrote 5 bytes
		return 5, nil
	}

	// write value
	// length check performed by xbinary
	wrote, err = xbinary.LittleEndian.PutInt64(b.buffer, b.pos+1, value)
	if err != nil {
		return 0, err
	}

	// write type code
	b.buffer[b.pos] = byte(Long64Code.OpCode)

	// set field offset
	b.offsets[field] = b.pos

	// incr pos
	b.pos += 9

	// wrote 9 bytes
	return 9, nil
}
