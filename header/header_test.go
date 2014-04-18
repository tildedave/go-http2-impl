package header

// http://tools.ietf.org/html/draft-ietf-httpbis-header-compression-07

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEncodeHeaderFieldAddsToHeaderTable(t *testing.T) {
	table := HeaderTable{}

	Encode(HeaderField{":method", "GET"}, &table)

	assert.Equal(t, table.Entries[0], HeaderField{":method", "GET"})
}

func TestEncodeHeaderFieldUsesHeaderTableSizeAsOffset(t *testing.T) {
	// http://tools.ietf.org/html/draft-ietf-httpbis-header-compression-07#page-35
	table := HeaderTable{}

	Encode(HeaderField{ ":method", "GET" }, &table)
	h := Encode(HeaderField{ ":scheme", "http" }, &table)
	assert.Equal(t, h, "\x87")
}

func TestEncodeHeaderFieldFromStaticTable(t *testing.T) {
	table := HeaderTable{}

	var encoded string

	encoded = Encode(HeaderField{":method", "GET"}, &table)
	assert.Equal(t, encoded, "\x82")
	encoded = Encode(HeaderField{":path", "/"}, &table)
	assert.Equal(t, encoded, "\x85")

	assert.Equal(t, table.Entries[0], HeaderField{":path", "/"})
	assert.Equal(t, table.Entries[1], HeaderField{":method", "GET"})
}

func TestEncodeHeaderFieldWithNameAndLiteralValue(t *testing.T) {
	h := Encode(HeaderField{":authority", "www.example.com"}, &HeaderTable{})

	assert.Equal(t, h, "\x41\x0f\x77\x77\x77\x2e\x65\x78\x61\x6d\x70\x6c\x65\x2e\x63\x6f\x6d")
}

func TestEncodeAnotherHeaderFieldWithNameAndLiteralValue(t *testing.T) {
	h := Encode(HeaderField{"user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/33.0.1750.70 Safari/537.36"}, &HeaderTable{})

	assert.Equal(t, h, "\x7A\x68Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/33.0.1750.70 Safari/537.36")
}

func TestEncodeHeaderFieldFromAppendixD3(t *testing.T) {
	// http://tools.ietf.org/html/draft-ietf-httpbis-header-compression-07#appendix-D.3
	var h string

	table := HeaderTable{}
	h = Encode(HeaderField{":method", "GET"}, &table)
	assert.Equal(t, h, "\x82")
	h = Encode(HeaderField{":scheme", "http"}, &table)
	assert.Equal(t, h, "\x87")
	h = Encode(HeaderField{":path", "/"}, &table)
	assert.Equal(t, h, "\x86")
	h = Encode(HeaderField{":authority", "www.example.com"}, &table)
	assert.Equal(t, h, "\x44\x0f\x77\x77\x77\x2e\x65\x78\x61\x6d\x70\x6c\x65\x2e\x63\x6f\x6d")
}
