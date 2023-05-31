(module
  
	;; $log_global_<TYPE> (new_value, current_value, ptr_debug, len_debug) => new_value
  (func $log_global_i32 (param $new i32) (param $current i32) (param $ptr i32) (param $len i32) (result i32)
    global.get $wt_color
    if
      i32.const offset($wt_ansi_watch)
      i32.const length($wt_ansi_watch)
      call $wt_print
    end

    i32.const offset($log_watch_global_1)
    i32.const length($log_watch_global_1)
    call $wt_print
    local.get $ptr
    local.get $len
    call $wt_print

    i32.const offset($log_watch_global_2)
    i32.const length($log_watch_global_2)
    call $wt_print

    local.get $current
    call $wt_format_i32_hex

    i32.const offset($db_number_i32)
    i32.const 8
    call $wt_print

    i32.const offset($log_watch_global_3)
    i32.const length($log_watch_global_3)
    call $wt_print

    local.get $new
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

    local.get $new
  )

  (func $log_global_i64 (param $new i64) (param $current i64) (param $ptr i32) (param $len i32) (result i64)
    local.get $new
  )

  (func $log_global_f32 (param $new f32) (param $current f32) (param $ptr i32) (param $len i32) (result f32)
    local.get $new
  )

  (func $log_global_f64 (param $new f64) (param $current f64) (param $ptr i32) (param $len i32) (result f64)
    local.get $new
  )

  ;;
  (func $log_local_i32 (param $new i32) (param $current i32) (param $ptr i32) (param $len i32) (result i32)
    global.get $wt_color
    if
      i32.const offset($wt_ansi_watch)
      i32.const length($wt_ansi_watch)
      call $wt_print
    end

    i32.const offset($log_watch_local_1)
    i32.const length($log_watch_local_1)
    call $wt_print
    local.get $ptr
    local.get $len
    call $wt_print

    i32.const offset($log_watch_local_2)
    i32.const length($log_watch_local_2)
    call $wt_print

    local.get $current
    call $wt_format_i32_hex

    i32.const offset($db_number_i32)
    i32.const 8
    call $wt_print

    i32.const offset($log_watch_local_3)
    i32.const length($log_watch_local_3)
    call $wt_print

    local.get $new
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

    local.get $new
  )

  (func $log_local_i64 (param $new i64) (param $current i64) (param $ptr i32) (param $len i32) (result i64)
    local.get $new
  )

  (func $log_local_f32 (param $new f32) (param $current f32) (param $ptr i32) (param $len i32) (result f32)
    local.get $new
  )

  (func $log_local_f64 (param $new f64) (param $current f64) (param $ptr i32) (param $len i32) (result f64)
    local.get $new
  )

  (data $log_watch_global_1 "GLOBAL ")
  (data $log_watch_global_2 " value ")
  (data $log_watch_global_3 "=>")

  (data $log_watch_local_1 "LOCAL ")
  (data $log_watch_local_2 " value ")
  (data $log_watch_local_3 "=>")

)