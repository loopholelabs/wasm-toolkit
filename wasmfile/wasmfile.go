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
	"bytes"
	"debug/dwarf"
	"fmt"
	"io/ioutil"
	"strings"
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

	// dwarf debugging data
	dwarfLoc    []byte
	dwarfData   *dwarf.Data
	lineNumbers map[uint64]LineInfo
	// debug info derived from dwarf
	functionDebug     map[int]string
	functionSignature map[int]string
	localNames        []*LocalNameData

	// custom names section data
	FunctionNames map[int]string
	globalNames   map[int]string
	dataNames     map[int]string
}

const WasmHeader uint32 = 0x6d736100
const WasmVersion uint32 = 0x00000001

type ValType byte

const (
	ValI32  ValType = 0x7f
	ValI64  ValType = 0x7e
	ValF32  ValType = 0x7d
	ValF64  ValType = 0x7c
	ValNone ValType = 0x40
)

var ValTypeToByte map[string]ValType
var ByteToValType map[ValType]string

func init() {
	ValTypeToByte = make(map[string]ValType)
	ValTypeToByte["i32"] = ValI32
	ValTypeToByte["i64"] = ValI64
	ValTypeToByte["f32"] = ValF32
	ValTypeToByte["f64"] = ValF64
	ValTypeToByte["none"] = ValNone

	ByteToValType = make(map[ValType]string)
	ByteToValType[ValI32] = "i32"
	ByteToValType[ValI64] = "i64"
	ByteToValType[ValF32] = "f32"
	ByteToValType[ValF64] = "f64"
	ByteToValType[ValNone] = "none"
}

const (
	LimitTypeMin    byte = 0x00
	LimitTypeMinMax byte = 0x01
)

type ExportType byte

const (
	ExportFunc   ExportType = 0
	ExportTable  ExportType = 1
	ExportMem    ExportType = 2
	ExportGlobal ExportType = 3
)

const FuncTypePrefix byte = 0x60

const TableTypeFuncref byte = 0x70

type SectionId byte

const (
	SectionCustom    SectionId = 0
	SectionType      SectionId = 1
	SectionImport    SectionId = 2
	SectionFunction  SectionId = 3
	SectionTable     SectionId = 4
	SectionMemory    SectionId = 5
	SectionGlobal    SectionId = 6
	SectionExport    SectionId = 7
	SectionStart     SectionId = 8
	SectionElem      SectionId = 9
	SectionCode      SectionId = 10
	SectionData      SectionId = 11
	SectionDataCount SectionId = 12
)

type FunctionEntry struct {
	TypeIndex int
}

type TypeEntry struct {
	Param  []ValType
	Result []ValType
}

type CustomEntry struct {
	Name string
	Data []byte
}

type ExportEntry struct {
	Name  string
	Type  ExportType
	Index int
}

type ImportEntry struct {
	Module string
	Name   string
	Type   ExportType
	Index  int
}

type TableEntry struct {
	TableType byte
	LimitMin  int
	LimitMax  int
}

type GlobalEntry struct {
	Type       ValType
	Mut        byte
	Expression []*Expression
}

type MemoryEntry struct {
	LimitMin int
	LimitMax int
}

type CodeEntry struct {
	Locals         []ValType
	PCValid        bool
	CodeSectionPtr uint64
	CodeSectionLen uint64
	Expression     []*Expression
}

type DataEntry struct {
	MemIndex int
	Offset   []*Expression
	Data     []byte
}

type ElemEntry struct {
	TableIndex int
	Offset     []*Expression
	Indexes    []uint64
}

type LineInfo struct {
	Filename   string
	Linenumber int
	Column     int
}

// Create a new WasmFile from a file
func New(filename string) (*WasmFile, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	wf := &WasmFile{}
	err = wf.DecodeBinary(data)
	return wf, err
}

