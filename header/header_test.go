package header

// http://tools.ietf.org/html/draft-ietf-httpbis-header-compression-07

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func twoHeaderFields(name1 string, value1 string, name2 string, value2 string) HeaderSet {
	return HeaderSet{[]HeaderField{HeaderField{name1, value1},
		HeaderField{name2, value2}}}
}

func TestEncodeHeaderFieldFromStaticTable(t *testing.T) {
	assert.Equal(t, Encode(HeaderField{":method", "GET"}), "\x82")
	assert.Equal(t, Encode(HeaderField{":method", "POST"}), "\x83")
	assert.Equal(t, Encode(HeaderField{":path", "/"}), "\x84")
}

func TestEncodeHeaderFieldWithNameAndLiteralValue(t *testing.T) {
	// From http://tools.ietf.org/html/draft-ietf-httpbis-header-compression-07#appendix-D.3
	h := Encode(HeaderField{":authority", "www.example.com"})

	assert.Equal(t, h, "\x44\x0f\x77\x77\x77\x2e\x65\x78\x61\x6d\x70\x6c\x65\x2e\x63\x6f\x6d")
}
