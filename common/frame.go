package common

import (
	"encoding/binary"
)

// baseFrame is used internally to marshal other types of frames
type baseFrame struct {
	Type uint8
	Flags uint8
	StreamIdentifier uint32
	Payload string
}

// http://tools.ietf.org/html/draft-ietf-httpbis-http2-11#section-6.1
type DATA struct {
	Data string
	Padding string

	Flags struct {
		END_STREAM bool
		END_SEGMENT bool
	}
}

// http://tools.ietf.org/html/draft-ietf-httpbis-http2-11#section-6.2
type HEADERS struct {
	PriorityGroupIdentifier uint32
	Weight uint8
	StreamDependency uint32
	HeaderBlockFragment string
	Padding string

	Flags struct {
		END_STREAM bool
		END_SEGMENT bool
		END_HEADERS bool
		PRIORITY_GROUP bool
		PRIORITY_DEPENDENCY bool
	}
}

// http://tools.ietf.org/html/draft-ietf-httpbis-http2-11#page-35
type SETTINGS_Parameter struct {
	Identifier uint8
	Value uint32
}

type SETTINGS struct {
	Parameters []SETTINGS_Parameter
}

// http://tools.ietf.org/html/draft-ietf-httpbis-http2-11#section-6.7
type PING struct {
	OpaqueData uint64
	Flags struct {
		ACK bool
	}
}

// http://tools.ietf.org/html/draft-ietf-httpbis-http2-11#section-6.8
type GOAWAY struct {
	LastStreamId uint32
	ErrorCode uint32
	AdditionalDebugData string
}

type Frame interface {
	Marshal() string
}

func (f baseFrame) Marshal() []byte {
	header := make([]byte, 8)
	binary.BigEndian.PutUint16(header, uint16(len(f.Payload)))
	header[2] = f.Type
	header[3] = f.Flags
	binary.BigEndian.PutUint32(header[4:8], f.StreamIdentifier)

	return append(header, f.Payload...)
}

func (f GOAWAY) Marshal() []byte {
	bf := baseFrame{}
	bf.Type = 0x7

	payload := make([]byte, 8 + len(f.AdditionalDebugData))
	binary.BigEndian.PutUint32(payload[0:4], f.LastStreamId)
	binary.BigEndian.PutUint32(payload[4:8], f.ErrorCode)
	copy(payload, f.AdditionalDebugData)

	bf.Payload = string(payload)

	return bf.Marshal()
}

func (f PING) Marshal() []byte {
	bf := baseFrame{}
	bf.Type = 0x6
	if (f.Flags.ACK) {
		bf.Flags = 0x1
	}

	payload := make([]byte, 8)
	binary.BigEndian.PutUint64(payload, f.OpaqueData)
	bf.Payload = string(payload)

	return bf.Marshal()
}

func paddingHeaders(bf* baseFrame, padding string) []byte {
	paddingHeaders := make([]byte, 0, 2)
	paddingLength := uint16(len(padding))
	if paddingLength > 0 {
		// set PADDING_LOW flag
		bf.Flags |= 0x08

		if paddingLength > 256 {
			// set PADDING_HIGH flag
			bf.Flags |= 0x10
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
	bf := baseFrame{}
	bf.Type = 0x0

	payload := paddingHeaders(&bf, f.Padding)
	payload = append(payload, f.Data...)
	payload = append(payload, f.Padding...)
	bf.Payload = string(payload)

	if (f.Flags.END_STREAM) {
		bf.Flags |= 0x01
	}
	if (f.Flags.END_SEGMENT) {
		bf.Flags |= 0x02
	}

	return bf.Marshal()
}

func (f HEADERS) Marshal() []byte {
	bf := baseFrame{}
	bf.Type = 0x1

	flagHeaders := make([]byte, 0, 5)
	if f.Flags.PRIORITY_GROUP {
		flagHeaders = flagHeaders[0:5]
		binary.BigEndian.PutUint32(flagHeaders, f.PriorityGroupIdentifier)
		flagHeaders[0] |= 0x80
		flagHeaders[4] = f.Weight
		bf.Flags |= 0x20

	} else if f.Flags.PRIORITY_DEPENDENCY {
		flagHeaders = flagHeaders[0:4]
		binary.BigEndian.PutUint32(flagHeaders, f.StreamDependency)
		flagHeaders[0] |= 0x80

		bf.Flags |= 0x40
	}

	payload := paddingHeaders(&bf, f.Padding)
	payload = append(payload, flagHeaders...)
	payload = append(payload, f.HeaderBlockFragment...)
	payload = append(payload, f.Padding...)

	bf.Payload = string(payload)

	if (f.Flags.END_STREAM) {
		bf.Flags |= 0x01
	}
	if (f.Flags.END_SEGMENT) {
		bf.Flags |= 0x02
	}
	if (f.Flags.END_HEADERS) {
		bf.Flags |= 0x04
	}

	return bf.Marshal()
}
