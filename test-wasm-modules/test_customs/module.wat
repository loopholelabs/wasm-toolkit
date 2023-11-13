(module
  (type (func (param i64 i64 i64 i64) (result i64)))
  (import "env" "hello" (func $hello (type 0)))

  (func $_start)

  (func $run (result i64)

    i64.const 0   ;; FID
    i64.const 1
    i64.const 2
    i64.const 3
    call $hello

    i64.const 1   ;; FID
    i64.const 1
    i64.const 2
    i64.const 3
    call $hello

    i64.const 2   ;; FID
    i64.const 1
    i64.const 2
    i64.const 3
    call $hello

    return
  )


  (func $resize (param i64) (param i64) (result i64)
    i64.const 123
    return
  )

  (memory $mem 2)           ;; How much memory we need in 64K pages to start with...

  (export "memory" (memory $mem))

  (export "_start" (func $_start))
  (export "run" (func $run))

  (export "resize" (func $resize))
)