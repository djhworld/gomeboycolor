package main

//Probably needlessly complex debug facility, but useful none the less

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"types"
	"utils"
)

type DebugCommandHandler func(*GameboyColor, ...string)

type DebugOptions struct {
	debuggerOn   bool
	breakWhen    types.Word
	debugFuncMap map[string]DebugCommandHandler
	debugHelpStr []string
}

func (g *DebugOptions) help() {
	fmt.Println("Commands are: -")
	for _, desc := range g.debugHelpStr {
		fmt.Println("	-", desc)
	}
}

func (g *DebugOptions) Init() {
	g.debuggerOn = false
	g.debugFuncMap = make(map[string]DebugCommandHandler)
	g.AddDebugFunc("p", "Print CPU state", func(gbc *GameboyColor, remaining ...string) {
		fmt.Println(gbc.cpu)
	})

	g.AddDebugFunc("r", "Reset", func(gbc *GameboyColor, remaining ...string) {
		gbc.Reset()
	})

	g.AddDebugFunc("?", "Print this help message", func(gbc *GameboyColor, remaining ...string) {
		g.help()
	})

	g.AddDebugFunc("help", "Print this help message", func(gbc *GameboyColor, remaining ...string) {
		g.help()
	})

	g.AddDebugFunc("d", "Disconnect from debugger", func(gbc *GameboyColor, remaining ...string) {
		gbc.debugOptions.debuggerOn = false
	})

	g.AddDebugFunc("s", "Step", func(gbc *GameboyColor, remaining ...string) {
		gbc.Step()
		fmt.Println(gbc.cpu)
	})

	g.AddDebugFunc("b", "Set breakpoint", func(gbc *GameboyColor, remaining ...string) {
		if len(remaining) == 0 {
			fmt.Println("You must provide a PC address to break on!")
			return
		}

		var arg string = remaining[0]

		if bp, err := ToMemoryAddress(arg); err != nil {
			fmt.Println("Could not parse memory address argument:", arg)
			fmt.Println("\t", err)
		} else {
			fmt.Println("Setting breakpoint to:", bp)
			g.breakWhen = bp
		}
	})

	g.AddDebugFunc("rm", "Read data from memory", func(gbc *GameboyColor, remaining ...string) {
		var startAddr types.Word
		switch len(remaining) {
		case 0:
			fmt.Println("You must provide at least a starting address to inspect")
			return
		default:
			addr, err := ToMemoryAddress(remaining[0])
			if err != nil {
				fmt.Println("Could not parse memory address: ", remaining[0])
				return
			}
			startAddr = addr
		}
		lb := startAddr - (startAddr & 0x000F)
		hb := startAddr + (0x0F - (startAddr & 0x000F))
		fmt.Print("\t\t")
		for w := lb; w <= hb; w++ {
			fmt.Printf("   %X ", byte(w%16))

		}
		fmt.Println()

		fmt.Printf("%s\t\t", lb)
		for w := lb; w <= hb; w++ {
			fmt.Print(utils.ByteToString(gbc.mmu.ReadByte(w)), " ")
		}
		fmt.Println()

	})

	g.AddDebugFunc("wm", "Write data to memory", func(gbc *GameboyColor, remaining ...string) {
		var value byte
		var toAddr types.Word
		switch len(remaining) {
		case 0:
			fmt.Println("You must provide a byte value and address to write to")
			return
		case 1:
			fmt.Println("You must provide a byte value to put in memory")
			return
		default:
			addr, err := ToMemoryAddress(remaining[0])
			if err != nil {
				fmt.Println("Could not parse memory address: ", remaining[0])
				return
			}
			val, err := utils.StringToByte(remaining[1])
			if err != nil {
				fmt.Println("Could not parse value: ", remaining[1], err)
				return
			}
			toAddr = addr
			value = val
		}

		fmt.Println("Writing", utils.ByteToString(value), "to", toAddr)
		gbc.mmu.WriteByte(toAddr, value)
	})

	g.AddDebugFunc("q", "Quit emulator", func(gbc *GameboyColor, remaining ...string) {
		os.Exit(0)
	})
}

func (g *DebugOptions) AddDebugFunc(command string, description string, f DebugCommandHandler) {
	g.debugFuncMap[command] = f
	g.debugHelpStr = append(g.debugHelpStr, utils.PadRight(command, 4, " ")+" = "+description)
}

func ToMemoryAddress(s string) (types.Word, error) {
	if len(s) > 4 {
		return 0x0, errors.New("Please enter an address between 0000 and FFFF")
	}

	result, err := strconv.ParseInt(s, 16, 64)

	return types.Word(result), err
}
