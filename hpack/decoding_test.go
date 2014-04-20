package hpack

// http://tools.ietf.org/html/draft-ietf-httpbis-header-compression-07

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDecode(t *testing.T)  {
	context := NewEncodingContext()

	headers := context.Decode("\x82")

	assert.Equal(t, len(headers.Headers), 1)
	assert.Equal(t, headers.Headers[0], HeaderField{":method", "GET"})
}
