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
			if refHeader == h {
				mustEncode = false
			}
		}

		if mustEncode {
			encoded += context.EncodeField(h)

			refset.Entries = append(refset.Entries, h)
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

func unpackLiteral(wireBytes []byte, i int) (string, int) {
	len := int(wireBytes[i] & 0x4F)
	str := string(wireBytes[i + 1:i + 1 + len])

	return str, i + 1 + len
}

func (context *EncodingContext) Decode(wire string) HeaderSet {
	headers := []HeaderField{}
	table := context.HeaderTable
	wireBytes := []byte(wire)

	for i := 0; i < len(wireBytes); i++ {
		if wireBytes[i] & IndexedMask == IndexedMask {
			index := wireBytes[i] & 0x4F
			header := table.HeaderAt(int(index))
			headers = append(headers, header)
			table.AddHeader(header)

			continue
		}

		if wireBytes[i] & LiteralIndexedMask == LiteralIndexedMask {


			nameIndex := wireBytes[i] & 0x2F

			if nameIndex == byte(0) {
				var name, value string

				name, i = unpackLiteral(wireBytes, i + 1)
				value, i = unpackLiteral(wireBytes, i)

				header := HeaderField{ name, value }
				headers = append(headers, header)
				table.AddHeader(header)

				i += 1000
			} else {
				var value string

				nameHeader := table.HeaderAt(int(nameIndex))

				value, i = unpackLiteral(wireBytes, i + 1)

				header := HeaderField{ nameHeader.Name, value }

				headers = append(headers, header)
				table.AddHeader(header)
			}
		}
	}

	return HeaderSet{ headers }
}
