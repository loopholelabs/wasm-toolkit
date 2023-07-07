package main

import (
	"fmt"
	"os"

	"github.com/loopholelabs/wasm-toolkit/pkg/wasm/debug"
	"github.com/loopholelabs/wasm-toolkit/pkg/wasm/wasmfile"
)

func main() {

	args := os.Args[1:]

	wasmfilename := args[0] // For now...

	fmt.Printf("Loading up wasm file %s\n", wasmfilename)
	wfile, err := wasmfile.New(wasmfilename)
	if err != nil {
		panic(err)
	}

	wfile.Debug = &debug.WasmDebug{}

	fmt.Printf("Parsing custom name section...\n")
	wfile.Debug.ParseNameSectionData(wfile.GetCustomSectionData("name"))

	fmt.Printf("Parsing custom dwarf debug sections...\n")
	err = wfile.Debug.ParseDwarf(wfile)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Parsing dwarf line numbers...\n")
	err = wfile.Debug.ParseDwarfLineNumbers()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Parsing dwarf local variables...\n")
	err = wfile.Debug.ParseDwarfVariables(wfile)
	if err != nil {
		panic(err)
	}

	outputWasm := "output.wasm"
	outputWat := "output.wat"

	fmt.Printf("Writing wasm out to %s...\n", outputWasm)
	f, err := os.Create(outputWasm)
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

	fmt.Printf("Writing wat out to %s...\n", outputWat)
	watf, err := os.Create(outputWat)
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
