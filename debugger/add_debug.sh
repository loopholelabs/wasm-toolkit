#!/bin/bash

wasm2wat $1 > tmp.wat
go run . tmp.wat > output.wat
wat2wasm output.wat
