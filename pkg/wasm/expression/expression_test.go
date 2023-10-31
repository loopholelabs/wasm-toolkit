package expression

import (
	"bytes"

	"github.com/loopholelabs/wasm-toolkit/pkg/wasm/types"
	"github.com/stretchr/testify/assert"

	"testing"
)

func verifyEncodeDecode(t *testing.T, expr *Expression) *Expression {
	var buf bytes.Buffer
	err := expr.EncodeBinary(&buf)
	assert.NoError(t, err)

	// Decode from binary
	expr2, n, err := NewExpression(buf.Bytes(), 0)

	assert.NoError(t, err)
	assert.Equal(t, n, buf.Len())
	assert.Equal(t, len(expr2), 1)

	equal := expr.Equals(expr2[0])

	assert.True(t, equal)

	return expr2[0]
}

func TestSimpleExpressions(t *testing.T) {
	for i := 0; i < 256; i++ {
		expr := &Expression{
			Opcode: Opcode(i),
		}
		if expr.HasNoArgs() && i != int(InstrToOpcode["end"]) {
			verifyEncodeDecode(t, expr)
		}
	}
}

func TestBrTable(t *testing.T) {
	expr := &Expression{
		Opcode:     InstrToOpcode["br_table"],
		Labels:     []int{1, 2, 3, 4},
		LabelIndex: 5,
	}

	expr2 := verifyEncodeDecode(t, expr)
	assert.Equal(t, expr2.Opcode, expr.Opcode)
}

func TestBr(t *testing.T) {
	expr := &Expression{
		Opcode:     InstrToOpcode["br"],
		LabelIndex: 7,
	}

	expr2 := verifyEncodeDecode(t, expr)
	assert.Equal(t, expr2.Opcode, expr.Opcode)
}

func TestBrIf(t *testing.T) {
	expr := &Expression{
		Opcode:     InstrToOpcode["br_if"],
		LabelIndex: 9,
	}

	expr2 := verifyEncodeDecode(t, expr)
	assert.Equal(t, expr2.Opcode, expr.Opcode)
}

func TestMemoryExpressions(t *testing.T) {
	for i := 0; i < 256; i++ {
		expr := &Expression{
			Opcode:    Opcode(i),
			MemAlign:  3,
			MemOffset: 9,
		}
		if expr.HasMemoryArgs() {
			expr2 := verifyEncodeDecode(t, expr)
			assert.Equal(t, expr.Opcode, expr2.Opcode)
			assert.Equal(t, expr.MemAlign, expr2.MemAlign)
			assert.Equal(t, expr.MemOffset, expr2.MemOffset)
		}
	}
}

func TestMemorySize(t *testing.T) {
	expr := &Expression{
		Opcode: InstrToOpcode["memory.size"],
	}

	expr2 := verifyEncodeDecode(t, expr)
	assert.Equal(t, expr2.Opcode, expr.Opcode)
}

func TestMemoryGrow(t *testing.T) {
	expr := &Expression{
		Opcode: InstrToOpcode["memory.grow"],
	}

	expr2 := verifyEncodeDecode(t, expr)
	assert.Equal(t, expr2.Opcode, expr.Opcode)
}

func TestBlockIfLoop(t *testing.T) {
	for _, c := range []string{"block", "if", "loop"} {
		expr := &Expression{
			Opcode: InstrToOpcode[c],
			Result: types.ValI32,
		}

		expr2 := verifyEncodeDecode(t, expr)
		assert.Equal(t, expr2.Opcode, expr.Opcode)
		assert.Equal(t, expr.Result, expr2.Result)
	}
}

func TestI32Const(t *testing.T) {
	for _, v := range []int32{1, -90, 123456, 90000000} {
		expr := &Expression{
			Opcode:   InstrToOpcode["i32.const"],
			I32Value: v,
		}

		expr2 := verifyEncodeDecode(t, expr)
		assert.Equal(t, expr2.Opcode, expr.Opcode)
		assert.Equal(t, expr.I32Value, expr2.I32Value)
	}
}

func TestI64Const(t *testing.T) {
	for _, v := range []int64{1, -90, 123456, 90000000, -0xffebd6e0} {
		expr := &Expression{
			Opcode:   InstrToOpcode["i64.const"],
			I64Value: v,
		}

		var buf bytes.Buffer

		err := expr.EncodeBinary(&buf)
		if err != nil {
			panic(err)
		}

		expr2 := verifyEncodeDecode(t, expr)
		assert.Equal(t, expr2.Opcode, expr.Opcode)
		assert.Equal(t, expr.I64Value, expr2.I64Value)
	}
}

