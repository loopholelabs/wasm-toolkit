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

	wasmfile "github.com/loopholelabs/wasm-toolkit/pkg/wasm"
	"github.com/loopholelabs/wasm-toolkit/pkg/wasm/debug"
	"github.com/spf13/cobra"
)

var (
	cmdWasm2Wat = &cobra.Command{
		Use:   "wasm2wat",
		Short: "Use wasm2wat to translate a wasm file to wat",
		Long:  `This will include any dwarf debug information available.`,
		Run:   runWasm2Wat,
	}
)

func init() {
	rootCmd.AddCommand(cmdWasm2Wat)
}

func runWasm2Wat(ccmd *cobra.Command, args []string) {
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
}
