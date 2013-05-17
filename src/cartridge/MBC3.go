package cartridge

import (
	"constants"
	"fmt"
	"log"
	"types"
)

//Represents MBC3
type MBC3 struct {
	Name            string
	romBank0        []byte
	romBanks        [][]byte
	ramBanks        [][]byte
	selectedROMBank int
	selectedRAMBank int
	hasRAM          bool
	ramEnabled      bool
	MaxMemMode      int
	ROMSize         int
	RAMSize         int
}

func NewMBC3(rom []byte, romSize int, ramSize int) *MBC3 {
	var m *MBC3 = new(MBC3)

	m.Name = "CARTRIDGE-MBC3"
	m.MaxMemMode = constants.SIXTEENMB_ROM_8KBRAM
	m.ROMSize = romSize
	m.RAMSize = ramSize

	m.ramEnabled = true
	if ramSize > 0 {
		m.hasRAM = true
		m.selectedRAMBank = 0
		m.ramBanks = populateRAMBanks(m.RAMSize / 0x2000)
	}

	m.selectedROMBank = 0
	m.romBank0 = rom[0x0000:0x4000]
	m.romBanks = populateROMBanks(rom, m.ROMSize/0x4000)

	log.Println(m)
	return m
}
func (m *MBC3) String() string {
	return fmt.Sprint(m.Name+": ROM Banks: ", len(m.romBanks), ", RAM Banks: ", len(m.ramBanks), ". ROM size: ", m.ROMSize, " bytes. RAM size: ", m.RAMSize, " bytes")
}

func (m *MBC3) Write(addr types.Word, value byte) {
	switch {
	case addr >= 0x0000 && addr <= 0x1FFF:
		//when in 4/32 mode...
		if m.MaxMemMode == constants.FOURMB_ROM_32KBRAM && m.hasRAM {
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
			log.Println(m.Name + ": Switching MBC3 mode to 16/8")
			m.MaxMemMode = constants.SIXTEENMB_ROM_8KBRAM
		} else {
			log.Println(m.Name + ": Switching MBC3 mode to 4/32")
			m.MaxMemMode = constants.FOURMB_ROM_32KBRAM
		}
	case addr >= 0xA000 && addr <= 0xBFFF:
		if m.hasRAM && m.ramEnabled {
			switch m.MaxMemMode {
			case constants.FOURMB_ROM_32KBRAM:
				m.ramBanks[m.selectedRAMBank][addr-0xA000] = value
			case constants.SIXTEENMB_ROM_8KBRAM:
				m.ramBanks[0][addr-0xA000] = value
			}
		}
	}
}

func (m *MBC3) Read(addr types.Word) byte {
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
		if m.hasRAM && m.ramEnabled {
			switch m.MaxMemMode {
			case constants.FOURMB_ROM_32KBRAM:
				return m.ramBanks[m.selectedRAMBank][addr-0xA000]
			case constants.SIXTEENMB_ROM_8KBRAM:
				return m.ramBanks[0][addr-0xA000]
			}
		}
	}

	return 0x00
}

func (m *MBC3) switchROMBank(bank int) {
	m.selectedROMBank = bank
}

func (m *MBC3) switchRAMBank(bank int) {
	m.selectedRAMBank = bank
}

func (m *MBC3) SaveRam(filename string) error {
	log.Println("Saving RAM")
	return nil
}

func (m *MBC3) LoadRam(filename string) error {
	log.Println("Loading RAM")
	return nil
}
