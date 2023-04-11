(module
  (global $debug_mem_size (mut i32) (i32.const 0))
  (global $debug_start_mem (mut i32) (i32.const 0)) 

  (func $debug_memory_size (result i32)
    memory.size
    global.get $debug_mem_size
    i32.sub  
  )

  (func $debug_memory_grow (param i32) (result i32)
    call $debug_memory_size

    local.get 0
    memory.grow
    drop

    global.get $debug_start_mem
    local.get 0
    i32.const 16
    i32.shl
    i32.add

    global.get $debug_start_mem

    global.get $debug_mem_size
    i32.const 16
    i32.shl

    memory.copy

    global.get $debug_start_mem
    local.get 0
    i32.const 16
    i32.shl
    i32.add
    global.set $debug_start_mem
  )

)
