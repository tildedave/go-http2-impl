package hpack

func encodeInteger(i int, prefixSize uint) string {
	maxEncoded := (1 << prefixSize) - 1
	if i < maxEncoded {
		return string(byte(i))
	}

	repr := make([]byte, 0)
	repr = append(repr, byte(maxEncoded))

	i -= maxEncoded
	for ; i >= 128; {
		repr = append(repr, byte(i % 128 + 128))
		i /= 128
	}
	repr = append(repr, byte(i))

	return string(repr)
}
