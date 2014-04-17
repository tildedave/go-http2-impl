package common

import (
	"encoding/binary"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestMarshalEmptyFrame(t *testing.T) {
	f := Frame{ 0, 0, 0, ""}
	marshalled_f := Marshal(f)

	length := binary.BigEndian.Uint16(marshalled_f) & 0x7F

	assert.Equal(t, length, uint16(0),
		"Length should have been nothing")
}

func TestMarshalFrameWithPayloadIncludesLength(t *testing.T) {
	f := Frame{ 0, 0, 0, "this is the payload of the frame"}
	marshalled_f := Marshal(f)

	length := binary.BigEndian.Uint16(marshalled_f) & 0x7F

	assert.Equal(t, int(length), len(f.Payload),
		"Length field in header should have been the length of payload")
}
