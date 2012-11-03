package main

import "testing"
import "github.com/stretchrcom/testify/assert"

const ZEROW Word = 0
const ZEROB byte = 0

var cpu *Z80 = NewCPU(NewMockMMU())

func reset() {
	cpu = nil
	cpu = NewCPU(NewMockMMU())
	cpu.Reset()
}

func TestIncrementPC(t *testing.T) {
	reset()
	cpu.PC = 0x02
	cpu.IncrementPC(1)
	assert.Equal(t, cpu.PC, Word(0x03))
}

// INSTRUCTIONS START
//-----------------------------------------------------------------------

//ADD A,r tests
//------------------------------------------
func TestAddA_r(t *testing.T) {
	reset()

	cpu.R.A = 0x01
	cpu.R.B = 0x02
	cpu.AddA_r(&cpu.R.B)
	assert.Equal(t, cpu.R.A, byte(0x03))
}

func TestAddA_rForZeroFlag(t *testing.T) {
	reset()
	cpu.R.A = 0x00
	cpu.R.B = 0x00

	//Check flag is set when result is zero
	cpu.AddA_r(&cpu.R.B)
	assert.Equal(t, (cpu.R.F & 0x80), byte(0x80))

	reset()

	//Check flag is NOT set when result is not zero
	cpu.R.B = 0x09
	cpu.AddA_r(&cpu.R.B)
	assert.Equal(t, (cpu.R.F & 0x80), byte(0x00))
}

func TestAddA_rForCarryFlag(t *testing.T) {
	reset()
	cpu.R.A = 0xFA
	cpu.R.E = 0x19

	//Check flag is set when result is zero
	cpu.AddA_r(&cpu.R.E)
	assert.Equal(t, (cpu.R.F & 0x40), byte(0x40))

	reset()
	cpu.R.A = 0xFA
	cpu.R.E = 0x02
	//Check flag is NOT set when result is not zero
	cpu.AddA_r(&cpu.R.E)
	assert.Equal(t, (cpu.R.F & 0x40), byte(0x00))
}

func TestAddA_rClockTimings(t *testing.T) {
	reset()
	cpu.R.B = 0x10
	cpu.AddA_r(&cpu.R.B)

	assert.Equal(t, cpu.LastInstrCycle.m, Word(1))
	assert.Equal(t, cpu.LastInstrCycle.t, Word(4))

}

//ADD A,(HL) tests
//------------------------------------------
func TestAddA_hl(t *testing.T) {
	reset()

	cpu.R.A = 0x05
	cpu.R.H = 0x22
	cpu.R.L = 0x23

	cpu.mmu.WriteByte(0x2223, 0x04)
	cpu.AddA_hl()
	assert.Equal(t, cpu.R.A, byte(0x09))
}

func TestAddA_hlForZeroFlag(t *testing.T) {
	reset()
	cpu.R.A = 0x00
	cpu.R.H = 0x03
	cpu.R.L = 0x02
	cpu.mmu.WriteByte(0x0302, 0x00)

	//Check flag is set when result is zero
	cpu.AddA_hl()
	assert.Equal(t, (cpu.R.F & 0x80), byte(0x80))

	reset()

	cpu.R.A = 0x00
	cpu.R.H = 0x03
	cpu.R.L = 0x02
	cpu.mmu.WriteByte(0x0302, 0x02)

	//Check flag is NOT set when result is not zero
	cpu.AddA_hl()
	assert.Equal(t, (cpu.R.F & 0x80), byte(0x00))
}

func TestAddA_hlForCarryFlag(t *testing.T) {
	var memoryAddr Word = 0x0302
	var H byte = 0x03
	var L byte = 0x02

	reset()
	cpu.R.A = 0xFE
	cpu.R.H = H
	cpu.R.L = L
	cpu.mmu.WriteByte(memoryAddr, 0x03)

	//Check flag is set when result overflows
	cpu.AddA_hl()
	assert.Equal(t, (cpu.R.F & 0x40), byte(0x40))

	reset()

	cpu.R.A = 0xFE
	cpu.R.H = H
	cpu.R.L = L
	cpu.mmu.WriteByte(memoryAddr, 0x01)

	//Check flag is set when result does not
	cpu.AddA_hl()
	assert.Equal(t, (cpu.R.F & 0x40), byte(0x00))
}

