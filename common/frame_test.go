package common

import (
	"encoding/binary"
	"testing"
	"github.com/stretchr/testify/assert"
)


type stringWrapper struct {
	data string
}

func lengthOfMarshalledFrame(marshalled_f []byte) uint16 {
	return binary.BigEndian.Uint16(marshalled_f[0:2]) & 0x7F
}

func TestMarshalEmptyFrame(t *testing.T) {
	f := baseFrame{}

	assert.Equal(t, lengthOfMarshalledFrame(f.Marshal()), uint16(0),
		"Length should have been nothing")
}

func TestMarshalFrameWithPayloadIncludesLength(t *testing.T) {
	f := baseFrame{}
	f.Payload = "this is the payload of the frame"

	marshalled_f := f.Marshal()

	assert.Equal(t, lengthOfMarshalledFrame(marshalled_f),
		uint16(len("this is the payload of the frame")),
		"Length field in header should have been the length of payload")
}

func TestMarshalFrameWithType(t *testing.T) {
	f := baseFrame{}
	f.Type = byte(8)

	assert.Equal(t, byte(8), f.Marshal()[2],
		"Type should have been marshalled as the third octet")
}

func TestMarshalFrameWithFlags(t *testing.T) {
	f := baseFrame{}
	f.Flags = byte(0xD)

	assert.Equal(t, byte(0xD), f.Marshal()[3],
		"Flags should have been marshalled as the fourth octet")
}

func TestMarshalFrameWithStreamIdentifier(t *testing.T) {
	f := baseFrame{}
	f.StreamIdentifier = 168036609

	marshalled_f := f.Marshal()

	stream_identifier := binary.BigEndian.Uint32(marshalled_f[4:8]) & uint32(0x7FFFFFFF)

	assert.Equal(t, stream_identifier, f.StreamIdentifier,
		"Stream identifier should have been marshalled as fifth through eighth octets")
}

func TestMarshalGOAWAYFrameSetsType7(t *testing.T) {
	f := GOAWAYFrame{}
	assert.Equal(t, f.Marshal()[2], byte(7),
		"Type should have been marshalled as 0x7")
}

func TestMarshalGOAWAYFrameSetsNoFlags(t *testing.T) {
	f := GOAWAYFrame{}
	assert.Equal(t, f.Marshal()[3], byte(0),
		"Should have set no flags")
}

func TestMarshalGOAWAYFrameWithNoAdditionalDebugInfoSetsLength(t *testing.T) {
	f := GOAWAYFrame{}
	assert.Equal(t, lengthOfMarshalledFrame(f.Marshal()), uint16(16),
		"Length should have been 16 octets")
}

func TestMarshalGOAWAYFrameWithDebugInfoSetsLength(t *testing.T) {
	f := GOAWAYFrame{}
	f.AdditionalDebugData = "This is some additional debug info to help you"

	expectedLength := len(f.AdditionalDebugData) + 16

	assert.Equal(t, lengthOfMarshalledFrame(f.Marshal()),
		uint16(expectedLength),
		"Length should included the additional debug data")
}

func TestMarshalGOAWAYFrameIncludesLastStreamId(t *testing.T) {
	f := GOAWAYFrame{}
	f.ErrorCode = 12487291

	marshalled_f := f.Marshal()

	t.Log(marshalled_f)

	lastStreamId := binary.BigEndian.Uint32(marshalled_f[8:16])
	assert.Equal(t, lastStreamId, f.LastStreamId,
		"Marshalled frame should have included last stream id")
}

func TestMarshalGOAWAYFrameIncludesErrorCode(t *testing.T) {
	f := GOAWAYFrame{}
	f.ErrorCode = 12487291

	marshalled_f := f.Marshal()

	t.Log(marshalled_f)

	errorCode := binary.BigEndian.Uint32(marshalled_f[16:24])
	assert.Equal(t, errorCode, f.ErrorCode,
		"Marshalled frame should have included error code")
}
