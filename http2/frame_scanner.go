package http2

import (
	"bufio"
	"io"
)

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
