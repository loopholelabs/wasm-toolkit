static FUNNY_GLOBAL:i32 = 99;


fn main() {
	println!("Hello from main.\n");
  let v = example_function(12, 34);

  println!("Example returned {v}");
}

#[export_name = "hello"]
fn hello() {
	println!("Hello from hello :)\n");
}

#[export_name = "example"]
fn example_function(value_x:i32, value_y:i32) -> i32 {
  let mut local_var:i32 = 78;

  // Watch local_var here...
  unsafe {
    let ptr = &local_var as *const i32;
    println!("Calling watch on {}", ptr as i32);
    watch(1, ptr, 4);
  }

  if value_x==0 {
    local_var+=1;
  }

  local_var = 0x1234;

  if value_y==0 {
    local_var+=1;
  }

  local_var = 0x999;

  return (value_x * value_y) + FUNNY_GLOBAL + local_var;
}

#[link(wasm_import_module = "scale")]
extern "C" {
    #[link_name = "watch"]
    fn watch(id: u32, ptr: *const i32, size: u32);
}

#[link(wasm_import_module = "scale")]
extern "C" {
    #[link_name = "unwatch"]
    fn unwatch(id: u32);
}
