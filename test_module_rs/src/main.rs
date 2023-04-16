
fn main() {
	println!("Hello from main.\n")
}

#[export_name = "hello"]
fn hello() {
	println!("Hello from hello :)\n")
}
