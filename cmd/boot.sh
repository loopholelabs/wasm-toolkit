#!/bin/bash

echo "Compiling tracer to wasm..."
tinygo build -o tracer_tiny.wasm -target=wasi .

echo "Adding tracing to the tracer..."
wasmtime --dir . tracer_tiny.wasm -- strace -i tracer_tiny.wasm -o tracer_tiny_traced.wasm --func 'loopholelabs' --all --color --timing --dwarf
