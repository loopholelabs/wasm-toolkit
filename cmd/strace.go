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

	"github.com/loopholelabs/wasm-toolkit/wasmfile"

	"github.com/spf13/cobra"
)

var (
	cmdStrace = &cobra.Command{
		Use:   "strace",
		Short: "Use strace to add tracing output to as wasm file",
		Long:  `This will output debug info to STDERR`,
		Run:   runStrace,
	}
)

func init() {
	rootCmd.AddCommand(cmdStrace)
	cmdStrace.Flags().StringVarP(&Input, "input", "i", "", "Input file name")
	cmdStrace.Flags().StringVarP(&Output, "output", "o", "output.wasm", "Output file name")
}

func runStrace(ccmd *cobra.Command, args []string) {
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
	/*
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
	*/
	// Add a payload to the wasm file
	memFunctions, err := wasmfile.NewFromWat(path.Join("wat_code", "memory.wat"))
	if err != nil {
		panic(err)
	}

	originalFunctionLength := len(wfile.Code)

	wfile.AddFuncsFrom(memFunctions)

	payload_size := 1

	data_ptr := wfile.Memory[0].LimitMin << 16

	wfile.SetGlobal("$debug_mem_size", wasmfile.ValI32, fmt.Sprintf("i32.const %d", payload_size)) // The size of our addition in 64k pages
	wfile.SetGlobal("$debug_start_mem", wasmfile.ValI32, fmt.Sprintf("i32.const %d", data_ptr))
	wfile.Memory[0].LimitMin = wfile.Memory[0].LimitMin + payload_size

	// Now we can start doing interesting things...

	// Add a payload to the wasm file
	debugFunctions, err := wasmfile.NewFromWat(path.Join("wat_code", "strace.wat"))
	if err != nil {
		panic(err)
	}

	wfile.AddDataFrom(int32(data_ptr), debugFunctions)
	wfile.AddFuncsFrom(debugFunctions)

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

			functionIndex := idx + len(wfile.Import)

			blockInstr := "block"
			f := wfile.Function[idx]
			t := wfile.Type[f.TypeIndex]
			if len(t.Result) > 0 {
				blockInstr = fmt.Sprintf("block (result %s)", wasmfile.ByteToValType[t.Result[0]])
			}

			// For now
			blockInstr = ""

			// Create some useful data...
			wfile.AddData(fmt.Sprintf("$function_name_%d", functionIndex), []byte(wfile.GetFunctionIdentifier(functionIndex, false)))

			startCode := fmt.Sprintf(`%s
			i32.const %d
			i32.const offset($function_name_%d)
			i32.const length($function_name_%d)
			call $debug_enter_func
			`, blockInstr, functionIndex, functionIndex, functionIndex)

			// Do parameters...
			for paramIndex, pt := range t.Param {
				startCode = fmt.Sprintf(`%s
					i32.const %d
					i32.const %d
					local.get %d
					call $debug_enter_%s
					`, startCode, functionIndex, paramIndex, paramIndex, wasmfile.ByteToValType[pt])
			}

			startCode = fmt.Sprintf(`%s
					i32.const %d
					call $debug_enter_end
					`, startCode, functionIndex)

			err = c.InsertFuncStart(wfile, startCode)
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
	}

	fmt.Printf("Writing wat out to %s...\n", Output)
	f, err := os.Create(Output)
	if err != nil {
		panic(err)
	}

	err = wfile.EncodeWat(f)
	if err != nil {
		panic(err)
	}

	err = f.Close()
	if err != nil {
		panic(err)
	}

	//
	f2, err := os.Create("debug.wasm")
	if err != nil {
		panic(err)
	}

	err = wfile.EncodeBinary(f2)
	if err != nil {
		panic(err)
	}

	err = f2.Close()
	if err != nil {
		panic(err)
	}

}
