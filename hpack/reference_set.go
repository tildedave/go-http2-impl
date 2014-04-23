package hpack

type ReferenceSet struct {
	Entries []*HeaderField
}

func (refset *ReferenceSet) Add(ref *HeaderField) {
	refset.Entries = append(refset.Entries, ref)
}

func (refset *ReferenceSet) Remove(h *HeaderField) {
	idx := -1
	for i, ref := range refset.Entries {
		if ref == h {
			idx = i
		}
	}
	if idx == -1 {
		return
	}

	entries := refset.Entries
	refset.Entries = append(entries[0:idx], entries[idx + 1:]...)
}
