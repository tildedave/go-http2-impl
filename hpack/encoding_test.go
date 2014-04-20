package hpack

// http://tools.ietf.org/html/draft-ietf-httpbis-header-compression-07

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEncodeHeaderFieldAddsToHeaderTable(t *testing.T) {
	context := NewEncodingContext()

	context.EncodeField(HeaderField{":method", "GET"})

	assert.Equal(t, context.HeaderTable.Entries[0],
		HeaderField{":method", "GET"})
}

func TestEncodeHeaderFieldUsesHeaderTableSizeAsOffset(t *testing.T) {
	// http://tools.ietf.org/html/draft-ietf-httpbis-header-compression-07#page-35
	context := NewEncodingContext()

	context.EncodeField(HeaderField{ ":method", "GET" })
	h := context.EncodeField(HeaderField{ ":scheme", "http" })

	assert.Equal(t, h, "\x87")
}

func TestEncodeHeaderFieldFromStaticTable(t *testing.T) {
	context := NewEncodingContext()

	var encoded string

	encoded = context.EncodeField(HeaderField{":method", "GET"})
	assert.Equal(t, encoded, "\x82")
	encoded = context.EncodeField(HeaderField{":path", "/"})
	assert.Equal(t, encoded, "\x85")

	assert.Equal(t, context.HeaderTable.Entries[0], HeaderField{":path", "/"})
	assert.Equal(t, context.HeaderTable.Entries[1], HeaderField{":method", "GET"})
}

func TestEncodeHeaderFieldWithNameAndLiteralValue(t *testing.T) {
	h := NewEncodingContext().EncodeField(HeaderField{":authority", "www.example.com"})

	assert.Equal(t, h, "\x41\x0f\x77\x77\x77\x2e\x65\x78\x61\x6d\x70\x6c\x65\x2e\x63\x6f\x6d")
}

func TestEncodeAnotherHeaderFieldWithNameAndLiteralValue(t *testing.T) {
	h := NewEncodingContext().EncodeField(HeaderField{"user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/33.0.1750.70 Safari/537.36"})

	assert.Equal(t, h, "\x7A\x68Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/33.0.1750.70 Safari/537.36")
}

func TestEncodeHeaderFieldFromAppendixD3(t *testing.T) {
	// http://tools.ietf.org/html/draft-ietf-httpbis-header-compression-07#appendix-D.3
	var h string

	context := NewEncodingContext()
	h = context.EncodeField(HeaderField{":method", "GET"})
	assert.Equal(t, h, "\x82")
	h = context.EncodeField(HeaderField{":scheme", "http"})
	assert.Equal(t, h, "\x87")
	h = context.EncodeField(HeaderField{":path", "/"})
	assert.Equal(t, h, "\x86")
	h = context.EncodeField(HeaderField{":authority", "www.example.com"})

	assert.Equal(t, h, "\x44\x0f\x77\x77\x77\x2e\x65\x78\x61\x6d\x70\x6c\x65\x2e\x63\x6f\x6d")
}

func TestEncodeHeaderFieldWithLiteralNameAndLiteralValue(t *testing.T) {
	var h string

	context := NewEncodingContext()
	h = context.EncodeField(HeaderField{"custom-header", "puppy-dogs"})

	assert.Equal(t, h, "\x40\x0dcustom-header\x0apuppy-dogs")

	assert.Equal(t, context.HeaderTable.Entries[0],
		HeaderField{"custom-header", "puppy-dogs"})
}

func TestEncodeHeaderSet(t *testing.T) {
	context := NewEncodingContext()

	h := context.Encode(HeaderSet{ []HeaderField{
		{":method", "GET"},
		{":scheme", "http"},
		{":path", "/"},
		{":authority", "www.example.com"},
	}})

	assert.Equal(t, h, "\x82\x87\x86\x44\x0f\x77\x77\x77\x2e\x65\x78\x61\x6d\x70\x6c\x65\x2e\x63\x6f\x6d")
}

func TestEncodeHeaderSetWithReferenceSet(t *testing.T) {
	context := NewEncodingContext()

	context.Encode(HeaderSet{ []HeaderField{
		{":method", "GET"},
		{":scheme", "http"},
		{":path", "/"},
		{":authority", "www.example.com"},
	}})

	h := context.Encode(HeaderSet{ []HeaderField{
		{":method", "GET"},
		{":scheme", "http"},
		{":path", "/"},
		{":authority", "www.example.com"},
		{"cache-control", "no-cache"},
	}})

	assert.Equal(t, h, "\x5c\x08\x6e\x6f\x2d\x63\x61\x63\x68\x65")
}

