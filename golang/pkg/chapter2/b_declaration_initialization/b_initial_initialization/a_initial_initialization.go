package a_initial_initialization

var (
	a, b, c int
)

// init() will initialize global variables after import statement finishes
func init() {
	a = b + c
	b = function()
	c = 1
}
func function() int {
	return c + 1
}