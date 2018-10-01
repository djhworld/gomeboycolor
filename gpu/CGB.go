package gpu

import (
	"fmt"

	"github.com/djhworld/gomeboycolor/types"
)

const OBJ_PRIORITY = 1
const BG_PRIORITY = 2

//Colour GB graphics register addresses
const (
	CGB_VRAM_BANK_SELECT        types.Word = 0xFF4F
	CGB_BGP_WRITESPEC_REGISTER             = 0xFF68
	CGB_BGP_WRITEDATA_REGISTER             = 0xFF69
	CGB_OBJP_WRITESPEC_REGISTER            = 0xFF6A
	CGB_OBJP_WRITEDATA_REGISTER            = 0xFF6B
)

//Represents the attribute data for a background tile
type CGBBackgroundTileAttrs struct {
	HasPriority      bool
	FlipHorizontally bool
	FlipVertically   bool
	BankNo           int
	PaletteNo        int
}

func newCGBBackgroundTileAttrs(attributeData byte) *CGBBackgroundTileAttrs {
	var cbc *CGBBackgroundTileAttrs = new(CGBBackgroundTileAttrs)
	cbc.HasPriority = (attributeData & 0x80) == 0x80
	cbc.FlipVertically = (attributeData & 0x40) == 0x40
	cbc.FlipHorizontally = (attributeData & 0x20) == 0x20
	cbc.BankNo = int((attributeData & 0x08) >> 3)
	cbc.PaletteNo = int(attributeData & 0x07)
	return cbc
}

func (cattr *CGBBackgroundTileAttrs) String() string {
	return fmt.Sprintf("%#v", cattr)
}

//Represents a color
type CGBColor types.Word

func (c CGBColor) ToRGB() types.RGB {
	return types.RGB{
		Red:   byte(c&0x001F) * 8,
		Green: byte(c&0x03E0>>5) * 8,
		Blue:  byte(c&0x7C00>>10) * 8}
}

func (c CGBColor) High() byte {
	return byte((c & 0xFF00) >> 8)
}

func (c CGBColor) Low() byte {
	return byte(c & 0x00FF)
}

//Represents a color palette (4 colors)
type CGBPalette [4]CGBColor

func (cp *CGBPalette) UpdateHigh(colorNo int, value byte) {
	cp[colorNo] = (cp[colorNo] & 0x00FF) | CGBColor(value)<<8
}

func (cp *CGBPalette) UpdateLow(colorNo int, value byte) {
	cp[colorNo] = (cp[colorNo] & 0xFF00) | CGBColor(value)
}

//Represents the write specification register for a color palette
type CGBPaletteSpecRegister struct {
	Value           byte
	PalleteNo       int
	PalleteDataNo   int
	High            bool
	IncrementOnNext bool
}

func (psr *CGBPaletteSpecRegister) Update(value byte) {
	psr.Value = value
	psr.PalleteNo = int((value & 0x38) >> 3)
	psr.PalleteDataNo = int((value & 0x06) >> 1)
	psr.High = (value & 0x01) == 0x01
	psr.IncrementOnNext = (value & 0x80) == 0x80
}

func (psr *CGBPaletteSpecRegister) Increment() {
	psr.Update(psr.Value + 1)
}

