package main

import (
	"bufio"
	"errors"
	"log"
	"os"
	"runtime"

	"github.com/djhworld/gomeboycolor/cartridge"
	"github.com/djhworld/gomeboycolor/config"
	"github.com/djhworld/gomeboycolor/gbc"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	if len(os.Args) != 2 {
		log.Fatalf("ERROR: %v", errors.New("Please specify the location of a ROM to boot"))
	}

	// 1. Setup configuration
	conf := &config.Config{
		Title:         "dummmy gomeboycolor",
		ScreenSize:    1,
		SkipBoot:      false,
		DisplayFPS:    true,
		ColorMode:     true,
		Debug:         false,
		BreakOn:       "0x0000",
		DumpState:     false,
		Headless:      false,
		FrameRateLock: 58,
	}

	romFile := os.Args[1]

	// 1. Initialise emulator
	emulator, err := createEmulator(romFile, conf)
	if err != nil {
		log.Fatalf("ERROR: %v", err)
	}

	// 2. Starts core emulator runtime in a goroutine
	go emulator.Run()

	// 3. Start the IO loop to run indefinitely to handle screen updates/keyboard input etc.
	emulator.RunIO()
}

func createEmulator(romFile string, conf *config.Config) (*gbc.GomeboyColor, error) {

	// 2. Load ROM file into a cartridge struct
	cart, err := createCartridge(os.Args[1])
	if err != nil {
		return nil, err
	}

	// 3. Create save store
	saveStore := NewNoopStore()

	// 4. Create IO handler
	ioHandler := NewTerminalIO(conf.FrameRateLock, conf.Headless, conf.DisplayFPS)

	// 5. Initialise emulator
	return gbc.Init(
		cart,
		saveStore,
		conf,
		ioHandler,
	)
}

func createCartridge(romFilename string) (*cartridge.Cartridge, error) {
	romContents, err := retrieveROM(romFilename)
	if err != nil {
		return nil, err
	}

	return cartridge.NewCartridge(romFilename, romContents)
}

func retrieveROM(filename string) ([]byte, error) {
	file, err := os.Open(filename)

	if err != nil {
		return nil, err
	}
	defer file.Close()

	stats, statsErr := file.Stat()
	if statsErr != nil {
		return nil, statsErr
	}

	size := stats.Size()
	bytes := make([]byte, size)

	bufr := bufio.NewReader(file)
	_, err = bufr.Read(bytes)

	return bytes, err
}
