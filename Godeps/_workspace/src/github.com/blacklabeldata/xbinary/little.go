package xbinary

import (
	"encoding/binary"
	"math"
)

type littleEndian struct{}

var LittleEndian littleEndian

func (littleEndian) Uint8(b []byte, index int) (uint8, error) {
	if len(b) < index || index < 0 {
		return 0, ErrOutOfRange
	}
	return uint8(b[index]), nil
}

func (littleEndian) Int8(b []byte, index int) (int8, error) {
	if len(b) < index || index < 0 {
		return 0, ErrOutOfRange
	}
	return int8(b[index]), nil
}

func (littleEndian) Uint16(b []byte, index int) (uint16, error) {
	if len(b) < index+1 || index < 0 {
		return 0, ErrOutOfRange
	}
	return binary.LittleEndian.Uint16(b[index:]), nil
}

func (littleEndian) Int16(b []byte, index int) (int16, error) {
	if len(b) < index+1 || index < 0 {
		return 0, ErrOutOfRange
	}
	return int16(binary.LittleEndian.Uint16(b[index:])), nil
}

func (littleEndian) Uint32(b []byte, index int) (uint32, error) {
	if len(b) < index+3 || index < 0 {
		return 0, ErrOutOfRange
	}
	return binary.LittleEndian.Uint32(b[index:]), nil
}

func (littleEndian) Int32(b []byte, index int) (int32, error) {
	if len(b) < index+3 || index < 0 {
		return 0, ErrOutOfRange
	}
	return int32(binary.LittleEndian.Uint32(b[index:])), nil
}

func (littleEndian) Uint64(b []byte, index int) (uint64, error) {
	if len(b) < index+7 || index < 0 {
		return 0, ErrOutOfRange
	}
	return binary.LittleEndian.Uint64(b[index:]), nil
}

func (littleEndian) Int64(b []byte, index int) (int64, error) {
	if len(b) < index+7 || index < 0 {
		return 0, ErrOutOfRange
	}
	return int64(binary.LittleEndian.Uint64(b[index:])), nil
}

func (littleEndian) Float32(b []byte, index int) (float32, error) {
	if len(b) < index+3 || index < 0 {
		return 0, ErrOutOfRange
	}
	return math.Float32frombits(binary.LittleEndian.Uint32(b[index:])), nil
}

func (l littleEndian) Float64(b []byte, index int) (float64, error) {
	val, err := l.Uint64(b, index)
	if err != nil {
		return 0.0, err
	}
	return math.Float64frombits(val), nil
}

func (littleEndian) PutUint8(b []byte, index int, v uint8) (uint64, error) {
	if len(b) < index || index < 0 {
		return 0, ErrOutOfRange
	}
	b[index] = byte(v)
	return 1, nil
}

func (littleEndian) PutUint16(b []byte, index int, v uint16) (uint64, error) {
	if len(b) < index+1 || index < 0 {
		return 0, ErrOutOfRange
	}
	binary.LittleEndian.PutUint16(b[index:], v)
	return 2, nil
}

func (littleEndian) PutUint32(b []byte, index int, v uint32) (uint64, error) {
	if len(b) < index+3 || index < 0 {
		return 0, ErrOutOfRange
	}
	binary.LittleEndian.PutUint32(b[index:], v)
	return 4, nil
}

func (littleEndian) PutUint64(b []byte, index int, v uint64) (uint64, error) {
	if len(b) < index+7 || index < 0 {
		return 0, ErrOutOfRange
	}
	binary.LittleEndian.PutUint64(b[index:], v)
	return 8, nil
}

func (littleEndian) PutInt8(b []byte, index int, v int8) (uint64, error) {
	if len(b) < index || index < 0 {
		return 0, ErrOutOfRange
	}
	b[index] = byte(v)
	return 1, nil
}

func (littleEndian) PutInt16(b []byte, index int, v int16) (uint64, error) {
	if len(b) < index+1 || index < 0 {
		return 0, ErrOutOfRange
	}
	binary.LittleEndian.PutUint16(b[index:], uint16(v))
	return 2, nil
}

func (littleEndian) PutInt32(b []byte, index int, v int32) (uint64, error) {
	if len(b) < index+3 || index < 0 {
		return 0, ErrOutOfRange
	}
	binary.LittleEndian.PutUint32(b[index:], uint32(v))
	return 4, nil
}

