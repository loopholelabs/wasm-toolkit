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
	"io"
	"math"

	"github.com/loopholelabs/wasm-toolkit/wasmfile/encoding"
)

func (e *Expression) EncodeBinary(w io.Writer) error {

	// First deal with simple opcodes (No args)
	if e.HasNoArgs() {
		_, err := w.Write([]byte{byte(e.Opcode)})
		return err
	} else if e.Opcode == InstrToOpcode["br_table"] {
		_, err := w.Write([]byte{byte(e.Opcode)})
		if err != nil {
			return err
		}
		err = encoding.WriteUvarint(w, uint64(len(e.Labels)))
		if err != nil {
			return err
		}
		for _, l := range e.Labels {
			err = encoding.WriteUvarint(w, uint64(l))
			if err != nil {
				return err
			}
		}
		return encoding.WriteUvarint(w, uint64(e.LabelIndex))
	} else if e.Opcode == InstrToOpcode["br"] ||
		e.Opcode == InstrToOpcode["br_if"] {
		_, err := w.Write([]byte{byte(e.Opcode)})
		if err != nil {
			return err
		}
		return encoding.WriteUvarint(w, uint64(e.LabelIndex))
	} else if e.HasMemoryArgs() {
		_, err := w.Write([]byte{byte(e.Opcode)})
		if err != nil {
			return err
		}
		err = encoding.WriteUvarint(w, uint64(e.MemAlign))
		return encoding.WriteUvarint(w, uint64(e.MemOffset))
	} else if e.Opcode == InstrToOpcode["memory.size"] ||
		e.Opcode == InstrToOpcode["memory.grow"] {
		_, err := w.Write([]byte{byte(e.Opcode)})
		if err != nil {
			return err
		}
		_, err = w.Write([]byte{byte(0x00)})
		return err
	} else if e.Opcode == InstrToOpcode["block"] ||
		e.Opcode == InstrToOpcode["if"] ||
		e.Opcode == InstrToOpcode["loop"] {
		_, err := w.Write([]byte{byte(e.Opcode), byte(e.Result)})
		if err != nil {
			return err
		}
		return err
	} else if e.Opcode == InstrToOpcode["i32.const"] {
		_, err := w.Write([]byte{byte(e.Opcode)})
		if err != nil {
			return err
		}
		return encoding.WriteVarint(w, int64(e.I32Value))
	} else if e.Opcode == InstrToOpcode["i64.const"] {
		_, err := w.Write([]byte{byte(e.Opcode)})
		if err != nil {
			return err
		}
		return encoding.WriteVarint(w, e.I64Value)
	} else if e.Opcode == InstrToOpcode["f32.const"] {
		_, err := w.Write([]byte{byte(e.Opcode)})
		if err != nil {
			return err
		}
		ival := math.Float32bits(e.F32Value)
		b := binary.LittleEndian.AppendUint32(make([]byte, 0), ival)
		_, err = w.Write(b)
		return err

	} else if e.Opcode == InstrToOpcode["f64.const"] {
		_, err := w.Write([]byte{byte(e.Opcode)})
		if err != nil {
			return err
		}
		ival := math.Float64bits(e.F64Value)
		b := binary.LittleEndian.AppendUint64(make([]byte, 0), ival)
		_, err = w.Write(b)
		return err
	} else if e.Opcode == InstrToOpcode["local.get"] ||
		e.Opcode == InstrToOpcode["local.set"] ||
		e.Opcode == InstrToOpcode["local.tee"] {
		_, err := w.Write([]byte{byte(e.Opcode)})
		if err != nil {
			return err
		}
		return encoding.WriteUvarint(w, uint64(e.LocalIndex))
	} else if e.Opcode == InstrToOpcode["global.get"] ||
		e.Opcode == InstrToOpcode["global.set"] {
		_, err := w.Write([]byte{byte(e.Opcode)})
		if err != nil {
			return err
		}
		return encoding.WriteUvarint(w, uint64(e.GlobalIndex))
	} else if e.Opcode == InstrToOpcode["call"] {
		_, err := w.Write([]byte{byte(e.Opcode)})
		if err != nil {
			return err
		}
		return encoding.WriteUvarint(w, uint64(e.FuncIndex))
	} else if e.Opcode == InstrToOpcode["call_indirect"] {
		_, err := w.Write([]byte{byte(e.Opcode)})
		if err != nil {
			return err
		}
		err = encoding.WriteUvarint(w, uint64(e.TypeIndex))
		if err != nil {
			return err
		}
		return encoding.WriteUvarint(w, uint64(e.TableIndex))
	} else if e.Opcode == ExtendedOpcodeFC {
		_, err := w.Write([]byte{byte(e.Opcode)})
		if err != nil {
			return err
		}
		err = encoding.WriteUvarint(w, uint64(e.OpcodeExt))
		if err != nil {
			return err
		}

		// Now deal with opcodeExt...
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
