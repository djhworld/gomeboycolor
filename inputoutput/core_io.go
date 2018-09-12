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

type Display interface {
	DrawFrame(*types.Screen)
	Stop()
}

// CoreIO contains all core functionality for running the IO event loop
// all sub types should extend this type
type CoreIO struct {
	KeyHandler          *KeyHandler
	ScreenOutputChannel chan *types.Screen
	AudioOutputChannel  chan int
	stopChannel         chan int
	display             Display
	headless            bool
	frameRateLock       int64
	onCloseHandler      func()
}

func newCoreIO(frameRateLock int64, headless bool, display Display) *CoreIO {
	i := new(CoreIO)
	i.KeyHandler = new(KeyHandler)
	i.ScreenOutputChannel = make(chan *types.Screen)
	i.AudioOutputChannel = make(chan int)
	i.stopChannel = make(chan int, 1)
	i.display = display
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
			i.display.Stop()
			i.onCloseHandler()
			done = true
		case data := <-i.ScreenOutputChannel:
			if !i.headless {
				<-fpsThrottler
				i.display.DrawFrame(data)
			}
		}
	}
}
