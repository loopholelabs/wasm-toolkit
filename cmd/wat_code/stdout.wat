(module
  (type (func (param i32 i32 i32 i32) (result i32)))
  (import "wasi_snapshot_preview1" "fd_write" (func $debug_fd_write (type 0)))

;; $wt_print - Print a string
  (func $wt_print (param $ptr i32) (param $len i32)
    (local $iovp i32)

    i32.const offset($iovec)

    local.tee $iovp
    local.get $ptr
    i32.store

    local.get $iovp
    local.get $len
    i32.store offset=4

    i32.const 2
    local.get $iovp
    i32.const 1
    i32.const offset($bytes_written)
    call $debug_fd_write
    drop
  )

;; $wt_format_i32 as hex into the buffer ($db_number_i32)
  (func $wt_format_i32_hex (param $num i32)
    (local $count i32)
    (local $ptr i32)
    (local $shift_val i32)

    i32.const offset($db_number_i32)
    local.set $ptr

    i32.const 28
    local.set $shift_val

    ;; 28 / 24 / 20 / 16 / 12 / 8 / 4 / 0

    i32.const 0
    local.set $count

    loop
      ;; Work out the value to store...

      local.get $ptr

      local.get $num
      local.get $shift_val
      i32.shr_u
      i32.const 15
      i32.and

      i32.const offset($db_hex)
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
      br_if 0
    end
  )

;; $wt_format_i32 as dec into the buffer ($db_number_i32)
  (func $wt_format_i32_dec (param $num i32)
    (local $count i32)
    (local $ptr i32)
    (local $divide_val i32)
    (local $in_number i32)
    (local $store_value i32)

    i32.const offset($db_number_i32)
    local.set $ptr

    i32.const 1000000000
    local.set $divide_val

    i32.const 0
    local.set $count

    loop
      ;; Work out the value to store...

      local.get $ptr

      local.get $num
      local.get $divide_val
      i32.div_u
      i32.const 10
      i32.rem_u

      i32.const offset($db_hex)
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
      br_if 0
    end
  )

;; $wt_format_i32 as dec into the buffer ($db_number_i32)
  (func $wt_format_i32_dec_nz (param $num i32)
    (local $count i32)
    (local $ptr i32)
    (local $divide_val i32)
    (local $in_number i32)
    (local $store_value i32)

    i32.const offset($db_number_i32)
    local.set $ptr

    i32.const 1000000000
    local.set $divide_val

    i32.const 0
    local.set $count

    loop
      ;; Work out the value to store...

      local.get $ptr

      local.get $num
      local.get $divide_val
      i32.div_u
      i32.const 10
      i32.rem_u
      local.tee $store_value
      i32.eqz
      if (result i32)
        local.get $in_number
        if (result i32)
          local.get $store_value
        else
          i32.const 16
        end
      else
        i32.const 1
        local.set $in_number
        local.get $store_value
      end

      i32.const offset($db_hex)
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
      br_if 0
    end
  )

;; $wt_format_i64 as dec into the buffer ($db_number_i64)
  (func $wt_format_i64_dec_nz (param $num i64)
    (local $count i32)
    (local $ptr i32)
    (local $divide_val i64)
    (local $in_number i32)
    (local $store_value i32)

    i32.const offset($db_number_i64)
    local.set $ptr

    i64.const 1000000000000000000
    local.set $divide_val

    i32.const 0
    local.set $count

    loop
      ;; Work out the value to store...

      local.get $ptr

      local.get $num
      local.get $divide_val
      i64.div_u
      i64.const 10
      i64.rem_u
      i32.wrap_i64
      local.tee $store_value
      i32.eqz
      if (result i32)
        local.get $in_number
        if (result i32)
          local.get $store_value
        else
          i32.const 16
        end
      else
        i32.const 1
        local.set $in_number
        local.get $store_value
      end

      i32.const offset($db_hex)
      i32.add

      i32.load8_u
      i32.store8

      local.get $ptr
      i32.const 1
      i32.add
      local.set $ptr

      local.get $divide_val
      i64.const 10
      i64.div_u
      local.set $divide_val

      local.get $count
      i32.const 1
      i32.add
      local.tee $count
      i32.const 19
      i32.lt_u
      br_if 0
    end
  )

;; $wt_format_i64 as dec into the buffer ($db_number_i64)
  (func $wt_format_i64_dec (param $num i64)
    (local $count i32)
    (local $ptr i32)
    (local $divide_val i64)
    (local $store_value i32)

    i32.const offset($db_number_i64)
    local.set $ptr

    i64.const 1000000000000000000
    local.set $divide_val

    i32.const 0
    local.set $count

    loop
      ;; Work out the value to store...

      local.get $ptr

      local.get $num
      local.get $divide_val
      i64.div_u
      i64.const 10
      i64.rem_u
      i32.wrap_i64
      i32.const offset($db_hex)
      i32.add

      i32.load8_u
      i32.store8

      local.get $ptr
      i32.const 1
      i32.add
      local.set $ptr

      local.get $divide_val
      i64.const 10
      i64.div_u
      local.set $divide_val

      local.get $count
      i32.const 1
      i32.add
      local.tee $count
      i32.const 19
      i32.lt_u
      br_if 0
    end
  )

  ;; $wt_print_hex prints the input as a series of hex bytes (non optimal)
  (func $wt_print_hex (param $ptr i32) (param $len i32)
    (local $p i32)
  
    block
      loop
        local.get $p
        local.get $len
        i32.ge_u
        br_if 1

        local.get $ptr
        local.get $p
        i32.add
        i32.load8_u
        i32.const 4
        i32.shr_u

        i32.const offset($db_hex)
        i32.add
        i32.const 1
        call $wt_print

        local.get $ptr
        local.get $p
        i32.add
        i32.load8_u
        i32.const 15
        i32.and

        i32.const offset($db_hex)
        i32.add
        i32.const 1
        call $wt_print

        local.get $p
        i32.const 1
        i32.add
        local.set $p
        br 0
      end
    end
  )

  ;; $wt_print_hex converts the input to a series of hex bytes
  (func $wt_conv_hex (param $ptr i32) (param $len i32) (param $dest i32)
    (local $p i32)
  
    block
      loop
        local.get $p
        local.get $len
        i32.ge_u
        br_if 1

        ;; first nibble
        local.get $p
        i32.const 1
        i32.shl
        local.get $dest
        i32.add

        local.get $ptr
        local.get $p
        i32.add
        i32.load8_u
        i32.const 4
        i32.shr_u

        i32.const offset($db_hex)
        i32.add

        i32.load8_u
        i32.store8

        ;; second nibble
        local.get $p
        i32.const 1
        i32.shl
        local.get $dest
        i32.add

        local.get $ptr
        local.get $p
        i32.add
        i32.load8_u
        i32.const 15
        i32.and

        i32.const offset($db_hex)
        i32.add

        i32.load8_u
        i32.store8 offset=1

        local.get $p
        i32.const 1
        i32.add
        local.set $p
        br 0
      end
    end
  )

  ;; Structures needed for fd_write
  (data $iovec 8)
  (data $bytes_written 4)

  ;; For number to string conversions
  (data $db_hex "0123456789abcdef ")
  (data $db_number_i32 10)
  (data $db_number_i64 19)

)