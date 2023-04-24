#!/bin/bash

echo "Building tinygo example module..."
tinygo build -opt=0 -o module_tinygo.wasm -scheduler=none -target=wasi -x ./test_module_go/main.go

echo "Building using gotip..."
GOARCH=wasm GOOS=wasip1 gotip build -o module_gotip.wasm test_module_gotip/main.go

echo "Building using rust..."
cd test_module_rs
rm -rf target
cargo build --release --target=wasm32-wasi
mv target/wasm32-wasi/release/hello_world.wasm ../module_rust.wasm
cd ..
