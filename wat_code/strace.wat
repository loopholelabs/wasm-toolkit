(module
  (type (func (param i32 i32 i32 i32) (result i32)))
  (import "wasi_snapshot_preview1" "fd_write" (func $debug_fd_write (type 0)))


  (func $debug_func_wasi_context (param $ptr i32) (param $len i32)
    global.get $debug_color
    if
      i32.const offset($debug_ansi_wasi_context)
      i32.const length($debug_ansi_wasi_context)
      call $debug_print
    end  

    local.get $ptr
    local.get $len
    call $debug_print  
  )

  (func $debug_func_wasi_done
    i32.const offset($dd_wasi_var_end)
    i32.const length($dd_wasi_var_end)
    call $debug_print

    global.get $debug_color
    if
      i32.const offset($debug_ansi_none)
      i32.const length($debug_ansi_none)
      call $debug_print
    end  
  )

  (func $debug_func_wasi_done_string
    i32.const offset($dd_wasi_var_end_string)
    i32.const length($dd_wasi_var_end_string)
    call $debug_print

    global.get $debug_color
    if
      i32.const offset($debug_ansi_none)
      i32.const length($debug_ansi_none)
      call $debug_print
    end  
  )

  (func $debug_enter_func (param $fid i32) (param $str_ptr i32) (param $str_len i32)
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
        call $debug_print

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
    call $debug_print

    local.get $str_ptr
    local.get $str_len
    call $debug_print

    i32.const offset($debug_param_start)
    i32.const length($debug_param_start)
    call $debug_print
  )

  (func $debug_exit_func (param $fid i32) (param $str_ptr i32) (param $str_len i32)
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
        call $debug_print

        local.get $count
        i32.const 1
        i32.sub
        local.set $count
        br 0
      end
    end

    i32.const offset($debug_exit)
    i32.const length($debug_exit)
    call $debug_print

    local.get $str_ptr
    local.get $str_len
    call $debug_print
  )

  (func $debug_param_separator
    i32.const offset($debug_param_sep)
    i32.const length($debug_param_sep)
    call $debug_print  
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
        call $debug_print

        local.get $count
        i32.const 1
        i32.sub
        local.set $count
        br 0
      end
    end

    global.get $debug_color
    if
      i32.const offset($debug_ansi_context)
      i32.const length($debug_ansi_context)
      call $debug_print
    end

    local.get $str_ptr
    local.get $str_len
    call $debug_print

    global.get $debug_color
    if
      i32.const offset($debug_ansi_none)
      i32.const length($debug_ansi_none)
      call $debug_print
    end

    i32.const offset($debug_newline)
    i32.const length($debug_newline)
    call $debug_print
)

  (func $debug_param_name (param $str_ptr i32) (param $str_len i32)
    global.get $debug_color
    if
      i32.const offset($debug_ansi_param_name)
      i32.const length($debug_ansi_param_name)
      call $debug_print
    end

    local.get $str_ptr
    local.get $str_len
    call $debug_print

    global.get $debug_color
    if
      i32.const offset($debug_ansi_none)
      i32.const length($debug_ansi_none)
      call $debug_print
    end

    i32.const offset($debug_param_name_end)
    i32.const length($debug_param_name_end)
    call $debug_print
  )

(func $debug_exit_func_wasi (param $value i32) (result i32)
    (local $err_offset i32)
    (local $err_length i32)
    
    global.get $debug_color
    if
      i32.const offset($debug_ansi_result)
      i32.const length($debug_ansi_result)
      call $debug_print
    end

    i32.const offset($debug_return_value)
    i32.const length($debug_return_value)
    call $debug_print
    i32.const offset($debug_value_i32)
    i32.const length($debug_value_i32)
    call $debug_print

    local.get $value
    call $db_format_i32_hex

    i32.const offset($db_number_i32)
    i32.const 8
    call $debug_print

    global.get $debug_color
    if
      i32.const offset($debug_ansi_wasi_context)
      i32.const length($debug_ansi_wasi_context)
      call $debug_print
    end

;; Lookup the wasi error message...
    local.get $value
    i32.const 77
    i32.lt_u
    if
      i32.const offset($debug_sp)
      i32.const length($debug_sp)
      call $debug_print

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

      call $debug_print
    end

    global.get $debug_color
    if
      i32.const offset($debug_ansi_none)
      i32.const length($debug_ansi_none)
      call $debug_print
    end

    i32.const offset($debug_newline)
    i32.const length($debug_newline)
    call $debug_print

    call $debug_summary_maybe
    local.get $value
  )

  (func $debug_exit_func_i32 (param $value i32) (result i32)
    global.get $debug_color
    if
      i32.const offset($debug_ansi_result)
      i32.const length($debug_ansi_result)
      call $debug_print
    end

    i32.const offset($debug_return_value)
    i32.const length($debug_return_value)
    call $debug_print
    i32.const offset($debug_value_i32)
    i32.const length($debug_value_i32)
    call $debug_print

    local.get $value
    call $db_format_i32_hex

    i32.const offset($db_number_i32)
    i32.const 8
    call $debug_print

    global.get $debug_color
    if
      i32.const offset($debug_ansi_none)
      i32.const length($debug_ansi_none)
      call $debug_print
    end

    i32.const offset($debug_newline)
    i32.const length($debug_newline)
    call $debug_print

    call $debug_summary_maybe
    local.get $value
  )

;; $debug_exit_func_i64
  (func $debug_exit_func_i64 (param $value i64) (result i64)
    global.get $debug_color
    if
      i32.const offset($debug_ansi_result)
      i32.const length($debug_ansi_result)
      call $debug_print
    end

    i32.const offset($debug_return_value)
    i32.const length($debug_return_value)
    call $debug_print
    i32.const offset($debug_value_i64)
    i32.const length($debug_value_i64)
    call $debug_print

    local.get $value
    i64.const 32
    i64.shr_u
    i32.wrap_i64
    call $db_format_i32_hex

    i32.const offset($db_number_i32)
    i32.const 8
    call $debug_print

    local.get $value
    i32.wrap_i64
    call $db_format_i32_hex

    i32.const offset($db_number_i32)
    i32.const 8
    call $debug_print

    global.get $debug_color
    if
      i32.const offset($debug_ansi_none)
      i32.const length($debug_ansi_none)
      call $debug_print
    end

    i32.const offset($debug_newline)
    i32.const length($debug_newline)
    call $debug_print

    call $debug_summary_maybe
    local.get $value
  )

;; $debug_exit_func_f32
  (func $debug_exit_func_f32 (param $value f32) (result f32)
    global.get $debug_color
    if
      i32.const offset($debug_ansi_result)
      i32.const length($debug_ansi_result)
      call $debug_print
    end

    i32.const offset($debug_return_value)
    i32.const length($debug_return_value)
    call $debug_print
    i32.const offset($debug_value_f32)
    i32.const length($debug_value_f32)
    call $debug_print

    ;; TODO

    global.get $debug_color
    if
      i32.const offset($debug_ansi_none)
      i32.const length($debug_ansi_none)
      call $debug_print
    end

    i32.const offset($debug_newline)
    i32.const length($debug_newline)
    call $debug_print

    call $debug_summary_maybe
    local.get $value
  )

;; $debug_exit_func_f64
  (func $debug_exit_func_f64 (param $value f64) (result f64)
    global.get $debug_color
    if
      i32.const offset($debug_ansi_result)
      i32.const length($debug_ansi_result)
      call $debug_print
    end

    i32.const offset($debug_return_value)
    i32.const length($debug_return_value)
    call $debug_print
    i32.const offset($debug_value_f64)
    i32.const length($debug_value_f64)
    call $debug_print

    ;; TODO

    global.get $debug_color
    if
      i32.const offset($debug_ansi_none)
      i32.const length($debug_ansi_none)
      call $debug_print
    end

    i32.const offset($debug_newline)
    i32.const length($debug_newline)
    call $debug_print

    call $debug_summary_maybe
    local.get $value
  )

;; $debug_exit_func_none
  (func $debug_exit_func_none
    i32.const offset($debug_newline)
    i32.const length($debug_newline)
    call $debug_print

    call $debug_summary_maybe
  )

  (func $debug_summary_maybe
  )

  (func $debug_enter_i32 (param $fid i32) (param $pid i32) (param $value i32)
    global.get $debug_color
    if
      i32.const offset($debug_ansi_param)
      i32.const length($debug_ansi_param)
      call $debug_print
    end

    i32.const offset($debug_value_i32)
    i32.const length($debug_value_i32)
    call $debug_print

    local.get $value
    call $db_format_i32_hex

    i32.const offset($db_number_i32)
    i32.const 8
    call $debug_print

    global.get $debug_color
    if
      i32.const offset($debug_ansi_none)
      i32.const length($debug_ansi_none)
      call $debug_print
    end
  )

  (func $debug_enter_i64 (param $fid i32) (param $pid i32) (param $value i64)
    global.get $debug_color
    if
      i32.const offset($debug_ansi_param)
      i32.const length($debug_ansi_param)
      call $debug_print
    end

    i32.const offset($debug_value_i64)
    i32.const length($debug_value_i64)
    call $debug_print

    local.get $value
    i64.const 32
    i64.shr_u
    i32.wrap_i64
    call $db_format_i32_hex

    i32.const offset($db_number_i32)
    i32.const 8
    call $debug_print

    local.get $value
    i32.wrap_i64
    call $db_format_i32_hex

    i32.const offset($db_number_i32)
    i32.const 8
    call $debug_print

    global.get $debug_color
    if
      i32.const offset($debug_ansi_none)
      i32.const length($debug_ansi_none)
      call $debug_print
    end
  )

  (func $debug_enter_f32 (param $fid i32) (param $pid i32) (param $value f32)
    global.get $debug_color
    if
      i32.const offset($debug_ansi_param)
      i32.const length($debug_ansi_param)
      call $debug_print
    end

    i32.const offset($debug_value_f64)
    i32.const length($debug_value_f64)
    call $debug_print

    ;; TODO
    global.get $debug_color
    if
      i32.const offset($debug_ansi_none)
      i32.const length($debug_ansi_none)
      call $debug_print
    end
  )

  (func $debug_enter_f64 (param $fid i32) (param $pid i32) (param $value f64)
    global.get $debug_color
    if
      i32.const offset($debug_ansi_param)
      i32.const length($debug_ansi_param)
      call $debug_print
    end

    i32.const offset($debug_value_f64)
    i32.const length($debug_value_f64)
    call $debug_print

    ;; TODO
    global.get $debug_color
    if
      i32.const offset($debug_ansi_none)
      i32.const length($debug_ansi_none)
      call $debug_print
    end
  )

  (func $debug_enter_end (param $fid i32)
    i32.const offset($debug_param_end)
    i32.const length($debug_param_end)
    call $debug_print

    i32.const offset($debug_newline)
    i32.const length($debug_newline)
    call $debug_print
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

  (func $debug_print (param $ptr i32) (param $len i32)
    (local $iovp i32)

    i32.const offset($iovec)

    local.tee $iovp
    local.get $ptr
    i32.store

    local.get $iovp
    local.get $len
    i32.store offset=4

    i32.const 2
    local.get $iovp
    i32.const 1
    i32.const offset($bytes_written)
    call $debug_fd_write
    drop
  )

;; $db_format_i32 as hex into the buffer ($db_number_i32)
  (func $db_format_i32_hex (param $num i32)
    (local $count i32)
    (local $ptr i32)
    (local $shift_val i32)

    i32.const offset($db_number_i32)
    local.set $ptr

    i32.const 28
    local.set $shift_val

    ;; 28 / 24 / 20 / 16 / 12 / 8 / 4 / 0

    i32.const 0
    local.set $count

    loop
      ;; Work out the value to store...

      local.get $ptr

      local.get $num
      local.get $shift_val
      i32.shr_u
      i32.const 15
      i32.and

      i32.const offset($db_hex)
      i32.add

      i32.load8_u

      i32.store8

      local.get $ptr
      i32.const 1
      i32.add
      local.set $ptr

      local.get $shift_val
      i32.const 4
      i32.sub
      local.set $shift_val

      local.get $count
      i32.const 1
      i32.add
      local.tee $count
      i32.const 8
      i32.lt_u
      br_if 0
    end
  )

;; $db_format_i32 as dec into the buffer ($db_number_i32)
  (func $db_format_i32_dec (param $num i32)
    (local $count i32)
    (local $ptr i32)
    (local $divide_val i32)
    (local $in_number i32)
    (local $store_value i32)

    i32.const offset($db_number_i32)
    local.set $ptr

    i32.const 1000000000
    local.set $divide_val

    i32.const 0
    local.set $count

    loop
      ;; Work out the value to store...

      local.get $ptr

      local.get $num
      local.get $divide_val
      i32.div_u
      i32.const 10
      i32.rem_u

      i32.const offset($db_hex)
      i32.add

      i32.load8_u
      i32.store8

      local.get $ptr
      i32.const 1
      i32.add
      local.set $ptr

      local.get $divide_val
      i32.const 10
      i32.div_u
      local.set $divide_val

      local.get $count
      i32.const 1
      i32.add
      local.tee $count
      i32.const 10
      i32.lt_u
      br_if 0
    end
  )

;; $db_format_i32 as dec into the buffer ($db_number_i32)
  (func $db_format_i32_dec_nz (param $num i32)
    (local $count i32)
    (local $ptr i32)
    (local $divide_val i32)
    (local $in_number i32)
    (local $store_value i32)

    i32.const offset($db_number_i32)
    local.set $ptr

    i32.const 1000000000
    local.set $divide_val

    i32.const 0
    local.set $count

    loop
      ;; Work out the value to store...

      local.get $ptr

      local.get $num
      local.get $divide_val
      i32.div_u
      i32.const 10
      i32.rem_u
      local.tee $store_value
      i32.eqz
      if (result i32)
        local.get $in_number
        if (result i32)
          local.get $store_value
        else
          i32.const 16
        end
      else
        i32.const 1
        local.set $in_number
        local.get $store_value
      end

      i32.const offset($db_hex)
      i32.add

      i32.load8_u
      i32.store8

      local.get $ptr
      i32.const 1
      i32.add
      local.set $ptr

      local.get $divide_val
      i32.const 10
      i32.div_u
      local.set $divide_val

      local.get $count
      i32.const 1
      i32.add
      local.tee $count
      i32.const 10
      i32.lt_u
      br_if 0
    end
  )

;; $db_format_i64 as dec into the buffer ($db_number_i64)
  (func $db_format_i64_dec_nz (param $num i64)
    (local $count i32)
    (local $ptr i32)
    (local $divide_val i64)
    (local $in_number i32)
    (local $store_value i32)

    i32.const offset($db_number_i64)
    local.set $ptr

    i64.const 1000000000000000000
    local.set $divide_val

    i32.const 0
    local.set $count

    loop
      ;; Work out the value to store...

      local.get $ptr

      local.get $num
      local.get $divide_val
      i64.div_u
      i64.const 10
      i64.rem_u
      i32.wrap_i64
      local.tee $store_value
      i32.eqz
      if (result i32)
        local.get $in_number
        if (result i32)
          local.get $store_value
        else
          i32.const 16
        end
      else
        i32.const 1
        local.set $in_number
        local.get $store_value
      end

      i32.const offset($db_hex)
      i32.add

      i32.load8_u
      i32.store8

      local.get $ptr
      i32.const 1
      i32.add
      local.set $ptr

      local.get $divide_val
      i64.const 10
      i64.div_u
      local.set $divide_val

      local.get $count
      i32.const 1
      i32.add
      local.tee $count
      i32.const 19
      i32.lt_u
      br_if 0
    end
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
        call $debug_print
									
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

  (data $db_hex "0123456789ABCDEF ")
  (data $db_number_i32 10)
  (data $db_number_i64 19)

  (data $helloworld "Hello world\0d\0a")
  (data $byeworld "Bye world\0d\0a")
  (data $iovec 8)
  (data $bytes_written 4)

  (data $debug_return_value " => ")
  (data $debug_value_i32 "i32:")
  (data $debug_value_i64 "i64:")
  (data $debug_value_f32 "f32:")
  (data $debug_value_f64 "f64:")

  (data $debug_param_name_end "=")

  (data $debug_param_start "(")
  (data $debug_param_sep ", ")
  (data $debug_param_end ")")

  (data $debug_newline "\0d\0a")
  (data $debug_enter "-> ")
  (data $debug_exit "<- ")
  (data $debug_single_sp " ")
  (data $debug_sp "  ")
  (data $debug_table_sep " | ")
  (data $debug_memory_change " => ")

  (data $debug_ansi_param "\1b[32m")
  (data $debug_ansi_result "\1b[31m")
  (data $debug_ansi_context "\1b[36m")
  (data $debug_ansi_param_name "\1b[33m")
  (data $debug_ansi_none "\1b[0m")

  (data $debug_ansi_wasi_context "\1b[35m")

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


  (data $wasi_errors 0)
  (data $wasi_error_messages 0)

  (global $debug_current_stack_depth (mut i32) (i32.const 0))

  (global $debug_color i32 (i32.const 0))

  (global $wasi_result_args_get_count (mut i32) (i32.const 0))
  (global $wasi_result_envs_get_count (mut i32) (i32.const 0))

)