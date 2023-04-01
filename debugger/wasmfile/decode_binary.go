package wasmfile

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

func (wf *WasmFile) DecodeBinary(data []byte) error {
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

		if sectionType == byte(SectionCustom) {
			wf.ParseSectionCustom(sectionData)
		} else if sectionType == byte(SectionType) {
			wf.ParseSectionType(sectionData)
		} else if sectionType == byte(SectionImport) {
			wf.ParseSectionImport(sectionData)
		} else if sectionType == byte(SectionFunction) {
			wf.ParseSectionFunction(sectionData)
		} else if sectionType == byte(SectionTable) {
			wf.ParseSectionTable(sectionData)
		} else if sectionType == byte(SectionMemory) {
			wf.ParseSectionMemory(sectionData)
		} else if sectionType == byte(SectionGlobal) {
			wf.ParseSectionGlobal(sectionData)
		} else if sectionType == byte(SectionExport) {
			wf.ParseSectionExport(sectionData)
		} else if sectionType == byte(SectionStart) {
			wf.ParseSectionStart(sectionData)
		} else if sectionType == byte(SectionElem) {
			wf.ParseSectionElem(sectionData)
		} else if sectionType == byte(SectionCode) {
			wf.ParseSectionCode(sectionData)
		} else if sectionType == byte(SectionData) {
			wf.ParseSectionData(sectionData)
		}

	}
	return nil
}

func (wf *WasmFile) ParseSectionData(data []byte) {
	ptr := 0
	dataVecLength, l := binary.Uvarint(data)
	ptr += l

	for i := 0; i < int(dataVecLength); i++ {
		memindex, l := binary.Uvarint(data[ptr:])
		ptr += l
		offset, l := NewExpression(data[ptr:], 0)
		ptr += l
		bytesLength, l := binary.Uvarint(data[ptr:])
		ptr += l
		dataBytes := data[ptr : ptr+int(bytesLength)]
		ptr += int(bytesLength)

		d := &DataEntry{
			MemIndex: int(memindex),
			Offset:   offset,
			Data:     dataBytes,
		}

		wf.Data = append(wf.Data, d)
	}
}

func (wf *WasmFile) ParseSectionCode(data []byte) {
	ptr := 0
	codeVecLength, l := binary.Uvarint(data)
	ptr += l

	for i := 0; i < int(codeVecLength); i++ {
		clen, l := binary.Uvarint(data[ptr:])
		ptr += l

		codeptr := uint64(ptr) // Start of the code
		code := data[ptr : ptr+int(clen)]
		ptr += int(clen)

		locals := make([]ValType, 0)

		vclen, l := binary.Uvarint(code)
		locptr := l

		for lo := 0; lo < int(vclen); lo++ {
			paramLen, ll := binary.Uvarint(code[locptr:])
			locptr += ll
			for lod := 0; lod < int(paramLen); lod++ {
				ty := code[locptr+lod]
				locals = append(locals, ValType(ty))
			}
			locptr += int(paramLen)
		}

		expression, _ := NewExpression(code[locptr:], codeptr+uint64(locptr))

		c := &CodeEntry{
			Locals:         locals,
			CodeSectionPtr: codeptr,
			ExprData:       code[locptr:], // TODO
			Expression:     expression,
		}
		wf.Code = append(wf.Code, c)
	}
}

func (wf *WasmFile) ParseSectionElem(data []byte) {
	ptr := 0
	elemVecLength, l := binary.Uvarint(data)
	ptr += l

	for i := 0; i < int(elemVecLength); i++ {
		tableIndex, l := binary.Uvarint(data[ptr:])
		ptr += l
		offset, l := NewExpression(data[ptr:], 0)
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
}

func (wf *WasmFile) ParseSectionImport(data []byte) {
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
			Type:   ExportType(importType),
			Index:  int(importIndex),
		}
		wf.Import = append(wf.Import, e)
	}
}

