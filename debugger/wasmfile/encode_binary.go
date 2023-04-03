package wasmfile

import (
	"bytes"
	"encoding/binary"
	"io"
)

func writeSectionHeader(w io.Writer, s byte, length int) error {
	sectionHeadBuffer := make([]byte, 10)
	sectionHeadBuffer[0] = s
	l := binary.PutUvarint(sectionHeadBuffer[1:], uint64(length))
	_, err := w.Write(sectionHeadBuffer[:l+1])
	return err
}

func writeString(w io.Writer, s string) error {
	data := []byte(s)
	err := writeUvarint(w, uint64(len(data)))
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}

func writeUvarint(w io.Writer, v uint64) error {
	b := binary.AppendUvarint(make([]byte, 0), v)
	_, err := w.Write(b)
	return err
}

func writeVarint(w io.Writer, v int64) error {
	b := AppendSleb128(make([]byte, 0), v)
	_, err := w.Write(b)
	return err
}

func (wf *WasmFile) EncodeBinary(w io.Writer) error {
	header := make([]byte, 8)
	binary.LittleEndian.PutUint32(header, WasmHeader)
	binary.LittleEndian.PutUint32(header[4:], WasmVersion)
	_, err := w.Write(header)
	if err != nil {
		return err
	}

	// Section Type
	if len(wf.Type) > 0 {
		var buf bytes.Buffer
		writeUvarint(&buf, uint64(len(wf.Type)))
		for _, t := range wf.Type {
			err = t.EncodeBinary(&buf)
			if err != nil {
				return err
			}
		}

		// Write a single type section
		writeSectionHeader(w, byte(SectionType), buf.Len())
		_, err = w.Write(buf.Bytes())
		if err != nil {
			return err
		}
	}

	// Section Import
	if len(wf.Import) > 0 {
		var buf bytes.Buffer
		writeUvarint(&buf, uint64(len(wf.Import)))
		for _, i := range wf.Import {
			err = i.EncodeBinary(&buf)
			if err != nil {
				return err
			}
		}

		// Write a single import section
		writeSectionHeader(w, byte(SectionImport), buf.Len())
		_, err = w.Write(buf.Bytes())
		if err != nil {
			return err
		}
	}

	// Section Function
	if len(wf.Function) > 0 {
		var buf bytes.Buffer
		writeUvarint(&buf, uint64(len(wf.Function)))
		for _, f := range wf.Function {
			err = f.EncodeBinary(&buf)
			if err != nil {
				return err
			}
		}

		// Write a single function section
		writeSectionHeader(w, byte(SectionFunction), buf.Len())
		_, err = w.Write(buf.Bytes())
		if err != nil {
			return err
		}
	}

	// Section Table
	if len(wf.Table) > 0 {
		var buf bytes.Buffer
		writeUvarint(&buf, uint64(len(wf.Table)))
		for _, t := range wf.Table {
			err = t.EncodeBinary(&buf)
			if err != nil {
				return err
			}
		}

		// Write a single table section
		writeSectionHeader(w, byte(SectionTable), buf.Len())
		_, err = w.Write(buf.Bytes())
		if err != nil {
			return err
		}
	}

	// Section Memory
	if len(wf.Memory) > 0 {
		var buf bytes.Buffer
		writeUvarint(&buf, uint64(len(wf.Memory)))
		for _, t := range wf.Memory {
			err = t.EncodeBinary(&buf)
			if err != nil {
				return err
			}
		}

		// Write a single memory section
		writeSectionHeader(w, byte(SectionMemory), buf.Len())
		_, err = w.Write(buf.Bytes())
		if err != nil {
			return err
		}
	}

	// Section Global
	if len(wf.Global) > 0 {
		var buf bytes.Buffer
		writeUvarint(&buf, uint64(len(wf.Global)))
		for _, t := range wf.Global {
			err = t.EncodeBinary(&buf)
			if err != nil {
				return err
			}
		}

		// Write a single memory section
		writeSectionHeader(w, byte(SectionGlobal), buf.Len())
		_, err = w.Write(buf.Bytes())
		if err != nil {
			return err
		}
	}

	// Section Export
	if len(wf.Export) > 0 {
		var buf bytes.Buffer
		writeUvarint(&buf, uint64(len(wf.Export)))
		for _, t := range wf.Export {
			err = t.EncodeBinary(&buf)
			if err != nil {
				return err
			}
		}

		// Write a single export section
		writeSectionHeader(w, byte(SectionExport), buf.Len())
		_, err = w.Write(buf.Bytes())
		if err != nil {
			return err
		}
	}

	// TODO StartSection

	// Section Elem
	if len(wf.Elem) > 0 {
		var buf bytes.Buffer
		writeUvarint(&buf, uint64(len(wf.Elem)))
		for _, t := range wf.Elem {
			err = t.EncodeBinary(&buf)
			if err != nil {
				return err
			}
		}

		// Write a single elem section
		writeSectionHeader(w, byte(SectionElem), buf.Len())
		_, err = w.Write(buf.Bytes())
		if err != nil {
			return err
		}
	}

	// Section DataCount
	var buf bytes.Buffer
	writeUvarint(&buf, uint64(len(wf.Data)))

	// Write a single data count section
	writeSectionHeader(w, byte(SectionDataCount), buf.Len())
	_, err = w.Write(buf.Bytes())
	if err != nil {
		return err
	}

	// Section Code
	if len(wf.Code) > 0 {
		var buf bytes.Buffer
		writeUvarint(&buf, uint64(len(wf.Code)))
		for _, c := range wf.Code {
			err = c.EncodeBinary(&buf)
			if err != nil {
				return err
			}
		}

		// Write a single code section
		writeSectionHeader(w, byte(SectionCode), buf.Len())
		_, err = w.Write(buf.Bytes())
		if err != nil {
			return err
		}
	}

	// Section Data
	if len(wf.Data) > 0 {
		var buf bytes.Buffer
		writeUvarint(&buf, uint64(len(wf.Data)))
		for _, t := range wf.Data {
			err = t.EncodeBinary(&buf)
			if err != nil {
				return err
			}
		}

		// Write a single data section
		writeSectionHeader(w, byte(SectionData), buf.Len())
		_, err = w.Write(buf.Bytes())
		if err != nil {
			return err
		}
	}

	// Section Custom
	if len(wf.Custom) > 0 {
		for _, c := range wf.Custom {
			var buf bytes.Buffer
			// Write the name, and the data...
			writeString(&buf, c.Name)
			// Now write the data into &buf
			_, err := buf.Write(c.Data)
			if err != nil {
				return err
			}

			// Write a single type section
			writeSectionHeader(w, byte(SectionCustom), buf.Len())
			_, err = w.Write(buf.Bytes())
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (ie *ImportEntry) EncodeBinary(w io.Writer) error {
	err := writeString(w, ie.Module)
	if err != nil {
		return err
	}
	err = writeString(w, ie.Name)
	if err != nil {
		return err
	}
	exportType := make([]byte, 1)
	exportType[0] = byte(ie.Type)
	_, err = w.Write(exportType)
	if err != nil {
		return err
	}

	return writeUvarint(w, uint64(ie.Index))
}

func (te *TypeEntry) EncodeBinary(w io.Writer) error {
	// Write out funcTypePrefix (only thing it supports atm)
	typePrefix := make([]byte, 1)
	typePrefix[0] = FuncTypePrefix
	_, err := w.Write(typePrefix)
	if err != nil {
		return err
	}

	err = writeUvarint(w, uint64(len(te.Param)))
	if err != nil {
		return err
	}

	params := make([]byte, len(te.Param))
	for i, p := range te.Param {
		params[i] = byte(p)
	}
	_, err = w.Write(params)
	if err != nil {
		return err
	}

	err = writeUvarint(w, uint64(len(te.Result)))
	if err != nil {
		return err
	}

	results := make([]byte, len(te.Result))
	for i, p := range te.Result {
		results[i] = byte(p)
	}
	_, err = w.Write(results)

	return err
}

func (f *FunctionEntry) EncodeBinary(w io.Writer) error {
	err := writeUvarint(w, uint64(f.TypeIndex))
	return err
}

func (c *TableEntry) EncodeBinary(w io.Writer) error {
	var buf bytes.Buffer

	buf.WriteByte(c.TableType)
	if c.LimitMax == 0 { // TODO: Fixme
		buf.WriteByte(LimitTypeMin)
		writeUvarint(&buf, uint64(c.LimitMin))
	} else {
		buf.WriteByte(LimitTypeMinMax)
		writeUvarint(&buf, uint64(c.LimitMin))
		writeUvarint(&buf, uint64(c.LimitMax))

	}

	_, err := w.Write(buf.Bytes())
	return err
}

func (c *MemoryEntry) EncodeBinary(w io.Writer) error {
	var buf bytes.Buffer

	if c.LimitMax == 0 { // TODO: Fixme
		buf.WriteByte(LimitTypeMin)
		writeUvarint(&buf, uint64(c.LimitMin))
	} else {
		buf.WriteByte(LimitTypeMinMax)
		writeUvarint(&buf, uint64(c.LimitMin))
		writeUvarint(&buf, uint64(c.LimitMax))

	}

	_, err := w.Write(buf.Bytes())
	return err
}

func (c *GlobalEntry) EncodeBinary(w io.Writer) error {
	var buf bytes.Buffer

	buf.WriteByte(byte(c.Type))
	buf.WriteByte(c.Mut)

	for _, e := range c.Expression {
		e.EncodeBinary(&buf)
	}
	buf.WriteByte(0x0b) // END

	_, err := w.Write(buf.Bytes())
	return err
}

func (c *ExportEntry) EncodeBinary(w io.Writer) error {
	var buf bytes.Buffer

	writeString(&buf, c.Name)
	buf.WriteByte(byte(c.Type))
	writeUvarint(&buf, uint64(c.Index))

	_, err := w.Write(buf.Bytes())
	return err
}

func (c *CodeEntry) EncodeBinary(w io.Writer) error {
	var buf bytes.Buffer

	writeUvarint(&buf, uint64(len(c.Locals)))
	for _, l := range c.Locals {
		writeUvarint(&buf, 1)
		buf.WriteByte(byte(l))
	}

	for _, e := range c.Expression {
		e.EncodeBinary(&buf)
	}
	buf.WriteByte(0x0b) // END

	err := writeUvarint(w, uint64(buf.Len()))
	if err != nil {
		return err
	}
	_, err = w.Write(buf.Bytes())
	return err
}

func (c *ElemEntry) EncodeBinary(w io.Writer) error {
	var buf bytes.Buffer

	err := writeUvarint(&buf, uint64(c.TableIndex))
	if err != nil {
		return err
	}

	for _, e := range c.Offset {
		e.EncodeBinary(&buf)
	}
	buf.WriteByte(0x0b) // END

	_, err = w.Write(buf.Bytes())
	if err != nil {
		return err
	}

	err = writeUvarint(w, uint64(len(c.Indexes)))
	if err != nil {
		return err
	}
	for _, ii := range c.Indexes {
		err = writeUvarint(w, ii)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *DataEntry) EncodeBinary(w io.Writer) error {
	var buf bytes.Buffer

	err := writeUvarint(&buf, uint64(c.MemIndex))
	if err != nil {
		return err
	}

	for _, e := range c.Offset {
		e.EncodeBinary(&buf)
	}
	buf.WriteByte(0x0b) // END

	_, err = w.Write(buf.Bytes())
	if err != nil {
		return err
	}

	err = writeUvarint(w, uint64(len(c.Data)))
	if err != nil {
		return err
	}
	_, err = w.Write(c.Data)
	return err
}
