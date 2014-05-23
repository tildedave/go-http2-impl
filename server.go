package main

import (
	"github.com/tildedave/go-http2-impl/frame"
)

type Conn interface {
	Write(b []byte) (n int, err error)
	Close() error
}

type Server struct {
	conn Conn
}

const preface = "PRI * HTTP/2.0\r\n\r\nSM\r\n\r\n"

func (s *Server) Respond(wire string) {
	if wire != preface {
		f := frame.GOAWAY{0, 1, "Did not include connection preface"}
		s.conn.Write(f.Marshal())
		s.conn.Close()
		return
	}

	s.conn.Write([]byte(preface))
	// TODO: SETTINGS frame
}
