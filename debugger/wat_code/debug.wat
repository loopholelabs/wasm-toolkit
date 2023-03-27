(module

;; $debug_exit_func
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

        i32.const offset.$debug_sp
        global.get $debug_start_mem
        i32.add
        i32.const length.$debug_sp
        call $debug_print

        local.get $count
        i32.const 1
        i32.sub
        local.set $count
        br 0
      end
    end

    i32.const offset.$debug_exit
    global.get $debug_start_mem
    i32.add
    i32.const length.$debug_exit
    call $debug_print

    local.get $fid
    i32.const 3
    i32.shl
    global.get $debug_start_mem
    i32.add
    i32.const offset.$debug_function_table
    i32.add
    i32.load
    global.get $debug_start_mem
    i32.add
    i32.const offset.$debug_function_data
    i32.add
    local.get 0
    i32.const 3
    i32.shl
    global.get $debug_start_mem
    i32.add
    i32.const offset.$debug_function_table
    i32.add
    i32.load offset=4
    call $debug_print
  )

;; $debug_exit_func_i32
  (func $debug_exit_func_i32 (param $value i32) (result i32)
    i32.const offset.$debug_return_value
    global.get $debug_start_mem
    i32.add
    i32.const length.$debug_return_value
    call $debug_print
    i32.const offset.$debug_value_i32
    global.get $debug_start_mem
    i32.add
    i32.const length.$debug_value_i32
    call $debug_print

    local.get $value
    call $db_format_i32_hex

    i32.const offset.$db_number_i32
    global.get $debug_start_mem
    i32.add
    i32.const length.$db_number_i32
    call $debug_print

    i32.const offset.$debug_newline
    global.get $debug_start_mem
    i32.add
    i32.const length.$debug_newline
    call $debug_print

    call $debug_summary_maybe
    local.get $value
  )

;; $debug_exit_func_i64
  (func $debug_exit_func_i64 (param $value i64) (result i64)
    i32.const offset.$debug_return_value
    global.get $debug_start_mem
    i32.add
    i32.const length.$debug_return_value
    call $debug_print
    i32.const offset.$debug_value_i64
    global.get $debug_start_mem
    i32.add
    i32.const length.$debug_value_i64
    call $debug_print

    local.get $value
    i64.const 32
    i64.shr_u
    i32.wrap_i64
    call $db_format_i32_hex

    i32.const offset.$db_number_i32
    global.get $debug_start_mem
    i32.add
    i32.const length.$db_number_i32
    call $debug_print

    local.get $value
    i32.wrap_i64
    call $db_format_i32_hex

    i32.const offset.$db_number_i32
    global.get $debug_start_mem
    i32.add
    i32.const length.$db_number_i32
    call $debug_print

    i32.const offset.$debug_newline
    global.get $debug_start_mem
    i32.add
    i32.const length.$debug_newline
    call $debug_print

    call $debug_summary_maybe
    local.get $value
  )

;; $debug_exit_func_f32
  (func $debug_exit_func_f32 (param $value f32) (result f32)
    i32.const offset.$debug_return_value
    global.get $debug_start_mem
    i32.add
    i32.const length.$debug_return_value
    call $debug_print
    i32.const offset.$debug_value_f32
    global.get $debug_start_mem
    i32.add
    i32.const length.$debug_value_f32
    call $debug_print

    ;; TODO

    i32.const offset.$debug_newline
    global.get $debug_start_mem
    i32.add
    i32.const length.$debug_newline
    call $debug_print

    call $debug_summary_maybe
    local.get $value
  )

;; $debug_exit_func_f64
  (func $debug_exit_func_f64 (param $value f64) (result f64)
    i32.const offset.$debug_return_value
    global.get $debug_start_mem
    i32.add
    i32.const length.$debug_return_value
    call $debug_print
    i32.const offset.$debug_value_f64
    global.get $debug_start_mem
    i32.add
    i32.const length.$debug_value_f64
    call $debug_print

    ;; TODO

    i32.const offset.$debug_newline
    global.get $debug_start_mem
    i32.add
    i32.const length.$debug_newline
    call $debug_print

    call $debug_summary_maybe
    local.get $value
  )

;; $debug_exit_func_none
  (func $debug_exit_func_none
    i32.const offset.$debug_newline
    global.get $debug_start_mem
    i32.add
    i32.const length.$debug_newline
    call $debug_print

    call $debug_summary_maybe
  )

;; $debug_summary_maybe
  (func $debug_summary_maybe
    block $summary
      global.get $debug_current_stack_depth
      i32.const 0
      i32.gt_u
      br_if $summary
      call $debug_summary
    end  
  )

;; $debug_summary_func
  (func $debug_summary_func (param $fid i32)
    (local $metric_count i32)
    local.get $fid
    i32.const 2
    i32.shl
    i32.const offset.$debug_function_metrics
    i32.add
    global.get $debug_start_mem
    i32.add
    i32.load
    local.tee $metric_count
    i32.eqz
    br_if 0

    local.get $metric_count
    call $db_format_i32_hex

    i32.const offset.$db_number_i32
    global.get $debug_start_mem
    i32.add
    i32.const length.$db_number_i32
    call $debug_print

    i32.const offset.$debug_table_sep
    global.get $debug_start_mem
    i32.add
    i32.const length.$debug_table_sep
    call $debug_print

    local.get $fid
    i32.const 3
    i32.shl
    global.get $debug_start_mem
    i32.add
    i32.const offset.$debug_function_table
    i32.add
    i32.load
    global.get $debug_start_mem
    i32.add
    i32.const offset.$debug_function_data
    i32.add
    local.get 0
    i32.const 3
    i32.shl
    global.get $debug_start_mem
    i32.add
    i32.const offset.$debug_function_table
    i32.add
    i32.load offset=4
    call $debug_print

    i32.const offset.$debug_newline
    global.get $debug_start_mem
    i32.add
    i32.const length.$debug_newline
    call $debug_print
  
  )

;; $debug_summary
  (func $debug_summary
      (local $f_id i32)
      i32.const offset.$debug_summary
      global.get $debug_start_mem
      i32.add
      i32.const length.$debug_summary
      call $debug_print

;; Go through all functions and print details for the function
      i32.const 0
      local.set $f_id

      loop $loop
        local.get $f_id
        call $debug_summary_func

        local.get $f_id
        i32.const 1
        i32.add
        local.tee $f_id
        global.get $debug_num_funcs
        i32.lt_u
        br_if $loop
      end
  )

  (data $debug_return_value " => ")
  (data $debug_value_i32 "i32:")
  (data $debug_value_i64 "i64:")
  (data $debug_value_f32 "f32:")
  (data $debug_value_f64 "f64:")

  (data $debug_param_start "(")
  (data $debug_param_sep ", ")
  (data $debug_param_end ")")

  (data $debug_newline "\0d\0a")
  (data $debug_enter "-> ")
  (data $debug_exit "<- ")
  (data $debug_sp "  ")
  (data $debug_table_sep " | ")

  (data $debug_summary "\0d\0a-- Summary of execution --\0d\0a")

  (global $debug_current_stack_depth (mut i32) (i32.const 0))
)
