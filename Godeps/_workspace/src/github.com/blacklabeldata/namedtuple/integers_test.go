package namedtuple

import (
	"testing"

	"github.com/blacklabeldata/xbinary"
	"github.com/stretchr/testify/assert"
)

func TestBuilderPutUint8Fail(t *testing.T) {
	// create test type
	User := createTestTupleType()

	// create builder
	buffer := make([]byte, 1)
	builder := NewBuilder(User, buffer)

	// fails type check
	wrote, err := builder.PutUint8("uuid", uint8(20))
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), wrote)

	// fails length check
	wrote, err = builder.PutUint8("age", uint8(20))
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), wrote)
}

func TestBuilderPutUint8Pass(t *testing.T) {
	// create test type
	User := createTestTupleType()

	// create builder
	buffer := make([]byte, 2)
	builder := NewBuilder(User, buffer)

	// successful write
	wrote, err := builder.PutUint8("age", uint8(20))
	assert.Nil(t, err)
	assert.Equal(t, uint64(2), wrote)

	// test data validity
	assert.Equal(t, UnsignedInt8Code.OpCode, uint8(builder.buffer[0]))
	assert.Equal(t, 20, int(builder.buffer[1]))

	// validate field offset
	assert.Equal(t, 0, builder.offsets["age"])
}

func TestBuilderPutInt8Fail(t *testing.T) {

	// create test type
	// integer test type
	TestType := New("testing", "int8")
	TestType.AddVersion(
		Field{"int8", true, Int8Field},
		Field{"uint8", true, Uint8Field},
	)

	// create builder
	buffer := make([]byte, 1)
	builder := NewBuilder(TestType, buffer)

	// fails type check
	wrote, err := builder.PutInt8("uint8", int8(20))
	// fmt.Println(err)
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), wrote)

	// fails length check
	wrote, err = builder.PutInt8("int8", int8(20))
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), wrote)
}

func TestBuilderPutInt8Pass(t *testing.T) {
	// create test type
	// integer test type
	TestType := New("testing", "int8")
	TestType.AddVersion(
		Field{"int8", true, Int8Field},
		Field{"uint8", true, Uint8Field},
	)

	// create builder
	buffer := make([]byte, 2)
	builder := NewBuilder(TestType, buffer)

	// successful write
	wrote, err := builder.PutInt8("int8", int8(20))
	assert.Nil(t, err)
	assert.Equal(t, uint64(2), wrote)

	// test data validity
	assert.Equal(t, Int8Code.OpCode, uint8(builder.buffer[0]))
	assert.Equal(t, 20, int(builder.buffer[1]))

	// validate field offset
	assert.Equal(t, 0, builder.offsets["int8"])
}

func TestBuilderPutUint16Fail_1(t *testing.T) {

	// create test type
	// integer test type
	TestType := New("testing", "uint16")
	TestType.AddVersion(
		Field{"int16", true, Int16Field},
		Field{"uint16", true, Uint16Field},
	)

	// create builder
	buffer := make([]byte, 1)
	builder := NewBuilder(TestType, buffer)

	// fails type check
	wrote, err := builder.PutUint16("int16", uint16(20))
	// fmt.Println(err)
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), wrote)

	// fails length check
	wrote, err = builder.PutUint16("uint16", uint16(20))
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), wrote)
}

func TestBuilderPutUint16Pass_1(t *testing.T) {

	// create test type
	// integer test type
	TestType := New("testing", "uint16")
	TestType.AddVersion(
		Field{"int16", true, Int16Field},
		Field{"uint16", true, Uint16Field},
	)

	// create builder
	buffer := make([]byte, 2)
	builder := NewBuilder(TestType, buffer)

	// successful write
	wrote, err := builder.PutUint16("uint16", uint16(20))
	assert.Nil(t, err)
	assert.Equal(t, uint64(2), wrote)

	// test data validity
	assert.Equal(t, UnsignedShort8Code.OpCode, uint8(builder.buffer[0]))
	assert.Equal(t, 20, int(builder.buffer[1]))

	// validate field offset
	assert.Equal(t, 0, builder.offsets["uint16"])
}

