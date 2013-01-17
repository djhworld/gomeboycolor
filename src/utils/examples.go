package utils

import "fmt"

func ExampleJoinBytes() {
	var hob byte = 0xF3
	var lob byte = 0xA9

	fmt.Printf("%X", JoinBytes(hob, lob))

	//Output: F3A9
}
