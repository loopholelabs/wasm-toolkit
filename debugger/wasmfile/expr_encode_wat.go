package wasmfile

import (
	"io"
)

// TODO
func (e *Expression) EncodeWat(w io.Writer) error {

	// First deal with simple opcodes (No args)
	if e.Opcode == instrToOpcode["unreachable"] ||
		e.Opcode == instrToOpcode["nop"] ||
		e.Opcode == instrToOpcode["return"] ||
		e.Opcode == instrToOpcode["drop"] ||
		e.Opcode == instrToOpcode["select"] ||
		e.Opcode == instrToOpcode["i32.eqz"] ||
		e.Opcode == instrToOpcode["i32.eq"] ||
		e.Opcode == instrToOpcode["i32.ne"] ||
		e.Opcode == instrToOpcode["i32.lt_s"] ||
		e.Opcode == instrToOpcode["i32.lt_u"] ||
		e.Opcode == instrToOpcode["i32.gt_s"] ||
		e.Opcode == instrToOpcode["i32.gt_u"] ||
		e.Opcode == instrToOpcode["i32.le_s"] ||
		e.Opcode == instrToOpcode["i32.le_u"] ||
		e.Opcode == instrToOpcode["i32.ge_s"] ||
		e.Opcode == instrToOpcode["i32.ge_u"] ||
		e.Opcode == instrToOpcode["i64.eqz"] ||
		e.Opcode == instrToOpcode["i64.eq"] ||
		e.Opcode == instrToOpcode["i64.ne"] ||
		e.Opcode == instrToOpcode["i64.lt_s"] ||
		e.Opcode == instrToOpcode["i64.lt_u"] ||
		e.Opcode == instrToOpcode["i64.gt_s"] ||
		e.Opcode == instrToOpcode["i64.gt_u"] ||
		e.Opcode == instrToOpcode["i64.le_s"] ||
		e.Opcode == instrToOpcode["i64.le_u"] ||
		e.Opcode == instrToOpcode["i64.ge_s"] ||
		e.Opcode == instrToOpcode["i64.ge_u"] ||
		e.Opcode == instrToOpcode["f32.eq"] ||
		e.Opcode == instrToOpcode["f32.ne"] ||
		e.Opcode == instrToOpcode["f32.lt"] ||
		e.Opcode == instrToOpcode["f32.gt"] ||
		e.Opcode == instrToOpcode["f32.le"] ||
		e.Opcode == instrToOpcode["f32.ge"] ||
		e.Opcode == instrToOpcode["f64.eq"] ||
		e.Opcode == instrToOpcode["f64.ne"] ||
		e.Opcode == instrToOpcode["f64.lt"] ||
		e.Opcode == instrToOpcode["f64.gt"] ||
		e.Opcode == instrToOpcode["f64.le"] ||
		e.Opcode == instrToOpcode["f64.ge"] ||

		e.Opcode == instrToOpcode["i32.clz"] ||
		e.Opcode == instrToOpcode["i32.ctz"] ||
		e.Opcode == instrToOpcode["i32.popcnt"] ||
		e.Opcode == instrToOpcode["i32.add"] ||
		e.Opcode == instrToOpcode["i32.sub"] ||
		e.Opcode == instrToOpcode["i32.mul"] ||
		e.Opcode == instrToOpcode["i32.div_s"] ||
		e.Opcode == instrToOpcode["i32.div_u"] ||
		e.Opcode == instrToOpcode["i32.rem_s"] ||
		e.Opcode == instrToOpcode["i32.rem_u"] ||
		e.Opcode == instrToOpcode["i32.and"] ||
		e.Opcode == instrToOpcode["i32.or"] ||
		e.Opcode == instrToOpcode["i32.xor"] ||
		e.Opcode == instrToOpcode["i32.shl"] ||
		e.Opcode == instrToOpcode["i32.shr_s"] ||
		e.Opcode == instrToOpcode["i32.shr_u"] ||
		e.Opcode == instrToOpcode["i32.rotl_s"] ||
		e.Opcode == instrToOpcode["i32.rotr_u"] ||

		e.Opcode == instrToOpcode["i64.clz"] ||
		e.Opcode == instrToOpcode["i64.ctz"] ||
		e.Opcode == instrToOpcode["i64.popcnt"] ||
		e.Opcode == instrToOpcode["i64.add"] ||
		e.Opcode == instrToOpcode["i64.sub"] ||
		e.Opcode == instrToOpcode["i64.mul"] ||
		e.Opcode == instrToOpcode["i64.div_s"] ||
		e.Opcode == instrToOpcode["i64.div_u"] ||
		e.Opcode == instrToOpcode["i64.rem_s"] ||
		e.Opcode == instrToOpcode["i64.rem_u"] ||
		e.Opcode == instrToOpcode["i64.and"] ||
		e.Opcode == instrToOpcode["i64.or"] ||
		e.Opcode == instrToOpcode["i64.xor"] ||
		e.Opcode == instrToOpcode["i64.shl"] ||
		e.Opcode == instrToOpcode["i64.shr_s"] ||
		e.Opcode == instrToOpcode["i64.shr_u"] ||
		e.Opcode == instrToOpcode["i64.rotl_s"] ||
		e.Opcode == instrToOpcode["i64.rotr_u"] ||

		e.Opcode == instrToOpcode["f32.abs"] ||
		e.Opcode == instrToOpcode["f32.neg"] ||
		e.Opcode == instrToOpcode["f32.ceil"] ||
		e.Opcode == instrToOpcode["f32.floor"] ||
		e.Opcode == instrToOpcode["f32.trunc"] ||
		e.Opcode == instrToOpcode["f32.nearest"] ||
		e.Opcode == instrToOpcode["f32.sqrt"] ||
		e.Opcode == instrToOpcode["f32.add"] ||
		e.Opcode == instrToOpcode["f32.sub"] ||
		e.Opcode == instrToOpcode["f32.mul"] ||
		e.Opcode == instrToOpcode["f32.div"] ||
		e.Opcode == instrToOpcode["f32.min"] ||
		e.Opcode == instrToOpcode["f32.max"] ||
		e.Opcode == instrToOpcode["f32.copysign"] ||

		e.Opcode == instrToOpcode["f64.abs"] ||
		e.Opcode == instrToOpcode["f64.neg"] ||
		e.Opcode == instrToOpcode["f64.ceil"] ||
		e.Opcode == instrToOpcode["f64.floor"] ||
		e.Opcode == instrToOpcode["f64.trunc"] ||
		e.Opcode == instrToOpcode["f64.nearest"] ||
		e.Opcode == instrToOpcode["f64.sqrt"] ||
		e.Opcode == instrToOpcode["f64.add"] ||
		e.Opcode == instrToOpcode["f64.sub"] ||
		e.Opcode == instrToOpcode["f64.mul"] ||
		e.Opcode == instrToOpcode["f64.div"] ||
		e.Opcode == instrToOpcode["f64.min"] ||
		e.Opcode == instrToOpcode["f64.max"] ||
		e.Opcode == instrToOpcode["f64.copysign"] ||

		e.Opcode == instrToOpcode["i32.wrap_i64"] ||
		e.Opcode == instrToOpcode["i32.trunc_f32_s"] ||
		e.Opcode == instrToOpcode["i32.trunc_f32_u"] ||
		e.Opcode == instrToOpcode["i32.trunc_f64_s"] ||
		e.Opcode == instrToOpcode["i32.trunc_f64_u"] ||
		e.Opcode == instrToOpcode["i64.extend_i32_s"] ||
		e.Opcode == instrToOpcode["i64.extend_i32_u"] ||
		e.Opcode == instrToOpcode["i64.trunc_f32_s"] ||
		e.Opcode == instrToOpcode["i64.trunc_f32_u"] ||
		e.Opcode == instrToOpcode["i64.trunc_f64_s"] ||
		e.Opcode == instrToOpcode["i64.trunc_f64_u"] ||
		e.Opcode == instrToOpcode["f32.convert_i32_s"] ||
		e.Opcode == instrToOpcode["f32.convert_i32_u"] ||
		e.Opcode == instrToOpcode["f32.convert_i64_s"] ||
		e.Opcode == instrToOpcode["f32.convert_i64_u"] ||
		e.Opcode == instrToOpcode["f32.demote_f64"] ||
		e.Opcode == instrToOpcode["f64.convert_i32_s"] ||
		e.Opcode == instrToOpcode["f64.convert_i32_u"] ||
		e.Opcode == instrToOpcode["f64.convert_i64_s"] ||
		e.Opcode == instrToOpcode["f64.convert_i64_u"] ||
		e.Opcode == instrToOpcode["f64.promote_f32"] ||
		e.Opcode == instrToOpcode["i32.reinterpret_f32"] ||
		e.Opcode == instrToOpcode["i64.reinterpret_f64"] ||
		e.Opcode == instrToOpcode["f32.reinterpret_i32"] ||
		e.Opcode == instrToOpcode["f64.reinterpret_i64"] ||

		e.Opcode == instrToOpcode["i32.extend8_s"] ||
		e.Opcode == instrToOpcode["i32.extend16_s"] ||
		e.Opcode == instrToOpcode["i64.extend8_s"] ||
		e.Opcode == instrToOpcode["i64.extend16_s"] ||
		e.Opcode == instrToOpcode["i64.extend32_s"] {

		_, err := w.Write([]byte{byte(e.Opcode)})
		return err
	} else if e.Opcode == instrToOpcode["br_table"] {
		/*
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
		*/
		panic("TODO")
	} else if e.Opcode == instrToOpcode["br"] ||
		e.Opcode == instrToOpcode["br_if"] {
		/*
			val, l := binary.Uvarint(data[ptr:])
			ptr += l
			exps = append(exps,
				&Expression{
					PC:         pc + uint64(opptr),
					Opcode:     Opcode(opcode),
					LabelIndex: int(val),
				})
		*/
		panic("TODO")
	} else if e.Opcode == instrToOpcode["i32.load"] ||
		e.Opcode == instrToOpcode["i64.load"] ||
		e.Opcode == instrToOpcode["f32.load"] ||
		e.Opcode == instrToOpcode["f64.load"] ||
		e.Opcode == instrToOpcode["i32.load8_s"] ||
		e.Opcode == instrToOpcode["i32.load8_u"] ||
		e.Opcode == instrToOpcode["i32.load16_s"] ||
		e.Opcode == instrToOpcode["i32.load16_u"] ||
		e.Opcode == instrToOpcode["i64.load8_s"] ||
		e.Opcode == instrToOpcode["i64.load8_u"] ||
		e.Opcode == instrToOpcode["i64.load16_s"] ||
		e.Opcode == instrToOpcode["i64.load16_u"] ||
		e.Opcode == instrToOpcode["i64.load32_s"] ||
		e.Opcode == instrToOpcode["i64.load32_u"] ||
		e.Opcode == instrToOpcode["i32.store"] ||
		e.Opcode == instrToOpcode["i64.store"] ||
		e.Opcode == instrToOpcode["f32.store"] ||
		e.Opcode == instrToOpcode["f64.store"] ||
		e.Opcode == instrToOpcode["i32.store8"] ||
		e.Opcode == instrToOpcode["i32.store16"] ||
		e.Opcode == instrToOpcode["i64.store8"] ||
		e.Opcode == instrToOpcode["i64.store16"] ||
		e.Opcode == instrToOpcode["i64.store32"] {
		/*
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
		*/
		panic("TODO")
	} else if e.Opcode == instrToOpcode["memory.size"] ||
		e.Opcode == instrToOpcode["memory.grow"] {
		/*
			//				memoryIndex := data[ptr]
			ptr++
			exps = append(exps,
				&Expression{
					PC:     pc + uint64(opptr),
					Opcode: Opcode(opcode),
				})
		*/
		panic("TODO")
	} else if e.Opcode == instrToOpcode["block"] ||
		e.Opcode == instrToOpcode["if"] ||
		e.Opcode == instrToOpcode["loop"] {
		/*
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
		*/
		panic("TODO")
	} else if e.Opcode == instrToOpcode["i32.const"] {
		/*
			val, l := DecodeSleb128(data[ptr:])
			ptr += int(l)
			exps = append(exps,
				&Expression{
					PC:       pc + uint64(opptr),
					Opcode:   Opcode(opcode),
					I32Value: int32(val),
				})
		*/
		panic("TODO")
	} else if e.Opcode == instrToOpcode["i64.const"] {
		/*
			val, l := DecodeSleb128(data[ptr:])
			ptr += int(l)
			exps = append(exps,
				&Expression{
					PC:       pc + uint64(opptr),
					Opcode:   Opcode(opcode),
					I64Value: int64(val),
				})
		*/
		panic("TODO")
	} else if e.Opcode == instrToOpcode["f32.const"] {
		/*
			ival := binary.LittleEndian.Uint32(data[ptr : ptr+4])
			val := math.Float32frombits(ival)
			ptr += 4
			exps = append(exps,
				&Expression{
					PC:       pc + uint64(opptr),
					Opcode:   Opcode(opcode),
					F32Value: float32(val),
				})
		*/
		panic("TODO")
	} else if e.Opcode == instrToOpcode["f64.const"] {
		/*
			ival := binary.LittleEndian.Uint64(data[ptr : ptr+8])
			val := math.Float64frombits(ival)
			ptr += 8
			exps = append(exps,
				&Expression{
					PC:       pc + uint64(opptr),
					Opcode:   Opcode(opcode),
					F64Value: float64(val),
				})
		*/
		panic("TODO")
	} else if e.Opcode == instrToOpcode["local.get"] ||
		e.Opcode == instrToOpcode["local.set"] ||
		e.Opcode == instrToOpcode["local.tee"] {
		/*
			val, l := binary.Uvarint(data[ptr:])
			ptr += l
			exps = append(exps,
				&Expression{
					PC:         pc + uint64(opptr),
					Opcode:     Opcode(opcode),
					LocalIndex: int(val),
				})
		*/
		panic("TODO")
	} else if e.Opcode == instrToOpcode["global.get"] ||
		e.Opcode == instrToOpcode["global.set"] {
		/*
			val, l := binary.Uvarint(data[ptr:])
			ptr += l
			exps = append(exps,
				&Expression{
					PC:          pc + uint64(opptr),
					Opcode:      Opcode(opcode),
					GlobalIndex: int(val),
				})
		*/
		panic("TODO")
	} else if e.Opcode == instrToOpcode["call"] {
		/*
			val, l := binary.Uvarint(data[ptr:])
			ptr += l
			exps = append(exps,
				&Expression{
					PC:        pc + uint64(opptr),
					Opcode:    Opcode(opcode),
					FuncIndex: int(val),
				})
		*/
		panic("TODO")
	} else if e.Opcode == instrToOpcode["call_indirect"] {
		/*
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
		*/
		panic("TODO")
	} else if e.Opcode == ExtendedOpcodeFC {
		/*
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
		*/
		panic("TODO")
	} else {
		panic("TODO UNKNOWN")
	}

}
