package http2

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/tildedave/go-hpack-impl/hpack"
	http2test "github.com/tildedave/go-http2-impl/testing"
	"strings"
	"testing"
)

func NewTestConn() (Conn, *http2test.FakeConn) {
	ioc := http2test.NewFakeConn()
	conn := NewConn(ioc)

	return conn, ioc
}

func TestWrite(t *testing.T) {
	conn, fakeConn := NewTestConn()

	request, _ := NewRequest("GET", "http://www.google.com/", strings.NewReader(""), conn)
	request.Write()

	s := NewFrameScanner(bytes.NewReader(fakeConn.Written))

	// expect one frame: HEADERS

	// TODO: What is the stream identifier?  (Depends on how connection generates
	// new stream ids.  Also, client initiated vs server initiated are different.)

	h := HEADERS{}
	h.Flags.END_HEADERS = true
	h.Flags.END_STREAM = true
	h.StreamId = 23

	context := hpack.NewEncodingContext()
	h.HeaderBlockFragment = context.Encode(hpack.HeaderSet{[]hpack.HeaderField{
		{":method", "GET"},
		{":path", "/"},
	}})

	assert.True(t, s.Scan())
	_, f, err := Unmarshal(s.Bytes())
	assert.Nil(t, err)
	assert.Equal(t, f, h)
}
