(module

  ;; wt_print_function_name - Given a function ID, print out the function name.
  (func $wt_print_function_name (param $fid i32)
    (local $ptr i32)
    i32.const offset($wt_all_function_names_locs)
    local.get $fid
    i32.const 3
    i32.shl
    i32.add
    local.tee $ptr

    ;; Get the address
    i32.load

    ;; We need to adjust this address manually
    i32.const offset($wt_all_function_names)
    i32.add

    ;; Get the offset
    local.get $ptr
    i32.load offset=4
    call $wt_print
  )

  ;; debug_enter_func - Called when a function is first entered.
  (func $debug_enter_func (param $fid i32)
    (local $count i32)
    global.get $debug_current_stack_depth
    local.set $count

    block
      loop
        local.get $count
        i32.eqz
        br_if 1

        i32.const offset($debug_sp)
        i32.const length($debug_sp)
        call $wt_print

        local.get $count
        i32.const 1
        i32.sub
        local.set $count
        br 0
      end
    end

    global.get $debug_current_stack_depth
    i32.const 1
    i32.add
    global.set $debug_current_stack_depth

    i32.const offset($debug_enter)
    i32.const length($debug_enter)
    call $wt_print

    local.get $fid
    call $wt_print_function_name

    i32.const offset($debug_param_start)
    i32.const length($debug_param_start)
    call $wt_print
  )

  ;; debug_exit_func - Called when we first exit a function
  (func $debug_exit_func (param $fid i32)
    (local $count i32)
  
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
        i32.const length($debug_sp)
        call $wt_print

        local.get $count
        i32.const 1
        i32.sub
        local.set $count
        br 0
      end
    end

    i32.const offset($debug_exit)
    i32.const length($debug_exit)
    call $wt_print

    local.get $fid
    call $wt_print_function_name
  )

  ;; debug_exit_func_i32 - Exit a func with an i32 result
  (func $debug_exit_func_i32 (param $value i32) (result i32)
    global.get $wt_color
    if
      i32.const offset($wt_ansi_result)
      i32.const length($wt_ansi_result)
      call $wt_print
    end

    i32.const offset($debug_return_value)
    i32.const length($debug_return_value)
    call $wt_print
    i32.const offset($debug_value_i32)
    i32.const length($debug_value_i32)
    call $wt_print

    local.get $value
    call $wt_format_i32_hex

    i32.const offset($db_number_i32)
    i32.const 8
    call $wt_print

    global.get $wt_color
    if
      i32.const offset($wt_ansi_none)
      i32.const length($wt_ansi_none)
      call $wt_print
    end

    i32.const offset($debug_newline)
    i32.const length($debug_newline)
    call $wt_print

    call $debug_summary_maybe
    local.get $value
  )

  ;; debug_exit_func_i64 - Exit a func with an i64 result
  (func $debug_exit_func_i64 (param $value i64) (result i64)
    global.get $wt_color
    if
      i32.const offset($wt_ansi_result)
      i32.const length($wt_ansi_result)
      call $wt_print
    end

    i32.const offset($debug_return_value)
    i32.const length($debug_return_value)
    call $wt_print
    i32.const offset($debug_value_i64)
    i32.const length($debug_value_i64)
    call $wt_print

    local.get $value
    i64.const 32
    i64.shr_u
    i32.wrap_i64
    call $wt_format_i32_hex

    i32.const offset($db_number_i32)
    i32.const 8
    call $wt_print

    local.get $value
    i32.wrap_i64
    call $wt_format_i32_hex

    i32.const offset($db_number_i32)
    i32.const 8
    call $wt_print

    global.get $wt_color
    if
      i32.const offset($wt_ansi_none)
      i32.const length($wt_ansi_none)
      call $wt_print
    end

    i32.const offset($debug_newline)
    i32.const length($debug_newline)
    call $wt_print

    call $debug_summary_maybe
    local.get $value
  )

  ;; debug_exit_func_f32 - Exit a func with an f32 result
  (func $debug_exit_func_f32 (param $value f32) (result f32)
    global.get $wt_color
    if
      i32.const offset($wt_ansi_result)
      i32.const length($wt_ansi_result)
      call $wt_print
    end

    i32.const offset($debug_return_value)
    i32.const length($debug_return_value)
    call $wt_print
    i32.const offset($debug_value_f32)
    i32.const length($debug_value_f32)
    call $wt_print

    ;; TODO

    global.get $wt_color
    if
      i32.const offset($wt_ansi_none)
      i32.const length($wt_ansi_none)
      call $wt_print
    end

    i32.const offset($debug_newline)
    i32.const length($debug_newline)
    call $wt_print

    call $debug_summary_maybe
    local.get $value
  )

  ;; debug_exit_func_f64 - Exit a func with an f64 result
  (func $debug_exit_func_f64 (param $value f64) (result f64)
    global.get $wt_color
    if
      i32.const offset($wt_ansi_result)
      i32.const length($wt_ansi_result)
      call $wt_print
    end

    i32.const offset($debug_return_value)
    i32.const length($debug_return_value)
    call $wt_print
    i32.const offset($debug_value_f64)
    i32.const length($debug_value_f64)
    call $wt_print

    ;; TODO

    global.get $wt_color
    if
      i32.const offset($wt_ansi_none)
      i32.const length($wt_ansi_none)
      call $wt_print
    end

    i32.const offset($debug_newline)
    i32.const length($debug_newline)
    call $wt_print

    call $debug_summary_maybe
    local.get $value
  )

  ;; debug_exit_func_none - Exit a func with no result
  (func $debug_exit_func_none
    i32.const offset($debug_newline)
    i32.const length($debug_newline)
    call $wt_print

    call $debug_summary_maybe
  )

  ;; debug_enter_i32 - Function entry with an i32
  (func $debug_enter_i32 (param $fid i32) (param $pid i32) (param $value i32)
    global.get $wt_color
    if
      i32.const offset($wt_ansi_param)
      i32.const length($wt_ansi_param)
      call $wt_print
    end

    i32.const offset($debug_value_i32)
    i32.const length($debug_value_i32)
    call $wt_print

    local.get $value
    call $wt_format_i32_hex

    i32.const offset($db_number_i32)
    i32.const 8
    call $wt_print

    global.get $wt_color
    if
      i32.const offset($wt_ansi_none)
      i32.const length($wt_ansi_none)
      call $wt_print
    end
  )

  ;; debug_enter_i64 - Function entry with an i64
  (func $debug_enter_i64 (param $fid i32) (param $pid i32) (param $value i64)
    global.get $wt_color
    if
      i32.const offset($wt_ansi_param)
      i32.const length($wt_ansi_param)
      call $wt_print
    end

    i32.const offset($debug_value_i64)
    i32.const length($debug_value_i64)
    call $wt_print

    local.get $value
    i64.const 32
    i64.shr_u
    i32.wrap_i64
    call $wt_format_i32_hex

    i32.const offset($db_number_i32)
    i32.const 8
    call $wt_print

    local.get $value
    i32.wrap_i64
    call $wt_format_i32_hex

    i32.const offset($db_number_i32)
    i32.const 8
    call $wt_print

    global.get $wt_color
    if
      i32.const offset($wt_ansi_none)
      i32.const length($wt_ansi_none)
      call $wt_print
    end
  )

  ;; debug_enter_f32 - Function entry with an f32
  (func $debug_enter_f32 (param $fid i32) (param $pid i32) (param $value f32)
    global.get $wt_color
    if
      i32.const offset($wt_ansi_param)
      i32.const length($wt_ansi_param)
      call $wt_print
    end

    i32.const offset($debug_value_f64)
    i32.const length($debug_value_f64)
    call $wt_print

    ;; TODO
    global.get $wt_color
    if
      i32.const offset($wt_ansi_none)
      i32.const length($wt_ansi_none)
      call $wt_print
    end
  )

  ;; debug_enter_f64 - Function entry with an f64
  (func $debug_enter_f64 (param $fid i32) (param $pid i32) (param $value f64)
    global.get $wt_color
    if
      i32.const offset($wt_ansi_param)
      i32.const length($wt_ansi_param)
      call $wt_print
    end

    i32.const offset($debug_value_f64)
    i32.const length($debug_value_f64)
    call $wt_print

    ;; TODO
    global.get $wt_color
    if
      i32.const offset($wt_ansi_none)
      i32.const length($wt_ansi_none)
      call $wt_print
    end
  )

  ;; debug_enter_end - Function entry completed
  (func $debug_enter_end (param $fid i32)
    i32.const offset($debug_param_end)
    i32.const length($debug_param_end)
    call $wt_print

    i32.const offset($debug_newline)
    i32.const length($debug_newline)
    call $wt_print
  )

  (data $debug_param_start "(")
  (data $debug_param_sep ", ")
  (data $debug_param_end ")")
  (data $debug_return_value " => ")
  (data $debug_value_i32 "i32:")
  (data $debug_value_i64 "i64:")
  (data $debug_value_f32 "f32:")
  (data $debug_value_f64 "f64:")

)