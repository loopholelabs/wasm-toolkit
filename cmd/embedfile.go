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

	"github.com/loopholelabs/wasm-toolkit/internal/wat"
	wasmfile "github.com/loopholelabs/wasm-toolkit/pkg/wasm"
	"github.com/loopholelabs/wasm-toolkit/pkg/wasm/debug"
	"github.com/loopholelabs/wasm-toolkit/pkg/wasm/types"

	"github.com/spf13/cobra"
)

var (
	cmdEmbedfile = &cobra.Command{
		Use:   "embedfile",
		Short: "Add a file to the wasm",
		Long:  `This will embed a file within the wasm`,
		Run:   runEmbedFile,
	}
)

var em_filename = "embedtest"
var em_content = "Yeah!"
var em_contentfile = ""

func init() {
	rootCmd.AddCommand(cmdEmbedfile)
	cmdEmbedfile.Flags().StringVar(&em_filename, "filename", "embedtest", "Embed filename")
	cmdEmbedfile.Flags().StringVar(&em_content, "content", "Hey! This isn't really a file. It's embedded in the wasm.", "Embed content")
	cmdEmbedfile.Flags().StringVar(&em_contentfile, "contentfile", "", "Embed content from file")
}

func runEmbedFile(ccmd *cobra.Command, args []string) {
	if Input == "" {
		panic("No input file")
	}

	fmt.Printf("Loading wasm file \"%s\"...\n", Input)
	wfile, err := wasmfile.New(Input)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Parsing custom name section...\n")
	wfile.Debug = &debug.WasmDebug{}
	wfile.Debug.ParseNameSectionData(wfile.GetCustomSectionData("name"))

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

	// TODO: Wrap file imports so we can do what we want to...

	originalFunctionLength := len(wfile.Code)

	wfile.AddFuncsFrom(memFunctions, func(m map[int]int) {})

	data_ptr := wfile.Memory[0].LimitMin << 16
	wfile.SetGlobal("$debug_start_mem", types.ValI32, fmt.Sprintf("i32.const %d", data_ptr))

	// Now we can start doing interesting things...

	em_content_data := []byte(em_content)

	if em_contentfile != "" {
		bytes, err := os.ReadFile(em_contentfile)
		if err != nil {
			panic(err)
		}
		em_content_data = bytes
	}

	// Add a payload to the wasm file
	embedFunctions := &wasmfile.WasmFile{}
	data, err = wat.Wat_content.ReadFile(path.Join("wat_code", "embed.wat"))
	if err != nil {
		panic(err)
	}
	err = embedFunctions.DecodeWat(data)
	if err != nil {
		panic(err)
	}

	wfile.AddDataFrom(int32(data_ptr), embedFunctions)

	wfile.AddData("$file_name", []byte(em_filename))
	wfile.AddData("$file_content", em_content_data)

	// Find out how much data we need for the payload
	total_payload_data := data_ptr
	if len(wfile.Data) > 0 {
		last_data := wfile.Data[len(wfile.Data)-1]
		total_payload_data = int(last_data.Offset[0].I32Value) + len(last_data.Data) - data_ptr
	}

	payload_size := (total_payload_data + 65535) >> 16
	fmt.Printf("Payload data of %d (%d pages)\n", total_payload_data, payload_size)

	wfile.SetGlobal("$debug_mem_size", types.ValI32, fmt.Sprintf("i32.const %d", payload_size)) // The size of our addition in 64k pages
	wfile.Memory[0].LimitMin = wfile.Memory[0].LimitMin + payload_size

	wfile.AddFuncsFrom(embedFunctions, func(m map[int]int) {}) // NB: This may mean inserting an import which changes all func numbers.

	// Redirect some imports...
	import_redirect_map := map[string]string{
		"wasi_snapshot_preview1:fd_prestat_get": "$wrap_fd_prestat_get",
		"wasi_snapshot_preview1:path_open":      "$wrap_path_open",
		"wasi_snapshot_preview1:fd_read":        "$wrap_fd_read",
	}

	for from, to := range import_redirect_map {
		fromId := wfile.LookupImport(from)
		toId := wfile.LookupFunctionID(to)

		fmt.Printf("Redirecting code from %d to %d\n", fromId, toId)

		for idx, c := range wfile.Code {
			if idx < originalFunctionLength {
				c.ModifyAllCalls(map[int]int{fromId: toId})
			}
		}

	}

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
