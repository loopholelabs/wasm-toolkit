package wasmfile

import (
	"debug/dwarf"
	"fmt"
	"io"
)

func (wf *WasmFile) ParseDwarf() error {
	debug_abbrev := wf.GetCustomSectionData(".debug_abbrev")
	debug_aranges := wf.GetCustomSectionData(".debug_aranges")
	debug_info := wf.GetCustomSectionData(".debug_info")
	debug_line := wf.GetCustomSectionData(".debug_line")
	debug_pubnames := wf.GetCustomSectionData(".debug_pubnames")
	debug_ranges := wf.GetCustomSectionData(".debug_ranges")
	debug_str := wf.GetCustomSectionData(".debug_str")

	debug_frame := make([]byte, 0) // call frame info

	dd, err := dwarf.New(debug_abbrev, debug_aranges, debug_frame, debug_info, debug_line, debug_pubnames, debug_ranges, debug_str)
	if err != nil {
		return err
	}

	wf.dwarfData = dd
	return nil
}

func (wf *WasmFile) ParseDwarfLineNumbers() error {
	wf.lineNumbers = make(map[uint64]LineInfo)

	entryReader := wf.dwarfData.Reader()

	for {
		// Read all entries in sequence
		entry, err := entryReader.Next()
		if entry == nil || err == io.EOF {
			// We've reached the end of DWARF entries
			break
		}

		if entry.Tag == dwarf.TagCompileUnit {
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

					wf.lineNumbers[ent.Address] = LineInfo{
						Filename:   ent.File.Name,
						Linenumber: ent.Line,
					}
					fmt.Printf("LINE %d %s:%d\n", ent.Address, ent.File.Name, ent.Line)
				}
			}
		}
	}
	return nil
}

func (wf *WasmFile) ParseDwarfVariables() error {
	entryReader := wf.dwarfData.Reader()

	for {
		// Read all entries in sequence
		entry, err := entryReader.Next()
		if entry == nil || err == io.EOF {
			// We've reached the end of DWARF entries
			break
		}

		fmt.Printf("ENTRY %v\n", entry)

		if entry.Tag == dwarf.TagVariable {

			//	log.Printf("TagVariable %v\n", entry)

		} else if entry.Tag == dwarf.TagFormalParameter {

			/*
				// Show some other dwarf detail...
				vname := "<unknown>"
				vtype := int64(-1)
				vloc := int64(-1)
				for _, field := range entry.Field {
					log.Printf("Field %v\n", field)
					if field.Attr == dwarf.AttrName {
						vname = field.Val.(string)
					} else if field.Attr == dwarf.AttrType {
						switch field.Val.(type) {
						case int64:
							vtype = field.Val.(int64)
						}
					} else if field.Attr == dwarf.AttrLocation {
						switch field.Val.(type) {
						case int64:
							vloc = field.Val.(int64)
						}
					}
				}

				fmt.Printf("FormalParameter %s - %d - %d\n", vname, vtype, vloc)
			*/
		}
	}
	return nil
}
