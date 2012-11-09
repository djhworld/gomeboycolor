package mmu

import (
	"github.com/djhworld/gomeboycolor/types"
	"github.com/stretchrcom/testify/assert"
	"testing"
)

func TestWriteByteToExternalRAM(t *testing.T) {
	//boundary tests

	//low
	var address types.Word = 0xA000
	var value byte = 0x83
	var normalisedLoc int = 0

	t.Logf("Writing %X to %X", value, address)
	gbc := new(GbcMMU)
	gbc.inBootMode = false
	gbc.WriteByte(address, value)
	assert.Equal(t, gbc.externalRAM[normalisedLoc], value)

	//middle
	address = 0xAFFF
	value = 0x33
	normalisedLoc = 4095

	t.Logf("Writing %X to %X", value, address)
	gbc = new(GbcMMU)
	gbc.inBootMode = false
	gbc.WriteByte(address, value)
	assert.Equal(t, gbc.externalRAM[normalisedLoc], value)

	//high
	address = 0xBFFF
	value = 0xA2
	normalisedLoc = 8191

	t.Logf("Writing %X to %X", value, address)
	gbc = new(GbcMMU)
	gbc.inBootMode = false
	gbc.WriteByte(address, value)
	assert.Equal(t, gbc.externalRAM[normalisedLoc], value)
}

func TestWriteByteToWorkingRAM(t *testing.T) {
	//boundary tests

	//low
	var address types.Word = 0xC000
	var value byte = 0x83
	var normalisedLoc int = 0
	var normalisedShadowLoc int = 0

	t.Logf("Writing %X to %X", value, address)
	gbc := new(GbcMMU)
	gbc.inBootMode = false
	gbc.WriteByte(address, value)
	assert.Equal(t, gbc.workingRAM[normalisedLoc], value)
	//check shadow 
	assert.Equal(t, gbc.workingRAMShadow[normalisedShadowLoc], value)

	//middle
	address = 0xCFFF
	value = 0x31
	normalisedLoc = 4095
	normalisedShadowLoc = 3583

	t.Logf("Writing %X to %X", value, address)
	gbc = new(GbcMMU)
	gbc.inBootMode = false
	gbc.WriteByte(address, value)
	assert.Equal(t, gbc.workingRAM[normalisedLoc], value)
	//check shadow 
	assert.Equal(t, gbc.workingRAMShadow[normalisedShadowLoc], value)

	//high
	address = 0xDFFF
	value = 0x87
	normalisedLoc = 8191

	t.Logf("Writing %X to %X", value, address)
	gbc = new(GbcMMU)
	gbc.inBootMode = false
	gbc.WriteByte(address, value)
	assert.Equal(t, gbc.workingRAM[normalisedLoc], value)
	//no shadow available as working ram shadow shaves off that last 512 bytes
}

func TestWriteByteToZeroPageRAM(t *testing.T) {
	//boundary tests

	//low
	var address types.Word = 0xFF80
	var value byte = 0x83
	var normalisedLoc int = 0

	t.Logf("Writing %X to %X", value, address)
	gbc := new(GbcMMU)
	gbc.inBootMode = false
	gbc.WriteByte(address, value)
	assert.Equal(t, gbc.zeroPageRAM[normalisedLoc], value)

	//middle
	address = 0xFFBF
	value = 0x33
	normalisedLoc = 63

	t.Logf("Writing %X to %X", value, address)
	gbc = new(GbcMMU)
	gbc.inBootMode = false
	gbc.WriteByte(address, value)
	assert.Equal(t, gbc.zeroPageRAM[normalisedLoc], value)

	//high
	address = 0xFFFF
	value = 0xA2
	normalisedLoc = 127

	t.Logf("Writing %X to %X", value, address)
	gbc = new(GbcMMU)
	gbc.inBootMode = false
	gbc.WriteByte(address, value)
	assert.Equal(t, gbc.zeroPageRAM[normalisedLoc], value)
}

