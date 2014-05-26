package http2

import (
	"io"
	"io/ioutil"
	"net/url"
)

type Request struct {
	Method string
	URL    *url.URL
	Body   io.ReadCloser
	Conn   Conn
}

func (r *Request) Write(w io.Writer) {
	// TODO
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
