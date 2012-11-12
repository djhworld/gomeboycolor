package utils

//Joins two bytes together to form a 16 bit integer 
func JoinBytes(hob, lob byte) uint16 {
	return (uint16(hob) << 8) ^ uint16(lob)
}

//Splits one 16 bit integer to two bytes
func SplitIntoBytes(bb uint16) (byte, byte) {
	return byte(bb >> 8), byte(bb & 0x00FF)
}

//swaps the nibbles of a byte around
func SwapNibbles(a byte) byte {
	return (a&0xF0)>>4 ^ ((a & 0x0F) << 4)
}
