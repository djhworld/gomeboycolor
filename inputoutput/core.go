package inputoutput

import (
	"log"
	"time"

	"github.com/djhworld/gomeboycolor/types"
)

const PREFIX string = "IO"
const SCREEN_WIDTH int = 160
const SCREEN_HEIGHT int = 144

type IOHandler interface {
	Init(title string, screenSize int, onCloseHandler func()) error
	GetKeyHandler() *KeyHandler
	GetScreenOutputChannel() chan *types.Screen
	Run()
}

type CoreIO struct {
	KeyHandler          *KeyHandler
	Display             *Display
	ScreenOutputChannel chan *types.Screen
	AudioOutputChannel  chan int
	headless            bool
	frameRateLock       int64
}

func newCoreIO(frameRateLock int64, headless bool) *CoreIO {
	var i *CoreIO = new(CoreIO)
	i.KeyHandler = new(KeyHandler)
	i.Display = new(Display)
	i.ScreenOutputChannel = make(chan *types.Screen)
	i.AudioOutputChannel = make(chan int)
	i.headless = headless
	i.frameRateLock = frameRateLock
	return i
}

func (i *CoreIO) GetScreenOutputChannel() chan *types.Screen {
	return i.ScreenOutputChannel
}

func (i *CoreIO) GetKeyHandler() *KeyHandler {
	return i.KeyHandler
}

//This will wait for updates to the display or audio and dispatch them accordingly
func (i *CoreIO) Run() {
	fpsLock := time.Second / time.Duration(i.frameRateLock)
	fpsThrottler := time.Tick(fpsLock)

	for {
		select {
		case data := <-i.ScreenOutputChannel:
			if !i.headless {
				<-fpsThrottler
				i.Display.drawFrame(data)
			}
		case data := <-i.AudioOutputChannel:
			if !i.headless {
				log.Println("Writing %d to audio!", data)
			}
		}
	}
}
