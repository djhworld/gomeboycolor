package inputoutput

import (
	"time"

	"github.com/djhworld/gomeboycolor/types"
)

const prefix string = "IO"
const screenWidth int = 160
const screenHeight int = 144

// IOHandler interface for handling all IO interations with the emulator
type IOHandler interface {
	Init(title string, screenSize int, onCloseHandler func()) error
	GetKeyHandler() *KeyHandler
	GetScreenOutputChannel() chan *types.Screen
	Run()
}

// CoreIO contains all core functionality for running the IO event loop
// all sub types should extend this type
type CoreIO struct {
	KeyHandler          *KeyHandler
	Display             *Display
	ScreenOutputChannel chan *types.Screen
	AudioOutputChannel  chan int
	stopChannel         chan int
	headless            bool
	frameRateLock       int64
	onCloseHandler      func()
}

func newCoreIO(frameRateLock int64, headless bool) *CoreIO {
	i := new(CoreIO)
	i.KeyHandler = new(KeyHandler)
	i.Display = new(Display)
	i.ScreenOutputChannel = make(chan *types.Screen)
	i.AudioOutputChannel = make(chan int)
	i.stopChannel = make(chan int, 1)
	i.headless = headless
	i.frameRateLock = frameRateLock
	i.onCloseHandler = nil
	return i
}

// GetScreenOutputChannel returns the channel to push screen
// change events to the IO event loop
func (i *CoreIO) GetScreenOutputChannel() chan *types.Screen {
	return i.ScreenOutputChannel
}

// GetKeyHandler returns the key handler component
// for managing interactions with the keyboard
func (i *CoreIO) GetKeyHandler() *KeyHandler {
	return i.KeyHandler
}

// Run runs the IO event loop
func (i *CoreIO) Run() {
	fpsLock := time.Second / time.Duration(i.frameRateLock)
	fpsThrottler := time.Tick(fpsLock)
	done := false

	for !done {
		select {
		case <-i.stopChannel:
			i.Display.destroy()
			i.onCloseHandler()
			done = true
		case data := <-i.ScreenOutputChannel:
			if !i.headless {
				<-fpsThrottler
				i.Display.drawFrame(data)
			}
		}
	}
}
