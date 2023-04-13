(module

;; $debug_enter_func
  (func $debug_enter_func (param $fid i32)
    (local $count i32)
    (local $metrics_ptr i32)
    global.get $debug_current_stack_depth
    local.set $count

    ;; Update metrics
    local.get $fid
    i32.const 4
    i32.shl
    i32.const offset($debug_function_metrics)
    i32.add
    global.get $debug_start_mem
    i32.add
    local.tee $metrics_ptr
    local.get $metrics_ptr
    i32.load
    i32.const 1
    i32.add
    i32.store

    ;; Store enter timestamp in metrics stack
    global.get $debug_timestamps_stack_pointer
    i32.const offset($debug_timestamps_stack)
    i32.add
    global.get $debug_start_mem
    i32.add
    call $debug_gettime
    i64.store

    global.get $debug_timestamps_stack_pointer
    i32.const 8
    i32.add
    global.set $debug_timestamps_stack_pointer

    ;; TODO: Detect stack overflow here...

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

    i32.const offset($debug_enter)
    global.get $debug_start_mem
    i32.add
    i32.const length($debug_enter)
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

    i32.const offset($debug_param_start)
    global.get $debug_start_mem
    i32.add
    i32.const length($debug_param_start)
    call $debug_print

    global.get $debug_current_stack_depth
    i32.const 1
    i32.add
    global.set $debug_current_stack_depth
  )

(func $debug_enter_i32 (param $fid i32) (param $pid i32) (param $value i32)
    local.get $pid
    i32.eqz
    if
    else
      i32.const offset($debug_param_sep)
      global.get $debug_start_mem
      i32.add
      i32.const length($debug_param_sep)
      call $debug_print
    end

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
)

(func $debug_enter_i64 (param $fid i32) (param $pid i32) (param $value i64)
    local.get $pid
    i32.eqz
    if
    else
      i32.const offset($debug_param_sep)
      global.get $debug_start_mem
      i32.add
      i32.const length($debug_param_sep)
      call $debug_print
    end

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
)

(func $debug_enter_f32 (param $fid i32) (param $pid i32) (param $value f32)
    local.get $pid
    i32.eqz
    if
    else
      i32.const offset($debug_param_sep)
      global.get $debug_start_mem
      i32.add
      i32.const length($debug_param_sep)
      call $debug_print
    end

    i32.const offset($debug_value_f32)
    global.get $debug_start_mem
    i32.add
    i32.const length($debug_value_f32)
    call $debug_print

    ;; TODO
)

(func $debug_enter_f64 (param $fid i32) (param $pid i32) (param $value f64)
    local.get $pid
    i32.eqz
    if
    else
      i32.const offset($debug_param_sep)
      global.get $debug_start_mem
      i32.add
      i32.const length($debug_param_sep)
      call $debug_print
    end

    i32.const offset($debug_value_f64)
    global.get $debug_start_mem
    i32.add
    i32.const length($debug_value_f64)
    call $debug_print

    ;; TODO
)

(func $debug_enter_end (param $fid i32)
    i32.const offset($debug_param_end)
    global.get $debug_start_mem
    i32.add
    i32.const length($debug_param_end)
    call $debug_print

    i32.const offset($debug_newline)
    global.get $debug_start_mem
    i32.add
    i32.const length($debug_newline)
    call $debug_print
)

)