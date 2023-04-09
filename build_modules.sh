#!/bin/bash
tinygo build -opt=0 -o module1.wasm -scheduler=none -target=wasi ./test_module_go/main.go
