package dma

import (
	"log"

	"github.com/djhworld/gomeboycolor/components"
	"github.com/djhworld/gomeboycolor/constants"
	"github.com/djhworld/gomeboycolor/mmu"
	"github.com/djhworld/gomeboycolor/types"
)

const (
	NAME                     string     = "HDMA"
	CGB_HDMA_SOURCE_HIGH_REG types.Word = 0xFF51
	CGB_HDMA_SOURCE_LOW_REG  types.Word = 0xFF52
	CGB_HDMA_DEST_HIGH_REG   types.Word = 0xFF53
	CGB_HDMA_DEST_LOW_REG    types.Word = 0xFF54
	CGB_HDMA_REG             types.Word = 0xFF55
)

type HDMARegisters struct {
	srcHigh byte
	srcLow  byte
	dstHigh byte
	dstLow  byte
}

type HDMATransfer struct {
	Source      types.Word
	Destination types.Word
	Length      int
}

type HDMA struct {
	running          bool
	displayOn        bool
	isHblankTransfer bool
	gpuMode          byte
	registers        HDMARegisters
	hdmaTransfer     HDMATransfer
	mmu              *mmu.GbcMMU
}

func NewHDMA(mmu *mmu.GbcMMU) *HDMA {
	h := new(HDMA)
	h.mmu = mmu
	h.Reset()
	return h
}

func (h *HDMA) Name() string {
	return NAME
}

func (h *HDMA) Step() {
	if !h.IsRunning() {
		return
	}

	// Write 1 block (16 bytes)
	for i := types.Word(0x0000); i < 0x0010; i++ {
		data := h.mmu.ReadByte(h.hdmaTransfer.Source + i)
		h.mmu.WriteByte(h.hdmaTransfer.Destination+i, data)
	}

	// Move src/destination forward to next block
	h.hdmaTransfer.Source += 0x0010
	h.hdmaTransfer.Destination += 0x0010

	h.hdmaTransfer.Length--
	if h.hdmaTransfer.Length <= 0 {
		h.running = false
		h.hdmaTransfer.Length = 0x7F
	} else if h.gpuMode == constants.HBLANK_MODE {
		// next HDMA will run on next HBLANK
		h.gpuMode = 0xFF
	}
}

func (h *HDMA) OnGPUModeChange(mode byte) {
	h.gpuMode = mode
}

func (h *HDMA) OnDisplayChange(on bool) {
	h.displayOn = on
}

// HDMA Source/Dest registers always return 0xFF when read
// Source: https://github.com/AntonioND/giibiiadvance/blob/master/docs/TCAGBD.pdf
func (h *HDMA) Read(address types.Word) byte {
	switch address {
	case CGB_HDMA_SOURCE_HIGH_REG:
	case CGB_HDMA_SOURCE_LOW_REG:
	case CGB_HDMA_DEST_HIGH_REG:
	case CGB_HDMA_DEST_LOW_REG:
		return 0xFF
	case CGB_HDMA_REG:
		if h.running {
			return 0x00
		} else {
			return (1 << 7) | byte(h.hdmaTransfer.Length)
		}
	default:
		log.Println(NAME + "- WARN - Unsupported register in HDMA")
	}
	return 0xFF
}

func (h *HDMA) Write(address types.Word, value byte) {
	switch address {
	case CGB_HDMA_SOURCE_HIGH_REG:
		h.registers.srcHigh = value
	case CGB_HDMA_SOURCE_LOW_REG:
		// lower 4 bits are ignored
		h.registers.srcLow = value & 0xF0
	case CGB_HDMA_DEST_HIGH_REG:
		// upper 3 bits are ignored
		h.registers.dstHigh = value & 0x1F
	case CGB_HDMA_DEST_LOW_REG:
		// lower 4 bits are ignored
		h.registers.dstLow = value & 0xF0
	case CGB_HDMA_REG:
		if h.running && (value&(1<<7)) == 0 {
			h.running = false
		} else {
			h.startTransfer(value)
		}
	}

}

func (h *HDMA) LinkIRQHandler(m components.IRQHandler) {

}

func (h *HDMA) IsRunning() bool {
	if !h.running {
		return false
	} else if h.isHblankTransfer && (h.gpuMode == constants.HBLANK_MODE || h.displayOn == false) {
		return true
	} else if !h.isHblankTransfer {
		return true
	} else {
		return false
	}
}

func (h *HDMA) computeTransferLocations() {
	h.hdmaTransfer.Source = (types.Word(h.registers.srcHigh) << 8) | types.Word(h.registers.srcLow)
	h.hdmaTransfer.Destination = (types.Word(h.registers.dstHigh) << 8) | types.Word(h.registers.dstLow)

	// align destination to be at least 0x8000 (VRAM)
	// Not doing this was causing a huge pain in some ROMS
	h.hdmaTransfer.Destination |= 0x8000
}

func (h *HDMA) startTransfer(value byte) {
	h.isHblankTransfer = (value & 0x80) != 0
	h.running = true
	h.hdmaTransfer.Length = int(value&0x7F) + 1
	h.computeTransferLocations()
}

func (h *HDMA) Reset() {
	h.hdmaTransfer = HDMATransfer{}
	h.registers = HDMARegisters{0xFF, 0xFF, 0xFF, 0xFF}
	h.isHblankTransfer = false
}
