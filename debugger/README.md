# Wasm trace

1. Run `./add_debug.sh <wasmfile>` to create a debug version of the wasm file. This inserts host callbacks.

2. Run the new wasm file and observe STDERR

## Features

* Full trace of all function calls within the wasm file

* Shows all parameters and return values

* Shows dwarf line number / source file information

* Monitors execution time and shows summary at the end of expensive functions

## TODO

* Use dwarf debug info on parameters