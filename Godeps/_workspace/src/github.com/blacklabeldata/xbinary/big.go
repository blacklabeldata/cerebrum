package xbinary

import (
	"encoding/binary"
	"math"
)

type bigEndian struct{}

var BigEndian bigEndian

func (bigEndian) Uint16(b []byte, index int) (uint16, error) {
	if len(b) < index+1 || index < 0 {
		return 0, ErrOutOfRange
	}
	return binary.BigEndian.Uint16(b[index:]), nil
}

func (bigEndian) Int16(b []byte, index int) (int16, error) {
	if len(b) < index+1 || index < 0 {
		return 0, ErrOutOfRange
	}
	return int16(binary.BigEndian.Uint16(b[index:])), nil
}

func (bigEndian) Uint32(b []byte, index int) (uint32, error) {
	if len(b) < index+3 || index < 0 {
		return 0, ErrOutOfRange
	}
	return binary.BigEndian.Uint32(b[index:]), nil
}

func (bigEndian) Int32(b []byte, index int) (int32, error) {
	if len(b) < index+3 || index < 0 {
		return 0, ErrOutOfRange
	}
	return int32(binary.BigEndian.Uint32(b[index:])), nil
}

func (bigEndian) Uint64(b []byte, index int) (uint64, error) {
	if len(b) < index+7 || index < 0 {
		return 0, ErrOutOfRange
	}
	return binary.BigEndian.Uint64(b[index:]), nil
}

func (bigEndian) Int64(b []byte, index int) (int64, error) {
	if len(b) < index+7 || index < 0 {
		return 0, ErrOutOfRange
	}
	return int64(binary.BigEndian.Uint64(b[index:])), nil
}

func (bigEndian) Float32(b []byte, index int) (float32, error) {
	if len(b) < index+3 || index < 0 {
		return 0, ErrOutOfRange
	}
	return math.Float32frombits(binary.BigEndian.Uint32(b[index:])), nil
}

func (l bigEndian) Float64(b []byte, index int) (float64, error) {
	val, err := l.Uint64(b, index)
	if err != nil {
		return 0.0, err
	}
	return math.Float64frombits(val), nil
}

func (bigEndian) PutUint16(b []byte, index int, v uint16) (uint64, error) {
	if len(b) < index+1 || index < 0 {
		return 0, ErrOutOfRange
	}
	binary.BigEndian.PutUint16(b[index:], v)
	return 2, nil
}

func (bigEndian) PutUint32(b []byte, index int, v uint32) (uint64, error) {
	if len(b) < index+3 || index < 0 {
		return 0, ErrOutOfRange
	}
	binary.BigEndian.PutUint32(b[index:], v)
	return 4, nil
}

func (bigEndian) PutUint64(b []byte, index int, v uint64) (uint64, error) {
	if len(b) < index+7 || index < 0 {
		return 0, ErrOutOfRange
	}
	binary.BigEndian.PutUint64(b[index:], v)
	return 8, nil
}

func (bigEndian) PutInt16(b []byte, index int, v int16) (uint64, error) {
	if len(b) < index+1 || index < 0 {
		return 0, ErrOutOfRange
	}
	binary.BigEndian.PutUint16(b[index:], uint16(v))
	return 2, nil
}

func (bigEndian) PutInt32(b []byte, index int, v int32) (uint64, error) {
	if len(b) < index+3 || index < 0 {
		return 0, ErrOutOfRange
	}
	binary.BigEndian.PutUint32(b[index:], uint32(v))
	return 4, nil
}

func (bigEndian) PutInt64(b []byte, index int, v int64) (uint64, error) {
	if len(b) < index+7 || index < 0 {
		return 0, ErrOutOfRange
	}
	binary.BigEndian.PutUint64(b[index:], uint64(v))
	return 8, nil
}

func (bigEndian) PutFloat64(b []byte, index int, v float64) (uint64, error) {
	if len(b) < index+7 || index < 0 {
		return 0, ErrOutOfRange
	}
	binary.BigEndian.PutUint64(b[index:], math.Float64bits(v))
	return 8, nil
}

func (bigEndian) PutFloat32(b []byte, index int, v float32) (uint64, error) {
	if len(b) < index+3 || index < 0 {
		return 0, ErrOutOfRange
	}
	binary.BigEndian.PutUint32(b[index:], math.Float32bits(v))
	return 4, nil
}

func (bigEndian) Uint16Array(b []byte, index int, dest *[]uint16) error {
	if (len(*dest)*2+index) > len(b) || index < 0 {
		return ErrOutOfRange
	}
	for i := 0; i < len(*dest); i++ {
		(*dest)[i] = binary.BigEndian.Uint16(b[index+i*2:])
	}
	return nil
}

func (bigEndian) Uint32Array(b []byte, index int, dest *[]uint32) error {
	if (len(*dest)*4+index) > len(b) || index < 0 {
		return ErrOutOfRange
	}
	for i := 0; i < len(*dest); i++ {
		(*dest)[i] = binary.BigEndian.Uint32(b[index+i*4:])
	}
	return nil
}

func (bigEndian) Uint64Array(b []byte, index int, dest *[]uint64) error {
	if (len(*dest)*8+index) > len(b) || index < 0 {
		return ErrOutOfRange
	}
	for i := 0; i < len(*dest); i++ {
		(*dest)[i] = binary.BigEndian.Uint64(b[index+i*8:])
	}
	return nil
}

