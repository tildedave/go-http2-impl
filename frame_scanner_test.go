package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFrameScannerReturnsAFrame(t *testing.T) {
	fakeConn := newFakeConn()
	b := PING{OpaqueData: 3957102}.Marshal()
	fakeConn.readData = [][]byte{b}

	s := NewFrameScanner(fakeConn)

	assert.True(t, s.Scan())
	assert.Equal(t, s.Bytes(), b)
}

func TestFrameScanner_IncompleteFrame(t *testing.T) {
	fakeConn := newFakeConn()
	b := PING{OpaqueData: 3957102}.Marshal()
	fakeConn.readData = [][]byte{b[0 : len(b)-1]}

	s := NewFrameScanner(fakeConn)

	assert.False(t, s.Scan())
}

func TestFrameScanner_IncompleteFrameThatIsLaterCompleted(t *testing.T) {
	fakeConn := newFakeConn()
	b := PING{OpaqueData: 3957102}.Marshal()
	fakeConn.readData = [][]byte{b[0 : len(b)-1], b[len(b)-1:]}

	s := NewFrameScanner(fakeConn)

	assert.True(t, s.Scan())
	assert.Equal(t, s.Bytes(), b)
}

func TestFrameScanner_TwoFrames(t *testing.T) {
	fakeConn := newFakeConn()
	b1 := PING{OpaqueData: 3957102}.Marshal()
	b2 := PING{OpaqueData: 12311}.Marshal()

	fakeConn.readData = [][]byte{b1, b2}

	s := NewFrameScanner(fakeConn)

	assert.True(t, s.Scan())
	assert.Equal(t, s.Bytes(), b1)
	assert.True(t, s.Scan())
	assert.Equal(t, s.Bytes(), b2)
	assert.False(t, s.Scan())
}

func TestFrameScanner_TwoFramesCombined(t *testing.T) {
	fakeConn := newFakeConn()
	b1 := PING{OpaqueData: 3957102}.Marshal()
	b2 := PING{OpaqueData: 12311}.Marshal()

	fakeConn.readData = [][]byte{append(b1, b2...)}

	s := NewFrameScanner(fakeConn)

	assert.True(t, s.Scan())
	assert.Equal(t, s.Bytes(), b1)
	assert.True(t, s.Scan())
	assert.Equal(t, s.Bytes(), b2)
	assert.False(t, s.Scan())
}

func TestFrameScanner_TwoFrames_Uneven(t *testing.T) {
	fakeConn := newFakeConn()
	b1 := PING{OpaqueData: 3957102}.Marshal()
	b2 := PING{OpaqueData: 12311}.Marshal()

	fakeConn.readData = [][]byte{b1[0:13], append(b1[13:], b2...)}

	s := NewFrameScanner(fakeConn)

	assert.True(t, s.Scan())
	assert.Equal(t, s.Bytes(), b1)
	assert.True(t, s.Scan())
	assert.Equal(t, s.Bytes(), b2)
	assert.False(t, s.Scan())
}
