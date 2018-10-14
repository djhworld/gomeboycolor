package cpu

// Timing tests inspired by http://www.devrs.com/gb/files/opcodes.html

import (
	"testing"

	"github.com/djhworld/gomeboycolor/cartridge"
	"github.com/djhworld/gomeboycolor/timer"
	"github.com/djhworld/gomeboycolor/types"
	"github.com/djhworld/gomeboycolor/utils"
)

type FlagState struct {
	C bool
	H bool
	N bool
	Z bool
}

func TestEmptyInstruction(t *testing.T) {
	c := setupCPU(nil)
	EMPTY_INSTRUCTION.Execute(c)
}

func TestInstructionStringer(t *testing.T) {
	i := &Instruction{0xF2, "foo", 1, 2, func(c *GbcCPU) {}}
	s := i.String()
	if s != "0xF2 foo" {
		t.Log("Got result:", s)
		t.FailNow()
	}
}

func Test8bitLoadTimings(t *testing.T) {
	for _, opcode := range []byte{0x40, 0x41, 0x42, 0x43, 0x44, 0x45, 0x47, 0x48, 0x49, 0x4A, 0x4B, 0x4C, 0x4D, 0x4F, 0x50, 0x51, 0x52, 0x53, 0x54, 0x55, 0x57, 0x58, 0x59, 0x5A, 0x5B, 0x5C, 0x5D, 0x5F, 0x60, 0x61, 0x62, 0x63, 0x64, 0x65, 0x67, 0x68, 0x69, 0x6A, 0x6B, 0x6C, 0x6D, 0x6F, 0x78, 0x79, 0x7A, 0x7B, 0x7C, 0x7D, 0x7F} {
		RunInstrAndAssertTimings(opcode, nil, 1, t)
	}

	for _, opcode := range []byte{0x02, 0x06, 0x0A, 0x0E, 0x12, 0x16, 0x1A, 0x1E, 0x22, 0x26, 0x2A, 0x2E, 0x32, 0x3A, 0x3E, 0x46, 0x4E, 0x56, 0x5E, 0x66, 0x6E, 0x70, 0x71, 0x72, 0x73, 0x74, 0x75, 0x77, 0x7E, 0xE2, 0xF2, 0xF9} {
		RunInstrAndAssertTimings(opcode, nil, 2, t)
	}

	for _, opcode := range []byte{0x01, 0x11, 0x21, 0x31, 0x36, 0xF8} {
		RunInstrAndAssertTimings(opcode, nil, 3, t)
	}

	for _, opcode := range []byte{0xEA, 0xFA} {
		RunInstrAndAssertTimings(opcode, nil, 4, t)
	}

	// LDH
	for _, opcode := range []byte{0xE0, 0xF0} {
		RunInstrAndAssertTimings(opcode, nil, 3, t)
	}

	// LDI
	for _, opcode := range []byte{0x22, 0x2A} {
		RunInstrAndAssertTimings(opcode, nil, 2, t)
	}

	// LDD
	for _, opcode := range []byte{0x32, 0x3A} {
		RunInstrAndAssertTimings(opcode, nil, 2, t)
	}
}

