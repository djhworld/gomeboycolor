package gpu

import (
	"components"
	"constants"
	"fmt"
	"inputoutput"
	"log"
	"types"
)

const NAME = "GPU"
const PREFIX = NAME + ":"

const DISPLAY_WIDTH int = 160
const DISPLAY_HEIGHT int = 144

const (
	TILEMAP0  types.Word = 0x9800
	TILEMAP1             = 0x9C00
	TILEDATA0            = 0x8800
	TILEDATA1            = 0x8000
)

const LCDC types.Word = 0xFF40
const STAT types.Word = 0xFF41
const SCROLLY types.Word = 0xFF42
const SCROLLX types.Word = 0xFF43
const LY types.Word = 0xFF44
const LYC types.Word = 0xFF45
const BGP types.Word = 0xFF47
const OBJECTPALETTE_0 types.Word = 0xFF48
const OBJECTPALETTE_1 types.Word = 0xFF49
const WX types.Word = 0xFF4B
const WY types.Word = 0xFF4A

const HBLANK byte = 0x00
const VBLANK byte = 0x01
const OAMREAD byte = 0x02
const VRAMREAD byte = 0x03
const Sprite8x16Mode byte = 0
const Sprite8x8Mode byte = 1

var GBColours []types.RGB = []types.RGB{
	types.RGB{Red: 235, Green: 235, Blue: 235},
	types.RGB{Red: 196, Green: 196, Blue: 196},
	types.RGB{Red: 96, Green: 96, Blue: 96},
	types.RGB{Red: 0, Green: 0, Blue: 0},
}

type RawTile [16]byte
type Tile [8][8]int
type Palette [4]types.RGB

