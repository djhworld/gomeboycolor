/**
 * Created with IntelliJ IDEA.
 * User: danielharper
 * Date: 15/03/2013
 * Time: 20:34
 * To change this template use File | Settings | File Templates.
 */
package timer

import (
	"types"
	"fmt"
	"log"
	"components"
	"utils"
)

const (
	DIV_REGISTER  types.Word = 0xFF04
	TIMA_REGISTER            = 0xFF05
	TMA_REGISTER             = 0xFF06
	TAC_REGISTER             = 0xFF07
)

/*
var ClockFrequencies map[int]int = map[int]int{
	0: 4194304/2^10,
	1: 4194304/2^4/1024,
	2: 4194304/2^6/1024,
	3: 4194304/2^8/1024,
}
*/

type Timer struct {
	timerOn              bool
	dividerRegister      byte //DIV 0xFF04
	timerCounterRegister byte //TIMA 0xFF05
	timerModuloRegister  byte //TMA 0xFF06
	timerControlRegister byte //TAC 0xFF07
	inputClockFrequency  int
	a int
}

func NewTimer() *Timer {
	var t *Timer = new(Timer)
	t.Reset()
	return t
}

func (timer *Timer) Name() string {
	return "Timer"
}

func (timer *Timer) Read(Address types.Word) byte {
	log.Printf("Reading from timer register %s", Address)
	switch Address {
	case DIV_REGISTER:
		timer.dividerRegister++ //TODO: this is obviously wrong, timers need to work properly
		return timer.dividerRegister
	case TIMA_REGISTER:
		return timer.timerCounterRegister
	case TMA_REGISTER:
		return timer.timerModuloRegister
	case TAC_REGISTER:
		return timer.timerControlRegister
	default:
		panic(fmt.Sprintln("Timer module is not set up to handle address", Address))
	}
	return 0x00
}

func (timer *Timer) Write(Address types.Word, Value byte) {
	log.Printf("Writing %s to timer register %s", utils.ByteToString(Value), Address)

	switch Address {
	case DIV_REGISTER:
		//writing any value sets this to 0
		timer.dividerRegister = 0x00
	case TIMA_REGISTER:
		timer.timerCounterRegister = Value
	case TMA_REGISTER:
		timer.timerModuloRegister = Value
	case TAC_REGISTER:
		timer.timerControlRegister = Value
		timer.timerOn = (timer.timerControlRegister & 0x04 == 0x04)
	default:
		panic(fmt.Sprintln("Timer module is not set up to handle address", Address))
	}
}

func (timer *Timer) LinkIRQHandler(m components.IRQHandler) {

}

func (timer *Timer) Reset() {
	log.Println("Resetting", timer.Name())
}
