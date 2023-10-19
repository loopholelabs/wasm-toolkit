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

	"github.com/loopholelabs/wasm-toolkit/pkg/wasm/debug"
	"github.com/loopholelabs/wasm-toolkit/pkg/wasm/wasmfile"

	"github.com/loopholelabs/wasm-toolkit/pkg/customs"

	"github.com/spf13/cobra"
)

var (
	cmdCustoms = &cobra.Command{
		Use:   "customs",
		Short: "Manipulate import/export",
		Long:  `This manipulates imports and exports`,
		Run:   runCustoms,
	}
)

func init() {
	rootCmd.AddCommand(cmdCustoms)

	//cmdCustoms.Flags().StringVar(&source_file, "filename", "", "Source filename")
}

func runCustoms(ccmd *cobra.Command, args []string) {
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

	c := customs.RemapMuxImport{
		Source: customs.Import{
			Module: "env", Name: "hello",
		},
		Mapper: map[uint64]customs.Import{
			0: {
				Module: "env", Name: "zero",
			},
			1: {
				Module: "env", Name: "one",
			},
			2: {
				Module: "env", Name: "two",
			},
		},
	}

	fmt.Printf("Remap %v\n", c)

	customs.MuxImport(wfile, c)

	ce := customs.RemapMuxExport{
		Source: "resize",
		Mapper: map[uint64]string{
			0: "resize_zero",
			1: "resize_one",
			2: "resize_two",
		},
	}

	fmt.Printf("Remap %v\n", ce)

	customs.MuxExport(wfile, ce)

	// DONE DONE DONE...

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

}
