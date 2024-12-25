package main

/*
	GO 1.3 RELEASE NOTES
*/

/*
	INTRO
Version 1.3 arrives six months later after 1.2 and contains no language changes. It focuses primarily on implementation work, providing precise garbage collection, a major refactoring of the compiler toolchain that results in faster builds, especially for large projects, significant performance improvements across the board.
*/

/*
	CHANGES TO THE MEMORY MODEL
	The Go 1.3 memory model adds a new rule concerning sending and receiving on buffered channels, to make explicit that a buffered channel can be used as a simple semaphore, using a send into the channel to acquire and a receive from the channel to release. This isn't a language change, just a clarification about an expected property of communication.
*/

/*
	CHANGES TO THE IMPLEMENTATION AND TOOLS
*/

/*
	1. Stack
	Go 1.3 has changed the implementation of goroutine stacks away from the `old`, `segmented` model to a contiguous model. When a goroutine needs more stack than is avaiable, its stack is transfered to a larger single block of memory. The overhead of this transfer amortizes well and eliminates the old `hot spot` problem when a calculation repeatedly steps across a segment boundary. Details including performance numbers are in this design document (https://tip.golang.org/s/contigstacks)
*/

/*
	2. Changes to the Garbage Collector
	For a while now, the GC has been precise when examining values in the heap. The Go 1.3 release adds equivalent precision to values on the stack. This means that a non-pointer Go value such as an integer will never be mistaken for a pointer and prevent unused memory from being reclaimed.

	Starting with Go 1.3 the runtime assumes that values with pointer type contain pointers and other values do not. This assumption is fundamental to the precise behavior of both stack expansion and garbage collection.

	Programs that use package `unsafe` to store pointers in integer-typed vals are also illegal but more difficult to diagnose during execution. Because the pointers are hidden from the runtime, a stack expansion or garbage collection may reclaim the memory point at, creating dangling pointers.

	UPD: Code that uses unsafe.Pointer to convert an integer-typed value held in memory into a pointer is illegal and must be rewritten. Such code can be identified by `go vet`
*/


/*
	PERFORMANCE
	The performance of Go binaries for this release has improved in many cases due to changes in the runtime and garbage collection, plus some changes to libraries. Significant instances include:
		- The runtime handles defers more efficiently, reducing the memory footprint by about 2 KB per goroutine that calls defer.
		- The GC has been sped up, using a concurrent sweeep algorithm, better parallelization, and larger pages. The cumulative effect can be 50-70% reduction in collector pause time.
		- The race detector is now 40% faster
		- The regular expression package `regexp` is now signficant faster for certain simple expressions due to the implementation of a second, on-pass execution engine, The choice of which engine to use is automatic.
	Also, the runtime now includes in stack dumps how long a gorotuine has been blocked, which can be useful information when debugging deadlocks or performance issues.
*/