func Test8bitALUTimings(t *testing.T) {
	// SUB|ADD|ADC|SBC|AND|OR|CP r
	for _, opcode := range []byte{0x80, 0x81, 0x82, 0x83, 0x84, 0x85, 0x87, 0x88, 0x89, 0x8A, 0x8B, 0x8C, 0x8D, 0x8F, 0x90, 0x91, 0x92, 0x93, 0x94, 0x95, 0x97, 0x98, 0x99, 0x9A, 0x9B, 0x9C, 0x9D, 0x9F, 0xA0, 0xA1, 0xA2, 0xA3, 0xA4, 0xA5, 0xA7, 0xA8, 0xA9, 0xAA, 0xAB, 0xAC, 0xAD, 0xAF, 0xB0, 0xB1, 0xB2, 0xB3, 0xB4, 0xB5, 0xB7, 0xB8, 0xB9, 0xBA, 0xBB, 0xBC, 0xBD, 0xBF} {
		RunInstrAndAssertTimings(opcode, nil, 1, t)
	}

	// SUB|ADD|ADC|SBC|AND|OR|CP (HL)
	for _, opcode := range []byte{0x86, 0x8E, 0x96, 0x9E, 0xA6, 0xAE, 0xB6, 0xBE} {
		RunInstrAndAssertTimings(opcode, nil, 2, t)
	}

	// SUB|ADD|ADC|SBC|AND|OR|CP d8
	for _, opcode := range []byte{0xC6, 0xCE, 0xD6, 0xDE, 0xE6, 0xEE, 0xF6, 0xFE} {
		RunInstrAndAssertTimings(opcode, nil, 2, t)
	}

	// INC r
	for _, opcode := range []byte{0x04, 0x0C, 0x14, 0x1C, 0x24, 0x2C, 0x3C} {
		RunInstrAndAssertTimings(opcode, nil, 1, t)
	}

	RunInstrAndAssertTimings(0x34, nil, 3, t) // INC (HL)

	// DEC r
	for _, opcode := range []byte{0x05, 0x0D, 0x15, 0x1D, 0x25, 0x2D, 0x3D} {
		RunInstrAndAssertTimings(opcode, nil, 1, t)
	}

	RunInstrAndAssertTimings(0x35, nil, 3, t) // DEC (HL)
}

func TestJumpTimings(t *testing.T) {
	RunInstrAndAssertTimings(0xC3, nil, 4, t) // JP nn

	RunInstrAndAssertTimings(0xC2, &FlagState{Z: true}, 3, t)  // JP cc, nn (Z - false)
	RunInstrAndAssertTimings(0xC2, &FlagState{Z: false}, 4, t) // JP cc, nn (Z - false)
	RunInstrAndAssertTimings(0xCA, &FlagState{Z: false}, 3, t) // JP cc, nn (Z - true)
	RunInstrAndAssertTimings(0xCA, &FlagState{Z: true}, 4, t)  // JP cc, nn (Z - true)

	RunInstrAndAssertTimings(0xD2, &FlagState{C: true}, 3, t)  // JP cc, nn (C - false)
	RunInstrAndAssertTimings(0xD2, &FlagState{C: false}, 4, t) // JP cc, nn (C - false)
	RunInstrAndAssertTimings(0xDA, &FlagState{C: false}, 3, t) // JP cc, nn (C - true)
	RunInstrAndAssertTimings(0xDA, &FlagState{C: true}, 4, t)  // JP cc, nn (C - true)

	RunInstrAndAssertTimings(0xE9, nil, 1, t) // JP (HL)

	RunInstrAndAssertTimings(0x18, nil, 3, t) // JR n

	RunInstrAndAssertTimings(0x20, &FlagState{Z: true}, 2, t)  // JR cc, nn (Z - false)
	RunInstrAndAssertTimings(0x20, &FlagState{Z: false}, 3, t) // JR cc, nn (Z - false)
	RunInstrAndAssertTimings(0x28, &FlagState{Z: false}, 2, t) // JR cc, nn (Z - true)
	RunInstrAndAssertTimings(0x28, &FlagState{Z: true}, 3, t)  // JR cc, nn (Z - true)

	RunInstrAndAssertTimings(0x30, &FlagState{C: true}, 2, t)  // JR cc, nn (C - false)
	RunInstrAndAssertTimings(0x30, &FlagState{C: false}, 3, t) // JR cc, nn (C - false)
	RunInstrAndAssertTimings(0x38, &FlagState{C: false}, 2, t) // JR cc, nn (C - true)
	RunInstrAndAssertTimings(0x38, &FlagState{C: true}, 3, t)  // JR cc, nn (C - true)

}

