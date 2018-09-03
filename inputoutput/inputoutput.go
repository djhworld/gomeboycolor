package inputoutput

import (
	"log"

	"github.com/djhworld/gomeboycolor/components"
	"github.com/djhworld/gomeboycolor/constants"
	"github.com/djhworld/gomeboycolor/types"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

const PREFIX string = "IO"
const ROW_1 byte = 0x10
const ROW_2 byte = 0x20
const SCREEN_WIDTH int = 160
const SCREEN_HEIGHT int = 144

var DefaultControlScheme ControlScheme = ControlScheme{
	glfw.KeyUp,
	glfw.KeyDown,
	glfw.KeyLeft,
	glfw.KeyRight,
	glfw.KeyZ,
	glfw.KeyX,
	glfw.KeyA,
	glfw.KeyS,
}

type ControlScheme struct {
	UP     glfw.Key
	DOWN   glfw.Key
	LEFT   glfw.Key
	RIGHT  glfw.Key
	A      glfw.Key
	B      glfw.Key
	START  glfw.Key
	SELECT glfw.Key
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
func (k *KeyHandler) KeyDown(key glfw.Key) {
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
func (k *KeyHandler) KeyUp(key glfw.Key) {
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

type IO struct {
	KeyHandler          *KeyHandler
	Display             *Display
	ScreenOutputChannel chan *types.Screen
	AudioOutputChannel  chan int
	headless            bool
}

func NewIO() *IO {
	var i *IO = new(IO)
	i.KeyHandler = new(KeyHandler)
	i.Display = new(Display)
	i.ScreenOutputChannel = make(chan *types.Screen)
	i.AudioOutputChannel = make(chan int)
	i.headless = false
	return i
}

func (i *IO) init(title string, screenSize int, headless bool, onCloseHandler func()) error {
	i.headless = headless
	var err error = nil

	if !i.headless {
		err = i.Display.init(title, screenSize, onCloseHandler)
		if err != nil {
			return err
		}

		i.KeyHandler.Init(DefaultControlScheme) //TODO: allow user to define controlscheme
		i.Display.window.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
			if action == glfw.Repeat {
				i.KeyHandler.KeyDown(key)
				return
			}

			if action == glfw.Press {
				i.KeyHandler.KeyDown(key)
			} else {
				i.KeyHandler.KeyUp(key)
			}
		})
	}

	return err
}

//This will wait for updates to the display or audio and dispatch them accordingly
func (i *IO) Run(title string, screenSize int, headless bool, onCloseHandler func()) error {
	if err := i.init(title, screenSize, headless, onCloseHandler); err != nil {
		return err
	}

	for {
		select {
		case data := <-i.ScreenOutputChannel:
			if !i.headless {
				i.Display.drawFrame(data)
			}
		case data := <-i.AudioOutputChannel:
			if !i.headless {
				log.Println("Writing %d to audio!", data)
			}
		}
	}
}

type Display struct {
	Name                 string
	ScreenSizeMultiplier int
	window               *glfw.Window
}

func (s *Display) init(title string, screenSizeMultiplier int, onCloseHandler func()) error {
	var err error

	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}

	s.Name = PREFIX + "-SCREEN"

	log.Printf("%s: Initialising display", s.Name)

	s.ScreenSizeMultiplier = screenSizeMultiplier
	log.Printf("%s: Set screen size multiplier to %dx", s.Name, s.ScreenSizeMultiplier)

	glfw.WindowHint(glfw.Resizable, glfw.False)
	window, err := glfw.CreateWindow(SCREEN_WIDTH*s.ScreenSizeMultiplier, SCREEN_HEIGHT*s.ScreenSizeMultiplier, "Testing", nil, nil)
	if err != nil {
		return err
	}

	window.SetTitle(title)

	//TODO fix desktop mode
	window.SetPos((s.ScreenSizeMultiplier)/2, (s.ScreenSizeMultiplier)/2)

	window.SetCloseCallback(func(w *glfw.Window) {
		w.Destroy()
		glfw.Terminate()
		onCloseHandler()
	})

	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		return err
	}

	gl.ClearColor(0.255, 0.255, 0.255, 0)

	//resize functionTrue
	/*
		onResize := func(window *glfw.Window, w, h int) {
			gl.Viewport(0, 0, int32(w), int32(h))
			gl.MatrixMode(gl.PROJECTION)
			gl.LoadIdentity()
			gl.Ortho(0, float64(w), float64(h), 0, -1, 1)
			gl.ClearColor(0.255, 0.255, 0.255, 0)
			gl.Clear(gl.COLOR_BUFFER_BIT)
			gl.MatrixMode(gl.MODELVIEW)
			gl.LoadIdentity()
		}

		window.SetSizeCallback(onResize)
	*/

	s.window = window

	return nil

}

func (s *Display) drawFrame(screenData *types.Screen) {
	gl.Viewport(0, 0, int32(SCREEN_WIDTH*s.ScreenSizeMultiplier)*2, int32(SCREEN_HEIGHT*s.ScreenSizeMultiplier)*2)
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	gl.Ortho(0, float64(SCREEN_WIDTH*s.ScreenSizeMultiplier), float64(SCREEN_HEIGHT*s.ScreenSizeMultiplier), 0, -1, 1)
	gl.ClearColor(0.255, 0.255, 0.255, 0)
	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()

	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.Disable(gl.DEPTH_TEST)
	gl.PointSize(float32(s.ScreenSizeMultiplier) * 2.0)
	//gl.PointSize(float32(s.ScreenSizeMultiplier))
	gl.Begin(gl.POINTS)
	for y := 0; y < SCREEN_HEIGHT; y++ {
		for x := 0; x < SCREEN_WIDTH; x++ {
			var pixel types.RGB = screenData[y][x]
			gl.Color3ub(pixel.Red, pixel.Green, pixel.Blue)
			gl.Vertex2i(int32(x*s.ScreenSizeMultiplier), int32(y*s.ScreenSizeMultiplier))
		}
	}

	gl.End()
	glfw.PollEvents()
	s.window.SwapBuffers()
}