/*
Calculate the priority of BG vs OBJ using rules defined in
the Gameboy Programming Manual - Page 67
*/
func calculateObjToBackgroundPriority(backgroundPriorityFlag, objPriorityFlag bool, bgDotData, objDotData int) int {
	if backgroundPriorityFlag { //BG has highest priority
		switch {
		case objDotData == 0x00 && bgDotData == 0x00: //bg gets priority
			return BG_PRIORITY
		case objDotData == 0x00 && bgDotData != 0x00: //bg gets priority
			return BG_PRIORITY
		case objDotData != 0x00 && bgDotData == 0x00: //object gets priority
			return OBJ_PRIORITY
		case objDotData != 0x00 && bgDotData != 0x00: //bg gets priority
			return BG_PRIORITY
		}
	} else { // bg depends on what obj priority is
		if objPriorityFlag { //OBJ has highest priority
			switch {
			case objDotData == 0x00 && bgDotData == 0x00: //bg gets priority
				return BG_PRIORITY
			case objDotData == 0x00 && bgDotData != 0x00: //bg gets priority
				return BG_PRIORITY
			case objDotData != 0x00 && bgDotData == 0x00: //object gets priority
				return OBJ_PRIORITY
			case objDotData != 0x00 && bgDotData != 0x00: //object gets priority
				return OBJ_PRIORITY
			}
		} else {
			switch { //give priority to BG
			case objDotData == 0x00 && bgDotData == 0x00: //bg gets priority
				return BG_PRIORITY
			case objDotData == 0x00 && bgDotData != 0x00: //bg gets priority
				return BG_PRIORITY
			case objDotData != 0x00 && bgDotData == 0x00: //object gets priority
				return OBJ_PRIORITY
			case objDotData != 0x00 && bgDotData != 0x00: //bg gets priority
				return BG_PRIORITY
			}
		}

	}

	return OBJ_PRIORITY
}

