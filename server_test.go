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

func newTestConn() (conn, *fakeConn) {
	ioc := newFakeConn()
	conn := newConn(ioc)

	return conn, ioc
}

func TestInitiateConnWithoutPreface(t *testing.T) {
	conn, fakeConn := newTestConn()

	f := GOAWAY{0, 1, "Did not include connection preface"}
	bytes := f.Marshal()

	fakeConn.readData = [][]byte{[]byte("not the preface")}
	conn.serve()

	assert.Equal(t, fakeConn.written, bytes)
	assert.True(t, fakeConn.closed, "Should have closed the connection")
}

func TestRespondWithThePreface(t *testing.T) {
	conn, fakeConn := newTestConn()

	// also needs to write settings frame too.
	preface := "PRI * HTTP/2.0\r\n\r\nSM\r\n\r\n"

	fakeConn.readData = [][]byte{[]byte("PRI * HTTP/2.0\r\n\r\nSM\r\n\r\n")}
	conn.serve()

	assert.Equal(t, fakeConn.written, []byte(preface))
	assert.False(t, fakeConn.closed)
}

func TestFrameScannerReturnsAFrame(t *testing.T) {
	fakeConn := newFakeConn()
	b := PING{OpaqueData: 3957102}.Marshal()
	fakeConn.readData = [][]byte{b}

	s := NewFrameScanner(fakeConn)

	assert.True(t, s.Scan())
	assert.Equal(t, s.Bytes(), b)
}

func TestFrameScanner_IncompleteFrame(t *testing.T) {
	fakeConn := newFakeConn()
	b := PING{OpaqueData: 3957102}.Marshal()
	fakeConn.readData = [][]byte{b[0 : len(b)-1]}

	s := NewFrameScanner(fakeConn)

	assert.False(t, s.Scan())
}

func TestFrameScanner_IncompleteFrameThatIsLaterCompleted(t *testing.T) {
	fakeConn := newFakeConn()
	b := PING{OpaqueData: 3957102}.Marshal()
	fakeConn.readData = [][]byte{b[0 : len(b)-1], b[len(b)-1:]}

	s := NewFrameScanner(fakeConn)

	assert.True(t, s.Scan())
	assert.Equal(t, s.Bytes(), b)
}

func TestFrameScanner_TwoFrames(t *testing.T) {
	fakeConn := newFakeConn()
	b1 := PING{OpaqueData: 3957102}.Marshal()
	b2 := PING{OpaqueData: 12311}.Marshal()

	fakeConn.readData = [][]byte{b1, b2}

	s := NewFrameScanner(fakeConn)

	assert.True(t, s.Scan())
	assert.Equal(t, s.Bytes(), b1)
	assert.True(t, s.Scan())
	assert.Equal(t, s.Bytes(), b2)
	assert.False(t, s.Scan())
}

func TestFrameScanner_TwoFramesCombined(t *testing.T) {
	fakeConn := newFakeConn()
	b1 := PING{OpaqueData: 3957102}.Marshal()
	b2 := PING{OpaqueData: 12311}.Marshal()

	fakeConn.readData = [][]byte{append(b1, b2...)}

	s := NewFrameScanner(fakeConn)

	assert.True(t, s.Scan())
	assert.Equal(t, s.Bytes(), b1)
	assert.True(t, s.Scan())
	assert.Equal(t, s.Bytes(), b2)
	assert.False(t, s.Scan())
}

func TestFrameScanner_TwoFrames_Uneven(t *testing.T) {
	fakeConn := newFakeConn()
	b1 := PING{OpaqueData: 3957102}.Marshal()
	b2 := PING{OpaqueData: 12311}.Marshal()

	fakeConn.readData = [][]byte{b1[0:13], append(b1[13:], b2...)}

	s := NewFrameScanner(fakeConn)

	assert.True(t, s.Scan())
	assert.Equal(t, s.Bytes(), b1)
	assert.True(t, s.Scan())
	assert.Equal(t, s.Bytes(), b2)
	assert.False(t, s.Scan())
}
