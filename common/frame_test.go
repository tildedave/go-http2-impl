package common

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestMarshalEmptyFrame(t *testing.T) {
	f := Marshal(Frame{ 0, 0, 0, ""})

	f[0] & 192
}

func TestMarshalFrameWithPayload(t *testing.T) {
	f := Marshal(Frame{ 0, 0, 0, ""})

	f[0] & 192
}
