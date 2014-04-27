package hpack

import (
	"fmt"
	"encoding/binary"
	"errors"
)

type HuffmanCode struct {
	bits uint32
	bitLength uint
}

var HuffmanTable = map[byte]HuffmanCode{
	byte(0): HuffmanCode{0x03ffffba, 26},
	byte(1): HuffmanCode{0x03ffffbb, 26},
	byte(2): HuffmanCode{0x03ffffbc, 26},
	byte(3): HuffmanCode{0x03ffffbd, 26},
	byte(4): HuffmanCode{0x03ffffbe, 26},
	byte(5): HuffmanCode{0x03ffffbf, 26},
	byte(6): HuffmanCode{0x03ffffc0, 26},
	byte(7): HuffmanCode{0x03ffffc1, 26},
	byte(8): HuffmanCode{0x03ffffc2, 26},
	byte(9): HuffmanCode{0x03ffffc3, 26},
	byte(10): HuffmanCode{0x03ffffc4, 26},
	byte(11): HuffmanCode{0x03ffffc5, 26},
	byte(12): HuffmanCode{0x03ffffc6, 26},
	byte(13): HuffmanCode{0x03ffffc7, 26},
	byte(14): HuffmanCode{0x03ffffc8, 26},
	byte(15): HuffmanCode{0x03ffffc9, 26},
	byte(16): HuffmanCode{0x03ffffca, 26},
	byte(17): HuffmanCode{0x03ffffcb, 26},
	byte(18): HuffmanCode{0x03ffffcc, 26},
	byte(19): HuffmanCode{0x03ffffcd, 26},
	byte(20): HuffmanCode{0x03ffffce, 26},
	byte(21): HuffmanCode{0x03ffffcf, 26},
	byte(22): HuffmanCode{0x03ffffd0, 26},
	byte(23): HuffmanCode{0x03ffffd1, 26},
	byte(24): HuffmanCode{0x03ffffd2, 26},
	byte(25): HuffmanCode{0x03ffffd3, 26},
	byte(26): HuffmanCode{0x03ffffd4, 26},
	byte(27): HuffmanCode{0x03ffffd5, 26},
	byte(28): HuffmanCode{0x03ffffd6, 26},
	byte(29): HuffmanCode{0x03ffffd7, 26},
	byte(30): HuffmanCode{0x03ffffd8, 26},
	byte(31): HuffmanCode{0x03ffffd9, 26},
	byte(32): HuffmanCode{0x06, 5},
	byte(33): HuffmanCode{0x1ffc, 13},
	byte(34): HuffmanCode{0x01f0, 9},
	byte(35): HuffmanCode{0x3ffc, 14},
	byte(36): HuffmanCode{0x7ffc, 15},
	byte(37): HuffmanCode{0x1e, 6},
	byte(38): HuffmanCode{0x64, 7},
	byte(39): HuffmanCode{0x1ffd, 13},
	byte(40): HuffmanCode{0x03fa, 10},
	byte(41): HuffmanCode{0x01f1, 9},
	byte(42): HuffmanCode{0x03fb, 10},
	byte(43): HuffmanCode{0x03fc, 10},
	byte(44): HuffmanCode{0x65, 7},
	byte(45): HuffmanCode{0x66, 7},
	byte(46): HuffmanCode{0x1f, 6},
	byte(47): HuffmanCode{0x07, 5},
	byte(48): HuffmanCode{0x00, 4},
	byte(49): HuffmanCode{0x01, 4},
	byte(50): HuffmanCode{0x02, 4},
	byte(51): HuffmanCode{0x08, 5},
	byte(52): HuffmanCode{0x20, 6},
	byte(53): HuffmanCode{0x21, 6},
	byte(54): HuffmanCode{0x22, 6},
	byte(55): HuffmanCode{0x23, 6},
	byte(56): HuffmanCode{0x24, 6},
	byte(57): HuffmanCode{0x25, 6},
	byte(58): HuffmanCode{0x26, 6},
	byte(59): HuffmanCode{0xec, 8},
	byte(60): HuffmanCode{0x01fffc, 17},
	byte(61): HuffmanCode{0x27, 6},
	byte(62): HuffmanCode{0x7ffd, 15},
	byte(63): HuffmanCode{0x03fd, 10},
	byte(64): HuffmanCode{0x7ffe, 15},
	byte(65): HuffmanCode{0x67, 7},
	byte(66): HuffmanCode{0xed, 8},
	byte(67): HuffmanCode{0xee, 8},
	byte(68): HuffmanCode{0x68, 7},
	byte(69): HuffmanCode{0xef, 8},
	byte(70): HuffmanCode{0x69, 7},
	byte(71): HuffmanCode{0x6a, 7},
	byte(72): HuffmanCode{0x01f2, 9},
	byte(73): HuffmanCode{0xf0, 8},
	byte(74): HuffmanCode{0x01f3, 9},
	byte(75): HuffmanCode{0x01f4, 9},
	byte(76): HuffmanCode{0x01f5, 9},
	byte(77): HuffmanCode{0x6b, 7},
	byte(78): HuffmanCode{0x6c, 7},
	byte(79): HuffmanCode{0xf1, 8},
	byte(80): HuffmanCode{0xf2, 8},
	byte(81): HuffmanCode{0x01f6, 9},
	byte(82): HuffmanCode{0x01f7, 9},
	byte(83): HuffmanCode{0x6d, 7},
	byte(84): HuffmanCode{0x28, 6},
	byte(85): HuffmanCode{0xf3, 8},
	byte(86): HuffmanCode{0x01f8, 9},
	byte(87): HuffmanCode{0x01f9, 9},
	byte(88): HuffmanCode{0xf4, 8},
	byte(89): HuffmanCode{0x01fa, 9},
	byte(90): HuffmanCode{0x01fb, 9},
	byte(91): HuffmanCode{0x07fc, 11},
	byte(92): HuffmanCode{0x03ffffda, 26},
	byte(93): HuffmanCode{0x07fd, 11},
	byte(94): HuffmanCode{0x3ffd, 14},
	byte(95): HuffmanCode{0x6e, 7},
	byte(96): HuffmanCode{0x03fffe, 18},
	byte(97): HuffmanCode{0x09, 5},
	byte(98): HuffmanCode{0x6f, 7},
	byte(99): HuffmanCode{0x0a, 5},
	byte(100): HuffmanCode{0x29, 6},
	byte(101): HuffmanCode{0x0b, 5},
	byte(102): HuffmanCode{0x70, 7},
	byte(103): HuffmanCode{0x2a, 6},
	byte(104): HuffmanCode{0x2b, 6},
	byte(105): HuffmanCode{0x0c, 5},
	byte(106): HuffmanCode{0xf5, 8},
	byte(107): HuffmanCode{0xf6, 8},
	byte(108): HuffmanCode{0x2c, 6},
	byte(109): HuffmanCode{0x2d, 6},
	byte(110): HuffmanCode{0x2e, 6},
	byte(111): HuffmanCode{0x0d, 5},
	byte(112): HuffmanCode{0x2f, 6},
	byte(113): HuffmanCode{0x01fc, 9},
	byte(114): HuffmanCode{0x30, 6},
	byte(115): HuffmanCode{0x31, 6},
	byte(116): HuffmanCode{0x0e, 5},
	byte(117): HuffmanCode{0x71, 7},
	byte(118): HuffmanCode{0x72, 7},
	byte(119): HuffmanCode{0x73, 7},
	byte(120): HuffmanCode{0x74, 7},
	byte(121): HuffmanCode{0x75, 7},
	byte(122): HuffmanCode{0xf7, 8},
	byte(123): HuffmanCode{0x01fffd, 17},
	byte(124): HuffmanCode{0x0ffc, 12},
	byte(125): HuffmanCode{0x01fffe, 17},
	byte(126): HuffmanCode{0x0ffd, 12},
	byte(127): HuffmanCode{0x03ffffdb, 26},
	byte(128): HuffmanCode{0x03ffffdc, 26},
	byte(129): HuffmanCode{0x03ffffdd, 26},
	byte(130): HuffmanCode{0x03ffffde, 26},
	byte(131): HuffmanCode{0x03ffffdf, 26},
	byte(132): HuffmanCode{0x03ffffe0, 26},
	byte(133): HuffmanCode{0x03ffffe1, 26},
	byte(134): HuffmanCode{0x03ffffe2, 26},
	byte(135): HuffmanCode{0x03ffffe3, 26},
	byte(136): HuffmanCode{0x03ffffe4, 26},
	byte(137): HuffmanCode{0x03ffffe5, 26},
	byte(138): HuffmanCode{0x03ffffe6, 26},
	byte(139): HuffmanCode{0x03ffffe7, 26},
	byte(140): HuffmanCode{0x03ffffe8, 26},
	byte(141): HuffmanCode{0x03ffffe9, 26},
	byte(142): HuffmanCode{0x03ffffea, 26},
	byte(143): HuffmanCode{0x03ffffeb, 26},
	byte(144): HuffmanCode{0x03ffffec, 26},
	byte(145): HuffmanCode{0x03ffffed, 26},
	byte(146): HuffmanCode{0x03ffffee, 26},
	byte(147): HuffmanCode{0x03ffffef, 26},
	byte(148): HuffmanCode{0x03fffff0, 26},
	byte(149): HuffmanCode{0x03fffff1, 26},
	byte(150): HuffmanCode{0x03fffff2, 26},
	byte(151): HuffmanCode{0x03fffff3, 26},
	byte(152): HuffmanCode{0x03fffff4, 26},
	byte(153): HuffmanCode{0x03fffff5, 26},
	byte(154): HuffmanCode{0x03fffff6, 26},
	byte(155): HuffmanCode{0x03fffff7, 26},
	byte(156): HuffmanCode{0x03fffff8, 26},
	byte(157): HuffmanCode{0x03fffff9, 26},
	byte(158): HuffmanCode{0x03fffffa, 26},
	byte(159): HuffmanCode{0x03fffffb, 26},
	byte(160): HuffmanCode{0x03fffffc, 26},
	byte(161): HuffmanCode{0x03fffffd, 26},
	byte(162): HuffmanCode{0x03fffffe, 26},
	byte(163): HuffmanCode{0x03ffffff, 26},
	byte(164): HuffmanCode{0x01ffff80, 25},
	byte(165): HuffmanCode{0x01ffff81, 25},
	byte(166): HuffmanCode{0x01ffff82, 25},
	byte(167): HuffmanCode{0x01ffff83, 25},
	byte(168): HuffmanCode{0x01ffff84, 25},
	byte(169): HuffmanCode{0x01ffff85, 25},
	byte(170): HuffmanCode{0x01ffff86, 25},
	byte(171): HuffmanCode{0x01ffff87, 25},
	byte(172): HuffmanCode{0x01ffff88, 25},
	byte(173): HuffmanCode{0x01ffff89, 25},
	byte(174): HuffmanCode{0x01ffff8a, 25},
	byte(175): HuffmanCode{0x01ffff8b, 25},
	byte(176): HuffmanCode{0x01ffff8c, 25},
	byte(177): HuffmanCode{0x01ffff8d, 25},
	byte(178): HuffmanCode{0x01ffff8e, 25},
	byte(179): HuffmanCode{0x01ffff8f, 25},
	byte(180): HuffmanCode{0x01ffff90, 25},
	byte(181): HuffmanCode{0x01ffff91, 25},
	byte(182): HuffmanCode{0x01ffff92, 25},
	byte(183): HuffmanCode{0x01ffff93, 25},
	byte(184): HuffmanCode{0x01ffff94, 25},
	byte(185): HuffmanCode{0x01ffff95, 25},
	byte(186): HuffmanCode{0x01ffff96, 25},
	byte(187): HuffmanCode{0x01ffff97, 25},
	byte(188): HuffmanCode{0x01ffff98, 25},
	byte(189): HuffmanCode{0x01ffff99, 25},
	byte(190): HuffmanCode{0x01ffff9a, 25},
	byte(191): HuffmanCode{0x01ffff9b, 25},
	byte(192): HuffmanCode{0x01ffff9c, 25},
	byte(193): HuffmanCode{0x01ffff9d, 25},
	byte(194): HuffmanCode{0x01ffff9e, 25},
	byte(195): HuffmanCode{0x01ffff9f, 25},
	byte(196): HuffmanCode{0x01ffffa0, 25},
	byte(197): HuffmanCode{0x01ffffa1, 25},
	byte(198): HuffmanCode{0x01ffffa2, 25},
	byte(199): HuffmanCode{0x01ffffa3, 25},
	byte(200): HuffmanCode{0x01ffffa4, 25},
	byte(201): HuffmanCode{0x01ffffa5, 25},
	byte(202): HuffmanCode{0x01ffffa6, 25},
	byte(203): HuffmanCode{0x01ffffa7, 25},
	byte(204): HuffmanCode{0x01ffffa8, 25},
	byte(205): HuffmanCode{0x01ffffa9, 25},
	byte(206): HuffmanCode{0x01ffffaa, 25},
	byte(207): HuffmanCode{0x01ffffab, 25},
	byte(208): HuffmanCode{0x01ffffac, 25},
	byte(209): HuffmanCode{0x01ffffad, 25},
	byte(210): HuffmanCode{0x01ffffae, 25},
	byte(211): HuffmanCode{0x01ffffaf, 25},
	byte(212): HuffmanCode{0x01ffffb0, 25},
	byte(213): HuffmanCode{0x01ffffb1, 25},
	byte(214): HuffmanCode{0x01ffffb2, 25},
	byte(215): HuffmanCode{0x01ffffb3, 25},
	byte(216): HuffmanCode{0x01ffffb4, 25},
	byte(217): HuffmanCode{0x01ffffb5, 25},
	byte(218): HuffmanCode{0x01ffffb6, 25},
	byte(219): HuffmanCode{0x01ffffb7, 25},
	byte(220): HuffmanCode{0x01ffffb8, 25},
	byte(221): HuffmanCode{0x01ffffb9, 25},
	byte(222): HuffmanCode{0x01ffffba, 25},
	byte(223): HuffmanCode{0x01ffffbb, 25},
	byte(224): HuffmanCode{0x01ffffbc, 25},
	byte(225): HuffmanCode{0x01ffffbd, 25},
	byte(226): HuffmanCode{0x01ffffbe, 25},
	byte(227): HuffmanCode{0x01ffffbf, 25},
	byte(228): HuffmanCode{0x01ffffc0, 25},
	byte(229): HuffmanCode{0x01ffffc1, 25},
	byte(230): HuffmanCode{0x01ffffc2, 25},
	byte(231): HuffmanCode{0x01ffffc3, 25},
	byte(232): HuffmanCode{0x01ffffc4, 25},
	byte(233): HuffmanCode{0x01ffffc5, 25},
	byte(234): HuffmanCode{0x01ffffc6, 25},
	byte(235): HuffmanCode{0x01ffffc7, 25},
	byte(236): HuffmanCode{0x01ffffc8, 25},
	byte(237): HuffmanCode{0x01ffffc9, 25},
	byte(238): HuffmanCode{0x01ffffca, 25},
	byte(239): HuffmanCode{0x01ffffcb, 25},
	byte(240): HuffmanCode{0x01ffffcc, 25},
	byte(241): HuffmanCode{0x01ffffcd, 25},
	byte(242): HuffmanCode{0x01ffffce, 25},
	byte(243): HuffmanCode{0x01ffffcf, 25},
	byte(244): HuffmanCode{0x01ffffd0, 25},
	byte(245): HuffmanCode{0x01ffffd1, 25},
	byte(246): HuffmanCode{0x01ffffd2, 25},
	byte(247): HuffmanCode{0x01ffffd3, 25},
	byte(248): HuffmanCode{0x01ffffd4, 25},
	byte(249): HuffmanCode{0x01ffffd5, 25},
	byte(250): HuffmanCode{0x01ffffd6, 25},
	byte(251): HuffmanCode{0x01ffffd7, 25},
	byte(252): HuffmanCode{0x01ffffd8, 25},
	byte(253): HuffmanCode{0x01ffffd9, 25},
	byte(254): HuffmanCode{0x01ffffda, 25},
	byte(255): HuffmanCode{0x01ffffdb, 25},
}

