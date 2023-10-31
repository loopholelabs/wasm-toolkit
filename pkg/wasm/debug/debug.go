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
)

type WasmDebug struct {
	// These come from the 'name' custom section
	FunctionNames map[int]string
	GlobalNames   map[int]string
	DataNames     map[int]string

	// dwarf debugging data
	DwarfLoc    *DwarfLocations
	DwarfData   *dwarf.Data
	LineNumbers map[uint64]LineInfo
	// debug info derived from dwarf
	FunctionDebug     map[int]string
	FunctionSignature map[int]string
	LocalNames        []*LocalNameData

	GlobalAddresses map[string]*GlobalNameData
}

func NewEmpty() *WasmDebug {
	wd := &WasmDebug{}
	wd.FunctionNames = make(map[int]string)
	wd.GlobalNames = make(map[int]string)
	wd.DataNames = make(map[int]string)

	wd.LineNumbers = make(map[uint64]LineInfo)
	wd.FunctionDebug = make(map[int]string)
	wd.FunctionSignature = make(map[int]string)
	wd.LocalNames = make([]*LocalNameData, 0)
	wd.GlobalAddresses = make(map[string]*GlobalNameData)

	return wd
}

type LocalNameData struct {
	StartPC uint64
	EndPC   uint64
	Index   int
	VarName string
	VarType string
}

type GlobalNameData struct {
	Name    string
	Address uint64
	Size    uint64
	Type    string
}

type CustomSectionProvider interface {
	GetCustomSectionData(name string) []byte
}

func (wd *WasmDebug) ParseDwarf(wf CustomSectionProvider) error {
	debug_abbrev := wf.GetCustomSectionData(".debug_abbrev")
	debug_aranges := wf.GetCustomSectionData(".debug_aranges")
	debug_info := wf.GetCustomSectionData(".debug_info")
	debug_line := wf.GetCustomSectionData(".debug_line")
	debug_pubnames := wf.GetCustomSectionData(".debug_pubnames")
	debug_ranges := wf.GetCustomSectionData(".debug_ranges")
	debug_str := wf.GetCustomSectionData(".debug_str")

	debug_loc := wf.GetCustomSectionData(".debug_loc")
	wd.DwarfLoc = NewDwarfLocations(debug_loc)

	debug_frame := make([]byte, 0) // call frame info

	dd, err := dwarf.New(debug_abbrev, debug_aranges, debug_frame, debug_info, debug_line, debug_pubnames, debug_ranges, debug_str)
	if err != nil {
		return nil // ok, but lets move on and ignore the error.
	}

	wd.DwarfData = dd
	return nil
}

// Renumber functions using a remap
func (wd *WasmDebug) RenumberFunctions(remap map[int]int) {
	// This modifies FunctionNames, functionDebug, functionSignature
	newFunctionNames := make(map[int]string)
	newFunctionDebug := make(map[int]string)
	newFunctionSignature := make(map[int]string)
	for o, n := range remap {
		v, ok := wd.FunctionNames[o]
		if ok {
			newFunctionNames[n] = v
		}
		v, ok = wd.FunctionDebug[o]
		if ok {
			newFunctionDebug[n] = v
		}
		v, ok = wd.FunctionSignature[o]
		if ok {
			newFunctionSignature[n] = v
		}
	}
	wd.FunctionNames = newFunctionNames
	wd.FunctionDebug = newFunctionDebug
	wd.FunctionSignature = newFunctionSignature
}
