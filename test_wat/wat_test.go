package testwat

import (
	"bytes"
	"context"
	"fmt"
	"path"
	"testing"

	"github.com/loopholelabs/wasm-toolkit/internal/wat"
	wasmfile "github.com/loopholelabs/wasm-toolkit/pkg/wasm"
	"github.com/loopholelabs/wasm-toolkit/pkg/wasm/types"
	"github.com/stretchr/testify/assert"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
)

func TestWatStdout(t *testing.T) {
	stdout := wasmfile.NewEmpty()
	data, err := wat.Wat_content.ReadFile(path.Join("wat_code", "stdout.wat"))
	assert.NoError(t, err)
	err = stdout.DecodeWat(data)
	assert.NoError(t, err)

	stdout_test := wasmfile.NewEmpty()
	data_test, err := wat.Wat_content.ReadFile(path.Join("wat_code", "stdout_test.wat"))
	assert.NoError(t, err)
	err = stdout_test.DecodeWat(data_test)
	assert.NoError(t, err)

	data_ptr := int32(0)
	data_base := int(data_ptr)

	wfile := wasmfile.NewEmpty()
	wfile.AddFuncsFrom(stdout_test, func(m map[int]int) {})
	data_ptr = wfile.AddDataFrom(data_ptr, stdout_test)
	wfile.AddExports(stdout_test)

	wfile.AddFuncsFrom(stdout, func(m map[int]int) {})
	data_ptr = wfile.AddDataFrom(data_ptr, stdout)

	// Resolve / link everything...
	for _, c := range wfile.Code {
		err = c.ResolveLengths(wfile)
		assert.NoError(t, err)

		err = c.ResolveRelocations(wfile, data_base)
		assert.NoError(t, err)

		err = c.ResolveGlobals(wfile)
		assert.NoError(t, err)

		err = c.ResolveFunctions(wfile)
		assert.NoError(t, err)
	}

	// Sort out memory...

	wfile.Memory = append(wfile.Memory, &wasmfile.MemoryEntry{LimitMin: 2, LimitMax: 2})
	wfile.Export = append(wfile.Export, &wasmfile.ExportEntry{
		Type:  types.ExportMem,
		Name:  "memory",
		Index: 0,
	})

	var buf bytes.Buffer
	err = wfile.EncodeBinary(&buf)

	wasmData := buf.Bytes()

	fmt.Printf("Got wasm file %d\n", len(wasmData))

	for _, e := range wfile.Export {
		fmt.Printf("Export %s %v %v\n", e.Name, e.Type, e.Index)
	}

	ctx := context.TODO()

	// TODO: Load up some wat files and test them.

	r := wazero.NewRuntime(ctx)
	defer r.Close(ctx)

	output := ""

	var mod api.Module

	r.NewHostModuleBuilder("wasi_snapshot_preview1").NewFunctionBuilder().
		WithFunc(func(ctx context.Context, mod api.Module, fd uint32, iov uint32, n uint32, byteswritten uint32) uint32 {
			mem := mod.Memory()

			// Get the data from memory...
			for i := uint32(0); i < n; i++ {
				ptr, _ := mem.ReadUint32Le(iov)
				len, _ := mem.ReadUint32Le(iov + 4)
				data, _ := mem.Read(ptr, len)
				output = output + string(data)
				iov += 8
			}
			return 0
		}).Export("fd_write").Instantiate(ctx)

	mod, err = r.Instantiate(ctx, wasmData)
	assert.NoError(t, err)

	f := mod.ExportedFunction("$test_stdout")

	_, err = f.Call(ctx)
	assert.NoError(t, err)

	assert.Equal(t, "Hello world", output)
}
