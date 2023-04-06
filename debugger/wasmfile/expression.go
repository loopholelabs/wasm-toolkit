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
)

type Opcode byte

// TODO
// "table.get"			- 0x25
// "table.set"			- 0x26
// "select <t*>"		- 0x1c
// "ref.null t"			- 0xd0
// "ref.is_null"		- 0xd1
// "ref.func x"			- 0xd2

// All vector instructions 0xfd -

func DecodeSleb128(b []byte) (s int64, n int) {
	result := int64(0)
	shift := 0
	ptr := 0
	for {
		by := b[ptr]
		ptr++
		result = result | (int64(by&0x7f) << shift)
		shift += 7
		if (by & 0x80) == 0 {
			if shift < 64 && (by&0x40) != 0 {
				return result | (^0 << shift), ptr
			}
			return result, ptr
		}
	}
}

func AppendSleb128(buf []byte, val int64) []byte {
	for {
		b := val & 0x7f
		val = val >> 7
		if (val == 0 && b&0x40 == 0) ||
			(val == -1 && b&0x40 != 0) {
			buf = append(buf, byte(b))
			return buf
		}
		buf = append(buf, byte(b|0x80))
	}
}

const ExtendedOpcodeFC = Opcode(0xfc)

var instrToOpcodeFC = map[string]int{
	"i32.trunc_sat_f32_s": 0,
	"i32.trunc_sat_f32_u": 1,
	"i32.trunc_sat_f64_s": 2,
	"i32.trunc_sat_f64_u": 3,
	"i64.trunc_sat_f32_s": 4,
	"i64.trunc_sat_f32_u": 5,
	"i64.trunc_sat_f64_s": 6,
	"i64.trunc_sat_f64_u": 7,
	"memory.init":         8,
	"data.drop":           9,
	"memory.copy":         10,
	"memory.fill":         11,
	"table.init":          12,
	"elem.drop":           13,
	"table.copy":          14,
	"table.grow":          15,
	"table.size":          16,
	"table.fill":          17,
}

