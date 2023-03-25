package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/loopholelabs/wasm-lnkr/wasm"
)

func main() {

	args := os.Args[1:]

	// Just a single arg for now, the input filename...

	watfile := args[0]

	log.Print("Loading mod...")
	module := wasm.NewModule(watfile)
	module.Parse()

	log.Print("Loading debug.wat...")
	debugModule := wasm.NewModule("debug.wat")
	debugModule.Parse()

	// Add some extra imports that we'll use to signal to the host about events.
	// debug_sf - Start function
	// debug_sfp_* - Parameter details
	// debug_sfp - Parameters finished
	// debug_ef - End function
	// debug_ef_* - End function with return value
	module.Imports = append(module.Imports, wasm.NewImport("(import \"env\" \"debug_sf\" (func $debug_sf (param i32)))"))
	module.Imports = append(module.Imports, wasm.NewImport("(import \"env\" \"debug_sfp_i32\" (func $debug_sfp_i32 (param i32 i32)))"))
	module.Imports = append(module.Imports, wasm.NewImport("(import \"env\" \"debug_sfp_i64\" (func $debug_sfp_i64 (param i32 i64)))"))
	module.Imports = append(module.Imports, wasm.NewImport("(import \"env\" \"debug_sfp_f32\" (func $debug_sfp_f32 (param i32 f32)))"))
	module.Imports = append(module.Imports, wasm.NewImport("(import \"env\" \"debug_sfp_f64\" (func $debug_sfp_f64 (param i32 f64)))"))
	module.Imports = append(module.Imports, wasm.NewImport("(import \"env\" \"debug_sfp\" (func $debug_sfp (param i32)))"))

	module.Imports = append(module.Imports, wasm.NewImport("(import \"env\" \"debug_ef\" (func $debug_ef (param i32)))"))
	module.Imports = append(module.Imports, wasm.NewImport("(import \"env\" \"debug_ef_i32\" (func $debug_ef_i32 (param i32 i32) (result i32)))"))
	module.Imports = append(module.Imports, wasm.NewImport("(import \"env\" \"debug_ef_i64\" (func $debug_ef_i64 (param i64 i32) (result i64)))"))
	module.Imports = append(module.Imports, wasm.NewImport("(import \"env\" \"debug_ef_f32\" (func $debug_ef_f32 (param f32 i32) (result f32)))"))
	module.Imports = append(module.Imports, wasm.NewImport("(import \"env\" \"debug_ef_f64\" (func $debug_ef_f64 (param f64 i32) (result f64)))"))

	// Go through each function, and add some debugging to it.
	for findex, f := range module.Funcs {
		ins := make([]string, 0)

		// Add a header to the function.
		ins = append(ins, fmt.Sprintf("block %s", f.Result))
		ins = append(ins, fmt.Sprintf("i32.const %d", findex))
		ins = append(ins, "call $debug_sf")

		// Go through the function parameters, and send them to the host
		pindex := 0
		for _, pa := range f.Params {
			if strings.HasPrefix(pa, "(param ") {
				bits := pa[7 : len(pa)-1]
				words := strings.Fields(bits)

				for _, w := range words {
					if w == "i32" {
						ins = append(ins, fmt.Sprintf("i32.const %d", findex))
						ins = append(ins, fmt.Sprintf("local.get %d", pindex))
						ins = append(ins, "call $debug_sfp_i32")
					} else if w == "i64" {
						ins = append(ins, fmt.Sprintf("i32.const %d", findex))
						ins = append(ins, fmt.Sprintf("local.get %d", pindex))
						ins = append(ins, "call $debug_sfp_i64")
					} else if w == "f32" {
						ins = append(ins, fmt.Sprintf("i32.const %d", findex))
						ins = append(ins, fmt.Sprintf("local.get %d", pindex))
						ins = append(ins, "call $debug_sfp_f32")
					} else if w == "f64" {
						ins = append(ins, fmt.Sprintf("i32.const %d", findex))
						ins = append(ins, fmt.Sprintf("local.get %d", pindex))
						ins = append(ins, "call $debug_sfp_f64")
					}
					pindex++
				}
			}
		}

		// Parameters are finished
		ins = append(ins, fmt.Sprintf("i32.const %d", findex))
		ins = append(ins, "call $debug_sfp")

		// Now insert function footers anywhere the function will return.
		for _, i := range f.Instructions {
			if strings.Trim(i, wasm.Whitespace) == "return" {
				// Insert a call before the return
				ins = append(ins, fmt.Sprintf("i32.const %d", findex))

				if f.Result == "(result i32)" {
					ins = append(ins, "call $debug_ef_i32")
				} else if f.Result == "(result i64)" {
					ins = append(ins, "call $debug_ef_i64")
				} else if f.Result == "(result f32)" {
					ins = append(ins, "call $debug_ef_f32")
				} else if f.Result == "(result f64)" {
					ins = append(ins, "call $debug_ef_f64")
				} else {
					ins = append(ins, "call $debug_ef")
				}
			}
			ins = append(ins, i)
		}

		// Add a footer at the end of the function. NB we add a block/end to prevent early function return.
		ins = append(ins, "end")
		ins = append(ins, fmt.Sprintf("i32.const %d", findex))
		if f.Result == "(result i32)" {
			ins = append(ins, "call $debug_ef_i32")
		} else if f.Result == "(result i64)" {
			ins = append(ins, "call $debug_ef_i64")
		} else if f.Result == "(result f32)" {
			ins = append(ins, "call $debug_ef_f32")
		} else if f.Result == "(result f64)" {
			ins = append(ins, "call $debug_ef_f64")
		} else {
			ins = append(ins, "call $debug_ef")
		}

		f.Instructions = ins
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

	// Add the debug functions
	for _, f := range debugModule.Funcs {
		module.Funcs = append(module.Funcs, f)
	}

	// Add the debug globals
	for _, g := range debugModule.Globals {
		module.Globals = append(module.Globals, g)
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
	for _, f := range module.Funcs {
		// First add the address and length to our table
		bs := make([]byte, 4)
		binary.LittleEndian.PutUint32(bs, uint32(len(functionNameData)))
		functionNameTable = append(functionNameTable, bs...)
		binary.LittleEndian.PutUint32(bs, uint32(len(f.Identifier)))
		functionNameTable = append(functionNameTable, bs...)

		// Now add the string onto our data
		functionNameData = append(functionNameData, []byte(f.Identifier)...)
	}

	// Now we need to write these as data items...
	allowed := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890 "

	debug_function_table_ptr := debug_mem_start
	debug_function_table := ""
	for _, v := range functionNameTable {
		if strings.Index(allowed, string(rune(v))) == -1 {
			debug_function_table = debug_function_table + fmt.Sprintf("\\%02x", v)
		} else {
			debug_function_table = debug_function_table + string(rune(v))
		}
	}

	debug_function_data_ptr := debug_mem_start + len(functionNameTable)
	debug_function_data := ""
	for _, v := range functionNameData {
		if strings.Index(allowed, string(rune(v))) == -1 {
			debug_function_data = debug_function_data + fmt.Sprintf("\\%02x", v)
		} else {
			debug_function_data = debug_function_data + string(rune(v))
		}
	}

	module.Datas = append(module.Datas, wasm.NewData(fmt.Sprintf("(data $debug_function_table (i32.const %d) \"%s\")", debug_function_table_ptr, debug_function_table)))
	module.Datas = append(module.Datas, wasm.NewData(fmt.Sprintf("(data $debug_function_data (i32.const %d) \"%s\")", debug_function_data_ptr, debug_function_data)))

	// Write out the new wat file
	fmt.Printf(";; #### MERGED ####\n%s", module.Write())

	// Write out the function names
	// TODO: Do this a better way
	for findex, f := range module.Funcs {
		fmt.Printf(";; functionNames[%d] = \"%s\";\n", findex, f.Identifier)
	}

}
