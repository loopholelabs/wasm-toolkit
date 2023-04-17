#!/bin/bash
rm main_strace.wasm

go run . strace -i ../main.wasm -o main_strace.wasm --all --color --func '^\$IMPORT.*|\_start|^\$main\..*' --timing true

#go run . strace -i ../module1.wasm -o module1_strace.wasm --all --color --func '^\$IMPORT.*|\_start' --timing true --watch main.some_global,main.another_global

wasmtime --dir . --env TEST=1 main_strace.wasm 1 2 3
