package cpu

import (
	"fmt"
	"log"

	"github.com/djhworld/gomeboycolor/constants"
	"github.com/djhworld/gomeboycolor/mmu"
	"github.com/djhworld/gomeboycolor/types"
	"github.com/djhworld/gomeboycolor/utils"
)

const NAME = "CPU"
const PREFIX = NAME + ":"

//flags
const (
	_ = iota
	C
	H
	N
	Z
)

type CurrentInstruction struct {
	*Instruction
	Operands [2]byte
}

type CPUFrame struct {
	PC                      types.Word // Program Counter
	SP                      types.Word // Stack Pointer
	R                       Registers
	InterruptsEnabled       bool
	CurrentInstruction      *CurrentInstruction
	LastInstrCycle          Clock
	PCJumped                bool
	Halted                  bool
	InterruptFlagBeforeHalt byte
}

func (cpu *GbcCPU) GetFrame() *CPUFrame {
	var frame *CPUFrame = new(CPUFrame)
	frame.PC = cpu.PC
	frame.SP = cpu.SP
	frame.R = cpu.R
	frame.InterruptsEnabled = cpu.InterruptsEnabled
	frame.CurrentInstruction = cpu.CurrentInstruction
	frame.LastInstrCycle = cpu.LastInstrCycle
	frame.PCJumped = cpu.PCJumped
	frame.Halted = cpu.Halted
	frame.InterruptFlagBeforeHalt = cpu.InterruptFlagBeforeHalt
	return frame
}

type Registers struct {
	A byte
	B byte
	C byte
	D byte
	E byte
	H byte
	L byte
	F byte // Flags Register
}

func (r Registers) String() string {
	formatByte := func(b byte) string {
		if b < 0x10 {
			return fmt.Sprintf("0x0%X", b)
		}
		return fmt.Sprintf("0x%X", b)
	}
	return fmt.Sprintf("A: %s  B: %s  C: %s  D: %s  E: %s  H: %s  L: %s", formatByte(r.A), formatByte(r.B), formatByte(r.C), formatByte(r.D), formatByte(r.E), formatByte(r.H), formatByte(r.L))
}

//See ZILOG z80 cpu manual p.80  (http://www.zilog.com/docs/z80/um0080.pdf)
type Clock struct {
	M int
	t int
}

func (c *Clock) T() int {
	return c.M * 4
}

func (c *Clock) Reset() {
	c.M = 0
	c.t = 0
}

func (c *Clock) String() string {
	return fmt.Sprintf("[M: %X, T: %X]", c.M, c.t)
}

type GbcCPU struct {
	PC                      types.Word // Program Counter
	SP                      types.Word // Stack Pointer
	R                       Registers
	InterruptsEnabled       bool
	CurrentInstruction      *CurrentInstruction
	LastInstrCycle          Clock
	mmu                     mmu.MemoryMappedUnit
	PCJumped                bool
	Halted                  bool
	InterruptFlagBeforeHalt byte
	Speed                   int
}

func NewCPU(m mmu.MemoryMappedUnit) *GbcCPU {
	cpu := new(GbcCPU)
	cpu.Reset()
	cpu.mmu = m
	log.Println(PREFIX, "Linked CPU to MMU")
	return cpu
}

func (cpu *GbcCPU) Reset() {
	log.Println(PREFIX, "Resetting", NAME)
	cpu.PC = 0
	cpu.SP = 0
	cpu.R.A = 0
	cpu.R.B = 0
	cpu.R.C = 0
	cpu.R.D = 0
	cpu.R.E = 0
	cpu.R.F = 0
	cpu.R.H = 0
	cpu.R.L = 0
	cpu.Speed = 1
	cpu.CurrentInstruction = &CurrentInstruction{Instruction: Instructions[0x00], Operands: [2]byte{}}
	cpu.InterruptsEnabled = true
	cpu.LastInstrCycle.Reset()
	cpu.PCJumped = false
	cpu.Halted = false
}

func (cpu *GbcCPU) FlagsString() string {
	var flags string = ""
	var minus string = "-"

	if cpu.IsFlagSet(Z) {
		flags += "Z"
	} else {
		flags += minus
	}

	if cpu.IsFlagSet(N) {
		flags += "N"
	} else {
		flags += minus
	}

	if cpu.IsFlagSet(H) {
		flags += "H"
	} else {
		flags += minus
	}

	if cpu.IsFlagSet(C) {
		flags += "C"
	} else {
		flags += minus
	}
	return flags
}

func (cpu *GbcCPU) String() string {
	return fmt.Sprint("PC: ", cpu.PC, "  SP: ", cpu.SP, "  ", cpu.R, "  ", cpu.FlagsString(), "  ", cpu.CurrentInstruction)
}

func (cpu *GbcCPU) ResetFlag(flag int) {
	switch flag {
	case Z:
		cpu.R.F = cpu.resetBit(7, cpu.R.F)
	case N:
		cpu.R.F = cpu.resetBit(6, cpu.R.F)
	case H:
		cpu.R.F = cpu.resetBit(5, cpu.R.F)
	case C:
		cpu.R.F = cpu.resetBit(4, cpu.R.F)
	default:
		log.Fatalf(PREFIX+" Unknown flag %c", flag)
	}
	cpu.R.F &= 0xF0
}

func (cpu *GbcCPU) SetFlag(flag int) {
	switch flag {
	case Z:
		cpu.R.F = cpu.setBit(7, cpu.R.F)
	case N:
		cpu.R.F = cpu.setBit(6, cpu.R.F)
	case H:
		cpu.R.F = cpu.setBit(5, cpu.R.F)
	case C:
		cpu.R.F = cpu.setBit(4, cpu.R.F)
	default:
		log.Fatalf(PREFIX+" Unknown flag %c", flag)
	}
	cpu.R.F &= 0xF0
}

func (cpu *GbcCPU) IsFlagSet(flag int) bool {
	switch flag {
	case Z:
		return cpu.R.F&0x80 == 0x80
	case N:
		return cpu.R.F&0x40 == 0x40
	case H:
		return cpu.R.F&0x20 == 0x20
	case C:
		return cpu.R.F&0x10 == 0x10
	default:
		log.Fatalf(PREFIX+" Unknown flag %c", flag)
	}
	return false
}

func (cpu *GbcCPU) IncrementPC(by int) {
	cpu.PC += types.Word(by)
}

func (cpu *GbcCPU) Step() int {
	cpu.LastInstrCycle.Reset()
	var opcode byte

	if !cpu.Halted {
		cpu.CheckForInterrupts()
		opcode = cpu.ReadByte(cpu.PC)

		if opcode == 0xCB {
			cpu.IncrementPC(1)
			opcode = cpu.ReadByte(cpu.PC)
			cpu.Compile(InstructionsCB[opcode])
		} else {
			cpu.Compile(Instructions[opcode])
		}

		cpu.CurrentInstruction.Execute(cpu)

		//this is put in place to check whether the PC has been altered by an instruction. If it has then don't
		//do any incrementing
		if cpu.PCJumped == false {
			cpu.IncrementPC(cpu.CurrentInstruction.OperandsSize + 1)
		}

		cpu.PCJumped = false

		//calculate cycles
		cpu.LastInstrCycle.M += cpu.CurrentInstruction.Cycles
	} else {
		iflagnow := cpu.mmu.ReadByte(constants.INTERRUPT_FLAG_ADDR)

		//if IF flag has changed then the cpu should resume...
		if cpu.InterruptFlagBeforeHalt != iflagnow {
			cpu.Halted = false
		}

		//Halt consumes 1 cpu cycle
		cpu.LastInstrCycle.M = 1
	}

	return cpu.LastInstrCycle.M
}