func TestBuilderPutUint16Fail_2(t *testing.T) {

	// create test type
	// integer test type
	TestType := New("testing", "uint16")
	TestType.AddVersion(
		Field{"int16", true, Int16Field},
		Field{"uint16", true, Uint16Field},
	)

	// create builder
	buffer := make([]byte, 1)
	builder := NewBuilder(TestType, buffer)

	// fails length check
	wrote, err := builder.PutUint16("uint16", uint16(300))
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), wrote)
}

func TestBuilderPutUint16Pass_2(t *testing.T) {

	// create test type
	// integer test type
	TestType := New("testing", "uint16")
	TestType.AddVersion(
		Field{"int16", true, Int16Field},
		Field{"uint16", true, Uint16Field},
	)

	// create builder
	buffer := make([]byte, 3)
	builder := NewBuilder(TestType, buffer)

	// successful write
	wrote, err := builder.PutUint16("uint16", uint16(300))
	assert.Nil(t, err)
	assert.Equal(t, uint64(3), wrote)

	// test data validity
	assert.Equal(t, UnsignedShort16Code.OpCode, uint8(builder.buffer[0]))

	value, err := xbinary.LittleEndian.Uint16(buffer, 1)
	assert.Equal(t, uint16(300), value)

	// validate field offset
	assert.Equal(t, 0, builder.offsets["uint16"])
}

//
func TestBuilderPutInt16Fail_1(t *testing.T) {

	// create test type
	// integer test type
	TestType := New("testing", "int16")
	TestType.AddVersion(
		Field{"int16", true, Int16Field},
		Field{"uint16", true, Uint16Field},
	)

	// create builder
	buffer := make([]byte, 1)
	builder := NewBuilder(TestType, buffer)

	// fails type check
	wrote, err := builder.PutInt16("uint16", int16(20))
	// fmt.Println(err)
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), wrote)

	// fails length check
	wrote, err = builder.PutInt16("int16", int16(20))
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), wrote)
}

func TestBuilderPutInt16Pass_1(t *testing.T) {

	// create test type
	// integer test type
	TestType := New("testing", "int16")
	TestType.AddVersion(
		Field{"int16", true, Int16Field},
		Field{"uint16", true, Uint16Field},
	)

	// create builder
	buffer := make([]byte, 2)
	builder := NewBuilder(TestType, buffer)

	// successful write
	wrote, err := builder.PutInt16("int16", int16(20))
	assert.Nil(t, err)
	assert.Equal(t, uint64(2), wrote)

	// test data validity
	assert.Equal(t, Short8Code.OpCode, uint8(builder.buffer[0]))
	assert.Equal(t, 20, int(builder.buffer[1]))

	// validate field offset
	assert.Equal(t, 0, builder.offsets["int16"])
}

func TestBuilderPutInt16Fail_2(t *testing.T) {

	// create test type
	// integer test type
	TestType := New("testing", "int16")
	TestType.AddVersion(
		Field{"int16", true, Int16Field},
		Field{"uint16", true, Uint16Field},
	)

	// create builder
	buffer := make([]byte, 1)
	builder := NewBuilder(TestType, buffer)

	// fails length check
	wrote, err := builder.PutInt16("int16", int16(300))
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), wrote)
}

func TestBuilderPutInt16Pass_2(t *testing.T) {

	// create test type
	// integer test type
	TestType := New("testing", "int16")
	TestType.AddVersion(
		Field{"int16", true, Int16Field},
		Field{"uint16", true, Uint16Field},
	)

	// create builder
	buffer := make([]byte, 3)
	builder := NewBuilder(TestType, buffer)

	// successful write
	wrote, err := builder.PutInt16("int16", int16(-300))
	assert.Nil(t, err)
	assert.Equal(t, uint64(3), wrote)

	// test data validity
	assert.Equal(t, Short16Code.OpCode, uint8(builder.buffer[0]))

	value, err := xbinary.LittleEndian.Int16(buffer, 1)
	assert.Equal(t, int16(-300), value)

	// validate field offset
	assert.Equal(t, 0, builder.offsets["int16"])
}

