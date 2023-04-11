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

func (wf *WasmFile) ParseName() error {
	wf.functionNames = make(map[int]string)
	wf.globalNames = make(map[int]string)
	wf.dataNames = make(map[int]string)

	nameData := wf.GetCustomSectionData("name")
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

				wf.functionNames[int(idx)] = fmt.Sprintf("$%s", string(nameValue))
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

				wf.globalNames[int(idx)] = fmt.Sprintf("$%s", string(nameValue))
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

				wf.dataNames[int(idx)] = fmt.Sprintf("$%s", string(nameValue))
			}

		} else {
			fmt.Printf("TODO: Name %d - %d\n", subsectionID, subsectionLength)
		}
	}
	return nil
}

func (wf *WasmFile) GetFunctionIdentifier(fid int, defaultEmpty bool) string {
	f, ok := wf.functionNames[fid]
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

func (wf *WasmFile) GetGlobalIdentifier(gid int) string {
	f, ok := wf.globalNames[gid]
	if ok {
		f = strings.ReplaceAll(f, "(", "_")
		f = strings.ReplaceAll(f, ")", "_")
		return f
	}
	return fmt.Sprintf("%d", gid)
}

func (wf *WasmFile) GetDataIdentifier(did int) string {
	f, ok := wf.dataNames[did]
	if ok {
		f = strings.ReplaceAll(f, "(", "_")
		f = strings.ReplaceAll(f, ")", "_")
		return f
	}
	return ""
}

func (wf *WasmFile) LookupDataId(n string) int {
	for idx, name := range wf.dataNames {
		if n == name {
			return idx
		}
	}
	return -1
}

func (wf *WasmFile) LookupGlobalID(n string) int {
	for idx, name := range wf.globalNames {
		if n == name {
			return idx
		}
	}
	return -1
}
