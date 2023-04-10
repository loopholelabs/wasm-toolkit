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
	"bufio"
	"fmt"
	"io"
)

func (e *Expression) EncodeWat(w io.Writer, prefix string, wf *WasmFile) error {
	comment := "" //fmt.Sprintf("    ;; PC=%d", e.PC) // TODO From line numbers, vars etc

	lineNumberData := wf.GetLineNumberInfo(e.PC)
	if lineNumberData != "" {
		comment = fmt.Sprintf(" ;; Src = %s", lineNumberData)
	}

	wr := bufio.NewWriter(w)

	defer func() {
		wr.Flush()
	}()

	// First deal with simple opcodes (No args)
	if e.Opcode == instrToOpcode["unreachable"] ||
		e.Opcode == instrToOpcode["nop"] ||
		e.Opcode == instrToOpcode["return"] ||
		e.Opcode == instrToOpcode["drop"] ||
		e.Opcode == instrToOpcode["select"] ||
		e.Opcode == instrToOpcode["end"] ||
		e.Opcode == instrToOpcode["else"] ||
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

		_, err := wr.WriteString(fmt.Sprintf("%s%s%s\n", prefix, opcodeToInstr[e.Opcode], comment))
		return err
	} else if e.Opcode == instrToOpcode["br_table"] {
		targets := ""
		for _, l := range e.Labels {
			targets = fmt.Sprintf("%s %d", targets, l)
		}
		defaultTarget := fmt.Sprintf(" %d", e.LabelIndex)
		_, err := wr.WriteString(fmt.Sprintf("%s%s%s%s%s\n", prefix, opcodeToInstr[e.Opcode], targets, defaultTarget, comment))
		return err
	} else if e.Opcode == instrToOpcode["br"] ||
		e.Opcode == instrToOpcode["br_if"] {
		target := fmt.Sprintf(" %d", e.LabelIndex)
		_, err := wr.WriteString(fmt.Sprintf("%s%s%s%s\n", prefix, opcodeToInstr[e.Opcode], target, comment))
		return err
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
		modAlign := fmt.Sprintf(" align=%d", 1<<e.MemAlign)
		modOffset := fmt.Sprintf(" offset=%d", e.MemOffset)
		if e.MemOffset == 0 {
			modOffset = ""
		}
		// TODO: Default align?
		/*
			if e.MemAlign == 0 {
				modAlign = ""
			}
		*/
		_, err := wr.WriteString(fmt.Sprintf("%s%s%s%s%s\n", prefix, opcodeToInstr[e.Opcode], modOffset, modAlign, comment))
		return err
	} else if e.Opcode == instrToOpcode["memory.size"] ||
		e.Opcode == instrToOpcode["memory.grow"] {
		_, err := wr.WriteString(fmt.Sprintf("%s%s%s\n", prefix, opcodeToInstr[e.Opcode], comment))
		return err
	} else if e.Opcode == instrToOpcode["block"] ||
		e.Opcode == instrToOpcode["if"] ||
		e.Opcode == instrToOpcode["loop"] {

		result := ""
		if e.Result != ValNone {
			result = fmt.Sprintf(" (result %s)", byteToValType[e.Result])
		}

		_, err := wr.WriteString(fmt.Sprintf("%s%s%s%s\n", prefix, opcodeToInstr[e.Opcode], result, comment))

		if e.InnerExpression != nil {
			for _, ie := range e.InnerExpression {
				err = ie.EncodeWat(wr, fmt.Sprintf("%s%s", prefix, "    "), wf)
				if err != nil {
					return err
				}
			}

			_, err = wr.WriteString(fmt.Sprintf("%s%s\n", prefix, "end"))
		}
		return err
	} else if e.Opcode == instrToOpcode["i32.const"] {
		value := fmt.Sprintf(" %d", e.I32Value)
		_, err := wr.WriteString(fmt.Sprintf("%s%s%s%s\n", prefix, opcodeToInstr[e.Opcode], value, comment))
		return err
	} else if e.Opcode == instrToOpcode["i64.const"] {
		value := fmt.Sprintf(" %d", e.I64Value)
		_, err := wr.WriteString(fmt.Sprintf("%s%s%s%s\n", prefix, opcodeToInstr[e.Opcode], value, comment))
		return err
	} else if e.Opcode == instrToOpcode["f32.const"] {
		value := fmt.Sprintf(" %f", e.F32Value)
		_, err := wr.WriteString(fmt.Sprintf("%s%s%s%s\n", prefix, opcodeToInstr[e.Opcode], value, comment))
		return err
	} else if e.Opcode == instrToOpcode["f64.const"] {
		value := fmt.Sprintf(" %f", e.F64Value)
		_, err := wr.WriteString(fmt.Sprintf("%s%s%s%s\n", prefix, opcodeToInstr[e.Opcode], value, comment))
		return err
	} else if e.Opcode == instrToOpcode["local.get"] ||
		e.Opcode == instrToOpcode["local.set"] ||
		e.Opcode == instrToOpcode["local.tee"] {
		tname := wf.GetLocalVarName(e.PC, e.LocalIndex)
		if tname != "" {
			comment = comment + " ;; Variable " + tname
		}
		localTarget := fmt.Sprintf(" %d", e.LocalIndex)
		_, err := wr.WriteString(fmt.Sprintf("%s%s%s%s\n", prefix, opcodeToInstr[e.Opcode], localTarget, comment))
		return err
	} else if e.Opcode == instrToOpcode["global.get"] ||
		e.Opcode == instrToOpcode["global.set"] {
		g := wf.GetGlobalIdentifier(e.GlobalIndex)
		globalTarget := fmt.Sprintf(" %s", g)
		_, err := wr.WriteString(fmt.Sprintf("%s%s%s%s\n", prefix, opcodeToInstr[e.Opcode], globalTarget, comment))
		return err
	} else if e.Opcode == instrToOpcode["call"] {
		f := wf.GetFunctionIdentifier(e.FuncIndex)
		callTarget := fmt.Sprintf(" %s", f)
		_, err := wr.WriteString(fmt.Sprintf("%s%s%s%s\n", prefix, opcodeToInstr[e.Opcode], callTarget, comment))
		return err
	} else if e.Opcode == instrToOpcode["call_indirect"] {
		typeIndex := fmt.Sprintf(" (type %d)", e.TypeIndex)
		_, err := wr.WriteString(fmt.Sprintf("%s%s%s%s\n", prefix, opcodeToInstr[e.Opcode], typeIndex, comment))
		return err
	} else if e.Opcode == ExtendedOpcodeFC {
		// Now deal with opcode2...
		if e.OpcodeExt == instrToOpcodeFC["memory.copy"] {
			_, err := wr.WriteString(fmt.Sprintf("%s%s%s\n", prefix, opcodeToInstrFC[e.OpcodeExt], comment))
			return err
		} else if e.OpcodeExt == instrToOpcodeFC["memory.fill"] {
			_, err := wr.WriteString(fmt.Sprintf("%s%s%s\n", prefix, opcodeToInstrFC[e.OpcodeExt], comment))
			return err
		} else if e.OpcodeExt == instrToOpcodeFC["i32.trunc_sat_f32_s"] ||
			e.OpcodeExt == instrToOpcodeFC["i32.trunc_sat_f32_u"] ||
			e.OpcodeExt == instrToOpcodeFC["i32.trunc_sat_f64_s"] ||
			e.OpcodeExt == instrToOpcodeFC["i32.trunc_sat_f64_u"] ||
			e.OpcodeExt == instrToOpcodeFC["i64.trunc_sat_f32_s"] ||
			e.OpcodeExt == instrToOpcodeFC["i64.trunc_sat_f32_u"] ||
			e.OpcodeExt == instrToOpcodeFC["i64.trunc_sat_f64_s"] ||
			e.OpcodeExt == instrToOpcodeFC["i64.trunc_sat_f64_u"] {
			_, err := wr.WriteString(fmt.Sprintf("%s%s%s\n", prefix, opcodeToInstrFC[e.OpcodeExt], comment))
			return err
		} else {
			return fmt.Errorf("Unsupported opcode 0xfc %d", e.OpcodeExt)
		}
	} else {
		return fmt.Errorf("Unsupported opcode %d", e.Opcode)
	}

}
