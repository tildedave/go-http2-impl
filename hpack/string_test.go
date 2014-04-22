package hpack

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEncodeStringLiteralThatFitsInPrefix(t *testing.T) {
	encoded := encodeLiteral("Call me Ishmael")

	assert.Equal(t, encoded, "\x0FCall me Ishmael")
}

func TestEncodeStringLiteralThatOverflowsPrefix(t *testing.T) {
	str := "Once upon a time and a very good time it was there was a moocow coming down along the road and this moocow that was coming down along the road met a nicens little boy named baby tuckoo...  His father told him that story: his father looked at him through a glass: he had a hairy face.  He was baby tuckoo. The moocow came down the road where Betty Byrne lived: she sold lemon platt. O, the wild rose blossoms On the little green place. He sang that song. That was his song. O, the green wothe botheth. When you wet the bed first it is warm then it gets cold. His mother put on the oilsheet. That had the queer smell."

	encoded := encodeLiteral(str)

	assert.Equal(t, encoded, "\x7f\xe8\x03" + str)
}

func TestDecodeStringLiteralThatFitsInPrefix(t *testing.T) {
	wire := []byte("\x0FCall me Ishmael")
	decoded := decodeLiteral(&wire)

	assert.Equal(t, decoded, "Call me Ishmael")
	assert.Equal(t, wire, []byte{})
}

func TestDecodeStringLiteralThatOverflowsPrefix(t *testing.T) {
	expectedStr := "Once upon a time and a very good time it was there was a moocow coming down along the road and this moocow that was coming down along the road met a nicens little boy named baby tuckoo...  His father told him that story: his father looked at him through a glass: he had a hairy face.  He was baby tuckoo. The moocow came down the road where Betty Byrne lived: she sold lemon platt. O, the wild rose blossoms On the little green place. He sang that song. That was his song. O, the green wothe botheth. When you wet the bed first it is warm then it gets cold. His mother put on the oilsheet. That had the queer smell."

	wire := []byte("\x7f\xe8\x03" + expectedStr)
	decoded := decodeLiteral(&wire)

	assert.Equal(t, decoded, expectedStr)
	assert.Equal(t, wire, []byte{})
}
