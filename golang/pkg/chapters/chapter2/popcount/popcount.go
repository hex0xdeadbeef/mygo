package popcount

var (
	pc [256]byte
)

// pc[i] is the population count of i var pc [256]byte
func init() {
	for i := range pc {
		pc[i] = pc[i/2] + byte(i&1)
	}
}

// PopCount returns the population count (number of set bits) of x.
func PopCount(x uint64) int {
	return int(pc[byte(x>>(0*8))] + pc[byte(x>>(1*8))] + pc[byte(x>>(2*8))] + pc[byte(x>>(3*8))] + pc[byte(x>>(4*8))] + pc[byte(x>>(5*8))] + pc[byte(x>>(6*8))] + pc[byte(x>>(7*8))])
}

func PopCountLoop(x uint64) int {
	var total int
	for i := 0; i < 8; i++ {
		currentEightBytes := x << (i * 8)
		total += int(pc[byte(currentEightBytes)])
	}
	return total
}

func PopCountWithOneBitShifting(x uint64) int {
	var total int
	for i := 0; i < 64; i++ {
		total += int(1 & (x >> i))
	}
	return total
}

func PopCountFast(x uint64) int {
	var total int
	for x != 0 {
		x &= (x - 1)
		total++
	}

	return total
}
