package main

/*
	GO 1.4 RELEASE NOTES
*/

/*
	CHANGES TO THE LANGUAGE
*/

/*
	1. For-range loops

	Up until Go 1.3, `for-range` loop had two forms:

		for i, v := range x {
			...
		}

		for i := range x {
			...
		}

	If one wasn't interested in the loop values, only the iteration itself, it was still necessary to mention a variable (probably the blank identifier, as in for _ = range x), because the form

		for range x {
			...
		}

	wasn't syntactically permitted.

	This situation seemed awkward, so as of Go 1.4 the variable-free form is now legal. The pattern arises rarely but the code can be cleaner when it does.
*/

/*
	2. Method calls on **T
	Given these declarations:

		type T int
		func (T) M() {}
		var x **T

	both `gc` and `gccgo` accepted the method call

		x.M()

	which is a double dereference of the pointer-to-pointer `x`. The Go specification allows a single derefence to be inserted automatically, but not two, so this call is erroneous according to the language definition. It has therefore been disallowed in Go 1.4, which is a breaking change, although very few programs will be affected.

	UPD: Code that depends on the old. erroneous behavior will no longer compile but is easy to fix by adding an explicit dereference.
*/

/*
	CHANGES TO THE IMPLEMENTATIONS AND TOOLS
*/

/*
	1. Changes to the runtime
	Prior to Go 1.4, the runtime (GC, concurrency support, interface management, maps, slices, strings, ... ) was mostly written in C, with some assembler support. In 1.4 much of the code has been translated to Go so that the GC can scan the stacks of programs in the runtime and get accurate information about what variables are active. This change was large but should have no semantic effect on programs.

	This rewrite allows the GC in 1.4 to be fully precise, meaning that it's aware of the location of all active pointers in the program. This means the heap will be smaller as there will be no false positives keeping non-pointers alive. Other related changes also reduce the heap size, which is smaller by 10%-30% overall relative to the previous release.

	A consequence is that stacks are no longer segmented, eliminating the `hot split` problem. When a stack limit is reached, a new, larger stack is allocated, all active frames for the goroutine are copied there, and any pointers into the stack are updated. Performance can be noticeably better in some cases and is always more predictable.

	The use of contiguous stacks means that stacks can start smaller without triggering performance issues, so the default starting size for a goroutine's stack in 1.4 has been reduced from 8192 bytes to 2048 bytes.

	As preparation for the concurrent GC scheduled for the 1.5 release, writes to pointer values in the heap are now done by a function call, called a write barrier, rather than directly from the function updating the value. In this next release, this will permit the GC to mediate writes to the heap while it's running. This change has no semantic effect on programs in Go 1.4, but was included in the release to test the compiler and the resulting performance.

	The implementation of interface values has been modified. In earlier releases, the interface contained a word that was either a pointer or a one-word scalar value, depending on the type of the concrete object stored. This implementation was problematical for the GC, so as of 1.4 interface values always hold a pointer. In running programs, most interface values were pointers anyway, so the effect is minimal, but programs that store integers (for example) in interfaces will see more allocations.

	As of Go 1.3, the runtime crashes if it finds a memory word that should contain a valid pointer but instead contains an obviously invalid pointer (for example, the value 3). Programs that store integers in pointer values may run afoul of this check and crash. In Go 1.4, setting the CODEBUG varibale invalidptr=0 disables the crash as a workaround, but we cannot guarantee that future releases will be able to avoid the crash; the correct fix is to rewrite code not to alias integers and pointers.
*/

/*
	PERFORMANCE
	Most programs will run about the same speed or slightly faster in 1.4 than in 1.3; some will be slightly slower. There are many changes, making it hard to be precise about what to expect.

	As mentioned above, much of the runtime was translated to Go from C, which led to some reduction in heap sizes. It also improved performance slightly because the Go compiler is better at optimization, due to things like inlining, than the C compiler used to build the runtime.

	The GC was sped up, leading to measurable improvements for garbage-heavy programs. On the other hand, the new write barriers slow things down again, typically by about the same amount but, depending on their behavior, some programs may be somewhat slower or faster.
*/