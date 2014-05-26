package http2

import (
	"github.com/tildedave/go-hpack-impl/hpack"
	"net"
)

// Internal structure for managing a server connection
// Connections own streams
type conn struct {
	context *hpack.EncodingContext
	ioc     net.Conn
}

type Conn interface {
	Read(p []byte) (n int, err error)
	Write(p []byte) (n int, err error)
	Close() error
	EncodeHeaderSet(hs hpack.HeaderSet) string
	HuffmanEncodeHeaderSet(hs hpack.HeaderSet) string
}

func NewConn(ioc net.Conn) *conn {
	return &conn{
		context: hpack.NewEncodingContext(),
		ioc:     ioc,
	}
}

func (c *conn) Read(p []byte) (n int, err error) {
	return c.ioc.Read(p)
}

func (c *conn) Write(p []byte) (n int, err error) {
	return c.ioc.Write(p)
}

func (c *conn) Close() (err error) {
	return c.ioc.Close()
}

func (c *conn) EncodeHeaderSet(hs hpack.HeaderSet) string {
	return c.context.Encode(hs)
}

func (c *conn) HuffmanEncodeHeaderSet(hs hpack.HeaderSet) string {
	return c.context.Encode(hs)
}
