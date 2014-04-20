package server

import (
	"http2-impl/common"
)

func Respond(data []byte) (common.Frame, bool) {
	f := common.GOAWAY{}

	return f, false
}