//
func TestBuilderPutUint32Fail_1(t *testing.T) {

	// create test type
	// integer test type
	TestType := New("testing", "uint32")
	TestType.AddVersion(
		Field{"int32", true, Int32Field},
		Field{"uint32", true, Uint32Field},
	)

	// create builder
	buffer := make([]byte, 1)
	builder := NewBuilder(TestType, buffer)

	// fails type check
	wrote, err := builder.PutUint32("int32", uint32(20))
	// fmt.Println(err)
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), wrote)

	// fails length check
	wrote, err = builder.PutUint32("uint32", uint32(20))
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), wrote)
}

func TestBuilderPutUint32Pass_1(t *testing.T) {

	// create test type
	// integer test type
	TestType := New("testing", "uint32")
	TestType.AddVersion(
		Field{"int32", true, Int32Field},
		Field{"uint32", true, Uint32Field},
	)

	// create builder
	buffer := make([]byte, 2)
	builder := NewBuilder(TestType, buffer)

	// successful write
	wrote, err := builder.PutUint32("uint32", uint32(20))
	assert.Nil(t, err)
	assert.Equal(t, uint64(2), wrote)

	// test data validity
	assert.Equal(t, UnsignedInt8Code.OpCode, uint8(builder.buffer[0]))
	assert.Equal(t, 20, int(builder.buffer[1]))

	// validate field offset
	assert.Equal(t, 0, builder.offsets["uint32"])
}

func TestBuilderPutUint32Fail_2(t *testing.T) {

	// create test type
	// integer test type
	TestType := New("testing", "uint32")
	TestType.AddVersion(
		Field{"int32", true, Int32Field},
		Field{"uint32", true, Uint32Field},
	)

	// create builder
	buffer := make([]byte, 1)
	builder := NewBuilder(TestType, buffer)

	// fails length check
	wrote, err := builder.PutUint32("uint32", uint32(300))
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), wrote)
}

func TestBuilderPutUint32Pass_2(t *testing.T) {

	// create test type
	// integer test type
	TestType := New("testing", "uint32")
	TestType.AddVersion(
		Field{"int32", true, Int32Field},
		Field{"uint32", true, Uint32Field},
	)

	// create builder
	buffer := make([]byte, 5)
	builder := NewBuilder(TestType, buffer)

	// successful write
	wrote, err := builder.PutUint32("uint32", uint32(300))
	assert.Nil(t, err)
	assert.Equal(t, uint64(3), wrote)

	// test data validity
	assert.Equal(t, UnsignedInt16Code.OpCode, uint8(builder.buffer[0]))

	value, err := xbinary.LittleEndian.Uint32(buffer, 1)
	assert.Equal(t, uint16(300), value)

	// validate field offset
	assert.Equal(t, 0, builder.offsets["uint32"])
}

func TestBuilderPutUint32Fail_3(t *testing.T) {

	// create test type
	// integer test type
	TestType := New("testing", "uint32")
	TestType.AddVersion(
		Field{"int32", true, Int32Field},
		Field{"uint32", true, Uint32Field},
	)

	// create builder
	buffer := make([]byte, 3)
	builder := NewBuilder(TestType, buffer)

	// fails length check
	wrote, err := builder.PutUint32("uint32", uint32(135000))
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), wrote)
}

