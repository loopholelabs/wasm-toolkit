# Go wasmtoolkit

## Features

* Encoding and decoding from `.wasm` binary (99% done)
* Encoding and decoding from `.wat` text (99% done)

## Quickstart

* wasm2wat - `./wasm-toolkit wasm2wat -i something.wasm -o something.wat`
* strace - `./wasm-toolkit strace -i something.wasm -o something-with-strace-stderr.wasm`

## TODO

* Add function timings and summary output
* Add import func wrapping
* Add dwarf debug info
