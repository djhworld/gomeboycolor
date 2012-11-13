package cpu

import "testing"
import "github.com/stretchrcom/testify/assert"
import "github.com/djhworld/gomeboycolor/types"

const ZEROW types.Word = 0
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
	assert.Equal(t, cpu.PC, types.Word(0x03))
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

	//Check flag is set when result is carried
	cpu.AddA_r(&cpu.R.E)
	assert.Equal(t, (cpu.R.F & 0x10), byte(0x10))

	reset()
	cpu.R.A = 0xFA
	cpu.R.E = 0x02

	//Check flag is NOT set when result is not carried
	cpu.AddA_r(&cpu.R.E)
	assert.Equal(t, (cpu.R.F & 0x10), byte(0x00))
}

func TestAddA_rForHalfCarryFlag(t *testing.T) {
	reset()
	cpu.R.A = 0xA9
	cpu.R.E = 0x1F

	//Check flag is set when result is half carried
	cpu.AddA_r(&cpu.R.E)
	assert.Equal(t, (cpu.R.F & 0x20), byte(0x20))

	reset()
	cpu.R.A = 0xFA
	cpu.R.E = 0x02

	//Check flag is NOT set when result is not half carried
	cpu.AddA_r(&cpu.R.E)
	assert.Equal(t, (cpu.R.F & 0x20), byte(0x00))
}

func TestAddA_rForSubtractFlagReset(t *testing.T) {
	reset()
	cpu.R.F = 0x40
	cpu.R.A = 0xFA
	cpu.R.E = 0x02

	cpu.AddA_r(&cpu.R.E)
	assert.Equal(t, cpu.R.F, byte(0x00))

	reset()
	cpu.R.F = 0x70
	cpu.R.A = 0xFA
	cpu.R.E = 0x02

	cpu.AddA_r(&cpu.R.E)
	assert.Equal(t, cpu.R.F, byte(0x30))

}

func TestAddA_rClockTimings(t *testing.T) {
	reset()
	cpu.R.B = 0x10
	cpu.AddA_r(&cpu.R.B)

	assert.Equal(t, cpu.LastInstrCycle.m, types.Word(1))
	assert.Equal(t, cpu.LastInstrCycle.t, types.Word(4))
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
	var memoryAddr types.Word = 0x0302
	var H byte = 0x03
	var L byte = 0x02

	reset()
	cpu.R.A = 0xFE
	cpu.R.H = H
	cpu.R.L = L
	cpu.mmu.WriteByte(memoryAddr, 0x03)

	//Check flag is set when result overflows
	cpu.AddA_hl()
	assert.Equal(t, (cpu.R.F & 0x10), byte(0x10))

	reset()

	cpu.R.A = 0xFE
	cpu.R.H = H
	cpu.R.L = L
	cpu.mmu.WriteByte(memoryAddr, 0x01)

	//Check flag is set when result does not
	cpu.AddA_hl()
	assert.Equal(t, (cpu.R.F & 0x10), byte(0x00))
}

func TestAddA_hlClockTimings(t *testing.T) {
	reset()
	cpu.AddA_hl()
	assert.Equal(t, cpu.LastInstrCycle.m, types.Word(2))
	assert.Equal(t, cpu.LastInstrCycle.t, types.Word(8))

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
	assert.Equal(t, cpu.PC, types.Word(0x000F))

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
	assert.Equal(t, (cpu.R.F & 0x10), byte(0x10))

	reset()

	cpu.R.A = 0xFE
	cpu.PC = 0x000E
	cpu.mmu.WriteByte(0x000E, 0x01)

	//Check flag is set when result does not
	cpu.AddA_n()
	assert.Equal(t, (cpu.R.F & 0x10), byte(0x00))
}

func TestAddA_nClockTimings(t *testing.T) {
	reset()
	cpu.AddA_n()
	assert.Equal(t, cpu.LastInstrCycle.m, types.Word(2))
	assert.Equal(t, cpu.LastInstrCycle.t, types.Word(8))
}

//LD r,n tests
//------------------------------------------
func TestLDrn(t *testing.T) {
	var expected byte = 0x0A
	var addr types.Word = 0x0001

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
	var addr types.Word = 0x0001
	var expected types.Word = 0x0002
	reset()
	cpu.PC = addr
	cpu.LDrn(&cpu.R.B)

	assert.Equal(t, cpu.PC, expected)
}

func TestLDrn_ClockTimings(t *testing.T) {
	reset()
	cpu.LDrn(&cpu.R.B)
	assert.Equal(t, cpu.LastInstrCycle.m, types.Word(2))
	assert.Equal(t, cpu.LastInstrCycle.t, types.Word(8))
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
	assert.Equal(t, cpu.LastInstrCycle.m, types.Word(1))
	assert.Equal(t, cpu.LastInstrCycle.t, types.Word(4))
}

//LD r,(HL) tests
//------------------------------------------
func TestLDr_hl(t *testing.T) {
	var expected byte = 0xAE
	var addr types.Word = 0x1002

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
	assert.Equal(t, cpu.LastInstrCycle.m, types.Word(2))
	assert.Equal(t, cpu.LastInstrCycle.t, types.Word(8))
}

//LD (HL),r tests
//------------------------------------------
func TestLDhl_r(t *testing.T) {
	var expected byte = 0xF3
	var addr types.Word = 0x1002

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
	assert.Equal(t, cpu.LastInstrCycle.m, types.Word(2))
	assert.Equal(t, cpu.LastInstrCycle.t, types.Word(8))
}

//LD (HL),n tests
//------------------------------------------
func TestLDhl_n(t *testing.T) {
	var expected byte = 0x3A
	var addr types.Word = 0x0201
	var HL types.Word = 0xFFEE

	//test B
	reset()
	cpu.R.H, cpu.R.L = 0xFF, 0xEE
	cpu.mmu.WriteByte(addr, expected)
	cpu.PC = addr
	cpu.LDhl_n()
	assert.Equal(t, cpu.mmu.ReadByte(HL), expected)
}

func TestLDhl_nPCIncrementCheck(t *testing.T) {
	var addr types.Word = 0x0001
	var expected types.Word = 0x0002
	reset()
	cpu.PC = addr
	cpu.LDhl_n()

	assert.Equal(t, cpu.PC, expected)
}

func TestLDhl_nClockTimings(t *testing.T) {
	reset()
	cpu.LDhl_n()
	assert.Equal(t, cpu.LastInstrCycle.m, types.Word(3))
	assert.Equal(t, cpu.LastInstrCycle.t, types.Word(12))
}

//LD r,(BC) tests
//------------------------------------------
func TestLDr_bc(t *testing.T) {
	var expected byte = 0xAE
	var addr types.Word = 0x1002

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
	assert.Equal(t, cpu.LastInstrCycle.m, types.Word(2))
	assert.Equal(t, cpu.LastInstrCycle.t, types.Word(8))
}

//LD r,(DE) tests
//------------------------------------------
func TestLDr_de(t *testing.T) {
	var expected byte = 0xAE
	var addr types.Word = 0x1002

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
	assert.Equal(t, cpu.LastInstrCycle.m, types.Word(2))
	assert.Equal(t, cpu.LastInstrCycle.t, types.Word(8))
}

//LD r,nn tests
//------------------------------------------
func TestLDr_nn(t *testing.T) {
	var expected byte = 0x3E
	var addr types.Word = 0x0002
	var valueAddr types.Word = 0x3334
	cpu.mmu.WriteByte(addr, 0x33)
	cpu.mmu.WriteByte(addr+1, 0x34)
	cpu.mmu.WriteByte(valueAddr, expected)
	cpu.PC = addr

	cpu.LDr_nn(&cpu.R.A)
	assert.Equal(t, cpu.R.A, expected)
}

func TestLDr_nnPCIncremented(t *testing.T) {
	var addr types.Word = 0x0002
	var expected types.Word = 0x0004
	cpu.PC = addr
	cpu.LDr_nn(&cpu.R.A)
	assert.Equal(t, cpu.PC, expected)
}

func TestLDr_nnClockTimings(t *testing.T) {
	reset()
	cpu.LDr_nn(&cpu.R.A)
	assert.Equal(t, cpu.LastInstrCycle.m, types.Word(4))
	assert.Equal(t, cpu.LastInstrCycle.t, types.Word(16))
}

//LD (BC),r tests
//------------------------------------------
func TestLDbc_r(t *testing.T) {
	var expected byte = 0xF3
	var addr types.Word = 0x1002

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
	assert.Equal(t, cpu.LastInstrCycle.m, types.Word(2))
	assert.Equal(t, cpu.LastInstrCycle.t, types.Word(8))
}

//LD (DE),r tests
//------------------------------------------
func TestLDde_r(t *testing.T) {
	var expected byte = 0xF3
	var addr types.Word = 0x1002

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
	assert.Equal(t, cpu.LastInstrCycle.m, types.Word(2))
	assert.Equal(t, cpu.LastInstrCycle.t, types.Word(8))
}

//LD nn,r tests
//------------------------------------------
func TestLDnn_r(t *testing.T) {
	reset()
	var expected byte = 0x7C
	var addr types.Word = 0x0005
	var valueAddr types.Word = 0x3031

	cpu.mmu.WriteByte(addr, 0x30)
	cpu.mmu.WriteByte(addr+1, 0x31)

	cpu.PC = addr
	cpu.R.A = expected
	cpu.LDnn_r(&cpu.R.A)

	assert.Equal(t, cpu.mmu.ReadByte(valueAddr), expected)
}

func TestLDnn_rPCIncremented(t *testing.T) {
	reset()
	var addr types.Word = 0x0002
	var expected types.Word = 0x0004
	cpu.PC = addr
	cpu.LDnn_r(&cpu.R.A)
	assert.Equal(t, cpu.PC, expected)
}

func TestLDnn_rClockTimings(t *testing.T) {
	reset()
	cpu.LDnn_r(&cpu.R.A)

	assert.Equal(t, cpu.LastInstrCycle.m, types.Word(4))
	assert.Equal(t, cpu.LastInstrCycle.t, types.Word(16))
}

//LD r,(C) tests
//------------------------------------------
func TestLDr_ffplusc(t *testing.T) {
	reset()
	var valueAddr types.Word = 0xFF03
	var expected byte = 0x03
	cpu.R.C = expected

	cpu.mmu.WriteByte(valueAddr, expected)

	cpu.LDr_ffplusc(&cpu.R.A)

	assert.Equal(t, cpu.R.A, expected)
}

func TestLDr_ffpluscClockTimings(t *testing.T) {
	reset()
	cpu.LDr_ffplusc(&cpu.R.A)
	assert.Equal(t, cpu.LastInstrCycle.m, types.Word(2))
	assert.Equal(t, cpu.LastInstrCycle.t, types.Word(8))
}

//LD (C),r tests
//------------------------------------------
func TestLDffplusc_r(t *testing.T) {
	reset()
	var valueAddr types.Word = 0xFF03
	cpu.R.C = 0x03

	var expected byte = 0x05
	cpu.R.A = expected

	cpu.LDffplusc_r(&cpu.R.A)
	assert.Equal(t, cpu.mmu.ReadByte(valueAddr), expected)

}

