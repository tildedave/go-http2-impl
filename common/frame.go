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
	AdditionalDebugData []byte
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

	marshalled_f := append(length, f.Type, f.Flags)
	marshalled_f = append(marshalled_f, streamIdentifier...)
	marshalled_f = append(marshalled_f, f.Payload...)

	return marshalled_f
}

func (f GOAWAYFrame) Marshal() []byte {
	bf := baseFrame{}
	bf.Type = 0x7

	lastStreamId := make([]byte, 8)
	binary.BigEndian.PutUint32(lastStreamId, f.LastStreamId)

	errorCode := make([]byte, 8)
	binary.BigEndian.PutUint32(errorCode, f.ErrorCode)

	payload := append(lastStreamId, errorCode...)
	payload = append(payload, f.AdditionalDebugData...)

	bf.Payload = payload

	return bf.Marshal()
}
