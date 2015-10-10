package namedtuple

import (
	"bytes"
	"io"
	"testing"

	"github.com/blacklabeldata/xbinary"
	"github.com/stretchr/testify/assert"
)

func TestNewDecoder(t *testing.T) {
	var buf []byte
	dec := NewDecoder(DefaultRegistry, bytes.NewReader(buf))
	assert.NotNil(t, dec)

	// Should be a decoder
	d, ok := dec.(decoder)
	assert.True(t, ok)
	assert.Equal(t, DefaultMaxSize, d.maxSize)
}

func TestNewDecoderSize(t *testing.T) {
	var buf []byte
	dec := NewDecoderSize(DefaultRegistry, 512, bytes.NewReader(buf))
	assert.NotNil(t, dec)

	// Should be a decoder
	d, ok := dec.(decoder)
	assert.True(t, ok)
	assert.Equal(t, uint64(512), d.maxSize)
}

func TestDecoderParseLength8(t *testing.T) {
	d := decoder{DefaultRegistry, 512, nil, nil}

	// Create buffer
	var buf = []byte{5}

	// Parse length
	length, err := d.parseLength(1, buf)
	assert.Nil(t, err)
	assert.Equal(t, uint64(5), length)
}

func TestDecoderParseLength8Fail(t *testing.T) {
	d := decoder{DefaultRegistry, 512, nil, nil}

	// Create buffer
	var buf []byte

	// Parse length
	length, err := d.parseLength(1, buf)
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), length)
}

func TestDecoderParseLength16(t *testing.T) {
	d := decoder{DefaultRegistry, 512, nil, nil}

	// Create buffer
	buf := make([]byte, 2)
	xbinary.LittleEndian.PutUint16(buf, 0, 512)

	// Parse length
	length, err := d.parseLength(2, buf)
	assert.Nil(t, err)
	assert.Equal(t, uint64(512), length)
}

func TestDecoderParseLength16Fail(t *testing.T) {
	d := decoder{DefaultRegistry, 512, nil, nil}

	// Create buffer
	var buf []byte

	// Parse length
	length, err := d.parseLength(2, buf)
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), length)
}

func TestDecoderParseLength32(t *testing.T) {
	d := decoder{DefaultRegistry, 512, nil, nil}

	// Create buffer
	buf := make([]byte, 4)
	xbinary.LittleEndian.PutUint32(buf, 0, 512*1024)

	// Parse length
	length, err := d.parseLength(4, buf)
	assert.Nil(t, err)
	assert.Equal(t, uint64(512*1024), length)
}

func TestDecoderParseLength32Fail(t *testing.T) {
	d := decoder{DefaultRegistry, 512, nil, nil}

	// Create buffer
	var buf []byte

	// Parse length
	length, err := d.parseLength(4, buf)
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), length)
}

func TestDecoderParseLength64(t *testing.T) {
	d := decoder{DefaultRegistry, 512, nil, nil}

	// Create buffer
	buf := make([]byte, 8)
	xbinary.LittleEndian.PutUint64(buf, 0, 512*10e9)

	// Parse length
	length, err := d.parseLength(8, buf)
	assert.Nil(t, err)
	assert.Equal(t, uint64(512*10e9), length)
}

func TestDecoderParseLength64Fail(t *testing.T) {
	d := decoder{DefaultRegistry, 512, nil, nil}

	// Create buffer
	var buf []byte

	// Parse length
	length, err := d.parseLength(8, buf)
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), length)
}

func TestDecoderParseLengthFail(t *testing.T) {
	d := decoder{DefaultRegistry, 512, nil, nil}

	// Create buffer
	var buf []byte

	// Parse length
	length, err := d.parseLength(0, buf)
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), length)
}

func TestDecoderDecodeFailReadByte(t *testing.T) {
	var buf []byte
	dec := NewDecoderSize(DefaultRegistry, 512, bytes.NewReader(buf))
	tup, err := dec.Decode()
	assert.NotNil(t, err)
	assert.Equal(t, EmptyTuple, tup)
}

func TestDecoderDecodeFailReadLength64(t *testing.T) {
	var buf = []byte{192}
	dec := NewDecoderSize(DefaultRegistry, 512, bytes.NewReader(buf))
	tup, err := dec.Decode()
	assert.NotNil(t, err)
	assert.Equal(t, EmptyTuple, tup)
}

func TestDecoderDecodeFailReadLength32(t *testing.T) {
	var buf = []byte{128}
	dec := NewDecoderSize(DefaultRegistry, 512, bytes.NewReader(buf))
	tup, err := dec.Decode()
	assert.NotNil(t, err)
	assert.Equal(t, EmptyTuple, tup)
}

