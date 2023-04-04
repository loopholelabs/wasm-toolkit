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
			for index, p := range typedata.Param {
				comment := ""
				vname := wf.GetLocalVarName(code.CodeSectionPtr, index)
				if vname != "" {
					comment = " ;; " + vname
				}

				params = fmt.Sprintf("%s\n        (param %s)%s", params, byteToValType[p], comment)
			}
		}

		if len(typedata.Result) > 0 {
			results = "        (result"
			for _, p := range typedata.Result {
				results = results + " " + byteToValType[p]
			}
			results = results + ")\n"
		}

		f := wf.GetFunctionIdentifier(index + len(wf.Import))

		// Encode it and send it out...
		// TODO: Function identifier
		d := wf.GetFunctionDebug(index + len(wf.Import))
		tdata := fmt.Sprintf("\n    (func %s (type %d) ;; function_index=%d\n%s%s\n%s", f, tindex, index, d, params, results)
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

		// Bit of a special case here. We know the function ends with an END opcode...
		lastAddr := code.CodeSectionPtr + code.CodeSectionLen - 1
		lineNumberData := wf.GetLineNumberInfo(lastAddr)
		comment := ""
		if lineNumberData != "" {
			comment = fmt.Sprintf(" ;; Src = %s", lineNumberData)
		}

		fmt.Printf("LineData %d %d %d %s\n", code.CodeSectionPtr, code.CodeSectionLen, lastAddr, lineNumberData)

		_, err = wr.WriteString(fmt.Sprintf("    )%s\n", comment))
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
