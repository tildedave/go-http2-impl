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
