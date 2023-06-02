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

package expression

import (
	"encoding/binary"
	"fmt"
	"math"

	"github.com/loopholelabs/wasm-toolkit/pkg/wasm/encoding"
	"github.com/loopholelabs/wasm-toolkit/pkg/wasm/types"
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

		expr := &Expression{
			PC:     pc + uint64(opptr),
			Opcode: Opcode(opcode),
		}

		// First deal with simple opcodes (No args)
		if expr.HasNoArgs() && expr.Opcode != InstrToOpcode["end"] {
			// Fall through and add it as it is.
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
			expr.Labels = labels
			expr.LabelIndex = int(defaultLabelIdx)
		} else if Opcode(opcode) == InstrToOpcode["br"] ||
			Opcode(opcode) == InstrToOpcode["br_if"] {
			val, l := binary.Uvarint(data[ptr:])
			ptr += l
			expr.LabelIndex = int(val)
		} else if expr.HasMemoryArgs() {
			align, l := binary.Uvarint(data[ptr:])
			ptr += l
			offset, l := binary.Uvarint(data[ptr:])
			ptr += l
			expr.MemAlign = int(align)
			expr.MemOffset = int(offset)
		} else if Opcode(opcode) == InstrToOpcode["memory.size"] ||
			Opcode(opcode) == InstrToOpcode["memory.grow"] {
			//				memoryIndex := data[ptr]
			// TODO: Use this to support multiple memories etc
			ptr++
		} else if Opcode(opcode) == InstrToOpcode["block"] ||
			Opcode(opcode) == InstrToOpcode["if"] ||
			Opcode(opcode) == InstrToOpcode["loop"] {
			// Read the blocktype
			valType := data[ptr]
			ptr++

			expr.Result = types.ValType(valType)
			nestCounter++
		} else if Opcode(opcode) == InstrToOpcode["end"] {
			nestCounter--
			if nestCounter == 0 {
				break
			}
		} else if Opcode(opcode) == InstrToOpcode["i32.const"] {
			val, l := encoding.DecodeSleb128(data[ptr:])
			ptr += int(l)
			expr.I32Value = int32(val)
		} else if Opcode(opcode) == InstrToOpcode["i64.const"] {
			val, l := encoding.DecodeSleb128(data[ptr:])
			ptr += int(l)
			expr.I64Value = int64(val)
		} else if Opcode(opcode) == InstrToOpcode["f32.const"] {
			ival := binary.LittleEndian.Uint32(data[ptr : ptr+4])
			val := math.Float32frombits(ival)
			ptr += 4
			expr.F32Value = float32(val)
		} else if Opcode(opcode) == InstrToOpcode["f64.const"] {
			ival := binary.LittleEndian.Uint64(data[ptr : ptr+8])
			val := math.Float64frombits(ival)
			ptr += 8
			expr.F64Value = float64(val)
		} else if Opcode(opcode) == InstrToOpcode["local.get"] ||
			Opcode(opcode) == InstrToOpcode["local.set"] ||
			Opcode(opcode) == InstrToOpcode["local.tee"] {
			val, l := binary.Uvarint(data[ptr:])
			ptr += l
			expr.LocalIndex = int(val)
		} else if Opcode(opcode) == InstrToOpcode["global.get"] ||
			Opcode(opcode) == InstrToOpcode["global.set"] {
			val, l := binary.Uvarint(data[ptr:])
			ptr += l
			expr.GlobalIndex = int(val)
		} else if Opcode(opcode) == InstrToOpcode["call"] {
			val, l := binary.Uvarint(data[ptr:])
			ptr += l
			expr.FuncIndex = int(val)
		} else if Opcode(opcode) == InstrToOpcode["call_indirect"] {
			typeIdx, l := binary.Uvarint(data[ptr:])
			ptr += l
			tableIdx, l := binary.Uvarint(data[ptr:])
			ptr += l
			expr.TypeIndex = int(typeIdx)
			expr.TableIndex = int(tableIdx)
		} else if Opcode(opcode) == ExtendedOpcodeFC {
			opcode2, l := binary.Uvarint(data[ptr:])
			ptr += l
			expr.OpcodeExt = int(opcode2)
			// Now deal with opcode2...
			if int(opcode2) == instrToOpcodeFC["memory.copy"] {
				// For now we expect two 0 bytes.
				ptr += 2
			} else if int(opcode2) == instrToOpcodeFC["memory.fill"] {
				// For now we expect one 0 byte.
				ptr++
			} else if int(opcode2) == instrToOpcodeFC["i32.trunc_sat_f32_s"] ||
				int(opcode2) == instrToOpcodeFC["i32.trunc_sat_f32_u"] ||
				int(opcode2) == instrToOpcodeFC["i32.trunc_sat_f64_s"] ||
				int(opcode2) == instrToOpcodeFC["i32.trunc_sat_f64_u"] ||
				int(opcode2) == instrToOpcodeFC["i64.trunc_sat_f32_s"] ||
				int(opcode2) == instrToOpcodeFC["i64.trunc_sat_f32_u"] ||
				int(opcode2) == instrToOpcodeFC["i64.trunc_sat_f64_s"] ||
				int(opcode2) == instrToOpcodeFC["i64.trunc_sat_f64_u"] {

			} else {
				return nil, 0, fmt.Errorf("Unsupported opcode 0xfc %d", opcode2)
			}

		} else {
			ptr--
			return nil, 0, fmt.Errorf("Unsupported opcode %d", data[ptr])
		}

		// Add it on to the list...
		exps = append(exps, expr)
	}
	return exps, ptr, nil
}