func TestCallTimings(t *testing.T) {
	RunInstrAndAssertTimings(0xCD, nil, 6, t) // CALL nn

	RunInstrAndAssertTimings(0xC4, &FlagState{Z: true}, 3, t)  // CALL cc, nn (Z - false)
	RunInstrAndAssertTimings(0xC4, &FlagState{Z: false}, 6, t) // CALL cc, nn (Z - false)
	RunInstrAndAssertTimings(0xCC, &FlagState{Z: false}, 3, t) // CALL cc, nn (Z - true)
	RunInstrAndAssertTimings(0xCC, &FlagState{Z: true}, 6, t)  // CALL cc, nn (Z - true)

	RunInstrAndAssertTimings(0xD4, &FlagState{C: true}, 3, t)  // CALL cc, nn (C - false)
	RunInstrAndAssertTimings(0xD4, &FlagState{C: false}, 6, t) // CALL cc, nn (C - false)
	RunInstrAndAssertTimings(0xDC, &FlagState{C: false}, 3, t) // CALL cc, nn (C - true)
	RunInstrAndAssertTimings(0xDC, &FlagState{C: true}, 6, t)  // CALL cc, nn (C - true)

}

func TestRestartTimings(t *testing.T) {
	opcodes := []byte{0xC7, 0xCF, 0xD7, 0xDF, 0xE7, 0xEF, 0xF7, 0xFF}

	for _, opcode := range opcodes {
		RunInstrAndAssertTimings(opcode, nil, 4, t) // RST n
	}
}

func TestReturnTimings(t *testing.T) {
	RunInstrAndAssertTimings(0xC9, nil, 4, t) // RET
	RunInstrAndAssertTimings(0xD9, nil, 4, t) // RETI

	RunInstrAndAssertTimings(0xC0, &FlagState{Z: true}, 2, t)  // RET cc, nn (Z - false)
	RunInstrAndAssertTimings(0xC0, &FlagState{Z: false}, 5, t) // RET cc, nn (Z - false)
	RunInstrAndAssertTimings(0xC8, &FlagState{Z: false}, 2, t) // RET cc, nn (Z - true)
	RunInstrAndAssertTimings(0xC8, &FlagState{Z: true}, 5, t)  // RET cc, nn (Z - true)

	RunInstrAndAssertTimings(0xD0, &FlagState{C: true}, 2, t)  // RET cc, nn (C - false)
	RunInstrAndAssertTimings(0xD0, &FlagState{C: false}, 5, t) // RET cc, nn (C - false)
	RunInstrAndAssertTimings(0xD8, &FlagState{C: false}, 2, t) // RET cc, nn (C - true)
	RunInstrAndAssertTimings(0xD8, &FlagState{C: true}, 5, t)  // RET cc, nn (C - true)

}

func TestMiscTimings(t *testing.T) {
	RunInstrAndAssertTimings(0x00, nil, 1, t) // NOP
	RunInstrAndAssertTimings(0x76, nil, 1, t) // HALT
	RunInstrAndAssertTimings(0x27, nil, 1, t) // DAA
	RunInstrAndAssertTimings(0x2F, nil, 1, t) // CPL
	RunInstrAndAssertTimings(0x3F, nil, 1, t) // CCF
	RunInstrAndAssertTimings(0x37, nil, 1, t) // SCF
	RunInstrAndAssertTimings(0xF3, nil, 1, t) // DI
	RunInstrAndAssertTimings(0xFB, nil, 1, t) // EI

	//TODO: this is supposed to be 1 cycle but we do the CGB speed check here....
	RunInstrAndAssertTimings(0x10, nil, 2, t) // STOP
}

func TestSwapTimings(t *testing.T) {
	RunCBInstrAndAssertTimings(0x36, nil, 4, t) // SWAP (HL)

	opcodes := []byte{
		0x30,
		0x31,
		0x32,
		0x33,
		0x34,
		0x35,
		0x37,
	}

	for _, opcode := range opcodes {
		RunCBInstrAndAssertTimings(opcode, nil, 2, t) // SWAP n
	}

}

