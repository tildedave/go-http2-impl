package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/tildedave/go-http2-impl/http2"
	http2test "github.com/tildedave/go-http2-impl/testing"
	"testing"
)

var _ = fmt.Printf // package fmt is now used

func NewTestConn() (http2.Conn, *http2test.FakeConn) {
	ioc := http2test.NewFakeConn()
	conn := http2.NewConn(ioc)

	return conn, ioc
}

func TestServeWithoutPreface(t *testing.T) {
	conn, fakeConn := NewTestConn()

	f := http2.GOAWAY{0, 1, "Did not include connection preface"}
	bytes := f.Marshal()

	fakeConn.ReadData = [][]byte{[]byte("not the preface")}
	serve(conn)

	assert.Equal(t, fakeConn.Written, bytes)
	assert.True(t, fakeConn.Closed, "Should have closed the connection")
}

func TestServeWithThePreface(t *testing.T) {
	conn, fakeConn := NewTestConn()

	preface := "PRI * HTTP/2.0\r\n\r\nSM\r\n\r\n"

	fakeConn.ReadData = [][]byte{[]byte("PRI * HTTP/2.0\r\n\r\nSM\r\n\r\n")}
	serve(conn)

	assert.Equal(t, fakeConn.Written[0:len(preface)], []byte(preface))
	assert.False(t, fakeConn.Closed)
}

func TestServeWithThePrefaceSendsSettingsFrame(t *testing.T) {
	conn, fakeConn := NewTestConn()

	preface := "PRI * HTTP/2.0\r\n\r\nSM\r\n\r\n"
	settingsFrame := http2.SETTINGS{}

	fakeConn.ReadData = [][]byte{[]byte("PRI * HTTP/2.0\r\n\r\nSM\r\n\r\n")}
	serve(conn)

	assert.Equal(t, fakeConn.Written[len(preface):], settingsFrame.Marshal())
	assert.False(t, fakeConn.Closed)
}
