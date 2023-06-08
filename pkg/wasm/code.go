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

package wasmfile

import "github.com/loopholelabs/wasm-toolkit/pkg/wasm/types"

/**
 * This will redirect all calls for an imported function to another function.
 * The imported function will be removed.
 * The target function must already exist, with a name.
 */
func (wf *WasmFile) RedirectImport(fromModule string, from string, to string) {

	fid := wf.LookupFunctionID(to)

	if fid == -1 {
		panic("Target function not found!")
	}

	remap := map[int]int{}
	remap_imports := map[int]int{}

	// Now we need to REMOVE the old imports.
	newImports := make([]*ImportEntry, 0)
	for n, i := range wf.Import {
		if i.Module == fromModule && i.Name == from {
			remap_imports[n] = fid
		} else {
			remap[n] = len(newImports)
			newImports = append(newImports, i)
		}
	}

	// Remap everything in the Code section because we're removing an import.
	for n := range wf.Code {
		remap[len(wf.Import)+n] = len(newImports) + n
	}

	// Remap the import, and THEN remap due to removing the imports.
	for _, c := range wf.Code {
		c.ModifyAllCalls(remap_imports)
		c.ModifyAllCalls(remap)
	}

	wf.Import = newImports

	// We also need to fixup any Elems sections
	for _, el := range wf.Elem {
		for idx, funcidx := range el.Indexes {
			newidx, ok := remap[int(funcidx)]
			if ok {
				el.Indexes[idx] = uint64(newidx)
			}
		}
	}

	// Fixup exports
	for _, ex := range wf.Export {
		if ex.Type == types.ExportFunc {
			newidx, ok := remap[ex.Index]
			if ok {
				ex.Index = newidx
			}
		}
	}

	wf.Debug.RenumberFunctions(remap)
}