func TestBuilderPutUint32Pass_3(t *testing.T) {

	// create test type
	// integer test type
	TestType := New("testing", "uint32")
	TestType.AddVersion(
		Field{"int32", true, Int32Field},
		Field{"uint32", true, Uint32Field},
	)

	// create builder
	buffer := make([]byte, 5)
	builder := NewBuilder(TestType, buffer)

	// successful write
	wrote, err := builder.PutUint32("uint32", uint32(135000))
	assert.Nil(t, err)
	assert.Equal(t, uint64(5), wrote)

	// test data validity
	assert.Equal(t, UnsignedInt32Code.OpCode, uint8(builder.buffer[0]))

	value, err := xbinary.LittleEndian.Uint32(buffer, 1)
	assert.Equal(t, uint32(135000), value)

	// validate field offset
	assert.Equal(t, 0, builder.offsets["uint32"])
}

//
func TestBuilderPutInt32Fail_1(t *testing.T) {

	// create test type
	// integer test type
	TestType := New("testing", "uint32")
	TestType.AddVersion(
		Field{"int32", true, Int32Field},
		Field{"uint32", true, Uint32Field},
	)

	// create builder
	buffer := make([]byte, 1)
	builder := NewBuilder(TestType, buffer)

	// fails type check
	wrote, err := builder.PutInt32("uint32", int32(20))
	// fmt.Println(err)
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), wrote)

	// fails length check
	wrote, err = builder.PutInt32("int32", int32(20))
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), wrote)
}

func TestBuilderPutInt32Pass_1(t *testing.T) {

	// create test type
	// integer test type
	TestType := New("testing", "uint32")
	TestType.AddVersion(
		Field{"int32", true, Int32Field},
		Field{"uint32", true, Uint32Field},
	)

	// create builder
	buffer := make([]byte, 2)
	builder := NewBuilder(TestType, buffer)

	// successful write
	wrote, err := builder.PutInt32("int32", int32(20))
	assert.Nil(t, err)
	assert.Equal(t, uint64(2), wrote)

	// test data validity
	assert.Equal(t, Int8Code.OpCode, uint8(builder.buffer[0]))
	assert.Equal(t, 20, int(builder.buffer[1]))

	// validate field offset
	assert.Equal(t, 0, builder.offsets["uint32"])
}

func TestBuilderPutInt32Fail_2(t *testing.T) {

	// create test type
	// integer test type
	TestType := New("testing", "uint32")
	TestType.AddVersion(
		Field{"int32", true, Int32Field},
		Field{"uint32", true, Uint32Field},
	)

	// create builder
	buffer := make([]byte, 1)
	builder := NewBuilder(TestType, buffer)

	// fails length check
	wrote, err := builder.PutInt32("int32", int32(300))
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), wrote)
}

func TestBuilderPutInt32Pass_2(t *testing.T) {

	// create test type
	// integer test type
	TestType := New("testing", "uint32")
	TestType.AddVersion(
		Field{"int32", true, Int32Field},
		Field{"uint32", true, Uint32Field},
	)

	// create builder
	buffer := make([]byte, 5)
	builder := NewBuilder(TestType, buffer)

	// successful write
	wrote, err := builder.PutInt32("int32", int32(300))
	assert.Nil(t, err)
	assert.Equal(t, uint64(3), wrote)

	// test data validity
	assert.Equal(t, Int16Code.OpCode, uint8(builder.buffer[0]))

	value, err := xbinary.LittleEndian.Int16(buffer, 1)
	assert.Equal(t, int16(300), value)

	// validate field offset
	assert.Equal(t, 0, builder.offsets["uint32"])
}

func TestBuilderPutInt32Fail_3(t *testing.T) {

	// create test type
	// integer test type
	TestType := New("testing", "uint32")
	TestType.AddVersion(
		Field{"int32", true, Int32Field},
		Field{"uint32", true, Uint32Field},
	)

	// create builder
	buffer := make([]byte, 3)
	builder := NewBuilder(TestType, buffer)

	// fails length check
	wrote, err := builder.PutInt32("int32", int32(135000))
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), wrote)
}

