package xbinary

import (
	"encoding/binary"
	// "fmt"
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestString(t *testing.T) {
	buf := make([]byte, 6)
	expected := []byte{0x67, 0x6f, 0x6c, 0x61, 0x6e, 0x67}

	// put string
	bytes, err := LittleEndian.PutString(buf, 0, "golang")
	assert.Equal(t, expected, buf)

	// put string
	bytes, err = LittleEndian.PutString(buf, -1, "golang")
	assert.Equal(t, uint64(0), bytes)
	assert.NotNil(t, err)
	assert.Equal(t, err, ErrOutOfRange)

	// put string
	bytes, err = LittleEndian.PutString(buf, 4, "golang")
	assert.Equal(t, uint64(0), bytes)
	assert.NotNil(t, err)
	assert.Equal(t, err, ErrOutOfRange)

	// get string
	name, err := LittleEndian.String(buf, 0, 6)
	assert.Equal(t, "golang", name)
	assert.Nil(t, err)

	// get string error
	name, err = LittleEndian.String(buf, 4, 6)
	assert.Equal(t, "", name)
	assert.NotNil(t, name)

	// get string error
	name, err = LittleEndian.String(buf, -1, 6)
	assert.Equal(t, "", name)
	assert.Equal(t, err, ErrOutOfRange)
	assert.NotNil(t, err)

	// get string error
	name, err = LittleEndian.String(buf, 4, 6)
	assert.Equal(t, "", name)
	assert.Equal(t, err, ErrOutOfRange)
	assert.NotNil(t, err)
}

func TestUint16(t *testing.T) {
	var buf = make([]byte, 2)
	var expected uint16 = 257

	// Put
	bytes, err := LittleEndian.PutUint16(buf, 0, expected)
	assert.Nil(t, err)
	assert.Equal(t, expected, binary.LittleEndian.Uint16(buf))
	assert.Equal(t, uint64(2), bytes)

	// Put error (index = -1)
	bytes, err = LittleEndian.PutUint16(buf, -1, expected)
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), bytes)

	// Put error
	bytes, err = LittleEndian.PutUint16(buf, 2, expected)
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), bytes)

	// Get
	val, err := LittleEndian.Uint16(buf, 0)
	assert.Nil(t, err)
	assert.Equal(t, expected, val)

	// Get Error (index = -1)
	val, err = LittleEndian.Uint16(buf, -1)
	assert.NotNil(t, err)
	assert.Equal(t, uint16(0), val)

	// Get Error (index = 2)
	val, err = LittleEndian.Uint16(buf, 2)
	assert.NotNil(t, err)
	assert.Equal(t, uint16(0), val)
}

func TestUint32(t *testing.T) {
	var buf = make([]byte, 4)
	var expected uint32 = 2 << 15

	// Put
	bytes, err := LittleEndian.PutUint32(buf, 0, expected)
	assert.Nil(t, err)
	assert.Equal(t, expected, binary.LittleEndian.Uint32(buf))
	assert.Equal(t, uint64(4), bytes)

	// Put error (index = -1)
	bytes, err = LittleEndian.PutUint32(buf, -1, expected)
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), bytes)

	// Put error (index = 2)
	bytes, err = LittleEndian.PutUint32(buf, 2, expected)
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), bytes)

	// Get
	val, err := LittleEndian.Uint32(buf, 0)
	assert.Nil(t, err)
	assert.Equal(t, expected, val)

	// Get Error (index = -1)
	val, err = LittleEndian.Uint32(buf, -1)
	assert.NotNil(t, err)
	assert.Equal(t, uint32(0), val)

	// Get Error (index = 2)
	val, err = LittleEndian.Uint32(buf, 2)
	assert.NotNil(t, err)
	assert.Equal(t, uint32(0), val)
}

