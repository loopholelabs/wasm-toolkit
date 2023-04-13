package main

import (
	"fmt"
	"io/ioutil"
)

func main() {
	files := []string{"test", "embedtest"}

	for _, f := range files {
		fmt.Printf("Open '%s'\n", f)

		data, err := ioutil.ReadFile(f)
		if err == nil {
			fmt.Printf("--------\n%s\n--------\n\n\n", data)
		} else {
			fmt.Printf("No such file\n")
		}
	}
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
func exampleFunction(x_value int32, y_value int32) int32 {
	var zoobs int32 = 45
	if x_value == 0 {
		return -1
	}
	if y_value == 0 {
		return -2
	}
	if x_value > 44 {
		zoobs = 1
	}

	return x_value * y_value * zoobs
}
