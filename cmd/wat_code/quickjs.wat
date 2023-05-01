(module
  (type (func))
  (import "env" "host_snapshot" (func $host_snapshot (type 0)))


  ;; Show a quickjs arg
  (func $otel_quickjs_arg (param $i i32) (param $v i64)
    (local $tag i32)
    (local $val i32)

    i32.const offset($ot_comma)
    i32.const length($ot_comma)
    call $otel_output_trace_data

    ;; Split it into tag / value
    local.get $v
    i64.const 32
    i64.shr_u
    i32.wrap_i64
    local.set $tag

    local.get $v
    i32.wrap_i64
    local.set $val

    ;; If it's a simple INT argument...
    local.get $tag
    i32.const 0
    i32.eq
    if
      local.get $i
      i32.const offset($ot_at_qjs_param_i)
      i32.const 12
      i32.add
      call $wt_conv_byte_dec
  
      local.get $val
      call $wt_format_i32_hex

      i32.const offset($ot_at_qjs_param_i)
      i32.const length($ot_at_qjs_param_i)
      i32.const offset($db_number_i32)
      i32.const 8
      call $otel_output_attr_string
      return
    end

    ;; If it's a simple string argument...
    local.get $tag
    i32.const 0xfffffff9
    i32.eq
    if
      local.get $i
      i32.const offset($ot_at_qjs_param_s)
      i32.const 12
      i32.add
      call $wt_conv_byte_dec
  
      ;; Deal with the val
      i32.const offset($ot_at_qjs_param_s)
      i32.const length($ot_at_qjs_param_s)

      ;; Internal QuickJS. MAY CHANGE      
      local.get $val
      i32.const 16
      i32.add

      local.get $val
      i32.load offset=4

      call $otel_output_attr_string
      return
    end

    local.get $i
    i32.const offset($ot_at_qjs_param_o)
    i32.const 12
    i32.add
    call $wt_conv_byte_dec

    local.get $v
    call $wt_format_i64_hex

    i32.const offset($ot_at_qjs_param_o)
    i32.const length($ot_at_qjs_param_o)
    i32.const offset($db_number_i64)
    i32.const 16
    call $otel_output_attr_string

  )

  (func $otel_quickjs_call (param $context i32) (param $v_func i64) (param $v_this i64) (param $argc i32) (param $argv i32)
    (local $par i32)
    (local $fn_string i64)

    local.get $v_func
    i64.const 0xffffffff0013aed0
    i64.eq
    if
      call $host_snapshot
    end

    block
      loop
        local.get $par
        local.get $argc
        i32.ge_u
        br_if 1

        ;; Get the js param and deal with it...
        local.get $par

        local.get $par
        i32.const 3
        i32.shl
        local.get $argv
        i32.add
        i64.load

        call $otel_quickjs_arg

        local.get $par
        i32.const 1
        i32.add
        local.set $par
        br 0
      end
    end


    i32.const offset($ot_comma)
    i32.const length($ot_comma)
    call $otel_output_trace_data

    local.get $argc
    call $wt_format_i32_hex

    i32.const offset($ot_at_qjs_argc)
    i32.const length($ot_at_qjs_argc)
    i32.const offset($db_number_i32)
    i32.const 8
    call $otel_output_attr_string

    ;; Now onto func stuff...

    local.get $v_func
    i32.wrap_i64
    i32.const 64
    i32.const offset($ot_at_qjs_data)
    call $wt_conv_hex

    i32.const offset($ot_comma)
    i32.const length($ot_comma)
    call $otel_output_trace_data

    i32.const offset($ot_at_qjs_func)
    i32.const length($ot_at_qjs_func)
    i32.const offset($ot_at_qjs_data)
    i32.const length($ot_at_qjs_data)
    call $otel_output_attr_string

    local.get $v_this
    i32.wrap_i64
    i32.const 64
    i32.const offset($ot_at_qjs_data)
    call $wt_conv_hex

    i32.const offset($ot_comma)
    i32.const length($ot_comma)
    call $otel_output_trace_data

    i32.const offset($ot_at_qjs_this)
    i32.const length($ot_at_qjs_this)
    i32.const offset($ot_at_qjs_data)
    i32.const length($ot_at_qjs_data)
    call $otel_output_attr_string

    ;; Get the 'name' property of the function...

    local.get $context
    local.get $v_func
    i32.const offset($atom_name)
    call $JS_GetPropertyStr

    ;;local.get $context
    ;;local.get $v_func
    ;;local.get $argc
    ;;local.get $argv
    ;;call $js_function_toString
    local.set $fn_string
    ;; i64 result...

    i32.const offset($ot_comma)
    i32.const length($ot_comma)
    call $otel_output_trace_data

    local.get $fn_string
    call $wt_format_i64_hex

    i32.const offset($ot_at_qjs_test)
    i32.const length($ot_at_qjs_test)
    i32.const offset($db_number_i64)
    i32.const 16
    call $otel_output_attr_string

    ;; Show it as string

    i32.const offset($ot_comma)
    i32.const length($ot_comma)
    call $otel_output_trace_data

    ;; Deal with the val
    i32.const offset($ot_at_qjs_test_s)
    i32.const length($ot_at_qjs_test_s)

    ;; Internal QuickJS. MAY CHANGE      
    local.get $fn_string
    i32.wrap_i64
    i32.const 16
    i32.add

    local.get $fn_string
    i32.wrap_i64
    i32.load offset=4

    call $otel_output_attr_string


  )

(func $qjs_otel_exit_func (param $context i32) (param $v_func i64)
  (local $time_end i64)
  (local $fn_name i64)

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

    i32.const offset($ot_resource_name)
    i32.const length($ot_resource_name)
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

    local.get $context
    local.get $v_func
    i32.const offset($atom_name)
    call $JS_GetPropertyStr
    local.set $fn_name

    ;; Internal QuickJS. MAY CHANGE      
    local.get $fn_name
    i32.wrap_i64
    i32.const 16
    i32.add

    local.get $fn_name
    i32.wrap_i64
    i32.load offset=4
    call $otel_output_trace_data

    i32.const offset($ot_speech)
    i32.const length($ot_speech)
    call $otel_output_trace_data

    ;; Print out the trace_id
    i32.const offset($ot_comma)
    i32.const length($ot_comma)
    call $otel_output_trace_data

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

  )

  (data $ot_at_qjs_func "quickjs_func")
  (data $ot_at_qjs_this "quickjs_this")
  (data $ot_at_qjs_func20 "quickjs_func20")
  (data $ot_at_qjs_data "00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")

  (data $ot_at_qjs_argc "quickjs_argc")
  (data $ot_at_qjs_argv "quickjs_argv")

  (data $ot_at_qjs_param_i "qjs_param_i_000")
  (data $ot_at_qjs_param_o "qjs_param_o_000")
  (data $ot_at_qjs_param_s "qjs_param_s_000")
  (data $ot_at_qjs_param_s_v "qjs_param_s_000_val")

  (data $ot_at_qjs_test "qjs_test")
  (data $ot_at_qjs_test_s "qjs_test_s")

  (data $atom_name "name\00")


)