package hpack

import (
	"fmt"
)

type EncodingContext struct {
	HeaderTable HeaderTable
	ReferenceSet ReferenceSet
	Update struct {
		ReferenceSetEmptying bool
		MaximumHeaderTableSizeChange int
	}
}

func NewEncodingContext() *EncodingContext {
	context := &EncodingContext{}
	context.HeaderTable.MaxSize = 1024

	return context
}

func (context *EncodingContext) AddHeader(h HeaderField) {
	ref := context.HeaderTable.AddHeader(h)
	refset := &context.ReferenceSet

	refset.Entries = append(refset.Entries, ref)
}

func (context *EncodingContext) EncodeField(h HeaderField) string {
	var idx int

	table := &context.HeaderTable
	idx = table.ContainsHeader(h)
	if idx != 0 {
		a := make([]byte, 1)
		a[0] = byte(idx)
		a[0] |= 0x80

		table.AddHeader(h)
		encodedHeaders := string(a)
		return encodedHeaders
	}

	idx = table.ContainsName(h.Name)
	if idx != 0 {
		a := make([]byte, 2)
		a[0] = byte(idx)
		a[0] |= 0x40
		a[1] = byte(len(h.Value))

		table.AddHeader(h)
		encodedHeaders := string(a) + h.Value
		return encodedHeaders
	}

	// Literal name, literal value
	table.AddHeader(h)
	encodedHeaders := ""
	encodedHeaders += string(0x40)
	encodedHeaders += string(byte(len(h.Name)))
	encodedHeaders += h.Name
	encodedHeaders += string(byte(len(h.Value)))
	encodedHeaders += h.Value

	return string(encodedHeaders)
}

func (context *EncodingContext) Encode(hs HeaderSet) string {
	encoded := ""

	if context.Update.ReferenceSetEmptying {
		context.ReferenceSet = ReferenceSet{}
		context.Update.ReferenceSetEmptying = false
		encoded += "\x30"
	}

	refset := &context.ReferenceSet

	for _, h := range hs.Headers {
		mustEncode := true

		for _, refHeader := range refset.Entries {
			if *refHeader == h {
				mustEncode = false
			}
		}

		if mustEncode {
			encoded += context.EncodeField(h)
			// Not the correct way to do this
			refset.Entries = append(refset.Entries, &HeaderField{h.Name, h.Value})
		}
	}
	return encoded
}

const (
	IndexedMask = 0x80
	LiteralIndexedMask = 0x40
	LiteralNeverIndexedMask = 0x10
	EncodingContextUpdateMask = 0x20
	LiteralNoIndexMask = 0x00
)

func fmtIsNotUnused() {
	fmt.Println("line to not complain about unused fmt import")
}

func unpackLiteral(wireBytes *[]byte) (string) {
	length := int((*wireBytes)[0] & 0x4F)
	str := string((*wireBytes)[1:1 + length])

	*wireBytes = (*wireBytes)[1 + length:]

	return str
}

func (context *EncodingContext) Decode(wire string) HeaderSet {
	headers := []HeaderField{}
	table := &context.HeaderTable
	wireBytes := []byte(wire)

	for ; len(wireBytes) > 0 ; {
		if wireBytes[0] & IndexedMask == IndexedMask {
			index := wireBytes[0] & 0x4F
			header := table.HeaderAt(int(index))
			headers = append(headers, header)
			context.AddHeader(header)

			wireBytes = wireBytes[1: ]

			continue
		}

		if wireBytes[0] & LiteralIndexedMask == LiteralIndexedMask {
			nameIndex := wireBytes[0] & 0x2F
			wireBytes = wireBytes[1:]

			if nameIndex == byte(0) {
				name := unpackLiteral(&wireBytes)
				value := unpackLiteral(&wireBytes)

				header := HeaderField{ name, value }
				headers = append(headers, header)
				context.AddHeader(header)
			} else {
				nameHeader := table.HeaderAt(int(nameIndex))

				value := unpackLiteral(&wireBytes)

				header := HeaderField{ nameHeader.Name, value }

				headers = append(headers, header)
				context.AddHeader(header)
			}
		}
	}

	return HeaderSet{ headers }
}