func TestLDffplusc_rClockTimings(t *testing.T) {
	reset()
	cpu.LDffplusc_r(&cpu.R.A)
	assert.Equal(t, cpu.LastInstrCycle.m, types.Word(2))
	assert.Equal(t, cpu.LastInstrCycle.t, types.Word(8))
}

//LDD r, (HL) tests
//------------------------------------------
func TestLDDr_hl(t *testing.T) {
	reset()
	var valueAddr types.Word = 0x9965
	var expected byte = 0x2D
	cpu.R.H, cpu.R.L = 0x99, 0x65
	cpu.mmu.WriteByte(valueAddr, expected)

	cpu.LDDr_hl(&cpu.R.A)
	assert.Equal(t, cpu.R.A, expected)
	assert.Equal(t, cpu.R.H, byte(0x99))
	assert.Equal(t, cpu.R.L, byte(0x64))

	//Test that decrementing decrements H register when L = 0xFF
	reset()
	valueAddr = 0x9900
	expected = 0x2D
	cpu.R.H, cpu.R.L = 0x99, 0x00
	cpu.mmu.WriteByte(valueAddr, expected)

	cpu.LDDr_hl(&cpu.R.A)
	assert.Equal(t, cpu.R.A, expected)
	assert.Equal(t, cpu.R.H, byte(0x98))
	assert.Equal(t, cpu.R.L, byte(0xFF))
}

func TestLDDr_hlClockTimings(t *testing.T) {
	reset()
	cpu.LDDr_hl(&cpu.R.A)
	assert.Equal(t, cpu.LastInstrCycle.m, types.Word(2))
	assert.Equal(t, cpu.LastInstrCycle.t, types.Word(8))
}

//LDI r, (HL) tests
//------------------------------------------
func TestLDIr_hl(t *testing.T) {
	reset()
	var valueAddr types.Word = 0x9965
	var expected byte = 0x2D
	cpu.R.H, cpu.R.L = 0x99, 0x65
	cpu.mmu.WriteByte(valueAddr, expected)

	cpu.LDIr_hl(&cpu.R.A)
	assert.Equal(t, cpu.R.A, expected)
	assert.Equal(t, cpu.R.H, byte(0x99))
	assert.Equal(t, cpu.R.L, byte(0x66))

	//Test that H register increments when L = 0xFF
	reset()
	valueAddr = 0x99FF
	expected = 0x2D
	cpu.R.H, cpu.R.L = 0x99, 0xFF
	cpu.mmu.WriteByte(valueAddr, expected)

	cpu.LDIr_hl(&cpu.R.A)
	assert.Equal(t, cpu.R.A, expected)
	assert.Equal(t, cpu.R.H, byte(0x9A))
	assert.Equal(t, cpu.R.L, byte(0x00))
}

func TestLDIr_hlClockTimings(t *testing.T) {
	reset()
	cpu.LDIr_hl(&cpu.R.A)
	assert.Equal(t, cpu.LastInstrCycle.m, types.Word(2))
	assert.Equal(t, cpu.LastInstrCycle.t, types.Word(8))
}

//LDD (HL), r tests
//------------------------------------------
func TestLDDhl_r(t *testing.T) {
	reset()
	var valueAddr types.Word = 0x7531
	var expected byte = 0x9E
	cpu.R.H = 0x75
	cpu.R.L = 0x31
	cpu.R.A = expected

	cpu.LDDhl_r(&cpu.R.A)

	assert.Equal(t, cpu.mmu.ReadByte(valueAddr), expected)
	assert.Equal(t, cpu.R.H, byte(0x75))
	assert.Equal(t, cpu.R.L, byte(0x30))

	//Test that decrementing decrements H register when L = 0xFF
	valueAddr = 0x7500
	expected = 0x9E
	cpu.R.H = 0x75
	cpu.R.L = 0x00
	cpu.R.A = expected

	cpu.LDDhl_r(&cpu.R.A)

	assert.Equal(t, cpu.mmu.ReadByte(valueAddr), expected)
	assert.Equal(t, cpu.R.H, byte(0x74))
	assert.Equal(t, cpu.R.L, byte(0xFF))
}

func TestLDDhl_rClockTimings(t *testing.T) {
	reset()
	cpu.LDDhl_r(&cpu.R.A)
	assert.Equal(t, cpu.LastInstrCycle.m, types.Word(2))
	assert.Equal(t, cpu.LastInstrCycle.t, types.Word(8))
}

//LDI (HL), r tests
//------------------------------------------
func TestLDIhl_r(t *testing.T) {
	reset()
	var valueAddr types.Word = 0x7531
	var expected byte = 0x9E
	cpu.R.H = 0x75
	cpu.R.L = 0x31
	cpu.R.A = expected

	cpu.LDIhl_r(&cpu.R.A)

	assert.Equal(t, cpu.mmu.ReadByte(valueAddr), expected)
	assert.Equal(t, cpu.R.H, byte(0x75))
	assert.Equal(t, cpu.R.L, byte(0x32))

	//Test that decrementing decrements H register when L = 0xFF
	valueAddr = 0x75FF
	expected = 0x9E
	cpu.R.H = 0x75
	cpu.R.L = 0xFF
	cpu.R.A = expected

	cpu.LDIhl_r(&cpu.R.A)

	assert.Equal(t, cpu.mmu.ReadByte(valueAddr), expected)
	assert.Equal(t, cpu.R.H, byte(0x76))
	assert.Equal(t, cpu.R.L, byte(0x00))
}

func TestLDIhl_rClockTimings(t *testing.T) {
	reset()
	cpu.LDIhl_r(&cpu.R.A)
	assert.Equal(t, cpu.LastInstrCycle.m, types.Word(2))
	assert.Equal(t, cpu.LastInstrCycle.t, types.Word(8))
}

//LDH n, r tests
//------------------------------------------
func TestLDHn_r(t *testing.T) {
	reset()
	var expected byte = 0x77
	var valueAddr types.Word = 0xFF03
	cpu.PC = 0x0003
	cpu.mmu.WriteByte(valueAddr, expected)
	cpu.LDHn_r(&cpu.R.A)

	assert.Equal(t, cpu.R.A, expected)
}

func TestLDHn_rCheckPCIncremented(t *testing.T) {
	reset()
	var expected types.Word = 0x0002
	cpu.PC = 0x0001
	cpu.LDHn_r(&cpu.R.A)
	assert.Equal(t, cpu.PC, expected)
}

func TestLDHn_rClockTimings(t *testing.T) {
	reset()
	cpu.LDHn_r(&cpu.R.A)
	assert.Equal(t, cpu.LastInstrCycle.m, types.Word(3))
	assert.Equal(t, cpu.LastInstrCycle.t, types.Word(12))
}

//LDH r, n tests
//------------------------------------------
func TestLDHr_n(t *testing.T) {
	reset()
	var valueAddr types.Word = 0xFF03
	var expected byte = 0x4A

	cpu.PC = 0x0003
	cpu.mmu.WriteByte(valueAddr, expected)

	assert.Equal(t, cpu.mmu.ReadByte(valueAddr), expected)
}

func TestLDHr_nCheckPCIncremented(t *testing.T) {
	reset()
	var expected types.Word = 0x0002
	cpu.PC = 0x0001
	cpu.LDHr_n(&cpu.R.A)
	assert.Equal(t, cpu.PC, expected)

}

func TestLDHr_nClockTimings(t *testing.T) {
	reset()
	cpu.LDHr_n(&cpu.R.A)
	assert.Equal(t, cpu.LastInstrCycle.m, types.Word(3))
	assert.Equal(t, cpu.LastInstrCycle.t, types.Word(12))
}

//LD n,nn tests
//------------------------------------------
func TestLDn_nn(t *testing.T) {
	reset()
	var expected1 byte = 0xA3
	var expected2 byte = 0xF0
	cpu.PC = 0x0003
	cpu.mmu.WriteByte(0x0003, expected1)
	cpu.mmu.WriteByte(0x0004, expected2)

	cpu.LDn_nn(&cpu.R.B, &cpu.R.C)

	assert.Equal(t, cpu.R.B, expected1)
	assert.Equal(t, cpu.R.C, expected2)
}

func TestLDn_nnCheckPCIncremented(t *testing.T) {
	reset()
	var expected types.Word = 0x0003
	cpu.PC = 0x0001
	cpu.LDn_nn(&cpu.R.B, &cpu.R.C)
	assert.Equal(t, cpu.PC, expected)
}

func TestLDn_nnClockTimings(t *testing.T) {
	reset()
	cpu.LDn_nn(&cpu.R.B, &cpu.R.C)
	assert.Equal(t, cpu.LastInstrCycle.m, types.Word(3))
	assert.Equal(t, cpu.LastInstrCycle.t, types.Word(12))
}

//LD SP,nn tests
//------------------------------------------
func TestLDSP_nn(t *testing.T) {
	reset()
	var expected types.Word = 0xA3F0
	cpu.PC = 0x0003
	cpu.mmu.WriteWord(0x0003, expected)

	cpu.LDSP_nn()

	assert.Equal(t, cpu.SP, expected)
}

func TestLDSP_nnCheckPCIncremented(t *testing.T) {
	reset()
	var expected types.Word = 0x0003
	cpu.PC = 0x0001
	cpu.LDSP_nn()
	assert.Equal(t, cpu.PC, expected)
}

func TestLDSP_nnClockTimings(t *testing.T) {
	reset()
	cpu.LDSP_nn()
	assert.Equal(t, cpu.LastInstrCycle.m, types.Word(3))
	assert.Equal(t, cpu.LastInstrCycle.t, types.Word(12))
}

//LD SP, rr tests
//------------------------------------------
func TestLDSP_rr(t *testing.T) {
	reset()
	var expected types.Word = 0x3987
	cpu.R.H = 0x39
	cpu.R.L = 0x87

	cpu.LDSP_rr(&cpu.R.H, &cpu.R.L)

	assert.Equal(t, cpu.SP, expected)

}

func TestLDSP_rrClockTimings(t *testing.T) {
	reset()
	cpu.LDSP_rr(&cpu.R.H, &cpu.R.L)

	assert.Equal(t, cpu.LastInstrCycle.m, types.Word(2))
	assert.Equal(t, cpu.LastInstrCycle.t, types.Word(8))
}

//LDHL SP, n tests
//------------------------------------------
func TestLDHLSP_n(t *testing.T) {
	var n byte = 0x09
	var expectedH byte = 0x30
	var expectedL byte = 0x3C

	reset()
	cpu.SP = 0x3033
	cpu.PC = 0x0003
	cpu.mmu.WriteByte(0x0003, n)
	cpu.LDHLSP_n()

	assert.Equal(t, cpu.R.H, expectedH)
	assert.Equal(t, cpu.R.L, expectedL)
}

func TestLDHLSP_nPCIncremented(t *testing.T) {
	reset()

	var expected types.Word = 0x0003
	cpu.PC = 0x0002
	cpu.LDHLSP_n()

	assert.Equal(t, cpu.PC, expected)
}