func TestUint64(t *testing.T) {
	var buf = make([]byte, 8)
	var expected uint64 = 2 << 31

	// Put
	bytes, err := LittleEndian.PutUint64(buf, 0, expected)
	assert.Nil(t, err)
	assert.Equal(t, expected, binary.LittleEndian.Uint64(buf))
	assert.Equal(t, uint64(8), bytes)

	// Put error (index = -1)
	bytes, err = LittleEndian.PutUint64(buf, -1, expected)
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), bytes)

	// Put error (index = 2)
	bytes, err = LittleEndian.PutUint64(buf, 2, expected)
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), bytes)

	// Get
	val, err := LittleEndian.Uint64(buf, 0)
	assert.Nil(t, err)
	assert.Equal(t, expected, val)

	// Get Error (index = -1)
	val, err = LittleEndian.Uint64(buf, -1)
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), val)

	// Get Error (index = 2)
	val, err = LittleEndian.Uint64(buf, 2)
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), val)
}

func TestInt16(t *testing.T) {
	var buf = make([]byte, 2)
	var expected int16 = -257

	// Put
	bytes, err := LittleEndian.PutInt16(buf, 0, expected)
	assert.Nil(t, err)
	assert.Equal(t, expected, int16(binary.LittleEndian.Uint16(buf)))
	assert.Equal(t, uint64(2), bytes)

	// Put error (index = -1)
	bytes, err = LittleEndian.PutInt16(buf, -1, expected)
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), bytes)

	// Put error
	bytes, err = LittleEndian.PutInt16(buf, 2, expected)
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), bytes)

	// Get
	val, err := LittleEndian.Int16(buf, 0)
	assert.Nil(t, err)
	assert.Equal(t, expected, val)

	// Get Error (index = -1)
	val, err = LittleEndian.Int16(buf, -1)
	assert.NotNil(t, err)
	assert.Equal(t, int16(0), val)

	// Get Error (index = 2)
	val, err = LittleEndian.Int16(buf, 2)
	assert.NotNil(t, err)
	assert.Equal(t, int16(0), val)
}

func TestInt32(t *testing.T) {
	var buf = make([]byte, 4)
	var expected int32 = -(2 << 15)

	// Put
	bytes, err := LittleEndian.PutInt32(buf, 0, expected)
	assert.Nil(t, err)
	assert.Equal(t, expected, int32(binary.LittleEndian.Uint32(buf)))
	assert.Equal(t, uint64(4), bytes)

	// Put error (index = -1)
	bytes, err = LittleEndian.PutInt32(buf, -1, expected)
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), bytes)

	// Put error
	bytes, err = LittleEndian.PutInt32(buf, 2, expected)
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), bytes)

	// Get
	val, err := LittleEndian.Int32(buf, 0)
	assert.Nil(t, err)
	assert.Equal(t, expected, val)

	// Get Error (index = -1)
	val, err = LittleEndian.Int32(buf, -1)
	assert.NotNil(t, err)
	assert.Equal(t, int32(0), val)

	// Get Error (index = 2)
	val, err = LittleEndian.Int32(buf, 2)
	assert.NotNil(t, err)
	assert.Equal(t, int32(0), val)
}

func TestInt64(t *testing.T) {
	var buf = make([]byte, 8)
	var expected int64 = -(2 << 31)

	// Put
	bytes, err := LittleEndian.PutInt64(buf, 0, expected)
	assert.Nil(t, err)
	assert.Equal(t, expected, int64(binary.LittleEndian.Uint64(buf)))
	assert.Equal(t, uint64(8), bytes)

	// Put error (index = -1)
	bytes, err = LittleEndian.PutInt64(buf, -1, expected)
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), bytes)

	// Put error
	bytes, err = LittleEndian.PutInt64(buf, 2, expected)
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), bytes)

	// Get
	val, err := LittleEndian.Int64(buf, 0)
	assert.Nil(t, err)
	assert.Equal(t, expected, val)

	// Get Error (index = -1)
	val, err = LittleEndian.Int64(buf, -1)
	assert.NotNil(t, err)
	assert.Equal(t, int64(0), val)

	// Get Error (index = 2)
	val, err = LittleEndian.Int64(buf, 2)
	assert.NotNil(t, err)
	assert.Equal(t, int64(0), val)
}

