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

func decodeInteger(wire *[]byte, prefixSize uint) uint {
	var m, i uint
	var mask byte

	mask = (1 << prefixSize) - 1

	w := (*wire)[0]
	*wire = (*wire)[1:]

	i = uint(w & mask)
	if i < uint(mask) {
		return i
	}

	for ;; {
		w = (*wire)[0]
		*wire = (*wire)[1:]

		i += uint(w & 127) * (1 << m)
		m += 7

		if w & 128 != 128 {
			return i
		}
	}

	return 0
}
