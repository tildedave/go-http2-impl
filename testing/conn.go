package testing

import (
	"github.com/tildedave/go-http2-impl/http2"
	"net"
	"time"
)

type dummyAddr string

func (a dummyAddr) Network() string {
	return string(a)
}

func (a dummyAddr) String() string {
	return string(a)
}

type FakeConn struct {
	ReadData [][]byte
	Written  []byte
	Closed   bool
}

func (c *FakeConn) Close() error {
	c.Closed = true
	return nil
}

func (c *FakeConn) Read(b []byte) (int, error) {
	if len(c.ReadData) == 0 {
		return 0, nil
	}
	n := copy(b, c.ReadData[0])
	c.ReadData = c.ReadData[1:]

	return n, nil
}

func (c *FakeConn) Write(b []byte) (int, error) {
	c.Written = append(c.Written, b...)
	return len(b), nil
}

func (c *FakeConn) LocalAddr() net.Addr                { return dummyAddr("local-addr") }
func (c *FakeConn) RemoteAddr() net.Addr               { return dummyAddr("remote-addr") }
func (c *FakeConn) SetDeadline(t time.Time) error      { return nil }
func (c *FakeConn) SetReadDeadline(t time.Time) error  { return nil }
func (c *FakeConn) SetWriteDeadline(t time.Time) error { return nil }

func (c *FakeConn) Clear() {
	c.ReadData = make([][]byte, 0)
	c.Written = make([]byte, 0)
}

func NewFakeConn() *FakeConn {
	conn := new(FakeConn)
	conn.Written = make([]byte, 0)

	return conn
}

func NewTestConn() (http2.Conn, *FakeConn) {
	ioc := NewFakeConn()
	conn := http2.NewConn(ioc)

	return conn, ioc
}