func TestFloat32(t *testing.T) {
	var buf = make([]byte, 4)
	var expected = float32(2 << 15)

	// Put
	bytes, err := LittleEndian.PutFloat32(buf, 0, expected)
	assert.Nil(t, err)
	assert.Equal(t, expected, math.Float32frombits(binary.LittleEndian.Uint32(buf)))
	assert.Equal(t, uint64(4), bytes)

	// Put error (index = -1)
	bytes, err = LittleEndian.PutFloat32(buf, -1, expected)
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), bytes)

	// Put error (index = 2)
	bytes, err = LittleEndian.PutFloat32(buf, 2, expected)
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), bytes)

	// Get
	val, err := LittleEndian.Float32(buf, 0)
	assert.Nil(t, err)
	assert.Equal(t, expected, val)

	// Get Error (index = -1)
	val, err = LittleEndian.Float32(buf, -1)
	assert.NotNil(t, err)
	assert.Equal(t, float32(0), val)

	// Get Error (index = 2)
	val, err = LittleEndian.Float32(buf, 2)
	assert.NotNil(t, err)
	assert.Equal(t, float32(0), val)
}

func TestFloat64(t *testing.T) {
	var buf = make([]byte, 8)
	var expected = float64(2 << 31)

	// Put
	bytes, err := LittleEndian.PutFloat64(buf, 0, expected)
	assert.Nil(t, err)
	assert.Equal(t, expected, math.Float64frombits(binary.LittleEndian.Uint64(buf)))
	assert.Equal(t, uint64(8), bytes)

	// Put error (index = -1)
	bytes, err = LittleEndian.PutFloat64(buf, -1, expected)
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), bytes)

	// Put error (index = 2)
	bytes, err = LittleEndian.PutFloat64(buf, 2, expected)
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), bytes)

	// Get
	val, err := LittleEndian.Float64(buf, 0)
	assert.Nil(t, err)
	assert.Equal(t, expected, val)

	// Get Error (index = -1)
	val, err = LittleEndian.Float64(buf, -1)
	assert.NotNil(t, err)
	assert.Equal(t, float64(0), val)

	// Get Error (index = 2)
	val, err = LittleEndian.Float64(buf, 2)
	assert.NotNil(t, err)
	assert.Equal(t, float64(0), val)
}

func TestUint16Array(t *testing.T) {
	var buf = make([]byte, 6)
	var expected = []uint16{0, 1, 2}

	// Put
	bytes, err := LittleEndian.PutUint16Array(buf, 0, expected)
	assert.Nil(t, err)
	assert.Equal(t, uint64(6), bytes)
	assert.Equal(t, expected[0], binary.LittleEndian.Uint16(buf))
	assert.Equal(t, expected[1], binary.LittleEndian.Uint16(buf[2:]))
	assert.Equal(t, expected[2], binary.LittleEndian.Uint16(buf[4:]))

	// Put error (index = -1)
	bytes, err = LittleEndian.PutUint16Array(buf, -1, expected)
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), bytes)

	// Put error
	bytes, err = LittleEndian.PutUint16Array(buf, 2, expected)
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), bytes)

	// Get
	dest := make([]uint16, 3)
	err = LittleEndian.Uint16Array(buf, 0, &dest)
	assert.Nil(t, err)
	assert.Equal(t, expected, dest)

	// Get Error (index = -1)
	dest = make([]uint16, 3)
	err = LittleEndian.Uint16Array(buf, -1, &dest)
	assert.NotNil(t, err)
	assert.Equal(t, dest[0], uint16(0))
	assert.Equal(t, dest[1], uint16(0))
	assert.Equal(t, dest[2], uint16(0))

	// Get Error (index = 2)
	err = LittleEndian.Uint16Array(buf, 2, &dest)
	assert.NotNil(t, err)
	assert.Equal(t, dest[0], uint16(0))
	assert.Equal(t, dest[1], uint16(0))
	assert.Equal(t, dest[2], uint16(0))
}

