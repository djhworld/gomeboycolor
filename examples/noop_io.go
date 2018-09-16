package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/djhworld/gomeboycolor/inputoutput"
	"github.com/djhworld/gomeboycolor/types"
)

// Setup the gameboy controls -> key code mappings
var DummyControlScheme inputoutput.ControlScheme = inputoutput.ControlScheme{
	UP:     1,
	DOWN:   2,
	LEFT:   3,
	RIGHT:  4,
	A:      5,
	B:      6,
	START:  7,
	SELECT: 8,
}

// NoopIO is simple IO handler that doesn't output to any device
type NoopIO struct {
	*inputoutput.CoreIO
	noopDisplay *noopDisplay
}

func NewNoopIO(frameRateLock int64, headless bool, displayFps bool) *NoopIO {
	log.Println("Creating NOOP based IO Handler")
	noopDisplay := new(noopDisplay)

	frameRateReporter := func(v float32) {
		if displayFps {
			log.Printf("Average frame rate\t%.2f\tfps", v)
		}
	}

	return &NoopIO{
		inputoutput.NewCoreIO(frameRateLock, headless, frameRateReporter, noopDisplay),
		noopDisplay,
	}
}

func (i *NoopIO) Init(title string, screenSize int, onCloseHandler func()) error {
	var err error
	i.OnCloseHandler = onCloseHandler

	// 1. Initialise Display
	err = i.noopDisplay.init(title, screenSize)
	if err != nil {
		return err
	}

	// 2. Handle cleanup on exit (e.g. ctrl+C)
	i.setupExitHandler()

	// 3. Initialise key handler with control scheme
	i.KeyHandler.Init(DummyControlScheme)

	// 4. Handle keyboard updates
	/*
		keyboardCallback := func(action string, key int) {
			if action == "keydown" {
				i.KeyHandler.KeyDown(key)
			} else {
				i.KeyHandler.KeyUp(key)
			}
		}
	*/

	return err
}

func (i *NoopIO) setupExitHandler() {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Println("Stopping")
		i.StopChannel <- 1
	}()
}

type noopDisplay struct {
}

func (s *noopDisplay) init(title string, screenSizeMultiplier int) error {
	return nil
}

func (s *noopDisplay) Stop() {
	// no-op
}

// DrawFrame receives the screen data that you can use to draw to your display
func (s *noopDisplay) DrawFrame(screenData *types.Screen) {
	for y := 0; y < inputoutput.SCREEN_HEIGHT; y++ {
		for x := 0; x < inputoutput.SCREEN_WIDTH; x++ {
			//var pixel types.RGB = screenData[y][x]
			//draw pixel to display
		}
	}
}
