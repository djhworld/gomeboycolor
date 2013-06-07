package mmu

import (
	"cartridge"
	"components"
	"constants"
	"errors"
	"fmt"
	"log"
	"sort"
	"types"
	"utils"
)

const PREFIX = "MMU"

const (
	CGB_WRAM_BANK_SELECT      types.Word = 0xFF70
	CGB_DOUBLE_SPEED_PREP_REG types.Word = 0xFF4D
)

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
	internalRAM       [8][4096]byte //0xC000 -> 0xDFFF (CGB Working RAM) (8x banks of 4KB)
	internalRAMShadow [7680]byte    //0xE000 -> 0xFDFF
	emptySpace        [52]byte      //0xFF4C -> 0xFF7F
	zeroPageRAM       [128]byte     //0xFF80 - 0xFFFE
	inBootMode        bool
	dmgStatusRegister byte
	DMARegister       byte
	interruptsEnabled byte
	interruptsFlag    byte
	peripheralIOMap   map[types.Word]components.Peripheral

	//CGB features
	cgbWramBankSelectedRegister       byte
	cgbDoubleSpeedPreparationRegister byte
	RunningColorGBHardware            bool
}

func NewGbcMMU() *GbcMMU {
	var mmu *GbcMMU = new(GbcMMU)
	mmu.peripheralIOMap = make(map[types.Word]components.Peripheral)
	mmu.Reset()
	return mmu
}

func (mmu *GbcMMU) Reset() {
	log.Println(PREFIX+": Resetting", PREFIX)
	mmu.inBootMode = true
	mmu.interruptsFlag = 0x00
	mmu.cgbWramBankSelectedRegister = 0x00
	mmu.cgbDoubleSpeedPreparationRegister = 0x00
	mmu.RunningColorGBHardware = false
}

func (mmu *GbcMMU) WriteByte(addr types.Word, value byte) {
	//Check peripherals first
	if p, ok := mmu.peripheralIOMap[addr]; ok {
		p.Write(addr, value)
		return
	}

	switch {
	case addr >= 0x0000 && addr <= 0x9FFF:
		mmu.cartridge.MBC.Write(addr, value)
	//Cartridge External RAM
	case addr >= 0xA000 && addr <= 0xBFFF:
		mmu.cartridge.MBC.Write(addr, value)
	//GB Internal RAM
	case addr >= 0xC000 && addr <= 0xDFFF:
		mmu.WriteToWorkingRAM(addr, value)
		//copy value to shadow if within shadow range
		if addr >= 0xC000 && addr <= 0xDDFF {
			mmu.internalRAMShadow[addr&(0xDDFF-0xC000)] = mmu.ReadByte(addr)
		}
	//INTERRUPT FLAG
	case addr == 0xFF0F:
		mmu.interruptsFlag = value
	//DMA transfer
	case addr == 0xFF46:
		dmaStartAddr := types.Word(value) << 8
		var i types.Word
		for i = 0; i < 0xA0; i++ {
			oamAddr := 0xFE00 + i
			oamData := mmu.ReadByte(dmaStartAddr + i)
			mmu.WriteByte(oamAddr, oamData)
		}
	//Empty but "unusable for I/O"
	case addr > 0xFF4C && addr <= 0xFF7F: //TODO: hmmmm
		switch addr {
		case CGB_DOUBLE_SPEED_PREP_REG:
			if mmu.RunningColorGBHardware == false {
				log.Fatalf("GB is not running in Color GB mode, cannot switch to double speed!")
			}
			mmu.cgbDoubleSpeedPreparationRegister = value
		case 0xFF50:
			mmu.dmgStatusRegister = value
		case 0xFF51:
		case 0xFF52:
		case 0xFF53:
		case 0xFF54:
		case 0xFF55:
			log.Printf("writing 0x%X to CGB HDMA transfer register %s!", value, addr)
		//Color GB Working RAM Bank Selection
		case CGB_WRAM_BANK_SELECT:
			if mmu.RunningColorGBHardware == false {
				log.Fatalf("Cannot write to %s in non-CGB mode!", CGB_WRAM_BANK_SELECT)
			}
			mmu.cgbWramBankSelectedRegister = value
		default:
			mmu.emptySpace[addr-0xFF4D] = value
		}
	//Zero page RAM
	case addr >= 0xFF80 && addr <= 0xFFFF:
		if addr == 0xFFFF {
			mmu.interruptsEnabled = value
		} else {
			mmu.zeroPageRAM[addr&(0xFFFF-0xFF80)] = value
		}
	default:
		//log.Printf("%s: WARNING - Attempting to write 0x%X to address %s, this is invalid/unimplemented", PREFIX, value, addr)
	}
}