func TestInt16Array(t *testing.T) {
	var buf = make([]byte, 6)
	var expected = []int16{0, -1, 2}

	// Put
	bytes, err := LittleEndian.PutInt16Array(buf, 0, expected)
	assert.Nil(t, err)
	assert.Equal(t, uint64(6), bytes)
	assert.Equal(t, expected[0], int16(binary.LittleEndian.Uint16(buf)))
	assert.Equal(t, expected[1], int16(binary.LittleEndian.Uint16(buf[2:])))
	assert.Equal(t, expected[2], int16(binary.LittleEndian.Uint16(buf[4:])))

	// Put error (index = -1)
	bytes, err = LittleEndian.PutInt16Array(buf, -1, expected)
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), bytes)

	// Put error
	bytes, err = LittleEndian.PutInt16Array(buf, 2, expected)
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), bytes)

	// Get
	dest := make([]int16, 3)
	err = LittleEndian.Int16Array(buf, 0, &dest)
	assert.Nil(t, err)
	assert.Equal(t, expected, dest)

	// Get Error (index = -1)
	dest = make([]int16, 3)
	err = LittleEndian.Int16Array(buf, -1, &dest)
	assert.NotNil(t, err)
	assert.Equal(t, dest[0], int16(0))
	assert.Equal(t, dest[1], int16(0))
	assert.Equal(t, dest[2], int16(0))

	// Get Error (index = 2)
	err = LittleEndian.Int16Array(buf, 2, &dest)
	assert.NotNil(t, err)
	assert.Equal(t, dest[0], int16(0))
	assert.Equal(t, dest[1], int16(0))
	assert.Equal(t, dest[2], int16(0))
}

func TestUint32Array(t *testing.T) {
	var buf = make([]byte, 12)
	var expected = []uint32{0, 1, 2}

	// Put
	bytes, err := LittleEndian.PutUint32Array(buf, 0, expected)
	assert.Nil(t, err)
	assert.Equal(t, uint64(12), bytes)
	assert.Equal(t, expected[0], binary.LittleEndian.Uint32(buf))
	assert.Equal(t, expected[1], binary.LittleEndian.Uint32(buf[4:]))
	assert.Equal(t, expected[2], binary.LittleEndian.Uint32(buf[8:]))

	// Put error (index = -1)
	bytes, err = LittleEndian.PutUint32Array(buf, -1, expected)
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), bytes)

	// Put error
	bytes, err = LittleEndian.PutUint32Array(buf, 2, expected)
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), bytes)

	// Get
	dest := make([]uint32, 3)
	err = LittleEndian.Uint32Array(buf, 0, &dest)
	assert.Nil(t, err)
	assert.Equal(t, expected, dest)

	// Get Error (index = -1)
	dest = make([]uint32, 3)
	err = LittleEndian.Uint32Array(buf, -1, &dest)
	assert.NotNil(t, err)
	assert.Equal(t, dest[0], uint32(0))
	assert.Equal(t, dest[1], uint32(0))
	assert.Equal(t, dest[2], uint32(0))

	// Get Error (index = 2)
	err = LittleEndian.Uint32Array(buf, 2, &dest)
	assert.NotNil(t, err)
	assert.Equal(t, dest[0], uint32(0))
	assert.Equal(t, dest[1], uint32(0))
	assert.Equal(t, dest[2], uint32(0))
}

