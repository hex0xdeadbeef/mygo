package main

/*
	GO 1.2 RELEASE NOTES
*/

/*
	CHANGES TO THE LANG
*/

/*
	1. Use of `nil`

	For example in Go 1.0:

		type T struct {
			X [1<<24]byte
			Field int32
		}

		func main() {
			var t *T
		}

	The `nil` pointer `t` could be used to access memory incorrectly: the expression `x.Field` could access memory at address 1<<24. To prevent such unsafe pointer, in Go 1.2 the compilers now guarantee that any indirection through a `nil` pointer, such as illustrated here but also in `nil` pointers to array, `nil` interface values, `nil` slices, and so on, will either panic or return a correct, safe non-nil value. In short, any expression that explicitly or implicitly requires evaluation of a `nil` address is an error. The implementation may inject extra tests into the compiled program to enforce this behavior (https://go.googlesource.com/proposal/+/refs/heads/master/design/4238-go12nil.md)

	UPD: Most code that depends on the old behavior is erroneous and will fail when run. Such programs will need to be updated by hand.
*/

func dereferencingUsage() {
	type (
		T2 struct {
			X     [1 << 24]byte
			Field int32
		}

		T struct {
			Field1 int32
			Field2 int32
		}
	)

	var x *T
	p1 := &x.Field1 // nil dereference in field selectionnilnessnilderef
	p2 := &x.Field2 // nil dereference in field selectionnilnessnilderef

	var x2 *T2
	p3 := &x2.X // nil dereference in field selectionnilnessnilderef

	_, _, _ = p1, p2, p3
}

/*
	2. Three index slices
	Go 1.2 adds the ability to specify the capacity as well as the length when using a slicing operation on an existing array or slice. A slicing operation creates a new slice by describing a contiguous section of an already-created array or slice:

		var array [10]int
		slice := array[2:4]

	The capacity of the slice is the maximum number of elems that the slice may hold, even after reslicing; it reflects the size of the underlying array. In this example, the capacity of the `slice` variable is 8.

	Go 1.2 adds a new syntax to allow a slicing op to specify the capacity as well as the length. A second colon introduces the capacity value, which must be less than or equal to the capacity of the source slice or array, adjusted for the origin. For instance:

		slice = array[2:4:7]

	sets the slice to have the same length as in the earlier example but its capacity is now only 5 elems (7-2). It is impossible to use this new slice value to access the last three elems of the original array.

	In this three-index notation, a missing first index ([:i:j]) defaults to zero but the other two indices must always be specified explicitly. It's possible that future of Go may introduce default values for these indices.

	UPD: This is a backwards-compatible change that affects no existing.
*/

/*
	CHANGES TO THE IMPLEMENTATION AND TOOLS
*/

/*
	1. Pre-emption in the scheduler
	In prior releases, a goroutine that was looping forever could starve out other goroutines on the same thread, a serious problem when GOMAXPROCS provided only one user thread. In Go 1.2, this is partially addressed: The scheduler is invoked occasionally upon entry to a function.
	This means that any loop that includes a (non-inlined) function call can be preempted, allowing other goroutines to run on the same thread.
*/

/*
	2. Limit on the number of threads.
	Go 1.2 introduces a configurable limit (default 10,000) to the total number of threads a single program may have in its address space, to avoid resource starvation issues in some environments. Note that goroutines are multiplexed onto threads (G : M) so this limit doesn't directly limit the number of goroutines, only the number that may be simultaneously blocked in a syscall. In practice, the limit is hard to reach.

	The new SetMaxThreads function in the runtime/debug package controls the thread count limit.

	UPD: Few functions will be affected by the limit, but if a program dies because it hits the limit, it could be modified to call SetMaxThreads to set a higher count. Even better would be to refactor the program to need fewer threads, reducing consumption of kernel resources.
*/

/*
	3. Stack size
	In Go 1.2, the minimum size of the stack when a goroutine is created has been lifted from 4 KB to 8 KB. Many programs were suffering performance problems with the old size, which had a tendency to introduce expensive stack-segment switching in performance-critical sections. The new number was determined by empirical testing.

	At the other end, the new function SetMaxStack in the runtime/debug package controls the maximum size of a single goroutine's stack. The default of 1 GB on 64-bit systems and 250 MB on 32-bit systems. Before Go 1.2, it was easy for a runaway recursion to consume all the memory on a machine.

	UPD: The increased minimum stack size may cause programs with many goroutines to use more memory. There's no workaround, but plans for future releases to include new stack management technology that should address the problem better.
*/

/*
	4. Test Coverage
	One major new feature of `go test` is that it can now compute and, with help from a new, separately installed `go tool cover` program, display test coverage results.

	The cover tool is part of the go.tools subrepository. it can be installed by running:

		go get code.google.com/p/go.tools/cmd/cover

	The cover tool does two things. First, when `go test` is given the `-cover` flag, it's run automatically to rewrite the source for the package and insert instrumentation. The test is then compiled and run as usual, and basic coverage statistics are reported.

	Second, for more detailed reports, different flags to `go test` can create a coverage profile file, which the cover program, invoked with `go tool cover`, can then analyze.
*/

/*
	PERFORMANCE
*/
