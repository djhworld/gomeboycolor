package inputoutput

import (
	"time"

	"github.com/djhworld/gomeboycolor/metric"
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
	GetAvgFrameRate() float32
	Run()
}

type Display interface {
	DrawFrame(*types.Screen)
	Stop()
}

// CoreIO contains all core functionality for running the IO event loop
// all sub types should extend this type
type CoreIO struct {
	keyHandler          *KeyHandler
	audioOutputChannel  chan int
	screenOutputChannel chan *types.Screen
	stopChannel         chan int
	display             Display
	headless            bool
	frameRateLock       int64
	onCloseHandler      func()
	frameRateCounter    *metric.FPSCounter
	frameRateReporter   func(float32)
}

func newCoreIO(frameRateLock int64, headless bool, frameRateReporter func(float32), display Display) *CoreIO {
	i := new(CoreIO)
	i.keyHandler = new(KeyHandler)
	i.audioOutputChannel = make(chan int)
	i.screenOutputChannel = make(chan *types.Screen)
	i.stopChannel = make(chan int, 1)
	i.display = display
	i.headless = headless
	i.frameRateLock = frameRateLock
	i.onCloseHandler = nil
	i.frameRateCounter = metric.NewFPSCounter()
	i.frameRateReporter = frameRateReporter
	return i
}

// GetScreenOutputChannel returns the channel to push screen
// change events to the IO event loop
func (i *CoreIO) GetScreenOutputChannel() chan *types.Screen {
	return i.screenOutputChannel
}

// GetKeyHandler returns the key handler component
// for managing interactions with the keyboard
func (i *CoreIO) GetKeyHandler() *KeyHandler {
	return i.keyHandler
}

func (i *CoreIO) GetAvgFrameRate() float32 {
	return i.frameRateCounter.Avg()
}

// Run runs the IO event loop
func (i *CoreIO) Run() {
	fpsLock := time.Second / time.Duration(i.frameRateLock)
	fpsThrottler := time.Tick(fpsLock)
	frameRateCountTicker := time.Tick(1 * time.Second)
	frameCount := 0
	isRunning := true

	for isRunning {
		select {
		case data := <-i.screenOutputChannel:
			<-fpsThrottler
			i.display.DrawFrame(data)
			frameCount++
		case <-i.stopChannel:
			i.display.Stop()
			i.onCloseHandler()
			isRunning = false
		case <-frameRateCountTicker:
			i.frameRateCounter.Add(frameCount)
			i.frameRateReporter(i.frameRateCounter.Avg())
			frameCount = 0
		}
	}
}
