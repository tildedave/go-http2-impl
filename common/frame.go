package common

import (
	"encoding/binary"
)

type Payload struct {
	Data string
}

type Frame struct {
	Type uint8
	Flags uint8
	StreamIdentifier uint32
	Payload Payload
}

func Marshal(f Frame) []byte {
	length := make([]byte, 2)
	binary.BigEndian.PutUint16(length, uint16(len(f.Payload.Data)))

	streamIdentifier := make([]byte, 4)
	binary.BigEndian.PutUint32(streamIdentifier, f.StreamIdentifier)

	return append(append(length, f.Type, f.Flags),
		      streamIdentifier...)
}