func TestAddA_hlClockTimings(t *testing.T) {
	reset()
	cpu.AddA_hl()
	assert.Equal(t, cpu.LastInstrCycle.m, Word(2))
	assert.Equal(t, cpu.LastInstrCycle.t, Word(8))

}

//ADD A,n tests
//------------------------------------------
func TestAddA_n(t *testing.T) {
	reset()

	cpu.R.A = 0x09
	cpu.PC = 0x000E
	cpu.mmu.WriteByte(0x000E, 0x01)

	cpu.AddA_n()
	assert.Equal(t, cpu.R.A, byte(0x0A))
}

func TestAddA_nCheckPC(t *testing.T) {
	reset()

	cpu.R.A = 0x09
	cpu.PC = 0x000E
	cpu.mmu.WriteByte(0x000E, 0x01)

	cpu.AddA_n()
	assert.Equal(t, cpu.PC, Word(0x000F))

}

func TestAddA_nForZeroFlag(t *testing.T) {
	reset()

	cpu.R.A = 0x00
	cpu.PC = 0x000E
	cpu.mmu.WriteByte(0x000E, 0x00)

	//Check flag is set when result is zero
	cpu.AddA_n()
	assert.Equal(t, (cpu.R.F & 0x80), byte(0x80))

	reset()
	cpu.R.A = 0x09
	cpu.PC = 0x000E
	cpu.mmu.WriteByte(0x000E, 0x01)

	//Check flag is NOT set when result is not zero
	cpu.AddA_n()
	assert.Equal(t, (cpu.R.F & 0x80), byte(0x00))
}

func TestAddA_nForCarryFlag(t *testing.T) {
	reset()

	cpu.R.A = 0xFE
	cpu.PC = 0x000E
	cpu.mmu.WriteByte(0x000E, 0x03)

	//Check flag is set when result overflows
	cpu.AddA_n()
	assert.Equal(t, (cpu.R.F & 0x40), byte(0x40))

	reset()

	cpu.R.A = 0xFE
	cpu.PC = 0x000E
	cpu.mmu.WriteByte(0x000E, 0x01)

	//Check flag is set when result does not
	cpu.AddA_n()
	assert.Equal(t, (cpu.R.F & 0x40), byte(0x00))
}

func TestAddA_nClockTimings(t *testing.T) {
	reset()
	cpu.AddA_n()
	assert.Equal(t, cpu.LastInstrCycle.m, Word(2))
	assert.Equal(t, cpu.LastInstrCycle.t, Word(8))
}

//LD r,n tests
//------------------------------------------
func TestLDrn(t *testing.T) {
	var expected byte = 0x0A
	var addr Word = 0x0001

	//test B
	reset()
	cpu.PC = addr
	cpu.mmu.WriteByte(addr, expected)
	cpu.LDrn(&cpu.R.B)
	assert.Equal(t, cpu.R.B, expected)

	//test C
	reset()
	cpu.PC = addr
	cpu.mmu.WriteByte(addr, expected)
	cpu.LDrn(&cpu.R.C)
	assert.Equal(t, cpu.R.C, expected)

	//test D
	reset()
	cpu.PC = addr
	cpu.mmu.WriteByte(addr, expected)
	cpu.LDrn(&cpu.R.D)
	assert.Equal(t, cpu.R.D, expected)

	//test E
	reset()
	cpu.PC = addr
	cpu.mmu.WriteByte(addr, expected)
	cpu.LDrn(&cpu.R.E)
	assert.Equal(t, cpu.R.E, expected)

	//test H
	reset()
	cpu.PC = addr
	cpu.mmu.WriteByte(addr, expected)
	cpu.LDrn(&cpu.R.H)
	assert.Equal(t, cpu.R.H, expected)

	//test L
	reset()
	cpu.PC = addr
	cpu.mmu.WriteByte(addr, expected)
	cpu.LDrn(&cpu.R.L)
	assert.Equal(t, cpu.R.L, expected)
}