func TestBitwiseTimings(t *testing.T) {

	// BIT n
	bitOpcodes := []byte{0x40, 0x41, 0x42, 0x43, 0x44, 0x45, 0x47, 0x48, 0x49, 0x4A, 0x4B, 0x4C, 0x4D, 0x4F, 0x50, 0x51, 0x52, 0x53, 0x54, 0x55, 0x57, 0x58, 0x59, 0x5A, 0x5B, 0x5C, 0x5D, 0x5F, 0x60, 0x61, 0x62, 0x63, 0x64, 0x65, 0x67, 0x68, 0x69, 0x6A, 0x6B, 0x6C, 0x6D, 0x6F, 0x70, 0x71, 0x72, 0x73, 0x74, 0x75, 0x77, 0x78, 0x79, 0x7A, 0x7B, 0x7C, 0x7D, 0x7F}
	for _, opcode := range bitOpcodes {
		RunCBInstrAndAssertTimings(opcode, nil, 2, t)
	}

	// BIT n, HL
	bitHLOpcodes := []byte{0x46, 0x4E, 0x56, 0x5E, 0x66, 0x6E, 0x76, 0x7E}
	for _, opcode := range bitHLOpcodes {
		RunCBInstrAndAssertTimings(opcode, nil, 3, t)
	}

	// SET n
	setOpcodes := []byte{0xC0, 0xC1, 0xC2, 0xC3, 0xC4, 0xC5, 0xC7, 0xC8, 0xC9, 0xCA, 0xCB, 0xCC, 0xCD, 0xCF, 0xD0, 0xD1, 0xD2, 0xD3, 0xD4, 0xD5, 0xD7, 0xD8, 0xD9, 0xDA, 0xDB, 0xDC, 0xDD, 0xDF, 0xE0, 0xE1, 0xE2, 0xE3, 0xE4, 0xE5, 0xE7, 0xE8, 0xE9, 0xEA, 0xEB, 0xEC, 0xED, 0xEF, 0xF0, 0xF1, 0xF2, 0xF3, 0xF4, 0xF5, 0xF7, 0xF8, 0xF9, 0xFA, 0xFB, 0xFC, 0xFD, 0xFF}
	for _, opcode := range setOpcodes {
		RunCBInstrAndAssertTimings(opcode, nil, 2, t)
	}

	// SET n, HL
	setHLOpcodes := []byte{0xC6, 0xCE, 0xD6, 0xDE, 0xE6, 0xEE, 0xF6, 0xFE}
	for _, opcode := range setHLOpcodes {
		RunCBInstrAndAssertTimings(opcode, nil, 4, t)
	}

	// RES n
	resOpcodes := []byte{0x80, 0x81, 0x82, 0x83, 0x84, 0x85, 0x87, 0x88, 0x89, 0x8A, 0x8B, 0x8C, 0x8D, 0x8F, 0x90, 0x91, 0x92, 0x93, 0x94, 0x95, 0x97, 0x98, 0x99, 0x9A, 0x9B, 0x9C, 0x9D, 0x9F, 0xA0, 0xA1, 0xA2, 0xA3, 0xA4, 0xA5, 0xA7, 0xA8, 0xA9, 0xAA, 0xAB, 0xAC, 0xAD, 0xAF, 0xB0, 0xB1, 0xB2, 0xB3, 0xB4, 0xB5, 0xB7, 0xB8, 0xB9, 0xBA, 0xBB, 0xBC, 0xBD, 0xBF}
	for _, opcode := range resOpcodes {
		RunCBInstrAndAssertTimings(opcode, nil, 2, t)
	}

	// RES n, HL
	resHLOpcodes := []byte{0x86, 0x8E, 0x96, 0x9E, 0xA6, 0xAE, 0xB6, 0xBE}
	for _, opcode := range resHLOpcodes {
		RunCBInstrAndAssertTimings(opcode, nil, 4, t)
	}

}

