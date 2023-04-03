package wasmfile

import (
	"bufio"
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

	// #### Write out Data

	// #### Write out Elem

	_, err = wr.WriteString(")\n")
	if err != nil {
		return err
	}

	err = wr.Flush()
	return err
}
