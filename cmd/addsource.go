/*
	Copyright 2023 Loophole Labs

	Licensed under the Apache License, Version 2.0 (the "License");
	you may not use this file except in compliance with the License.
	You may obtain a copy of the License at

		   http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.
*/

package main

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/loopholelabs/wasm-toolkit/internal/wat"
	"github.com/loopholelabs/wasm-toolkit/wasmfile"

	"github.com/spf13/cobra"
)

var (
	cmdAddSource = &cobra.Command{
		Use:   "addsource",
		Short: "Add some source code to an interpreter wasm",
		Long:  `This will embed some source code into the wasm`,
		Run:   runAddSource,
	}
)

var source_file string

func init() {
	rootCmd.AddCommand(cmdAddSource)
	cmdAddSource.Flags().StringVar(&source_file, "filename", "", "Source filename")
}

func runAddSource(ccmd *cobra.Command, args []string) {
	if Input == "" {
		panic("No input file")
	}

	fmt.Printf("Loading wasm file \"%s\"...\n", Input)
	wfile, err := wasmfile.New(Input)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Parsing custom name section...\n")
	err = wfile.ParseName()
	if err != nil {
		panic(err)
	}

	originalFunctionLength := len(wfile.Code)

	// Add a payload to the wasm file
	memFunctions := &wasmfile.WasmFile{}
	data, err := wat.Wat_content.ReadFile(path.Join("wat_code", "memory.wat"))
	if err != nil {
		panic(err)
	}
	err = memFunctions.DecodeWat(data)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Adding functions from memory.wat...\n")
	wfile.AddFuncsFrom(memFunctions, func(m map[int]int) {})

	data_ptr := wfile.Memory[0].LimitMin << 16
	wfile.SetGlobal("$debug_start_mem", wasmfile.ValI32, fmt.Sprintf("i32.const %d", data_ptr))

	// Now we can start doing what we want...

	bytes, err := os.ReadFile(source_file)
	if err != nil {
		panic(err)
	}

	// Now we just need to adjust the imported functions get_source_len and get_source_ptr and then remove them.

	// Add a payload to the wasm file
	replacedFunctions := &wasmfile.WasmFile{}
	data, err = wat.Wat_content.ReadFile(path.Join("wat_code", "addsource.wat"))
	if err != nil {
		panic(err)
	}
	err = replacedFunctions.DecodeWat(data)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Adding functions from addsource.wat...\n")
	wfile.AddFuncsFrom(replacedFunctions, func(m map[int]int) {})

	wfile.AddDataFrom(int32(data_ptr), replacedFunctions)

	fmt.Printf("Adding source data %d bytes...\n", len(bytes))
	wfile.AddData("$source_data", bytes)

	// Now we need to remap any calls to the new functions

	fid_get_source_len := wfile.LookupFunctionID("$get_source_len")
	fid_get_source_ptr := wfile.LookupFunctionID("$get_source")

	remap := map[int]int{}
	remap_imports := map[int]int{}

	// Now we need to REMOVE the old imports.
	newImports := make([]*wasmfile.ImportEntry, 0)
	for n, i := range wfile.Import {
		if i.Module == "env" && i.Name == "get_source_len" {
			remap_imports[n] = fid_get_source_len
		} else if i.Module == "env" && i.Name == "get_source" {
			remap_imports[n] = fid_get_source_ptr
		} else {
			remap[n] = len(newImports)
			// Keep them for now...
			newImports = append(newImports, i)
		}

	}

	// Remap everything in the Code section because we're removing 2 imports.
	for n, _ := range wfile.Code {
		remap[len(wfile.Import)+n] = len(newImports) + n
	}

	// Remap the imports, and THEN remap due to removing the imports.
	for idx, c := range wfile.Code {
		if idx < originalFunctionLength {
			c.ModifyAllCalls(remap_imports)
			c.ModifyAllCalls(remap)
		}
	}

	wfile.Import = newImports

	// We also need to fixup any Elems sections
	for _, el := range wfile.Elem {
		for idx, funcidx := range el.Indexes {
			newidx, ok := remap[int(funcidx)]
			if ok {
				el.Indexes[idx] = uint64(newidx)
			}
		}
	}

	// Fixup exports
	for _, ex := range wfile.Export {
		if ex.Type == wasmfile.ExportFunc {
			newidx, ok := remap[ex.Index]
			if ok {
				ex.Index = newidx
			}
		}
	}

	// Now we need to remap any calls to the new functions

	// Does this work ok?
	wfile.Renumber_functions(remap)

	// Find out how much data we need for the payload
	total_payload_data := data_ptr
	if len(wfile.Data) > 0 {
		last_data := wfile.Data[len(wfile.Data)-1]
		total_payload_data = int(last_data.Offset[0].I32Value) + len(last_data.Data) - data_ptr
	}

	payload_size := (total_payload_data + 65535) >> 16

	wfile.SetGlobal("$debug_mem_size", wasmfile.ValI32, fmt.Sprintf("i32.const %d", payload_size)) // The size of our addition in 64k pages
	wfile.Memory[0].LimitMin = wfile.Memory[0].LimitMin + payload_size

	// Pass on the fact of if source_file is gzip or not.
	source_gzipped := 0
	if strings.HasSuffix(source_file, ".gz") {
		source_gzipped = 1
	}
	wfile.SetGlobal("$source_gzipped", wasmfile.ValI32, fmt.Sprintf("i32.const %d", source_gzipped))

	// Adjust any memory.size / memory.grow calls
	for idx, c := range wfile.Code {
		if idx < originalFunctionLength {
			err = c.ReplaceInstr(wfile, "memory.grow", "call $debug_memory_grow")
			if err != nil {
				panic(err)
			}
			err = c.ReplaceInstr(wfile, "memory.size", "call $debug_memory_size")
			if err != nil {
				panic(err)
			}
		} else {
			// Do any relocation adjustments...
			err = c.InsertAfterRelocating(wfile, `global.get $debug_start_mem
																						i32.add`)
			if err != nil {
				panic(err)
			}
		}

		err = c.ResolveLengths(wfile)
		if err != nil {
			panic(err)
		}

		err = c.ResolveRelocations(wfile, data_ptr)
		if err != nil {
			panic(err)
		}

		err = c.ResolveGlobals(wfile)
		if err != nil {
			panic(err)
		}

		err = c.ResolveFunctions(wfile)
		if err != nil {
			panic(err)
		}
	}

	fmt.Printf("Writing wasm out to %s...\n", Output)
	f, err := os.Create(Output)
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
	/*
	   fmt.Printf("Writing debug.wat\n")
	   f2, err := os.Create("debug.wat")

	   	if err != nil {
	   		panic(err)
	   	}

	   err = wfile.EncodeWat(f2)

	   	if err != nil {
	   		panic(err)
	   	}

	   err = f2.Close()

	   	if err != nil {
	   		panic(err)
	   	}
	*/
}
