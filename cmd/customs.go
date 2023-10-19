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

var muxDefImport = ""
var muxDefExport = ""

func init() {
	rootCmd.AddCommand(cmdCustoms)

	cmdCustoms.Flags().StringVar(&muxDefImport, "muximport", "", "Definition for mux import")
	cmdCustoms.Flags().StringVar(&muxDefExport, "muxexport", "", "Definition for mux export")
}

// Example:
//	--muximport "env/hello,0:env/zero,1:env/one,2:env/two"
//	--muxexport "resize,0:resize_zero,1:resize_one,2:resize_two"

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

	if muxDefImport != "" {
		ci, err := customs.ParseRemapMuxImport(muxDefImport)
		if err != nil {
			panic(err)
		}
		err = customs.MuxImport(wfile, *ci)
		if err != nil {
			panic(err)
		}
	}

	if muxDefExport != "" {
		ci, err := customs.ParseRemapMuxExport(muxDefExport)
		if err != nil {
			panic(err)
		}
		err = customs.MuxExport(wfile, *ci)
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
