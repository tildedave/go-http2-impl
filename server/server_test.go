package server

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestRespondWithoutPrefaceClosesConnection(t *testing.T) {
	_, closedConn := Respond([]byte("not the preface"))
	if closedConn != false {
		t.Errorf("Should have closed the connection")
	}
}

func TestRespondWithoutPrefaceReturnsGOAWAY(t *testing.T) {
	response, _ := Respond([]byte("not the preface"))
	//    All frames begin with an 8-octet header followed by a payload
	// of between 0 and 16,383 octets.

	assert.Equal(t, response.Type, uint8(2),
		"Server should set GOAWAY type 0x00000010")

	additionalDebugData := response.Payload[8:]

	assert.Equal(t, additionalDebugData, "Client did not send prefix",
		"Server should indicate the client needs to send prefix")
}