func (littleEndian) PutInt64(b []byte, index int, v int64) (uint64, error) {
	if len(b) < index+7 || index < 0 {
		return 0, ErrOutOfRange
	}
	binary.LittleEndian.PutUint64(b[index:], uint64(v))
	return 8, nil
}

func (littleEndian) PutFloat64(b []byte, index int, v float64) (uint64, error) {
	if len(b) < index+7 || index < 0 {
		return 0, ErrOutOfRange
	}
	binary.LittleEndian.PutUint64(b[index:], math.Float64bits(v))
	return 8, nil
}

func (littleEndian) PutFloat32(b []byte, index int, v float32) (uint64, error) {
	if len(b) < index+3 || index < 0 {
		return 0, ErrOutOfRange
	}
	binary.LittleEndian.PutUint32(b[index:], math.Float32bits(v))
	return 4, nil
}

func (littleEndian) Uint16Array(b []byte, index int, dest *[]uint16) error {
	if (len(*dest)*2+index) > len(b) || index < 0 {
		return ErrOutOfRange
	}
	for i := 0; i < len(*dest); i++ {
		(*dest)[i] = binary.LittleEndian.Uint16(b[index+i*2:])
	}
	return nil
}

func (littleEndian) Uint32Array(b []byte, index int, dest *[]uint32) error {
	if (len(*dest)*4+index) > len(b) || index < 0 {
		return ErrOutOfRange
	}
	for i := 0; i < len(*dest); i++ {
		(*dest)[i] = binary.LittleEndian.Uint32(b[index+i*4:])
	}
	return nil
}

func (littleEndian) Uint64Array(b []byte, index int, dest *[]uint64) error {
	if (len(*dest)*8+index) > len(b) || index < 0 {
		return ErrOutOfRange
	}
	for i := 0; i < len(*dest); i++ {
		(*dest)[i] = binary.LittleEndian.Uint64(b[index+i*8:])
	}
	return nil
}

func (littleEndian) Int16Array(b []byte, index int, dest *[]int16) error {
	if (len(*dest)*2+index) > len(b) || index < 0 {
		return ErrOutOfRange
	}
	for i := 0; i < len(*dest); i++ {
		(*dest)[i] = int16(binary.LittleEndian.Uint16(b[index+i*2:]))
	}
	return nil
}

func (littleEndian) Int32Array(b []byte, index int, dest *[]int32) error {
	if (len(*dest)*4+index) > len(b) || index < 0 {
		return ErrOutOfRange
	}
	for i := 0; i < len(*dest); i++ {
		(*dest)[i] = int32(binary.LittleEndian.Uint32(b[index+i*4:]))
	}
	return nil
}

func (littleEndian) Int64Array(b []byte, index int, dest *[]int64) error {
	if (len(*dest)*8+index) > len(b) || index < 0 {
		return ErrOutOfRange
	}
	for i := 0; i < len(*dest); i++ {
		(*dest)[i] = int64(binary.LittleEndian.Uint64(b[index+i*8:]))
	}
	return nil
}
func (littleEndian) Float32Array(b []byte, index int, dest *[]float32) error {
	if (len(*dest)*4+index) > len(b) || index < 0 {
		return ErrOutOfRange
	}
	for i := 0; i < len(*dest); i++ {
		(*dest)[i] = math.Float32frombits(binary.LittleEndian.Uint32(b[index+i*4:]))
	}
	return nil
}

func (littleEndian) Float64Array(b []byte, index int, dest *[]float64) error {
	if (len(*dest)*8+index) > len(b) || index < 0 {
		return ErrOutOfRange
	}
	for i := 0; i < len(*dest); i++ {
		(*dest)[i] = math.Float64frombits(binary.LittleEndian.Uint64(b[index+i*8:]))
	}
	return nil
}

func (littleEndian) PutUint16Array(b []byte, index int, value []uint16) (uint64, error) {
	var reqsize = (len(value)*2 + index)
	if reqsize > len(b) || index < 0 {
		return 0, ErrOutOfRange
	}
	for i := 0; i < len(value); i++ {
		binary.LittleEndian.PutUint16(b[index+i*2:], value[i])
	}
	return uint64(len(value) * 2), nil
}

