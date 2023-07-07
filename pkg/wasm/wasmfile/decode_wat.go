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
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/loopholelabs/wasm-toolkit/pkg/wasm/debug"
	"github.com/loopholelabs/wasm-toolkit/pkg/wasm/encoding"
	"github.com/loopholelabs/wasm-toolkit/pkg/wasm/expression"
	"github.com/loopholelabs/wasm-toolkit/pkg/wasm/types"
)

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

func (wf *WasmFile) RegisterNextFunctionName(n string) {
	idx := len(wf.Debug.FunctionNames)
	wf.Debug.FunctionNames[idx] = n
}

func (wf *WasmFile) RegisterNextGlobalName(n string) {
	idx := len(wf.Debug.GlobalNames)
	wf.Debug.GlobalNames[idx] = n
}

func (wf *WasmFile) RegisterNextDataName(n string) {
	idx := len(wf.Debug.DataNames)
	wf.Debug.DataNames[idx] = n
}

func (wf *WasmFile) DecodeWat(data []byte) (err error) {
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
	// Parse the wat file and fill in all the data...
	wf.Debug = &debug.WasmDebug{}
	wf.Debug.FunctionNames = make(map[int]string)
	wf.Debug.GlobalNames = make(map[int]string)
	wf.Debug.DataNames = make(map[int]string)

	text := string(data)

	// Read the module
	moduleText, _ := encoding.ReadElement(text)

	moduleType, _ := encoding.ReadToken(moduleText[1:])

	if moduleType != "module" {
		return errors.New("Invalud module. Expected 'module'")
	}

	// Now read all the individual elements from within the module...

	text = text[len(moduleType)+1:]
	bodytext := text // Save it for a second pass

	for {
		text = strings.TrimLeft(text, encoding.Whitespace) // Skip to next bit
		// End of the module?
		if text[0] == ')' {
			break
		}

		// Skip any single line comments
		for {
			if strings.HasPrefix(text, ";;") {
				// Skip to end of line
				p := strings.Index(text, "\n")
				if p == -1 {
					panic("TODO: Comment without newline")
				}
				text = text[p+1:]
				text = strings.TrimLeft(text, encoding.Whitespace) // Skip to next bit
			} else {
				break
			}
		}

		e, _ := encoding.ReadElement(text)
		eType, _ := encoding.ReadToken(e[1:])

		if eType == "data" {
			de := &DataEntry{}
			err = de.DecodeWat(e, wf)
			wf.Data = append(wf.Data, de)
		} else if eType == "elem" {
			ee := &ElemEntry{}
			err = ee.DecodeWat(e, wf)
			wf.Elem = append(wf.Elem, ee)
		} else if eType == "func" {
			ee := &FunctionEntry{}
			err = ee.DecodeWat(e, wf)
			wf.Function = append(wf.Function, ee)
		} else if eType == "global" {
			ge := &GlobalEntry{}
			err = ge.DecodeWat(e, wf)
			wf.Global = append(wf.Global, ge)
		} else if eType == "import" {
			ie := &ImportEntry{}
			err = ie.DecodeWat(e, wf)
			wf.Import = append(wf.Import, ie)
		} else if eType == "memory" {
			ee := &MemoryEntry{}
			err = ee.DecodeWat(e)
			wf.Memory = append(wf.Memory, ee)
		} else if eType == "table" {
			ee := &TableEntry{}
			err = ee.DecodeWat(e)
			wf.Table = append(wf.Table, ee)
		} else if eType == "type" {
			ee := &TypeEntry{}
			err = ee.DecodeWat(e)
			wf.Type = append(wf.Type, ee)
		} else if eType == "export" {
			// Deal with it in 2nd pass
		} else {
			panic(fmt.Sprintf("Unknown element \"%s\"", eType))
		}
		if err != nil {
			return err
		}

		// Skip over this element
		text = text[len(e):]
	}

	// Second pass
	text = bodytext

	for {
		text = strings.TrimLeft(text, encoding.Whitespace) // Skip to next bit
		// End of the module?
		if text[0] == ')' {
			break
		}

		// Skip any single line comments
		for {
			if strings.HasPrefix(text, ";;") {
				// Skip to end of line
				p := strings.Index(text, "\n")
				if p == -1 {
					panic("TODO: Comment without newline")
				}
				text = text[p+1:]
				text = strings.TrimLeft(text, encoding.Whitespace) // Skip to next bit
			} else {
				break
			}
		}

		e, _ := encoding.ReadElement(text)
		eType, _ := encoding.ReadToken(e[1:])

		if eType == "export" {
			ee := &ExportEntry{}
			err = ee.DecodeWat(e, wf)
			wf.Export = append(wf.Export, ee)
		} else if eType == "func" {
			ce := &CodeEntry{}
			err = ce.DecodeWat(e, wf)
			wf.Code = append(wf.Code, ce)
		}
		if err != nil {
			return err
		}

		// Skip over this element
		text = text[len(e):]
	}

	return nil
}

