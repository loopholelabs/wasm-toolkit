(module
  (type (func (param i32)))
  (type (func (param i32 i32)))
  (type (func (result i32)))

  (import "scale" "get_invocation_id" (func $get_invocation_id (type 0)))
  (import "scale" "send_otel_trace_json" (func $send_otel_trace_json (type 1)))
  (import "scale" "get_service_name_len" (func $get_service_name_len (type 2)))
  (import "scale" "get_service_name" (func $get_service_name (type 0)))

  (func $cache_service_name
    global.get $service_name_set
    i32.eqz
    if
      call $get_service_name_len
      global.set $service_name_len
      i32.const offset($service_name)
      call $get_service_name
      i32.const 1
      global.set $service_name_set
    end
  )

  (data $service_name 256)
  (global $service_name_len (mut i32) (i32.const 0))
  (global $service_name_set (mut i32) (i32.const 0))
)