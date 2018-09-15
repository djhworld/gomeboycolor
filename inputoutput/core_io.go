package inputoutput

import (
	"time"

	"github.com/djhworld/gomeboycolor/metric"
	"github.com/djhworld/gomeboycolor/types"
)

const PREFIX string = "IO"
const SCREEN_WIDTH int = 160
const SCREEN_HEIGHT int = 144

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
	OnCloseHandler func()
	KeyHandler     *KeyHandler
	StopChannel    chan int
	Headless       bool

	audioOutputChannel  chan int
	screenOutputChannel chan *types.Screen
	display             Display
	frameRateLock       int64
	frameRateCounter    *metric.FPSCounter
	frameRateReporter   func(float32)
}

func NewCoreIO(frameRateLock int64, headless bool, frameRateReporter func(float32), display Display) *CoreIO {
	i := new(CoreIO)
	i.KeyHandler = new(KeyHandler)
	i.StopChannel = make(chan int, 1)
	i.Headless = headless
	i.OnCloseHandler = nil

	i.screenOutputChannel = make(chan *types.Screen)
	i.audioOutputChannel = make(chan int)
	i.display = display
	i.frameRateLock = frameRateLock
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
	return i.KeyHandler
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
		case <-i.StopChannel:
			i.display.Stop()
			i.OnCloseHandler()
			isRunning = false
		case <-frameRateCountTicker:
			i.frameRateCounter.Add(frameCount)
			i.frameRateReporter(i.frameRateCounter.Avg())
			frameCount = 0
		}
	}
}
