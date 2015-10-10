package namedtuple

import (
	"bytes"
	"testing"

	"github.com/blacklabeldata/xbinary"
	"github.com/stretchr/testify/assert"
)

func TestProtocolVersionMask(t *testing.T) {
	assert.Equal(t, 63, ProtocolVersionMask, "ProtocolVersionMask must be 63")
}

func TestProtocolSizeEnumMask(t *testing.T) {
	assert.Equal(t, 192, ProtocolSizeEnumMask, "ProtocolSizeEnumMask must be 192")
}

func TestParseProtocolHeaderByteCount(t *testing.T) {
	var header uint8 = 192

	byteCount, version := ParseProtocolHeader(header)
	assert.Equal(t, 0, int(version), "Version should be 0")
	assert.Equal(t, 8, int(byteCount), "Byte count should be 8")

	header = 128
	byteCount, version = ParseProtocolHeader(header)
	assert.Equal(t, 0, int(version), "Version should be 0")
	assert.Equal(t, 4, int(byteCount), "Byte count should be 4")

	header = 64
	byteCount, version = ParseProtocolHeader(header)
	assert.Equal(t, 0, int(version), "Version should be 0")
	assert.Equal(t, 2, int(byteCount), "Byte count should be 2")

	header = 0
	byteCount, version = ParseProtocolHeader(header)
	assert.Equal(t, 0, int(version), "Version should be 0")
	assert.Equal(t, 1, int(byteCount), "Byte count should be 1")
}

func TestParseProtocolHeaderVersion(t *testing.T) {
	var header uint8 = 255

	byteCount, version := ParseProtocolHeader(header)
	assert.Equal(t, 63, int(version), "Version should be 63")
	assert.Equal(t, 8, int(byteCount), "Byte count should be 0")

	header = 32
	byteCount, version = ParseProtocolHeader(header)
	assert.Equal(t, 32, int(version), "Version should be 32")
	assert.Equal(t, 1, int(byteCount), "Byte count should be 1")

	header = 1
	byteCount, version = ParseProtocolHeader(header)
	assert.Equal(t, 1, int(version), "Version should be 1")
	assert.Equal(t, 1, int(byteCount), "Byte count should be 1")

	header = 0
	byteCount, version = ParseProtocolHeader(header)
	assert.Equal(t, 0, int(version), "Version should be 0")
	assert.Equal(t, 1, int(byteCount), "Byte count should be 1")
}

func TestTupleHeaderSize(t *testing.T) {
	header := TupleHeader{}
	header.FieldSize = 1
	header.FieldCount = 4
	assert.Equal(t, VersionOneTupleHeaderSize+1*4, header.Size(), "TupleHeader should be 17")
}

func createTestUserType() TupleType {
	// fields
	uuid := Field{"uuid", true, StringField}
	username := Field{"username", true, StringField}
	age := Field{"age", false, Uint8Field}

	// create tuple type
	User := New("testing", "user")
	User.AddVersion(uuid, username, age)
	return User
}

func TestTupleHeaderWriteToSmall(t *testing.T) {
	User := createTestUserType()
	header := TupleHeader{
		ProtocolVersion: 0,
		TupleVersion:    1,
		NamespaceHash:   User.NamespaceHash,
		Hash:            User.Hash,
		FieldCount:      3,
		FieldSize:       1,
		ContentLength:   64,
		Offsets:         []uint64{1, 2, 3},
		Type:            User,
	}

	var buf []byte
	writer := bytes.NewBuffer(buf)
	_, err := header.WriteTo(writer)
	assert.Nil(t, err)

	buffer := writer.Bytes()

	// Check length
	assert.Equal(t, 16, len(buffer))

	// check tuple version
	byteCount, version := ParseProtocolHeader(buffer[0])
	assert.Equal(t, uint8(1), byteCount, "Byte count should be 1")
	assert.Equal(t, uint8(1), version, "Version should be 1")

	// Read hashes an
	dst := make([]uint32, 3)
	err = xbinary.LittleEndian.Uint32Array(buffer, 1, &dst)

	// Check hashes
	assert.Equal(t, header.NamespaceHash, dst[0])
	assert.Equal(t, header.Hash, dst[1])

	// Check field count
	assert.Equal(t, header.FieldCount, dst[2])

	// Check offsets
	assert.Equal(t, uint8(1), buffer[13])
	assert.Equal(t, uint8(2), buffer[14])
	assert.Equal(t, uint8(3), buffer[15])
}

