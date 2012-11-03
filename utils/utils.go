package utils

func JoinBytes(hob, lob byte) uint16 {
	return (uint16(hob) << 8) ^ uint16(lob)
}
