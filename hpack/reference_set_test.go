package hpack

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddToReferenceSet(t *testing.T) {
	h := HeaderField{ ":status", "302" }
	refset := ReferenceSet{}
	refset.Add(&h)

	assert.Equal(t, len(refset.Entries), 1)
	assert.Equal(t, &h, refset.Entries[0])
}

func TestRemoveFromReferenceSet(t *testing.T) {
	h := HeaderField{ ":status", "302" }
	refset := ReferenceSet{ []*HeaderField{ &h } }
	refset.Remove(&h)

	assert.Equal(t, len(refset.Entries), 0)
}
