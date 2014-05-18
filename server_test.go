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

func TestRespondWithoutPreface(t *testing.T) {
	conn := new(MockConn)

	f := frame.GOAWAY{0, 1, "Did not include connection preface"}
	bytes := f.Marshal()

	conn.On("Close").Return(nil)
	conn.On("Write", bytes).Return(len(bytes), nil)

	server := Server{}
	server.conn = conn

	server.Respond("not the preface")

	conn.Mock.AssertExpectations(t)
}
