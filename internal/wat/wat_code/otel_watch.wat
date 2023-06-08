(module

  (func $log_mem_common_end
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

  (func $log_mem_common_start (param $debug_ptr i32) (param $debug_len i32)
    (local $time_end i64)

    ;; Output a trace for this watch event...
    call $debug_gettime
    local.set $time_end

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

    local.get $time_end
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
    local.get $debug_ptr
    local.get $debug_len
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

    ;; Print out the span_id (random)
    i32.const offset($ot_comma)
    i32.const length($ot_comma)
    call $otel_output_trace_data

    i32.const offset($ot_watch_span_id)
    i32.const 8
    call $debug_random_get
    drop

    i32.const offset($ot_watch_span_id)
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
    i32.const offset($ot_at_type_mem)
    i32.const length($ot_at_type_mem)
    call $otel_output_attr_string  
  )

  (func $log_mem_i32.store (param $address i32) (param $offset i32) (param $size i32) (param $debug_ptr i32) (param $debug_len i32) (result i32)
    (local $memrange i32)
    (local $time_end i32)
    local.get $address
    local.get $offset
    i32.add
    local.get $address
    local.get $offset
    i32.add
    local.get $size
    i32.const 3
    i32.shr_u
    i32.add
    call $check_dynamic_watches
    local.tee $memrange
    i32.eqz
    if
      local.get $address
      return
    end

    local.get $debug_ptr
    local.get $debug_len
    call $log_mem_common_start

    ;; Our attributes here...

    i32.const offset($ot_comma)
    i32.const length($ot_comma)
    call $otel_output_trace_data

    local.get $memrange
    i32.load offset=8
    call $wt_format_i32_hex
    i32.const offset($watch_mem_out_id)
    i32.const length($watch_mem_out_id)
    i32.const offset($db_number_i32)
    i32.const 8
    call $otel_output_attr_string


    i32.const offset($ot_comma)
    i32.const length($ot_comma)
    call $otel_output_trace_data

    local.get $address
    call $wt_format_i32_hex

    i32.const offset($watch_mem_out_address)
    i32.const length($watch_mem_out_address)
    i32.const offset($db_number_i32)
    i32.const 8
    call $otel_output_attr_string

    i32.const offset($ot_comma)
    i32.const length($ot_comma)
    call $otel_output_trace_data

    local.get $offset
    call $wt_format_i32_hex

    i32.const offset($watch_mem_out_offset)
    i32.const length($watch_mem_out_offset)
    i32.const offset($db_number_i32)
    i32.const 8
    call $otel_output_attr_string

    i32.const offset($ot_comma)
    i32.const length($ot_comma)
    call $otel_output_trace_data

    local.get $address
    local.get $offset
    i32.add
    i32.load
    ;; Current value...
    call $wt_format_i32_hex

    local.get $size
    i32.const 32
    i32.eq
    if
      i32.const offset($watch_mem_out_oldvalue)
      i32.const length($watch_mem_out_oldvalue)
      i32.const offset($db_number_i32)
      i32.const 8
      call $otel_output_attr_string
    end

    local.get $size
    i32.const 16
    i32.eq
    if
      i32.const offset($watch_mem_out_oldvalue)
      i32.const length($watch_mem_out_oldvalue)
      i32.const offset($db_number_i32)
      i32.const 4
      i32.add
      i32.const 4
      call $otel_output_attr_string
    end

    local.get $size
    i32.const 8
    i32.eq
    if
      i32.const offset($watch_mem_out_oldvalue)
      i32.const length($watch_mem_out_oldvalue)
      i32.const offset($db_number_i32)
      i32.const 6
      i32.add
      i32.const 2
      call $otel_output_attr_string
    end  

    i32.const offset($ot_comma)
    i32.const length($ot_comma)
    call $otel_output_trace_data
    
    global.get $log_memory_value_i32
    call $wt_format_i32_hex

    local.get $size
    i32.const 32
    i32.eq
    if
      i32.const offset($watch_mem_out_newvalue)
      i32.const length($watch_mem_out_newvalue)
      i32.const offset($db_number_i32)
      i32.const 8
      call $otel_output_attr_string
    end

    local.get $size
    i32.const 16
    i32.eq
    if
      i32.const offset($watch_mem_out_newvalue)
      i32.const length($watch_mem_out_newvalue)
      i32.const offset($db_number_i32)
      i32.const 4
      i32.add
      i32.const 4
      call $otel_output_attr_string
    end

    local.get $size
    i32.const 8
    i32.eq
    if
      i32.const offset($watch_mem_out_newvalue)
      i32.const length($watch_mem_out_newvalue)
      i32.const offset($db_number_i32)
      i32.const 6
      i32.add
      i32.const 2
      call $otel_output_attr_string
    end  


    call $log_mem_common_end

    local.get $address
  )

  (func $log_mem_i64.store (param $address i32) (param $offset i32) (param $size i32) (param $debug_ptr i32) (param $debug_len i32) (result i32)
    (local $memrange i32)
    local.get $address
    local.get $offset
    i32.add
    local.get $address
    local.get $offset
    i32.add
    local.get $size
    i32.const 3
    i32.shr_u
    i32.add
    call $check_dynamic_watches
    local.tee $memrange
    i32.eqz
    if
      local.get $address
      return
    end

    ;; Output a trace...

    local.get $address
  )

  (func $check_dynamic_watches (param $start i32) (param $end i32) (result i32)
    (local $p i32)
    (local $cptr i32)

    block
      loop
        local.get $p
        global.get $wt_dynamic_mem_ranges_len
        i32.ge_u
        br_if 1

        ;; Now check the range pointed at by index $p
        local.get $p
        i32.const 4
        i32.shl
        i32.const offset($wt_dynamic_mem_ranges)
        i32.add
        local.set $cptr

        ;; Check the dynamic memory range here
        local.get $start
        local.get $cptr
        i32.load
        i32.ge_u
        if
          ;; $start is greater than the range start
          local.get $end
          local.get $cptr
          i32.load offset=4
          i32.le_u
          if
            ;; We are within the range
            local.get $cptr
            ;; Return ptr to the memory range info
            return
          end
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

  (func $watch_add (param $id i32) (param $addr i32) (param $len i32)
    (local $p i32)

  ;; Add a new entry in $wt_dynamic_mem_ranges
    global.get $wt_dynamic_mem_ranges_len
    i32.const 4
    i32.shl
    i32.const offset($wt_dynamic_mem_ranges)
    i32.add
    local.tee $p
    i32.const length($wt_dynamic_mem_ranges)
    i32.const offset($wt_dynamic_mem_ranges)
    i32.add

    i32.ge_u
    if
      ;; We ran out of watch space...
;;      i32.const offset($error_watch_overflow)
;;      i32.const length($error_watch_overflow)
;;      call $wt_print
      unreachable
    end
    ;; Store things here...
    local.get $p
    local.get $addr
    i32.store
    
    local.get $p
    local.get $addr
    local.get $len
    i32.add
    i32.store offset=4

    ;; Now store the tag. Len=0 means to use the ptr as an ID.
    local.get $p
    local.get $id
    i32.store offset=8

    local.get $p
    i32.const 0
    i32.store offset=12

    global.get $wt_dynamic_mem_ranges_len
    i32.const 1
    i32.add
    global.set $wt_dynamic_mem_ranges_len
  )

  (func $watch_del (param $id i32)
    ;; Remove an entry from $wt_dynamic_mem_ranges
    ;; TODO
  )

  (data $wt_dynamic_mem_ranges 1600)
  (global $wt_dynamic_mem_ranges_len (mut i32) (i32.const 0))

  (global $log_memory_value_i32 (mut i32) (i32.const 0))
  (global $log_memory_value_i64 (mut i64) (i64.const 0))

  (data $ot_watch_span_id 8)

  (data $watch_mem_out_id "watch_id")
  (data $watch_mem_out_address "address")
  (data $watch_mem_out_offset "offset")

  (data $watch_mem_out_newvalue "value_new")
  (data $watch_mem_out_oldvalue "value_old")

)