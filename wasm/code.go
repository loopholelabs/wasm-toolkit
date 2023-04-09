package wasm

// Here, we're going to do some code analysis.
// Essentially, we want to find instances of
// local.get 0
//
// load / store instruction
//
// Then we can say that param0 is of type pointer, rather than just an i32
//
// If the param is passed to another function, then we need to check that function onwards etc

func IsPointer(f Func, pindex int) bool {
	// Find each instance of "local.get pindex" and trace forwards to see what happens to the value.
	// Ideally we would also want to trace it going through another local/global before use.

	return false
}
