package cartridge

import (
	"errors"
	"fmt"
	"strings"
	"utils"
)

const (
	MBC_0          = 0x00
	MBC_1          = 0x01
	MBC_1_RAM      = 0x02
	MBC_1_RAM_BATT = 0x03
)

type CartridgeType struct {
	ID          byte
	Description string
}

var CartridgeTypes map[byte]CartridgeType = map[byte]CartridgeType{
	MBC_0:          CartridgeType{MBC_0, "ROM ONLY"},
	MBC_1:          CartridgeType{MBC_1, "ROM+MBC1"},
	MBC_1_RAM:      CartridgeType{MBC_1_RAM, "ROM+MBC1+RAM"},
	MBC_1_RAM_BATT: CartridgeType{MBC_1_RAM_BATT, "ROM+MBC1+RAM+BATT"},
}

type Cartridge struct {
	Title      string
	IsColourGB bool
	Type       CartridgeType
	ROMSize    int
	RAMSize    int
	IsJapanese bool
	MBC        MemoryBankController
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
	}

	c.Title = strings.TrimSpace(string(rom[0x0134:0x0142]))
	c.IsColourGB = (rom[0x0143] == 0x80)

	ctype := rom[0x0147]
	//validate
	if v, ok := CartridgeTypes[ctype]; !ok {
		return errors.New(fmt.Sprintf("Unknown cartridge type: %X for ROM", ctype))
	} else {
		c.Type = v
	}

	if romSize := rom[0x0148]; romSize > 0x06 {
		return errors.New(fmt.Sprintf("Handling for ROM size id: 0x%X is currently unimplemented", romSize))
	} else {
		c.ROMSize = 0x8000 << romSize
	}

	switch rom[0x0149] {
	case 0x00:
		c.RAMSize = 0
	case 0x01:
		c.RAMSize = 2048
	case 0x02:
		c.RAMSize = 8192
	case 0x03:
		c.RAMSize = 32768
	case 0x04:
		c.RAMSize = 131072
	}

	c.IsJapanese = (rom[0x014A] == 0x00)

	switch c.Type.ID {
	case MBC_0:
		c.MBC = NewMBC0(rom)
	case MBC_1, MBC_1_RAM, MBC_1_RAM_BATT:
		c.MBC = NewMBC1(rom, c.ROMSize, c.RAMSize)

	default:
		return errors.New("Error: Cartridge type " + utils.ByteToString(c.Type.ID) + " is currently unsupported")
	}

	return nil
}

func (c *Cartridge) String() string {
	startingString := "Gameboy"
	if c.IsColourGB {
		startingString += " Color"
	}

	var destinationRegion string
	if c.IsJapanese {
		destinationRegion = "Japanese"
	} else {
		destinationRegion = "Non-Japanese"
	}

	var header []string = []string{
		fmt.Sprintf(utils.PadRight("Title:", 19, " ")+"%s", c.Title),
		fmt.Sprintf(utils.PadRight("Type:", 19, " ")+"%s %s", c.Type.Description, utils.ByteToString(c.Type.ID)),
		fmt.Sprintf(utils.PadRight("ROM Size:", 19, " ")+"%d", c.ROMSize),
		fmt.Sprintf(utils.PadRight("RAM Size:", 19, " ")+"%d", c.RAMSize),
		fmt.Sprintf(utils.PadRight("Destination code:", 19, " ")+"%s", destinationRegion),
	}

	return fmt.Sprint("\n", startingString, " Cartridge\n--------------------------------------\n") +
		fmt.Sprintln(strings.Join(header, "\n")) +
		fmt.Sprint("--------------------------------------\n")
}
