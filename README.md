# Go wasmtoolkit

## Features

* Encoding and decoding from `.wasm` binary (99% done)
* Encoding and decoding from `.wat` text (99% done)

## Quickstart

* wasm2wat - `./wasm-toolkit wasm2wat -i something.wasm -o something.wat`
* wat2wasm - `./wasm-toolkit wat2wasm -i something.wat -o something.wasm`
* strace - `./wasm-toolkit strace -i something.wasm -o something-with-strace-stderr.wasm`
* embedfile - `./wasm-toolkit embedfile -i something.wasm -o something_embed.wasm --filename embedtest --content "This is some file data :)"`

## Strace

This will wrap imports, log all data only for imports

`./wasm-toolkit strace -i ../module1.wasm -o module1_debug.wasm --all --func ".*IMPORT.*"`

This will add logging for all functions starting with `$runtime`

`./wasm-toolkit strace -i ../module1.wasm -o module1_debug.wasm --all --func '^\$runtime.*'`

## Embed file

![alt text](https://raw.githubusercontent.com/loopholelabs/wasm-toolkit/jm/new/embed.png)

## Example output

On the left is an strace like output. On the right is a wat output with debugging info.

![alt text](https://raw.githubusercontent.com/loopholelabs/wasm-toolkit/jm/new/output.png)

## TODO

* Add function timings and summary output
* Add more dwarf debug info
