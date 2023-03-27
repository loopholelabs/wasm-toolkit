(module
  (import "wasi_snapshot_preview1" "fd_write" (func $debug_fd_write (param i32 i32 i32 i32) (result i32)))

;; $debug_print
  (func $debug_print (param $ptr i32) (param $len i32)
    (local $iovp i32)

    i32.const offset.$iovec
    global.get $debug_start_mem
    i32.add
    local.set $iovp

    local.get $iovp
    local.get $ptr
    i32.store

    local.get $iovp
    i32.const 4
    i32.add
    local.get $len
    i32.store

    i32.const 1
    local.get $iovp
    i32.const 1
    i32.const offset.$bytes_written
    global.get $debug_start_mem
    i32.add
    call $debug_fd_write
    drop
  )

;; $db_format_i32 as hex into the buffer ($db_number_i32)
  (func $db_format_i32_hex (param $num i32)
    (local $count i32)
    (local $ptr i32)
    (local $shift_val i32)

    i32.const offset.$db_number_i32
    global.get $debug_start_mem
    i32.add
    local.set $ptr

    i32.const 28
    local.set $shift_val

    ;; 28 / 24 / 20 / 16 / 12 / 8 / 4 / 0

    i32.const 0
    local.set $count

    loop $ldb
      ;; Work out the value to store...

      local.get $ptr

      local.get $num
      local.get $shift_val
      i32.shr_u
      i32.const 0x0f
      i32.and

      i32.const offset.$db_hex
      i32.add
      global.get $debug_start_mem
      i32.add

      i32.load8_u

      i32.store8

      local.get $ptr
      i32.const 1
      i32.add
      local.set $ptr

      local.get $shift_val
      i32.const 4
      i32.sub
      local.set $shift_val

      local.get $count
      i32.const 1
      i32.add
      local.tee $count
      i32.const 8
      i32.lt_u
      br_if $ldb
    end
  )

;; $db_format_i32 as dec into the buffer ($db_number_i32)
  (func $db_format_i32_dec (param $num i32)
    (local $count i32)
    (local $ptr i32)
    (local $divide_val i32)
    (local $in_number i32)
    (local $store_value i32)

    i32.const offset.$db_number_i32
    global.get $debug_start_mem
    i32.add
    local.set $ptr

    i32.const 1000000000
    local.set $divide_val

    i32.const 0
    local.set $count

    loop $ldb
      ;; Work out the value to store...

      local.get $ptr

      local.get $num
      local.get $divide_val
      i32.div_u
      i32.const 10
      i32.rem_u

      i32.const offset.$db_hex
      i32.add
      global.get $debug_start_mem
      i32.add

      i32.load8_u
      i32.store8

      local.get $ptr
      i32.const 1
      i32.add
      local.set $ptr

      local.get $divide_val
      i32.const 10
      i32.div_u
      local.set $divide_val

      local.get $count
      i32.const 1
      i32.add
      local.tee $count
      i32.const 10
      i32.lt_u
      br_if $ldb
    end
  )

  (data $db_hex "0123456789ABCDEF")
  (data $db_number_i32 10)

  (data $iovec 8)
  (data $bytes_written 4)
)
