package namedtuple

import (
	"testing"

	"github.com/blacklabeldata/xbinary"
	"github.com/stretchr/testify/assert"
)

// Float32
func TestPutFloat32Fail(t *testing.T) {

	// create test type
	// float test type
	TestType := New("testing", "float")
	TestType.AddVersion(
		Field{"float32", true, Float32Field},
		Field{"float64", true, Float64Field},
	)

	// create builder
	buffer := make([]byte, 1)
	builder := NewBuilder(TestType, buffer)

	// fails type check
	wrote, err := builder.PutFloat32("float64", float32(3.14159))
	// fmt.Println(err)
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), wrote)

	// fails length check
	wrote, err = builder.PutFloat32("float32", float32(3.14159))
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), wrote)
}

func TestPutFloat32Pass(t *testing.T) {

	// create test type
	// float test type
	TestType := New("testing", "float")
	TestType.AddVersion(
		Field{"float32", true, Float32Field},
		Field{"float64", true, Float64Field},
	)

	// create builder
	buffer := make([]byte, 5)
	builder := NewBuilder(TestType, buffer)

	// successful write
	wrote, err := builder.PutFloat32("float32", float32(3.14159))
	assert.Nil(t, err)
	assert.Equal(t, uint64(5), wrote)

	// test data validity
	assert.Equal(t, FloatCode.OpCode, uint8(builder.buffer[0]))

	value, err := xbinary.LittleEndian.Float32(buffer, 1)
	assert.Equal(t, float32(3.14159), value)

	// validate field offset
	assert.Equal(t, 0, builder.offsets["float32"])
}

// Float64
func TestPutFloat64Fail(t *testing.T) {

	// create test type
	// float test type
	TestType := New("testing", "float")
	TestType.AddVersion(
		Field{"float32", true, Float32Field},
		Field{"float64", true, Float64Field},
	)

	// create builder
	buffer := make([]byte, 1)
	builder := NewBuilder(TestType, buffer)

	// fails type check
	wrote, err := builder.PutFloat64("float32", float64(3.14159))
	// fmt.Println(err)
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), wrote)

	// fails length check
	wrote, err = builder.PutFloat64("float64", float64(3.14159))
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), wrote)
}

func TestPutFloat64Pass(t *testing.T) {

	// create test type
	// float test type
	TestType := New("testing", "float")
	TestType.AddVersion(
		Field{"float32", true, Float32Field},
		Field{"float64", true, Float64Field},
	)

	// create builder
	buffer := make([]byte, 9)
	builder := NewBuilder(TestType, buffer)

	// successful write
	wrote, err := builder.PutFloat64("float64", float64(3.14159))
	assert.Nil(t, err)
	assert.Equal(t, uint64(9), wrote)

	// test data validity
	assert.Equal(t, DoubleCode.OpCode, uint8(builder.buffer[0]))

	value, err := xbinary.LittleEndian.Float64(buffer, 1)
	assert.Equal(t, float64(3.14159), value)

	// validate field offset
	assert.Equal(t, 0, builder.offsets["float64"])
}
