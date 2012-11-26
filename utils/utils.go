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

//If you pass a number in between 0-7 it returns the 
//value in relation to the position of the bit.
func BitToValue(b byte) byte {
	switch b {
	case 0x00:
		return 0x01
	case 0x01:
		return 0x02
	case 0x02:
		return 0x04
	case 0x03:
		return 0x08
	case 0x04:
		return 0x10
	case 0x05:
		return 0x20
	case 0x06:
		return 0x40
	case 0x07:
		return 0x80
	}
	return 0x01
}

func CalculateCycles(timings [16][16]int, instruction byte) int {
	i, j := int(instruction&0xF0)>>4, int(instruction&0x0F)
	return timings[i][j]
}

