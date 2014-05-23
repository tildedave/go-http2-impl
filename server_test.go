package main

import (
	"testing"
	"github.com/tildedave/go-http2-impl/frame"
	"github.com/stretchr/testify/mock"
)

type MockConn struct{
	mock.Mock
}

func (c *MockConn) Close() error {
	args := c.Mock.Called()
	return args.Error(0)
}

func (c *MockConn) Write(b []byte) (int, error) {
	args := c.Mock.Called(b)
	return args.Int(0), args.Error(1)
}


func NewTestServer() (Server, *MockConn) {
	conn := new(MockConn)
	s := Server{ conn }

	return s, conn
}

func TestRespondWithoutPreface(t *testing.T) {
	server, conn := NewTestServer()

	f := frame.GOAWAY{0, 1, "Did not include connection preface"}
	bytes := f.Marshal()

	conn.On("Close").Return(nil)
	conn.On("Write", bytes).Return(len(bytes), nil)

	server.Respond("not the preface")

	conn.Mock.AssertExpectations(t)
}

func TestRespondWithThePreface(t *testing.T) {
	server, conn := NewTestServer()

	// also needs to write settings frame too.
	preface := "PRI * HTTP/2.0\r\n\r\nSM\r\n\r\n"
	conn.On("Write", []byte(preface)).Return(len(preface), nil)

	server.Respond("PRI * HTTP/2.0\r\n\r\nSM\r\n\r\n")

	conn.Mock.AssertExpectations(t)
}
