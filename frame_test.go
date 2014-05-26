package main

import (
	"encoding/binary"
	"github.com/stretchr/testify/assert"
	"testing"
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
	f := base{}

	assert.Equal(t, frameLength(f.Marshal()), uint16(0),
		"Length should have been nothing")
}

func TestMarshalFrame_WithPayloadIncludesLength(t *testing.T) {
	f := base{}
	f.Payload = "this is the payload of the frame"

	marshalled := f.Marshal()

	assert.Equal(t, frameLength(marshalled),
		uint16(len("this is the payload of the frame")),
		"Length field in header should have been the length of payload")
}

func TestMarshalFrame_WithType(t *testing.T) {
	f := base{}
	f.Type = uint8(8)

	assert.Equal(t, uint8(8), frameType(f.Marshal()),
		"Type should have been marshalled as the third octet")
}

func TestMarshalFrame_WithFlags(t *testing.T) {
	f := base{}
	f.Flags = uint8(0xD)

	assert.Equal(t, uint8(0xD), frameFlags(f.Marshal()),
		"Flags should have been marshalled as the fourth octet")
}

func TestMarshalFrame_WithStreamId(t *testing.T) {
	f := base{}
	f.StreamId = 168036609

	marshalled := f.Marshal()

	stream_identifier := binary.BigEndian.Uint32(marshalled[4:8]) & uint32(0x7FFFFFFF)

	assert.Equal(t, stream_identifier, f.StreamId,
		"Stream identifier should have been marshalled as fifth through eighth octets")
}

func TestMarshalGOAWAY(t *testing.T) {
	f := GOAWAY{}
	f.ErrorCode = 12487291
	f.AdditionalDebugData = "This is some additional debug info to help you"

	marshalled := f.Marshal()

	assert.Equal(t, frameType(marshalled), uint8(7),
		"Type should have been marshalled as 0x7")
	assert.Equal(t, frameFlags(marshalled), uint8(0),
		"Should have set no flags")
	assert.Equal(t, frameLength(marshalled), uint16(8+len(f.AdditionalDebugData)),
		"Length should have been 8 octets")

	lastStreamId := binary.BigEndian.Uint32(marshalled[8:12])
	assert.Equal(t, lastStreamId, f.LastStreamId,
		"Marshalled frame should have included last stream id")
}

func TestMarshalGOAWAY_WithDebugInfoSetsLength(t *testing.T) {
	f := GOAWAY{}
	f.AdditionalDebugData = "This is some additional debug info to help you"

	expectedLength := len(f.AdditionalDebugData) + 8

	assert.Equal(t, frameLength(f.Marshal()),
		uint16(expectedLength),
		"Length should included the additional debug data")
}

