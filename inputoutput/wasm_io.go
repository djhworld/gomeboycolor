// +build wasm

package inputoutput

import (
	"encoding/base64"
	"log"

	"syscall/js"

	"github.com/djhworld/gomeboycolor/types"
)

const (
	JsKeyUp    = 38
	JsKeyDown  = 40
	JsKeyLeft  = 37
	JsKeyRight = 39
	JsKeyZ     = 90
	JsKeyX     = 88
	JsKeyA     = 65
	JsKeyS     = 83
)

var DefaultControlScheme ControlScheme = ControlScheme{
	UP:     JsKeyUp,
	DOWN:   JsKeyDown,
	LEFT:   JsKeyLeft,
	RIGHT:  JsKeyRight,
	A:      JsKeyZ,
	B:      JsKeyX,
	START:  JsKeyA,
	SELECT: JsKeyS,
}

// WebAssemblyIO is for running the emulator on a WASM environment
type WebAssemblyIO struct {
	*CoreIO
}

func webWorkerMessageCallback(args []js.Value) {
	return
}

func NewWebAssemblyIO(frameRateLock int64, headless bool) *WebAssemblyIO {
	log.Println("Creating WebAssembly based IO Handler")

	return &WebAssemblyIO{
		newCoreIO(frameRateLock, headless),
	}
}

func (i *WebAssemblyIO) Init(title string, screenSize int, onCloseHandler func()) error {
	var err error = nil

	var messageCB js.Callback
	messageCB = js.NewCallback(func(args []js.Value) {
		input := args[0].Get("data")
		switch input.Index(0).Int() {
		case 0:
			i.KeyHandler.KeyUp(input.Index(1).Int())
		case 1:
			i.KeyHandler.KeyDown(input.Index(1).Int())
		case 9:
			log.Println("Resetting")
			messageCB.Release()
			onCloseHandler()
		}
	})

	if !i.headless {
		err = i.Display.init(title, screenSize, onCloseHandler)
		if err != nil {
			return err
		}

	}

	i.KeyHandler.Init(DefaultControlScheme) //TODO: allow user to define controlscheme

	self := js.Global().Get("self")
	self.Call("addEventListener", "message", messageCB)

	return err
}

type Display struct {
	Name                 string
	ScreenSizeMultiplier int
	Frame                []byte
}

//TODO on close handler?
func (s *Display) init(title string, screenSizeMultiplier int, onCloseHandler func()) error {
	s.Name = prefix + "-SCREEN"

	log.Printf("%s: Initialising display", s.Name)

	s.ScreenSizeMultiplier = screenSizeMultiplier
	log.Printf("%s: Set screen size multiplier to %dx", s.Name, s.ScreenSizeMultiplier)

	s.Frame = make([]byte, screenWidth*screenHeight*4)

	return nil

}

func (s *Display) drawFrame(screenData *types.Screen) {
	i := 0

	for y := 0; y < screenHeight; y++ {
		for x := 0; x < screenWidth; x++ {
			var pixel types.RGB = screenData[y][x]
			s.Frame[i] = pixel.Red
			s.Frame[i+1] = pixel.Green
			s.Frame[i+2] = pixel.Blue
			s.Frame[i+3] = 255
			i += 4
		}
	}

	js.Global().Call("postMessage", base64.StdEncoding.EncodeToString(s.Frame))
}
