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
	length := make([]byte, 2)
	binary.BigEndian.PutUint16(length, uint16(len(f.Payload)))

	streamIdentifier := make([]byte, 4)
	binary.BigEndian.PutUint32(streamIdentifier, f.StreamIdentifier)

	marshalled_f := append(length, f.Type, f.Flags)
	marshalled_f = append(marshalled_f, streamIdentifier...)
	marshalled_f = append(marshalled_f, f.Payload...)

	return marshalled_f
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
