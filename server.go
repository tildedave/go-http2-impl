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

func (s *Server) Respond(wire string) {
	f := frame.GOAWAY{0, 1, "Did not include connection preface"}
	s.conn.Write(f.Marshal())
	s.conn.Close()
}
