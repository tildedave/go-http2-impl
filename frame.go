package main

import (
	"encoding/binary"
	"fmt"
)

var _ = fmt.Printf // package fmt is now used

// base is used internally to marshal other types of frames
type base struct {
	Type     uint8
	Flags    uint8
	StreamId uint32
	Payload  string
}

// http://tools.ietf.org/html/draft-ietf-httpbis-http2-11#section-6.1
type DATA struct {
	StreamId uint32
	Data     string
	Padding  string

	Flags struct {
		END_STREAM  bool
		END_SEGMENT bool
	}
}

// http://tools.ietf.org/html/draft-ietf-httpbis-http2-11#section-6.2
type HEADERS struct {
	StreamId            uint32
	PriorityGroupId     uint32
	Weight              uint8
	StreamDependency    uint32
	HeaderBlockFragment string
	Padding             string

	Flags struct {
		END_STREAM          bool
		END_SEGMENT         bool
		END_HEADERS         bool
		PRIORITY_GROUP      bool
		PRIORITY_DEPENDENCY bool
		EXCLUSIVE           bool
	}
}

// http://tools.ietf.org/html/draft-ietf-httpbis-http2-11#section-6.3
type PRIORITY struct {
	StreamId         uint32
	PriorityGroupId  uint32
	Weight           uint8
	StreamDependency uint32
	Flags            struct {
		PRIORITY_GROUP      bool
		PRIORITY_DEPENDENCY bool
		EXCLUSIVE           bool
	}
}

// http://tools.ietf.org/html/draft-ietf-httpbis-http2-11#section-6.4
type RST_STREAM struct {
	StreamId  uint32
	ErrorCode uint32
}

type PUSH_PROMISE struct {
	StreamId            uint32
	PromisedStreamId    uint32
	HeaderBlockFragment string
	Padding             string
	Flags               struct {
		END_HEADERS bool
	}
}

const (
	SETTINGS_HEADER_TABLE_SIZE      = 1
	SETTINGS_ENABLE_PUSH            = 2
	SETTINGS_MAX_CONCURRENT_STREAMS = 3
	SETTINGS_INITIAL_WINDOW_SIZE    = 4
)

