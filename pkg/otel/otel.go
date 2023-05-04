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

package otel

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"path"
	"regexp"

	"github.com/loopholelabs/wasm-toolkit/internal/wat"
	"github.com/loopholelabs/wasm-toolkit/wasmfile"
)

type Otel_config struct {
	Func_regexp string
	Quickjs     bool
	Scale_api   bool
}

/**
 * Add otel to a wasm.
 *
 */
func AddOtel(wasmInput []byte, config Otel_config) ([]byte, error) {
	// First parse the wasm binary
	wfile := &wasmfile.WasmFile{}
	err := wfile.DecodeBinary(wasmInput)
	if err != nil {
		return nil, err
	}

	// Parse custom name section
	err = wfile.ParseName()
	if err != nil {
		return nil, err
	}

	// Parsing custom dwarf debug section
	err = wfile.ParseDwarf()
	if err != nil {
		return nil, err
	}

	// Wrap import functions
	wasi_functions := wrapImports(wfile)

	originalFunctionLength := len(wfile.Code)

	data_ptr := wfile.Memory[0].LimitMin << 16

	// Load up the individual wat files, and add them in
	files := []string{
		"memory.wat",
		"stdout.wat",
		"otel.wat"}

	if config.Quickjs {
		files = append(files, "quickjs.wat")
	}

	if config.Scale_api {
		files = append(files, "otel_scale.wat")
	} else {
		files = append(files, "otel_stderr.wat")
	}

	// Add the code from wat files...
	ptr := int32(data_ptr)
	for _, file := range files {
		data, err := wat.Wat_content.ReadFile(path.Join("wat_code", file))
		if err != nil {
			return nil, err
		}
		mod := &wasmfile.WasmFile{}
		err = mod.DecodeWat(data)
		if err != nil {
			return nil, err
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

	wfile.SetGlobal("$debug_start_mem", wasmfile.ValI32, fmt.Sprintf("i32.const %d", data_ptr))

	// Parse the dwarf stuff *here*
	err = wfile.ParseDwarfLineNumbers()
	if err != nil {
		return nil, err
	}

	err = wfile.ParseDwarfVariables()
	if err != nil {
		return nil, err
	}

	// Add the wasi error info
	addWasiErrorInfo(wfile)
	// Add function info
	addFunctionInfo(wfile)

	// Now do function adjustments
	for idx, c := range wfile.Code {
		if idx < originalFunctionLength {
			// First deal with our memory adjustment
			err = c.ReplaceInstr(wfile, "memory.grow", "call $debug_memory_grow")
			if err != nil {
				return nil, err
			}
			err = c.ReplaceInstr(wfile, "memory.size", "call $debug_memory_size")
			if err != nil {
				return nil, err
			}

			functionIndex := idx + len(wfile.Import)
			fidentifier := wfile.GetFunctionIdentifier(functionIndex, false)

			match, err := regexp.MatchString(config.Func_regexp, fidentifier)
			if err != nil {
				return nil, err
			}

			//IF the function matches, patch it.
			if match {
				f := wfile.Function[idx]
				t := wfile.Type[f.TypeIndex]

				// First deal with args. Create a mirror set of locals and copy params there.
				// This is so that they're definitely available unmodified at the function exit.

				local_index_local := len(t.Param)
				local_index_mirrored_params := local_index_local + len(c.Locals)

				new_locals := make([]wasmfile.ValType, 0)
				new_locals = append(new_locals, c.Locals...)
				for _, vt := range t.Param {
					new_locals = append(new_locals, vt)
				}
				// Now set the new locals, and setup some code to mirror the params into these new locals
				c.Locals = new_locals

				blockInstr := "block"
				if len(t.Result) > 0 {
					blockInstr = fmt.Sprintf("block (result %s)", wasmfile.ByteToValType[t.Result[0]])
				}

				startCode := fmt.Sprintf(`%s
					i32.const %d
					call $otel_enter_func`, blockInstr, functionIndex)

				// Copy the params to our locals for func exit
				for idx := range t.Param {
					target_idx := local_index_mirrored_params + idx
					startCode = fmt.Sprintf(`%s
						local.get %d
						local.set %d`, startCode, idx, target_idx)
				}

				err = c.InsertFuncStart(wfile, startCode)
				if err != nil {
					return nil, err
				}

				// Now do the endCode
				endCode := ""

				// If we're doing quickJS
				if config.Quickjs && fidentifier == "$JS_CallInternal" {

					endCode = fmt.Sprintf(`%s
						local.get %d
						local.get %d
						call $qjs_otel_exit_func`, endCode,
						local_index_mirrored_params,   // context
						local_index_mirrored_params+1, // funcObj
					)

					endCode = fmt.Sprintf(`%s
						local.get %d
						local.get %d
						local.get %d
						local.get %d
						local.get %d
						call $otel_quickjs_call
						`, endCode,
						local_index_mirrored_params,   // context
						local_index_mirrored_params+1, // funcObj
						local_index_mirrored_params+2, // thisObj
						local_index_mirrored_params+4, // argc
						local_index_mirrored_params+5) // args

					endCode = fmt.Sprintf(`%s
						i32.const %d
						call $otel_exit_func_done`, endCode, functionIndex)
				} else {

					endCode = fmt.Sprintf(`%s
						i32.const %d
						call $otel_exit_func`, endCode, functionIndex)

					// Process params for exit
					for idx, vt := range t.Param {
						target_idx := local_index_mirrored_params + idx
						endCode = fmt.Sprintf(`%s
							i32.const %d
							i32.const %d
							local.get %d
							call $otel_exit_func_%s`, endCode, functionIndex, idx, target_idx, wasmfile.ByteToValType[vt])
					}

					// Add result values to the trace
					if len(t.Result) == 1 {
						rt := t.Result[0]
						endCode = fmt.Sprintf(`%s
							i32.const %d
							call $otel_exit_func_result_%s`, endCode, functionIndex, wasmfile.ByteToValType[rt])
					}

					endCode = fmt.Sprintf(`%s
						i32.const %d
						call $otel_exit_func_done`, endCode, functionIndex)
				}

				err = c.ReplaceInstr(wfile, "return", endCode+"\nreturn")
				if err != nil {
					return nil, err
				}

				err = c.InsertFuncEnd(wfile, "end\n"+endCode)
				if err != nil {
					return nil, err
				}

			}
		}

		// We need to fixup any instructions that refer to data that may move.
		err = c.InsertAfterRelocating(wfile, `global.get $debug_start_mem
			i32.add`)
		if err != nil {
			return nil, err
		}

		err = c.ResolveLengths(wfile)
		if err != nil {
			return nil, err
		}

		err = c.ResolveRelocations(wfile, data_ptr)
		if err != nil {
			return nil, err
		}

		err = c.ResolveGlobals(wfile)
		if err != nil {
			return nil, err
		}

		err = c.ResolveFunctions(wfile)
		if err != nil {
			return nil, err
		}
	}

	// Find out how much data we need for the payload
	total_payload_data := data_ptr
	if len(wfile.Data) > 0 {
		last_data := wfile.Data[len(wfile.Data)-1]
		total_payload_data = int(last_data.Offset[0].I32Value) + len(last_data.Data) - data_ptr
	}

	payload_size := (total_payload_data + 65535) >> 16

	wfile.SetGlobal("$debug_mem_size", wasmfile.ValI32, fmt.Sprintf("i32.const %d", payload_size)) // The size of our addition in 64k pages
	wfile.Memory[0].LimitMin = wfile.Memory[0].LimitMin + payload_size

	var buf bytes.Buffer
	err = wfile.EncodeBinary(&buf)

	return buf.Bytes(), nil
}

/**
 * wrap all imports with new functions
 *
 */
func wrapImports(wfile *wasmfile.WasmFile) map[int]string {
	// Keep track of wasi import wrappers so that we can add context to them later.
	wasi_functions := make(map[int]string)

	// Wrap all imports...
	// Then they will get included in normal debug logging and or timing
	for idx, i := range wfile.Import {
		newidx := len(wfile.Import) + len(wfile.Code)

		// First we create a func wrapper, then adjust all calls
		f := &wasmfile.FunctionEntry{
			TypeIndex: i.Index,
		}

		t := wfile.Type[i.Index]

		// Load the params...
		expr := make([]*wasmfile.Expression, 0)
		for idx := range t.Param {
			expr = append(expr, &wasmfile.Expression{
				Opcode:     wasmfile.InstrToOpcode["local.get"],
				LocalIndex: idx,
			})
		}

		expr = append(expr, &wasmfile.Expression{
			Opcode:    wasmfile.InstrToOpcode["call"],
			FuncIndex: idx,
		})

		c := &wasmfile.CodeEntry{
			Locals:     make([]wasmfile.ValType, 0),
			Expression: expr,
		}

		// Fixup any calls to call the new function wrapper instead of the import.
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

		// Set the function name
		wfile.FunctionNames[newidx] = fmt.Sprintf("$IMPORT_%s_%s", i.Module, i.Name)

		// Add the new function/code.
		wfile.Function = append(wfile.Function, f)
		wfile.Code = append(wfile.Code, c)
	}

	return wasi_functions
}

/**
 * Add wasi error info into the wasm file as $wasi_errors / $wasi_error_messages
 *
 */
func addWasiErrorInfo(wfile *wasmfile.WasmFile) {
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

	wfile.AddData("$wasi_errors", []byte(data_wasi_err))
	wfile.AddData("$wasi_error_messages", []byte(data_wasi_err_ptrs))
}

/**
 * Add function info to the wasm file.
 *
 */
func addFunctionInfo(wfile *wasmfile.WasmFile) {
	// Get a function name map, and add it as data...
	data_function_names := make([]byte, 0)
	data_function_names_locs := make([]byte, 0)
	data_function_sigs := make([]byte, 0)
	data_function_sigs_locs := make([]byte, 0)
	data_function_srcs := make([]byte, 0)
	data_function_srcs_locs := make([]byte, 0)

	num_functions := len(wfile.Import) + len(wfile.Code)

	for idx := 0; idx < num_functions; idx++ {
		functionIndex := idx
		name := wfile.GetFunctionIdentifier(functionIndex, false)
		signature := wfile.GetFunctionSignature(functionIndex)
		debug := ""
		if idx >= len(wfile.Import) {
			c := wfile.Code[idx-len(wfile.Import)]
			debug = wfile.GetLineNumberRange(c)
		}

		data_function_names_locs = binary.LittleEndian.AppendUint32(data_function_names_locs, uint32(len(data_function_names)))
		data_function_names_locs = binary.LittleEndian.AppendUint32(data_function_names_locs, uint32(len([]byte(name))))
		data_function_names = append(data_function_names, []byte(name)...)

		data_function_sigs_locs = binary.LittleEndian.AppendUint32(data_function_sigs_locs, uint32(len(data_function_sigs)))
		data_function_sigs_locs = binary.LittleEndian.AppendUint32(data_function_sigs_locs, uint32(len([]byte(signature))))
		data_function_sigs = append(data_function_sigs, []byte(signature)...)

		data_function_srcs_locs = binary.LittleEndian.AppendUint32(data_function_srcs_locs, uint32(len(data_function_srcs)))
		data_function_srcs_locs = binary.LittleEndian.AppendUint32(data_function_srcs_locs, uint32(len([]byte(debug))))
		data_function_srcs = append(data_function_srcs, []byte(debug)...)
	}

	wfile.AddData("$wt_all_function_names", []byte(data_function_names))
	wfile.AddData("$wt_all_function_names_locs", []byte(data_function_names_locs))
	wfile.AddData("$wt_all_function_sigs", []byte(data_function_sigs))
	wfile.AddData("$wt_all_function_sigs_locs", []byte(data_function_sigs_locs))
	wfile.AddData("$wt_all_function_srcs", []byte(data_function_srcs))
	wfile.AddData("$wt_all_function_srcs_locs", []byte(data_function_srcs_locs))
	wfile.SetGlobal("$wt_all_function_length", wasmfile.ValI32, fmt.Sprintf("i32.const %d", num_functions))
}
