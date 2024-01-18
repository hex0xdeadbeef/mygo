package b_declaration_initialization

var (
	a, b, c int
)

// init() will initialize global variables after import statement finishes
func init() {
	a = b + c      // 3
	b = function() // 2
	c = 1          // 1
}

// This function access the value of global initialized variable c (Compiler understands that it already has a value)
func function() int {
	return c + 1
}
