package main

//Probably needlessly complex debug facility, but useful none the less

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"types"
	"utils"
)

type DebugRule struct {
	description  string
	ruleFunction func(g *GameboyColor) bool
}

func (d DebugRule) String() string {
	return d.description
}

type DebugRuleEngine struct {
	DebugRuleChain []DebugRule
}

func NewDebugRuleEngine() *DebugRuleEngine {
	var ruleEngine *DebugRuleEngine = new(DebugRuleEngine)
	ruleEngine.DebugRuleChain = []DebugRule{}
	return ruleEngine
}

func (pc *DebugRuleEngine) Set(s string) error {
	split := strings.Split(s, ",")
	var rules []Rule = []Rule{}
	for _, r := range split {
		rules = append(rules, Rule(r))
	}
	return pc.Parse(rules)
}

func (pc *DebugRuleEngine) String() string {
	var result string

	result += "\n"
	for i, r := range pc.DebugRuleChain {
		result += fmt.Sprintln("\t", i, "-", r)
	}

	return result
}

func (pc *DebugRuleEngine) Parse(rules []Rule) error {
	pc.DebugRuleChain = make([]DebugRule, 0)

	for _, rule := range rules {
		fn, err := rule.Parse()
		if err != nil {
			return err
		}

		pc.DebugRuleChain = append(pc.DebugRuleChain, fn)
	}

	return nil
}

type DebugCommandHandler func(*GameboyColor, ...string)

type DebugOptions struct {
	debuggerOn   bool
	ruleEngine   *DebugRuleEngine
	debugFuncMap map[string]DebugCommandHandler
}

func debugHelp() {
	fmt.Println("Commands are: -")
	fmt.Println("    - c (continue)")
	fmt.Println("    - s (step over)")
	fmt.Println("    - dump")
	fmt.Println("        + cpu (dump cpu state)")
	fmt.Println("            + registers (dump registers)")
	fmt.Println("            + flags (dump flags)")
	fmt.Println("        + mmu (dump memory)")
	fmt.Println("            + peripheral-map (dump memory addresses peripherals are interested in)")
	fmt.Println("    - exit (quit application)")
	fmt.Println("    - disconnect (stop debugging application and continue execution of application)")
}

func (g *DebugOptions) Init() {
	g.debuggerOn = false
	g.debugFuncMap = make(map[string]DebugCommandHandler)
	g.AddDebugFunc("dump", Dump)

	g.AddDebugFunc("exit", func(gbc *GameboyColor, remaining ...string) {
		os.Exit(0)
	})
	g.AddDebugFunc("disconnect", func(gbc *GameboyColor, remaining ...string) {
		gbc.debugOptions.debuggerOn = false
	})
	g.AddDebugFunc("s", func(gbc *GameboyColor, remaining ...string) {
		gbc.Step()
		fmt.Println(gbc.cpu.CurrentInstruction)
	})
}

func (g *DebugOptions) AddDebugFunc(command string, f DebugCommandHandler) {
	g.debugFuncMap[command] = f
}

func Dump(gbc *GameboyColor, remaining ...string) {
	if len(remaining) == 0 {
		debugHelp()
		return
	}

	switch remaining[0] {
	case "cpu":
		remaining = remaining[1:]
		DumpCPU(gbc, remaining...)
	case "mmu":
		remaining = remaining[1:]
		DumpMMU(gbc, remaining...)
	default:
		debugHelp()
	}
}

func DumpCPU(gbc *GameboyColor, remaining ...string) {
	if len(remaining) == 0 {
		debugHelp()
		return
	}

	switch remaining[0] {
	case "state":
		fmt.Println(gbc.cpu)
	case "registers":
		fmt.Println(gbc.cpu.R)
	case "flags":
		fmt.Println(gbc.cpu.FlagsString())
	default:
		debugHelp()
	}
}

func DumpMMU(gbc *GameboyColor, remaining ...string) {
	if len(remaining) == 0 {
		debugHelp()
		return
	}

	switch remaining[0] {
	case "range":
		DumpMemory(gbc)
	case "peripheral-map":
		gbc.mmu.PrintPeripheralMap()
	default:
		debugHelp()
	}
}

func DumpMemory(gbc *GameboyColor, remaining ...string) {
	b := bufio.NewWriter(os.Stdout)
	r := bufio.NewReader(os.Stdin)

	var memorystart string
	var memoryend string
	var start types.Word
	var end types.Word
	var err error

	fmt.Fprint(b, "Memory start address: 0x")
	b.Flush()

	memorystart, _ = r.ReadString('\n')
	memorystart = strings.Replace(memorystart, "\n", "", -1)
	start, err = ToMemoryAddress(memorystart)

	if err != nil {
		fmt.Fprintln(b, "Could not parse memory address", memorystart, ", please enter in the form 0x1234")
		fmt.Fprintln(b, err)
		b.Flush()
		return
	}

	fmt.Fprint(b, "Memory end address: 0x")
	b.Flush()

	memoryend, _ = r.ReadString('\n')
	memoryend = strings.Replace(memoryend, "\n", "", -1)
	end, err = ToMemoryAddress(memoryend)

	if err != nil {
		fmt.Fprintln(b, "Could not parse memory address", memoryend, ", please enter in the form 0x1234")
		fmt.Fprintln(b, err)
		b.Flush()
		return
	}

	printMemory := func(addr types.Word) {
		fmt.Fprintf(b, "[%s] -> 0x%X\n", addr, gbc.mmu.ReadByte(addr))
		b.Flush()
	}

	if start == end {
		printMemory(start)
	} else {
		for i := start; i < end; i++ {
			printMemory(i)
		}
		printMemory(end)
	}

	fmt.Fprintln(b)
	b.Flush()
}

