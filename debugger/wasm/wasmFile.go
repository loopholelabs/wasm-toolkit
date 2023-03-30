package wasm

import (
	"bytes"
	"debug/dwarf"
	"encoding/binary"
	"io"
	"io/ioutil"
)

type WasmFileInfo struct {
	filename          string
	dwarfData         *dwarf.Data
	FunctionLocations map[int]*FunctionInfo
}

type FunctionInfo struct {
	codeStart uint64
	codeEnd   uint64
	LineMin   int
	LineMax   int
	LineFile  string
}

func (fi *FunctionInfo) addLineNumberData(f string, l int) {
	fi.LineFile = f
	// TODO: More than one file?
	if fi.LineMax == -1 || l > fi.LineMax {
		fi.LineMax = l
	}
	if fi.LineMin == -1 || l < fi.LineMin {
		fi.LineMin = l
	}
}

// Read a wasmfile into some usable stuff
func NewWasmFile(f string) *WasmFileInfo {
	return &WasmFileInfo{
		filename:          f,
		FunctionLocations: make(map[int]*FunctionInfo, 0),
	}
}

func (wf *WasmFileInfo) FindFunctionByPtr(ptr uint64) int {
	for i, fi := range wf.FunctionLocations {
		if ptr >= fi.codeStart && ptr <= fi.codeEnd {
			return i
		}
	}
	return -1
}

func (wf *WasmFileInfo) ReadDwarf() error {
	data, err := ioutil.ReadFile(wf.filename)
	if err != nil {
		return err
	}

	// TODO: Check the wasm header/version
	data = data[8:]

	rr := bytes.NewReader(data)

	debug_info := make([]byte, 0)
	//debug_pubtypes := make([]byte, 0)
	//debug_loc := make([]byte, 0)
	debug_ranges := make([]byte, 0)
	debug_aranges := make([]byte, 0)
	debug_abbrev := make([]byte, 0)
	debug_line := make([]byte, 0)
	debug_str := make([]byte, 0)
	debug_pubnames := make([]byte, 0)

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

		if sectionType == 0 {
			srr := bytes.NewReader(sectionData)
			nameLength, err := binary.ReadUvarint(srr)
			if err != nil {
				return err
			}
			nameData := make([]byte, nameLength)
			srr.Read(nameData)

			// Read the rest
			customData, err := ioutil.ReadAll(srr)
			if err != nil {
				return err
			}

			name := string(nameData)

			if name == ".debug_info" {
				debug_info = customData
			} else if name == ".debug_pubtypes" {
				//debug_pubtypes = customData
			} else if name == ".debug_loc" {
				//debug_loc = customData
			} else if name == ".debug_ranges" {
				debug_ranges = customData
			} else if name == ".debug_aranges" {
				debug_aranges = customData
			} else if name == ".debug_abbrev" {
				debug_abbrev = customData
			} else if name == ".debug_line" {
				debug_line = customData
			} else if name == ".debug_str" {
				debug_str = customData
			} else if name == ".debug_pubnames" {
				debug_pubnames = customData
			}
		}
	}

	/*
		fmt.Printf(".debug_info %d\n", len(debug_info))         // core dwarf info
		fmt.Printf(".debug_pubtypes %d\n", len(debug_pubtypes)) // Lookup for global types
		fmt.Printf(".debug_loc %d\n", len(debug_loc))           // location lists DW_AT
		fmt.Printf(".debug_ranges %d\n", len(debug_ranges))     // address ranges DW_AT
		fmt.Printf(".debug_aranges %d\n", len(debug_aranges))   // map addresses -> compilation units
		fmt.Printf(".debug_abbrev %d\n", len(debug_abbrev))     // Abbrev for all CUs
		fmt.Printf(".debug_line %d\n", len(debug_line))         // line numbers
		fmt.Printf(".debug_str %d\n", len(debug_str))           // str table used in .debug_info
		fmt.Printf(".debug_pubnames %d\n", len(debug_pubnames)) // lookup for global obj/fun
	*/

	debug_frame := make([]byte, 0) // call frame info

	dd, err := dwarf.New(debug_abbrev, debug_aranges, debug_frame, debug_info, debug_line, debug_pubnames, debug_ranges, debug_str)
	if err != nil {
		return err
	}

	wf.dwarfData = dd
	return nil
}

func (wf *WasmFileInfo) ReadFunctionInfo() error {
	data, err := ioutil.ReadFile(wf.filename)
	if err != nil {
		return err
	}

	// TODO: Check the wasm header/version
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

		if sectionType == 10 {
			ptr := 0
			codeLength, l := binary.Uvarint(sectionData)
			if err != nil {
				return err
			}
			ptr += l

			for i := 0; i < int(codeLength); i++ {
				codeptr := uint64(ptr)
				clen, l := binary.Uvarint(sectionData[ptr:])
				if err != nil {
					return err
				}
				ptr += l

				// Save the function data
				wf.FunctionLocations[i] = &FunctionInfo{
					codeStart: codeptr,
					codeEnd:   codeptr + clen,
					LineMin:   -1,
					LineMax:   -1,
					LineFile:  "<unknown>",
				}

				ptr += int(clen)
			}
		}
	}
	return nil
}

func (wf *WasmFileInfo) ParseDwarfLineNumbers() error {
	entryReader := wf.dwarfData.Reader()

	for {
		// Read all entries in sequence
		entry, err := entryReader.Next()
		if entry == nil || err == io.EOF {
			// We've reached the end of DWARF entries
			break
		}

		if entry.Tag == dwarf.TagCompileUnit {

			// Go through fields
			/*
				for _, field := range entry.Field {
					if field.Attr == dwarf.AttrName {
						fmt.Println(field.Val.(string))
					}
				}
			*/

			liner, err := wf.dwarfData.LineReader(entry)

			if err != nil {
				return err
			}
			if liner != nil {
				ent := dwarf.LineEntry{}
				for {
					err = liner.Next(&ent)
					if err == io.EOF {
						break
					}
					fnid := wf.FindFunctionByPtr(ent.Address)
					if fnid != -1 {
						// Add the line number info to the function
						wf.FunctionLocations[fnid].addLineNumberData(ent.File.Name, ent.Line)
						//fmt.Printf("Line %d: %s %d -> fnid %d\n", ent.Address, ent.File.Name, ent.Line, fnid)
					} else {
						// Unknown function...
					}
				}
			}
		}
	}
	return nil
}
