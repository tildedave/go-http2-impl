package hpack

import (
	"encoding/binary"
)

type HuffmanCode struct {
	bits string
	bitLength uint
}

var HuffmanTable = map[byte]HuffmanCode{
	byte(0): HuffmanCode{  "\x03\xff\xff\xba",  26 },
	byte(1): HuffmanCode{  "\x03\xff\xff\xbb",  26 },
	byte(2): HuffmanCode{  "\x03\xff\xff\xbc",  26 },
	byte(3): HuffmanCode{  "\x03\xff\xff\xbd",  26 },
	byte(4): HuffmanCode{  "\x03\xff\xff\xbe",  26 },
	byte(5): HuffmanCode{  "\x03\xff\xff\xbf",  26 },
	byte(6): HuffmanCode{  "\x03\xff\xff\xc0",  26 },
	byte(7): HuffmanCode{  "\x03\xff\xff\xc1",  26 },
	byte(8): HuffmanCode{  "\x03\xff\xff\xc2",  26 },
	byte(9): HuffmanCode{  "\x03\xff\xff\xc3",  26 },
	byte(10): HuffmanCode{  "\x03\xff\xff\xc4",  26 },
	byte(11): HuffmanCode{  "\x03\xff\xff\xc5",  26 },
	byte(12): HuffmanCode{  "\x03\xff\xff\xc6",  26 },
	byte(13): HuffmanCode{  "\x03\xff\xff\xc7",  26 },
	byte(14): HuffmanCode{  "\x03\xff\xff\xc8",  26 },
	byte(15): HuffmanCode{  "\x03\xff\xff\xc9",  26 },
	byte(16): HuffmanCode{  "\x03\xff\xff\xca",  26 },
	byte(17): HuffmanCode{  "\x03\xff\xff\xcb",  26 },
	byte(18): HuffmanCode{  "\x03\xff\xff\xcc",  26 },
	byte(19): HuffmanCode{  "\x03\xff\xff\xcd",  26 },
	byte(20): HuffmanCode{  "\x03\xff\xff\xce",  26 },
	byte(21): HuffmanCode{  "\x03\xff\xff\xcf",  26 },
	byte(22): HuffmanCode{  "\x03\xff\xff\xd0",  26 },
	byte(23): HuffmanCode{  "\x03\xff\xff\xd1",  26 },
	byte(24): HuffmanCode{  "\x03\xff\xff\xd2",  26 },
	byte(25): HuffmanCode{  "\x03\xff\xff\xd3",  26 },
	byte(26): HuffmanCode{  "\x03\xff\xff\xd4",  26 },
	byte(27): HuffmanCode{  "\x03\xff\xff\xd5",  26 },
	byte(28): HuffmanCode{  "\x03\xff\xff\xd6",  26 },
	byte(29): HuffmanCode{  "\x03\xff\xff\xd7",  26 },
	byte(30): HuffmanCode{  "\x03\xff\xff\xd8",  26 },
	byte(31): HuffmanCode{  "\x03\xff\xff\xd9",  26 },
	byte(32): HuffmanCode{  "\x06",   5 },
	byte(33): HuffmanCode{  "\x1f\xfc",  13 },
	byte(34): HuffmanCode{  "\x01\xf0",   9 },
	byte(35): HuffmanCode{  "\x3f\xfc",  14 },
	byte(36): HuffmanCode{  "\x7f\xfc",  15 },
	byte(37): HuffmanCode{  "\x1e",   6 },
	byte(38): HuffmanCode{  "\x64",   7 },
	byte(39): HuffmanCode{  "\x1f\xfd",  13 },
	byte(40): HuffmanCode{  "\x03\xfa",  10 },
	byte(41): HuffmanCode{  "\x01\xf1",   9 },
	byte(42): HuffmanCode{  "\x03\xfb",  10 },
	byte(43): HuffmanCode{  "\x03\xfc",  10 },
	byte(44): HuffmanCode{  "\x65",   7 },
	byte(45): HuffmanCode{  "\x66",   7 },
	byte(46): HuffmanCode{  "\x1f",   6 },
	byte(47): HuffmanCode{  "\x07",   5 },
	byte(48): HuffmanCode{  "\x00",   4 },
	byte(49): HuffmanCode{  "\x01",   4 },
	byte(50): HuffmanCode{  "\x02",   4 },
	byte(51): HuffmanCode{  "\x08",   5 },
	byte(52): HuffmanCode{  "\x20",   6 },
	byte(53): HuffmanCode{  "\x21",   6 },
	byte(54): HuffmanCode{  "\x22",   6 },
	byte(55): HuffmanCode{  "\x23",   6 },
	byte(56): HuffmanCode{  "\x24",   6 },
	byte(57): HuffmanCode{  "\x25",   6 },
	byte(58): HuffmanCode{  "\x26",   6 },
	byte(59): HuffmanCode{  "\xec",   8 },
	byte(60): HuffmanCode{  "\x01\xff\xfc",  17 },
	byte(61): HuffmanCode{  "\x27",   6 },
	byte(62): HuffmanCode{  "\x7f\xfd",  15 },
	byte(63): HuffmanCode{  "\x03\xfd",  10 },
	byte(64): HuffmanCode{  "\x7f\xfe",  15 },
	byte(65): HuffmanCode{  "\x67",   7 },
	byte(66): HuffmanCode{  "\xed",   8 },
	byte(67): HuffmanCode{  "\xee",   8 },
	byte(68): HuffmanCode{  "\x68",   7 },
	byte(69): HuffmanCode{  "\xef",   8 },
	byte(70): HuffmanCode{  "\x69",   7 },
	byte(71): HuffmanCode{  "\x6a",   7 },
	byte(72): HuffmanCode{  "\x01\xf2",   9 },
	byte(73): HuffmanCode{  "\xf0",   8 },
	byte(74): HuffmanCode{  "\x01\xf3",   9 },
	byte(75): HuffmanCode{  "\x01\xf4",   9 },
	byte(76): HuffmanCode{  "\x01\xf5",   9 },
	byte(77): HuffmanCode{  "\x6b",   7 },
	byte(78): HuffmanCode{  "\x6c",   7 },
	byte(79): HuffmanCode{  "\xf1",   8 },
	byte(80): HuffmanCode{  "\xf2",   8 },
	byte(81): HuffmanCode{  "\x01\xf6",   9 },
	byte(82): HuffmanCode{  "\x01\xf7",   9 },
	byte(83): HuffmanCode{  "\x6d",   7 },
	byte(84): HuffmanCode{  "\x28",   6 },
	byte(85): HuffmanCode{  "\xf3",   8 },
	byte(86): HuffmanCode{  "\x01\xf8",   9 },
	byte(87): HuffmanCode{  "\x01\xf9",   9 },
	byte(88): HuffmanCode{  "\xf4",   8 },
	byte(89): HuffmanCode{  "\x01\xfa",   9 },
	byte(90): HuffmanCode{  "\x01\xfb",   9 },
	byte(91): HuffmanCode{  "\x07\xfc",  11 },
	byte(92): HuffmanCode{  "\x03\xff\xff\xda",  26 },
	byte(93): HuffmanCode{  "\x07\xfd",  11 },
	byte(94): HuffmanCode{  "\x3f\xfd",  14 },
	byte(95): HuffmanCode{  "6e",   7 },
	byte(96): HuffmanCode{  "\x03\xff\xfe",  18 },
	byte(97): HuffmanCode{  "\x09",   5 },
	byte(98): HuffmanCode{  "\x6f",   7 },
	byte(99): HuffmanCode{  "\x0a",   5 },
	byte(100): HuffmanCode{  "\x29",   6 },
	byte(101): HuffmanCode{  "\x0b",   5 },
	byte(102): HuffmanCode{  "\x70",   7 },
	byte(103): HuffmanCode{  "\x2a",   6 },
	byte(104): HuffmanCode{  "\x2b",   6 },
	byte(105): HuffmanCode{  "\x0c",   5 },
	byte(106): HuffmanCode{  "\xf5",   8 },
	byte(107): HuffmanCode{  "\xf6",   8 },
	byte(108): HuffmanCode{  "\x2c",   6 },
	byte(109): HuffmanCode{  "\x2d",   6 },
	byte(110): HuffmanCode{  "\x2e",   6 },
	byte(111): HuffmanCode{  "\x0d",   5 },
	byte(112): HuffmanCode{ "\x2f",   6 },
	byte(113): HuffmanCode{  "\x01\xfc",   9 },
	byte(114): HuffmanCode{ "\x30",   6 },
	byte(115): HuffmanCode{  "\x31",   6 },
	byte(116): HuffmanCode{  "\x0e",   5 },
	byte(117): HuffmanCode{  "\x71",   7 },
	byte(118): HuffmanCode{  "\x72",   7 },
	byte(119): HuffmanCode{  "\x73",   7 },
	byte(120): HuffmanCode{  "\x74",   7 },
	byte(121): HuffmanCode{  "\x75",   7 },
	byte(122): HuffmanCode{  "\xf7",   8 },
	byte(123): HuffmanCode{  "\x01\xff\xfd",  17 },
	byte(124): HuffmanCode{  "\x0f\xfc",  12 },
	byte(125): HuffmanCode{  "\x01\xff\xfe",  17 },
	byte(126): HuffmanCode{  "\x0f\xfd",  12 },
	byte(127): HuffmanCode{  "\x03\xff\xff\xdb",  26 },
	byte(128): HuffmanCode{  "\x03\xff\xff\xdc",  26 },
	byte(129): HuffmanCode{  "\x03\xff\xff\xdd",  26 },
	byte(130): HuffmanCode{  "\x03\xff\xff\xde",  26 },
	byte(131): HuffmanCode{  "\x03\xff\xff\xdf",  26 },
	byte(132): HuffmanCode{  "\x03\xff\xff\xe0",  26 },
	byte(133): HuffmanCode{  "\x03\xff\xff\xe1",  26 },
	byte(134): HuffmanCode{  "\x03\xff\xff\xe2",  26 },
	byte(135): HuffmanCode{  "\x03\xff\xff\xe3",  26 },
	byte(136): HuffmanCode{  "\x03\xff\xff\xe4",  26 },
	byte(137): HuffmanCode{  "\x03\xff\xff\xe5",  26 },
	byte(138): HuffmanCode{  "\x03\xff\xff\xe6",  26 },
	byte(139): HuffmanCode{  "\x03\xff\xff\xe7",  26 },
	byte(140): HuffmanCode{  "\x03\xff\xff\xe8",  26 },
	byte(141): HuffmanCode{  "\x03\xff\xff\xe9",  26 },
	byte(142): HuffmanCode{  "\x03\xff\xff\xea",  26 },
	byte(143): HuffmanCode{  "\x03\xff\xff\xeb",  26 },
	byte(144): HuffmanCode{  "\x03\xff\xff\xec",  26 },
	byte(145): HuffmanCode{  "\x03\xff\xff\xed",  26 },
	byte(146): HuffmanCode{  "\x03\xff\xff\xee",  26 },
	byte(147): HuffmanCode{  "\x03\xff\xff\xef",  26 },
	byte(148): HuffmanCode{  "\x03\xff\xff\xf0",  26 },
	byte(149): HuffmanCode{  "\x03\xff\xff\xf1",  26 },
	byte(150): HuffmanCode{  "\x03\xff\xff\xf2",  26 },
	byte(151): HuffmanCode{  "\x03\xff\xff\xf3",  26 },
	byte(152): HuffmanCode{  "\x03\xff\xff\xf4",  26 },
	byte(153): HuffmanCode{  "\x03\xff\xff\xf5",  26 },
	byte(154): HuffmanCode{  "\x03\xff\xff\xf6",  26 },
	byte(155): HuffmanCode{  "\x03\xff\xff\xf7",  26 },
	byte(156): HuffmanCode{  "\x03\xff\xff\xf8",  26 },
	byte(157): HuffmanCode{  "\x03\xff\xff\xf9",  26 },
	byte(158): HuffmanCode{  "\x03\xff\xff\xfa",  26 },
	byte(159): HuffmanCode{  "\x03\xff\xff\xfb",  26 },
	byte(160): HuffmanCode{  "\x03\xff\xff\xfc",  26 },
	byte(161): HuffmanCode{  "\x03\xff\xff\xfd",  26 },
	byte(162): HuffmanCode{  "\x03\xff\xff\xfe",  26 },
	byte(163): HuffmanCode{  "\x03\xff\xff\xff",  26 },
	byte(164): HuffmanCode{  "\x01\xff\xff\x80",  25 },
	byte(165): HuffmanCode{  "\x01\xff\xff\x81",  25 },
	byte(166): HuffmanCode{  "\x01\xff\xff\x82",  25 },
	byte(167): HuffmanCode{  "\x01\xff\xff\x83",  25 },
	byte(168): HuffmanCode{  "\x01\xff\xff\x84",  25 },
	byte(169): HuffmanCode{  "\x01\xff\xff\x85",  25 },
	byte(170): HuffmanCode{  "\x01\xff\xff\x86",  25 },
	byte(171): HuffmanCode{  "\x01\xff\xff\x87",  25 },
	byte(172): HuffmanCode{  "\x01\xff\xff\x88",  25 },
	byte(173): HuffmanCode{  "\x01\xff\xff\x89",  25 },
	byte(174): HuffmanCode{  "\x01\xff\xff\x8a",  25 },
	byte(175): HuffmanCode{  "\x01\xff\xff\x8b",  25 },
	byte(176): HuffmanCode{  "\x01\xff\xff\x8c",  25 },
	byte(177): HuffmanCode{  "\x01\xff\xff\x8d",  25 },
	byte(178): HuffmanCode{  "\x01\xff\xff\x8e",  25 },
	byte(179): HuffmanCode{  "\x01\xff\xff\x8f",  25 },
	byte(180): HuffmanCode{  "\x01\xff\xff\x90",  25 },
	byte(181): HuffmanCode{  "\x01\xff\xff\x91",  25 },
	byte(182): HuffmanCode{  "\x01\xff\xff\x92",  25 },
	byte(183): HuffmanCode{  "\x01\xff\xff\x93",  25 },
	byte(184): HuffmanCode{  "\x01\xff\xff\x94",  25 },
	byte(185): HuffmanCode{  "\x01\xff\xff\x95",  25 },
	byte(186): HuffmanCode{  "\x01\xff\xff\x96",  25 },
	byte(187): HuffmanCode{  "\x01\xff\xff\x97",  25 },
	byte(188): HuffmanCode{  "\x01\xff\xff\x98",  25 },
	byte(189): HuffmanCode{  "\x01\xff\xff\x99",  25 },
	byte(190): HuffmanCode{  "\x01\xff\xff\x9a",  25 },
	byte(191): HuffmanCode{  "\x01\xff\xff\x9b",  25 },
	byte(192): HuffmanCode{  "\x01\xff\xff\x9c",  25 },
	byte(193): HuffmanCode{  "\x01\xff\xff\x9d",  25 },
	byte(194): HuffmanCode{  "\x01\xff\xff\x9e",  25 },
	byte(195): HuffmanCode{  "\x01\xff\xff\x9f",  25 },
	byte(196): HuffmanCode{  "\x01\xff\xff\xa0",  25 },
	byte(197): HuffmanCode{  "\x01\xff\xff\xa1",  25 },
	byte(198): HuffmanCode{  "\x01\xff\xff\xa2",  25 },
	byte(199): HuffmanCode{  "\x01\xff\xff\xa3",  25 },
	byte(200): HuffmanCode{  "\x01\xff\xff\xa4",  25 },
	byte(201): HuffmanCode{  "\x01\xff\xff\xa5",  25 },
	byte(202): HuffmanCode{  "\x01\xff\xff\xa6",  25 },
	byte(203): HuffmanCode{  "\x01\xff\xff\xa7",  25 },
	byte(204): HuffmanCode{  "\x01\xff\xff\xa8",  25 },
	byte(205): HuffmanCode{  "\x01\xff\xff\xa9",  25 },
	byte(206): HuffmanCode{  "\x01\xff\xff\xaa",  25 },
	byte(207): HuffmanCode{  "\x01\xff\xff\xab",  25 },
	byte(208): HuffmanCode{  "\x01\xff\xff\xac",  25 },
	byte(209): HuffmanCode{  "\x01\xff\xff\xad",  25 },
	byte(210): HuffmanCode{  "\x01\xff\xff\xae",  25 },
	byte(211): HuffmanCode{  "\x01\xff\xff\xaf",  25 },
	byte(212): HuffmanCode{  "\x01\xff\xff\xb0",  25 },
	byte(213): HuffmanCode{  "\x01\xff\xff\xb1",  25 },
	byte(214): HuffmanCode{  "\x01\xff\xff\xb2",  25 },
	byte(215): HuffmanCode{  "\x01\xff\xff\xb3",  25 },
	byte(216): HuffmanCode{  "\x01\xff\xff\xb4",  25 },
	byte(217): HuffmanCode{  "\x01\xff\xff\xb5",  25 },
	byte(218): HuffmanCode{  "\x01\xff\xff\xb6",  25 },
	byte(219): HuffmanCode{  "\x01\xff\xff\xb7",  25 },
	byte(220): HuffmanCode{  "\x01\xff\xff\xb8",  25 },
	byte(221): HuffmanCode{  "\x01\xff\xff\xb9",  25 },
	byte(222): HuffmanCode{  "\x01\xff\xff\xba",  25 },
	byte(223): HuffmanCode{  "\x01\xff\xff\xbb",  25 },
	byte(224): HuffmanCode{  "\x01\xff\xff\xbc",  25 },
	byte(225): HuffmanCode{  "\x01\xff\xff\xbd",  25 },
	byte(226): HuffmanCode{  "\x01\xff\xff\xbe",  25 },
	byte(227): HuffmanCode{  "\x01\xff\xff\xbf",  25 },
	byte(228): HuffmanCode{  "\x01\xff\xff\xc0",  25 },
	byte(229): HuffmanCode{  "\x01\xff\xff\xc1",  25 },
	byte(230): HuffmanCode{  "\x01\xff\xff\xc2",  25 },
	byte(231): HuffmanCode{  "\x01\xff\xff\xc3",  25 },
	byte(232): HuffmanCode{  "\x01\xff\xff\xc4",  25 },
	byte(233): HuffmanCode{  "\x01\xff\xff\xc5",  25 },
	byte(234): HuffmanCode{  "\x01\xff\xff\xc6",  25 },
	byte(235): HuffmanCode{  "\x01\xff\xff\xc7",  25 },
	byte(236): HuffmanCode{  "\x01\xff\xff\xc8",  25 },
	byte(237): HuffmanCode{  "\x01\xff\xff\xc9",  25 },
	byte(238): HuffmanCode{  "\x01\xff\xff\xca",  25 },
	byte(239): HuffmanCode{  "\x01\xff\xff\xcb",  25 },
	byte(240): HuffmanCode{  "\x01\xff\xff\xcc",  25 },
	byte(241): HuffmanCode{  "\x01\xff\xff\xcd",  25 },
	byte(242): HuffmanCode{  "\x01\xff\xff\xce",  25 },
	byte(243): HuffmanCode{  "\x01\xff\xff\xcf",  25 },
	byte(244): HuffmanCode{  "\x01\xff\xff\xd0",  25 },
	byte(245): HuffmanCode{  "\x01\xff\xff\xd1",  25 },
	byte(246): HuffmanCode{  "\x01\xff\xff\xd2",  25 },
	byte(247): HuffmanCode{  "\x01\xff\xff\xd3",  25 },
	byte(248): HuffmanCode{  "\x01\xff\xff\xd4",  25 },
	byte(249): HuffmanCode{  "\x01\xff\xff\xd5",  25 },
	byte(250): HuffmanCode{  "\x01\xff\xff\xd6",  25 },
	byte(251): HuffmanCode{  "\x01\xff\xff\xd7",  25 },
	byte(252): HuffmanCode{  "\x01\xff\xff\xd8",  25 },
	byte(253): HuffmanCode{  "\x01\xff\xff\xd9",  25 },
	byte(254): HuffmanCode{  "\x01\xff\xff\xda",  25 },
	byte(255): HuffmanCode{  "\x01\xff\xff\xdb",  25 },
}

