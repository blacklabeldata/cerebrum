package xbinary

import (
	"encoding/binary"
	// "fmt"
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestStringBigEndian(t *testing.T) {
	buf := make([]byte, 6)
	expected := []byte{0x67, 0x6f, 0x6c, 0x61, 0x6e, 0x67}

	// put string
	bytes, err := BigEndian.PutString(buf, 0, "golang")
	assert.Equal(t, expected, buf)

	// put string
	bytes, err = BigEndian.PutString(buf, -1, "golang")
	assert.Equal(t, uint64(0), bytes)
	assert.NotNil(t, err)
	assert.Equal(t, err, ErrOutOfRange)

	// put string
	bytes, err = BigEndian.PutString(buf, 4, "golang")
	assert.Equal(t, uint64(0), bytes)
	assert.NotNil(t, err)
	assert.Equal(t, err, ErrOutOfRange)

	// get string
	name, err := BigEndian.String(buf, 0, 6)
	assert.Equal(t, "golang", name)
	assert.Nil(t, err)

	// get string error
	name, err = BigEndian.String(buf, 4, 6)
	assert.Equal(t, "", name)
	assert.NotNil(t, name)

	// get string error
	name, err = BigEndian.String(buf, -1, 6)
	assert.Equal(t, "", name)
	assert.Equal(t, err, ErrOutOfRange)
	assert.NotNil(t, err)
}

func TestUint16BigEndian(t *testing.T) {
	var buf = make([]byte, 2)
	var expected uint16 = 257

	// Put
	bytes, err := BigEndian.PutUint16(buf, 0, expected)
	assert.Nil(t, err)
	assert.Equal(t, expected, binary.BigEndian.Uint16(buf))
	assert.Equal(t, uint64(2), bytes)

	// Put error (index = -1)
	bytes, err = BigEndian.PutUint16(buf, -1, expected)
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), bytes)

	// Put error
	bytes, err = BigEndian.PutUint16(buf, 2, expected)
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), bytes)

	// Get
	val, err := BigEndian.Uint16(buf, 0)
	assert.Nil(t, err)
	assert.Equal(t, expected, val)

	// Get Error (index = -1)
	val, err = BigEndian.Uint16(buf, -1)
	assert.NotNil(t, err)
	assert.Equal(t, uint16(0), val)

	// Get Error (index = 2)
	val, err = BigEndian.Uint16(buf, 2)
	assert.NotNil(t, err)
	assert.Equal(t, uint16(0), val)
}

func TestUint32BigEndian(t *testing.T) {
	var buf = make([]byte, 4)
	var expected uint32 = 2 << 15

	// Put
	bytes, err := BigEndian.PutUint32(buf, 0, expected)
	assert.Nil(t, err)
	assert.Equal(t, expected, binary.BigEndian.Uint32(buf))
	assert.Equal(t, uint64(4), bytes)

	// Put error (index = -1)
	bytes, err = BigEndian.PutUint32(buf, -1, expected)
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), bytes)

	// Put error (index = 2)
	bytes, err = BigEndian.PutUint32(buf, 2, expected)
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), bytes)

	// Get
	val, err := BigEndian.Uint32(buf, 0)
	assert.Nil(t, err)
	assert.Equal(t, expected, val)

	// Get Error (index = -1)
	val, err = BigEndian.Uint32(buf, -1)
	assert.NotNil(t, err)
	assert.Equal(t, uint32(0), val)

	// Get Error (index = 2)
	val, err = BigEndian.Uint32(buf, 2)
	assert.NotNil(t, err)
	assert.Equal(t, uint32(0), val)
}

