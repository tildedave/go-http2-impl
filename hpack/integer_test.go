package hpack

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEncodeIntegerThatFitsInPrefix(t *testing.T) {
	encodedInteger := encodeInteger(10, 5)

	assert.Equal(t, len(encodedInteger), 1)
	assert.Equal(t, byte(encodedInteger[0]) & 0x1f, byte(0x0a))
}

func TestEncodeIntegerThatOverflowsPrefix(t *testing.T) {
	encodedInteger := encodeInteger(34, 5)

	assert.Equal(t, len(encodedInteger), 2)
	assert.Equal(t, byte(encodedInteger[0]) & 0x1f, byte(0x1f))
	assert.Equal(t, byte(encodedInteger[1]), byte(0x03))
}

func TestEncodeIntegerThatOverflowsPrefixTwice(t *testing.T) {
	encodedInteger := encodeInteger(1337, 5)

	t.Log([]byte(encodedInteger))

	assert.Equal(t, len(encodedInteger), 3)
	assert.Equal(t, byte(encodedInteger[0]) & 0x1f, byte(0x1f))
	assert.Equal(t, byte(encodedInteger[1]), byte(0x9a))
	assert.Equal(t, byte(encodedInteger[2]), byte(0x0a))
}
