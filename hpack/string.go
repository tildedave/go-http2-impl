package hpack

func encodeLiteral(literal string) string {
	return encodeInteger(len(literal), 7) + literal
}

func encodeLiteralHuffman(literal string) string {
	str := EncodeHuffman(literal)
	lenBytes := []byte(encodeInteger(len(str), 7))
	lenBytes[0] |= 0x80

	return string(lenBytes) + str
}

func decodeLiteral(wire *[]byte) string {
	len := decodeInteger(wire, 7)

	decoded := (*wire)[0:len]
	*wire = (*wire)[len:]

	return string(decoded)
}
