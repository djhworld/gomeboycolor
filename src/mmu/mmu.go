package mmu

import (
	"cartridge"
	"errors"
	"log"
	"types"
	"utils"
	"sort"
	"fmt"
)

const PREFIX = "MMU"

var ROMIsBiggerThanRegion error = errors.New("ROM is bigger than addressable region")

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
	internalRAM       [8192]byte //0xC000 -> 0xDFFF
	internalRAMShadow [7680]byte //0xE000 -> 0xFDFF
	zeroPageRAM       [128]byte  //0xFF80 - 0xFFFE
	inBootMode        bool
	dmgStatusRegister byte
	peripheralIOMap   map[types.Word]Peripheral
}

func NewGbcMMU() *GbcMMU {
	var mmu *GbcMMU = new(GbcMMU)
	mmu.peripheralIOMap = make(map[types.Word]Peripheral)
	mmu.Reset()
	return mmu
}

func (mmu *GbcMMU) Reset() {
	mmu.inBootMode = true
}

//TODO: NEED TO HANDLE WRITES TO ROM SPACE SO CAN CALCULATE ROM BANKS ETC
func (mmu *GbcMMU) WriteByte(addr types.Word, value byte) {
	//Check peripherals first
	//Graphics sprite information 0xFE00 - 0xFE9F
	//Graphics VRAM: 0x8000 - 0x9FFF
	//Graphics Registers: 0xFF40-0xFF49, 0xFF51-0xFF70
	if p, ok := mmu.peripheralIOMap[addr]; ok {
		p.Write(addr, value)
		return
	}

	switch {
	case addr >= 0x0000 && addr <= 0x9FFF:
		log.Printf("Attempted to write value %X to %s, hmmmm", value, addr)
		return
	//Cartridge External RAM
	case addr >= 0xA000 && addr <= 0xBFFF:
		switch mmu.cartridge.Type.ID {
		case cartridge.MBC0:
			mmu.cartridge.WriteByteToSwitchableRAM(0, addr, value)
		default:
			log.Fatalf("Write is not set up for address %s for cartridge type %s", addr, mmu.cartridge.Type.Description)
		}
	//GB Internal RAM
	case addr >= 0xC000 && addr <= 0xDFFF:
		mmu.internalRAM[addr&(0xDFFF-0xC000)] = value
		//copy value to shadow if within shadow range
		if addr >= 0xC000 && addr <= 0xDDFF {
			mmu.internalRAMShadow[addr&(0xDDFF-0xC000)] = value
		}
	//INTERRUPT FLAG
	case addr == 0xFF0F:
		log.Printf("Attempting to write %X to interrupt flag, this is currently unimplemented", value)
	//DMG flag
	case addr == 0xFF50:
		mmu.dmgStatusRegister = value
	//Zero page RAM
	case addr >= 0xFF80 && addr <= 0xFFFF:
		mmu.zeroPageRAM[addr&(0xFFFF-0xFF80)] = value
	default:
		log.Printf("%s: WARNING - Attempting to write to address %s, this is invalid/unimplemented", PREFIX, addr)
	}
}

func (mmu *GbcMMU) ReadByte(addr types.Word) byte {
	//Check peripherals first
	//Graphics sprite information 0xFE00 - 0xFE9F
	//Graphics VRAM: 0x8000 - 0x9FFF
	//Graphics Registers: 0xFF40-0xFF49, 0xFF51-0xFF70
	if p, ok := mmu.peripheralIOMap[addr]; ok {
		return p.Read(addr)
	}

	switch {
	//ROM Bank 0
	case addr >= 0x0000 && addr <= 0x3FFF:
		if mmu.inBootMode && addr < 0x0100 {
			//in bios mode, read from bios
			return mmu.bios[addr]
		}
		return mmu.cartridge.ReadByteFromROM(addr)
	//ROM Bank 1 (switchable)
	case addr >= 0x4000 && addr <= 0x7FFF:
		switch mmu.cartridge.Type.ID {
		case cartridge.MBC0:
			return mmu.cartridge.ReadByteFromSwitchableROM(0, addr)
		default:
			log.Fatalf("Read is not set up for address %s for cartridge type %s", addr, mmu.cartridge.Type.Description)
		}
	//RAM Bank (switchable)
	case addr >= 0xA000 && addr <= 0xBFFF:
		switch mmu.cartridge.Type.ID {
		case cartridge.MBC0:
			return mmu.cartridge.ReadByteFromSwitchableRAM(0, addr)
		default:
			log.Fatalf("Read is not set up for address %s for cartridge type %s", addr, mmu.cartridge.Type.Description)
		}
	//GB Internal RAM
	case addr >= 0xC000 && addr <= 0xDFFF:
		return mmu.internalRAM[addr&(0xDFFF-0xC000)]
	//GB Internal RAM shadow
	case addr >= 0xE000 && addr <= 0xFDFF:
		return mmu.internalRAMShadow[addr&(0xFDFF-0xE000)]
	//INTERRUPT FLAG
	case addr == 0xFF0F:
		log.Printf("Attempting to read from interrupt flag, this is currently unimplemented")
	//DMG FLAG
	case addr == 0xFF50:
		return mmu.dmgStatusRegister
	//Zero page RAM
	case addr >= 0xFF80 && addr <= 0xFFFF:
		return mmu.zeroPageRAM[addr&(0xFFFF-0xFF80)]
	default:
		log.Printf("%s: WARNING - Attempting to read from address %s, this is invalid/unimplemented", PREFIX, addr)
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

func (mmu *GbcMMU) ConnectPeripheral(p Peripheral, startAddr, endAddr types.Word) {
	log.Printf("%s: Connecting MMU to %s on address range %s to %s", PREFIX, p.Name(), startAddr, endAddr)
	for addr := startAddr; addr <= endAddr; addr++ {
		mmu.peripheralIOMap[addr] = p
	}
}

func (mmu *GbcMMU) PrintPeripheralMap() {
	var addrs types.Words
	for k, _ := range mmu.peripheralIOMap {
		addrs = append(addrs, k)
	}

	sort.Sort(addrs)

	for i, addr := range addrs {
		peripheral := mmu.peripheralIOMap[addr]

		fmt.Printf("[%s] -> %s   ", addr, peripheral.Name())
		if i % 8 == 0 {
			fmt.Println()
		}
	}

	fmt.Println()
}

func (mmu *GbcMMU) LoadBIOS(data []byte) (bool, error) {
	log.Println(PREFIX+": Loading", len(data), "byte BIOS ROM into MMU")
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
	log.Printf("%s: Loaded cartridge into MMU: -\n%s\n", PREFIX, cart)

}
