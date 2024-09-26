package main

import (
	"fmt"
	"runtime"
)

/*
	STACK TRICKS
1. Every type of an object has the specific class of object and realtively the maximum size of the class the object will be allocated either on the stack or the heap.
2. The cases when an object is moved to the heap:
	1) The object is used by many goroutines
	2) The object is returned from function as a pointer
	3) The object is passed to a function as an interface parameter
3. More allocations in the stack lead to blowing the stack frame up, then reallocation of the stack gets more expansive.
4. The backdoor of this situation is to create slice directly by pointing the index and value of the last element in the slice.
*/

func SliceAllocationA() {
	var (
		s = make([]byte, (1<<6)*(1<<10))
	)

	_ = s
}

func SliceAllocationB() {
	var (
		s = make([]byte, 1<<28)
	)

	_ = s
}

func SliceAllocationC() {
	printAlloc()

	var (
		// referencing to the last index and assigning a value to it
		s = []byte{1 << 28: 1}
	)

	printAlloc()

	s = append(s, s...)

	printAlloc()

	_ = s

	runtime.KeepAlive(&s)
}

func printAlloc() {
	var (
		m runtime.MemStats
	)

	runtime.ReadMemStats(&m)
	fmt.Printf("%d MB\n", m.Alloc/(1<<20))
}
