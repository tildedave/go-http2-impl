package hpack

type HeaderField struct {
	Name  string
	Value string
}

func (h HeaderField) Size() int {
	return len(h.Name) + len(h.Value) + 32
}
