package server

import (
	"http2-impl/common"
)

func Respond(data []byte) (common.Frame, bool) {
	f := common.Frame{ 2, 0, 0, "        Client did not send prefix" }

	return f, false
}
