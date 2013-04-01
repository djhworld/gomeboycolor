/**
 * Created with IntelliJ IDEA.
 * User: danielharper
 * Date: 15/03/2013
 * Time: 20:34
 * To change this template use File | Settings | File Templates.
 */
package timer

import (
	"components"
	"log"
	"types"
	"utils"
)

type Timer struct {
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
	timer.a++
	return byte(timer.a)
}

func (timer *Timer) Write(Address types.Word, Value byte) {
	log.Printf("Writing %s to timer register %s", Address, utils.ByteToString(Value))

}

func (timer *Timer) LinkIRQHandler(m components.IRQHandler) {

}

func (timer *Timer) Reset() {
	log.Println("Resetting", timer.Name())
}
