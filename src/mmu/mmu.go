package mmu

import (
	"cartridge"
	"errors"
	"log"
	"types"
	"utils"
)

const PREFIX = "MMU:"

var ROMIsBiggerThanRegion error = errors.New("ROM is bigger than addressable region")

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
	LoadBIOS(data []byte) (bool, error)
	LoadCartridge(cart *cartridge.Cartridge)
	Reset()
}

type GbcMMU struct {
	bios              [256]byte //0x0000 -> 0x00FF
	cartridge         *cartridge.Cartridge
	externalRAM       [8192]byte //0xA000 -> 0xBFFF
	workingRAM        [8192]byte //0xC000 -> 0xDFFF
	workingRAMShadow  [7680]byte //0xE000 -> 0xFDFF
	zeroPageRAM       [128]byte  //0xFF80 - 0xFFFF
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

//TODO: THIS NEEDS REVIEWING AS IT IS WRONG.
func (mmu *GbcMMU) WriteByte(addr types.Word, value byte) {
	switch {
	case addr >= 0x2000 && addr <= 0x7FFF:
		log.Println("Hmmm")
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
	case addr >= 0xE000 && addr <= 0xFF7F:
		//Graphics sprite information
		if addr >= 0xFE00 && addr <= 0xFE9F {
			mmu.peripherals[GPU].Write(addr, value)
		} else if addr >= 0xFEA0 && addr <= 0xFF7F {
			//DMG status register (i.e. in boot mode)
			if addr == 0xFF50 {
				mmu.dmgStatusRegister = value
			}

			//Graphics registers
			if a := addr & 0x00F0; a >= 0x40 && a <= 0x70 {
				mmu.peripherals[GPU].Write(addr, value)
			}
		}
	//Zero page RAM
	case addr >= 0xFF80 && addr <= 0xFFFF:
		mmu.zeroPageRAM[addr&(0xFFFF-0xFF80)] = value
	default:
		log.Printf(PREFIX+" Address %X is unwritable/unknown", addr)
		panic("Address is unwritable/unknown")
	}
}

func (mmu *GbcMMU) ReadByte(addr types.Word) byte {
	switch {
	//ROM Bank 0
	case addr >= 0x0000 && addr <= 0x3FFF:
		//boot area/ROM after boot
		if mmu.inBootMode && addr <= 0x00FF {
			return mmu.bios[addr]
		}

		return mmu.cartridge.ROM[addr]
	//ROM Bank 1 (switchable)
	case addr >= 0x4000 && addr <= 0x7FFF:
		switch mmu.cartridge.Type.ID {
		case cartridge.MBC0:
			return mmu.cartridge.ROM[addr]
		default:
			log.Fatalf("Read is not set up for address 0x%X for cartridge type %s", addr, mmu.cartridge.Type.Description)
		}
	//Graphics VRAM
	case addr >= 0x8000 && addr <= 0x9FFF:
		return mmu.peripherals[GPU].Read(addr)
	//RAM Bank (switchable)
	case addr >= 0xA000 && addr <= 0xBFFF:
		switch mmu.cartridge.Type.ID {
		case cartridge.MBC0:
			return mmu.externalRAM[addr&(0xBFFF-0xA000)]
		default:
			log.Fatalf("Read is not set up for address 0x%X for cartridge type %s", addr, mmu.cartridge.Type.Description)
		}
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
	case addr >= 0xFF00 && addr <= 0xFF4C:
		if addr == 0xFF50 {
			return mmu.dmgStatusRegister
		}
		//Graphics registers
		if a := addr & 0x00F0; a >= 0x40 && a <= 0x70 {
			return mmu.peripherals[GPU].Read(addr)
		}
	//Zero page RAM
	case addr >= 0xFF80 && addr <= 0xFFFF:
		if addr == 0xFFFF {
			log.Println(PREFIX, "WARNING - Attempting to read from interrupt register - this is unimplemented!!")
		}
		return mmu.zeroPageRAM[addr&(0xFFFF-0xFF80)]
	}

	return 0x00
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

func (mmu *GbcMMU) LoadBIOS(data []byte) (bool, error) {
	log.Println(PREFIX+" Loading ", len(data), "byte BIOS ROM into MMU")
	if len(data) > len(mmu.bios) {
		return false, ROMIsBiggerThanRegion
	}

	for i, b := range data {
		mmu.bios[i] = b
	}
	return true, nil
}

func (mmu *GbcMMU) LoadCartridge(cart *cartridge.Cartridge) {
	mmu.cartridge = cart
	log.Printf("%s Loaded cartridge: -\n%s\n", PREFIX, cart)
}
