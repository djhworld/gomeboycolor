package inputoutput

import (
	"log"

	"github.com/djhworld/gomeboycolor/components"
	"github.com/djhworld/gomeboycolor/constants"
	"github.com/djhworld/gomeboycolor/types"
)

const ROW_1 byte = 0x10
const ROW_2 byte = 0x20

type ControlScheme struct {
	UP     int
	DOWN   int
	LEFT   int
	RIGHT  int
	A      int
	B      int
	START  int
	SELECT int
}

type KeyHandler struct {
	controlScheme ControlScheme
	colSelect     byte
	rows          [2]byte
	irqHandler    components.IRQHandler
}

func (k *KeyHandler) Init(cs ControlScheme) {
	k.controlScheme = cs
	k.Reset()
}

func (k *KeyHandler) Name() string {
	return PREFIX + "-KEYB"
}

func (k *KeyHandler) Reset() {
	log.Printf("%s: Resetting", k.Name())
	k.rows[0], k.rows[1] = 0x0F, 0x0F
	k.colSelect = 0x00
}

func (k *KeyHandler) LinkIRQHandler(m components.IRQHandler) {
	k.irqHandler = m
	log.Printf("%s: Linked IRQ Handler to Keyboard Handler", k.Name())
}

func (k *KeyHandler) Read(addr types.Word) byte {
	var value byte

	switch k.colSelect {
	case ROW_1:
		value = k.rows[1]
	case ROW_2:
		value = k.rows[0]
	default:
		value = 0x00
	}

	return value
}

func (k *KeyHandler) Write(addr types.Word, value byte) {
	k.colSelect = value & 0x30
}

//released sets bit for key to 0
func (k *KeyHandler) KeyDown(key int) {
	k.irqHandler.RequestInterrupt(constants.JOYP_HILO_IRQ)
	switch key {
	case k.controlScheme.UP:
		k.rows[0] &= 0xB
	case k.controlScheme.DOWN:
		k.rows[0] &= 0x7
	case k.controlScheme.LEFT:
		k.rows[0] &= 0xD
	case k.controlScheme.RIGHT:
		k.rows[0] &= 0xE
	case k.controlScheme.A:
		k.rows[1] &= 0xE
	case k.controlScheme.B:
		k.rows[1] &= 0xD
	case k.controlScheme.START:
		k.rows[1] &= 0x7
	case k.controlScheme.SELECT:
		k.rows[1] &= 0xB
	}
}

//released sets bit for key to 1
func (k *KeyHandler) KeyUp(key int) {
	switch key {
	case k.controlScheme.UP:
		k.rows[0] |= 0x4
	case k.controlScheme.DOWN:
		k.rows[0] |= 0x8
	case k.controlScheme.LEFT:
		k.rows[0] |= 0x2
	case k.controlScheme.RIGHT:
		k.rows[0] |= 0x1
	case k.controlScheme.A:
		k.rows[1] |= 0x1
	case k.controlScheme.B:
		k.rows[1] |= 0x2
	case k.controlScheme.START:
		k.rows[1] |= 0x8
	case k.controlScheme.SELECT:
		k.rows[1] |= 0x4
	}
}