func (cpu *GbcCPU) CheckForInterrupts() bool {
	if cpu.InterruptsEnabled {
		var ie byte = cpu.mmu.ReadByte(constants.INTERRUPT_ENABLED_FLAG_ADDR)
		var iflag byte = cpu.mmu.ReadByte(constants.INTERRUPT_FLAG_ADDR)
		var interrupt byte = iflag & ie
		if interrupt != 0x00 {
			switch {
			case interrupt&constants.V_BLANK_IRQ == constants.V_BLANK_IRQ:
				cpu.mmu.WriteByte(constants.INTERRUPT_FLAG_ADDR, iflag&0xFE)
				cpu.pushWordToStack(cpu.PC)
				cpu.PC = types.Word(constants.V_BLANK_IR_ADDR)
				cpu.InterruptsEnabled = false
				return true
			case interrupt&constants.LCD_IRQ == constants.LCD_IRQ:
				cpu.mmu.WriteByte(constants.INTERRUPT_FLAG_ADDR, iflag&0xFD)
				cpu.pushWordToStack(cpu.PC)
				cpu.PC = types.Word(constants.LCD_IR_ADDR)
				cpu.InterruptsEnabled = false
				return true
			case interrupt&constants.TIMER_OVERFLOW_IRQ == constants.TIMER_OVERFLOW_IRQ:
				cpu.mmu.WriteByte(constants.INTERRUPT_FLAG_ADDR, iflag&0xFB)
				cpu.pushWordToStack(cpu.PC)
				cpu.PC = types.Word(constants.TIMER_OVERFLOW_IR_ADDR)
				cpu.InterruptsEnabled = false
				return true
			case interrupt&constants.JOYP_HILO_IRQ == constants.JOYP_HILO_IRQ:
				log.Println("JOYP!")
				cpu.mmu.WriteByte(constants.INTERRUPT_FLAG_ADDR, iflag&0xEF)
				cpu.pushWordToStack(cpu.PC)
				cpu.PC = types.Word(constants.JOYP_HILO_IR_ADDR)
				cpu.InterruptsEnabled = false
				return true
			default:
				log.Fatalf("Unknown interrupt = %d", interrupt)
			}
		}
	}

	return false
}

//Checks to see if the CPU speed should change to double (CGB only)
func (cpu *GbcCPU) SetCPUSpeed() {
	var speedPrepRegister byte = cpu.mmu.ReadByte(mmu.CGB_DOUBLE_SPEED_PREP_REG)
	if speedPrepRegister&0x01 == 0x01 {
		switch cpu.Speed {
		case 2:
			cpu.Speed = 1
			cpu.mmu.WriteByte(mmu.CGB_DOUBLE_SPEED_PREP_REG, 0x00)
		case 1:
			cpu.Speed = 2
			cpu.mmu.WriteByte(mmu.CGB_DOUBLE_SPEED_PREP_REG, 0x80)
		default:
			panic(fmt.Sprint("Unsupported CPU speed ", cpu.Speed, " this should not happen!"))
		}
		log.Printf("CPU: Setting CPU speed to %dx speed", cpu.Speed)
	}
}

func (cpu *GbcCPU) Compile(instruction *Instruction) {
	cpu.CurrentInstruction.Instruction = instruction
	switch instruction.OperandsSize {
	case 1:
		cpu.CurrentInstruction.Operands[0] = cpu.mmu.ReadByte(cpu.PC + 1)
	case 2:
		cpu.CurrentInstruction.Operands[0] = cpu.mmu.ReadByte(cpu.PC + 1)
		cpu.CurrentInstruction.Operands[1] = cpu.mmu.ReadByte(cpu.PC + 2)
	}
}

func (cpu *GbcCPU) pushByteToStack(b byte) {
	cpu.SP--
	cpu.WriteByte(cpu.SP, b)
}

func (cpu *GbcCPU) pushWordToStack(word types.Word) {
	hs, ls := utils.SplitIntoBytes(uint16(word))
	cpu.pushByteToStack(hs)
	cpu.pushByteToStack(ls)
}

func (cpu *GbcCPU) popByteFromStack() byte {
	var b byte = cpu.ReadByte(cpu.SP)
	cpu.SP++
	return b
}

func (cpu *GbcCPU) popWordFromStack() types.Word {
	ls := cpu.popByteFromStack()
	hs := cpu.popByteFromStack()

	return types.Word(utils.JoinBytes(hs, ls))
}

func (cpu *GbcCPU) ReadByte(addr types.Word) byte {
	return cpu.mmu.ReadByte(addr)
}

func (cpu *GbcCPU) WriteByte(addr types.Word, value byte) {
	cpu.mmu.WriteByte(addr, value)
}

// INSTRUCTION HELPERS
//-----------------------------------------------------------------------

//General add function
//side effect is flags are set on the cpu accordingly
func (cpu *GbcCPU) addBytes(a, b byte) byte {
	var calculation byte = a + b

	cpu.ResetFlag(N)

	if calculation == 0x00 {
		cpu.SetFlag(Z)
	} else {
		cpu.ResetFlag(Z)
	}

	if (calculation^b^a)&0x10 == 0x10 {
		cpu.SetFlag(H)
	} else {
		cpu.ResetFlag(H)
	}

	if calculation < a {
		cpu.SetFlag(C)
	} else {
		cpu.ResetFlag(C)
	}

	return calculation
}

//General add words function
//side effect is flags are set on the cpu accordingly
func (cpu *GbcCPU) addWords(a, b types.Word) types.Word {
	var calculation types.Word = a + b

	if calculation < a {
		cpu.SetFlag(C)
	} else {
		cpu.ResetFlag(C)
	}

	if (calculation^b^a)&0x1000 == 0x1000 {
		cpu.SetFlag(H)
	} else {
		cpu.ResetFlag(H)
	}

	cpu.ResetFlag(N)

	return calculation
}

//General subtraction function
//side effect is flags are set on the cpu accordingly
func (cpu *GbcCPU) subBytes(a, b byte) byte {

	cpu.SetFlag(N)

	if (int(a) & 0xF) < (int(b) & 0xF) {
		cpu.SetFlag(H)
	} else {
		cpu.ResetFlag(H)
	}

	if (int(a) & 0xFF) < (int(b) & 0xFF) {
		cpu.SetFlag(C)
	} else {
		cpu.ResetFlag(C)
	}

	var calculation byte = a - b
	if calculation == 0x00 {
		cpu.SetFlag(Z)
	} else {
		cpu.ResetFlag(Z)
	}

	return calculation
}

//General and function
//side effect is flags are set on the cpu accordingly
func (cpu *GbcCPU) andBytes(a, b byte) byte {
	var calculation byte = a & b

	cpu.SetFlag(H)
	cpu.ResetFlag(N)
	cpu.ResetFlag(C)

	if calculation == 0x00 {
		cpu.SetFlag(Z)
	} else {
		cpu.ResetFlag(Z)
	}

	return calculation
}

//General or function
//side effect is flags are set on the cpu accordingly
func (cpu *GbcCPU) orBytes(a, b byte) byte {
	var calculation byte = a | b

	cpu.ResetFlag(H)
	cpu.ResetFlag(N)
	cpu.ResetFlag(C)

	if calculation == 0x00 {
		cpu.SetFlag(Z)
	} else {
		cpu.ResetFlag(Z)
	}

	return calculation
}

