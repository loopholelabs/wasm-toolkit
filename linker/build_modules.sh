#!/bin/bash
tinygo build -o module1.wasm -scheduler=none -target=wasi --no-debug ./mod1/main.go
tinygo build -o module2.wasm -scheduler=none -target=wasi --no-debug ./mod2/main.go

# Now convert them to wat

wasm2wat module1.wasm > module1.wat
wasm2wat module2.wasm > module2.wat
