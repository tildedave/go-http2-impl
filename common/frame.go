package common

import (
	"encoding/binary"
)

type Frame struct {
	Type uint8
	Flags uint8
	StreamIdentifier uint32
	Payload Marshaller
}

type Data_Frame struct {
	PadHigh uint8
	PadLow uint8
}

type GOAWAY_Frame struct {
	LastStreamId uint32
	ErrorCode uint32
	AdditionalDebugData string
}

type Marshaller interface {
	Marshal() []byte
}

func (f Frame) Marshal() []byte {
	payload := f.Payload.Marshal()

	length := make([]byte, 2)
	binary.BigEndian.PutUint16(length, uint16(len(payload)))

	streamIdentifier := make([]byte, 4)
	binary.BigEndian.PutUint32(streamIdentifier, f.StreamIdentifier)

	return append(append(length, f.Type, f.Flags),
		      streamIdentifier...)
}
