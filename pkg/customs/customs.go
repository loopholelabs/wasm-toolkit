package customs

import (
	"fmt"

	"github.com/loopholelabs/wasm-toolkit/pkg/wasm/expression"
	"github.com/loopholelabs/wasm-toolkit/pkg/wasm/types"
	"github.com/loopholelabs/wasm-toolkit/pkg/wasm/wasmfile"
)

type Import struct {
	Module string
	Name   string
}

type RemapMuxImport struct {
	Source Import
	Mapper map[uint64]Import
}

type RemapMuxExport struct {
	Source string
	Mapper map[uint64]string
}

/**
 *
 */
func MuxExport(wfile *wasmfile.WasmFile, c RemapMuxExport) {
	// For each item in the map we need to add a new function, and an export.

	var sourceType *wasmfile.TypeEntry = nil
	var sourceTypeId int = 0
	var sourceFunctionId int = 0
	newExports := make([]*wasmfile.ExportEntry, 0)
	for _, e := range wfile.Export {
		if e.Type == types.ExportFunc && e.Name == c.Source {
			// We found our source
			// Now look up the function
			sourceFunctionId = e.Index
			p := e.Index - len(wfile.Import)
			sourceTypeId = wfile.Function[p].TypeIndex
			sourceType = wfile.Type[sourceTypeId]
		} else {
			newExports = append(newExports, e)
		}
	}
	if sourceType == nil {
		panic("Invalid wasm file")
	}

	// Now create a new type without the mux uint64 id
	newType := sourceType.Clone()
	// Remove the uint64 ID
	newType.Param = newType.Param[1:]

	newTypeId := wfile.AddTypeMaybe(newType)

	for id, nname := range c.Mapper {
		// First create a new function to use...
		wfile.Function = append(wfile.Function, &wasmfile.FunctionEntry{
			TypeIndex: newTypeId,
		})

		exp := make([]*expression.Expression, 0)

		exp = append(exp, &expression.Expression{
			Opcode:   expression.InstrToOpcode["i64.const"],
			I64Value: int64(id),
		})
		// Now load the params...

		// Push on the arguments
		for pid := range newType.Param {
			exp = append(exp,
				&expression.Expression{
					Opcode:     expression.InstrToOpcode["local.get"],
					LocalIndex: pid,
				},
			)
		}
		// Do the call...
		// NB the return type is the same
		exp = append(exp,
			&expression.Expression{
				Opcode:    expression.InstrToOpcode["call"],
				FuncIndex: sourceFunctionId,
			},
		)

		cod := &wasmfile.CodeEntry{
			Locals:         []types.ValType{},
			PCValid:        false,
			CodeSectionPtr: 0,
			CodeSectionLen: 0,
			Expression:     exp,
		}

		wfile.Code = append(wfile.Code, cod)

		newExports = append(newExports, &wasmfile.ExportEntry{
			Type:  types.ExportFunc,
			Name:  nname,
			Index: 0, // TODO: Point to the new function we just created
		})
	}

	// Update the exports...
	wfile.Export = newExports
}

/**
 * This will do the following
 * - Add new imports
 * - Remove an existing import
 * - Reroute calls to the existing import to the new imports, using the first param(uint64) as index
 *
 * So for example, given a mux map of {1: "callImportONE", 2: "callImportTWO"}
 *   callImport(1,...); callImport(2,...);
 *
 * would get replaced with
 *   callImportONE(...); callImportTWO(...);
 */
