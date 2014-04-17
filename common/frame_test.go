package common

import (
	"encoding/binary"
	"testing"
	"github.com/stretchr/testify/assert"
)


type stringWrapper struct {
	data string
}

// This is named wrong
func (s stringWrapper) Marshal() []byte {
	return []byte(s.data)
}

func testFrame() Frame {
	f := Frame{}
	f.Payload = stringWrapper{}

	return f
}

func TestMarshalEmptyFrame(t *testing.T) {
	f := testFrame()

	marshalled_f := f.Marshal()

	length := binary.BigEndian.Uint16(marshalled_f[0:2]) & 0x7F

	assert.Equal(t, length, uint16(0),
		"Length should have been nothing")
}

func TestMarshalFrameWithPayloadIncludesLength(t *testing.T) {
	f := testFrame()
	f.Payload = stringWrapper{ "this is the payload of the frame" }

	marshalled_f := f.Marshal()

	length := binary.BigEndian.Uint16(marshalled_f[0:2]) & 0x7F

	assert.Equal(t, int(length), len("this is the payload of the frame"),
		"Length field in header should have been the length of payload")
}

func TestMarshalFrameWithType(t *testing.T) {
	f := testFrame()
	f.Type = byte(8)

	assert.Equal(t, byte(8), f.Marshal()[2],
		"Type should have been marshalled as the third byte")
}

func TestMarshalFrameWithFlags(t *testing.T) {
	f := testFrame()
	f.Flags = byte(0xD)

	assert.Equal(t, byte(0xD), f.Marshal()[3],
		"Flags should have been marshalled as the fourth byte")
}

func TestMarshalFrameWithStreamIdentifier(t *testing.T) {
	f := testFrame()
	f.StreamIdentifier = 168036609

	marshalled_f := f.Marshal()

	stream_identifier := binary.BigEndian.Uint32(marshalled_f[4:8]) & uint32(0x7FFFFFFF)

	assert.Equal(t, stream_identifier, f.StreamIdentifier,
		"Stream identifier should have been marshalled as fifth through eighth bytes")
}