func TestDecoderDecodeFailReadLength16(t *testing.T) {
	var buf = []byte{64}
	dec := NewDecoderSize(DefaultRegistry, 512, bytes.NewReader(buf))
	tup, err := dec.Decode()
	assert.NotNil(t, err)
	assert.Equal(t, EmptyTuple, tup)
}

func TestDecoderDecodeFailReadLength8(t *testing.T) {
	var buf = []byte{0}
	dec := NewDecoderSize(DefaultRegistry, 512, bytes.NewReader(buf))
	tup, err := dec.Decode()
	assert.NotNil(t, err)
	assert.Equal(t, EmptyTuple, tup)
}

func TestDecoderDecodeFailParseLength(t *testing.T) {
	// A correct byte buffer would include 5 bytes
	// By setting the header to 128, we're telling the
	// decoder to expect 4 bytes to follow.The decoder should
	// return an error as there are not enought bytes.
	buf := make([]byte, 4)
	buf[0] = 128

	dec := NewDecoderSize(DefaultRegistry, 512, bytes.NewReader(buf))
	tup, err := dec.Decode()
	assert.NotNil(t, err)
	assert.Equal(t, EmptyTuple, tup)
}

func TestDecoderDecodeFailMaxLength(t *testing.T) {
	buf := make([]byte, 5)
	buf[0] = 128
	xbinary.LittleEndian.PutUint32(buf, 1, 1024)

	dec := NewDecoderSize(DefaultRegistry, 512, bytes.NewReader(buf))
	tup, err := dec.Decode()
	assert.NotNil(t, err)
	assert.Equal(t, ErrTupleExceedsMaxSize, err)
	assert.Equal(t, EmptyTuple, tup)
}

func TestDecoderDecodeFailEOF(t *testing.T) {
	buf := make([]byte, 5)
	buf[0] = 128
	buf[0] |= 1
	xbinary.LittleEndian.PutUint32(buf, 1, 256)

	dec := NewDecoderSize(DefaultRegistry, 512, bytes.NewReader(buf))
	tup, err := dec.Decode()
	assert.NotNil(t, err)
	assert.Equal(t, io.EOF, err)
	assert.Equal(t, EmptyTuple, tup)
}

func TestDecoderDecodeFailUnknownVersion(t *testing.T) {
	// 32 bytes include header (2 bytes) and content length (30 bytes)
	buf := make([]byte, 32)
	buf[0] = 0
	buf[1] = 30

	dec := NewDecoderSize(DefaultRegistry, 512, bytes.NewReader(buf))
	tup, err := dec.Decode()
	assert.NotNil(t, err)
	assert.Equal(t, ErrInvalidProtocolVersion, err)
	assert.Equal(t, EmptyTuple, tup)
}

func TestDecoder_ParseVersionOneTuple(t *testing.T) {
	var buf []byte
	dec := NewDecoder(DefaultRegistry, bytes.NewReader(buf))
	assert.NotNil(t, dec)

	// Should be a decoder
	d, ok := dec.(decoder)
	assert.True(t, ok)

	// Parse tuple
	tup, err := d.parseVersionOneTuple(0, 0, 0)
	assert.Equal(t, EmptyTuple, tup)
	assert.NotNil(t, err)
	assert.Equal(t, ErrTupleLengthTooSmall, err)
}

func TestDecoderParseVersionOneTupleFailUnknownType(t *testing.T) {
	// 32 bytes include header (2 bytes) and content length (30 bytes)
	buf := make([]byte, 32)
	buf[0] = 1
	buf[1] = 30

	dec := NewDecoder(DefaultRegistry, bytes.NewReader(buf))
	assert.NotNil(t, dec)

	// Parse tuple
	tup, err := dec.Decode()
	assert.Equal(t, EmptyTuple, tup)
	assert.NotNil(t, err)
	assert.Equal(t, ErrUnknownTupleType, err)
}

func TestReadFieldOffsets8(t *testing.T) {
	buf := make([]byte, 13+5)
	buf[13], buf[14], buf[15], buf[16], buf[17] = 1, 2, 3, 4, 5

	// Read offsets
	offsets, err := readFieldOffsets(1, 5, buf)
	assert.Nil(t, err)
	assert.Equal(t, []uint64{1, 2, 3, 4, 5}, offsets)
}

func TestReadFieldOffsets8Fail(t *testing.T) {
	buf := make([]byte, 13)

	// Read offsets
	_, err := readFieldOffsets(1, 5, buf)
	assert.NotNil(t, err)
	assert.Equal(t, ErrTupleLengthTooSmall, err)
}

