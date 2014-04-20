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

	assert.Equal(t, len(encodedInteger), 3)
	assert.Equal(t, byte(encodedInteger[0]) & 0x1f, byte(0x1f))
	assert.Equal(t, byte(encodedInteger[1]), byte(0x9a))
	assert.Equal(t, byte(encodedInteger[2]), byte(0x0a))
}

func TestDecodeIntegerThatFitsInPrefix(t *testing.T) {
	wire := []byte{0x0a}
	decodedInteger := decodeInteger(&wire, 5)

	assert.Equal(t, decodedInteger, uint(10))
	assert.Equal(t, wire, []byte{})
}

func TestDecodeIntegerThatOverflowsPrefix(t *testing.T) {
	wire := []byte{0x1f, 0x03}
	decodedInteger := decodeInteger(&wire, 5)

	assert.Equal(t, decodedInteger, uint(34))
	assert.Equal(t, wire, []byte{})
}

func TestDecodeIntegerThatOverflowsPrefixTwice(t *testing.T) {
	wire := []byte{0x1f, 0x9a, 0x0a}
	decodedInteger := decodeInteger(&wire, 5)

	assert.Equal(t, decodedInteger, uint(1337))
	assert.Equal(t, wire, []byte{})
}
