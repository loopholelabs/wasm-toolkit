# Wasm trace

1. Run `./go.sh` to create a debug version of the wasm file. This inserts host callbacks.

2. Copy the list of function names from the output, into `scale-ts/src/runtime/module.ts`
   (This step will be removed soon)

3. Go into scale-ts and run `./go.sh ../http-handle-go-debug.wasm`

## Features

* Full trace of all function calls within the wasm file

* Shows all parameters and return values

* For any i32 params and returns, it also shows a snippet of the memory they may be pointing to

* Memory diffs - show what memory has changed at the end of each function

## TODO

* Show global var diffs at function end

* Remove the need for all the code in the host - most/all of this can be done within the wasm file and simply output to wasi/stderr
