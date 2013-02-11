/**
 * Created with IntelliJ IDEA.
 * User: danielharper
 * Date: 29/01/2013
 * Time: 20:16
 * To change this template use File | Settings | File Templates.
 */
package cpu

/*
import (
	"cartridge"
	"github.com/stretchrcom/testify/assert"
	"testing"
	"types"
	"utils"
)

var NoOperandsInstr Instruction = Instruction{0x00, "Test instruction", 0, 0, [2]byte{}}
var OneOperandsInstr Instruction = Instruction{0x00, "Test instruction", 1, 0, [2]byte{}}
var TwoOperandsInstr Instruction = Instruction{0x00, "Test instruction", 2, 0, [2]byte{}}

//OR A r
func TestOrA_r(t *testing.T) {
	cpu := NewCPU()
	cpu.R.A = 0x1A
	cpu.R.B = 0x0F

	cpu.OrA_r(&cpu.R.B)

	var expected byte = 0x1F
	assert.Equal(t, cpu.R.A, expected)
}

func TestOrA_rFlagsReset(t *testing.T) {
	cpu := NewCPU()
	cpu.R.A = 0x1A
	cpu.R.B = 0x0F

	cpu.SetFlag(N)
	cpu.SetFlag(H)
	cpu.SetFlag(C)
	cpu.OrA_r(&cpu.R.B)

	assert.False(t, cpu.IsFlagSet(N))
	assert.False(t, cpu.IsFlagSet(H))
	assert.False(t, cpu.IsFlagSet(C))
}

func TestOrA_rZeroFlag(t *testing.T) {
	cpu := NewCPU()
	cpu.R.A = 0x00

	assert.False(t, cpu.IsFlagSet(Z))
	cpu.OrA_r(&cpu.R.A)

	assert.True(t, cpu.IsFlagSet(Z))
}

//OR A hl
func TestOrA_hl(t *testing.T) {
	cpu := NewCPU()
	cpu.LinkMMU(NewMockMMU())
	cpu.R.H = 0x11
	cpu.R.L = 0x1F
	var hlAddr types.Word = types.Word(utils.JoinBytes(cpu.R.H, cpu.R.L))
	cpu.WriteByte(hlAddr, 0x0F)
	cpu.R.A = 0x1A

	cpu.OrA_hl()

	var expected byte = 0x1F
	assert.Equal(t, cpu.R.A, expected)
}

func TestOrA_hlFlagsReset(t *testing.T) {
	cpu := NewCPU()
	cpu.LinkMMU(NewMockMMU())
	cpu.R.H = 0x11
	cpu.R.L = 0x1F
	var hlAddr types.Word = types.Word(utils.JoinBytes(cpu.R.H, cpu.R.L))
	cpu.WriteByte(hlAddr, 0x0F)
	cpu.R.A = 0x1A

	cpu.SetFlag(N)
	cpu.SetFlag(H)
	cpu.SetFlag(C)
	cpu.OrA_hl()

	assert.False(t, cpu.IsFlagSet(N))
	assert.False(t, cpu.IsFlagSet(H))
	assert.False(t, cpu.IsFlagSet(C))
}

func TestOrA_hlZeroFlag(t *testing.T) {
	cpu := NewCPU()
	cpu.LinkMMU(NewMockMMU())
	cpu.R.H = 0x11
	cpu.R.L = 0x1F
	var hlAddr types.Word = types.Word(utils.JoinBytes(cpu.R.H, cpu.R.L))
	cpu.WriteByte(hlAddr, 0x00)
	cpu.R.A = 0x00

	assert.False(t, cpu.IsFlagSet(Z))

	cpu.OrA_hl()

	assert.True(t, cpu.IsFlagSet(Z))
}

//OR A n
func TestOrA_n(t *testing.T) {
	cpu := NewCPU()
	cpu.CurrentInstruction = OneOperandsInstr
	cpu.CurrentInstruction.Operands[0] = 0x0F

	cpu.R.A = 0x1A

	cpu.OrA_n()

	var expected byte = 0x1F
	assert.Equal(t, cpu.R.A, expected)
}

func TestOrA_nFlagsReset(t *testing.T) {
	cpu := NewCPU()
	cpu.CurrentInstruction = OneOperandsInstr
	cpu.CurrentInstruction.Operands[0] = 0x0F

	cpu.R.A = 0x1A

	cpu.SetFlag(N)
	cpu.SetFlag(H)
	cpu.SetFlag(C)
	cpu.OrA_n()

	assert.False(t, cpu.IsFlagSet(N))
	assert.False(t, cpu.IsFlagSet(H))
	assert.False(t, cpu.IsFlagSet(C))
}

func TestOrA_nZeroFlag(t *testing.T) {
	cpu := NewCPU()
	cpu.CurrentInstruction = OneOperandsInstr
	cpu.CurrentInstruction.Operands[0] = 0x00

	cpu.R.A = 0x00

	assert.False(t, cpu.IsFlagSet(Z))

	cpu.OrA_n()

	assert.True(t, cpu.IsFlagSet(Z))
}

func TestDI(t *testing.T) {
	cpu := NewCPU()
	cpu.InterruptsEnabled = true

	cpu.DI()

	var expected bool = false
	assert.Equal(t, cpu.InterruptsEnabled, expected)
}

//SUB A r
func TestSubA_r(t *testing.T) {
	cpu := NewCPU()
	cpu.R.A = 0x0A
	cpu.R.B = 0x02

	cpu.SubA_r(&cpu.R.B)

	var expected byte = 0x08
	assert.Equal(t, cpu.R.A, expected)
	assert.False(t, cpu.IsFlagSet(Z))
}

func TestSubA_rNFlagSet(t *testing.T) {
	cpu := NewCPU()
	cpu.R.A = 0x0A
	cpu.R.B = 0x02

	cpu.SubA_r(&cpu.R.B)

	assert.True(t, cpu.IsFlagSet(N))
}

func TestSubA_rZeroFlagSet(t *testing.T) {
	cpu := NewCPU()
	cpu.R.A = 0x0A

	cpu.SubA_r(&cpu.R.A)

	assert.True(t, cpu.IsFlagSet(Z))
}

func TestSubA_rCarryFlagSetOnNoBorrow(t *testing.T) {
	cpu := NewCPU()
	cpu.R.A = 0xA0
	cpu.R.B = 0x01

	cpu.SubA_r(&cpu.R.B)

	assert.True(t, cpu.IsFlagSet(C))
}

func TestSubA_rCarryFlagResetOnBorrow(t *testing.T) {
	cpu := NewCPU()
	cpu.R.A = 0xA0
	cpu.R.B = 0xFF

	cpu.SubA_r(&cpu.R.B)

	assert.False(t, cpu.IsFlagSet(C))
}

func TestSubA_rHalfCarryFlagSetOnNoBorrow(t *testing.T) {
	cpu := NewCPU()
	cpu.R.A = 0xA0
	cpu.R.B = 0x01

	cpu.SubA_r(&cpu.R.B)

	assert.True(t, cpu.IsFlagSet(H))
}

func TestSubA_rHalfCarryFlagResetOnBorrow(t *testing.T) {
	cpu := NewCPU()
	cpu.R.A = 0xA1
	cpu.R.B = 0x01

	cpu.SubA_r(&cpu.R.B)

	assert.False(t, cpu.IsFlagSet(H))
}

//SUB A (HL)
func TestSubA_hl(t *testing.T) {
	cpu := NewCPU()
	cpu.LinkMMU(NewMockMMU())
	cpu.R.A = 0x0A
	cpu.R.H, cpu.R.L = 0x00, 0x01
	cpu.mmu.WriteByte(types.Word(utils.JoinBytes(cpu.R.H, cpu.R.L)), 0x02)

	cpu.SubA_hl()

	var expected byte = 0x08
	assert.Equal(t, cpu.R.A, expected)
	assert.False(t, cpu.IsFlagSet(Z))
}

func TestSubA_hlNFlagSet(t *testing.T) {
	cpu := NewCPU()
	cpu.LinkMMU(NewMockMMU())
	cpu.R.A = 0x0A
	cpu.R.H, cpu.R.L = 0x00, 0x01
	cpu.mmu.WriteByte(types.Word(utils.JoinBytes(cpu.R.H, cpu.R.L)), 0x02)

	cpu.SubA_hl()

	assert.True(t, cpu.IsFlagSet(N))
}

func TestSubA_hlZeroFlagSet(t *testing.T) {
	cpu := NewCPU()
	cpu.LinkMMU(NewMockMMU())
	cpu.R.A = 0x0A
	cpu.R.H, cpu.R.L = 0x00, 0x01
	cpu.mmu.WriteByte(types.Word(utils.JoinBytes(cpu.R.H, cpu.R.L)), cpu.R.A)

	cpu.SubA_hl()

	assert.True(t, cpu.IsFlagSet(Z))
}

func TestSubA_hlCarryFlagSetOnNoBorrow(t *testing.T) {
	cpu := NewCPU()
	cpu.LinkMMU(NewMockMMU())
	cpu.R.A = 0xA0
	cpu.R.H, cpu.R.L = 0x00, 0x01
	cpu.mmu.WriteByte(types.Word(utils.JoinBytes(cpu.R.H, cpu.R.L)), 0x01)

	cpu.SubA_hl()

	assert.True(t, cpu.IsFlagSet(C))
}

func TestSubA_hlCarryFlagResetOnBorrow(t *testing.T) {
	cpu := NewCPU()
	cpu.LinkMMU(NewMockMMU())
	cpu.R.A = 0xA0
	cpu.R.H, cpu.R.L = 0x00, 0x01
	cpu.mmu.WriteByte(types.Word(utils.JoinBytes(cpu.R.H, cpu.R.L)), 0xFF)

	cpu.SubA_hl()

	assert.False(t, cpu.IsFlagSet(C))
}

func TestSubA_hlHalfCarryFlagSetOnNoBorrow(t *testing.T) {
	cpu := NewCPU()
	cpu.LinkMMU(NewMockMMU())
	cpu.R.A = 0xA0
	cpu.R.H, cpu.R.L = 0x00, 0x01
	cpu.mmu.WriteByte(types.Word(utils.JoinBytes(cpu.R.H, cpu.R.L)), 0x01)

	cpu.SubA_hl()

	assert.True(t, cpu.IsFlagSet(H))
}

func TestSubA_hlHalfCarryFlagResetOnBorrow(t *testing.T) {
	cpu := NewCPU()
	cpu.LinkMMU(NewMockMMU())
	cpu.R.A = 0xA1
	cpu.R.H, cpu.R.L = 0x00, 0x01
	cpu.mmu.WriteByte(types.Word(utils.JoinBytes(cpu.R.H, cpu.R.L)), 0x01)

	cpu.SubA_hl()

	assert.False(t, cpu.IsFlagSet(H))
}

//SUB A, n
func TestSubA_n(t *testing.T) {
	cpu := NewCPU()
	cpu.R.A = 0x0A
	cpu.CurrentInstruction = OneOperandsInstr
	cpu.CurrentInstruction.Operands[0] = 0x02

	cpu.SubA_n()

	var expected byte = 0x08
	assert.Equal(t, cpu.R.A, expected)
	assert.False(t, cpu.IsFlagSet(Z))
}

func TestSubA_nNFlagSet(t *testing.T) {
	cpu := NewCPU()
	cpu.R.A = 0x0A
	cpu.CurrentInstruction = OneOperandsInstr
	cpu.CurrentInstruction.Operands[0] = 0x02

	cpu.SubA_n()

	assert.True(t, cpu.IsFlagSet(N))
}

func TestSubA_nZeroFlagSet(t *testing.T) {
	cpu := NewCPU()
	cpu.R.A = 0x0A
	cpu.CurrentInstruction = OneOperandsInstr
	cpu.CurrentInstruction.Operands[0] = cpu.R.A

	cpu.SubA_n()

	assert.True(t, cpu.IsFlagSet(Z))
}

func TestSubA_nCarryFlagSetOnNoBorrow(t *testing.T) {
	cpu := NewCPU()
	cpu.R.A = 0xA0
	cpu.CurrentInstruction = OneOperandsInstr
	cpu.CurrentInstruction.Operands[0] = 0x01

	cpu.SubA_n()

	assert.True(t, cpu.IsFlagSet(C))
}

func TestSubA_nCarryFlagResetOnBorrow(t *testing.T) {
	cpu := NewCPU()
	cpu.R.A = 0xA0
	cpu.CurrentInstruction = OneOperandsInstr
	cpu.CurrentInstruction.Operands[0] = 0xFF

	cpu.SubA_n()

	assert.False(t, cpu.IsFlagSet(C))
}

func TestSubA_nHalfCarryFlagSetOnNoBorrow(t *testing.T) {
	cpu := NewCPU()
	cpu.R.A = 0xA0
	cpu.CurrentInstruction = OneOperandsInstr
	cpu.CurrentInstruction.Operands[0] = 0x01

	cpu.SubA_n()

	assert.True(t, cpu.IsFlagSet(H))
}

func TestSubA_nHalfCarryFlagResetOnBorrow(t *testing.T) {
	cpu := NewCPU()
	cpu.R.A = 0xA1
	cpu.CurrentInstruction = OneOperandsInstr
	cpu.CurrentInstruction.Operands[0] = 0x01

	cpu.SubA_n()

	assert.False(t, cpu.IsFlagSet(H))
}

//LD A r
func TestLDrr(t *testing.T) {
	cpu := NewCPU()
	cpu.R.A = 0xAA
	cpu.R.D = 0x05

	cpu.LDrr(&cpu.R.D, &cpu.R.A)

	assert.Equal(t, cpu.R.D, byte(0xAA))
}

//INC r
func TestInc_r(t *testing.T) {
	cpu := NewCPU()
	cpu.R.B = 0x01

	cpu.Inc_r(&cpu.R.B)

	var expected byte = 0x02
	assert.Equal(t, cpu.R.B, expected)
}

func TestInc_rNFlagReset(t *testing.T) {
	cpu := NewCPU()
	cpu.R.B = 0x01

	cpu.SetFlag(N)
	cpu.Inc_r(&cpu.R.B)

	assert.False(t, cpu.IsFlagSet(N))
}

func TestInc_rEnsureCarryFlagIsUnaffected(t *testing.T) {
	cpu := NewCPU()
	cpu.R.B = 0x01

	cpu.SetFlag(C)
	cpu.Inc_r(&cpu.R.B)

	assert.True(t, cpu.IsFlagSet(C))

	cpu = NewCPU()
	cpu.R.B = 0x01

	cpu.ResetFlag(C)
	cpu.Inc_r(&cpu.R.B)

	assert.False(t, cpu.IsFlagSet(C))
}

func TestInc_rCheckForHalfCarryOnBorrow(t *testing.T) {
	cpu := NewCPU()
	cpu.R.B = 0x1F

	assert.False(t, cpu.IsFlagSet(H))

	cpu.Inc_r(&cpu.R.B)

	assert.True(t, cpu.IsFlagSet(H))
}

func TestInc_rCheckForZeroFlag(t *testing.T) {
	cpu := NewCPU()
	cpu.R.B = 0xFF

	assert.False(t, cpu.IsFlagSet(Z))

	cpu.Inc_r(&cpu.R.B)

	assert.True(t, cpu.IsFlagSet(Z))
}

//INC hl
func TestInc_hl(t *testing.T) {
	cpu := NewCPU()
	cpu.LinkMMU(NewMockMMU())
	cpu.R.H, cpu.R.L = 0x00, 0x01
	hlAddr := types.Word(utils.JoinBytes(cpu.R.H, cpu.R.L))
	cpu.mmu.WriteByte(hlAddr, 0x01)

	cpu.Inc_hl()

	var expected byte = 0x02
	assert.Equal(t, cpu.mmu.ReadByte(hlAddr), expected)
}

func TestInc_hlNFlagReset(t *testing.T) {
	cpu := NewCPU()
	cpu.LinkMMU(NewMockMMU())
	cpu.R.H, cpu.R.L = 0x00, 0x01
	hlAddr := types.Word(utils.JoinBytes(cpu.R.H, cpu.R.L))
	cpu.mmu.WriteByte(hlAddr, 0x01)

	cpu.SetFlag(N)
	cpu.Inc_hl()

	assert.False(t, cpu.IsFlagSet(N))
}

func TestInc_hlEnsureCarryFlagIsUnaffected(t *testing.T) {
	cpu := NewCPU()
	cpu.LinkMMU(NewMockMMU())
	cpu.R.H, cpu.R.L = 0x00, 0x01
	hlAddr := types.Word(utils.JoinBytes(cpu.R.H, cpu.R.L))
	cpu.mmu.WriteByte(hlAddr, 0x01)

	cpu.SetFlag(C)
	cpu.Inc_hl()

	assert.True(t, cpu.IsFlagSet(C))

	cpu = NewCPU()
	cpu.LinkMMU(NewMockMMU())
	cpu.R.H, cpu.R.L = 0x00, 0x01
	hlAddr = types.Word(utils.JoinBytes(cpu.R.H, cpu.R.L))
	cpu.mmu.WriteByte(hlAddr, 0x01)

	cpu.ResetFlag(C)
	cpu.Inc_hl()

	assert.False(t, cpu.IsFlagSet(C))
}

func TestInc_hlCheckForHalfCarryOnBorrow(t *testing.T) {
	cpu := NewCPU()
	cpu.LinkMMU(NewMockMMU())
	cpu.R.H, cpu.R.L = 0x00, 0x01
	hlAddr := types.Word(utils.JoinBytes(cpu.R.H, cpu.R.L))
	cpu.mmu.WriteByte(hlAddr, 0x1F)

	assert.False(t, cpu.IsFlagSet(H))

	cpu.Inc_hl()

	assert.True(t, cpu.IsFlagSet(H))
}

func TestInc_hlCheckForZeroFlag(t *testing.T) {
	cpu := NewCPU()
	cpu.LinkMMU(NewMockMMU())
	cpu.R.H, cpu.R.L = 0x00, 0x01
	hlAddr := types.Word(utils.JoinBytes(cpu.R.H, cpu.R.L))
	cpu.mmu.WriteByte(hlAddr, 0xFF)

	assert.False(t, cpu.IsFlagSet(Z))

	cpu.Inc_hl()

	assert.True(t, cpu.IsFlagSet(Z))
}

//INC rr
func TestInc_rr(t *testing.T) {
	cpu := NewCPU()
	cpu.R.B = 0x01
	cpu.R.C = 0x02

	cpu.Inc_rr(&cpu.R.B, &cpu.R.C)

	var expected types.Word = 0x0103
	var result types.Word = types.Word(utils.JoinBytes(cpu.R.B, cpu.R.C))
	assert.Equal(t, result, expected)
}

func TestInc_rrNoFlagsAffected(t *testing.T) {
	cpu := NewCPU()
	cpu.R.B = 0x01

	cpu.SetFlag(Z)
	cpu.SetFlag(H)
	cpu.SetFlag(N)
	cpu.SetFlag(C)
	cpu.Inc_rr(&cpu.R.B, &cpu.R.C)

	assert.True(t, cpu.IsFlagSet(Z))
	assert.True(t, cpu.IsFlagSet(H))
	assert.True(t, cpu.IsFlagSet(N))
	assert.True(t, cpu.IsFlagSet(C))
}

//DEC rr
func TestDec_rr(t *testing.T) {
	cpu := NewCPU()
	cpu.R.B = 0x01
	cpu.R.C = 0x02

	cpu.Dec_rr(&cpu.R.B, &cpu.R.C)

	var expected types.Word = 0x0101
	var result types.Word = types.Word(utils.JoinBytes(cpu.R.B, cpu.R.C))
	assert.Equal(t, result, expected)
}

func TestDec_rrNoFlagsAffected(t *testing.T) {
	cpu := NewCPU()
	cpu.R.B = 0x01

	cpu.SetFlag(Z)
	cpu.SetFlag(H)
	cpu.SetFlag(N)
	cpu.SetFlag(C)
	cpu.Dec_rr(&cpu.R.B, &cpu.R.C)

	assert.True(t, cpu.IsFlagSet(Z))
	assert.True(t, cpu.IsFlagSet(H))
	assert.True(t, cpu.IsFlagSet(N))
	assert.True(t, cpu.IsFlagSet(C))
}

//NOP
func TestNOP(t *testing.T) {
	cpu := NewCPU()
	cpu.NOP()

	assert.True(t, cpu.PC == types.Word(0x0000))
	assert.True(t, cpu.SP == types.Word(0x0000))
	assert.True(t, cpu.R.A == 0x00)
	assert.True(t, cpu.R.B == 0x00)
	assert.True(t, cpu.R.C == 0x00)
	assert.True(t, cpu.R.D == 0x00)
	assert.True(t, cpu.R.E == 0x00)
	assert.True(t, cpu.R.H == 0x00)
	assert.True(t, cpu.R.L == 0x00)
}

//LD A, (C)
func TestLDr_ffplusc(t *testing.T) {
	cpu := NewCPU()
	cpu.LinkMMU(NewMockMMU())
	cpu.R.A = 0x01
	cpu.R.C = 0x0A
	var expected byte = 0xAA
	cpu.WriteByte(0xFF00+types.Word(cpu.R.C), expected)

	cpu.LDr_ffplusc(&cpu.R.A)

	assert.Equal(t, cpu.R.A, expected)
}

//RET CC
func TestRetCC_NZ(t *testing.T) {
	cpu := NewCPU()
	cpu.LinkMMU(NewMockMMU())
	var expectedPC types.Word = 0x0102
	cpu.pushWordToStack(expectedPC)
	cpu.PC = 0x0FFF

	cpu.Retcc(Z, false)

	assert.Equal(t, cpu.PC, expectedPC)

	//and the opposite
	cpu = NewCPU()
	cpu.LinkMMU(NewMockMMU())
	cpu.PC = 0x0FFF
	expectedPC = cpu.PC
	cpu.pushWordToStack(0xFFFF)

	cpu.SetFlag(Z)
	cpu.Retcc(Z, false)

	assert.Equal(t, cpu.PC, expectedPC)
}

func TestRetCC_Z(t *testing.T) {
	cpu := NewCPU()
	cpu.LinkMMU(NewMockMMU())
	var expectedPC types.Word = 0x0102
	cpu.pushWordToStack(expectedPC)
	cpu.PC = 0x0FFF

	cpu.SetFlag(Z)
	cpu.Retcc(Z, true)

	assert.Equal(t, cpu.PC, expectedPC)

	//and the opposite
	cpu = NewCPU()
	cpu.LinkMMU(NewMockMMU())
	cpu.PC = 0x0FFF
	expectedPC = cpu.PC
	cpu.pushWordToStack(0xFFFF)

	cpu.Retcc(Z, true)

	assert.Equal(t, cpu.PC, expectedPC)
}

func TestRetCC_NC(t *testing.T) {
	cpu := NewCPU()
	cpu.LinkMMU(NewMockMMU())
	var expectedPC types.Word = 0x0102
	cpu.pushWordToStack(expectedPC)
	cpu.PC = 0x0FFF

	cpu.Retcc(C, false)

	assert.Equal(t, cpu.PC, expectedPC)

	//and the opposite
	cpu = NewCPU()
	cpu.LinkMMU(NewMockMMU())
	cpu.PC = 0x0FFF
	expectedPC = cpu.PC
	cpu.pushWordToStack(0xFFFF)

	cpu.SetFlag(C)
	cpu.Retcc(C, false)

	assert.Equal(t, cpu.PC, expectedPC)
}

func TestRetCC_C(t *testing.T) {
	cpu := NewCPU()
	cpu.LinkMMU(NewMockMMU())
	var expectedPC types.Word = 0x0102
	cpu.pushWordToStack(expectedPC)
	cpu.PC = 0x0FFF

	cpu.SetFlag(C)
	cpu.Retcc(C, true)

	assert.Equal(t, cpu.PC, expectedPC)

	//and the opposite
	cpu = NewCPU()
	cpu.LinkMMU(NewMockMMU())
	cpu.PC = 0x0FFF
	expectedPC = cpu.PC
	cpu.pushWordToStack(0xFFFF)

	cpu.Retcc(C, true)

	assert.Equal(t, cpu.PC, expectedPC)
}

//ADD A r
func TestAddA_r(t *testing.T) {
	cpu := NewCPU()
	cpu.R.A = 0x05
	cpu.R.B = 0x04

	cpu.AddA_r(&cpu.R.B)

	var expected byte = 0x09
	assert.Equal(t, cpu.R.A, expected)
}

func TestAddA_rNFlagReset(t *testing.T) {
	cpu := NewCPU()
	cpu.R.A = 0x05
	cpu.R.B = 0x04

	cpu.SetFlag(N)
	cpu.AddA_r(&cpu.R.B)

	assert.False(t, cpu.IsFlagSet(N))
}

func TestAddA_rZeroFlag(t *testing.T) {
	cpu := NewCPU()
	cpu.R.A = 0xFF
	cpu.R.B = 0x01

	assert.False(t, cpu.IsFlagSet(Z))

	cpu.AddA_r(&cpu.R.B)

	assert.True(t, cpu.IsFlagSet(Z))
}

func TestAddA_rCarryFlag(t *testing.T) {
	cpu := NewCPU()
	cpu.R.A = 0xFF
	cpu.R.B = 0x05

	assert.False(t, cpu.IsFlagSet(C))

	cpu.AddA_r(&cpu.R.B)

	assert.True(t, cpu.IsFlagSet(C))
}

func TestAddA_rHalfCarryFlag(t *testing.T) {
	cpu := NewCPU()
	cpu.R.A = 0x1F
	cpu.R.B = 0x01

	assert.False(t, cpu.IsFlagSet(H))

	cpu.AddA_r(&cpu.R.B)

	assert.True(t, cpu.IsFlagSet(H))
}

//ADD A (HL)
func TestAddA_hl(t *testing.T) {
	cpu := NewCPU()
	cpu.LinkMMU(NewMockMMU())
	cpu.R.A = 0x05
	cpu.R.H, cpu.R.L = 0x00, 0x01
	hlAddr := types.Word(utils.JoinBytes(cpu.R.H, cpu.R.L))
	cpu.mmu.WriteByte(hlAddr, 0x04)

	cpu.AddA_hl()

	var expected byte = 0x09
	assert.Equal(t, cpu.R.A, expected)
}

func TestAddA_hlNFlagReset(t *testing.T) {
	cpu := NewCPU()
	cpu.LinkMMU(NewMockMMU())
	cpu.R.A = 0x05
	cpu.R.H, cpu.R.L = 0x00, 0x01
	hlAddr := types.Word(utils.JoinBytes(cpu.R.H, cpu.R.L))
	cpu.mmu.WriteByte(hlAddr, 0x04)

	cpu.SetFlag(N)
	cpu.AddA_hl()

	assert.False(t, cpu.IsFlagSet(N))
}

func TestAddA_hlZeroFlag(t *testing.T) {
	cpu := NewCPU()
	cpu.LinkMMU(NewMockMMU())
	cpu.R.A = 0xFF
	cpu.R.H, cpu.R.L = 0x00, 0x01
	hlAddr := types.Word(utils.JoinBytes(cpu.R.H, cpu.R.L))
	cpu.mmu.WriteByte(hlAddr, 0x01)

	assert.False(t, cpu.IsFlagSet(Z))

	cpu.AddA_hl()

	assert.True(t, cpu.IsFlagSet(Z))
}

func TestAddA_hlCarryFlag(t *testing.T) {
	cpu := NewCPU()
	cpu.LinkMMU(NewMockMMU())
	cpu.R.A = 0xFF
	cpu.R.H, cpu.R.L = 0x00, 0x01
	hlAddr := types.Word(utils.JoinBytes(cpu.R.H, cpu.R.L))
	cpu.mmu.WriteByte(hlAddr, 0x05)

	assert.False(t, cpu.IsFlagSet(C))

	cpu.AddA_hl()

	assert.True(t, cpu.IsFlagSet(C))
}

func TestAddA_hlHalfCarryFlag(t *testing.T) {
	cpu := NewCPU()
	cpu.LinkMMU(NewMockMMU())
	cpu.R.A = 0x1F
	cpu.R.H, cpu.R.L = 0x00, 0x01
	hlAddr := types.Word(utils.JoinBytes(cpu.R.H, cpu.R.L))
	cpu.mmu.WriteByte(hlAddr, 0x01)

	assert.False(t, cpu.IsFlagSet(H))

	cpu.AddA_hl()

	assert.True(t, cpu.IsFlagSet(H))
}

//ADD A n
func TestAddA_n(t *testing.T) {
	cpu := NewCPU()
	cpu.CurrentInstruction = OneOperandsInstr
	cpu.CurrentInstruction.Operands[0] = 0x04
	cpu.R.A = 0x05

	cpu.AddA_n()

	var expected byte = 0x09
	assert.Equal(t, cpu.R.A, expected)
}

func TestAddA_nNFlagReset(t *testing.T) {
	cpu := NewCPU()
	cpu.CurrentInstruction = OneOperandsInstr
	cpu.CurrentInstruction.Operands[0] = 0x04
	cpu.R.A = 0x05

	cpu.SetFlag(N)
	cpu.AddA_n()

	assert.False(t, cpu.IsFlagSet(N))
}

func TestAddA_nZeroFlag(t *testing.T) {
	cpu := NewCPU()
	cpu.CurrentInstruction = OneOperandsInstr
	cpu.CurrentInstruction.Operands[0] = 0x01
	cpu.R.A = 0xFF

	assert.False(t, cpu.IsFlagSet(Z))

	cpu.AddA_n()

	assert.True(t, cpu.IsFlagSet(Z))
}

func TestAddA_nCarryFlag(t *testing.T) {
	cpu := NewCPU()
	cpu.R.A = 0xFF
	cpu.CurrentInstruction = OneOperandsInstr
	cpu.CurrentInstruction.Operands[0] = 0x05

	assert.False(t, cpu.IsFlagSet(C))

	cpu.AddA_n()

	assert.True(t, cpu.IsFlagSet(C))
}

func TestAddA_nHalfCarryFlag(t *testing.T) {
	cpu := NewCPU()
	cpu.R.A = 0x1F

	cpu.CurrentInstruction = OneOperandsInstr
	cpu.CurrentInstruction.Operands[0] = 0x01

	assert.False(t, cpu.IsFlagSet(H))

	cpu.AddA_n()

	assert.True(t, cpu.IsFlagSet(H))
}

//XOR A r
func TestXorA_r(t *testing.T) {
	cpu := NewCPU()
	cpu.R.A = 0x0A
	cpu.R.B = 0x0F

	cpu.XorA_r(&cpu.R.B)

	var expected byte = 0x05
	assert.Equal(t, cpu.R.A, expected)
}

func TestXorA_rFlagsReset(t *testing.T) {
	cpu := NewCPU()
	cpu.R.A = 0x0A
	cpu.R.B = 0x0F

	cpu.SetFlag(N)
	cpu.SetFlag(H)
	cpu.SetFlag(C)
	cpu.XorA_r(&cpu.R.B)

	assert.False(t, cpu.IsFlagSet(N))
	assert.False(t, cpu.IsFlagSet(H))
	assert.False(t, cpu.IsFlagSet(C))
}

func TestXorA_rZeroFlag(t *testing.T) {
	cpu := NewCPU()
	cpu.R.A = 0xFF

	assert.False(t, cpu.IsFlagSet(Z))
	cpu.XorA_r(&cpu.R.A)

	assert.True(t, cpu.IsFlagSet(Z))
}

//XOR A hl
func TestXorA_hl(t *testing.T) {
	cpu := NewCPU()
	cpu.LinkMMU(NewMockMMU())
	cpu.R.H = 0x11
	cpu.R.L = 0x1F
	var hlAddr types.Word = types.Word(utils.JoinBytes(cpu.R.H, cpu.R.L))
	cpu.WriteByte(hlAddr, 0x0F)
	cpu.R.A = 0x0A

	cpu.XorA_hl()

	var expected byte = 0x05
	assert.Equal(t, cpu.R.A, expected)
}

func TestXorA_hlFlagsReset(t *testing.T) {
	cpu := NewCPU()
	cpu.LinkMMU(NewMockMMU())
	cpu.R.H = 0x11
	cpu.R.L = 0x1F
	var hlAddr types.Word = types.Word(utils.JoinBytes(cpu.R.H, cpu.R.L))
	cpu.WriteByte(hlAddr, 0x0F)
	cpu.R.A = 0x0A

	cpu.SetFlag(N)
	cpu.SetFlag(H)
	cpu.SetFlag(C)
	cpu.XorA_hl()

	assert.False(t, cpu.IsFlagSet(N))
	assert.False(t, cpu.IsFlagSet(H))
	assert.False(t, cpu.IsFlagSet(C))
}

func TestXorA_hlZeroFlag(t *testing.T) {
	cpu := NewCPU()
	cpu.LinkMMU(NewMockMMU())
	cpu.R.H = 0x11
	cpu.R.L = 0x1F
	var hlAddr types.Word = types.Word(utils.JoinBytes(cpu.R.H, cpu.R.L))
	cpu.WriteByte(hlAddr, 0x0A)
	cpu.R.A = 0x0A

	assert.False(t, cpu.IsFlagSet(Z))

	cpu.XorA_hl()

	assert.True(t, cpu.IsFlagSet(Z))
}

//XOR A n
func TestXorA_n(t *testing.T) {
	cpu := NewCPU()
	cpu.CurrentInstruction = OneOperandsInstr
	cpu.CurrentInstruction.Operands[0] = 0x0F

	cpu.R.A = 0x0A

	cpu.XorA_n()

	var expected byte = 0x05
	assert.Equal(t, cpu.R.A, expected)
}

func TestXorA_nFlagsReset(t *testing.T) {
	cpu := NewCPU()
	cpu.CurrentInstruction = OneOperandsInstr
	cpu.CurrentInstruction.Operands[0] = 0x0F

	cpu.R.A = 0x0A

	cpu.SetFlag(N)
	cpu.SetFlag(H)
	cpu.SetFlag(C)
	cpu.XorA_n()

	assert.False(t, cpu.IsFlagSet(N))
	assert.False(t, cpu.IsFlagSet(H))
	assert.False(t, cpu.IsFlagSet(C))
}

func TestXorA_nZeroFlag(t *testing.T) {
	cpu := NewCPU()
	cpu.CurrentInstruction = OneOperandsInstr
	cpu.CurrentInstruction.Operands[0] = 0x0A

	cpu.R.A = 0x0A

	assert.False(t, cpu.IsFlagSet(Z))

	cpu.XorA_n()

	assert.True(t, cpu.IsFlagSet(Z))
}

//ADC A r
func TestAddCA_r(t *testing.T) {
	cpu := NewCPU()
	cpu.R.A = 0x05
	cpu.R.B = 0x04

	cpu.AddCA_r(&cpu.R.B)

	var expected byte = 0x09
	assert.Equal(t, cpu.R.A, expected)
}

func TestAddCA_rWithCarry(t *testing.T) {
	cpu := NewCPU()
	cpu.R.A = 0x05
	cpu.R.B = 0x04
	cpu.SetFlag(C)

	cpu.AddCA_r(&cpu.R.B)

	var expected byte = 0x0A
	assert.Equal(t, cpu.R.A, expected)
}

func TestAddCA_rNFlagReset(t *testing.T) {
	cpu := NewCPU()
	cpu.R.A = 0x05
	cpu.R.B = 0x04

	cpu.SetFlag(N)
	cpu.AddCA_r(&cpu.R.B)

	assert.False(t, cpu.IsFlagSet(N))
}

func TestAddCA_rZeroFlag(t *testing.T) {
	cpu := NewCPU()
	cpu.R.A = 0xFF
	cpu.R.B = 0x01

	assert.False(t, cpu.IsFlagSet(Z))

	cpu.AddCA_r(&cpu.R.B)

	assert.True(t, cpu.IsFlagSet(Z))
}

func TestAddCA_rCarryFlag(t *testing.T) {
	cpu := NewCPU()
	cpu.R.A = 0xFF
	cpu.R.B = 0x05

	assert.False(t, cpu.IsFlagSet(C))

	cpu.AddCA_r(&cpu.R.B)

	assert.True(t, cpu.IsFlagSet(C))
}

func TestAddCA_rHalfCarryFlag(t *testing.T) {
	cpu := NewCPU()
	cpu.R.A = 0x1F
	cpu.R.B = 0x01

	assert.False(t, cpu.IsFlagSet(H))

	cpu.AddCA_r(&cpu.R.B)

	assert.True(t, cpu.IsFlagSet(H))
}

//ADC A (HL)
func TestAddCA_hl(t *testing.T) {
	cpu := NewCPU()
	cpu.LinkMMU(NewMockMMU())
	cpu.R.A = 0x05
	cpu.R.H, cpu.R.L = 0x00, 0x01
	hlAddCr := types.Word(utils.JoinBytes(cpu.R.H, cpu.R.L))
	cpu.mmu.WriteByte(hlAddCr, 0x04)

	cpu.AddCA_hl()

	var expected byte = 0x09
	assert.Equal(t, cpu.R.A, expected)
}

func TestAddCA_hlWithCarry(t *testing.T) {
	cpu := NewCPU()
	cpu.LinkMMU(NewMockMMU())
	cpu.R.A = 0x05
	cpu.R.H, cpu.R.L = 0x00, 0x01
	hlAddCr := types.Word(utils.JoinBytes(cpu.R.H, cpu.R.L))
	cpu.mmu.WriteByte(hlAddCr, 0x04)
	cpu.SetFlag(C)

	cpu.AddCA_hl()

	var expected byte = 0x0A
	assert.Equal(t, cpu.R.A, expected)
}

func TestAddCA_hlNFlagReset(t *testing.T) {
	cpu := NewCPU()
	cpu.LinkMMU(NewMockMMU())
	cpu.R.A = 0x05
	cpu.R.H, cpu.R.L = 0x00, 0x01
	hlAddCr := types.Word(utils.JoinBytes(cpu.R.H, cpu.R.L))
	cpu.mmu.WriteByte(hlAddCr, 0x04)

	cpu.SetFlag(N)
	cpu.AddCA_hl()

	assert.False(t, cpu.IsFlagSet(N))
}

func TestAddCA_hlZeroFlag(t *testing.T) {
	cpu := NewCPU()
	cpu.LinkMMU(NewMockMMU())
	cpu.R.A = 0xFF
	cpu.R.H, cpu.R.L = 0x00, 0x01
	hlAddCr := types.Word(utils.JoinBytes(cpu.R.H, cpu.R.L))
	cpu.mmu.WriteByte(hlAddCr, 0x01)

	assert.False(t, cpu.IsFlagSet(Z))

	cpu.AddCA_hl()

	assert.True(t, cpu.IsFlagSet(Z))
}

func TestAddCA_hlCarryFlag(t *testing.T) {
	cpu := NewCPU()
	cpu.LinkMMU(NewMockMMU())
	cpu.R.A = 0xFF
	cpu.R.H, cpu.R.L = 0x00, 0x01
	hlAddCr := types.Word(utils.JoinBytes(cpu.R.H, cpu.R.L))
	cpu.mmu.WriteByte(hlAddCr, 0x05)

	assert.False(t, cpu.IsFlagSet(C))

	cpu.AddCA_hl()

	assert.True(t, cpu.IsFlagSet(C))
}

func TestAddCA_hlHalfCarryFlag(t *testing.T) {
	cpu := NewCPU()
	cpu.LinkMMU(NewMockMMU())
	cpu.R.A = 0x1F
	cpu.R.H, cpu.R.L = 0x00, 0x01
	hlAddCr := types.Word(utils.JoinBytes(cpu.R.H, cpu.R.L))
	cpu.mmu.WriteByte(hlAddCr, 0x01)

	assert.False(t, cpu.IsFlagSet(H))

	cpu.AddCA_hl()

	assert.True(t, cpu.IsFlagSet(H))
}

//ADC A n
func TestAddCA_n(t *testing.T) {
	cpu := NewCPU()
	cpu.CurrentInstruction = OneOperandsInstr
	cpu.CurrentInstruction.Operands[0] = 0x04
	cpu.R.A = 0x05

	cpu.AddCA_n()

	var expected byte = 0x09
	assert.Equal(t, cpu.R.A, expected)
}

func TestAddCA_nWithCarry(t *testing.T) {
	cpu := NewCPU()
	cpu.CurrentInstruction = OneOperandsInstr
	cpu.CurrentInstruction.Operands[0] = 0x04
	cpu.R.A = 0x05
	cpu.SetFlag(C)

	cpu.AddCA_n()

	var expected byte = 0x0A
	assert.Equal(t, cpu.R.A, expected)
}

func TestAddCA_nNFlagReset(t *testing.T) {
	cpu := NewCPU()
	cpu.CurrentInstruction = OneOperandsInstr
	cpu.CurrentInstruction.Operands[0] = 0x04
	cpu.R.A = 0x05

	cpu.SetFlag(N)
	cpu.AddCA_n()

	assert.False(t, cpu.IsFlagSet(N))
}

func TestAddCA_nZeroFlag(t *testing.T) {
	cpu := NewCPU()
	cpu.CurrentInstruction = OneOperandsInstr
	cpu.CurrentInstruction.Operands[0] = 0x01
	cpu.R.A = 0xFF

	assert.False(t, cpu.IsFlagSet(Z))

	cpu.AddCA_n()

	assert.True(t, cpu.IsFlagSet(Z))
}

func TestAddCA_nCarryFlag(t *testing.T) {
	cpu := NewCPU()
	cpu.R.A = 0xFF
	cpu.CurrentInstruction = OneOperandsInstr
	cpu.CurrentInstruction.Operands[0] = 0x05

	assert.False(t, cpu.IsFlagSet(C))

	cpu.AddCA_n()

	assert.True(t, cpu.IsFlagSet(C))
}

func TestAddCA_nHalfCarryFlag(t *testing.T) {
	cpu := NewCPU()
	cpu.R.A = 0x1F

	cpu.CurrentInstruction = OneOperandsInstr
	cpu.CurrentInstruction.Operands[0] = 0x01

	assert.False(t, cpu.IsFlagSet(H))

	cpu.AddCA_n()

	assert.True(t, cpu.IsFlagSet(H))
}

//SBC A r
func TestSubAC_r(t *testing.T) {
	cpu := NewCPU()
	cpu.R.A = 0x0A
	cpu.R.B = 0x02

	cpu.SubAC_r(&cpu.R.B)

	var expected byte = 0x08
	assert.Equal(t, cpu.R.A, expected)
	assert.False(t, cpu.IsFlagSet(Z))
}

func TestSubAC_rNFlagSet(t *testing.T) {
	cpu := NewCPU()
	cpu.R.A = 0x0A
	cpu.R.B = 0x02

	cpu.SubAC_r(&cpu.R.B)

	assert.True(t, cpu.IsFlagSet(N))
}

func TestSubAC_rZeroFlagSet(t *testing.T) {
	cpu := NewCPU()
	cpu.R.A = 0x0A

	cpu.SubAC_r(&cpu.R.A)

	assert.True(t, cpu.IsFlagSet(Z))
}

func TestSubAC_rCarryFlagSetOnNoBorrow(t *testing.T) {
	cpu := NewCPU()
	cpu.R.A = 0xA0
	cpu.R.B = 0x01

	cpu.SubAC_r(&cpu.R.B)

	assert.True(t, cpu.IsFlagSet(C))
}

func TestSubAC_rCarryFlagResetOnBorrow(t *testing.T) {
	cpu := NewCPU()
	cpu.R.A = 0xA0
	cpu.R.B = 0xFF

	cpu.SubAC_r(&cpu.R.B)

	assert.False(t, cpu.IsFlagSet(C))
}

func TestSubAC_rHalfCarryFlagSetOnNoBorrow(t *testing.T) {
	cpu := NewCPU()
	cpu.R.A = 0xA0
	cpu.R.B = 0x01

	cpu.SubAC_r(&cpu.R.B)

	assert.True(t, cpu.IsFlagSet(H))
}

func TestSubAC_rHalfCarryFlagResetOnBorrow(t *testing.T) {
	cpu := NewCPU()
	cpu.R.A = 0xA1
	cpu.R.B = 0x01

	cpu.SubAC_r(&cpu.R.B)

	assert.False(t, cpu.IsFlagSet(H))
}

//SBC A (HL)
func TestSubAC_hl(t *testing.T) {
	cpu := NewCPU()
	cpu.LinkMMU(NewMockMMU())
	cpu.R.A = 0x0A
	cpu.R.H, cpu.R.L = 0x00, 0x01
	cpu.mmu.WriteByte(types.Word(utils.JoinBytes(cpu.R.H, cpu.R.L)), 0x02)

	cpu.SubAC_hl()

	var expected byte = 0x08
	assert.Equal(t, cpu.R.A, expected)
	assert.False(t, cpu.IsFlagSet(Z))
}

func TestSubAC_hlNFlagSet(t *testing.T) {
	cpu := NewCPU()
	cpu.LinkMMU(NewMockMMU())
	cpu.R.A = 0x0A
	cpu.R.H, cpu.R.L = 0x00, 0x01
	cpu.mmu.WriteByte(types.Word(utils.JoinBytes(cpu.R.H, cpu.R.L)), 0x02)

	cpu.SubAC_hl()

	assert.True(t, cpu.IsFlagSet(N))
}

func TestSubAC_hlZeroFlagSet(t *testing.T) {
	cpu := NewCPU()
	cpu.LinkMMU(NewMockMMU())
	cpu.R.A = 0x0A
	cpu.R.H, cpu.R.L = 0x00, 0x01
	cpu.mmu.WriteByte(types.Word(utils.JoinBytes(cpu.R.H, cpu.R.L)), cpu.R.A)

	cpu.SubAC_hl()

	assert.True(t, cpu.IsFlagSet(Z))
}

func TestSubAC_hlCarryFlagSetOnNoBorrow(t *testing.T) {
	cpu := NewCPU()
	cpu.LinkMMU(NewMockMMU())
	cpu.R.A = 0xA0
	cpu.R.H, cpu.R.L = 0x00, 0x01
	cpu.mmu.WriteByte(types.Word(utils.JoinBytes(cpu.R.H, cpu.R.L)), 0x01)

	cpu.SubAC_hl()

	assert.True(t, cpu.IsFlagSet(C))
}

func TestSubAC_hlCarryFlagResetOnBorrow(t *testing.T) {
	cpu := NewCPU()
	cpu.LinkMMU(NewMockMMU())
	cpu.R.A = 0xA0
	cpu.R.H, cpu.R.L = 0x00, 0x01
	cpu.mmu.WriteByte(types.Word(utils.JoinBytes(cpu.R.H, cpu.R.L)), 0xFF)

	cpu.SubAC_hl()

	assert.False(t, cpu.IsFlagSet(C))
}

func TestSubAC_hlHalfCarryFlagSetOnNoBorrow(t *testing.T) {
	cpu := NewCPU()
	cpu.LinkMMU(NewMockMMU())
	cpu.R.A = 0xA0
	cpu.R.H, cpu.R.L = 0x00, 0x01
	cpu.mmu.WriteByte(types.Word(utils.JoinBytes(cpu.R.H, cpu.R.L)), 0x01)

	cpu.SubAC_hl()

	assert.True(t, cpu.IsFlagSet(H))
}

func TestSubAC_hlHalfCarryFlagResetOnBorrow(t *testing.T) {
	cpu := NewCPU()
	cpu.LinkMMU(NewMockMMU())
	cpu.R.A = 0xA1
	cpu.R.H, cpu.R.L = 0x00, 0x01
	cpu.mmu.WriteByte(types.Word(utils.JoinBytes(cpu.R.H, cpu.R.L)), 0x01)

	cpu.SubAC_hl()

	assert.False(t, cpu.IsFlagSet(H))
}

//SBC A, n
func TestSubAC_n(t *testing.T) {
	cpu := NewCPU()
	cpu.R.A = 0x0A
	cpu.CurrentInstruction = OneOperandsInstr
	cpu.CurrentInstruction.Operands[0] = 0x02

	cpu.SubAC_n()

	var expected byte = 0x08
	assert.Equal(t, cpu.R.A, expected)
	assert.False(t, cpu.IsFlagSet(Z))
}

func TestSubAC_nNFlagSet(t *testing.T) {
	cpu := NewCPU()
	cpu.R.A = 0x0A
	cpu.CurrentInstruction = OneOperandsInstr
	cpu.CurrentInstruction.Operands[0] = 0x02

	cpu.SubAC_n()

	assert.True(t, cpu.IsFlagSet(N))
}

func TestSubAC_nZeroFlagSet(t *testing.T) {
	cpu := NewCPU()
	cpu.R.A = 0x0A
	cpu.CurrentInstruction = OneOperandsInstr
	cpu.CurrentInstruction.Operands[0] = cpu.R.A

	cpu.SubAC_n()

	assert.True(t, cpu.IsFlagSet(Z))
}

func TestSubAC_nCarryFlagSetOnNoBorrow(t *testing.T) {
	cpu := NewCPU()
	cpu.R.A = 0xA0
	cpu.CurrentInstruction = OneOperandsInstr
	cpu.CurrentInstruction.Operands[0] = 0x01

	cpu.SubAC_n()

	assert.True(t, cpu.IsFlagSet(C))
}

func TestSubAC_nCarryFlagResetOnBorrow(t *testing.T) {
	cpu := NewCPU()
	cpu.R.A = 0xA0
	cpu.CurrentInstruction = OneOperandsInstr
	cpu.CurrentInstruction.Operands[0] = 0xFF

	cpu.SubAC_n()

	assert.False(t, cpu.IsFlagSet(C))
}

func TestSubAC_nHalfCarryFlagSetOnNoBorrow(t *testing.T) {
	cpu := NewCPU()
	cpu.R.A = 0xA0
	cpu.CurrentInstruction = OneOperandsInstr
	cpu.CurrentInstruction.Operands[0] = 0x01

	cpu.SubAC_n()

	assert.True(t, cpu.IsFlagSet(H))
}

func TestSubAC_nHalfCarryFlagResetOnBorrow(t *testing.T) {
	cpu := NewCPU()
	cpu.R.A = 0xA1
	cpu.CurrentInstruction = OneOperandsInstr
	cpu.CurrentInstruction.Operands[0] = 0x01

	cpu.SubAC_n()

	assert.False(t, cpu.IsFlagSet(H))
}

//AND r
func TestAndA_r(t *testing.T) {
	cpu := NewCPU()
	cpu.R.A = 0xDE
	cpu.R.B = 0x8A

	cpu.AndA_r(&cpu.R.B)

	var expected byte = 0x8A

	assert.Equal(t, cpu.R.A, expected)
}

func TestAndA_rFlagsReset(t *testing.T) {
	cpu := NewCPU()
	cpu.R.A = 0xDE
	cpu.R.B = 0x8A

	cpu.SetFlag(N)
	cpu.SetFlag(C)

	cpu.AndA_r(&cpu.R.B)

	assert.False(t, cpu.IsFlagSet(N))
	assert.False(t, cpu.IsFlagSet(C))
}

func TestAndA_rHFlagSet(t *testing.T) {
	cpu := NewCPU()
	cpu.R.A = 0xDE
	cpu.R.B = 0x8A

	assert.False(t, cpu.IsFlagSet(H))

	cpu.AndA_r(&cpu.R.B)

	assert.True(t, cpu.IsFlagSet(H))
}

func TestAndA_rZeroFlagSet(t *testing.T) {
	cpu := NewCPU()
	cpu.R.A = 0xDE
	cpu.R.B = 0x00

	assert.False(t, cpu.IsFlagSet(Z))

	cpu.AndA_r(&cpu.R.B)

	assert.True(t, cpu.IsFlagSet(Z))
}

//AND hl
func TestAndA_hl(t *testing.T) {
	cpu := NewCPU()
	cpu.LinkMMU(NewMockMMU())
	cpu.R.H = 0x11
	cpu.R.L = 0x1F
	var hlAddr types.Word = types.Word(utils.JoinBytes(cpu.R.H, cpu.R.L))
	cpu.WriteByte(hlAddr, 0x8A)
	cpu.R.A = 0xDE

	cpu.AndA_hl()

	var expected byte = 0x8A

	assert.Equal(t, cpu.R.A, expected)
}

func TestAndA_hlFlagsReset(t *testing.T) {
	cpu := NewCPU()
	cpu.LinkMMU(NewMockMMU())
	cpu.R.H = 0x11
	cpu.R.L = 0x1F
	var hlAddr types.Word = types.Word(utils.JoinBytes(cpu.R.H, cpu.R.L))
	cpu.WriteByte(hlAddr, 0x8A)
	cpu.R.A = 0xDE

	cpu.SetFlag(N)
	cpu.SetFlag(C)

	cpu.AndA_hl()

	assert.False(t, cpu.IsFlagSet(N))
	assert.False(t, cpu.IsFlagSet(C))
}

func TestAndA_hlHFlagSet(t *testing.T) {
	cpu := NewCPU()
	cpu.LinkMMU(NewMockMMU())
	cpu.R.H = 0x11
	cpu.R.L = 0x1F
	var hlAddr types.Word = types.Word(utils.JoinBytes(cpu.R.H, cpu.R.L))
	cpu.WriteByte(hlAddr, 0x8A)
	cpu.R.A = 0xDE

	assert.False(t, cpu.IsFlagSet(H))

	cpu.AndA_hl()

	assert.True(t, cpu.IsFlagSet(H))
}

func TestAndA_hlZeroFlagSet(t *testing.T) {
	cpu := NewCPU()
	cpu.LinkMMU(NewMockMMU())
	cpu.R.H = 0x11
	cpu.R.L = 0x1F
	var hlAddr types.Word = types.Word(utils.JoinBytes(cpu.R.H, cpu.R.L))
	cpu.WriteByte(hlAddr, 0x00)
	cpu.R.A = 0xDE

	assert.False(t, cpu.IsFlagSet(Z))

	cpu.AndA_hl()

	assert.True(t, cpu.IsFlagSet(Z))
}

//AND n
func TestAndA_n(t *testing.T) {
	cpu := NewCPU()
	cpu.CurrentInstruction = OneOperandsInstr
	cpu.CurrentInstruction.Operands[0] = 0x8A
	cpu.R.A = 0xDE

	cpu.AndA_n()

	var expected byte = 0x8A

	assert.Equal(t, cpu.R.A, expected)
}

func TestAndA_nFlagsReset(t *testing.T) {
	cpu := NewCPU()
	cpu.CurrentInstruction = OneOperandsInstr
	cpu.CurrentInstruction.Operands[0] = 0x8A
	cpu.R.A = 0xDE

	cpu.SetFlag(N)
	cpu.SetFlag(C)

	cpu.AndA_n()

	assert.False(t, cpu.IsFlagSet(N))
	assert.False(t, cpu.IsFlagSet(C))
}

func TestAndA_nHFlagSet(t *testing.T) {
	cpu := NewCPU()
	cpu.CurrentInstruction = OneOperandsInstr
	cpu.CurrentInstruction.Operands[0] = 0x8A
	cpu.R.A = 0xDE

	assert.False(t, cpu.IsFlagSet(H))

	cpu.AndA_n()

	assert.True(t, cpu.IsFlagSet(H))
}

func TestAndA_nZeroFlagSet(t *testing.T) {
	cpu := NewCPU()
	cpu.CurrentInstruction = OneOperandsInstr
	cpu.CurrentInstruction.Operands[0] = 0x00
	cpu.R.A = 0xDE

	assert.False(t, cpu.IsFlagSet(Z))

	cpu.AndA_n()

	assert.True(t, cpu.IsFlagSet(Z))
}

//DEC r
func TestDec_r(t *testing.T) {
	cpu := NewCPU()
	cpu.R.C = 0x11

	cpu.Dec_r(&cpu.R.C)

	var expected byte = 0x10
	assert.Equal(t, cpu.R.C, expected)
}

func TestDec_rNFlagSet(t *testing.T) {
	cpu := NewCPU()
	cpu.R.C = 0x11

	assert.False(t, cpu.IsFlagSet(N))

	cpu.Dec_r(&cpu.R.C)

	assert.True(t, cpu.IsFlagSet(N))
}

func TestDec_rZeroFlagSet(t *testing.T) {
	cpu := NewCPU()
	cpu.R.C = 0x01

	assert.False(t, cpu.IsFlagSet(Z))

	cpu.Dec_r(&cpu.R.C)

	assert.True(t, cpu.IsFlagSet(Z))
}

func TestDec_rHFlagSetOnBorrow(t *testing.T) {
	cpu := NewCPU()
	cpu.R.C = 0xF0

	assert.False(t, cpu.IsFlagSet(H))

	cpu.Dec_r(&cpu.R.C)

	assert.True(t, cpu.IsFlagSet(H))
}

func TestDec_rCarryFlagUnaffected(t *testing.T) {
	cpu := NewCPU()
	cpu.R.C = 0xF1

	cpu.SetFlag(C)
	cpu.Dec_r(&cpu.R.C)

	assert.True(t, cpu.IsFlagSet(C))
}

//DEC hl
func TestDec_hl(t *testing.T) {
	cpu := NewCPU()
	cpu.LinkMMU(NewMockMMU())
	cpu.R.H = 0x11
	cpu.R.L = 0x1F
	var hlAddr types.Word = types.Word(utils.JoinBytes(cpu.R.H, cpu.R.L))
	cpu.WriteByte(hlAddr, 0x11)

	cpu.Dec_hl()

	var expected byte = 0x10
	assert.Equal(t, cpu.mmu.ReadByte(hlAddr), expected)
}

func TestDec_hlNFlagSet(t *testing.T) {
	cpu := NewCPU()
	cpu.LinkMMU(NewMockMMU())
	cpu.R.H = 0x11
	cpu.R.L = 0x1F
	var hlAddr types.Word = types.Word(utils.JoinBytes(cpu.R.H, cpu.R.L))
	cpu.WriteByte(hlAddr, 0x11)

	assert.False(t, cpu.IsFlagSet(N))

	cpu.Dec_hl()

	assert.True(t, cpu.IsFlagSet(N))
}

func TestDec_hlZeroFlagSet(t *testing.T) {
	cpu := NewCPU()
	cpu.LinkMMU(NewMockMMU())
	cpu.R.H = 0x11
	cpu.R.L = 0x1F
	var hlAddr types.Word = types.Word(utils.JoinBytes(cpu.R.H, cpu.R.L))
	cpu.WriteByte(hlAddr, 0x01)

	assert.False(t, cpu.IsFlagSet(Z))

	cpu.Dec_hl()

	assert.True(t, cpu.IsFlagSet(Z))
}

func TestDec_hlHFlagSetOnBorrow(t *testing.T) {
	cpu := NewCPU()
	cpu.LinkMMU(NewMockMMU())
	cpu.R.H = 0x11
	cpu.R.L = 0x1F
	var hlAddr types.Word = types.Word(utils.JoinBytes(cpu.R.H, cpu.R.L))
	cpu.WriteByte(hlAddr, 0xF0)

	assert.False(t, cpu.IsFlagSet(H))

	cpu.Dec_hl()

	assert.True(t, cpu.IsFlagSet(H))
}

func TestDec_hlCarryFlagUnaffected(t *testing.T) {
	cpu := NewCPU()
	cpu.LinkMMU(NewMockMMU())
	cpu.R.H = 0x11
	cpu.R.L = 0x1F
	var hlAddr types.Word = types.Word(utils.JoinBytes(cpu.R.H, cpu.R.L))
	cpu.WriteByte(hlAddr, 0xF1)

	cpu.SetFlag(C)
	cpu.Dec_hl()

	assert.True(t, cpu.IsFlagSet(C))
}

//SWAP r
func TestSwap_r(t *testing.T) {
	cpu := NewCPU()
	cpu.R.E = 0xD4

	cpu.Swap_r(&cpu.R.E)

	var expected byte = 0x4D
	assert.Equal(t, expected, cpu.R.E)
}

func TestSwap_rFlagsReset(t *testing.T) {
	cpu := NewCPU()
	cpu.R.E = 0xD4

	cpu.SetFlag(N)
	cpu.SetFlag(H)
	cpu.SetFlag(C)

	cpu.Swap_r(&cpu.R.E)

	assert.False(t, cpu.IsFlagSet(N))
	assert.False(t, cpu.IsFlagSet(H))
	assert.False(t, cpu.IsFlagSet(C))
}

func TestSwap_rZeroFlag(t *testing.T) {
	cpu := NewCPU()
	cpu.R.E = 0x00

	assert.False(t, cpu.IsFlagSet(Z))

	cpu.Swap_r(&cpu.R.E)

	assert.True(t, cpu.IsFlagSet(Z))
}

//SWAP hl
func TestSwap_hl(t *testing.T) {
	cpu := NewCPU()
	cpu.LinkMMU(NewMockMMU())
	cpu.R.H, cpu.R.L = 0x01, 0x06
	var hlAddr types.Word = types.Word(utils.JoinBytes(cpu.R.H, cpu.R.L))
	cpu.mmu.WriteByte(hlAddr, 0xD4)

	cpu.Swap_hl()

	var expected byte = 0x4D
	assert.Equal(t, expected, cpu.mmu.ReadByte(hlAddr))
}

func TestSwap_hlFlagsReset(t *testing.T) {
	cpu := NewCPU()
	cpu.LinkMMU(NewMockMMU())
	cpu.R.H, cpu.R.L = 0x01, 0x06
	var hlAddr types.Word = types.Word(utils.JoinBytes(cpu.R.H, cpu.R.L))
	cpu.mmu.WriteByte(hlAddr, 0xD4)

	cpu.SetFlag(N)
	cpu.SetFlag(H)
	cpu.SetFlag(C)

	cpu.Swap_hl()

	assert.False(t, cpu.IsFlagSet(N))
	assert.False(t, cpu.IsFlagSet(H))
	assert.False(t, cpu.IsFlagSet(C))
}

func TestSwap_hlZeroFlag(t *testing.T) {
	cpu := NewCPU()
	cpu.LinkMMU(NewMockMMU())
	cpu.R.H, cpu.R.L = 0x01, 0x06
	var hlAddr types.Word = types.Word(utils.JoinBytes(cpu.R.H, cpu.R.L))
	cpu.mmu.WriteByte(hlAddr, 0x00)

	assert.False(t, cpu.IsFlagSet(Z))

	cpu.Swap_hl()

	assert.True(t, cpu.IsFlagSet(Z))
}

//CPA r
func TestCP_r(t *testing.T) {
	cpu := NewCPU()
	cpu.R.A = 0xAA
	cpu.R.B = 0x24

	cpu.CPA_r(&cpu.R.B)

	//check nothing has changed
	var expectedA byte = 0xAA
	var expectedB byte = 0x24

	assert.Equal(t, cpu.R.A, expectedA)
	assert.Equal(t, cpu.R.B, expectedB)
}

func TestCP_rNFlagSet(t *testing.T) {
	cpu := NewCPU()

	assert.False(t, cpu.IsFlagSet(N))

	cpu.CPA_r(&cpu.R.B)

	assert.True(t, cpu.IsFlagSet(N))
}

func TestCP_rHFlag(t *testing.T) {
	cpu := NewCPU()
	cpu.R.A = 0xA0
	cpu.R.B = 0x01

	assert.False(t, cpu.IsFlagSet(H))

	cpu.CPA_r(&cpu.R.B)

	assert.True(t, cpu.IsFlagSet(H))
}

func TestCP_rCarryFlag(t *testing.T) {
	cpu := NewCPU()
	cpu.R.A = 0xA0
	cpu.R.B = 0x01

	assert.False(t, cpu.IsFlagSet(C))

	cpu.CPA_r(&cpu.R.B)

	assert.True(t, cpu.IsFlagSet(C))
}

func TestCP_rZeroFlag(t *testing.T) {
	cpu := NewCPU()
	cpu.R.A = 0x02
	cpu.R.B = 0x02

	assert.False(t, cpu.IsFlagSet(Z))

	cpu.CPA_r(&cpu.R.B)

	assert.True(t, cpu.IsFlagSet(Z))
}

//CPA hl
func TestCPA_hl(t *testing.T) {
	cpu := NewCPU()
	cpu.LinkMMU(NewMockMMU())
	cpu.R.A = 0xAA
	cpu.R.H, cpu.R.L = 0x00, 0x01
	var hlAddr types.Word = types.Word(utils.JoinBytes(cpu.R.H, cpu.R.L))
	cpu.mmu.WriteByte(hlAddr, 0x24)

	cpu.CPA_hl()

	//check nothing has changed
	var expectedA byte = 0xAA
	var expectedHL byte = 0x24

	assert.Equal(t, cpu.R.A, expectedA)
	assert.Equal(t, cpu.mmu.ReadByte(hlAddr), expectedHL)
}

func TestCPA_hlNFlagSet(t *testing.T) {
	cpu := NewCPU()
	cpu.LinkMMU(NewMockMMU())
	cpu.R.A = 0xAA
	cpu.R.H, cpu.R.L = 0x00, 0x01
	var hlAddr types.Word = types.Word(utils.JoinBytes(cpu.R.H, cpu.R.L))
	cpu.mmu.WriteByte(hlAddr, 0x24)

	assert.False(t, cpu.IsFlagSet(N))

	cpu.CPA_hl()

	assert.True(t, cpu.IsFlagSet(N))
}

func TestCPA_hlHFlag(t *testing.T) {
	cpu := NewCPU()
	cpu.LinkMMU(NewMockMMU())
	cpu.R.A = 0xA0
	cpu.R.H, cpu.R.L = 0x00, 0x01
	var hlAddr types.Word = types.Word(utils.JoinBytes(cpu.R.H, cpu.R.L))
	cpu.mmu.WriteByte(hlAddr, 0x01)

	assert.False(t, cpu.IsFlagSet(H))

	cpu.CPA_hl()

	assert.True(t, cpu.IsFlagSet(H))
}

func TestCPA_hlCarryFlag(t *testing.T) {
	cpu := NewCPU()
	cpu.LinkMMU(NewMockMMU())
	cpu.R.A = 0xA0
	cpu.R.H, cpu.R.L = 0x00, 0x01
	var hlAddr types.Word = types.Word(utils.JoinBytes(cpu.R.H, cpu.R.L))
	cpu.mmu.WriteByte(hlAddr, 0x01)

	assert.False(t, cpu.IsFlagSet(C))

	cpu.CPA_hl()

	assert.True(t, cpu.IsFlagSet(C))
}

func TestCPA_hlZeroFlag(t *testing.T) {
	cpu := NewCPU()
	cpu.LinkMMU(NewMockMMU())
	cpu.R.A = 0x02
	cpu.R.H, cpu.R.L = 0x00, 0x01
	var hlAddr types.Word = types.Word(utils.JoinBytes(cpu.R.H, cpu.R.L))
	cpu.mmu.WriteByte(hlAddr, 0x02)

	assert.False(t, cpu.IsFlagSet(Z))

	cpu.CPA_hl()

	assert.True(t, cpu.IsFlagSet(Z))
}

//CPA r
func TestCPA_n(t *testing.T) {
	cpu := NewCPU()
	cpu.CurrentInstruction = OneOperandsInstr
	cpu.R.A = 0xAA
	cpu.CurrentInstruction.Operands[0] = 0x24

	cpu.CPA_n()

	//check nothing has changed
	var expectedA byte = 0xAA

	assert.Equal(t, cpu.R.A, expectedA)
}

func TestCPA_nNFlagSet(t *testing.T) {
	cpu := NewCPU()
	cpu.CurrentInstruction = OneOperandsInstr
	cpu.R.A = 0xAA
	cpu.CurrentInstruction.Operands[0] = 0x24

	assert.False(t, cpu.IsFlagSet(N))

	cpu.CPA_n()

	assert.True(t, cpu.IsFlagSet(N))
}

func TestCPA_nHFlag(t *testing.T) {
	cpu := NewCPU()
	cpu.CurrentInstruction = OneOperandsInstr
	cpu.R.A = 0xA0
	cpu.CurrentInstruction.Operands[0] = 0x01

	assert.False(t, cpu.IsFlagSet(H))

	cpu.CPA_n()

	assert.True(t, cpu.IsFlagSet(H))
}

func TestCPA_nCarryFlag(t *testing.T) {
	cpu := NewCPU()
	cpu.CurrentInstruction = OneOperandsInstr
	cpu.R.A = 0xA0
	cpu.CurrentInstruction.Operands[0] = 0x01

	assert.False(t, cpu.IsFlagSet(C))

	cpu.CPA_n()

	assert.True(t, cpu.IsFlagSet(C))
}

func TestCPA_nZeroFlag(t *testing.T) {
	cpu := NewCPU()
	cpu.CurrentInstruction = OneOperandsInstr
	cpu.R.A = 0x02
	cpu.CurrentInstruction.Operands[0] = 0x02

	assert.False(t, cpu.IsFlagSet(Z))

	cpu.CPA_n()

	assert.True(t, cpu.IsFlagSet(Z))
}

//CPL
func TestCPL(t *testing.T) {
	cpu := NewCPU()
	cpu.R.A = 0xA9

	cpu.CPL()

	var expected byte = 0x56
	assert.Equal(t, cpu.R.A, expected)
}

func TestCPLFlagsSet(t *testing.T) {
	cpu := NewCPU()
	cpu.R.A = 0xA9

	assert.False(t, cpu.IsFlagSet(N))
	assert.False(t, cpu.IsFlagSet(H))

	cpu.CPL()

	assert.True(t, cpu.IsFlagSet(N))
	assert.True(t, cpu.IsFlagSet(H))
}

func TestCPLFlagsUnaffected(t *testing.T) {
	cpu := NewCPU()
	cpu.R.A = 0xA9

	cpu.SetFlag(Z)
	cpu.SetFlag(C)

	cpu.CPL()

	assert.True(t, cpu.IsFlagSet(Z))
	assert.True(t, cpu.IsFlagSet(C))

	cpu = NewCPU()
	cpu.R.A = 0xA9

	cpu.CPL()

	assert.False(t, cpu.IsFlagSet(Z))
	assert.False(t, cpu.IsFlagSet(C))
}

// RLCA
func TestRLCAWithBit7ToCarry(t *testing.T) {
	cpu := NewCPU()
	cpu.R.A = 0xA5

	assert.False(t, cpu.IsFlagSet(C))

	cpu.RLCA()

	var expectedA byte = 0x4B
	assert.Equal(t, cpu.R.A, expectedA)
	assert.True(t, cpu.IsFlagSet(C))
}

func TestRLCAWithNoBit7ToCarry(t *testing.T) {
	cpu := NewCPU()
	cpu.R.A = 0x15

	assert.False(t, cpu.IsFlagSet(C))

	cpu.RLCA()

	var expectedA byte = 0x2A
	assert.Equal(t, cpu.R.A, expectedA)
	assert.False(t, cpu.IsFlagSet(C))
}

func TestRLCAFlagsReset(t *testing.T) {
	cpu := NewCPU()

	cpu.SetFlag(N)
	cpu.SetFlag(H)

	cpu.RLCA()

	assert.False(t, cpu.IsFlagSet(N))
	assert.False(t, cpu.IsFlagSet(H))
}

func TestRLCAZeroFlag(t *testing.T) {
	cpu := NewCPU()

	cpu.R.A = 0x00

	assert.False(t, cpu.IsFlagSet(Z))

	cpu.RLCA()

	assert.True(t, cpu.IsFlagSet(Z))
}

// RLC r
func TestRLC_rWithBit7ToCarry(t *testing.T) {
	cpu := NewCPU()
	cpu.R.C = 0xA5

	assert.False(t, cpu.IsFlagSet(C))

	cpu.Rlc_r(&cpu.R.C)

	var expectedA byte = 0x4B
	assert.Equal(t, cpu.R.C, expectedA)
	assert.True(t, cpu.IsFlagSet(C))
}

func TestRLC_rWithNoBit7ToCarry(t *testing.T) {
	cpu := NewCPU()
	cpu.R.C = 0x15

	assert.False(t, cpu.IsFlagSet(C))

	cpu.Rlc_r(&cpu.R.C)

	var expectedC byte = 0x2A
	assert.Equal(t, cpu.R.C, expectedC)
	assert.False(t, cpu.IsFlagSet(C))
}

func TestRLC_rFlagsReset(t *testing.T) {
	cpu := NewCPU()

	cpu.SetFlag(N)
	cpu.SetFlag(H)

	cpu.Rlc_r(&cpu.R.C)

	assert.False(t, cpu.IsFlagSet(N))
	assert.False(t, cpu.IsFlagSet(H))
}

func TestRLC_rZeroFlag(t *testing.T) {
	cpu := NewCPU()

	cpu.R.C = 0x00

	assert.False(t, cpu.IsFlagSet(Z))

	cpu.Rlc_r(&cpu.R.C)

	assert.True(t, cpu.IsFlagSet(Z))
}

// RLC hl
func TestRlc_hlWithBit7ToCarry(t *testing.T) {
	cpu := NewCPU()
	cpu.LinkMMU(NewMockMMU())
	cpu.R.H, cpu.R.L = 0x00, 0x01
	hlAddr := types.Word(utils.JoinBytes(cpu.R.H, cpu.R.L))
	cpu.mmu.WriteByte(hlAddr, 0xA5)

	assert.False(t, cpu.IsFlagSet(C))

	cpu.Rlc_hl()

	var expected byte = 0x4B
	assert.Equal(t, cpu.mmu.ReadByte(hlAddr), expected)
	assert.True(t, cpu.IsFlagSet(C))
}

func TestRlc_hlWithNoBit7ToCarry(t *testing.T) {
	cpu := NewCPU()
	cpu.LinkMMU(NewMockMMU())
	cpu.R.H, cpu.R.L = 0x00, 0x01
	hlAddr := types.Word(utils.JoinBytes(cpu.R.H, cpu.R.L))
	cpu.mmu.WriteByte(hlAddr, 0x15)

	assert.False(t, cpu.IsFlagSet(C))

	cpu.Rlc_hl()

	var expected byte = 0x2A
	assert.Equal(t, cpu.mmu.ReadByte(hlAddr), expected)
	assert.False(t, cpu.IsFlagSet(C))
}

func TestRlc_hlFlagsReset(t *testing.T) {
	cpu := NewCPU()
	cpu.LinkMMU(NewMockMMU())

	cpu.SetFlag(N)
	cpu.SetFlag(H)

	cpu.Rlc_hl()

	assert.False(t, cpu.IsFlagSet(N))
	assert.False(t, cpu.IsFlagSet(H))
}

// RRCA
func TestRRCAWithBit0ToCarry(t *testing.T) {
	cpu := NewCPU()
	cpu.R.A = 0xA5

	assert.False(t, cpu.IsFlagSet(C))

	cpu.RRCA()

	var expectedA byte = 0xD2
	assert.Equal(t, cpu.R.A, expectedA)
	assert.True(t, cpu.IsFlagSet(C))
}

func TestRRCAWithNoBit0ToCarry(t *testing.T) {
	cpu := NewCPU()
	cpu.R.A = 0xA2

	assert.False(t, cpu.IsFlagSet(C))

	cpu.RRCA()

	var expectedA byte = 0x51
	assert.Equal(t, cpu.R.A, expectedA)
	assert.False(t, cpu.IsFlagSet(C))
}

func TestRRCAFlagsReset(t *testing.T) {
	cpu := NewCPU()

	cpu.SetFlag(N)
	cpu.SetFlag(H)

	cpu.RRCA()

	assert.False(t, cpu.IsFlagSet(N))
	assert.False(t, cpu.IsFlagSet(H))
}

func TestRRCAZeroFlag(t *testing.T) {
	cpu := NewCPU()

	cpu.R.A = 0x00

	assert.False(t, cpu.IsFlagSet(Z))

	cpu.RRCA()

	assert.True(t, cpu.IsFlagSet(Z))
}

// RRC r
func TestRRC_rWithBit0ToCarry(t *testing.T) {
	cpu := NewCPU()
	cpu.R.C = 0xA5

	assert.False(t, cpu.IsFlagSet(C))

	cpu.Rrc_r(&cpu.R.C)

	var expectedA byte = 0xD2
	assert.Equal(t, cpu.R.C, expectedA)
	assert.True(t, cpu.IsFlagSet(C))
}

func TestRRC_rWithNoBit0ToCarry(t *testing.T) {
	cpu := NewCPU()
	cpu.R.C = 0xA2

	assert.False(t, cpu.IsFlagSet(C))

	cpu.Rrc_r(&cpu.R.C)

	var expectedC byte = 0x51
	assert.Equal(t, cpu.R.C, expectedC)
	assert.False(t, cpu.IsFlagSet(C))
}

func TestRRC_rFlagsReset(t *testing.T) {
	cpu := NewCPU()

	cpu.SetFlag(N)
	cpu.SetFlag(H)

	cpu.Rrc_r(&cpu.R.C)

	assert.False(t, cpu.IsFlagSet(N))
	assert.False(t, cpu.IsFlagSet(H))
}

func TestRRC_rZeroFlag(t *testing.T) {
	cpu := NewCPU()

	cpu.R.C = 0x00

	assert.False(t, cpu.IsFlagSet(Z))

	cpu.Rrc_r(&cpu.R.C)

	assert.True(t, cpu.IsFlagSet(Z))
}

// RRC hl
func TestRrc_hlWithBit0ToCarry(t *testing.T) {
	cpu := NewCPU()
	cpu.LinkMMU(NewMockMMU())
	cpu.R.H, cpu.R.L = 0x00, 0x01
	hlAddr := types.Word(utils.JoinBytes(cpu.R.H, cpu.R.L))
	cpu.mmu.WriteByte(hlAddr, 0xA5)

	assert.False(t, cpu.IsFlagSet(C))

	cpu.Rrc_hl()

	var expected byte = 0xD2
	assert.Equal(t, cpu.mmu.ReadByte(hlAddr), expected)
	assert.True(t, cpu.IsFlagSet(C))
}

func TestRrc_hlWithNoBit0ToCarry(t *testing.T) {
	cpu := NewCPU()
	cpu.LinkMMU(NewMockMMU())
	cpu.R.H, cpu.R.L = 0x00, 0x01
	hlAddr := types.Word(utils.JoinBytes(cpu.R.H, cpu.R.L))
	cpu.mmu.WriteByte(hlAddr, 0xA2)

	assert.False(t, cpu.IsFlagSet(C))

	cpu.Rrc_hl()

	var expected byte = 0x51
	assert.Equal(t, cpu.mmu.ReadByte(hlAddr), expected)
	assert.False(t, cpu.IsFlagSet(C))
}

func TestRrc_hlFlagsReset(t *testing.T) {
	cpu := NewCPU()
	cpu.LinkMMU(NewMockMMU())

	cpu.SetFlag(N)
	cpu.SetFlag(H)

	cpu.Rrc_hl()

	assert.False(t, cpu.IsFlagSet(N))
	assert.False(t, cpu.IsFlagSet(H))
}

func TestRrc_hlZeroFlag(t *testing.T) {
	cpu := NewCPU()
	cpu.LinkMMU(NewMockMMU())
	cpu.R.H, cpu.R.L = 0x00, 0x01
	hlAddr := types.Word(utils.JoinBytes(cpu.R.H, cpu.R.L))
	cpu.mmu.WriteByte(hlAddr, 0x00)

	assert.False(t, cpu.IsFlagSet(Z))

	cpu.Rrc_hl()

	assert.True(t, cpu.IsFlagSet(Z))
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
*/
