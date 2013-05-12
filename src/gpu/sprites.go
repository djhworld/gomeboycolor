package gpu

//Thought it would be appropriate to put the sprite type code in a different file

import (
	"fmt"
	"types"
)

type Sprite interface {
	UpdateSprite(addr types.Word, value byte)
	GetTileID(no int) int
	SpriteAttributes() *SpriteAttributes
}

//8x8 Sprites!
type Sprite8x8 struct {
	SpriteAttrs *SpriteAttributes
	TileID      int
}

func NewSprite8x8() *Sprite8x8 {
	var sprite *Sprite8x8 = new(Sprite8x8)
	sprite.SpriteAttrs = new(SpriteAttributes)
	return sprite
}

func (s *Sprite8x8) String() string {
	return fmt.Sprintf("Size: 8x8\nTile ID: %d\nAttributes: %s", s.TileID, s.SpriteAttrs)
}

func (s *Sprite8x8) SpriteAttributes() *SpriteAttributes {
	return s.SpriteAttrs
}

func (s *Sprite8x8) UpdateSprite(addr types.Word, value byte) {
	var spriteAttrId int = int(addr % 4)
	if spriteAttrId == 2 {
		s.TileID = int(value)
	} else {
		s.SpriteAttrs.Update(spriteAttrId, value)
	}
}

func (s *Sprite8x8) GetTileID(no int) int {
	if no > 0 {
		panic("8x8 sprites only consist of one tile")
	}
	return s.TileID
}

// 8x16 SPRITES!
type Sprite8x16 struct {
	SpriteAttrs *SpriteAttributes
	TileIDs     [2]int
}

func NewSprite8x16() *Sprite8x16 {
	var sprite *Sprite8x16 = new(Sprite8x16)
	sprite.SpriteAttrs = new(SpriteAttributes)
	return sprite
}

func (s *Sprite8x16) String() string {
	return fmt.Sprintf("Size: 8x16\nTile IDs: %v\nAttributes: %s", s.TileIDs, s.SpriteAttrs)
}

func (s *Sprite8x16) UpdateSprite(addr types.Word, value byte) {
	var spriteAttrId int = int(addr % 4)
	if spriteAttrId == 2 {
		s.TileIDs[0] = int(value)
		s.TileIDs[1] = int(value + 1)
	} else {
		s.SpriteAttrs.Update(spriteAttrId, value)
	}
}

func (s *Sprite8x16) GetTileID(no int) int {
	if no > 1 {
		panic("8x16 sprites only consist of two tiles")
	}
	return s.TileIDs[no]
}

func (s *Sprite8x16) SpriteAttributes() *SpriteAttributes {
	return s.SpriteAttrs
}

//Sprite attributes
type SpriteAttributes struct {
	Y                      int
	X                      int
	SpriteHasPriority      bool
	ShouldFlipVertically   bool
	ShouldFlipHorizontally bool
	PaletteSelected        int
}

func (sa *SpriteAttributes) Update(attributeId int, fromValue byte) {
	switch attributeId {
	case 0:
		sa.Y = int(fromValue)
	case 1:
		sa.X = int(fromValue)
	case 3:
		if (fromValue & 0x80) == 0x80 {
			sa.SpriteHasPriority = true
		} else {
			sa.SpriteHasPriority = false
		}

		if (fromValue & 0x40) == 0x40 {
			sa.ShouldFlipVertically = true
		} else {
			sa.ShouldFlipVertically = false
		}

		if (fromValue & 0x20) == 0x20 {
			sa.ShouldFlipHorizontally = true
		} else {
			sa.ShouldFlipHorizontally = false
		}

		if (fromValue & 0x10) == 0x10 {
			sa.PaletteSelected = 1
		} else {
			sa.PaletteSelected = 0
		}
	}
}

func (s *SpriteAttributes) String() string {
	return fmt.Sprintf("[X: %d | Y: %d | Sprite has priority? %v | Flip sprite vertically? %v | Flip sprite horizontally? %v | Palette no: %d]", s.X, s.Y, s.SpriteHasPriority, s.ShouldFlipVertically, s.ShouldFlipHorizontally, s.PaletteSelected)
}
