(module
  

;; $debug_summary_func
  (func $debug_summary_func (param $fid i32)
    (local $metric_count i32)
    (local $metrics_ptr i32)
    local.get $fid
    i32.const 4
    i32.shl
    i32.const offset.$debug_function_metrics
    i32.add
    global.get $debug_start_mem
    i32.add
    local.tee $metrics_ptr
    i32.load
    local.tee $metric_count
    i32.eqz
    br_if 0

    local.get $metric_count
    call $db_format_i32_dec_nz

    i32.const offset.$db_number_i32
    global.get $debug_start_mem
    i32.add
    i32.const 10
    ;;length.$db_number_i32
    call $debug_print

    i32.const offset.$debug_table_sep
    global.get $debug_start_mem
    i32.add
    i32.const length.$debug_table_sep
    call $debug_print

;; Now print out total time in ns
    local.get $metrics_ptr
    i64.load offset=4
    call $db_format_i64_dec_nz

    i32.const offset.$db_number_i64
    global.get $debug_start_mem
    i32.add
    i32.const 19
    ;;length.$db_number_i32
    call $debug_print

    i32.const offset.$debug_table_sep
    global.get $debug_start_mem
    i32.add
    i32.const length.$debug_table_sep
    call $debug_print

    local.get $fid
    i32.const 3
    i32.shl
    global.get $debug_start_mem
    i32.add
    i32.const offset.$debug_function_table
    i32.add
    i32.load
    global.get $debug_start_mem
    i32.add
    i32.const offset.$debug_function_data
    i32.add
    local.get 0
    i32.const 3
    i32.shl
    global.get $debug_start_mem
    i32.add
    i32.const offset.$debug_function_table
    i32.add
    i32.load offset=4
    call $debug_print

    i32.const offset.$debug_newline
    global.get $debug_start_mem
    i32.add
    i32.const length.$debug_newline
    call $debug_print
  
  )

;; $debug_summary
  (func $debug_summary
      (local $f_id i32)
      i32.const offset.$debug_summary
      global.get $debug_start_mem
      i32.add
      i32.const length.$debug_summary
      call $debug_print

;; Go through all functions and print details for the function

    block $all_done
      loop $loop
        call $debug_find_expensive_function
        local.tee $f_id
        i32.const -1
        i32.eq
        br_if $all_done

        local.get $f_id
        call $debug_summary_func

        ;; Clear the time metric so that it doesn't get returned again
        local.get $f_id
        i32.const 4
        i32.shl
        i32.const offset.$debug_function_metrics
        i32.add
        global.get $debug_start_mem
        i32.add
        i64.const 0
        i64.store offset=4

        br $loop
      end
    end
  )

  (func $debug_find_expensive_function (result i32)
    (local $f_id i32)
    (local $best_id i32)
    (local $best_val i64)
    (local $metrics_ptr i32)

    i32.const offset.$debug_function_metrics
    global.get $debug_start_mem
    i32.add
    local.set $metrics_ptr

    i32.const 0
    local.set $f_id

    i32.const -1
    local.set $best_id

    i64.const 0
    local.set $best_val

    loop $loop
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
      global.get $debug_num_funcs
      i32.lt_u
      br_if $loop
    end

    local.get $best_id
  )

  (data $debug_summary "\0d\0a-- Summary of execution --\0d\0aCount      | Time (ns)           | Function\0d\0a-----------+---------------------+\0d\0a")
)