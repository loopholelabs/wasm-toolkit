(module

;; $debug_exit_func
  (func $debug_exit_func (param $fid i32)
    (local $count i32)
    (local $metrics_ptr i32)

    ;; Update metrics
    local.get $fid
    i32.const 4
    i32.shl
    i32.const offset($debug_function_metrics)
    i32.add
    global.get $debug_start_mem
    i32.add
    local.tee $metrics_ptr

    call $debug_gettime

    ;; Pop timestamp off the timestamp stack
    global.get $debug_timestamps_stack_pointer
    i32.const 8
    i32.sub
    global.set $debug_timestamps_stack_pointer

    global.get $debug_timestamps_stack_pointer
    i32.const offset($debug_timestamps_stack)
    i32.add
    global.get $debug_start_mem
    i32.add
    i64.load

    i64.sub

    local.get $metrics_ptr
    i64.load offset=4
    i64.add

    i64.store offset=4

    global.get $debug_current_stack_depth
    i32.const 1
    i32.sub
    global.set $debug_current_stack_depth

    global.get $debug_current_stack_depth
    local.set $count

    block
      loop
        local.get $count
        i32.eqz
        br_if 1

        i32.const offset($debug_sp)
        global.get $debug_start_mem
        i32.add
        i32.const length($debug_sp)
        call $debug_print

        local.get $count
        i32.const 1
        i32.sub
        local.set $count
        br 0
      end
    end

    i32.const offset($debug_exit)
    global.get $debug_start_mem
    i32.add
    i32.const length($debug_exit)
    call $debug_print

    local.get $fid
    i32.const 3
    i32.shl
    global.get $debug_start_mem
    i32.add
    i32.const offset($debug_function_table)
    i32.add
    i32.load
    global.get $debug_start_mem
    i32.add
    i32.const offset($debug_function_data)
    i32.add
    local.get 0
    i32.const 3
    i32.shl
    global.get $debug_start_mem
    i32.add
    i32.const offset($debug_function_table)
    i32.add
    i32.load offset=4
    call $debug_print
  )

;; $debug_exit_func_i32
  (func $debug_exit_func_i32 (param $value i32) (result i32)
    i32.const offset($debug_return_value)
    global.get $debug_start_mem
    i32.add
    i32.const length($debug_return_value)
    call $debug_print
    i32.const offset($debug_value_i32)
    global.get $debug_start_mem
    i32.add
    i32.const length($debug_value_i32)
    call $debug_print

    local.get $value
    call $db_format_i32_hex

    i32.const offset($db_number_i32)
    global.get $debug_start_mem
    i32.add
    i32.const length($db_number_i32)
    call $debug_print

    i32.const offset($debug_newline)
    global.get $debug_start_mem
    i32.add
    i32.const length($debug_newline)
    call $debug_print

    call $debug_summary_maybe
    local.get $value
  )

;; $debug_exit_func_i64
  (func $debug_exit_func_i64 (param $value i64) (result i64)
    i32.const offset($debug_return_value)
    global.get $debug_start_mem
    i32.add
    i32.const length($debug_return_value)
    call $debug_print
    i32.const offset($debug_value_i64)
    global.get $debug_start_mem
    i32.add
    i32.const length($debug_value_i64)
    call $debug_print

    local.get $value
    i64.const 32
    i64.shr_u
    i32.wrap_i64
    call $db_format_i32_hex

    i32.const offset($db_number_i32)
    global.get $debug_start_mem
    i32.add
    i32.const length($db_number_i32)
    call $debug_print

    local.get $value
    i32.wrap_i64
    call $db_format_i32_hex

    i32.const offset($db_number_i32)
    global.get $debug_start_mem
    i32.add
    i32.const length($db_number_i32)
    call $debug_print

    i32.const offset($debug_newline)
    global.get $debug_start_mem
    i32.add
    i32.const length($debug_newline)
    call $debug_print

    call $debug_summary_maybe
    local.get $value
  )

;; $debug_exit_func_f32
  (func $debug_exit_func_f32 (param $value f32) (result f32)
    i32.const offset($debug_return_value)
    global.get $debug_start_mem
    i32.add
    i32.const length($debug_return_value)
    call $debug_print
    i32.const offset($debug_value_f32)
    global.get $debug_start_mem
    i32.add
    i32.const length($debug_value_f32)
    call $debug_print

    ;; TODO

    i32.const offset($debug_newline)
    global.get $debug_start_mem
    i32.add
    i32.const length($debug_newline)
    call $debug_print

    call $debug_summary_maybe
    local.get $value
  )

;; $debug_exit_func_f64
  (func $debug_exit_func_f64 (param $value f64) (result f64)
    i32.const offset($debug_return_value)
    global.get $debug_start_mem
    i32.add
    i32.const length($debug_return_value)
    call $debug_print
    i32.const offset($debug_value_f64)
    global.get $debug_start_mem
    i32.add
    i32.const length($debug_value_f64)
    call $debug_print

    ;; TODO

    i32.const offset($debug_newline)
    global.get $debug_start_mem
    i32.add
    i32.const length($debug_newline)
    call $debug_print

    call $debug_summary_maybe
    local.get $value
  )

;; $debug_exit_func_none
  (func $debug_exit_func_none
    i32.const offset($debug_newline)
    global.get $debug_start_mem
    i32.add
    i32.const length($debug_newline)
    call $debug_print

    call $debug_summary_maybe
  )

)