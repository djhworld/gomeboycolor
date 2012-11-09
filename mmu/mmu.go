package main

import (
	"github.com/djhworld/gomeboycolor/utils"
)

/*
This represents the memory mapped unit for our emulator. 
The regions of memory are arrays of fixed size that are selected
based on the memory address given
*/

var ROMIsBiggerThanRegion error
var ROMWillOverextendAddressableRegion error

type ROMType byte

const (
	BOOT = 0x00
	ROM  = 0x01
)

type MMU interface {
	WriteByte(address Word, value byte)
	WriteWord(address Word, value Word)
	ReadByte(address Word) byte
	ReadWord(address Word) Word
	SetInBootMode(mode bool)
	LoadROM(startAddr Word, rt ROMType, data []byte) (bool, error)
}

type GbcMMU struct {
	boot             [256]byte   //0x0000 -> 0x00FF
	ROM              [32768]byte // 0x0000 -> 0x7FFF
	externalRAM      [8192]byte  //0xA000 -> 0xBFFF
	workingRAM       [8192]byte  //0xC000 -> 0xDFFF
	workingRAMShadow [7680]byte  //0xE000 -> 0xFDFF
	zeroPageRAM      [128]byte   //0xFF80 - 0xFFFF
	inBootMode       bool
}

func (mmu *GbcMMU) ReadByte(addr Word) byte {
	switch {
	//boot area/ROM after boot
	case addr >= 0x0000 && addr <= 0x00FF:
		if mmu.inBootMode {
			return mmu.boot[addr]
		} else {
			return mmu.ROM[addr]
		}
	//ROM Bank 0
	case addr >= 0x1000 && addr <= 0x3FFF:
		return mmu.ROM[addr]
	//ROM Bank 1
	case addr >= 0x4000 && addr <= 0x7FFF:
		return mmu.ROM[addr]
	//Graphics VRAM
	case addr >= 0x8000 && addr <= 0x9FFF:
		//TODO - needs GPU setup
		return 0x00
	//Cartridge External RAM
	case addr >= 0xA000 && addr <= 0xBFFF:
		return mmu.externalRAM[addr&(0xBFFF-0xA000)]
	//GB Working RAM 
	case addr >= 0xC000 && addr <= 0xDFFF:
		return mmu.workingRAM[addr&(0xDFFF-0xC000)]
	//GB Working RAM shadow
	case addr >= 0xE000 && addr <= 0xFDFF:
		return mmu.workingRAM[addr&(0xFDFF-0xE000)]
	//Graphics sprite information
	case addr >= 0xFE00 && addr <= 0xFE9F:
		//TODO - needs GPU setup
		return 0
	//Mem. mapped IO
	case addr >= 0xFF00 && addr <= 0xFF7F:
		//TODO
		return 0
	//Zero page RAM
	case addr >= 0xFF80 && addr <= 0xFFFF:
		return mmu.zeroPageRAM[addr&(0xFFFF-0xFF80)]

	}

	return 0
}

func (mmu *GbcMMU) WriteByte(addr Word, value byte) {
	switch {
	//Graphics VRAM
	case addr >= 0x8000 && addr <= 0x9FFF:
		//TODO - needs gpu setup
		return
	//Cartridge External RAM
	case addr >= 0xA000 && addr <= 0xBFFF:
		mmu.externalRAM[addr&(0xBFFF-0xA000)] = value
	//GB Working RAM 
	case addr >= 0xC000 && addr <= 0xDFFF:
		mmu.workingRAM[addr&(0xDFFF-0xC000)] = value
		//copy value to shadow if within shadow range
		if addr >= 0xC000 && addr <= 0xDDFF {
			mmu.workingRAMShadow[addr&(0xDDFF-0xC000)] = value
		}
	//Graphics sprite information
	case addr >= 0xFE00 && addr <= 0xFE9F:
		//TODO - needs GPU setup
		return
	//Mem. mapped IO
	case addr >= 0xFF00 && addr <= 0xFF7F:
		//TODO
		return
	//Zero page RAM
	case addr >= 0xFF80 && addr <= 0xFFFF:
		mmu.zeroPageRAM[addr&(0xFFFF-0xFF80)] = value
	default:
		panic("Address is unwritable/uknown")
	}
}

func (mmu *GbcMMU) ReadWord(addr Word) Word {
	var b1 byte = mmu.ReadByte(addr)
	var b2 byte = mmu.ReadByte(addr + 1)
	return Word(utils.JoinBytes(b1, b2))
}

func (mmu *GbcMMU) WriteWord(addr Word, value Word) {
	b1, b2 := utils.SplitIntoBytes(uint16(value))
	mmu.WriteByte(addr, b1)
	mmu.WriteByte(addr+1, b2)
}

func (mmu *GbcMMU) SetInBootMode(mode bool) {
	mmu.inBootMode = mode
}

func (mmu *GbcMMU) LoadROM(startAddr Word, rt ROMType, data []byte) (bool, error) {
	switch rt {
	case BOOT:
		if len(data) > len(mmu.boot) {
			return false, ROMIsBiggerThanRegion
		}

		if Word(len(data)) > Word(len(data))-startAddr {
			return false, ROMWillOverextendAddressableRegion
		}

		for i, b := range data {
			mmu.boot[startAddr+Word(i)] = b
		}
		//TODO: ROM
	}

	return true, nil
}
