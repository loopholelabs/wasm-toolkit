(module


  (func $otel_watch_global_go_string (param $name_ptr i32) (param $name_len i32) (param $ptr i32) (param $len i32)
    i32.const offset($ot_comma)
    i32.const length($ot_comma)
    call $otel_output_trace_data

    local.get $name_ptr
    local.get $name_len

    ;; Get the data pointer
    local.get $ptr
    i32.load

    ;; Get the length
    local.get $ptr
    i32.load offset=4

    call $otel_output_attr_hexdata
  )

  (func $otel_watch_global_go_byte_slice (param $name_ptr i32) (param $name_len i32) (param $ptr i32) (param $len i32)
    i32.const offset($ot_comma)
    i32.const length($ot_comma)
    call $otel_output_trace_data


    local.get $name_ptr
    local.get $name_len

    ;; Get the data pointer
    local.get $ptr
    i32.load

    ;; Get the length
    local.get $ptr
    i32.load offset=4

    call $otel_output_attr_hexdata
  )

)