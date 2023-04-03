package wasmfile

import (
	"encoding/binary"
	"fmt"
	"math"
)

func NewExpression(data []byte, pc uint64) ([]*Expression, int) {
	exps := make([]*Expression, 0)
	ptr := 0

	for {
		opptr := ptr
		opcode := data[ptr]
		ptr++

		// First deal with simple opcodes (No args)
		if Opcode(opcode) == instrToOpcode["unreachable"] ||
			Opcode(opcode) == instrToOpcode["nop"] ||
			Opcode(opcode) == instrToOpcode["return"] ||
			Opcode(opcode) == instrToOpcode["drop"] ||
			Opcode(opcode) == instrToOpcode["select"] ||
			Opcode(opcode) == instrToOpcode["i32.eqz"] ||
			Opcode(opcode) == instrToOpcode["i32.eq"] ||
			Opcode(opcode) == instrToOpcode["i32.ne"] ||
			Opcode(opcode) == instrToOpcode["i32.lt_s"] ||
			Opcode(opcode) == instrToOpcode["i32.lt_u"] ||
			Opcode(opcode) == instrToOpcode["i32.gt_s"] ||
			Opcode(opcode) == instrToOpcode["i32.gt_u"] ||
			Opcode(opcode) == instrToOpcode["i32.le_s"] ||
			Opcode(opcode) == instrToOpcode["i32.le_u"] ||
			Opcode(opcode) == instrToOpcode["i32.ge_s"] ||
			Opcode(opcode) == instrToOpcode["i32.ge_u"] ||
			Opcode(opcode) == instrToOpcode["i64.eqz"] ||
			Opcode(opcode) == instrToOpcode["i64.eq"] ||
			Opcode(opcode) == instrToOpcode["i64.ne"] ||
			Opcode(opcode) == instrToOpcode["i64.lt_s"] ||
			Opcode(opcode) == instrToOpcode["i64.lt_u"] ||
			Opcode(opcode) == instrToOpcode["i64.gt_s"] ||
			Opcode(opcode) == instrToOpcode["i64.gt_u"] ||
			Opcode(opcode) == instrToOpcode["i64.le_s"] ||
			Opcode(opcode) == instrToOpcode["i64.le_u"] ||
			Opcode(opcode) == instrToOpcode["i64.ge_s"] ||
			Opcode(opcode) == instrToOpcode["i64.ge_u"] ||
			Opcode(opcode) == instrToOpcode["f32.eq"] ||
			Opcode(opcode) == instrToOpcode["f32.ne"] ||
			Opcode(opcode) == instrToOpcode["f32.lt"] ||
			Opcode(opcode) == instrToOpcode["f32.gt"] ||
			Opcode(opcode) == instrToOpcode["f32.le"] ||
			Opcode(opcode) == instrToOpcode["f32.ge"] ||
			Opcode(opcode) == instrToOpcode["f64.eq"] ||
			Opcode(opcode) == instrToOpcode["f64.ne"] ||
			Opcode(opcode) == instrToOpcode["f64.lt"] ||
			Opcode(opcode) == instrToOpcode["f64.gt"] ||
			Opcode(opcode) == instrToOpcode["f64.le"] ||
			Opcode(opcode) == instrToOpcode["f64.ge"] ||

			Opcode(opcode) == instrToOpcode["i32.clz"] ||
			Opcode(opcode) == instrToOpcode["i32.ctz"] ||
			Opcode(opcode) == instrToOpcode["i32.popcnt"] ||
			Opcode(opcode) == instrToOpcode["i32.add"] ||
			Opcode(opcode) == instrToOpcode["i32.sub"] ||
			Opcode(opcode) == instrToOpcode["i32.mul"] ||
			Opcode(opcode) == instrToOpcode["i32.div_s"] ||
			Opcode(opcode) == instrToOpcode["i32.div_u"] ||
			Opcode(opcode) == instrToOpcode["i32.rem_s"] ||
			Opcode(opcode) == instrToOpcode["i32.rem_u"] ||
			Opcode(opcode) == instrToOpcode["i32.and"] ||
			Opcode(opcode) == instrToOpcode["i32.or"] ||
			Opcode(opcode) == instrToOpcode["i32.xor"] ||
			Opcode(opcode) == instrToOpcode["i32.shl"] ||
			Opcode(opcode) == instrToOpcode["i32.shr_s"] ||
			Opcode(opcode) == instrToOpcode["i32.shr_u"] ||
			Opcode(opcode) == instrToOpcode["i32.rotl_s"] ||
			Opcode(opcode) == instrToOpcode["i32.rotr_u"] ||

			Opcode(opcode) == instrToOpcode["i64.clz"] ||
			Opcode(opcode) == instrToOpcode["i64.ctz"] ||
			Opcode(opcode) == instrToOpcode["i64.popcnt"] ||
			Opcode(opcode) == instrToOpcode["i64.add"] ||
			Opcode(opcode) == instrToOpcode["i64.sub"] ||
			Opcode(opcode) == instrToOpcode["i64.mul"] ||
			Opcode(opcode) == instrToOpcode["i64.div_s"] ||
			Opcode(opcode) == instrToOpcode["i64.div_u"] ||
			Opcode(opcode) == instrToOpcode["i64.rem_s"] ||
			Opcode(opcode) == instrToOpcode["i64.rem_u"] ||
			Opcode(opcode) == instrToOpcode["i64.and"] ||
			Opcode(opcode) == instrToOpcode["i64.or"] ||
			Opcode(opcode) == instrToOpcode["i64.xor"] ||
			Opcode(opcode) == instrToOpcode["i64.shl"] ||
			Opcode(opcode) == instrToOpcode["i64.shr_s"] ||
			Opcode(opcode) == instrToOpcode["i64.shr_u"] ||
			Opcode(opcode) == instrToOpcode["i64.rotl_s"] ||
			Opcode(opcode) == instrToOpcode["i64.rotr_u"] ||

			Opcode(opcode) == instrToOpcode["f32.abs"] ||
			Opcode(opcode) == instrToOpcode["f32.neg"] ||
			Opcode(opcode) == instrToOpcode["f32.ceil"] ||
			Opcode(opcode) == instrToOpcode["f32.floor"] ||
			Opcode(opcode) == instrToOpcode["f32.trunc"] ||
			Opcode(opcode) == instrToOpcode["f32.nearest"] ||
			Opcode(opcode) == instrToOpcode["f32.sqrt"] ||
			Opcode(opcode) == instrToOpcode["f32.add"] ||
			Opcode(opcode) == instrToOpcode["f32.sub"] ||
			Opcode(opcode) == instrToOpcode["f32.mul"] ||
			Opcode(opcode) == instrToOpcode["f32.div"] ||
			Opcode(opcode) == instrToOpcode["f32.min"] ||
			Opcode(opcode) == instrToOpcode["f32.max"] ||
			Opcode(opcode) == instrToOpcode["f32.copysign"] ||

			Opcode(opcode) == instrToOpcode["f64.abs"] ||
			Opcode(opcode) == instrToOpcode["f64.neg"] ||
			Opcode(opcode) == instrToOpcode["f64.ceil"] ||
			Opcode(opcode) == instrToOpcode["f64.floor"] ||
			Opcode(opcode) == instrToOpcode["f64.trunc"] ||
			Opcode(opcode) == instrToOpcode["f64.nearest"] ||
			Opcode(opcode) == instrToOpcode["f64.sqrt"] ||
			Opcode(opcode) == instrToOpcode["f64.add"] ||
			Opcode(opcode) == instrToOpcode["f64.sub"] ||
			Opcode(opcode) == instrToOpcode["f64.mul"] ||
			Opcode(opcode) == instrToOpcode["f64.div"] ||
			Opcode(opcode) == instrToOpcode["f64.min"] ||
			Opcode(opcode) == instrToOpcode["f64.max"] ||
			Opcode(opcode) == instrToOpcode["f64.copysign"] ||

			Opcode(opcode) == instrToOpcode["i32.wrap_i64"] ||
			Opcode(opcode) == instrToOpcode["i32.trunc_f32_s"] ||
			Opcode(opcode) == instrToOpcode["i32.trunc_f32_u"] ||
			Opcode(opcode) == instrToOpcode["i32.trunc_f64_s"] ||
			Opcode(opcode) == instrToOpcode["i32.trunc_f64_u"] ||
			Opcode(opcode) == instrToOpcode["i64.extend_i32_s"] ||
			Opcode(opcode) == instrToOpcode["i64.extend_i32_u"] ||
			Opcode(opcode) == instrToOpcode["i64.trunc_f32_s"] ||
			Opcode(opcode) == instrToOpcode["i64.trunc_f32_u"] ||
			Opcode(opcode) == instrToOpcode["i64.trunc_f64_s"] ||
			Opcode(opcode) == instrToOpcode["i64.trunc_f64_u"] ||
			Opcode(opcode) == instrToOpcode["f32.convert_i32_s"] ||
			Opcode(opcode) == instrToOpcode["f32.convert_i32_u"] ||
			Opcode(opcode) == instrToOpcode["f32.convert_i64_s"] ||
			Opcode(opcode) == instrToOpcode["f32.convert_i64_u"] ||
			Opcode(opcode) == instrToOpcode["f32.demote_f64"] ||
			Opcode(opcode) == instrToOpcode["f64.convert_i32_s"] ||
			Opcode(opcode) == instrToOpcode["f64.convert_i32_u"] ||
			Opcode(opcode) == instrToOpcode["f64.convert_i64_s"] ||
			Opcode(opcode) == instrToOpcode["f64.convert_i64_u"] ||
			Opcode(opcode) == instrToOpcode["f64.promote_f32"] ||
			Opcode(opcode) == instrToOpcode["i32.reinterpret_f32"] ||
			Opcode(opcode) == instrToOpcode["i64.reinterpret_f64"] ||
			Opcode(opcode) == instrToOpcode["f32.reinterpret_i32"] ||
			Opcode(opcode) == instrToOpcode["f64.reinterpret_i64"] ||

			Opcode(opcode) == instrToOpcode["i32.extend8_s"] ||
			Opcode(opcode) == instrToOpcode["i32.extend16_s"] ||
			Opcode(opcode) == instrToOpcode["i64.extend8_s"] ||
			Opcode(opcode) == instrToOpcode["i64.extend16_s"] ||
			Opcode(opcode) == instrToOpcode["i64.extend32_s"] {
			exps = append(exps,
				&Expression{
					PC:     pc + uint64(opptr),
					Opcode: Opcode(opcode),
				})

		} else if Opcode(opcode) == instrToOpcode["br_table"] {
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
		} else if Opcode(opcode) == instrToOpcode["br"] ||
			Opcode(opcode) == instrToOpcode["br_if"] {
			val, l := binary.Uvarint(data[ptr:])
			ptr += l
			exps = append(exps,
				&Expression{
					PC:         pc + uint64(opptr),
					Opcode:     Opcode(opcode),
					LabelIndex: int(val),
				})
		} else if Opcode(opcode) == instrToOpcode["i32.load"] ||
			Opcode(opcode) == instrToOpcode["i64.load"] ||
			Opcode(opcode) == instrToOpcode["f32.load"] ||
			Opcode(opcode) == instrToOpcode["f64.load"] ||
			Opcode(opcode) == instrToOpcode["i32.load8_s"] ||
			Opcode(opcode) == instrToOpcode["i32.load8_u"] ||
			Opcode(opcode) == instrToOpcode["i32.load16_s"] ||
			Opcode(opcode) == instrToOpcode["i32.load16_u"] ||
			Opcode(opcode) == instrToOpcode["i64.load8_s"] ||
			Opcode(opcode) == instrToOpcode["i64.load8_u"] ||
			Opcode(opcode) == instrToOpcode["i64.load16_s"] ||
			Opcode(opcode) == instrToOpcode["i64.load16_u"] ||
			Opcode(opcode) == instrToOpcode["i64.load32_s"] ||
			Opcode(opcode) == instrToOpcode["i64.load32_u"] ||
			Opcode(opcode) == instrToOpcode["i32.store"] ||
			Opcode(opcode) == instrToOpcode["i64.store"] ||
			Opcode(opcode) == instrToOpcode["f32.store"] ||
			Opcode(opcode) == instrToOpcode["f64.store"] ||
			Opcode(opcode) == instrToOpcode["i32.store8"] ||
			Opcode(opcode) == instrToOpcode["i32.store16"] ||
			Opcode(opcode) == instrToOpcode["i64.store8"] ||
			Opcode(opcode) == instrToOpcode["i64.store16"] ||
			Opcode(opcode) == instrToOpcode["i64.store32"] {
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
		} else if Opcode(opcode) == instrToOpcode["memory.size"] ||
			Opcode(opcode) == instrToOpcode["memory.grow"] {
			//				memoryIndex := data[ptr]
			ptr++
			exps = append(exps,
				&Expression{
					PC:     pc + uint64(opptr),
					Opcode: Opcode(opcode),
				})
		} else if Opcode(opcode) == instrToOpcode["block"] ||
			Opcode(opcode) == instrToOpcode["if"] ||
			Opcode(opcode) == instrToOpcode["loop"] {
			// Read the blocktype, and then read Expression
			valType := data[ptr]
			ptr++

			ex, l := NewExpression(data[ptr:], pc+uint64(ptr))
			ptr += l

			exps = append(exps,
				&Expression{
					PC:              pc + uint64(opptr),
					Opcode:          Opcode(opcode),
					Result:          ValType(valType),
					InnerExpression: ex,
				})

		} else if Opcode(opcode) == instrToOpcode["i32.const"] {
			val, l := DecodeSleb128(data[ptr:])
			ptr += int(l)
			exps = append(exps,
				&Expression{
					PC:       pc + uint64(opptr),
					Opcode:   Opcode(opcode),
					I32Value: int32(val),
				})
		} else if Opcode(opcode) == instrToOpcode["i64.const"] {
			val, l := DecodeSleb128(data[ptr:])
			ptr += int(l)
			exps = append(exps,
				&Expression{
					PC:       pc + uint64(opptr),
					Opcode:   Opcode(opcode),
					I64Value: int64(val),
				})
		} else if Opcode(opcode) == instrToOpcode["f32.const"] {
			ival := binary.LittleEndian.Uint32(data[ptr : ptr+4])
			val := math.Float32frombits(ival)
			ptr += 4
			exps = append(exps,
				&Expression{
					PC:       pc + uint64(opptr),
					Opcode:   Opcode(opcode),
					F32Value: float32(val),
				})
		} else if Opcode(opcode) == instrToOpcode["f64.const"] {
			ival := binary.LittleEndian.Uint64(data[ptr : ptr+8])
			val := math.Float64frombits(ival)
			ptr += 8
			exps = append(exps,
				&Expression{
					PC:       pc + uint64(opptr),
					Opcode:   Opcode(opcode),
					F64Value: float64(val),
				})
		} else if Opcode(opcode) == instrToOpcode["local.get"] ||
			Opcode(opcode) == instrToOpcode["local.set"] ||
			Opcode(opcode) == instrToOpcode["local.tee"] {
			val, l := binary.Uvarint(data[ptr:])
			ptr += l
			exps = append(exps,
				&Expression{
					PC:         pc + uint64(opptr),
					Opcode:     Opcode(opcode),
					LocalIndex: int(val),
				})
		} else if Opcode(opcode) == instrToOpcode["global.get"] ||
			Opcode(opcode) == instrToOpcode["global.set"] {
			val, l := binary.Uvarint(data[ptr:])
			ptr += l
			exps = append(exps,
				&Expression{
					PC:          pc + uint64(opptr),
					Opcode:      Opcode(opcode),
					GlobalIndex: int(val),
				})
		} else if Opcode(opcode) == instrToOpcode["call"] {
			val, l := binary.Uvarint(data[ptr:])
			ptr += l
			exps = append(exps,
				&Expression{
					PC:        pc + uint64(opptr),
					Opcode:    Opcode(opcode),
					FuncIndex: int(val),
				})
		} else if Opcode(opcode) == instrToOpcode["call_indirect"] {
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
				panic(fmt.Sprintf("Unsupported opcode 0xfc %d", opcode2))
			}

		} else if Opcode(opcode) == instrToOpcode["end"] {
			return exps, ptr
		} else {
			ptr--
			panic(fmt.Sprintf("TODO: Expression %x\n", data[ptr:ptr+16]))
		}
	}
}