//General xor function
//side effect is flags are set on the cpu accordingly
func (cpu *GbcCPU) xorBytes(a, b byte) byte {
	var calculation byte = a ^ b

	cpu.ResetFlag(H)
	cpu.ResetFlag(N)
	cpu.ResetFlag(C)

	if calculation == 0x00 {
		cpu.SetFlag(Z)
	} else {
		cpu.ResetFlag(Z)
	}

	return calculation
}

//General increment function
//Side effect is CPU flags get set accordingly
func (cpu *GbcCPU) incByte(value byte) byte {
	var calculation byte = value + 0x01

	//N should be reset
	cpu.ResetFlag(N)

	//set zero flag
	if calculation == 0 {
		cpu.SetFlag(Z)
	} else {
		cpu.ResetFlag(Z)
	}

	//set half carry flag
	if (calculation^0x01^value)&0x10 == 0x10 {
		cpu.SetFlag(H)
	} else {
		cpu.ResetFlag(H)
	}

	return calculation
}

//General decrement function
//Side effect is CPU flags get set accordingly
func (cpu *GbcCPU) decByte(a byte) byte {
	var calculation byte = a - 1
	cpu.SetFlag(N)

	if calculation == 0 {
		cpu.SetFlag(Z)
	} else {
		cpu.ResetFlag(Z)
	}

	//Set half carry flag if needed
	if (calculation^0x01^a)&0x10 == 0x10 {
		cpu.SetFlag(H)
	} else {
		cpu.ResetFlag(H)
	}

	return calculation
}

//General swap function
//Side effect is CPU flags get set accordingly
func (cpu *GbcCPU) swapByte(a byte) byte {
	var calculation byte = utils.SwapNibbles(a)
	cpu.ResetFlag(N)
	cpu.ResetFlag(H)
	cpu.ResetFlag(C)

	if calculation == 0x00 {
		cpu.SetFlag(Z)
	} else {
		cpu.ResetFlag(Z)
	}

	return calculation
}

//General bit - test function
//Side effect is CPU flags get set accordingly
func (cpu *GbcCPU) bitTest(bit byte, a byte) {
	if (a >> bit & 1) == 0 {
		cpu.SetFlag(Z)
	} else {
		cpu.ResetFlag(Z)
	}

	cpu.ResetFlag(N)
	cpu.SetFlag(H)
}

//General set bit function
func (cpu *GbcCPU) setBit(bit byte, a byte) byte {
	return a | (1 << uint(bit))
}

//General reset bit function
func (cpu *GbcCPU) resetBit(bit byte, a byte) byte {
	return a & ^(1 << uint(bit))
}

// INSTRUCTIONS START
//-----------------------------------------------------------------------

//NOP
//No operation
func (cpu *GbcCPU) NOP() {

}

//HALT
//Halt CPU
func (cpu *GbcCPU) HALT() {
	//get and store the state of the IF FLAG now so we know when it changes
	cpu.InterruptFlagBeforeHalt = cpu.mmu.ReadByte(constants.INTERRUPT_FLAG_ADDR)
	cpu.Halted = true
}

//STOP
func (cpu *GbcCPU) Stop() {
	log.Println("CPU: Stopping...")
	//After a stop instruction is executed, CGB hardware should check to see if the CPU speed should change
	cpu.SetCPUSpeed()
}

//DI
//Disable interrupts
func (cpu *GbcCPU) DI() {
	cpu.InterruptsEnabled = false
}

//EI
//Enable interrupts
func (cpu *GbcCPU) EI() {
	cpu.InterruptsEnabled = true
}

//LD r,n
//Load value (n) from memory address in the PC into register (r) and increment PC by 1
func (cpu *GbcCPU) LDrn(r *byte) {
	*r = cpu.CurrentInstruction.Operands[0]
}

//LD r,r
//Load value from register (r2) into register (r1)
func (cpu *GbcCPU) LDrr(r1 *byte, r2 *byte) {
	*r1 = *r2
}

//LD rr, r
//Load value from register (r) into memory address located at register pair (RR)
func (cpu *GbcCPU) LDrr_r(hs *byte, ls *byte, r *byte) {
	var RR types.Word = types.Word(utils.JoinBytes(*hs, *ls))
	cpu.WriteByte(RR, *r)
}

//LD r, rr
//Load value from memory address located in register pair (RR) into register (r)
func (cpu *GbcCPU) LDr_rr(hs *byte, ls *byte, r *byte) {
	var RR types.Word = types.Word(utils.JoinBytes(*hs, *ls))
	*r = cpu.ReadByte(RR)
}

//LD nn,r
//Load value from register (r) and put it in memory address (nn) taken from the next 2 bytes of memory from the PC. Increment the PC by 2
func (cpu *GbcCPU) LDnn_r(r *byte) {
	var ls byte = cpu.CurrentInstruction.Operands[0]
	var hs byte = cpu.CurrentInstruction.Operands[1]
	var resultAddr types.Word = types.Word(utils.JoinBytes(hs, ls))
	cpu.WriteByte(resultAddr, *r)
}

//LD r, nn
//Load the value in memory address defined from the next two bytes relative to the PC and store it in register (r). Increment the PC by 2
func (cpu *GbcCPU) LDr_nn(r *byte) {
	var ls byte = cpu.CurrentInstruction.Operands[0]
	var hs byte = cpu.CurrentInstruction.Operands[1]
	var nn types.Word = types.Word(utils.JoinBytes(hs, ls))
	*r = cpu.ReadByte(nn)
}

//LD (HL),n
//Load the value (n) from the memory address in the PC and put it in the memory address designated by register pair (HL)
func (cpu *GbcCPU) LDhl_n() {
	var HL types.Word = types.Word(utils.JoinBytes(cpu.R.H, cpu.R.L))
	var value byte = cpu.CurrentInstruction.Operands[0]
	cpu.WriteByte(HL, value)
}

//LD r,(C)
//Load the value from memory addressed 0xFF00 + value in register C. Store it in register (r)
func (cpu *GbcCPU) LDr_ffplusc(r *byte) {
	var valueAddr types.Word = 0xFF00 + types.Word(cpu.R.C)
	*r = cpu.ReadByte(valueAddr)
}

//LD (C),r
//Load the value from register (r) and store it in memory addressed 0xFF00 + value in register C.
func (cpu *GbcCPU) LDffplusc_r(r *byte) {
	var valueAddr types.Word = 0xFF00 + types.Word(cpu.R.C)
	cpu.WriteByte(valueAddr, *r)
}

//LDD r, (HL)
//Load the value from memory addressed in register pair (HL) and store it in register R. Decrement the HL registers
func (cpu *GbcCPU) LDDr_hl(r *byte) {
	var HL types.Word = types.Word(utils.JoinBytes(cpu.R.H, cpu.R.L))
	*r = cpu.ReadByte(HL)

	HL -= 1

	cpu.R.H, cpu.R.L = utils.SplitIntoBytes(uint16(HL))
}

//LDD (HL), r
//Load the value in register (r) and store in memory addressed in register pair (HL). Decrement the HL registers
func (cpu *GbcCPU) LDDhl_r(r *byte) {
	var HL types.Word = types.Word(utils.JoinBytes(cpu.R.H, cpu.R.L))
	cpu.WriteByte(HL, *r)

	HL -= 1

	cpu.R.H, cpu.R.L = utils.SplitIntoBytes(uint16(HL))
}

