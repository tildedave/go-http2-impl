package hpack

type ReferenceSet struct {
	Entries map[*HeaderField]int
}

func NewReferenceSet() *ReferenceSet {
	refset := ReferenceSet{}
	refset.Entries = make(map[*HeaderField]int)

	return &refset
}

func (refset *ReferenceSet) Add(ref *HeaderField) {
	refset.Entries[ref] = 1
}

func (refset *ReferenceSet) Remove(h *HeaderField) {
	delete(refset.Entries, h)
}

func (refset *ReferenceSet) Clear() {
	refset.Entries = make(map[*HeaderField]int)
}

func (refset *ReferenceSet) Contains(ref *HeaderField) bool {
	return refset.Entries[ref] != 0
}

func (refset *ReferenceSet) ContainsHeader(h HeaderField) bool {
	for refHeader, _ := range refset.Entries {
		if *refHeader == h {
			return true
		}
	}

	return false
}
