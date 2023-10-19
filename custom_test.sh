#!/bin/bash

cd cmd
rm cmd
go build .
cd ..

rm test_customs/out.wasm
rm debug.wat

# Run
./cmd/cmd customs -i test_customs/module.wasm -o test_customs/out.wasm \
--muximport "env/hello,0:env/zero,1:env/one,2:env/two" \
--muxexport "resize,0:resize_zero,1:resize_one,2:resize_two"

# Show it
wasm2wat test_customs/out.wasm