func MuxImport(wfile *wasmfile.WasmFile, c RemapMuxImport) {

	// Make sure there's a type
	var sourceType *wasmfile.TypeEntry = nil
	var sourceTypeId int = 0
	for _, i := range wfile.Import {
		if i.Module == c.Source.Module && i.Name == c.Source.Name {
			sourceTypeId = i.Index
			sourceType = wfile.Type[i.Index]
		}
	}
	if sourceType == nil {
		panic("Source type not found")
	}
	// Now create a new type without the mux uint64 id
	newType := sourceType.Clone()
	// Remove the uint64 ID
	newType.Param = newType.Param[1:]

	tid := wfile.AddTypeMaybe(newType)

	remap := map[int]int{}

	newImports := make([]*wasmfile.ImportEntry, 0)

	sourceId := 0

	for id, i := range wfile.Import {
		if i.Module == c.Source.Module && i.Name == c.Source.Name {
			// Remove
			remap[id] = 0 // TODO: Remap to the new mux function
			sourceId = id // We'll update this in a bit
		} else {
			newImports = append(newImports, i)
			remap[id] = len(newImports)
		}
	}

	// Add the new imports, and keep track of which ID each is for the demux function.
	newFunctions := make(map[int]Import)
	for _, i := range c.Mapper {
		newFunctions[len(newImports)] = i
		newImports = append(newImports, &wasmfile.ImportEntry{
			Module: i.Module,
			Name:   i.Name,
			Type:   types.ExportFunc,
			Index:  tid,
		})
	}

	for n := range wfile.Code {
		remap[len(wfile.Import)+n] = len(newImports) + n
	}

	wfile.Import = newImports

	// Adjust to our new function (Added soon)
	remap[sourceId] = len(wfile.Import) + len(wfile.Code)

	for _, c := range wfile.Code {
		c.ModifyAllCalls(remap)
	}

	wfile.Debug.RenumberFunctions(remap)

	// Add debug name for the new imports
	for id, ii := range newFunctions {
		wfile.Debug.FunctionNames[id] = fmt.Sprintf("mux_%s_%s", ii.Module, ii.Name)
	}

	exp := make([]*expression.Expression, 0)

	// Write the code here to call the correct function

	for mid, re := range c.Mapper {

		exp = append(exp,
			&expression.Expression{
				Opcode: expression.InstrToOpcode["block"],
				Result: types.ValNone,
			},
			&expression.Expression{
				Opcode:     expression.InstrToOpcode["local.get"],
				LocalIndex: 0,
			},
			&expression.Expression{
				Opcode:   expression.InstrToOpcode["i64.const"],
				I64Value: int64(mid),
			},
			&expression.Expression{
				Opcode: expression.InstrToOpcode["i64.ne"],
			},
			&expression.Expression{
				Opcode:     expression.InstrToOpcode["br_if"],
				LabelIndex: 0,
			},
		)

		// Push on the arguments
		for pid := range newType.Param {
			exp = append(exp,
				&expression.Expression{
					Opcode:     expression.InstrToOpcode["local.get"],
					LocalIndex: pid,
				},
			)
		}
		// Do the call...
		exp = append(exp,
			&expression.Expression{
				Opcode:               expression.InstrToOpcode["call"],
				FunctionId:           fmt.Sprintf("mux_%s_%s", re.Module, re.Name),
				FunctionNeedsLinking: true,
			},
			&expression.Expression{
				Opcode: expression.InstrToOpcode["return"],
			},
		)

		exp = append(exp,
			&expression.Expression{
				Opcode:     expression.InstrToOpcode["end"],
				LocalIndex: 0,
			},
		)

	}

	// Default fallthrough value
	exp = append(exp, &expression.Expression{
		Opcode: expression.InstrToOpcode["unreachable"],
	},
		&expression.Expression{
			Opcode:   expression.InstrToOpcode["i64.const"],
			I64Value: 123,
		},
	)
	// Now add the new demux function
	wfile.Function = append(wfile.Function, &wasmfile.FunctionEntry{
		TypeIndex: sourceTypeId,
	})

	cod := &wasmfile.CodeEntry{
		Locals:         []types.ValType{},
		PCValid:        false,
		CodeSectionPtr: 0,
		CodeSectionLen: 0,
		Expression:     exp,
	}

	// Resolve the functions
	err := cod.ResolveFunctions(wfile)
	if err != nil {
		panic(err)
	}
	wfile.Code = append(wfile.Code, cod)

	// Remap exports
	for _, ee := range wfile.Export {
		if ee.Type == types.ExportFunc {
			newid, ok := remap[ee.Index]
			if !ok {
				panic("Remap failed for some odd reason")
			}
			ee.Index = newid
		}
	}

}
