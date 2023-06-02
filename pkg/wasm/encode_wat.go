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
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/loopholelabs/wasm-toolkit/pkg/wasm/types"
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
				params = params + " " + types.ByteToValType[p]
			}
			params = params + ")"
		}

		if len(t.Result) > 0 {
			results = " (result"
			for _, p := range t.Result {
				results = results + " " + types.ByteToValType[p]
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

	// #### Write out Import
	for index, t := range wf.Import {
		exp := ""
		if t.Type == types.ExportFunc {
			exp = fmt.Sprintf("(func %s (type %d))", wf.Debug.GetFunctionIdentifier(index, true), t.Index)
		} else if t.Type == types.ExportGlobal {
			exp = fmt.Sprintf("(global %d)", t.Index)
		} else if t.Type == types.ExportMem {
			exp = fmt.Sprintf("(memory %d)", t.Index)
		} else if t.Type == types.ExportTable {
			exp = fmt.Sprintf("(table %d)", t.Index)
		}

		edata := fmt.Sprintf("    (import \"%s\" \"%s\" %s)\n", t.Module, t.Name, exp)
		_, err = wr.WriteString(edata)
		if err != nil {
			return err
		}
	}

	// #### Write out Global
	for index, g := range wf.Global {
		t := types.ByteToValType[g.Type]
		if g.Mut == 0x01 {
			t = fmt.Sprintf("(mut %s)", t)
		}

		var buf bytes.Buffer
		for _, ee := range g.Expression {
			err := ee.EncodeWat(&buf, "", wf, wf.Debug)
			if err != nil {
				return err
			}
		}

		gname := wf.Debug.GetGlobalIdentifier(index, true)

		edata := fmt.Sprintf("    (global %s %s (%s))\n", gname, t, buf.Bytes())
		_, err = wr.WriteString(edata)
		if err != nil {
			return err
		}
	}

	// #### Write out Memory
	for _, m := range wf.Memory {
		limits := fmt.Sprintf("%d", m.LimitMin)
		if m.LimitMax != 0 {
			limits = fmt.Sprintf("%s %d", limits, m.LimitMax)
		}

		mdata := fmt.Sprintf("    (memory %s)\n", limits)
		_, err = wr.WriteString(mdata)
		if err != nil {
			return err
		}
	}

	// #### Write out Table
	for _, t := range wf.Table {
		limits := fmt.Sprintf("%d", t.LimitMin)
		if t.LimitMax != 0 {
			limits = fmt.Sprintf("%s %d", limits, t.LimitMax)
		}

		tabType := "funcref"

		mdata := fmt.Sprintf("    (table %s %s)\n", limits, tabType)
		_, err = wr.WriteString(mdata)
		if err != nil {
			return err
		}
	}

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

				params = fmt.Sprintf("%s\n        (param %s)%s", params, types.ByteToValType[p], comment)
			}
		}

		if len(typedata.Result) > 0 {
			results = "        (result"
			for _, p := range typedata.Result {
				results = results + " " + types.ByteToValType[p]
			}
			results = results + ")\n"
		}

		f := wf.Debug.GetFunctionIdentifier(index+len(wf.Import), true)

		// Encode it and send it out...
		d := wf.GetFunctionDebug(index + len(wf.Import))
		tdata := fmt.Sprintf("\n    (func %s (type %d) ;; function_index=%d\n%s%s\n%s", f, tindex, index, d, params, results)
		_, err = wr.WriteString(tdata)
		if err != nil {
			return err
		}

		// Write out locals...
		for _, l := range code.Locals {
			_, err = wr.WriteString(fmt.Sprintf("        (local %s)\n", types.ByteToValType[l]))
			if err != nil {
				return err
			}
		}

		var buf bytes.Buffer
		for _, e := range code.Expression {
			err = e.EncodeWat(&buf, "        ", wf, wf.Debug)
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

		_, err = wr.WriteString(fmt.Sprintf("    )%s\n", comment))
		if err != nil {
			return err
		}

	}

	// #### Write out Export
	for _, t := range wf.Export {
		exp := ""
		if t.Type == types.ExportFunc {
			exp = fmt.Sprintf("(func %s)", wf.Debug.GetFunctionIdentifier(t.Index, false))
		} else if t.Type == types.ExportGlobal {
			exp = fmt.Sprintf("(global %d)", t.Index)
		} else if t.Type == types.ExportMem {
			exp = fmt.Sprintf("(memory %d)", t.Index)
		} else if t.Type == types.ExportTable {
			exp = fmt.Sprintf("(table %d)", t.Index)
		}

		edata := fmt.Sprintf("    (export \"%s\" %s)\n", t.Name, exp)
		_, err = wr.WriteString(edata)
		if err != nil {
			return err
		}
	}

	// #### Write out Data
	for index, d := range wf.Data {
		id := wf.Debug.GetDataIdentifier(index)

		var buf bytes.Buffer
		for _, ee := range d.Offset {
			err := ee.EncodeWat(&buf, "", wf, wf.Debug)
			if err != nil {
				return err
			}
		}

		dat := d.GetStringEncodedData()

		ddata := fmt.Sprintf("    (data %s (%s) \"%s\")\n", id, strings.Trim(buf.String(), " \t\r\n"), dat)
		_, err = wr.WriteString(ddata)
		if err != nil {
			return err
		}
	}

	// #### Write out Elem
	for _, e := range wf.Elem {

		var buf bytes.Buffer
		for _, ee := range e.Offset {
			err := ee.EncodeWat(&buf, "", wf, wf.Debug)
			if err != nil {
				return err
			}
		}

		funcs := ""
		for _, f := range e.Indexes {
			fid := wf.Debug.GetFunctionIdentifier(int(f), false)
			funcs = funcs + " " + fid
		}

		ddata := fmt.Sprintf("    (elem (%s) func%s)\n", strings.Trim(buf.String(), " \t\r\n"), funcs)
		_, err = wr.WriteString(ddata)
		if err != nil {
			return err
		}
	}

	_, err = wr.WriteString(")\n")
	if err != nil {
		return err
	}

	err = wr.Flush()
	return err
}

func (d *DataEntry) GetStringEncodedData() string {
	allowed := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ "
	var buf bytes.Buffer
	for _, b := range d.Data {
		if strings.Index(allowed, string(rune(b))) != -1 {
			buf.WriteRune(rune(b))
		} else {
			buf.WriteString(fmt.Sprintf("\\%02x", b))
		}
	}
	return buf.String()
}