var HuffmanEOS = HuffmanCode{ "\x01\xff\xff\xdc", 25 }

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

func padToUint32(a HuffmanCode) uint32 {
	bits := []byte(a.bits)

	asInt32 := make([]byte, 4)
	asInt32[0] = 0
	asInt32[1] = 0
	asInt32[2] = 0
	asInt32[3] = 0

	if len(bits) == 0 {
		// nothing
	} else if len(bits) <= 1 {
		asInt32[3] = bits[0]
	} else if len(bits) <= 2 {
		asInt32[3] = bits[1]
		asInt32[2] = bits[0]
	} else if len(bits) <= 3 {
		asInt32[3] = bits[2]
		asInt32[2] = bits[1]
		asInt32[1] = bits[0]
	} else if len(bits) <= 4 {
		asInt32[3] = bits[3]
		asInt32[2] = bits[2]
		asInt32[1] = bits[1]
		asInt32[0] = bits[0]
	}

	return binary.BigEndian.Uint32(asInt32)
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
	paddedA := padToUint32(a)
	paddedA = paddedA << uint(32 - a.bitLength)
	// we now have 32 - a.bitLength bits left in our uint
	// two options:
	// 1) b fits in the rest of the uint32, in which case we return another code
	//    (ugh code must be aligned back to LSB)
	//  2) b overflows uint32, in which case we
	//   i) return the uint32
	//   ii) return a code aligned to LSB

	if b.bitLength < 32 - a.bitLength {
		// b fits
		paddedB := padToUint32(b)
		paddedB = paddedB << uint(32 - a.bitLength - b.bitLength)
		remaining := uint(32 - a.bitLength - b.bitLength)

		code := make([]byte, 4)
		binary.BigEndian.PutUint32(code, (paddedA | paddedB) >> remaining)

		return "", HuffmanCode{ string(code), a.bitLength + b.bitLength }
	} else {
		// overflow
		overflow := make([]byte, 4)
		paddedB := padToUint32(b)
		overflowBits := uint(b.bitLength - (32 - a.bitLength))

		binary.BigEndian.PutUint32(overflow, paddedA | (paddedB >> overflowBits))
		code := make([]byte, 4)
		binary.BigEndian.PutUint32(code, paddedB & ((1 << overflowBits) - 1))

		return string(overflow), HuffmanCode{ string(code), overflowBits }
	}
}
