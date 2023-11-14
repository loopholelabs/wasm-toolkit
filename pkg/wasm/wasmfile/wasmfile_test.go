package wasmfile

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmptyEncodeDecodeBinary(t *testing.T) {
	wf := NewEmpty()

	var enc bytes.Buffer

	err := wf.EncodeBinary(&enc)
	assert.NoError(t, err)

	data := enc.Bytes()

	assert.Equal(t, 11, len(data))
	// Smallest wasm binary - simple header, and a SectionDataCount
	assert.Equal(t, "0061736d010000000c0100", hex.EncodeToString(data))

	wf2 := NewEmpty()
	err = wf2.DecodeBinary(data)
	assert.NoError(t, err)

	// Make sure the decoded wasm file is empty.
	assert.Equal(t, 0, len(wf2.Function))
	assert.Equal(t, 0, len(wf2.Type))
	assert.Equal(t, 0, len(wf2.Custom))
	assert.Equal(t, 0, len(wf2.Export))
	assert.Equal(t, 0, len(wf2.Import))
	assert.Equal(t, 0, len(wf2.Table))
	assert.Equal(t, 0, len(wf2.Global))
	assert.Equal(t, 0, len(wf2.Memory))
	assert.Equal(t, 0, len(wf2.Code))
	assert.Equal(t, 0, len(wf2.Data))
	assert.Equal(t, 0, len(wf2.Elem))
}

func TestEmptyEncodeDecodeWat(t *testing.T) {
	wf := NewEmpty()

	var wat bytes.Buffer
	err := wf.EncodeWat(&wat)
	assert.NoError(t, err)
	watData := string(wat.Bytes())
	assert.Equal(t, "(module\n)\n", watData)

	wf2 := NewEmpty()
	err = wf2.DecodeWat(wat.Bytes())
	assert.NoError(t, err)

	// Make sure the decoded wasm file is empty.
	assert.Equal(t, 0, len(wf2.Function))
	assert.Equal(t, 0, len(wf2.Type))
	assert.Equal(t, 0, len(wf2.Custom))
	assert.Equal(t, 0, len(wf2.Export))
	assert.Equal(t, 0, len(wf2.Import))
	assert.Equal(t, 0, len(wf2.Table))
	assert.Equal(t, 0, len(wf2.Global))
	assert.Equal(t, 0, len(wf2.Memory))
	assert.Equal(t, 0, len(wf2.Code))
	assert.Equal(t, 0, len(wf2.Data))
	assert.Equal(t, 0, len(wf2.Elem))
}
