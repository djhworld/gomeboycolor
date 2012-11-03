package main

type MMU interface {
	WriteByte(address Word, value byte)
	WriteWord(address Word, value Word)

	ReadByte(address Word) byte
	ReadWord(address Word) Word
}


