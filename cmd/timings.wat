(module
  (type (func (param i32 i64 i32) (result i32)))
  (import "wasi_snapshot_preview1" "clock_time_get" (func $debug_clock_time_get (type 0)))

  (func $debug_gettime (result i64)
    i32.const 0
    i64.const 1000
    i32.const offset($debug_clock_loc)
    call $debug_clock_time_get
    drop
    i32.const offset($debug_clock_loc)
    i64.load
  )

  (func $timings_enter_func (param $fid i32)
    (local $metrics_ptr i32)
;; Update metrics

    ;; metrics entry
    ;; 4 bytes i32  Function call counter
    ;; 8 bytes i64  Total time

    ;; Inc call counter
    i32.const offset($metrics_data)
    local.get $fid
    i32.const 4
    i32.shl
    i32.add
    local.tee $metrics_ptr
    local.get $metrics_ptr
    i32.load
    i32.const 1
    i32.add
    i32.store

    ;; Store enter timestamp in metrics stack
    global.get $debug_timestamps_stack_pointer
    i32.const offset($debug_timestamps_stack)
    i32.add
    call $debug_gettime
    i64.store

    global.get $debug_timestamps_stack_pointer
    i32.const 8
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

  (func $timings_exit_func (param $fid i32)
    (local $metrics_ptr i32)

    ;; Update metrics
    local.get $fid
    i32.const 4
    i32.shl
    i32.const offset($metrics_data)
    i32.add
    local.tee $metrics_ptr

    call $debug_gettime

    ;; Pop timestamp off the timestamp stack
    global.get $debug_timestamps_stack_pointer
    i32.const 8
    i32.sub
    global.set $debug_timestamps_stack_pointer

    global.get $debug_timestamps_stack_pointer
    i32.const offset($debug_timestamps_stack)
    i32.add
    i64.load

    i64.sub

    local.get $metrics_ptr
    i64.load offset=4
    i64.add
    i64.store offset=4  
  )

  (func $debug_summary_maybe
;;    global.get $debug_current_stack_depth
;;    i32.eqz
;;    if
;;      call $show_timings_summary
;;    end
  )


;; $debug_summary_func
  (func $debug_summary_func (param $fid i32)
    (local $metric_count i32)
    (local $metrics_ptr i32)
    local.get $fid
    i32.const 4
    i32.shl
    i32.const offset($metrics_data)
    i32.add
    local.tee $metrics_ptr
    i32.load
    local.tee $metric_count
    i32.eqz
    br_if 0

    local.get $metric_count
    call $wt_format_i32_dec_nz

    i32.const offset($db_number_i32)
    i32.const 10
    call $wt_print

    i32.const offset($debug_table_sep)
    i32.const length($debug_table_sep)
    call $wt_print

;; Now print out total time in ns
    local.get $metrics_ptr
    i64.load offset=4
    call $wt_format_i64_dec_nz

    i32.const offset($db_number_i64)
    i32.const 19
    call $wt_print

    i32.const offset($debug_table_sep)
    i32.const length($debug_table_sep)
    call $wt_print

    local.get $fid
    call $wt_print_function_name

    i32.const offset($debug_newline)
    i32.const length($debug_newline)
    call $wt_print
  
  )

  (func $show_timings_summary
      (local $f_id i32)

      global.get $debug_do_timings
      i32.eqz
      br_if 0

      i32.const offset($debug_summary)
      i32.const length($debug_summary)
      call $wt_print

;; Go through all functions and print details for the function

    block
      loop
        call $debug_find_expensive_function
        local.tee $f_id
        i32.const -1
        i32.eq
          br_if 1

        local.get $f_id
        call $debug_summary_func

        ;; Clear the time metric so that it doesn't get returned again
        local.get $f_id
        i32.const 4
        i32.shl
        i32.const offset($metrics_data)
        i32.add
        i64.const 0
        i64.store offset=4

        br 0
      end
    end
  )

  (func $debug_find_expensive_function (result i32)
    (local $f_id i32)
    (local $best_id i32)
    (local $best_val i64)
    (local $metrics_ptr i32)

    i32.const offset($metrics_data)
    local.set $metrics_ptr

    i32.const 0
    local.set $f_id

    i32.const -1
    local.set $best_id

    i64.const 0
    local.set $best_val

    loop
      local.get $metrics_ptr
      i64.load offset=4
      local.get $best_val
      i64.gt_u
      if
        local.get $metrics_ptr
        i64.load offset=4
        local.set $best_val
        local.get $f_id
        local.set $best_id
      end
      ;; Check if it's better

      local.get $metrics_ptr
      i32.const 16
      i32.add
      local.set $metrics_ptr

      local.get $f_id
      i32.const 1
      i32.add
      local.tee $f_id
      global.get $wt_all_function_length
      i32.lt_u
      br_if 0
    end

    local.get $best_id
  )


  (data $debug_clock_loc 8)

  (data $debug_summary "\0d\0a-- Summary of execution --\0d\0aCount      | Time (ns)           | Function\0d\0a-----------+---------------------+\0d\0a")

  ;; Only allow 100 function stack for now
  (data $debug_timestamps_stack 800)

  (global $debug_timestamps_stack_pointer (mut i32) (i32.const 0))

)