func TestLDrn_PCIncrementCheck(t *testing.T) {
	var addr Word = 0x0001
	var expected Word = 0x0002
	reset()
	cpu.PC = addr
	cpu.LDrn(&cpu.R.B)

	assert.Equal(t, cpu.PC, expected)
}

func TestLDrn_ClockTimings(t *testing.T) {
	reset()
	cpu.LDrn(&cpu.R.B)
	assert.Equal(t, cpu.LastInstrCycle.m, Word(2))
	assert.Equal(t, cpu.LastInstrCycle.t, Word(8))
}

//LD r,r tests
//------------------------------------------
func TestLDrr(t *testing.T) {
	var expected byte = 0x0A

	//test A <- A
	reset()
	cpu.R.A = expected
	cpu.LDrr(&cpu.R.A, &cpu.R.A)
	assert.Equal(t, cpu.R.A, expected)

	//test A <- B
	reset()
	cpu.R.A = 0x03
	cpu.R.B = expected
	cpu.LDrr(&cpu.R.A, &cpu.R.B)
	assert.Equal(t, cpu.R.A, expected)

	//test A <- C
	reset()
	cpu.R.A = 0x03
	cpu.R.C = expected
	cpu.LDrr(&cpu.R.A, &cpu.R.C)
	assert.Equal(t, cpu.R.A, expected)

	//test A <- D
	reset()
	cpu.R.A = 0x03
	cpu.R.D = expected
	cpu.LDrr(&cpu.R.A, &cpu.R.D)
	assert.Equal(t, cpu.R.A, expected)

	//test A <- E
	reset()
	cpu.R.A = 0x03
	cpu.R.E = expected
	cpu.LDrr(&cpu.R.A, &cpu.R.E)
	assert.Equal(t, cpu.R.A, expected)

	//test A <- H
	reset()
	cpu.R.A = 0x03
	cpu.R.H = expected
	cpu.LDrr(&cpu.R.A, &cpu.R.H)
	assert.Equal(t, cpu.R.A, expected)

	//test A <- L
	reset()
	cpu.R.A = 0x03
	cpu.R.L = expected
	cpu.LDrr(&cpu.R.A, &cpu.R.L)
	assert.Equal(t, cpu.R.A, expected)

	//test B <- B
	reset()
	cpu.R.B = 0x03
	cpu.R.B = expected
	cpu.LDrr(&cpu.R.B, &cpu.R.B)
	assert.Equal(t, cpu.R.B, expected)

	//test B <- C
	reset()
	cpu.R.B = 0x03
	cpu.R.C = expected
	cpu.LDrr(&cpu.R.B, &cpu.R.C)
	assert.Equal(t, cpu.R.B, expected)

	//test B <- D
	reset()
	cpu.R.B = 0x03
	cpu.R.D = expected
	cpu.LDrr(&cpu.R.B, &cpu.R.D)
	assert.Equal(t, cpu.R.B, expected)

	//test B <- E
	reset()
	cpu.R.B = 0x03
	cpu.R.E = expected
	cpu.LDrr(&cpu.R.B, &cpu.R.E)
	assert.Equal(t, cpu.R.B, expected)

	//test B <- H
	reset()
	cpu.R.B = 0x03
	cpu.R.H = expected
	cpu.LDrr(&cpu.R.B, &cpu.R.H)
	assert.Equal(t, cpu.R.B, expected)

	//test B <- L
	reset()
	cpu.R.B = 0x03
	cpu.R.L = expected
	cpu.LDrr(&cpu.R.B, &cpu.R.L)
	assert.Equal(t, cpu.R.B, expected)

	//test C <- B
	reset()
	cpu.R.C = 0x03
	cpu.R.B = expected
	cpu.LDrr(&cpu.R.C, &cpu.R.B)
	assert.Equal(t, cpu.R.C, expected)

	//test C <- C
	reset()
	cpu.R.C = 0x03
	cpu.R.C = expected
	cpu.LDrr(&cpu.R.C, &cpu.R.C)
	assert.Equal(t, cpu.R.C, expected)

	//test C <- D
	reset()
	cpu.R.C = 0x03
	cpu.R.D = expected
	cpu.LDrr(&cpu.R.C, &cpu.R.D)
	assert.Equal(t, cpu.R.C, expected)

	//test C <- E
	reset()
	cpu.R.C = 0x03
	cpu.R.E = expected
	cpu.LDrr(&cpu.R.C, &cpu.R.E)
	assert.Equal(t, cpu.R.C, expected)

	//test C <- H
	reset()
	cpu.R.C = 0x03
	cpu.R.H = expected
	cpu.LDrr(&cpu.R.C, &cpu.R.H)
	assert.Equal(t, cpu.R.C, expected)

	//test C <- L
	reset()
	cpu.R.C = 0x03
	cpu.R.L = expected
	cpu.LDrr(&cpu.R.C, &cpu.R.L)
	assert.Equal(t, cpu.R.C, expected)

}