//LDI r, (HL)
//Load the value from memory addressed in register pair (HL) and store it in register R. Increment the HL registers
func (cpu *GbcCPU) LDIr_hl(r *byte) {
	var HL types.Word = types.Word(utils.JoinBytes(cpu.R.H, cpu.R.L))
	*r = cpu.ReadByte(HL)

	HL += 1

	cpu.R.H, cpu.R.L = utils.SplitIntoBytes(uint16(HL))
}

//LDI (HL), r
//Load the value in register (r) and store in memory addressed in register pair (HL). Increment the HL registers
func (cpu *GbcCPU) LDIhl_r(r *byte) {
	var HL types.Word = types.Word(utils.JoinBytes(cpu.R.H, cpu.R.L))
	cpu.WriteByte(HL, *r)

	HL += 1

	cpu.R.H, cpu.R.L = utils.SplitIntoBytes(uint16(HL))
}

//LDH n, r
func (cpu *GbcCPU) LDHn_r(r *byte) {
	var n byte = cpu.CurrentInstruction.Operands[0]
	cpu.WriteByte(types.Word(0xFF00)+types.Word(n), *r)
}

//LDH r, n
//Load value (n) in register (r) and store it in memory address FF00+PC. Increment PC by 1
func (cpu *GbcCPU) LDHr_n(r *byte) {
	var n byte = cpu.CurrentInstruction.Operands[0]
	*r = cpu.ReadByte((types.Word(0xFF00) + types.Word(n)))
}

//LD n, nn
func (cpu *GbcCPU) LDn_nn(r1, r2 *byte) {
	var ls byte = cpu.CurrentInstruction.Operands[0]
	var hs byte = cpu.CurrentInstruction.Operands[1]

	//LS nibble first
	*r1 = hs
	*r2 = ls
}

//LD SP, nn
func (cpu *GbcCPU) LDSP_nn() {
	var ls byte = cpu.CurrentInstruction.Operands[0]
	var hs byte = cpu.CurrentInstruction.Operands[1]

	cpu.SP = types.Word(utils.JoinBytes(hs, ls))
}

//LD nn, SP
func (cpu *GbcCPU) LDnn_SP() {
	var ls byte = cpu.CurrentInstruction.Operands[0]
	var hs byte = cpu.CurrentInstruction.Operands[1]
	var addr types.Word = types.Word(utils.JoinBytes(hs, ls))

	cpu.WriteByte(addr+1, byte(cpu.SP&0xFF00>>8))
	cpu.WriteByte(addr, byte(cpu.SP&0x00FF))
}

//LD SP, rr
func (cpu *GbcCPU) LDSP_hl() {
	cpu.SP = types.Word(utils.JoinBytes(cpu.R.H, cpu.R.L))
}

//LDHL SP, n
func (cpu *GbcCPU) LDHLSP_n() {
	var n byte = cpu.CurrentInstruction.Operands[0]

	var HL types.Word

	if n > 127 {
		HL = cpu.SP - types.Word(-n)
	} else {
		HL = cpu.SP + types.Word(n)
	}

	var check types.Word = types.Word(cpu.SP ^ types.Word(n) ^ ((cpu.SP + types.Word(n)) & 0xffff))

	//set carry flag
	if (check & 0x100) == 0x100 {
		cpu.SetFlag(C)
	} else {
		cpu.ResetFlag(C)
	}

	if (check & 0x10) == 0x10 {
		cpu.SetFlag(H)
	} else {
		cpu.ResetFlag(H)
	}

	//reset flags
	cpu.ResetFlag(Z)
	cpu.ResetFlag(N)

	cpu.R.H, cpu.R.L = utils.SplitIntoBytes(uint16(HL))
}

//PUSH nn
//Push register pair nn onto the stack and decrement the SP twice
func (cpu *GbcCPU) Push_nn(r1, r2 *byte) {
	word := types.Word(utils.JoinBytes(*r1, *r2))
	cpu.pushWordToStack(word)
}

//POP nn
//Pop the stack twice onto register pair nn
func (cpu *GbcCPU) Pop_nn(r1, r2 *byte) {
	*r1, *r2 = utils.SplitIntoBytes(uint16(cpu.popWordFromStack()))
}

//POP AF
//Pop the stack twice onto register pair AF
func (cpu *GbcCPU) Pop_AF() {
	cpu.R.A, cpu.R.F = utils.SplitIntoBytes(uint16(cpu.popWordFromStack()))
	cpu.R.F &= 0xF0
}

//ADD A,r
//Add the value in register (r) to register A
func (cpu *GbcCPU) AddA_r(r *byte) {
	cpu.R.A = cpu.addBytes(cpu.R.A, *r)
}

//ADD A,(HL)
//Add the value in memory addressed in register pair (HL) to register A
func (cpu *GbcCPU) AddA_hl() {
	var HL types.Word = types.Word(utils.JoinBytes(cpu.R.H, cpu.R.L))
	var value byte = cpu.ReadByte(HL)

	cpu.R.A = cpu.addBytes(cpu.R.A, value)
}

//ADD A,n
//Add the value in memory addressed PC to register A. Increment the PC by 1
func (cpu *GbcCPU) AddA_n() {
	var value byte = cpu.CurrentInstruction.Operands[0]
	cpu.R.A = cpu.addBytes(cpu.R.A, value)
}

//ADDC A,r
func (cpu *GbcCPU) AddCA_r(r *byte) {
	var carry int = 0
	if cpu.IsFlagSet(C) {
		carry = 1
	}

	cpu.ResetFlag(N)

	if ((int(cpu.R.A) & 0xF) + (int(*r) & 0xF) + carry) > 0xF {
		cpu.SetFlag(H)
	} else {
		cpu.ResetFlag(H)
	}

	if ((int(cpu.R.A) & 0xFF) + (int(*r) & 0xFF) + carry) > 0xFF {
		cpu.SetFlag(C)
	} else {
		cpu.ResetFlag(C)
	}

	cpu.R.A += *r + byte(carry)

	if cpu.R.A == 0 {
		cpu.SetFlag(Z)
	} else {
		cpu.ResetFlag(Z)
	}
}

//ADDC A,(HL)
func (cpu *GbcCPU) AddCA_hl() {
	var HL types.Word = types.Word(utils.JoinBytes(cpu.R.H, cpu.R.L))
	var hlValue byte = cpu.ReadByte(HL)

	var carry int = 0
	if cpu.IsFlagSet(C) {
		carry = 1
	}

	cpu.ResetFlag(N)

	if ((int(cpu.R.A) & 0xF) + (int(hlValue) & 0xF) + carry) > 0xF {
		cpu.SetFlag(H)
	} else {
		cpu.ResetFlag(H)
	}

	if ((int(cpu.R.A) & 0xFF) + (int(hlValue) & 0xFF) + carry) > 0xFF {
		cpu.SetFlag(C)
	} else {
		cpu.ResetFlag(C)
	}

	cpu.R.A += hlValue + byte(carry)

	if cpu.R.A == 0 {
		cpu.SetFlag(Z)
	} else {
		cpu.ResetFlag(Z)
	}
}

//ADDC A,n
func (cpu *GbcCPU) AddCA_n() {
	var value byte = cpu.CurrentInstruction.Operands[0]
	var carry int = 0
	if cpu.IsFlagSet(C) {
		carry = 1
	}

	cpu.ResetFlag(N)

	if ((int(cpu.R.A) & 0xF) + (int(value) & 0xF) + carry) > 0xF {
		cpu.SetFlag(H)
	} else {
		cpu.ResetFlag(H)
	}

	if ((int(cpu.R.A) & 0xFF) + (int(value) & 0xFF) + carry) > 0xFF {
		cpu.SetFlag(C)
	} else {
		cpu.ResetFlag(C)
	}

	cpu.R.A += value + byte(carry)

	if cpu.R.A == 0 {
		cpu.SetFlag(Z)
	} else {
		cpu.ResetFlag(Z)
	}
}