func TestLDHLSP_nFlags(t *testing.T) {
	reset()

	//half carry flag
	var expected byte = 0x30
	cpu.R.F = 0x10
	cpu.PC = 0x0003
	cpu.SP = 0x0003
	cpu.mmu.WriteByte(0x0003, 0x9F)

	cpu.LDHLSP_n()
	assert.Equal(t, cpu.R.F, expected)

	//carry flag
	reset()
	expected = 0x10
	cpu.R.F = 0x00
	cpu.PC = 0x0003
	cpu.SP = 0xFFF3
	cpu.mmu.WriteByte(0x0003, 0x90)

	cpu.LDHLSP_n()
	assert.Equal(t, cpu.R.F, expected)
}

func TestLDHLSP_nClockTimings(t *testing.T) {
	reset()
	cpu.LDHLSP_n()

	assert.Equal(t, cpu.LastInstrCycle.m, types.Word(3))
	assert.Equal(t, cpu.LastInstrCycle.t, types.Word(12))
}

//LD nn, SP tests
//------------------------------------------
func TestLDnn_SP(t *testing.T) {
	reset()
	var expectedSP types.Word = 0x3AF2
	cpu.PC = 0x0009
	cpu.mmu.WriteWord(0x0009, 0xF003)
	cpu.SP = expectedSP
	cpu.LDnn_SP()
	assert.Equal(t, cpu.mmu.ReadWord(0xF003), expectedSP)

}

func TestLDnn_SPCheckPCIncremented(t *testing.T) {
	reset()
	var expected types.Word = 0x000B
	cpu.PC = 0x0009
	cpu.LDnn_SP()
	assert.Equal(t, cpu.PC, expected)
}

func TestLDnn_SPClockTimings(t *testing.T) {
	reset()
	cpu.LDnn_SP()
	assert.Equal(t, cpu.LastInstrCycle.m, byte(5))
	assert.Equal(t, cpu.LastInstrCycle.t, byte(20))
}

//NOP
//------------------------------------------
func TestNOP(t *testing.T) {
	reset()
	cpu.NOP()
	assert.Equal(t, cpu.LastInstrCycle.m, types.Word(1))
	assert.Equal(t, cpu.LastInstrCycle.t, types.Word(4))
}

//PUSH nn tests
//------------------------------------------
func TestPushnn(t *testing.T) {
	var expectedSP types.Word = 0x0002
	reset()
	cpu.SP = 0x0004
	cpu.R.B = 0xFA
	cpu.R.C = 0xD9

	cpu.Push_nn(&cpu.R.B, &cpu.R.C)
	assert.Equal(t, cpu.SP, expectedSP)
	assert.Equal(t, cpu.mmu.ReadByte(0x0003), cpu.R.B)
	assert.Equal(t, cpu.mmu.ReadByte(0x0002), cpu.R.C)
}

func TestPushnnClockTimings(t *testing.T) {
	reset()
	cpu.Push_nn(&cpu.R.B, &cpu.R.C)
	assert.Equal(t, cpu.LastInstrCycle.m, byte(3))
	assert.Equal(t, cpu.LastInstrCycle.t, byte(12))
}

//POP nn tests
//------------------------------------------
func TestPopnn(t *testing.T) {
	var expectedSP types.Word = 0x0004
	reset()
	cpu.SP = 0x0004
	cpu.R.B = 0xFA
	cpu.R.C = 0xD9

	cpu.Push_nn(&cpu.R.B, &cpu.R.C)
	cpu.Pop_nn(&cpu.R.H, &cpu.R.L)

	assert.Equal(t, cpu.SP, expectedSP)
	assert.Equal(t, cpu.R.H, cpu.R.B)
	assert.Equal(t, cpu.R.L, cpu.R.C)
}

func TestPopnnClockTimings(t *testing.T) {
	reset()
	cpu.Pop_nn(&cpu.R.B, &cpu.R.C)
	assert.Equal(t, cpu.LastInstrCycle.m, byte(3))
	assert.Equal(t, cpu.LastInstrCycle.t, byte(12))
}

//ADDC A,r tests
//------------------------------------------
func TestADDCA_r(t *testing.T) {
	var expectedA byte = 0x03

	//with carry flag set
	reset()
	cpu.R.A = 0x03
	cpu.R.B = 0xFF
	cpu.R.F = 0x10
	cpu.AddCA_r(&cpu.R.B)
	assert.Equal(t, cpu.R.A, expectedA)

	//with carry flag not set
	expectedA = 0x02
	reset()
	cpu.R.A = 0x03
	cpu.R.B = 0xFF
	cpu.AddCA_r(&cpu.R.B)
	assert.Equal(t, cpu.R.A, expectedA)
}

func TestADDCA_rZeroFlagSet(t *testing.T) {
	//with carry flag set
	reset()
	cpu.R.A = 0x00
	cpu.R.B = 0x00
	cpu.R.F = 0x10
	cpu.AddCA_r(&cpu.R.B)
	assert.Equal(t, cpu.IsFlagSet(Z), false)

	//with carry flag set
	reset()
	cpu.R.A = 0x00
	cpu.R.B = 0x00
	cpu.R.F = 0x00
	cpu.AddCA_r(&cpu.R.B)
	assert.Equal(t, cpu.IsFlagSet(Z), true)
}

func TestADDCA_rCarryFlagSet(t *testing.T) {
	//with carry flag set
	reset()
	cpu.R.A = 0x0f
	cpu.R.B = 0x33
	cpu.R.F = 0x10
	cpu.AddCA_r(&cpu.R.B)
	assert.Equal(t, cpu.IsFlagSet(C), true)

	//with carry flag not set
	reset()
	cpu.R.A = 0xFE
	cpu.R.B = 0x33
	cpu.R.F = 0x00
	cpu.AddCA_r(&cpu.R.B)
	assert.Equal(t, cpu.IsFlagSet(C), true)
}

func TestADDCA_rHalfCarryFlagSet(t *testing.T) {
	//with carry flag set
	reset()
	cpu.R.A = 0x0f
	cpu.R.B = 0x33
	cpu.R.F = 0x10
	cpu.AddCA_r(&cpu.R.B)
	assert.Equal(t, cpu.IsFlagSet(H), true)

	//with carry flag not set
	reset()
	cpu.R.A = 0x01
	cpu.R.B = 0x33
	cpu.R.F = 0x00
	cpu.AddCA_r(&cpu.R.B)
	assert.Equal(t, cpu.IsFlagSet(H), false)

}

func TestADDCA_rSubtractFlagReset(t *testing.T) {
	//with carry flag set
	reset()
	cpu.AddCA_r(&cpu.R.B)
	assert.Equal(t, cpu.IsFlagSet(N), false)
}

func TestADDCA_rClockTimings(t *testing.T) {
	reset()
	cpu.AddCA_r(&cpu.R.B)
	assert.Equal(t, cpu.LastInstrCycle.m, byte(1))
	assert.Equal(t, cpu.LastInstrCycle.t, byte(4))
}

//ADDC A,(HL) tests
//------------------------------------------
func TestADDCA_hl(t *testing.T) {
	var expectedA byte = 0x05

	//with carry flag set
	reset()
	cpu.R.F = 0x10
	cpu.R.A = 0x03
	cpu.R.H = 0xFF
	cpu.R.L = 0x11
	cpu.mmu.WriteByte(0xFF11, 0x01)
	cpu.AddCA_hl()
	assert.Equal(t, cpu.R.A, expectedA)

	//with carry flag not set
	reset()
	expectedA = 0x04
	cpu.R.A = 0x03
	cpu.R.H = 0xFF
	cpu.R.L = 0x11
	cpu.mmu.WriteByte(0xFF11, 0x01)
	cpu.AddCA_hl()
	assert.Equal(t, cpu.R.A, expectedA)
}

func TestADDCA_hlZeroFlagSet(t *testing.T) {
	//with carry flag set
	reset()
	cpu.R.F = 0x10
	cpu.R.A = 0x00
	cpu.R.H = 0xFF
	cpu.R.L = 0x11
	cpu.mmu.WriteByte(0xFF11, 0x00)
	cpu.AddCA_hl()
	assert.Equal(t, cpu.IsFlagSet(Z), false)

	//with carry flag not set
	reset()
	cpu.R.A = 0x00
	cpu.R.H = 0xFF
	cpu.R.L = 0x11
	cpu.mmu.WriteByte(0xFF11, 0x00)
	cpu.AddCA_hl()
	assert.Equal(t, cpu.IsFlagSet(Z), true)

}

func TestADDCA_hlCarryFlagSet(t *testing.T) {
	//with carry flag set
	reset()
	cpu.R.F = 0x10
	cpu.R.A = 0xFC
	cpu.R.H = 0xFF
	cpu.R.L = 0x11
	cpu.mmu.WriteByte(0xFF11, 0x02)
	cpu.AddCA_hl()
	assert.Equal(t, cpu.IsFlagSet(C), true)

	//with carry flag not set
	reset()
	cpu.R.A = 0xFC
	cpu.R.H = 0xFF
	cpu.R.L = 0x11
	cpu.mmu.WriteByte(0xFF11, 0x06)
	cpu.AddCA_hl()
	assert.Equal(t, cpu.IsFlagSet(C), true)

}

func TestADDCA_hlHalfCarryFlagSet(t *testing.T) {
	//with carry flag set
	reset()
	cpu.R.F = 0x10
	cpu.R.A = 0xF0
	cpu.R.H = 0xFF
	cpu.R.L = 0x11
	cpu.mmu.WriteByte(0xFF11, 0x03)
	cpu.AddCA_hl()
	assert.Equal(t, cpu.IsFlagSet(H), false)

	//with carry flag not set
	reset()
	cpu.R.A = 0xA9
	cpu.R.H = 0xFF
	cpu.R.L = 0x11
	cpu.mmu.WriteByte(0xFF11, 0x07)
	cpu.AddCA_hl()
	assert.Equal(t, cpu.IsFlagSet(H), true)

}

func TestADDCA_hlSubtractFlagReset(t *testing.T) {
	reset()
	cpu.R.H = 0xFF
	cpu.R.L = 0x11
	cpu.mmu.WriteByte(0xFF11, 0x03)
	cpu.AddCA_hl()
	assert.Equal(t, cpu.IsFlagSet(N), false)
}

func TestADDCA_hlClockTimings(t *testing.T) {
	reset()
	cpu.AddCA_hl()
	assert.Equal(t, cpu.LastInstrCycle.m, byte(2))
	assert.Equal(t, cpu.LastInstrCycle.t, byte(8))
}

//ADDC A,n tests
//------------------------------------------
func TestADDCA_n(t *testing.T) {
	var expectedA byte = 0x05

	//with carry flag set
	reset()
	cpu.R.F = 0x10
	cpu.R.A = 0x03
	cpu.PC = 0x0001
	cpu.mmu.WriteByte(0x0001, 0x01)
	cpu.AddCA_n()
	assert.Equal(t, cpu.R.A, expectedA)

	//with carry flag not set
	reset()
	expectedA = 0x04
	cpu.R.A = 0x03
	cpu.PC = 0x0001
	cpu.mmu.WriteByte(0x0001, 0x01)
	cpu.AddCA_n()
	assert.Equal(t, cpu.R.A, expectedA)

}

