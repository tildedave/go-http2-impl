package hpack

// TODO: headers of arbitrary length with integer encoding algorithm
// TODO: hpack test cases https://github.com/http2jp/hpack-test-case

func Decode(headers string, table *HeaderTable) ([]HeaderField, int) {
	var decodedHeaders []HeaderField

	return decodedHeaders, 0
}
