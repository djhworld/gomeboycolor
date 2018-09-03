package main

import (
	"log"
	"runtime"

	"github.com/djhworld/gomeboycolor/cartridge"
	"github.com/djhworld/gomeboycolor/config"
	"github.com/djhworld/gomeboycolor/gbc"
	"github.com/djhworld/gomeboycolor/saves"
)

func StartEmulator(cartridge *cartridge.Cartridge, saveStore saves.Store, conf *config.Config) error {
	log.Println("Starting emulator")

	emulator, err := gbc.Init(cartridge, saveStore, conf)
	if err != nil {
		return err
	}

	//Starts emulator code in a goroutine
	go emulator.Run()

	//lock the OS thread here
	runtime.LockOSThread()

	//set the IO controller to run indefinitely (it waits for screen updates)
	emulator.RunIO()
	return nil
}
