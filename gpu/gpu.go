package gpu

import (
	"github.com/djhworld/gomeboycolor/types"
	"github.com/go-gl/gl"
	"github.com/go-gl/glfw"
	"log"
)

const PREFIX = "GPU:"

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
	screen [144][160]int
	vram   [8192]byte

	mode    byte
	clock   int
	ly      int
	lcdc    byte
	lyc     byte
	stat    byte
	scrollY byte
	scrollX byte
	bgp     byte

	bgrdOn          bool
	spritesOn       bool
	windowOn        bool
	displayOn       bool
	tileDataSelect  bool
	coincidenceFlag bool
	spriteSize      byte

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

func (g *GPU) Init(title string) error {
	log.Println(PREFIX, "Initialising display")
	var err error

	err = glfw.Init()
	if err != nil {
		return err
	}

	err = glfw.OpenWindow(DISPLAY_WIDTH, DISPLAY_HEIGHT, 0, 0, 0, 0, 0, 0, glfw.Windowed)
	if err != nil {
		return err
	}
	glfw.SetWindowTitle(title)

	//resize function
	onResize := func(w, h int) {
		gl.MatrixMode(gl.PROJECTION)
		gl.LoadIdentity()
		gl.Viewport(0, 0, w, h)
		gl.Ortho(0, float64(w), float64(h), 0, -1, 1)
		gl.ClearColor(0, 0, 0, 0)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		gl.MatrixMode(gl.MODELVIEW)
		gl.LoadIdentity()
	}

	glfw.SetWindowPos(700, 400)
	glfw.SetWindowSizeCallback(onResize)
	gl.ClearColor(0.255, 0.255, 0.255, 0)
	return nil
}

func (g *GPU) Reset() {
	g.Write(LCDC, 0x00)
	g.screen = *new([144][160]int)
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
			//vblank is over
			g.DrawFrame()
			glfw.SwapBuffers()

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
		g.UpdateTile(addr, value)
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
			g.coincidenceFlag = g.lyc == byte(g.ly)
		case SCROLLY:
			g.scrollY = value
		case SCROLLX:
			g.scrollX = value
		case LYC:
			//TODO:
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
		case BGP:
			return g.bgp
		default:
			log.Fatalf("GPU register address %X unknown")
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

		//draw the pixel to the screen data buffer (running through the palette)
		g.screen[g.ly][x] = g.palette[g.tiledata[tileId][tileY][tileX]]

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

func (g *GPU) DrawFrame() {
	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.Enable(gl.POINT_SMOOTH)
	gl.PointSize(1)
	gl.Begin(gl.POINTS)
	for y := 0; y < DISPLAY_HEIGHT; y++ {
		for x := 0; x < DISPLAY_WIDTH; x++ {
			switch g.screen[y][x] {
			case 0:
				gl.Color3ub(235, 235, 235)
			case 1, 2, 3:
				gl.Color3ub(0, 0, 0)
			}

			gl.Vertex2i(x, y)
		}
	}

	gl.End()
}

func (g *GPU) DumpScreen() [144][160]int {
	return g.screen
}
