wasm2wat http-handler-TestRuntimeHTTPSignatureGo.wasm > http-handler-TestRuntimeHTTPSignatureGo.wat
go run . http-handler-TestRuntimeHTTPSignatureGo.wat > http-handle-go-debug.wat
wat2wasm http-handle-go-debug.wat
