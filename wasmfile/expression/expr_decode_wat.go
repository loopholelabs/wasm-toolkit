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
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/loopholelabs/wasm-toolkit/wasmfile/encoding"
	"github.com/loopholelabs/wasm-toolkit/wasmfile/types"
)

func (e *Expression) DecodeWat(s string, localNames map[string]int) error {
	s = encoding.SkipComment(s)
	s = strings.Trim(s, encoding.Whitespace)

	opcode, s := encoding.ReadToken(s)

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
		opcode == "i32.rotl" ||
		opcode == "i32.rotr" ||

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
		opcode == "i64.rotl" ||
		opcode == "i64.rotr" ||

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

		e.Opcode = InstrToOpcode[opcode]
		return nil
	} else if opcode == "br_table" {
		e.Opcode = InstrToOpcode[opcode]
		e.Labels = make([]int, 0)
		var err error
		var br_target string
		var li int
		for {
			s = strings.Trim(s, encoding.Whitespace)
			if len(s) == 0 || strings.HasPrefix(s, ";;") {
				break
			}
			br_target, s = encoding.ReadToken(s)
			li, err = strconv.Atoi(br_target)
			if err != nil {
				return err
			}
			e.Labels = append(e.Labels, li)
		}

		// Remove the last label and put it into default
		e.LabelIndex = e.Labels[len(e.Labels)-1]
		e.Labels = e.Labels[:len(e.Labels)-1]
		return nil
	} else if opcode == "br" ||
		opcode == "br_if" {
		e.Opcode = InstrToOpcode[opcode]
		var err error
		var br_target string
		br_target, s = encoding.ReadToken(s)
		e.LabelIndex, err = strconv.Atoi(br_target)
		if err != nil {
			return err
		}
		return nil
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
		e.Opcode = InstrToOpcode[opcode]
		for {
			var t string
			s = strings.Trim(s, encoding.Whitespace)
			if len(s) == 0 {
				break
			}
			if strings.HasPrefix(s, ";;") {
				break
			}
			t, s = encoding.ReadToken(s)
			// Optional align=<V>
			// Optional offset=<V>
			if strings.HasPrefix(t, "align=") {
				v, err := strconv.Atoi(t[6:])
				if err != nil {
					return err
				}
				if v == 1 {
					e.MemAlign = 0
				} else if v == 2 {
					e.MemAlign = 1
				} else if v == 4 {
					e.MemAlign = 2
				} else if v == 8 {
					e.MemAlign = 3
				} else if v == 16 {
					e.MemAlign = 4
				} else {
					return fmt.Errorf("Invalid align %d", v)
				}
				return nil
			} else if strings.HasPrefix(t, "offset=") {
				v, err := strconv.Atoi(t[7:])
				if err != nil {
					return err
				}
				e.MemOffset = v
			} else {
				return errors.New("Error parsing memory operands")
			}
		}
		return nil
	} else if opcode == "memory.size" ||
		opcode == "memory.grow" {
		e.Opcode = InstrToOpcode[opcode]
		return nil
	} else if opcode == "block" ||
		opcode == "if" ||
		opcode == "loop" ||
		opcode == "else" ||
		opcode == "end" {
		e.Opcode = InstrToOpcode[opcode]
		e.Result = types.ValNone
		// Optional result type...
		s = strings.Trim(s, encoding.Whitespace)
		if len(s) == 0 {
			return nil
		}
		if s[0] == '(' {
			// eg (result i32)
			var ok bool
			var rtype string
			rtype, s = encoding.ReadElement(s)
			if strings.HasPrefix(rtype, "(result") {
				rtype = strings.Trim(rtype[7:len(rtype)-1], encoding.Whitespace)
				e.Result, ok = types.ValTypeToByte[rtype]
				if !ok {
					return errors.New("Error parsing block result")
				}
			}
		}
		return nil
	} else if opcode == "i32.const" {
		e.Opcode = InstrToOpcode[opcode]
		s = strings.Trim(s, encoding.Whitespace)
		v, _ := encoding.ReadToken(s)
		if strings.HasPrefix(v, "offset(") {
			dname := v[7 : len(v)-1]
			e.DataOffsetNeedsLinking = true
			e.DataOffsetNeedsAdjusting = true
			e.I32DataId = dname
			return nil
		} else if strings.HasPrefix(v, "reloffset(") {
			dname := v[7 : len(v)-1]
			e.DataOffsetNeedsLinking = true
			e.DataOffsetNeedsAdjusting = false
			e.I32DataId = dname
			return nil
		} else if strings.HasPrefix(v, "length(") {
			// Lookup the data length...
			dname := v[7 : len(v)-1]
			e.DataLengthNeedsLinking = true
			e.I32DataId = dname
			return nil
		}
		// Support other bases...
		if strings.HasPrefix(v, "0x") {
			vv, err := strconv.ParseUint(v[2:], 16, 32)
			if err != nil {
				return err
			}
			e.I32Value = int32(vv)
			return nil
		}

		vv, err := strconv.Atoi(v)
		if err != nil {
			return err
		}
		e.I32Value = int32(vv)
		return nil
	} else if opcode == "i64.const" {
		e.Opcode = InstrToOpcode[opcode]
		s = strings.Trim(s, encoding.Whitespace)
		v, _ := encoding.ReadToken(s)
		// Support other bases...
		if strings.HasPrefix(v, "0x") {
			vv, err := strconv.ParseUint(v[2:], 16, 64)
			if err != nil {
				return err
			}

			e.I64Value = int64(vv)

			return nil
		}
		vv, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return err
		}
		e.I64Value = int64(vv)
		return nil
	} else if opcode == "f32.const" {
		s = strings.Trim(s, encoding.Whitespace)
		v, _ := encoding.ReadToken(s)
		vv, err := strconv.ParseFloat(v, 32)
		if err != nil {
			return err
		}
		e.Opcode = InstrToOpcode[opcode]
		e.F32Value = float32(vv)
		return nil
	} else if opcode == "f64.const" {
		s = strings.Trim(s, encoding.Whitespace)
		v, _ := encoding.ReadToken(s)
		vv, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return err
		}
		e.Opcode = InstrToOpcode[opcode]
		e.F64Value = float64(vv)
		return nil
	} else if opcode == "local.get" ||
		opcode == "local.set" ||
		opcode == "local.tee" {
		e.Opcode = InstrToOpcode[opcode]
		var target string
		var lid int
		var err error
		target, s = encoding.ReadToken(s)
		if localNames != nil && strings.HasPrefix(target, "$") {
			// Find the id for it...
			lid, ok := localNames[target]
			if !ok {
				return fmt.Errorf("Local name %s not found", target)
			}
			e.LocalIndex = lid
			return nil
		}
		lid, err = strconv.Atoi(target)
		if err != nil {
			return err
		}
		e.LocalIndex = lid
		return nil
	} else if opcode == "global.get" ||
		opcode == "global.set" {
		e.Opcode = InstrToOpcode[opcode]
		var target string
		var gid int
		var err error
		target, s = encoding.ReadToken(s)
		if target[0] == '$' {
			e.GlobalNeedsLinking = true
			e.GlobalId = target
			return nil
		} else {
			gid, err = strconv.Atoi(target)
			if err != nil {
				return err
			}
			e.GlobalIndex = gid
			return nil
		}
	} else if opcode == "call" {
		e.Opcode = InstrToOpcode[opcode]
		var target string
		var fid int
		var err error
		target, s = encoding.ReadToken(s)
		if target[0] == '$' {
			e.FunctionNeedsLinking = true
			e.FunctionId = target
			return nil
		} else {
			fid, err = strconv.Atoi(target)
			if err != nil {
				return err
			}
			e.FuncIndex = fid
			return nil
		}
	} else if opcode == "call_indirect" {
		e.Opcode = InstrToOpcode[opcode]
		s = strings.Trim(s, encoding.Whitespace)
		if s[0] == '(' {
			typeInfo, _ := encoding.ReadElement(s)
			if strings.HasPrefix(typeInfo, "(type") {
				typeInfo = strings.Trim(typeInfo[5:len(typeInfo)-1], encoding.Whitespace)
				var err error
				e.TypeIndex, err = strconv.Atoi(typeInfo)
				if err != nil {
					return err
				}
			} else {
				return errors.New("Error parsing call_indirect")
			}
		} else {
			return errors.New("Error parsing call_indirect")
		}
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
