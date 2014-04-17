package common

import (
	"encoding/binary"
	"testing"
	"github.com/stretchr/testify/assert"
)


type stringWrapper struct {
	data string
}

func frameType(marshalled_f []byte) uint8 {
	return uint8(marshalled_f[2])
}

func frameLength(marshalled_f []byte) uint16 {
	return binary.BigEndian.Uint16(marshalled_f[0:2]) & 0x7F
}

func frameFlags(marshalled_f []byte) uint8 {
	return uint8(marshalled_f[3])
}

func TestMarshalEmptyFrame(t *testing.T) {
	f := baseFrame{}

	assert.Equal(t, frameLength(f.Marshal()), uint16(0),
		"Length should have been nothing")
}

func TestMarshalFrameWithPayloadIncludesLength(t *testing.T) {
	f := baseFrame{}
	f.Payload = "this is the payload of the frame"

	marshalled_f := f.Marshal()

	assert.Equal(t, frameLength(marshalled_f),
		uint16(len("this is the payload of the frame")),
		"Length field in header should have been the length of payload")
}

func TestMarshalFrameWithType(t *testing.T) {
	f := baseFrame{}
	f.Type = uint8(8)

	assert.Equal(t, uint8(8), frameType(f.Marshal()),
		"Type should have been marshalled as the third octet")
}

func TestMarshalFrameWithFlags(t *testing.T) {
	f := baseFrame{}
	f.Flags = uint8(0xD)

	assert.Equal(t, uint8(0xD), frameFlags(f.Marshal()),
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
	assert.Equal(t, frameType(f.Marshal()), uint8(7),
		"Type should have been marshalled as 0x7")
}

func TestMarshalGOAWAYFrameSetsNoFlags(t *testing.T) {
	f := GOAWAYFrame{}
	assert.Equal(t, f.Marshal()[3], uint8(0),
		"Should have set no flags")
}

func TestMarshalGOAWAYFrameWithNoAdditionalDebugInfoSetsLength(t *testing.T) {
	f := GOAWAYFrame{}
	assert.Equal(t, frameLength(f.Marshal()), uint16(8),
		"Length should have been 8 octets")
}

func TestMarshalGOAWAYFrameWithDebugInfoSetsLength(t *testing.T) {
	f := GOAWAYFrame{}
	f.AdditionalDebugData = "This is some additional debug info to help you"

	expectedLength := len(f.AdditionalDebugData) + 8

	assert.Equal(t, frameLength(f.Marshal()),
		uint16(expectedLength),
		"Length should included the additional debug data")
}

func TestMarshalGOAWAYFrameIncludesLastStreamId(t *testing.T) {
	f := GOAWAYFrame{}
	f.ErrorCode = 12487291

	marshalled_f := f.Marshal()

	t.Log(marshalled_f)

	lastStreamId := binary.BigEndian.Uint32(marshalled_f[8:12])
	assert.Equal(t, lastStreamId, f.LastStreamId,
		"Marshalled frame should have included last stream id")
}

func TestMarshalGOAWAYFrameIncludesErrorCode(t *testing.T) {
	f := GOAWAYFrame{}
	f.ErrorCode = 12487291

	marshalled_f := f.Marshal()

	t.Log(marshalled_f)

	errorCode := binary.BigEndian.Uint32(marshalled_f[12:16])
	assert.Equal(t, errorCode, f.ErrorCode,
		"Marshalled frame should have included error code")
}

func TestMarshalPingFrameSetsType6(t *testing.T) {
	f := PingFrame{}
	assert.Equal(t, frameType(f.Marshal()), uint8(6),
		"Ping frame must have had a type of 0x6")
}

func TestMarshalPingFrameSetsOpaqueData(t *testing.T) {
	f := PingFrame{}
	f.OpaqueData = 219748174981749872
	marshalled_f := f.Marshal()
	opaqueData := binary.BigEndian.Uint64(marshalled_f[8:16])

	assert.Equal(t, opaqueData, f.OpaqueData,
		"Ping frame should have included opaque data")
}

func TestMarshalPingFrameHasLengthOf8(t *testing.T) {
	f := PingFrame{}
	assert.Equal(t, frameLength(f.Marshal()), uint8(8),
		"Ping frame must have had a length field value of 8")
}

func TestMarshalPingFrameIncludesAckIfSet(t *testing.T) {
	f := PingFrame{}
	f.Flags.ACK = true

	marshalled_f := f.Marshal()


	assert.Equal(t, frameFlags(marshalled_f) & 0x1, uint8(1),
		"Ping frame with ACK flag should have had 0x1 flag bit set")
}

func TestMarshalPingFrameDoesNotIncludeAckIfUnset(t *testing.T) {
	f := PingFrame{}
	f.Flags.ACK = false

	marshalled_f := f.Marshal()

	assert.Equal(t, frameFlags(marshalled_f) & 0x1, uint8(0),
		"Ping frame without ACK flag should not have had 0x1 flag bit set")
}
