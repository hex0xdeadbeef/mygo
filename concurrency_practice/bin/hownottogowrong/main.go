package main

import (
	"fmt"
	"runtime"
	"time"
)

func main() {
	// testA()
	testB()
}

func testA() {
	// Use only one kernel
	runtime.GOMAXPROCS(1)

	x := 0

	// Put the goroutine into the queue
	go func() {
		// There's no switch points
		for {
			x++
		}
	}()

	// Switch context to the main() goroutine. After start we will wind up here and then the second goroutine starts.
	time.Sleep(500 * time.Millisecond)

	fmt.Println(x)
}

func testB() {
	// Use only one kernel
	runtime.GOMAXPROCS(1)

	x := 0

	// Put the goroutine into the queue
	go func() {
		for {
			// Explicit points the compiler to arrange a context-switch point
			runtime.Gosched()
			x++
		}
	}()

	// Switch context to the main() goroutine. After start we will wind up here and then the second goroutine starts.
	time.Sleep(500 * time.Millisecond)

	fmt.Println(x)
}
