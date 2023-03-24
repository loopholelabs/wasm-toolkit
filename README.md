# Wasm toolkit

Pre-requisites: Make sure you have wat2wasm and wasm2wat from webt https://github.com/WebAssembly/wabt


## debugger

Can be used to modify a wasm file to add debug trace logging. Currently it requires a host with the correct calls, but later it will be done within the wasm and sent to wasi/stderr

## linker

Can be used to link 2 wasm files into a single wasm file. This is a work in progress. Most things work, but there's some more work to be done.