//SUB A,r
func (cpu *GbcCPU) SubA_r(r *byte) {
	cpu.R.A = cpu.subBytes(cpu.R.A, *r)
}

//SUB A,hl
func (cpu *GbcCPU) SubA_hl() {
	var HL types.Word = types.Word(utils.JoinBytes(cpu.R.H, cpu.R.L))
	var value byte = cpu.ReadByte(HL)

	cpu.R.A = cpu.subBytes(cpu.R.A, value)
}

//SUB A,n
func (cpu *GbcCPU) SubA_n() {
	var value byte = cpu.CurrentInstruction.Operands[0]
	cpu.R.A = cpu.subBytes(cpu.R.A, value)
}

//SBC A,r
func (cpu *GbcCPU) SubAC_r(r *byte) {
	var un int = int(*r) & 0xff
	var tmpa int = int(cpu.R.A) & 0xff
	var ua int = int(cpu.R.A) & 0xff
	ua -= un
	if cpu.IsFlagSet(C) {
		ua -= 1
	}

	if ua < 0 {
		cpu.R.F = 0x50
	} else {
		cpu.R.F = 0x40
	}
	ua &= 0xff
	if ua == 0 {
		cpu.SetFlag(Z)
	} else {
		cpu.ResetFlag(Z)
	}

	if ((ua ^ un ^ tmpa) & 0x10) == 0x10 {
		cpu.SetFlag(H)
	} else {
		cpu.ResetFlag(H)
	}

	cpu.R.A = byte(ua)
}

//SBC A, (HL)
func (cpu *GbcCPU) SubAC_hl() {
	var HL types.Word = types.Word(utils.JoinBytes(cpu.R.H, cpu.R.L))
	var value byte = cpu.ReadByte(HL)

	var un int = int(value) & 0xff
	var tmpa int = int(cpu.R.A) & 0xff
	var ua int = int(cpu.R.A) & 0xff

	ua -= un

	if cpu.IsFlagSet(C) {
		ua -= 1
	}

	if ua < 0 {
		cpu.R.F = 0x50
	} else {
		cpu.R.F = 0x40
	}

	ua &= 0xff
	if ua == 0 {
		cpu.SetFlag(Z)
	} else {
		cpu.ResetFlag(Z)
	}

	if ((ua ^ un ^ tmpa) & 0x10) == 0x10 {
		cpu.SetFlag(H)
	} else {
		cpu.ResetFlag(H)
	}

	cpu.R.A = byte(ua)
}

//SBC A, n
func (cpu *GbcCPU) SubAC_n() {
	var value byte = cpu.CurrentInstruction.Operands[0]
	var un int = int(value) & 0xff
	var tmpa int = int(cpu.R.A) & 0xff
	var ua int = int(cpu.R.A) & 0xff

	ua -= un

	if cpu.IsFlagSet(C) {
		ua -= 1
	}

	if ua < 0 {
		cpu.R.F = 0x50
	} else {
		cpu.R.F = 0x40
	}

	ua &= 0xff
	if ua == 0 {
		cpu.SetFlag(Z)
	} else {
		cpu.ResetFlag(Z)
	}

	if ((ua ^ un ^ tmpa) & 0x10) == 0x10 {
		cpu.SetFlag(H)
	} else {
		cpu.ResetFlag(H)
	}

	cpu.R.A = byte(ua)
}

//AND A, r
func (cpu *GbcCPU) AndA_r(r *byte) {
	cpu.R.A = cpu.andBytes(cpu.R.A, *r)
}

//AND A, (HL)
func (cpu *GbcCPU) AndA_hl() {
	var HL types.Word = types.Word(utils.JoinBytes(cpu.R.H, cpu.R.L))
	var value byte = cpu.ReadByte(HL)
	cpu.R.A = cpu.andBytes(cpu.R.A, value)
}

//AND A, n
func (cpu *GbcCPU) AndA_n() {
	var value byte = cpu.CurrentInstruction.Operands[0]
	cpu.R.A = cpu.andBytes(cpu.R.A, value)
}

//OR A, r
func (cpu *GbcCPU) OrA_r(r *byte) {
	cpu.R.A = cpu.orBytes(cpu.R.A, *r)
}

//OR A, (HL)
func (cpu *GbcCPU) OrA_hl() {
	var HL types.Word = types.Word(utils.JoinBytes(cpu.R.H, cpu.R.L))
	var value byte = cpu.ReadByte(HL)
	cpu.R.A = cpu.orBytes(cpu.R.A, value)
}

//OR A, n
func (cpu *GbcCPU) OrA_n() {
	var value byte = cpu.CurrentInstruction.Operands[0]
	cpu.R.A = cpu.orBytes(cpu.R.A, value)
}

//XOR A, r
func (cpu *GbcCPU) XorA_r(r *byte) {
	cpu.R.A = cpu.xorBytes(cpu.R.A, *r)
}

//XOR A, (HL)
func (cpu *GbcCPU) XorA_hl() {
	var HL types.Word = types.Word(utils.JoinBytes(cpu.R.H, cpu.R.L))
	var value byte = cpu.ReadByte(HL)
	cpu.R.A = cpu.xorBytes(cpu.R.A, value)
}

//XOR A, n
func (cpu *GbcCPU) XorA_n() {
	var value byte = cpu.CurrentInstruction.Operands[0]
	cpu.R.A = cpu.xorBytes(cpu.R.A, value)
}

//CP A, r
func (cpu *GbcCPU) CPA_r(r *byte) {
	cpu.subBytes(cpu.R.A, *r)
}

//CP A, (HL)
func (cpu *GbcCPU) CPA_hl() {
	var hlAddr types.Word = types.Word(utils.JoinBytes(cpu.R.H, cpu.R.L))
	var hlValue byte = cpu.ReadByte(hlAddr)
	cpu.subBytes(cpu.R.A, hlValue)
}

//CP A, n
func (cpu *GbcCPU) CPA_n() {
	var value byte = cpu.CurrentInstruction.Operands[0]
	cpu.subBytes(cpu.R.A, value)
}

//INC r
func (cpu *GbcCPU) Inc_r(r *byte) {
	*r = cpu.incByte(*r)
}

//INC (HL)
func (cpu *GbcCPU) Inc_hl() {
	var hlAddr types.Word = types.Word(utils.JoinBytes(cpu.R.H, cpu.R.L))
	var hlValue byte = cpu.ReadByte(hlAddr)
	var result byte = cpu.incByte(hlValue)
	cpu.WriteByte(hlAddr, result)
}

//DEC r
func (cpu *GbcCPU) Dec_r(r *byte) {
	*r = cpu.decByte(*r)
}

//DEC (HL)
func (cpu *GbcCPU) Dec_hl() {
	var hlAddr types.Word = types.Word(utils.JoinBytes(cpu.R.H, cpu.R.L))
	var hlValue byte = cpu.ReadByte(hlAddr)
	var result byte = cpu.decByte(hlValue)
	cpu.WriteByte(hlAddr, result)
}

