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

	"github.com/loopholelabs/wasm-toolkit/pkg/otel"
	"github.com/spf13/cobra"
)

var (
	cmdOtel = &cobra.Command{
		Use:   "otel",
		Short: "Add tracing output to as wasm file, output as otel json",
		Long:  `This will output to STDERR`,
		Run:   runOtel,
	}
)

var otel_func_regex = ".*"
var otel_quickjs = false
var is_scale_host = false

func init() {
	rootCmd.AddCommand(cmdOtel)
	cmdOtel.Flags().StringVarP(&otel_func_regex, "func", "f", ".*", "Func name regexp")
	cmdOtel.Flags().BoolVarP(&otel_quickjs, "qjs", "j", false, "Do quickjs otel")
	cmdOtel.Flags().BoolVarP(&is_scale_host, "scale", "s", false, "Is scale host")
}

func runOtel(ccmd *cobra.Command, args []string) {
	if Input == "" {
		panic("No input file")
	}

	fmt.Printf("Loading wasm file \"%s\"...\n", Input)
	data, err := os.ReadFile(Input)
	if err != nil {
		panic(err)
	}

	config := otel.Otel_config{
		Func_regexp: otel_func_regex,
		Quickjs:     otel_quickjs,
		Scale_api:   is_scale_host,
	}
	newdata, err := otel.AddOtel(data, config)

	fmt.Printf("Writing wasm out to %s...\n", Output)
	os.WriteFile(Output, newdata, 0660)
}
