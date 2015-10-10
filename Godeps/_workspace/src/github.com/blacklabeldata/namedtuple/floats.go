package namedtuple

import "github.com/blacklabeldata/xbinary"

// PutFloat32 writes a 32-bit float for the given string field. The field type must be `Float32Field`, otherwise an error is returned. The type code is written first then the value. Upon success, the number of bytes written is returned along with a nil error.
func (b *TupleBuilder) PutFloat32(field string, value float32) (wrote uint64, err error) {

	// field type should be
	if err = b.typeCheck(field, Float32Field); err != nil {
		return 0, err
	}

	// write value
	// length check performed by xbinary
	wrote, err = xbinary.LittleEndian.PutFloat32(b.buffer, b.pos+1, value)
	if err != nil {
		return 0, err
	}

	// write type code
	b.buffer[b.pos] = byte(FloatCode.OpCode)

	// set field offset
	b.offsets[field] = b.pos

	// incr pos
	b.pos += 5

	return 5, nil
}

// PutFloat64 writes a 64-bit float (or double in some languages) for the given string field. The field type must be `Float64Field`, otherwise an error is returned. The type code is written first then the value. Upon success, the number of bytes written is returned along with a nil error.
func (b *TupleBuilder) PutFloat64(field string, value float64) (wrote uint64, err error) {

	// field type should be
	if err = b.typeCheck(field, Float64Field); err != nil {
		return 0, err
	}

	// write value
	// length check performed by xbinary
	wrote, err = xbinary.LittleEndian.PutFloat64(b.buffer, b.pos+1, value)
	if err != nil {
		return 0, err
	}

	// write type code
	b.buffer[b.pos] = byte(DoubleCode.OpCode)

	// set field offset
	b.offsets[field] = b.pos

	// incr pos
	b.pos += 9

	return 9, nil
}
