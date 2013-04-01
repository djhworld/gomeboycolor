package cartridge

import (
	"log"
	"types"
)

const (
	SIXTEENMB_ROM_8KBRAM = iota
	FOURMB_ROM_32KBRAM
)

//Represents MBC1
type MBC1 struct {
	Name            string
	romBank0        []byte
	romBanks        [][]byte
	ramBanks        [][]byte
	selectedROMBank int
	selectedRAMBank int
	hasRAM          bool
	ramEnabled      bool
	NoOfROMBanks    int
	NoOfRAMBanks    int
	MaxMemMode      int
}

func NewMBC1(rom []byte, romSize int, ramSize int) *MBC1 {
	var m *MBC1 = new(MBC1)

	m.Name = "CARTRIDGE-MBC1"
	m.MaxMemMode = SIXTEENMB_ROM_8KBRAM

	if ramSize > 0 {
		m.hasRAM = true
		m.selectedRAMBank = 0
		m.ramBanks = make([][]byte, 4)
		for i := 0; i < 4; i++ {
			m.ramBanks[i] = make([]byte, 0x2000)
		}
	}
	m.ramEnabled = true

	m.selectedROMBank = 0
	m.romBank0 = rom[0x0000:0x4000]

	m.NoOfROMBanks = romSize / 0x4000

	m.romBanks = populateROMBanks(rom, m.NoOfROMBanks)

	log.Println(m.Name+": Initialised memory bank controller with", len(m.romBanks), "ROM banks,", len(m.ramBanks), "RAM banks")
	return m
}

func (m *MBC1) Write(addr types.Word, value byte) {
	switch {
	case addr >= 0x0000 && addr <= 0x1FFF:
		//when in 4/32 mode...
		if m.MaxMemMode == FOURMB_ROM_32KBRAM && m.hasRAM {
			if r := value & 0x0F; r == 0x0A {
				log.Println(m.Name + ": Enabling RAM")
				m.ramEnabled = true
			} else {
				log.Println(m.Name + ": Disabling RAM")
				m.ramEnabled = false
			}
		}
	case addr >= 0x2000 && addr <= 0x3FFF:
		m.switchROMBank(int(value & 0x1F))
	case addr >= 0x4000 && addr <= 0x5FFF:
		m.switchRAMBank(int(value & 0x03))
	case addr >= 0x6000 && addr <= 0x7FFF:
		if mode := value & 0x01; mode == 0x00 {
			log.Println(m.Name + ": Switching MBC1 mode to 16/8")
			m.MaxMemMode = SIXTEENMB_ROM_8KBRAM
		} else {
			log.Println(m.Name + ": Switching MBC1 mode to 4/32")
			m.MaxMemMode = FOURMB_ROM_32KBRAM
		}
	case addr >= 0xA000 && addr <= 0xBFFF:
		if m.hasRAM && m.ramEnabled {
			m.ramBanks[m.selectedRAMBank][addr-0xA000] = value
		}
	}
}

func (m *MBC1) Read(addr types.Word) byte {
	//ROM Bank 0
	if addr < 0x4000 {
		return m.romBank0[addr]
	}

	//Switchable ROM BANK
	if addr >= 0x4000 && addr < 0x8000 {
		return m.romBanks[m.selectedROMBank][addr-0x4000]
	}

	//Upper bounds of memory map.
	if addr >= 0xA000 && addr <= 0xC000 {
		//		if m.MaxMemMode == FOURMB_ROM_32KBRAM {
		return m.ramBanks[m.selectedRAMBank][addr-0xA000]
		//		} else {
		//TODO: sort this out for 16/8
		//log.Fatalf("Unimplemented for 16/8 mode (Addr = %s)", addr)
		//			return 0x00
		//		}
	}

	return 0x00
}

func (m *MBC1) switchROMBank(bank int) {
	m.selectedROMBank = bank
}

func (m *MBC1) switchRAMBank(bank int) {
	m.selectedRAMBank = bank
}

func populateROMBanks(rom []byte, noOfBanks int) [][]byte {
	romBanks := make([][]byte, noOfBanks)

	//ROM Bank 0 and 1 are the same
	romBanks[0] = rom[0x4000:0x8000]
	var chunk int = 0x4000
	for i := 1; i < noOfBanks; i++ {
		romBanks[i] = rom[chunk : chunk+0x4000]
		chunk += 0x4000
	}

	return romBanks
}