var instrToOpcode = map[string]Opcode{
	// Control
	"unreachable":   Opcode(0x00),
	"nop":           Opcode(0x01),
	"block":         Opcode(0x02),
	"loop":          Opcode(0x03),
	"if":            Opcode(0x04),
	"else":          Opcode(0x05),
	"end":           Opcode(0x0b),
	"br":            Opcode(0x0c),
	"br_if":         Opcode(0x0d),
	"br_table":      Opcode(0x0e),
	"return":        Opcode(0x0f),
	"call":          Opcode(0x10),
	"call_indirect": Opcode(0x11),

	// Parametric
	"drop":   Opcode(0x1a),
	"select": Opcode(0x1b),

	// Variable
	"local.get":  Opcode(0x20),
	"local.set":  Opcode(0x21),
	"local.tee":  Opcode(0x22),
	"global.get": Opcode(0x23),
	"global.set": Opcode(0x24),

	// Memory
	"i32.load":     Opcode(0x28),
	"i64.load":     Opcode(0x29),
	"f32.load":     Opcode(0x2a),
	"f64.load":     Opcode(0x2b),
	"i32.load8_s":  Opcode(0x2c),
	"i32.load8_u":  Opcode(0x2d),
	"i32.load16_s": Opcode(0x2e),
	"i32.load16_u": Opcode(0x2f),
	"i64.load8_s":  Opcode(0x30),
	"i64.load8_u":  Opcode(0x31),
	"i64.load16_s": Opcode(0x32),
	"i64.load16_u": Opcode(0x33),
	"i64.load32_s": Opcode(0x34),
	"i64.load32_u": Opcode(0x35),
	"i32.store":    Opcode(0x36),
	"i64.store":    Opcode(0x37),
	"f32.store":    Opcode(0x38),
	"f64.store":    Opcode(0x39),
	"i32.store8":   Opcode(0x3a),
	"i32.store16":  Opcode(0x3b),
	"i64.store8":   Opcode(0x3c),
	"i64.store16":  Opcode(0x3d),
	"i64.store32":  Opcode(0x3e),
	"memory.size":  Opcode(0x3f),
	"memory.grow":  Opcode(0x40),

	// Numeric
	"i32.const":           Opcode(0x41),
	"i64.const":           Opcode(0x42),
	"f32.const":           Opcode(0x43),
	"f64.const":           Opcode(0x44),
	"i32.eqz":             Opcode(0x45),
	"i32.eq":              Opcode(0x46),
	"i32.ne":              Opcode(0x47),
	"i32.lt_s":            Opcode(0x48),
	"i32.lt_u":            Opcode(0x49),
	"i32.gt_s":            Opcode(0x4a),
	"i32.gt_u":            Opcode(0x4b),
	"i32.le_s":            Opcode(0x4c),
	"i32.le_u":            Opcode(0x4d),
	"i32.ge_s":            Opcode(0x4e),
	"i32.ge_u":            Opcode(0x4f),
	"i64.eqz":             Opcode(0x50),
	"i64.eq":              Opcode(0x51),
	"i64.ne":              Opcode(0x52),
	"i64.lt_s":            Opcode(0x53),
	"i64.lt_u":            Opcode(0x54),
	"i64.gt_s":            Opcode(0x55),
	"i64.gt_u":            Opcode(0x56),
	"i64.le_s":            Opcode(0x57),
	"i64.le_u":            Opcode(0x58),
	"i64.ge_s":            Opcode(0x59),
	"i64.ge_u":            Opcode(0x5a),
	"f32.eq":              Opcode(0x5b),
	"f32.ne":              Opcode(0x5c),
	"f32.lt":              Opcode(0x5d),
	"f32.gt":              Opcode(0x5e),
	"f32.le":              Opcode(0x5f),
	"f32.ge":              Opcode(0x60),
	"f64.eq":              Opcode(0x61),
	"f64.ne":              Opcode(0x62),
	"f64.lt":              Opcode(0x63),
	"f64.gt":              Opcode(0x64),
	"f64.le":              Opcode(0x65),
	"f64.ge":              Opcode(0x66),
	"i32.clz":             Opcode(0x67),
	"i32.ctz":             Opcode(0x68),
	"i32.popcnt":          Opcode(0x69),
	"i32.add":             Opcode(0x6a),
	"i32.sub":             Opcode(0x6b),
	"i32.mul":             Opcode(0x6c),
	"i32.div_s":           Opcode(0x6d),
	"i32.div_u":           Opcode(0x6e),
	"i32.rem_s":           Opcode(0x6f),
	"i32.rem_u":           Opcode(0x70),
	"i32.and":             Opcode(0x71),
	"i32.or":              Opcode(0x72),
	"i32.xor":             Opcode(0x73),
	"i32.shl":             Opcode(0x74),
	"i32.shr_s":           Opcode(0x75),
	"i32.shr_u":           Opcode(0x76),
	"i32.rotl":            Opcode(0x77),
	"i32.rotr":            Opcode(0x78),
	"i64.clz":             Opcode(0x79),
	"i64.ctz":             Opcode(0x7a),
	"i64.popcnt":          Opcode(0x7b),
	"i64.add":             Opcode(0x7c),
	"i64.sub":             Opcode(0x7d),
	"i64.mul":             Opcode(0x7e),
	"i64.div_s":           Opcode(0x7f),
	"i64.div_u":           Opcode(0x80),
	"i64.rem_s":           Opcode(0x81),
	"i64.rem_u":           Opcode(0x82),
	"i64.and":             Opcode(0x83),
	"i64.or":              Opcode(0x84),
	"i64.xor":             Opcode(0x85),
	"i64.shl":             Opcode(0x86),
	"i64.shr_s":           Opcode(0x87),
	"i64.shr_u":           Opcode(0x88),
	"i64.rotl":            Opcode(0x89),
	"i64.rotr":            Opcode(0x8a),
	"f32.abs":             Opcode(0x8b),
	"f32.neg":             Opcode(0x8c),
	"f32.ceil":            Opcode(0x8d),
	"f32.floor":           Opcode(0x8e),
	"f32.trunc":           Opcode(0x8f),
	"f32.nearest":         Opcode(0x90),
	"f32.sqrt":            Opcode(0x91),
	"f32.add":             Opcode(0x92),
	"f32.sub":             Opcode(0x93),
	"f32.mul":             Opcode(0x94),
	"f32.div":             Opcode(0x95),
	"f32.min":             Opcode(0x96),
	"f32.max":             Opcode(0x97),
	"f32.copysign":        Opcode(0x98),
	"f64.abs":             Opcode(0x99),
	"f64.neg":             Opcode(0x9a),
	"f64.ceil":            Opcode(0x9b),
	"f64.floor":           Opcode(0x9c),
	"f64.trunc":           Opcode(0x9d),
	"f64.nearest":         Opcode(0x9e),
	"f64.sqrt":            Opcode(0x9f),
	"f64.add":             Opcode(0xa0),
	"f64.sub":             Opcode(0xa1),
	"f64.mul":             Opcode(0xa2),
	"f64.div":             Opcode(0xa3),
	"f64.min":             Opcode(0xa4),
	"f64.max":             Opcode(0xa5),
	"f64.copysign":        Opcode(0xa6),
	"i32.wrap_i64":        Opcode(0xa7),
	"i32.trunc_f32_s":     Opcode(0xa8),
	"i32.trunc_f32_u":     Opcode(0xa9),
	"i32.trunc_f64_s":     Opcode(0xaa),
	"i32.trunc_f64_u":     Opcode(0xab),
	"i64.extend_i32_s":    Opcode(0xac),
	"i64.extend_i32_u":    Opcode(0xad),
	"i64.trunc_f32_s":     Opcode(0xae),
	"i64.trunc_f32_u":     Opcode(0xaf),
	"i64.trunc_f64_s":     Opcode(0xb0),
	"i64.trunc_f64_u":     Opcode(0xb1),
	"f32.convert_i32_s":   Opcode(0xb2),
	"f32.convert_i32_u":   Opcode(0xb3),
	"f32.convert_i64_s":   Opcode(0xb4),
	"f32.convert_i64_u":   Opcode(0xb5),
	"f32.demote_f64":      Opcode(0xb6),
	"f64.convert_i32_s":   Opcode(0xb7),
	"f64.convert_i32_u":   Opcode(0xb8),
	"f64.convert_i64_s":   Opcode(0xb9),
	"f64.convert_i64_u":   Opcode(0xba),
	"f64.promote_f32":     Opcode(0xbb),
	"i32.reinterpret_f32": Opcode(0xbc),
	"i64.reinterpret_f64": Opcode(0xbd),
	"f32.reinterpret_i32": Opcode(0xbe),
	"f64.reinterpret_i64": Opcode(0xbf),
	"i32.extend8_s":       Opcode(0xc0),
	"i32.extend16_s":      Opcode(0xc1),
	"i64.extend8_s":       Opcode(0xc2),
	"i64.extend16_s":      Opcode(0xc3),
	"i64.extend32_s":      Opcode(0xc4),
}

