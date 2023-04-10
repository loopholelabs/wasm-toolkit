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

var Input string

func init() {
	rootCmd.AddCommand(cmdWasm2Wat)
	cmdWasm2Wat.Flags().StringVarP(&Input, "input", "i", "", "Input file name")
}

func runWasm2Wat(ccmd *cobra.Command, args []string) {
	if Input == "" {
		panic("No input file")
	}
	// executes what cmdOne is supposed to do
	fmt.Printf("Wasm2wat...%s \n", Input)

}