func TestLDrr_ClockTimings(t *testing.T) {
	reset()
	cpu.LDrr(&cpu.R.A, &cpu.R.B)
	assert.Equal(t, cpu.LastInstrCycle.m, Word(1))
	assert.Equal(t, cpu.LastInstrCycle.t, Word(4))
}

//LD r,(HL) tests
//------------------------------------------
func TestLDr_hl(t *testing.T) {
	var expected byte = 0xAE
	var addr Word = 0x1002

	reset()
	cpu.R.H, cpu.R.L = 0x10, 0x02

	cpu.mmu.WriteByte(addr, expected)

	//A <- (HL)
	cpu.LDr_hl(&cpu.R.A)
	assert.Equal(t, cpu.R.A, expected)

	//B <- (HL)
	cpu.LDr_hl(&cpu.R.B)
	assert.Equal(t, cpu.R.B, expected)

	//C <- (HL)
	cpu.LDr_hl(&cpu.R.C)
	assert.Equal(t, cpu.R.C, expected)
}

func TestLDr_hlClockTimings(t *testing.T) {
	reset()
	cpu.LDr_hl(&cpu.R.A)
	assert.Equal(t, cpu.LastInstrCycle.m, Word(2))
	assert.Equal(t, cpu.LastInstrCycle.t, Word(8))
}

//LD (HL),r tests
//------------------------------------------
func TestLDhl_r(t *testing.T) {
	var expected byte = 0xF3
	var addr Word = 0x1002

	// (HL) <- B
	reset()
	cpu.R.H, cpu.R.L = 0x10, 0x02
	cpu.R.B = expected
	cpu.LDhl_r(&cpu.R.B)
	assert.Equal(t, cpu.mmu.ReadByte(addr), expected)

	// (HL) <- C
	reset()
	cpu.R.H, cpu.R.L = 0x10, 0x02
	cpu.R.C = expected
	cpu.LDhl_r(&cpu.R.C)
	assert.Equal(t, cpu.mmu.ReadByte(addr), expected)

	// (HL) <- D
	reset()
	cpu.R.H, cpu.R.L = 0x10, 0x02
	cpu.R.D = expected
	cpu.LDhl_r(&cpu.R.D)
	assert.Equal(t, cpu.mmu.ReadByte(addr), expected)
}

func TestLDhl_rClockTimings(t *testing.T) {
	reset()
	cpu.LDhl_r(&cpu.R.B)
	assert.Equal(t, cpu.LastInstrCycle.m, Word(2))
	assert.Equal(t, cpu.LastInstrCycle.t, Word(8))
}

//LD (HL),n tests
//------------------------------------------
func TestLDhl_n(t *testing.T) {
	var expected byte = 0x3A
	var addr Word = 0x0201
	var HL Word = 0xFFEE

	//test B
	reset()
	cpu.R.H, cpu.R.L = 0xFF, 0xEE
	cpu.mmu.WriteByte(addr, expected)
	cpu.PC = addr
	cpu.LDhl_n()
	assert.Equal(t, cpu.mmu.ReadByte(HL), expected)
}

func TestLDhl_nPCIncrementCheck(t *testing.T) {
	var addr Word = 0x0001
	var expected Word = 0x0002
	reset()
	cpu.PC = addr
	cpu.LDhl_n()

	assert.Equal(t, cpu.PC, expected)
}

func TestLDhl_nClockTimings(t *testing.T) {
	reset()
	cpu.LDhl_n()
	assert.Equal(t, cpu.LastInstrCycle.m, Word(3))
	assert.Equal(t, cpu.LastInstrCycle.t, Word(12))
}

