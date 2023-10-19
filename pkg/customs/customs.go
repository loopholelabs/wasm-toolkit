package customs

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

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

const DEF_SEPARATOR = ","
const ID_SEPARATOR = ":"
const MODULE_SEPARATOR = "/"

func ParseRemapMuxExport(e string) (*RemapMuxExport, error) {
	// eg "hello,0:zero,1:one,2:two"

	eBits := strings.Split(e, DEF_SEPARATOR)
	eDef := &RemapMuxExport{Mapper: map[uint64]string{}}
	for i, v := range eBits {
		if i == 0 {
			eDef.Source = v
		} else {
			// Split on :
			vals := strings.Split(v, ID_SEPARATOR)
			if len(vals) != 2 {
				return nil, errors.New("Invalid definition")
			}
			id, err := strconv.ParseUint(vals[0], 10, 64)
			if err != nil {
				return nil, errors.New("Invalid definition")
			}
			eDef.Mapper[id] = vals[1]
		}
	}
	return eDef, nil
}

func ParseRemapMuxImport(e string) (*RemapMuxImport, error) {
	// eg "env/hello,0:env/zero,1:env/one,2:env/two"

	eBits := strings.Split(e, DEF_SEPARATOR)
	eDef := &RemapMuxImport{Mapper: map[uint64]Import{}}
	for i, v := range eBits {
		if i == 0 {
			b := strings.Split(v, MODULE_SEPARATOR)
			if len(b) != 2 {
				return nil, errors.New("Invalid definition")
			}
			eDef.Source = Import{
				Module: b[0],
				Name:   b[1],
			}
		} else {
			// Split on :
			vals := strings.Split(v, ID_SEPARATOR)
			if len(vals) != 2 {
				return nil, errors.New("Invalid definition")
			}
			id, err := strconv.ParseUint(vals[0], 10, 64)
			if err != nil {
				return nil, errors.New("Invalid definition")
			}
			b := strings.Split(vals[1], MODULE_SEPARATOR)
			if len(b) != 2 {
				return nil, errors.New("Invalid definition")
			}
			eDef.Mapper[id] = Import{
				Module: b[0],
				Name:   b[1],
			}
		}
	}
	return eDef, nil
}

/**
 *
 */
func MuxExport(wfile *wasmfile.WasmFile, c RemapMuxExport) error {
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
			if p >= len(wfile.Function) {
				return errors.New("Invalid wasm file")
			}
			sourceTypeId = wfile.Function[p].TypeIndex
			if sourceTypeId >= len(wfile.Type) {
				return errors.New("Invalid wasm file")
			}
			sourceType = wfile.Type[sourceTypeId]
		} else {
			newExports = append(newExports, e)
		}
	}
	if sourceType == nil {
		return errors.New("Source function not found")
	}

	// Now create a new type without the mux uint64 id
	newType := sourceType.Clone()
	if len(newType.Param) < 1 || newType.Param[0] != types.ValI64 {
		return errors.New("First param of mux must be i64")
	}
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

		newFunctionId := len(wfile.Import) + len(wfile.Code)
		wfile.Code = append(wfile.Code, cod)

		newExports = append(newExports, &wasmfile.ExportEntry{
			Type:  types.ExportFunc,
			Name:  nname,
			Index: newFunctionId,
		})
	}

	// Update the exports...
	wfile.Export = newExports
	return nil
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
func MuxImport(wfile *wasmfile.WasmFile, c RemapMuxImport) error {

	// Make sure there's a type
	var sourceType *wasmfile.TypeEntry = nil
	var sourceTypeId int = 0
	for _, i := range wfile.Import {
		if i.Module == c.Source.Module && i.Name == c.Source.Name {
			sourceTypeId = i.Index
			if i.Index >= len(wfile.Type) {
				return errors.New("Invalid wasm file")
			}
			sourceType = wfile.Type[i.Index]
		}
	}
	if sourceType == nil {
		return errors.New("Source function not found")
	}
	// Now create a new type without the mux uint64 id
	newType := sourceType.Clone()
	if len(newType.Param) < 1 || newType.Param[0] != types.ValI64 {
		return errors.New("First param of mux must be i64")
	}
	// Remove the uint64 ID
	newType.Param = newType.Param[1:]

	tid := wfile.AddTypeMaybe(newType)

	remap := map[int]int{}

	newImports := make([]*wasmfile.ImportEntry, 0)

	sourceId := 0

	for id, i := range wfile.Import {
		if i.Module == c.Source.Module && i.Name == c.Source.Name {
			// Remove
			remap[id] = 0 // This gets set later
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
		return err
	}
	wfile.Code = append(wfile.Code, cod)

	// Remap exports
	for _, ee := range wfile.Export {
		if ee.Type == types.ExportFunc {
			newid, ok := remap[ee.Index]
			if !ok {
				return errors.New("Invalid wasm file")
			}
			ee.Index = newid
		}
	}
	return nil
}
