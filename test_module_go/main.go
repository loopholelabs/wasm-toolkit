package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"unsafe"

	_ "unsafe"
)

var some_global int32 = 0x1234
var another_global int32 = 1

func main() {
	fmt.Printf("TEST environment var is %v\n", os.Getenv("TEST"))
	os.WriteFile("newfile", []byte("Testing"), 660)

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

	os.Rename("newfile", "newfile2")
	os.Remove("newfile2")

	exampleFunction(12, 46)

	fmt.Printf("some_global is %d, another_global is %d\n", some_global, another_global)

	var tt = make([]byte, 10)
	ptr := &tt[0]
	watch(4, uint32(uintptr(unsafe.Pointer(ptr))), 10)

	tt[0] = 0x12
	tt[1] = 0x34

	fmt.Printf("Data %x\n", tt)

	// panic("Something happened")
}

//export hello
//go:linkname hello
func hello() {
	some_global = 1
	v := make([]byte, 1024)
	fmt.Printf("Module says hello() I have a byte array len %d\n", len(v))

	s := exampleFunction(55, 66)
	fmt.Printf("returned: %d\n", s)
}

//go:noinline
func exampleFunction(x_value int32, y_value int32) int32 {
	var jm_boob int32 = 1

	watch(1, uint32(uintptr(unsafe.Pointer(&jm_boob))), 4)
	watch(2, uint32(uintptr(unsafe.Pointer(&some_global))), 4)
	watch(3, uint32(uintptr(unsafe.Pointer(&another_global))), 4)

	some_global = 2
	another_global = 0x999
	var zoobs int32 = 45
	if x_value == 0 {
		jm_boob++
	}
	jm_boob++
	if y_value == 0 {
		return -2 + jm_boob
	}

	unwatch(1)

	jm_boob++

	watch(1, uint32(uintptr(unsafe.Pointer(&jm_boob))), 4)

	if x_value > 44 {
		zoobs = 1
		jm_boob++
	}

	anotherFunction(&jm_boob)

	another_global = 0x2222

	//unwatch(0)

	return (x_value * y_value * zoobs) + jm_boob
}

func anotherFunction(d *int32) {
	*d = 0x1234
}

//go:wasm-module scale
//export watch
func watch(id uint32, addr uint32, len uint32)

//go:wasm-module scale
//export unwatch
func unwatch(id uint32)