func ToMemoryAddress(s string) (types.Word, error) {
	if len(s) > 4 {
		return 0x0, errors.New("Please enter an address between 0000 and FFFF")
	}

	result, err := strconv.ParseInt(s, 16, 64)

	return types.Word(result), err
}

type Rule string

func (r *Rule) Parse() (DebugRule, error) {
	var rule string = string(*r)
	rule = strings.ToLower(rule)
	exp, _ := regexp.Compile("(gbc.+|cpu.+|memory.+)(==|>|<|>=|<=)(0x[0-9a-f][0-9a-f][0-9a-f][0-9a-f]|0x[0-9a-f][0-9a-f])")
	matched := exp.MatchString(rule)

	var debugRuleFunc = DebugRule{rule, nil}

	if !matched {
		return debugRuleFunc, errors.New("Debug rule '" + rule + "' cannot be parsed. Format should be " + exp.String())
	}

	var statementInfo []string = exp.FindAllStringSubmatch(rule, -1)[0]
	var option []string = strings.Split(statementInfo[1], ".")
	var operator string = statementInfo[2]
	var valueStr string = statementInfo[3]
	valueStr = strings.Replace(valueStr, "0x", "", -1)
	val, _ := strconv.ParseInt(valueStr, 16, 64)
	var value types.Word = types.Word(val)

	switch option[0] {
	case "gbc":
		var sub string = ""

		if len(option) > 1 {
			sub = option[1]
		}

		switch sub {
		case "boot-mode":
			debugRuleFunc.ruleFunction = func(g *GameboyColor) bool {
				if byte(value) > 0 {
					return g.inBootMode == true
				}
				return g.inBootMode == false
			}
		default:
			return debugRuleFunc, errors.New("Unknown gameboy component: '" + sub + "'")
		}
	case "cpu":
		var sub string = ""

		if len(option) > 1 {
			sub = option[1]
		}

		switch sub {
		case "register":
			var register string = ""

			if len(option) > 2 {
				register = option[2]
			}

			switch register {
			case "a":
				debugRuleFunc.ruleFunction = func(g *GameboyColor) bool {
					return utils.CompareBytes(g.cpu.R.A, byte(value), operator)
				}
			case "b":
				debugRuleFunc.ruleFunction = func(g *GameboyColor) bool {
					return utils.CompareBytes(g.cpu.R.B, byte(value), operator)
				}
			case "c":
				debugRuleFunc.ruleFunction = func(g *GameboyColor) bool {
					return utils.CompareBytes(g.cpu.R.C, byte(value), operator)
				}
			case "d":
				debugRuleFunc.ruleFunction = func(g *GameboyColor) bool {
					return utils.CompareBytes(g.cpu.R.D, byte(value), operator)
				}
			case "e":
				debugRuleFunc.ruleFunction = func(g *GameboyColor) bool {
					return utils.CompareBytes(g.cpu.R.E, byte(value), operator)
				}
			case "f":
				debugRuleFunc.ruleFunction = func(g *GameboyColor) bool {
					return utils.CompareBytes(g.cpu.R.F, byte(value), operator)
				}
			case "h":
				debugRuleFunc.ruleFunction = func(g *GameboyColor) bool {
					return utils.CompareBytes(g.cpu.R.H, byte(value), operator)
				}
			case "l":
				debugRuleFunc.ruleFunction = func(g *GameboyColor) bool {
					return utils.CompareBytes(g.cpu.R.L, byte(value), operator)
				}
			default:
				return debugRuleFunc, errors.New("Unknown cpu register: '" + register + "'")
			}
		case "pc": //Program counter
			debugRuleFunc.ruleFunction = func(g *GameboyColor) bool {
				return utils.CompareWords(uint16(g.cpu.PC), uint16(value), operator)
			}
		case "sp": //Stack pointer
			debugRuleFunc.ruleFunction = func(g *GameboyColor) bool {
				return utils.CompareWords(uint16(g.cpu.SP), uint16(value), operator)
			}
		case "instruction":
			debugRuleFunc.ruleFunction = func(g *GameboyColor) bool {
				return utils.CompareBytes(g.cpu.CurrentInstruction.Opcode, byte(value), operator)
			}
		default:
			return debugRuleFunc, errors.New("Unknown cpu component: '" + sub + "'")
		}
	case "memory":
		fmt.Println("Asking for memory", value)
	default:
		return debugRuleFunc, errors.New("Unknown gameboy component: '" + option[0] + "'")
	}

	return debugRuleFunc, nil
}
