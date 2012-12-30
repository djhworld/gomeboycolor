package mmu

import "types"

type Peripheral interface {
	Read(Address types.Word) byte
	Write(Address types.Word, Value byte)
}