func (bigEndian) Int16Array(b []byte, index int, dest *[]int16) error {
	if (len(*dest)*2+index) > len(b) || index < 0 {
		return ErrOutOfRange
	}
	for i := 0; i < len(*dest); i++ {
		(*dest)[i] = int16(binary.BigEndian.Uint16(b[index+i*2:]))
	}
	return nil
}

func (bigEndian) Int32Array(b []byte, index int, dest *[]int32) error {
	if (len(*dest)*4+index) > len(b) || index < 0 {
		return ErrOutOfRange
	}
	for i := 0; i < len(*dest); i++ {
		(*dest)[i] = int32(binary.BigEndian.Uint32(b[index+i*4:]))
	}
	return nil
}

func (bigEndian) Int64Array(b []byte, index int, dest *[]int64) error {
	if (len(*dest)*8+index) > len(b) || index < 0 {
		return ErrOutOfRange
	}
	for i := 0; i < len(*dest); i++ {
		(*dest)[i] = int64(binary.BigEndian.Uint64(b[index+i*8:]))
	}
	return nil
}
func (bigEndian) Float32Array(b []byte, index int, dest *[]float32) error {
	if (len(*dest)*4+index) > len(b) || index < 0 {
		return ErrOutOfRange
	}
	for i := 0; i < len(*dest); i++ {
		(*dest)[i] = math.Float32frombits(binary.BigEndian.Uint32(b[index+i*4:]))
	}
	return nil
}

func (bigEndian) Float64Array(b []byte, index int, dest *[]float64) error {
	if (len(*dest)*8+index) > len(b) || index < 0 {
		return ErrOutOfRange
	}
	for i := 0; i < len(*dest); i++ {
		(*dest)[i] = math.Float64frombits(binary.BigEndian.Uint64(b[index+i*8:]))
	}
	return nil
}

func (bigEndian) PutUint16Array(b []byte, index int, value []uint16) (uint64, error) {
	var reqsize = (len(value)*2 + index)
	if reqsize > len(b) || index < 0 {
		return 0, ErrOutOfRange
	}
	for i := 0; i < len(value); i++ {
		binary.BigEndian.PutUint16(b[index+i*2:], value[i])
	}
	return uint64(len(value) * 2), nil
}

func (bigEndian) PutUint32Array(b []byte, index int, value []uint32) (uint64, error) {
	var reqsize = (len(value)*4 + index)
	if reqsize > len(b) || index < 0 {
		return 0, ErrOutOfRange
	}
	for i := 0; i < len(value); i++ {
		binary.BigEndian.PutUint32(b[index+i*4:], value[i])
	}
	return uint64(reqsize - index), nil
}

func (bigEndian) PutUint64Array(b []byte, index int, value []uint64) (uint64, error) {
	var reqsize = (len(value)*8 + index)
	if reqsize > len(b) || index < 0 {
		return 0, ErrOutOfRange
	}
	for i := 0; i < len(value); i++ {
		binary.BigEndian.PutUint64(b[index+i*8:], value[i])
	}
	return uint64(reqsize - index), nil
}

func (bigEndian) PutInt16Array(b []byte, index int, value []int16) (uint64, error) {
	var reqsize = (len(value)*2 + index)
	if reqsize > len(b) || index < 0 {
		return 0, ErrOutOfRange
	}
	for i := 0; i < len(value); i++ {
		binary.BigEndian.PutUint16(b[index+i*2:], uint16(value[i]))
	}
	return uint64(len(value) * 2), nil
}

func (bigEndian) PutInt32Array(b []byte, index int, value []int32) (uint64, error) {
	var reqsize = (len(value)*4 + index)
	if reqsize > len(b) || index < 0 {
		return 0, ErrOutOfRange
	}
	for i := 0; i < len(value); i++ {
		binary.BigEndian.PutUint32(b[index+i*4:], uint32(value[i]))
	}
	return uint64(reqsize - index), nil
}

func (bigEndian) PutInt64Array(b []byte, index int, value []int64) (uint64, error) {
	var reqsize = (len(value)*8 + index)
	if reqsize > len(b) || index < 0 {
		return 0, ErrOutOfRange
	}
	for i := 0; i < len(value); i++ {
		binary.BigEndian.PutUint64(b[index+i*8:], uint64(value[i]))
	}
	return uint64(reqsize - index), nil
}

func (bigEndian) PutFloat32Array(b []byte, index int, value []float32) (uint64, error) {
	var reqsize = (len(value)*4 + index)
	if reqsize > len(b) || index < 0 {
		return 0, ErrOutOfRange
	}
	for i := 0; i < len(value); i++ {
		binary.BigEndian.PutUint32(b[index+i*4:], math.Float32bits(value[i]))
	}
	return uint64(reqsize - index), nil
}

func (bigEndian) PutFloat64Array(b []byte, index int, value []float64) (uint64, error) {
	var reqsize = (len(value)*8 + index)
	if reqsize > len(b) || index < 0 {
		return 0, ErrOutOfRange
	}
	for i := 0; i < len(value); i++ {
		binary.BigEndian.PutUint64(b[index+i*8:], math.Float64bits(value[i]))
	}
	return uint64(reqsize - index), nil
}

func (bigEndian) String(b []byte, index, size int) (string, error) {
	if index+size > len(b) || index < 0 {
		return "", ErrOutOfRange
	}
	return string(b[index:]), nil
}

func (bigEndian) PutString(b []byte, index int, value string) (uint64, error) {
	if index+len(value) > len(b) || index < 0 {
		return 0, ErrOutOfRange
	}
	copy(b[index:index+len(value)], []byte(value))
	return uint64(len(value)), nil
}