func TestReadFieldOffsets16(t *testing.T) {
	buf := make([]byte, 13+6)
	xbinary.LittleEndian.PutUint16(buf, 13, 1)
	xbinary.LittleEndian.PutUint16(buf, 15, 2)
	xbinary.LittleEndian.PutUint16(buf, 17, 3)

	// Read offsets
	offsets, err := readFieldOffsets(2, 3, buf)
	assert.Nil(t, err)
	assert.Equal(t, []uint64{1, 2, 3}, offsets)
}

func TestReadFieldOffsets16Fail(t *testing.T) {
	var buf []byte

	// Read offsets
	_, err := readFieldOffsets(2, 5, buf)
	assert.NotNil(t, err)
	assert.Equal(t, xbinary.ErrOutOfRange, err)
}

func TestReadFieldOffsets32(t *testing.T) {
	buf := make([]byte, 13+12)
	xbinary.LittleEndian.PutUint32(buf, 13, 1)
	xbinary.LittleEndian.PutUint32(buf, 17, 2)
	xbinary.LittleEndian.PutUint32(buf, 21, 3)

	// Read offsets
	offsets, err := readFieldOffsets(4, 3, buf)
	assert.Nil(t, err)
	assert.Equal(t, []uint64{1, 2, 3}, offsets)
}

func TestReadFieldOffsets32Fail(t *testing.T) {
	var buf []byte

	// Read offsets
	_, err := readFieldOffsets(4, 5, buf)
	assert.NotNil(t, err)
	assert.Equal(t, xbinary.ErrOutOfRange, err)
}

func TestReadFieldOffsets64(t *testing.T) {
	buf := make([]byte, 13+24)
	xbinary.LittleEndian.PutUint64(buf, 13, 1)
	xbinary.LittleEndian.PutUint64(buf, 21, 2)
	xbinary.LittleEndian.PutUint64(buf, 29, 3)

	// Read offsets
	offsets, err := readFieldOffsets(8, 3, buf)
	assert.Nil(t, err)
	assert.Equal(t, []uint64{1, 2, 3}, offsets)
}

func TestReadFieldOffsets64Fail(t *testing.T) {
	var buf []byte

	// Read offsets
	_, err := readFieldOffsets(8, 5, buf)
	assert.NotNil(t, err)
	assert.Equal(t, xbinary.ErrOutOfRange, err)
}

func TestReadFieldOffsetsFail(t *testing.T) {
	var buf []byte

	// Read offsets
	_, err := readFieldOffsets(0, 5, buf)
	assert.NotNil(t, err)
	assert.Equal(t, ErrInvalidLength, err)
}

func TestDecode(t *testing.T) {
	// Create encoder
	var buf []byte
	out := bytes.NewBuffer(buf)
	encoder := NewEncoder(out)

	// Create message type
	msgBuffer := make([]byte, 256)
	Message := createTestMessageType()
	t.Logf("Namespace Hash: %x", Message.NamespaceHash)
	t.Logf("Type Hash: %x", Message.Hash)

	// Create location type
	locBuffer := make([]byte, 256)
	Location := createTestLocationType()

	// Create location builder
	locBuilder := Location.Builder(locBuffer)
	locBuilder.PutFloat32("lon", 150.5)
	locBuilder.PutFloat32("lat", 50.5)
	locBuilder.PutFloat32("alt", 9022)
	loc, err := locBuilder.Build()
	assert.Nil(t, err, "Error should be nil")

	// Create message builder
	msgBuilder := Message.Builder(msgBuffer)
	msgBuilder.PutString("payload", "Vacation in Miami, FL")
	msgBuilder.PutString("userid", "eliquious")
	_, err = msgBuilder.PutTuple("loc", loc)
	assert.Nil(t, err, "Error should be nil")

	// Build message
	msg, err := msgBuilder.Build()
	assert.Nil(t, err, "Error should be nil")

	// Encode message
	err = encoder.Encode(msg)
	assert.Nil(t, err, "Error should be nil")

	// b := out.Bytes()
	// t.Logf("Output: %d - %#v", len(b), b)

	// Create Registry
	reg := NewRegistry()
	reg.Register(Location)
	reg.Register(Message)
	// t.Logf("Registry: ", reg.content)

	// Create decoder
	dec := NewDecoder(reg, bytes.NewReader(out.Bytes()))
	message, err := dec.Decode()
	assert.Nil(t, err)
	assert.True(t, bytes.Equal(msg.data, message.data))
	assert.Equal(t, msg.Header, message.Header)
}
