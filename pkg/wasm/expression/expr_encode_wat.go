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
	"bufio"
	"fmt"
	"io"

	"github.com/loopholelabs/wasm-toolkit/pkg/wasm/types"
)

type WasmDebugContext interface {
	GetLineNumberInfo(pc uint64) string
	GetGlobalIdentifier(globalIdx int, defaultEmpty bool) string
	GetFunctionIdentifier(funcIdx int, defaultEmpty bool) string
	GetLocalVarName(pc uint64, localIdx int) string
}

func (e *Expression) EncodeWat(w io.Writer, prefix string, wd WasmDebugContext) error {
	comment := "" //fmt.Sprintf("    ;; PC=%d", e.PC) // TODO From line numbers, vars etc

	lineNumberData := wd.GetLineNumberInfo(e.PC)
	if lineNumberData != "" {
		comment = fmt.Sprintf(" ;; Src = %s", lineNumberData)
	}

	wr := bufio.NewWriter(w)

	defer func() {
		wr.Flush()
	}()

	// First deal with simple opcodes (No args)
	if e.HasNoArgs() {
		_, err := wr.WriteString(fmt.Sprintf("%s%s%s\n", prefix, opcodeToInstr[e.Opcode], comment))
		return err
	} else if e.Opcode == InstrToOpcode["br_table"] {
		targets := ""
		for _, l := range e.Labels {
			targets = fmt.Sprintf("%s %d", targets, l)
		}
		defaultTarget := fmt.Sprintf(" %d", e.LabelIndex)
		_, err := wr.WriteString(fmt.Sprintf("%s%s%s%s%s\n", prefix, opcodeToInstr[e.Opcode], targets, defaultTarget, comment))
		return err
	} else if e.Opcode == InstrToOpcode["br"] ||
		e.Opcode == InstrToOpcode["br_if"] {
		target := fmt.Sprintf(" %d", e.LabelIndex)
		_, err := wr.WriteString(fmt.Sprintf("%s%s%s%s\n", prefix, opcodeToInstr[e.Opcode], target, comment))
		return err
	} else if e.HasMemoryArgs() {
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
	} else if e.Opcode == InstrToOpcode["memory.size"] ||
		e.Opcode == InstrToOpcode["memory.grow"] {
		_, err := wr.WriteString(fmt.Sprintf("%s%s%s\n", prefix, opcodeToInstr[e.Opcode], comment))
		return err
	} else if e.Opcode == InstrToOpcode["block"] ||
		e.Opcode == InstrToOpcode["if"] ||
		e.Opcode == InstrToOpcode["loop"] {

		result := ""
		if e.Result != types.ValNone {
			result = fmt.Sprintf(" (result %s)", types.ByteToValType[e.Result])
		}

		_, err := wr.WriteString(fmt.Sprintf("%s%s%s%s\n", prefix, opcodeToInstr[e.Opcode], result, comment))

		return err
	} else if e.Opcode == InstrToOpcode["i32.const"] {
		value := fmt.Sprintf(" %d", e.I32Value)
		_, err := wr.WriteString(fmt.Sprintf("%s%s%s%s\n", prefix, opcodeToInstr[e.Opcode], value, comment))
		return err
	} else if e.Opcode == InstrToOpcode["i64.const"] {
		value := fmt.Sprintf(" %d", e.I64Value)
		_, err := wr.WriteString(fmt.Sprintf("%s%s%s%s\n", prefix, opcodeToInstr[e.Opcode], value, comment))
		return err
	} else if e.Opcode == InstrToOpcode["f32.const"] {
		value := fmt.Sprintf(" %f", e.F32Value)
		if value == " +Inf" || value == " -Inf" {
			value = " inf"
		} else if value == " NaN" {
			value = " nan"
		}

		_, err := wr.WriteString(fmt.Sprintf("%s%s%s%s\n", prefix, opcodeToInstr[e.Opcode], value, comment))
		return err
	} else if e.Opcode == InstrToOpcode["f64.const"] {
		value := fmt.Sprintf(" %f", e.F64Value)
		if value == " +Inf" || value == " -Inf" {
			value = " inf"
		} else if value == " NaN" {
			value = " nan"
		}

		_, err := wr.WriteString(fmt.Sprintf("%s%s%s%s\n", prefix, opcodeToInstr[e.Opcode], value, comment))
		return err
	} else if e.Opcode == InstrToOpcode["local.get"] ||
		e.Opcode == InstrToOpcode["local.set"] ||
		e.Opcode == InstrToOpcode["local.tee"] {
		tname := wd.GetLocalVarName(e.PC, e.LocalIndex)
		//
		if tname == "" {
			tname = wd.GetLocalVarName(e.PCNext, e.LocalIndex)
		}

		if tname != "" {
			comment = comment + " ;; Variable " + tname
		}
		localTarget := fmt.Sprintf(" %d", e.LocalIndex)
		_, err := wr.WriteString(fmt.Sprintf("%s%s%s%s\n", prefix, opcodeToInstr[e.Opcode], localTarget, comment))
		return err
	} else if e.Opcode == InstrToOpcode["global.get"] ||
		e.Opcode == InstrToOpcode["global.set"] {
		g := wd.GetGlobalIdentifier(e.GlobalIndex, false)
		globalTarget := fmt.Sprintf(" %s", g)
		_, err := wr.WriteString(fmt.Sprintf("%s%s%s%s\n", prefix, opcodeToInstr[e.Opcode], globalTarget, comment))
		return err
	} else if e.Opcode == InstrToOpcode["call"] {
		f := wd.GetFunctionIdentifier(e.FuncIndex, false)
		callTarget := fmt.Sprintf(" %s", f)
		_, err := wr.WriteString(fmt.Sprintf("%s%s%s%s\n", prefix, opcodeToInstr[e.Opcode], callTarget, comment))
		return err
	} else if e.Opcode == InstrToOpcode["call_indirect"] {
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
