# Go wasmtoolkit

* TODO: Not everything is implemented yet. Wasm compatibility is 99%. You may be the 1%.

## Features

* Encoding and decoding from `.wasm` binary (99% done)
* Encoding and decoding from `.wat` text (99% done)

* Adding strace like output to a wasm file, which gets output to `STDERR` as it runs.
  * All function calls, with arguments and return
  * All calls to imported functions, with args and return
  * Dwarf source file and line number ranges
  * Dwarf paramater names
  * Wasi preview1 call and return values
  * Function call count and timings summary
  * Watch globals by name (i32 only so far)

* POC Embedding a file into a wasm which is then available to the module.

## Quickstart

* wasm2wat - `./wasm-toolkit wasm2wat -i something.wasm -o something.wat`
* wat2wasm - `./wasm-toolkit wat2wasm -i something.wat -o something.wasm`
* strace - `./wasm-toolkit strace -i something.wasm -o something-with-strace-stderr.wasm`
* embedfile - `./wasm-toolkit embedfile -i something.wasm -o something_embed.wasm --filename embedtest --content "This is some file data :)"`

## Strace

![alt text](https://raw.githubusercontent.com/loopholelabs/wasm-toolkit/master/strace.png)

This will wrap imports, log all data only for imports

`./wasm-toolkit strace -i ../module1.wasm -o module1_debug.wasm --all --func ".*IMPORT.*"`

This will add logging for all functions starting with `$runtime`

`./wasm-toolkit strace -i ../module1.wasm -o module1_debug.wasm --all --func '^\$runtime.*'`

You can also compile wasm-toolkit to wasm and add tracing to it :)

## Embed file

![alt text](https://raw.githubusercontent.com/loopholelabs/wasm-toolkit/master/embed.png)

## Example output

On the left is an strace like output. On the right is a wat output with debugging info.

![alt text](https://raw.githubusercontent.com/loopholelabs/wasm-toolkit/master/output.png)

## TODO

* Dwarf watched globals - automatically determine type/length from dwarf data, support many types.
* Hook unreachable calls and add far more context to aid debugging. eg show all watched variables.
