package cartridge

import (
	"errors"
	"fmt"
	"strings"
)

const (
	MBC0          = 0x00
	MBC1          = 0x01
	MBC1_RAM      = 0x02
	MBC1_RAM_BATT = 0x03
)

type Cartridge struct {
	Type           CartridgeType
	Title          string
	IsColourGBCart bool
	ROM            []byte
}

func NewCartridge(rom []byte) (*Cartridge, error) {
	var cart *Cartridge = new(Cartridge)

	ctype := rom[0x0147]
	//validate
	if v, ok := CartridgeTypes[ctype]; !ok {
		return nil, errors.New(fmt.Sprintf("Unknown cartridge type: %X for ROM", ctype))
	} else {
		cart.Type = v
	}

	if rom[0x0143] == 0x80 {
		cart.IsColourGBCart = true
	}

	cart.Title = strings.TrimSpace(string(rom[0x0134:0x0142]))
	cart.ROM = rom
	return cart, nil
}

func (c *Cartridge) String() string {
	startingString := "Gameboy"
	if c.IsColourGBCart {
		startingString += " Color"
	}
	return fmt.Sprint("\n", startingString, " Cartridge\n--------------------------------------\n") +
		fmt.Sprintf("Title: \"%s\"\n", c.Title) +
		fmt.Sprintf("Type: %s (0x%X)\n", c.Type.Description, c.Type.ID) +
		fmt.Sprintf("Size: %d bytes\n", len(c.ROM)) +
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
