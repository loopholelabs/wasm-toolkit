#!/bin/bash

wasm2wat $1 > tmp.wat
go run . tmp.wat $1 > output.wat
wat2wasm output.wat