func TestEncodeHeaderSetWithReferenceSetEmptying(t *testing.T) {
	context := NewEncodingContext()

	context.Encode(HeaderSet{ []HeaderField{
		{":method", "GET"},
		{":scheme", "http"},
		{":path", "/"},
		{":authority", "www.example.com"},
	}})

	assert.Equal(t, context.HeaderTable.Size(), 180)

	context.Encode(HeaderSet{ []HeaderField{
		{":method", "GET"},
		{":scheme", "http"},
		{":path", "/"},
		{":authority", "www.example.com"},
		{"cache-control", "no-cache"},
	}})

	assert.Equal(t, context.HeaderTable.Size(), 233)

	context.Update.ReferenceSetEmptying = true
	h := context.Encode(HeaderSet{ []HeaderField{
		{":method", "GET"},
		{":scheme", "https"},
		{":path", "/index.html"},
		{":authority", "www.example.com"},
		{"custom-key", "custom-value"},
	}})

	assert.Equal(t, context.Update.ReferenceSetEmptying, false)
	assert.Equal(t, h, "\x30\x85\x8c\x8b\x84\x40\x0a\x63\x75\x73\x74\x6f\x6d\x2d\x6b\x65\x79\x0c\x63\x75\x73\x74\x6f\x6d\x2d\x76\x61\x6c\x75\x65")
	assert.Equal(t, context.HeaderTable.Size(), 379)
}

func TestEncodeHeaderSetWithEviction(t *testing.T) {
	var h string

	context := NewEncodingContext()
	context.HeaderTable.MaxSize = 256

	h = context.Encode(HeaderSet{ []HeaderField{
		{":status", "302"},
		{"cache-control", "private"},
		{"date", "Mon, 21 Oct 2013 20:13:21 GMT"},
		{"location", "https://www.example.com"},
	}})

	assert.Equal(t, h, "\x48\x03\x33\x30\x32\x59\x07\x70\x72\x69\x76\x61\x74\x65\x63\x1d\x4d\x6f\x6e\x2c\x20\x32\x31\x20\x4f\x63\x74\x20\x32\x30\x31\x33\x20\x32\x30\x3a\x31\x33\x3a\x32\x31\x20\x47\x4d\x54\x71\x17\x68\x74\x74\x70\x73\x3a\x2f\x2f\x77\x77\x77\x2e\x65\x78\x61\x6d\x70\x6c\x65\x2e\x63\x6f\x6d")
	assert.Equal(t, context.HeaderTable.Size(), 222)

	h = context.Encode(HeaderSet { []HeaderField{
		{":status", "200"},
		{"cache-control", "private"},
		{"date", "Mon, 21 Oct 2013 20:13:21 GMT"},
		{"location", "https://www.example.com"},
	}})

	assert.Equal(t, h, "\x8c")
	assert.Equal(t, len(context.HeaderTable.Entries), 4, "Should have evicted header to make room")
	assert.Equal(t, context.HeaderTable.Size(), 222)

	h = context.Encode(HeaderSet { []HeaderField{
		{":status", "200"},
		{"cache-control", "private"},
		{"date", "Mon, 21 Oct 2013 20:13:22 GMT"},
		{"location", "https://www.example.com"},
		{"content-encoding", "gzip"},
		{"set-cookie", "foo=ASDJKHQKBZXOQWEOPIUAXQWEOIU; max-age=3600; version=1"},
	}})

	// need to figure out how the encoder indicates reference set management
	// to the decoder here
	t.SkipNow()
	assert.Equal(t, h, "\x84\x84\x43\x1d\x4d\x6f\x6e\x2c\x20\x32\x31\x20\x4f\x63\x74\x20\x32\x30\x31\x33\x20\x32\x30\x3a\x31\x33\x3a\x32\x32\x20\x47\x4d\x54\x5e\x04\x67\x7a\x69\x70\x84\x84\x83\x83\x7b\x38\x66\x6f\x6f\x3d\x41\x53\x44\x4a\x4b\x48\x51\x4b\x42\x5a\x58\x4f\x51\x57\x45\x4f\x50\x49\x55\x41\x58\x51\x57\x45\x4f\x49\x55\x3b\x20\x6d\x61\x78\x2d\x61\x67\x65\x3d\x33\x36\x30\x30\x3b\x20\x76\x65\x72\x73\x69\x6f\x6e\x3d\x31")
}
