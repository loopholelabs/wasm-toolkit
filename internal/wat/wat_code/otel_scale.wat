(module
  (type (func (param i32)))
  (type (func (param i32 i32)))

  (import "scale" "get_invocation_id" (func $get_invocation_id (type 0)))
  (import "scale" "send_otel_trace_json" (func $send_otel_trace_json (type 1)))

)