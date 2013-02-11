package gpu

import (
	"inputoutput"
	"log"
	"types"
)

const NAME = "GPU"
const PREFIX = NAME + ":"

const DISPLAY_WIDTH int = 160
const DISPLAY_HEIGHT int = 144

const TILEMAP0 = 0x9800
const TILEMAP1 = 0x9C00

const LCDC types.Word = 0xFF40
const STAT types.Word = 0xFF41
const SCROLLY types.Word = 0xFF42
const SCROLLX types.Word = 0xFF43
const LY types.Word = 0xFF44
const LYC types.Word = 0xFF45
const BGP types.Word = 0xFF47
const WX types.Word = 0xFF4B
const WY types.Word = 0xFF4A

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

type RawTile [16]byte
type Tile [8][8]int

type GPU struct {
	screenData [144][160]int
	screen     inputoutput.Screen
	vram       [8192]byte
	oamRam     [160]byte

	mode            byte
	clock           int
	ly              int
	lcdc            byte
	lyc             byte
	stat            byte
	scrollY         byte
	scrollX         byte
	bgp             byte
	coincidenceFlag bool

	bgrdOn         bool
	spritesOn      bool
	windowOn       bool
	displayOn      bool
	tileDataSelect bool
	spriteSize     byte

	bgTilemap     types.Word
	windowTilemap types.Word
	rawTiledata   [384]RawTile
	tiledata      [384]Tile

	palette [4]int
}

func NewGPU() *GPU {
	var g *GPU = new(GPU)
	g.Reset()
	return g
}

func (g *GPU) LinkScreen(screen inputoutput.Screen) {
	g.screen = screen
	log.Println(PREFIX, "Linked screen to GPU")
}

func (g *GPU) Name() string {
	return NAME
}

func (g *GPU) Reset() {
	g.Write(LCDC, 0x00)
	g.screenData = *new([144][160]int)
	g.mode = 0
	g.ly = 0
	g.clock = 0
}

func (g *GPU) Step(t int) {
	g.clock += t
	if g.ly < DISPLAY_HEIGHT {
		if g.clock >= 204 {
			g.mode = HBLANK
			if g.clock >= 456 {
				g.RenderLine()
				g.clock = 0
				g.ly += 1
			}
		} else if g.clock >= 172 {
			g.mode = VRAMREAD
		} else {
			g.mode = OAMREAD
		}
	} else {
		g.mode = VBLANK
		g.Write(LCDC, g.lcdc^0x80)

		//for each step revert back 10 lines
		if g.ly < 154 {
			if g.clock >= 456 {
				g.clock = 0
				g.ly += 1
			}
		} else {
			//vblank is over, draw to screen
			g.screen.DrawFrame(&g.screenData)
			g.clock = 0
			g.ly = 0
		}
	}

	g.coincidenceFlag = false
	if byte(g.ly) == g.lyc && (g.stat&0x40) == 0x40 {
		g.coincidenceFlag = true
	}
}

//Called from mmu
func (g *GPU) Write(addr types.Word, value byte) {
	switch {
	case addr >= 0x8000 && addr <= 0x9FFF:
		g.vram[addr&0x1FFF] = value
		g.UpdateTile(addr, value)
	case addr >= 0xFE00 && addr <= 0xFE9F:
		g.oamRam[addr&0x009F] = value
	default:
		switch addr {
		case LCDC:
			g.lcdc = value

			g.displayOn = value&0x80 == 0x80 //bit 7

			if value&0x40 == 0x40 { //bit 6 
				g.windowTilemap = TILEMAP1
			} else {
				g.windowTilemap = TILEMAP0
			}

			g.windowOn = value&0x20 == 0x20       //bit 5 
			g.tileDataSelect = value&0x10 == 0x10 //bit 4

			if value&0x08 == 0x08 { //bit 3
				g.bgTilemap = TILEMAP1
			} else {
				g.bgTilemap = TILEMAP0
			}

			if value&0x04 == 0x04 { //bit 2
				g.spriteSize = 128
			} else {
				g.spriteSize = 64
			}

			g.spritesOn = value&0x02 == 0x02 //bit 1
			g.bgrdOn = value&0x01 == 0x01    //bit 0
		case STAT:
			g.stat = value
		case SCROLLY:
			g.scrollY = value
		case SCROLLX:
			g.scrollX = value
		case WX:
			log.Println(PREFIX, "Writing to WX!")
		case WY:
			log.Println(PREFIX, "Writing to WY!")
		case LYC:
			g.lyc = value
		case BGP:
			g.bgp = value
			g.palette[0] = int(value & 0x03)
			g.palette[1] = int((value >> 2) & 0x03)
			g.palette[2] = int((value >> 4) & 0x03)
			g.palette[3] = int((value >> 6) & 0x03)
		}
	}
}

