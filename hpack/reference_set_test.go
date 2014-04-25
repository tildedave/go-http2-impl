package hpack

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestContains(t *testing.T) {
	h := HeaderField{ ":status", "302" }
	h2 := HeaderField{ ":authority", "www.example.com" }

	refset := NewReferenceSet()
	refset.Add(&h)

	assert.True(t, refset.Contains(&h))
	assert.False(t, refset.Contains(&h2))
}

func TestContainsHeader(t *testing.T) {
	h := HeaderField{ ":status", "302" }
	h2 := HeaderField{ ":status", "302" }
	h3 := HeaderField{ ":method", "POST" }

	refset := NewReferenceSet()
	refset.Add(&h)

	assert.True(t, refset.ContainsHeader(h2))
	assert.False(t, refset.ContainsHeader(h3))
}

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
