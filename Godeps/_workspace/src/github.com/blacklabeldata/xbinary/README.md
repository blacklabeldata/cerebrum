# xbinary

[![Coverage Status](https://coveralls.io/repos/blacklabeldata/xbinary/badge.svg?branch=master&service=github)](https://coveralls.io/github/blacklabeldata/xbinary?branch=master) [![Build Status](https://travis-ci.org/blacklabeldata/xbinary.svg)](https://travis-ci.org/blacklabeldata/xbinary)

Extends the functionality provided by the 'binary' package in the Go standard lib. Adds
capability to read/write all signed and unsigned integer types as well as both float32 and
 float 64. Arrays for all numeric types are also supported.

--
    import "github.com/blacklabeldata/xbinary"


## Usage

```go
var BigEndian bigEndian
```

```go
var LittleEndian littleEndian
```

### Errors

```go
var ErrOutOfRange = fmt.Errorf("Index out of range")
```

#### type ExtendedBuffer


ExtendedBuffer interface simply wraps common binary encoding/decoding.


```go
type ExtendedBuffer interface {

	// Reads a uint16 at the specified index in a byte array.
	// Returns an ErrOutOfRange if the index provided is outside the bounds of the byte array.
	Uint16(b []byte, index int) (uint16, error)

	// Reads a uint32 at the specified index in a byte array.
	// Returns an ErrOutOfRange if the index provided is outside the bounds of the byte array.
	Uint32(b []byte, index int) (uint32, error)

	// Reads a uint64 at the specified index in a byte array.
	// Returns an ErrOutOfRange if the index provided is outside the bounds of the byte array.
	Uint64(b []byte, index int) (uint64, error)

	// Reads a int16 at the specified index in a byte array.
	// Returns an ErrOutOfRange if the index provided is outside the bounds of the byte array.
	Int16(b []byte, index int) (int16, error)

	// Reads a int32 at the specified index in a byte array.
	// Returns an ErrOutOfRange if the index provided is outside the bounds of the byte array.
	Int32(b []byte, index int) (int32, error)

	// Reads a int64 at the specified index in a byte array.
	// Returns an ErrOutOfRange if the index provided is outside the bounds of the byte array.
	Int64(b []byte, index int) (int64, error)

	// Reads a float32 at the specified index in a byte array.
	// Returns an ErrOutOfRange if the index provided is outside the bounds of the byte array.
	Float32(b []byte, index int) (float32, error)

	// Reads a float64 at the specified index in a byte array.
	// Returns an ErrOutOfRange if the index provided is outside the bounds of the byte array.
	Float64(b []byte, index int) (float64, error)

	// Reads a string at the specified index in a byte array.
	// Returns an ErrOutOfRange if the index provided is outside the bounds of the byte array.
	String(b []byte, index, size int) (string, error)

	// Inserts a uint16 into the byte array at the given index.
	// Returns the number of bytes written and an ErrOutOfRange
	// if there is no more space in the byte array.
	PutUint16(b []byte, index int, value uint16) (uint64, error)

	// Inserts a Uint32 into the byte array at the given index.
	// Returns the number of bytes written and an ErrOutOfRange
	// if there is no more space in the byte array.
	PutUint32(b []byte, index int, value uint32) (uint64, error)

	// Inserts a Uint64 into the byte array at the given index.
	// Returns the number of bytes written and an ErrOutOfRange
	// if there is no more space in the byte array.
	PutUint64(b []byte, index int, value uint64) (uint64, error)

	// Inserts a int16 into the byte array at the given index.
	// Returns the number of bytes written and an ErrOutOfRange
	// if there is no more space in the byte array.
	PutInt16(b []byte, index int, value int16) (uint64, error)

	// Inserts a int32 into the byte array at the given index.
	// Returns the number of bytes written and an ErrOutOfRange
	// if there is no more space in the byte array.
	PutInt32(b []byte, index int, value int32) (uint64, error)

	// Inserts a int64 into the byte array at the given index.
	// Returns the number of bytes written and an ErrOutOfRange
	// if there is no more space in the byte array.
	PutInt64(b []byte, index int, value int64) (uint64, error)

	// Inserts a Float32 into the byte array at the given index.
	// Returns the number of bytes written and an ErrOutOfRange
	// if there is no more space in the byte array.
	PutFloat32(b []byte, index int, value float32) (uint64, error)

	// Inserts a Float64 into the byte array at the given index.
	// Returns the number of bytes written and an ErrOutOfRange
	// if there is no more space in the byte array.
	PutFloat64(b []byte, index int, value float64) (uint64, error)

	// Inserts a String into the byte array at the given index.
	// Returns the number of bytes written and an ErrOutOfRange
	// if there is no more space in the byte array.
	PutString(b []byte, index int, value string) (uint64, error)

	// Fills in the given Uint16 array at the given offset into the byte array.
	// If the index given and the total size of the array is too large to fit
	// in the byte array, an ErrOutOfRange error will be returned.
	Uint16Array(b []byte, index int, dest *[]uint16) error

	// Fills in the given Uint32 array at the given offset into the byte array.
	// If the index given and the total size of the array is too large to fit
	// in the byte array, an ErrOutOfRange error will be returned.
	Uint32Array(b []byte, index int, dest *[]uint32) error

	// Fills in the given Uint64 array at the given offset into the byte array.
	// If the index given and the total size of the array is too large to fit
	// in the byte array, an ErrOutOfRange error will be returned.
	Uint64Array(b []byte, index int, dest *[]uint64) error

	// Fills in the given int16 array at the given offset into the byte array.
	// If the index given and the total size of the array is too large to fit
	// in the byte array, an ErrOutOfRange error will be returned.
	Int16Array(b []byte, index int, dest *[]uint16) error

	// Fills in the given int32 array at the given offset into the byte array.
	// If the index given and the total size of the array is too large to fit
	// in the byte array, an ErrOutOfRange error will be returned.
	Int32Array(b []byte, index int, dest *[]uint32) error

	// Fills in the given int64 array at the given offset into the byte array.
	// If the index given and the total size of the array is too large to fit
	// in the byte array, an ErrOutOfRange error will be returned.
	Int64Array(b []byte, index int, dest *[]uint64) error

	// Fills in the given Float32 array at the given offset into the byte array.
	// If the index given and the total size of the array is too large to fit
	// in the byte array, an ErrOutOfRange error will be returned.
	Float32Array(b []byte, index int, dest *[]float32) error

	// Fills in the given Float64 array at the given offset into the byte array.
	// If the index given and the total size of the array is too large to fit
	// in the byte array, an ErrOutOfRange error will be returned.
	Float64Array(b []byte, index int, dest *[]float64) error

	// Inserts the Uint16 array into the byte array at the given offset.
	// The total number of bytes written are returned. If the index and
	// the size given array doesn't fit in the remaining byte array an
	// ErrOutOfRange error is returned.
	PutUint16Array(b []byte, index int, value []uint16) (uint64, error)

	// Inserts the Uint32 array into the byte array at the given offset.
	// The total number of bytes written are returned. If the index and
	// the size given array doesn't fit in the remaining byte array an
	// ErrOutOfRange error is returned.
	PutUint32Array(b []byte, index int, value []uint32) (uint64, error)

	// Inserts the Uint64 array into the byte array at the given offset.
	// The total number of bytes written are returned. If the index and
	// the size given array doesn't fit in the remaining byte array an
	// ErrOutOfRange error is returned.
	PutUint64Array(b []byte, index int, value []uint64) (uint64, error)

	// Inserts the Float32 array into the byte array at the given offset.
	// The total number of bytes written are returned. If the index and
	// the size given array doesn't fit in the remaining byte array an
	// ErrOutOfRange error is returned.
	PutFloat32Array(b []byte, index int, value []float32) (uint64, error)

	// Inserts the Float64 array into the byte array at the given offset.
	// The total number of bytes written are returned. If the index and
	// the size given array doesn't fit in the remaining byte array an
	// ErrOutOfRange error is returned.
	PutFloat64Array(b []byte, index int, value []float64) (uint64, error)
}
```