var HuffmanEOS = HuffmanCode{0x01ffffdc, 25}

type huffmanNode struct {
	value uint8
	isLeaf bool
	left, right *huffmanNode
}

func newHuffmanNode() *huffmanNode {
	return &huffmanNode{0, false, nil, nil}
}

func insertCode(parent *huffmanNode, code HuffmanCode, val uint8) {
	if code.bitLength == 0 {
		// this must be the place
		parent.value = val
		parent.isLeaf = true
	} else {
		// determine if we need to add this to the left (0) or right (1)
		// of the parent

		var next *huffmanNode

		code.bitLength -= 1
		mask := uint32(1 << code.bitLength)
		if code.bits & mask == mask {
			// right (1)
			next = parent.right
			if next == nil {
				next = newHuffmanNode()
				parent.right = next
			}
		} else {
			// left (0)
			next = parent.left
			if next == nil {
				next = newHuffmanNode()
				parent.left = next
			}
		}
		insertCode(next, code, val)
	}
}

func lookupCode(parent *huffmanNode, code HuffmanCode) *huffmanNode {
	if parent == nil {
		return nil
	}

	if parent.isLeaf {
		return parent
	}

	if code.bitLength == 0 {
		return nil
	}

	code.bitLength -= 1
	mask := uint32(1 << code.bitLength)

	if code.bits & mask == mask {
		return lookupCode(parent.right, code)
	} else {
		return lookupCode(parent.left, code)
	}
}

