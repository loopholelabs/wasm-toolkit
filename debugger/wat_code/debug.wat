(module

;; $debug_summary_maybe
  (func $debug_summary_maybe
    block $summary
      global.get $debug_current_stack_depth
      i32.const 0
      i32.gt_u
      br_if $summary
      call $debug_summary
    end  
  )

  (data $debug_return_value " => ")
  (data $debug_value_i32 "i32:")
  (data $debug_value_i64 "i64:")
  (data $debug_value_f32 "f32:")
  (data $debug_value_f64 "f64:")

  (data $debug_param_start "(")
  (data $debug_param_sep ", ")
  (data $debug_param_end ")")

  (data $debug_newline "\0d\0a")
  (data $debug_enter "-> ")
  (data $debug_exit "<- ")
  (data $debug_single_sp " ")
  (data $debug_sp "  ")
  (data $debug_table_sep " | ")
  (data $debug_memory_change " => ")

  (global $debug_current_stack_depth (mut i32) (i32.const 0))
)
