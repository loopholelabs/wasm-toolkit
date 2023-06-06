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

const DW_OP_addr = 0x03
const DW_OP_deref = 0x06
const DW_OP_const1u = 0x08
const DW_OP_const1s = 0x09
const DW_OP_const2u = 0x0a
const DW_OP_const2s = 0x0b
const DW_OP_const4u = 0x0c
const DW_OP_const4s = 0x0d
const DW_OP_const8u = 0x0e
const DW_OP_const8s = 0x0f

const DW_OP_constu = 0x10
const DW_OP_consts = 0x11
const DW_OP_dup = 0x12
const DW_OP_drop = 0x13
const DW_OP_over = 0x14
const DW_OP_pick = 0x15
const DW_OP_swap = 0x16
const DW_OP_rot = 0x17
const DW_OP_xderef = 0x18
const DW_OP_abs = 0x19
const DW_OP_and = 0x1a
const DW_OP_div = 0x1b
const DW_OP_minus = 0x1c
const DW_OP_mod = 0x1d
const DW_OP_mul = 0x1e
const DW_OP_neg = 0x1f

const DW_OP_not = 0x20
const DW_OP_or = 0x21
const DW_OP_plus = 0x22
const DW_OP_plus_uconst = 0x23
const DW_OP_shl = 0x24
const DW_OP_shr = 0x25
const DW_OP_shra = 0x26
const DW_OP_xor = 0x27
const DW_OP_bra = 0x28
const DW_OP_eq = 0x29
const DW_OP_ge = 0x2a
const DW_OP_gt = 0x2b
const DW_OP_le = 0x2c
const DW_OP_lt = 0x2d
const DW_OP_ne = 0x2e
const DW_OP_skip = 0x2f

const DW_OP_regx = 0x90
const DW_OP_fbreg = 0x91
const DW_OP_bregx = 0x92
const DW_OP_piece = 0x93
const DW_OP_deref_size = 0x94
const DW_OP_xderef_size = 0x95
const DW_OP_nop = 0x96
const DW_OP_push_object_address = 0x97
const DW_OP_call2 = 0x98
const DW_OP_call4 = 0x99
const DW_OP_call_ref = 0x9a
const DW_OP_form_tls_address = 0x9b
const DW_OP_call_frame_cfa = 0x9c
const DW_OP_bit_piece = 0x9d
const DW_OP_implicit_value = 0x9e
const DW_OP_stack_value = 0x9f

const DW_OP_implicit_pointer = 0xa0
const DW_OP_addrx = 0xa1
const DW_OP_constx = 0xa2
const DW_OP_entry_value = 0xa3
const DW_OP_const_type = 0xa4
const DW_OP_regval_type = 0xa5
const DW_OP_deref_type = 0xa6
const DW_OP_xderef_type = 0xa7
const DW_OP_convert = 0xa8
const DW_OP_reinterpret = 0xa9

const DW_OP_lo_user = 0xe0
const DW_OP_GNU_push_tls_address = 0xe0

const DW_OP_WASM_location = 0xed

const DW_OP_GNU_implicit_pointer = 0xf2
const DW_OP_GNU_entry_value = 0xf3
const DW_OP_GNU_const_type = 0xf4
const DW_OP_GNU_regval_type = 0xf5
const DW_OP_GNU_deref_type = 0xf6
const DW_OP_GNU_convert = 0xf7
const DW_OP_GNU_parameter_ref = 0xfa
const DW_OP_hi_user = 0xff

// DW_OP_WASM_location types
const DW_Location_Local = 0
const DW_Location_Global = 1
const DW_Location_Stack = 2 // 0 = bottom of the stack
const DW_Location_Global_i32 = 3

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