func TestBuilderPutInt32Pass_3(t *testing.T) {

	// create test type
	// integer test type
	TestType := New("testing", "uint32")
	TestType.AddVersion(
		Field{"int32", true, Int32Field},
		Field{"uint32", true, Uint32Field},
	)

	// create builder
	buffer := make([]byte, 5)
	builder := NewBuilder(TestType, buffer)

	// successful write
	wrote, err := builder.PutInt32("int32", int32(135000))
	assert.Nil(t, err)
	assert.Equal(t, uint64(5), wrote)

	// test data validity
	assert.Equal(t, Int32Code.OpCode, uint8(builder.buffer[0]))

	value, err := xbinary.LittleEndian.Int32(buffer, 1)
	assert.Equal(t, int32(135000), value)

	// validate field offset
	assert.Equal(t, 0, builder.offsets["uint32"])
}

//
func TestBuilderPutUint64Fail_1(t *testing.T) {

	// create test type
	// integer test type
	TestType := New("testing", "uint64")
	TestType.AddVersion(
		Field{"int64", true, Int64Field},
		Field{"uint64", true, Uint64Field},
	)

	// create builder
	buffer := make([]byte, 1)
	builder := NewBuilder(TestType, buffer)

	// fails type check
	wrote, err := builder.PutUint64("int64", uint64(20))
	// fmt.Println(err)
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), wrote)

	// fails length check
	wrote, err = builder.PutUint64("uint64", uint64(20))
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), wrote)
}

func TestBuilderPutUint64Pass_1(t *testing.T) {

	// create test type
	// integer test type
	TestType := New("testing", "uint64")
	TestType.AddVersion(
		Field{"int64", true, Int64Field},
		Field{"uint64", true, Uint64Field},
	)

	// create builder
	buffer := make([]byte, 2)
	builder := NewBuilder(TestType, buffer)

	// successful write
	wrote, err := builder.PutUint64("uint64", uint64(20))
	assert.Nil(t, err)
	assert.Equal(t, uint64(2), wrote)

	// test data validity
	assert.Equal(t, UnsignedLong8Code.OpCode, uint8(builder.buffer[0]))
	assert.Equal(t, 20, int(builder.buffer[1]))

	// validate field offset
	assert.Equal(t, 0, builder.offsets["uint64"])
}

func TestBuilderPutUint64Fail_2(t *testing.T) {

	// create test type
	// integer test type
	TestType := New("testing", "uint64")
	TestType.AddVersion(
		Field{"int64", true, Int64Field},
		Field{"uint64", true, Uint64Field},
	)

	// create builder
	buffer := make([]byte, 1)
	builder := NewBuilder(TestType, buffer)

	// fails length check
	wrote, err := builder.PutUint64("uint64", uint64(300))
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), wrote)
}

func TestBuilderPutUint64Pass_2(t *testing.T) {

	// create test type
	// integer test type
	TestType := New("testing", "uint64")
	TestType.AddVersion(
		Field{"int64", true, Int64Field},
		Field{"uint64", true, Uint64Field},
	)

	// create builder
	buffer := make([]byte, 5)
	builder := NewBuilder(TestType, buffer)

	// successful write
	wrote, err := builder.PutUint64("uint64", uint64(300))
	assert.Nil(t, err)
	assert.Equal(t, uint64(3), wrote)

	// test data validity
	assert.Equal(t, UnsignedLong16Code.OpCode, uint8(builder.buffer[0]))

	value, err := xbinary.LittleEndian.Uint16(buffer, 1)
	assert.Equal(t, uint16(300), value)

	// validate field offset
	assert.Equal(t, 0, builder.offsets["uint64"])
}

func TestBuilderPutUint64Fail_3(t *testing.T) {

	// create test type
	// integer test type
	TestType := New("testing", "uint64")
	TestType.AddVersion(
		Field{"int64", true, Int64Field},
		Field{"uint64", true, Uint64Field},
	)

	// create builder
	buffer := make([]byte, 3)
	builder := NewBuilder(TestType, buffer)

	// fails length check
	wrote, err := builder.PutUint64("uint64", uint64(135000))
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), wrote)
}

