package hpack

// http://tools.ietf.org/html/draft-ietf-httpbis-header-compression-07

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDecodeSingleHeader(t *testing.T)  {
	context := NewEncodingContext()

	headers := context.Decode("\x82")

	assert.Equal(t, len(headers.Headers), 1)
	assert.Equal(t, headers.Headers[0], HeaderField{":method", "GET"})
}

func TestDecodeMultipleHeaders(t *testing.T)  {
	context := NewEncodingContext()

	headers := context.Decode("\x82\x87")

	assert.Equal(t, len(headers.Headers), 2)
	assert.Equal(t, headers.Headers[0], HeaderField{":method", "GET"})
	assert.Equal(t, headers.Headers[1], HeaderField{":scheme", "http"})
}

func TestDecodeLiteralHeaderWithIndexedName(t *testing.T) {
	context := NewEncodingContext()

	headers := context.Decode("\x82\x87\x86\x44\x0fwww.example.com")

	assert.Equal(t, len(headers.Headers), 4)
	assert.Equal(t, headers.Headers[3], HeaderField{":authority", "www.example.com"})
}

func TestDecodeLiteralHeaderWithLiteralName(t *testing.T) {
	context := NewEncodingContext()

	headers := context.Decode("\x40\x0dcustom-header\x0apuppy-dogs")

	assert.Equal(t, len(headers.Headers), 1)
	assert.Equal(t, headers.Headers[0], HeaderField{"custom-header", "puppy-dogs"})
}
