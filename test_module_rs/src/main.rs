static FUNNY_GLOBAL:i32 = 99;


fn main() {
	println!("Hello from main.\n");
  let _ = example_function(12, 34);
}

#[export_name = "hello"]
fn hello() {
	println!("Hello from hello :)\n");
}

#[export_name = "example"]
fn example_function(value_x:i32, value_y:i32) -> i32 {
  return (value_x * value_y) + FUNNY_GLOBAL;
}
