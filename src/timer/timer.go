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
	"fmt"
	"log"
	"time"
	"types"
	"constants"
)

const (
	DIV_REGISTER  types.Word = 0xFF04
	TIMA_REGISTER            = 0xFF05
	TMA_REGISTER             = 0xFF06
	TAC_REGISTER             = 0xFF07
)

const (
	NAME = "TIMER"
)

var ClockFrequencies map[int]time.Duration = map[int]time.Duration{
	4096:   time.Duration(244141 * time.Nanosecond),
	16384:  time.Duration(61035 * time.Nanosecond),
	65536:  time.Duration(15259 * time.Nanosecond),
	262144: time.Duration(3815 * time.Nanosecond),
}

type TimerCounter struct {
	Name       string
	Value      byte
	ticker     *time.Ticker
	duration   time.Duration
	irqHandler components.IRQHandler
	timerOn    bool
}

func NewTimerCounter(name string, d time.Duration, irqHandler components.IRQHandler) *TimerCounter {
	c := new(TimerCounter)
	c.Value = 0x00
	c.Name = name
	c.duration = d
	c.ticker = time.NewTicker(c.duration)
	c.irqHandler = irqHandler
	c.timerOn = true
	log.Println(NAME+":", "Started a new timer", name, "with a duration of", d)
	return c
}

func (c *TimerCounter) Run() {
	for {
		select {
		case <-c.ticker.C:
			if c.timerOn {

				c.Value++
				if c.Value == 0x00 && c.irqHandler != nil {
					c.irqHandler.RequestInterrupt(constants.TIMER_OVERFLOW_IRQ)
				}
			}
		}
	}
}

func (c *TimerCounter) Reset() {
	c.Value = 0x00
}

func (c *TimerCounter) Stop() {
	c.ticker.Stop()
}

type Timer struct {
	timaCounter  *TimerCounter
	divCounter   *TimerCounter
	tacFrequency *time.Duration
	tacRegister  byte
	tmaRegister  byte
	irqHandler   components.IRQHandler
}

func NewTimer() *Timer {
	var t *Timer = new(Timer)
	t.divCounter = NewTimerCounter("DIV", ClockFrequencies[16384], nil)
	go t.divCounter.Run()
	t.ResetTIMACounter()
	return t
}

func (timer *Timer) Name() string {
	return NAME
}

func (timer *Timer) Read(Address types.Word) byte {
	switch Address {
	case DIV_REGISTER:
		return timer.divCounter.Value
	case TIMA_REGISTER:
		return timer.timaCounter.Value
	case TMA_REGISTER:
		return timer.tmaRegister
	case TAC_REGISTER:
		return timer.tacRegister
	default:
		panic(fmt.Sprintln("Timer module is not set up to handle address", Address))
	}
	return 0x00
}

func (timer *Timer) Write(address types.Word, value byte) {
	switch address {
	case DIV_REGISTER:
		timer.divCounter.Reset()
	case TIMA_REGISTER:
		timer.timaCounter.Value = value
	case TMA_REGISTER:
		timer.tmaRegister = value
	case TAC_REGISTER:
		timer.tacRegister = value
		timer.ResetTIMACounter()
	default:
		panic(fmt.Sprintln("Timer module is not set up to handle address", address))
	}
}

func (timer *Timer) ResetTIMACounter() {
	var frequency byte = timer.tacRegister & 0x03

	if timer.timaCounter != nil {
		timer.timaCounter.Stop()
	}

	switch frequency {
	case 0x00:
		timer.timaCounter = NewTimerCounter("TIMA", ClockFrequencies[4096], timer.irqHandler)
	case 0x01:
		timer.timaCounter = NewTimerCounter("TIMA", ClockFrequencies[262144], timer.irqHandler)
	case 0x10:
		timer.timaCounter = NewTimerCounter("TIMA", ClockFrequencies[65536], timer.irqHandler)
	case 0x11:
		timer.timaCounter = NewTimerCounter("TIMA", ClockFrequencies[16384], timer.irqHandler)
	}

	timer.timaCounter.timerOn = timer.tacRegister&0x04 == 0x04

	go timer.timaCounter.Run()
}

func (timer *Timer) LinkIRQHandler(m components.IRQHandler) {
	timer.irqHandler = m
}

func (timer *Timer) Reset() {
	log.Println("Resetting", timer.Name())
}
