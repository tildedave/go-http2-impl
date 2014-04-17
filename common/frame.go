package common

type Frame struct {
	Type uint8
	Flags uint8
	StreamIdentifier uint32
	Payload string
}

func Marshal(f Frame) []byte {
	length := len(f.Payload)

	return []byte{ byte(length >> 8), byte(length) }
}