func (e *TypeEntry) DecodeWat(d string) error {
	//   (type (;0;) (func (param i32 i32 i32 i32) (result i32)))

	s := strings.Trim(d[5:len(d)-1], encoding.Whitespace)
	s = encoding.SkipComment(s)
	fspec, s := encoding.ReadElement(s)
	if fspec == "(func)" {
		// Special case, nothing else to do.
		return nil
	}
	if strings.HasPrefix(fspec, "(func ") && fspec[len(fspec)-1] == ')' {
		fspec = fspec[6 : len(fspec)-1]
		for {
			var el string
			fspec = encoding.SkipComment(fspec)
			fspec = strings.Trim(fspec, encoding.Whitespace)
			if len(fspec) == 0 {
				break
			}
			el, fspec = encoding.ReadElement(fspec)
			if strings.HasPrefix(el, "(param ") {
				// Now read each type
				el = el[7 : len(el)-1]
				for {
					var ptype string
					el = encoding.SkipComment(el)
					el = strings.Trim(el, encoding.Whitespace)
					if len(el) == 0 {
						break
					}
					ptype, el = encoding.ReadToken(el)
					b, ok := types.ValTypeToByte[ptype]
					if !ok {
						return fmt.Errorf("Unknown param type (%s)", ptype)
					}
					e.Param = append(e.Param, b)
				}
			} else if strings.HasPrefix(el, "(result ") {
				// atm we only support one return type.
				rtype := strings.Trim(el[8:len(el)-1], encoding.Whitespace)
				b, ok := types.ValTypeToByte[rtype]
				if !ok {
					return fmt.Errorf("Unknown result type (%s)", rtype)
				}
				e.Result = append(e.Result, b)
			} else {
				return fmt.Errorf("Unknown spec in type %s", el)
			}
		}
	} else {
		return errors.New("Only support type func atm")
	}
	return nil
}

func (e *TableEntry) DecodeWat(d string) error {
	//  (table (;0;) 3 3 funcref)

	s := strings.Trim(d[6:len(d)-1], encoding.Whitespace)
	s = encoding.SkipComment(s)
	// Should be a number next (min)
	var mmin string
	var mmax string
	var err error
	mmin, s = encoding.ReadToken(s)
	e.LimitMin, err = strconv.Atoi(mmin)
	if err != nil {
		return err
	}

	s = encoding.SkipComment(s)
	s = strings.Trim(s, encoding.Whitespace)
	mmax, s = encoding.ReadToken(s)
	e.LimitMax, err = strconv.Atoi(mmax)
	if err != nil {
		return err
	}

	tabtype, s := encoding.ReadToken(s)
	if tabtype != "funcref" {
		return errors.New("Only table funcref supported atm")
	}
	e.TableType = types.TableTypeFuncref
	return nil
}

