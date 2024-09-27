package main

/*
GOROUTINE STACK ---------------------------------------------------------------------------------------------------------------------------------
1. The lowes size boundaries of goroutines (Gs) in different Go versions in KB:
	1) Go 1.2: 8
	2) Go 1.4: 2

2. If the Go is able to increase the size of a stack, it means that Go can define when the size of memory is being allocated must not be changed.

3. To define how Go manages memory we should use "go run -gcflags -S main.go".
	1) runtime.morestack_noctxt(...) manages the size of a stack if it needs more memory
		1) CALL used in asm code to use this method to increase the stack size when it needs it.
		2) NOSPLIT means that the control of overflow of a stack isn't needed


	2) First of all the size of the current stack is estimated counting on the fields (lo uintptr, hi uintptr) of the stack struct:
		type stack {
			lo uintptr
			hi uintptr
		}
	as a difference between hi and lo

	3) After 2) the estimated size is doubled up and the check of overflowing is performed. The max size depends on an architecture.
		1) 64-bit system - 1 GB
		2) 32-bit system - 250 MB


	4) The size of a stack is defined by its boundaries.

	5) The old stack is copied into a new memory area and the old area is freed.

	6) The instruction copystack() copies all the old stack and transmit all the addresses to the new stack.

4. While garbage collecting the size of a goroutine stack will be reduced if it's possible.
	1) To do this the function shrinkstack() is invoked

5. From the 1.3 version Go works with "contiguous" stack.
	1) "contiguous" stack is the strategy when the stack is copied into a stack that is larger than old one.

*/