func TestADDCA_nZeroFlagSet(t *testing.T) {
	//with carry flag set
	reset()
	cpu.R.F = 0x10
	cpu.R.A = 0x00
	cpu.PC = 0x0001
	cpu.mmu.WriteByte(0x0001, 0x00)
	cpu.AddCA_n()
	assert.Equal(t, cpu.IsFlagSet(Z), false)

	//with carry flag not set
	reset()
	cpu.R.A = 0x00
	cpu.PC = 0x0001
	cpu.mmu.WriteByte(0x0001, 0x00)
	cpu.AddCA_n()
	assert.Equal(t, cpu.IsFlagSet(Z), true)

}

func TestADDCA_nCarryFlagSet(t *testing.T) {
	//with carry flag set
	reset()
	cpu.R.F = 0x10
	cpu.R.A = 0xFC
	cpu.PC = 0x0001
	cpu.mmu.WriteByte(0x0001, 0x01)
	cpu.AddCA_n()
	assert.Equal(t, cpu.IsFlagSet(C), true)

	//with carry flag not set
	reset()
	cpu.R.A = 0xFE
	cpu.PC = 0x0001
	cpu.mmu.WriteByte(0x0001, 0x0F)
	cpu.AddCA_n()
	assert.Equal(t, cpu.IsFlagSet(C), true)

}

func TestADDCA_nHalfCarryFlagSet(t *testing.T) {
	//with carry flag set
	reset()
	cpu.R.F = 0x10
	cpu.R.A = 0x89
	cpu.PC = 0x0001
	cpu.mmu.WriteByte(0x0001, 0x01)
	cpu.AddCA_n()
	assert.Equal(t, cpu.IsFlagSet(H), false)

	//with carry flag not set
	reset()
	cpu.R.A = 0x89
	cpu.PC = 0x0001
	cpu.mmu.WriteByte(0x0001, 0x0F)
	cpu.AddCA_n()
	assert.Equal(t, cpu.IsFlagSet(H), true)
}

func TestADDCA_nSubtractFlagReset(t *testing.T) {
	reset()
	cpu.R.F = 0x10
	cpu.R.A = 0x89
	cpu.PC = 0x0001
	cpu.mmu.WriteByte(0x0001, 0x01)
	cpu.AddCA_n()
	assert.Equal(t, cpu.IsFlagSet(N), false)

}

func TestADDCA_nPCIncremented(t *testing.T) {
	reset()
	var expectedPC types.Word = 0x0002
	cpu.PC = 0x0001
	cpu.AddCA_n()
	assert.Equal(t, cpu.PC, expectedPC)
}

func TestADDCA_nClockTimings(t *testing.T) {
	reset()
	cpu.AddCA_n()
	assert.Equal(t, cpu.LastInstrCycle.m, byte(2))
	assert.Equal(t, cpu.LastInstrCycle.t, byte(8))
}

// NOTE: Just going to encorporate all checks and things in one test from now it, it's becoming a PITA to write 4-5 tests per instruction!

//SUB A,r tests 
//------------------------------------------
func TestSUBA_r(t *testing.T) {
	reset()
	var expectedA byte

	//Check subtraction works
	expectedA = 0x03
	cpu.R.A = 0x05
	cpu.R.B = 0x02
	cpu.SubA_r(&cpu.R.B)
	assert.Equal(t, cpu.R.A, expectedA)
	//Check timings are correct
	assert.Equal(t, cpu.LastInstrCycle.m, byte(1))
	assert.Equal(t, cpu.LastInstrCycle.t, byte(4))
	//Check N flag is set 
	assert.Equal(t, cpu.IsFlagSet(N), true, "Subract flag (N) is not set!")
	//Check other flags are not set 
	assert.Equal(t, cpu.IsFlagSet(Z), false, "Zero flag (Z) should not be set!")
	assert.Equal(t, cpu.IsFlagSet(H), false, "Half Carry flag (H) should not be set!")
	assert.Equal(t, cpu.IsFlagSet(C), false, "Carry flag (C) should not be set!")

	//Check zero flag is set
	reset()
	expectedA = 0x00
	cpu.R.A = 0x05
	cpu.R.B = 0x05
	cpu.SubA_r(&cpu.R.B)
	assert.Equal(t, cpu.R.A, expectedA)
	assert.Equal(t, cpu.IsFlagSet(Z), true, "Zero flag (Z) is not set!")
	assert.Equal(t, cpu.IsFlagSet(N), true, "Subract flag (N) is not set!")

	//Check half carry flag is set
	reset()
	expectedA = 0xe8
	cpu.R.A = 0xf1
	cpu.R.B = 0x09
	cpu.SubA_r(&cpu.R.B)
	assert.Equal(t, cpu.R.A, expectedA)
	assert.Equal(t, cpu.IsFlagSet(H), true, "Half Carry flag (H) is not set!")
	assert.Equal(t, cpu.IsFlagSet(N), true, "Subract flag (N) is not set!")
	assert.Equal(t, cpu.IsFlagSet(C), false, "Carry flag (C) should not be set!")
	assert.Equal(t, cpu.IsFlagSet(Z), false, "Zero flag (Z) should not be set!")

	//Check carry flag is set
	reset()
	expectedA = 0xfe
	cpu.R.A = 0x15
	cpu.R.B = 0x17
	cpu.SubA_r(&cpu.R.B)
	assert.Equal(t, cpu.R.A, expectedA)
	assert.Equal(t, cpu.IsFlagSet(C), true, "Carry flag (C) is not set!")
	assert.Equal(t, cpu.IsFlagSet(N), true, "Subract flag (N) is not set!")
	assert.Equal(t, cpu.IsFlagSet(Z), false, "Zero flag (Z) should not be set!")
	assert.Equal(t, cpu.IsFlagSet(H), true, "Half Carry flag (H) should not be set!")
}

//SUB A,(HL) tests 
//------------------------------------------
func TestSUBA_hl(t *testing.T) {
	reset()
	var expectedA byte

	//Check subtraction works
	expectedA = 0x03
	cpu.R.A = 0x05
	cpu.R.H = 0x01
	cpu.R.L = 0x02
	cpu.mmu.WriteByte(0x0102, 0x02)
	cpu.SubA_hl()
	assert.Equal(t, cpu.R.A, expectedA)
	//Check timings are correct
	assert.Equal(t, cpu.LastInstrCycle.m, byte(2))
	assert.Equal(t, cpu.LastInstrCycle.t, byte(8))
	//Check N flag is set 
	assert.Equal(t, cpu.IsFlagSet(N), true, "Subract flag (N) is not set!")
	//Check other flags are not set 
	assert.Equal(t, cpu.IsFlagSet(Z), false, "Zero flag (Z) should not be set!")
	assert.Equal(t, cpu.IsFlagSet(H), false, "Half Carry flag (H) should not be set!")
	assert.Equal(t, cpu.IsFlagSet(C), false, "Carry flag (C) should not be set!")

	//Check zero flag is set
	reset()
	expectedA = 0x00
	cpu.R.A = 0x05
	cpu.R.H = 0x01
	cpu.R.L = 0x02
	cpu.mmu.WriteByte(0x0102, 0x05)
	cpu.SubA_hl()
	assert.Equal(t, cpu.R.A, expectedA)
	assert.Equal(t, cpu.IsFlagSet(Z), true, "Zero flag (Z) is not set!")
	assert.Equal(t, cpu.IsFlagSet(N), true, "Subract flag (N) is not set!")

	//Check half carry flag is set
	reset()
	expectedA = 0xe8
	cpu.R.A = 0xf1
	cpu.R.H = 0x01
	cpu.R.L = 0x02
	cpu.mmu.WriteByte(0x0102, 0x09)
	cpu.SubA_hl()
	assert.Equal(t, cpu.R.A, expectedA)
	assert.Equal(t, cpu.IsFlagSet(H), true, "Half Carry flag (H) is not set!")
	assert.Equal(t, cpu.IsFlagSet(N), true, "Subract flag (N) is not set!")
	assert.Equal(t, cpu.IsFlagSet(C), false, "Carry flag (C) should not be set!")
	assert.Equal(t, cpu.IsFlagSet(Z), false, "Zero flag (Z) should not be set!")

	//Check carry flag is set
	reset()
	expectedA = 0xfe
	cpu.R.A = 0x15
	cpu.R.H = 0x01
	cpu.R.L = 0x02
	cpu.mmu.WriteByte(0x0102, 0x17)
	cpu.SubA_hl()
	assert.Equal(t, cpu.R.A, expectedA)
	assert.Equal(t, cpu.IsFlagSet(C), true, "Carry flag (C) is not set!")
	assert.Equal(t, cpu.IsFlagSet(N), true, "Subract flag (N) is not set!")
	assert.Equal(t, cpu.IsFlagSet(Z), false, "Zero flag (Z) should not be set!")
	assert.Equal(t, cpu.IsFlagSet(H), true, "Half Carry flag (H) should not be set!")
}

//SUB A,n tests 
//------------------------------------------
func TestSUBA_n(t *testing.T) {
	var expectedA byte
	var expectedPC types.Word

	//Check subtraction works
	reset()
	expectedA = 0x03
	expectedPC = 0x0002
	cpu.R.A = 0x05
	cpu.PC = 0x0001
	cpu.mmu.WriteByte(cpu.PC, 0x02)
	cpu.SubA_n()
	assert.Equal(t, cpu.R.A, expectedA)

	//check PC has been incremented
	assert.Equal(t, cpu.PC, expectedPC)
	//Check timings are correct
	assert.Equal(t, cpu.LastInstrCycle.m, byte(2))
	assert.Equal(t, cpu.LastInstrCycle.t, byte(8))
	//Check N flag is set 
	assert.Equal(t, cpu.IsFlagSet(N), true, "Subract flag (N) is not set!")
	//Check other flags are not set 
	assert.Equal(t, cpu.IsFlagSet(Z), false, "Zero flag (Z) should not be set!")
	assert.Equal(t, cpu.IsFlagSet(H), false, "Half Carry flag (H) should not be set!")
	assert.Equal(t, cpu.IsFlagSet(C), false, "Carry flag (C) should not be set!")

	//Check zero flag is set
	reset()
	expectedA = 0x00
	cpu.R.A = 0x05
	cpu.PC = 0x0001
	cpu.mmu.WriteByte(cpu.PC, 0x05)
	cpu.SubA_n()
	assert.Equal(t, cpu.R.A, expectedA)
	assert.Equal(t, cpu.IsFlagSet(Z), true, "Zero flag (Z) is not set!")
	assert.Equal(t, cpu.IsFlagSet(N), true, "Subract flag (N) is not set!")

	//Check half carry flag is set
	reset()
	expectedA = 0xe8
	cpu.R.A = 0xf1
	cpu.PC = 0x0001
	cpu.mmu.WriteByte(cpu.PC, 0x09)
	cpu.SubA_n()
	assert.Equal(t, cpu.R.A, expectedA)
	assert.Equal(t, cpu.IsFlagSet(H), true, "Half Carry flag (H) is not set!")
	assert.Equal(t, cpu.IsFlagSet(N), true, "Subract flag (N) is not set!")
	assert.Equal(t, cpu.IsFlagSet(C), false, "Carry flag (C) should not be set!")
	assert.Equal(t, cpu.IsFlagSet(Z), false, "Zero flag (Z) should not be set!")

	//Check carry flag is set
	reset()
	expectedA = 0xfe
	cpu.R.A = 0x15
	cpu.PC = 0x0001
	cpu.mmu.WriteByte(cpu.PC, 0x17)
	cpu.SubA_n()
	assert.Equal(t, cpu.R.A, expectedA)
	assert.Equal(t, cpu.IsFlagSet(C), true, "Carry flag (C) is not set!")
	assert.Equal(t, cpu.IsFlagSet(N), true, "Subract flag (N) is not set!")
	assert.Equal(t, cpu.IsFlagSet(Z), false, "Zero flag (Z) should not be set!")
	assert.Equal(t, cpu.IsFlagSet(H), true, "Half Carry flag (H) should not be set!")
}

