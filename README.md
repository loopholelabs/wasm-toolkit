# Go wasmtoolkit

## Features

* Encoding and decoding from `.wasm` binary (99% done)
* Encoding and decoding from `.wat` text (99% done)

## Quickstart

* wasm2wat - `./wasm-toolkit wasm2wat -i something.wasm -o something.wat`

## TODO

* Easy ways to hook/wrap instructions, calls
* Add function entry/exit code
* Data relocation
* wat extensions such as `offset` and `length` of data
