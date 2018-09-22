package main

import (
	"log"
	"time"

	"github.com/djhworld/gomeboycolor/inputoutput"
	"github.com/djhworld/gomeboycolor/types"
	"github.com/gdamore/tcell"
)

// Setup the gameboy controls -> key code mappings
var DummyControlScheme inputoutput.ControlScheme = inputoutput.ControlScheme{
	UP:     int(tcell.KeyUp),
	DOWN:   int(tcell.KeyDown),
	LEFT:   int(tcell.KeyLeft),
	RIGHT:  int(tcell.KeyRight),
	A:      122,
	B:      120,
	START:  97,
	SELECT: 115,
}

// TerminalIO is simple IO handler that outputs screen data to a terminal
type TerminalIO struct {
	*inputoutput.CoreIO
	terminalDisplay *terminalDisplay
}

func NewTerminalIO(frameRateLock int64, headless bool, displayFps bool) *TerminalIO {
	log.Println("Creating TERMINAL based IO Handler")

	frameRateReporter := func(v float32) {
		if displayFps {
			//log.Printf("Average frame rate\t%.2f\tfps", v)
		}
	}

	terminalDisplay := new(terminalDisplay)

	return &TerminalIO{
		inputoutput.NewCoreIO(frameRateLock, headless, frameRateReporter, terminalDisplay),
		terminalDisplay,
	}
}

func (i *TerminalIO) Init(title string, screenSize int, onCloseHandler func()) error {
	i.OnCloseHandler = onCloseHandler

	// 1. Initialise key handler with control scheme
	i.KeyHandler.Init(DummyControlScheme)

	// 2. Initialise tcell screen
	screen, err := tcell.NewScreen()
	if err != nil {
		return err
	}

	if err = screen.Init(); err != nil {
		return err
	}

	screen.Clear()

	// 3. Initialise display
	err = i.terminalDisplay.init(screen)
	if err != nil {
		return err
	}

	// 4. Setup key handler
	go func() {
	loop:
		for {
			ev := screen.PollEvent()
			switch ev := ev.(type) {
			case *tcell.EventKey:
				switch ev.Key() {
				case tcell.KeyEscape:
					i.StopChannel <- 1
					break loop
				case tcell.KeyRune:
					i.KeyHandler.KeyDown(int(ev.Rune()))
					time.Sleep(50 * time.Millisecond)
					i.KeyHandler.KeyUp(int(ev.Rune()))
				default:
					i.KeyHandler.KeyDown(int(ev.Key()))
					time.Sleep(50 * time.Millisecond)
					i.KeyHandler.KeyUp(int(ev.Key()))
				}
			}
		}
	}()

	return err
}

type terminalDisplay struct {
	screen tcell.Screen
}

func (n *terminalDisplay) init(screen tcell.Screen) error {
	n.screen = screen
	return nil
}

func (n *terminalDisplay) Stop() {
	n.screen.Fini()
}

// DrawFrame receives the screen data that you can use to draw to your display
func (n *terminalDisplay) DrawFrame(screenData *types.Screen) {
	st := tcell.StyleDefault
	for y := 0; y < inputoutput.SCREEN_HEIGHT; y++ {
		for x := 0; x < inputoutput.SCREEN_WIDTH; x++ {
			var pixel types.RGB = screenData[y][x]
			color := tcell.NewRGBColor(int32(pixel.Red), int32(pixel.Green), int32(pixel.Blue))
			st = st.Background(color)
			n.screen.SetCell(x, y, st, ' ')
		}
	}
	n.screen.Show()
}
