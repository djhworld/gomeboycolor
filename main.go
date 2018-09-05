// +build !wasm

package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/djhworld/gomeboycolor/cartridge"
	"github.com/djhworld/gomeboycolor/config"
	"github.com/djhworld/gomeboycolor/gbc"
	"github.com/djhworld/gomeboycolor/saves"
)

const TITLE string = "gomeboycolor"

var VERSION string

const (
	SKIP_BOOT_FLAG   string = "skipboot"
	SCREEN_SIZE_FLAG        = "size"
	SHOW_FPS_FLAG           = "showfps"
	TITLE_FLAG              = "title"
	DUMP_FLAG               = "dump"
	DEBUGGER_ON_FLAG        = "debug"
	BREAK_WHEN_FLAG         = "b"
	COLOR_MODE_FLAG         = "color"
	HELP_FLAG               = "help"
	HEADLESS_FLAG           = "headless"
)

var title *string = flag.String(TITLE_FLAG, TITLE, "Title to use")
var showFps *bool = flag.Bool(SHOW_FPS_FLAG, false, "Calculate and display frames per second")
var screenSizeMultiplier *int = flag.Int(SCREEN_SIZE_FLAG, 1, "Screen size multiplier")
var skipBoot *bool = flag.Bool(SKIP_BOOT_FLAG, false, "Skip boot sequence")
var colorMode *bool = flag.Bool(COLOR_MODE_FLAG, true, "Emulates Gameboy Color Hardware")
var help *bool = flag.Bool(HELP_FLAG, false, "Show this help message")
var headless *bool = flag.Bool(HEADLESS_FLAG, false, "Run emulator without output")

//debug stuff...
var dumpState *bool = flag.Bool(DUMP_FLAG, false, "Print state of machine after each cycle (WARNING - WILL RUN SLOW)")
var debug *bool = flag.Bool(DEBUGGER_ON_FLAG, false, "Enable debugger")
var breakOn *string = flag.String(BREAK_WHEN_FLAG, "0x0000", "Break into debugger when PC equals a given value between 0x0000 and 0xFFFF")

func PrintHelp() {
	fmt.Println("\nUsage: -\n")
	fmt.Println("To launch the emulator, simply run and pass it the location of your ROM file, e.g. ")
	fmt.Println("\n\tgomeboycolor location/of/romfile.gbc\n")
	fmt.Println("Flags: -\n")
	fmt.Println("	-help			->	Show this help message")
	fmt.Println("	-skipboot		->	Disables the boot sequence and will boot you straight into the ROM you have provided. Defaults to false")
	fmt.Println("	-color			->	Turns color GB features on. Defaults to true")
	fmt.Println("	-showfps		->	Prints average frames per second to the console. Defaults to false")
	fmt.Println("	-dump			-> 	Dump CPU state after every cycle. Will be very SLOW and resource intensive. Defaults to false")
	fmt.Println("	-size=(1-6)		->	Set screen size. Defaults to 1.")
	fmt.Println("	-headless		->	Runs emulator without output")
	fmt.Println("	-title=(title)		->	Change window title. Defaults to 'gomeboycolor'.")
	fmt.Println("\nYou can pass an option argument to the boolean flags if you want to enable that particular option. e.g. to disable the boot screen you would do the following")
	fmt.Println("\n\tgomeboycolor -skipboot=false location/of/romfile.gbc\n")
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Printf("%s. %s\n", TITLE, VERSION)
	fmt.Println("Copyright (c) 2013. Daniel James Harper.")
	fmt.Println("http://djhworld.github.io/gomeboycolor")
	fmt.Println(strings.Repeat("*", 120))

	flag.Usage = PrintHelp

	flag.Parse()

	if *help {
		PrintHelp()
		os.Exit(1)
	}

	if flag.NArg() != 1 {
		log.Fatalf("Please specify the location of a ROM to boot")
		return
	}

	//Parse and validate settings file (if found)
	conf := &config.Config{
		Title:      TITLE,
		ScreenSize: *screenSizeMultiplier,
		SkipBoot:   *skipBoot,
		DisplayFPS: *showFps,
		ColorMode:  *colorMode,
		Debug:      *debug,
		BreakOn:    *breakOn,
		DumpState:  *dumpState,
		Headless:   *headless,
	}
	fmt.Println(conf)

	cart, err := createCartridge(flag.Arg(0))
	if err != nil {
		log.Println(err)
		return
	}

	//TODO make this configurable...?
	saveStore := saves.NewFileSystemStore("")

	log.Println("Starting emulator")

	emulator, err := gbc.Init(cart, saveStore, conf)
	if err != nil {
		log.Println(err)
		return
	}

	//Starts emulator code in a goroutine
	go emulator.Run()

	//lock the OS thread here
	runtime.LockOSThread()

	//set the IO controller to run indefinitely (it waits for screen updates)
	emulator.RunIO()
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

	var size int64 = stats.Size()
	bytes := make([]byte, size)

	bufr := bufio.NewReader(file)
	_, err = bufr.Read(bytes)

	return bytes, err
}
