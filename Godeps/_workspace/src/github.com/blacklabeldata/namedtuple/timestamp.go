package namedtuple

import (
	"math"
	"time"

	"github.com/blacklabeldata/xbinary"
)

// PutTimestamp writes a 64-bit signed integer (using `time.UnixNano()`) for the given `time.Time` value. The timestamp field type must be a `TimestampField`. If the buffer does not have enough space available an error is returned. Upon successful write, the number of bytes written will be returned as well as a nil error.
func (b *TupleBuilder) PutTimestamp(field string, value time.Time) (wrote uint64, err error) {

	// field type should be
	if err = b.typeCheck(field, TimestampField); err != nil {
		return 0, err
	}

	// write value
	// length check performed by xbinary
	wrote, err = xbinary.LittleEndian.PutInt64(b.buffer, b.pos+1, value.UnixNano())
	if err != nil {
		return 0, err
	}

	// write type code
	b.buffer[b.pos] = byte(TimestampCode.OpCode)

	// set field offset
	b.offsets[field] = b.pos

	// incr pos
	b.pos += 9

	// wrote 9 bytes
	return 9, nil
}

func (b *TupleBuilder) PutTimestampArray(field string, times []time.Time) (wrote int, err error) {

	// field type should be
	if err = b.typeCheck(field, TimestampArrayField); err != nil {
		return 0, err
	}

	// convert times to int64
	var value = make([]int64, len(times))
	for i := 0; i < len(times); i++ {
		value[i] = times[i].UnixNano()
	}

	size := len(value)
	if size < math.MaxUint8 {

		// write length
		if _, err = xbinary.LittleEndian.PutInt64Array(b.buffer, b.pos+2, value); err != nil {
			return 2, err
		}

		// write type code
		b.buffer[b.pos] = byte(TimestampArray8Code.OpCode)

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
		b.buffer[b.pos] = byte(TimestampArray16Code.OpCode)

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
		b.buffer[b.pos] = byte(TimestampArray32Code.OpCode)

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
		b.buffer[b.pos] = byte(TimestampArray64Code.OpCode)

		wrote += 9 + size
	}

	b.offsets[field] = b.pos
	b.pos += wrote
	return
}
