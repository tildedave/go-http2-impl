package main

import (
	"github.com/tildedave/go-http2-impl/frame"
	"bufio"
	"strings"
)

type Conn interface {
	Read(b []byte) (n int, err error)
	Write(b []byte) (n int, err error)
	Close() error
}

type Server struct {
}

const preface = "PRI * HTTP/2.0\r\n\r\nSM\r\n\r\n"

func (s *Server) InitiateConn(conn Conn) error {
	scanner := bufio.NewScanner(conn)
	stopped := scanner.Scan()
	str := ""

	if stopped != false {
		str += scanner.Text() + "\r\n"
		if !strings.HasPrefix(preface, str) {
			f := frame.GOAWAY{0, 1, "Did not include connection preface"}
			conn.Write(f.Marshal())
			conn.Close()
			return nil
		}
	}

	conn.Write([]byte(preface))
	// TODO: SETTINGS frame

	return nil;
}
