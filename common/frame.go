package common

import (
	"encoding/binary"
)

// base is used internally to marshal other types of frames
type base struct {
	Type             uint8
	Flags            uint8
	StreamIdentifier uint32
	Payload          string
}

// http://tools.ietf.org/html/draft-ietf-httpbis-http2-11#section-6.1
type DATA struct {
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
	Marshal() string
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
