package namedtuple

import (
	// "fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// func createUintTestType() TupleType {

// 	// unsigned integer test type
// 	UintTestType := New("uint")

// 	// Integers
// 	UintTestType.AddVersion(
// 		Field{"uint8-8", true, Uint8Field},
// 		Field{"uint16-8", true, Uint16Field},
// 		Field{"uint16-16", true, Uint16Field},
// 		Field{"uint32-8", true, Uint32Field},
// 		Field{"uint32-16", true, Uint32Field},
// 		Field{"uint32-32", true, Uint32Field},
// 		Field{"uint64-8", true, Uint64Field},
// 		Field{"uint64-16", true, Uint64Field},
// 		Field{"uint64-32", true, Uint64Field},
// 		Field{"uint64-64", true, Uint64Field},
// 	)
// 	// Arrays
// 	UintTestType.AddVersion(
// 		Field{"uint8-8-array", true, Uint8FieldArray},
// 		Field{"uint16-8-array", true, Uint16FieldArray},
// 		Field{"uint16-1-array6", true, Uint16FieldArray},
// 		Field{"uint32-8-array", true, Uint32FieldArray},
// 		Field{"uint32-1-array6", true, Uint32FieldArray},
// 		Field{"uint32-3-array2", true, Uint32FieldArray},
// 		Field{"uint64-8-array", true, Uint64FieldArray},
// 		Field{"uint64-1-array6", true, Uint64FieldArray},
// 		Field{"uint64-3-array2", true, Uint64FieldArray},
// 		Field{"uint64-6-array4", true, Uint64FieldArray},
// 	)
// 	return UintTestType
// }

func TestNewBuilder(t *testing.T) {

	// create test type
	User := createTestTupleType()

	// create builder
	buffer := make([]byte, 1024)
	builder := NewBuilder(User, buffer)

	// verify type
	assert.Equal(t, User, builder.tupleType)

	// verify type fields
	assert.Equal(t, len(User.fields), len(builder.fields))
	for name := range builder.fields {

		// make sure the type has the same fields as the builder
		assert.True(t, User.Contains(name))
	}
}

func TestBuilderAvailableEmpty(t *testing.T) {
	// create test type
	User := createTestTupleType()

	// create builder
	buffer := make([]byte, 1024)
	builder := NewBuilder(User, buffer)

	// verify available == 1024
	assert.Equal(t, 1024, builder.available())

}

func TestBuilderTypeCheck(t *testing.T) {
	// create test type
	User := createTestTupleType()

	// create builder
	buffer := make([]byte, 1024)
	builder := NewBuilder(User, buffer)

	// testing correct fields
	assert.Nil(t, builder.typeCheck("uuid", StringField))
	assert.Nil(t, builder.typeCheck("username", StringField))
	assert.Nil(t, builder.typeCheck("age", Uint8Field))
	assert.Nil(t, builder.typeCheck("location", TupleField))

	// testing invalid field
	assert.NotNil(t, builder.typeCheck("school", StringField))

	// testing invalid type
	assert.NotNil(t, builder.typeCheck("uuid", TimestampField))
}

// building
func TestBuild(t *testing.T) {

	// type
	User := createTestTupleType()

	// create builder
	buffer := make([]byte, 1024)
	builder := NewBuilder(User, buffer)

	// fields
	builder.PutString("username", "value")
	builder.PutString("uuid", "value")
	builder.PutUint8("age", 25)

	// tuple
	user, err := builder.Build()
	assert.Nil(t, err)
	assert.NotNil(t, user)

	// check data length
	assert.Equal(t, 7+7+2, len(user.data), "Length of user tuple")
}

func TestTupleTypeNew(t *testing.T) {

	User := New("testing", "User")

	assert.Equal(t, "User", User.Name)
	assert.Equal(t, 0, User.NumVersions())

	hash := syncHash.Hash([]byte("User"))
	assert.Equal(t, hash, User.Hash)
}

func TestTupleTypeAddVersion(t *testing.T) {

	// fields
	uuid := Field{"uuid", true, StringField}
	username := Field{"username", true, StringField}
	age := Field{"age", false, Uint8Field}
	location := Field{"location", false, TupleField}

	// create tuple type
	User := New("testing", "user")
	User.AddVersion(uuid, username, age)
	User.AddVersion(location)

	// verify versions were added
	vs := User.Versions()
	assert.Equal(t, 2, User.NumVersions())
	assert.Equal(t, len(vs), User.NumVersions())

	// verify fields
	// version 1
	assert.Equal(t, 1, int(vs[0].Num))
	assert.Equal(t, uuid, vs[0].Fields[0])
	assert.Equal(t, username, vs[0].Fields[1])
	assert.Equal(t, age, vs[0].Fields[2])

	// version 2
	assert.Equal(t, 2, int(vs[1].Num))
	assert.Equal(t, location, vs[1].Fields[0])
}

func TestTupleTypeFieldOffset(t *testing.T) {

	// fields
	uuid := Field{"uuid", true, StringField}
	username := Field{"username", true, StringField}
	age := Field{"age", false, Uint8Field}
	location := Field{"location", false, TupleField}

	// create tuple type
	User := New("testing", "user")
	User.AddVersion(uuid, username, age)
	User.AddVersion(location)

	// uuid field
	offset, exists := User.Offset("uuid")
	assert.Equal(t, 0, offset)
	assert.Equal(t, true, exists)

	// username field
	offset, exists = User.Offset("username")
	assert.Equal(t, 1, offset)
	assert.Equal(t, true, exists)

	// age field
	offset, exists = User.Offset("age")
	assert.Equal(t, 2, offset)
	assert.Equal(t, true, exists)

	// location field
	offset, exists = User.Offset("location")
	assert.Equal(t, 3, offset)
	assert.Equal(t, true, exists)

	// bad field
	offset, exists = User.Offset("bad")
	assert.Equal(t, 0, offset)
	assert.Equal(t, false, exists)
}