func TestUint64BigEndian(t *testing.T) {
	var buf = make([]byte, 8)
	var expected uint64 = 2 << 31

	// Put
	bytes, err := BigEndian.PutUint64(buf, 0, expected)
	assert.Nil(t, err)
	assert.Equal(t, expected, binary.BigEndian.Uint64(buf))
	assert.Equal(t, uint64(8), bytes)

	// Put error (index = -1)
	bytes, err = BigEndian.PutUint64(buf, -1, expected)
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), bytes)

	// Put error (index = 2)
	bytes, err = BigEndian.PutUint64(buf, 2, expected)
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), bytes)

	// Get
	val, err := BigEndian.Uint64(buf, 0)
	assert.Nil(t, err)
	assert.Equal(t, expected, val)

	// Get Error (index = -1)
	val, err = BigEndian.Uint64(buf, -1)
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), val)

	// Get Error (index = 2)
	val, err = BigEndian.Uint64(buf, 2)
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), val)
}

func TestInt16BigEndian(t *testing.T) {
	var buf = make([]byte, 2)
	var expected int16 = -257

	// Put
	bytes, err := BigEndian.PutInt16(buf, 0, expected)
	assert.Nil(t, err)
	assert.Equal(t, expected, int16(binary.BigEndian.Uint16(buf)))
	assert.Equal(t, uint64(2), bytes)

	// Put error (index = -1)
	bytes, err = BigEndian.PutInt16(buf, -1, expected)
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), bytes)

	// Put error
	bytes, err = BigEndian.PutInt16(buf, 2, expected)
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), bytes)

	// Get
	val, err := BigEndian.Int16(buf, 0)
	assert.Nil(t, err)
	assert.Equal(t, expected, val)

	// Get Error (index = -1)
	val, err = BigEndian.Int16(buf, -1)
	assert.NotNil(t, err)
	assert.Equal(t, int16(0), val)

	// Get Error (index = 2)
	val, err = BigEndian.Int16(buf, 2)
	assert.NotNil(t, err)
	assert.Equal(t, int16(0), val)
}

func TestInt32BigEndian(t *testing.T) {
	var buf = make([]byte, 4)
	var expected int32 = -(2 << 15)

	// Put
	bytes, err := BigEndian.PutInt32(buf, 0, expected)
	assert.Nil(t, err)
	assert.Equal(t, expected, int32(binary.BigEndian.Uint32(buf)))
	assert.Equal(t, uint64(4), bytes)

	// Put error (index = -1)
	bytes, err = BigEndian.PutInt32(buf, -1, expected)
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), bytes)

	// Put error
	bytes, err = BigEndian.PutInt32(buf, 2, expected)
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), bytes)

	// Get
	val, err := BigEndian.Int32(buf, 0)
	assert.Nil(t, err)
	assert.Equal(t, expected, val)

	// Get Error (index = -1)
	val, err = BigEndian.Int32(buf, -1)
	assert.NotNil(t, err)
	assert.Equal(t, int32(0), val)

	// Get Error (index = 2)
	val, err = BigEndian.Int32(buf, 2)
	assert.NotNil(t, err)
	assert.Equal(t, int32(0), val)
}

func TestInt64BigEndian(t *testing.T) {
	var buf = make([]byte, 8)
	var expected int64 = -(2 << 31)

	// Put
	bytes, err := BigEndian.PutInt64(buf, 0, expected)
	assert.Nil(t, err)
	assert.Equal(t, expected, int64(binary.BigEndian.Uint64(buf)))
	assert.Equal(t, uint64(8), bytes)

	// Put error (index = -1)
	bytes, err = BigEndian.PutInt64(buf, -1, expected)
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), bytes)

	// Put error
	bytes, err = BigEndian.PutInt64(buf, 2, expected)
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), bytes)

	// Get
	val, err := BigEndian.Int64(buf, 0)
	assert.Nil(t, err)
	assert.Equal(t, expected, val)

	// Get Error (index = -1)
	val, err = BigEndian.Int64(buf, -1)
	assert.NotNil(t, err)
	assert.Equal(t, int64(0), val)

	// Get Error (index = 2)
	val, err = BigEndian.Int64(buf, 2)
	assert.NotNil(t, err)
	assert.Equal(t, int64(0), val)
}

