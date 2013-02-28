package utils

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func ByteToString(b byte) string {
	var zeroes string
	if b < 0x10 {
		zeroes += "0"
	}
	return fmt.Sprintf("0x%s%X", zeroes, b)
}

func StringToByte(s string) (byte, error) {
	s = strings.Replace(s, "0x", "", 1)
	if len(s) > 2 {
		return 0x0, errors.New("Please enter an address between 00 and FF")
	}

	result, err := strconv.ParseInt(s, 16, 64)

	return byte(result), err
}

func StringToWord(s string) (uint16, error) {
	s = strings.Replace(s, "0x", "", 1)

	if len(s) > 4 {
		return 0x0, errors.New("Please enter an address between 0000 and FFFF")
	}

	result, err := strconv.ParseInt(s, 16, 64)

	return uint16(result), err
}

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

func CompareBytes(a, b byte, operator string) bool {
	switch operator {
	case "==":
		return a == b
	case ">":
		return a > b
	case "<":
		return a < b
	case ">=":
		return a >= b
	case "<=":
		return a <= b
	}
	return false
}

func CompareWords(a, b uint16, operator string) bool {
	switch operator {
	case "==":
		return a == b
	case ">":
		return a > b
	case "<":
		return a < b
	case ">=":
		return a >= b
	case "<=":
		return a <= b
	}
	return false
}

func PadRight(s string, maxLen int, padStr string) string {
	if l := len(s); l < maxLen {
		padAmount := maxLen - l
		s = s + strings.Repeat(padStr, padAmount)
	}
	return s
}
