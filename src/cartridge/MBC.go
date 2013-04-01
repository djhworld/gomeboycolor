/**
 * Created with IntelliJ IDEA.
 * User: danielharper
 * Date: 20/03/2013
 * Time: 22:10
 * To change this template use File | Settings | File Templates.
 */
package cartridge

import "types"

type MemoryBankController interface {
	Write(addr types.Word, value byte)
	Read(addr types.Word) byte
	switchROMBank(bank int)
	switchRAMBank(bank int)
}