func TestF32Const(t *testing.T) {
	expr := &Expression{
		Opcode:   InstrToOpcode["f32.const"],
		F32Value: 123.456,
	}

	expr2 := verifyEncodeDecode(t, expr)
	assert.Equal(t, expr2.Opcode, expr.Opcode)
	assert.Equal(t, expr.F32Value, expr2.F32Value)
}

func TestF64Const(t *testing.T) {
	expr := &Expression{
		Opcode:   InstrToOpcode["f64.const"],
		F64Value: 123.456,
	}

	expr2 := verifyEncodeDecode(t, expr)
	assert.Equal(t, expr2.Opcode, expr.Opcode)
	assert.Equal(t, expr.F64Value, expr2.F64Value)
}

func TestLocals(t *testing.T) {
	for _, c := range []string{"local.get", "local.set", "local.tee"} {
		expr := &Expression{
			Opcode:     InstrToOpcode[c],
			LocalIndex: 7,
		}

		expr2 := verifyEncodeDecode(t, expr)
		assert.Equal(t, expr2.Opcode, expr.Opcode)
		assert.Equal(t, expr.LocalIndex, expr2.LocalIndex)
	}
}

func TestGlobals(t *testing.T) {
	for _, c := range []string{"global.get", "global.set"} {
		expr := &Expression{
			Opcode:      InstrToOpcode[c],
			GlobalIndex: 7,
		}

		expr2 := verifyEncodeDecode(t, expr)
		assert.Equal(t, expr2.Opcode, expr.Opcode)
		assert.Equal(t, expr.GlobalIndex, expr2.GlobalIndex)
	}
}

func TestCall(t *testing.T) {
	expr := &Expression{
		Opcode:    InstrToOpcode["call"],
		FuncIndex: 123,
	}

	expr2 := verifyEncodeDecode(t, expr)
	assert.Equal(t, expr2.Opcode, expr.Opcode)
	assert.Equal(t, expr.FuncIndex, expr2.FuncIndex)
}

func TestCallIndirect(t *testing.T) {
	expr := &Expression{
		Opcode:     InstrToOpcode["call_indirect"],
		TypeIndex:  123,
		TableIndex: 7,
	}

	expr2 := verifyEncodeDecode(t, expr)
	assert.Equal(t, expr2.Opcode, expr.Opcode)
	assert.Equal(t, expr.TypeIndex, expr2.TypeIndex)
	assert.Equal(t, expr.TableIndex, expr2.TableIndex)
}

func TestMemoryCopy(t *testing.T) {
	expr := &Expression{
		Opcode:    ExtendedOpcodeFC,
		OpcodeExt: instrToOpcodeFC["memory.copy"],
	}

	expr2 := verifyEncodeDecode(t, expr)
	assert.Equal(t, expr2.Opcode, expr.Opcode)
	assert.Equal(t, expr2.OpcodeExt, expr.OpcodeExt)
}

func TestMemoryFill(t *testing.T) {
	expr := &Expression{
		Opcode:    ExtendedOpcodeFC,
		OpcodeExt: instrToOpcodeFC["memory.fill"],
	}

	expr2 := verifyEncodeDecode(t, expr)
	assert.Equal(t, expr2.Opcode, expr.Opcode)
	assert.Equal(t, expr2.OpcodeExt, expr.OpcodeExt)
}

func TestTruncSat(t *testing.T) {
	for _, c := range []string{
		"i32.trunc_sat_f32_s",
		"i32.trunc_sat_f32_u",
		"i32.trunc_sat_f64_s",
		"i32.trunc_sat_f64_u",
		"i64.trunc_sat_f32_s",
		"i64.trunc_sat_f32_u",
		"i64.trunc_sat_f64_s",
		"i64.trunc_sat_f64_u",
	} {
		expr := &Expression{
			Opcode:    ExtendedOpcodeFC,
			OpcodeExt: instrToOpcodeFC[c],
		}

		expr2 := verifyEncodeDecode(t, expr)
		assert.Equal(t, expr2.Opcode, expr.Opcode)
		assert.Equal(t, expr2.OpcodeExt, expr.OpcodeExt)
	}
}
