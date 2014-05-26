package main

import (
	"bufio"
	"fmt"
	"github.com/tildedave/go-hpack-impl/hpack"
	"io"
	"net"
	"strings"
)

var _ = fmt.Printf // package fmt is now used

// Internal structure for managing a server connection
// Connections own streams
type conn struct {
	context *hpack.EncodingContext
	ioc     net.Conn
}

func newConn(ioc net.Conn) conn {
	return conn{
		context: hpack.NewEncodingContext(),
		ioc:     ioc,
	}
}

const preface = "PRI * HTTP/2.0\r\n\r\nSM\r\n\r\n"

func (c *conn) serve() error {
	scanner := bufio.NewScanner(c.ioc)
	str := ""

	// TODO: connection upgrade from HTTP 1.0
	for stopped := scanner.Scan(); stopped != false; stopped = scanner.Scan() {
		str += scanner.Text() + "\r\n"
		if !strings.HasPrefix(preface, str) {
			f := GOAWAY{0, 1, "Did not include connection preface"}
			c.ioc.Write(f.Marshal())
			c.ioc.Close()

			return nil
		}

		if preface == str {
			break
		}
	}

	c.ioc.Write([]byte(preface))
	// TODO: SETTINGS frame

	return nil
}

func NewFrameScanner(r io.Reader) *bufio.Scanner {
	s := bufio.NewScanner(r)
	s.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		advance, f, err := Unmarshal(data)
		if f != nil || err != nil {
			return advance, data[0:advance], err
		}

		return 0, nil, nil
	})
	return s
}
