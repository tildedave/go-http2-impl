package server

import "testing"

func TestRespondWithoutPrefaceClosesConnection(t *testing.T) {
	_, closedConn := Respond("something that is not the preface")
	if closedConn != false {
		t.Errorf("Should have closed the connection")
	}
}
