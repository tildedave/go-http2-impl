package http2

import (
	"github.com/tildedave/go-hpack-impl/hpack"
	"net"
)

// Internal structure for managing a server connection
// Connections own streams
type serverConn struct {
	context      *hpack.EncodingContext
	ioc          net.Conn
	lastStreamId int
}

type Conn interface {
	Read(p []byte) (n int, err error)
	Write(p []byte) (n int, err error)
	Close() error
	EncodeHeaderSet(hs hpack.HeaderSet) string
	NextStreamId() uint32
}

func NewServerConn(ioc net.Conn) *serverConn {
	return &serverConn{
		context:      hpack.NewEncodingContext(),
		ioc:          ioc,
		lastStreamId: 0,
	}
}

func (c *serverConn) Read(p []byte) (n int, err error) {
	return c.ioc.Read(p)
}

func (c *serverConn) Write(p []byte) (n int, err error) {
	return c.ioc.Write(p)
}

func (c *serverConn) Close() (err error) {
	return c.ioc.Close()
}

func (c *serverConn) EncodeHeaderSet(hs hpack.HeaderSet) string {
	// TODO: EncodeHuffman instead
	return c.context.Encode(hs)
}

func (c *serverConn) NextStreamId() uint32 {
	// TODO: overflow
	nextStreamId := uint32(c.lastStreamId + 2)
	c.lastStreamId += 2

	return nextStreamId
}