type GPU struct {
	screenData            [144][160]types.RGB
	screen                inputoutput.Screen
	irqHandler            components.IRQHandler
	vram                  [8192]byte
	oamRam                [160]byte
	vBlankInterruptThrown bool
	lcdInterruptThrown    bool

	mode    byte
	clock   int
	ly      int
	lcdc    byte
	lyc     byte
	stat    byte
	scrollY byte
	scrollX byte
	windowX byte
	windowY byte
	bgp     byte
	obp0    byte
	obp1    byte

	bgrdOn         bool
	spritesOn      bool
	windowOn       bool
	displayOn      bool
	tileDataSelect types.Word
	spriteSizeMode byte

	bgTilemap     types.Word
	windowTilemap types.Word
	rawTiledata   [512]RawTile
	tiledata      [512]Tile
	sprites8x8    [40]Sprite
	sprites8x16   [40]Sprite

	bgPalette      Palette
	objectPalettes [2]Palette
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

func (g *GPU) LinkIRQHandler(m components.IRQHandler) {
	g.irqHandler = m
	log.Println(PREFIX, "Linked IRQ Handler to GPU")
}

func (g *GPU) Name() string {
	return NAME
}

func (g *GPU) Reset() {
	log.Println("Resetting", g.Name())
	g.Write(LCDC, 0x00)
	g.screenData = *new([144][160]types.RGB)
	g.mode = 0
	g.ly = 0
	g.clock = 0
	g.vBlankInterruptThrown = false
	g.lcdInterruptThrown = false

	for i := 0; i < 40; i++ {
		g.sprites8x8[i] = NewSprite8x8()
		g.sprites8x16[i] = NewSprite8x16()
	}
}

func (g *GPU) Step(t int) {
	if !g.displayOn {
		g.ly = 0
		g.clock = 456
		g.mode = HBLANK
	} else {
		if g.ly >= 144 {
			g.mode = VBLANK
		} else if g.clock >= 456-80 {
			g.mode = OAMREAD
		} else if g.clock >= 456-80-172 {
			g.mode = VRAMREAD
		} else {
			g.mode = HBLANK
		}
	}

	if !g.displayOn {
		return
	}

	g.clock -= t

	if g.clock <= 0 {
		if g.ly < 144 {
			if g.displayOn {
				if g.bgrdOn {
					g.RenderBackgroundScanline()
				}

				if g.windowOn {
					g.RenderWindowScanline()
				}
			}
		} else if g.ly == 144 {
			if g.spritesOn {
				g.RenderSprites()
			}

			//throw vblank interrupt
			if g.vBlankInterruptThrown == false {
				g.irqHandler.RequestInterrupt(constants.V_BLANK_IRQ)
				g.vBlankInterruptThrown = true
			}
			g.screen.DrawFrame(&g.screenData)

		} else if g.ly > 153 {
			g.vBlankInterruptThrown = false
			g.ly = 0
		}

		g.clock += 456
		g.ly += 1

		if byte(g.ly) == g.lyc {
			g.stat |= 0x04
		}

		g.CheckForLCDCSTATInterrupt()
	}

}

func (g *GPU) CheckForLCDCSTATInterrupt() {
	switch {
	case byte(g.ly) == g.lyc && (g.Read(STAT)&0x40) == 0x40:
		g.irqHandler.RequestInterrupt(constants.LCD_IRQ)
	case g.mode == OAMREAD && (g.Read(STAT)&0x20) == 0x20:
		g.irqHandler.RequestInterrupt(constants.LCD_IRQ)
	case g.mode == VBLANK && (g.Read(STAT)&0x10) == 0x10:
		g.irqHandler.RequestInterrupt(constants.LCD_IRQ)
	case g.mode == HBLANK && (g.Read(STAT)&0x08) == 0x08:
		g.irqHandler.RequestInterrupt(constants.LCD_IRQ)
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
		g.UpdateSprite(addr, value)
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

			g.windowOn = value&0x20 == 0x20 //bit 5

			if value&0x10 == 0x10 { //bit 4
				g.tileDataSelect = TILEDATA1
			} else {
				g.tileDataSelect = TILEDATA0
			}

			if value&0x08 == 0x08 { //bit 3
				g.bgTilemap = TILEMAP1
			} else {
				g.bgTilemap = TILEMAP0
			}

			if value&0x04 == 0x04 { //bit 2
				g.spriteSizeMode = Sprite8x16Mode
			} else {
				g.spriteSizeMode = Sprite8x8Mode
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
			g.windowX = value
		case WY:
			g.windowY = value
		case LY:
			g.ly = 0
		case LYC:
			g.lyc = value
		case BGP:
			g.bgp = value
			g.bgPalette = g.byteToPalette(value)
		case OBJECTPALETTE_0:
			g.obp0 = value
			g.objectPalettes[0] = g.byteToPalette(value)
		case OBJECTPALETTE_1:
			g.obp1 = value
			g.objectPalettes[1] = g.byteToPalette(value)
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
			return (g.mode | g.stat&0xF8)
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
		case OBJECTPALETTE_0:
			return g.obp0
		case OBJECTPALETTE_1:
			return g.obp1
		case WX:
			return g.windowX
		case WY:
			return g.windowY
		default:
			log.Printf(PREFIX+" WARNING: register address %s unknown", addr)
			return 0x00
		}
	}
}

func (g *GPU) UpdateSprite(addr types.Word, value byte) {
	var spriteId types.Word = (addr & 0x00FF) / 4
	if g.spriteSizeMode == Sprite8x8Mode {
		g.sprites8x8[spriteId].UpdateSprite(addr, value)
	} else {
		g.sprites8x16[spriteId].UpdateSprite(addr, value)
	}
}

//Update the tile at address with value
func (g *GPU) UpdateTile(addr types.Word, value byte) {
	//get the ID of the tile being updated (between 0 and 383)
	var tileId types.Word = ((addr & 0x1FFF) >> 4) & 511
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

func (g *GPU) RenderBackgroundScanline() {
	//find where in the tile map we are related to the current scan line + scroll Y (wraps around)
	screenYAdjusted := int(g.ly) + int(g.scrollY)
	var initialTilemapOffset types.Word = g.bgTilemap + types.Word(screenYAdjusted)%256/8*32
	var initialLineOffset types.Word = types.Word(g.scrollX) / 8 % 32

	//find where in the tile we are
	initialTileX := int(g.scrollX) % 8
	initialTileY := screenYAdjusted % 8

	//screen will always draw from X = 0
	g.DrawScanline(initialTilemapOffset, initialLineOffset, 0, initialTileX, initialTileY)
}

func (g *GPU) RenderWindowScanline() {
	screenYAdjusted := g.ly - int(g.windowY)

	if (g.windowX >= 0 && g.windowX < 167) && (g.windowY >= 0 && g.windowY < 144) && screenYAdjusted >= 0 {
		var initialTilemapOffset types.Word = g.windowTilemap + types.Word(screenYAdjusted)/8*32
		var initialLineOffset types.Word = 0
		var screenXAdjusted int = int((g.windowX - 7) % 255)

		//find where in the tile we are
		initialTileX := screenXAdjusted % 8
		initialTileY := screenYAdjusted % 8

		g.DrawScanline(initialTilemapOffset, initialLineOffset, screenXAdjusted, initialTileX, initialTileY)
	}
}

func (g *GPU) DrawScanline(tilemapOffset, lineOffset types.Word, screenX, tileX, tileY int) {
	//get tile to start from
	tileId := g.calculateTileNo(tilemapOffset, lineOffset)

	for ; screenX < DISPLAY_WIDTH; screenX++ {
		//draw the pixel to the screenData data buffer (running through the bgPalette)
		color := g.bgPalette[g.tiledata[tileId][tileY][tileX]]
		g.screenData[g.ly][screenX] = color

		//move along line in tile until you reach the end
		tileX++
		if tileX == 8 {
			tileX = 0
			lineOffset = (lineOffset + 1) % 32

			//get next tile in line
			tileId = g.calculateTileNo(tilemapOffset, lineOffset)
		}
	}
}

//method to calculate the tilenumber within the tilemap
func (g *GPU) calculateTileNo(tilemapOffset types.Word, lineOffset types.Word) int {
	tileId := int(g.Read(types.Word(tilemapOffset + lineOffset)))

	//if tile data is 0 then it is signed
	if g.tileDataSelect == TILEDATA0 {
		if tileId < 128 {
			tileId += 256
		}
	}
	return tileId
}

func (g *GPU) RenderSprites() {
	if g.spriteSizeMode == Sprite8x8Mode {
		for _, sprite := range g.sprites8x8 {
			if sprite.SpriteAttributes().X != 0x00 && sprite.SpriteAttributes().Y != 0x00 {
				g.DrawSpriteTile(sprite, sprite.GetTileID(0), 0, 0)
			}
		}
	} else {
		for _, sprite := range g.sprites8x16 {
			if sprite.SpriteAttributes().X != 0x00 && sprite.SpriteAttributes().Y != 0x00 {
				g.DrawSpriteTile(sprite, sprite.GetTileID(0), 0, 0)
				g.DrawSpriteTile(sprite, sprite.GetTileID(1), 0, 8) //8 pixels down for second tile
			}
		}
	}
}

//TODO: Sprite precedence rules
// Draws a tile for the given sprite. Only draws one tile
func (g *GPU) DrawSpriteTile(s Sprite, tileId int, screenXOffset int, screenYOffset int) {
	tile := g.FormatTileForSprite(s, tileId)
	if s.SpriteAttributes().X >= 0 && s.SpriteAttributes().Y >= 0 {
		sx, sy := s.SpriteAttributes().X-8, s.SpriteAttributes().Y-16
		for y := 0; y < 8; y++ {
			for x := 0; x < 8; x++ {
				if tile[y][x] != 0 {
					tilecolor := g.objectPalettes[s.SpriteAttributes().PaletteSelected][tile[y][x]]
					adjX, adjY := sx+x+screenXOffset, sy+y+screenYOffset
					if (adjY < DISPLAY_HEIGHT && adjY >= 0) && (adjX < DISPLAY_WIDTH && adjX >= 0) {
						if s.SpriteAttributes().SpriteHasPriority == false {
							g.screenData[adjY][adjX] = tilecolor
						} else if g.screenData[adjY][adjX] == (types.RGB{Red: 235, Green: 235, Blue: 235}) {
							g.screenData[adjY][adjX] = tilecolor
						}
					}
				}
			}
		}
	}
}

func (g *GPU) FormatTileForSprite(s Sprite, tileId int) *Tile {
	t := &g.tiledata[tileId]

	if s.SpriteAttributes().ShouldFlipHorizontally {
		var t2 *Tile = new(Tile)

		for y := 0; y < 8; y++ {
			for x := 0; x < 8; x++ {
				t2[y][x] = t[y][7-x]
			}
		}

		return t2
	}

	if s.SpriteAttributes().ShouldFlipVertically {
		var t2 *Tile = new(Tile)

		for y := 0; y < 8; y++ {
			for x := 0; x < 8; x++ {
				t2[y][x] = t[7-y][x]
			}
		}

		return t2
	}

	return t
}

//debug helpers
func (g *GPU) DumpTiles() [512][8][8]types.RGB {
	fmt.Println("Dumping", len(g.tiledata), "tiles")
	var out [512][8][8]types.RGB
	for i, tile := range g.tiledata {
		for y := 0; y < 8; y++ {
			for x := 0; x < 8; x++ {
				cr := GBColours[tile[y][x]]
				out[i][y][x] = cr
			}
		}
	}

	return out
}

func (g *GPU) Dump8x8Sprites() [40][8][8]types.RGB {
	fmt.Println("Dumping", len(g.sprites8x8), "sprites")
	var out [40][8][8]types.RGB
	for i, spr := range g.sprites8x8 {
		for y := 0; y < 8; y++ {
			for x := 0; x < 8; x++ {
				tileId := spr.GetTileID(0)
				tile := g.tiledata[tileId]
				cr := GBColours[tile[y][x]]
				out[i][y][x] = cr
			}
		}
	}
	return out
}

func (g *GPU) DumpTilemap(tileMapAddr types.Word, tileDataSigned bool) [256][256]types.RGB {
	fmt.Print("Dumping Tilemap ", tileMapAddr)
	if tileDataSigned {
		fmt.Println(" (signed)")
	} else {
		fmt.Println(" (unsigned)")
	}

	var result [256][256]types.RGB
	var tileMapAddrOffset types.Word = tileMapAddr
	var rx int = 0
	var ry int = 0

	for lineX := 0; lineX < 32; lineX++ {
		for tileY := 0; tileY < 8; tileY++ {
			for lineY := 0; lineY < 32; lineY++ {
				tileId := int(g.Read(tileMapAddrOffset + types.Word(lineY)))
				if tileDataSigned {
					if tileId < 128 {
						tileId += 256
					}
				}
				tile := g.tiledata[tileId]
				for tileX := 0; tileX < 8; tileX++ {
					cr := GBColours[tile[tileY][tileX]]
					result[rx][ry] = cr
					rx++
				}
			}
			rx = 0
			ry++
		}
		tileMapAddrOffset += types.Word(32)
	}
	return result
}

func (g *GPU) byteToPalette(b byte) Palette {
	var palette Palette
	palette[0] = GBColours[int(b&0x03)]
	palette[1] = GBColours[int((b>>2)&0x03)]
	palette[2] = GBColours[int((b>>4)&0x03)]
	palette[3] = GBColours[(int(b>>6) & 0x03)]
	return palette
}
