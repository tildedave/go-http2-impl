package hpack

import (
	"container/list"
	"errors"
	"fmt"
)

type EncodingContext struct {
	HeaderTable HeaderTable
	ReferenceSet *ReferenceSet
	Update struct {
		ReferenceSetEmptying bool
		MaximumHeaderTableSizeChange int
	}
}

func NewEncodingContext() *EncodingContext {
	context := &EncodingContext{}
	context.HeaderTable.MaxSize = 1024
	context.ReferenceSet = NewReferenceSet()

	return context
}

func (context *EncodingContext) AddHeader(h HeaderField) {
	ref := context.HeaderTable.AddHeader(h)
	if ref != nil {
		context.ReferenceSet.Add(ref)
	}
}

func (context *EncodingContext) EncodeField(h HeaderField) string {
	var idx int

	table := &context.HeaderTable
	idx = table.ContainsHeader(h)
	if idx != 0 {
		a := []byte(encodeInteger(idx, 7))
		a[0] |= 0x80

		context.AddHeader(h)

		return string(a)
	}

	idx = table.ContainsName(h.Name)
	if idx != 0 {
		a := []byte(encodeInteger(idx, 6))
		a[0] |= 0x40

		context.AddHeader(h)
		return string(a) + encodeLiteral(h.Value)
	}

	// Literal name, literal value
	context.AddHeader(h)
	return string(0x40) + encodeLiteral(h.Name) + encodeLiteral(h.Value)
}

func (context *EncodingContext) encodeMissingHeaders(hs HeaderSet) string {
	// Headers absent from reference set are encoded as their
	// indexed encoding

	encoded := ""
	refset := context.ReferenceSet

	absent := list.New()
	for refHeader, _ := range refset.Entries {
		present := false
		for _, h := range hs.Headers {
			if *refHeader == h {
				present = true
				break
			}
		}
		if !present {
			absent.PushBack(refHeader)
		}
	}

	for e := absent.Front(); e != nil; e = e.Next() {
		// find which index this was at in the header
		// table, don't need to pop it from there tho

		refHeader := e.Value.(*HeaderField)
		idx := context.HeaderTable.ContainsHeader(*refHeader)
		a := []byte(encodeInteger(idx, 7))
		a[0] |= 0x80

		encoded += string(a)

		context.ReferenceSet.Remove(refHeader)
	}

	return encoded
}

func (context *EncodingContext) Encode(hs HeaderSet) string {
	encoded := ""

	// TODO: ideally encodeMissingHeaders detects this
	if context.Update.ReferenceSetEmptying {
		context.ReferenceSet = NewReferenceSet()
		context.Update.ReferenceSetEmptying = false
		encoded += "\x30"
	}

	encoded += context.encodeMissingHeaders(hs)

	refset := context.ReferenceSet
	for _, h := range hs.Headers {
		mustEncode := true

		for refHeader, _ := range refset.Entries {
			if *refHeader == h {
				mustEncode = false
			}
		}

		if mustEncode {
			encoded += context.EncodeField(h)
		}
	}
	return encoded
}

const (
	IndexedMask = 0x80
	LiteralIndexedMask = 0x40
	LiteralNeverIndexMask = 0x10
	ContextUpdateMask = 0x20
	LiteralNoIndexMask = 0x00
)

func fmtIsNotUnused() {
	fmt.Println("line to not complain about unused fmt import")
}

func decodeLiteralHeader(wireBytes *[]byte, indexBits uint, table *HeaderTable) (HeaderField) {
	nameIndex := decodeInteger(wireBytes, indexBits)
	if nameIndex == uint(0) {
		name := decodeLiteral(wireBytes)
		value := decodeLiteral(wireBytes)

		return HeaderField{ name, value }
	}

	nameHeader := table.HeaderAt(int(nameIndex))
	value := decodeLiteral(wireBytes)
	return HeaderField{ nameHeader.Name, value }
}

func (context *EncodingContext) Decode(wire string) (hs HeaderSet, err error) {
	headers := []HeaderField{}
	wireBytes := []byte(wire)

	table := &context.HeaderTable
	refset := context.ReferenceSet

	for ; len(wireBytes) > 0 ; {
		if wireBytes[0] & ContextUpdateMask == ContextUpdateMask {
			if wireBytes[0] & 0x30 == 0x30 {
				// empty reference set
				refset.Clear()
			}
			wireBytes = wireBytes[1: ]
			continue
		}

		if wireBytes[0] & IndexedMask == IndexedMask {
			index := decodeInteger(&wireBytes, 7)
			header := table.HeaderAt(int(index))
			headers = append(headers, header)
			context.AddHeader(header)

			continue
		}

		if wireBytes[0] & LiteralIndexedMask == LiteralIndexedMask {
			header := decodeLiteralHeader(&wireBytes, 6, table)
			headers = append(headers, header)
			context.AddHeader(header)
			continue
		}

		if wireBytes[0] & LiteralNeverIndexMask == LiteralNeverIndexMask {
			header := decodeLiteralHeader(&wireBytes, 4, table)
			headers = append(headers, header)
			continue
		}

		if wireBytes[0] & LiteralNoIndexMask == LiteralNoIndexMask {
			header := decodeLiteralHeader(&wireBytes, 4, table)
			headers = append(headers, header)
			continue
		}

		return HeaderSet{}, errors.New("Could not decode")
	}

	for h, _ := range refset.Entries {
		found := false

		for _, emitted := range headers {

			if emitted == *h {
				found = true
			}
		}

		if !found {
			headers = append(headers, *h)
		}
	}

	return HeaderSet{ headers }, nil
}