func TestMarshalPING(t *testing.T) {
	f := PING{}
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

func TestMarshalPINGIncludesAckIfSet(t *testing.T) {
	f := PING{}
	f.Flags.ACK = true

	marshalled := f.Marshal()

	assert.Equal(t, frameFlags(marshalled)&0x1, uint8(1),
		"Ping frame with ACK flag should have had 0x1 flag bit set")
}

func TestMarshalPINGDoesNotIncludeAckIfUnset(t *testing.T) {
	f := PING{}
	f.Flags.ACK = false

	marshalled := f.Marshal()

	assert.Equal(t, frameFlags(marshalled)&0x1, uint8(0),
		"Ping frame without ACK flag should not have had 0x1 flag bit set")
}

func TestMarshalDATA_WithoutPadding(t *testing.T) {
	f := DATA{}
	f.Data = "This is the data associated with the data frame"

	marshalled := f.Marshal()

	assert.Equal(t, frameType(marshalled), uint8(0x0),
		"Data frame should have type 0x0")
	assert.Equal(t, frameLength(marshalled), uint16(len(f.Data)))

	assert.Equal(t, frameFlags(marshalled)&0x08, byte(0),
		"Padding low flag should not have been set")
	assert.Equal(t, frameFlags(marshalled)&0x10, byte(0),
		"Padding high flag should not have been set")

	assert.Equal(t, []byte(f.Data), marshalled[8:], "Data did not match")
}

func TestMarshalDATA_WithEndStreamFlag(t *testing.T) {
	f := DATA{}
	f.Flags.END_STREAM = true

	marshalled := f.Marshal()
	assert.Equal(t, frameFlags(marshalled)&0x1, uint8(0x1),
		"Data frame should have end stream flag set")
}

func TestMarshalDATA_WithEndSegmentFlag(t *testing.T) {
	f := DATA{}
	f.Flags.END_SEGMENT = true

	marshalled := f.Marshal()
	assert.Equal(t, frameFlags(marshalled)&0x2, uint8(0x2),
		"Data frame should have end segment flag set")
}

func TestMarshalDATA_WithSmallAmountOfPadding(t *testing.T) {
	f := DATA{}
	f.Data = "This is the data associated with the frame"
	f.Padding = "This padding is less than 256 bytes"

	marshalled := f.Marshal()
	expectedLength := uint16(len(f.Data) + len(f.Padding) + 1)

	assert.Equal(t, frameLength(marshalled), expectedLength,
		"Length did not include the data, the padding, and the padding header fields")

	assert.Equal(t, frameFlags(marshalled)&0x08, byte(0x08),
		"Padding low flag should have been set")
	assert.Equal(t, marshalled[8], uint8(len(f.Padding)),
		"Padding low should have been the length of the padding")
	assert.Equal(t, marshalled[9:9+len(f.Data)], []byte(f.Data),
		"Data did not match")
	assert.Equal(t, marshalled[9+len(f.Data):], []byte(f.Padding),
		"Padding did not match")
}

func TestMarshalDATA_WithPaddingHighSet(t *testing.T) {
	f := DATA{}
	f.Data = "This is the data associated with the data frame"

	paddingLength := 310
	for i := 0; i < paddingLength; i++ {
		f.Padding += "a"
	}

	marshalled := f.Marshal()

	assert.Equal(t, frameFlags(marshalled)&0x08, byte(0x08),
		"Padding low flag should have been set")
	assert.Equal(t, frameFlags(marshalled)&0x10, byte(0x10),
		"Padding high flag should have been set")

	assert.Equal(t, binary.BigEndian.Uint16(marshalled[8:10]),
		uint16(len(f.Padding)),
		"Padding length should have been equal to length of padding")
	assert.Equal(t, marshalled[10:10+len(f.Data)], []byte(f.Data),
		"Data did not match")
	assert.Equal(t, marshalled[10+len(f.Data):], []byte(f.Padding),
		"Padding did not match")
}

func TestMarshalHEADERS(t *testing.T) {
	f := HEADERS{}
	f.HeaderBlockFragment = "accept-encoding:gzip"

	marshalled := f.Marshal()

	assert.Equal(t, frameType(marshalled), byte(0x01),
		"Type of HEADERS frame should have been 0x01")

	assert.Equal(t,
		marshalled[8:],
		[]byte("accept-encoding:gzip"))
}

func TestMarshalHEADERS_WithPriorityGroup(t *testing.T) {
	f := HEADERS{}
	f.PriorityGroupId = 21984080
	f.Weight = 123
	f.HeaderBlockFragment = "accept-encoding:gzip"
	f.Flags.PRIORITY_GROUP = true

	marshalled := f.Marshal()

	assert.Equal(t, frameFlags(marshalled)&0x20, byte(0x20),
		"Flag for PRIORITY_GROUP should have been set")

	assert.Equal(t, marshalled[8]&0x80, byte(0x80),
		"R bit for PRIORITY_GROUP should have been set")

	assert.Equal(t,
		binary.BigEndian.Uint32(marshalled[8:12])^0x80000000,
		f.PriorityGroupId,
		"Priority group identifier did not match")

	assert.Equal(t, marshalled[12], byte(f.Weight), "Weight did not match")

	assert.Equal(t,
		marshalled[13:],
		[]byte("accept-encoding:gzip"))
}

func TestMarshalHEADERS_WithPriorityDependency(t *testing.T) {
	f := HEADERS{}
	f.StreamDependency = 39781097
	f.Flags.PRIORITY_DEPENDENCY = true
	f.Flags.EXCLUSIVE = true
	f.HeaderBlockFragment = "accept-encoding:gzip"

	marshalled := f.Marshal()

	assert.Equal(t, frameFlags(marshalled)&0x40, byte(0x40),
		"Flag for PRIORITY_DEPENDENCY should have been set")

	assert.Equal(t, marshalled[8]&0x80, byte(0x80),
		"E bit for PRIORITY_DEPENDENCY should have been set")

	assert.Equal(t,
		binary.BigEndian.Uint32(marshalled[8:12])^0x80000000,
		f.StreamDependency,
		"Stream dependency did not match")

	assert.Equal(t,
		marshalled[12:],
		[]byte("accept-encoding:gzip"))
}

func TestMarshalHEADERS_WithSmallAmountOfPadding(t *testing.T) {
	f := HEADERS{}
	f.HeaderBlockFragment = "content-type:application/json"
	f.Padding = "This is less than 256 padding"

	marshalled := f.Marshal()

	assert.Equal(t,
		marshalled[8],
		byte(len(f.Padding)))

	assert.Equal(t, frameFlags(marshalled)&0x08, byte(0x08),
		"Padding low flag should have been set")
	assert.Equal(t, marshalled[8], byte(len(f.Padding)),
		"Padding low length should have been set")
	assert.Equal(t, marshalled[9:9+len(f.HeaderBlockFragment)],
		[]byte(f.HeaderBlockFragment),
		"Header block fragment should have matched")
	assert.Equal(t, marshalled[9+len(f.HeaderBlockFragment):],
		[]byte(f.Padding),
		"Padding should have matched")
}

func TestMarshalHEADERS_WithPaddingHighSet(t *testing.T) {
	f := HEADERS{}
	f.HeaderBlockFragment = "content-type:application/json"

	paddingLength := 371
	for i := 0; i < paddingLength; i++ {
		f.Padding += "b"
	}

	marshalled := f.Marshal()

	assert.Equal(t, binary.BigEndian.Uint16(marshalled[8:10]),
		uint16(len(f.Padding)),
		"Padding length should have been equal to length of padding")

	assert.Equal(t, frameFlags(marshalled)&0x08, byte(0x08),
		"Padding low flag should have been set")
	assert.Equal(t, frameFlags(marshalled)&0x10, byte(0x10),
		"Padding high flag should have been set")

	assert.Equal(t, marshalled[10:10+len(f.HeaderBlockFragment)],
		[]byte(f.HeaderBlockFragment),
		"Header block fragment should have matched")
	assert.Equal(t, marshalled[10+len(f.HeaderBlockFragment):],
		[]byte(f.Padding),
		"Padding should have matched")
}

func TestMarshalHEADERS_WithEndStreamFlag(t *testing.T) {
	f := HEADERS{}
	f.Flags.END_STREAM = true

	marshalled := f.Marshal()
	assert.Equal(t, frameFlags(marshalled)&0x1, uint8(0x1),
		"Headers frame should have end stream flag set")
}

func TestMarshalHEADERS_WithEndSegmentFlag(t *testing.T) {
	f := HEADERS{}
	f.Flags.END_SEGMENT = true

	marshalled := f.Marshal()
	assert.Equal(t, frameFlags(marshalled)&0x2, uint8(0x2),
		"Headers frame should have end segment flag set")
}

func TestMarshalHEADERS_WithEndHeadersFlag(t *testing.T) {
	f := HEADERS{}
	f.Flags.END_HEADERS = true

	marshalled := f.Marshal()
	assert.Equal(t, frameFlags(marshalled)&0x4, uint8(0x4),
		"Headers frame should have end headers flag set")
}

func TestMarshalPRIORITY_WithPriorityDependency(t *testing.T) {
	f := PRIORITY{}
	f.Flags.PRIORITY_DEPENDENCY = true
	f.StreamId = 1111
	f.StreamDependency = 123456

	marshalled := f.Marshal()
	assert.Equal(t, frameType(marshalled), uint8(0x2),
		"Expected frame type of priority to be 0x2")
	assert.Equal(t, frameFlags(marshalled)&0x40, uint8(0x40),
		"Headers frame should have had priority group set")
	assert.Equal(t, binary.BigEndian.Uint32(marshalled[4:8]), f.StreamId,
		"Stream identifier was not correct in the payload")
	assert.Equal(t, binary.BigEndian.Uint32(marshalled[8:12]), f.StreamDependency,
		"Stream dependency was not correct in the payload")
}

func TestMarshalPRIORITY_WithExclusivePriorityDependency(t *testing.T) {
	f := PRIORITY{}
	f.Flags.PRIORITY_DEPENDENCY = true
	f.StreamDependency = 123456
	f.Flags.EXCLUSIVE = true

	marshalled := f.Marshal()
	assert.Equal(t, frameType(marshalled), uint8(0x2),
		"Expected frame type of priority to be 0x2")
	assert.Equal(t, frameFlags(marshalled)&0x40, uint8(0x40),
		"Headers frame should have had priority group set")

	dependency := binary.BigEndian.Uint32([]byte{
		marshalled[8] & 0x7F,
		marshalled[9],
		marshalled[10],
		marshalled[11],
	})
	assert.Equal(t, dependency, f.StreamDependency,
		"Stream dependency was not correct in the payload")
	assert.Equal(t, marshalled[8]&0x80, uint8(0x80))
}

func TestMarshalPRIORITY_WithPriorityGroup(t *testing.T) {
	f := PRIORITY{}
	f.Flags.PRIORITY_GROUP = true
	f.PriorityGroupId = 912742
	f.Weight = 111

	marshalled := f.Marshal()
	assert.Equal(t, frameType(marshalled), uint8(0x2),
		"Expected frame type of priority to be 0x2")
	assert.Equal(t, frameFlags(marshalled)&0x20, uint8(0x20),
		"Headers frame should have had priority group set")
	assert.Equal(t, binary.BigEndian.Uint32(marshalled[8:12]),
		f.PriorityGroupId,
		"Priority Group Id was not correct in the payload")
	assert.Equal(t, marshalled[12], f.Weight,
		"Weight was not correct in the payload")
}

func TestMarshalRST_STREAM(t *testing.T) {
	f := RST_STREAM{}
	f.StreamId = 123
	f.ErrorCode = 2

	marshalled := f.Marshal()

	assert.Equal(t, frameType(marshalled), uint8(0x3),
		"Expected frame type of settings to be 0x3")
	assert.Equal(t, frameLength(marshalled), uint16(4),
		"Expected frame length to be 4 octets")
	assert.Equal(t, binary.BigEndian.Uint32(marshalled[4:8]), f.StreamId,
		"Expected stream identifier to match")
	assert.Equal(t, binary.BigEndian.Uint32(marshalled[8:12]), f.ErrorCode,
		"Expected error code to match")
}

func TestMarshalSETTINGS(t *testing.T) {
	f := SETTINGS{}
	f.Parameters = []Parameter{{uint8(1), uint32(1298431729)},
		{uint8(2), uint32(1478921795)}}

	marshalled := f.Marshal()

	assert.Equal(t, frameType(marshalled), uint8(0x4),
		"Expected frame type of settings to be 0x4")
	assert.Equal(t, frameLength(marshalled), uint16(10),
		"Expected frame length to be 10 octets (two parameters)")

	assert.Equal(t, marshalled[8],
		f.Parameters[0].Id,
		"Expected first frame parameter to match")
	assert.Equal(t, binary.BigEndian.Uint32(marshalled[9:13]),
		f.Parameters[0].Value,
		"Expected first frame value to match")
	assert.Equal(t, marshalled[13],
		f.Parameters[1].Id,
		"Expected second frame parameter to match")
	assert.Equal(t, binary.BigEndian.Uint32(marshalled[14:18]),
		f.Parameters[1].Value,
		"Expected second frame value to match")
}

func TestMarshalSETTINGS_WithACKFlag(t *testing.T) {
	f := SETTINGS{}
	f.Flags.ACK = true

	marshalled := f.Marshal()

	assert.Equal(t, frameFlags(marshalled)&0x1, uint8(0x1),
		"Expected frame to have set ACK flag of 0x1")
}

func TestMarshalPUSH_PROMISE(t *testing.T) {
	f := PUSH_PROMISE{}
	f.StreamId = 123
	f.PromisedStreamId = 456
	f.HeaderBlockFragment = "fragment of header block"
	f.Padding = "aaaaaaaaaaaaa"

	marshalled := f.Marshal()

	assert.Equal(t, frameType(marshalled), uint8(0x5),
		"Expected frame type for push promise frame to be 0x5")
	assert.Equal(t, frameFlags(marshalled)&0x8, uint8(0x8),
		"PAD_LOW flag should have been set")
	assert.Equal(t, frameFlags(marshalled)&0x10, uint8(0x0),
		"PAD_HIGH flag should not have been set")
	assert.Equal(t, marshalled[8], uint8(len(f.Padding)))
	assert.Equal(t, binary.BigEndian.Uint32(marshalled[9:13]), f.PromisedStreamId)
	assert.Equal(t, string(marshalled[13:13+len(f.HeaderBlockFragment)]), f.HeaderBlockFragment)
	assert.Equal(t, string(marshalled[13+len(f.HeaderBlockFragment):]), f.Padding)
}
func TestMarshalPUSH_PROMISE_WithEndHeadersFlag(t *testing.T) {
	f := PUSH_PROMISE{}
	f.StreamId = 123
	f.Flags.END_HEADERS = true
	marshalled := f.Marshal()

	assert.Equal(t, frameFlags(marshalled)&0x4, uint8(0x4),
		"END_HEADERS flag should have been set")
}

func TestMarshalCONTINUATION(t *testing.T) {
	f := CONTINUATION{}
	f.StreamId = 123
	f.HeaderBlockFragment = "fragment of header block"
	f.Padding = "aaaaaaaaaaaaa"

	marshalled := f.Marshal()

	assert.Equal(t, frameType(marshalled), uint8(0x9),
		"Expected frame type for push promise frame to be 0x9")
	assert.Equal(t, frameFlags(marshalled)&0x8, uint8(0x8),
		"PAD_LOW flag should have been set")
	assert.Equal(t, frameFlags(marshalled)&0x10, uint8(0x0),
		"PAD_HIGH flag should not have been set")
	assert.Equal(t, marshalled[8], uint8(len(f.Padding)))
	assert.Equal(t, string(marshalled[9:9+len(f.HeaderBlockFragment)]), f.HeaderBlockFragment)
	assert.Equal(t, string(marshalled[9+len(f.HeaderBlockFragment):]), f.Padding)
}

func TestMarshalCONTINUATION_WithEndHeadersFlag(t *testing.T) {
	f := CONTINUATION{}
	f.StreamId = 123
	f.Flags.END_HEADERS = true

	marshalled := f.Marshal()

	assert.Equal(t, frameFlags(marshalled)&0x4, uint8(0x4),
		"END_HEADERS flag should have been set")
}

func TestUnmarshalDATA_WithSmallPadding(t *testing.T) {
	f := DATA{
		StreamId: 37,
		Data:     "This is the data associated with the data frame",
		Padding:  "This padding is less than 256 bytes",
	}
	b := f.Marshal()
	uf, err := Unmarshal(b)

	assert.Nil(t, err)
	assert.IsType(t, DATA{}, uf)
	assert.Equal(t, uf, f)
}

func TestUnmarshalDATA_WithLargePadding(t *testing.T) {
	f := DATA{
		StreamId: 37,
		Data:     "This is the data associated with the data frame",
		Padding:  "",
	}
	paddingLength := 310
	for i := 0; i < paddingLength; i++ {
		f.Padding += string(0x00)
	}

	b := f.Marshal()
	uf, err := Unmarshal(b)

	assert.Nil(t, err)
	assert.IsType(t, DATA{}, uf)
	assert.Equal(t, uf, f)
}

func TestUnmarshalDATA_WithEndStream(t *testing.T) {
	f := DATA{StreamId: 123}
	f.Flags.END_STREAM = true

	b := f.Marshal()
	uf, err := Unmarshal(b)

	assert.Nil(t, err)
	assert.IsType(t, DATA{}, uf)
	assert.True(t, uf.(DATA).Flags.END_STREAM)
}

func TestUnmarshalDATA_WithEndSegment(t *testing.T) {
	f := DATA{StreamId: 123}
	f.Flags.END_SEGMENT = true

	b := f.Marshal()
	uf, err := Unmarshal(b)

	assert.Nil(t, err)
	assert.IsType(t, DATA{}, uf)
	assert.True(t, uf.(DATA).Flags.END_SEGMENT)
}

func assertUnmarshalError(t *testing.T, b []byte, expectedError error) {
	uf, err := Unmarshal(b)

	assert.Nil(t, uf)
	assert.Equal(t, err, expectedError)
}

func TestUnmarshalDATA_NoStreamId(t *testing.T) {
	f := DATA{}
	f.StreamId = 0

	assertUnmarshalError(t, f.Marshal(), ConnectionError{PROTOCOL_ERROR, "DATA frame must have stream identifier"})
}

func TestUnmarshalDATA_IncompatiblePaddingFlags(t *testing.T) {
	f := DATA{StreamId: 123, Data: "dagljkjagldka"}
	b := f.Marshal()
	b[3] = 0x10

	assertUnmarshalError(t, b, ConnectionError{PROTOCOL_ERROR, "PAD_HIGH was set but PAD_LOW was not set"})
}

func TestUnmarshalPING(t *testing.T) {
	f := PING{}
	f.OpaqueData = 2198179

	b := f.Marshal()
	uf, err := Unmarshal(b)

	assert.Nil(t, err)
	assert.IsType(t, PING{}, uf)
	assert.Equal(t, f, uf)
}

func TestUnmarshalPING_WithACK(t *testing.T) {
	f := PING{}
	f.OpaqueData = 2198179
	f.Flags.ACK = true

	b := f.Marshal()
	uf, err := Unmarshal(b)

	assert.Nil(t, err)
	assert.IsType(t, PING{}, uf)
	assert.Equal(t, f, uf)
}

func TestUnmarshalPING_WithStreamId(t *testing.T) {
	f := PING{}
	f.OpaqueData = 2198179

	b := f.Marshal()
	b[4] = 10
	assertUnmarshalError(t, b, ConnectionError{PROTOCOL_ERROR, "PING frame must not have stream identifier"})
}

func TestUnmarshalPING_WithBadLength(t *testing.T) {
	f := PING{}
	f.OpaqueData = 2198179
	b := f.Marshal()
	b[1] = 7

	assertUnmarshalError(t, b, ConnectionError{FRAME_SIZE_ERROR, "PING payload must have length of 8"})
}

func TestUnmarshalGOAWAY(t *testing.T) {
	f := GOAWAY{
		LastStreamId:        0,
		ErrorCode:           PROTOCOL_ERROR,
		AdditionalDebugData: "Malformed frame",
	}
	b := f.Marshal()

	uf, err := Unmarshal(b)

	assert.Nil(t, err)
	assert.IsType(t, GOAWAY{}, uf)
	assert.Equal(t, f, uf)
}

func TestUnmarshalHEADERS(t *testing.T) {
	f := HEADERS{}
	f.StreamId = 2139480
	f.HeaderBlockFragment = "accept-encoding:gzip"
	f.Flags.END_HEADERS = true

	b := f.Marshal()
	uf, err := Unmarshal(b)

	assert.Nil(t, err)
	assert.IsType(t, HEADERS{}, uf)
	assert.Equal(t, f, uf)
}

func TestUnmarshalHEADERS_WithPriorityGroup(t *testing.T) {
	f := HEADERS{}
	f.StreamId = 2139480
	f.PriorityGroupId = 21984080
	f.Weight = 123
	f.HeaderBlockFragment = "accept-encoding:gzip"
	f.Flags.PRIORITY_GROUP = true

	b := f.Marshal()
	uf, err := Unmarshal(b)

	assert.Nil(t, err)
	assert.IsType(t, HEADERS{}, uf)
	assert.Equal(t, f, uf)
}

func TestUnmarshalHEADERS_WithPriorityDependency(t *testing.T) {
	f := HEADERS{}
	f.StreamId = 2139480
	f.StreamDependency = 39781097
	f.Flags.PRIORITY_DEPENDENCY = true
	f.Flags.EXCLUSIVE = true
	f.HeaderBlockFragment = "accept-encoding:gzip"

	b := f.Marshal()
	uf, err := Unmarshal(b)

	assert.Nil(t, err)
	assert.IsType(t, HEADERS{}, uf)
	assert.Equal(t, f, uf)
}

func TestUnmarshalHEADERS_WithNoStreamId(t *testing.T) {
	f := HEADERS{}

	b := f.Marshal()

	assertUnmarshalError(t, b, ConnectionError{PROTOCOL_ERROR, "HEADERS frame must have stream identifier"})
}

func TestUnmarshalHEADERS_WithConflictingPriorityGroupAndDependencies(t *testing.T) {
	f := HEADERS{}
	f.StreamId = 1
	b := f.Marshal()
	b[3] |= 0x20
	b[3] |= 0x40

	assertUnmarshalError(t, b, ConnectionError{PROTOCOL_ERROR, "Cannot set both PRIORITY_GROUP and PRIORITY_DEPENDENCY flags"})
}

func TestUnmarshalPRIORITY_WithExclusivePriorityDependency(t *testing.T) {
	f := PRIORITY{}
	f.StreamId = 111
	f.StreamDependency = 123
	f.Flags.PRIORITY_DEPENDENCY = true
	f.Flags.EXCLUSIVE = true

	b := f.Marshal()
	uf, err := Unmarshal(b)

	assert.Nil(t, err)
	assert.IsType(t, PRIORITY{}, uf)
	assert.Equal(t, f, uf)
}

func TestUnmarshalPRIORITY_WithPriorityGroup(t *testing.T) {
	f := PRIORITY{}
	f.StreamId = 111
	f.PriorityGroupId = 555
	f.Weight = 123
	f.Flags.PRIORITY_GROUP = true

	b := f.Marshal()
	uf, err := Unmarshal(b)

	assert.Nil(t, err)
	assert.IsType(t, PRIORITY{}, uf)
	assert.Equal(t, f, uf)
}

func TestUnmarshalPRIORITY_WithNoStreamId(t *testing.T) {
	f := PRIORITY{}
	b := f.Marshal()

	assertUnmarshalError(t, b, ConnectionError{PROTOCOL_ERROR, "PRIORITY frame must have stream identifier"})
}

func TestUnmarshalPRIORITY_WithConflictingFlags(t *testing.T) {
	f := PRIORITY{}
	f.StreamId = 1
	b := f.Marshal()
	b[3] |= 0x20
	b[3] |= 0x40

	assertUnmarshalError(t, b, ConnectionError{PROTOCOL_ERROR, "Cannot set both PRIORITY_GROUP and PRIORITY_DEPENDENCY flags"})
}

func TestUnmarshalRST_STREAM(t *testing.T) {
	f := RST_STREAM{}
	f.ErrorCode = 12390
	f.StreamId = 123

	b := f.Marshal()
	uf, err := Unmarshal(b)

	assert.Nil(t, err)
	assert.IsType(t, RST_STREAM{}, uf)
	assert.Equal(t, f, uf)
}

func TestUnmarshalRST_STREAM_WithNoStreamIdenfier(t *testing.T) {
	f := RST_STREAM{}

	assertUnmarshalError(t, f.Marshal(), ConnectionError{PROTOCOL_ERROR, "RST_STREAM frame must have stream identifier"})
}

func TestUnmarshalSETTINGS(t *testing.T) {
	f := SETTINGS{}
	f.Parameters = []Parameter{
		{SETTINGS_HEADER_TABLE_SIZE, 512},
		{SETTINGS_INITIAL_WINDOW_SIZE, 120000},
	}

	b := f.Marshal()
	uf, err := Unmarshal(b)

	assert.Nil(t, err)
	assert.IsType(t, SETTINGS{}, uf)
	assert.Equal(t, f, uf)
}

func TestUnmarshalSETTINGS_WithAck(t *testing.T) {
	f := SETTINGS{}
	f.Flags.ACK = true

	b := f.Marshal()
	uf, err := Unmarshal(b)

	assert.Nil(t, err)
	assert.IsType(t, SETTINGS{}, uf)
	assert.Equal(t, f, uf)
}

func TestUnmarshalSETTINGS_WithInvalidId(t *testing.T) {
	f := SETTINGS{}
	f.Parameters = []Parameter{{15, 512}}

	assertUnmarshalError(t, f.Marshal(), ConnectionError{PROTOCOL_ERROR, "Settings frame specified invalid identifier: 15"})
}

func TestUnmarshalSETTINGS_WithAckAndPayload(t *testing.T) {
	f := SETTINGS{}
	f.Parameters = []Parameter{{SETTINGS_HEADER_TABLE_SIZE, 512}}
	f.Flags.ACK = true

	assertUnmarshalError(t, f.Marshal(), ConnectionError{FRAME_SIZE_ERROR, "Payload of Settings frame with ACK flag must be empty"})
}

func TestUnmarshalPUSH_PROMISE(t *testing.T) {
	f := PUSH_PROMISE{}
	f.StreamId = 123
	f.PromisedStreamId = 456
	f.HeaderBlockFragment = "fragment of header block"
	f.Padding = "aaaaaaaaaaaaa"
	f.Flags.END_HEADERS = true

	b := f.Marshal()
	uf, err := Unmarshal(b)

	assert.Nil(t, err)
	assert.IsType(t, PUSH_PROMISE{}, uf)
	assert.Equal(t, f, uf)
}

func TestUnmarshalPUSH_PROMISE_NoStreamId(t *testing.T) {
	f := PUSH_PROMISE{}

	assertUnmarshalError(t, f.Marshal(), ConnectionError{PROTOCOL_ERROR, "PUSH_PROMISE frame must have stream identifier"})
}

func TestUnmarshalCONTINUATION(t *testing.T) {
	f := CONTINUATION{}
	f.StreamId = 123
	f.HeaderBlockFragment = "fragment of header block"
	f.Padding = "aaaaaaaaaaaaa"
	f.Flags.END_HEADERS = true

	b := f.Marshal()
	uf, err := Unmarshal(b)

	assert.Nil(t, err)
	assert.IsType(t, CONTINUATION{}, uf)
	assert.Equal(t, f, uf)
}

func TestUnmarshalCONTINUATION_NoStreamId(t *testing.T) {
	f := CONTINUATION{}

	assertUnmarshalError(t, f.Marshal(), ConnectionError{PROTOCOL_ERROR, "CONTINUATION frame must have stream identifier"})
}

func TestUnmarshalIncompleteHeader(t *testing.T) {
	f := PING{}
	f.OpaqueData = 2198179

	b := f.Marshal()[0:3]
	uf, err := Unmarshal(b)

	assert.Nil(t, err)
	assert.Nil(t, uf)
}

func TestUnmarshalIncompletePayload(t *testing.T) {
	f := PING{}
	f.OpaqueData = 2198179

	b := f.Marshal()[0:11]
	uf, err := Unmarshal(b)

	assert.Nil(t, err)
	assert.Nil(t, uf)
}