func TestBuilderPutUint64Pass_3(t *testing.T) {

	// create test type
	// integer test type
	TestType := New("testing", "uint64")
	TestType.AddVersion(
		Field{"int64", true, Int64Field},
		Field{"uint64", true, Uint64Field},
	)

	// create builder
	buffer := make([]byte, 5)
	builder := NewBuilder(TestType, buffer)

	// successful write
	wrote, err := builder.PutUint64("uint64", uint64(135000))
	assert.Nil(t, err)
	assert.Equal(t, uint64(5), wrote)

	// test data validity
	assert.Equal(t, UnsignedLong32Code.OpCode, uint8(builder.buffer[0]))

	value, err := xbinary.LittleEndian.Uint32(buffer, 1)
	assert.Equal(t, uint64(135000), value)

	// validate field offset
	assert.Equal(t, 0, builder.offsets["uint64"])
}

func TestBuilderPutUint64Fail_4(t *testing.T) {

	// create test type
	// integer test type
	TestType := New("testing", "uint64")
	TestType.AddVersion(
		Field{"int64", true, Int64Field},
		Field{"uint64", true, Uint64Field},
	)

	// create builder
	buffer := make([]byte, 5)
	builder := NewBuilder(TestType, buffer)

	// fails length check
	wrote, err := builder.PutUint64("uint64", uint64(17179869184)) // 2^34
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), wrote)
}

func TestBuilderPutUint64Pass_4(t *testing.T) {

	// create test type
	// integer test type
	TestType := New("testing", "uint64")
	TestType.AddVersion(
		Field{"int64", true, Int64Field},
		Field{"uint64", true, Uint64Field},
	)

	// create builder
	buffer := make([]byte, 9)
	builder := NewBuilder(TestType, buffer)

	// successful write
	wrote, err := builder.PutUint64("uint64", uint64(17179869184)) // 2^34
	assert.Nil(t, err)
	assert.Equal(t, uint64(9), wrote)

	// test data validity
	assert.Equal(t, UnsignedLong64Code.OpCode, uint8(builder.buffer[0]))

	value, err := xbinary.LittleEndian.Uint64(buffer, 1)
	assert.Equal(t, uint64(17179869184), value)

	// validate field offset
	assert.Equal(t, 0, builder.offsets["uint64"])
}

//
func TestBuilderPutInt64Fail_1(t *testing.T) {

	// create test type
	// integer test type
	TestType := New("testing", "uint64")
	TestType.AddVersion(
		Field{"int64", true, Int64Field},
		Field{"uint64", true, Uint64Field},
	)

	// create builder
	buffer := make([]byte, 1)
	builder := NewBuilder(TestType, buffer)

	// fails type check
	wrote, err := builder.PutInt64("uint64", int64(20))
	// fmt.Println(err)
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), wrote)

	// fails length check
	wrote, err = builder.PutInt64("int64", int64(20))
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), wrote)
}

func TestBuilderPutInt64Pass_1(t *testing.T) {

	// create test type
	// integer test type
	TestType := New("testing", "uint64")
	TestType.AddVersion(
		Field{"int64", true, Int64Field},
		Field{"uint64", true, Uint64Field},
	)

	// create builder
	buffer := make([]byte, 2)
	builder := NewBuilder(TestType, buffer)

	// successful write
	wrote, err := builder.PutInt64("int64", int64(20))
	assert.Nil(t, err)
	assert.Equal(t, uint64(2), wrote)

	// test data validity
	assert.Equal(t, Long8Code.OpCode, uint8(builder.buffer[0]))
	assert.Equal(t, 20, int(builder.buffer[1]))

	// validate field offset
	assert.Equal(t, 0, builder.offsets["uint64"])
}

