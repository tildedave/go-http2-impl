package common

type Frame struct {
	Type uint8
	Flags uint8
	StreamIdentifier uint32
	Payload string
}

func Marshal(f Frame) []byte {
	length := len(f.Payload)

	return []byte{ byte(length >> 8), byte(length), f.Type, f.Flags,
		byte(f.StreamIdentifier >> 24),
		byte(f.StreamIdentifier >> 16),
		byte(f.StreamIdentifier >> 8),
		byte(f.StreamIdentifier)}
}
