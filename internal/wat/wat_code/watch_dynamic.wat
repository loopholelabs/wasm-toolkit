(module

  (func $check_dynamic_watches (param $start i32) (param $end i32) (result i32)
    (local $p i32)
    (local $cptr i32)

    block
      loop
        local.get $p
        global.get $wt_dynamic_mem_ranges_len
        i32.ge_u
        br_if 1

        ;; Now check the range pointed at by index $p
        local.get $p
        i32.const 4
        i32.shl
        i32.const offset($wt_dynamic_mem_ranges)
        i32.add
        local.set $cptr

        ;; Check the dynamic memory range here
        local.get $start
        local.get $cptr
        i32.load
        i32.ge_u
        if
          ;; $start is greater than the range start
          local.get $end
          local.get $cptr
          i32.load offset=4
          i32.le_u
          if
            ;; We are within the range
            local.get $cptr
            ;; Return ptr to the memory range info
            return
          end
        end

        local.get $p
        i32.const 1
        i32.add
        local.set $p
        br 0
      end
    end
    i32.const 0
  )

  (func $watch_add (param $id i32) (param $addr i32) (param $len i32)
    (local $p i32)

  ;; Add a new entry in $wt_dynamic_mem_ranges
    global.get $wt_dynamic_mem_ranges_len
    i32.const 4
    i32.shl
    i32.const offset($wt_dynamic_mem_ranges)
    i32.add
    local.tee $p
    i32.const length($wt_dynamic_mem_ranges)
    i32.const offset($wt_dynamic_mem_ranges)
    i32.add

    i32.ge_u
    if
      ;; We ran out of watch space...
;;      i32.const offset($error_watch_overflow)
;;      i32.const length($error_watch_overflow)
;;      call $wt_print
      unreachable
    end
    ;; Store things here...
    local.get $p
    local.get $addr
    i32.store
    
    local.get $p
    local.get $addr
    local.get $len
    i32.add
    i32.store offset=4

    ;; Now store the tag. Len=0 means to use the ptr as an ID.
    local.get $p
    local.get $id
    i32.store offset=8

    local.get $p
    i32.const 0
    i32.store offset=12

    global.get $wt_dynamic_mem_ranges_len
    i32.const 1
    i32.add
    global.set $wt_dynamic_mem_ranges_len
  )

  (func $watch_del (param $id i32)
    (local $p i32)
    (local $cptr i32)
    (local $eptr i32)
    ;; Remove an entry from $wt_dynamic_mem_ranges

    block
      loop
        local.get $p
        global.get $wt_dynamic_mem_ranges_len
        i32.ge_u
        br_if 1

        ;; Check if it's the correct ID...
        local.get $p
        i32.const 4
        i32.shl
        i32.const offset($wt_dynamic_mem_ranges)
        i32.add
        local.set $cptr

        local.get $cptr
        i32.load offset=8
        local.get $id
        i32.eq
        if
          ;; We found the correct ID. Now swap with the END element, and dec $wt_dynamic_mem_ranges_len
          global.get $wt_dynamic_mem_ranges_len
          i32.const 1
          i32.sub
          i32.const 4
          i32.shl
          i32.const offset($wt_dynamic_mem_ranges)
          i32.add
          local.set $eptr

          ;; Now copy elements if required...
          local.get $cptr
          local.get $eptr
          i32.ne
          if
            local.get $cptr
            local.get $eptr
            i32.const 16
            memory.copy ;; dst / src / size
          end

          global.get $wt_dynamic_mem_ranges_len
          i32.const 1
          i32.sub
          global.set $wt_dynamic_mem_ranges_len
          return
        end

        local.get $p
        i32.const 1
        i32.add
        local.set $p
        br 0
      end
    end

    ;; The ID was not found - do nothing.
  )

  (data $wt_dynamic_mem_ranges 1600)
  (global $wt_dynamic_mem_ranges_len (mut i32) (i32.const 0))

)