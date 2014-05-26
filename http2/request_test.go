package http2

import (
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

	request, _ := NewRequest("GET", "http://www.google.com", strings.NewReader(""), conn)
	request.Write()

	// Should output to the fakeConn.
	t.Log(conn, fakeConn, request)
}
