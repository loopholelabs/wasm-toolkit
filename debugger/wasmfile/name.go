package wasmfile

import (
	"encoding/binary"
	"fmt"
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

				wf.functionNames[int(idx)] = string(nameValue)
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

				wf.globalNames[int(idx)] = string(nameValue)
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

				wf.dataNames[int(idx)] = string(nameValue)
			}

		} else {
			fmt.Printf("TODO: Name %d - %d\n", subsectionID, subsectionLength)
		}
	}
	return nil
}

func (wf *WasmFile) GetFunctionIdentifier(fid int) string {
	f, ok := wf.functionNames[fid]
	if ok {
		return fmt.Sprintf("$%s", f)
	}
	return fmt.Sprintf("%d", fid)
}

func (wf *WasmFile) GetGlobalIdentifier(gid int) string {
	f, ok := wf.globalNames[gid]
	if ok {
		return fmt.Sprintf("$%s", f)
	}
	return fmt.Sprintf("%d", gid)
}

func (wf *WasmFile) GetDataIdentifier(did int) string {
	f, ok := wf.dataNames[did]
	if ok {
		return fmt.Sprintf("$%s", f)
	}
	return fmt.Sprintf("%d", did)
}
