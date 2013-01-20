package cartridge

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"types"
)

const (
	MBC0          = 0x00
	MBC1          = 0x01
	MBC1_RAM      = 0x02
	MBC1_RAM_BATT = 0x03
)

type Cartridge struct {
	Type              CartridgeType
	Title             string
	IsColourGBCart    bool
	Size              int
	romBank0          []byte
	switchableROMBank [128][]byte
	switchableRAMBank [16][]byte
}

func NewCartridge(rom []byte) (*Cartridge, error) {
	var cart *Cartridge = new(Cartridge)
	err := cart.Init(rom)

	if err != nil {
		return nil, err
	}

	return cart, nil
}

func (c *Cartridge) Init(rom []byte) error {
	if size := len(rom); size < 32768 {
		return errors.New(fmt.Sprintf("ROM size %d is too small", size))
	} else {
		c.Size = size
	}

	ctype := rom[0x0147]
	//validate
	if v, ok := CartridgeTypes[ctype]; !ok {
		return errors.New(fmt.Sprintf("Unknown cartridge type: %X for ROM", ctype))
	} else {
		c.Type = v
	}

	if rom[0x0143] == 0x80 {
		c.IsColourGBCart = true
	}

	c.Title = strings.TrimSpace(string(rom[0x0134:0x0142]))

	//always static
	c.romBank0 = rom[0x0000:0x4000]

	switch c.Type.ID {
	case MBC0:
		c.switchableROMBank[0] = rom[0x4000:0x8000]
		c.switchableRAMBank[0] = make([]byte, 8192)
	}

	return nil
}

func (c *Cartridge) ReadByteFromROM(addr types.Word) byte {
	if addr > 0x3FFF {
		log.Fatalf("Cannot read from ROM, address: %X is invalid", addr)
	}
	return c.romBank0[addr]
}

func (c *Cartridge) ReadByteFromSwitchableROM(bank int, addr types.Word) byte {
	if addr < 0x4000 || addr > 0x7FFF {
		log.Fatalf("Cannot read from switchable ROM bank %d, address: %X is invalid", bank, addr)
	}
	return c.switchableROMBank[bank][addr&(0x7FFF-0x4000)]
}

func (c *Cartridge) ReadByteFromSwitchableRAM(bank int, addr types.Word) byte {
	if addr < 0xA000 || addr > 0xBFFF {
		log.Fatalf("Cannot read from switchable ROM bank %d, address: %X is invalid", bank, addr)
	}
	return c.switchableRAMBank[bank][addr&(0xBFFF-0xA000)]
}

func (c *Cartridge) WriteByteToSwitchableRAM(bank int, addr types.Word, value byte) {
	if addr < 0xA000 || addr > 0xBFFF {
		log.Fatalf("Cannot write to switchable ROM bank %d, address: %X is invalid", bank, addr)
	}
	c.switchableRAMBank[bank][addr&(0xBFFF-0xA000)] = value
}

func (c *Cartridge) String() string {
	startingString := "Gameboy"
	if c.IsColourGBCart {
		startingString += " Color"
	}
	return fmt.Sprint("\n", startingString, " Cartridge\n--------------------------------------\n") +
		fmt.Sprintf("Title: \"%s\"\n", c.Title) +
		fmt.Sprintf("Type: %s (0x%X)\n", c.Type.Description, c.Type.ID) +
		fmt.Sprintf("Size: %d bytes\n", c.Size) +
		fmt.Sprint("--------------------------------------\n")
}

type CartridgeType struct {
	ID          byte
	Description string
}

var CartridgeTypes map[byte]CartridgeType = map[byte]CartridgeType{
	MBC0:          CartridgeType{MBC0, "ROM ONLY"},
	MBC1:          CartridgeType{MBC1, "ROM+MBC1"},
	MBC1_RAM:      CartridgeType{MBC1_RAM, "ROM+MBC1+RAM"},
	MBC1_RAM_BATT: CartridgeType{MBC1_RAM_BATT, "ROM+MBC1+RAM+BATT"},
}