// --------------- 16 bit operations ---------------
//ADD HL,rr
func (cpu *GbcCPU) Addhl_rr(r1, r2 *byte) {
	var HL types.Word = types.Word(utils.JoinBytes(cpu.R.H, cpu.R.L))
	var RR types.Word = types.Word(utils.JoinBytes(*r1, *r2))
	var result types.Word = cpu.addWords(HL, RR)
	cpu.R.H, cpu.R.L = utils.SplitIntoBytes(uint16(result))
}

//ADD HL,SP
func (cpu *GbcCPU) Addhl_sp() {
	var HL types.Word = types.Word(utils.JoinBytes(cpu.R.H, cpu.R.L))
	var result types.Word = cpu.addWords(HL, cpu.SP)
	cpu.R.H, cpu.R.L = utils.SplitIntoBytes(uint16(result))
}

//ADD SP,n
func (cpu *GbcCPU) Addsp_n() {
	var n byte = cpu.CurrentInstruction.Operands[0]

	var calculation types.Word

	// immediate value is signed
	if n > 127 {
		calculation = cpu.SP - types.Word(-n)
	} else {
		calculation = cpu.SP + types.Word(n)
	}

	var check types.Word = types.Word(cpu.SP ^ types.Word(n) ^ ((cpu.SP + types.Word(n)) & 0xffff))

	//set carry flag
	if (check & 0x100) == 0x100 {
		cpu.SetFlag(C)
	} else {
		cpu.ResetFlag(C)
	}

	if (check & 0x10) == 0x10 {
		cpu.SetFlag(H)
	} else {
		cpu.ResetFlag(H)
	}

	//reset flags
	cpu.ResetFlag(Z)
	cpu.ResetFlag(N)

	cpu.SP = calculation
}

//INC rr
func (cpu *GbcCPU) Inc_rr(r1, r2 *byte) {
	var RR types.Word = types.Word(utils.JoinBytes(*r1, *r2))
	RR += 1
	*r1, *r2 = utils.SplitIntoBytes(uint16(RR))
}

//INC SP
func (cpu *GbcCPU) Inc_sp() {
	cpu.SP = (cpu.SP + 1) & 0xFFFF
}

//DEC rr
func (cpu *GbcCPU) Dec_rr(r1, r2 *byte) {
	var RR types.Word = types.Word(utils.JoinBytes(*r1, *r2))
	RR -= 1
	*r1, *r2 = utils.SplitIntoBytes(uint16(RR))
}

//DEC SP
func (cpu *GbcCPU) Dec_sp() {
	cpu.SP = (cpu.SP - 1) & 0xFFFF
}

//CPL
func (cpu *GbcCPU) CPL() {
	cpu.R.A = ^cpu.R.A
	cpu.SetFlag(N)
	cpu.SetFlag(H)
}

//CCF
func (cpu *GbcCPU) CCF() {
	if cpu.IsFlagSet(C) {
		cpu.ResetFlag(C)
	} else {
		cpu.SetFlag(C)
	}

	cpu.ResetFlag(N)
	cpu.ResetFlag(H)
}

//SCF
func (cpu *GbcCPU) SCF() {
	cpu.SetFlag(C)

	cpu.ResetFlag(N)
	cpu.ResetFlag(H)
}

//DAA - this instruction was a complete PITA to implement
//thankfully DParrot from here http://forums.nesdev.com/viewtopic.php?t=9088
//provided a correct solution that passes the blargg tests
func (cpu *GbcCPU) Daa() {
	var a types.Word = types.Word(cpu.R.A)

	if cpu.IsFlagSet(N) == false {
		if cpu.IsFlagSet(H) || a&0x0F > 9 {
			a += 0x06
		}

		if cpu.IsFlagSet(C) || a > 0x9F {
			a += 0x60
		}
	} else {
		if cpu.IsFlagSet(H) {
			a = (a - 6) & 0xFF
		}

		if cpu.IsFlagSet(C) {
			a -= 0x60
		}
	}

	cpu.ResetFlag(H)

	if a&0x100 == 0x100 {
		cpu.SetFlag(C)
	}

	a &= 0xFF

	if a == 0 {
		cpu.SetFlag(Z)
	} else {
		cpu.ResetFlag(Z)
	}

	cpu.R.A = byte(a)
}

//SWAP r
func (cpu *GbcCPU) Swap_r(r *byte) {
	*r = cpu.swapByte(*r)
}

//SWAP (HL)
func (cpu *GbcCPU) Swap_hl() {
	var hlAddr types.Word = types.Word(utils.JoinBytes(cpu.R.H, cpu.R.L))
	var hlValue byte = cpu.ReadByte(hlAddr)
	var result = cpu.swapByte(hlValue)
	cpu.WriteByte(hlAddr, result)
}

//RLCA
func (cpu *GbcCPU) RLCA() {
	var bit7 bool = false

	if cpu.R.A&0x80 == 0x80 {
		bit7 = true
	}

	var calculation byte = cpu.R.A << 1

	if bit7 {
		cpu.SetFlag(C)
		calculation ^= 0x01
	} else {
		cpu.ResetFlag(C)
	}

	//reset flags
	cpu.ResetFlag(Z)
	cpu.ResetFlag(N)
	cpu.ResetFlag(H)

	cpu.R.A = calculation
}

//RLA
func (cpu *GbcCPU) RLA() {
	var bit7 bool = false
	var calculation byte = cpu.R.A

	if calculation&0x80 == 0x80 {
		bit7 = true
	}

	calculation = calculation << 1

	if cpu.IsFlagSet(C) {
		calculation ^= 0x01
	}

	if bit7 {
		cpu.SetFlag(C)
	} else {
		cpu.ResetFlag(C)
	}

	cpu.ResetFlag(Z)
	cpu.ResetFlag(N)
	cpu.ResetFlag(H)

	cpu.R.A = calculation
}

//RLC r
func (cpu *GbcCPU) Rlc_r(r *byte) {
	var bit7 bool = false

	if *r&0x80 == 0x80 {
		bit7 = true
	}

	var calculation byte = *r << 1

	if bit7 {
		cpu.SetFlag(C)
		calculation ^= 0x01
	} else {
		cpu.ResetFlag(C)
	}

	if calculation == 0x00 {
		cpu.SetFlag(Z)
	} else {
		cpu.ResetFlag(Z)

	}

	//reset flags
	cpu.ResetFlag(N)
	cpu.ResetFlag(H)

	*r = calculation
}

//RLC (HL)
func (cpu *GbcCPU) Rlc_hl() {
	var hlAddr types.Word = types.Word(utils.JoinBytes(cpu.R.H, cpu.R.L))
	var hlValue byte = cpu.mmu.ReadByte(hlAddr)

	var bit7 bool = false

	if hlValue&0x80 == 0x80 {
		bit7 = true
	}

	var calculation byte = hlValue << 1

	if bit7 {
		cpu.SetFlag(C)
		calculation ^= 0x01
	} else {
		cpu.ResetFlag(C)
	}

	if calculation == 0x00 {
		cpu.SetFlag(Z)
	} else {
		cpu.ResetFlag(Z)
	}

	//reset flags
	cpu.ResetFlag(N)
	cpu.ResetFlag(H)

	cpu.mmu.WriteByte(hlAddr, calculation)
}

//RL r
func (cpu *GbcCPU) Rl_r(r *byte) {
	var bit7 bool = false
	var calculation byte = *r

	if calculation&0x80 == 0x80 {
		bit7 = true
	}

	calculation = calculation << 1

	if cpu.IsFlagSet(C) {
		calculation ^= 0x01
	}

	if bit7 {
		cpu.SetFlag(C)
	} else {
		cpu.ResetFlag(C)
	}

	if calculation == 0x00 {
		cpu.SetFlag(Z)
	} else {
		cpu.ResetFlag(Z)
	}

	cpu.ResetFlag(N)
	cpu.ResetFlag(H)

	*r = calculation
}

