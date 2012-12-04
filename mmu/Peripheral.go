package mmu

import "github.com/djhworld/gomeboycolor/types"

type Peripheral interface {
	Read(Address types.Word) byte
	Write(Address types.Word, Value byte)
}