func (e *MemoryEntry) DecodeWat(d string) error {
	// (memory (;0;) 2)

	s := strings.Trim(d[7:len(d)-1], encoding.Whitespace)
	s = encoding.SkipComment(s)
	// Should be a number next (min)
	var mmin string
	var mmax string
	var err error
	mmin, s = encoding.ReadToken(s)
	e.LimitMin, err = strconv.Atoi(mmin)
	if err != nil {
		return err
	}

	s = encoding.SkipComment(s)
	s = strings.Trim(s, encoding.Whitespace)
	if len(s) > 0 {
		mmax, s = encoding.ReadToken(s)
		e.LimitMax, err = strconv.Atoi(mmax)
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *ImportEntry) DecodeWat(d string, wf *WasmFile) error {
	//  (import "wasi_snapshot_preview1" "fd_write" (func $runtime.fd_write (type 0)))
	var err error

	s := strings.TrimLeft(d[7:len(d)-1], encoding.Whitespace)

	e.Module, s = encoding.ReadString(s)
	e.Module = e.Module[1 : len(e.Module)-1]
	e.Name, s = encoding.ReadString(s)
	e.Name = e.Name[1 : len(e.Name)-1]

	var idata, typedata, tdata string
	idata, s = encoding.ReadElement(s)
	iType, _ := encoding.ReadToken(idata[1:])
	if iType == "func" {
		idata = strings.Trim(idata[5:len(idata)-1], encoding.Whitespace)
		// Read the (optional) function name ID
		if idata[0] != '(' {
			var fname string
			fname, idata = encoding.ReadToken(idata)
			idata = strings.Trim(idata, encoding.Whitespace)
			wf.RegisterNextFunctionName(fname)
		}
		// Now read the type...
		typedata, _ = encoding.ReadElement(idata)
		tdata, s = encoding.ReadToken(typedata[1:])
		if tdata == "type" {
			typedata = strings.Trim(typedata[5:len(typedata)-1], encoding.Whitespace)
			// Read the value
			e.Index, err = strconv.Atoi(typedata)
			if err != nil {
				return err
			}
		} else {
			panic("Issue parsing import func")
		}

	} else {
		panic("TODO: Import other than func")
	}

	return nil
}

func (e *GlobalEntry) DecodeWat(d string, wf *WasmFile) error {
	//  (global $__stack_pointer (mut i32) (i32.const 65536))

	s := strings.Trim(d[7:len(d)-1], encoding.Whitespace)
	if s[0] == '$' {
		// We have an identifier, lets use it
		var id string
		id, s = encoding.ReadToken(s)
		wf.RegisterNextGlobalName(id)
	}

	// Next we either have a type, or (mut <type>)
	var ty string
	var ok bool
	if s[0] == '(' {
		var mutty string
		mutty, s = encoding.ReadElement(s)
		if strings.HasPrefix(mutty, "(mut ") && mutty[len(mutty)-1] == ')' {
			e.Mut = 1
			ty = mutty[5 : len(mutty)-1]
		} else {
			return fmt.Errorf("Cannot parse global %s", d)
		}
	} else {
		ty, s = encoding.ReadToken(s)
		e.Mut = 0
	}

	e.Type, ok = types.ValTypeToByte[ty]
	if !ok {
		return fmt.Errorf("Invalid type in global %s", ty)
	}

	s = strings.Trim(s, encoding.Whitespace)
	expr, _ := encoding.ReadElement(s)
	// Read the expression
	expr = expr[1 : len(expr)-1]
	// TODO: Support proper expressions. For now we only support a single instruction
	e.Expression = make([]*expression.Expression, 0)
	ex := &expression.Expression{}
	err := ex.DecodeWat(expr, nil)
	if err != nil {
		return err
	}
	e.Expression = append(e.Expression, ex)

	return nil
}

func (e *CodeEntry) DecodeWat(d string, wf *WasmFile) error {
	e.Locals = make([]types.ValType, 0)

	s := strings.Trim(d[5:len(d)-1], encoding.Whitespace)

	// Optional Identifier
	if s[0] == '$' {
		_, s = encoding.ReadToken(s)
	}

	localNames := make(map[string]int)
	localIndex := 0

	// FIXME: If the func only has (type) and not explicit (param) (local)

	for {
		// Skip comments...

		s = strings.Trim(s, encoding.Whitespace)
		if len(s) == 0 {
			break
		}

		if strings.HasPrefix(s, ";;") {
			// Skip this line
			line_end := strings.Index(s, "\n")
			if line_end == -1 {
				return nil // All done?
			}
			s = s[line_end:]
		} else if s[0] == '(' {
			var el string
			el, s = encoding.ReadElement(s)
			eType, _ := encoding.ReadToken(el[1:])
			if eType == "type" {
			} else if eType == "param" {
				// Might have a name here...
				el = strings.Trim(el[6:len(el)-1], encoding.Whitespace)
				if el[0] == '$' {
					var name string
					name, el = encoding.ReadToken(el)
					localNames[name] = localIndex
				}

				// Read the tokens one by one for each param
				for {
					el = strings.Trim(el, encoding.Whitespace)
					if len(el) == 0 {
						break
					}
					_, el = encoding.ReadToken(el)
					localIndex++
				}

				// TODO: Use to sanity check
			} else if eType == "result" {
				// TODO: Use to sanity check
			} else if eType == "local" {
				// eg (local $hello i32)
				// eg (local i64 i64)

				ltypes := strings.Trim(el[6:len(el)-1], encoding.Whitespace)
				for {
					ltypes = strings.Trim(ltypes, encoding.Whitespace)
					if len(ltypes) == 0 {
						break
					}
					var tok string
					tok, ltypes = encoding.ReadToken(ltypes)
					if tok[0] == '$' {
						// preRegister a name
						localNames[tok] = localIndex
					} else {
						l, ok := types.ValTypeToByte[tok]
						if !ok {
							return fmt.Errorf("Invalid local type %s", tok)
						}
						e.Locals = append(e.Locals, l)
						localIndex++
					}
				}
			}
		} else {
			break
		}
	}

	// Then just read instructions...
	for {
		s = strings.Trim(s, encoding.Whitespace)
		if len(s) == 0 {
			break
		}
		lend := strings.Index(s, "\n")
		if lend == -1 {
			lend = len(s)
		}
		ecode := s[:lend]
		s = s[lend:]

		// Ignore any ;; comments
		cend := strings.Index(ecode, ";;")
		if cend != -1 {
			ecode = ecode[:cend]
		}

		ecode = strings.Trim(ecode, encoding.Whitespace)

		if len(ecode) > 0 {
			newe := &expression.Expression{}
			err := newe.DecodeWat(ecode, localNames)
			if err != nil {
				return err
			}
			e.Expression = append(e.Expression, newe)
		}
	}

	return nil
}

func (e *FunctionEntry) DecodeWat(d string, wf *WasmFile) error {
	s := strings.TrimLeft(d[5:len(d)-1], encoding.Whitespace)
	// eg (func $write (type 7) (param i32 i32 i32) (result i32)

	// Optional Identifier
	if s[0] == '$' {
		var fname string
		fname, s = encoding.ReadToken(s)
		// Store the name for lookups...
		wf.RegisterNextFunctionName(fname)
	}

	newTypeEntry := &TypeEntry{}

	for {
		s = strings.Trim(s, encoding.Whitespace)
		if len(s) == 0 {
			break
		}

		if s[0] == '(' {
			var el string
			var err error
			el, s = encoding.ReadElement(s)
			if strings.HasPrefix(el, "(type ") {
				el = strings.Trim(el[5:len(el)-1], encoding.Whitespace)
				e.TypeIndex, err = strconv.Atoi(el)
				return err
			} else if strings.HasPrefix(el, "(param ") {
				// Now read each type
				el = el[7 : len(el)-1]
				// Could be a name here...
				if el[0] == '$' {
					_, el = encoding.ReadToken(el)
				}
				for {
					var ptype string
					el = encoding.SkipComment(el)
					el = strings.Trim(el, encoding.Whitespace)
					if len(el) == 0 {
						break
					}
					ptype, el = encoding.ReadToken(el)
					b, ok := types.ValTypeToByte[ptype]
					if !ok {
						return fmt.Errorf("Unknown param type (%s)", ptype)
					}
					newTypeEntry.Param = append(newTypeEntry.Param, b)
				}
			} else if strings.HasPrefix(el, "(result ") {
				// atm we only support one return type.
				rtype := strings.Trim(el[8:len(el)-1], encoding.Whitespace)
				b, ok := types.ValTypeToByte[rtype]
				if !ok {
					return fmt.Errorf("Unknown result type (%s)", rtype)
				}
				newTypeEntry.Result = append(newTypeEntry.Result, b)
			}
		} else {
			break
		}
	}

	// Add a new Type for this function, or use an existing one...
	for idx, t := range wf.Type {
		// Check if we can use it or not...
		if t.Equals(newTypeEntry) {
			e.TypeIndex = idx
			return nil
		}
	}

	// Need to add a new type
	e.TypeIndex = len(wf.Type)
	wf.Type = append(wf.Type, newTypeEntry)
	return nil

}

func (e *ExportEntry) DecodeWat(d string, wf *WasmFile) error {
	//  (export "memory" (memory 0))
	//  (export "hello" (func $hello))

	s := strings.TrimLeft(d[7:len(d)-1], encoding.Whitespace)

	e.Name, s = encoding.ReadString(s)
	e.Name = e.Name[1 : len(e.Name)-1]
	s = strings.Trim(s, encoding.Whitespace)
	el, _ := encoding.ReadElement(s)
	etype, erest := encoding.ReadToken(el[1:])
	erest = erest[:len(erest)-1]
	if etype == "memory" {
		e.Type = types.ExportMem
		idx, err := strconv.Atoi(erest)
		if err != nil {
			return err
		}
		e.Index = idx
	} else if etype == "func" {
		e.Type = types.ExportFunc
		if strings.HasPrefix(erest, "$") {
			fname, _ := encoding.ReadToken(erest)
			fid := wf.Debug.LookupFunctionID(fname)
			if fid == -1 {
				return fmt.Errorf("Function %s not found in export", fname)
			}
			e.Index = fid
			// Parse the ID and look it up
		} else {
			idx, err := strconv.Atoi(erest)
			if err != nil {
				return err
			}
			e.Index = idx

		}
	} else {
		return errors.New("TODO: Support other exports")
	}

	return nil
}

func (e *ElemEntry) DecodeWat(d string, wf *WasmFile) error {
	// (elem (;0;) (i32.const 1) func $runtime.memequal $runtime.hash32)
	e.TableIndex = 0 // For now only one table

	s := strings.Trim(d[5:len(d)-1], encoding.Whitespace)
	s = encoding.SkipComment(s)

	expr, s := encoding.ReadElement(s)
	// Read the expression
	expr = expr[1 : len(expr)-1]
	// TODO: Support proper expressions. For now we only support a single instruction
	e.Offset = make([]*expression.Expression, 0)
	ex := &expression.Expression{}
	err := ex.DecodeWat(expr, nil)
	if err != nil {
		return err
	}
	e.Offset = append(e.Offset, ex)

	s = strings.Trim(s, encoding.Whitespace)
	var elemType string
	elemType, s = encoding.ReadToken(s)
	if elemType == "func" {
		for {
			s = strings.Trim(s, encoding.Whitespace)
			if len(s) == 0 {
				break
			}
			var fid string
			var findex int
			fid, s = encoding.ReadToken(s)
			if strings.HasPrefix(fid, "$") {
				findex = wf.Debug.LookupFunctionID(fid)
				if findex == -1 {
					return fmt.Errorf("Function not found %s", fid)
				}
			} else {
				findex, err = strconv.Atoi(fid)
				if err != nil {
					return err
				}
			}
			e.Indexes = append(e.Indexes, uint64(findex))
		}
	} else {
		return fmt.Errorf("Unknown type for elem %s", elemType)
	}

	return nil
}

func (e *DataEntry) DecodeWat(d string, wf *WasmFile) error {
	//	* (data $.data (i32.const 66160) "x\9c\19\f6\dc\02\01\00\00\00\00\00\9c\03\01\00\c1\82\01\00\00\00\00\00\04\00\00\00\0c\00\00\00\01\00\00\00\00\00\00\00\01\00\00\00\00\00\00\00\02\00\00\00\a8\02\01\00\98\01\00\00\01\00\00\00\ff\01\01\00\0b\00\00\00\00\00\00\00 \01\01\00\13\00\00\003\01\01\00\13"))
	//	* (data $.data 10)
	//	* (data $.data "hello world")

	s := strings.Trim(d[5:len(d)-1], encoding.Whitespace)
	var id string

	if s[0] == '$' {
		id, s = encoding.ReadToken(s)
		wf.RegisterNextDataName(id)
	}
	s = strings.Trim(s, encoding.Whitespace)
	if s[0] == '(' {
		// Must have a specific Offset already set
		var expr string
		expr, s = encoding.ReadElement(s)
		// Read the expression
		expr = expr[1 : len(expr)-1]
		// TODO: Support proper expressions. For now we only support a single instruction
		e.Offset = make([]*expression.Expression, 0)
		ex := &expression.Expression{}
		err := ex.DecodeWat(expr, nil)
		if err != nil {
			return err
		}
		e.Offset = append(e.Offset, ex)
	} else {
		// Assume this data should go right after the last bit of data... (Aligned)
		data_ptr := int32(0)
		if len(wf.Data) > 0 {
			prev := wf.Data[len(wf.Data)-1]
			data_ptr = prev.Offset[0].I32Value + int32(len(prev.Data))
		}

		// Align it...
		data_ptr = (data_ptr + 3) & -4

		e.Offset = []*expression.Expression{
			{
				Opcode:   expression.InstrToOpcode["i32.const"],
				I32Value: data_ptr,
			},
		}
	}

	s = strings.Trim(s, encoding.Whitespace)

	if s[0] == '"' {
		// Parse the data
		s = s[1 : len(s)-1]
		for {
			if len(s) == 0 {
				break
			}
			// TODO: \r\n\t
			if s[0] == '\\' {
				// Parse the byte value...
				bval := s[1:3]
				bv, err := strconv.ParseInt(bval, 16, 32)
				if err != nil {
					return err
				}
				e.Data = append(e.Data, byte(bv))
				s = s[3:]
			} else {
				e.Data = append(e.Data, byte(s[0]))
				s = s[1:]
			}
		}
	} else {
		// Assume it's a number...
		length, err := strconv.Atoi(s)
		if err != nil {
			return err
		}
		e.Data = make([]byte, length)
	}

	return nil
}
