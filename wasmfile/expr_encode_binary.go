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
	"io"
	"math"
)

func (e *Expression) EncodeBinary(w io.Writer) error {

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
		_, err := w.Write([]byte{byte(e.Opcode)})
		if err != nil {
			return err
		}
		err = writeUvarint(w, uint64(len(e.Labels)))
		if err != nil {
			return err
		}
		for _, l := range e.Labels {
			err = writeUvarint(w, uint64(l))
			if err != nil {
				return err
			}
		}
		return writeUvarint(w, uint64(e.LabelIndex))
	} else if e.Opcode == instrToOpcode["br"] ||
		e.Opcode == instrToOpcode["br_if"] {
		_, err := w.Write([]byte{byte(e.Opcode)})
		if err != nil {
			return err
		}
		return writeUvarint(w, uint64(e.LabelIndex))
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
		_, err := w.Write([]byte{byte(e.Opcode)})
		if err != nil {
			return err
		}
		err = writeUvarint(w, uint64(e.MemAlign))
		return writeUvarint(w, uint64(e.MemOffset))
	} else if e.Opcode == instrToOpcode["memory.size"] ||
		e.Opcode == instrToOpcode["memory.grow"] {
		_, err := w.Write([]byte{byte(e.Opcode)})
		if err != nil {
			return err
		}

		_, err = w.Write([]byte{byte(0x00)})
		return err
	} else if e.Opcode == instrToOpcode["block"] ||
		e.Opcode == instrToOpcode["if"] ||
		e.Opcode == instrToOpcode["loop"] {
		_, err := w.Write([]byte{byte(e.Opcode), byte(e.Result)})
		if err != nil {
			return err
		}

		for _, ie := range e.InnerExpression {
			err = ie.EncodeBinary(w)
			if err != nil {
				return err
			}
		}

		_, err = w.Write([]byte{byte(instrToOpcode["end"])})
		return err
	} else if e.Opcode == instrToOpcode["i32.const"] {
		_, err := w.Write([]byte{byte(e.Opcode)})
		if err != nil {
			return err
		}
		return writeVarint(w, int64(e.I32Value))
	} else if e.Opcode == instrToOpcode["i64.const"] {
		_, err := w.Write([]byte{byte(e.Opcode)})
		if err != nil {
			return err
		}
		return writeVarint(w, e.I64Value)
	} else if e.Opcode == instrToOpcode["f32.const"] {
		_, err := w.Write([]byte{byte(e.Opcode)})
		if err != nil {
			return err
		}
		ival := math.Float32bits(e.F32Value)
		b := binary.LittleEndian.AppendUint32(make([]byte, 0), ival)
		_, err = w.Write(b)
		return err

	} else if e.Opcode == instrToOpcode["f64.const"] {
		_, err := w.Write([]byte{byte(e.Opcode)})
		if err != nil {
			return err
		}
		ival := math.Float64bits(e.F64Value)
		b := binary.LittleEndian.AppendUint64(make([]byte, 0), ival)
		_, err = w.Write(b)
		return err
	} else if e.Opcode == instrToOpcode["local.get"] ||
		e.Opcode == instrToOpcode["local.set"] ||
		e.Opcode == instrToOpcode["local.tee"] {
		_, err := w.Write([]byte{byte(e.Opcode)})
		if err != nil {
			return err
		}
		return writeUvarint(w, uint64(e.LocalIndex))
	} else if e.Opcode == instrToOpcode["global.get"] ||
		e.Opcode == instrToOpcode["global.set"] {
		_, err := w.Write([]byte{byte(e.Opcode)})
		if err != nil {
			return err
		}
		return writeUvarint(w, uint64(e.GlobalIndex))
	} else if e.Opcode == instrToOpcode["call"] {
		_, err := w.Write([]byte{byte(e.Opcode)})
		if err != nil {
			return err
		}
		return writeUvarint(w, uint64(e.FuncIndex))
	} else if e.Opcode == instrToOpcode["call_indirect"] {
		_, err := w.Write([]byte{byte(e.Opcode)})
		if err != nil {
			return err
		}
		err = writeUvarint(w, uint64(e.TypeIndex))
		if err != nil {
			return err
		}
		return writeUvarint(w, uint64(e.TableIndex))
	} else if e.Opcode == ExtendedOpcodeFC {
		_, err := w.Write([]byte{byte(e.Opcode)})
		if err != nil {
			return err
		}
		err = writeUvarint(w, uint64(e.OpcodeExt))
		if err != nil {
			return err
		}

		// Now deal with opcode2...
		if e.OpcodeExt == instrToOpcodeFC["memory.copy"] {
			_, err := w.Write([]byte{byte(0), byte(0)})
			return err
		} else if e.OpcodeExt == instrToOpcodeFC["memory.fill"] {
			_, err := w.Write([]byte{byte(0)})
			return err
		} else if e.OpcodeExt == instrToOpcodeFC["i32.trunc_sat_f32_s"] ||
			e.OpcodeExt == instrToOpcodeFC["i32.trunc_sat_f32_u"] ||
			e.OpcodeExt == instrToOpcodeFC["i32.trunc_sat_f64_s"] ||
			e.OpcodeExt == instrToOpcodeFC["i32.trunc_sat_f64_u"] ||
			e.OpcodeExt == instrToOpcodeFC["i64.trunc_sat_f32_s"] ||
			e.OpcodeExt == instrToOpcodeFC["i64.trunc_sat_f32_u"] ||
			e.OpcodeExt == instrToOpcodeFC["i64.trunc_sat_f64_s"] ||
			e.OpcodeExt == instrToOpcodeFC["i64.trunc_sat_f64_u"] {
			return nil
		} else {
			return fmt.Errorf("Unsupported opcode 0xfc %d", e.OpcodeExt)
		}
	} else {
		return fmt.Errorf("Unsupported opcode %d", e.Opcode)
	}

}
