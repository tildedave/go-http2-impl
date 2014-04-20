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

func decodeInteger(encodedInteger string, prefixSize uint) uint {
	var mask, m, i uint

	mask = (1 << prefixSize) - 1
	i = uint(byte(encodedInteger[0]) & byte(mask))
	if i < mask {
		return i
	}

	for ;; {
		encodedInteger = encodedInteger[1:]
		i += uint(byte(encodedInteger[0]) & 127) * (1 << m)
		m += 7

		if encodedInteger[0] & 128 != 128 {
			return i
		}
	}

	return 0
}
