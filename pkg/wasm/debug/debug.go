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
