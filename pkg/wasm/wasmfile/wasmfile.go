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
	"fmt"

	"github.com/loopholelabs/wasm-toolkit/pkg/wasm/debug"
	"github.com/loopholelabs/wasm-toolkit/pkg/wasm/expression"
	"github.com/loopholelabs/wasm-toolkit/pkg/wasm/types"
)

type WasmFile struct {
	// Each section of the wasm file
	Function []*FunctionEntry
	Type     []*TypeEntry
	Custom   []*CustomEntry
	Export   []*ExportEntry
	Import   []*ImportEntry
	Table    []*TableEntry
	Global   []*GlobalEntry
	Memory   []*MemoryEntry
	Code     []*CodeEntry
	Data     []*DataEntry
	Elem     []*ElemEntry

	Debug *debug.WasmDebug
}

const WasmHeader uint32 = 0x6d736100
const WasmVersion uint32 = 0x00000001

// FunctionEntry
type FunctionEntry struct {
	TypeIndex int
}

// TypeEntry
type TypeEntry struct {
	Param  []types.ValType
	Result []types.ValType
}

// CustomEntry
type CustomEntry struct {
	Name string
	Data []byte
}

// ExportEntry
type ExportEntry struct {
	Name  string
	Type  types.ExportType
	Index int
}

// ImportEntry
type ImportEntry struct {
	Module string
	Name   string
	Type   types.ExportType
	Index  int
}

// TableEntry
type TableEntry struct {
	TableType byte
	LimitMin  int
	LimitMax  int
}

// GlobalEntry
type GlobalEntry struct {
	Type       types.ValType
	Mut        byte
	Expression []*expression.Expression
}

// MemoryEntry
type MemoryEntry struct {
	LimitMin int
	LimitMax int
}

// CodeEntry
type CodeEntry struct {
	Locals         []types.ValType
	PCValid        bool
	CodeSectionPtr uint64
	CodeSectionLen uint64
	Expression     []*expression.Expression
}

// DataEntry
type DataEntry struct {
	MemIndex int
	Offset   []*expression.Expression
	Data     []byte
}

// ElemEntry
type ElemEntry struct {
	TableIndex int
	Offset     []*expression.Expression
	Indexes    []uint64
}

func NewEmpty() *WasmFile {
	return &WasmFile{
		Debug: debug.NewEmpty(),
	}
}

func (wf *WasmFile) GetCustomSectionData(name string) []byte {
	for _, c := range wf.Custom {
		if c.Name == name {
			return c.Data
		}
	}
	return nil
}

func (wf *WasmFile) FindFunction(pc uint64) int {
	for index, c := range wf.Code {

		if c.PCValid && pc >= c.CodeSectionPtr && pc <= (c.CodeSectionPtr+c.CodeSectionLen) {
			return len(wf.Import) + index
		}
	}
	return -1
}

func (wf *WasmFile) LookupImport(n string) int {
	for idx, i := range wf.Import {
		iname := fmt.Sprintf("%s:%s", i.Module, i.Name)
		if iname == n {
			return idx
		}
	}
	return -1
}

func (t *TypeEntry) Clone() *TypeEntry {
	newType := &TypeEntry{
		Param:  make([]types.ValType, 0),
		Result: make([]types.ValType, 0),
	}
	for _, v := range t.Result {
		newType.Result = append(newType.Result, v)
	}
	for _, v := range t.Param {
		newType.Param = append(newType.Param, v)
	}
	return newType
}
