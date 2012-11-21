package mmu

import (
	"errors"
	"github.com/djhworld/gomeboycolor/types"
	"github.com/djhworld/gomeboycolor/utils"
	"log"
)

/*
This represents the memory mapped unit for our emulator. 
The regions of memory are arrays of fixed size that are selected
based on the memory address given
*/

var ROMIsBiggerThanRegion error = errors.New("ROM is bigger than addressable region")
var ROMWillOverextendAddressableRegion = errors.New("ROM will overextend addressable region based on start address and ROM size")

const (
	BOOT    = 0x00
	CARTROM = 0x01
)

type MemoryMappedUnit interface {
	WriteByte(address types.Word, value byte)
	WriteWord(address types.Word, value types.Word)
	ReadByte(address types.Word) byte
	ReadWord(address types.Word) types.Word
	SetInBootMode(mode bool)
	LoadROM(startAddr types.Word, rt types.ROMType, data []byte) (bool, error)
	Reset()
}

type GbcMMU struct {
	boot             [256]byte   //0x0000 -> 0x00FF
	cartrom          [32768]byte // 0x0000 -> 0x7FFF
	externalRAM      [8192]byte  //0xA000 -> 0xBFFF
	workingRAM       [8192]byte  //0xC000 -> 0xDFFF
	workingRAMShadow [7680]byte  //0xE000 -> 0xFDFF
	mmIO             [128]byte   //0xFF00 -> 0xFF7F //TODO 
	zeroPageRAM      [128]byte   //0xFF80 - 0xFFFF
	inBootMode       bool
}

func (mmu *GbcMMU) Reset() {
	mmu.inBootMode = true
}

func (mmu *GbcMMU) ReadByte(addr types.Word) byte {
	switch {
	//boot area/ROM after boot
	case addr >= 0x0000 && addr <= 0x00FF:
		if mmu.inBootMode {
			return mmu.boot[addr]
		} else {
			return mmu.cartrom[addr]
		}
	//ROM Bank 0
	case addr >= 0x0100 && addr <= 0x3FFF:
		return mmu.cartrom[addr]
	//ROM Bank 1
	case addr >= 0x4000 && addr <= 0x7FFF:
		return mmu.cartrom[addr]
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
		return mmu.mmIO[addr&(0xFF7F-0xFF00)]
	//Zero page RAM
	case addr >= 0xFF80 && addr <= 0xFFFF:
		return mmu.zeroPageRAM[addr&(0xFFFF-0xFF80)]

	}

	return 0
}

func (mmu *GbcMMU) WriteByte(addr types.Word, value byte) {
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
		mmu.mmIO[addr&(0xFF7F-0xFF00)] = value
	//Zero page RAM
	case addr >= 0xFF80 && addr <= 0xFFFF:
		mmu.zeroPageRAM[addr&(0xFFFF-0xFF80)] = value
	default:
		log.Printf("Address %X is unwritable/unknown", addr)
		panic("Address is unwritable/unknown")
	}
}

func (mmu *GbcMMU) ReadWord(addr types.Word) types.Word {
	var b1 byte = mmu.ReadByte(addr)
	var b2 byte = mmu.ReadByte(addr + 1)
	return types.Word(utils.JoinBytes(b1, b2))
}

func (mmu *GbcMMU) WriteWord(addr types.Word, value types.Word) {
	b1, b2 := utils.SplitIntoBytes(uint16(value))
	mmu.WriteByte(addr, b1)
	mmu.WriteByte(addr+1, b2)
}

func (mmu *GbcMMU) SetInBootMode(mode bool) {
	mmu.inBootMode = mode
}

func (mmu *GbcMMU) LoadROM(startAddr types.Word, rt types.ROMType, data []byte) (bool, error) {
	doBoundaryChecks := func(mem int) error {
		if len(data) > mem {
			return ROMIsBiggerThanRegion
		}

		if startAddr+types.Word(len(data)) > types.Word(mem) {
			return ROMWillOverextendAddressableRegion
		}

		return nil
	}

	switch rt {
	case BOOT:
		if err := doBoundaryChecks(len(mmu.boot)); err != nil {
			return false, err
		}

		for i, b := range data {
			mmu.boot[startAddr+types.Word(i)] = b
		}

	case CARTROM:
		if err := doBoundaryChecks(len(mmu.cartrom)); err != nil {
			return false, err
		}

		for i, b := range data {
			mmu.cartrom[startAddr+types.Word(i)] = b
		}
	}

	return true, nil
}
