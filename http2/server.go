package http2

import (
	"bufio"
	"fmt"
	"strings"
)

var _ = fmt.Printf // package fmt is now used

const preface = "PRI * HTTP/2.0\r\n\r\nSM\r\n\r\n"

func serve(conn Conn) error {
	scanner := bufio.NewScanner(conn)
	str := ""

	// TODO: connection upgrade from HTTP 1.0
	for stopped := scanner.Scan(); stopped != false; stopped = scanner.Scan() {
		str += scanner.Text() + "\r\n"
		if !strings.HasPrefix(preface, str) {
			f := GOAWAY{0, 1, "Did not include connection preface"}
			conn.Write(f.Marshal())
			conn.Close()

			return nil
		}

		if preface == str {
			break
		}
	}

	conn.Write([]byte(preface))
	conn.Write(SETTINGS{}.Marshal())
	// wait for ACK

	return nil
}
