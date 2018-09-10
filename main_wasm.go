// +build wasm

package main

import (
	"fmt"
	"log"
	"strings"
	"syscall/js"
	"encoding/base64"

	"github.com/djhworld/gomeboycolor/cartridge"
	"github.com/djhworld/gomeboycolor/inputoutput"
	"github.com/djhworld/gomeboycolor/config"
	"github.com/djhworld/gomeboycolor/gbc"
	"github.com/djhworld/gomeboycolor/saves"
)

const TITLE string = "gomeboycolor"

var VERSION string

func main() {
	fmt.Printf("%s. %s\n", TITLE, VERSION)
	fmt.Println("Copyright (c) 2013. Daniel James Harper.")
	fmt.Println("http://djhworld.github.io/gomeboycolor")
	fmt.Println(strings.Repeat("*", 120))

	//Parse and validate settings file (if found)
	conf := &config.Config{
		Title:      TITLE,
		ScreenSize: 1,
		SkipBoot:   false,
		DisplayFPS: true,
		ColorMode:  true,
		Debug:      false,
		BreakOn:    "0x0000",
		DumpState:  false,
		Headless:   false,
	}

	fmt.Println("config = ", conf)

	cart, err := createCartridge("tetris")
	if err != nil {
		log.Println(err)
		return
	}

	//TODO make this configurable...?
	saveStore := saves.NewNoopStore()

	log.Println("Starting emulator")

	emulator, err := gbc.Init(cart, saveStore, conf, inputoutput.NewWebAssemblyIO(35, conf.Headless))
	if err != nil {
		log.Println(err)
		return
	}

	//set the IO controller to run indefinitely (it waits for screen updates)
	go emulator.RunIO()

	//Starts emulator code
	emulator.Run()
}

func createCartridge(romFilename string) (*cartridge.Cartridge, error) {
	romContents, err := retrieveROM(romFilename)
	if err != nil {
		return nil, err
	}

	return cartridge.NewCartridge(romFilename, romContents)
}

func retrieveROM(filename string) ([]byte, error) {
	romContents := js.Global().Get("self").Get("ROM").String()
	return base64.StdEncoding.DecodeString(romContents)
}