func TestInt32Array(t *testing.T) {
	var buf = make([]byte, 12)
	var expected = []int32{0, -1, 2}

	// Put
	bytes, err := LittleEndian.PutInt32Array(buf, 0, expected)
	assert.Nil(t, err)
	assert.Equal(t, uint64(12), bytes)
	assert.Equal(t, expected[0], int32(binary.LittleEndian.Uint32(buf)))
	assert.Equal(t, expected[1], int32(binary.LittleEndian.Uint32(buf[4:])))
	assert.Equal(t, expected[2], int32(binary.LittleEndian.Uint32(buf[8:])))

	// Put error (index = -1)
	bytes, err = LittleEndian.PutInt32Array(buf, -1, expected)
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), bytes)

	// Put error
	bytes, err = LittleEndian.PutInt32Array(buf, 2, expected)
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), bytes)

	// Get
	dest := make([]int32, 3)
	err = LittleEndian.Int32Array(buf, 0, &dest)
	assert.Nil(t, err)
	assert.Equal(t, expected, dest)

	// Get Error (index = -1)
	dest = make([]int32, 3)
	err = LittleEndian.Int32Array(buf, -1, &dest)
	assert.NotNil(t, err)
	assert.Equal(t, dest[0], int32(0))
	assert.Equal(t, dest[1], int32(0))
	assert.Equal(t, dest[2], int32(0))

	// Get Error (index = 2)
	err = LittleEndian.Int32Array(buf, 2, &dest)
	assert.NotNil(t, err)
	assert.Equal(t, dest[0], int32(0))
	assert.Equal(t, dest[1], int32(0))
	assert.Equal(t, dest[2], int32(0))
}

func TestUint64Array(t *testing.T) {
	var buf = make([]byte, 24)
	var expected = []uint64{0, 1, 2}

	// Put
	bytes, err := LittleEndian.PutUint64Array(buf, 0, expected)
	assert.Nil(t, err)
	assert.Equal(t, uint64(24), bytes)
	assert.Equal(t, expected[0], binary.LittleEndian.Uint64(buf))
	assert.Equal(t, expected[1], binary.LittleEndian.Uint64(buf[8:]))
	assert.Equal(t, expected[2], binary.LittleEndian.Uint64(buf[16:]))
	// fmt.Println(expected)
	// fmt.Println(buf)

	// Put error (index = -1)
	bytes, err = LittleEndian.PutUint64Array(buf, -1, expected)
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), bytes)

	// Put error
	bytes, err = LittleEndian.PutUint64Array(buf, 2, expected)
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), bytes)

	// Get
	dest := make([]uint64, 3)
	err = LittleEndian.Uint64Array(buf, 0, &dest)
	assert.Nil(t, err)
	// fmt.Println(expected)
	// fmt.Println(dest)
	assert.Equal(t, expected, dest)

	// Get Error (index = -1)
	dest = make([]uint64, 3)
	err = LittleEndian.Uint64Array(buf, -1, &dest)
	assert.NotNil(t, err)
	assert.Equal(t, dest[0], uint64(0))
	assert.Equal(t, dest[1], uint64(0))
	assert.Equal(t, dest[2], uint64(0))

	// Get Error (index = 2)
	err = LittleEndian.Uint64Array(buf, 2, &dest)
	assert.NotNil(t, err)
	assert.Equal(t, dest[0], uint64(0))
	assert.Equal(t, dest[1], uint64(0))
	assert.Equal(t, dest[2], uint64(0))
}

func TestInt64Array(t *testing.T) {
	var buf = make([]byte, 24)
	var expected = []int64{0, -1, 2}

	// Put
	bytes, err := LittleEndian.PutInt64Array(buf, 0, expected)
	assert.Nil(t, err)
	assert.Equal(t, uint64(24), bytes)
	assert.Equal(t, expected[0], int64(binary.LittleEndian.Uint64(buf)))
	assert.Equal(t, expected[1], int64(binary.LittleEndian.Uint64(buf[8:])))
	assert.Equal(t, expected[2], int64(binary.LittleEndian.Uint64(buf[16:])))

	// Put error (index = -1)
	bytes, err = LittleEndian.PutInt64Array(buf, -1, expected)
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), bytes)

	// Put error
	bytes, err = LittleEndian.PutInt64Array(buf, 2, expected)
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), bytes)

	// Get
	dest := make([]int64, 3)
	err = LittleEndian.Int64Array(buf, 0, &dest)
	assert.Nil(t, err)
	assert.Equal(t, expected, dest)

	// Get Error (index = -1)
	dest = make([]int64, 3)
	err = LittleEndian.Int64Array(buf, -1, &dest)
	assert.NotNil(t, err)
	assert.Equal(t, dest[0], int64(0))
	assert.Equal(t, dest[1], int64(0))
	assert.Equal(t, dest[2], int64(0))

	// Get Error (index = 2)
	err = LittleEndian.Int64Array(buf, 2, &dest)
	assert.NotNil(t, err)
	assert.Equal(t, dest[0], int64(0))
	assert.Equal(t, dest[1], int64(0))
	assert.Equal(t, dest[2], int64(0))
}

