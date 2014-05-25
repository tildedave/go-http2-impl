package main

import (
	"fmt"
	"encoding/binary"
)

var _ = fmt.Printf  // package fmt is now used

// base is used internally to marshal other types of frames
type base struct {
	Type             uint8
	Flags            uint8
	StreamIdentifier uint32
	Payload          string
}

// http://tools.ietf.org/html/draft-ietf-httpbis-http2-11#section-6.1
type DATA struct {
	StreamIdentifier uint32
	Data    string
	Padding string

	Flags struct {
		END_STREAM  bool
		END_SEGMENT bool
	}
}

// http://tools.ietf.org/html/draft-ietf-httpbis-http2-11#section-6.2
type HEADERS struct {
	PriorityGroupIdentifier uint32
	Weight                  uint8
	StreamDependency        uint32
	HeaderBlockFragment     string
	Padding                 string

	Flags struct {
		END_STREAM          bool
		END_SEGMENT         bool
		END_HEADERS         bool
		PRIORITY_GROUP      bool
		PRIORITY_DEPENDENCY bool
	}
}

// http://tools.ietf.org/html/draft-ietf-httpbis-http2-11#page-35
type Parameter struct {
	Identifier uint8
	Value      uint32
}

type SETTINGS struct {
	Parameters []Parameter
	Flags      struct {
		ACK bool
	}
}

// http://tools.ietf.org/html/draft-ietf-httpbis-http2-11#section-6.7
type PING struct {
	OpaqueData uint64
	Flags      struct {
		ACK bool
	}
}

// http://tools.ietf.org/html/draft-ietf-httpbis-http2-11#section-6.8
type GOAWAY struct {
	LastStreamId        uint32
	ErrorCode           uint32
	AdditionalDebugData string
}

type Frame interface {
	Marshal() []byte
}

func (f base) Marshal() []byte {
	header := make([]byte, 8)
	binary.BigEndian.PutUint16(header, uint16(len(f.Payload)))
	header[2] = f.Type
	header[3] = f.Flags
	binary.BigEndian.PutUint32(header[4:8], f.StreamIdentifier)

	return append(header, f.Payload...)
}

func (f GOAWAY) Marshal() []byte {
	b := base{}
	b.Type = 0x7

	payload := make([]byte, 8+len(f.AdditionalDebugData))
	binary.BigEndian.PutUint32(payload[0:4], f.LastStreamId)
	binary.BigEndian.PutUint32(payload[4:8], f.ErrorCode)
	copy(payload, f.AdditionalDebugData)

	b.Payload = string(payload)

	return b.Marshal()
}

func (f PING) Marshal() []byte {
	b := base{}
	b.Type = 0x6
	if f.Flags.ACK {
		b.Flags = 0x1
	}

	payload := make([]byte, 8)
	binary.BigEndian.PutUint64(payload, f.OpaqueData)
	b.Payload = string(payload)

	return b.Marshal()
}

func paddingHeaders(b *base, padding string) []byte {
	paddingHeaders := make([]byte, 0, 2)
	paddingLength := uint16(len(padding))
	if paddingLength > 0 {
		// set PADDING_LOW flag
		b.Flags |= 0x08

		if paddingLength > 256 {
			// set PADDING_HIGH flag
			b.Flags |= 0x10
			paddingHeaders = paddingHeaders[0:2]
			binary.BigEndian.PutUint16(paddingHeaders, paddingLength)

		} else {
			paddingHeaders = paddingHeaders[0:1]
			paddingHeaders[0] = uint8(paddingLength)
		}
	}

	return paddingHeaders
}

func (f DATA) Marshal() []byte {
	b := base{}
	b.Type = 0x0
	b.StreamIdentifier = f.StreamIdentifier

	payload := paddingHeaders(&b, f.Padding)
	payload = append(payload, f.Data...)
	payload = append(payload, f.Padding...)
	b.Payload = string(payload)

	if f.Flags.END_STREAM {
		b.Flags |= 0x01
	}
	if f.Flags.END_SEGMENT {
		b.Flags |= 0x02
	}

	return b.Marshal()
}