func TestWriteByteToBootRegion(t *testing.T) {
	gbc := new(GbcMMU)
	gbc.inBootMode = true

	//should panic as you can't write to ROM!
	assert.Panics(t, func() {
		gbc.WriteByte(0x0001, 0xFE)
	}, "Should have panicked!")
}

func TestWriteByteToROMRegion(t *testing.T) {
	gbc := new(GbcMMU)
	gbc.inBootMode = false

	//should panic as you can't write to ROM!
	assert.Panics(t, func() {
		gbc.WriteByte(0x3FFE, 0xFE)
	}, "Should have paniciked!")
}

func TestRegionBoundaries(t *testing.T) {
	gbc := new(GbcMMU)
	gbc.boot[0] = 1
	gbc.boot[255] = 1

	gbc.cartrom[0] = 1
	gbc.cartrom[32767] = 1

	gbc.externalRAM[0] = 1
	gbc.externalRAM[8191] = 1

	gbc.workingRAM[0] = 1
	gbc.workingRAM[8191] = 1

	gbc.workingRAMShadow[0] = 1
	gbc.workingRAMShadow[7679] = 1

	gbc.zeroPageRAM[0] = 1
	gbc.zeroPageRAM[127] = 1

}

//TODO: READ BYTE/WORD operations - will require some mechanism to load ROM into memory first?

func TestLoadBootROM(t *testing.T) {
	var startAddr types.Word = 0
	var ROM []byte = []byte{0x03, 0x77, 0x04, 0xFF, 0xA3, 0xA2, 0xB3}
	gbc := new(GbcMMU)
	gbc.LoadROM(startAddr, BOOT, ROM)
	//check whether start address -> end of ROM is equal to ROM
	assert.Equal(t, gbc.boot[startAddr:len(ROM)], ROM)

	//check that error is returned if ROM is loaded that will over extend BOOT region
	gbc = new(GbcMMU)
	startAddr = 253
	ok, err := gbc.LoadROM(startAddr, BOOT, ROM)
	assert.False(t, ok)
	assert.NotNil(t, err)
	assert.Equal(t, ROMWillOverextendAddressableRegion, err)

	//check that error is returned if ROM is loaded that will over extend BOOT region
	gbc = new(GbcMMU)
	startAddr = 0
	ok, err = gbc.LoadROM(startAddr, BOOT, make([]byte, 3000))
	assert.False(t, ok)
	assert.NotNil(t, err)
	assert.Equal(t, ROMIsBiggerThanRegion, err)
}

func TestLoadCartROM(t *testing.T) {
	var startAddr types.Word = 0
	var rom []byte = []byte{0x03, 0x77, 0x04, 0xFF, 0xA3, 0xA2, 0xB3, 0xFF, 0x2C}
	gbc := new(GbcMMU)
	gbc.LoadROM(startAddr, CARTROM, rom)
	//check whether start address -> end of ROM is equal to ROM
	assert.Equal(t, gbc.cartrom[startAddr:len(rom)], rom)

	//check that error is returned if ROM is loaded that will over extend BOOT region
	gbc = new(GbcMMU)
	startAddr = 32765
	ok, err := gbc.LoadROM(startAddr, CARTROM, rom)
	assert.False(t, ok)
	assert.NotNil(t, err)
	assert.Equal(t, ROMWillOverextendAddressableRegion, err)

	//check that error is returned if ROM is loaded that will over extend BOOT region
	gbc = new(GbcMMU)
	startAddr = 0
	ok, err = gbc.LoadROM(startAddr, CARTROM, make([]byte, 42765))
	assert.False(t, ok)
	assert.NotNil(t, err)
	assert.Equal(t, ROMIsBiggerThanRegion, err)
}

func TestImplementsInterface(t *testing.T) {
	gbc := new(GbcMMU)
	assert.Implements(t, (*MemoryMappedUnit)(nil), gbc)
}
