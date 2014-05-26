package main

import (
	"bufio"
	"fmt"
	"github.com/tildedave/go-hpack-impl/hpack"
	"io"
	"strings"
)

var _ = fmt.Printf // package fmt is now used

type Conn interface {
	Read(b []byte) (n int, err error)
	Write(b []byte) (n int, err error)
	Close() error
}

type Server struct {
	EncodingContext *hpack.EncodingContext
}

const preface = "PRI * HTTP/2.0\r\n\r\nSM\r\n\r\n"

func (s *Server) InitiateConn(conn Conn) error {
	scanner := bufio.NewScanner(conn)
	str := ""

	// TODO: connection upgrade from HTTP 1.0
	for stopped := scanner.Scan(); stopped != false; stopped = scanner.Scan() {
		str += scanner.Text() + "\r\n"
		if !strings.HasPrefix(preface, str) {
			f := GOAWAY{0, 1, "Did not include connection preface"}
			conn.Write(f.Marshal())
			conn.Close()
			return nil
		}

		if preface == str {
			break
		}
	}

	conn.Write([]byte(preface))
	s.EncodingContext = hpack.NewEncodingContext()
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
