/**
 * Just a place holder for now
 * User: danielharper
 * Date: 28/02/2013
 * Time: 21:26
 */
package apu

import (
	"components"
	"log"
	"types"
)

type APU struct{}

func NewAPU() *APU {
	var a *APU = new(APU)
	a.Reset()
	return a
}

func (apu *APU) Name() string {
	return "APU"
}

func (apu *APU) Read(Address types.Word) byte {
	return 0x00
}

func (apu *APU) Write(Address types.Word, Value byte) {

}

func (apu *APU) LinkIRQHandler(m components.IRQHandler) {

}

func (apu *APU) Reset() {
	log.Println("Resetting", apu.Name())
}
