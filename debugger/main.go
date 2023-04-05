package main

import (
	"os"

	"github.com/loopholelabs/wasm-lnkr/wasmfile"
)

func main() {

	args := os.Args[1:]

	wasmfilename := args[0] // For now...

	wfile, err := wasmfile.New(wasmfilename)

	err = wfile.ParseName()
	if err != nil {
		panic(err)
	}

	err = wfile.ParseDwarf()
	if err != nil {
		panic(err)
	}

	err = wfile.ParseDwarfLineNumbers()
	if err != nil {
		panic(err)
	}

	err = wfile.ParseDwarfVariables()
	if err != nil {
		panic(err)
	}

	f, err := os.Create("test.wasm")
	if err != nil {
		panic(err)
	}

	err = wfile.EncodeBinary(f)
	if err != nil {
		panic(err)
	}

	err = f.Close()
	if err != nil {
		panic(err)
	}

	watf, err := os.Create("test.wat")
	if err != nil {
		panic(err)
	}

	err = wfile.EncodeWat(watf)
	if err != nil {
		panic(err)
	}

	err = watf.Close()
	if err != nil {
		panic(err)
	}

}
