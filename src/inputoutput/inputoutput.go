package inputoutput

import (
	"components"
	"constants"
	"github.com/go-gl/gl"
	"github.com/go-gl/glfw"
	"log"
	"types"
)

const PREFIX string = "IO"
const ROW_1 byte = 0x10
const ROW_2 byte = 0x20
const SCREEN_WIDTH int = 160
const SCREEN_HEIGHT int = 144

var DefaultControlScheme ControlScheme = ControlScheme{glfw.KeyUp, glfw.KeyDown, glfw.KeyLeft, glfw.KeyRight, 90, 88, 294, 288}

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

//Clients just need to talk to the interface to draw frames. Screen data is a pointer for performance reasons
type Screen interface {
	DrawFrame(screenData *[144][160]types.RGB)
}

type IOConfig struct {
	Title                string
	ScreenSizeMultiplier int
	OnCloseHandler       func()
	ControlScheme        ControlScheme
}

type IO struct {
	KeyHandler *KeyHandler
	Display    *Display
}

func NewIO() *IO {
	var i *IO = new(IO)
	i.KeyHandler = new(KeyHandler)
	i.Display = new(Display)
	return i
}

func (i *IO) Init(config IOConfig) error {
	var err error

	err = glfw.Init()
	if err != nil {
		return err
	}

	err = i.Display.init(config.Title, config.ScreenSizeMultiplier)
	if err != nil {
		return err
	}

	i.KeyHandler.Init(config.ControlScheme)
	glfw.SetKeyCallback(func(key, state int) {
		if state == glfw.KeyPress {
			i.KeyHandler.KeyDown(key)
		} else {
			i.KeyHandler.KeyUp(key)
		}
	})

	glfw.SetWindowCloseCallback(func() int {
		glfw.CloseWindow()
		glfw.Terminate()
		config.OnCloseHandler()
		return 0
	})

	return nil
}

type Display struct {
	Name                 string
	ScreenSizeMultiplier int
}

func (s *Display) init(title string, screenSizeMultiplier int) error {
	s.Name = PREFIX + "-SCREEN"

	log.Printf("%s: Initialising display", s.Name)
	var err error

	s.ScreenSizeMultiplier = screenSizeMultiplier
	log.Printf("%s: Set screen size multiplier to %dx", s.Name, s.ScreenSizeMultiplier)

	glfw.OpenWindowHint(glfw.WindowNoResize, 1)
	err = glfw.OpenWindow(SCREEN_WIDTH*s.ScreenSizeMultiplier, SCREEN_HEIGHT*s.ScreenSizeMultiplier, 0, 0, 0, 0, 0, 0, glfw.Windowed)
	if err != nil {
		return err
	}

	glfw.SetWindowTitle(title)

	//resize function
	onResize := func(w, h int) {
		gl.Viewport(0, 0, w, h)
		gl.MatrixMode(gl.PROJECTION)
		gl.LoadIdentity()
		gl.Ortho(0, float64(w), float64(h), 0, -1, 1)
		gl.ClearColor(0.255, 0.255, 0.255, 0)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		gl.MatrixMode(gl.MODELVIEW)
		gl.LoadIdentity()
	}

	glfw.SetWindowSizeCallback(onResize)
	glfw.SetWindowPos(100, 400)

	gl.ClearColor(0.255, 0.255, 0.255, 0)

	return nil

}

func (s *Display) DrawFrame(screenData *[144][160]types.RGB) {
	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.Disable(gl.DEPTH_TEST)
	gl.PointSize(float32(s.ScreenSizeMultiplier) + 1.0)
	gl.Begin(gl.POINTS)
	for y := 0; y < SCREEN_HEIGHT; y++ {
		for x := 0; x < SCREEN_WIDTH; x++ {
			var pixel types.RGB = screenData[y][x]
			gl.Color3ub(pixel.Red, pixel.Green, pixel.Blue)
			gl.Vertex2i(x*s.ScreenSizeMultiplier, y*s.ScreenSizeMultiplier)
		}
	}

	gl.End()
	glfw.SwapBuffers()
}