func TestFloat32BigEndian(t *testing.T) {
	var buf = make([]byte, 4)
	var expected = float32(2 << 15)

	// Put
	bytes, err := BigEndian.PutFloat32(buf, 0, expected)
	assert.Nil(t, err)
	assert.Equal(t, expected, math.Float32frombits(binary.BigEndian.Uint32(buf)))
	assert.Equal(t, uint64(4), bytes)

	// Put error (index = -1)
	bytes, err = BigEndian.PutFloat32(buf, -1, expected)
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), bytes)

	// Put error (index = 2)
	bytes, err = BigEndian.PutFloat32(buf, 2, expected)
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), bytes)

	// Get
	val, err := BigEndian.Float32(buf, 0)
	assert.Nil(t, err)
	assert.Equal(t, expected, val)

	// Get Error (index = -1)
	val, err = BigEndian.Float32(buf, -1)
	assert.NotNil(t, err)
	assert.Equal(t, float32(0), val)

	// Get Error (index = 2)
	val, err = BigEndian.Float32(buf, 2)
	assert.NotNil(t, err)
	assert.Equal(t, float32(0), val)
}

func TestFloat64BigEndian(t *testing.T) {
	var buf = make([]byte, 8)
	var expected = float64(2 << 31)

	// Put
	bytes, err := BigEndian.PutFloat64(buf, 0, expected)
	assert.Nil(t, err)
	assert.Equal(t, expected, math.Float64frombits(binary.BigEndian.Uint64(buf)))
	assert.Equal(t, uint64(8), bytes)

	// Put error (index = -1)
	bytes, err = BigEndian.PutFloat64(buf, -1, expected)
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), bytes)

	// Put error (index = 2)
	bytes, err = BigEndian.PutFloat64(buf, 2, expected)
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), bytes)

	// Get
	val, err := BigEndian.Float64(buf, 0)
	assert.Nil(t, err)
	assert.Equal(t, expected, val)

	// Get Error (index = -1)
	val, err = BigEndian.Float64(buf, -1)
	assert.NotNil(t, err)
	assert.Equal(t, float64(0), val)

	// Get Error (index = 2)
	val, err = BigEndian.Float64(buf, 2)
	assert.NotNil(t, err)
	assert.Equal(t, float64(0), val)
}

func TestUint16ArrayBigEndian(t *testing.T) {
	var buf = make([]byte, 6)
	var expected = []uint16{0, 1, 2}

	// Put
	bytes, err := BigEndian.PutUint16Array(buf, 0, expected)
	assert.Nil(t, err)
	assert.Equal(t, uint64(6), bytes)
	assert.Equal(t, expected[0], binary.BigEndian.Uint16(buf))
	assert.Equal(t, expected[1], binary.BigEndian.Uint16(buf[2:]))
	assert.Equal(t, expected[2], binary.BigEndian.Uint16(buf[4:]))

	// Put error (index = -1)
	bytes, err = BigEndian.PutUint16Array(buf, -1, expected)
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), bytes)

	// Put error
	bytes, err = BigEndian.PutUint16Array(buf, 2, expected)
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), bytes)

	// Get
	dest := make([]uint16, 3)
	err = BigEndian.Uint16Array(buf, 0, &dest)
	assert.Nil(t, err)
	assert.Equal(t, expected, dest)

	// Get Error (index = -1)
	dest = make([]uint16, 3)
	err = BigEndian.Uint16Array(buf, -1, &dest)
	assert.NotNil(t, err)
	assert.Equal(t, dest[0], uint16(0))
	assert.Equal(t, dest[1], uint16(0))
	assert.Equal(t, dest[2], uint16(0))

	// Get Error (index = 2)
	err = BigEndian.Uint16Array(buf, 2, &dest)
	assert.NotNil(t, err)
	assert.Equal(t, dest[0], uint16(0))
	assert.Equal(t, dest[1], uint16(0))
	assert.Equal(t, dest[2], uint16(0))
}

