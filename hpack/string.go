package hpack

import "fmt"

func encodeLiteral(literal string) string {
	return encodeInteger(len(literal), 7) + literal
}

func encodeLiteralHuffman(literal string) string {
	str := EncodeHuffman(literal)

	fmt.Println("encoding stuff", literal, []byte(str), len(str))

	lenBytes := []byte(encodeInteger(len(str), 7))
	lenBytes[0] |= 0x80

	return string(lenBytes) + str
}

func decodeLiteral(wire *[]byte) string {
	isHuffman := (*wire)[0] >> 7 == byte(0x01)
	len := decodeInteger(wire, 7)

	toDecode := (*wire)[0:len]
	*wire = (*wire)[len:]

	if isHuffman {
		if huffmanTree == nil {
			buildHuffmanTree()
		}

		decoded, _ := decodeStringHuffman(string(toDecode))
		fmt.Println("it was" , decoded)
		return decoded
	} else {
		return string(toDecode)
	}
}
