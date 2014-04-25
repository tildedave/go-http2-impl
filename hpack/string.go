package hpack

func encodeLiteral(literal string) string {
	return encodeInteger(len(literal), 7) + literal
}

func decodeLiteral(wire *[]byte) string {
	len := decodeInteger(wire, 7)

	decoded := (*wire)[0:len]
	*wire = (*wire)[len:]

	return string(decoded)
}
