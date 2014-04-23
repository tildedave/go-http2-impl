package hpack

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddToReferenceSet(t *testing.T) {
	h := HeaderField{ ":status", "302" }
	refset := NewReferenceSet()
	refset.Add(&h)

	assert.Equal(t, len(refset.Entries), 1)
	assert.Equal(t, 1, refset.Entries[&h])
}

func TestRemoveFromReferenceSet(t *testing.T) {
	h := HeaderField{ ":status", "302" }
	refset := NewReferenceSet()
	refset.Entries[&h] = 1
	refset.Remove(&h)

	assert.Equal(t, len(refset.Entries), 0)
}
