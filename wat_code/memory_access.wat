(module

  (func $debug_i32.store (param $addr i32) (param $value i32) (param $offset i32) (param $align i32)
    local.get $addr
    i32.const 4
    call $debug_log_write_pre
    local.get $addr
    local.get $offset
    i32.add
    local.get $value
    i32.store  
    local.get $addr
    i32.const 4
    call $debug_log_write_post
  )

  (func $debug_i32.store8 (param $addr i32) (param $value i32) (param $offset i32) (param $align i32)
    local.get $addr
    i32.const 1
    call $debug_log_write_pre
    local.get $addr
    local.get $offset
    i32.add
    local.get $value
    i32.store8
    local.get $addr
    i32.const 1
    call $debug_log_write_post
  )

  (func $debug_i32.store16 (param $addr i32) (param $value i32) (param $offset i32) (param $align i32)
    local.get $addr
    i32.const 2
    call $debug_log_write_pre
    local.get $addr
    local.get $offset
    i32.add
    local.get $value
    i32.store16
    local.get $addr
    i32.const 2
    call $debug_log_write_post
  )

  (func $debug_i64.store (param $addr i32) (param $value i64) (param $offset i32) (param $align i32)
    local.get $addr
    i32.const 8
    call $debug_log_write_pre
    local.get $addr
    local.get $offset
    i32.add
    local.get $value
    i64.store
    local.get $addr
    i32.const 8
    call $debug_log_write_post
  )

  (func $debug_i64.store8 (param $addr i32) (param $value i64) (param $offset i32) (param $align i32)
    local.get $addr
    i32.const 1
    call $debug_log_write_pre
    local.get $addr
    local.get $offset
    i32.add
    local.get $value
    i64.store8
    local.get $addr
    i32.const 1
    call $debug_log_write_post
  )

  (func $debug_i64.store16 (param $addr i32) (param $value i64) (param $offset i32) (param $align i32)
    local.get $addr
    i32.const 2
    call $debug_log_write_pre
    local.get $addr
    local.get $offset
    i32.add
    local.get $value
    i64.store16
    local.get $addr
    i32.const 2
    call $debug_log_write_post
  )

  (func $debug_i64.store32 (param $addr i32) (param $value i64) (param $offset i32) (param $align i32)
    local.get $addr
    i32.const 4
    call $debug_log_write_pre
    local.get $addr
    local.get $offset
    i32.add
    local.get $value
    i64.store32
    local.get $addr
    i32.const 4
    call $debug_log_write_post
  )

  (func $debug_f32.store (param $addr i32) (param $value f32) (param $offset i32) (param $align i32)
    local.get $addr
    i32.const 4
    call $debug_log_write_pre
    local.get $addr
    local.get $offset
    i32.add
    local.get $value
    f32.store
    local.get $addr
    i32.const 4
    call $debug_log_write_post
  )

  (func $debug_f64.store (param $addr i32) (param $value f64) (param $offset i32) (param $align i32)
    local.get $addr
    i32.const 8
    call $debug_log_write_pre
    local.get $addr
    local.get $offset
    i32.add
    local.get $value
    f64.store
    local.get $addr
    i32.const 8
    call $debug_log_write_post
  )

  ;; These two log functions get called before and after writes...
  (func $debug_log_write_pre (param $addr i32) (param $len i32)

    i32.const offset.$debug_mem_access_start
    global.get $debug_start_mem
    i32.add
    i32.const length.$debug_mem_access_start
    call $debug_print

    ;; Print address
    local.get $addr
    call $db_format_i32_hex

    i32.const offset.$db_number_i32
    global.get $debug_start_mem
    i32.add
    i32.const 8
    ;;length.$db_number_i32
    call $debug_print

    i32.const offset.$debug_single_sp
    global.get $debug_start_mem
    i32.add
    i32.const length.$debug_single_sp
    call $debug_print

    local.get $addr
    local.get $len
    call $show_bytes

    i32.const offset.$debug_memory_change
    global.get $debug_start_mem
    i32.add
    i32.const length.$debug_memory_change
    call $debug_print

  )

  (func $debug_log_write_post (param $addr i32) (param $len i32)

    local.get $addr
    local.get $len
    call $show_bytes

    i32.const offset.$debug_newline
    global.get $debug_start_mem
    i32.add
    i32.const length.$debug_newline
    call $debug_print
  )

  (func $show_bytes (param $addr i32) (param $len i32)
    (local $count i32)
    i32.const 0
    local.set $count

    loop $loop

    local.get $count
    i32.eqz
    if
    else
      i32.const offset.$debug_single_sp
      global.get $debug_start_mem
      i32.add
      i32.const length.$debug_single_sp
      call $debug_print
    end

    local.get $count
    local.get $addr
    i32.add
    i32.load8_u

    ;; Now print it out...
    call $db_format_i32_hex

    ;; 00000000
    i32.const offset.$db_number_i32
    global.get $debug_start_mem
    i32.add
    i32.const 6
    i32.add
    i32.const 2
    call $debug_print

    local.get $count
    i32.const 1
    i32.add
    local.tee $count
    local.get $len
    i32.lt_u

    br_if $loop
    end
  )

  (data $debug_mem_access_start "# Memory store @ ")
)