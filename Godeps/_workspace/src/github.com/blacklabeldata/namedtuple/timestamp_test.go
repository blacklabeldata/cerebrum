package namedtuple

import (
	"testing"
	"time"

	"github.com/blacklabeldata/xbinary"
	"github.com/stretchr/testify/assert"
)

// time testing
func TestPutTimestampFail(t *testing.T) {

	// create test type
	// float test type
	TestType := New("testing", "time")
	TestType.AddVersion(
		Field{"timestamp", true, TimestampField},
		Field{"float64", true, Float64Field},
	)

	// create builder
	buffer := make([]byte, 1)
	builder := NewBuilder(TestType, buffer)

	// fails type check
	wrote, err := builder.PutTimestamp("float64", time.Now())
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), wrote)

	// fails length check
	wrote, err = builder.PutTimestamp("timestamp", time.Now())
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), wrote)
}

func TestPutTimestampPass(t *testing.T) {

	// create test type
	// float test type
	TestType := New("testing", "time")
	TestType.AddVersion(
		Field{"timestamp", true, TimestampField},
		Field{"float64", true, Float64Field},
	)

	// create builder
	buffer := make([]byte, 9)
	builder := NewBuilder(TestType, buffer)

	// successful write
	now := time.Now()
	wrote, err := builder.PutTimestamp("timestamp", now)
	assert.Nil(t, err)
	assert.Equal(t, uint64(9), wrote)

	// test data validity
	assert.Equal(t, TimestampCode.OpCode, uint8(builder.buffer[0]))

	value, err := xbinary.LittleEndian.Int64(buffer, 1)
	assert.Equal(t, now.UnixNano(), value)

	// validate field offset
	assert.Equal(t, 0, builder.offsets["timestamp"])
}