//RL (HL)
func (cpu *GbcCPU) Rl_hl() {
	var HL types.Word = types.Word(utils.JoinBytes(cpu.R.H, cpu.R.L))
	var value byte = cpu.ReadByte(HL)

	var bit7 bool = false
	var calculation byte = value

	if calculation&0x80 == 0x80 {
		bit7 = true
	}

	calculation = calculation << 1

	if cpu.IsFlagSet(C) {
		calculation ^= 0x01
	}

	if bit7 {
		cpu.SetFlag(C)
	} else {
		cpu.ResetFlag(C)
	}

	if calculation == 0x00 {
		cpu.SetFlag(Z)
	} else {
		cpu.ResetFlag(Z)
	}

	cpu.ResetFlag(N)
	cpu.ResetFlag(H)
	cpu.WriteByte(HL, calculation)
}

//RRCA
func (cpu *GbcCPU) RRCA() {
	var bit0 bool = false

	if cpu.R.A&0x01 == 0x01 {
		bit0 = true
	}

	var calculation byte = cpu.R.A >> 1

	if bit0 {
		cpu.SetFlag(C)
		calculation ^= 0x80
	} else {
		cpu.ResetFlag(C)
	}

	//reset flags
	cpu.ResetFlag(Z)
	cpu.ResetFlag(N)
	cpu.ResetFlag(H)

	cpu.R.A = calculation
}

//RRA
func (cpu *GbcCPU) RRA() {
	var bit0 bool = false
	var calculation byte = cpu.R.A

	if calculation&0x01 == 0x01 {
		bit0 = true
	}

	calculation = calculation >> 1

	if cpu.IsFlagSet(C) {
		calculation ^= 0x80
	}

	if bit0 {
		cpu.SetFlag(C)
	} else {
		cpu.ResetFlag(C)
	}

	cpu.ResetFlag(Z)
	cpu.ResetFlag(N)
	cpu.ResetFlag(H)

	cpu.R.A = calculation
}

//RRC r
func (cpu *GbcCPU) Rrc_r(r *byte) {
	var bit0 bool = false

	if *r&0x01 == 0x01 {
		bit0 = true
	}

	var calculation byte = *r >> 1

	if bit0 {
		cpu.SetFlag(C)
		calculation ^= 0x80
	} else {
		cpu.ResetFlag(C)
	}

	//reset flags
	if calculation == 0x00 {
		cpu.SetFlag(Z)
	} else {
		cpu.ResetFlag(Z)
	}

	cpu.ResetFlag(N)
	cpu.ResetFlag(H)

	*r = calculation
}

//RRC (HL)
func (cpu *GbcCPU) Rrc_hl() {
	var hlAddr types.Word = types.Word(utils.JoinBytes(cpu.R.H, cpu.R.L))
	var hlValue byte = cpu.mmu.ReadByte(hlAddr)
	var bit0 bool = false

	if hlValue&0x01 == 0x01 {
		bit0 = true
	}

	var calculation byte = hlValue >> 1

	if bit0 {
		cpu.SetFlag(C)
		calculation ^= 0x80
	} else {
		cpu.ResetFlag(C)
	}

	//reset flags
	if calculation == 0x00 {
		cpu.SetFlag(Z)
	} else {
		cpu.ResetFlag(Z)
	}

	cpu.ResetFlag(N)
	cpu.ResetFlag(H)

	cpu.mmu.WriteByte(hlAddr, calculation)
}

//RR r
func (cpu *GbcCPU) Rr_r(r *byte) {
	var bit0 bool = false
	var calculation byte = *r

	if calculation&0x01 == 0x01 {
		bit0 = true
	}

	calculation = calculation >> 1

	if cpu.IsFlagSet(C) {
		calculation ^= 0x80
	}

	if bit0 {
		cpu.SetFlag(C)
	} else {
		cpu.ResetFlag(C)
	}

	if calculation == 0x00 {
		cpu.SetFlag(Z)
	} else {
		cpu.ResetFlag(Z)
	}

	cpu.ResetFlag(N)
	cpu.ResetFlag(H)

	*r = calculation
}

//RR (HL)
func (cpu *GbcCPU) Rr_hl() {
	var HLAddr types.Word = types.Word(utils.JoinBytes(cpu.R.H, cpu.R.L))
	var value byte = cpu.ReadByte(HLAddr)

	var bit0 bool = false
	var calculation byte = value

	if calculation&0x01 == 0x01 {
		bit0 = true
	}

	calculation = calculation >> 1

	if cpu.IsFlagSet(C) {
		calculation ^= 0x80
	}

	if bit0 {
		cpu.SetFlag(C)
	} else {
		cpu.ResetFlag(C)
	}

	if calculation == 0x00 {
		cpu.SetFlag(Z)
	} else {
		cpu.ResetFlag(Z)
	}

	cpu.ResetFlag(N)
	cpu.ResetFlag(H)
	cpu.WriteByte(HLAddr, calculation)
}

//SLA r
func (cpu *GbcCPU) Sla_r(r *byte) {
	var bit7 bool = false
	var calculation byte = *r

	if calculation&0x80 == 0x80 {
		bit7 = true
	}

	calculation = calculation << 1

	cpu.ResetFlag(N)
	cpu.ResetFlag(H)

	if calculation == 0 {
		cpu.SetFlag(Z)
	} else {
		cpu.ResetFlag(Z)
	}

	if bit7 {
		cpu.SetFlag(C)
	} else {
		cpu.ResetFlag(C)
	}

	*r = calculation
}

//SLA (HL)
func (cpu *GbcCPU) Sla_hl() {
	var HL types.Word = types.Word(utils.JoinBytes(cpu.R.H, cpu.R.L))
	var value byte = cpu.ReadByte(HL)
	var calculation byte = value
	var bit7 bool = false

	if calculation&0x80 == 0x80 {
		bit7 = true
	}

	calculation = calculation << 1

	cpu.ResetFlag(N)
	cpu.ResetFlag(H)

	if calculation == 0 {
		cpu.SetFlag(Z)
	} else {
		cpu.ResetFlag(Z)
	}

	if bit7 {
		cpu.SetFlag(C)
	} else {
		cpu.ResetFlag(C)
	}

	cpu.WriteByte(HL, calculation)
}

//SRA r
func (cpu *GbcCPU) Sra_r(r *byte) {
	var bit0 bool = false
	var calculation byte = *r

	if calculation&0x01 == 0x01 {
		bit0 = true
	}

	calculation = (calculation >> 1) | (calculation & 0x80)

	cpu.ResetFlag(N)
	cpu.ResetFlag(H)

	if calculation == 0 {
		cpu.SetFlag(Z)
	} else {
		cpu.ResetFlag(Z)
	}

	if bit0 {
		cpu.SetFlag(C)
	} else {
		cpu.ResetFlag(C)
	}

	*r = calculation
}

