package hpack

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHeaderAtReturnsHeaderInTable(t *testing.T) {
	table := HeaderTable{}
	headerField := HeaderField{ ":method", "GET" }

	table.Entries = append(table.Entries, headerField)

	assert.Equal(t, table.HeaderAt(1), headerField)
}

func TestHeaderAtReturnsStaticHeader(t *testing.T) {
	table := HeaderTable{}

	assert.Equal(t, table.HeaderAt(3), HeaderField{ ":method", "POST" })
}

func TestHeaderAtReturnsStaticHeaderWithNonEmptyHeaderTable(t *testing.T) {
	table := HeaderTable{}
	table.Entries = append(table.Entries, HeaderField{ ":authority", "mine" })

	assert.Equal(t, table.HeaderAt(4), HeaderField{ ":method", "POST" })
}
