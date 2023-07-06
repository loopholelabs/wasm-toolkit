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

package debug

import (
	"debug/dwarf"
	"fmt"
	"io"
	"strings"
)

func (wd *WasmDebug) GetLocalVarName(pc uint64, index int) string {
	for _, lnd := range wd.LocalNames {
		if lnd.Index == index && (pc >= lnd.StartPC && pc <= lnd.EndPC) {
			return lnd.VarName
		}
	}
	return ""
}

func (wd *WasmDebug) GetLocalVarType(pc uint64, index int) string {
	for _, lnd := range wd.LocalNames {
		if lnd.Index == index && (pc >= lnd.StartPC && pc <= lnd.EndPC) {
			return lnd.VarType
		}
	}
	return ""
}

func (wd *WasmDebug) GetFunctionDebug(fid int) string {
	de, ok := wd.FunctionDebug[fid]
	if ok {
		return de
	}
	return ""
}

func (wd *WasmDebug) SetFunctionSignature(fid int, de string) {
	if wd.FunctionSignature == nil {
		wd.FunctionSignature = make(map[int]string)
	}
	wd.FunctionSignature[fid] = de
}

func (wd *WasmDebug) GetFunctionSignature(fid int) string {
	de, ok := wd.FunctionSignature[fid]
	if ok {
		return de
	}
	return ""
}

type FunctionFinder interface {
	FindFunction(uint64) int
}

func (wd *WasmDebug) ParseDwarfGlobals() {
	wd.GlobalAddresses = make(map[string]*GlobalNameData)

	if wd.DwarfData == nil {
		return
	}

	entryReader := wd.DwarfData.Reader()

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
					ty, err := wd.DwarfData.Type(offset)
					if err == nil {
						vsize = ty.Size()
						vtype = ty.String()
					}
				}
			}

			if vaddr != nil && vname != "" {
				// Parse the expression
				ld := &LocationData{
					Expression: vaddr,
				}
				addr, err := ld.GetAddress()
				if err == nil {

					globalInfo := &GlobalNameData{
						Name:    vname,
						Address: uint64(addr),
						Size:    uint64(vsize),
						Type:    vtype,
					}
					wd.GlobalAddresses[vname] = globalInfo
				} else {
					// TODO
					// fmt.Printf("Variable but not simple expr... %s %x\n", vname, vaddr)
				}
			}
		}

		if entry.Tag == dwarf.TagSubprogram {
			entryReader.SkipChildren()
		}
	}
}

func (wd *WasmDebug) ParseDwarfVariables(wf FunctionFinder) error {
	wd.ParseDwarfGlobals()

	wd.FunctionDebug = make(map[int]string)
	if wd.FunctionSignature == nil {
		wd.FunctionSignature = make(map[int]string)
	}
	wd.LocalNames = make([]*LocalNameData, 0)

	if wd.DwarfData == nil {
		return nil
	}

	entryReader := wd.DwarfData.Reader()

	for {
		// Read all entries in sequence
		entry, err := entryReader.Next()
		if entry == nil || err == io.EOF {
			// We've reached the end of DWARF entries
			break
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

			log := false
			if strings.HasPrefix(spname, "main.") ||
				spname == "main" ||
				spname == "example_function" {
				//fmt.Printf("TagSubprogram %s %d\n", spname, sploc)
				log = true
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

					if log {
						//fmt.Printf(" Entry %v\n", entry)
					}

					vname := "<unknown>"
					vtype := ""
					vloc := int64(-1)
					vlocbytes := make([]byte, 0)
					for _, field := range entry.Field {
						if log {
							//fmt.Printf(" .. %v\n", field)
						}
						if field.Attr == dwarf.AttrName {
							vname = field.Val.(string)
						} else if field.Attr == dwarf.AttrType {
							switch field.Val.(type) {
							case dwarf.Offset:
								t := field.Val.(dwarf.Offset)
								ty, err := wd.DwarfData.Type(t)
								if err == nil {
									vtype = ty.String()
								}
							}
						} else if field.Attr == dwarf.AttrLocation {
							switch field.Val.(type) {
							case int64:
								vloc = field.Val.(int64)
							case []byte:
								vlocbytes = field.Val.([]byte)
							}
						}
					}

					if entry.Tag == dwarf.TagFormalParameter {
						if vloc != -1 {
							locdata := wd.DwarfLoc.ReadLocation(uint64(vloc))
							for _, ld := range locdata {
								// We have code ptr range here...
								if log {
									//fmt.Printf("  = LocationData %d %d %x\n", ld.StartAddress, ld.EndAddress, ld.Expression)
								}
								locs := ld.ExtractWasmLocations()
								for _, l := range locs {
									if l.IsLocal {
										wd.LocalNames = append(wd.LocalNames, &LocalNameData{
											StartPC: uint64(ld.StartAddress), //sploc),
											EndPC:   uint64(ld.EndAddress),   //sploc),
											Index:   int(l.Index),
											VarName: vname,
											VarType: vtype,
										})
										if log {
											//fmt.Printf("LocationLocal %s %s (%d-%d) %d local %d\n", spname, vname, ld.StartAddress, ld.EndAddress, sploc, l.Index)
										}
									}
								}
							}
						} else {
							ld := &LocationData{
								Expression: vlocbytes,
							}
							locs := ld.ExtractWasmLocations()
							for _, l := range locs {
								if l.IsLocal {
									wd.LocalNames = append(wd.LocalNames, &LocalNameData{
										StartPC: uint64(ld.StartAddress), //sploc),
										EndPC:   uint64(ld.EndAddress),   //sploc),
										Index:   int(l.Index),
										VarName: vname,
										VarType: vtype,
									})
									if log {
										//fmt.Printf("LocationLocal %s %s (%d-%d) %d local %d\n", spname, vname, ld.StartAddress, ld.EndAddress, sploc, l.Index)
									}
								}
							}
						}
						if len(params) > 0 {
							params = params + ", "
						}
						params = fmt.Sprintf("%s%s(%s)", params, vname, vtype)
					} else if entry.Tag == dwarf.TagVariable {

						if log {
							//fmt.Printf("  - Variable %v | %s %d [%x]\n", entry, vname, vloc, vlocbytes)
						}

						if vloc != -1 {
							locdata := wd.DwarfLoc.ReadLocation(uint64(vloc))
							for _, ld := range locdata {

								if log {
									//fmt.Printf("  = LOC %d-%d : %x\n", ld.StartAddress, ld.EndAddress, ld.Expression)
								}

								locs := ld.ExtractWasmLocations()
								for _, l := range locs {
									if l.IsLocal {
										// Store in the locals lookup...
										wd.LocalNames = append(wd.LocalNames, &LocalNameData{
											StartPC: uint64(ld.StartAddress),
											EndPC:   uint64(ld.EndAddress),
											Index:   int(l.Index),
											VarName: vname,
										})

										//										fmt.Printf("LocationLocalVariable %s %s %d-%d  local %d\n", spname, vname, ld.StartAddress, ld.EndAddress, l.Index)
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
				wd.FunctionSignature[fid] = fmt.Sprintf("%s(%s)", spname, params)
				wd.FunctionDebug[fid] = function_debug
			}
		}
	}

	return nil
}
