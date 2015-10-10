package namedtuple

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	// "github.com/blacklabeldata/xbinary"
)

func createTestLocationType() TupleType {
	// fields
	lon := Field{"lon", true, Float32Field}
	lat := Field{"lat", true, Float32Field}
	alt := Field{"alt", false, Float32Field}

	// create tuple type
	Location := New("testing", "Location")
	Location.AddVersion(lon, lat, alt)
	return Location
}

func createTestMessageType() TupleType {
	// fields
	userid := Field{"userid", true, StringField}
	payload := Field{"payload", true, StringField}
	loc := Field{"loc", false, TupleField}

	// create tuple type
	Message := New("testing", "Message")
	Message.AddVersion(userid, payload, loc)
	return Message
}

func TestEncode(t *testing.T) {
	var buf []byte
	out := bytes.NewBuffer(buf)
	encoder := NewEncoder(out)

	msgBuffer := make([]byte, 256)
	Message := createTestMessageType()

	locBuffer := make([]byte, 256)
	Location := createTestLocationType()

	locBuilder := Location.Builder(locBuffer)
	locBuilder.PutFloat32("lon", 150.5)
	locBuilder.PutFloat32("lat", 50.5)
	locBuilder.PutFloat32("alt", 9022)
	loc, err := locBuilder.Build()
	assert.Nil(t, err, "Error should be nil")
	// t.Logf("Tuple : %#v", loc)
	// t.Logf("Tuple Size : %d", loc.Size())
	// t.Logf("Tuple Header Size : %d", loc.Header.Size())

	msgBuilder := Message.Builder(msgBuffer)
	msgBuilder.PutString("payload", "Vacation in Miami, FL")
	msgBuilder.PutString("userid", "eliquious")
	_, err = msgBuilder.PutTuple("loc", loc)
	assert.Nil(t, err, "Error should be nil")

	msg, err := msgBuilder.Build()
	assert.Nil(t, err, "Error should be nil")

	// t.Logf("Tuple : %#v", msg)
	// t.Logf("Tuple Size : %d", msg.Size())
	// t.Logf("Tuple Header Size : %d", msg.Header.Size())

	err = encoder.Encode(msg)
	assert.Nil(t, err, "Error should be nil")

	// b := out.Bytes()
	// t.Logf("Output: %d", len(b), b)
}
