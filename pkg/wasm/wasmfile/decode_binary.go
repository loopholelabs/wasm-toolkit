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
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/loopholelabs/wasm-toolkit/pkg/wasm/expression"
	"github.com/loopholelabs/wasm-toolkit/pkg/wasm/types"
)

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

/**
 * Decode a wasm binary into a WasmFile
 *
 */
func (wf *WasmFile) DecodeBinary(data []byte) (err error) {
	/*
		defer func() {
			r := recover()
			if r != nil {
				switch x := r.(type) {
				case string:
					err = errors.New(x)
				case error:
					err = x
				default:
					err = errors.New("unknown panic")
				}
			}
		}()
	*/
	hd := binary.LittleEndian.Uint32(data)
	vr := binary.LittleEndian.Uint32(data[4:])

	if hd != WasmHeader || vr != WasmVersion {
		return fmt.Errorf("Invalid header/version %x/%x", hd, vr)
	}

	data = data[8:]

	rr := bytes.NewReader(data)

	for {
		sectionType, err := rr.ReadByte()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		sectionLength, err := binary.ReadUvarint(rr)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		sectionData := make([]byte, sectionLength)

		_, err = rr.Read(sectionData)
		if err == io.EOF {
			break
		}

		// Process each section

		if sectionType == byte(types.SectionCustom) {
			err = wf.ParseSectionCustom(sectionData)
		} else if sectionType == byte(types.SectionType) {
			err = wf.ParseSectionType(sectionData)
		} else if sectionType == byte(types.SectionImport) {
			err = wf.ParseSectionImport(sectionData)
		} else if sectionType == byte(types.SectionFunction) {
			err = wf.ParseSectionFunction(sectionData)
		} else if sectionType == byte(types.SectionTable) {
			err = wf.ParseSectionTable(sectionData)
		} else if sectionType == byte(types.SectionMemory) {
			err = wf.ParseSectionMemory(sectionData)
		} else if sectionType == byte(types.SectionGlobal) {
			err = wf.ParseSectionGlobal(sectionData)
		} else if sectionType == byte(types.SectionExport) {
			err = wf.ParseSectionExport(sectionData)
		} else if sectionType == byte(types.SectionStart) {
			err = wf.ParseSectionStart(sectionData)
		} else if sectionType == byte(types.SectionElem) {
			err = wf.ParseSectionElem(sectionData)
		} else if sectionType == byte(types.SectionCode) {
			err = wf.ParseSectionCode(sectionData)
		} else if sectionType == byte(types.SectionData) {
			err = wf.ParseSectionData(sectionData)
		} else if sectionType == byte(types.SectionDataCount) {
			err = wf.ParseSectionDataCount(sectionData)
		} else {
			return fmt.Errorf("Unknown section %d", sectionType)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

/**
 * Parse a DataCount section
 *
 */
func (wf *WasmFile) ParseSectionDataCount(data []byte) error {
	/*
		ptr := 0
		dataCount, l := binary.Uvarint(data)
	*/
	// For now, we don't care...
	return nil
}

func getDataContext(data []byte) []byte {
	l := len(data)
	if l > 16 {
		l = 16
	}
	return data[:l]
}

/**
 * Parse a Data section
 *
 */
func (wf *WasmFile) ParseSectionData(data []byte) error {
	ptr := 0
	dataVecLength, l := binary.Uvarint(data)
	if l <= 0 {
		return fmt.Errorf("Error decoding SectionData dataVecLength %x", getDataContext(data))
	}
	ptr += l

	for i := 0; i < int(dataVecLength); i++ {
		memindex, l := binary.Uvarint(data[ptr:])
		if l <= 0 {
			return fmt.Errorf("Error decoding SectionData memindex %x", getDataContext(data))
		}
		ptr += l
		offset, l, err := expression.NewExpression(data[ptr:], 0)
		if err != nil {
			return err
		}
		ptr += l
		bytesLength, l := binary.Uvarint(data[ptr:])
		if l <= 0 {
			return fmt.Errorf("Error decoding SectionData bytesLength %x", getDataContext(data))
		}
		ptr += l
		if ptr+int(bytesLength) > len(data) {
			return fmt.Errorf("Error decoding SectionData not enough data %d > %d", ptr+int(bytesLength), len(data))
		}
		dataBytes := data[ptr : ptr+int(bytesLength)]
		ptr += int(bytesLength)

		d := &DataEntry{
			MemIndex: int(memindex),
			Offset:   offset,
			Data:     dataBytes,
		}

		wf.Data = append(wf.Data, d)
	}
	return nil
}

/**
 * Parse a Code section
 *
 */
func (wf *WasmFile) ParseSectionCode(data []byte) error {
	ptr := 0
	codeVecLength, l := binary.Uvarint(data)
	if l <= 0 {
		return fmt.Errorf("Error decoding SectionCode codeVecLength %x", getDataContext(data))
	}
	ptr += l

	for i := 0; i < int(codeVecLength); i++ {
		clen, l := binary.Uvarint(data[ptr:])
		if l <= 0 {
			return fmt.Errorf("Error decoding SectionCode clen %x", getDataContext(data))
		}
		ptr += l
		if ptr+int(clen) > len(data) {
			return fmt.Errorf("Error decoding SectionCode not enough data %d > %d", ptr+int(clen), len(data))
		}

		codeptr := uint64(ptr) // Start of the code
		code := data[ptr : ptr+int(clen)]
		ptr += int(clen)

		locals := make([]types.ValType, 0)

		vclen, l := binary.Uvarint(code)
		if l <= 0 {
			return fmt.Errorf("Error decoding SectionCode vclen %x", getDataContext(data))
		}
		locptr := l

		for lo := 0; lo < int(vclen); lo++ {
			paramLen, ll := binary.Uvarint(code[locptr:])
			if l <= 0 {
				return fmt.Errorf("Error decoding SectionCode paramLen %x", getDataContext(data))
			}
			locptr += ll
			ty := code[locptr]
			locptr++

			for lod := 0; lod < int(paramLen); lod++ {
				locals = append(locals, types.ValType(ty))
			}
		}

		expression, _, err := expression.NewExpression(code[locptr:], codeptr+uint64(locptr))
		if err != nil {
			return err
		}

		c := &CodeEntry{
			Locals:         locals,
			PCValid:        true,
			CodeSectionPtr: codeptr,
			CodeSectionLen: clen,
			Expression:     expression,
		}
		wf.Code = append(wf.Code, c)
	}
	return nil
}

/**
 * Parse an Elem section
 *
 */
func (wf *WasmFile) ParseSectionElem(data []byte) error {
	ptr := 0
	elemVecLength, l := binary.Uvarint(data)
	ptr += l

	for i := 0; i < int(elemVecLength); i++ {
		tableIndex, l := binary.Uvarint(data[ptr:])
		ptr += l
		offset, l, err := expression.NewExpression(data[ptr:], 0)
		if err != nil {
			return err
		}

		ptr += l
		funcVecLength, l := binary.Uvarint(data[ptr:])
		ptr += l
		indexes := make([]uint64, 0)
		for f := 0; f < int(funcVecLength); f++ {
			funcIndex, l := binary.Uvarint(data[ptr:])
			ptr += l
			indexes = append(indexes, funcIndex)
		}
		e := &ElemEntry{
			TableIndex: int(tableIndex),
			Offset:     offset,
			Indexes:    indexes,
		}
		wf.Elem = append(wf.Elem, e)
	}
	return nil
}

/**
 * Parse an Import section
 *
 */
func (wf *WasmFile) ParseSectionImport(data []byte) error {
	ptr := 0
	importVecLength, l := binary.Uvarint(data)
	ptr += l

	for i := 0; i < int(importVecLength); i++ {
		modLength, l := binary.Uvarint(data[ptr:])
		ptr += l
		mod := data[ptr : ptr+int(modLength)]
		ptr += int(modLength)
		nameLength, l := binary.Uvarint(data[ptr:])
		ptr += l
		name := data[ptr : ptr+int(nameLength)]
		ptr += int(nameLength)
		importType := data[ptr]
		ptr++
		importIndex, l := binary.Uvarint(data[ptr:])
		ptr += l
		e := &ImportEntry{
			Module: string(mod),
			Name:   string(name),
			Type:   types.ExportType(importType),
			Index:  int(importIndex),
		}
		wf.Import = append(wf.Import, e)
	}
	return nil
}

/**
 * Parse a Function section
 *
 */
func (wf *WasmFile) ParseSectionFunction(data []byte) error {
	ptr := 0
	funcVecLength, l := binary.Uvarint(data)
	ptr += l

	for i := 0; i < int(funcVecLength); i++ {
		id, l := binary.Uvarint(data[ptr:])

		f := &FunctionEntry{
			TypeIndex: int(id),
		}
		wf.Function = append(wf.Function, f)
		ptr += l
	}
	return nil
}

/**
 * Parse a Table section
 *
 */
func (wf *WasmFile) ParseSectionTable(data []byte) error {
	ptr := 0
	tableVecLength, l := binary.Uvarint(data)
	ptr += l

	for i := 0; i < int(tableVecLength); i++ {
		tableType := data[ptr]
		ptr++
		limitMax := uint64(0)
		limitMin := uint64(0)
		var l int
		if data[ptr] == types.LimitTypeMin {
			ptr++
			limitMin, l = binary.Uvarint(data[ptr:])
			ptr += l
		} else if data[ptr] == types.LimitTypeMinMax {
			ptr++
			limitMin, l = binary.Uvarint(data[ptr:])
			ptr += l
			limitMax, l = binary.Uvarint(data[ptr:])
			ptr += l
		} else {
			return fmt.Errorf("Invalid limit type in TableSection %d", data[ptr])
		}
		t := &TableEntry{
			TableType: tableType,
			LimitMin:  int(limitMin),
			LimitMax:  int(limitMax),
		}
		wf.Table = append(wf.Table, t)
	}
	return nil
}

/**
 * Parse a Memory section
 *
 */
func (wf *WasmFile) ParseSectionMemory(data []byte) error {
	ptr := 0
	memoryVecLength, l := binary.Uvarint(data)
	ptr += l

	for i := 0; i < int(memoryVecLength); i++ {
		limitMax := uint64(0)
		limitMin := uint64(0)
		var l int
		if data[ptr] == types.LimitTypeMin {
			ptr++
			limitMin, l = binary.Uvarint(data[ptr:])
			ptr += l
		} else if data[ptr] == types.LimitTypeMinMax {
			ptr++
			limitMin, l = binary.Uvarint(data[ptr:])
			ptr += l
			limitMax, l = binary.Uvarint(data[ptr:])
			ptr += l
		} else {
			return fmt.Errorf("Invalid limit type in MemorySection %d", data[ptr])
		}
		m := &MemoryEntry{
			LimitMin: int(limitMin),
			LimitMax: int(limitMax),
		}
		wf.Memory = append(wf.Memory, m)
	}
	return nil
}

/**
 * Parse a Global section
 *
 */
func (wf *WasmFile) ParseSectionGlobal(data []byte) error {
	ptr := 0
	globalVecLength, l := binary.Uvarint(data)
	ptr += l

	for i := 0; i < int(globalVecLength); i++ {
		valType := data[ptr]
		ptr++
		valMut := data[ptr]
		ptr++
		// Read the init expression
		expression, n, err := expression.NewExpression(data[ptr:], 0)
		if err != nil {
			return err
		}

		ptr += n

		g := &GlobalEntry{
			Type:       types.ValType(valType),
			Mut:        valMut,
			Expression: expression,
		}
		wf.Global = append(wf.Global, g)
	}
	return nil
}

/**
 * Parse an Export section
 *
 */
func (wf *WasmFile) ParseSectionExport(data []byte) error {
	ptr := 0
	exportVecLength, l := binary.Uvarint(data)
	ptr += l

	for i := 0; i < int(exportVecLength); i++ {
		nameLength, l := binary.Uvarint(data[ptr:])
		ptr += l
		name := data[ptr : ptr+int(nameLength)]
		ptr += int(nameLength)
		exportType := data[ptr]
		ptr++
		exportIndex, l := binary.Uvarint(data[ptr:])
		ptr += l
		e := &ExportEntry{
			Name:  string(name),
			Type:  types.ExportType(exportType),
			Index: int(exportIndex),
		}
		wf.Export = append(wf.Export, e)
	}
	return nil
}

/**
 * Parse a Start section
 *
 */
func (wf *WasmFile) ParseSectionStart(data []byte) error {
	return nil
	//return fmt.Errorf("TODO: ParseSectionStart %d\n", len(data))
}

/**
 * Parse a Custom section
 *
 */
func (wf *WasmFile) ParseSectionCustom(data []byte) error {
	ptr := 0
	nameLength, l := binary.Uvarint(data)
	ptr += l

	nameData := data[ptr : ptr+int(nameLength)]
	ptr += int(nameLength)

	c := &CustomEntry{
		Name: string(nameData),
		Data: data[ptr:],
	}

	wf.Custom = append(wf.Custom, c)
	return nil
}

/**
 * Parse a Type section
 *
 */
func (wf *WasmFile) ParseSectionType(data []byte) error {
	ptr := 0
	typeVecLength, l := binary.Uvarint(data)
	ptr += l

	for i := 0; i < int(typeVecLength); i++ {
		t := &TypeEntry{
			Param:  make([]types.ValType, 0),
			Result: make([]types.ValType, 0),
		}

		// Read a functype
		if data[ptr] == types.FuncTypePrefix {
			ptr++
			// Now read param / result vectors
			paramVecLength, l := binary.Uvarint(data[ptr:])
			ptr += l
			for p := 0; p < int(paramVecLength); p++ {
				t.Param = append(t.Param, types.ValType(data[ptr]))
				ptr++
			}
			resultVecLength, l := binary.Uvarint(data[ptr:])
			ptr += l
			for p := 0; p < int(resultVecLength); p++ {
				t.Result = append(t.Result, types.ValType(data[ptr]))
				ptr++
			}
			wf.Type = append(wf.Type, t)
		} else {
			return fmt.Errorf("Invalid type %d", data[ptr])
		}
	}
	return nil
}
