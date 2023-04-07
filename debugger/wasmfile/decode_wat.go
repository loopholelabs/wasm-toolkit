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

	text := string(data)

	// Read the module
	moduleText, _ := ReadElement(text)

	moduleType, _ := ReadToken(moduleText[1:])

	if moduleType != "module" {
		return errors.New("Invalud module. Expected 'module'")
	}

	// Now read all the individual elements from within the module...

	text = text[len(moduleType)+1:]

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
			err = de.DecodeWat(e)
			wf.Data = append(wf.Data, de)
		} else if eType == "elem" {
			ee := &ElemEntry{}
			err = ee.DecodeWat(e)
			wf.Elem = append(wf.Elem, ee)
		} else if eType == "export" {
			ee := &ExportEntry{}
			err = ee.DecodeWat(e)
			wf.Export = append(wf.Export, ee)
		} else if eType == "func" {
			ee := &FunctionEntry{}
			err = ee.DecodeWat(e)
			wf.Function = append(wf.Function, ee)
			ce := &CodeEntry{}
			err = ce.DecodeWat(e)
			wf.Code = append(wf.Code, ce)
		} else if eType == "global" {
			ge := &GlobalEntry{}
			err = ge.DecodeWat(e)
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
		} else {
			panic(fmt.Sprintf("Unknown element \"%s\"", eType))
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
	fmt.Printf("TODO: Decode Type\n")
	return nil
}

func (e *TableEntry) DecodeWat(d string) error {
	fmt.Printf("TODO: Decode Table\n")
	return nil
}

func (e *MemoryEntry) DecodeWat(d string) error {
	fmt.Printf("TODO: Decode Memory\n")
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

func (e *GlobalEntry) DecodeWat(d string) error {
	fmt.Printf("TODO: Decode Global\n")
	return nil
}

func (e *CodeEntry) DecodeWat(d string) error {
	fmt.Printf("TODO: Decode Code\n")
	return nil
}

func (e *FunctionEntry) DecodeWat(d string) error {
	fmt.Printf("TODO: Decode Function\n")
	return nil
}

func (e *ExportEntry) DecodeWat(d string) error {
	fmt.Printf("TODO: Decode Export\n")
	return nil
}

func (e *ElemEntry) DecodeWat(d string) error {
	fmt.Printf("TODO: Decode Elem\n")
	return nil
}

func (e *DataEntry) DecodeWat(d string) error {
	fmt.Printf("TODO: Decode Data\n")
	return nil
}