//LD r,(BC) tests
//------------------------------------------
func TestLDr_bc(t *testing.T) {
	var expected byte = 0xAE
	var addr Word = 0x1002

	reset()
	cpu.R.B, cpu.R.C = 0x10, 0x02

	cpu.mmu.WriteByte(addr, expected)

	//A <- (BC)
	cpu.LDr_bc(&cpu.R.A)
	assert.Equal(t, cpu.R.A, expected)
}

func TestLDr_bcClockTimings(t *testing.T) {
	reset()
	cpu.LDr_bc(&cpu.R.A)
	assert.Equal(t, cpu.LastInstrCycle.m, Word(2))
	assert.Equal(t, cpu.LastInstrCycle.t, Word(8))
}

//LD r,(DE) tests
//------------------------------------------
func TestLDr_de(t *testing.T) {
	var expected byte = 0xAE
	var addr Word = 0x1002

	reset()
	cpu.R.D, cpu.R.E = 0x10, 0x02

	cpu.mmu.WriteByte(addr, expected)

	//A <- (DE)
	cpu.LDr_de(&cpu.R.A)
	assert.Equal(t, cpu.R.A, expected)
}

func TestLDr_deClockTimings(t *testing.T) {
	reset()
	cpu.LDr_de(&cpu.R.A)
	assert.Equal(t, cpu.LastInstrCycle.m, Word(2))
	assert.Equal(t, cpu.LastInstrCycle.t, Word(8))
}

//LD r,nn tests
//------------------------------------------
func TestLDr_nn(t *testing.T) {
	var expected byte = 0x3E
	var addr Word = 0x0002
	var valueAddr Word = 0x3334
	cpu.mmu.WriteByte(addr, 0x33)
	cpu.mmu.WriteByte(addr+1, 0x34)
	cpu.mmu.WriteByte(valueAddr, expected)
	cpu.PC = addr

	cpu.LDr_nn(&cpu.R.A)
	assert.Equal(t, cpu.R.A, expected)
}

func TestLDr_nnPCIncremented(t *testing.T) {
	var addr Word = 0x0002
	var expected Word = 0x0004
	cpu.PC = addr
	cpu.LDr_nn(&cpu.R.A)
	assert.Equal(t, cpu.PC, expected)
}

func TestLDr_nnClockTimings(t *testing.T) {
	reset()
	cpu.LDr_nn(&cpu.R.A)
	assert.Equal(t, cpu.LastInstrCycle.m, Word(4))
	assert.Equal(t, cpu.LastInstrCycle.t, Word(16))
}

//LD (BC),r tests
//------------------------------------------
func TestLDbc_r(t *testing.T) {
	var expected byte = 0xF3
	var addr Word = 0x1002

	// (BC) <- A
	reset()
	cpu.R.B, cpu.R.C = 0x10, 0x02
	cpu.R.A = expected
	cpu.LDbc_r(&cpu.R.A)
	assert.Equal(t, cpu.mmu.ReadByte(addr), expected)
}

func TestLDbc_rClockTimings(t *testing.T) {
	reset()
	cpu.LDbc_r(&cpu.R.A)
	assert.Equal(t, cpu.LastInstrCycle.m, Word(2))
	assert.Equal(t, cpu.LastInstrCycle.t, Word(8))
}

//LD (DE),r tests
//------------------------------------------
func TestLDde_r(t *testing.T) {
	var expected byte = 0xF3
	var addr Word = 0x1002

	// (DE) <- A
	reset()
	cpu.R.D, cpu.R.E = 0x10, 0x02
	cpu.R.A = expected
	cpu.LDde_r(&cpu.R.A)
	assert.Equal(t, cpu.mmu.ReadByte(addr), expected)
}

func TestLDde_rClockTimings(t *testing.T) {
	reset()
	cpu.LDde_r(&cpu.R.A)
	assert.Equal(t, cpu.LastInstrCycle.m, Word(2))
	assert.Equal(t, cpu.LastInstrCycle.t, Word(8))
}

