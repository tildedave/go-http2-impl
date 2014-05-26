package http2

import (
	"fmt"
	"github.com/tildedave/go-hpack-impl/hpack"
	"io"
	"io/ioutil"
	"net/url"
	"strings"
)

var _ = fmt.Println // fmt is now used

type Header map[string][]string

func (h Header) Add(key, value string) {
	key = strings.ToLower(key)
	if val, ok := h[key]; ok {
		h[key] = append(val, value)
	} else {
		h[key] = []string{value}
	}
}

func (h Header) Set(key, value string) {
	key = strings.ToLower(key)
	h[key] = []string{value}
}

func (h Header) Get(key string) string {
	if len(h[key]) != 0 {
		return h[key][0]
	}

	return ""
}

type Request struct {
	Method string
	URL    *url.URL
	Body   io.ReadCloser
	Conn   Conn
	Header Header
}

func (r *Request) Write() {
	headers := make([]hpack.HeaderField, len(r.Header)+2)
	headers[0] = hpack.HeaderField{":method", r.Method}
	headers[1] = hpack.HeaderField{":path", r.URL.Path}
	i := 2
	for key, vals := range r.Header {
		for _, val := range vals {
			headers[i] = hpack.HeaderField{key, val}
			i++
		}
	}

	data := r.Conn.EncodeHeaderSet(hpack.HeaderSet{headers})

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
	req.Header = make(Header)

	return req, nil
}