//Called from mmu
func (g *GPU) Read(addr types.Word) byte {
	switch {
	case addr >= 0x8000 && addr <= 0x9FFF:
		return g.vram[addr&0x1FFF]
	case addr >= 0xFE00 && addr <= 0xFE9F:
		return g.oamRam[addr&0x009F]
	default:
		switch addr {
		case LCDC:
			return g.lcdc
		case STAT:
			g.stat = 0x00
			if g.coincidenceFlag {
				g.stat ^= 0x44
			}

			switch g.mode {
			case HBLANK:
				g.stat &^= 0x33
			case VBLANK:
				g.stat ^= 0x11
			case OAMREAD:
				g.stat ^= 0x22
			case VRAMREAD:
				g.stat ^= 0x33
			}

			return g.stat
		case SCROLLY:
			return g.scrollY
		case SCROLLX:
			return g.scrollX
		case LY:
			return byte(g.ly)
		case LYC:
			return g.lyc
		case BGP:
			return g.bgp
		default:
			log.Printf(PREFIX+" WARNING: register address %s unknown", addr)
			return 0x00
		}
	}
	return 0x00
}

//Update the tile at address with value
func (g *GPU) UpdateTile(addr types.Word, value byte) {

	//get the ID of the tile being updated (between 0 and 383)
	var tileId types.Word = (addr & 0x17FF) >> 4
	g.rawTiledata[tileId][addr%16] = value

	recalcTile := func(rawtile RawTile) Tile {
		var tile Tile
		for tileY := 0; tileY < 8; tileY++ {
			lineLo, lineHi := int(rawtile[tileY*2]), int(rawtile[(tileY*2)+1])
			var tileX uint = 0
			for ; tileX < 8; tileX++ {
				tile[tileY][tileX] = ((lineLo >> (7 - tileX) & 1) | (lineHi>>(7-tileX)&1)<<1)
			}
		}

		return tile
	}

	g.tiledata[tileId] = recalcTile(g.rawTiledata[tileId])
}

func (g *GPU) RenderLine() {
	var mapoffset types.Word = g.bgTilemap + ((types.Word(g.ly+int(g.scrollY)))>>3)<<5
	var lineoffset types.Word = (types.Word(g.scrollX) >> 3) % 32
	tileY := (g.ly + int(g.scrollY)) % 8
	tileX := int(g.scrollX) % 8

	//get the ID of the tile being drawn
	//TODO: calculate value if in tileset #1
	tileId := g.Read(types.Word(mapoffset + lineoffset))

	for x := 0; x < DISPLAY_WIDTH; x++ {

		//draw the pixel to the screenData data buffer (running through the palette)
		g.screenData[g.ly][x] = g.palette[g.tiledata[tileId][tileY][tileX]]

		//move along line in tile until you reach the end
		tileX++
		if tileX == 8 {
			tileX = 0
			lineoffset = (lineoffset + 1) % 32
			//get next tile in line
			tileId = g.Read(types.Word(mapoffset + lineoffset))
		}
	}
}

func (g *GPU) DumpScreen() [144][160]int {
	return g.screenData
}