func TestInt16ArrayBigEndian(t *testing.T) {
	var buf = make([]byte, 6)
	var expected = []int16{0, -1, 2}

	// Put
	bytes, err := BigEndian.PutInt16Array(buf, 0, expected)
	assert.Nil(t, err)
	assert.Equal(t, uint64(6), bytes)
	assert.Equal(t, expected[0], int16(binary.BigEndian.Uint16(buf)))
	assert.Equal(t, expected[1], int16(binary.BigEndian.Uint16(buf[2:])))
	assert.Equal(t, expected[2], int16(binary.BigEndian.Uint16(buf[4:])))

	// Put error (index = -1)
	bytes, err = BigEndian.PutInt16Array(buf, -1, expected)
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), bytes)

	// Put error
	bytes, err = BigEndian.PutInt16Array(buf, 2, expected)
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), bytes)

	// Get
	dest := make([]int16, 3)
	err = BigEndian.Int16Array(buf, 0, &dest)
	assert.Nil(t, err)
	assert.Equal(t, expected, dest)

	// Get Error (index = -1)
	dest = make([]int16, 3)
	err = BigEndian.Int16Array(buf, -1, &dest)
	assert.NotNil(t, err)
	assert.Equal(t, dest[0], int16(0))
	assert.Equal(t, dest[1], int16(0))
	assert.Equal(t, dest[2], int16(0))

	// Get Error (index = 2)
	err = BigEndian.Int16Array(buf, 2, &dest)
	assert.NotNil(t, err)
	assert.Equal(t, dest[0], int16(0))
	assert.Equal(t, dest[1], int16(0))
	assert.Equal(t, dest[2], int16(0))
}

func TestUint32ArrayBigEndian(t *testing.T) {
	var buf = make([]byte, 12)
	var expected = []uint32{0, 1, 2}

	// Put
	bytes, err := BigEndian.PutUint32Array(buf, 0, expected)
	assert.Nil(t, err)
	assert.Equal(t, uint64(12), bytes)
	assert.Equal(t, expected[0], binary.BigEndian.Uint32(buf))
	assert.Equal(t, expected[1], binary.BigEndian.Uint32(buf[4:]))
	assert.Equal(t, expected[2], binary.BigEndian.Uint32(buf[8:]))

	// Put error (index = -1)
	bytes, err = BigEndian.PutUint32Array(buf, -1, expected)
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), bytes)

	// Put error
	bytes, err = BigEndian.PutUint32Array(buf, 2, expected)
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), bytes)

	// Get
	dest := make([]uint32, 3)
	err = BigEndian.Uint32Array(buf, 0, &dest)
	assert.Nil(t, err)
	assert.Equal(t, expected, dest)

	// Get Error (index = -1)
	dest = make([]uint32, 3)
	err = BigEndian.Uint32Array(buf, -1, &dest)
	assert.NotNil(t, err)
	assert.Equal(t, dest[0], uint32(0))
	assert.Equal(t, dest[1], uint32(0))
	assert.Equal(t, dest[2], uint32(0))

	// Get Error (index = 2)
	err = BigEndian.Uint32Array(buf, 2, &dest)
	assert.NotNil(t, err)
	assert.Equal(t, dest[0], uint32(0))
	assert.Equal(t, dest[1], uint32(0))
	assert.Equal(t, dest[2], uint32(0))
}

