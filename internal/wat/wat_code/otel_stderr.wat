(module
  (type (func (param i32 i32) (result i32)))
  (import "wasi_snapshot_preview1" "random_get" (func $debug_random_get (type 1)))

  (func $get_invocation_id (param $ptr i32)
    global.get $trace_id_set
    i32.eqz
    if
      i32.const offset($trace_id_stderr)
      i32.const length($trace_id_stderr)
      call $debug_random_get
      drop
      i32.const 1
      global.set $trace_id_set
    end

    ;; dest / src / len
    local.get $ptr
    i32.const offset($trace_id_stderr)
    i32.const 16
    memory.copy
  )

  (func $send_otel_trace_json (param $ptr i32) (param $len i32)
    local.get $ptr
    local.get $len
    call $wt_print  
  )

  (func $cache_service_name
    ;; Nothing to do...
  )

  (global $trace_id_set (mut i32) (i32.const 0))
  (data $trace_id_stderr 16)

  (data $service_name "unknown")
  (global $service_name_len (mut i32) (i32.const 7))
)