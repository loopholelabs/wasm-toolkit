package wasmfile

import (
	"encoding/binary"
	"fmt"
	"math"
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
	l := len(b)
	if l > 10 {
		l = 10
	}

	for i := 0; i < l; i++ {
		s |= int64(b[i]&0x7f) << (7 * i)
		if b[i]&0x80 == 0 {
			// If it's signed
			if b[i]&0x40 != 0 {
				s |= ^0 << (7 * (i + 1))
			}
			n = i + 1
			return
		}
	}
	panic("Error decoding leb128")
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
