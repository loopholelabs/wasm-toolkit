#!/bin/bash

# Add tracing output to the wasm module
go run . otel -i ../module_tinygo.wasm -o module_tinygo_otel.wasm --func '.*'

# Run the wasm module and send any tracing data off to jaeger
(wasmtime --dir . module_tinygo_otel.wasm 3>&1 1>&2- 2>&3- ) | tee trace.log | awk '{system("curl -X POST -H '\''Content-Type: application/json'\'' -d '\''" $0 "'\'' http://localhost:4318/v1/traces")}'
