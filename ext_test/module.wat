;;
;; Test out our extension
;;
;;

(module
  ;; Import the extension functions
  (type (func (param i64) (param i32) (param i32) (result i64)))
  (import "env" "ext_5c7d22390f9101d459292d76c11b5e9f66c327b1766aae34b9cc75f9f40e8206_New" (func $fetch_New (type 0)))
  (import "env" "ext_5c7d22390f9101d459292d76c11b5e9f66c327b1766aae34b9cc75f9f40e8206_HttpConnector_Fetch" (func $fetch_HttpConnector_Fetch (type 0)))

  ;; _start function unusued
  (func $_start)

  ;; initialize function unused
  (func $initialize (result i64)
    i64.const 0
  )

;; Our main run function...
  (func $run (result i64)
    (local $header_size i64)
    (local $h_item_size i64)
    (local $str_ptr_len i64)

    ;; Point to our param buffer
    i32.const 60000
    global.set $read_ptr

    ;; We need to encode an HttpConfig
    i64.const 123         ;; A timeout value
    i32.const 12          ;; KIND.int32   
    call $encode_int

    ;; Now call fetch_New
    i64.const 0           ;; Not used (not an instance call)
    i32.const 60000       ;; -> param buffer
    global.get $read_ptr
    i32.const 60000
    i32.sub               ;; len(param buffer)
    call $fetch_New
    ;; Save the instance ID for later
    global.set $fetch_instance_id

    ;; Now call fetch_HttpConnector_Fetch
    ;; Point to our param buffer
    i32.const 60000
    global.set $read_ptr

    i32.const 90200       ;; ->url
    i32.const 19          ;; len(url)
    i32.const 5           ;; Kind.string
    call $encode_string

    global.get $fetch_instance_id
    i32.const 60000       ;; -> param buffer
    global.get $read_ptr
    i32.const 60000
    i32.sub               ;; len(param buffer)
    call $fetch_HttpConnector_Fetch
    drop                  ;; Doesn't return interface, so we can drop


    ;; We now have an HttpResponse in the buffer. We need to find the body

    i32.const 30004
    global.set $read_ptr
    call $decode_map
    local.set $header_size

    ;; Now loop thru headers
    ;; Now we need to loop through reading [String, StringList] elements
    block $done_headers
      loop $read_headers
        local.get $header_size
        i64.eqz
        br_if $done_headers               ;; Quit if it's zero

        local.get $header_size
        i64.const 1
        i64.sub
        local.set $header_size

        call $decode_string   ;; Read the key
        drop

        call $decode_array    ;; Read the array header
        local.set $h_item_size

        ;; Now loop through reading the strings
        loop $read_values
          local.get $h_item_size
          i64.eqz
          br_if $read_headers     ;; Quit if it's zero

          local.get $h_item_size
          i64.const 1
          i64.sub
          local.set $h_item_size

          call $decode_string     ;; Read the header value
          drop

          br $read_values
        end

        br $read_headers
      end
    end

    ;; Now there's a status

    block
      call $decode_int
      
      i64.const 200
      i64.eq
      br_if 0
      unreachable               ;; The status wasn't 200. Panic()!
    end

    ;; Now we're at the body which is an array of bytes...
    call $decode_string
    ;; Returns a packed (ptr << 32) | len

    local.set $str_ptr_len

    ;; dst, src, size
    i32.const 90022       ;; Just after our return string
    local.get $str_ptr_len
    i64.const 32
    i64.shr_u
    i32.wrap_i64

    local.get $str_ptr_len
    i32.wrap_i64

    memory.copy           ;; Copy the extension body into our own string


    ;; Now lets form our own response
    i32.const 4
    global.set $read_ptr    ;; Reset read ptr to start of buffer...

    i32.const 90000         ;; -> "Hello world from wasm!"
    i32.const 22

    local.get $str_ptr_len
    i32.wrap_i64

    i32.add                 ;; Add the extension response body len

    i32.const 0x05          ;; KIND.string
    call $encode_string     ;; body

    ;; Now pack (ptr << 32) | len
    i64.const 4             ;; ptr
    i64.const 32
    i64.shl                 ;; ptr << 32

    global.get $read_ptr
    i64.extend_i32_u
    i64.const 4             ;; sub start of buffer ptr...
    i64.sub

    i64.or                  ;; return (ptr << 32) | len
    return
  )

;; Resize the global buffer
  (func $resize (param i32) (result i32)
    i32.const 4
    global.set $read_ptr    ;; Reset read ptr to start of buffer...

    i32.const 0
    local.get 0             ;; Get the size arg
    i32.store               ;; Store it in memory
    i32.const 4             ;; Return a ptr of 4, which is just after the size we stored.
  )

;; Resize the fetch ext buffer
  (func $fetch_resize (param i32) (result i32)
    i32.const 30000
    local.get 0
    i32.store
    i32.const 30004
  )

  ;; Encode an array header
  (func $encode_array (param $val_kind i32) (param $len i32)
    global.get $read_ptr
    i32.const 0x01
    i32.store8                ;; Write it

    global.get $read_ptr
    i32.const 1
    i32.add
    global.set $read_ptr

    global.get $read_ptr
    local.get $val_kind
    i32.store8                ;; Write it

    global.get $read_ptr
    i32.const 1
    i32.add
    global.set $read_ptr

    local.get $len
    i64.extend_i32_u
    i32.const 0x0a          ;; KIND.Uint32
    call $encode_uint
  )

  ;; Read and skip an array header
  ;; Returns the size of the array
  (func $decode_array
    (result i64)
    ;; Kind | ValKind | int(size)

    global.get $read_ptr
    i32.const 2
    i32.add                 ;; Skip past the KIND byte
    global.set $read_ptr

    call $decode_uint        ;; Decode the length
    return
  )

  ;; Encode a map header
  (func $encode_map (param $key_kind i32) (param $val_kind i32) (param $len i32)
    global.get $read_ptr
    i32.const 0x02
    i32.store8                ;; Write it

    global.get $read_ptr
    i32.const 1
    i32.add
    global.set $read_ptr

    global.get $read_ptr
    local.get $key_kind
    i32.store8                ;; Write it

    global.get $read_ptr
    i32.const 1
    i32.add
    global.set $read_ptr

    global.get $read_ptr
    local.get $val_kind
    i32.store8                ;; Write it

    global.get $read_ptr
    i32.const 1
    i32.add
    global.set $read_ptr

    local.get $len
    i64.extend_i32_u
    i32.const 0x0a          ;; KIND.Uint32
    call $encode_uint
  )

  ;; Read and skip a map header
  ;; Returns the size of the map
  (func $decode_map
    (result i64)
    ;; Kind | KeyKind | ValKind | int(size)

    global.get $read_ptr
    i32.const 3
    i32.add                 ;; Skip past the KIND byte
    global.set $read_ptr

    call $decode_uint        ;; Decode the length
    return
  )

  ;; Encode a string or uint8array
  (func $encode_string (param $ptr i32) (param $len i32) (param $kind i32)
    global.get $read_ptr
    local.get $kind
    i32.store8                ;; Write it

    global.get $read_ptr
    i32.const 1
    i32.add
    global.set $read_ptr

    local.get $len
    i64.extend_i32_u
    i32.const 0x0a          ;; KIND.Uint32
    call $encode_uint

    ;; Now copy the data
    global.get $read_ptr
    local.get $ptr
    local.get $len
    memory.copy

    global.get $read_ptr
    local.get $len
    i32.add
    global.set $read_ptr
  )

  ;; Read and skip a string
  ;; Returns a packed (ptr << 32) | len
  (func $decode_string
    (result i64)
    (local $string_len i32)
    global.get $read_ptr
    i32.const 1
    i32.add                 ;; Skip past the KIND byte
    global.set $read_ptr

    call $decode_uint        ;; Decode the length
    i32.wrap_i64
    local.tee $string_len

    global.get $read_ptr
    i64.extend_i32_u
    i64.const 32
    i64.shl
    local.get $string_len
    i64.extend_i32_u
    i64.or                  ;; Setup return value of (ptr << 32) | len

    local.get $string_len
    global.get $read_ptr    ;; Skip past the bytes
    i32.add
    global.set $read_ptr

    return
  )

  ;; Write a varint (signed)
  (func $encode_int (param $value i64) (param $kind i32)
    local.get $value
    i64.const 1
    i64.shl
    ;; TODO: Honour sign flag
    local.get $kind
    call $encode_uint
  )

  ;; Write a varint
  (func $encode_uint (param $value i64) (param $kind i32)
    (local $current_ptr i32)
    (local $current_val i64)
    (local $current_shift i64)
    (local $current_byte i64)
    
    i64.const 0
    local.set $current_shift

    local.get $value
    local.set $current_val

    global.get $read_ptr
    local.get $kind
    i32.store8                ;; Write it

    global.get $read_ptr
    i32.const 1
    i32.add
    local.set $current_ptr    ;; Store it locally
    block $done_write_int
      loop $write_int
        local.get $current_ptr

        local.get $current_val
        i64.const 0x7f
        i64.and

        local.get $current_val
        i64.const 0x7f
        i64.gt_u
        i64.extend_i32_u
        i64.const 7
        i64.shl
        i64.or

        i64.store8              ;; Store value

        local.get $current_ptr
        i32.const 1
        i32.add
        local.set $current_ptr  ;; $current_ptr++

        local.get $current_val
        i64.const 0x80
        i64.lt_u
        br_if $done_write_int

        local.get $current_val
        i64.const 7
        i64.shr_u
        local.set $current_val

        local.get $current_shift
        i64.const 7
        i64.add
        local.set $current_shift

        br $write_int
      end
    end

    local.get $current_ptr
    global.set $read_ptr
    return
  )

  ;; Read and skip a varint signed
  (func $decode_int (result i64)
    call $decode_uint
    i64.const 1
    i64.shr_u
    ;; TODO: Honour sign flag
    return
  )

  ;; Read and skip a varint
  ;; Returns the int value
  (func $decode_uint    (result i64)
    (local $current_ptr i32)
    (local $current_val i64)
    (local $current_shift i64)
    (local $current_byte i64)

    i64.const 0
    local.set $current_shift

    i64.const 0
    local.set $current_val  ;; Start with 0

    global.get $read_ptr
    i32.const 1
    i32.add                 ;; Skip past the KIND byte
    local.set $current_ptr  ;; Store it locally
    loop $read_int
      local.get $current_ptr
      i64.load8_u           ;; Load up a byte
      local.tee $current_byte

      i64.const 0x7f
      i64.and
      local.get $current_shift
      i64.shl
      local.get $current_val
      i64.or
      local.set $current_val

      local.get $current_shift
      i64.const 7
      i64.add
      local.set $current_shift    ;; current_shift += 7

                            ;; Move forward 1 byte
      local.get $current_ptr
      i32.const 1
      i32.add
      local.set $current_ptr

      local.get $current_byte
      i32.wrap_i64
      i32.const 7
      i32.shr_u
      ;; If the top bit was set, loop
      br_if $read_int
    end

    local.get $current_ptr
    global.set $read_ptr

    local.get $current_val  ;; Return the value
    return
  )

  (func $debug (param i32))

  (global $read_ptr (mut i32) (i32.const 0))

  (global $fetch_instance_id (mut i64) (i64.const 0))

  (memory $mem 2)           ;; How much memory we need in 64K pages to start with...

  (export "memory" (memory $mem))

  (export "_start" (func $_start))
  (export "run" (func $run))
  (export "resize" (func $resize))
  (export "initialize" (func $initialize))

  (export "ext_5c7d22390f9101d459292d76c11b5e9f66c327b1766aae34b9cc75f9f40e8206_Resize" (func $fetch_resize))

  ;; Some strings we need...
  (data $.data (i32.const 90000) "Hello world from wasm ")

  (data $id (i32.const 90100) "    ")

  (data $url (i32.const 90200) "https://ifconfig.me")

  ;; Memory
  ;; 0      read buffer
  ;; 30000  fetch read buffer
  ;; 60000  params buffer
  ;; 90000  data
)