(module

  ;; Color output settings
  (data $wt_ansi_param "\1b[32m")
  (data $wt_ansi_result "\1b[31m")
  (data $wt_ansi_context "\1b[36m")
  (data $wt_ansi_param_name "\1b[33m")
  (data $wt_ansi_none "\1b[0m")
  (data $wt_ansi_wasi_context "\1b[35m")

  ;; Flag - is color enabled or not?
  (global $wt_color i32 (i32.const 0))
)