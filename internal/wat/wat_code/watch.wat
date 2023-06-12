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
    call $wt_format_i64_hex

    i32.const offset($db_number_i64)
    i32.const 16
    call $wt_print

    i32.const offset($log_watch_global_3)
    i32.const length($log_watch_global_3)
    call $wt_print

    local.get $new
    call $wt_format_i64_hex

    i32.const offset($db_number_i64)
    i32.const 16
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
    call $wt_format_i64_hex

    i32.const offset($db_number_i64)
    i32.const 16
    call $wt_print

    i32.const offset($log_watch_local_3)
    i32.const length($log_watch_local_3)
    call $wt_print

    local.get $new
    call $wt_format_i64_hex

    i32.const offset($db_number_i64)
    i32.const 16
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

  (func $log_local_f32 (param $new f32) (param $current f32) (param $ptr i32) (param $len i32) (result f32)
    local.get $new
  )

  (func $log_local_f64 (param $new f64) (param $current f64) (param $ptr i32) (param $len i32) (result f64)
    local.get $new
  )

;; MEMORY ACCESS
  (func $log_mem_filter (param $start i32) (param $end i32) (result i32)
    (local $c i32)

    ;; First check the dynamic watch ranges...
    local.get $start
    local.get $end
    call $check_dynamic_watches
    local.tee $c
    if
      local.get $c
      return
    end

    block
      loop
        local.get $c
        i32.const length($wt_mem_ranges)
        i32.ge_u
        br_if 1

        ;; Check if we're inside the memory range...
        local.get $start
        local.get $c
        i32.const offset($wt_mem_ranges)
        i32.add
        i32.load
        i32.ge_u
        if

          ;; $start is greater than the range start
          local.get $end
          local.get $c
          i32.const offset($wt_mem_ranges)
          i32.add
          i32.load offset=4
          i32.le_u
          if
            ;; We are within the range
            local.get $c
            i32.const offset($wt_mem_ranges)
            i32.add
            ;; Return ptr to the memory range info
            return
          end
        end

        local.get $c
        i32.const 16
        i32.add
        local.set $c
        br 0
      end
    end

    ;; Default to NOTHING
    i32.const 0
  )

  (func $log_print_i32_size (param $size i32)
    local.get $size
    i32.const 32
    i32.eq
    if
      i32.const offset($db_number_i32)
      i32.const 8
      call $wt_print
    end

    local.get $size
    i32.const 16
    i32.eq
    if
      i32.const offset($db_number_i32)
      i32.const 4
      i32.add
      i32.const 4
      call $wt_print
    end

    local.get $size
    i32.const 8
    i32.eq
    if
      i32.const offset($db_number_i32)
      i32.const 6
      i32.add
      i32.const 2
      call $wt_print
    end  
  )

  (func $log_print_i64_size (param $size i32)
    local.get $size
    i32.const 64
    i32.eq
    if
      i32.const offset($db_number_i64)
      i32.const 16
      call $wt_print
    end

    local.get $size
    i32.const 32
    i32.eq
    if
      i32.const offset($db_number_i64)
      i32.const 8
      i32.add
      i32.const 8
      call $wt_print
    end

    local.get $size
    i32.const 16
    i32.eq
    if
      i32.const offset($db_number_i64)
      i32.const 12
      i32.add
      i32.const 4
      call $wt_print
    end

    local.get $size
    i32.const 8
    i32.eq
    if
      i32.const offset($db_number_i64)
      i32.const 14
      i32.add
      i32.const 2
      call $wt_print
    end  
  )

  (func $log_mem_i32.store (param $address i32) (param $offset i32) (param $size i32) (param $debug_ptr i32) (param $debug_len i32) (result i32)
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
    call $log_mem_filter
    local.tee $memrange
    i32.eqz
    if
      local.get $address
      return
    end

    global.get $wt_color
    if
      i32.const offset($wt_ansi_watch)
      i32.const length($wt_ansi_watch)
      call $wt_print
    end

    i32.const offset($log_watch_memory_0)
    i32.const length($log_watch_memory_0)
    call $wt_print

    local.get $memrange
    if
      local.get $memrange
      i32.load offset=12
      if
        local.get $memrange
        i32.load offset=8
        i32.const offset($wt_mem_tags)
        i32.add
        local.get $memrange
        i32.load offset=12
        call $wt_print
      else
        ;; It's an ID instead...
        local.get $memrange
        i32.load offset=8
        i32.const offset($wt_id_tag)
        i32.const 4
        i32.add
        call $wt_conv_byte_dec
        i32.const offset($wt_id_tag)
        i32.const length($wt_id_tag)
        call $wt_print
      end
    end

    ;;i32.const offset($log_watch_memory_1)
    ;;i32.const length($log_watch_memory_1)
    local.get $debug_ptr
    local.get $debug_len
    call $wt_print

    i32.const offset($log_watch_memory_1)
    i32.const length($log_watch_memory_1)
    call $wt_print

    local.get $address
    call $wt_format_i32_hex

    i32.const offset($db_number_i32)
    i32.const 8
    call $wt_print

    local.get $offset
    if
      i32.const offset($log_watch_memory_1b)
      i32.const length($log_watch_memory_1b)
      call $wt_print

      local.get $offset
      call $wt_format_i32_hex

      i32.const offset($db_number_i32)
      i32.const 8
      call $wt_print
    end

    i32.const offset($log_watch_memory_2)
    i32.const length($log_watch_memory_2)
    call $wt_print

    local.get $address
    local.get $offset
    i32.add
    i32.load
    call $wt_format_i32_hex

    local.get $size
    call $log_print_i32_size

    i32.const offset($log_watch_memory_3)
    i32.const length($log_watch_memory_3)
    call $wt_print

    global.get $log_memory_value_i32
    call $wt_format_i32_hex

    local.get $size
    call $log_print_i32_size

    global.get $wt_color
    if
      i32.const offset($wt_ansi_none)
      i32.const length($wt_ansi_none)
      call $wt_print
    end

    i32.const offset($debug_newline)
    i32.const length($debug_newline)
    call $wt_print

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
    call $log_mem_filter
    local.tee $memrange
    i32.eqz
    if
      local.get $address
      return
    end

    global.get $wt_color
    if
      i32.const offset($wt_ansi_watch)
      i32.const length($wt_ansi_watch)
      call $wt_print
    end

    i32.const offset($log_watch_memory_0)
    i32.const length($log_watch_memory_0)
    call $wt_print

    local.get $memrange
    if
      local.get $memrange
      i32.load offset=8
      i32.const offset($wt_mem_tags)
      i32.add
      local.get $memrange
      i32.load offset=12
      call $wt_print
    end

    ;;i32.const offset($log_watch_memory_1)
    ;;i32.const length($log_watch_memory_1)
    local.get $debug_ptr
    local.get $debug_len
    call $wt_print

    i32.const offset($log_watch_memory_1)
    i32.const length($log_watch_memory_1)
    call $wt_print

    local.get $address
    call $wt_format_i32_hex

    i32.const offset($db_number_i32)
    i32.const 8
    call $wt_print

    local.get $offset
    if
      i32.const offset($log_watch_memory_1b)
      i32.const length($log_watch_memory_1b)
      call $wt_print

      local.get $offset
      call $wt_format_i32_hex

      i32.const offset($db_number_i32)
      i32.const 8
      call $wt_print
    end

    i32.const offset($log_watch_memory_2)
    i32.const length($log_watch_memory_2)
    call $wt_print

    local.get $address
    local.get $offset
    i32.add
    i64.load
    call $wt_format_i64_hex

    local.get $size
    call $log_print_i64_size

    i32.const offset($log_watch_memory_3)
    i32.const length($log_watch_memory_3)
    call $wt_print

    global.get $log_memory_value_i64
    call $wt_format_i64_hex

    local.get $size
    call $log_print_i64_size

    global.get $wt_color
    if
      i32.const offset($wt_ansi_none)
      i32.const length($wt_ansi_none)
      call $wt_print
    end

    i32.const offset($debug_newline)
    i32.const length($debug_newline)
    call $wt_print

    local.get $address
  )

  (data $wt_id_tag "tag_+++")

  (data $log_watch_global_1 "GLOBAL ")
  (data $log_watch_global_2 " value ")
  (data $log_watch_global_3 "=>")

  (data $log_watch_local_1 "LOCAL ")
  (data $log_watch_local_2 " value ")
  (data $log_watch_local_3 "=>")


  (data $log_watch_memory_0 "MEMORY ")
  (data $log_watch_memory_1 " | ")
  (data $log_watch_memory_1b "+")
  (data $log_watch_memory_2 " value ")
  (data $log_watch_memory_3 "=>")

  (global $log_memory_value_i32 (mut i32) (i32.const 0))
  (global $log_memory_value_i64 (mut i64) (i64.const 0))

  (data $debug_here "DEBUGOUT_")

)