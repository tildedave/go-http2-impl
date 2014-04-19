package hpack

// http://tools.ietf.org/html/draft-ietf-httpbis-header-compression-07

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEncodeHeaderFieldAddsToHeaderTable(t *testing.T) {
	table := HeaderTable{}

	HeaderField{":method", "GET"}.Encode(&table)

	assert.Equal(t, table.Entries[0], HeaderField{":method", "GET"})
}

func TestEncodeHeaderFieldUsesHeaderTableSizeAsOffset(t *testing.T) {
	// http://tools.ietf.org/html/draft-ietf-httpbis-header-compression-07#page-35
	table := HeaderTable{}

	HeaderField{ ":method", "GET" }.Encode(&table)
	h := HeaderField{ ":scheme", "http" }.Encode(&table)
	assert.Equal(t, h, "\x87")
}

func TestEncodeHeaderFieldFromStaticTable(t *testing.T) {
	table := HeaderTable{}

	var encoded string

	encoded = HeaderField{":method", "GET"}.Encode(&table)
	assert.Equal(t, encoded, "\x82")
	encoded = HeaderField{":path", "/"}.Encode(&table)
	assert.Equal(t, encoded, "\x85")

	assert.Equal(t, table.Entries[0], HeaderField{":path", "/"})
	assert.Equal(t, table.Entries[1], HeaderField{":method", "GET"})
}

func TestEncodeHeaderFieldWithNameAndLiteralValue(t *testing.T) {
	h := HeaderField{":authority", "www.example.com"}.Encode(&HeaderTable{})

	assert.Equal(t, h, "\x41\x0f\x77\x77\x77\x2e\x65\x78\x61\x6d\x70\x6c\x65\x2e\x63\x6f\x6d")
}

func TestEncodeAnotherHeaderFieldWithNameAndLiteralValue(t *testing.T) {
	h := HeaderField{"user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/33.0.1750.70 Safari/537.36"}.Encode(&HeaderTable{})

	assert.Equal(t, h, "\x7A\x68Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/33.0.1750.70 Safari/537.36")
}

func TestEncodeHeaderFieldFromAppendixD3(t *testing.T) {
	// http://tools.ietf.org/html/draft-ietf-httpbis-header-compression-07#appendix-D.3
	var h string

	table := HeaderTable{}
	h = HeaderField{":method", "GET"}.Encode(&table)
	assert.Equal(t, h, "\x82")
	h = HeaderField{":scheme", "http"}.Encode(&table)
	assert.Equal(t, h, "\x87")
	h = HeaderField{":path", "/"}.Encode(&table)
	assert.Equal(t, h, "\x86")
	h = HeaderField{":authority", "www.example.com"}.Encode(&table)
	assert.Equal(t, h, "\x44\x0f\x77\x77\x77\x2e\x65\x78\x61\x6d\x70\x6c\x65\x2e\x63\x6f\x6d")
}

func TestEncodeHeaderFieldWithLiteralNameAndLiteralValue(t *testing.T) {
	var h string

	table := HeaderTable{}
	h = HeaderField{"custom-header", "puppy-dogs"}.Encode(&table)
	assert.Equal(t, h, "\x40\x0dcustom-header\x0apuppy-dogs")

	assert.Equal(t, table.Entries[0], HeaderField{"custom-header", "puppy-dogs"})
}

func TestEncodeHeaderSet(t *testing.T) {
	table := HeaderTable{}
	refset := ReferenceSet{}

	h := HeaderSet{ []HeaderField{
		{":method", "GET"},
		{":scheme", "http"},
		{":path", "/"},
		{":authority", "www.example.com"},
	}}.Encode(&table, &refset)

	assert.Equal(t, h, "\x82\x87\x86\x44\x0f\x77\x77\x77\x2e\x65\x78\x61\x6d\x70\x6c\x65\x2e\x63\x6f\x6d")
}

func TestEncodeHeaderSetWithReferenceSet(t *testing.T) {
	table := HeaderTable{}
	refset := ReferenceSet{}

	HeaderSet{ []HeaderField{
		{":method", "GET"},
		{":scheme", "http"},
		{":path", "/"},
		{":authority", "www.example.com"},
	}}.Encode(&table, &refset)

	h := HeaderSet{ []HeaderField{
		{":method", "GET"},
		{":scheme", "http"},
		{":path", "/"},
		{":authority", "www.example.com"},
		{"cache-control", "no-cache"},
	}}.Encode(&table, &refset)

	assert.Equal(t, h, "\x5c\x08\x6e\x6f\x2d\x63\x61\x63\x68\x65")
}

func TestEncodeHeaderSetWithReferenceSetAndThreeRequests(t *testing.T) {
	table := HeaderTable{}
	refset := ReferenceSet{}

	HeaderSet{ []HeaderField{
		{":method", "GET"},
		{":scheme", "http"},
		{":path", "/"},
		{":authority", "www.example.com"},
	}}.Encode(&table, &refset)

	HeaderSet{ []HeaderField{
		{":method", "GET"},
		{":scheme", "http"},
		{":path", "/"},
		{":authority", "www.example.com"},
		{"cache-control", "no-cache"},
	}}.Encode(&table, &refset)

	h := HeaderSet{ []HeaderField{
		{":method", "GET"},
		{":scheme", "https"},
		{":path", "/index.html"},
		{":authority", "www.example.com"},
		{"custom-key", "custom-value"},
	}}.Encode(&table, &ReferenceSet{})

	assert.Equal(t, h, "\x30\x85\x8c\x8b\x84\x40\x0a\x63\x75\x73\x74\x6f\x6d\x2d\x6b\x65\x79\x0c\x63\x75\x73\x74\x6f\x6d\x2d\x76\x61\x6c\x75\x65")
}

/*
func TestDecodeHeaderWithIndexedNameAndValue(t *testing.T) {
	table := HeaderTable{}
	decoded, _ := Decode("\x82".Encode(&table)
	assert.Equal(t, []HeaderField{{":method", "GET"}}, decoded)
}
*/
