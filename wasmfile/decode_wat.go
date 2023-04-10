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
	"strconv"
	"strings"
)

func (wf *WasmFile) LookupFunctionID(n string) int {
	for idx, name := range wf.functionNames {
		if n == name {
			return idx
		}
	}
	return -1
}

func (wf *WasmFile) RegisterNextFunctionName(n string) {
	idx := len(wf.Function) + 1
	wf.functionNames[idx] = n
}

func (wf *WasmFile) RegisterNextGlobalName(n string) {
	idx := len(wf.Global) + 1
	wf.globalNames[idx] = n
}

func (wf *WasmFile) RegisterNextDataName(n string) {
	idx := len(wf.Data) + 1
	wf.dataNames[idx] = n
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

	wf.functionNames = make(map[int]string)
	wf.globalNames = make(map[int]string)
	wf.dataNames = make(map[int]string)

	text := string(data)

	// Read the module
	moduleText, _ := ReadElement(text)

	moduleType, _ := ReadToken(moduleText[1:])

	if moduleType != "module" {
		return errors.New("Invalud module. Expected 'module'")
	}

	// Now read all the individual elements from within the module...

	text = text[len(moduleType)+1:]
	bodytext := text // Save it for a second pass

	for {
		text = strings.TrimLeft(text, Whitespace) // Skip to next bit
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
				text = strings.TrimLeft(text, Whitespace) // Skip to next bit
			} else {
				break
			}
		}

		e, _ := ReadElement(text)
		eType, _ := ReadToken(e[1:])

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
			err = ie.DecodeWat(e)
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
		text = strings.TrimLeft(text, Whitespace) // Skip to next bit
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
				text = strings.TrimLeft(text, Whitespace) // Skip to next bit
			} else {
				break
			}
		}

		e, _ := ReadElement(text)
		eType, _ := ReadToken(e[1:])

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

	s := strings.Trim(d[5:len(d)-1], Whitespace)
	s = SkipComment(s)
	fspec, s := ReadElement(s)
	if fspec == "(func)" {
		// Special case, nothing else to do.
		return nil
	}
	if strings.HasPrefix(fspec, "(func ") && fspec[len(fspec)-1] == ')' {
		fspec = fspec[6 : len(fspec)-1]
		for {
			var el string
			fspec = SkipComment(fspec)
			fspec = strings.Trim(fspec, Whitespace)
			if len(fspec) == 0 {
				break
			}
			el, fspec = ReadElement(fspec)
			if strings.HasPrefix(el, "(param ") {
				// Now read each type
				el = el[7 : len(el)-1]
				for {
					var ptype string
					el = SkipComment(el)
					el = strings.Trim(el, Whitespace)
					if len(el) == 0 {
						break
					}
					ptype, el = ReadToken(el)
					b, ok := valTypeToByte[ptype]
					if !ok {
						return fmt.Errorf("Unknown param type (%s)", ptype)
					}
					e.Param = append(e.Param, b)
				}
			} else if strings.HasPrefix(el, "(result ") {
				// atm we only support one return type.
				rtype := strings.Trim(el[8:len(el)-1], Whitespace)
				b, ok := valTypeToByte[rtype]
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

	s := strings.Trim(d[6:len(d)-1], Whitespace)
	s = SkipComment(s)
	// Should be a number next (min)
	var mmin string
	var mmax string
	var err error
	mmin, s = ReadToken(s)
	e.LimitMin, err = strconv.Atoi(mmin)
	if err != nil {
		return err
	}

	s = SkipComment(s)
	s = strings.Trim(s, Whitespace)
	mmax, s = ReadToken(s)
	e.LimitMax, err = strconv.Atoi(mmax)
	if err != nil {
		return err
	}

	tabtype, s := ReadToken(s)
	if tabtype != "funcref" {
		return errors.New("Only table funcref supported atm")
	}
	return nil
}

func (e *MemoryEntry) DecodeWat(d string) error {
	// (memory (;0;) 2)

	s := strings.Trim(d[7:len(d)-1], Whitespace)
	s = SkipComment(s)
	// Should be a number next (min)
	var mmin string
	var mmax string
	var err error
	mmin, s = ReadToken(s)
	e.LimitMin, err = strconv.Atoi(mmin)
	if err != nil {
		return err
	}

	s = SkipComment(s)
	s = strings.Trim(s, Whitespace)
	if len(s) > 0 {
		mmax, s = ReadToken(s)
		e.LimitMax, err = strconv.Atoi(mmax)
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *ImportEntry) DecodeWat(d string) error {
	//  (import "wasi_snapshot_preview1" "fd_write" (func $runtime.fd_write (type 0)))
	var err error

	s := strings.TrimLeft(d[7:len(d)-1], Whitespace)

	e.Module, s = ReadString(s)
	e.Name, s = ReadString(s)

	var idata, typedata, tdata string
	idata, s = ReadElement(s)
	iType, _ := ReadToken(idata[1:])
	if iType == "func" {
		idata = strings.Trim(idata[5:len(idata)-1], Whitespace)
		// Read the (optional) function name ID
		if idata[0] != '(' {
			_, idata = ReadToken(idata)
			idata = strings.Trim(idata, Whitespace)
		}
		// Now read the type...
		typedata, _ = ReadElement(idata)
		tdata, s = ReadToken(typedata[1:])
		if tdata == "type" {
			typedata = strings.Trim(typedata[5:len(typedata)-1], Whitespace)
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

	s := strings.Trim(d[7:len(d)-1], Whitespace)
	if s[0] == '$' {
		// We have an identifier, lets use it
		var id string
		id, s = ReadToken(s)
		wf.RegisterNextGlobalName(id)
	}

	// Next we either have a type, or (mut <type>)
	var ty string
	var ok bool
	if s[0] == '(' {
		var mutty string
		mutty, s = ReadElement(s)
		if strings.HasPrefix(mutty, "(mut ") && mutty[len(mutty)-1] == ')' {
			e.Mut = 1
			ty = mutty[5 : len(mutty)-1]
		} else {
			return fmt.Errorf("Cannot parse global %s", d)
		}
	} else {
		ty, s = ReadToken(s)
		e.Mut = 0
	}

	e.Type, ok = valTypeToByte[ty]
	if !ok {
		return fmt.Errorf("Invalid type in global %s", ty)
	}

	s = strings.Trim(s, Whitespace)
	expr, _ := ReadElement(s)
	// Read the expression
	expr = expr[1 : len(expr)-1]
	// TODO: Support proper expressions. For now we only support a single instruction
	e.Expression = make([]*Expression, 0)
	ex := &Expression{}
	err := ex.DecodeWat(expr, wf)
	if err != nil {
		return err
	}
	e.Expression = append(e.Expression, ex)

	return nil
}

func (e *CodeEntry) DecodeWat(d string, wf *WasmFile) error {
	e.Locals = make([]ValType, 0)

	s := strings.Trim(d[5:len(d)-1], Whitespace)

	// Optional Identifier
	if s[0] == '$' {
		_, s = ReadToken(s)
	}

	localNames := make(map[int]string)

	for {
		// Skip comments...

		s = strings.Trim(s, Whitespace)
		if strings.HasPrefix(s, ";;") {
			// Skip this line
			line_end := strings.Index(s, "\n")
			if line_end == -1 {
				return nil // All done?
			}
			s = s[line_end:]
		} else if s[0] == '(' {
			var el string
			el, s = ReadElement(s)
			eType, _ := ReadToken(el[1:])
			if eType == "type" {
			} else if eType == "param" {
				// TODO: Use to sanity check
			} else if eType == "result" {
				// TODO: Use to sanity check
			} else if eType == "local" {
				// eg (local $hello i32)
				// eg (local i64 i64)

				types := strings.Trim(el[6:len(el)-1], Whitespace)
				for {
					types = strings.Trim(types, Whitespace)
					if len(types) == 0 {
						break
					}
					var tok string
					tok, types = ReadToken(types)
					if tok[0] == '$' {
						// preRegister a name
						localNames[len(e.Locals)] = tok
					} else {
						l, ok := valTypeToByte[tok]
						if !ok {
							return fmt.Errorf("Invalid local type %s", tok)
						}
						e.Locals = append(e.Locals, l)
					}
				}
			}
		} else {
			break
		}
	}

	// Then just read instructions...
	for {
		s = strings.Trim(s, Whitespace)
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

		ecode = strings.Trim(ecode, Whitespace)

		if len(ecode) > 0 {
			newe := &Expression{}
			err := newe.DecodeWat(ecode, wf)
			if err != nil {
				return err
			}
			e.Expression = append(e.Expression, newe)
		}
	}

	return nil
}

func (e *FunctionEntry) DecodeWat(d string, wf *WasmFile) error {
	s := strings.TrimLeft(d[5:len(d)-1], Whitespace)
	// eg (func $write (type 7) (param i32 i32 i32) (result i32)

	// Optional Identifier
	if s[0] == '$' {
		var fname string
		fname, s = ReadToken(s)
		// Store the name for lookups...
		wf.RegisterNextFunctionName(fname)
	}

	for {
		s = strings.Trim(s, Whitespace)
		if s[0] == '(' {
			var el string
			var err error
			el, s = ReadElement(s)
			eType, _ := ReadToken(el[1:])
			if eType == "type" {
				el = strings.Trim(el[5:len(el)-1], Whitespace)
				e.TypeIndex, err = strconv.Atoi(el)
				return err
			}
		} else {
			return errors.New("Error parsing func. Did not find a type.")
		}
	}
}

func (e *ExportEntry) DecodeWat(d string, wf *WasmFile) error {
	//  (export "memory" (memory 0))
	//  (export "hello" (func $hello))

	s := strings.TrimLeft(d[7:len(d)-1], Whitespace)

	e.Name, s = ReadString(s)
	s = strings.Trim(s, Whitespace)
	el, _ := ReadElement(s)
	etype, erest := ReadToken(el[1:])
	erest = erest[:len(erest)-1]
	if etype == "memory" {
		e.Type = ExportMem
		idx, err := strconv.Atoi(erest)
		if err != nil {
			return err
		}
		e.Index = idx
	} else if etype == "func" {
		e.Type = ExportFunc
		if strings.HasPrefix(erest, "$") {
			fname, _ := ReadToken(erest)
			fid := wf.LookupFunctionID(fname)
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

	s := strings.Trim(d[5:len(d)-1], Whitespace)
	s = SkipComment(s)

	expr, s := ReadElement(s)
	// Read the expression
	expr = expr[1 : len(expr)-1]
	// TODO: Support proper expressions. For now we only support a single instruction
	e.Offset = make([]*Expression, 0)
	ex := &Expression{}
	err := ex.DecodeWat(expr, wf)
	if err != nil {
		return err
	}
	e.Offset = append(e.Offset, ex)

	s = strings.Trim(s, Whitespace)
	var elemType string
	elemType, s = ReadToken(s)
	if elemType == "func" {
		for {
			s = strings.Trim(s, Whitespace)
			if len(s) == 0 {
				break
			}
			var fid string
			var findex int
			fid, s = ReadToken(s)
			if strings.HasPrefix(fid, "$") {
				findex = wf.LookupFunctionID(fid)
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

	s := strings.Trim(d[5:len(d)-1], Whitespace)

	if s[0] == '$' {
		var id string
		id, s = ReadToken(s)
		wf.RegisterNextDataName(id)
	}
	s = strings.Trim(s, Whitespace)
	if s[0] == '(' {
		// Must have a specific Offset already set
		var expr string
		expr, s = ReadElement(s)
		// Read the expression
		expr = expr[1 : len(expr)-1]
		// TODO: Support proper expressions. For now we only support a single instruction
		e.Offset = make([]*Expression, 0)
		ex := &Expression{}
		err := ex.DecodeWat(expr, wf)
		if err != nil {
			return err
		}
		e.Offset = append(e.Offset, ex)
	}

	s = strings.Trim(s, Whitespace)

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
