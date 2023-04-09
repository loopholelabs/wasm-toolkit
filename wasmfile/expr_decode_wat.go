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
	"fmt"
	"strings"
)

func (e *Expression) DecodeWat(s string) error {
	s = SkipComment(s)
	s = strings.Trim(s, Whitespace)

	opcode, s := ReadToken(s)

	// First deal with simple opcodes (No args)
	if opcode == "unreachable" ||
		opcode == "nop" ||
		opcode == "return" ||
		opcode == "drop" ||
		opcode == "select" ||
		opcode == "i32.eqz" ||
		opcode == "i32.eq" ||
		opcode == "i32.ne" ||
		opcode == "i32.lt_s" ||
		opcode == "i32.lt_u" ||
		opcode == "i32.gt_s" ||
		opcode == "i32.gt_u" ||
		opcode == "i32.le_s" ||
		opcode == "i32.le_u" ||
		opcode == "i32.ge_s" ||
		opcode == "i32.ge_u" ||
		opcode == "i64.eqz" ||
		opcode == "i64.eq" ||
		opcode == "i64.ne" ||
		opcode == "i64.lt_s" ||
		opcode == "i64.lt_u" ||
		opcode == "i64.gt_s" ||
		opcode == "i64.gt_u" ||
		opcode == "i64.le_s" ||
		opcode == "i64.le_u" ||
		opcode == "i64.ge_s" ||
		opcode == "i64.ge_u" ||
		opcode == "f32.eq" ||
		opcode == "f32.ne" ||
		opcode == "f32.lt" ||
		opcode == "f32.gt" ||
		opcode == "f32.le" ||
		opcode == "f32.ge" ||
		opcode == "f64.eq" ||
		opcode == "f64.ne" ||
		opcode == "f64.lt" ||
		opcode == "f64.gt" ||
		opcode == "f64.le" ||
		opcode == "f64.ge" ||

		opcode == "i32.clz" ||
		opcode == "i32.ctz" ||
		opcode == "i32.popcnt" ||
		opcode == "i32.add" ||
		opcode == "i32.sub" ||
		opcode == "i32.mul" ||
		opcode == "i32.div_s" ||
		opcode == "i32.div_u" ||
		opcode == "i32.rem_s" ||
		opcode == "i32.rem_u" ||
		opcode == "i32.and" ||
		opcode == "i32.or" ||
		opcode == "i32.xor" ||
		opcode == "i32.shl" ||
		opcode == "i32.shr_s" ||
		opcode == "i32.shr_u" ||
		opcode == "i32.rotl_s" ||
		opcode == "i32.rotr_u" ||

		opcode == "i64.clz" ||
		opcode == "i64.ctz" ||
		opcode == "i64.popcnt" ||
		opcode == "i64.add" ||
		opcode == "i64.sub" ||
		opcode == "i64.mul" ||
		opcode == "i64.div_s" ||
		opcode == "i64.div_u" ||
		opcode == "i64.rem_s" ||
		opcode == "i64.rem_u" ||
		opcode == "i64.and" ||
		opcode == "i64.or" ||
		opcode == "i64.xor" ||
		opcode == "i64.shl" ||
		opcode == "i64.shr_s" ||
		opcode == "i64.shr_u" ||
		opcode == "i64.rotl_s" ||
		opcode == "i64.rotr_u" ||

		opcode == "f32.abs" ||
		opcode == "f32.neg" ||
		opcode == "f32.ceil" ||
		opcode == "f32.floor" ||
		opcode == "f32.trunc" ||
		opcode == "f32.nearest" ||
		opcode == "f32.sqrt" ||
		opcode == "f32.add" ||
		opcode == "f32.sub" ||
		opcode == "f32.mul" ||
		opcode == "f32.div" ||
		opcode == "f32.min" ||
		opcode == "f32.max" ||
		opcode == "f32.copysign" ||

		opcode == "f64.abs" ||
		opcode == "f64.neg" ||
		opcode == "f64.ceil" ||
		opcode == "f64.floor" ||
		opcode == "f64.trunc" ||
		opcode == "f64.nearest" ||
		opcode == "f64.sqrt" ||
		opcode == "f64.add" ||
		opcode == "f64.sub" ||
		opcode == "f64.mul" ||
		opcode == "f64.div" ||
		opcode == "f64.min" ||
		opcode == "f64.max" ||
		opcode == "f64.copysign" ||

		opcode == "i32.wrap_i64" ||
		opcode == "i32.trunc_f32_s" ||
		opcode == "i32.trunc_f32_u" ||
		opcode == "i32.trunc_f64_s" ||
		opcode == "i32.trunc_f64_u" ||
		opcode == "i64.extend_i32_s" ||
		opcode == "i64.extend_i32_u" ||
		opcode == "i64.trunc_f32_s" ||
		opcode == "i64.trunc_f32_u" ||
		opcode == "i64.trunc_f64_s" ||
		opcode == "i64.trunc_f64_u" ||
		opcode == "f32.convert_i32_s" ||
		opcode == "f32.convert_i32_u" ||
		opcode == "f32.convert_i64_s" ||
		opcode == "f32.convert_i64_u" ||
		opcode == "f32.demote_f64" ||
		opcode == "f64.convert_i32_s" ||
		opcode == "f64.convert_i32_u" ||
		opcode == "f64.convert_i64_s" ||
		opcode == "f64.convert_i64_u" ||
		opcode == "f64.promote_f32" ||
		opcode == "i32.reinterpret_f32" ||
		opcode == "i64.reinterpret_f64" ||
		opcode == "f32.reinterpret_i32" ||
		opcode == "f64.reinterpret_i64" ||

		opcode == "i32.extend8_s" ||
		opcode == "i32.extend16_s" ||
		opcode == "i64.extend8_s" ||
		opcode == "i64.extend16_s" ||
		opcode == "i64.extend32_s" {

		e.Opcode = instrToOpcode[opcode]
		return nil
	} else if opcode == "br_table" {
		/*
			targets := ""
			for _, l := range e.Labels {
				targets = fmt.Sprintf("%s %d", targets, l)
			}
			defaultTarget := fmt.Sprintf(" %d", e.LabelIndex)
			_, err := wr.WriteString(fmt.Sprintf("%s%s%s%s%s\n", prefix, opcodeToInstr[e.Opcode], targets, defaultTarget, comment))
			return err
		*/
	} else if opcode == "br" ||
		opcode == "br_if" {
		/*
			target := fmt.Sprintf(" %d", e.LabelIndex)
			_, err := wr.WriteString(fmt.Sprintf("%s%s%s%s\n", prefix, opcodeToInstr[e.Opcode], target, comment))
			return err
		*/
	} else if opcode == "i32.load" ||
		opcode == "i64.load" ||
		opcode == "f32.load" ||
		opcode == "f64.load" ||
		opcode == "i32.load8_s" ||
		opcode == "i32.load8_u" ||
		opcode == "i32.load16_s" ||
		opcode == "i32.load16_u" ||
		opcode == "i64.load8_s" ||
		opcode == "i64.load8_u" ||
		opcode == "i64.load16_s" ||
		opcode == "i64.load16_u" ||
		opcode == "i64.load32_s" ||
		opcode == "i64.load32_u" ||
		opcode == "i32.store" ||
		opcode == "i64.store" ||
		opcode == "f32.store" ||
		opcode == "f64.store" ||
		opcode == "i32.store8" ||
		opcode == "i32.store16" ||
		opcode == "i64.store8" ||
		opcode == "i64.store16" ||
		opcode == "i64.store32" {
		/*
			modAlign := fmt.Sprintf(" align=%d", 1<<e.MemAlign)
			modOffset := fmt.Sprintf(" offset=%d", e.MemOffset)
			if e.MemOffset == 0 {
				modOffset = ""
			}
			// TODO: Default align?
			_, err := wr.WriteString(fmt.Sprintf("%s%s%s%s%s\n", prefix, opcodeToInstr[e.Opcode], modOffset, modAlign, comment))
			return err
		*/
	} else if opcode == "memory.size" ||
		opcode == "memory.grow" {
		/*
			_, err := wr.WriteString(fmt.Sprintf("%s%s%s\n", prefix, opcodeToInstr[e.Opcode], comment))
			return err
		*/
	} else if opcode == "block" ||
		opcode == "if" ||
		opcode == "loop" {
		/*
			result := ""
			if e.Result != ValNone {
				result = fmt.Sprintf(" (result %s)", byteToValType[e.Result])
			}

			_, err := wr.WriteString(fmt.Sprintf("%s%s%s%s\n", prefix, opcodeToInstr[e.Opcode], result, comment))

			for _, ie := range e.InnerExpression {
				err = ie.EncodeWat(wr, fmt.Sprintf("%s%s", prefix, "    "), wf)
				if err != nil {
					return err
				}
			}

			_, err = wr.WriteString(fmt.Sprintf("%s%s\n", prefix, "end"))
			return err
		*/
	} else if opcode == "i32.const" {
		/*
			value := fmt.Sprintf(" %d", e.I32Value)
			_, err := wr.WriteString(fmt.Sprintf("%s%s%s%s\n", prefix, opcodeToInstr[e.Opcode], value, comment))
			return err
		*/
	} else if opcode == "i64.const" {
		/*
			value := fmt.Sprintf(" %d", e.I64Value)
			_, err := wr.WriteString(fmt.Sprintf("%s%s%s%s\n", prefix, opcodeToInstr[e.Opcode], value, comment))
			return err
		*/
	} else if opcode == "f32.const" {
		/*
			value := fmt.Sprintf(" %f", e.F32Value)
			_, err := wr.WriteString(fmt.Sprintf("%s%s%s%s\n", prefix, opcodeToInstr[e.Opcode], value, comment))
			return err
		*/
	} else if opcode == "f64.const" {
		/*
			value := fmt.Sprintf(" %f", e.F64Value)
			_, err := wr.WriteString(fmt.Sprintf("%s%s%s%s\n", prefix, opcodeToInstr[e.Opcode], value, comment))
			return err
		*/
	} else if opcode == "local.get" ||
		opcode == "local.set" ||
		opcode == "local.tee" {
		/*
			tname := wf.GetLocalVarName(e.PC, e.LocalIndex)
			if tname != "" {
				comment = comment + " ;; Variable " + tname
			}
			localTarget := fmt.Sprintf(" %d", e.LocalIndex)
			_, err := wr.WriteString(fmt.Sprintf("%s%s%s%s\n", prefix, opcodeToInstr[e.Opcode], localTarget, comment))
			return err
		*/
	} else if opcode == "global.get" ||
		opcode == "global.set" {
		/*
			g := wf.GetGlobalIdentifier(e.GlobalIndex)
			globalTarget := fmt.Sprintf(" %s", g)
			_, err := wr.WriteString(fmt.Sprintf("%s%s%s%s\n", prefix, opcodeToInstr[e.Opcode], globalTarget, comment))
			return err
		*/
	} else if opcode == "call" {
		/*
			f := wf.GetFunctionIdentifier(e.FuncIndex)
			callTarget := fmt.Sprintf(" %s", f)
			_, err := wr.WriteString(fmt.Sprintf("%s%s%s%s\n", prefix, opcodeToInstr[e.Opcode], callTarget, comment))
			return err
		*/
	} else if opcode == "call_indirect" {
		/*
			typeIndex := fmt.Sprintf(" (type %d)", e.TypeIndex)
			_, err := wr.WriteString(fmt.Sprintf("%s%s%s%s\n", prefix, opcodeToInstr[e.Opcode], typeIndex, comment))
			return err
		*/
	} else if opcode == "memory.copy" {
		e.Opcode = ExtendedOpcodeFC
		e.OpcodeExt = instrToOpcodeFC[opcode]
	} else if opcode == "memory.fill" {
		e.Opcode = ExtendedOpcodeFC
		e.OpcodeExt = instrToOpcodeFC[opcode]
	} else if opcode == "i32.trunc_sat_f32_s" ||
		opcode == "i32.trunc_sat_f32_u" ||
		opcode == "i32.trunc_sat_f64_s" ||
		opcode == "i32.trunc_sat_f64_u" ||
		opcode == "i64.trunc_sat_f32_s" ||
		opcode == "i64.trunc_sat_f32_u" ||
		opcode == "i64.trunc_sat_f64_s" ||
		opcode == "i64.trunc_sat_f64_u" {
		e.Opcode = ExtendedOpcodeFC
		e.OpcodeExt = instrToOpcodeFC[opcode]
	} else {
		return fmt.Errorf("Unsupported opcode %s", opcode)
	}

	return nil
}
