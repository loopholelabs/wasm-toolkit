(module

  (func $get_source_len (result i32)
    i32.const length($source_data)  
  )

  (func $get_source (param $dest i32) (result i32)
    local.get $dest
    i32.const offset($source_data)  
    i32.const length($source_data)
    memory.copy
    global.get $source_gzipped
  )

  (data $dummy_data 1)
  (global $source_gzipped (mut i32) (i32.const 0))
)