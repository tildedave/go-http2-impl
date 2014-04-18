package header

// http://tools.ietf.org/html/draft-ietf-httpbis-header-compression-07

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func singleHeaderField(key string, value string) HeaderSet {
	return HeaderSet{[]HeaderField{HeaderField{key, value}}}
}

func TestEncodeNoHeadersIsNothing(t *testing.T) {
	assert.Equal(t, Encode(HeaderSet{}), []byte{})
}

func TestEncodeIndexedHeaderFieldFromStaticTable(t *testing.T) {
	assert.Equal(t, Encode(singleHeaderField(":method", "GET")), []byte{0x82})
	assert.Equal(t, Encode(singleHeaderField(":method", "POST")), []byte{0x83})
	assert.Equal(t, Encode(singleHeaderField(":path", "/")), []byte{0x84})
}