//SRA (HL)
func (cpu *GbcCPU) Sra_hl() {
	var HL types.Word = types.Word(utils.JoinBytes(cpu.R.H, cpu.R.L))
	var value byte = cpu.ReadByte(HL)
	var calculation byte = value
	var bit0 bool = false

	if calculation&0x01 == 0x01 {
		bit0 = true
	}

	calculation = (calculation >> 1) | (calculation & 0x80)

	cpu.ResetFlag(N)
	cpu.ResetFlag(H)

	if calculation == 0 {
		cpu.SetFlag(Z)
	} else {
		cpu.ResetFlag(Z)
	}

	if bit0 {
		cpu.SetFlag(C)
	} else {
		cpu.ResetFlag(C)
	}

	cpu.WriteByte(HL, calculation)
}

//SRL r
func (cpu *GbcCPU) Srl_r(r *byte) {
	var bit0 bool = false
	var calculation byte = *r

	if calculation&0x01 == 0x01 {
		bit0 = true
	}

	calculation = calculation >> 1

	cpu.ResetFlag(N)
	cpu.ResetFlag(H)

	if calculation == 0 {
		cpu.SetFlag(Z)
	} else {
		cpu.ResetFlag(Z)
	}

	if bit0 {
		cpu.SetFlag(C)
	} else {
		cpu.ResetFlag(C)
	}

	*r = calculation
}

//SRL (HL)
func (cpu *GbcCPU) Srl_hl() {
	var HL types.Word = types.Word(utils.JoinBytes(cpu.R.H, cpu.R.L))
	var value byte = cpu.ReadByte(HL)
	var calculation byte = value
	var bit0 bool = false

	if calculation&0x01 == 0x01 {
		bit0 = true
	}

	calculation = calculation >> 1

	cpu.ResetFlag(N)
	cpu.ResetFlag(H)

	if calculation == 0 {
		cpu.SetFlag(Z)
	} else {
		cpu.ResetFlag(Z)
	}

	if bit0 {
		cpu.SetFlag(C)
	} else {
		cpu.ResetFlag(C)
	}

	cpu.WriteByte(HL, calculation)
}

//BIT b, r
func (cpu *GbcCPU) Bitb_r(b byte, r *byte) {
	cpu.bitTest(b, *r)
}

//BIT b,(HL)
func (cpu *GbcCPU) Bitb_hl(b byte) {
	var HL types.Word = types.Word(utils.JoinBytes(cpu.R.H, cpu.R.L))
	var value byte = cpu.ReadByte(HL)
	cpu.bitTest(b, value)
}

// SET b, r
func (cpu *GbcCPU) Setb_r(b byte, r *byte) {
	*r = cpu.setBit(b, *r)
}

// SET b, (HL)
func (cpu *GbcCPU) Setb_hl(b byte) {
	var HL types.Word = types.Word(utils.JoinBytes(cpu.R.H, cpu.R.L))
	var HLValue byte = cpu.ReadByte(HL)
	var result byte = cpu.setBit(b, HLValue)

	cpu.WriteByte(HL, result)
}

// RES b, r
func (cpu *GbcCPU) Resb_r(b byte, r *byte) {
	*r = cpu.resetBit(b, *r)
}

// RES b, (HL)
func (cpu *GbcCPU) Resb_hl(b byte) {
	var HL types.Word = types.Word(utils.JoinBytes(cpu.R.H, cpu.R.L))
	var HLValue byte = cpu.ReadByte(HL)
	var result byte = cpu.resetBit(b, HLValue)

	cpu.WriteByte(HL, result)
}

//JP nn
func (cpu *GbcCPU) JP_nn() {
	var ls byte = cpu.CurrentInstruction.Operands[0]
	var hs byte = cpu.CurrentInstruction.Operands[1]
	cpu.PC = types.Word(utils.JoinBytes(hs, ls))
	cpu.PCJumped = true
}

//JP (HL)
func (cpu *GbcCPU) JP_hl() {
	var HL types.Word = types.Word(utils.JoinBytes(cpu.R.H, cpu.R.L))
	cpu.PC = HL
	cpu.PCJumped = true
}

//JP cc, nn
func (cpu *GbcCPU) JPcc_nn(flag int, jumpWhen bool) {
	var ls byte = cpu.CurrentInstruction.Operands[0]
	var hs byte = cpu.CurrentInstruction.Operands[1]

	if cpu.IsFlagSet(flag) == jumpWhen {
		cpu.PCJumped = true
		cpu.PC = types.Word(utils.JoinBytes(hs, ls))
		cpu.CurrentInstruction.Cycles = 4
	} else {
		cpu.CurrentInstruction.Cycles = 3
	}
}

//JR n
func (cpu *GbcCPU) JR_n() {
	var n byte = cpu.CurrentInstruction.Operands[0]
	if n != 0x00 {
		cpu.PC += types.Word(cpu.CurrentInstruction.OperandsSize + 1)

		if n > 127 {
			cpu.PC -= types.Word(-n)
		} else {
			cpu.PC += types.Word(n)
		}

		cpu.PCJumped = true
	}
}

//JR cc, nn
func (cpu *GbcCPU) JRcc_nn(flag int, jumpWhen bool) {
	var n byte = cpu.CurrentInstruction.Operands[0]

	if cpu.IsFlagSet(flag) == jumpWhen {
		if n != 0x00 {
			cpu.PC += types.Word(cpu.CurrentInstruction.OperandsSize + 1)

			if n > 127 {
				cpu.PC -= types.Word(-n)
			} else {
				cpu.PC += types.Word(n)
			}

			cpu.PCJumped = true
		}
		cpu.CurrentInstruction.Cycles = 3
	} else {
		cpu.CurrentInstruction.Cycles = 2
	}
}

// CALL nn
//Push address of next instruction onto stack and then jump to address nn
func (cpu *GbcCPU) Call_nn() {
	var ls byte = cpu.CurrentInstruction.Operands[0]
	var hs byte = cpu.CurrentInstruction.Operands[1]
	var nextInstr types.Word = cpu.PC + 3
	cpu.pushWordToStack(nextInstr)
	cpu.PC = types.Word(utils.JoinBytes(hs, ls))
	cpu.PCJumped = true
}

// CALL cc,nn
func (cpu *GbcCPU) Callcc_nn(flag int, callWhen bool) {
	var ls byte = cpu.CurrentInstruction.Operands[0]
	var hs byte = cpu.CurrentInstruction.Operands[1]
	var nextInstr types.Word = cpu.PC + 3

	if cpu.IsFlagSet(flag) == callWhen {
		cpu.pushWordToStack(nextInstr)
		cpu.PC = types.Word(utils.JoinBytes(hs, ls))
		cpu.PCJumped = true
		cpu.CurrentInstruction.Cycles = 6
	} else {
		cpu.CurrentInstruction.Cycles = 3
	}
}

// RET
func (cpu *GbcCPU) Ret() {
	cpu.PC = cpu.popWordFromStack()
	cpu.PCJumped = true
}

// RET cc
func (cpu *GbcCPU) Retcc(flag int, returnWhen bool) {
	if cpu.IsFlagSet(flag) == returnWhen {
		cpu.PC = cpu.popWordFromStack()
		cpu.PCJumped = true
		cpu.CurrentInstruction.Cycles = 5
	} else {
		cpu.CurrentInstruction.Cycles = 2
	}
}

// RETI
func (cpu *GbcCPU) Ret_i() {
	cpu.PC = cpu.popWordFromStack()
	cpu.InterruptsEnabled = true
	cpu.PCJumped = true
}

// RST n
func (cpu *GbcCPU) Rst(n byte) {
	cpu.pushWordToStack(cpu.PC + 1)
	cpu.PC = types.Word(n)
	cpu.PCJumped = true
}

//-----------------------------------------------------------------------
//INSTRUCTIONS END
