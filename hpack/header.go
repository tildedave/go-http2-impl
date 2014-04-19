package hpack

type HeaderField struct {
	Name  string
	Value string
}

type HeaderTable struct {
	Entries []HeaderField
	MaxSize int
}

type ReferenceSet struct {
	Entries []HeaderField
}

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

type HeaderSet struct {
	Headers []HeaderField
}

var StaticTable = []HeaderField{
	HeaderField{":authority", ""},
	HeaderField{":method", "GET"},
	HeaderField{":method", "POST"},
	HeaderField{":path", "/"},
	HeaderField{":path", "/index.html"},
	HeaderField{":scheme", "http"},
	HeaderField{":scheme", "https"},
	HeaderField{":status", "200"},
	HeaderField{":status", "204"},
	HeaderField{":status", "206"},
	HeaderField{":status", "304"},
	HeaderField{":status", "400"},
	HeaderField{":status", "404"},
	HeaderField{":status", "500"},
	HeaderField{"accept-charset", ""},
	HeaderField{"accept-encoding", ""},
	HeaderField{"accept-language", ""},
	HeaderField{"accept-ranges", ""},
	HeaderField{"accept", ""},
	HeaderField{"access-control-allow-origin", ""},
	HeaderField{"age", ""},
	HeaderField{"allow", ""},
	HeaderField{"authorization", ""},
	HeaderField{"cache-control", ""},
	HeaderField{"content-disposition", ""},
	HeaderField{"content-encoding", ""},
	HeaderField{"content-language", ""},
	HeaderField{"content-length", ""},
	HeaderField{"content-location", ""},
	HeaderField{"content-range", ""},
	HeaderField{"content-type", ""},
	HeaderField{"cookie", ""},
	HeaderField{"date", ""},
	HeaderField{"etag", ""},
	HeaderField{"expect", ""},
	HeaderField{"expires", ""},
	HeaderField{"from", ""},
	HeaderField{"host", ""},
	HeaderField{"if-match", ""},
	HeaderField{"if-modified-since", ""},
	HeaderField{"if-none-match", ""},
	HeaderField{"if-range", ""},
	HeaderField{"if-unmodified-since", ""},
	HeaderField{"last-modified", ""},
	HeaderField{"link", ""},
	HeaderField{"location", ""},
	HeaderField{"max-forwards", ""},
	HeaderField{"proxy-authenticate", ""},
	HeaderField{"proxy-authorization", ""},
	HeaderField{"range", ""},
	HeaderField{"referer", ""},
	HeaderField{"refresh", ""},
	HeaderField{"retry-after", ""},
	HeaderField{"server", ""},
	HeaderField{"set-cookie", ""},
	HeaderField{"strict-transport-security", ""},
	HeaderField{"transfer-encoding", ""},
	HeaderField{"user-agent", ""},
	HeaderField{"vary", ""},
	HeaderField{"via", ""},
	HeaderField{"www-authenticate", ""},
}

func (t *HeaderTable) AddHeader(header HeaderField) {
	for _, table_h := range t.Entries {
		if table_h == header {
			return
		}
	}

	for ; t.Size() + header.Size() > t.MaxSize ; {
		t.Entries = t.Entries[0:len(t.Entries) - 1]
	}

	t.Entries = append([]HeaderField{ header }, t.Entries...)
}

func (t HeaderTable) ContainsHeader(h HeaderField) int {
	for idx, table_h := range t.Entries {
		if table_h == h {
			return idx + 1
		}
	}

	for idx, table_h := range StaticTable {
		if table_h == h {
			return idx + len(t.Entries) + 1
		}
	}

	return 0
}

func (t HeaderTable) ContainsName(name string) int {
	idx := t.ContainsHeader(HeaderField{name, ""})
	if idx != 0 {
		return idx
	}

	for idx, table_h := range StaticTable {
		if table_h.Name == name {
			return idx + len(t.Entries) + 1
		}
	}

	return 0
}

func (h HeaderField) Size() int {
	return len(h.Name) + len(h.Value) + 32
}

func (t HeaderTable) Size() int {
	size := 0
	for _, header := range t.Entries {
		size += header.Size()
	}
	return size
}

func (h HeaderField) Encode(context *EncodingContext) string {
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

func (hs HeaderSet) Encode(context *EncodingContext) string {
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
			encoded += h.Encode(context)

			refset.Entries = append(refset.Entries, h)
		}
	}
	return encoded
}

// TODO: headers of arbitrary length with integer encoding algorithm
// TODO: header table size -- need header table representation that more closely matches
// TODO: is this the right representation for emission?
// TODO: hpack test cases https://github.com/http2jp/hpack-test-case

func Decode(headers string, table *HeaderTable) ([]HeaderField, int) {
	var decodedHeaders []HeaderField

	return decodedHeaders, 0
}
