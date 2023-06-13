(module

  (func $test_stdout
    i32.const offset($test_hello_world)
    i32.const length($test_hello_world)
    call $wt_print
  )

  (export "$test_stdout" (func $test_stdout))

  (data $test_hello_world "Hello world")
)