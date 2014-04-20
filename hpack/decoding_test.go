package hpack

// http://tools.ietf.org/html/draft-ietf-httpbis-header-compression-07

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDecodeSingleHeader(t *testing.T)  {
	context := NewEncodingContext()

	headers, _ := context.Decode("\x82")

	assert.Equal(t, len(headers.Headers), 1)
	assert.Equal(t, headers.Headers[0], HeaderField{":method", "GET"})
}

func TestDecodeMultipleHeaders(t *testing.T)  {
	context := NewEncodingContext()

	headers, _ := context.Decode("\x82\x87")

	assert.Equal(t, len(headers.Headers), 2)
	assert.Equal(t, headers.Headers[0], HeaderField{":method", "GET"})
	assert.Equal(t, headers.Headers[1], HeaderField{":scheme", "http"})
}

func TestDecodeLiteralHeaderWithIndexedName(t *testing.T) {
	context := NewEncodingContext()

	headers, _ := context.Decode("\x82\x87\x86\x44\x0fwww.example.com")

	assert.Equal(t, len(headers.Headers), 4)
	assert.Equal(t, headers.Headers[3], HeaderField{":authority", "www.example.com"})
}

func TestDecodeLiteralHeaderWithLiteralName(t *testing.T) {
	context := NewEncodingContext()

	headers, _ := context.Decode("\x40\x0dcustom-header\x0apuppy-dogs")

	assert.Equal(t, len(headers.Headers), 1)
	assert.Equal(t, headers.Headers[0], HeaderField{"custom-header", "puppy-dogs"})
}

func TestDecodeWithReferenceSet(t *testing.T) {
	context := NewEncodingContext()

	// :method: GET
	context.Decode("\x82")
	// :scheme: http
	headers, _ := context.Decode("\x87")

	assert.Equal(t, len(headers.Headers), 2)
	assert.Equal(t, headers.Headers[0], HeaderField{ ":scheme", "http" })
	assert.Equal(t, headers.Headers[1], HeaderField{ ":method", "GET" })
	assert.Equal(t, len(context.HeaderTable.Entries), 2)
}

func TestDecodeWithReferenceSetEmpty(t *testing.T) {
	context := NewEncodingContext()

	// :method: GET
	context.Decode("\x82")
	// :scheme: http
	headers, _ := context.Decode("\x30\x87")

	assert.Equal(t, len(headers.Headers), 1)
	assert.Equal(t, headers.Headers[0], HeaderField{ ":scheme", "http" })
	assert.Equal(t, len(context.HeaderTable.Entries), 2)
}

func TestDecodeWithLiteralFieldWithoutIndexing(t *testing.T) {
	context := NewEncodingContext()

	headers, _ := context.Decode("\x04\x0c\x2f\x73\x61\x6d\x70\x6c\x65\x2f\x70\x61\x74\x68")

	assert.Equal(t, len(headers.Headers), 1)
	assert.Equal(t, headers.Headers[0], HeaderField{ ":path", "/sample/path"})
	assert.Equal(t, len(context.HeaderTable.Entries), 0)
}
