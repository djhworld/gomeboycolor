package mmu

import "types"

type Peripheral interface {
	Name() string
	Read(Address types.Word) byte
	Write(Address types.Word, Value byte)
}