func TestFloat32Array(t *testing.T) {
	var buf = make([]byte, 12)
	var expected = []float32{0.0, 1.0, 2.0}

	// Put
	bytes, err := LittleEndian.PutFloat32Array(buf, 0, expected)
	assert.Nil(t, err)
	assert.Equal(t, uint64(12), bytes)
	assert.Equal(t, expected[0], math.Float32frombits(binary.LittleEndian.Uint32(buf)))
	assert.Equal(t, expected[1], math.Float32frombits(binary.LittleEndian.Uint32(buf[4:])))
	assert.Equal(t, expected[2], math.Float32frombits(binary.LittleEndian.Uint32(buf[8:])))

	// Put error (index = -1)
	bytes, err = LittleEndian.PutFloat32Array(buf, -1, expected)
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), bytes)

	// Put error
	bytes, err = LittleEndian.PutFloat32Array(buf, 2, expected)
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), bytes)

	// Get
	dest := make([]float32, 3)
	err = LittleEndian.Float32Array(buf, 0, &dest)
	assert.Nil(t, err)
	assert.Equal(t, expected, dest)

	// Get Error (index = -1)
	dest = make([]float32, 3)
	err = LittleEndian.Float32Array(buf, -1, &dest)
	assert.NotNil(t, err)
	assert.Equal(t, dest[0], float32(0.0))
	assert.Equal(t, dest[1], float32(0.0))
	assert.Equal(t, dest[2], float32(0.0))

	// Get Error (index = 2)
	err = LittleEndian.Float32Array(buf, 2, &dest)
	assert.NotNil(t, err)
	assert.Equal(t, dest[0], float32(0.0))
	assert.Equal(t, dest[1], float32(0.0))
	assert.Equal(t, dest[2], float32(0.0))
}

func TestFloat64Array(t *testing.T) {
	var buf = make([]byte, 24)
	var expected = []float64{0.0, 1.0, 2.0}

	// Put
	bytes, err := LittleEndian.PutFloat64Array(buf, 0, expected)
	assert.Nil(t, err)
	assert.Equal(t, uint64(24), bytes)
	assert.Equal(t, expected[0], math.Float64frombits(binary.LittleEndian.Uint64(buf)))
	assert.Equal(t, expected[1], math.Float64frombits(binary.LittleEndian.Uint64(buf[8:])))
	assert.Equal(t, expected[2], math.Float64frombits(binary.LittleEndian.Uint64(buf[16:])))

	// Put error (index = -1)
	bytes, err = LittleEndian.PutFloat64Array(buf, -1, expected)
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), bytes)

	// Put error
	bytes, err = LittleEndian.PutFloat64Array(buf, 2, expected)
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), bytes)

	// Get
	dest := make([]float64, 3)
	err = LittleEndian.Float64Array(buf, 0, &dest)
	assert.Nil(t, err)
	assert.Equal(t, expected, dest)

	// Get Error (index = -1)
	dest = make([]float64, 3)
	err = LittleEndian.Float64Array(buf, -1, &dest)
	assert.NotNil(t, err)
	assert.Equal(t, dest[0], float64(0.0))
	assert.Equal(t, dest[1], float64(0.0))
	assert.Equal(t, dest[2], float64(0.0))

	// Get Error (index = 2)
	err = LittleEndian.Float64Array(buf, 2, &dest)
	assert.NotNil(t, err)
	assert.Equal(t, dest[0], float64(0.0))
	assert.Equal(t, dest[1], float64(0.0))
	assert.Equal(t, dest[2], float64(0.0))
}
