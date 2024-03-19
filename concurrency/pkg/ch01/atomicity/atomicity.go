package atomicity

import "fmt"

// There are three ops:
// 1. Retrieve the value of i
// 2. Increment the value
// 3. Rewrite i

func AtomicExecution() {
	var i int

	i++

	fmt.Println(i)
}

func NonAtomicExecution() {
	var i int

	go func() {
		i++
	}()

	i += 2

	fmt.Println(i)
}