var CGB_BACKGROUND_TILE_ATTRS []*CGBBackgroundTileAttrs = []*CGBBackgroundTileAttrs{
	newCGBBackgroundTileAttrs(byte(0)),
	newCGBBackgroundTileAttrs(byte(1)),
	newCGBBackgroundTileAttrs(byte(2)),
	newCGBBackgroundTileAttrs(byte(3)),
	newCGBBackgroundTileAttrs(byte(4)),
	newCGBBackgroundTileAttrs(byte(5)),
	newCGBBackgroundTileAttrs(byte(6)),
	newCGBBackgroundTileAttrs(byte(7)),
	newCGBBackgroundTileAttrs(byte(8)),
	newCGBBackgroundTileAttrs(byte(9)),
	newCGBBackgroundTileAttrs(byte(10)),
	newCGBBackgroundTileAttrs(byte(11)),
	newCGBBackgroundTileAttrs(byte(12)),
	newCGBBackgroundTileAttrs(byte(13)),
	newCGBBackgroundTileAttrs(byte(14)),
	newCGBBackgroundTileAttrs(byte(15)),
	newCGBBackgroundTileAttrs(byte(16)),
	newCGBBackgroundTileAttrs(byte(17)),
	newCGBBackgroundTileAttrs(byte(18)),
	newCGBBackgroundTileAttrs(byte(19)),
	newCGBBackgroundTileAttrs(byte(20)),
	newCGBBackgroundTileAttrs(byte(21)),
	newCGBBackgroundTileAttrs(byte(22)),
	newCGBBackgroundTileAttrs(byte(23)),
	newCGBBackgroundTileAttrs(byte(24)),
	newCGBBackgroundTileAttrs(byte(25)),
	newCGBBackgroundTileAttrs(byte(26)),
	newCGBBackgroundTileAttrs(byte(27)),
	newCGBBackgroundTileAttrs(byte(28)),
	newCGBBackgroundTileAttrs(byte(29)),
	newCGBBackgroundTileAttrs(byte(30)),
	newCGBBackgroundTileAttrs(byte(31)),
	newCGBBackgroundTileAttrs(byte(32)),
	newCGBBackgroundTileAttrs(byte(33)),
	newCGBBackgroundTileAttrs(byte(34)),
	newCGBBackgroundTileAttrs(byte(35)),
	newCGBBackgroundTileAttrs(byte(36)),
	newCGBBackgroundTileAttrs(byte(37)),
	newCGBBackgroundTileAttrs(byte(38)),
	newCGBBackgroundTileAttrs(byte(39)),
	newCGBBackgroundTileAttrs(byte(40)),
	newCGBBackgroundTileAttrs(byte(41)),
	newCGBBackgroundTileAttrs(byte(42)),
	newCGBBackgroundTileAttrs(byte(43)),
	newCGBBackgroundTileAttrs(byte(44)),
	newCGBBackgroundTileAttrs(byte(45)),
	newCGBBackgroundTileAttrs(byte(46)),
	newCGBBackgroundTileAttrs(byte(47)),
	newCGBBackgroundTileAttrs(byte(48)),
	newCGBBackgroundTileAttrs(byte(49)),
	newCGBBackgroundTileAttrs(byte(50)),
	newCGBBackgroundTileAttrs(byte(51)),
	newCGBBackgroundTileAttrs(byte(52)),
	newCGBBackgroundTileAttrs(byte(53)),
	newCGBBackgroundTileAttrs(byte(54)),
	newCGBBackgroundTileAttrs(byte(55)),
	newCGBBackgroundTileAttrs(byte(56)),
	newCGBBackgroundTileAttrs(byte(57)),
	newCGBBackgroundTileAttrs(byte(58)),
	newCGBBackgroundTileAttrs(byte(59)),
	newCGBBackgroundTileAttrs(byte(60)),
	newCGBBackgroundTileAttrs(byte(61)),
	newCGBBackgroundTileAttrs(byte(62)),
	newCGBBackgroundTileAttrs(byte(63)),
	newCGBBackgroundTileAttrs(byte(64)),
	newCGBBackgroundTileAttrs(byte(65)),
	newCGBBackgroundTileAttrs(byte(66)),
	newCGBBackgroundTileAttrs(byte(67)),
	newCGBBackgroundTileAttrs(byte(68)),
	newCGBBackgroundTileAttrs(byte(69)),
	newCGBBackgroundTileAttrs(byte(70)),
	newCGBBackgroundTileAttrs(byte(71)),
	newCGBBackgroundTileAttrs(byte(72)),
	newCGBBackgroundTileAttrs(byte(73)),
	newCGBBackgroundTileAttrs(byte(74)),
	newCGBBackgroundTileAttrs(byte(75)),
	newCGBBackgroundTileAttrs(byte(76)),
	newCGBBackgroundTileAttrs(byte(77)),
	newCGBBackgroundTileAttrs(byte(78)),
	newCGBBackgroundTileAttrs(byte(79)),
	newCGBBackgroundTileAttrs(byte(80)),
	newCGBBackgroundTileAttrs(byte(81)),
	newCGBBackgroundTileAttrs(byte(82)),
	newCGBBackgroundTileAttrs(byte(83)),
	newCGBBackgroundTileAttrs(byte(84)),
	newCGBBackgroundTileAttrs(byte(85)),
	newCGBBackgroundTileAttrs(byte(86)),
	newCGBBackgroundTileAttrs(byte(87)),
	newCGBBackgroundTileAttrs(byte(88)),
	newCGBBackgroundTileAttrs(byte(89)),
	newCGBBackgroundTileAttrs(byte(90)),
	newCGBBackgroundTileAttrs(byte(91)),
	newCGBBackgroundTileAttrs(byte(92)),
	newCGBBackgroundTileAttrs(byte(93)),
	newCGBBackgroundTileAttrs(byte(94)),
	newCGBBackgroundTileAttrs(byte(95)),
	newCGBBackgroundTileAttrs(byte(96)),
	newCGBBackgroundTileAttrs(byte(97)),
	newCGBBackgroundTileAttrs(byte(98)),
	newCGBBackgroundTileAttrs(byte(99)),
	newCGBBackgroundTileAttrs(byte(100)),
	newCGBBackgroundTileAttrs(byte(101)),
	newCGBBackgroundTileAttrs(byte(102)),
	newCGBBackgroundTileAttrs(byte(103)),
	newCGBBackgroundTileAttrs(byte(104)),
	newCGBBackgroundTileAttrs(byte(105)),
	newCGBBackgroundTileAttrs(byte(106)),
	newCGBBackgroundTileAttrs(byte(107)),
	newCGBBackgroundTileAttrs(byte(108)),
	newCGBBackgroundTileAttrs(byte(109)),
	newCGBBackgroundTileAttrs(byte(110)),
	newCGBBackgroundTileAttrs(byte(111)),
	newCGBBackgroundTileAttrs(byte(112)),
	newCGBBackgroundTileAttrs(byte(113)),
	newCGBBackgroundTileAttrs(byte(114)),
	newCGBBackgroundTileAttrs(byte(115)),
	newCGBBackgroundTileAttrs(byte(116)),
	newCGBBackgroundTileAttrs(byte(117)),
	newCGBBackgroundTileAttrs(byte(118)),
	newCGBBackgroundTileAttrs(byte(119)),
	newCGBBackgroundTileAttrs(byte(120)),
	newCGBBackgroundTileAttrs(byte(121)),
	newCGBBackgroundTileAttrs(byte(122)),
	newCGBBackgroundTileAttrs(byte(123)),
	newCGBBackgroundTileAttrs(byte(124)),
	newCGBBackgroundTileAttrs(byte(125)),
	newCGBBackgroundTileAttrs(byte(126)),
	newCGBBackgroundTileAttrs(byte(127)),
	newCGBBackgroundTileAttrs(byte(128)),
	newCGBBackgroundTileAttrs(byte(129)),
	newCGBBackgroundTileAttrs(byte(130)),
	newCGBBackgroundTileAttrs(byte(131)),
	newCGBBackgroundTileAttrs(byte(132)),
	newCGBBackgroundTileAttrs(byte(133)),
	newCGBBackgroundTileAttrs(byte(134)),
	newCGBBackgroundTileAttrs(byte(135)),
	newCGBBackgroundTileAttrs(byte(136)),
	newCGBBackgroundTileAttrs(byte(137)),
	newCGBBackgroundTileAttrs(byte(138)),
	newCGBBackgroundTileAttrs(byte(139)),
	newCGBBackgroundTileAttrs(byte(140)),
	newCGBBackgroundTileAttrs(byte(141)),
	newCGBBackgroundTileAttrs(byte(142)),
	newCGBBackgroundTileAttrs(byte(143)),
	newCGBBackgroundTileAttrs(byte(144)),
	newCGBBackgroundTileAttrs(byte(145)),
	newCGBBackgroundTileAttrs(byte(146)),
	newCGBBackgroundTileAttrs(byte(147)),
	newCGBBackgroundTileAttrs(byte(148)),
	newCGBBackgroundTileAttrs(byte(149)),
	newCGBBackgroundTileAttrs(byte(150)),
	newCGBBackgroundTileAttrs(byte(151)),
	newCGBBackgroundTileAttrs(byte(152)),
	newCGBBackgroundTileAttrs(byte(153)),
	newCGBBackgroundTileAttrs(byte(154)),
	newCGBBackgroundTileAttrs(byte(155)),
	newCGBBackgroundTileAttrs(byte(156)),
	newCGBBackgroundTileAttrs(byte(157)),
	newCGBBackgroundTileAttrs(byte(158)),
	newCGBBackgroundTileAttrs(byte(159)),
	newCGBBackgroundTileAttrs(byte(160)),
	newCGBBackgroundTileAttrs(byte(161)),
	newCGBBackgroundTileAttrs(byte(162)),
	newCGBBackgroundTileAttrs(byte(163)),
	newCGBBackgroundTileAttrs(byte(164)),
	newCGBBackgroundTileAttrs(byte(165)),
	newCGBBackgroundTileAttrs(byte(166)),
	newCGBBackgroundTileAttrs(byte(167)),
	newCGBBackgroundTileAttrs(byte(168)),
	newCGBBackgroundTileAttrs(byte(169)),
	newCGBBackgroundTileAttrs(byte(170)),
	newCGBBackgroundTileAttrs(byte(171)),
	newCGBBackgroundTileAttrs(byte(172)),
	newCGBBackgroundTileAttrs(byte(173)),
	newCGBBackgroundTileAttrs(byte(174)),
	newCGBBackgroundTileAttrs(byte(175)),
	newCGBBackgroundTileAttrs(byte(176)),
	newCGBBackgroundTileAttrs(byte(177)),
	newCGBBackgroundTileAttrs(byte(178)),
	newCGBBackgroundTileAttrs(byte(179)),
	newCGBBackgroundTileAttrs(byte(180)),
	newCGBBackgroundTileAttrs(byte(181)),
	newCGBBackgroundTileAttrs(byte(182)),
	newCGBBackgroundTileAttrs(byte(183)),
	newCGBBackgroundTileAttrs(byte(184)),
	newCGBBackgroundTileAttrs(byte(185)),
	newCGBBackgroundTileAttrs(byte(186)),
	newCGBBackgroundTileAttrs(byte(187)),
	newCGBBackgroundTileAttrs(byte(188)),
	newCGBBackgroundTileAttrs(byte(189)),
	newCGBBackgroundTileAttrs(byte(190)),
	newCGBBackgroundTileAttrs(byte(191)),
	newCGBBackgroundTileAttrs(byte(192)),
	newCGBBackgroundTileAttrs(byte(193)),
	newCGBBackgroundTileAttrs(byte(194)),
	newCGBBackgroundTileAttrs(byte(195)),
	newCGBBackgroundTileAttrs(byte(196)),
	newCGBBackgroundTileAttrs(byte(197)),
	newCGBBackgroundTileAttrs(byte(198)),
	newCGBBackgroundTileAttrs(byte(199)),
	newCGBBackgroundTileAttrs(byte(200)),
	newCGBBackgroundTileAttrs(byte(201)),
	newCGBBackgroundTileAttrs(byte(202)),
	newCGBBackgroundTileAttrs(byte(203)),
	newCGBBackgroundTileAttrs(byte(204)),
	newCGBBackgroundTileAttrs(byte(205)),
	newCGBBackgroundTileAttrs(byte(206)),
	newCGBBackgroundTileAttrs(byte(207)),
	newCGBBackgroundTileAttrs(byte(208)),
	newCGBBackgroundTileAttrs(byte(209)),
	newCGBBackgroundTileAttrs(byte(210)),
	newCGBBackgroundTileAttrs(byte(211)),
	newCGBBackgroundTileAttrs(byte(212)),
	newCGBBackgroundTileAttrs(byte(213)),
	newCGBBackgroundTileAttrs(byte(214)),
	newCGBBackgroundTileAttrs(byte(215)),
	newCGBBackgroundTileAttrs(byte(216)),
	newCGBBackgroundTileAttrs(byte(217)),
	newCGBBackgroundTileAttrs(byte(218)),
	newCGBBackgroundTileAttrs(byte(219)),
	newCGBBackgroundTileAttrs(byte(220)),
	newCGBBackgroundTileAttrs(byte(221)),
	newCGBBackgroundTileAttrs(byte(222)),
	newCGBBackgroundTileAttrs(byte(223)),
	newCGBBackgroundTileAttrs(byte(224)),
	newCGBBackgroundTileAttrs(byte(225)),
	newCGBBackgroundTileAttrs(byte(226)),
	newCGBBackgroundTileAttrs(byte(227)),
	newCGBBackgroundTileAttrs(byte(228)),
	newCGBBackgroundTileAttrs(byte(229)),
	newCGBBackgroundTileAttrs(byte(230)),
	newCGBBackgroundTileAttrs(byte(231)),
	newCGBBackgroundTileAttrs(byte(232)),
	newCGBBackgroundTileAttrs(byte(233)),
	newCGBBackgroundTileAttrs(byte(234)),
	newCGBBackgroundTileAttrs(byte(235)),
	newCGBBackgroundTileAttrs(byte(236)),
	newCGBBackgroundTileAttrs(byte(237)),
	newCGBBackgroundTileAttrs(byte(238)),
	newCGBBackgroundTileAttrs(byte(239)),
	newCGBBackgroundTileAttrs(byte(240)),
	newCGBBackgroundTileAttrs(byte(241)),
	newCGBBackgroundTileAttrs(byte(242)),
	newCGBBackgroundTileAttrs(byte(243)),
	newCGBBackgroundTileAttrs(byte(244)),
	newCGBBackgroundTileAttrs(byte(245)),
	newCGBBackgroundTileAttrs(byte(246)),
	newCGBBackgroundTileAttrs(byte(247)),
	newCGBBackgroundTileAttrs(byte(248)),
	newCGBBackgroundTileAttrs(byte(249)),
	newCGBBackgroundTileAttrs(byte(250)),
	newCGBBackgroundTileAttrs(byte(251)),
	newCGBBackgroundTileAttrs(byte(252)),
	newCGBBackgroundTileAttrs(byte(253)),
	newCGBBackgroundTileAttrs(byte(254)),
	newCGBBackgroundTileAttrs(byte(255)),
}