// Create a new WasmFile from a file
func NewFromWat(filename string) (*WasmFile, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	wf := &WasmFile{}
	err = wf.DecodeWat(data)
	return wf, err
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

func (wf *WasmFile) SetGlobal(name string, t ValType, expr string) {
	ex := make([]*Expression, 0)
	e := &Expression{}
	e.DecodeWat(expr, wf, nil)
	ex = append(ex, e)

	idx := wf.LookupGlobalID(name)
	if idx == -1 {
		panic("Global not found")
	}

	wf.Global[idx].Type = t
	wf.Global[idx].Expression = ex
}

func (wf *WasmFile) AddTypeMaybe(te *TypeEntry) int {
	for idx, t := range wf.Type {
		if t.Equals(te) {
			return idx
		}
	}
	wf.Type = append(wf.Type, te)
	return len(wf.Type) - 1
}

func (wf *WasmFile) AddDataFrom(addr int32, wfSource *WasmFile) int32 {
	ptr := addr
	for _, d := range wfSource.Data {
		// Relocate the data
		d.Offset = []*Expression{
			{
				Opcode:   InstrToOpcode["i32.const"],
				I32Value: ptr,
			},
		}

		wf.Data = append(wf.Data, d)
		ptr += int32(len(d.Data))
		ptr = (ptr + 3) & -4
	}
	return ptr
}

func (wf *WasmFile) AddData(name string, data []byte) {
	ptr := int32(0)
	if len(wf.Data) > 0 {
		prev := wf.Data[len(wf.Data)-1]
		ptr = prev.Offset[0].I32Value + int32(len(prev.Data))
	}

	// Align things...
	ptr = (ptr + 3) & -4

	idx := len(wf.Data)
	wf.Data = append(wf.Data, &DataEntry{
		MemIndex: 0,
		Offset: []*Expression{
			{
				Opcode:   InstrToOpcode["i32.const"],
				I32Value: ptr,
			},
		},
		Data: data,
	})
	wf.dataNames[idx] = name
}

func (wf *WasmFile) AddFuncsFrom(wfSource *WasmFile) {
	globalModification := make(map[int]int)
	for idx, g := range wfSource.Global {
		newidx := len(wf.Global)
		globalModification[idx] = newidx
		wf.Global = append(wf.Global, g)
		name := wfSource.GetGlobalIdentifier(idx)

		wf.globalNames[newidx] = name
	}

	callModification := make(map[int]int) // old fid -> new fid

	// Deal with any imports
	for idx, i := range wfSource.Import {
		// Check if it's already being imported as something else...
		var newidx = -1
		for nidx, i2 := range wf.Import {
			if i.Module == i2.Module && i.Name == i2.Name {
				newidx = nidx
				break
			}
		}
		if newidx != -1 {
			callModification[idx] = newidx
		} else {
			// Need to add a new import then... (This means relocating every call as well)
			callModification[idx] = len(wf.Import)
			newidx := len(wf.Import)

			// TODO: Might need to add a type if there isn't one already
			t := wfSource.Type[i.Type]
			i.Index = wf.AddTypeMaybe(t)

			wf.Import = append(wf.Import, i)
			// Need to fix up any existing code, and the function names table
			newmap := make(map[int]string, 0)
			for idx, name := range wf.FunctionNames {
				if idx >= newidx {
					newmap[idx+1] = name
				} else {
					newmap[idx] = name
				}
			}
			name := wfSource.GetFunctionIdentifier(idx, true)
			if name != "" {
				newmap[newidx] = name
			}
			wf.FunctionNames = newmap

			rmap := make(map[int]int)
			for i := 0; i < len(wf.Code)+len(wf.Import); i++ {
				// Relocate everything at or above newidx
				if i >= newidx {
					rmap[i] = i + 1
				}
			}

			// Modify any exports
			for _, ex := range wf.Export {
				if ex.Type == ExportFunc && ex.Index >= newidx {
					ex.Index++
				}
			}

			for _, ce := range wf.Code {
				ce.ModifyAllCalls(rmap)
			}

		}
	}

	fmt.Printf("Imports sorted out... %v\n", callModification) // 0->0 ok

	for idx, f := range wfSource.Function {
		t := wfSource.Type[f.TypeIndex]
		name := wfSource.GetFunctionIdentifier(len(wfSource.Import)+idx, true)

		newidx := len(wf.Import) + len(wf.Function)

		// Add the functions in, copying the type if needed...
		wf.Function = append(wf.Function, f)
		f.TypeIndex = wf.AddTypeMaybe(t)

		// Add the function name if there is one
		if name != "" {
			wf.FunctionNames[newidx] = name
		}

		callModification[len(wfSource.Import)+idx] = newidx
	}

	// Now add the code
	for _, c := range wfSource.Code {

		c.ModifyAllCalls(callModification)
		c.ModifyAllGlobals(globalModification)
		wf.Code = append(wf.Code, c)
	}

}

func (ce *CodeEntry) ModifyAllGlobals(m map[int]int) {
	for _, e := range ce.Expression {
		newid, ok := m[e.GlobalIndex]
		if ok {
			e.GlobalIndex = newid
		}
	}
}

func (ce *CodeEntry) ModifyAllCalls(m map[int]int) {
	for _, e := range ce.Expression {
		newid, ok := m[e.FuncIndex]
		if ok {
			e.FuncIndex = newid
		}
	}
}

func (wf *WasmFile) ExpressionFromWat(d string) ([]*Expression, error) {
	newex := make([]*Expression, 0)
	lines := strings.Split(d, "\n")
	for _, toline := range lines {
		cptr := strings.Index(toline, ";;")
		if cptr != -1 {
			toline = toline[:cptr]
		}
		toline = strings.Trim(toline, Whitespace)
		if len(toline) > 0 {
			newe := &Expression{}
			err := newe.DecodeWat(toline, wf, nil)
			if err != nil {
				return newex, err
			}
			newex = append(newex, newe)
		}
	}
	return newex, nil
}

func (ce *CodeEntry) ReplaceInstr(wf *WasmFile, from string, to string) error {

	newex, err := wf.ExpressionFromWat(to)
	if err != nil {
		return err
	}

	// Now we need to find where to replace this code...
	adjustedExpression := make([]*Expression, 0)
	for _, e := range ce.Expression {
		var buf bytes.Buffer
		e.EncodeWat(&buf, "", wf)
		cd := buf.String()
		cend := strings.Index(cd, ";;")
		if cend != -1 {
			cd = cd[:cend]
		}

		if strings.Trim(cd, Whitespace) == from {
			// Replace it!
			for _, ne := range newex {
				adjustedExpression = append(adjustedExpression, ne)
			}
		} else {
			adjustedExpression = append(adjustedExpression, e)
		}
	}
	ce.Expression = adjustedExpression
	return nil
}

func (ce *CodeEntry) InsertFuncStart(wf *WasmFile, to string) error {
	newex, err := wf.ExpressionFromWat(to)
	if err != nil {
		return err
	}

	// Now we need to find where to replace this code...
	adjustedExpression := make([]*Expression, 0)
	for _, e := range newex {
		adjustedExpression = append(adjustedExpression, e)
	}

	for _, e := range ce.Expression {
		adjustedExpression = append(adjustedExpression, e)
	}
	ce.Expression = adjustedExpression
	return nil
}

func (ce *CodeEntry) InsertFuncEnd(wf *WasmFile, to string) error {
	newex, err := wf.ExpressionFromWat(to)
	if err != nil {
		return err
	}

	ce.Expression = append(ce.Expression, newex...)
	return nil
}

func (ce *CodeEntry) InsertAfterRelocating(wf *WasmFile, to string) error {
	newex, err := wf.ExpressionFromWat(to)
	if err != nil {
		return err
	}

	// Now we need to find where to insert the code
	adjustedExpression := make([]*Expression, 0)
	for _, e := range ce.Expression {
		adjustedExpression = append(adjustedExpression, e)
		if e.Relocating {
			for _, ne := range newex {
				adjustedExpression = append(adjustedExpression, ne)
			}
		}
	}
	ce.Expression = adjustedExpression
	return nil
}

func (te *TypeEntry) Equals(te2 *TypeEntry) bool {
	if len(te.Param) != len(te2.Param) || len(te.Result) != len(te2.Result) {
		return false
	}
	for idx, v := range te.Param {
		if v != te2.Param[idx] {
			return false
		}
	}
	for idx, v := range te.Result {
		if v != te2.Result[idx] {
			return false
		}
	}
	return true
}
