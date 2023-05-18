(module
  (type (func (param i32 i64 i32) (result i32)))
  (type (func (param i32 i32) (result i32)))
  (import "wasi_snapshot_preview1" "clock_time_get" (func $debug_clock_time_get (type 0)))
  (import "wasi_snapshot_preview1" "random_get" (func $debug_random_get (type 1)))

  ;; Get the current timestamp as an i64
  (func $debug_gettime (result i64)
    i32.const 0
    i64.const 1000
    i32.const offset($debug_clock_loc)
    call $debug_clock_time_get
    drop
    i32.const offset($debug_clock_loc)
    i64.load
  )

  (func $otel_output_trace_data (param $ptr i32) (param $len i32)
    ;; Copy it to the output buffer...

    ;; Detect buffer overflow here...
    global.get $otel_output_buffer_ptr
    local.get $len
    i32.add
    i32.const length($otel_output_buffer)
    i32.ge_u
    if
      i32.const offset($error_buffer_overflow)
      i32.const length($error_buffer_overflow)
      call $wt_print   
      unreachable
    end

    global.get $otel_output_buffer_ptr
    i32.const offset($otel_output_buffer)
    i32.add
    local.get $ptr
    local.get $len
    ;; dest / src / size
    memory.copy

    global.get $otel_output_buffer_ptr
    local.get $len
    i32.add
    global.set $otel_output_buffer_ptr
  )

  (func $otel_output_trace_data_flush
    ;; Flush
    i32.const offset($otel_output_buffer)
    global.get $otel_output_buffer_ptr
    call $send_otel_trace_json
    ;; call $wt_print

    ;; TODO: We can now do a single host call here to say we have some otel data.

    ;; Reset
    i32.const 0
    global.set $otel_output_buffer_ptr
  )

  ;; Enter a function for otel (Just pushes the time data onto a stack, creates a random spanID)
  ;;
  (func $otel_enter_func (param $fid i32)

    ;; Store enter timestamp in metrics stack
    global.get $debug_timestamps_stack_pointer
    i32.const offset($debug_timestamps_stack)
    i32.add
    call $debug_gettime
    i64.store

    ;; Create a random span ID
    global.get $debug_timestamps_stack_pointer
    i32.const offset($debug_timestamps_stack)
    i32.add
    i32.const 8
    i32.add
    i32.const 8
    call $debug_random_get
    drop

    ;; Accept it onto the stack...
    global.get $debug_timestamps_stack_pointer
    i32.const 16
    i32.add
    global.set $debug_timestamps_stack_pointer

    global.get $debug_timestamps_stack_pointer
    i32.const length($debug_timestamps_stack)
    i32.ge_u
    ;; Detect stack overflow
    if
      i32.const offset($error_stack_overflow)
      i32.const length($error_stack_overflow)
      call $wt_print
      unreachable
    end
  )

  ;; ORDER OF FUNCTION CALLS

  ;; otel_exit_func
  ;; * otel_exit_func_<TYPE>
  ;; <Optional> otel_exit_func_result_<TYPE>
  ;; otel_exit_func_done

  (func $otel_output_attr_string (param $name i32) (param $name_len i32) (param $val i32) (param $val_len i32)
    i32.const offset($ot_attr_start)
    i32.const length($ot_attr_start)
    call $otel_output_trace_data

    local.get $name
    local.get $name_len
    call $otel_output_trace_data

    i32.const offset($ot_attr_mid)
    i32.const length($ot_attr_mid)
    call $otel_output_trace_data

    i32.const offset($ot_attr_string_start)
    i32.const length($ot_attr_string_start)
    call $otel_output_trace_data

    local.get $val
    local.get $val_len
    call $otel_output_trace_data

    i32.const offset($ot_attr_string_end)
    i32.const length($ot_attr_string_end)
    call $otel_output_trace_data

    i32.const offset($ot_attr_end)
    i32.const length($ot_attr_end)
    call $otel_output_trace_data
  )

  (func $otel_output_attr_var_string (param $name i32) (param $name_len i32) (param $var i32) (param $var_len i32) (param $val i32) (param $val_len i32)
    i32.const offset($ot_attr_start)
    i32.const length($ot_attr_start)
    call $otel_output_trace_data

    local.get $name
    local.get $name_len
    call $otel_output_trace_data

    i32.const offset($ot_attr_mid)
    i32.const length($ot_attr_mid)
    call $otel_output_trace_data

    i32.const offset($ot_attr_string_start)
    i32.const length($ot_attr_string_start)
    call $otel_output_trace_data

    local.get $var
    local.get $var_len
    call $otel_output_trace_data

    i32.const offset($ot_equals)
    i32.const length($ot_equals)
    call $otel_output_trace_data

    local.get $val
    local.get $val_len
    call $otel_output_trace_data

    i32.const offset($ot_attr_string_end)
    i32.const length($ot_attr_string_end)
    call $otel_output_trace_data

    i32.const offset($ot_attr_end)
    i32.const length($ot_attr_end)
    call $otel_output_trace_data
  )

  ;; Output some hex data
  (func $otel_output_attr_hexdata (param $name i32) (param $name_len i32) (param $val i32) (param $val_len i32)
    (local $i i32)
    i32.const offset($ot_attr_start)
    i32.const length($ot_attr_start)
    call $otel_output_trace_data

    local.get $name
    local.get $name_len
    call $otel_output_trace_data

    i32.const offset($ot_attr_mid)
    i32.const length($ot_attr_mid)
    call $otel_output_trace_data

    i32.const offset($ot_attr_string_start)
    i32.const length($ot_attr_string_start)
    call $otel_output_trace_data

    ;; Loop round and output hex bytes...
    block
      loop
        local.get $i
        local.get $val_len
        i32.ge_u
        br_if 1

        ;; Now print out the hex byte...
        local.get $val
        local.get $i
        i32.add
        i32.load8_u

        ;; Byte is on stack ready to output.
        i32.const 4
        i32.shr_u

        i32.const offset($db_hex)
        i32.add
        i32.const 1
        call $otel_output_trace_data

        local.get $val
        local.get $i
        i32.add
        i32.load8_u
        i32.const 15
        i32.and

        i32.const offset($db_hex)
        i32.add
        i32.const 1
        call $otel_output_trace_data

        i32.const offset($ot_space)
        i32.const length($ot_space)
        call $otel_output_trace_data

        local.get $i
        i32.const 1
        i32.add
        local.set $i
        br 0
      end
    end

    i32.const offset($ot_encoded_speech)
    i32.const length($ot_encoded_speech)
    call $otel_output_trace_data


    i32.const 0
    local.set $i

    ;; Now output printable characters
    block
      loop
        local.get $i
        local.get $val_len
        i32.ge_u
        br_if 1

        ;; Now print out the hex byte...
        local.get $val
        local.get $i
        i32.add
        i32.load8_u

        ;; Check if it's safe to print
        call $is_printable
        if
          local.get $val
          local.get $i
          i32.add
          i32.const 1
          call $otel_output_trace_data
        else
          i32.const offset($ot_nonprintable)
          i32.const length($ot_nonprintable)
          call $otel_output_trace_data
        end

        local.get $i
        i32.const 1
        i32.add
        local.set $i
        br 0
      end
    end

    i32.const offset($ot_encoded_speech)
    i32.const length($ot_encoded_speech)
    call $otel_output_trace_data

    i32.const offset($ot_attr_string_end)
    i32.const length($ot_attr_string_end)
    call $otel_output_trace_data

    i32.const offset($ot_attr_end)
    i32.const length($ot_attr_end)
    call $otel_output_trace_data
  )

  (func $is_printable (param $val i32) (result i32)
    (local $p i32)
    block
      loop
        local.get $p
        i32.const length($_printable_characters)
        i32.ge_u
        br_if 1

        i32.const offset($_printable_characters)
        local.get $p
        i32.add
        i32.load8_u
        local.get $val
        i32.eq
        if
          i32.const 1
          return
        end

        local.get $p
        i32.const 1
        i32.add
        local.set $p
        br 0
      end
    end
    i32.const 0
  )

  (func $otel_exit_func_result_i32 (param $val i32) (param $fid i32) (result i32)
    i32.const offset($ot_comma)
    i32.const length($ot_comma)
    call $otel_output_trace_data

    local.get $val
    call $wt_format_i32_hex

    i32.const offset($ot_at_result)
    i32.const length($ot_at_result)
    i32.const offset($db_number_i32)
    i32.const 8
    call $otel_output_attr_string

    local.get $val
  )

  (func $otel_exit_func_result_i64 (param $val i64) (param $fid i32) (result i64)
    i32.const offset($ot_comma)
    i32.const length($ot_comma)
    call $otel_output_trace_data

    local.get $val
    call $wt_format_i64_hex

    i32.const offset($ot_at_result)
    i32.const length($ot_at_result)
    i32.const offset($db_number_i64)
    i32.const 16
    call $otel_output_attr_string

    local.get $val
  )

  (func $otel_exit_func_result_f32 (param $val f32) (param $fid i32) (result f32)
    i32.const offset($ot_comma)
    i32.const length($ot_comma)
    call $otel_output_trace_data

  ;; TODO

    i32.const offset($ot_at_result)
    i32.const length($ot_at_result)
    i32.const offset($ot_at_todo)
    i32.const length($ot_at_todo)
    call $otel_output_attr_string

    local.get $val
  )

  (func $otel_exit_func_result_f64 (param $val f64) (param $fid i32) (result f64)
    i32.const offset($ot_comma)
    i32.const length($ot_comma)
    call $otel_output_trace_data

  ;; TODO

    i32.const offset($ot_at_result)
    i32.const length($ot_at_result)
    i32.const offset($ot_at_todo)
    i32.const length($ot_at_todo)
    call $otel_output_attr_string

    local.get $val
  )

  (func $otel_watch_global (param $name_ptr i32) (param $name_len i32) (param $ptr i32) (param $len i32)
    i32.const offset($ot_comma)
    i32.const length($ot_comma)
    call $otel_output_trace_data

    local.get $len
    i32.const 4
    i32.eq
    if
      local.get $ptr
      i32.load
      call $wt_format_i32_hex

      local.get $name_ptr
      local.get $name_len
      i32.const offset($db_number_i32)
      i32.const 8
      call $otel_output_attr_string
      return
    end

    local.get $len
    i32.const 8
    i32.eq
    if
      local.get $ptr
      i64.load
      call $wt_format_i64_hex

      local.get $name_ptr
      local.get $name_len
      i32.const offset($db_number_i64)
      i32.const 16
      call $otel_output_attr_string
      return
    end


    local.get $name_ptr
    local.get $name_len
    local.get $ptr
    local.get $len
    call $otel_output_attr_hexdata
  )

  ;; Exit function with param i32
  (func $otel_exit_func_i32 (param $fid i32) (param $pid i32) (param $val i32)
    i32.const offset($ot_comma)
    i32.const length($ot_comma)
    call $otel_output_trace_data

    local.get $val
    call $wt_format_i32_hex

    local.get $pid
    i32.const offset($ot_at_param)
    i32.const 6
    i32.add
    call $wt_conv_byte_dec

    i32.const offset($ot_at_param)
    i32.const length($ot_at_param)
    i32.const offset($db_number_i32)
    i32.const 8
    call $otel_output_attr_string
  )

  ;; Exit function with param i64
  (func $otel_exit_func_i64 (param $fid i32) (param $pid i32) (param $val i64)
    i32.const offset($ot_comma)
    i32.const length($ot_comma)
    call $otel_output_trace_data

    local.get $val
    call $wt_format_i64_hex

    local.get $pid
    i32.const offset($ot_at_param)
    i32.const 6
    i32.add
    call $wt_conv_byte_dec

    i32.const offset($ot_at_param)
    i32.const length($ot_at_param)
    i32.const offset($db_number_i64)
    i32.const 16
    call $otel_output_attr_string
  )

  ;; Exit function with param f32
  (func $otel_exit_func_f32 (param $fid i32) (param $pid i32) (param $val f32)
    i32.const offset($ot_comma)
    i32.const length($ot_comma)
    call $otel_output_trace_data

    local.get $pid
    i32.const offset($ot_at_param)
    i32.const 6
    i32.add
    call $wt_conv_byte_dec

    i32.const offset($ot_at_param)
    i32.const length($ot_at_param)
    i32.const offset($ot_at_todo)
    i32.const length($ot_at_todo)
    call $otel_output_attr_string
  )

  ;; Exit function with param f64
  (func $otel_exit_func_f64 (param $fid i32) (param $pid i32) (param $val f64)
    i32.const offset($ot_comma)
    i32.const length($ot_comma)
    call $otel_output_trace_data

    local.get $pid
    i32.const offset($ot_at_param)
    i32.const 6
    i32.add
    call $wt_conv_byte_dec

    i32.const offset($ot_at_param)
    i32.const length($ot_at_param)
    i32.const offset($ot_at_todo)
    i32.const length($ot_at_todo)
    call $otel_output_attr_string
  )

  ;; Exit function with param i32
  (func $otel_exit_func_var_i32 (param $fid i32) (param $pid i32) (param $val i32) (param $var_ptr i32) (param $var_len i32)
    i32.const offset($ot_comma)
    i32.const length($ot_comma)
    call $otel_output_trace_data

    local.get $val
    call $wt_format_i32_hex

    local.get $pid
    i32.const offset($ot_at_param)
    i32.const 6
    i32.add
    call $wt_conv_byte_dec

    i32.const offset($ot_at_param)
    i32.const length($ot_at_param)
    local.get $var_ptr
    local.get $var_len
    i32.const offset($db_number_i32)
    i32.const 8
    call $otel_output_attr_var_string
  )

  ;; Exit function with param i64
  (func $otel_exit_func_var_i64 (param $fid i32) (param $pid i32) (param $val i64) (param $var_ptr i32) (param $var_len i32)
    i32.const offset($ot_comma)
    i32.const length($ot_comma)
    call $otel_output_trace_data

    local.get $val
    call $wt_format_i64_hex

    local.get $pid
    i32.const offset($ot_at_param)
    i32.const 6
    i32.add
    call $wt_conv_byte_dec

    i32.const offset($ot_at_param)
    i32.const length($ot_at_param)
    local.get $var_ptr
    local.get $var_len
    i32.const offset($db_number_i64)
    i32.const 16
    call $otel_output_attr_var_string
  )

  ;; Exit function with param f32
  (func $otel_exit_func_var_f32 (param $fid i32) (param $pid i32) (param $val f32) (param $var_ptr i32) (param $var_len i32)
    i32.const offset($ot_comma)
    i32.const length($ot_comma)
    call $otel_output_trace_data

    local.get $pid
    i32.const offset($ot_at_param)
    i32.const 6
    i32.add
    call $wt_conv_byte_dec

    i32.const offset($ot_at_param)
    i32.const length($ot_at_param)
    local.get $var_ptr
    local.get $var_len
    i32.const offset($ot_at_todo)
    i32.const length($ot_at_todo)
    call $otel_output_attr_var_string
  )

  ;; Exit function with param f64
  (func $otel_exit_func_var_f64 (param $fid i32) (param $pid i32) (param $val f64) (param $var_ptr i32) (param $var_len i32)
    i32.const offset($ot_comma)
    i32.const length($ot_comma)
    call $otel_output_trace_data

    local.get $pid
    i32.const offset($ot_at_param)
    i32.const 6
    i32.add
    call $wt_conv_byte_dec

    i32.const offset($ot_at_param)
    i32.const length($ot_at_param)
    local.get $var_ptr
    local.get $var_len
    i32.const offset($ot_at_todo)
    i32.const length($ot_at_todo)
    call $otel_output_attr_var_string
  )

  ;; Exit a function all done
  (func $otel_exit_func_done (param $fid i32)

    i32.const offset($ot_attributes_end)
    i32.const length($ot_attributes_end)
    call $otel_output_trace_data

    i32.const offset($ot_end)
    i32.const length($ot_end)
    call $otel_output_trace_data

    i32.const offset($debug_newline)
    i32.const length($debug_newline)
    call $otel_output_trace_data

    call $otel_output_trace_data_flush
  )

  ;; Exit a function. This is where the otel stuff gets sent out.
  (func $otel_exit_func (param $fid i32)
    (local $time_end i64)

    call $debug_gettime
    local.set $time_end

    ;; Pop timestamp off the timestamp stack
    global.get $debug_timestamps_stack_pointer
    i32.const 16
    i32.sub
    global.set $debug_timestamps_stack_pointer

    i32.const offset($ot_start)
    i32.const length($ot_start)
    call $otel_output_trace_data

    i32.const offset($ot_resource)
    i32.const length($ot_resource)
    call $otel_output_trace_data

;; Service name
    call $cache_service_name

    i32.const offset($service_name)
    global.get $service_name_len
    call $otel_output_trace_data

    i32.const offset($ot_resource_end)
    i32.const length($ot_resource_end)
    call $otel_output_trace_data

    i32.const offset($ot_start_scope_spans)
    i32.const length($ot_start_scope_spans)
    call $otel_output_trace_data

    ;; Print out the start time
    i32.const offset($ot_start_time)
    i32.const length($ot_start_time)
    call $otel_output_trace_data
    global.get $debug_timestamps_stack_pointer
    i32.const offset($debug_timestamps_stack)
    i32.add
    i64.load

    call $wt_format_i64_dec

    i32.const offset($db_number_i64)
    i32.const length($db_number_i64)
    call $otel_output_trace_data

    ;; Print out the end time
    i32.const offset($ot_comma)
    i32.const length($ot_comma)
    call $otel_output_trace_data
    i32.const offset($ot_end_time)
    i32.const length($ot_end_time)
    call $otel_output_trace_data
    local.get $time_end
    call $wt_format_i64_dec

    i32.const offset($db_number_i64)
    i32.const length($db_number_i64)
    call $otel_output_trace_data

    ;; Print out the name
    i32.const offset($ot_comma)
    i32.const length($ot_comma)
    call $otel_output_trace_data
    i32.const offset($ot_name)
    i32.const length($ot_name)
    call $otel_output_trace_data
    i32.const offset($ot_speech)
    i32.const length($ot_speech)
    call $otel_output_trace_data    
    local.get 0
    call $otel_output_trace_data_function_name    
    i32.const offset($ot_speech)
    i32.const length($ot_speech)
    call $otel_output_trace_data

    ;; Print out the trace_id
    i32.const offset($ot_comma)
    i32.const length($ot_comma)
    call $otel_output_trace_data

    i32.const offset($trace_id)
    call $get_invocation_id

    i32.const offset($trace_id)
    i32.const 16
    i32.const offset($ot_trace_id)
    i32.const 12
    i32.add
    call $wt_conv_hex

    i32.const offset($ot_trace_id)
    i32.const length($ot_trace_id)
    call $otel_output_trace_data

    ;; Print out the span_id
    i32.const offset($ot_comma)
    i32.const length($ot_comma)
    call $otel_output_trace_data

    global.get $debug_timestamps_stack_pointer
    i32.const offset($debug_timestamps_stack)
    i32.add
    i32.const 8
    i32.add
    i32.const 8
    i32.const offset($ot_span_id)
    i32.const 11
    i32.add
    call $wt_conv_hex

    i32.const offset($ot_span_id)
    i32.const length($ot_span_id)
    call $otel_output_trace_data

    ;; Get the parent if there is one...
    global.get $debug_timestamps_stack_pointer
    i32.const 0
    i32.ne
    if
      ;; If there is a parent span, print it out here...
      i32.const offset($ot_comma)
      i32.const length($ot_comma)
      call $otel_output_trace_data

      global.get $debug_timestamps_stack_pointer
      i32.const offset($debug_timestamps_stack)
      i32.add
      i32.const 8
      i32.sub
      i32.const 8
      i32.const offset($ot_parent_span_id)
      i32.const 18
      i32.add
      call $wt_conv_hex

      i32.const offset($ot_parent_span_id)
      i32.const length($ot_parent_span_id)
      call $otel_output_trace_data
    end

    i32.const offset($ot_comma)
    i32.const length($ot_comma)
    call $otel_output_trace_data

    i32.const offset($ot_attributes_start)
    i32.const length($ot_attributes_start)
    call $otel_output_trace_data

    ;; Output a dummy attribute for now...
    i32.const offset($ot_at_type)
    i32.const length($ot_at_type)
    i32.const offset($ot_at_type_fun)
    i32.const length($ot_at_type_fun)
    call $otel_output_attr_string

    ;; Output the function signature...
    i32.const offset($ot_comma)
    i32.const length($ot_comma)
    call $otel_output_trace_data

    i32.const offset($ot_attr_start)
    i32.const length($ot_attr_start)
    call $otel_output_trace_data

    i32.const offset($ot_at_fn_sig)
    i32.const length($ot_at_fn_sig)
    call $otel_output_trace_data

    i32.const offset($ot_attr_mid)
    i32.const length($ot_attr_mid)
    call $otel_output_trace_data

    i32.const offset($ot_attr_string_start)
    i32.const length($ot_attr_string_start)
    call $otel_output_trace_data

    local.get $fid
    call $otel_output_trace_data_function_signature

    i32.const offset($ot_attr_string_end)
    i32.const length($ot_attr_string_end)
    call $otel_output_trace_data

    i32.const offset($ot_attr_end)
    i32.const length($ot_attr_end)
    call $otel_output_trace_data

    ;; Output the function source...
    i32.const offset($ot_comma)
    i32.const length($ot_comma)
    call $otel_output_trace_data

    i32.const offset($ot_attr_start)
    i32.const length($ot_attr_start)
    call $otel_output_trace_data

    i32.const offset($ot_at_fn_src)
    i32.const length($ot_at_fn_src)
    call $otel_output_trace_data

    i32.const offset($ot_attr_mid)
    i32.const length($ot_attr_mid)
    call $otel_output_trace_data

    i32.const offset($ot_attr_string_start)
    i32.const length($ot_attr_string_start)
    call $otel_output_trace_data

    local.get $fid
    call $otel_output_trace_data_function_source

    i32.const offset($ot_attr_string_end)
    i32.const length($ot_attr_string_end)
    call $otel_output_trace_data

    i32.const offset($ot_attr_end)
    i32.const length($ot_attr_end)
    call $otel_output_trace_data

  )

  ;; wt_print_function_name - Given a function ID, print out the function name.
  (func $otel_output_trace_data_function_name (param $fid i32)
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
    call $otel_output_trace_data
  )

  ;; wt_print_function_signature - Given a function ID, print out the function signature.
  (func $otel_output_trace_data_function_signature (param $fid i32)
    (local $ptr i32)
    i32.const offset($wt_all_function_sigs_locs)
    local.get $fid
    i32.const 3
    i32.shl
    i32.add
    local.tee $ptr

    ;; Get the address
    i32.load

    ;; We need to adjust this address manually
    i32.const offset($wt_all_function_sigs)
    i32.add

    ;; Get the offset
    local.get $ptr
    i32.load offset=4
    call $otel_output_trace_data
  )

  ;; wt_print_function_src - Given a function ID, print out the function src.
  (func $otel_output_trace_data_function_source (param $fid i32)
    (local $ptr i32)
    i32.const offset($wt_all_function_srcs_locs)
    local.get $fid
    i32.const 3
    i32.shl
    i32.add
    local.tee $ptr

    ;; Get the address
    i32.load

    ;; We need to adjust this address manually
    i32.const offset($wt_all_function_srcs)
    i32.add

    ;; Get the offset
    local.get $ptr
    i32.load offset=4
    call $otel_output_trace_data
  )

;; Prelude
  (data $ot_start "{\22resource_spans\22:[{")
  (data $ot_resource "\22resource\22:{\22attributes\22:[{\22key\22:\22service.name\22,\22value\22:{\22stringValue\22:\22")
  (data $ot_resource_end "\22}}]},")

  (data $ot_start_scope_spans "\22scope_spans\22:[{\22spans\22:[{")

  (data $ot_end "}]}]}]}")

;; "trace_id": "<DATA>",
  (data $ot_trace_id "\22trace_id\22:\22--------------------------------\22")
;; "span_id": "<DATA>",
  (data $ot_span_id "\22span_id\22:\22----------------\22")
;; "parent_span_id": "<DATA>",
  (data $ot_parent_span_id "\22parent_span_id\22:\22----------------\22")
;; "name": "<DATA>"
  (data $ot_name "\22name\22:")
;; "start_time_unix_nano"
  (data $ot_start_time "\22start_time_unix_nano\22:")
;; "end_time_unix_nano"
  (data $ot_end_time "\22end_time_unix_nano\22:")
;; "attributes"
  (data $ot_attributes_start "\22attributes\22:[")

  (data $ot_attr_start "{\22key\22:\22")
  (data $ot_attr_mid "\22,\22value\22:{")
  (data $ot_attr_string_start "\22stringValue\22:\22")
  (data $ot_attr_string_end "\22")
  (data $ot_attr_end "}}")

  (data $ot_attributes_end "]")

  (data $ot_at_type "type")
  (data $ot_at_type_fun "wasm function")
  (data $ot_at_type_quickjs_fun "quickjs function")

  (data $ot_at_fn_sig "function")
  (data $ot_at_fn_src "source")

  (data $ot_at_result "result")
  (data $ot_at_param "param_000")

  (data $ot_at_todo "TODO")

  (data $ot_comma ",")
  (data $ot_space " ")
  (data $ot_equals "=")
  (data $ot_nonprintable ".")

  (data $ot_speech "\22")

  (data $ot_encoded_speech "\5C\22")

  (data $_printable_characters "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ ")

  (data $trace_id 16)

  (data $error_stack_overflow "Error: The timings stack overflowed. You win some you lose some I guess.\0d\0a")
  (data $error_buffer_overflow "Error: The output buffer overflowed. You win some you lose some I guess.\0d\0a")
  (data $debug_newline "\0a")

  (data $debug_clock_loc 8)

  (data $otel_output_buffer 4096)

  (global $otel_output_buffer_ptr (mut i32) (i32.const 0))

  ;; Timestamp stack (8 bytes per entry)
  (data $debug_timestamps_stack 80000)

  (global $debug_timestamps_stack_pointer (mut i32) (i32.const 0))

  (global $wt_all_function_length i32 (i32.const 0))

)