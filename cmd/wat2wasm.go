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

	"github.com/loopholelabs/wasm-toolkit/pkg/wasm/wasmfile"

	"github.com/spf13/cobra"
)

var (
	cmdWat2Wasm = &cobra.Command{
		Use:   "wat2wasm",
		Short: "Use wat2wasm to translate a wat file to wasm",
		Long:  ``,
		Run:   runWat2Wasm,
	}
)

func init() {
	rootCmd.AddCommand(cmdWat2Wasm)
}

func runWat2Wasm(ccmd *cobra.Command, args []string) {
	if Input == "" {
		panic("No input file")
	}

	fmt.Printf("Loading wat file \"%s\"...\n", Input)
	wfile, err := wasmfile.NewFromWat(Input)
	if err != nil {
		panic(err)
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
}