func TestInt32ArrayBigEndian(t *testing.T) {
	var buf = make([]byte, 12)
	var expected = []int32{0, -1, 2}

	// Put
	bytes, err := BigEndian.PutInt32Array(buf, 0, expected)
	assert.Nil(t, err)
	assert.Equal(t, uint64(12), bytes)
	assert.Equal(t, expected[0], int32(binary.BigEndian.Uint32(buf)))
	assert.Equal(t, expected[1], int32(binary.BigEndian.Uint32(buf[4:])))
	assert.Equal(t, expected[2], int32(binary.BigEndian.Uint32(buf[8:])))

	// Put error (index = -1)
	bytes, err = BigEndian.PutInt32Array(buf, -1, expected)
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), bytes)

	// Put error
	bytes, err = BigEndian.PutInt32Array(buf, 2, expected)
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), bytes)

	// Get
	dest := make([]int32, 3)
	err = BigEndian.Int32Array(buf, 0, &dest)
	assert.Nil(t, err)
	assert.Equal(t, expected, dest)

	// Get Error (index = -1)
	dest = make([]int32, 3)
	err = BigEndian.Int32Array(buf, -1, &dest)
	assert.NotNil(t, err)
	assert.Equal(t, dest[0], int32(0))
	assert.Equal(t, dest[1], int32(0))
	assert.Equal(t, dest[2], int32(0))

	// Get Error (index = 2)
	err = BigEndian.Int32Array(buf, 2, &dest)
	assert.NotNil(t, err)
	assert.Equal(t, dest[0], int32(0))
	assert.Equal(t, dest[1], int32(0))
	assert.Equal(t, dest[2], int32(0))
}

func TestUint64ArrayBigEndian(t *testing.T) {
	var buf = make([]byte, 24)
	var expected = []uint64{0, 1, 2}

	// Put
	bytes, err := BigEndian.PutUint64Array(buf, 0, expected)
	assert.Nil(t, err)
	assert.Equal(t, uint64(24), bytes)
	assert.Equal(t, expected[0], binary.BigEndian.Uint64(buf))
	assert.Equal(t, expected[1], binary.BigEndian.Uint64(buf[8:]))
	assert.Equal(t, expected[2], binary.BigEndian.Uint64(buf[16:]))
	// fmt.Println(expected)
	// fmt.Println(buf)

	// Put error (index = -1)
	bytes, err = BigEndian.PutUint64Array(buf, -1, expected)
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), bytes)

	// Put error
	bytes, err = BigEndian.PutUint64Array(buf, 2, expected)
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), bytes)

	// Get
	dest := make([]uint64, 3)
	err = BigEndian.Uint64Array(buf, 0, &dest)
	assert.Nil(t, err)
	// fmt.Println(expected)
	// fmt.Println(dest)
	assert.Equal(t, expected, dest)

	// Get Error (index = -1)
	dest = make([]uint64, 3)
	err = BigEndian.Uint64Array(buf, -1, &dest)
	assert.NotNil(t, err)
	assert.Equal(t, dest[0], uint64(0))
	assert.Equal(t, dest[1], uint64(0))
	assert.Equal(t, dest[2], uint64(0))

	// Get Error (index = 2)
	err = BigEndian.Uint64Array(buf, 2, &dest)
	assert.NotNil(t, err)
	assert.Equal(t, dest[0], uint64(0))
	assert.Equal(t, dest[1], uint64(0))
	assert.Equal(t, dest[2], uint64(0))
}

func TestInt64ArrayBigEndian(t *testing.T) {
	var buf = make([]byte, 24)
	var expected = []int64{0, -1, 2}

	// Put
	bytes, err := BigEndian.PutInt64Array(buf, 0, expected)
	assert.Nil(t, err)
	assert.Equal(t, uint64(24), bytes)
	assert.Equal(t, expected[0], int64(binary.BigEndian.Uint64(buf)))
	assert.Equal(t, expected[1], int64(binary.BigEndian.Uint64(buf[8:])))
	assert.Equal(t, expected[2], int64(binary.BigEndian.Uint64(buf[16:])))

	// Put error (index = -1)
	bytes, err = BigEndian.PutInt64Array(buf, -1, expected)
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), bytes)

	// Put error
	bytes, err = BigEndian.PutInt64Array(buf, 2, expected)
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), bytes)

	// Get
	dest := make([]int64, 3)
	err = BigEndian.Int64Array(buf, 0, &dest)
	assert.Nil(t, err)
	assert.Equal(t, expected, dest)

	// Get Error (index = -1)
	dest = make([]int64, 3)
	err = BigEndian.Int64Array(buf, -1, &dest)
	assert.NotNil(t, err)
	assert.Equal(t, dest[0], int64(0))
	assert.Equal(t, dest[1], int64(0))
	assert.Equal(t, dest[2], int64(0))

	// Get Error (index = 2)
	err = BigEndian.Int64Array(buf, 2, &dest)
	assert.NotNil(t, err)
	assert.Equal(t, dest[0], int64(0))
	assert.Equal(t, dest[1], int64(0))
	assert.Equal(t, dest[2], int64(0))
}

