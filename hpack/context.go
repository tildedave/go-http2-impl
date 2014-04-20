package hpack

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
	IndexedHeaderMask = 0x80
	LiteralHeaderIncrementalIndexMask = 0x40
	LiteralHeaderNeverIndexMask = 0x10
	EncodingContextUpdateMask = 0x20
	LiteralHeaderNoIndexingMask = 0x00
)

func (context *EncodingContext) Decode(wire string) HeaderSet {
	headers := []HeaderField{}
	wireBytes := []byte(wire)

	if wireBytes[0] & IndexedHeaderMask == IndexedHeaderMask {
		index := wireBytes[0] & 0x4F
		headers = append(headers, context.HeaderTable.HeaderAt(int(index)))
	}

	return HeaderSet{ headers }
}