func TestTupleHeaderWriteToMedium(t *testing.T) {
	User := createTestUserType()
	header := TupleHeader{
		ProtocolVersion: 0,
		TupleVersion:    1,
		NamespaceHash:   User.NamespaceHash,
		Hash:            User.Hash,
		FieldCount:      3,
		FieldSize:       2,
		ContentLength:   64,
		Offsets:         []uint64{1, 2, 3},
		Type:            User,
	}

	var buf []byte
	writer := bytes.NewBuffer(buf)
	_, err := header.WriteTo(writer)
	assert.Nil(t, err)

	buffer := writer.Bytes()

	// Check length
	assert.Equal(t, 19, len(buffer))

	// check tuple version
	byteCount, version := ParseProtocolHeader(buffer[0])
	assert.Equal(t, uint8(2), byteCount, "Byte count should be 1")
	assert.Equal(t, uint8(1), version, "Version should be 1")

	// Read hashes an
	dst := make([]uint32, 3)
	err = xbinary.LittleEndian.Uint32Array(buffer, 1, &dst)

	// Check hashes
	assert.Equal(t, header.NamespaceHash, dst[0])
	assert.Equal(t, header.Hash, dst[1])

	// Check field count
	assert.Equal(t, header.FieldCount, dst[2])

	// Check offsets
	offsets := make([]uint16, 3)
	err = xbinary.LittleEndian.Uint16Array(buffer, 13, &offsets)
	assert.Nil(t, err)
	assert.Equal(t, uint16(1), offsets[0])
	assert.Equal(t, uint16(2), offsets[1])
	assert.Equal(t, uint16(3), offsets[2])
}

func TestTupleHeaderWriteToLarge(t *testing.T) {
	User := createTestUserType()
	header := TupleHeader{
		ProtocolVersion: 0,
		TupleVersion:    1,
		NamespaceHash:   User.NamespaceHash,
		Hash:            User.Hash,
		FieldCount:      3,
		FieldSize:       4,
		ContentLength:   64,
		Offsets:         []uint64{1, 2, 3},
		Type:            User,
	}

	var buf []byte
	writer := bytes.NewBuffer(buf)
	_, err := header.WriteTo(writer)
	assert.Nil(t, err)

	buffer := writer.Bytes()

	// Check length
	assert.Equal(t, 25, len(buffer))

	// check tuple version
	byteCount, version := ParseProtocolHeader(buffer[0])
	assert.Equal(t, uint8(4), byteCount, "Byte count should be 1")
	assert.Equal(t, uint8(1), version, "Version should be 1")

	// Read hashes an
	dst := make([]uint32, 3)
	err = xbinary.LittleEndian.Uint32Array(buffer, 1, &dst)

	// Check hashes
	assert.Equal(t, header.NamespaceHash, dst[0])
	assert.Equal(t, header.Hash, dst[1])

	// Check field count
	assert.Equal(t, header.FieldCount, dst[2])

	// Check offsets
	offsets := make([]uint32, 3)
	err = xbinary.LittleEndian.Uint32Array(buffer, 13, &offsets)
	assert.Nil(t, err)
	assert.Equal(t, uint32(1), offsets[0])
	assert.Equal(t, uint32(2), offsets[1])
	assert.Equal(t, uint32(3), offsets[2])
}

func TestTupleHeaderWriteToExtraLarge(t *testing.T) {
	User := createTestUserType()
	header := TupleHeader{
		ProtocolVersion: 0,
		TupleVersion:    1,
		NamespaceHash:   User.NamespaceHash,
		Hash:            User.Hash,
		FieldCount:      3,
		FieldSize:       8,
		ContentLength:   64,
		Offsets:         []uint64{1, 2, 3},
		Type:            User,
	}

	var buf []byte
	writer := bytes.NewBuffer(buf)
	_, err := header.WriteTo(writer)
	assert.Nil(t, err)

	buffer := writer.Bytes()

	// Check length
	assert.Equal(t, 37, len(buffer))

	// check tuple version
	byteCount, version := ParseProtocolHeader(buffer[0])
	assert.Equal(t, uint8(8), byteCount, "Byte count should be 1")
	assert.Equal(t, uint8(1), version, "Version should be 1")

	// Read hashes an
	dst := make([]uint32, 3)
	err = xbinary.LittleEndian.Uint32Array(buffer, 1, &dst)

	// Check hashes
	assert.Equal(t, header.NamespaceHash, dst[0])
	assert.Equal(t, header.Hash, dst[1])

	// Check field count
	assert.Equal(t, header.FieldCount, dst[2])

	// Check offsets
	offsets := make([]uint64, 3)
	err = xbinary.LittleEndian.Uint64Array(buffer, 13, &offsets)
	assert.Nil(t, err)
	assert.Equal(t, uint64(1), offsets[0])
	assert.Equal(t, uint64(2), offsets[1])
	assert.Equal(t, uint64(3), offsets[2])
}
