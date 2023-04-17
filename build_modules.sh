#!/bin/bash
tinygo build -opt=0 -o module1.wasm -scheduler=none -target=wasi -x ./test_module_go/main.go

GOARCH=wasm GOOS=wasip1 gotip build -o main.wasm -work -x -gcflags="-N -dwarf=true -dwarflocationlists=false" -ldflags "-w=false" test_module_gotip/main.go

# GOOS=js GOARCH=wasm go build -a -work -x -gcflags="-N -dwarf=true -dwarflocationlists=false" -ldflags="-w=false" -ldflags="-extldflags=\"-g\"" -o mmm . >build2.log 2>&1
