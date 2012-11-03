package utils

//Joins two bytes together to form a 16 bit integer 
func JoinBytes(hob, lob byte) uint16 {
	return (uint16(hob) << 8) ^ uint16(lob)
}
