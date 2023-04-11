(module
  (import "wasi_snapshot_preview1" "clock_time_get" (func $debug_clock_time_get (param i32 i64 i32) (result i32)))

  (func $debug_gettime (result i64)
    i32.const 0
    i64.const 1000
    i32.const offset($debug_clock_loc)
    global.get $debug_start_mem
    i32.add
    call $debug_clock_time_get
    drop
    i32.const offset($debug_clock_loc)
    global.get $debug_start_mem
    i32.add
    i64.load
  )

  (data $debug_clock_loc 8)

  ;; Only allow 800 function stack for now
  (data $debug_timestamps_stack 800)

  (global $debug_timestamps_stack_pointer (mut i32) (i32.const 0))
)