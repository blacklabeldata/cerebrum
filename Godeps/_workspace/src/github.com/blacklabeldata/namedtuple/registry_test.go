package namedtuple

import (
	// "fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	// "time"
)

func createTestTupleType() TupleType {
	// fields
	uuid := Field{"uuid", true, StringField}
	username := Field{"username", true, StringField}
	age := Field{"age", false, Uint8Field}
	location := Field{"location", false, TupleField}

	lat := Field{"lat", false, Float32Field}
	lon := Field{"lon", false, Float32Field}
	alt := Field{"alt", false, Float64Field}

	// create tuple type
	User := New("testing", "user")
	User.AddVersion(uuid, username, age)
	User.AddVersion(location)
	User.AddVersion(lat, lon, alt)
	return User
}

func TestRegistry(t *testing.T) {

	// create new empty registry
	reg := NewRegistry()

	// make sure it's empty
	assert.Equal(t, 0, reg.Size())
	assert.Equal(t, len(reg.content), reg.Size())
}

func TestRegistryRegister(t *testing.T) {

	// create new empty registry
	reg := NewRegistry()

	// create type
	User := createTestTupleType()

	// add User type
	reg.Register(User)

	// make sure it's not empty
	assert.Equal(t, 1, reg.Size())
	assert.Equal(t, len(reg.content), reg.Size())
}

func TestRegistryUnregister(t *testing.T) {

	// create new empty registry
	reg := NewRegistry()

	// create type
	User := createTestTupleType()

	// add User type
	reg.Register(User)

	// make sure it's not empty
	assert.Equal(t, 1, reg.Size())
	assert.Equal(t, len(reg.content), reg.Size())

	// remove User type
	reg.Unregister(User)

	// make sure it's not empty
	assert.Equal(t, 0, reg.Size())
	assert.Equal(t, len(reg.content), reg.Size())
}

func TestRegistryContainsTrue(t *testing.T) {

	// create new empty registry
	reg := NewRegistry()

	// create type
	User := createTestTupleType()

	// add User type
	reg.Register(User)
	hash := reg.typeSignature(User.Namespace, User.Name)

	// make sure it contains the User type
	assert.Equal(t, User, reg.content[hash])

	// test contains function
	assert.Equal(t, true, reg.Contains(User))

	// test contains hash function
	assert.Equal(t, true, reg.ContainsHash(hash))

	// test contains name function
	assert.Equal(t, true, reg.ContainsName(User.Namespace, User.Name))
}

func TestRegistryContainsFalse(t *testing.T) {

	// create new empty registry
	reg := NewRegistry()

	// create type
	User := createTestTupleType()
	hash := reg.typeSignature(User.Namespace, User.Name)

	// DO NOT add User type
	// reg.Register(User)

	// make sure it DOES NOT contains the User type
	assert.Equal(t, TupleType{}, reg.content[hash])

	// test contains function
	assert.Equal(t, false, reg.Contains(User))

	// test contains hash function
	assert.Equal(t, false, reg.ContainsHash(hash))

	// test contains name function
	assert.Equal(t, false, reg.ContainsName(User.Namespace, User.Name))
}

func TestRegistryGetTrue(t *testing.T) {

	// create new empty registry
	reg := NewRegistry()

	// create type
	User := createTestTupleType()

	// add User type
	reg.Register(User)

	// make sure the registry contains the same User type
	tupleType, exists := reg.Get(User.Namespace, User.Name)
	assert.Equal(t, User, tupleType)
	assert.Equal(t, true, exists)
}

func TestRegistryGetFalse(t *testing.T) {

	// create new empty registry
	reg := NewRegistry()

	// create type
	User := createTestTupleType()

	// DO NOT add User type
	// reg.Register(User)

	// make sure the registry contains the same User type
	tupleType, exists := reg.Get(User.Namespace, User.Name)
	assert.Equal(t, TupleType{}, tupleType)
	assert.Equal(t, false, exists)
}

func TestRegistryGetWithHash(t *testing.T) {

	// create new empty registry
	reg := NewRegistry()

	// create type
	User := createTestTupleType()
	reg.Register(User)

	// make sure the registry contains the same User type
	tupleType, exists := reg.GetWithHash(User.NamespaceHash, User.Hash)
	assert.Equal(t, User, tupleType)
	assert.Equal(t, true, exists)
}

func TestRegistryGetWithHashFalse(t *testing.T) {

	// create new empty registry
	reg := NewRegistry()

	// create type
	User := createTestTupleType()

	// DO NOT add User type
	// reg.Register(User)

	// make sure the registry contains the same User type
	tupleType, exists := reg.GetWithHash(User.NamespaceHash, User.Hash)
	assert.Equal(t, TupleType{}, tupleType)
	assert.Equal(t, false, exists)
}
