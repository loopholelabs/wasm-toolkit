package main

import (
	"fmt"
)

func main() {
	fmt.Printf("Module says main()\n")
}

//export hello
//go:linkname hello
func hello() {
	v := make([]byte, 1024)
	fmt.Printf("Module says hello() I have a byte array len %d\n", len(v))

	s := exampleFunction(55, 66)
	fmt.Printf("returned: %d\n", s)
}

//go:noinline
func exampleFunction(x int32, y int32) int32 {
	var zoobs int32 = 45
	if x == 0 {
		return -1
	}
	if y == 0 {
		return -2
	}
	if x > 44 {
		zoobs = 1
	}

	return x * y * zoobs
}
