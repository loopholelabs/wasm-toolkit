(module

  (func $test_stdout
    i32.const offset($test_hello_world)
    i32.const length($test_hello_world)
    call $wt_print

    ;; Check a few things basically work...

    i32.const 0x12345678
    call $wt_format_i32_hex
    i32.const offset($db_number_i32)
    i32.const 8
    call $wt_print

    i64.const 0x123456789abcdef0
    call $wt_format_i64_hex
    i32.const offset($db_number_i64)
    i32.const 16
    call $wt_print

    i32.const offset($test_data)
    i32.const length($test_data)
    call $wt_print_hex
  )

  (export "$test_stdout" (func $test_stdout))

  (data $test_hello_world "Hello world")
  (data $test_data "abcdefgh")
)