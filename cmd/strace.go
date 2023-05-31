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
	"encoding/binary"
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/loopholelabs/wasm-toolkit/internal/wat"
	"github.com/loopholelabs/wasm-toolkit/wasmfile"
	"github.com/loopholelabs/wasm-toolkit/wasmfile/expression"
	"github.com/loopholelabs/wasm-toolkit/wasmfile/types"

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

var include_imports = false
var include_timings = false
var include_line_numbers = false
var include_func_signatures = false
var include_param_names = false
var include_all = false
var func_regex = ".*"
var cfg_color = false
var watch_globals = ""
var config_parse_dwarf = false

// If true, then we'll hook access to globals / locals, and output debug info...
var config_log_globals = false
var config_log_locals = false

func init() {
	rootCmd.AddCommand(cmdStrace)
	cmdStrace.Flags().StringVarP(&func_regex, "func", "f", ".*", "Func name regexp")
	cmdStrace.Flags().BoolVar(&include_line_numbers, "linenumbers", false, "Include line number info")
	cmdStrace.Flags().BoolVar(&include_func_signatures, "funcsignatures", false, "Include function signatures")
	cmdStrace.Flags().BoolVar(&include_param_names, "paramnames", false, "Include param names")
	cmdStrace.Flags().BoolVar(&include_timings, "timing", false, "Include timing summary")
	cmdStrace.Flags().BoolVar(&include_imports, "imports", false, "Include imports")
	cmdStrace.Flags().BoolVar(&include_all, "all", false, "Include everything")

	cmdStrace.Flags().BoolVar(&cfg_color, "color", false, "Output ANSI color in the log")
	cmdStrace.Flags().BoolVar(&config_parse_dwarf, "dwarf", false, "Parse dwarf line numbers and variables")

	cmdStrace.Flags().StringVarP(&watch_globals, "watch", "w", "", "List of globals to watch (, separated)")

	cmdStrace.Flags().BoolVar(&config_log_globals, "logglobals", false, "Log wasm global access")
	cmdStrace.Flags().BoolVar(&config_log_locals, "loglocals", false, "Log wasm local access")
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

	fmt.Printf("Parsing custom dwarf debug sections...\n")
	err = wfile.ParseDwarf()
	if err != nil {
		panic(err)
	}

	// Keep track of wasi import wrappers so that we can add context to them later.
	wasi_functions := make(map[int]string)

	// Wrap all imports if we need to...
	// Then they will get included in normal debug logging and or timing
	if include_all || include_imports {
		for idx, i := range wfile.Import {

			newidx := len(wfile.Import) + len(wfile.Code)

			// First we create a func wrapper, then adjust all calls
			f := &wasmfile.FunctionEntry{
				TypeIndex: i.Index,
			}

			t := wfile.Type[i.Index]

			// Load the params...
			expr := make([]*expression.Expression, 0)
			for idx := range t.Param {
				expr = append(expr, &expression.Expression{
					Opcode:     expression.InstrToOpcode["local.get"],
					LocalIndex: idx,
				})
			}

			expr = append(expr, &expression.Expression{
				Opcode:    expression.InstrToOpcode["call"],
				FuncIndex: idx,
			})

			c := &wasmfile.CodeEntry{
				Locals:     make([]types.ValType, 0),
				Expression: expr,
			}

			// Fixup any calls
			for _, c := range wfile.Code {
				c.ModifyAllCalls(map[int]int{idx: newidx})
			}

			// If they're wasi calls. Add function signatures etc
			if i.Module == "wasi_snapshot_preview1" {
				wasi_functions[newidx] = i.Name
				de, ok := wasmfile.Debug_wasi_snapshot_preview1[i.Name]
				if ok {
					wfile.SetFunctionSignature(newidx, de)
				}
			}

			wfile.FunctionNames[newidx] = fmt.Sprintf("$IMPORT_%s_%s", i.Module, i.Name) //wfile.GetFunctionIdentifier(idx, false))

			wfile.Function = append(wfile.Function, f)
			wfile.Code = append(wfile.Code, c)
		}
	}

	originalFunctionLength := len(wfile.Code)

	data_ptr := wfile.Memory[0].LimitMin << 16

	data_wasi_err := make([]byte, 0)
	data_wasi_err_ptrs := make([]byte, 0)

	errors_by_id := make([]string, 77)
	for m, v := range wasmfile.Wasi_errors {
		errors_by_id[v] = m
	}

	for _, m := range errors_by_id {
		data_wasi_err = binary.LittleEndian.AppendUint32(data_wasi_err, uint32(len(data_wasi_err_ptrs)))
		data_wasi_err = binary.LittleEndian.AppendUint32(data_wasi_err, uint32(len([]byte(m))))
		data_wasi_err_ptrs = append(data_wasi_err_ptrs, []byte(m)...)
	}

	/*
		datamap["$wasi_errors"] = we_data
		datamap["$wasi_error_messages"] = er_data
	*/

	//Wasi_errors

	// Load up the individual wat files, and add them in
	files := []string{
		"memory.wat",
		"stdout.wat",
		"strace.wat",
		"color.wat",
		"timings.wat",
		"watch.wat",
		"function_enter_exit.wat"}

	ptr := int32(data_ptr)
	for _, file := range files {
		fmt.Printf(" - Adding code from %s...\n", file)
		data, err := wat.Wat_content.ReadFile(path.Join("wat_code", file))
		if err != nil {
			panic(err)
		}

		mod := &wasmfile.WasmFile{}
		err = mod.DecodeWat(data)

		if err != nil {
			panic(err)
		}
		ptr = wfile.AddDataFrom(ptr, mod)
		wfile.AddFuncsFrom(mod, func(remap map[int]int) {
			// Fixup
			// wasi_functions
			newmap := make(map[int]string)
			for f, t := range remap {
				n, ok := wasi_functions[f]
				if ok {
					newmap[t] = n
				}
			}
			wasi_functions = newmap
		})
	}

	fmt.Printf("All wat code added...\n")

	wfile.SetGlobal("$debug_start_mem", types.ValI32, fmt.Sprintf("i32.const %d", data_ptr))

	if config_parse_dwarf {

		// Parse the dwarf stuff *here* incase the above messed up function IDs
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

	}

	// Get watch code
	watch_code := GetWatchCode(wfile)

	// Pass some config into wasm
	if include_timings {
		wfile.SetGlobal("$debug_do_timings", types.ValI32, fmt.Sprintf("i32.const 1"))
	}

	if cfg_color {
		wfile.SetGlobal("$wt_color", types.ValI32, fmt.Sprintf("i32.const 1"))
	}

	// Get a function name map, and add it as data...
	data_function_names := make([]byte, 0)
	data_function_locs := make([]byte, 0)
	data_metrics_data := make([]byte, 0)
	for idx := range wfile.Import {
		functionIndex := idx
		name := wfile.GetFunctionIdentifier(functionIndex, false)

		data_function_locs = binary.LittleEndian.AppendUint32(data_function_locs, uint32(len(data_function_names)))
		data_function_locs = binary.LittleEndian.AppendUint32(data_function_locs, uint32(len([]byte(name))))

		data_function_names = append(data_function_names, []byte(name)...)

		// Just add another 16 bytes on for now...
		data_metrics_data = append(data_metrics_data, make([]byte, 16)...)
	}

	for idx := range wfile.Code {
		functionIndex := len(wfile.Import) + idx
		name := wfile.GetFunctionIdentifier(functionIndex, false)

		data_function_locs = binary.LittleEndian.AppendUint32(data_function_locs, uint32(len(data_function_names)))
		data_function_locs = binary.LittleEndian.AppendUint32(data_function_locs, uint32(len([]byte(name))))

		data_function_names = append(data_function_names, []byte(name)...)

		// Just add another 16 bytes on for now...
		data_metrics_data = append(data_metrics_data, make([]byte, 16)...)
	}

	// Add those data elements into the mix...

	wfile.AddData("$wasi_errors", []byte(data_wasi_err))
	wfile.AddData("$wasi_error_messages", []byte(data_wasi_err_ptrs))

	wfile.AddData("$wt_all_function_names", []byte(data_function_names))
	wfile.AddData("$wt_all_function_names_locs", []byte(data_function_locs))
	wfile.AddData("$metrics_data", []byte(data_metrics_data))
	wfile.SetGlobal("$wt_all_function_length", types.ValI32, fmt.Sprintf("i32.const %d", len(wfile.Import)+len(wfile.Code)))

	fmt.Printf("Patching functions matching regexp \"%s\"\n", func_regex)

	// Adjust any memory.size / memory.grow calls
	for idx, c := range wfile.Code {
		fmt.Printf("Processing functions [%d/%d]\n", idx, len(wfile.Code))

		fn := wfile.Function[idx]
		functionType := wfile.Type[fn.TypeIndex]

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
			fidentifier := wfile.GetFunctionIdentifier(functionIndex, false)

			match, err := regexp.MatchString(func_regex, fidentifier)
			if err != nil {
				panic(err)
			}

			if match {
				fmt.Printf("Patching function[%d] %s\n", idx, fidentifier)
				// If it's a wasi call, then output some detail here...
				wasi_name, is_wasi := wasi_functions[functionIndex]

				if is_wasi {
					fmt.Printf(" (Wasi call to %s)\n", wasi_name)
				}

				blockInstr := "block"
				f := wfile.Function[idx]
				t := wfile.Type[f.TypeIndex]
				if len(t.Result) > 0 {
					blockInstr = fmt.Sprintf("block (result %s)", types.ByteToValType[t.Result[0]])
				}

				startCode := fmt.Sprintf(`%s
			i32.const %d
			call $debug_enter_func
			`, blockInstr, functionIndex)

				// Do parameters...
				for paramIndex, pt := range t.Param {
					if paramIndex > 0 {
						startCode = fmt.Sprintf(`%s
					call $debug_param_separator
					`, startCode)
					}

					// NB This assumes CodeSectionPtr to be correct...
					if include_all || include_param_names {
						if c.PCValid {
							vname := wfile.GetLocalVarName(c.CodeSectionPtr, paramIndex)
							if vname != "" {
								wfile.AddData(fmt.Sprintf("$dd_param_name_%d_%d", functionIndex, paramIndex), []byte(vname))
								startCode = fmt.Sprintf(`%s
					i32.const offset($dd_param_name_%d_%d)
					i32.const length($dd_param_name_%d_%d)
					call $debug_param_name
					`, startCode, functionIndex, paramIndex, functionIndex, paramIndex)
							}
						}
					}
					startCode = fmt.Sprintf(`%s
					i32.const %d
					i32.const %d
					local.get %d
					call $debug_enter_%s
					`, startCode, functionIndex, paramIndex, paramIndex, types.ByteToValType[pt])
				}

				startCode = fmt.Sprintf(`%s
					i32.const %d
					call $debug_enter_end
					`, startCode, functionIndex)

				// Now add a bit of debug....
				funcSig := wfile.GetFunctionSignature(functionIndex)
				if funcSig != "" && (include_all || include_func_signatures) {
					wfile.AddData(fmt.Sprintf("$dd_function_debug_sig_%d", functionIndex), []byte(funcSig))
					startCode = fmt.Sprintf(`%s
					i32.const offset($dd_function_debug_sig_%d)
					i32.const length($dd_function_debug_sig_%d)
					call $debug_func_context`, startCode, functionIndex, functionIndex)
				}

				lineRange := wfile.GetLineNumberRange(c)
				if lineRange != "" && (include_all || include_line_numbers) {
					wfile.AddData(fmt.Sprintf("$dd_function_debug_lines_%d", functionIndex), []byte(lineRange))
					startCode = fmt.Sprintf(`%s
					i32.const offset($dd_function_debug_lines_%d)
					i32.const length($dd_function_debug_lines_%d)
					call $debug_func_context
					`, startCode, functionIndex, functionIndex)
				}

				// Add some code to show function parameter values...
				startCode = fmt.Sprintf(`%s
					%s`, startCode, wasmfile.GetWasiParamCodeEnter(wasi_name))

				if include_timings {
					startCode = fmt.Sprintf(`%s
					i32.const %d
					call $timings_enter_func
					`, startCode, functionIndex)
				}

				// Add any watches
				if watch_globals != "" {
					startCode = fmt.Sprintf(`%s
					%s`, startCode, watch_code)
				}

				err = c.InsertFuncStart(wfile, startCode)
				if err != nil {
					panic(err)
				}

				rt := types.ValNone
				if len(t.Result) == 1 {
					rt = t.Result[0]
				}

				endCode := ""

				if include_timings {
					endCode = fmt.Sprintf(`%s
					i32.const %d
					call $timings_exit_func
					`, endCode, functionIndex)
				}

				endCode = fmt.Sprintf(`%s
				i32.const %d
				call $debug_exit_func`, endCode, functionIndex)

				if is_wasi && rt == types.ValI32 {
					// We also want to output the error message
					endCode = fmt.Sprintf(`%s
					call $debug_exit_func_wasi
					%s`, endCode, wasmfile.GetWasiParamCodeExit(wasi_name))

				} else {
					endCode = fmt.Sprintf(`%s
					call $debug_exit_func_%s`, endCode, types.ByteToValType[rt])
				}

				// Add any watches
				if watch_globals != "" {
					endCode = fmt.Sprintf(`%s
					%s`, endCode, watch_code)
				}

				err = c.ReplaceInstr(wfile, "return", endCode+"\nreturn")
				if err != nil {
					panic(err)
				}

				err = c.InsertFuncEnd(wfile, "end\n"+endCode)
				if err != nil {
					panic(err)
				}

				// Add local / global logging...
				if config_log_globals || config_log_locals {
					newCode := make([]*expression.Expression, 0)
					for _, e := range c.Expression {
						if config_log_globals &&
							e.Opcode == expression.InstrToOpcode["global.set"] &&
							!e.GlobalNeedsLinking {
							g := wfile.Global[e.GlobalIndex]
							gtype := types.ByteToValType[g.Type]
							linei := wfile.GetLineNumberBefore(c, e.PC)
							// Add some debug data for this global.set
							gdebug := fmt.Sprintf("global.set %d %s:%x %s", e.GlobalIndex, fidentifier, e.PC, linei)
							wfile.AddData(fmt.Sprintf("$dd_global_set_%d", e.PC), []byte(gdebug))

							wcode := fmt.Sprintf(`
							global.get %d
							i32.const offset($dd_global_set_%d)
							i32.const length($dd_global_set_%d)
							call $log_global_%s
							`, e.GlobalIndex, e.PC, e.PC, gtype)

							// $log_global_<TYPE> (new_value, current_value, ptr_debug, len_debug) => new_value

							wcex, err := expression.ExpressionFromWat(wcode)
							if err != nil {
								panic(err)
							}
							newCode = append(newCode, wcex...)

						} else if config_log_locals &&
							(e.Opcode == expression.InstrToOpcode["local.set"] || e.Opcode == expression.InstrToOpcode["local.tee"]) {

							vname := wfile.GetLocalVarName(e.PC, e.LocalIndex)

							var ltype string
							var debugPrefix string
							if e.LocalIndex >= len(functionType.Param) {
								l := c.Locals[e.LocalIndex-len(functionType.Param)]
								ltype = types.ByteToValType[l]
								debugPrefix = "local.set"
							} else {
								// Mutating/reusing params
								l := functionType.Param[e.LocalIndex]
								ltype = types.ByteToValType[l]
								debugPrefix = "local.set(param)"
							}
							linei := wfile.GetLineNumberBefore(c, e.PC)
							ldebug := fmt.Sprintf("%s %d %s %s:%x %s", debugPrefix, e.LocalIndex, vname, fidentifier, e.PC, linei)
							wfile.AddData(fmt.Sprintf("$dd_local_set_%d", e.PC), []byte(ldebug))

							wcode := fmt.Sprintf(`
								local.get %d
								i32.const offset($dd_local_set_%d)
								i32.const length($dd_local_set_%d)
								call $log_local_%s
								`, e.LocalIndex, e.PC, e.PC, ltype)

							// $log_local_<TYPE> (new_value, current_value, ptr_debug, len_debug) => new_value

							wcex, err := expression.ExpressionFromWat(wcode)
							if err != nil {
								panic(err)
							}
							newCode = append(newCode, wcex...)

						}
						newCode = append(newCode, e)
					}
					c.Expression = newCode
				}
			}
		}

		err = c.InsertAfterRelocating(wfile, `global.get $debug_start_mem
		i32.add`)
		if err != nil {
			panic(err)
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

func GetWatchCode(wf *wasmfile.WasmFile) string {
	if watch_globals == "" {
		return ""
	}

	code := ""
	watches := strings.Split(watch_globals, ",")
	for widx, w := range watches {
		// Lookup the address...
		ginfo, ok := wf.GlobalAddresses[w]
		if !ok {
			fmt.Printf("WARNING: I can't find the global %s\n", w)
			for n := range wf.GlobalAddresses {
				fmt.Printf(" - Global %s\n", n)
			}
			panic("Global name not found")
		} else {
			// Insert some code to show global...
			wf.AddData(fmt.Sprintf("$watch_name_%d", widx), []byte(w))

			code = fmt.Sprintf(`%s
				i32.const offset($watch_name_%d)
				i32.const length($watch_name_%d)
				i32.const %d
				call $wt_watch_i32
			`, code, widx, widx, uint32(ginfo.Address))
		}
	}
	return code
}
