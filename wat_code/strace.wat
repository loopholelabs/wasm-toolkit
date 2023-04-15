(module

  (func $wt_watch_i32 (param $name_ptr i32) (param $name_len i32) (param $addr i32)
    global.get $wt_color
    if
      i32.const offset($wt_ansi_watch)
      i32.const length($wt_ansi_watch)
      call $wt_print
    end

    i32.const offset($wt_debug_watch)
    i32.const length($wt_debug_watch)
    call $wt_print
    local.get 0
    local.get 1
    call $wt_print

    i32.const offset($wt_debug_watch_eq)
    i32.const length($wt_debug_watch_eq)
    call $wt_print

    local.get 2
    i32.load

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
  )




  (func $debug_func_wasi_context (param $ptr i32) (param $len i32)
    global.get $wt_color
    if
      i32.const offset($wt_ansi_wasi_context)
      i32.const length($wt_ansi_wasi_context)
      call $wt_print
    end  

    local.get $ptr
    local.get $len
    call $wt_print  
  )

  (func $debug_func_wasi_done
    i32.const offset($dd_wasi_var_end)
    i32.const length($dd_wasi_var_end)
    call $wt_print

    global.get $wt_color
    if
      i32.const offset($wt_ansi_none)
      i32.const length($wt_ansi_none)
      call $wt_print
    end  
  )

  (func $debug_func_wasi_done_string
    i32.const offset($dd_wasi_var_end_string)
    i32.const length($dd_wasi_var_end_string)
    call $wt_print

    global.get $wt_color
    if
      i32.const offset($wt_ansi_none)
      i32.const length($wt_ansi_none)
      call $wt_print
    end  
  )


  (func $debug_param_separator
    i32.const offset($debug_param_sep)
    i32.const length($debug_param_sep)
    call $wt_print  
  )

  (func $debug_func_context (param $str_ptr i32) (param $str_len i32)
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

    global.get $wt_color
    if
      i32.const offset($wt_ansi_context)
      i32.const length($wt_ansi_context)
      call $wt_print
    end

    local.get $str_ptr
    local.get $str_len
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
  )

  (func $debug_param_name (param $str_ptr i32) (param $str_len i32)
    global.get $wt_color
    if
      i32.const offset($wt_ansi_param_name)
      i32.const length($wt_ansi_param_name)
      call $wt_print
    end

    local.get $str_ptr
    local.get $str_len
    call $wt_print

    global.get $wt_color
    if
      i32.const offset($wt_ansi_none)
      i32.const length($wt_ansi_none)
      call $wt_print
    end

    i32.const offset($debug_param_name_end)
    i32.const length($debug_param_name_end)
    call $wt_print
  )

(func $debug_exit_func_wasi (param $value i32) (result i32)
    (local $err_offset i32)
    (local $err_length i32)
    
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
      i32.const offset($wt_ansi_wasi_context)
      i32.const length($wt_ansi_wasi_context)
      call $wt_print
    end

;; Lookup the wasi error message...
    local.get $value
    i32.const 77
    i32.lt_u
    if
      i32.const offset($debug_sp)
      i32.const length($debug_sp)
      call $wt_print

      i32.const offset($wasi_errors)
      local.get $value
      i32.const 3
      i32.shl
      i32.add
      i32.load
      i32.const offset($wasi_error_messages)
      i32.add

      i32.const offset($wasi_errors)
      local.get $value
      i32.const 3
      i32.shl
      i32.add
      i32.load offset=4
      ;; Length

      call $wt_print
    end

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


  (func $debug_strlen (param $ptr i32) (result i32)
    (local $count i32)

    block
      loop
        local.get $count
        local.get $ptr
        i32.add
        i32.load8_u
        i32.eqz
        br_if 1

        local.get $count
        i32.const 1
        i32.add
        local.set $count
        br 0
      end
    end

    local.get $count
  )

  (func $dd_wasi_get_something (param $argv i32) (param $argvBuf i32) (param $len i32) (param $str_ptr i32) (param $str_len i32)
    (local $count i32)

    block
      loop
        local.get $count
        local.get $len
        i32.eq
        br_if 1

        ;; Print out the arg here...
        local.get $str_ptr
        local.get $str_len
				call $debug_func_wasi_context

        local.get $argv
        i32.load

        local.get $argv
        i32.load
        call $debug_strlen
        call $wt_print
									
				call $debug_func_wasi_done_string

        local.get $argv
        i32.const 4
        i32.add
        local.set $argv

        local.get $count
        i32.const 1
        i32.add
        local.set $count
        br 0
      end
    end
  )

  (data $debug_param_name_end "=")

  (data $debug_newline "\0d\0a")
  (data $debug_enter "-> ")
  (data $debug_exit "<- ")
  (data $debug_single_sp " ")
  (data $debug_sp "  ")
  (data $debug_table_sep " | ")
  (data $debug_memory_change " => ")

  (data $wt_debug_watch " * Watched global ")
  (data $wt_debug_watch_eq " = ")

  (data $dd_wasi_res_path " =>path = \22")
  (data $dd_wasi_res_bytes " =>bytes = ")
  (data $dd_wasi_res_numargs " =>num_args = ")
  (data $dd_wasi_res_sizeargs " =>size_args = ")
  (data $dd_wasi_res_args " =>args = \22")
  (data $dd_wasi_res_numenvs " =>num_envs = ")
  (data $dd_wasi_res_sizeenvs " =>size_envs = ")
  (data $dd_wasi_res_envs " =>envs = \22")
  (data $dd_wasi_res_timestamp " =>timestamp = ")

  (data $dd_wasi_var_path "   path = \22")
  (data $dd_wasi_var_rename "\22 -> \22")
  (data $dd_wasi_var_end_string "\22\0d\0a")
  (data $dd_wasi_var_end "\0d\0a")

  (data $error_stack_overflow "Error: The timings stack overflowed. You win some you lose some I guess.\0d\0a")

  (global $debug_current_stack_depth (mut i32) (i32.const 0))

  (global $wasi_result_args_get_count (mut i32) (i32.const 0))
  (global $wasi_result_envs_get_count (mut i32) (i32.const 0))

  (global $wt_all_function_length i32 (i32.const 0))

  (global $debug_do_timings i32 (i32.const 0))
)