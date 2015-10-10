package namedtuple

import (
	"testing"

	"github.com/blacklabeldata/xbinary"
	"github.com/stretchr/testify/assert"
)

// String
func TestPutStringFail_1(t *testing.T) {
	// create test type
	// float test type
	TestType := New("testing", "string")
	TestType.AddVersion(
		Field{"string", true, StringField},
		Field{"bool", true, BooleanField},
	)

	// create builder
	buffer := make([]byte, 1)
	builder := NewBuilder(TestType, buffer)

	// fails type check
	wrote, err := builder.PutString("bool", "namedtuple")
	assert.NotNil(t, err)
	assert.Equal(t, 0, wrote)

	// fails length check
	wrote, err = builder.PutString("string", "namedtuple")
	assert.NotNil(t, err)
	assert.Equal(t, 0, wrote)
}

func TestPutStringPass_1(t *testing.T) {
	// create test type
	// float test type
	TestType := New("testing", "string")
	TestType.AddVersion(
		Field{"string", true, StringField},
		Field{"bool", true, BooleanField},
	)

	// create builder
	buffer := make([]byte, 12)
	builder := NewBuilder(TestType, buffer)

	// successful write
	wrote, err := builder.PutString("string", "namedtuple")
	assert.Nil(t, err)
	assert.Equal(t, 12, wrote)

	// test data validity
	assert.Equal(t, String8Code.OpCode, uint8(builder.buffer[0]))
	assert.Equal(t, 10, int(builder.buffer[1]))

	value, err := xbinary.LittleEndian.String(buffer, 2, 10)
	assert.Equal(t, "namedtuple", value)

	// validate field offset
	assert.Equal(t, 0, builder.offsets["string"])
}
func TestPutStringFail_2(t *testing.T) {
	// create test type
	// float test type
	TestType := New("testing", "string")
	TestType.AddVersion(
		Field{"string", true, StringField},
		Field{"bool", true, BooleanField},
	)

	// create builder
	buffer := make([]byte, 1)
	builder := NewBuilder(TestType, buffer)

	// fails length check
	wrote, err := builder.PutString("string", string(make([]byte, 300)))
	assert.NotNil(t, err)
	assert.Equal(t, 0, wrote)

	// create builder
	buffer = make([]byte, 3)
	builder = NewBuilder(TestType, buffer)

	// fails length check
	wrote, err = builder.PutString("string", string(make([]byte, 300)))
	assert.NotNil(t, err)
	assert.Equal(t, 0, wrote)
}

func TestPutStringPass_2(t *testing.T) {
	// create test type
	// float test type
	TestType := New("testing", "string")
	TestType.AddVersion(
		Field{"string", true, StringField},
		Field{"bool", true, BooleanField},
	)

	// create builder
	buffer := make([]byte, 303)
	builder := NewBuilder(TestType, buffer)

	// successful write
	input := string(make([]byte, 300))
	wrote, err := builder.PutString("string", input)
	assert.Nil(t, err)
	assert.Equal(t, 303, wrote)

	// test data validity
	assert.Equal(t, String16Code.OpCode, uint8(builder.buffer[0]))

	value, err := xbinary.LittleEndian.Uint16(buffer, 1)
	assert.Equal(t, 300, int(value))

	output, err := xbinary.LittleEndian.String(buffer, 3, int(value))
	assert.Equal(t, input, output)

	// validate field offset
	assert.Equal(t, 0, builder.offsets["string"])
}

func TestPutStringFail_3(t *testing.T) {
	// create test type
	// float test type
	TestType := New("testing", "string")
	TestType.AddVersion(
		Field{"string", true, StringField},
		Field{"bool", true, BooleanField},
	)

	// create builder
	buffer := make([]byte, 1)
	builder := NewBuilder(TestType, buffer)

	// fails length check
	wrote, err := builder.PutString("string", string(make([]byte, 135000)))
	assert.NotNil(t, err)
	assert.Equal(t, 0, wrote)

	// create builder
	buffer = make([]byte, 4)
	builder = NewBuilder(TestType, buffer)

	// fails length check
	wrote, err = builder.PutString("string", string(make([]byte, 135000)))
	assert.NotNil(t, err)
	assert.Equal(t, 0, wrote)
}

func TestPutStringPass_3(t *testing.T) {
	// create test type
	// float test type
	TestType := New("testing", "string")
	TestType.AddVersion(
		Field{"string", true, StringField},
		Field{"bool", true, BooleanField},
	)

	// create builder
	buffer := make([]byte, 135005)
	builder := NewBuilder(TestType, buffer)

	// successful write
	input := string(make([]byte, 135000))
	wrote, err := builder.PutString("string", input)
	assert.Nil(t, err)
	assert.Equal(t, 135005, wrote)

	// test data validity
	assert.Equal(t, String32Code.OpCode, uint8(builder.buffer[0]))

	value, err := xbinary.LittleEndian.Uint32(buffer, 1)
	assert.Equal(t, 135000, int(value))

	output, err := xbinary.LittleEndian.String(buffer, 5, int(value))
	assert.Equal(t, input, output)

	// validate field offset
	assert.Equal(t, 0, builder.offsets["string"])
}
