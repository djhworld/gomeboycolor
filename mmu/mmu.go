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

const PREFIX = "MMU:"

const (
	BOOT    = 0x00
	CARTROM = 0x01
)

//Peripherals
const (
	_ = iota
	GPU
	APU
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
	boot              [256]byte   //0x0000 -> 0x00FF
	cartrom           [32768]byte // 0x0000 -> 0x7FFF
	externalRAM       [8192]byte  //0xA000 -> 0xBFFF
	workingRAM        [8192]byte  //0xC000 -> 0xDFFF
	workingRAMShadow  [7680]byte  //0xE000 -> 0xFDFF
	zeroPageRAM       [128]byte   //0xFF80 - 0xFFFF
	inBootMode        bool
	peripherals       map[byte]Peripheral
	dmgStatusRegister byte
}

func NewGbcMMU() *GbcMMU {
	var mmu *GbcMMU = new(GbcMMU)
	mmu.peripherals = make(map[byte]Peripheral)
	mmu.Reset()
	return mmu
}

func (mmu *GbcMMU) Reset() {
	mmu.inBootMode = true
}

func (mmu *GbcMMU) WriteByte(addr types.Word, value byte) {
	switch {
	//Graphics VRAM
	case addr >= 0x8000 && addr <= 0x9FFF:
		mmu.peripherals[GPU].Write(addr, value)
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
		mmu.peripherals[GPU].Write(addr, value)
	//Mem. mapped IO
	case addr >= 0xFF00 && addr <= 0xFF7F:
		//DMG status register (i.e. in boot mode)
		if addr == 0xFF50 {
			mmu.dmgStatusRegister = value
		}

		//Graphics registers
		if a := addr & 0x00F0; a >= 0x40 && a <= 0x70 {
			mmu.peripherals[GPU].Write(addr, value)
		}
	//Zero page RAM
	case addr >= 0xFF80 && addr <= 0xFFFF:
		mmu.zeroPageRAM[addr&(0xFFFF-0xFF80)] = value
	default:
		log.Printf("Address %X is unwritable/unknown", addr)
		panic("Address is unwritable/unknown")
	}
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
		return mmu.peripherals[GPU].Read(addr)
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
		return mmu.peripherals[GPU].Read(addr)
	//Mem. mapped IO
	case addr >= 0xFF00 && addr <= 0xFF7F:
		if addr == 0xFF50 {
			return mmu.dmgStatusRegister
		}
		//Graphics registers
		if a := addr & 0x00F0; a >= 0x40 && a <= 0x70 {
			return mmu.peripherals[GPU].Read(addr)
		}
	//Zero page RAM
	case addr >= 0xFF80 && addr <= 0xFFFF:
		return mmu.zeroPageRAM[addr&(0xFFFF-0xFF80)]

	}

	return 0
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

func (mmu *GbcMMU) LinkGPU(p Peripheral) {
	log.Println(PREFIX, "Linked GPU to MMU")
	mmu.peripherals[GPU] = p
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

		log.Printf(PREFIX+" Writing data to boot rom sector (start address: 0x%X)", startAddr)
		for i, b := range data {
			mmu.boot[startAddr+types.Word(i)] = b
		}

	case CARTROM:
		if err := doBoundaryChecks(len(mmu.cartrom)); err != nil {
			return false, err
		}

		log.Printf(PREFIX+" Writing data to cartridge rom sector (start address: 0x%X)", startAddr)
		for i, b := range data {
			mmu.cartrom[startAddr+types.Word(i)] = b
		}
	}

	return true, nil
}