var huffmanTree *huffmanNode

func buildHuffmanTree() {
	// goal is to build a relationship between huffman codes
	// decode will pass in a prefix, one byte at a time
	// prefix navigates through < and >
	// if it hits a leaf node (==), done

	huffmanTree = newHuffmanNode()

	for i, code := range HuffmanTable {
		insertCode(huffmanTree, code, i)
	}
}

func decodeHuffmanHelper(wire *[]byte, parent *huffmanNode) (string, error) {
	var code HuffmanCode
	var node *huffmanNode
	var remainingInOctet uint8

	encoded := ""

	remainingInOctet = 8
	a := uint8((*wire)[0])
	*wire = (*wire)[1:]

	for ; len(*wire) > 0 ; {
		code = HuffmanCode{}
		node = nil

		// 0xEE
		// 1110 1110

		for ; node == nil ; {
			if remainingInOctet == 0 {
				// consume next
				if len(*wire) == 0 {
					// confirm EOS

					eos := HuffmanEOS.bits >> uint(HuffmanEOS.bitLength - code.bitLength)
					if eos & code.bits == eos {
						return encoded, nil
					}

					return "", errors.New("Sequence did not terminate with EOS marker")
				}

				remainingInOctet = 8
				a = uint8((*wire)[0])
				*wire = (*wire)[1:]
			}

			code.bitLength++
			remainingInOctet -= 1
			nextBit := uint32(a >> (remainingInOctet)) & 0x0001
			code.bits = (code.bits << 1) | nextBit

			node = lookupCode(huffmanTree, code)
		}

		encoded += string(node.value)
	}

	return encoded, nil
}

