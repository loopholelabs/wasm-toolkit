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

import (
	"debug/dwarf"
	"fmt"
	"io"

	"github.com/loopholelabs/wasm-toolkit/pkg/wasm/debug"
)

func (wf *WasmFile) ParseDwarf() error {
	debug_abbrev := wf.GetCustomSectionData(".debug_abbrev")
	debug_aranges := wf.GetCustomSectionData(".debug_aranges")
	debug_info := wf.GetCustomSectionData(".debug_info")
	debug_line := wf.GetCustomSectionData(".debug_line")
	debug_pubnames := wf.GetCustomSectionData(".debug_pubnames")
	debug_ranges := wf.GetCustomSectionData(".debug_ranges")
	debug_str := wf.GetCustomSectionData(".debug_str")

	debug_loc := wf.GetCustomSectionData(".debug_loc")
	wf.Debug.DwarfLoc = debug.NewDwarfLocations(debug_loc)

	debug_frame := make([]byte, 0) // call frame info

	dd, err := dwarf.New(debug_abbrev, debug_aranges, debug_frame, debug_info, debug_line, debug_pubnames, debug_ranges, debug_str)
	if err != nil {
		return nil // ok, but lets move on and ignore the error.
	}

	wf.Debug.DwarfData = dd
	return nil
}

func (wf *WasmFile) GetLocalVarName(pc uint64, index int) string {
	for _, lnd := range wf.Debug.LocalNames {
		if lnd.Index == index && (pc >= lnd.StartPC && pc <= lnd.EndPC) {
			return lnd.VarName
		}
	}
	return ""
}

func (wf *WasmFile) GetLocalVarType(pc uint64, index int) string {
	for _, lnd := range wf.Debug.LocalNames {
		if lnd.Index == index && (pc >= lnd.StartPC && pc <= lnd.EndPC) {
			return lnd.VarType
		}
	}
	return ""
}

func (wf *WasmFile) GetFunctionDebug(fid int) string {
	de, ok := wf.Debug.FunctionDebug[fid]
	if ok {
		return de
	}
	return ""
}

func (wf *WasmFile) SetFunctionSignature(fid int, de string) {
	if wf.Debug.FunctionSignature == nil {
		wf.Debug.FunctionSignature = make(map[int]string)
	}
	wf.Debug.FunctionSignature[fid] = de
}

func (wf *WasmFile) GetFunctionSignature(fid int) string {
	de, ok := wf.Debug.FunctionSignature[fid]
	if ok {
		return de
	}
	return ""
}