//AND A,r tests 
//------------------------------------------
func TestAndA_r(t *testing.T) {
	reset()
	var expectedA byte

	//test instruction
	expectedA = 0x0A
	cpu.R.A = 0xAA
	cpu.R.B = 0x0F
	cpu.AndA_r(&cpu.R.B)
	assert.Equal(t, cpu.R.A, expectedA)
	//Check timings are correct
	assert.Equal(t, cpu.LastInstrCycle.m, byte(1))
	assert.Equal(t, cpu.LastInstrCycle.t, byte(4))
	//H flag should be set
	assert.Equal(t, cpu.IsFlagSet(H), true)
	//result should not be zero so Z flag should not be set
	assert.Equal(t, cpu.IsFlagSet(Z), false)
	//N and C flags should not be set (or should be reset)
	assert.Equal(t, cpu.IsFlagSet(N), false)
	assert.Equal(t, cpu.IsFlagSet(C), false)

	//test zero flag set
	reset()
	expectedA = 0x00
	cpu.R.A = 0xAA
	cpu.R.B = 0x00
	cpu.AndA_r(&cpu.R.B)
	assert.Equal(t, cpu.R.A, expectedA)
	assert.Equal(t, cpu.IsFlagSet(Z), true)
	assert.Equal(t, cpu.IsFlagSet(H), true)
	assert.Equal(t, cpu.IsFlagSet(N), false)
	assert.Equal(t, cpu.IsFlagSet(C), false)

}

//AND A,(HL) tests 
//------------------------------------------
func TestAndA_hl(t *testing.T) {
	reset()
	var expectedA byte

	//test instruction
	expectedA = 0x0A
	cpu.R.A = 0xAA
	cpu.R.H = 0x01
	cpu.R.L = 0x02
	cpu.mmu.WriteByte(0x0102, 0x0F)
	cpu.AndA_hl()
	assert.Equal(t, cpu.R.A, expectedA)
	//Check timings are correct
	assert.Equal(t, cpu.LastInstrCycle.m, byte(2))
	assert.Equal(t, cpu.LastInstrCycle.t, byte(8))
	//H flag should be set
	assert.Equal(t, cpu.IsFlagSet(H), true)
	//result should not be zero so Z flag should not be set
	assert.Equal(t, cpu.IsFlagSet(Z), false)
	//N and C flags should not be set (or should be reset)
	assert.Equal(t, cpu.IsFlagSet(N), false)
	assert.Equal(t, cpu.IsFlagSet(C), false)

	//test zero flag set
	reset()
	expectedA = 0x00
	cpu.R.A = 0xAA
	cpu.R.H = 0x01
	cpu.R.L = 0x02
	cpu.mmu.WriteByte(0x0102, 0x00)
	cpu.AndA_hl()
	assert.Equal(t, cpu.R.A, expectedA)
	assert.Equal(t, cpu.IsFlagSet(Z), true)
	assert.Equal(t, cpu.IsFlagSet(H), true)
	assert.Equal(t, cpu.IsFlagSet(N), false)
	assert.Equal(t, cpu.IsFlagSet(C), false)

}

//AND A, n tests 
//------------------------------------------
func TestAndA_n(t *testing.T) {
	reset()
	var expectedA byte
	var expectedPC types.Word

	//test instruction
	expectedA = 0x0A
	expectedPC = 0x0103
	cpu.R.A = 0xAA
	cpu.PC = 0x0102
	cpu.mmu.WriteByte(cpu.PC, 0x0F)
	cpu.AndA_n()
	assert.Equal(t, cpu.R.A, expectedA)
	//check PC incremented 
	assert.Equal(t, cpu.PC, expectedPC)
	//Check timings are correct
	assert.Equal(t, cpu.LastInstrCycle.m, byte(2))
	assert.Equal(t, cpu.LastInstrCycle.t, byte(8))
	//H flag should be set
	assert.Equal(t, cpu.IsFlagSet(H), true)
	//result should not be zero so Z flag should not be set
	assert.Equal(t, cpu.IsFlagSet(Z), false)
	//N and C flags should not be set (or should be reset)
	assert.Equal(t, cpu.IsFlagSet(N), false)
	assert.Equal(t, cpu.IsFlagSet(C), false)

	//test zero flag set
	reset()
	expectedA = 0x00
	cpu.R.A = 0xAA
	cpu.PC = 0x0102
	cpu.mmu.WriteByte(cpu.PC, 0x00)
	cpu.AndA_n()
	assert.Equal(t, cpu.R.A, expectedA)
	assert.Equal(t, cpu.IsFlagSet(Z), true)
	assert.Equal(t, cpu.IsFlagSet(H), true)
	assert.Equal(t, cpu.IsFlagSet(N), false)
	assert.Equal(t, cpu.IsFlagSet(C), false)
}

//OR A,r tests 
//------------------------------------------
func TestOrA_r(t *testing.T) {
	reset()
	var expectedA byte

	//test instruction
	expectedA = 0xAF
	cpu.R.A = 0xAA
	cpu.R.B = 0x0F
	cpu.OrA_r(&cpu.R.B)
	assert.Equal(t, cpu.R.A, expectedA)
	//Check timings are correct
	assert.Equal(t, cpu.LastInstrCycle.m, byte(1))
	assert.Equal(t, cpu.LastInstrCycle.t, byte(4))
	//all flags should be disabled
	assert.Equal(t, cpu.IsFlagSet(H), false)
	assert.Equal(t, cpu.IsFlagSet(Z), false)
	assert.Equal(t, cpu.IsFlagSet(N), false)
	assert.Equal(t, cpu.IsFlagSet(C), false)

	//test zero flag set
	reset()
	expectedA = 0x00
	cpu.R.A = 0x00
	cpu.R.B = 0x00
	cpu.OrA_r(&cpu.R.B)
	assert.Equal(t, cpu.R.A, expectedA)
	assert.Equal(t, cpu.IsFlagSet(Z), true)
	assert.Equal(t, cpu.IsFlagSet(H), false)
	assert.Equal(t, cpu.IsFlagSet(N), false)
	assert.Equal(t, cpu.IsFlagSet(C), false)

}

//OR A,(HL) tests 
//------------------------------------------
func TestOrA_hl(t *testing.T) {
	reset()
	var expectedA byte

	//test instruction
	expectedA = 0xAF
	cpu.R.A = 0xAA
	cpu.R.H = 0x01
	cpu.R.L = 0x02
	cpu.mmu.WriteByte(0x0102, 0x0F)
	cpu.OrA_hl()
	assert.Equal(t, cpu.R.A, expectedA)
	//Check timings are correct
	assert.Equal(t, cpu.LastInstrCycle.m, byte(2))
	assert.Equal(t, cpu.LastInstrCycle.t, byte(8))
	//all flags should be disabled
	assert.Equal(t, cpu.IsFlagSet(H), false)
	assert.Equal(t, cpu.IsFlagSet(Z), false)
	assert.Equal(t, cpu.IsFlagSet(N), false)
	assert.Equal(t, cpu.IsFlagSet(C), false)

	//test zero flag set
	reset()
	expectedA = 0x00
	cpu.R.A = 0x00
	cpu.R.H = 0x01
	cpu.R.L = 0x02
	cpu.mmu.WriteByte(0x0102, 0x00)
	cpu.OrA_hl()
	assert.Equal(t, cpu.R.A, expectedA)
	assert.Equal(t, cpu.IsFlagSet(Z), true)
	assert.Equal(t, cpu.IsFlagSet(H), false)
	assert.Equal(t, cpu.IsFlagSet(N), false)
	assert.Equal(t, cpu.IsFlagSet(C), false)

}

//OR A, n tests 
//------------------------------------------
func TestOrA_n(t *testing.T) {
	reset()
	var expectedA byte
	var expectedPC types.Word

	//test instruction
	expectedA = 0xAF
	expectedPC = 0x0103
	cpu.R.A = 0xAA
	cpu.PC = 0x0102
	cpu.mmu.WriteByte(cpu.PC, 0x0F)
	cpu.OrA_n()
	assert.Equal(t, cpu.R.A, expectedA)
	//check PC incremented 
	assert.Equal(t, cpu.PC, expectedPC)
	//Check timings are correct
	assert.Equal(t, cpu.LastInstrCycle.m, byte(2))
	assert.Equal(t, cpu.LastInstrCycle.t, byte(8))
	//all flags should be disabled
	assert.Equal(t, cpu.IsFlagSet(H), false)
	assert.Equal(t, cpu.IsFlagSet(Z), false)
	assert.Equal(t, cpu.IsFlagSet(N), false)
	assert.Equal(t, cpu.IsFlagSet(C), false)

	//test zero flag set
	reset()
	expectedA = 0x00
	cpu.R.A = 0x00
	cpu.PC = 0x0102
	cpu.mmu.WriteByte(cpu.PC, 0x00)
	cpu.OrA_n()
	assert.Equal(t, cpu.R.A, expectedA)
	assert.Equal(t, cpu.IsFlagSet(Z), true)
	assert.Equal(t, cpu.IsFlagSet(H), false)
	assert.Equal(t, cpu.IsFlagSet(N), false)
	assert.Equal(t, cpu.IsFlagSet(C), false)
}

//XOR A,r tests 
//------------------------------------------
func TestXorA_r(t *testing.T) {
	reset()
	var expectedA byte

	//test instruction
	expectedA = 0xA5
	cpu.R.A = 0xAA
	cpu.R.B = 0x0F
	cpu.XorA_r(&cpu.R.B)
	assert.Equal(t, cpu.R.A, expectedA)
	//Check timings are correct
	assert.Equal(t, cpu.LastInstrCycle.m, byte(1))
	assert.Equal(t, cpu.LastInstrCycle.t, byte(4))
	//all flags should be disabled
	assert.Equal(t, cpu.IsFlagSet(H), false)
	assert.Equal(t, cpu.IsFlagSet(Z), false)
	assert.Equal(t, cpu.IsFlagSet(N), false)
	assert.Equal(t, cpu.IsFlagSet(C), false)

	//test zero flag set
	reset()
	expectedA = 0x00
	cpu.R.A = 0x00
	cpu.R.B = 0x00
	cpu.XorA_r(&cpu.R.B)
	assert.Equal(t, cpu.R.A, expectedA)
	assert.Equal(t, cpu.IsFlagSet(Z), true)
	assert.Equal(t, cpu.IsFlagSet(H), false)
	assert.Equal(t, cpu.IsFlagSet(N), false)
	assert.Equal(t, cpu.IsFlagSet(C), false)

}

