package utils

import "testing"
import "github.com/stretchrcom/testify/assert"

func TestSwapNibbles(t *testing.T) {
	assert.Equal(t, SwapNibbles(0x30), byte(0x03))
	assert.Equal(t, SwapNibbles(0xF2), byte(0x2F))
	assert.Equal(t, SwapNibbles(0x11), byte(0x11))
	assert.Equal(t, SwapNibbles(0xC4), byte(0x4C))
	assert.Equal(t, SwapNibbles(0xAE), byte(0xEA))
}
func TestJoinBytes(t *testing.T) {
	assert.Equal(t, JoinBytes(0x03, 0xFF), uint16(0x03FF))
	assert.Equal(t, JoinBytes(0x00, 0x00), uint16(0x0000))
	assert.Equal(t, JoinBytes(0x03, 0x03), uint16(0x0303))
	assert.Equal(t, JoinBytes(0xFF, 0xFF), uint16(0xFFFF))
	assert.Equal(t, JoinBytes(0xFE, 0xFF), uint16(0xFEFF))
	assert.Equal(t, JoinBytes(0xFF, 0x00), uint16(0xFF00))
}