//LD nn,r tests
//------------------------------------------
func TestLDnn_r(t *testing.T) {
	reset()
	var expected byte = 0x7C
	var addr Word = 0x0005
	var valueAddr Word = 0x3031

	cpu.mmu.WriteByte(addr, 0x30)
	cpu.mmu.WriteByte(addr+1, 0x31)

	cpu.PC = addr
	cpu.R.A = expected
	cpu.LDnn_r(&cpu.R.A)

	assert.Equal(t, cpu.mmu.ReadByte(valueAddr), expected)
}

func TestLDnn_rPCIncremented(t *testing.T) {
	reset()
	var addr Word = 0x0002
	var expected Word = 0x0004
	cpu.PC = addr
	cpu.LDnn_r(&cpu.R.A)
	assert.Equal(t, cpu.PC, expected)
}

func TestLDnn_rClockTimings(t *testing.T) {
	reset()
	cpu.LDnn_r(&cpu.R.A)

	assert.Equal(t, cpu.LastInstrCycle.m, Word(4))
	assert.Equal(t, cpu.LastInstrCycle.t, Word(16))
}

//LD r,(C) tests
//------------------------------------------
func TestLDr_ffplusc(t *testing.T) {
	reset()
	var valueAddr Word = 0xFF03
	var expected byte = 0x03
	cpu.R.C = expected

	cpu.mmu.WriteByte(valueAddr, expected)

	cpu.LDr_ffplusc(&cpu.R.A)

	assert.Equal(t, cpu.R.A, expected)
}

func TestLDr_ffpluscClockTimings(t *testing.T) {
	reset()
	cpu.LDr_ffplusc(&cpu.R.A)
	assert.Equal(t, cpu.LastInstrCycle.m, Word(2))
	assert.Equal(t, cpu.LastInstrCycle.t, Word(8))
}

//LD (C),r tests
//------------------------------------------
func TestLDffplusc_r(t *testing.T) {
	reset()
	var valueAddr Word = 0xFF03
	cpu.R.C = 0x03

	var expected byte = 0x05
	cpu.R.A = expected

	cpu.LDffplusc_r(&cpu.R.A)
	assert.Equal(t, cpu.mmu.ReadByte(valueAddr), expected)

}

func TestLDffplusc_rClockTimings(t *testing.T) {
	reset()
	cpu.LDffplusc_r(&cpu.R.A)
	assert.Equal(t, cpu.LastInstrCycle.m, Word(2))
	assert.Equal(t, cpu.LastInstrCycle.t, Word(8))
}

//-----------------------------------------------------------------------
//INSTRUCTIONS END

func TestReset(t *testing.T) {
	cpu.PC = 500
	cpu.R.B = 0x10
	cpu.SP = 0x4374
	cpu.R.H = 0x22
	cpu.Reset()
	assert.Equal(t, cpu.PC, ZEROW)
	assert.Equal(t, cpu.SP, ZEROW)
	assert.Equal(t, cpu.R.A, ZEROB)
	assert.Equal(t, cpu.R.B, ZEROB)
	assert.Equal(t, cpu.R.C, ZEROB)
	assert.Equal(t, cpu.R.D, ZEROB)
	assert.Equal(t, cpu.R.E, ZEROB)
	assert.Equal(t, cpu.R.F, ZEROB)
	assert.Equal(t, cpu.R.H, ZEROB)
	assert.Equal(t, cpu.R.L, ZEROB)
	assert.Equal(t, cpu.MachineCycles.m, ZEROW)
	assert.Equal(t, cpu.MachineCycles.t, ZEROW)
	assert.Equal(t, cpu.LastInstrCycle.m, ZEROW)
	assert.Equal(t, cpu.LastInstrCycle.t, ZEROW)
}

type MockMMU struct {
	memory map[Word]byte
}

func NewMockMMU() *MockMMU {
	var m *MockMMU = new(MockMMU)
	m.memory = make(map[Word]byte)
	return m
}

func (m *MockMMU) WriteByte(address Word, value byte) {
	m.memory[address] = value
}

func (m *MockMMU) WriteWord(address Word, value Word) {
}

func (m *MockMMU) ReadByte(address Word) byte {
	return m.memory[address]
}

func (m *MockMMU) ReadWord(address Word) Word {
	a, b := m.memory[address], m.memory[address+1]
	return (Word(a) << 8) ^ Word(b)
}
