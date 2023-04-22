(module
  (type (func (param i32 i64 i32) (result i32)))
  (type (func (param i32 i32) (result i32)))
  (import "wasi_snapshot_preview1" "clock_time_get" (func $debug_clock_time_get (type 0)))
  (import "wasi_snapshot_preview1" "random_get" (func $debug_random_get (type 1)))

  (func $debug_gettime (result i64)
    i32.const 0
    i64.const 1000
    i32.const offset($debug_clock_loc)
    call $debug_clock_time_get
    drop
    i32.const offset($debug_clock_loc)
    i64.load
  )

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
    i32.const 800
    i32.ge_u
    ;; Detect stack overflow
    if
      i32.const offset($error_stack_overflow)
      i32.const length($error_stack_overflow)
      call $wt_print
      unreachable
    end
  )

  (func $otel_exit_func (param $fid i32)
    (local $time_end i64)

    call $debug_gettime
    local.set $time_end

    ;; Set the trace_id if it hasn't already been set
    ;; TODO: Clear this

    global.get $trace_id_set
    i32.eqz
    if
      i32.const 1
      global.set $trace_id_set

      i32.const offset($trace_id)
      i32.const length($trace_id)
      call $debug_random_get
      drop

    end

    ;; Pop timestamp off the timestamp stack
    global.get $debug_timestamps_stack_pointer
    i32.const 16
    i32.sub
    global.set $debug_timestamps_stack_pointer


    i32.const offset($ot_start)
    i32.const length($ot_start)
    call $wt_print

    i32.const offset($ot_resource)
    i32.const length($ot_resource)
    call $wt_print

    i32.const offset($ot_resource_name)
    i32.const length($ot_resource_name)
    call $wt_print

    i32.const offset($ot_resource_end)
    i32.const length($ot_resource_end)
    call $wt_print

    i32.const offset($ot_start_scope_spans)
    i32.const length($ot_start_scope_spans)
    call $wt_print

    ;; Print out the start time
;;    i32.const offset($ot_comma)
;;    i32.const length($ot_comma)
;;    call $wt_print
    i32.const offset($ot_start_time)
    i32.const length($ot_start_time)
    call $wt_print
    global.get $debug_timestamps_stack_pointer
    i32.const offset($debug_timestamps_stack)
    i32.add
    i64.load

    call $wt_format_i64_dec

    i32.const offset($db_number_i64)
    i32.const length($db_number_i64)
    call $wt_print

    ;; Print out the end time
    i32.const offset($ot_comma)
    i32.const length($ot_comma)
    call $wt_print
    i32.const offset($ot_end_time)
    i32.const length($ot_end_time)
    call $wt_print
    local.get $time_end
    call $wt_format_i64_dec

    i32.const offset($db_number_i64)
    i32.const length($db_number_i64)
    call $wt_print

    ;; Print out the name
    i32.const offset($ot_comma)
    i32.const length($ot_comma)
    call $wt_print
    i32.const offset($ot_name)
    i32.const length($ot_name)
    call $wt_print
    i32.const offset($ot_speech)
    i32.const length($ot_speech)
    call $wt_print
    local.get 0
    call $wt_print_function_name    
    i32.const offset($ot_speech)
    i32.const length($ot_speech)
    call $wt_print

    ;; Print out the trace_id
    i32.const offset($ot_comma)
    i32.const length($ot_comma)
    call $wt_print

    i32.const offset($trace_id)
    i32.const 16
    i32.const offset($ot_trace_id)
    i32.const 12
    i32.add
    call $wt_conv_hex

    i32.const offset($ot_trace_id)
    i32.const length($ot_trace_id)
    call $wt_print

    ;; Print out the span_id
    i32.const offset($ot_comma)
    i32.const length($ot_comma)
    call $wt_print

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
    call $wt_print

    ;; Get the parent if there is one...
    global.get $debug_timestamps_stack_pointer
    i32.const 0
    i32.ne
    if
      ;; If there is a parent span, print it out here...
      i32.const offset($ot_comma)
      i32.const length($ot_comma)
      call $wt_print

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
      call $wt_print

    end

    i32.const offset($ot_end)
    i32.const length($ot_end)
    call $wt_print

    i32.const offset($debug_newline)
    i32.const length($debug_newline)
    call $wt_print
  )

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

  (data $ot_resource_name "loop-wasm-otel")

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

  (data $ot_comma ",")
  (data $ot_speech "\22")

  (data $trace_id 16)

  (data $error_stack_overflow "Error: The timings stack overflowed. You win some you lose some I guess.\0d\0a")
  (data $debug_newline "\0a")

  (data $debug_clock_loc 8)

  ;; Only allow 100 function stack for now
  (data $debug_timestamps_stack 800)

  (global $debug_timestamps_stack_pointer (mut i32) (i32.const 0))

  (global $wt_all_function_length i32 (i32.const 0))

  (global $trace_id_set (mut i32) (i32.const 0))
)