func TestBuilderPutInt64Fail_2(t *testing.T) {

	// create test type
	// integer test type
	TestType := New("testing", "uint64")
	TestType.AddVersion(
		Field{"int64", true, Int64Field},
		Field{"uint64", true, Uint64Field},
	)

	// create builder
	buffer := make([]byte, 1)
	builder := NewBuilder(TestType, buffer)

	// fails length check
	wrote, err := builder.PutInt64("int64", int64(300))
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), wrote)
}

func TestBuilderPutInt64Pass_2(t *testing.T) {

	// create test type
	// integer test type
	TestType := New("testing", "uint64")
	TestType.AddVersion(
		Field{"int64", true, Int64Field},
		Field{"uint64", true, Uint64Field},
	)

	// create builder
	buffer := make([]byte, 5)
	builder := NewBuilder(TestType, buffer)

	// successful write
	wrote, err := builder.PutInt64("int64", int64(300))
	assert.Nil(t, err)
	assert.Equal(t, uint64(3), wrote)

	// test data validity
	assert.Equal(t, Long16Code.OpCode, uint8(builder.buffer[0]))

	value, err := xbinary.LittleEndian.Int16(buffer, 1)
	assert.Equal(t, int16(300), value)

	// validate field offset
	assert.Equal(t, 0, builder.offsets["uint64"])
}

func TestBuilderPutInt64Fail_3(t *testing.T) {

	// create test type
	// integer test type
	TestType := New("testing", "uint64")
	TestType.AddVersion(
		Field{"int64", true, Int64Field},
		Field{"uint64", true, Uint64Field},
	)

	// create builder
	buffer := make([]byte, 3)
	builder := NewBuilder(TestType, buffer)

	// fails length check
	wrote, err := builder.PutInt64("int64", int64(135000))
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), wrote)
}

func TestBuilderPutInt64Pass_3(t *testing.T) {

	// create test type
	// integer test type
	TestType := New("testing", "uint64")
	TestType.AddVersion(
		Field{"int64", true, Int64Field},
		Field{"uint64", true, Uint64Field},
	)

	// create builder
	buffer := make([]byte, 5)
	builder := NewBuilder(TestType, buffer)

	// successful write
	wrote, err := builder.PutInt64("int64", int64(135000))
	assert.Nil(t, err)
	assert.Equal(t, uint64(5), wrote)

	// test data validity
	assert.Equal(t, Long32Code.OpCode, uint8(builder.buffer[0]))

	value, err := xbinary.LittleEndian.Int32(buffer, 1)
	assert.Equal(t, int64(135000), value)

	// validate field offset
	assert.Equal(t, 0, builder.offsets["uint64"])
}

func TestBuilderPutInt64Fail_4(t *testing.T) {

	// create test type
	// integer test type
	TestType := New("testing", "uint64")
	TestType.AddVersion(
		Field{"int64", true, Int64Field},
		Field{"uint64", true, Uint64Field},
	)

	// create builder
	buffer := make([]byte, 5)
	builder := NewBuilder(TestType, buffer)

	// fails length check
	wrote, err := builder.PutInt64("int64", int64(17179869184)) // 2^34
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), wrote)
}

func TestBuilderPutInt64Pass_4(t *testing.T) {

	// create test type
	// integer test type
	TestType := New("testing", "uint64")
	TestType.AddVersion(
		Field{"int64", true, Int64Field},
		Field{"uint64", true, Uint64Field},
	)

	// create builder
	buffer := make([]byte, 9)
	builder := NewBuilder(TestType, buffer)

	// successful write
	wrote, err := builder.PutInt64("int64", int64(17179869184)) // 2^34
	assert.Nil(t, err)
	assert.Equal(t, uint64(9), wrote)

	// test data validity
	assert.Equal(t, Long64Code.OpCode, uint8(builder.buffer[0]))

	value, err := xbinary.LittleEndian.Int64(buffer, 1)
	assert.Equal(t, int64(17179869184), value)

	// validate field offset
	assert.Equal(t, 0, builder.offsets["uint64"])
}
