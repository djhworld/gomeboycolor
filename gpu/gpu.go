package gpu

import (
	"github.com/djhworld/gomeboycolor/types"
	"github.com/go-gl/glfw"
	"log"
)

const DISPLAY_WIDTH int = 160
const DISPLAY_HEIGHT int = 144

const LCDC types.Word = 0xFF40
const STAT types.Word = 0xFF41
const SCROLLY types.Word = 0xFF42
const SCROLLX types.Word = 0xFF43
const LY types.Word = 0xFF44
const LYC types.Word = 0xFF45

const HBLANK byte = 0x00
const VBLANK byte = 0x01
const OAMREAD byte = 0x02
const VRAMREAD byte = 0x03

type RGBA struct {
	red   byte
	green byte
	blue  byte
	alpha byte
}

type Tile [16]byte

type GPU struct {
	screen   [160][144]RGBA
	mode     byte
	clock    int
	ly       int
	lcdc     byte
	stat     byte
	scrollY  byte
	scrollX  byte
	lyc      byte
	vram     [8192]byte
	tiledata [384]Tile
}

func NewGPU() *GPU {
	return new(GPU)
}

func (g *GPU) Init(title string) error {
	var err error

	err = glfw.Init()
	if err != nil {
		return err
	}

	err = glfw.OpenWindow(DISPLAY_WIDTH, DISPLAY_HEIGHT, 0, 0, 0, 0, 0, 0, glfw.Windowed)
	if err != nil {
		return err
	}

	glfw.SetWindowPos(800, 0)
	return nil
}

func (g *GPU) Reset() {
	g.screen = *new([160][144]RGBA)
	g.mode = 0
	g.ly = 0
	g.clock = 0

	for i := 0; i < len(g.tiledata); i++ {
		for j := 0; j < 16; j++ {
			g.tiledata[i][j] = 0
		}
	}

	for i := 0; i < 160; i++ {
		for j := 0; j < 144; j++ {
			g.screen[i][j] = RGBA{100, 100, 100, 255}
		}
	}
}

func (g *GPU) Step(t int) {
	g.clock += t
	log.Println("GPU Clock =", g.clock)
	log.Println("GPU Mode =", g.mode)
	log.Println("GPU Line =", g.ly)
	if g.ly < DISPLAY_HEIGHT {
		if g.clock >= 204 {
			g.mode = HBLANK
			if g.clock >= 456 {
				//TODO: draw line
				g.clock = 0
				g.ly += 1
			}
		} else if g.clock >= 172 {
			g.mode = VRAMREAD
		} else {
			g.mode = OAMREAD
		}
	} else {
		//TODO: blit frame to screen
		g.mode = VBLANK
		//for each step revert back 10 lines
		if g.ly < 154 {
			if g.clock >= 456 {
				g.clock = 0
				g.ly += 1
			}
		} else {
			g.clock = 0
			g.ly = 0
		}
	}
}

//Called from mmu
func (g *GPU) Write(addr types.Word, value byte) {
	switch {
	case addr >= 0x8000 && addr <= 0x9FFF:
		g.vram[addr&0x1FFF] = value

		//update is tile data
		if addr <= 0x97FF {
			g.UpdateTile(addr, value)
		}
	default:
		switch addr {
		case LCDC:
			log.Printf("LCDC being written to %X ---> %X!!!!!!", addr, value)
			//TODO: calculate
		case STAT:
			log.Printf("STAT being written to %X ---> %X!!!!!!", addr, value)
			//TODO: calculate
		case SCROLLY:
			log.Printf("SCROLL Y being written to %X ---> %X!!!!!!", addr, value)
			//TODO: calculate
		case SCROLLX:
			log.Printf("SCROLL X being written to %X ---> %X!!!!!!", addr, value)
			//TODO: calculate
		case LY:
			log.Printf("LY being written to %X ---> %X!!!!!!", addr, value)
			//TODO: calculate
		case LYC:
			log.Printf("LYC being written to %X ---> %X!!!!!!", addr, value)
			//TODO: calculate
		}
	}
}

//Called from mmu
func (g *GPU) Read(addr types.Word) byte {
	log.Printf("GPU is being read from %X", addr)

	switch {
	case addr >= 0x8000 && addr <= 0x9FFF:
		return g.vram[addr&0x1FFF]
	default:
		switch addr {
		case LCDC:
			return g.lcdc
		case STAT:
			return g.stat
		case SCROLLY:
			return g.scrollY
		case SCROLLX:
			return g.scrollX
		case LY:
			return byte(g.ly)
		case LYC:
			return g.lyc
		default:
			log.Fatalf("GPU register address %X unknown")
		}
	}
	return 0x00
}

func (g *GPU) UpdateTile(addr types.Word, value byte) {
	var tileId types.Word = (addr & 0x17FF) >> 4
	g.tiledata[int(tileId)][tileId%16] = value
}
