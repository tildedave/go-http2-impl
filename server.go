package main

import (
	"github.com/tildedave/go-http2-impl/frame"
)

func Respond(data []byte) (frame.Frame, bool) {
	f := frame.GOAWAY{}

	return f, false
}
