// +build wasm

package inputoutput

import (
	"encoding/base64"
	"log"

	"syscall/js"

	"github.com/djhworld/gomeboycolor/types"
	"github.com/djhworld/gomeboycolor/webworker"
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

// WebIO is for running the emulator in a web environment
type WebIO struct {
	*CoreIO
	html5Display *html5CanvasDisplay
}

func NewWebIO(frameRateLock int64, headless bool) *WebIO {
	log.Println("Creating Web based IO Handler")
	html5Display := new(html5CanvasDisplay)

	return &WebIO{
		newCoreIO(frameRateLock, headless, html5Display),
		html5Display,
	}
}

func (i *WebIO) Init(title string, screenSize int, onCloseHandler func()) error {
	var err error = nil
	i.onCloseHandler = onCloseHandler

	var messageCB js.Callback
	messageCB = js.NewCallback(func(args []js.Value) {
		input := args[0].Get("data")
		switch input.Index(0).String() {
		case "ku":
			i.keyHandler.KeyUp(input.Index(1).Int())
		case "kd":
			i.keyHandler.KeyDown(input.Index(1).Int())
		case "stop":
			messageCB.Release()
			i.stopChannel <- 1
		}
	})

	if !i.headless {
		err = i.html5Display.init(title)
		if err != nil {
			return err
		}

		i.keyHandler.Init(DefaultControlScheme) //TODO: allow user to define controlscheme
	}

	self := js.Global().Get("self")
	self.Call("addEventListener", "message", messageCB, false)

	return err
}

type html5CanvasDisplay struct {
	Name      string
	imageData []byte
}

//TODO on close handler?
func (s *html5CanvasDisplay) init(title string) error {
	s.Name = prefix + "-SCREEN"
	log.Printf("%s: Initialising display", s.Name)

	imageDataLen := screenWidth * screenHeight * 4
	s.imageData = make([]byte, imageDataLen, imageDataLen)

	return nil
}

func (s *html5CanvasDisplay) Stop() {
	// noop
}

func (s *html5CanvasDisplay) DrawFrame(screenData *types.Screen) {
	i := 0

	for y := 0; y < screenHeight; y++ {
		for x := 0; x < screenWidth; x++ {
			var pixel types.RGB = screenData[y][x]
			s.imageData[i] = pixel.Red
			s.imageData[i+1] = pixel.Green
			s.imageData[i+2] = pixel.Blue
			s.imageData[i+3] = 255
			i += 4
		}
	}

	webworker.SendScreenUpdate(base64.StdEncoding.EncodeToString(s.imageData))
}
