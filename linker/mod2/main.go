package main

import "fmt"

func main() {
	fmt.Printf("Module [2] says main()")
}

//export hello
//go:linkname hello
func hello() {
	fmt.Printf("Module [2] says hello().")
}
