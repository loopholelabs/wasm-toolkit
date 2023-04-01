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
	lenBuffer := make([]byte, 10)
	l := binary.PutUvarint(lenBuffer, v)
	_, err := w.Write(lenBuffer[:l])
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
