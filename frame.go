package main

import (
	"encoding/binary"
	"fmt"
)

var _ = fmt.Printf // package fmt is now used

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
	Data             string
	Padding          string

	Flags struct {
		END_STREAM  bool
		END_SEGMENT bool
	}
}

// http://tools.ietf.org/html/draft-ietf-httpbis-http2-11#section-6.2
type HEADERS struct {
	StreamIdentifier        uint32
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
		EXCLUSIVE           bool
	}
}

// http://tools.ietf.org/html/draft-ietf-httpbis-http2-11#section-6.3
type PRIORITY struct {
	PriorityGroupIdentifier uint32
	Weight                  uint8
	StreamDependency        uint32
	Flags                   struct {
		PRIORITY_GROUP      bool
		PRIORITY_DEPENDENCY bool
		EXCLUSIVE           bool
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
	binary.BigEndian.PutUint32(header[4:8], f.StreamIdentifier)

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
	b.StreamIdentifier = f.StreamIdentifier

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
		binary.BigEndian.PutUint32(payload, f.PriorityGroupIdentifier)
		payload[4] = f.Weight
	}
	b.Payload = string(payload)

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
	payloadLen := binary.BigEndian.Uint16([]byte{(*wire)[0] & 0x3F, (*wire)[1]})
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
		if streamIdentifier == 0 {
			return nil, ConnectionError{
				PROTOCOL_ERROR,
				"Data payload must have stream identifier",
			}
		}
		f, err = unmarshalDataPayload(frameFlags, streamIdentifier, string(toDecode))
	case 0x1:
		if streamIdentifier == 0 {
			return nil, ConnectionError{
				PROTOCOL_ERROR,
				"Headers payload must have stream identifier",
			}
		}
		f, err = unmarshalHeadersPayload(frameFlags, streamIdentifier, string(toDecode))
	case 0x4:
		f, err = unmarshalSettingsPayload(frameFlags, string(toDecode))
	case 0x6:
		if streamIdentifier != 0 {
			return nil, ConnectionError{
				PROTOCOL_ERROR,
				"Ping payload must not have stream identifier",
			}
		}
		if payloadLen != 8 {
			return nil, ConnectionError{
				FRAME_SIZE_ERROR,
				"Ping payload must have length of 8",
			}
		}
		f, err = unmarshalPingPayload(frameFlags, string(toDecode))
	case 0x7:
		f, err = unmarshalGoAwayPayload(string(toDecode))
	}

	if err != nil {
		return nil, err
	}
	return f, nil
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

func unmarshalDataPayload(frameFlags uint8, streamIdentifier uint32, payload string) (Frame, error) {
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
	f.StreamIdentifier = streamIdentifier

	return f, nil
}

func unmarshalGoAwayPayload(payload string) (Frame, error) {
	lastStreamId := binary.BigEndian.Uint32([]byte{
		payload[0] & 0x7F,
		payload[1],
		payload[2],
		payload[3],
	})

	return GOAWAY{
		LastStreamId:        lastStreamId,
		ErrorCode:           binary.BigEndian.Uint32([]byte(payload[4:8])),
		AdditionalDebugData: payload[8:],
	}, nil
}

func unmarshalHeadersPayload(frameFlags uint8, streamIdentifier uint32, payload string) (Frame, error) {
	f := HEADERS{}
	f.StreamIdentifier = streamIdentifier

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
		f.PriorityGroupIdentifier = binary.BigEndian.Uint32([]byte{
			payload[0] & 0x7F,
			payload[1],
			payload[2],
			payload[3],
		})
		f.Weight = payload[4]
		f.Flags.PRIORITY_GROUP = true
		payload = payload[5:]
	}
	if flagIsSet(frameFlags, 0x40) {
		// Priority dependency fields are present
		f.StreamDependency = binary.BigEndian.Uint32([]byte{
			payload[0] & 0x7F,
			payload[1],
			payload[2],
			payload[3],
		})
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
