package main

import (
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"github.com/loopholelabs/wasm-lnkr/wasm"
)

func main() {

	args := os.Args[1:]

	// Just a single arg for now, the input filename...
	watfile := args[0]

	wasmfile := args[1] // For now...

	w := wasm.NewWasmFile(wasmfile)

	err := w.ReadDwarf()
	if err != nil {
		panic(err)
	}

	err = w.ReadFunctionInfo()
	if err != nil {
		panic(err)
	}

	err = w.ParseDwarfLineNumbers()
	if err != nil {
		panic(err)
	}

	log.Print("Dwarf function info loaded up...")

	for fnid, fi := range w.FunctionLocations {
		log.Printf("Function %d - %s (%d-%d)\n", fnid, fi.LineFile, fi.LineMin, fi.LineMax)
	}

	hook_stores := false

	log.Print("Loading mod...")
	module := wasm.NewModule(watfile)
	module.Parse()

	// First of all, lets wrap imports in functions
	for _, i := range module.Imports {
		ff, callName := module.WrapImport(i)
		// Add it on, and do any translations in other functions...
		module.ReplaceInstruction(fmt.Sprintf("call %s", i.GetFuncName()), fmt.Sprintf("call %s", callName))
		// NB Add it *after* so that the call inside the function doesn't get replaced.
		module.Funcs = append(module.Funcs, ff)
	}

	// Go through each function, and add some debugging to it.
	for findex, f := range module.Funcs {
		ins := make([]string, 0)

		// Add a header to the function.
		ins = append(ins, fmt.Sprintf(
			"block %s", f.Result),
			fmt.Sprintf("i32.const %d", findex),
			"call $debug_enter_func")

		// Go through the function parameters, and do stuff with them
		pindex := 0
		for _, pa := range f.Params {
			if strings.HasPrefix(pa, "(param ") {
				bits := pa[7 : len(pa)-1]
				words := strings.Fields(bits)

				for _, w := range words {
					if w == "i32" {
						ins = append(ins, fmt.Sprintf("i32.const %d", findex),
							fmt.Sprintf("i32.const %d", pindex),
							fmt.Sprintf("local.get %d", pindex),
							"call $debug_enter_i32")
					} else if w == "i64" {
						ins = append(ins, fmt.Sprintf("i32.const %d", findex),
							fmt.Sprintf("i32.const %d", pindex),
							fmt.Sprintf("local.get %d", pindex),
							"call $debug_enter_i64")
					} else if w == "f32" {
						ins = append(ins, fmt.Sprintf("i32.const %d", findex),
							fmt.Sprintf("i32.const %d", pindex),
							fmt.Sprintf("local.get %d", pindex),
							"call $debug_enter_f32")
					} else if w == "f64" {
						ins = append(ins, fmt.Sprintf("i32.const %d", findex),
							fmt.Sprintf("i32.const %d", pindex),
							fmt.Sprintf("local.get %d", pindex),
							"call $debug_enter_f64")
					}
					pindex++
				}
			}
		}

		// Parameters are finished
		ins = append(ins, fmt.Sprintf("i32.const %d", findex),
			"call $debug_enter_end")

		// Now insert function footers anywhere the function will return.
		for _, i := range f.Instructions {
			if strings.Trim(i, wasm.Whitespace) == "return" {
				ins = append(ins, fmt.Sprintf("i32.const %d", findex),
					"call $debug_exit_func")
				if f.Result == "(result i32)" {
					ins = append(ins, "call $debug_exit_func_i32")
				} else if f.Result == "(result i64)" {
					ins = append(ins, "call $debug_exit_func_i64")
				} else if f.Result == "(result f32)" {
					ins = append(ins, "call $debug_exit_func_f32")
				} else if f.Result == "(result f64)" {
					ins = append(ins, "call $debug_exit_func_f64")
				} else {
					ins = append(ins, "call $debug_exit_func_none")
				}
			}
			ins = append(ins, i)
		}

		// Add a footer at the end of the function. NB we add a block/end to prevent early function return.
		ins = append(ins, "end")

		ins = append(ins, fmt.Sprintf("i32.const %d", findex),
			"call $debug_exit_func")
		if f.Result == "(result i32)" {
			ins = append(ins, "call $debug_exit_func_i32")
		} else if f.Result == "(result i64)" {
			ins = append(ins, "call $debug_exit_func_i64")
		} else if f.Result == "(result f32)" {
			ins = append(ins, "call $debug_exit_func_f32")
		} else if f.Result == "(result f64)" {
			ins = append(ins, "call $debug_exit_func_f64")
		} else {
			ins = append(ins, "call $debug_exit_func_none")
		}

		f.Instructions = ins

		// Optionally hook all memory stores...
		if hook_stores {
			f.FixMemoryInstrOffsetAlign("i32.store", "$debug_i32.store")
			f.FixMemoryInstrOffsetAlign("i32.store8", "$debug_i32.store8")
			f.FixMemoryInstrOffsetAlign("i32.store16", "$debug_i32.store16")
			f.FixMemoryInstrOffsetAlign("i64.store", "$debug_i64.store")
			f.FixMemoryInstrOffsetAlign("i64.store8", "$debug_i64.store8")
			f.FixMemoryInstrOffsetAlign("i64.store16", "$debug_i64.store16")
			f.FixMemoryInstrOffsetAlign("i64.store32", "$debug_i64.store32")
			f.FixMemoryInstrOffsetAlign("f32.store", "$debug_f32.store")
			f.FixMemoryInstrOffsetAlign("f64.store", "$debug_f64.store")
		}
	}

	// Hook memory.size and memory.grow
	for _, f := range module.Funcs {
		for _, i := range f.Instructions {
			if i == "memory.size" {
				i = "call $debug_memory.size"
			} else if i == "memory.grow" {
				i = "call $debug_memory.grow"
			}
		}
	}

	// Add a page of memory for our debug info
	debug_mem_size := 1
	debug_mem_start := module.Memorys[0].Size << 16

	module.Memorys[0].Size = module.Memorys[0].Size + debug_mem_size

	module.Globals = append(module.Globals, wasm.NewGlobal(fmt.Sprintf("(global $debug_mem_size i32 (i32.const %d))", debug_mem_size)))
	module.Globals = append(module.Globals, wasm.NewGlobal(fmt.Sprintf("(global $debug_start_mem (mut i32) (i32.const %d))", debug_mem_start)))

	// Write the function names and build a table of addresses / sizes
	functionNameTable := make([]byte, 0)
	functionNameData := make([]byte, 0)
	functionNameMetrics := make([]byte, 0)
	for findex, f := range module.Funcs {
		functioninfo, ok := w.FunctionLocations[findex]
		fullFunction := f.Identifier
		if ok {
			// Use the augmented dwarf info
			fullFunction = fmt.Sprintf("%s %s(%d-%d)", fullFunction, functioninfo.LineFile, functioninfo.LineMin, functioninfo.LineMax)
		}

		// First add the address and length to our table
		bs := make([]byte, 4)
		binary.LittleEndian.PutUint32(bs, uint32(len(functionNameData)))
		functionNameTable = append(functionNameTable, bs...)
		binary.LittleEndian.PutUint32(bs, uint32(len(fullFunction)))
		functionNameTable = append(functionNameTable, bs...)

		functionNameMetrics = append(functionNameMetrics, make([]byte, 16)...)
		// i32 Count
		// i64 running time

		// Now add the string onto our data
		functionNameData = append(functionNameData, []byte(fullFunction)...)
	}

	// Now we need to write these as data items...

	constants := make(map[string]int, 0)

	module.Globals = append(module.Globals, wasm.NewGlobal(fmt.Sprintf("(global $debug_num_funcs i32 (i32.const %d))", len(module.Funcs))))

	data_ptr := 0
	debug_function_table := wasm.EncodeData(functionNameTable)
	module.Datas = append(module.Datas, wasm.NewData(fmt.Sprintf("(data $debug_function_table (i32.const %d) %s)", data_ptr+debug_mem_start, debug_function_table)))
	constants["offset.$debug_function_table"] = 0
	constants["length.$debug_function_table"] = len(functionNameTable)
	data_ptr += len(functionNameTable)

	debug_function_data := wasm.EncodeData(functionNameData)
	module.Datas = append(module.Datas, wasm.NewData(fmt.Sprintf("(data $debug_function_data (i32.const %d) %s)", data_ptr+debug_mem_start, debug_function_data)))
	constants["offset.$debug_function_data"] = data_ptr
	constants["length.$debug_function_data"] = len(functionNameData)
	data_ptr += len(functionNameData)

	debug_function_metrics := wasm.EncodeData(functionNameMetrics)
	module.Datas = append(module.Datas, wasm.NewData(fmt.Sprintf("(data $debug_function_metrics (i32.const %d) %s)", data_ptr+debug_mem_start, debug_function_metrics)))
	constants["offset.$debug_function_metrics"] = data_ptr
	constants["length.$debug_function_metrics"] = len(functionNameMetrics)
	data_ptr += len(functionNameMetrics)

	// TODO: Load everything in debug dir and parse include it

	include_dir := "wat_code"

	files, err := ioutil.ReadDir(include_dir)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		// Include the wat_code...

		if !file.IsDir() {

			log.Printf("Loading %s...\n", path.Join(include_dir, file.Name()))
			debugModule := wasm.NewModule(path.Join(include_dir, file.Name()))
			debugModule.Parse()

			// Add the debug functions
			for _, f := range debugModule.Funcs {
				module.Funcs = append(module.Funcs, f)
			}

			// Add the debug globals
			for _, g := range debugModule.Globals {
				module.Globals = append(module.Globals, g)
			}

			// Add the debug imports
			for _, g := range debugModule.Imports {
				module.Imports = append(module.Imports, g)
			}

			// Now add any more datas on, and create consts for them. Replace in any code...
			for _, dd := range debugModule.Datas {
				dd.Location = fmt.Sprintf("(i32.const %d)", data_ptr+debug_mem_start)
				module.Datas = append(module.Datas, dd)
				// Replace any refs to len / offset in the code
				length := dd.DataLength()

				constants[fmt.Sprintf("offset.%s", dd.Identifier)] = data_ptr
				constants[fmt.Sprintf("length.%s", dd.Identifier)] = length

				data_ptr += length
				// Align to 8
				data_ptr = (data_ptr + 7) & 0xfffffff8
			}
		}
	}

	// Replace some constants in the code
	for k, v := range constants {
		module.ReplaceConst(k, v)
	}

	// This signifies the top of used memory
	module.Globals = append(module.Globals, wasm.NewGlobal(fmt.Sprintf("(global $debug_memory_top i32 (i32.const %d))", data_ptr)))

	// Write out the new wat file
	fmt.Printf(";; #### MERGED ####\n%s", module.Write())

	// Write out the function names
	// TODO: Do this a better way
	for findex, f := range module.Funcs {
		fmt.Printf(";; functionNames[%d] = \"%s\";\n", findex, f.Identifier)
	}

}
