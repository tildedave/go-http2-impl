package server

import (
	"go-http2-impl/common"
)

func Respond(data []byte) (common.Frame, bool) {
	f := common.GOAWAY{}

	return f, false
}
