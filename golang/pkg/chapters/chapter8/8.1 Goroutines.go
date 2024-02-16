package chapter8

import (
	"fmt"
	"time"
)

func FibonacciUsing(n int) {
	// The first parallel action
	go spinner(time.Millisecond * 100)

	// The second parallel action
	fibN := fibonacci(n)
	fmt.Printf("\rFibonacci (%d) = %d\n", n, fibN)
}

func spinner(delay time.Duration) {
	for {
		for _, r := range `-\|/` {
			fmt.Printf("\r%c", r)
			time.Sleep(delay)
		}
	}
}

func fibonacci(x int) int {
	if x < 2 {
		return x
	}
	return fibonacci(x-1) + fibonacci(x-2)
}
