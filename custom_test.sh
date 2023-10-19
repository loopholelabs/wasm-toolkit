#!/bin/bash

cd cmd
rm cmd
go build .
cd ..

rm test_customs/out.wasm
rm debug.wat

# Run
./cmd/cmd customs -i test_customs/module.wasm -o test_customs/out.wasm

# Show it
wasm2wat test_customs/out.wasm
