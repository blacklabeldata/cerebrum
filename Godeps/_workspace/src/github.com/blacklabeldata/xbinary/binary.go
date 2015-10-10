package xbinary

import (
	"fmt"
)

// Errors
var ErrOutOfRange = fmt.Errorf("Index out of range")

// ExtendedBuffer interface simply wraps common binary encoding/decoding.
type ExtendedBuffer interface {

	// Reads a uint8 at the specified index in a byte array.
	// Returns an ErrOutOfRange if the index provided is outside the bounds of the byte array.
	Uint8(b []byte, index int) (uint8, error)

	// Reads a uint16 at the specified index in a byte array.
	// Returns an ErrOutOfRange if the index provided is outside the bounds of the byte array.
	Uint16(b []byte, index int) (uint16, error)

	// Reads a uint32 at the specified index in a byte array.
	// Returns an ErrOutOfRange if the index provided is outside the bounds of the byte array.
	Uint32(b []byte, index int) (uint32, error)

	// Reads a uint64 at the specified index in a byte array.
	// Returns an ErrOutOfRange if the index provided is outside the bounds of the byte array.
	Uint64(b []byte, index int) (uint64, error)

	// Reads a int8 at the specified index in a byte array.
	// Returns an ErrOutOfRange if the index provided is outside the bounds of the byte array.
	Int8(b []byte, index int) (int8, error)

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
	PutUint8(b []byte, index int, value uint8) (int, error)

	// Inserts a uint16 into the byte array at the given index.
	// Returns the number of bytes written and an ErrOutOfRange
	// if there is no more space in the byte array.
	PutUint16(b []byte, index int, value uint16) (int, error)

	// Inserts a Uint32 into the byte array at the given index.
	// Returns the number of bytes written and an ErrOutOfRange
	// if there is no more space in the byte array.
	PutUint32(b []byte, index int, value uint32) (int, error)

	// Inserts a Uint64 into the byte array at the given index.
	// Returns the number of bytes written and an ErrOutOfRange
	// if there is no more space in the byte array.
	PutUint64(b []byte, index int, value uint64) (int, error)

	// Inserts a int8 into the byte array at the given index.
	// Returns the number of bytes written and an ErrOutOfRange
	// if there is no more space in the byte array.
	PutInt8(b []byte, index int, value int8) (int, error)

	// Inserts a int16 into the byte array at the given index.
	// Returns the number of bytes written and an ErrOutOfRange
	// if there is no more space in the byte array.
	PutInt16(b []byte, index int, value int16) (int, error)

	// Inserts a int32 into the byte array at the given index.
	// Returns the number of bytes written and an ErrOutOfRange
	// if there is no more space in the byte array.
	PutInt32(b []byte, index int, value int32) (int, error)

	// Inserts a int64 into the byte array at the given index.
	// Returns the number of bytes written and an ErrOutOfRange
	// if there is no more space in the byte array.
	PutInt64(b []byte, index int, value int64) (int, error)

	// Inserts a Float32 into the byte array at the given index.
	// Returns the number of bytes written and an ErrOutOfRange
	// if there is no more space in the byte array.
	PutFloat32(b []byte, index int, value float32) (int, error)

	// Inserts a Float64 into the byte array at the given index.
	// Returns the number of bytes written and an ErrOutOfRange
	// if there is no more space in the byte array.
	PutFloat64(b []byte, index int, value float64) (int, error)

	// Inserts a String into the byte array at the given index.
	// Returns the number of bytes written and an ErrOutOfRange
	// if there is no more space in the byte array.
	PutString(b []byte, index int, value string) (int, error)

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
	Int16Array(b []byte, index int, dest *[]int16) error

	// Fills in the given int32 array at the given offset into the byte array.
	// If the index given and the total size of the array is too large to fit
	// in the byte array, an ErrOutOfRange error will be returned.
	Int32Array(b []byte, index int, dest *[]int32) error

	// Fills in the given int64 array at the given offset into the byte array.
	// If the index given and the total size of the array is too large to fit
	// in the byte array, an ErrOutOfRange error will be returned.
	Int64Array(b []byte, index int, dest *[]int64) error

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
	PutUint16Array(b []byte, index int, value []uint16) (int, error)

	// Inserts the Uint32 array into the byte array at the given offset.
	// The total number of bytes written are returned. If the index and
	// the size given array doesn't fit in the remaining byte array an
	// ErrOutOfRange error is returned.
	PutUint32Array(b []byte, index int, value []uint32) (int, error)

	// Inserts the Uint64 array into the byte array at the given offset.
	// The total number of bytes written are returned. If the index and
	// the size given array doesn't fit in the remaining byte array an
	// ErrOutOfRange error is returned.
	PutUint64Array(b []byte, index int, value []uint64) (int, error)

	// Inserts the Float32 array into the byte array at the given offset.
	// The total number of bytes written are returned. If the index and
	// the size given array doesn't fit in the remaining byte array an
	// ErrOutOfRange error is returned.
	PutFloat32Array(b []byte, index int, value []float32) (int, error)

	// Inserts the Float64 array into the byte array at the given offset.
	// The total number of bytes written are returned. If the index and
	// the size given array doesn't fit in the remaining byte array an
	// ErrOutOfRange error is returned.
	PutFloat64Array(b []byte, index int, value []float64) (int, error)
}
