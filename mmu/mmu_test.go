package main

import (
	"github.com/stretchrcom/testify/assert"
	"testing"
)

func TestWriteByteToExternalRAM(t *testing.T) {
	//boundary tests

	//low
	var address Word = 0xA000
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
	var address Word = 0xC000
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
	var address Word = 0xFF80
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

	gbc.ROM[0] = 1
	gbc.ROM[32767] = 1

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

func TestImplementsInterface(t *testing.T) {
	gbc := new(GbcMMU)
	assert.Implements(t, (*MMU)(nil), gbc)
}
