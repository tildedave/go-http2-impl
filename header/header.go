package header

type HeaderField struct {
	Name  string
	Value string
}

type HeaderTable struct {
	Entries []HeaderField
}

type ReferenceSet struct {
	Entries []*HeaderField
}

type HeaderSet struct {
	Headers []HeaderField
}

var StaticTable = map[int]HeaderField{
	1:  {":authority", ""},
	2:  {":method", "GET"},
	3:  {":method", "POST"},
	4:  {":path", "/"},
	5:  {":path", "/index.html"},
	6:  {":scheme", "http"},
	7:  {":scheme", "https"},
	8:  {":status", "200"},
	9:  {":status", "204"},
	10: {":status", "206"},
	11: {":status", "304"},
	12: {":status", "400"},
	13: {":status", "404"},
	14: {":status", "500"},
	15: {"accept-charset", ""},
	16: {"accept-encoding", ""},
	17: {"accept-language", ""},
	18: {"accept-ranges", ""},
	19: {"accept", ""},
	20: {"access-control-allow-origin", ""},
	21: {"age", ""},
	22: {"allow", ""},
	23: {"authorization", ""},
	24: {"cache-control", ""},
	25: {"content-disposition", ""},
	26: {"content-encoding", ""},
	27: {"content-language", ""},
	28: {"content-length", ""},
	29: {"content-location", ""},
	30: {"content-range", ""},
	31: {"content-type", ""},
	32: {"cookie", ""},
	33: {"date", ""},
	34: {"etag", ""},
	35: {"expect", ""},
	36: {"expires", ""},
	37: {"from", ""},
	38: {"host", ""},
	39: {"if-match", ""},
	40: {"if-modified-since", ""},
	41: {"if-none-match", ""},
	42: {"if-range", ""},
	43: {"if-unmodified-since", ""},
	44: {"last-modified", ""},
	45: {"link", ""},
	46: {"location", ""},
	47: {"max-forwards", ""},
	48: {"proxy-authenticate", ""},
	49: {"proxy-authorization", ""},
	50: {"range", ""},
	51: {"referer", ""},
	52: {"refresh", ""},
	53: {"retry-after", ""},
	54: {"server", ""},
	55: {"set-cookie", ""},
	56: {"strict-transport-security", ""},
	57: {"transfer-encoding", ""},
	58: {"user-agent", ""},
	59: {"vary", ""},
	60: {"via", ""},
	61: {"www-authenticate", ""},
}

var StaticTableReverse = map[HeaderField]int{
	HeaderField{":authority", ""}: 1,
	HeaderField{":method", "GET"}: 2,
	HeaderField{":method", "POST"}: 3,
	HeaderField{":path", "/"}: 4,
	HeaderField{":path", "/index.html"}: 5,
	HeaderField{":scheme", "http"}: 6,
	HeaderField{":scheme", "https"}: 7,
	HeaderField{":status", "200"}: 8,
	HeaderField{":status", "204"}: 9,
	HeaderField{":status", "206"}: 10,
	HeaderField{":status", "304"}: 11,
	HeaderField{":status", "400"}: 12,
	HeaderField{":status", "404"}: 13,
	HeaderField{":status", "500"}: 14,
	HeaderField{"accept-charset", ""}: 15,
	HeaderField{"accept-encoding", ""}: 16,
	HeaderField{"accept-language", ""}: 17,
	HeaderField{"accept-ranges", ""}: 18,
	HeaderField{"accept", ""}: 19,
	HeaderField{"access-control-allow-origin", ""}: 20,
	HeaderField{"age", ""}: 21,
	HeaderField{"allow", ""}: 22,
	HeaderField{"authorization", ""}: 23,
	HeaderField{"cache-control", ""}: 24,
	HeaderField{"content-disposition", ""}: 25,
	HeaderField{"content-encoding", ""}: 26,
	HeaderField{"content-language", ""}: 27,
	HeaderField{"content-length", ""}: 28,
	HeaderField{"content-location", ""}: 29,
	HeaderField{"content-range", ""}: 30,
	HeaderField{"content-type", ""}: 31,
	HeaderField{"cookie", ""}: 32,
	HeaderField{"date", ""}: 33,
	HeaderField{"etag", ""}: 34,
	HeaderField{"expect", ""}: 35,
	HeaderField{"expires", ""}: 36,
	HeaderField{"from", ""}: 37,
	HeaderField{"host", ""}: 38,
	HeaderField{"if-match", ""}: 39,
	HeaderField{"if-modified-since", ""}: 40,
	HeaderField{"if-none-match", ""}: 41,
	HeaderField{"if-range", ""}: 42,
	HeaderField{"if-unmodified-since", ""}: 43,
	HeaderField{"last-modified", ""}: 44,
	HeaderField{"link", ""}: 45,
	HeaderField{"location", ""}: 46,
	HeaderField{"max-forwards", ""}: 47,
	HeaderField{"proxy-authenticate", ""}: 48,
	HeaderField{"proxy-authorization", ""}: 49,
	HeaderField{"range", ""}: 50,
	HeaderField{"referer", ""}: 51,
	HeaderField{"refresh", ""}: 52,
	HeaderField{"retry-after", ""}: 53,
	HeaderField{"server", ""}: 54,
	HeaderField{"set-cookie", ""}: 55,
	HeaderField{"strict-transport-security", ""}: 56,
	HeaderField{"transfer-encoding", ""}: 57,
	HeaderField{"user-agent", ""}: 58,
	HeaderField{"vary", ""}: 59,
	HeaderField{"via", ""}: 60,
	HeaderField{"www-authenticate", ""}: 61,
}

func addHeaderToTable(header HeaderField, table *HeaderTable) {
	entries := make([]HeaderField, len(table.Entries) + 1)
	copy(entries[1:], table.Entries)

	entries[0] = header
	table.Entries = entries
}

func (h HeaderField) Encode(table *HeaderTable) string {
	var encodedHeaders []byte

	idx := StaticTableReverse[h]

	// Indexed name and value
	if idx != 0 {
		a := make([]byte, 1)
		a[0] = byte(idx + len(table.Entries))
		a[0] |= 0x80

		addHeaderToTable(h, table)
		encodedHeaders = append(encodedHeaders, a...)
		return string(encodedHeaders)
	}

	idx = StaticTableReverse[HeaderField{h.Name, ""}]

	// Indexed name, literal value
	if idx != 0 {
		a := make([]byte, 2)
		a[0] = byte(idx + len(table.Entries))
		a[0] |= 0x40
		a[1] = byte(len(h.Value))

		addHeaderToTable(h, table)
		encodedHeaders = append(encodedHeaders, a...)
		encodedHeaders = append(encodedHeaders, h.Value...)
		return string(encodedHeaders)
	}

	// Literal name, literal value
	encodedHeaders = append(encodedHeaders, 0x40)
	encodedHeaders = append(encodedHeaders, byte(len(h.Name)))
	encodedHeaders = append(encodedHeaders, h.Name...)
	encodedHeaders = append(encodedHeaders, byte(len(h.Value)))
	encodedHeaders = append(encodedHeaders, h.Value...)
	addHeaderToTable(h, table)

	return string(encodedHeaders)
}

// TODO: headers of arbitrary length with integer encoding algorithm
// TODO: header table size -- need header table representation that more closely matches
// TODO: reference set
// TODO: is this the right representation for emission?
// TODO: hpack test cases https://github.com/http2jp/hpack-test-case
func Decode(headers string, table *HeaderTable) ([]HeaderField, int) {
	var decodedHeaders []HeaderField

	return decodedHeaders, 0
}
