package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/loopholelabs/wasm-lnkr/wasm"
)

func main() {

	args := os.Args[1:]

	wat1 := args[0]
	wat2 := args[1]

	//fmt.Printf("Starting lnkr, merging [%s] with [%s]...\n\n", wat1, wat2)

	log.Print("Loading mod1...")

	mod1 := wasm.NewModule(wat1)
	mod1.Parse()

	log.Print("Loading mod2...")

	mod2 := wasm.NewModule(wat2)
	mod2.Parse()
	/*
		fmt.Printf("Memory = %d / %d\n", mod1.Memorys[0].Size, mod2.Memorys[0].Size)

		// TODO: Load the wat files, and merge them

		fmt.Printf("Import %s / %s = %s\n", mod1.Imports[0].Identifier1, mod1.Imports[0].Identifier2, mod1.Imports[0].Import)
		fmt.Printf("Export %s = %s\n", mod1.Exports[0].Identifier, mod1.Exports[0].Export)
		fmt.Printf("Data %s / %s / %s\n", mod1.Datas[0].Identifier, mod1.Datas[0].Location, mod1.Datas[0].Data)
		fmt.Printf("Global %s / %s / %s\n", mod1.Globals[0].Identifier, mod1.Globals[0].Type, mod1.Globals[0].Value)
		fmt.Printf("Elem %s / %s / %s\n", mod1.Elems[0].Identifier, mod1.Elems[0].Type, mod1.Elems[0].Elems[0])
		fmt.Printf("Table %s / %d / %d / %s\n", mod1.Tables[0].Identifier, mod1.Tables[0].Size, mod1.Tables[0].Limit, mod1.Tables[0].Type)

		for _, f := range mod1.Funcs {
			fmt.Printf(" Function %s (%d)\n", f.Identifier, len(f.Instructions))
		}

		fmt.Printf("=== FUNCTION ===\n%s\n", mod1.Funcs[0].Write())
	*/
	newmod := wasm.NewModule("combined")

	// Start merging...
	prefix1 := "mod1."
	prefix2 := "mod2."

	// Merge exports
	newmod.Exports = wasm.MergeExports(prefix1, mod1.Exports, prefix2, mod2.Exports)

	// Merge globals
	newmod.Globals = wasm.MergeGlobals(prefix1, mod1.Globals, prefix2, mod2.Globals)

	// Merge imports
	newmod.Imports = wasm.MergeImports(mod1.Imports, mod2.Imports, len(mod1.Types))

	// See if we're importing env:debug already
	need_debug_import := true

	for _, ii := range newmod.Imports {
		if ii.Identifier1 == "\"env\"" && ii.Identifier2 == "\"h_debug\"" {
			need_debug_import = false
		}
	}

	if need_debug_import {
		newmod.Imports = append(newmod.Imports, wasm.NewImport("(import \"env\" \"h_debug\" (func $h_debug (param i32)))"))
	}

	// Merge memory
	if len(mod1.Memorys) != 1 || len(mod2.Memorys) != 1 {
		panic("Modules can only have one memory")
	}
	newmod.Memorys = append(newmod.Memorys, wasm.NewMemory(fmt.Sprintf("(memory (;0;) %d)", mod1.Memorys[0].Size+mod2.Memorys[0].Size)))

	// Merge data
	offset2 := mod1.Memorys[0].Size * 65536

	newmod.Datas = wasm.MergeDatas(prefix1, mod1.Datas, offset2, prefix2, mod2.Datas)

	// Merge type
	newmod.Types = wasm.MergeTypes(prefix1, mod1.Types, prefix2, mod2.Types)

	// TODO: Merge these
	totalSize := 0
	totalLimit := 0
	m2tab_offset := 0
	if len(mod1.Tables) == 1 {
		totalSize += mod1.Tables[0].Size
		totalLimit += mod1.Tables[0].Limit
		m2tab_offset = mod1.Tables[0].Limit
	}

	if len(mod2.Tables) == 1 {
		totalSize += mod2.Tables[0].Size
		totalLimit += mod2.Tables[0].Limit
	}

	// Create a large table
	newmod.Tables = append(newmod.Tables, wasm.NewTable(fmt.Sprintf("(table (;0;) %d %d funcref)", totalSize, totalLimit)))

	for _, e := range mod1.Elems {
		for ind, ff := range e.Elems {
			e.Elems[ind] = "$" + prefix1 + ff[1:]
		}
		newmod.Elems = append(newmod.Elems, e)
	}

	// Now add mod2 elems, but relocated
	for _, e := range mod2.Elems {
		log.Printf("m2 Elem %s", e.Offset)
		if strings.HasPrefix(e.Offset, "(i32.const ") {
			v, err := strconv.Atoi(e.Offset[11 : len(e.Offset)-1])
			if err != nil {
				panic("Error translating elem")
			}
			e.Offset = fmt.Sprintf("(i32.const %d)", v+m2tab_offset)
		}
		for ind, ff := range e.Elems {
			e.Elems[ind] = "$" + prefix2 + ff[1:]
		}
		newmod.Elems = append(newmod.Elems, e)
	}

	//newmod.Elems = append(newmod.Elems, wasm.NewElem("(elem (;0;) (i32.const 1) func $mod1.runtime.memequal $mod1.runtime.hash32)"))

	// Merge table
	// Merge elem
	// TODO we need to create a new table with enough entries
	// TODO we need to merge the elem entries
	// Then when it comes to processing func, we will need to adjust any call_indirect in mod2

	// Merge func

	for _, f := range mod1.Funcs {
		// Rewrite the function itself
		f.Identifier = "$" + prefix1 + f.Identifier[1:]

		f.PrefixCalls(prefix1, mod1.Funcs)
		f.PrefixGlobals(prefix1)
		f.FixMemoryInstr("memory.size", "$m1_memory.size")
		f.FixMemoryInstr("memory.grow", "$m1_memory.grow")

		newmod.Funcs = append(newmod.Funcs, f)

	}

	for _, f := range mod2.Funcs {
		// Rewrite the function itself
		f.Identifier = "$" + prefix2 + f.Identifier[1:]

		// Update type index
		f.AdjustType(len(mod1.Types))

		f.FixCallIndirectType(len(mod1.Types))

		f.PrefixCalls(prefix2, mod2.Funcs)
		f.PrefixGlobals(prefix2)
		f.FixMemoryInstr("memory.size", "$m2_memory.size")
		f.FixMemoryInstr("memory.grow", "$m2_memory.grow")
		f.FixMemoryInstr("memory.copy", "$m2_memory.copy")
		f.FixMemoryInstr("memory.fill", "$m2_memory.fill")

		f.FixMemoryInstrOffsetAlign("i32.store", "$m2_i32.store")
		f.FixMemoryInstrOffsetAlign("i32.store8", "$m2_i32.store8")
		f.FixMemoryInstrOffsetAlign("i32.store16", "$m2_i32.store16")

		f.FixMemoryInstrOffsetAlign("i64.store", "$m2_i64.store")
		f.FixMemoryInstrOffsetAlign("i64.store8", "$m2_i64.store8")
		f.FixMemoryInstrOffsetAlign("i64.store16", "$m2_i64.store16")
		f.FixMemoryInstrOffsetAlign("i64.store32", "$m2_i64.store32")

		// TODO: Need to fixup call_indirect
		//		f.FixMemoryInstr("call_indirect", "$m2_call.indirect")

		f.AdjustOutcalls(newmod.Imports)

		f.AdjustLoad("$mm_offset")

		// TODO: Fix f32/f64 store

		newmod.Funcs = append(newmod.Funcs, f)
	}

	// Link any calls between modules
	for _, f := range mod1.Funcs {
		f.LinkCalls(mod2.Funcs)
	}

	// Link any calls between modules
	for _, f := range mod2.Funcs {
		f.LinkCalls(mod1.Funcs)
	}

	// Routines to deal with memory sizes
	newmod.Globals = append(newmod.Globals, wasm.NewGlobal(fmt.Sprintf("(global $mm_offset (mut i32) (i32.const %d))", offset2)))

	newmod.Globals = append(newmod.Globals, wasm.NewGlobal("(global $mm_current_mod (mut i32) (i32.const 0))"))

	newmod.Funcs = append(newmod.Funcs, wasm.NewFunc(`(func $m1_memory.size (result i32)
		global.get $mm_offset
		i32.const 16
		i32.shr_u
	)`))
	newmod.Funcs = append(newmod.Funcs, wasm.NewFunc(`(func $m2_memory.size (result i32)
		memory.size
		global.get $mm_offset
		i32.const 16
		i32.shr_u
		i32.sub
	)`))

	// TODO: memory.grow
	newmod.Funcs = append(newmod.Funcs, wasm.NewFunc(`(func $m1_memory.grow (param i32) (result i32)
		i32.const 1	
		call $h_debug
		unreachable
	)`))
	newmod.Funcs = append(newmod.Funcs, wasm.NewFunc(`(func $m2_memory.grow (param i32) (result i32)
		i32.const 2
		call $h_debug
		unreachable
	)`))

	newmod.Funcs = append(newmod.Funcs, wasm.NewFunc(`(func $m2_memory.copy (param i32 i32 i32)
		local.get 0		;; src
		global.get $mm_offset
		i32.add
		local.get 1		;; dst
		global.get $mm_offset
		i32.add
		local.get 2		;; size
		memory.copy
	)`))
	newmod.Funcs = append(newmod.Funcs, wasm.NewFunc(`(func $m2_memory.fill (param i32 i32 i32)
		local.get 0		;; src
		global.get $mm_offset
		i32.add
		local.get 1		;; value
		local.get 2		;; size
		memory.fill
	)`))
	newmod.Funcs = append(newmod.Funcs, wasm.NewFunc(`(func $m2_i32.store (param i32 i32 i32 i32)
		local.get 0
		global.get $mm_offset
		i32.add
		local.get 2
		i32.add

		local.get 3			;; align
		i32.const 1
		i32.sub
		i32.const -1
		i32.xor
		i32.and

		local.get 1
		i32.store
	)`))
	newmod.Funcs = append(newmod.Funcs, wasm.NewFunc(`(func $m2_i32.store8 (param i32 i32 i32 i32)
		local.get 0
		global.get $mm_offset
		i32.add
		local.get 2
		i32.add

		local.get 3			;; align
		i32.const 1
		i32.sub
		i32.const -1
		i32.xor
		i32.and

		local.get 1
		i32.store8
	)`))
	newmod.Funcs = append(newmod.Funcs, wasm.NewFunc(`(func $m2_i32.store16 (param i32 i32 i32 i32)
		local.get 0
		global.get $mm_offset
		i32.add
		local.get 2
		i32.add

		local.get 3			;; align
		i32.const 1
		i32.sub
		i32.const -1
		i32.xor
		i32.and

		local.get 1
		i32.store16
	)`))

	newmod.Funcs = append(newmod.Funcs, wasm.NewFunc(`(func $m2_i64.store (param i32 i64 i32 i32)
		local.get 0
		global.get $mm_offset
		i32.add
		local.get 2
		i32.add

		local.get 3			;; align
		i32.const 1
		i32.sub
		i32.const -1
		i32.xor
		i32.and

		local.get 1
		i64.store
	)`))
	newmod.Funcs = append(newmod.Funcs, wasm.NewFunc(`(func $m2_i64.store8 (param i32 i64 i32 i32)
		local.get 0
		global.get $mm_offset
		i32.add
		local.get 2
		i32.add

		local.get 3			;; align
		i32.const 1
		i32.sub
		i32.const -1
		i32.xor
		i32.and

		local.get 1
		i64.store8
	)`))
	newmod.Funcs = append(newmod.Funcs, wasm.NewFunc(`(func $m2_i64.store16 (param i32 i64 i32 i32)
		local.get 0
		global.get $mm_offset
		i32.add
		local.get 2
		i32.add

		local.get 3			;; align
		i32.const 1
		i32.sub
		i32.const -1
		i32.xor
		i32.and

		local.get 1
		i64.store16
	)`))
	newmod.Funcs = append(newmod.Funcs, wasm.NewFunc(`(func $m2_i64.store32 (param i32 i64 i32 i32)
		local.get 0
		global.get $mm_offset
		i32.add
		local.get 2			;; offset
		i32.add

		local.get 3			;; align
		i32.const 1
		i32.sub
		i32.const -1
		i32.xor
		i32.and

		local.get 1
		i64.store32
	)`))

	newmod.Funcs = append(newmod.Funcs, wasm.NewFunc(`(func $m2_call.indirect (param i32)
		i32.const 7788
		call $h_debug
	)`))

	newmod.Funcs = append(newmod.Funcs, wasm.NewFunc(`(func $m2_swapout
    (local $current_ptr1 i32)
    (local $current_ptr2 i32)
    (local $current_val1 i32)
    (local $current_val2 i32)

		global.get $mm_current_mod
		i32.const 0
		i32.eq
		br_if 0					;; if it's currently 0, then return

    global.get $mm_offset
    local.set $current_ptr2

    i32.const 0
    local.set $current_ptr1

    loop $swap_mem
      local.get $current_ptr1
      i32.load
      local.set $current_val1

      local.get $current_ptr2
      i32.load
      local.set $current_val2

      local.get $current_ptr1
      local.get $current_val2
      i32.store

      local.get $current_ptr2
      local.get $current_val1
      i32.store      

      ;; Now advance pointers
      local.get $current_ptr1
      i32.const 4
      i32.add
      local.set $current_ptr1

      local.get $current_ptr2
      i32.const 4
      i32.add
      local.set $current_ptr2

      local.get $current_ptr1
      global.get $mm_offset
      i32.lt_u
      br_if $swap_mem
    end 

		i32.const 0
		global.set $mm_current_mod

	)`))

	newmod.Funcs = append(newmod.Funcs, wasm.NewFunc(`(func $m2_swapin
    (local $current_ptr1 i32)
    (local $current_ptr2 i32)
    (local $current_val1 i32)
    (local $current_val2 i32)

		global.get $mm_current_mod
		i32.const 1
		i32.eq
		br_if 0					;; if it's currently 1, then return

    global.get $mm_offset
    local.set $current_ptr2

    i32.const 0
    local.set $current_ptr1

    loop $swap_mem
      local.get $current_ptr1
      i32.load
      local.set $current_val1

      local.get $current_ptr2
      i32.load
      local.set $current_val2

      local.get $current_ptr1
      local.get $current_val2
      i32.store

      local.get $current_ptr2
      local.get $current_val1
      i32.store      

      ;; Now advance pointers
      local.get $current_ptr1
      i32.const 4
      i32.add
      local.set $current_ptr1

      local.get $current_ptr2
      i32.const 4
      i32.add
      local.set $current_ptr2

      local.get $current_ptr1
      global.get $mm_offset
      i32.lt_u
      br_if $swap_mem
    end 

		i32.const 1
		global.set $mm_current_mod

	)`))

	fmt.Printf(";; #### MERGED ####\n%s", newmod.Write())
}
