// +build !wasm

package inputoutput

import (
	"log"

	"github.com/djhworld/gomeboycolor/types"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

var DefaultControlScheme ControlScheme = ControlScheme{
	int(glfw.KeyUp),
	int(glfw.KeyDown),
	int(glfw.KeyLeft),
	int(glfw.KeyRight),
	int(glfw.KeyZ),
	int(glfw.KeyX),
	int(glfw.KeyA),
	int(glfw.KeyS),
}

// GlfwIO is for running the emulator using GLFW.
// libglfw3 will be required on the system
type GlfwIO struct {
	*CoreIO
	glfwDisplay *glfwDisplay
}

func NewGlfwIO(frameRateLock int64, headless bool, displayFps bool) *GlfwIO {
	log.Println("Creating GLFW based IO Handler")
	glfwDisplay := new(glfwDisplay)

	frameRateReporter := func(v float32) {
		if displayFps {
			log.Printf("Average frame rate\t%.2f\tfps", v)
		}
	}

	return &GlfwIO{
		newCoreIO(frameRateLock, headless, frameRateReporter, glfwDisplay),
		glfwDisplay,
	}
}

func (i *GlfwIO) Init(title string, screenSize int, onCloseHandler func()) error {
	var err error
	i.onCloseHandler = onCloseHandler

	if !i.headless {
		err = i.glfwDisplay.init(title, screenSize)
		if err != nil {
			return err
		}

		i.glfwDisplay.window.SetCloseCallback(func(w *glfw.Window) {
			i.stopChannel <- 1
		})

		i.keyHandler.Init(DefaultControlScheme) //TODO: allow user to define controlscheme

		i.glfwDisplay.window.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
			if action == glfw.Repeat {
				i.keyHandler.KeyDown(int(key))
				return
			}

			if action == glfw.Press {
				i.keyHandler.KeyDown(int(key))
			} else {
				i.keyHandler.KeyUp(int(key))
			}
		})
	}

	return err
}

type glfwDisplay struct {
	Name                 string
	ScreenSizeMultiplier int
	window               *glfw.Window
}

func (s *glfwDisplay) init(title string, screenSizeMultiplier int) error {
	var err error

	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}

	s.Name = prefix + "-SCREEN"

	log.Printf("%s: Initialising display", s.Name)

	s.ScreenSizeMultiplier = screenSizeMultiplier
	log.Printf("%s: Set screen size multiplier to %dx", s.Name, s.ScreenSizeMultiplier)

	glfw.WindowHint(glfw.Resizable, glfw.False)
	window, err := glfw.CreateWindow(screenWidth*s.ScreenSizeMultiplier, screenHeight*s.ScreenSizeMultiplier, "Testing", nil, nil)
	if err != nil {
		return err
	}

	window.SetTitle(title)

	vidMode := glfw.GetPrimaryMonitor().GetVideoMode()

	window.SetPos(vidMode.Width/3, vidMode.Height/3)

	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		return err
	}

	gl.ClearColor(0.255, 0.255, 0.255, 0)

	s.window = window

	return nil

}

func (s *glfwDisplay) Stop() {
	log.Println("Stopping display")
	s.window.Destroy()
	glfw.Terminate()
}

func (s *glfwDisplay) DrawFrame(screenData *types.Screen) {
	fw, fh := s.window.GetFramebufferSize()
	gl.Viewport(0, 0, int32(fw), int32(fh))
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	gl.Ortho(0, float64(screenWidth*s.ScreenSizeMultiplier), float64(screenHeight*s.ScreenSizeMultiplier), 0, -1, 1)
	gl.ClearColor(0.255, 0.255, 0.255, 0)
	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()

	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.Disable(gl.DEPTH_TEST)
	gl.PointSize(float32(s.ScreenSizeMultiplier) * 2.0)
	gl.Begin(gl.POINTS)
	for y := 0; y < screenHeight; y++ {
		for x := 0; x < screenWidth; x++ {
			var pixel types.RGB = screenData[y][x]
			gl.Color3ub(pixel.Red, pixel.Green, pixel.Blue)
			gl.Vertex2i(int32(x*s.ScreenSizeMultiplier), int32(y*s.ScreenSizeMultiplier))
		}
	}

	gl.End()
	glfw.PollEvents()
	s.window.SwapBuffers()
}
