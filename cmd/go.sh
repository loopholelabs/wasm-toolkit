#!/bin/bash

filename=$(basename -- "$1")
filename="${filename%.*}"

# Add tracing output to the wasm module
go run . otel -i $1 -o ${filename}_otel.wasm --func '.*'

# Run the wasm module and send any tracing data off to jaeger
(wasmtime --dir . ${filename}_otel.wasm 3>&1 1>&2- 2>&3- ) > trace.log

echo "Importing traces into jaeger..."
cat trace.log | awk '{system("curl -X POST -H '\''Content-Type: application/json'\'' -d '\''" $0 "'\'' http://localhost:4318/v1/traces")}'