func (littleEndian) PutUint32Array(b []byte, index int, value []uint32) (uint64, error) {
	var reqsize = (len(value)*4 + index)
	if reqsize > len(b) || index < 0 {
		return 0, ErrOutOfRange
	}
	for i := 0; i < len(value); i++ {
		binary.LittleEndian.PutUint32(b[index+i*4:], value[i])
	}
	return uint64(reqsize - index), nil
}

func (littleEndian) PutUint64Array(b []byte, index int, value []uint64) (uint64, error) {
	var reqsize = (len(value)*8 + index)
	if reqsize > len(b) || index < 0 {
		return 0, ErrOutOfRange
	}
	for i := 0; i < len(value); i++ {
		binary.LittleEndian.PutUint64(b[index+i*8:], value[i])
	}
	return uint64(reqsize - index), nil
}

func (littleEndian) PutInt16Array(b []byte, index int, value []int16) (uint64, error) {
	var reqsize = (len(value)*2 + index)
	if reqsize > len(b) || index < 0 {
		return 0, ErrOutOfRange
	}
	for i := 0; i < len(value); i++ {
		binary.LittleEndian.PutUint16(b[index+i*2:], uint16(value[i]))
	}
	return uint64(len(value) * 2), nil
}

func (littleEndian) PutInt32Array(b []byte, index int, value []int32) (uint64, error) {
	var reqsize = (len(value)*4 + index)
	if reqsize > len(b) || index < 0 {
		return 0, ErrOutOfRange
	}
	for i := 0; i < len(value); i++ {
		binary.LittleEndian.PutUint32(b[index+i*4:], uint32(value[i]))
	}
	return uint64(reqsize - index), nil
}

func (littleEndian) PutInt64Array(b []byte, index int, value []int64) (uint64, error) {
	var reqsize = (len(value)*8 + index)
	if reqsize > len(b) || index < 0 {
		return 0, ErrOutOfRange
	}
	for i := 0; i < len(value); i++ {
		binary.LittleEndian.PutUint64(b[index+i*8:], uint64(value[i]))
	}
	return uint64(reqsize - index), nil
}

func (littleEndian) PutFloat32Array(b []byte, index int, value []float32) (uint64, error) {
	var reqsize = (len(value)*4 + index)
	if reqsize > len(b) || index < 0 {
		return 0, ErrOutOfRange
	}
	for i := 0; i < len(value); i++ {
		binary.LittleEndian.PutUint32(b[index+i*4:], math.Float32bits(value[i]))
	}
	return uint64(reqsize - index), nil
}

func (littleEndian) PutFloat64Array(b []byte, index int, value []float64) (uint64, error) {
	var reqsize = (len(value)*8 + index)
	if reqsize > len(b) || index < 0 {
		return 0, ErrOutOfRange
	}
	for i := 0; i < len(value); i++ {
		binary.LittleEndian.PutUint64(b[index+i*8:], math.Float64bits(value[i]))
	}
	return uint64(reqsize - index), nil
}

func (littleEndian) String(b []byte, index, size int) (string, error) {
	if index+size > len(b) || index < 0 {
		return "", ErrOutOfRange
	}
	return string(b[index:]), nil
}

func (littleEndian) PutString(b []byte, index int, value string) (uint64, error) {
	if index+len(value) > len(b) || index < 0 {
		return 0, ErrOutOfRange
	}
	copy(b[index:index+len(value)], []byte(value))
	return uint64(len(value)), nil
}

func (littleEndian) PutUint8Array(b []byte, index int, value []uint8) (uint64, error) {
	var reqsize = (len(value) + index)
	if reqsize > len(b) || index < 0 {
		return 0, ErrOutOfRange
	}
	for i := 0; i < len(value); i++ {
		b[index+i] = byte(value[i])
	}
	return uint64(len(value)), nil
}

func (littleEndian) PutInt8Array(b []byte, index int, value []int8) (uint64, error) {
	var reqsize = (len(value) + index)
	if reqsize > len(b) || index < 0 {
		return 0, ErrOutOfRange
	}
	for i := 0; i < len(value); i++ {
		b[index+i] = byte(value[i])
	}
	return uint64(len(value)), nil
}
