package popcount

var (
	pc         [256]byte // it has the initialized zero values that correspond to its type
	testValues []int     // It has nil value
	jetbrains  string    // It has nil value
)

/* We can place unlimited init() statements in a source file*/

// There is init() function that initializes the variables at the package-level scope
func init() {
	for i, _ := range pc { // range is [0, 256). We get the index and the value here
		pc[i] = pc[i/2] + byte(i&1) // This statement counts the amount of the engaged bits for each value dynamically
	}
}

func init() {
	jetbrains = "govno"
}

func PopCountLoop(x uint64) int {
	var total int
	for i := 0; i < 8; i++ {
		byteNumber := pc[byte(x>>(i*8))] // There will be converted i-th next 8 bits from the left hand side
		total += int(byteNumber)
	}
	return total
}

func PopCounSimple(x uint64) int {
	return int((pc[byte(x>>(0*8))] +
		pc[byte(x>>(1*8))] +
		pc[byte(x>>(2*8))] +
		pc[byte(x>>(3*8))] +
		pc[byte(x>>(4*8))] +
		pc[byte(x>>(5*8))] +
		pc[byte(x>>(6*8))] +
		pc[byte(x>>(7*8))]))
}

func PopCountWithOneBitShifting(x uint64) int {
	var total int
	for i := 0; i < 64; i++ {
		total += int(x >> (i * 64))
	}
	return total
}

func PopCountWithSubractionPrecedingValue(x int64) int {
	var total int
	for x != 0 {
		total++
		x &= x - 1
	}
	return total
}