func (mmu *GbcMMU) ReadByte(addr types.Word) byte {
	//Check peripherals first
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
		return mmu.cartridge.MBC.Read(addr)
	//ROM Bank 1 (switchable)
	case addr >= 0x4000 && addr <= 0x7FFF:
		return mmu.cartridge.MBC.Read(addr)
	//RAM Bank (switchable)
	case addr >= 0xA000 && addr <= 0xBFFF:
		return mmu.cartridge.MBC.Read(addr)
	//GB Internal RAM
	case addr >= 0xC000 && addr <= 0xDFFF:
		return mmu.ReadFromWorkingRAM(addr)
	//GB Internal RAM shadow
	case addr >= 0xE000 && addr <= 0xFDFF:
		return mmu.internalRAMShadow[addr&(0xFDFF-0xE000)]
	//DMA register
	case addr == 0xFF46:
		return mmu.DMARegister
	//INTERRUPT FLAG
	case addr == 0xFF0F:
		return mmu.interruptsFlag
	//Empty but "unusable for I/O"
	case addr >= 0xFF4C && addr <= 0xFF7F:
		switch addr {
		case 0xFF50:
			return mmu.dmgStatusRegister
		case CGB_DOUBLE_SPEED_PREP_REG:
			return mmu.cgbDoubleSpeedPreparationRegister
		case CGB_WRAM_BANK_SELECT:
			return mmu.cgbWramBankSelectedRegister
		default:
			return mmu.emptySpace[addr-0xFF4C]
		}
	//Zero page RAM
	case addr >= 0xFF80 && addr <= 0xFFFF:
		if addr == 0xFFFF {
			return mmu.interruptsEnabled
		} else {
			return mmu.zeroPageRAM[addr&(0xFFFF-0xFF80)]
		}
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

//When the MMU is in boot mode, the area below 0x0100 is reserved for the BIOS
func (mmu *GbcMMU) SetInBootMode(mode bool) {
	mmu.inBootMode = mode
}

func (mmu *GbcMMU) ConnectPeripheral(p components.Peripheral, startAddr, endAddr types.Word) {
	if startAddr == endAddr {
		log.Printf("%s: Connecting MMU to %s on address %s", PREFIX, p.Name(), startAddr)
		mmu.peripheralIOMap[startAddr] = p
	} else {
		log.Printf("%s: Connecting MMU to %s on address range %s to %s", PREFIX, p.Name(), startAddr, endAddr)
		for addr := startAddr; addr <= endAddr; addr++ {
			mmu.peripheralIOMap[addr] = p
		}
	}
}

//Helper method for connecting peripherals that don't look at contiguous chunks of memory
func (mmu *GbcMMU) ConnectPeripheralOn(p components.Peripheral, addrs ...types.Word) {
	log.Printf("%s: Connecting MMU to %s to address(es): %s", PREFIX, p.Name(), addrs)
	for _, addr := range addrs {
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
		if i%8 == 0 {
			fmt.Println()
		}
	}

	fmt.Println()
}

//Puts BIOS ROM into special area in MMU
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

func (mmu *GbcMMU) IsCartridgeColor() bool {
	return mmu.cartridge.IsColourGB
}

func (mmu *GbcMMU) SaveCartridgeRam(savesDir string) {
	err := mmu.cartridge.SaveRam(savesDir)
	if err != nil {
		log.Println("Error occured attempting to save RAM to disk: ", err)
	}
}

func (mmu *GbcMMU) LoadCartridgeRam(savesDir string) {
	err := mmu.cartridge.LoadRam(savesDir)
	if err != nil {
		log.Println("Error occured attempting to load RAM from disk: ", err)
	}
}

func (mmu *GbcMMU) WriteToWorkingRAM(addr types.Word, value byte) {
	bankAddr := addr & 0x0FFF

	//First area of working RAM is always bank 0 for CGB and Non CGB
	if addr >= 0xC000 && addr <= 0xCFFF {
		mmu.internalRAM[0][bankAddr] = value
	} else if addr >= 0xD000 && addr <= 0xDFFF {
		// In color GB mode the internal RAM is 8x4KB banks (switchable by register 0xFF70)
		if mmu.RunningColorGBHardware {
			bankSelected := int(mmu.cgbWramBankSelectedRegister & 0x07)
			switch {
			//0 and 1 will select bank 1
			case bankSelected <= 1:
				mmu.internalRAM[1][bankAddr] = value
			case bankSelected > 1:
				mmu.internalRAM[bankSelected][bankAddr] = value
			}
		} else {
			//Non-CGB mode is just 8KB of RAM
			mmu.internalRAM[1][bankAddr] = value
		}
	} else {
		log.Fatalf("Address %s is invalid for CGB working RAM!", addr)
	}
}

func (mmu *GbcMMU) ReadFromWorkingRAM(addr types.Word) byte {
	bankAddr := addr & 0x0FFF

	//First area of working RAM is always bank 0 for CGB and Non CGB
	if addr >= 0xC000 && addr <= 0xCFFF {
		return mmu.internalRAM[0][bankAddr]
	} else if addr >= 0xD000 && addr <= 0xDFFF {
		// In color GB mode the internal RAM is 8x4KB banks (switchable by register 0xFF70)
		if mmu.RunningColorGBHardware {
			bankSelected := int(mmu.cgbWramBankSelectedRegister & 0x07)
			switch {
			//0 and 1 will select bank 1
			case bankSelected <= 1:
				return mmu.internalRAM[1][bankAddr]
			case bankSelected > 1:
				return mmu.internalRAM[bankSelected][bankAddr]
			}
		} else {
			//Non-CGB mode is just 8KB of RAM
			return mmu.internalRAM[1][bankAddr]
		}
	} else {
		log.Fatalf("Address %s is invalid for CGB working RAM!", addr)
	}

	return 0x00
}

//USE SHARED CONSTANTS FOR FLAGS AND STUFF TOO - for reuse in the CPU
func (mmu *GbcMMU) RequestInterrupt(interrupt byte) {
	oldVal := mmu.ReadByte(constants.INTERRUPT_FLAG_ADDR)
	switch interrupt {
	case constants.V_BLANK_IRQ:
		mmu.WriteByte(constants.INTERRUPT_FLAG_ADDR, oldVal|constants.V_BLANK_IRQ)
	case constants.LCD_IRQ:
		mmu.WriteByte(constants.INTERRUPT_FLAG_ADDR, oldVal|constants.LCD_IRQ)
	case constants.TIMER_OVERFLOW_IRQ:
		mmu.WriteByte(constants.INTERRUPT_FLAG_ADDR, oldVal|constants.TIMER_OVERFLOW_IRQ)
	case constants.JOYP_HILO_IRQ:
		mmu.WriteByte(constants.INTERRUPT_FLAG_ADDR, oldVal|constants.JOYP_HILO_IRQ)
	default:
		log.Println(PREFIX, "WARNING - interrupt", interrupt, "is currently unimplemented")
	}
}
