package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
	"time"
)

var _ = fmt.Printf // package fmt is now used

type dummyAddr string

func (a dummyAddr) Network() string {
	return string(a)
}

func (a dummyAddr) String() string {
	return string(a)
}

type fakeConn struct {
	readData [][]byte
	written  []byte
	closed   bool
}

func (c *fakeConn) Close() error {
	c.closed = true
	return nil
}

func (c *fakeConn) Read(b []byte) (int, error) {
	if len(c.readData) == 0 {
		return 0, nil
	}
	n := copy(b, c.readData[0])
	c.readData = c.readData[1:]

	return n, nil
}

func (c *fakeConn) Write(b []byte) (int, error) {
	c.written = append(c.written, b...)
	return len(b), nil
}

func (c *fakeConn) LocalAddr() net.Addr                { return dummyAddr("local-addr") }
func (c *fakeConn) RemoteAddr() net.Addr               { return dummyAddr("remote-addr") }
func (c *fakeConn) SetDeadline(t time.Time) error      { return nil }
func (c *fakeConn) SetReadDeadline(t time.Time) error  { return nil }
func (c *fakeConn) SetWriteDeadline(t time.Time) error { return nil }

func (c *fakeConn) Clear() {
	c.readData = make([][]byte, 0)
	c.written = make([]byte, 0)
}

func newFakeConn() *fakeConn {
	conn := new(fakeConn)
	conn.written = make([]byte, 0)

	return conn
}

func newTestConn() (conn, *fakeConn) {
	ioc := newFakeConn()
	conn := newConn(ioc)

	return conn, ioc
}

func TestServeWithoutPreface(t *testing.T) {
	conn, fakeConn := newTestConn()

	f := GOAWAY{0, 1, "Did not include connection preface"}
	bytes := f.Marshal()

	fakeConn.readData = [][]byte{[]byte("not the preface")}
	conn.serve()

	assert.Equal(t, fakeConn.written, bytes)
	assert.True(t, fakeConn.closed, "Should have closed the connection")
}

func TestServeWithThePreface(t *testing.T) {
	conn, fakeConn := newTestConn()

	preface := "PRI * HTTP/2.0\r\n\r\nSM\r\n\r\n"

	fakeConn.readData = [][]byte{[]byte("PRI * HTTP/2.0\r\n\r\nSM\r\n\r\n")}
	conn.serve()

	assert.Equal(t, fakeConn.written[0:len(preface)], []byte(preface))
	assert.False(t, fakeConn.closed)
}

func TestServeWithThePrefaceSendsSettingsFrame(t *testing.T) {
	conn, fakeConn := newTestConn()

	preface := "PRI * HTTP/2.0\r\n\r\nSM\r\n\r\n"
	settingsFrame := SETTINGS{}

	fakeConn.readData = [][]byte{[]byte("PRI * HTTP/2.0\r\n\r\nSM\r\n\r\n")}
	conn.serve()

	assert.Equal(t, fakeConn.written[len(preface):], settingsFrame.Marshal())
	assert.False(t, fakeConn.closed)
}