func TestFloat32ArrayBigEndian(t *testing.T) {
	var buf = make([]byte, 12)
	var expected = []float32{0.0, 1.0, 2.0}

	// Put
	bytes, err := BigEndian.PutFloat32Array(buf, 0, expected)
	assert.Nil(t, err)
	assert.Equal(t, uint64(12), bytes)
	assert.Equal(t, expected[0], math.Float32frombits(binary.BigEndian.Uint32(buf)))
	assert.Equal(t, expected[1], math.Float32frombits(binary.BigEndian.Uint32(buf[4:])))
	assert.Equal(t, expected[2], math.Float32frombits(binary.BigEndian.Uint32(buf[8:])))

	// Put error (index = -1)
	bytes, err = BigEndian.PutFloat32Array(buf, -1, expected)
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), bytes)

	// Put error
	bytes, err = BigEndian.PutFloat32Array(buf, 2, expected)
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), bytes)

	// Get
	dest := make([]float32, 3)
	err = BigEndian.Float32Array(buf, 0, &dest)
	assert.Nil(t, err)
	assert.Equal(t, expected, dest)

	// Get Error (index = -1)
	dest = make([]float32, 3)
	err = BigEndian.Float32Array(buf, -1, &dest)
	assert.NotNil(t, err)
	assert.Equal(t, dest[0], float32(0.0))
	assert.Equal(t, dest[1], float32(0.0))
	assert.Equal(t, dest[2], float32(0.0))

	// Get Error (index = 2)
	err = BigEndian.Float32Array(buf, 2, &dest)
	assert.NotNil(t, err)
	assert.Equal(t, dest[0], float32(0.0))
	assert.Equal(t, dest[1], float32(0.0))
	assert.Equal(t, dest[2], float32(0.0))
}

func TestFloat64ArrayBigEndian(t *testing.T) {
	var buf = make([]byte, 24)
	var expected = []float64{0.0, 1.0, 2.0}

	// Put
	bytes, err := BigEndian.PutFloat64Array(buf, 0, expected)
	assert.Nil(t, err)
	assert.Equal(t, uint64(24), bytes)
	assert.Equal(t, expected[0], math.Float64frombits(binary.BigEndian.Uint64(buf)))
	assert.Equal(t, expected[1], math.Float64frombits(binary.BigEndian.Uint64(buf[8:])))
	assert.Equal(t, expected[2], math.Float64frombits(binary.BigEndian.Uint64(buf[16:])))

	// Put error (index = -1)
	bytes, err = BigEndian.PutFloat64Array(buf, -1, expected)
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), bytes)

	// Put error
	bytes, err = BigEndian.PutFloat64Array(buf, 2, expected)
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), bytes)

	// Get
	dest := make([]float64, 3)
	err = BigEndian.Float64Array(buf, 0, &dest)
	assert.Nil(t, err)
	assert.Equal(t, expected, dest)

	// Get Error (index = -1)
	dest = make([]float64, 3)
	err = BigEndian.Float64Array(buf, -1, &dest)
	assert.NotNil(t, err)
	assert.Equal(t, dest[0], float64(0.0))
	assert.Equal(t, dest[1], float64(0.0))
	assert.Equal(t, dest[2], float64(0.0))

	// Get Error (index = 2)
	err = BigEndian.Float64Array(buf, 2, &dest)
	assert.NotNil(t, err)
	assert.Equal(t, dest[0], float64(0.0))
	assert.Equal(t, dest[1], float64(0.0))
	assert.Equal(t, dest[2], float64(0.0))
}
