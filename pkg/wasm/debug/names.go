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
	"encoding/binary"
	"fmt"
	"strings"
)

const subsectionModuleNames = 0
const subsectionFunctionNames = 1
const subsectionLocalNames = 2
const subsectionLabelNames = 3
const subsectionTypeNames = 4
const subsectionTableNames = 5
const subsectionMemoryNames = 6
const subsectionGlobalNames = 7
const subsectionDataNames = 9

/**
 * Parse the custom name section from a wasm file
 *
 */
func (wd *WasmDebug) ParseNameSectionData(nameData []byte) {
	wd.FunctionNames = make(map[int]string)
	wd.GlobalNames = make(map[int]string)
	wd.DataNames = make(map[int]string)

	ptr := 0

	for {
		if ptr == len(nameData) {
			break
		}
		// Read the subsection data...
		subsectionID := nameData[ptr]
		ptr++
		subsectionLength, l := binary.Uvarint(nameData[ptr:])
		ptr += l
		data := nameData[ptr : ptr+int(subsectionLength)]
		ptr += int(subsectionLength)

		if subsectionID == subsectionFunctionNames {
			// Now read all the function names...
			nameVecLength, l := binary.Uvarint(data)
			data = data[l:]

			for i := 0; i < int(nameVecLength); i++ {
				idx, l := binary.Uvarint(data)
				data = data[l:]
				nameLength, l := binary.Uvarint(data)
				data = data[l:]
				nameValue := data[:nameLength]
				data = data[nameLength:]

				// Make sure it's unique and doesn't already exist...
				dupidx := 1
				for {
					fname := fmt.Sprintf("$%s", string(nameValue))
					if dupidx > 1 {
						fname = fmt.Sprintf("%s_%d", fname, dupidx)
					}
					exists := false
					for _, n := range wd.FunctionNames {
						if n == fname {
							exists = true
							break
						}
					}

					if !exists {
						wd.FunctionNames[int(idx)] = fname
						break
					}
					dupidx++
				}
			}

		} else if subsectionID == subsectionGlobalNames {
			// Now read all the global names...
			nameVecLength, l := binary.Uvarint(data)
			data = data[l:]

			for i := 0; i < int(nameVecLength); i++ {
				idx, l := binary.Uvarint(data)
				data = data[l:]
				nameLength, l := binary.Uvarint(data)
				data = data[l:]
				nameValue := data[:nameLength]
				data = data[nameLength:]

				wd.GlobalNames[int(idx)] = fmt.Sprintf("$%s", string(nameValue))
			}
		} else if subsectionID == subsectionDataNames {
			nameVecLength, l := binary.Uvarint(data)
			data = data[l:]

			for i := 0; i < int(nameVecLength); i++ {
				idx, l := binary.Uvarint(data)
				data = data[l:]
				nameLength, l := binary.Uvarint(data)
				data = data[l:]
				nameValue := data[:nameLength]
				data = data[nameLength:]

				wd.DataNames[int(idx)] = fmt.Sprintf("$%s", string(nameValue))
			}
		} else {
			//fmt.Printf("TODO: Name %d - %d\n", subsectionID, subsectionLength)
		}
	}

}

func (wd *WasmDebug) GetFunctionIdentifier(fid int, defaultEmpty bool) string {
	f, ok := wd.FunctionNames[fid]
	if ok {
		f = strings.ReplaceAll(f, "(", "_")
		f = strings.ReplaceAll(f, ")", "_")
		f = strings.ReplaceAll(f, "{", "_")
		f = strings.ReplaceAll(f, "}", "_")
		f = strings.ReplaceAll(f, "[", "_")
		f = strings.ReplaceAll(f, "]", "_")
		f = strings.ReplaceAll(f, ",", "_")
		return f
	}
	if defaultEmpty {
		return ""
	}
	return fmt.Sprintf("%d", fid)
}

func (wd *WasmDebug) GetGlobalIdentifier(gid int, defaultEmpty bool) string {
	f, ok := wd.GlobalNames[gid]
	if ok {
		f = strings.ReplaceAll(f, "(", "_")
		f = strings.ReplaceAll(f, ")", "_")
		return f
	}
	if defaultEmpty {
		return ""
	}
	return fmt.Sprintf("%d", gid)
}

func (wd *WasmDebug) GetDataIdentifier(did int) string {
	f, ok := wd.DataNames[did]
	if ok {
		f = strings.ReplaceAll(f, "(", "_")
		f = strings.ReplaceAll(f, ")", "_")
		return f
	}
	return ""
}

func (wd *WasmDebug) LookupDataId(n string) int {
	for idx, name := range wd.DataNames {
		if n == name {
			return idx
		}
	}
	return -1
}

func (wd *WasmDebug) LookupGlobalID(n string) int {
	for idx, name := range wd.GlobalNames {
		if n == name {
			return idx
		}
	}
	return -1
}
