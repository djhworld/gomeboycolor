/**
 * Created with IntelliJ IDEA.
 * User: danielharper
 * Date: 02/03/2013
 * Time: 10:32
 * To change this template use File | Settings | File Templates.
 */
package constants

import "types"

//interrupt handler addresses
const (
	V_BLANK_IR_ADDR byte = 0x40
	LCD_IR_ADDR          = 0x48
)

const (
	V_BLANK_IRQ byte = 0x01
	LCD_IRQ          = 0x02
)

//memory addresses
const (
	INTERRUPT_ENABLED_FLAG_ADDR types.Word = 0xFFFF
	INTERRUPT_FLAG_ADDR                    = 0xFF0F
)
