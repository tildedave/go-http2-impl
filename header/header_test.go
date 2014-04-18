package header

// http://tools.ietf.org/html/draft-ietf-httpbis-header-compression-07

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEncodeIndexedHeaderFieldFromStaticTable(t *testing.T) {
	set := HeaderSet{[]HeaderField{HeaderField{":method", "GET"}}}
	encoded := Encode(set)

	assert.Equal(t, encoded, []byte{0x82})
}
