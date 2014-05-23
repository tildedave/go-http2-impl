package main

import (
	"testing"
	"github.com/tildedave/go-http2-impl/frame"
	"github.com/stretchr/testify/assert"
)

type MockConn struct{
	written []byte
	closed bool
}

func (c *MockConn) Close() error {
	c.closed = true
	return nil
}

func (c *MockConn) Write(b []byte) (int, error) {
	c.written = append(c.written, b...)
	return len(b), nil
}


func NewTestServer() (Server, *MockConn) {
	conn := new(MockConn)
	conn.written = make([]byte, 0)

	s := Server{ conn }

	return s, conn
}

func TestRespondWithoutPreface(t *testing.T) {
	server, conn := NewTestServer()

	f := frame.GOAWAY{0, 1, "Did not include connection preface"}
	bytes := f.Marshal()

	server.Respond("not the preface")

	assert.Equal(t, conn.written, bytes)
	assert.True(t, conn.closed, "Should have closed the connection")
}

func TestRespondWithThePreface(t *testing.T) {
	server, conn := NewTestServer()

	// also needs to write settings frame too.
	preface := "PRI * HTTP/2.0\r\n\r\nSM\r\n\r\n"

	server.Respond("PRI * HTTP/2.0\r\n\r\nSM\r\n\r\n")

	assert.Equal(t, conn.written, []byte(preface))
	assert.False(t, conn.closed)
}
