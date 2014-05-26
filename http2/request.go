package http2

import (
	"fmt"
	"github.com/tildedave/go-hpack-impl/hpack"
	"io"
	"io/ioutil"
	"net/url"
)

var _ = fmt.Println // fmt is now used

type Request struct {
	Method string
	URL    *url.URL
	Body   io.ReadCloser
	Conn   Conn
}

func (r *Request) Write() {
	hs := hpack.HeaderSet{[]hpack.HeaderField{
		{":method", r.Method},
		{":path", r.URL.Path},
	}}
	data := r.Conn.EncodeHeaderSet(hs)

	h := HEADERS{}
	h.HeaderBlockFragment = data
	h.Flags.END_STREAM = true
	h.Flags.END_HEADERS = true
	h.StreamId = 23

	r.Conn.Write(h.Marshal())
}

func NewRequest(method, urlStr string, body io.Reader, conn Conn) (*Request, error) {
	// Somewhat stolen from
	// http://golang.org/src/pkg/net/http/request.go?s=13983:14055#L434
	// See if I can reuse this later maybe

	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	rc, ok := body.(io.ReadCloser)
	if !ok && body != nil {
		rc = ioutil.NopCloser(body)
	}

	req := &Request{
		Method: method,
		URL:    u,
		Body:   rc,
		Conn:   conn,
	}

	return req, nil
}
