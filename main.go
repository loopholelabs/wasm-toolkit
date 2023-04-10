package main

import (
	"bytes"
	"fmt"
	"os"

	"github.com/loopholelabs/wasm-lnkr/wasmfile"
)

func main() {

	args := os.Args[1:]

	wasmfilename := args[0] // For now...

	fmt.Printf("Loading up wasm file %s\n", wasmfilename)
	wfile, err := wasmfile.New(wasmfilename)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Parsing custom name section...\n")
	err = wfile.ParseName()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Parsing custom dwarf debug sections...\n")
	err = wfile.ParseDwarf()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Parsing dwarf line numbers...\n")
	err = wfile.ParseDwarfLineNumbers()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Parsing dwarf local variables...\n")
	err = wfile.ParseDwarfVariables()
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

	// Do a second test from the wat file...

	outputWasm2 := "output2.wasm"
	outputWat2 := "output2.wat"

	fmt.Printf("Checking up wat file %s\n", outputWat)
	wf2, err := wasmfile.NewFromWat(outputWat)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Writing wasm out to %s...\n", outputWasm2)
	f2, err := os.Create(outputWasm2)
	if err != nil {
		panic(err)
	}

	err = wf2.EncodeBinary(f2)
	if err != nil {
		panic(err)
	}

	err = f2.Close()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Writing wat out to %s...\n", outputWat2)
	watf2, err := os.Create(outputWat2)
	if err != nil {
		panic(err)
	}

	err = wf2.EncodeWat(watf2)
	if err != nil {
		panic(err)
	}

	err = watf2.Close()
	if err != nil {
		panic(err)
	}

	// TODO: Verify wfile2
	for idx, f1 := range wfile.Table {
		f2 := wf2.Table[idx]

		var buf1 bytes.Buffer
		f1.EncodeBinary(&buf1)
		var buf2 bytes.Buffer
		f2.EncodeBinary(&buf2)

		data1 := buf1.Bytes()
		data2 := buf2.Bytes()

		var issue = false
		for i, b1 := range data1 {
			b2 := data2[i]
			if b1 != b2 {
				issue = true
				break
			}
		}

		if issue {
			fmt.Printf("Differs %d %x %x\n", idx, data1, data2)
		}

	}

	// Compare the code...
	for idx, c1 := range wfile.Code {
		c2 := wf2.Code[idx]

		var buf1 bytes.Buffer
		c1.EncodeBinary(&buf1)
		var buf2 bytes.Buffer
		c2.EncodeBinary(&buf2)

		data1 := buf1.Bytes()
		data2 := buf2.Bytes()

		var issue = false
		for i, b1 := range data1 {
			b2 := data2[i]
			if b1 != b2 {
				issue = true
				break
			}
		}

		if issue {
			fmt.Printf("Differs %d %x %x\n", idx, data1, data2)
		}
	}

}