func (wf *WasmFile) ParseDwarfVariables() error {
	wf.Debug.FunctionDebug = make(map[int]string)
	if wf.Debug.FunctionSignature == nil {
		wf.Debug.FunctionSignature = make(map[int]string)
	}
	wf.Debug.LocalNames = make([]*debug.LocalNameData, 0)

	wf.Debug.GlobalAddresses = make(map[string]*debug.GlobalNameData)

	if wf.Debug.DwarfData == nil {
		return nil
	}

	entryReader := wf.Debug.DwarfData.Reader()

	for {
		// Read all entries in sequence
		entry, err := entryReader.Next()
		if entry == nil || err == io.EOF {
			// We've reached the end of DWARF entries
			break
		}

		if entry.Tag == dwarf.TagVariable {
			// Parse the location address
			vname := ""
			var vaddr []byte
			vsize := int64(0)
			vtype := ""
			for _, field := range entry.Field {
				if field.Attr == dwarf.AttrName {
					vname = field.Val.(string)
				} else if field.Attr == dwarf.AttrLocation {
					// Parse the expression
					switch field.Val.(type) {
					case []byte:
						vaddr = field.Val.([]byte)
					}
				} else if field.Attr == dwarf.AttrType {
					offset := field.Val.(dwarf.Offset)
					ty, err := wf.Debug.DwarfData.Type(offset)
					if err == nil {
						vsize = ty.Size()
						vtype = ty.String()
					}
				}
			}

			if vaddr != nil && vname != "" {
				// Parse the expression
				ld := &debug.LocationData{
					Expression: vaddr,
				}
				addr, err := ld.GetAddress()
				if err == nil {

					globalInfo := &debug.GlobalNameData{
						Name:    vname,
						Address: uint64(addr),
						Size:    uint64(vsize),
						Type:    vtype,
					}
					wf.Debug.GlobalAddresses[vname] = globalInfo
				} else {
					// TODO
					// fmt.Printf("Variable but not simple expr... %s %x\n", vname, vaddr)
				}
			}
		}

		if entry.Tag == dwarf.TagSubprogram {
			spname := "<unknown>"
			sploc := uint64(0)
			for _, field := range entry.Field {
				//				log.Printf("Field %v\n", field)
				if field.Attr == dwarf.AttrName {
					spname = field.Val.(string)
				} else if field.Attr == dwarf.AttrLowpc {
					switch field.Val.(type) {
					case uint64:
						sploc = field.Val.(uint64)
					}
				}
			}

			params := ""
			locals := ""
			if entry.Children {
				// Read the children...
				for {
					entry, err := entryReader.Next()
					if err != nil {
						return err
					}
					if entry.Tag == 0 {
						break
					}

					vname := "<unknown>"
					vtype := ""
					vloc := int64(-1)
					for _, field := range entry.Field {
						if field.Attr == dwarf.AttrName {
							vname = field.Val.(string)
						} else if field.Attr == dwarf.AttrType {
							switch field.Val.(type) {
							case dwarf.Offset:
								t := field.Val.(dwarf.Offset)
								ty, err := wf.Debug.DwarfData.Type(t)
								if err == nil {
									vtype = ty.String()
								}
							}
						} else if field.Attr == dwarf.AttrLocation {
							switch field.Val.(type) {
							case int64:
								vloc = field.Val.(int64)
							}
						}
					}

					fmt.Printf("DwarfEntry tag=%v vname=%s entry=%v\n", entry.Tag, vname, entry)

					if entry.Tag == dwarf.TagFormalParameter {
						if vloc != -1 {
							locdata := wf.Debug.DwarfLoc.ReadLocation(uint64(vloc))
							for _, ld := range locdata {
								// We have code ptr range here...

								locs := ld.ExtractWasmLocations()
								for _, l := range locs {
									if l.IsLocal {
										// Store in the locals lookup...
										wf.Debug.LocalNames = append(wf.Debug.LocalNames, &debug.LocalNameData{
											StartPC: uint64(sploc), //ld.startAddress),
											EndPC:   uint64(sploc), //ld.endAddress),
											Index:   int(l.Index),
											VarName: vname,
											VarType: vtype,
										})

										fmt.Printf("LocationLocal %s %s (%d-%d) %d local %d\n", spname, vname, ld.StartAddress, ld.EndAddress, sploc, l.Index)
									}
								}
							}
						}
						if len(params) > 0 {
							params = params + ", "
						}
						params = fmt.Sprintf("%s%s(%s)", params, vname, vtype)
					} else if entry.Tag == dwarf.TagVariable {

						//						fmt.Printf("TagVariable %v\n", entry)

						if vloc != -1 {
							locdata := wf.Debug.DwarfLoc.ReadLocation(uint64(vloc))
							for _, ld := range locdata {
								// We have code ptr range here...

								fmt.Printf("Var Data %s is %d %d %x\n", vname, ld.StartAddress, ld.EndAddress, ld.Expression)

								locs := ld.ExtractWasmLocations()
								for _, l := range locs {
									if l.IsLocal {
										// Store in the locals lookup...
										wf.Debug.LocalNames = append(wf.Debug.LocalNames, &debug.LocalNameData{
											StartPC: uint64(ld.StartAddress),
											EndPC:   uint64(ld.EndAddress),
											Index:   int(l.Index),
											VarName: vname,
										})

										fmt.Printf("LocationLocalVariable %s %s %d-%d  local %d\n", spname, vname, ld.StartAddress, ld.EndAddress, l.Index)
									}
								}
							}
						}
						locals = fmt.Sprintf("%s;; local %s %s\n", locals, vname, vtype)
					}

				}
			}

			function_debug := fmt.Sprintf(";; %s(%s)\n%s", spname, params, locals)

			fid := wf.FindFunction(sploc)

			if fid != -1 {
				wf.Debug.FunctionSignature[fid] = fmt.Sprintf("%s(%s)", spname, params)
				wf.Debug.FunctionDebug[fid] = function_debug
			}
		}
	}
	return nil
}