func Test16bitArithmeticTimings(t *testing.T) {

	RunInstrAndAssertTimings(0x09, nil, 2, t) //ADD HL, BC
	RunInstrAndAssertTimings(0x19, nil, 2, t) //ADD HL, DE
	RunInstrAndAssertTimings(0x29, nil, 2, t) //ADD HL, HL
	RunInstrAndAssertTimings(0x39, nil, 2, t) //ADD HL, SP

	RunInstrAndAssertTimings(0xE8, nil, 4, t) //ADD SP, e

	RunInstrAndAssertTimings(0x03, nil, 2, t) //INC HL, BC
	RunInstrAndAssertTimings(0x13, nil, 2, t) //INC HL, DE
	RunInstrAndAssertTimings(0x23, nil, 2, t) //INC HL, HL
	RunInstrAndAssertTimings(0x33, nil, 2, t) //INC HL, SP

	RunInstrAndAssertTimings(0x0B, nil, 2, t) //DEC HL, BC
	RunInstrAndAssertTimings(0x1B, nil, 2, t) //DEC HL, DE
	RunInstrAndAssertTimings(0x2B, nil, 2, t) //DEC HL, HL
	RunInstrAndAssertTimings(0x3B, nil, 2, t) //DEC HL, SP
}

func Test16bitLoadTimings(t *testing.T) {
	// PUSH
	for _, opcode := range []byte{0xC5, 0xD5, 0xE5, 0xF5} {
		RunInstrAndAssertTimings(opcode, nil, 4, t)
	}

	// POP
	for _, opcode := range []byte{0xC1, 0xD1, 0xE1, 0xF1} {
		RunInstrAndAssertTimings(opcode, nil, 3, t)
	}

	RunInstrAndAssertTimings(0x01, nil, 3, t) // LD BC,nn
	RunInstrAndAssertTimings(0x11, nil, 3, t) // LD DE,nn
	RunInstrAndAssertTimings(0x21, nil, 3, t) // LD HL,nn
	RunInstrAndAssertTimings(0x31, nil, 3, t) // LD SP,nn

	RunInstrAndAssertTimings(0x08, nil, 5, t) // LD nn,SP

	RunInstrAndAssertTimings(0xF9, nil, 2, t) // LD SP,HL
	RunInstrAndAssertTimings(0xF8, nil, 3, t) // LD HL,(SP+e)

}

func TestShiftTimings(t *testing.T) {
	// SLA r
	for _, opcode := range []byte{0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x27} {
		RunCBInstrAndAssertTimings(opcode, nil, 2, t)
	}

	// SRA r
	for _, opcode := range []byte{0x28, 0x29, 0x2A, 0x2B, 0x2C, 0x2D, 0x2F} {
		RunCBInstrAndAssertTimings(opcode, nil, 2, t)
	}

	// SRL r
	for _, opcode := range []byte{0x38, 0x39, 0x3A, 0x3B, 0x3C, 0x3D, 0x3F} {
		RunCBInstrAndAssertTimings(opcode, nil, 2, t)
	}

	RunCBInstrAndAssertTimings(0x26, nil, 4, t) // SLA (HL)
	RunCBInstrAndAssertTimings(0x2E, nil, 4, t) // SRA (HL)
	RunCBInstrAndAssertTimings(0x3E, nil, 4, t) // SRL (HL)
}

func TestRotateTimings(t *testing.T) {
	RunInstrAndAssertTimings(0x07, nil, 1, t) // RLCA
	RunInstrAndAssertTimings(0x17, nil, 1, t) // RLA
	RunInstrAndAssertTimings(0x0F, nil, 1, t) // RRCA
	RunInstrAndAssertTimings(0x1F, nil, 1, t) // RRA

	// RLC r
	for _, opcode := range []byte{0x0, 0x1, 0x2, 0x3, 0x4, 0x5, 0x7} {
		RunCBInstrAndAssertTimings(opcode, nil, 2, t)
	}

	// RL r
	for _, opcode := range []byte{0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x17} {
		RunCBInstrAndAssertTimings(opcode, nil, 2, t)
	}

	// RRC r
	for _, opcode := range []byte{0x8, 0x9, 0xA, 0xB, 0xC, 0xD, 0xF} {
		RunCBInstrAndAssertTimings(opcode, nil, 2, t)
	}

	// RRC r
	for _, opcode := range []byte{0x8, 0x9, 0xA, 0xB, 0xC, 0xD, 0xF} {
		RunCBInstrAndAssertTimings(opcode, nil, 2, t)
	}

	// RR r
	for _, opcode := range []byte{0x18, 0x19, 0x1A, 0x1B, 0x1C, 0x1D, 0x1F} {
		RunCBInstrAndAssertTimings(opcode, nil, 2, t)
	}

	RunCBInstrAndAssertTimings(0x06, nil, 4, t) // RLC (HL)
	RunCBInstrAndAssertTimings(0x16, nil, 4, t) // RL (HL)
	RunCBInstrAndAssertTimings(0x0E, nil, 4, t) // RRC (HL)
	RunCBInstrAndAssertTimings(0x1E, nil, 4, t) // RR (HL)
}

