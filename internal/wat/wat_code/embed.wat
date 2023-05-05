(module
  (type (func (param i32 i32 i32 i32) (result i32)))
  (type (func (param i32 i32) (result i32)))
  (type (func (param i32 i32 i32 i32 i32 i64 i64 i32 i32) (result i32)))
  (import "wasi_snapshot_preview1" "fd_write" (func $debug_fd_write (type 0)))
  (import "wasi_snapshot_preview1" "fd_prestat_get" (func $fd_prestat_get (type 1)))
  (import "wasi_snapshot_preview1" "path_open" (func $path_open (type 2)))
  (import "wasi_snapshot_preview1" "fd_read" (func $fd_read (type 0)))

  (func $wrap_fd_prestat_get (param $fd i32) (param $ptr i32) (result i32)
    local.get 0
    local.get 1
    call $fd_prestat_get
  )

  (func $wrap_fd_read (param $fd i32) (param $iovs i32) (param $iovsLen i32) (param $nread i32) (result i32)
    (local $bytes i32)
    (local $iov_offset i32)
    (local $ioptr i32)
    (local $iolen i32)
    (local $current_iov i32)
    (local $count i32)

    i32.const 0
    local.set $iov_offset

    i32.const 0
    local.set $bytes

    local.get $fd
    i32.const 90001
    i32.eq
    if

      block
        loop
        local.get $iov_offset
        local.get $iovsLen
        i32.ge_u
        br_if 1
        ;; Now go through the iov and copy as much data as we can...

          local.get $iov_offset
          i32.const 3
          i32.shl
          local.get $iovs
          i32.add
          local.set $current_iov

          local.get $current_iov
          i32.load
          local.set $ioptr

          local.get $current_iov
          i32.load offset=4
          local.set $iolen

          ;; ioptr and iolen are now setup

          i32.const 0
          local.set $count

          block
            loop
              local.get $count
              local.get $iolen
              i32.ge_u
              br_if 1

              ;; Now do the actual write if we have any data to write...

              ;; If there is no more data, then return
              global.get $debug_fp
              i32.const length($file_content)
              i32.ge_u
              if
                local.get $nread
                local.get $bytes
                i32.store

                ;; WASI_SUCCESS
                i32.const 0
                return
              end

              ;; preload addr for write
              local.get $count
              local.get $ioptr
              i32.add

              ;; Write a byte
              i32.const offset($file_content)
              global.get $debug_fp
              i32.add
              i32.load8_u

              i32.store8

              global.get $debug_fp
              i32.const 1
              i32.add
              global.set $debug_fp

              local.get $bytes
              i32.const 1
              i32.add
              local.set $bytes

              local.get $count
              i32.const 1
              i32.add
              local.set $count
              br 0
            end
          end

          local.get $iov_offset
          i32.const 1
          i32.add
          local.set $iov_offset
          br 0
        end
      end


      local.get $nread
      local.get $bytes
      i32.store

      ;; WASI_SUCCESS
      i32.const 0
      return
    end

    local.get 0
    local.get 1
    local.get 2
    local.get 3
    call $fd_read
  )

  (func $wrap_path_open (param $dirfd i32) (param $dirflags i32) (param $pathPtr i32) (param $pathLen i32) (param $oflags i32) (param $fsRightsBase i64) (param $fsRightsInheriting i64) (param $fsFlags i32) (param $fd i32) (result i32)
    (local $count i32)
    ;; TODO: Check if the path is something we know about...

    i32.const 0
    local.set $count

    local.get $pathLen
    i32.const length($file_name)
    i32.eq
    if

      block
        loop
          local.get $count
          i32.const length($file_name)
          i32.eq
          if
            ;; Do the stuff

            ;; Set the FD
            local.get $fd
            i32.const 90001
            i32.store

            ;; Reset the read ptr
            i32.const 0
            global.set $debug_fp

;;            i32.const offset($debug_path_open)
;;            i32.const length($debug_path_open)
;;            call $debug_print

;;            local.get $pathPtr
;;            local.get $pathLen
;;            call $debug_print

;;            i32.const offset($debug_newline)
;;            i32.const length($debug_newline)
;;            call $debug_print

            ;; WASI_ESUCCESS
            i32.const 0
            return

          end

          ;; cmp bytes
          local.get $pathPtr
          local.get $count
          i32.add
          i32.load8_u

          i32.const offset($file_name)
          local.get $count
          i32.add
          i32.load8_u
          
          i32.ne
          br_if 1 

          local.get $count
          i32.const 1
          i32.add
          local.set $count
          br 0
        end
      end

    end

    local.get 0
    local.get 1
    local.get 2
    local.get 3
    local.get 4
    local.get 5
    local.get 6
    local.get 7
    local.get 8
    call $path_open
  )



  (func $debug_print (param $ptr i32) (param $len i32)
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

;; $db_format_i32 as hex into the buffer ($db_number_i32)
  (func $db_format_i32_hex (param $num i32)
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

;; $db_format_i32 as dec into the buffer ($db_number_i32)
  (func $db_format_i32_dec (param $num i32)
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

;; $db_format_i32 as dec into the buffer ($db_number_i32)
  (func $db_format_i32_dec_nz (param $num i32)
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

;; $db_format_i64 as dec into the buffer ($db_number_i64)
  (func $db_format_i64_dec_nz (param $num i64)
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

  (data $db_hex "0123456789ABCDEF ")
  (data $db_number_i32 10)
  (data $db_number_i64 19)

  (data $iovec 8)
  (data $bytes_written 4)

  (data $debug_newline "\0d\0a")
  (data $debug_path_open "+ path_open ")

  (data $debug_here "HERE\0d\0a")

;; TODO: Support multiple files here...
  (global $debug_fp (mut i32) (i32.const 0))

)