//XOR A,(HL) tests 
//------------------------------------------
func TestXorA_hl(t *testing.T) {
	reset()
	var expectedA byte

	//test instruction
	expectedA = 0xA5
	cpu.R.A = 0xAA
	cpu.R.H = 0x01
	cpu.R.L = 0x02
	cpu.mmu.WriteByte(0x0102, 0x0F)
	cpu.XorA_hl()
	assert.Equal(t, cpu.R.A, expectedA)
	//Check timings are correct
	assert.Equal(t, cpu.LastInstrCycle.m, byte(2))
	assert.Equal(t, cpu.LastInstrCycle.t, byte(8))
	//all flags should be disabled
	assert.Equal(t, cpu.IsFlagSet(H), false)
	assert.Equal(t, cpu.IsFlagSet(Z), false)
	assert.Equal(t, cpu.IsFlagSet(N), false)
	assert.Equal(t, cpu.IsFlagSet(C), false)

	//test zero flag set
	reset()
	expectedA = 0x00
	cpu.R.A = 0x00
	cpu.R.H = 0x01
	cpu.R.L = 0x02
	cpu.mmu.WriteByte(0x0102, 0x00)
	cpu.XorA_hl()
	assert.Equal(t, cpu.R.A, expectedA)
	assert.Equal(t, cpu.IsFlagSet(Z), true)
	assert.Equal(t, cpu.IsFlagSet(H), false)
	assert.Equal(t, cpu.IsFlagSet(N), false)
	assert.Equal(t, cpu.IsFlagSet(C), false)

}

//XOR A, n tests 
//------------------------------------------
func TestXorA_n(t *testing.T) {
	reset()
	var expectedA byte
	var expectedPC types.Word

	//test instruction
	expectedA = 0xA5
	expectedPC = 0x0103
	cpu.R.A = 0xAA
	cpu.PC = 0x0102
	cpu.mmu.WriteByte(cpu.PC, 0x0F)
	cpu.XorA_n()
	assert.Equal(t, cpu.R.A, expectedA)
	//check PC incremented 
	assert.Equal(t, cpu.PC, expectedPC)
	//Check timings are correct
	assert.Equal(t, cpu.LastInstrCycle.m, byte(2))
	assert.Equal(t, cpu.LastInstrCycle.t, byte(8))
	//all flags should be disabled
	assert.Equal(t, cpu.IsFlagSet(H), false)
	assert.Equal(t, cpu.IsFlagSet(Z), false)
	assert.Equal(t, cpu.IsFlagSet(N), false)
	assert.Equal(t, cpu.IsFlagSet(C), false)

	//test zero flag set
	reset()
	expectedA = 0x00
	cpu.R.A = 0x00
	cpu.PC = 0x0102
	cpu.mmu.WriteByte(cpu.PC, 0x00)
	cpu.XorA_n()
	assert.Equal(t, cpu.R.A, expectedA)
	assert.Equal(t, cpu.IsFlagSet(Z), true)
	assert.Equal(t, cpu.IsFlagSet(H), false)
	assert.Equal(t, cpu.IsFlagSet(N), false)
	assert.Equal(t, cpu.IsFlagSet(C), false)
}

//CP A, r tests 
//------------------------------------------
func TestCPA_r(t *testing.T) {

	//check for same (zero flag)
	reset()
	cpu.R.A = 0x05
	cpu.R.B = 0x05
	cpu.CPA_r(&cpu.R.B)
	//Check timings are correct
	assert.Equal(t, cpu.LastInstrCycle.m, byte(1))
	assert.Equal(t, cpu.LastInstrCycle.t, byte(4))
	assert.Equal(t, cpu.IsFlagSet(Z), true)
	assert.Equal(t, cpu.IsFlagSet(N), true)
	assert.Equal(t, cpu.IsFlagSet(H), false)
	assert.Equal(t, cpu.IsFlagSet(C), false)

	//check for carry
	reset()
	cpu.R.A = 0x05
	cpu.R.B = 0xAA
	cpu.CPA_r(&cpu.R.B)
	//Check timings are correct
	assert.Equal(t, cpu.IsFlagSet(Z), false)
	assert.Equal(t, cpu.IsFlagSet(N), true)
	assert.Equal(t, cpu.IsFlagSet(H), false)
	assert.Equal(t, cpu.IsFlagSet(C), true)
}

//CP A, (HL) tests 
//------------------------------------------
func TestCPA_hl(t *testing.T) {
	reset()
	cpu.R.A = 0x05
	cpu.R.H = 0x03
	cpu.R.L = 0xAA
	cpu.mmu.WriteByte(0x03AA, 0x05)
	cpu.CPA_hl()

	//Check timings are correct
	assert.Equal(t, cpu.LastInstrCycle.m, byte(2))
	assert.Equal(t, cpu.LastInstrCycle.t, byte(8))
	assert.Equal(t, cpu.IsFlagSet(Z), true)
	assert.Equal(t, cpu.IsFlagSet(N), true)
	assert.Equal(t, cpu.IsFlagSet(H), false)
	assert.Equal(t, cpu.IsFlagSet(C), false)

	//check for carry
	reset()
	cpu.R.A = 0x05
	cpu.R.H = 0x03
	cpu.R.L = 0xAA
	cpu.mmu.WriteByte(0x03AA, 0xAA)
	cpu.CPA_hl()
	//Check timings are correct
	assert.Equal(t, cpu.IsFlagSet(Z), false)
	assert.Equal(t, cpu.IsFlagSet(N), true)
	assert.Equal(t, cpu.IsFlagSet(H), false)
	assert.Equal(t, cpu.IsFlagSet(C), true)
}

//CP A, n tests 
//------------------------------------------
func TestCPA_n(t *testing.T) {
	var expectedPC types.Word = 0x0002
	reset()
	cpu.R.A = 0x05
	cpu.PC = 0x0001
	cpu.mmu.WriteByte(cpu.PC, 0x05)
	cpu.CPA_n()

	//Check PC incremented
	assert.Equal(t, cpu.PC, expectedPC)
	//Check timings are correct
	assert.Equal(t, cpu.LastInstrCycle.m, byte(2))
	assert.Equal(t, cpu.LastInstrCycle.t, byte(8))
	assert.Equal(t, cpu.IsFlagSet(Z), true)
	assert.Equal(t, cpu.IsFlagSet(N), true)
	assert.Equal(t, cpu.IsFlagSet(H), false)
	assert.Equal(t, cpu.IsFlagSet(C), false)

	//check for carry
	reset()
	cpu.R.A = 0x05
	cpu.PC = 0x0001
	cpu.mmu.WriteByte(cpu.PC, 0xAA)
	cpu.CPA_n()
	//Check timings are correct
	assert.Equal(t, cpu.IsFlagSet(Z), false)
	assert.Equal(t, cpu.IsFlagSet(N), true)
	assert.Equal(t, cpu.IsFlagSet(H), false)
	assert.Equal(t, cpu.IsFlagSet(C), true)
}

//INC r tests 
//------------------------------------------
func TestInc_r(t *testing.T) {
	var expectedB byte = 0x02
	reset()
	//set N flag as instruction should reset it
	cpu.SetFlag(N)

	cpu.R.B = 0x01
	cpu.Inc_r(&cpu.R.B)
	assert.Equal(t, cpu.R.B, expectedB)
	assert.Equal(t, cpu.IsFlagSet(Z), false)
	assert.Equal(t, cpu.IsFlagSet(H), false)
	assert.Equal(t, cpu.IsFlagSet(N), false)
	assert.Equal(t, cpu.IsFlagSet(C), false)
	//Check timings are correct
	assert.Equal(t, cpu.LastInstrCycle.m, byte(1))
	assert.Equal(t, cpu.LastInstrCycle.t, byte(4))

	reset()
	cpu.R.B = 0xFF
	cpu.Inc_r(&cpu.R.B)
	//check zero flag
	assert.Equal(t, cpu.IsFlagSet(Z), true)

	reset()
	cpu.R.B = 0x0F
	cpu.Inc_r(&cpu.R.B)
	//check half carry flag
	assert.Equal(t, cpu.IsFlagSet(H), true)

}

//INC (HL) tests 
//------------------------------------------
func TestInc_hl(t *testing.T) {
	var expectedVal byte = 0x02
	reset()
	//set N flag as instruction should reset it
	cpu.SetFlag(N)

	cpu.R.H = 0x03
	cpu.R.L = 0xFF
	cpu.mmu.WriteByte(0x03FF, 0x01)
	cpu.Inc_hl()
	assert.Equal(t, cpu.mmu.ReadByte(0x03FF), expectedVal)
	assert.Equal(t, cpu.IsFlagSet(Z), false)
	assert.Equal(t, cpu.IsFlagSet(H), false)
	assert.Equal(t, cpu.IsFlagSet(N), false)
	assert.Equal(t, cpu.IsFlagSet(C), false)
	//Check timings are correct
	assert.Equal(t, cpu.LastInstrCycle.m, byte(3))
	assert.Equal(t, cpu.LastInstrCycle.t, byte(12))

	reset()
	cpu.R.H = 0x03
	cpu.R.L = 0xFF
	cpu.mmu.WriteByte(0x03FF, 0xFF)
	cpu.Inc_hl()
	//check zero flag
	assert.Equal(t, cpu.IsFlagSet(Z), true)

	reset()
	cpu.R.H = 0x03
	cpu.R.L = 0xFF
	cpu.mmu.WriteByte(0x03FF, 0x0F)
	cpu.Inc_hl()
	//check half carry flag
	assert.Equal(t, cpu.IsFlagSet(H), true)
}

//DEC r tests 
//------------------------------------------
func TestDec_r(t *testing.T) {
	var expectedB byte = 0x01
	reset()
	//set C flag as instruction should not do anything to it
	cpu.SetFlag(C)

	cpu.R.B = 0x02
	cpu.Dec_r(&cpu.R.B)
	assert.Equal(t, cpu.R.B, expectedB)
	assert.Equal(t, cpu.IsFlagSet(N), true)
	assert.Equal(t, cpu.IsFlagSet(C), true)
	assert.Equal(t, cpu.IsFlagSet(Z), false)
	assert.Equal(t, cpu.IsFlagSet(H), false)
	//Check timings are correct
	assert.Equal(t, cpu.LastInstrCycle.m, byte(1))
	assert.Equal(t, cpu.LastInstrCycle.t, byte(4))

	reset()
	cpu.R.B = 0x01
	cpu.Dec_r(&cpu.R.B)
	//check zero flag
	assert.Equal(t, cpu.IsFlagSet(Z), true)
}

