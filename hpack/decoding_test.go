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

func TestDecodeWithLiteralFieldWithNeverIndexing(t *testing.T) {
	context := NewEncodingContext()

	headers, _ := context.Decode("\x14\x0c\x2f\x73\x61\x6d\x70\x6c\x65\x2f\x70\x61\x74\x68")

	assert.Equal(t, len(headers.Headers), 1)
	assert.Equal(t, headers.Headers[0], HeaderField{ ":path", "/sample/path"})
	assert.Equal(t, len(context.HeaderTable.Entries), 0)
}

func TestDecodeRemovesHeadersBasedOnDirective(t *testing.T) {
	var headers HeaderSet

	context := NewEncodingContext()

	headers, _ = context.Decode("\x82\x87")
	assert.Equal(t, len(headers.Headers), 2)

	t.Log(context.ReferenceSet, context.HeaderTable)

	headers, _ = context.Decode("\x81\x89")
	assert.Equal(t, len(headers.Headers), 2)

	t.Log(headers)
	t.Log(context)
}

func TestDecodeWithEviction(t *testing.T) {
	var headers HeaderSet

	context := NewEncodingContext()
	context.HeaderTable.MaxSize = 256

	headers, _ = context.Decode("\x48\x03\x33\x30\x32\x59\x07\x70\x72\x69\x76\x61\x74\x65\x63\x1d\x4d\x6f\x6e\x2c\x20\x32\x31\x20\x4f\x63\x74\x20\x32\x30\x31\x33\x20\x32\x30\x3a\x31\x33\x3a\x32\x31\x20\x47\x4d\x54\x71\x17\x68\x74\x74\x70\x73\x3a\x2f\x2f\x77\x77\x77\x2e\x65\x78\x61\x6d\x70\x6c\x65\x2e\x63\x6f\x6d")

	assert.Equal(t, len(headers.Headers), 4)
	assert.Equal(t, headers.Headers, []HeaderField{
		HeaderField{":status", "302" },
		HeaderField{"cache-control", "private"},
		HeaderField{"date", "Mon, 21 Oct 2013 20:13:21 GMT"},
		HeaderField{"location", "https://www.example.com"},
	})

	headers, _ = context.Decode("\x84\x8c")

	assert.Equal(t, len(headers.Headers), 4)
	assert.Equal(t, headers.Headers, []HeaderField{
		HeaderField{":status", "200" },
		HeaderField{"cache-control", "private"},
		HeaderField{"date", "Mon, 21 Oct 2013 20:13:21 GMT"},
		HeaderField{"location", "https://www.example.com"},
	})

	headers, _ = context.Decode("\x84\x84\x43\x1d\x4d\x6f\x6e\x2c\x20\x32\x31\x20\x4f\x63\x74\x20\x32\x30\x31\x33\x20\x32\x30\x3a\x31\x33\x3a\x32\x32\x20\x47\x4d\x54\x5e\x04\x67\x7a\x69\x70\x84\x84\x83\x83\x7b\x38\x66\x6f\x6f\x3d\x41\x53\x44\x4a\x4b\x48\x51\x4b\x42\x5a\x58\x4f\x51\x57\x45\x4f\x50\x49\x55\x41\x58\x51\x57\x45\x4f\x49\x55\x3b\x20\x6d\x61\x78\x2d\x61\x67\x65\x3d\x33\x36\x30\x30\x3b\x20\x76\x65\x72\x73\x69\x6f\x6e\x3d\x31")

	assert.Equal(t, len(headers.Headers), 6)
	assert.Equal(t, headers.Headers, []HeaderField{
		HeaderField{"cache-control", "private"},
		HeaderField{"date", "Mon, 21 Oct 2013 20:13:22 GMT"},
		HeaderField{"content-encoding", "gzip"},
		HeaderField{"location", "https://www.example.com"},
		HeaderField{":status", "200" },
		HeaderField{"set-cookie", "foo=ASDJKHQKBZXOQWEOPIUAXQWEOIU; max-age=3600; version=1"},
	})
}
