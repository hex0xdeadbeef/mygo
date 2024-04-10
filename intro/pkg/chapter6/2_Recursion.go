package chapter6

func recursiveFactorialComputation(number int) int {
	if number == 1 {
		return 1
	} else {
		return number * recursiveFactorialComputation(number-1)
	}
}
