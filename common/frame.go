package common

import (
	"encoding/binary"
)

// baseFrame is used internally to marshal other types of frames
type baseFrame struct {
	Type uint8
	Flags uint8
	StreamIdentifier uint32
	Payload []byte
}

type DataFrame struct {
	PadHigh uint8
	PadLow uint8
}

type GOAWAYFrame struct {
	LastStreamId uint32
	ErrorCode uint32
	AdditionalDebugData string
}

type Frame interface {
	Marshal() []byte
	Length() uint32
}


func (f baseFrame) Marshal() []byte {
	length := make([]byte, 2)
	binary.BigEndian.PutUint16(length, uint16(len(f.Payload)))

	streamIdentifier := make([]byte, 4)
	binary.BigEndian.PutUint32(streamIdentifier, f.StreamIdentifier)

	return append(append(length, f.Type, f.Flags),
		      streamIdentifier...)
}