//DEC (HL) tests 
//------------------------------------------
func TestDec_hl(t *testing.T) {
	var expectedVal byte = 0x01
	reset()

	//set C flag as instruction should not do anything to it
	cpu.SetFlag(C)

	cpu.R.H = 0x03
	cpu.R.L = 0xFF
	cpu.mmu.WriteByte(0x03FF, 0x02)
	cpu.Dec_hl()
	assert.Equal(t, cpu.mmu.ReadByte(0x03FF), expectedVal)
	assert.Equal(t, cpu.IsFlagSet(Z), false)
	assert.Equal(t, cpu.IsFlagSet(H), false)
	assert.Equal(t, cpu.IsFlagSet(N), true)
	assert.Equal(t, cpu.IsFlagSet(C), true)
	//Check timings are correct
	assert.Equal(t, cpu.LastInstrCycle.m, byte(3))
	assert.Equal(t, cpu.LastInstrCycle.t, byte(12))

	reset()
	cpu.R.H = 0x03
	cpu.R.L = 0xFF
	cpu.mmu.WriteByte(0x03FF, 0x01)
	cpu.Dec_hl()
	//check zero flag
	assert.Equal(t, cpu.IsFlagSet(Z), true)
}

//ADD HL,rr tests 
//------------------------------------------
func TestAddhl_rr(t *testing.T) {
	reset()
	var expectedH byte = 0x41
	var expectedL byte = 0x03

	cpu.SetFlag(N)
	cpu.R.H = 0x01
	cpu.R.L = 0x02
	cpu.R.B = 0x40
	cpu.R.C = 0x01

	cpu.Addhl_rr(&cpu.R.B, &cpu.R.C)
	assert.Equal(t, cpu.R.H, expectedH)
	assert.Equal(t, cpu.R.L, expectedL)
	//Check N flag is reset
	assert.Equal(t, cpu.IsFlagSet(N), false)
	//Check timings are correct
	assert.Equal(t, cpu.LastInstrCycle.m, byte(2))
	assert.Equal(t, cpu.LastInstrCycle.t, byte(8))

	//carry flag
	reset()
	expectedH = 0x00
	expectedL = 0xFD

	cpu.R.H = 0xFF
	cpu.R.L = 0xFE
	cpu.R.B = 0x00
	cpu.R.C = 0xFF

	cpu.Addhl_rr(&cpu.R.B, &cpu.R.C)
	assert.Equal(t, cpu.R.H, expectedH)
	assert.Equal(t, cpu.R.L, expectedL)
	//Check N flag is reset
	assert.Equal(t, cpu.IsFlagSet(N), false)
	//Check carry flag is set
	assert.Equal(t, cpu.IsFlagSet(C), true)

}

//ADD HL,SP tests 
//------------------------------------------
func TestAddhl_sp(t *testing.T) {
	reset()
	var expectedH byte = 0x41
	var expectedL byte = 0x03

	cpu.SetFlag(N)
	cpu.R.H = 0x01
	cpu.R.L = 0x02
	cpu.SP = 0x4001

	cpu.Addhl_sp()
	assert.Equal(t, cpu.R.H, expectedH)
	assert.Equal(t, cpu.R.L, expectedL)
	//Check N flag is reset
	assert.Equal(t, cpu.IsFlagSet(N), false)
	//Check timings are correct
	assert.Equal(t, cpu.LastInstrCycle.m, byte(2))
	assert.Equal(t, cpu.LastInstrCycle.t, byte(8))

	//carry flag
	reset()
	expectedH = 0x00
	expectedL = 0xFD

	cpu.R.H = 0xFF
	cpu.R.L = 0xFE
	cpu.SP = 0x00FF

	cpu.Addhl_sp()
	assert.Equal(t, cpu.R.H, expectedH)
	assert.Equal(t, cpu.R.L, expectedL)
	//Check N flag is reset
	assert.Equal(t, cpu.IsFlagSet(N), false)
	//Check carry flag is set
	assert.Equal(t, cpu.IsFlagSet(C), true)
}

//ADD SP,n tests 
//------------------------------------------
func TestAddsp_n(t *testing.T) {
	var expectedSP types.Word = 0x0006
	var expectedPC types.Word = 0x1001
	reset()
	cpu.SetFlag(Z)
	cpu.SetFlag(N)
	cpu.PC = 0x1000
	cpu.SP = 0x0003
	cpu.mmu.WriteByte(cpu.PC, 0x0003)

	cpu.Addsp_n()
	assert.Equal(t, cpu.SP, expectedSP)
	//Check Z and N flags are reset
	assert.Equal(t, cpu.IsFlagSet(N), false)
	assert.Equal(t, cpu.IsFlagSet(Z), false)
	//Check timings are correct
	assert.Equal(t, cpu.LastInstrCycle.m, byte(4))
	assert.Equal(t, cpu.LastInstrCycle.t, byte(16))
	//Check PC incremented
	assert.Equal(t, cpu.PC, expectedPC)

	//check carry flag
	expectedSP = 0x0006
	reset()
	cpu.SP = 0xFFFE
	cpu.mmu.WriteByte(cpu.PC, 0x0008)

	cpu.Addsp_n()
	assert.Equal(t, cpu.SP, expectedSP)
	//Check carry flag is set
	assert.Equal(t, cpu.IsFlagSet(C), true)

	//Check timings are correct
	assert.Equal(t, cpu.LastInstrCycle.m, byte(4))
	assert.Equal(t, cpu.LastInstrCycle.t, byte(16))
}

//INC rr tests 
//------------------------------------------
func TestInc_rr(t *testing.T) {
	var expectedH byte = 0x03
	var expectedL byte = 0xFF

	reset()
	cpu.R.H = 0x03
	cpu.R.L = 0xFE

	cpu.Inc_rr(&cpu.R.H, &cpu.R.L)

	assert.Equal(t, expectedH, cpu.R.H)
	assert.Equal(t, expectedL, cpu.R.L)

	//Check timings are correct
	assert.Equal(t, cpu.LastInstrCycle.m, byte(2))
	assert.Equal(t, cpu.LastInstrCycle.t, byte(8))
}

//INC SP tests 
//------------------------------------------
func TestInc_sp(t *testing.T) {
	var expectedSP types.Word = 0x03FF

	reset()
	cpu.SP = 0x03FE
	cpu.Inc_sp()

	assert.Equal(t, expectedSP, cpu.SP)

	//Check timings are correct
	assert.Equal(t, cpu.LastInstrCycle.m, byte(2))
	assert.Equal(t, cpu.LastInstrCycle.t, byte(8))
}

//DEC rr tests 
//------------------------------------------
func TestDec_rr(t *testing.T) {
	var expectedH byte = 0x03
	var expectedL byte = 0xFD

	reset()
	cpu.R.H = 0x03
	cpu.R.L = 0xFE

	cpu.Dec_rr(&cpu.R.H, &cpu.R.L)

	assert.Equal(t, expectedH, cpu.R.H)
	assert.Equal(t, expectedL, cpu.R.L)

	//Check timings are correct
	assert.Equal(t, cpu.LastInstrCycle.m, byte(2))
	assert.Equal(t, cpu.LastInstrCycle.t, byte(8))
}

//DEC SP tests 
//------------------------------------------
func TestDec_sp(t *testing.T) {
	var expectedSP types.Word = 0x03FD

	reset()
	cpu.SP = 0x03FE
	cpu.Dec_sp()

	assert.Equal(t, expectedSP, cpu.SP)

	//Check timings are correct
	assert.Equal(t, cpu.LastInstrCycle.m, byte(2))
	assert.Equal(t, cpu.LastInstrCycle.t, byte(8))
}

//CPL Tests
func TestCPL(t *testing.T) {
	var expectedA byte = 0x55
	reset()
	cpu.R.A = 0xAA
	cpu.CPL()

	assert.Equal(t, cpu.R.A, expectedA)
	assert.Equal(t, cpu.IsFlagSet(N), true)
	assert.Equal(t, cpu.IsFlagSet(H), true)

	//Check timings are correct
	assert.Equal(t, cpu.LastInstrCycle.m, byte(1))
	assert.Equal(t, cpu.LastInstrCycle.t, byte(4))

	expectedA = 0x00
	reset()
	cpu.R.A = 0xFF
	cpu.CPL()
	assert.Equal(t, cpu.R.A, expectedA)

	expectedA = 0xFF
	reset()
	cpu.R.A = 0x00
	cpu.CPL()
	assert.Equal(t, cpu.R.A, expectedA)
}

//CCF Tests
func TestCCF(t *testing.T) {
	reset()
	cpu.SetFlag(C)
	cpu.CCF()
	assert.Equal(t, cpu.IsFlagSet(C), false)

	reset()
	cpu.CCF()
	assert.Equal(t, cpu.IsFlagSet(C), true)

	//ensure n and h flags are reset
	reset()
	cpu.SetFlag(N)
	cpu.SetFlag(H)
	cpu.CCF()
	assert.Equal(t, cpu.IsFlagSet(N), false)
	assert.Equal(t, cpu.IsFlagSet(H), false)

	//Check timings are correct
	assert.Equal(t, cpu.LastInstrCycle.m, byte(1))
	assert.Equal(t, cpu.LastInstrCycle.t, byte(4))

}

//SCF Tests
func TestSCF(t *testing.T) {
	reset()
	cpu.SCF()
	assert.Equal(t, cpu.IsFlagSet(C), true)

	//ensure n and h flags are reset
	reset()
	cpu.SetFlag(N)
	cpu.SetFlag(H)
	cpu.SCF()
	assert.Equal(t, cpu.IsFlagSet(N), false)
	assert.Equal(t, cpu.IsFlagSet(H), false)

	//Check timings are correct
	assert.Equal(t, cpu.LastInstrCycle.m, byte(1))
	assert.Equal(t, cpu.LastInstrCycle.t, byte(4))
}

//SWAP r Tests
func TestSwap_r(t *testing.T) {
	var expectedA byte = 0xBF
	reset()
	cpu.R.A = 0xFB
	cpu.SetFlag(N)
	cpu.SetFlag(H)
	cpu.SetFlag(C)
	cpu.Swap_r(&cpu.R.A)

	assert.Equal(t, cpu.R.A, expectedA)
	assert.Equal(t, cpu.IsFlagSet(Z), false)
	//ensure flags are reset
	assert.Equal(t, cpu.IsFlagSet(N), false)
	assert.Equal(t, cpu.IsFlagSet(H), false)
	assert.Equal(t, cpu.IsFlagSet(C), false)

	//Check timings are correct
	assert.Equal(t, cpu.LastInstrCycle.m, byte(2))
	assert.Equal(t, cpu.LastInstrCycle.t, byte(8))

	//check zero flag
	expectedA = 0x00
	reset()
	cpu.R.A = 0x00
	cpu.Swap_r(&cpu.R.A)

	assert.Equal(t, cpu.R.A, expectedA)
	assert.Equal(t, cpu.IsFlagSet(Z), true)
}

