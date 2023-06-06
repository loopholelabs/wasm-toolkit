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

import "encoding/binary"

type DwarfLocations struct {
	data []byte
}

type LocationData struct {
	StartAddress uint32
	EndAddress   uint32
	Expression   []byte
}

func NewDwarfLocations(d []byte) *DwarfLocations {
	return &DwarfLocations{
		data: d,
	}
}

func (dl *DwarfLocations) ReadLocation(p uint64) []*LocationData {
	baseAddress := uint32(0)
	ld := make([]*LocationData, 0)

	ptr := p
	for {
		low := binary.LittleEndian.Uint32(dl.data[ptr:])
		ptr += 4
		high := binary.LittleEndian.Uint32(dl.data[ptr:])
		ptr += 4
		if low == 0 && high == 0 {
			break
		}
		if low == 0xffffffff {
			baseAddress = high
		} else {
			// Read expr len
			explen := binary.LittleEndian.Uint16(dl.data[ptr:])
			ptr += 2
			expr := dl.data[ptr : ptr+uint64(explen)]
			ptr += uint64(explen)
			ld = append(ld, &LocationData{
				StartAddress: baseAddress + low,
				EndAddress:   baseAddress + high,
				Expression:   expr,
			})
		}
	}
	return ld
}

const DW_OP_WASM_location = 0xed
const DW_Location_Local = 0
const DW_Location_Global = 1
const DW_Location_Stack = 2 // 0 = bottom of the stack
const DW_Location_Global_i32 = 3

const DW_OP_addr = 0x03

const DW_OP_stack_value = 0x9f
const DW_OP_piece = 0x93

type WasmLocation struct {
	IsLocal  bool
	IsGlobal bool
	IsStack  bool
	Index    uint64
}

func (ld *LocationData) ExtractWasmLocations() []*WasmLocation {
	locs := make([]*WasmLocation, 0)
	data := ld.Expression
	for {
		if len(data) == 0 {
			break
		}
		opcode := data[0]
		data = data[1:]
		if opcode == DW_OP_stack_value {
			// Fine...
		} else if opcode == DW_OP_piece {
			_, l := binary.Uvarint(data)
			data = data[l:]
		} else if opcode == DW_OP_WASM_location {
			t := data[0]
			data = data[1:]
			var index uint64
			if t == 3 {
				index = uint64(binary.LittleEndian.Uint32(data))
				data = data[4:]
			} else {
				var l int
				index, l = binary.Uvarint(data)
				data = data[l:]
			}
			locs = append(locs, &WasmLocation{
				IsLocal:  t == DW_Location_Local,
				IsGlobal: t == DW_Location_Global || t == DW_Location_Global_i32,
				IsStack:  t == DW_Location_Stack,
				Index:    index,
			})

		} else {
			// FIXME: Deal with other dwarf opcodes
			//			fmt.Printf("WARN: Unknown dwarf expression opcode %d %x\n", opcode, orgdata)
			return locs
		}
	}
	return locs
}
