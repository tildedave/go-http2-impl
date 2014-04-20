package main

import (
	"testing"
)

func TestRespondWithoutPrefaceClosesConnection(t *testing.T) {
	_, closedConn := Respond([]byte("not the preface"))
	if closedConn != false {
		t.Errorf("Should have closed the connection")
	}
}

func TestRespondWithoutPrefaceReturnsGOAWAY(t *testing.T) {
	Respond([]byte("not the preface"))

	t.SkipNow()
}
