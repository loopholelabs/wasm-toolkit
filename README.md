# WASM Toolkit

A Toolkit for working with WebAssembly

[![License: Apache 2.0](https://img.shields.io/badge/License-Apache%202.0-brightgreen.svg)](https://www.apache.org/licenses/LICENSE-2.0)
[![Discord](https://dcbadge.vercel.app/api/server/JYmFhtdPeu?style=flat)](https://loopholelabs.io/discord)
![Go Version](https://img.shields.io/badge/go%20version-%3E=1.18-61CFDD.svg)
[![Go Reference](https://pkg.go.dev/badge/github.com/loopholelabs/wasm-toolkit.svg)](https://pkg.go.dev/github.com/loopholelabs/wasm-toolkit)

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

* wasm2wat but including dwarf debug information - line numbers, variable names, etc

* POC Embedding a file into a wasm which is then available to the module.

## Quickstart

* wasm2wat - `./wasm-toolkit wasm2wat -i something.wasm -o something.wat`
* wat2wasm - `./wasm-toolkit wat2wasm -i something.wat -o something.wasm`
* strace - `./wasm-toolkit strace -i something.wasm -o something-with-strace-stderr.wasm`
* embedfile - `./wasm-toolkit embedfile -i something.wasm -o something_embed.wasm --filename embedtest --content "This is some file data :)"`

## Strace

### Simple example

`./wasm-toolkit strace -i ../module1.wasm -o module1_strace.wasm --all --color --func '^\$IMPORT.*|\_start'`

![alt text](https://raw.githubusercontent.com/loopholelabs/wasm-toolkit/master/screenshots/strace1.png)

### Filter functions

`./wasm-toolkit strace -i ../module1.wasm -o module1_strace.wasm --all --color --func '^\$main'`

![alt text](https://raw.githubusercontent.com/loopholelabs/wasm-toolkit/master/screenshots/strace2.png)

### Profiling

`./wasm-toolkit strace -i ../module1.wasm -o module1_strace.wasm --all --color --func '.*' --timing true`

![alt text](https://raw.githubusercontent.com/loopholelabs/wasm-toolkit/master/screenshots/strace3.png)

### Watch global variables

`./wasm-toolkit strace -i ../module1.wasm -o module1_strace.wasm --all --color --func '^\$main' --watch main.some_global,main.another_global`

![alt text](https://raw.githubusercontent.com/loopholelabs/wasm-toolkit/master/screenshots/strace4.png)

![alt text](https://raw.githubusercontent.com/loopholelabs/wasm-toolkit/master/screenshots/strace5.png)


You can also compile wasm-toolkit to wasm and add tracing to it :)

## Embed file (POC)

![alt text](https://raw.githubusercontent.com/loopholelabs/wasm-toolkit/master/embed.png)

## Example output

On the left is an strace like output. On the right is a wat output with debugging info.

![alt text](https://raw.githubusercontent.com/loopholelabs/wasm-toolkit/master/output.png)

## TODO

* Dwarf watched globals - automatically determine type/length from dwarf data, support many types.
* Hook unreachable calls and add far more context to aid debugging. eg show all watched variables.
* Output to open telemetry standard.
* More dwarf data - check up on the rust generated data.
* Implement remainder of wasm instructions, fix bugs in wat parsing (comments etc).
* More work on linking / composability.
* Finish wasi file embed.

## Contributing

Bug reports and pull requests are welcome on GitHub at [https://github.com/loopholelabs/wasm-toolkit][gitrepo]. For more contribution information check out [the contribution guide](https://github.com/loopholelabs/wasm-toolkit/blob/master/CONTRIBUTING.md).

## License

The WASM Toolkit project is available as open source under the terms of the [Apache License, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0).

## Code of Conduct

Everyone interacting in the WASM Toolkit project's codebases, issue trackers, chat rooms and mailing lists is expected to follow the [CNCF Code of Conduct](https://github.com/cncf/foundation/blob/master/code-of-conduct.md).

## Project Managed By:

[![https://loopholelabs.io][loopholelabs]](https://loopholelabs.io)

[gitrepo]: https://github.com/loopholelabs/wasm-toolkit
[loopholelabs]: https://cdn.loopholelabs.io/loopholelabs/LoopholeLabsLogo.svg
[loophomepage]: https://loopholelabs.io
