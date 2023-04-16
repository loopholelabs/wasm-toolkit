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
	"math"
)

func NewExpression(data []byte, pc uint64) ([]*Expression, int, error) {

	// This determines when we are finished
	nestCounter := 1

	exps := make([]*Expression, 0)
	ptr := 0

	for {
		if ptr == len(data) {
			break // All done
		}
		opptr := ptr
		opcode := data[ptr]
		ptr++

		// First deal with simple opcodes (No args)
		if Opcode(opcode) == InstrToOpcode["unreachable"] ||
			Opcode(opcode) == InstrToOpcode["nop"] ||
			Opcode(opcode) == InstrToOpcode["return"] ||
			Opcode(opcode) == InstrToOpcode["drop"] ||
			Opcode(opcode) == InstrToOpcode["select"] ||

			Opcode(opcode) == InstrToOpcode["else"] ||

			Opcode(opcode) == InstrToOpcode["i32.eqz"] ||
			Opcode(opcode) == InstrToOpcode["i32.eq"] ||
			Opcode(opcode) == InstrToOpcode["i32.ne"] ||
			Opcode(opcode) == InstrToOpcode["i32.lt_s"] ||
			Opcode(opcode) == InstrToOpcode["i32.lt_u"] ||
			Opcode(opcode) == InstrToOpcode["i32.gt_s"] ||
			Opcode(opcode) == InstrToOpcode["i32.gt_u"] ||
			Opcode(opcode) == InstrToOpcode["i32.le_s"] ||
			Opcode(opcode) == InstrToOpcode["i32.le_u"] ||
			Opcode(opcode) == InstrToOpcode["i32.ge_s"] ||
			Opcode(opcode) == InstrToOpcode["i32.ge_u"] ||
			Opcode(opcode) == InstrToOpcode["i64.eqz"] ||
			Opcode(opcode) == InstrToOpcode["i64.eq"] ||
			Opcode(opcode) == InstrToOpcode["i64.ne"] ||
			Opcode(opcode) == InstrToOpcode["i64.lt_s"] ||
			Opcode(opcode) == InstrToOpcode["i64.lt_u"] ||
			Opcode(opcode) == InstrToOpcode["i64.gt_s"] ||
			Opcode(opcode) == InstrToOpcode["i64.gt_u"] ||
			Opcode(opcode) == InstrToOpcode["i64.le_s"] ||
			Opcode(opcode) == InstrToOpcode["i64.le_u"] ||
			Opcode(opcode) == InstrToOpcode["i64.ge_s"] ||
			Opcode(opcode) == InstrToOpcode["i64.ge_u"] ||
			Opcode(opcode) == InstrToOpcode["f32.eq"] ||
			Opcode(opcode) == InstrToOpcode["f32.ne"] ||
			Opcode(opcode) == InstrToOpcode["f32.lt"] ||
			Opcode(opcode) == InstrToOpcode["f32.gt"] ||
			Opcode(opcode) == InstrToOpcode["f32.le"] ||
			Opcode(opcode) == InstrToOpcode["f32.ge"] ||
			Opcode(opcode) == InstrToOpcode["f64.eq"] ||
			Opcode(opcode) == InstrToOpcode["f64.ne"] ||
			Opcode(opcode) == InstrToOpcode["f64.lt"] ||
			Opcode(opcode) == InstrToOpcode["f64.gt"] ||
			Opcode(opcode) == InstrToOpcode["f64.le"] ||
			Opcode(opcode) == InstrToOpcode["f64.ge"] ||

			Opcode(opcode) == InstrToOpcode["i32.clz"] ||
			Opcode(opcode) == InstrToOpcode["i32.ctz"] ||
			Opcode(opcode) == InstrToOpcode["i32.popcnt"] ||
			Opcode(opcode) == InstrToOpcode["i32.add"] ||
			Opcode(opcode) == InstrToOpcode["i32.sub"] ||
			Opcode(opcode) == InstrToOpcode["i32.mul"] ||
			Opcode(opcode) == InstrToOpcode["i32.div_s"] ||
			Opcode(opcode) == InstrToOpcode["i32.div_u"] ||
			Opcode(opcode) == InstrToOpcode["i32.rem_s"] ||
			Opcode(opcode) == InstrToOpcode["i32.rem_u"] ||
			Opcode(opcode) == InstrToOpcode["i32.and"] ||
			Opcode(opcode) == InstrToOpcode["i32.or"] ||
			Opcode(opcode) == InstrToOpcode["i32.xor"] ||
			Opcode(opcode) == InstrToOpcode["i32.shl"] ||
			Opcode(opcode) == InstrToOpcode["i32.shr_s"] ||
			Opcode(opcode) == InstrToOpcode["i32.shr_u"] ||
			Opcode(opcode) == InstrToOpcode["i32.rotl"] ||
			Opcode(opcode) == InstrToOpcode["i32.rotr"] ||

			Opcode(opcode) == InstrToOpcode["i64.clz"] ||
			Opcode(opcode) == InstrToOpcode["i64.ctz"] ||
			Opcode(opcode) == InstrToOpcode["i64.popcnt"] ||
			Opcode(opcode) == InstrToOpcode["i64.add"] ||
			Opcode(opcode) == InstrToOpcode["i64.sub"] ||
			Opcode(opcode) == InstrToOpcode["i64.mul"] ||
			Opcode(opcode) == InstrToOpcode["i64.div_s"] ||
			Opcode(opcode) == InstrToOpcode["i64.div_u"] ||
			Opcode(opcode) == InstrToOpcode["i64.rem_s"] ||
			Opcode(opcode) == InstrToOpcode["i64.rem_u"] ||
			Opcode(opcode) == InstrToOpcode["i64.and"] ||
			Opcode(opcode) == InstrToOpcode["i64.or"] ||
			Opcode(opcode) == InstrToOpcode["i64.xor"] ||
			Opcode(opcode) == InstrToOpcode["i64.shl"] ||
			Opcode(opcode) == InstrToOpcode["i64.shr_s"] ||
			Opcode(opcode) == InstrToOpcode["i64.shr_u"] ||
			Opcode(opcode) == InstrToOpcode["i64.rotl"] ||
			Opcode(opcode) == InstrToOpcode["i64.rotr"] ||

			Opcode(opcode) == InstrToOpcode["f32.abs"] ||
			Opcode(opcode) == InstrToOpcode["f32.neg"] ||
			Opcode(opcode) == InstrToOpcode["f32.ceil"] ||
			Opcode(opcode) == InstrToOpcode["f32.floor"] ||
			Opcode(opcode) == InstrToOpcode["f32.trunc"] ||
			Opcode(opcode) == InstrToOpcode["f32.nearest"] ||
			Opcode(opcode) == InstrToOpcode["f32.sqrt"] ||
			Opcode(opcode) == InstrToOpcode["f32.add"] ||
			Opcode(opcode) == InstrToOpcode["f32.sub"] ||
			Opcode(opcode) == InstrToOpcode["f32.mul"] ||
			Opcode(opcode) == InstrToOpcode["f32.div"] ||
			Opcode(opcode) == InstrToOpcode["f32.min"] ||
			Opcode(opcode) == InstrToOpcode["f32.max"] ||
			Opcode(opcode) == InstrToOpcode["f32.copysign"] ||

			Opcode(opcode) == InstrToOpcode["f64.abs"] ||
			Opcode(opcode) == InstrToOpcode["f64.neg"] ||
			Opcode(opcode) == InstrToOpcode["f64.ceil"] ||
			Opcode(opcode) == InstrToOpcode["f64.floor"] ||
			Opcode(opcode) == InstrToOpcode["f64.trunc"] ||
			Opcode(opcode) == InstrToOpcode["f64.nearest"] ||
			Opcode(opcode) == InstrToOpcode["f64.sqrt"] ||
			Opcode(opcode) == InstrToOpcode["f64.add"] ||
			Opcode(opcode) == InstrToOpcode["f64.sub"] ||
			Opcode(opcode) == InstrToOpcode["f64.mul"] ||
			Opcode(opcode) == InstrToOpcode["f64.div"] ||
			Opcode(opcode) == InstrToOpcode["f64.min"] ||
			Opcode(opcode) == InstrToOpcode["f64.max"] ||
			Opcode(opcode) == InstrToOpcode["f64.copysign"] ||

			Opcode(opcode) == InstrToOpcode["i32.wrap_i64"] ||
			Opcode(opcode) == InstrToOpcode["i32.trunc_f32_s"] ||
			Opcode(opcode) == InstrToOpcode["i32.trunc_f32_u"] ||
			Opcode(opcode) == InstrToOpcode["i32.trunc_f64_s"] ||
			Opcode(opcode) == InstrToOpcode["i32.trunc_f64_u"] ||
			Opcode(opcode) == InstrToOpcode["i64.extend_i32_s"] ||
			Opcode(opcode) == InstrToOpcode["i64.extend_i32_u"] ||
			Opcode(opcode) == InstrToOpcode["i64.trunc_f32_s"] ||
			Opcode(opcode) == InstrToOpcode["i64.trunc_f32_u"] ||
			Opcode(opcode) == InstrToOpcode["i64.trunc_f64_s"] ||
			Opcode(opcode) == InstrToOpcode["i64.trunc_f64_u"] ||
			Opcode(opcode) == InstrToOpcode["f32.convert_i32_s"] ||
			Opcode(opcode) == InstrToOpcode["f32.convert_i32_u"] ||
			Opcode(opcode) == InstrToOpcode["f32.convert_i64_s"] ||
			Opcode(opcode) == InstrToOpcode["f32.convert_i64_u"] ||
			Opcode(opcode) == InstrToOpcode["f32.demote_f64"] ||
			Opcode(opcode) == InstrToOpcode["f64.convert_i32_s"] ||
			Opcode(opcode) == InstrToOpcode["f64.convert_i32_u"] ||
			Opcode(opcode) == InstrToOpcode["f64.convert_i64_s"] ||
			Opcode(opcode) == InstrToOpcode["f64.convert_i64_u"] ||
			Opcode(opcode) == InstrToOpcode["f64.promote_f32"] ||
			Opcode(opcode) == InstrToOpcode["i32.reinterpret_f32"] ||
			Opcode(opcode) == InstrToOpcode["i64.reinterpret_f64"] ||
			Opcode(opcode) == InstrToOpcode["f32.reinterpret_i32"] ||
			Opcode(opcode) == InstrToOpcode["f64.reinterpret_i64"] ||

			Opcode(opcode) == InstrToOpcode["i32.extend8_s"] ||
			Opcode(opcode) == InstrToOpcode["i32.extend16_s"] ||
			Opcode(opcode) == InstrToOpcode["i64.extend8_s"] ||
			Opcode(opcode) == InstrToOpcode["i64.extend16_s"] ||
			Opcode(opcode) == InstrToOpcode["i64.extend32_s"] {
			exps = append(exps,
				&Expression{
					PC:     pc + uint64(opptr),
					Opcode: Opcode(opcode),
				})

		} else if Opcode(opcode) == InstrToOpcode["br_table"] {
			numLabels, l := binary.Uvarint(data[ptr:])
			ptr += l
			labels := make([]int, 0)
			for ll := 0; ll < int(numLabels); ll++ {
				labelIdx, l := binary.Uvarint(data[ptr:])
				ptr += l
				labels = append(labels, int(labelIdx))
			}
			defaultLabelIdx, l := binary.Uvarint(data[ptr:])
			ptr += l
			exps = append(exps,
				&Expression{
					PC:         pc + uint64(opptr),
					Opcode:     Opcode(opcode),
					Labels:     labels,
					LabelIndex: int(defaultLabelIdx),
				})
		} else if Opcode(opcode) == InstrToOpcode["br"] ||
			Opcode(opcode) == InstrToOpcode["br_if"] {
			val, l := binary.Uvarint(data[ptr:])
			ptr += l
			exps = append(exps,
				&Expression{
					PC:         pc + uint64(opptr),
					Opcode:     Opcode(opcode),
					LabelIndex: int(val),
				})
		} else if Opcode(opcode) == InstrToOpcode["i32.load"] ||
			Opcode(opcode) == InstrToOpcode["i64.load"] ||
			Opcode(opcode) == InstrToOpcode["f32.load"] ||
			Opcode(opcode) == InstrToOpcode["f64.load"] ||
			Opcode(opcode) == InstrToOpcode["i32.load8_s"] ||
			Opcode(opcode) == InstrToOpcode["i32.load8_u"] ||
			Opcode(opcode) == InstrToOpcode["i32.load16_s"] ||
			Opcode(opcode) == InstrToOpcode["i32.load16_u"] ||
			Opcode(opcode) == InstrToOpcode["i64.load8_s"] ||
			Opcode(opcode) == InstrToOpcode["i64.load8_u"] ||
			Opcode(opcode) == InstrToOpcode["i64.load16_s"] ||
			Opcode(opcode) == InstrToOpcode["i64.load16_u"] ||
			Opcode(opcode) == InstrToOpcode["i64.load32_s"] ||
			Opcode(opcode) == InstrToOpcode["i64.load32_u"] ||
			Opcode(opcode) == InstrToOpcode["i32.store"] ||
			Opcode(opcode) == InstrToOpcode["i64.store"] ||
			Opcode(opcode) == InstrToOpcode["f32.store"] ||
			Opcode(opcode) == InstrToOpcode["f64.store"] ||
			Opcode(opcode) == InstrToOpcode["i32.store8"] ||
			Opcode(opcode) == InstrToOpcode["i32.store16"] ||
			Opcode(opcode) == InstrToOpcode["i64.store8"] ||
			Opcode(opcode) == InstrToOpcode["i64.store16"] ||
			Opcode(opcode) == InstrToOpcode["i64.store32"] {
			align, l := binary.Uvarint(data[ptr:])
			ptr += l
			offset, l := binary.Uvarint(data[ptr:])
			ptr += l
			exps = append(exps,
				&Expression{
					PC:        pc + uint64(opptr),
					Opcode:    Opcode(opcode),
					MemAlign:  int(align),
					MemOffset: int(offset),
				})
		} else if Opcode(opcode) == InstrToOpcode["memory.size"] ||
			Opcode(opcode) == InstrToOpcode["memory.grow"] {
			//				memoryIndex := data[ptr]
			ptr++
			exps = append(exps,
				&Expression{
					PC:     pc + uint64(opptr),
					Opcode: Opcode(opcode),
				})
		} else if Opcode(opcode) == InstrToOpcode["block"] ||
			Opcode(opcode) == InstrToOpcode["if"] ||
			Opcode(opcode) == InstrToOpcode["loop"] {
			// Read the blocktype
			valType := data[ptr]
			ptr++

			exps = append(exps,
				&Expression{
					PC:     pc + uint64(opptr),
					Opcode: Opcode(opcode),
					Result: ValType(valType),
				})
			nestCounter++
		} else if Opcode(opcode) == InstrToOpcode["end"] {
			nestCounter--
			if nestCounter == 0 {
				break
			}
			exps = append(exps,
				&Expression{
					PC:     pc + uint64(opptr),
					Opcode: Opcode(opcode),
				})
		} else if Opcode(opcode) == InstrToOpcode["i32.const"] {
			val, l := DecodeSleb128(data[ptr:])
			ptr += int(l)
			exps = append(exps,
				&Expression{
					PC:       pc + uint64(opptr),
					Opcode:   Opcode(opcode),
					I32Value: int32(val),
				})
		} else if Opcode(opcode) == InstrToOpcode["i64.const"] {
			val, l := DecodeSleb128(data[ptr:])
			ptr += int(l)
			exps = append(exps,
				&Expression{
					PC:       pc + uint64(opptr),
					Opcode:   Opcode(opcode),
					I64Value: int64(val),
				})
		} else if Opcode(opcode) == InstrToOpcode["f32.const"] {
			ival := binary.LittleEndian.Uint32(data[ptr : ptr+4])
			val := math.Float32frombits(ival)
			ptr += 4
			exps = append(exps,
				&Expression{
					PC:       pc + uint64(opptr),
					Opcode:   Opcode(opcode),
					F32Value: float32(val),
				})
		} else if Opcode(opcode) == InstrToOpcode["f64.const"] {
			ival := binary.LittleEndian.Uint64(data[ptr : ptr+8])
			val := math.Float64frombits(ival)
			ptr += 8
			exps = append(exps,
				&Expression{
					PC:       pc + uint64(opptr),
					Opcode:   Opcode(opcode),
					F64Value: float64(val),
				})
		} else if Opcode(opcode) == InstrToOpcode["local.get"] ||
			Opcode(opcode) == InstrToOpcode["local.set"] ||
			Opcode(opcode) == InstrToOpcode["local.tee"] {
			val, l := binary.Uvarint(data[ptr:])
			ptr += l
			exps = append(exps,
				&Expression{
					PC:         pc + uint64(opptr),
					Opcode:     Opcode(opcode),
					LocalIndex: int(val),
				})
		} else if Opcode(opcode) == InstrToOpcode["global.get"] ||
			Opcode(opcode) == InstrToOpcode["global.set"] {
			val, l := binary.Uvarint(data[ptr:])
			ptr += l
			exps = append(exps,
				&Expression{
					PC:          pc + uint64(opptr),
					Opcode:      Opcode(opcode),
					GlobalIndex: int(val),
				})
		} else if Opcode(opcode) == InstrToOpcode["call"] {
			val, l := binary.Uvarint(data[ptr:])
			ptr += l
			exps = append(exps,
				&Expression{
					PC:        pc + uint64(opptr),
					Opcode:    Opcode(opcode),
					FuncIndex: int(val),
				})
		} else if Opcode(opcode) == InstrToOpcode["call_indirect"] {
			typeIdx, l := binary.Uvarint(data[ptr:])
			ptr += l
			tableIdx, l := binary.Uvarint(data[ptr:])
			ptr += l
			exps = append(exps,
				&Expression{
					PC:         pc + uint64(opptr),
					Opcode:     Opcode(opcode),
					TypeIndex:  int(typeIdx),
					TableIndex: int(tableIdx),
				})
		} else if Opcode(opcode) == ExtendedOpcodeFC {
			opcode2, l := binary.Uvarint(data[ptr:])
			ptr += l
			// Now deal with opcode2...
			if int(opcode2) == instrToOpcodeFC["memory.copy"] {
				// For now we expect two 0 bytes.
				ptr += 2
				exps = append(exps,
					&Expression{
						PC:        pc + uint64(opptr),
						Opcode:    Opcode(opcode),
						OpcodeExt: int(opcode2),
					})
			} else if int(opcode2) == instrToOpcodeFC["memory.fill"] {
				// For now we expect one 0 byte.
				ptr++
				exps = append(exps,
					&Expression{
						PC:        pc + uint64(opptr),
						Opcode:    Opcode(opcode),
						OpcodeExt: int(opcode2),
					})
			} else if int(opcode2) == instrToOpcodeFC["i32.trunc_sat_f32_s"] ||
				int(opcode2) == instrToOpcodeFC["i32.trunc_sat_f32_u"] ||
				int(opcode2) == instrToOpcodeFC["i32.trunc_sat_f64_s"] ||
				int(opcode2) == instrToOpcodeFC["i32.trunc_sat_f64_u"] ||
				int(opcode2) == instrToOpcodeFC["i64.trunc_sat_f32_s"] ||
				int(opcode2) == instrToOpcodeFC["i64.trunc_sat_f32_u"] ||
				int(opcode2) == instrToOpcodeFC["i64.trunc_sat_f64_s"] ||
				int(opcode2) == instrToOpcodeFC["i64.trunc_sat_f64_u"] {

				exps = append(exps,
					&Expression{
						PC:        pc + uint64(opptr),
						Opcode:    Opcode(opcode),
						OpcodeExt: int(opcode2),
					})

			} else {
				return nil, 0, fmt.Errorf("Unsupported opcode 0xfc %d", opcode2)
			}

		} else {
			ptr--
			return nil, 0, fmt.Errorf("Unsupported opcode %d", data[ptr])
		}
	}
	return exps, ptr, nil
}
