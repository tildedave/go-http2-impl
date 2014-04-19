package hpack

// http://tools.ietf.org/html/draft-ietf-httpbis-header-compression-07

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEncodeHeaderFieldAddsToHeaderTable(t *testing.T) {
	context := EncodingContext{}

	HeaderField{":method", "GET"}.Encode(&context)

	assert.Equal(t, context.HeaderTable.Entries[0],
		HeaderField{":method", "GET"})
}

func TestEncodeHeaderFieldUsesHeaderTableSizeAsOffset(t *testing.T) {
	// http://tools.ietf.org/html/draft-ietf-httpbis-header-compression-07#page-35
	context := EncodingContext{}

	HeaderField{ ":method", "GET" }.Encode(&context)
	h := HeaderField{ ":scheme", "http" }.Encode(&context)
	assert.Equal(t, h, "\x87")
}

func TestEncodeHeaderFieldFromStaticTable(t *testing.T) {
	context := EncodingContext{}

	var encoded string

	encoded = HeaderField{":method", "GET"}.Encode(&context)
	assert.Equal(t, encoded, "\x82")
	encoded = HeaderField{":path", "/"}.Encode(&context)
	assert.Equal(t, encoded, "\x85")

	assert.Equal(t, context.HeaderTable.Entries[0], HeaderField{":path", "/"})
	assert.Equal(t, context.HeaderTable.Entries[1], HeaderField{":method", "GET"})
}

func TestEncodeHeaderFieldWithNameAndLiteralValue(t *testing.T) {
	h := HeaderField{":authority", "www.example.com"}.Encode(&EncodingContext{})

	assert.Equal(t, h, "\x41\x0f\x77\x77\x77\x2e\x65\x78\x61\x6d\x70\x6c\x65\x2e\x63\x6f\x6d")
}

func TestEncodeAnotherHeaderFieldWithNameAndLiteralValue(t *testing.T) {
	h := HeaderField{"user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/33.0.1750.70 Safari/537.36"}.Encode(&EncodingContext{})

	assert.Equal(t, h, "\x7A\x68Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/33.0.1750.70 Safari/537.36")
}

func TestEncodeHeaderFieldFromAppendixD3(t *testing.T) {
	// http://tools.ietf.org/html/draft-ietf-httpbis-header-compression-07#appendix-D.3
	var h string

	context := EncodingContext{}
	h = HeaderField{":method", "GET"}.Encode(&context)
	assert.Equal(t, h, "\x82")
	h = HeaderField{":scheme", "http"}.Encode(&context)
	assert.Equal(t, h, "\x87")
	h = HeaderField{":path", "/"}.Encode(&context)
	assert.Equal(t, h, "\x86")
	h = HeaderField{":authority", "www.example.com"}.Encode(&context)
	assert.Equal(t, h, "\x44\x0f\x77\x77\x77\x2e\x65\x78\x61\x6d\x70\x6c\x65\x2e\x63\x6f\x6d")
}

func TestEncodeHeaderFieldWithLiteralNameAndLiteralValue(t *testing.T) {
	var h string

	context := EncodingContext{}
	h = HeaderField{"custom-header", "puppy-dogs"}.Encode(&context)
	assert.Equal(t, h, "\x40\x0dcustom-header\x0apuppy-dogs")

	assert.Equal(t, context.HeaderTable.Entries[0],
		HeaderField{"custom-header", "puppy-dogs"})
}

func TestEncodeHeaderSet(t *testing.T) {
	context := EncodingContext{}

	h := HeaderSet{ []HeaderField{
		{":method", "GET"},
		{":scheme", "http"},
		{":path", "/"},
		{":authority", "www.example.com"},
	}}.Encode(&context)

	assert.Equal(t, h, "\x82\x87\x86\x44\x0f\x77\x77\x77\x2e\x65\x78\x61\x6d\x70\x6c\x65\x2e\x63\x6f\x6d")
}

func TestEncodeHeaderSetWithReferenceSet(t *testing.T) {
	context := EncodingContext{}

	HeaderSet{ []HeaderField{
		{":method", "GET"},
		{":scheme", "http"},
		{":path", "/"},
		{":authority", "www.example.com"},
	}}.Encode(&context)

	h := HeaderSet{ []HeaderField{
		{":method", "GET"},
		{":scheme", "http"},
		{":path", "/"},
		{":authority", "www.example.com"},
		{"cache-control", "no-cache"},
	}}.Encode(&context)

	assert.Equal(t, h, "\x5c\x08\x6e\x6f\x2d\x63\x61\x63\x68\x65")
}

func TestEncodeHeaderSetWithReferenceSetAndThreeRequests(t *testing.T) {
	context := EncodingContext{}

	HeaderSet{ []HeaderField{
		{":method", "GET"},
		{":scheme", "http"},
		{":path", "/"},
		{":authority", "www.example.com"},
	}}.Encode(&context)

	HeaderSet{ []HeaderField{
		{":method", "GET"},
		{":scheme", "http"},
		{":path", "/"},
		{":authority", "www.example.com"},
		{"cache-control", "no-cache"},
	}}.Encode(&context)

	context.ReferenceSet = ReferenceSet{}
	h := HeaderSet{ []HeaderField{
		{":method", "GET"},
		{":scheme", "https"},
		{":path", "/index.html"},
		{":authority", "www.example.com"},
		{"custom-key", "custom-value"},
	}}.Encode(&context)

	assert.Equal(t, h, "\x85\x8c\x8b\x84\x40\x0a\x63\x75\x73\x74\x6f\x6d\x2d\x6b\x65\x79\x0c\x63\x75\x73\x74\x6f\x6d\x2d\x76\x61\x6c\x75\x65")
}

/*
func TestDecodeHeaderWithIndexedNameAndValue(t *testing.T) {
	context := EncodingContext{}
	decoded, _ := Decode("\x82".Encode(&table)
	assert.Equal(t, []HeaderField{{":method", "GET"}}, decoded)
}
*/
