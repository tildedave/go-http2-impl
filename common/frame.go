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

type GOAWAYFrame struct {
	LastStreamId uint32
	ErrorCode uint32
	AdditionalDebugData string
}

type DataFrame struct {
	Data string
	Padding string

	Flags struct {
		END_STREAM bool
		END_SEGMENT bool
	}
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

	return bf.Marshal()
}
