package main

import (
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
	"time"
)

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

func newFakeConn() *fakeConn {
	conn := new(fakeConn)
	conn.written = make([]byte, 0)

	return conn
}

func NewTestServer() (Server, *fakeConn) {
	conn := newFakeConn()
	s := Server{}

	return s, conn
}

func TestInitiateConnWithoutPreface(t *testing.T) {
	server, conn := NewTestServer()

	f := GOAWAY{0, 1, "Did not include connection preface"}
	bytes := f.Marshal()

	conn.readData = [][]byte{[]byte("not the preface")}
	server.InitiateConn(conn)

	assert.Equal(t, conn.written, bytes)
	assert.True(t, conn.closed, "Should have closed the connection")
}

func TestRespondWithThePreface(t *testing.T) {
	server, conn := NewTestServer()

	// also needs to write settings frame too.
	preface := "PRI * HTTP/2.0\r\n\r\nSM\r\n\r\n"

	conn.readData = [][]byte{[]byte("PRI * HTTP/2.0\r\n\r\nSM\r\n\r\n")}
	server.InitiateConn(conn)

	assert.Equal(t, conn.written, []byte(preface))
	assert.False(t, conn.closed)
}

func TestFrameScannerReturnsAFrame(t *testing.T) {
	conn := newFakeConn()
	b := PING{OpaqueData: 3957102}.Marshal()
	conn.readData = [][]byte{b}

	s := NewFrameScanner(conn)

	assert.True(t, s.Scan())
	assert.Equal(t, s.Bytes(), b)
}

func TestFrameScanner_IncompleteFrame(t *testing.T) {
	conn := newFakeConn()
	b := PING{OpaqueData: 3957102}.Marshal()
	conn.readData = [][]byte{b[0 : len(b)-1]}

	s := NewFrameScanner(conn)

	assert.False(t, s.Scan())
}

func TestFrameScanner_IncompleteFrameThatIsLaterCompleted(t *testing.T) {
	conn := newFakeConn()
	b := PING{OpaqueData: 3957102}.Marshal()
	conn.readData = [][]byte{b[0 : len(b)-1], b[len(b)-1:]}

	s := NewFrameScanner(conn)

	assert.True(t, s.Scan())
	assert.Equal(t, s.Bytes(), b)
}

func TestFrameScanner_TwoFrames(t *testing.T) {
	conn := newFakeConn()
	b1 := PING{OpaqueData: 3957102}.Marshal()
	b2 := PING{OpaqueData: 12311}.Marshal()

	conn.readData = [][]byte{b1, b2}

	s := NewFrameScanner(conn)

	assert.True(t, s.Scan())
	assert.Equal(t, s.Bytes(), b1)
	assert.True(t, s.Scan())
	assert.Equal(t, s.Bytes(), b2)
	assert.False(t, s.Scan())
}

func TestFrameScanner_TwoFramesCombined(t *testing.T) {
	conn := newFakeConn()
	b1 := PING{OpaqueData: 3957102}.Marshal()
	b2 := PING{OpaqueData: 12311}.Marshal()

	conn.readData = [][]byte{append(b1, b2...)}

	s := NewFrameScanner(conn)

	assert.True(t, s.Scan())
	assert.Equal(t, s.Bytes(), b1)
	assert.True(t, s.Scan())
	assert.Equal(t, s.Bytes(), b2)
	assert.False(t, s.Scan())
}

func TestFrameScanner_TwoFrames_Uneven(t *testing.T) {
	conn := newFakeConn()
	b1 := PING{OpaqueData: 3957102}.Marshal()
	b2 := PING{OpaqueData: 12311}.Marshal()

	conn.readData = [][]byte{b1[0:13], append(b1[13:], b2...)}

	s := NewFrameScanner(conn)

	assert.True(t, s.Scan())
	assert.Equal(t, s.Bytes(), b1)
	assert.True(t, s.Scan())
	assert.Equal(t, s.Bytes(), b2)
	assert.False(t, s.Scan())
}
