package mmu

import (
	"github.com/stretchrcom/testify/assert"
	"log"
	"testing"
	"types"
)

func TestLoadBiosROM(t *testing.T) {
	var ROM []byte = []byte{0x03, 0x77, 0x04, 0xFF, 0xA3, 0xA2, 0xB3}
	mmu := NewGbcMMU()
	mmu.LoadROM(ROM, BIOS)
	assert.Equal(t, mmu.bios, ROM)

	//Check limit on size for BIOS
	mmu := NewGbcMMU()
	ok, err := gbc.LoadROM(make([]byte, 3000), BIOS)
	assert.False(t, ok)
	assert.NotNil(t, err)
	assert.Equal(t, ROMIsBiggerThanRegion, err)
}

func TestLoadCartROM(t *testing.T) {
	var startAddr types.Word = 0
	var ROM []byte = []byte{0x03, 0x77, 0x04, 0xFF, 0xA3, 0xA2, 0xB3}
	mmu := NewGbcMMU()
	mmu.LoadROM(ROM, CARTROM)

	assert.Equal(t, mmu.cartrom, ROM)
}

func TestImplementsInterface(t *testing.T) {
	gbc := NewGbcMMU()
	assert.Implements(t, (*MemoryMappedUnit)(nil), gbc)
}
