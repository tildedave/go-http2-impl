package header

// http://tools.ietf.org/html/draft-ietf-httpbis-header-compression-07

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func singleHeaderField(name string, value string) HeaderSet {
	return HeaderSet{[]HeaderField{HeaderField{name, value}}}
}

func twoHeaderFields(name1 string, value1 string, name2 string, value2 string) HeaderSet {
	return HeaderSet{[]HeaderField{HeaderField{name1, value1},
		HeaderField{name2, value2}}}
}

func TestEncodeNoHeadersIsNothing(t *testing.T) {
	assert.Equal(t, Encode(HeaderSet{}), "")
}

func TestEncodeIndexedHeaderFieldFromStaticTable(t *testing.T) {
	assert.Equal(t, Encode(singleHeaderField(":method", "GET")), "\x82")
	assert.Equal(t, Encode(singleHeaderField(":method", "POST")), "\x83")
	assert.Equal(t, Encode(singleHeaderField(":path", "/")), "\x84")
}

func TestEncodeIndexedHeaderFieldsFromStaticTable(t *testing.T) {
	assert.Equal(t, Encode(twoHeaderFields(":method", "GET", ":path", "/")),
		"\x82\x84")
}
