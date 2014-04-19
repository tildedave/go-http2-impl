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

/*
func TestDecodeHeaderWithIndexedNameAndValue(t *testing.T) {
	table := HeaderTable{}
	decoded, _ := Decode("\x82".Encode(&table)
	assert.Equal(t, []HeaderField{{":method", "GET"}}, decoded)
}
*/
