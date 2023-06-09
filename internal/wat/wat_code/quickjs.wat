(module


  ;; Show a quickjs arg
  (func $otel_quickjs_arg (param $context i32) (param $i i32) (param $v i64)
    local.get $i
    i32.const offset($js_param)
    i32.const 6
    i32.add
    call $wt_conv_byte_dec

    local.get $context
    i32.const offset($js_param)
    i32.const length($js_param)
    local.get $v
    call $otel_quickjs_prop
  )

  ;; Show a quickjs value
  (func $otel_quickjs_prop (param $context i32) (param $key_ptr i32) (param $key_len i32) (param $v i64)
    (local $tag i32)
    (local $val i32)

    (local $exception i64)

    i32.const offset($ot_comma)
    i32.const length($ot_comma)
    call $otel_output_trace_data

    ;; JS_JSONStringify

    ;; Convert it to a string
;;      i32.const 0
;;      global.set $trace_enable
;;    local.get $context
;;    local.get $v
;;    call $JS_ToQuotedString

;;      i32.const 1
;;      global.set $trace_enable
    ;;local.set $str_val

    ;; For testing...
    ;; local.set $v

    ;; Split it into tag / value
    local.get $v
    i64.const 32
    i64.shr_u
    i32.wrap_i64
    local.set $tag

    local.get $v
    i32.wrap_i64
    local.set $val

    ;; 0 TAG_INT
    local.get $tag
    i32.const 0
    i32.eq
    if
      local.get $val
      call $wt_format_i32_hex

      local.get $key_ptr
      local.get $key_len
      i32.const offset($js_type_int)
      i32.const length($js_type_int)
      i32.const offset($db_number_i32)
      i32.const 8
      call $otel_output_attr_var_string
      return
    end

    ;; 1 TAG_BOOL
    local.get $tag
    i32.const 1
    i32.eq
    if
    end

    ;; 2 TAG_NULL
    local.get $tag
    i32.const 2
    i32.eq
    if
    end

    ;; 3 TAG_UNDEFINED
    local.get $tag
    i32.const 3
    i32.eq
    if
      local.get $key_ptr
      local.get $key_len
      i32.const offset($js_type_undefined)
      i32.const length($js_type_undefined)
      call $otel_output_attr_string
      return
    end

    ;; 4 TAG_UNINITIALIZED
    local.get $tag
    i32.const 4
    i32.eq
    if
    end

    ;; 5 TAG_CATCH_OFFSET
    local.get $tag
    i32.const 5
    i32.eq
    if
    end

    ;; 6 TAG_EXCEPTION
    local.get $tag
    i32.const 6
    i32.eq
    if
      global.get $trace_enable
      i32.eqz
      if
        unreachable
      end

;;      i32.const 0
;;      global.set $trace_enable

      local.get $context
      call $JS_GetException
      drop
;;      i32.const 1
;;      global.set $trace_enable

      local.get $v
      ;; TODO Get useful info out of the object, or toString it...

      call $wt_format_i64_hex

      local.get $key_ptr
      local.get $key_len
      i32.const offset($js_type_exception)
      i32.const length($js_type_exception)
      i32.const offset($db_number_i64)
      i32.const 16
      call $otel_output_attr_var_string
      return
    end

    ;; 7 TAG_FLOAT64
    local.get $tag
    i32.const 7
    i32.eq
    if
    end

    ;; -1 TAG_OBJECT
    local.get $tag
    i32.const 0xffffffff
    i32.eq
    if
    end

    ;; -7 TAG_STRING
    local.get $tag
    i32.const 0xfffffff9
    i32.eq
    if  
      local.get $key_ptr
      local.get $key_len
      i32.const offset($js_type_string)
      i32.const length($js_type_string)

      ;; Internal QuickJS. MAY CHANGE      
      local.get $val
      i32.const 16
      i32.add
      local.get $val
      i32.load offset=4
      call $otel_output_attr_var_string
      return
    end

    ;; -8 TAG_SYMBOL
    local.get $tag
    i32.const 0xfffffff8
    i32.eq
    if
    end

    ;; -9 TAG_BIG_FLOAT
    local.get $tag
    i32.const 0xfffffff7
    i32.eq
    if
    end

    ;; -10 TAG_BIG_INT
    local.get $tag
    i32.const 0xfffffff6
    i32.eq
    if
    end

    ;; -11 TAG_BIG_DECIMAL
    local.get $tag
    i32.const 0xfffffff5
    i32.eq
    if
    end

    ;; Fallback generic output
    local.get $v
    call $wt_format_i64_hex

    local.get $key_ptr
    local.get $key_len
    i32.const offset($db_number_i64)
    i32.const 16
    call $otel_output_attr_string

  )

  (func $otel_quickjs_call (param $ret i64) (param $context i32) (param $v_func i64) (param $v_this i64) (param $argc i32) (param $argv i32) (result i64)
    (local $par i32)
    (local $fn_string i64)

    ;; First output the javascript arguments...
    block
      loop
        local.get $par
        local.get $argc
        i32.ge_u
        br_if 1

        local.get $context
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

    ;; Now deal with the result...

    local.get $context
    i32.const offset($js_result)
    i32.const length($js_result)
    local.get $ret
    call $otel_quickjs_prop

    local.get $ret
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

      global.get $trace_enable
      i32.eqz
      if
        unreachable
      end

      i32.const 0
      global.set $trace_enable


    local.get $context
    local.get $v_func
    i32.const offset($atom_name)
    call $JS_GetPropertyStr
    local.set $fn_name

      i32.const 1
      global.set $trace_enable


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
    i32.const offset($ot_at_type_quickjs_fun)
    i32.const length($ot_at_type_quickjs_fun)
    call $otel_output_attr_string

  )

  (data $js_param "param_000")
  (data $js_result "result")

  (data $js_type_int "(int)")
  (data $js_type_string "(string)")
  (data $js_type_exception "(exception)")
  (data $js_type_undefined "(undefined)")

  (data $atom_name "name\00")
)