var opcodeToInstr map[Opcode]string

var opcodeToInstrFC map[int]string

func init() {
	opcodeToInstr = make(map[Opcode]string)
	for s, o := range instrToOpcode {
		opcodeToInstr[o] = s
	}

	opcodeToInstrFC = make(map[int]string)
	for s, o := range instrToOpcodeFC {
		opcodeToInstrFC[o] = s
	}
}

type Expression struct {
	PC              uint64
	Opcode          Opcode
	OpcodeExt       int
	I32Value        int32
	I64Value        int64
	F32Value        float32
	F64Value        float64
	FuncIndex       int
	LocalIndex      int
	GlobalIndex     int
	LabelIndex      int
	TypeIndex       int
	TableIndex      int
	Labels          []int
	Result          ValType
	InnerExpression []*Expression
	MemAlign        int
	MemOffset       int
}

func (e *Expression) Show(prefix string, wf WasmFile) {
	opcode := opcodeToInstr[e.Opcode]
	if e.Opcode == ExtendedOpcodeFC {
		opcode = opcodeToInstrFC[e.OpcodeExt]
	}

	// See if we have any line info...
	lineInfo := ""
	li, ok := wf.lineNumbers[e.PC]
	if ok {
		// TODO: Read source file...
		lineInfo = fmt.Sprintf("%s:%d", li.Filename, li.Linenumber)
	}

	fmt.Printf("%d: %s %s ;; %s\n", e.PC, prefix, opcode, lineInfo)
	if e.InnerExpression != nil {
		for _, ie := range e.InnerExpression {
			ie.Show(prefix+" ", wf)
		}
	}
}