// http://tools.ietf.org/html/draft-ietf-httpbis-http2-11#page-35
type Parameter struct {
	Id    uint8
	Value uint32
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

// http://tools.ietf.org/html/draft-ietf-httpbis-http2-11#section-6.9
type WINDOW_UPDATE struct {
	StreamId            uint32
	WindowSizeIncrement uint32
}

// http://tools.ietf.org/html/draft-ietf-httpbis-http2-11#section-6.10
type CONTINUATION struct {
	StreamId            uint32
	HeaderBlockFragment string
	Padding             string
	Flags               struct {
		END_HEADERS bool
	}
}

// TODO: ALTSVC frame

type Frame interface {
	Marshal() []byte
}

const (
	NO_ERROR            = 0
	PROTOCOL_ERROR      = 1
	INTERNAL_ERROR      = 2
	FLOW_CONTROL_ERROR  = 3
	SETTINGS_TIMEOUT    = 4
	STREAM_CLOSED       = 5
	FRAME_SIZE_ERROR    = 6
	REFUSED_STREAM      = 7
	CANCEL              = 8
	COMPRESSION_ERROR   = 9
	CONNECT_ERROR       = 10
	ENHANCE_YOUR_CALM   = 11
	INADEQUATE_SECURITY = 12
)

type ConnectionError struct {
	Code    uint8
	Message string
}

func (e ConnectionError) Error() string {
	return fmt.Sprintf("ConnectionError: %s (%d)", e.Message, e.Code)
}

func (f base) Marshal() []byte {
	header := make([]byte, 8)
	binary.BigEndian.PutUint16(header, uint16(len(f.Payload)))
	header[2] = f.Type
	header[3] = f.Flags
	binary.BigEndian.PutUint32(header[4:8], f.StreamId)

	return append(header, f.Payload...)
}

func (f GOAWAY) Marshal() []byte {
	b := base{}
	b.Type = 0x7

	payload := make([]byte, 8+len(f.AdditionalDebugData))
	binary.BigEndian.PutUint32(payload[0:4], f.LastStreamId)
	binary.BigEndian.PutUint32(payload[4:8], f.ErrorCode)
	copy(payload[8:], f.AdditionalDebugData)

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
	b.StreamId = f.StreamId

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
	b.StreamId = f.StreamId

	flagHeaders := make([]byte, 0, 5)
	if f.Flags.PRIORITY_GROUP {
		flagHeaders = flagHeaders[0:5]
		binary.BigEndian.PutUint32(flagHeaders, f.PriorityGroupId)
		flagHeaders[0] |= 0x80
		flagHeaders[4] = f.Weight
		b.Flags |= 0x20

	} else if f.Flags.PRIORITY_DEPENDENCY {
		flagHeaders = flagHeaders[0:4]
		binary.BigEndian.PutUint32(flagHeaders, f.StreamDependency)
		if f.Flags.EXCLUSIVE {
			flagHeaders[0] |= 0x80
		}

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

func (f PRIORITY) Marshal() []byte {
	b := base{}
	b.Type = 0x2
	b.StreamId = f.StreamId

	var payload []byte
	if f.Flags.PRIORITY_DEPENDENCY {
		b.Flags |= 0x40
		payload = make([]byte, 4)
		binary.BigEndian.PutUint32(payload, f.StreamDependency)
		if f.Flags.EXCLUSIVE {
			payload[0] |= 0x80
		}
	} else if f.Flags.PRIORITY_GROUP {
		b.Flags |= 0x20
		payload = make([]byte, 5)
		binary.BigEndian.PutUint32(payload, f.PriorityGroupId)
		payload[4] = f.Weight
	}
	b.Payload = string(payload)

	return b.Marshal()
}

func (f RST_STREAM) Marshal() []byte {
	b := base{}
	b.Type = 0x3
	b.StreamId = f.StreamId

	payload := make([]byte, 4)
	binary.BigEndian.PutUint32(payload, f.ErrorCode)
	b.Payload = string(payload)

	return b.Marshal()
}

func (f SETTINGS) Marshal() []byte {
	b := base{}
	b.Type = 0x4

	payload := make([]byte, len(f.Parameters)*5)
	for i, parameter := range f.Parameters {
		payload[i*5] = parameter.Id
		binary.BigEndian.PutUint32(payload[i*5+1:], parameter.Value)
	}
	b.Payload = string(payload)

	if f.Flags.ACK {
		b.Flags |= 0x1
	}

	return b.Marshal()
}

func (f PUSH_PROMISE) Marshal() []byte {
	b := base{}
	b.Type = 0x5
	b.StreamId = f.StreamId

	if f.Flags.END_HEADERS {
		b.Flags |= 0x4
	}

	headers := paddingHeaders(&b, f.Padding)
	payload := make([]byte, 4+len(f.HeaderBlockFragment)+len(f.Padding))
	binary.BigEndian.PutUint32(payload[0:4], f.PromisedStreamId)
	copy(payload[4:4+len(f.HeaderBlockFragment)], f.HeaderBlockFragment)
	copy(payload[4+len(f.HeaderBlockFragment):], f.Padding)

	b.Payload = string(append(headers, payload...))

	return b.Marshal()
}

func (f WINDOW_UPDATE) Marshal() []byte {
	b := base{}
	b.Type = 0x8
	b.StreamId = f.StreamId

	payload := make([]byte, 4)
	binary.BigEndian.PutUint32(payload, f.WindowSizeIncrement&0x7FFFFFFF)

	b.Payload = string(payload)

	return b.Marshal()
}

func (f CONTINUATION) Marshal() []byte {
	b := base{}
	b.Type = 0x9
	b.StreamId = f.StreamId

	if f.Flags.END_HEADERS {
		b.Flags |= 0x4
	}

	headers := paddingHeaders(&b, f.Padding)
	payload := make([]byte, len(f.HeaderBlockFragment)+len(f.Padding))
	copy(payload[0:len(f.HeaderBlockFragment)], f.HeaderBlockFragment)
	copy(payload[len(f.HeaderBlockFragment):], f.Padding)

	b.Payload = string(append(headers, payload...))

	return b.Marshal()
}

func Unmarshal(wire []byte) (advance int, f Frame, err error) {
	if len(wire) < 8 {
		// Incomplete header
		return 0, nil, nil
	}
	payloadLen := binary.BigEndian.Uint16([]byte{wire[0] & 0x3F, wire[1]})
	frameType := wire[2]
	frameFlags := wire[3]
	streamId := uint31(string(wire[4:8]))

	if uint16(len(wire)) < payloadLen+8 {
		// Incomplete payload
		return 0, nil, nil
	}

	advance = int(payloadLen + 8)
	toDecode := string(wire[8:advance])

	switch frameType {
	case 0x0:
		if streamId == 0 {
			return advance, nil, ConnectionError{
				PROTOCOL_ERROR,
				"DATA frame must have stream identifier",
			}
		}
		f, err = unmarshalDataPayload(frameFlags, streamId, toDecode)
	case 0x1:
		if streamId == 0 {
			return advance, nil, ConnectionError{
				PROTOCOL_ERROR,
				"HEADERS frame must have stream identifier",
			}
		}
		f, err = unmarshalHeadersPayload(frameFlags, streamId, toDecode)
	case 0x2:
		if streamId == 0 {
			return advance, nil, ConnectionError{
				PROTOCOL_ERROR,
				"PRIORITY frame must have stream identifier",
			}
		}
		f, err = unmarshalPriorityPayload(frameFlags, streamId, toDecode)
	case 0x3:
		if streamId == 0 {
			return advance, nil, ConnectionError{
				PROTOCOL_ERROR,
				"RST_STREAM frame must have stream identifier",
			}
		}
		f, err = unmarshalRstStreamPayload(streamId, toDecode)
	case 0x4:
		f, err = unmarshalSettingsPayload(frameFlags, toDecode)
	case 0x5:
		if streamId == 0 {
			return advance, nil, ConnectionError{
				PROTOCOL_ERROR,
				"PUSH_PROMISE frame must have stream identifier",
			}
		}
		f, err = unmarshalPushPromisePayload(frameFlags, streamId, toDecode)
	case 0x6:
		if streamId != 0 {
			return advance, nil, ConnectionError{
				PROTOCOL_ERROR,
				"PING frame must not have stream identifier",
			}
		}
		if payloadLen != 8 {
			return advance, nil, ConnectionError{
				FRAME_SIZE_ERROR,
				"PING payload must have length of 8",
			}
		}
		f, err = unmarshalPingPayload(frameFlags, toDecode)
	case 0x7:
		f, err = unmarshalGoAwayPayload(toDecode)
	case 0x8:
		f, err = unmarshalWindowUpdatePayload(streamId, toDecode)
	case 0x9:
		if streamId == 0 {
			return advance, nil, ConnectionError{
				PROTOCOL_ERROR,
				"CONTINUATION frame must have stream identifier",
			}
		}
		f, err = unmarshalContinuationPayload(frameFlags, streamId, toDecode)
	}

	if err != nil {
		return advance, nil, err
	}
	return advance, f, nil
}

func flagIsSet(flags uint8, mask uint8) bool {
	return flags&mask == mask
}

func unmarshalPingPayload(frameFlags uint8, payload string) (Frame, error) {
	f := PING{}
	f.OpaqueData = binary.BigEndian.Uint64([]byte(payload))
	if flagIsSet(frameFlags, 0x1) {
		f.Flags.ACK = true
	}

	return f, nil
}

func decodePaddingLength(frameFlags uint8, payload *string) (uint16, error) {
	paddingLengthBytes := []byte{0x00, 0x00}
	if flagIsSet(frameFlags, 0x10) {
		// padHigh is present
		paddingLengthBytes[0] = (*payload)[0]
		*payload = (*payload)[1:]
		if !flagIsSet(frameFlags, 0x08) {
			return 0, ConnectionError{PROTOCOL_ERROR, "PAD_HIGH was set but PAD_LOW was not set"}
		}
	}
	if flagIsSet(frameFlags, 0x08) {
		// padLow is present
		paddingLengthBytes[1] = (*payload)[0]
		*payload = (*payload)[1:]
	}
	paddingLength := binary.BigEndian.Uint16(paddingLengthBytes)

	if paddingLength > uint16(len(*payload)) {
		return 0, ConnectionError{PROTOCOL_ERROR, "Padding length exceeded length of payload"}
	}

	return paddingLength, nil
}

func unmarshalDataPayload(frameFlags uint8, streamId uint32, payload string) (Frame, error) {
	// Check flags for pad high/pad low

	f := DATA{}

	if flagIsSet(frameFlags, 0x1) {
		f.Flags.END_STREAM = true
	}
	if flagIsSet(frameFlags, 0x2) {
		f.Flags.END_SEGMENT = true
	}
	paddingLength, err := decodePaddingLength(frameFlags, &payload)
	if err != nil {
		return nil, err
	}

	dataLength := uint16(len(payload)) - paddingLength
	f.Data = payload[0:dataLength]

	payload = payload[dataLength:]
	f.Padding = payload[0:paddingLength]
	f.StreamId = streamId

	return f, nil
}

func unmarshalGoAwayPayload(payload string) (Frame, error) {
	lastStreamId := uint31(payload[0:4])

	return GOAWAY{
		LastStreamId:        lastStreamId,
		ErrorCode:           binary.BigEndian.Uint32([]byte(payload[4:8])),
		AdditionalDebugData: payload[8:],
	}, nil
}

func unmarshalHeadersPayload(frameFlags uint8, streamId uint32, payload string) (Frame, error) {
	f := HEADERS{}
	f.StreamId = streamId

	paddingLength, err := decodePaddingLength(frameFlags, &payload)
	if err != nil {
		return nil, err
	}

	if flagIsSet(frameFlags, 0x1) {
		f.Flags.END_STREAM = true
	}
	if flagIsSet(frameFlags, 0x2) {
		f.Flags.END_SEGMENT = true
	}
	if flagIsSet(frameFlags, 0x4) {
		f.Flags.END_HEADERS = true
	}
	if flagIsSet(frameFlags, 0x20) && flagIsSet(frameFlags, 0x40) {
		return nil, ConnectionError{
			PROTOCOL_ERROR,
			"Cannot set both PRIORITY_GROUP and PRIORITY_DEPENDENCY flags",
		}
	}
	if flagIsSet(frameFlags, 0x20) {
		// Priority group fields (priority group identifier and weight) are
		// present
		f.PriorityGroupId = uint31(payload[0:4])
		f.Weight = payload[4]
		f.Flags.PRIORITY_GROUP = true
		payload = payload[5:]
	}
	if flagIsSet(frameFlags, 0x40) {
		// Priority dependency fields are present
		f.StreamDependency = uint31(payload[0:4])
		f.Flags.PRIORITY_DEPENDENCY = true
		if flagIsSet(payload[0], 0x80) {
			f.Flags.EXCLUSIVE = true
		}
		payload = payload[4:]
	}

	payloadLength := uint16(len(payload)) - paddingLength

	f.HeaderBlockFragment = payload[0:payloadLength]
	f.Padding = payload[payloadLength:]

	return f, nil
}

func unmarshalPriorityPayload(frameFlags uint8, streamId uint32, payload string) (Frame, error) {
	f := PRIORITY{}
	f.StreamId = streamId

	if flagIsSet(frameFlags, 0x20) && flagIsSet(frameFlags, 0x40) {
		return nil, ConnectionError{
			PROTOCOL_ERROR,
			"Cannot set both PRIORITY_GROUP and PRIORITY_DEPENDENCY flags",
		}
	}
	if flagIsSet(frameFlags, 0x20) {
		f.Flags.PRIORITY_GROUP = true
		f.PriorityGroupId = uint31(payload[0:4])
		f.Weight = uint8(payload[4])
	}
	if flagIsSet(frameFlags, 0x40) {
		f.Flags.PRIORITY_DEPENDENCY = true
		f.StreamDependency = uint31(payload[0:4])
		if flagIsSet(payload[0], 0x80) {
			f.Flags.EXCLUSIVE = true
		}
	}
	return f, nil
}

func unmarshalRstStreamPayload(streamId uint32, payload string) (Frame, error) {
	f := RST_STREAM{}
	f.StreamId = streamId
	f.ErrorCode = binary.BigEndian.Uint32([]byte(payload))

	return f, nil
}

func unmarshalSettingsPayload(frameFlags uint8, payload string) (Frame, error) {
	f := SETTINGS{}
	if flagIsSet(frameFlags, 0x1) {
		f.Flags.ACK = true
		if len(payload) > 0 {
			return nil, ConnectionError{
				FRAME_SIZE_ERROR,
				"Payload of Settings frame with ACK flag must be empty",
			}
		}
	}

	for len(payload) > 0 {
		if len(payload) < 5 {
			return nil, ConnectionError{
				FRAME_SIZE_ERROR,
				"Improperly constructed Settings frame",
			}
		}
		id := payload[0]
		if id == 0 || id > 4 {
			return nil, ConnectionError{
				PROTOCOL_ERROR,
				fmt.Sprintf("Settings frame specified invalid identifier: %d", id),
			}
		}
		f.Parameters = append(f.Parameters, Parameter{
			id,
			binary.BigEndian.Uint32([]byte(payload[1:5])),
		})
		payload = payload[5:]
	}

	return f, nil
}

func unmarshalPushPromisePayload(frameFlags uint8, streamId uint32, payload string) (Frame, error) {
	f := PUSH_PROMISE{}
	f.StreamId = streamId
	if flagIsSet(frameFlags, 0x4) {
		f.Flags.END_HEADERS = true
	}

	paddingLength, err := decodePaddingLength(frameFlags, &payload)
	if err != nil {
		return nil, err
	}

	f.PromisedStreamId = uint31(payload[0:4])
	payload = payload[4:]
	headerBlockLength := uint16(len(payload)) - paddingLength
	f.HeaderBlockFragment = payload[0:headerBlockLength]
	f.Padding = payload[headerBlockLength:]

	return f, nil
}

func unmarshalWindowUpdatePayload(streamId uint32, payload string) (Frame, error) {
	f := WINDOW_UPDATE{}
	f.StreamId = streamId
	f.WindowSizeIncrement = uint31(payload)

	return f, nil
}

func unmarshalContinuationPayload(frameFlags uint8, streamId uint32, payload string) (Frame, error) {
	f := CONTINUATION{}
	f.StreamId = streamId
	if flagIsSet(frameFlags, 0x4) {
		f.Flags.END_HEADERS = true
	}

	paddingLength, err := decodePaddingLength(frameFlags, &payload)
	if err != nil {
		return nil, err
	}

	headerBlockLength := uint16(len(payload)) - paddingLength
	f.HeaderBlockFragment = payload[0:headerBlockLength]
	f.Padding = payload[headerBlockLength:]

	return f, nil
}

func uint31(payload string) uint32 {
	return binary.BigEndian.Uint32([]byte{
		payload[0] & 0x7F,
		payload[1],
		payload[2],
		payload[3],
	})
}
