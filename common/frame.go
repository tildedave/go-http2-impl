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

type PingFrame struct {
	OpaqueData uint64
	Flags struct {
		ACK bool
	}
}

// http://tools.ietf.org/html/draft-ietf-httpbis-http2-11#section-6.1
type DataFrame struct {
	Data string
	Padding string

	Flags struct {
		END_STREAM bool
		END_SEGMENT bool
	}
}

// http://tools.ietf.org/html/draft-ietf-httpbis-http2-11#section-6.2
type HeadersFrame struct {

	PriorityGroupIdentifier uint32
	Weight uint8
	StreamDependency uint32
	HeaderBlockFragment string
	Padding string

	// Not sure what E and R bit flags are for -- seems to be covered
	// by header flag already?

	Flags struct {
		END_STREAM bool
		END_SEGMENT bool
		END_HEADERS bool
		PRIORITY_GROUP bool
		PRIORITY_DEPENDENCY bool
	}
}

// http://tools.ietf.org/html/draft-ietf-httpbis-http2-11#section-6.8
type GOAWAYFrame struct {
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

func (f GOAWAYFrame) Marshal() []byte {
	bf := baseFrame{}
	bf.Type = 0x7

	payload := make([]byte, 8 + len(f.AdditionalDebugData))
	binary.BigEndian.PutUint32(payload[0:4], f.LastStreamId)
	binary.BigEndian.PutUint32(payload[4:8], f.ErrorCode)
	copy(payload, f.AdditionalDebugData)

	bf.Payload = string(payload)

	return bf.Marshal()
}

func (f PingFrame) Marshal() []byte {
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

func (f DataFrame) Marshal() []byte {
	bf := baseFrame{}
	bf.Type = 0x0

	paddingLength := uint16(len(f.Padding))

	payload := make([]byte, 2 + len(f.Data) + len(f.Padding))
	binary.BigEndian.PutUint16(payload, paddingLength)

	copy(payload[2:], []byte(f.Data))
	copy(payload[2 + len(f.Data):], []byte(f.Padding))
	bf.Payload = string(payload)

	if paddingLength > 0 {
		// set PADDING_LOW flag
		bf.Flags |= 0x08

		if paddingLength > 256 {
			// set PADDING_HIGH flag
			bf.Flags |= 0x10
		}
	}

	if (f.Flags.END_STREAM) {
		bf.Flags |= 0x01
	}

	if (f.Flags.END_SEGMENT) {
		bf.Flags |= 0x02
	}

	return bf.Marshal()
}

func (f HeadersFrame) Marshal() []byte {
	bf := baseFrame{}
	bf.Type = 0x1

	return bf.Marshal()
}