func (f HEADERS) Marshal() []byte {
	b := base{}
	b.Type = 0x1

	flagHeaders := make([]byte, 0, 5)
	if f.Flags.PRIORITY_GROUP {
		flagHeaders = flagHeaders[0:5]
		binary.BigEndian.PutUint32(flagHeaders, f.PriorityGroupIdentifier)
		flagHeaders[0] |= 0x80
		flagHeaders[4] = f.Weight
		b.Flags |= 0x20

	} else if f.Flags.PRIORITY_DEPENDENCY {
		flagHeaders = flagHeaders[0:4]
		binary.BigEndian.PutUint32(flagHeaders, f.StreamDependency)
		flagHeaders[0] |= 0x80

		b.Flags |= 0x40
	}

	payload := paddingHeaders(&b, f.Padding)
	payload = append(payload, flagHeaders...)
	payload = append(payload, f.HeaderBlockFragment...)
	payload = append(payload, f.Padding...)

	b.Payload = string(payload)

	if f.Flags.END_STREAM {
		b.Flags |= 0x01
	}
	if f.Flags.END_SEGMENT {
		b.Flags |= 0x02
	}
	if f.Flags.END_HEADERS {
		b.Flags |= 0x04
	}

	return b.Marshal()
}

func (f SETTINGS) Marshal() []byte {
	b := base{}
	b.Type = 0x4

	payload := make([]byte, len(f.Parameters)*5)
	for i, parameter := range f.Parameters {
		payload[i*5] = parameter.Identifier
		binary.BigEndian.PutUint32(payload[i*5+1:], parameter.Value)
	}
	b.Payload = string(payload)

	if f.Flags.ACK {
		b.Flags |= 0x1
	}

	return b.Marshal()
}

func Unmarshal(wire *[]byte) (Frame, error) {
	// TODO: validation that all this is well formed.
	payloadLen  := binary.BigEndian.Uint16([]byte{(*wire)[0] & 0x3F, (*wire)[1]})
	frameType := (*wire)[2]
	frameFlags := (*wire)[3]
	streamIdentifier := binary.BigEndian.Uint32([]byte{
		(*wire)[4] & 0x7F,
		(*wire)[5],
		(*wire)[6],
		(*wire)[7],
	})
	*wire = (*wire)[8:]
	toDecode := (*wire)[0:payloadLen]
	*wire = (*wire)[payloadLen:]

	var err error
	var f Frame

	switch frameType {
	case 0x0:
		f, err = unmarshalDataPayload(frameFlags, streamIdentifier, string(toDecode))
	case 0x6:
		// TODO: if stream id is set, send PROTOCOL_ERROR
		// TODO: if length field is not 8, send FRAME_SIZE_ERROR
		f, err = unmarshalPingPayload(frameFlags, string(toDecode))
	}

	if err != nil {
		return nil, err
	}
	return f, nil
}

func flagIsSet(flags uint8, mask uint8) bool {
	return flags & mask == mask
}

func unmarshalPingPayload(frameFlags uint8, payload string) (Frame, error) {
	f := PING{}
	f.OpaqueData = binary.BigEndian.Uint64([]byte(payload))

	return f, nil
}

func unmarshalDataPayload(frameFlags uint8, streamIdentifier uint32, payload string) (Frame, error) {
	// Check flags for pad high/pad low

	f := DATA{}

	paddingLengthBytes := []byte{0x00, 0x00}
	if flagIsSet(frameFlags, 0x1) {
		f.Flags.END_STREAM = true
	}
	if flagIsSet(frameFlags, 0x2) {
		f.Flags.END_SEGMENT = true
	}
	if flagIsSet(frameFlags, 0x10) {
		// padHigh is present
		paddingLengthBytes[0] = payload[0]
		payload = payload[1:]
	}
	if flagIsSet(frameFlags, 0x08) {
		// padLow is present
		paddingLengthBytes[1] = payload[0]
		payload = payload[1:]
	}
	paddingLength := binary.BigEndian.Uint16(paddingLengthBytes)

	dataLength := uint16(len(payload)) - paddingLength
	f.Data = payload[0:dataLength]

	payload = payload[dataLength:]
	f.Padding = payload[0:paddingLength]
	f.StreamIdentifier = streamIdentifier

	return f, nil
}
