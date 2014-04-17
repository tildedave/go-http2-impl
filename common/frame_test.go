package common

import (
	"encoding/binary"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestMarshalEmptyFrame(t *testing.T) {
	f := Frame{ 0, 0, 0, ""}
	marshalled_f := Marshal(f)

	length := binary.BigEndian.Uint16(marshalled_f[0:2]) & 0x7F

	assert.Equal(t, length, uint16(0),
		"Length should have been nothing")
}

func TestMarshalFrameWithPayloadIncludesLength(t *testing.T) {
	f := Frame{ 0, 0, 0, "this is the payload of the frame"}
	marshalled_f := Marshal(f)

	length := binary.BigEndian.Uint16(marshalled_f[0:2]) & 0x7F

	assert.Equal(t, int(length), len(f.Payload),
		"Length field in header should have been the length of payload")
}

func TestMarshalFrameWithType(t *testing.T) {
	f := Frame{}
	f.Type = byte(8)

	assert.Equal(t, byte(8), Marshal(f)[2],
		"Type should have been marshalled as the third byte")
}

func TestMarshalFrameWithFlags(t *testing.T) {
	f := Frame{}
	f.Flags = byte(0xD)

	assert.Equal(t, byte(0xD), Marshal(f)[3],
		"Flags should have been marshalled as the fourth byte")
}

func TestMarshalFrameWithStreamIdentifier(t *testing.T) {
	f := Frame{}
	f.StreamIdentifier = 168036609
	marshalled_f := Marshal(f)

	stream_identifier := binary.BigEndian.Uint32(marshalled_f[4:8]) & uint32(0x7FFFFFFF)

	assert.Equal(t, stream_identifier, f.StreamIdentifier,
		"Stream identifier should have been marshalled as fifth through eighth bytes")
}
