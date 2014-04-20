package hpack

type HeaderTable struct {
	Entries []HeaderField
	MaxSize int
}

func (t *HeaderTable) AddHeader(header HeaderField) {
	for _, table_h := range t.Entries {
		if table_h == header {
			return
		}
	}

	for ; t.Size() + header.Size() > t.MaxSize ; {
		t.Entries = t.Entries[0:len(t.Entries) - 1]
	}

	t.Entries = append([]HeaderField{ header }, t.Entries...)
}

func (t HeaderTable) ContainsHeader(h HeaderField) int {
	for idx, table_h := range t.Entries {
		if table_h == h {
			return idx + 1
		}
	}

	for idx, table_h := range StaticTable {
		if table_h == h {
			return idx + len(t.Entries) + 1
		}
	}

	return 0
}

func (t HeaderTable) ContainsName(name string) int {
	idx := t.ContainsHeader(HeaderField{name, ""})
	if idx != 0 {
		return idx
	}

	for idx, table_h := range StaticTable {
		if table_h.Name == name {
			return idx + len(t.Entries) + 1
		}
	}

	return 0
}

func (t HeaderTable) Size() int {
	size := 0
	for _, header := range t.Entries {
		size += header.Size()
	}
	return size
}