func DecodeHuffman(wire *[]byte) (string, error) {
	fmt.Println("decode")
	return decodeHuffmanHelper(wire, huffmanTree)
}

func EncodeHuffman(str string) string {
	var overflow string

	partialCode := HuffmanCode{}
	encoded := ""

	for _, b := range []byte(str) {
		overflow, partialCode = combineHuffman(partialCode, HuffmanTable[b])
		encoded += overflow
	}

	if partialCode.bitLength > 0 {
		overflow, partialCode = combineHuffman(partialCode, HuffmanEOS)
		encoded += overflow
	}

	return encoded
}

func combineHuffman(a HuffmanCode, b HuffmanCode) (string, HuffmanCode) {
	// align a to MSB
	// fill in b
	// if this overflows 32 bits, return a string
	// return the combined huffman encoding

	// 1f 6 bits
	// 00011111
	// first 2 in octet are garbage

	// 12 bits
	// in table as
	// 0000 0001 1111 1111
	// align to
	// 0001 1111 1111 0000

	// align a to the start of the octet
	// but if we aren't at the start of octet, align a otherwise

	// 32 bits
	paddedA := a.bits << uint(32 - a.bitLength)

	// we now have 32 - a.bitLength bits left in our uint
	// two options:
	// 1) b fits in the rest of the uint32, in which case we return another code
	//    (code must be aligned back to LSB)
	//  2) b overflows uint32, in which case we
	//   i) return the uint32
	//   ii) return a code aligned to LSB

	if b.bitLength < 32 - a.bitLength {
		// b fits
		paddedB := b.bits << uint(32 - a.bitLength - b.bitLength)
		remaining := uint(32 - a.bitLength - b.bitLength)

		return "", HuffmanCode{
			(paddedA | paddedB) >> remaining,
			a.bitLength + b.bitLength,
		}
	} else {
		// overflow
		overflow := make([]byte, 4)
		overflowBits := uint(b.bitLength - (32 - a.bitLength))

		binary.BigEndian.PutUint32(overflow, paddedA | (b.bits >> overflowBits))

		return string(overflow), HuffmanCode{
			b.bits & ((1 << overflowBits) - 1),
			overflowBits,
		}
	}
}
