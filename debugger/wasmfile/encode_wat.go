package wasmfile

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
)

func (wf *WasmFile) EncodeWat(w io.Writer) error {
	wr := bufio.NewWriter(w)

	_, err := wr.WriteString("(module\n")
	if err != nil {
		return err
	}

	// #### Write out Type
	for index, t := range wf.Type {
		params := ""
		results := ""

		if len(t.Param) > 0 {
			params = " (param"
			for _, p := range t.Param {
				params = params + " " + byteToValType[p]
			}
			params = params + ")"
		}

		if len(t.Result) > 0 {
			results = " (result"
			for _, p := range t.Result {
				results = results + " " + byteToValType[p]
			}
			results = results + ")"
		}

		// Encode it and send it out...
		tdata := fmt.Sprintf("    (type (func%s%s)) ;; type_id=%d\n", params, results, index)
		_, err = wr.WriteString(tdata)
		if err != nil {
			return err
		}

	}

	// TODO: Encode the sections as wat

	// #### Write out Export

	// #### Write out Import

	// #### Write out Global

	// #### Write out Memory

	// #### Write out Table

	// #### Write out Function/Code
	for index, code := range wf.Code {
		function := wf.Function[index]
		tindex := function.TypeIndex
		typedata := wf.Type[tindex]

		params := ""
		results := ""

		if len(typedata.Param) > 0 {
			params = " (param"
			for _, p := range typedata.Param {
				params = params + " " + byteToValType[p]
			}
			params = params + ")"
		}

		if len(typedata.Result) > 0 {
			results = " (result"
			for _, p := range typedata.Result {
				results = results + " " + byteToValType[p]
			}
			results = results + ")"
		}

		f := wf.GetFunctionIdentifier(index + len(wf.Import))

		// Encode it and send it out...
		// TODO: Function identifier
		tdata := fmt.Sprintf("\n    (func %s (type %d)%s%s    ;; function_index=%d\n", f, tindex, params, results, index)
		_, err = wr.WriteString(tdata)
		if err != nil {
			return err
		}

		// Write out locals...
		for _, l := range code.Locals {
			_, err = wr.WriteString(fmt.Sprintf("        (local %s)\n", byteToValType[l]))
			if err != nil {
				return err
			}
		}

		var buf bytes.Buffer
		for _, e := range code.Expression {
			err = e.EncodeWat(&buf, "        ", wf)
			if err != nil {
				return err
			}
		}

		_, err = wr.Write(buf.Bytes())
		if err != nil {
			return err
		}

		_, err = wr.WriteString("    )\n")
		if err != nil {
			return err
		}

	}

	// #### Write out Data

	// #### Write out Elem

	_, err = wr.WriteString(")\n")
	if err != nil {
		return err
	}

	err = wr.Flush()
	return err
}