func (wf *WasmFile) ParseSectionFunction(data []byte) {
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
}

func (wf *WasmFile) ParseSectionTable(data []byte) {
	ptr := 0
	tableVecLength, l := binary.Uvarint(data)
	ptr += l

	for i := 0; i < int(tableVecLength); i++ {
		tableType := data[ptr]
		ptr++
		limitMax := uint64(0)
		limitMin := uint64(0)
		var l int
		if data[ptr] == LimitTypeMin {
			ptr++
			limitMin, l = binary.Uvarint(data[ptr:])
			ptr += l
		} else if data[ptr] == LimitTypeMinMax {
			ptr++
			limitMin, l = binary.Uvarint(data[ptr:])
			ptr += l
			limitMax, l = binary.Uvarint(data[ptr:])
			ptr += l
		} else {
			panic("Invalid limit type")
		}
		t := &TableEntry{
			TableType: tableType,
			LimitMin:  int(limitMin),
			LimitMax:  int(limitMax),
		}
		wf.Table = append(wf.Table, t)
	}
}

func (wf *WasmFile) ParseSectionMemory(data []byte) {
	ptr := 0
	memoryVecLength, l := binary.Uvarint(data)
	ptr += l

	for i := 0; i < int(memoryVecLength); i++ {
		limitMax := uint64(0)
		limitMin := uint64(0)
		var l int
		if data[ptr] == LimitTypeMin {
			ptr++
			limitMin, l = binary.Uvarint(data[ptr:])
			ptr += l
		} else if data[ptr] == LimitTypeMinMax {
			ptr++
			limitMin, l = binary.Uvarint(data[ptr:])
			ptr += l
			limitMax, l = binary.Uvarint(data[ptr:])
			ptr += l
		} else {
			panic("Invalid limit type")
		}
		m := &MemoryEntry{
			LimitMin: int(limitMin),
			LimitMax: int(limitMax),
		}
		wf.Memory = append(wf.Memory, m)
	}
}

func (wf *WasmFile) ParseSectionGlobal(data []byte) {
	ptr := 0
	globalVecLength, l := binary.Uvarint(data)
	ptr += l

	for i := 0; i < int(globalVecLength); i++ {
		valType := data[ptr]
		ptr++
		valMut := data[ptr]
		ptr++
		g := &GlobalEntry{
			Type: ValType(valType),
			Mut:  valMut,
		}
		wf.Global = append(wf.Global, g)
	}
}

func (wf *WasmFile) ParseSectionExport(data []byte) {
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
			Type:  ExportType(exportType),
			Index: int(exportIndex),
		}
		wf.Export = append(wf.Export, e)
	}
}

func (wf *WasmFile) ParseSectionStart(data []byte) {
	fmt.Printf("ParseSectionStart %d\n", len(data))
	panic("TODO: ParseSectionStart")
}

func (wf *WasmFile) ParseSectionCustom(data []byte) {
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
}

func (wf *WasmFile) ParseSectionType(data []byte) {
	ptr := 0
	typeVecLength, l := binary.Uvarint(data)
	ptr += l

	for i := 0; i < int(typeVecLength); i++ {
		t := &TypeEntry{
			Param:  make([]ValType, 0),
			Result: make([]ValType, 0),
		}

		// Read a functype
		if data[ptr] == FuncTypePrefix {
			ptr++
			// Now read param / result vectors
			paramVecLength, l := binary.Uvarint(data[ptr:])
			ptr += l
			for p := 0; p < int(paramVecLength); p++ {
				t.Param = append(t.Param, ValType(data[ptr]))
				ptr++
			}
			resultVecLength, l := binary.Uvarint(data[ptr:])
			ptr += l
			for p := 0; p < int(resultVecLength); p++ {
				t.Result = append(t.Result, ValType(data[ptr]))
				ptr++
			}
			wf.Type = append(wf.Type, t)
		} else {
			panic("Invalid type")
		}
	}
}