//SWAP (HL) Tests
func TestSwap_hl(t *testing.T) {
	var addr types.Word = 0x1001
	var expected byte = 0xBF
	reset()
	cpu.SetFlag(N)
	cpu.SetFlag(H)
	cpu.SetFlag(C)
	cpu.R.H = 0x10
	cpu.R.L = 0x01
	cpu.mmu.WriteByte(addr, 0xFB)
	cpu.Swap_hl()

	assert.Equal(t, cpu.mmu.ReadByte(addr), expected)
	assert.Equal(t, cpu.IsFlagSet(Z), false)
	//ensure flags are reset
	assert.Equal(t, cpu.IsFlagSet(N), false)
	assert.Equal(t, cpu.IsFlagSet(H), false)
	assert.Equal(t, cpu.IsFlagSet(C), false)

	//Check timings are correct
	assert.Equal(t, cpu.LastInstrCycle.m, byte(4))
	assert.Equal(t, cpu.LastInstrCycle.t, byte(16))

	addr = 0x1001
	expected = 0x00
	reset()
	cpu.R.H = 0x10
	cpu.R.L = 0x01
	cpu.mmu.WriteByte(addr, 0x00)
	cpu.Swap_hl()

	assert.Equal(t, cpu.mmu.ReadByte(addr), expected)
	assert.Equal(t, cpu.IsFlagSet(Z), true)
}

//RLCA tests
func TestRLCA(t *testing.T) {
	var expectedA byte = 0x08
	reset()
	cpu.SetFlag(N)
	cpu.SetFlag(H)
	cpu.R.A = 0x04
	cpu.RLCA()

	assert.Equal(t, cpu.R.A, expectedA)
	assert.Equal(t, cpu.IsFlagSet(Z), false)
	assert.Equal(t, cpu.IsFlagSet(C), false)
	//ensure flags are reset
	assert.Equal(t, cpu.IsFlagSet(N), false)
	assert.Equal(t, cpu.IsFlagSet(H), false)

	//Check timings are correct
	assert.Equal(t, cpu.LastInstrCycle.m, byte(1))
	assert.Equal(t, cpu.LastInstrCycle.t, byte(4))

	//check carry flag
	reset()
	cpu.R.A = 0x89
	cpu.RLCA()
	assert.Equal(t, cpu.IsFlagSet(C), true)

	//check zero flag
	reset()
	cpu.R.A = 0x00
	cpu.RLCA()
	assert.Equal(t, cpu.IsFlagSet(Z), true)
}

//RLA tests
func TestRLA(t *testing.T) {
	var expectedA byte = 0x08
	reset()
	cpu.SetFlag(N)
	cpu.SetFlag(H)
	cpu.R.A = 0x04
	cpu.RLA()

	assert.Equal(t, cpu.R.A, expectedA)
	assert.Equal(t, cpu.IsFlagSet(Z), false)
	assert.Equal(t, cpu.IsFlagSet(C), false)
	//ensure flags are reset
	assert.Equal(t, cpu.IsFlagSet(N), false)
	assert.Equal(t, cpu.IsFlagSet(H), false)

	//Check timings are correct
	assert.Equal(t, cpu.LastInstrCycle.m, byte(1))
	assert.Equal(t, cpu.LastInstrCycle.t, byte(4))

	//check carry flag
	reset()
	cpu.R.A = 0x89
	cpu.RLA()
	assert.Equal(t, cpu.IsFlagSet(C), true)

	//check carry flag with carry flag already set
	reset()
	expectedA = 0x85
	cpu.SetFlag(C)
	cpu.R.A = 0x42
	cpu.RLA()
	assert.Equal(t, cpu.R.A, expectedA)
	assert.Equal(t, cpu.IsFlagSet(C), false)

	//check zero flag
	reset()
	cpu.R.A = 0x00
	cpu.RLA()
	assert.Equal(t, cpu.IsFlagSet(Z), true)
}

//RRCA tests
func TestRRCA(t *testing.T) {
	var expectedA byte = 0x02
	reset()
	cpu.SetFlag(N)
	cpu.SetFlag(H)
	cpu.R.A = 0x04
	cpu.RRCA()

	assert.Equal(t, cpu.R.A, expectedA)
	assert.Equal(t, cpu.IsFlagSet(Z), false)
	assert.Equal(t, cpu.IsFlagSet(C), false)
	//ensure flags are reset
	assert.Equal(t, cpu.IsFlagSet(N), false)
	assert.Equal(t, cpu.IsFlagSet(H), false)

	//Check timings are correct
	assert.Equal(t, cpu.LastInstrCycle.m, byte(1))
	assert.Equal(t, cpu.LastInstrCycle.t, byte(4))

	//check carry flag
	reset()
	cpu.R.A = 0x33
	cpu.RRCA()
	assert.Equal(t, cpu.IsFlagSet(C), true)

	//check zero flag
	reset()
	cpu.R.A = 0x00
	cpu.RRCA()
	assert.Equal(t, cpu.IsFlagSet(Z), true)
}

//RRA tests
func TestRRA(t *testing.T) {
	var expectedA byte = 0x40
	reset()
	cpu.SetFlag(N)
	cpu.SetFlag(H)
	cpu.R.A = 0x80
	cpu.RRA()

	assert.Equal(t, cpu.R.A, expectedA)
	assert.Equal(t, cpu.IsFlagSet(Z), false)
	assert.Equal(t, cpu.IsFlagSet(C), false)
	//ensure flags are reset
	assert.Equal(t, cpu.IsFlagSet(N), false)
	assert.Equal(t, cpu.IsFlagSet(H), false)

	//Check timings are correct
	assert.Equal(t, cpu.LastInstrCycle.m, byte(1))
	assert.Equal(t, cpu.LastInstrCycle.t, byte(4))

	//check carry flag
	reset()
	cpu.R.A = 0x01
	cpu.RRA()
	assert.Equal(t, cpu.IsFlagSet(C), true)

	//check carry flag with carry flag already set
	reset()
	expectedA = 0xA2
	cpu.SetFlag(C)
	cpu.R.A = 0x44
	cpu.RRA()
	assert.Equal(t, cpu.R.A, expectedA)
	assert.Equal(t, cpu.IsFlagSet(C), false)

	//check zero flag
	reset()
	cpu.R.A = 0x00
	cpu.RRA()
	assert.Equal(t, cpu.IsFlagSet(Z), true)
}

//BIT b,r tests
func TestBitb_r(t *testing.T) {
	var expectedPC types.Word = 0x0002
	reset()
	cpu.PC = 0x0001
	cpu.R.A = 0x31
	cpu.SetFlag(N)
	cpu.mmu.WriteByte(cpu.PC, 0x01)
	cpu.Bitb_r(&cpu.R.A)

	assert.Equal(t, cpu.PC, expectedPC)
	//ensure flag is set
	assert.Equal(t, cpu.IsFlagSet(H), true)
	//ensure flag reset
	assert.Equal(t, cpu.IsFlagSet(N), false)
	//Check timings are correct
	assert.Equal(t, cpu.LastInstrCycle.m, byte(2))
	assert.Equal(t, cpu.LastInstrCycle.t, byte(8))

	var expected []bool = []bool{true, true, false, false, true, true, true, false}

	var results []bool = make([]bool, 8)
	var j byte = 0x07
	for i := 0x00; i <= 0x07; i++ {
		reset()
		cpu.PC = 0x0001
		cpu.R.A = 0x31
		cpu.SetFlag(N)
		cpu.mmu.WriteByte(cpu.PC, byte(i))
		cpu.Bitb_r(&cpu.R.A)
		results[j] = cpu.IsFlagSet(Z)
		j--
	}
	assert.Equal(t, expected, results)
}

//BIT b,(HL) tests
func TestBitb_hl(t *testing.T) {
	var expectedPC types.Word = 0x0002
	var addr types.Word = 0x3044
	reset()
	cpu.PC = 0x0001
	cpu.R.H = 0x30
	cpu.R.L = 0x44
	cpu.SetFlag(N)
	cpu.mmu.WriteByte(addr, 0x21)
	cpu.mmu.WriteByte(cpu.PC, 0x01)
	cpu.Bitb_hl()

	assert.Equal(t, cpu.PC, expectedPC)
	//ensure flag is set
	assert.Equal(t, cpu.IsFlagSet(H), true)
	//ensure flag reset
	assert.Equal(t, cpu.IsFlagSet(N), false)
	//Check timings are correct
	assert.Equal(t, cpu.LastInstrCycle.m, byte(4))
	assert.Equal(t, cpu.LastInstrCycle.t, byte(16))

	var expected []bool = []bool{true, true, false, false, true, true, true, false}

	var results []bool = make([]bool, 8)
	var j byte = 0x07
	for i := 0x00; i <= 0x07; i++ {
		reset()
		cpu.PC = 0x0001
		cpu.R.H = 0x30
		cpu.R.L = 0x44
		cpu.SetFlag(N)
		cpu.mmu.WriteByte(addr, 0x31)
		cpu.mmu.WriteByte(cpu.PC, byte(i))
		cpu.Bitb_hl()
		results[j] = cpu.IsFlagSet(Z)
		j--
	}

	assert.Equal(t, expected, results)
}

//JP nn tests
func TestJP_nn(t *testing.T) {
	var expectedPC types.Word = 0x5672
	reset()
	cpu.PC = 0x2001
	cpu.mmu.WriteWord(cpu.PC, expectedPC)
	cpu.JP_nn()

	assert.Equal(t, cpu.PC, expectedPC)

	//Check timings are correct
	assert.Equal(t, cpu.LastInstrCycle.m, byte(3))
	assert.Equal(t, cpu.LastInstrCycle.t, byte(12))
}

//JP (HL) tests
func TestJP_hl(t *testing.T) {
	var addr types.Word = 0x2001
	var expectedPC types.Word = 0x5672
	reset()
	cpu.R.H = 0x20
	cpu.R.L = 0x01
	cpu.mmu.WriteWord(addr, expectedPC)
	cpu.JP_hl()

	assert.Equal(t, cpu.PC, expectedPC)

	//Check timings are correct
	assert.Equal(t, cpu.LastInstrCycle.m, byte(1))
	assert.Equal(t, cpu.LastInstrCycle.t, byte(4))
}

//JP cc nn tests
func TestJPcc_nn(t *testing.T) {
	var expectedPC types.Word = 0x5672
	reset()
	cpu.PC = 0x2001
	cpu.mmu.WriteWord(cpu.PC, expectedPC)
	cpu.SetFlag(Z)
	cpu.JPcc_nn(Z, true)
	assert.Equal(t, cpu.PC, expectedPC)

	expectedPC = 0x2001
	reset()
	cpu.PC = 0x2001
	cpu.mmu.WriteWord(cpu.PC, expectedPC)
	cpu.SetFlag(Z)
	cpu.JPcc_nn(Z, false)
	assert.Equal(t, cpu.PC, expectedPC)

	expectedPC = 0x5672
	reset()
	cpu.PC = 0x2001
	cpu.mmu.WriteWord(cpu.PC, expectedPC)
	cpu.JPcc_nn(Z, false)
	assert.Equal(t, cpu.PC, expectedPC)

	//Check timings are correct
	assert.Equal(t, cpu.LastInstrCycle.m, byte(3))
	assert.Equal(t, cpu.LastInstrCycle.t, byte(12))
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

func (m *MockMMU) LoadROM(startAddr types.Word, rt types.ROMType, data []byte) (bool, error) {
	return true, nil
}

func (m *MockMMU) SetInBootMode(mode bool) {
}