func assertTimings(c *GbcCPU, t *testing.T, instr byte, expectedTiming int, isCB bool) {
	cputicks := c.Step()

	if isCB {
		t.Log("0xCB "+utils.ByteToString(instr)+" ("+c.CurrentInstruction.Description+")", "testing that instruction runs for", expectedTiming, "cycles")
	} else {
		t.Log(utils.ByteToString(instr)+" ("+c.CurrentInstruction.Description+")", "testing that instruction runs for", expectedTiming, "cycles")
	}

	if cputicks != expectedTiming {
		if isCB {
			t.Log("FAILED -----> instruction 0xCB", utils.ByteToString(instr)+" ("+c.CurrentInstruction.Description+")", "Expected", expectedTiming, "but got", cputicks)
		} else {
			t.Log("FAILED -----> instruction", utils.ByteToString(instr)+" ("+c.CurrentInstruction.Description+")", "Expected", expectedTiming, "but got", cputicks)

		}
		t.Fail()
	}
}

func setupCPU(flags *FlagState) *GbcCPU {
	var c *GbcCPU = NewCPU(NewMockMMU(), timer.NewTimer())
	if flags != nil {
		if flags.Z {
			c.SetFlag(Z)
		}
		if flags.C {
			c.SetFlag(C)
		}
		if flags.N {
			c.SetFlag(N)
		}
		if flags.H {
			c.SetFlag(H)
		}
	} else {
		c.ResetFlag(Z)
		c.ResetFlag(H)
		c.ResetFlag(N)
		c.ResetFlag(C)
		c.R.F = 0x00
	}
	c.PC = 0x0000

	return c
}

func RunCBInstrAndAssertTimings(instr byte, flagState *FlagState, expectedTiming int, t *testing.T) {
	c := setupCPU(flagState)
	c.WriteByte(c.PC, 0xCB)
	c.WriteByte(c.PC+1, instr)
	assertTimings(c, t, instr, expectedTiming, true)
}

func RunInstrAndAssertTimings(instr byte, flagState *FlagState, expectedTiming int, t *testing.T) {
	c := setupCPU(flagState)
	c.WriteByte(c.PC, instr)
	assertTimings(c, t, instr, expectedTiming, false)
}

type MockMMU struct {
	memory map[types.Word]byte
}

func NewMockMMU() *MockMMU {
	var m *MockMMU = new(MockMMU)
	m.memory = make(map[types.Word]byte)
	return m
}

func (m *MockMMU) WriteByte(address types.Word, value byte) {
	m.memory[address] = value
}

func (m *MockMMU) WriteWord(address types.Word, value types.Word) {
	m.memory[address] = byte(value >> 8)
	m.memory[address+1] = byte(value & 0x00FF)
}

func (m *MockMMU) ReadByte(address types.Word) byte {
	return m.memory[address]
}

func (m *MockMMU) ReadWord(address types.Word) types.Word {
	a, b := m.memory[address], m.memory[address+1]
	return (types.Word(a) << 8) ^ types.Word(b)
}

func (m *MockMMU) SetInBootMode(mode bool) {
}

func (m *MockMMU) Reset() {
	m.memory = make(map[types.Word]byte)
}

func (m *MockMMU) LoadBIOS(data []byte) (bool, error) {
	return true, nil
}

func (m *MockMMU) LoadCartridge(cart *cartridge.Cartridge) {
}
