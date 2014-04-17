package common

type Frame struct {
	Type uint8
	Flags uint8
	StreamIdentifier uint32
	Payload string
}

func Marshal(f Frame) string {
	return "akjdgl"
}
