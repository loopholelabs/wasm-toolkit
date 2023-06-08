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

package debug

import (
	"debug/dwarf"
	"fmt"
	"io"
)

type LineInfo struct {
	Filename   string
	Linenumber int
	Column     int
}

func (wd *WasmDebug) ParseDwarfLineNumbers() error {
	wd.LineNumbers = make(map[uint64]LineInfo)

	if wd.DwarfData == nil {
		return nil
	}
	entryReader := wd.DwarfData.Reader()

	for {
		// Read all entries in sequence
		entry, err := entryReader.Next()
		if entry == nil || err == io.EOF {
			// We've reached the end of DWARF entries
			break
		}

		if entry.Tag == dwarf.TagCompileUnit {
			liner, err := wd.DwarfData.LineReader(entry)

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

					wd.LineNumbers[ent.Address] = LineInfo{
						Filename:   ent.File.Name,
						Linenumber: ent.Line,
						Column:     ent.Column,
					}
				}
			}
		}
	}

	return nil
}

func (wd *WasmDebug) GetLineNumberInfo(pc uint64) string {
	// See if we have any line info...
	lineInfo := ""
	li, ok := wd.LineNumbers[pc]
	if ok {
		lineInfo = fmt.Sprintf("%s:%d.%d", li.Filename, li.Linenumber, li.Column)
	}
	return lineInfo
}

func (wd *WasmDebug) GetLineNumberBefore(start uint64, codePC uint64) string {
	for pc := codePC; pc >= start; pc-- {
		l := wd.GetLineNumberInfo(pc)
		if l != "" {
			return l
		}
	}
	return ""
}

func (wd *WasmDebug) GetLineNumberRange(start uint64, end uint64) string {
	// Collect all the ranges together...
	ranges := make(map[string][]int)

	for pc := start; pc <= end; pc++ {
		// Look it up...
		li, ok := wd.LineNumbers[pc]
		if ok {
			m, ok2 := ranges[li.Filename]
			if ok2 {
				// Add it on...
				ranges[li.Filename] = append(m, li.Linenumber)
			} else {
				ranges[li.Filename] = []int{li.Linenumber}
			}
		}
	}

	// Now lets bring things together...
	info := ""

	for filename, rg := range ranges {
		min := -1
		max := -1
		for _, v := range rg {
			if (min == -1) || (v < min) {
				min = v
			}
			if (max == -1) || (v > max) {
				max = v
			}
		}
		if info != "" {
			info = fmt.Sprintf("%s,", info)
		}
		info = fmt.Sprintf("%s%s(%d-%d)", info, filename, min, max)
	}

	return info
}
