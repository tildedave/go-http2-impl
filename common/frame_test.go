package common

import (
	"encoding/binary"
	"testing"
	"github.com/stretchr/testify/assert"
)


type stringWrapper struct {
	data string
}

func frameType(marshalled []byte) uint8 {
	return uint8(marshalled[2])
}

func frameLength(marshalled []byte) uint16 {
	return binary.BigEndian.Uint16(marshalled[0:2]) & 0x7F
}

func frameFlags(marshalled []byte) uint8 {
	return uint8(marshalled[3])
}

func TestMarshalEmptyFrame(t *testing.T) {
	f := baseFrame{}

	assert.Equal(t, frameLength(f.Marshal()), uint16(0),
		"Length should have been nothing")
}

func TestMarshalFrameWithPayloadIncludesLength(t *testing.T) {
	f := baseFrame{}
	f.Payload = "this is the payload of the frame"

	marshalled := f.Marshal()

	assert.Equal(t, frameLength(marshalled),
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

	marshalled := f.Marshal()

	stream_identifier := binary.BigEndian.Uint32(marshalled[4:8]) & uint32(0x7FFFFFFF)

	assert.Equal(t, stream_identifier, f.StreamIdentifier,
		"Stream identifier should have been marshalled as fifth through eighth octets")
}

func TestMarshalGOAWAYFrame(t *testing.T) {
	f := GOAWAYFrame{}
	f.ErrorCode = 12487291

	marshalled := f.Marshal()

	assert.Equal(t, frameType(marshalled), uint8(7),
		"Type should have been marshalled as 0x7")
	assert.Equal(t, frameFlags(marshalled), uint8(0),
		"Should have set no flags")
	assert.Equal(t, frameLength(marshalled), uint16(8),
		"Length should have been 8 octets")

	lastStreamId := binary.BigEndian.Uint32(marshalled[8:12])
	assert.Equal(t, lastStreamId, f.LastStreamId,
		"Marshalled frame should have included last stream id")
}

func TestMarshalGOAWAYFrameWithDebugInfoSetsLength(t *testing.T) {
	f := GOAWAYFrame{}
	f.AdditionalDebugData = "This is some additional debug info to help you"

	expectedLength := len(f.AdditionalDebugData) + 8

	assert.Equal(t, frameLength(f.Marshal()),
		uint16(expectedLength),
		"Length should included the additional debug data")
}

func TestMarshalPingFrame(t *testing.T) {
	f := PingFrame{}
	f.OpaqueData = 219748174981749872

	marshalled := f.Marshal()
	assert.Equal(t, frameType(marshalled), uint8(6),
		"Ping frame must have had a type of 0x6")

	opaqueData := binary.BigEndian.Uint64(marshalled[8:16])
	assert.Equal(t, opaqueData, f.OpaqueData,
		"Ping frame should have included opaque data")

	assert.Equal(t, frameLength(marshalled), uint8(8),
		"Ping frame must have had a length field value of 8")
}

func TestMarshalPingFrameIncludesAckIfSet(t *testing.T) {
	f := PingFrame{}
	f.Flags.ACK = true

	marshalled := f.Marshal()

	assert.Equal(t, frameFlags(marshalled) & 0x1, uint8(1),
		"Ping frame with ACK flag should have had 0x1 flag bit set")
}

func TestMarshalPingFrameDoesNotIncludeAckIfUnset(t *testing.T) {
	f := PingFrame{}
	f.Flags.ACK = false

	marshalled := f.Marshal()

	assert.Equal(t, frameFlags(marshalled) & 0x1, uint8(0),
		"Ping frame without ACK flag should not have had 0x1 flag bit set")
}

func TestMarshalDataFrameWithoutPadding(t *testing.T) {
	f := DataFrame{}
	f.Data = "This is the data associated with the data frame"

	marshalled := f.Marshal()

	assert.Equal(t, frameType(marshalled), uint8(0x0),
		"Data frame should have type 0x0")
	assert.Equal(t, frameLength(marshalled), uint16(len(f.Data)))

	assert.Equal(t, frameFlags(marshalled) & 0x08, byte(0),
		"Padding low flag should not have been set")
	assert.Equal(t, frameFlags(marshalled) & 0x10, byte(0),
		"Padding high flag should not have been set")

	assert.Equal(t, []byte(f.Data), marshalled[8:], "Data did not match")
}

func TestMarshalDataFrameWithEndStreamFlag(t *testing.T) {
	f := DataFrame{}
	f.Flags.END_STREAM = true

	marshalled := f.Marshal()
	assert.Equal(t, frameFlags(marshalled) & 0x1, uint8(0x1),
		"Data frame should have end stream flag set")
}

func TestMarshalDataFrameWithEndSegmentFlag(t *testing.T) {
	f := DataFrame{}
	f.Flags.END_SEGMENT = true

	marshalled := f.Marshal()
	assert.Equal(t, frameFlags(marshalled) & 0x2, uint8(0x2),
		"Data frame should have end segment flag set")
}

func TestMarshalDataFrameWithSmallAmountOfPadding(t *testing.T) {
	f := DataFrame{}
	f.Data = "This is the data associated with the frame"
	f.Padding = "This padding is less than 256 bytes"

	marshalled := f.Marshal()
	expectedLength := uint16(len(f.Data) + len(f.Padding) + 1)

	assert.Equal(t, frameLength(marshalled), expectedLength,
		"Length did not include the data, the padding, and the padding header fields")

	assert.Equal(t, frameFlags(marshalled) & 0x08, byte(0x08),
		"Padding low flag should have been set")
	assert.Equal(t, marshalled[8], uint8(len(f.Padding)),
		"Padding low should have been the length of the padding")
	assert.Equal(t, marshalled[9:9 + len(f.Data)], []byte(f.Data),
		"Data did not match")
	assert.Equal(t, marshalled[9 + len(f.Data):], []byte(f.Padding),
		"Padding did not match")
}

func TestMarshalDataFrameWithPaddingHighSet(t *testing.T) {
	f := DataFrame{}
	f.Data = "This is the data associated with the data frame"

	paddingLength := 310
	for i := 0; i < paddingLength; i++ {
		f.Padding += "a"
	}

	marshalled := f.Marshal()

	assert.Equal(t, frameFlags(marshalled) & 0x08, byte(0x08),
		"Padding low flag should have been set")
	assert.Equal(t, frameFlags(marshalled) & 0x10, byte(0x10),
		"Padding high flag should have been set")

	assert.Equal(t, binary.BigEndian.Uint16(marshalled[8:10]),
		uint16(len(f.Padding)),
		"Padding length should have been equal to length of padding")
	assert.Equal(t, marshalled[10:10 + len(f.Data)], []byte(f.Data),
		"Data did not match")
	assert.Equal(t, marshalled[10 + len(f.Data):], []byte(f.Padding),
		"Padding did not match")
}

func TestMarshalHeadersFrame(t *testing.T) {
	f := HeadersFrame{}
	f.HeaderBlockFragment = "accept-encoding:gzip"

	marshalled := f.Marshal()

	assert.Equal(t, frameType(marshalled), byte(0x01),
		"Type of HEADERS frame should have been 0x01")

	assert.Equal(t,
		marshalled[8:],
		[]byte("accept-encoding:gzip"))
}

func TestMarshalHeadersFrameWithPriorityGroup(t *testing.T) {
	f := HeadersFrame{}
	f.PriorityGroupIdentifier = 21984080
	f.Weight = 123
	f.HeaderBlockFragment = "accept-encoding:gzip"
	f.Flags.PRIORITY_GROUP = true

	marshalled := f.Marshal()

	assert.Equal(t, frameFlags(marshalled) & 0x20, byte(0x20),
		"Flag for PRIORITY_GROUP should have been set")

	assert.Equal(t, marshalled[8] & 0x80, byte(0x80),
		"R bit for PRIORITY_GROUP should have been set")

	assert.Equal(t,
		binary.BigEndian.Uint32(marshalled[8:12]) ^ 0x80000000,
		f.PriorityGroupIdentifier,
		"Priority group identifier did not match")

	assert.Equal(t, marshalled[12], byte(f.Weight), "Weight did not match")

	assert.Equal(t,
		marshalled[13:],
		[]byte("accept-encoding:gzip"))
}

func TestMarshalHeadersFrameWithPriorityDependency(t *testing.T) {
	f := HeadersFrame{}
	f.StreamDependency = 39781097
	f.Flags.PRIORITY_DEPENDENCY = true
	f.HeaderBlockFragment = "accept-encoding:gzip"

	marshalled := f.Marshal()

	assert.Equal(t, frameFlags(marshalled) & 0x40, byte(0x40),
		"Flag for PRIORITY_DEPENDENCY should have been set")

	assert.Equal(t, marshalled[8] & 0x80, byte(0x80),
		"E bit for PRIORITY_DEPENDENCY should have been set")

	assert.Equal(t,
		binary.BigEndian.Uint32(marshalled[8:12]) ^ 0x80000000,
		f.StreamDependency,
		"Stream dependency did not match")

	assert.Equal(t,
		marshalled[12:],
		[]byte("accept-encoding:gzip"))
}

func TestMarshalHeadersFrameWithSmallAmountOfPadding(t *testing.T) {
	f := HeadersFrame{}
	f.HeaderBlockFragment = "content-type:application/json"
	f.Padding = "This is less than 256 padding"

	marshalled := f.Marshal()

	assert.Equal(t,
		marshalled[8],
		byte(len(f.Padding)))

	assert.Equal(t, frameFlags(marshalled) & 0x08, byte(0x08),
		"Padding low flag should have been set")
	assert.Equal(t, marshalled[8], byte(len(f.Padding)),
		"Padding low length should have been set")
	assert.Equal(t, marshalled[9:9 + len(f.HeaderBlockFragment)],
		[]byte(f.HeaderBlockFragment),
		"Header block fragment should have matched")
	assert.Equal(t, marshalled[9 + len(f.HeaderBlockFragment):],
		[]byte(f.Padding),
		"Padding should have matched")
}

func TestMarshalHeadersFrameWithPaddingHighSet(t *testing.T) {
	f := HeadersFrame{}
	f.HeaderBlockFragment = "content-type:application/json"

	paddingLength := 371
	for i := 0; i < paddingLength; i++ {
		f.Padding += "b"
	}

	marshalled := f.Marshal()

	assert.Equal(t, binary.BigEndian.Uint16(marshalled[8:10]),
		uint16(len(f.Padding)),
		"Padding length should have been equal to length of padding")

	assert.Equal(t, frameFlags(marshalled) & 0x08, byte(0x08),
		"Padding low flag should have been set")
	assert.Equal(t, frameFlags(marshalled) & 0x10, byte(0x10),
		"Padding high flag should have been set")

	assert.Equal(t, marshalled[10:10 + len(f.HeaderBlockFragment)],
		[]byte(f.HeaderBlockFragment),
		"Header block fragment should have matched")
	assert.Equal(t, marshalled[10 + len(f.HeaderBlockFragment):],
		[]byte(f.Padding),
		"Padding should have matched")
}

func TestMarshalHeadersFrameWithEndStreamFlag(t *testing.T) {
	f := HeadersFrame{}
	f.Flags.END_STREAM = true

	marshalled := f.Marshal()
	assert.Equal(t, frameFlags(marshalled) & 0x1, uint8(0x1),
		"Headers frame should have end stream flag set")
}

func TestMarshalHeadersFrameWithEndSegmentFlag(t *testing.T) {
	f := HeadersFrame{}
	f.Flags.END_SEGMENT = true

	marshalled := f.Marshal()
	assert.Equal(t, frameFlags(marshalled) & 0x2, uint8(0x2),
		"Headers frame should have end segment flag set")
}

func TestMarshalHeadersFrameWithEndHeadersFlag(t *testing.T) {
	f := HeadersFrame{}
	f.Flags.END_HEADERS = true

	marshalled := f.Marshal()
	assert.Equal(t, frameFlags(marshalled) & 0x4, uint8(0x4),
		"Headers frame should have end headers flag set")
}
