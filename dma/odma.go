package dma

import (
	"github.com/djhworld/gomeboycolor/components"
	"github.com/djhworld/gomeboycolor/mmu"
	"github.com/djhworld/gomeboycolor/types"
)

const (
	OAMDMA_NAME  string     = "OAM-DMA"
	DMA_TRANSFER types.Word = 0xFF46
)

type OAMDMA struct {
	running      bool
	cycles       int
	transferFrom types.Word
	mmu          *mmu.GbcMMU
}

func NewOAMDMA(mmu *mmu.GbcMMU) *OAMDMA {
	o := new(OAMDMA)
	o.mmu = mmu
	o.Reset()
	return o
}

func (o *OAMDMA) Name() string {
	return OAMDMA_NAME
}

func (o *OAMDMA) Step(cycles int) {
	o.cycles += cycles

	if o.running {
		if o.cycles >= 648 {

			o.running = false
			o.cycles = 0
			o.doInstantDMATransfer(o.transferFrom, 0xFE00, 10, 16)
		}
	}
}

func (o *OAMDMA) Read(address types.Word) byte {
	return 0x00
}

func (o *OAMDMA) Write(address types.Word, value byte) {
	o.transferFrom = types.Word(value) << 8
	o.running = true
}

func (o *OAMDMA) LinkIRQHandler(m components.IRQHandler) {
}

func (o *OAMDMA) IsRunning() bool {
	return o.running
}

func (o *OAMDMA) Reset() {
}

func (o *OAMDMA) doInstantDMATransfer(startAddress, destinationAddr types.Word, blocks, blockSize int) {
	length := types.Word(blockSize * blocks)
	var i types.Word = 0x0000
	for ; i < length; i++ {
		data := o.mmu.ReadByte(startAddress + i)
		o.mmu.WriteByte(destinationAddr+i, data)
	}
}
