package main

/*
	GO 1.5 RELEASE NOTES
*/

/*
	INTRO
	The latest Go release, version 1.5, is a significant release, including major architectual changes to the implementation. Despite that, we expect almost all Go programs to continue compile and run as before, because the release still maintains the Go 1 promise compatibility.

	The biggest developments in the implementations are:
		- The compiler and runtime are now written entirely in Go (with a little assembler). C is no longer involved in the implementation, and so the C compiler that was once necessary for building the distribution is gone.
		- The GC is now concurrent and provides dramatically lower pause times by running, when possible, in parallel with other goroutines.
		- By default, Go programs with GOMAXPROCS set to the number of cores available; in prior releases it defaulted to 1.
		- Support for internal packages is now provided for all repositories, not just the Go source.
		- The `go` command now provides experimental support for `vendoring` external dependencies.
		- A new `go tool trace` command supports fine-grained of program execution.
		- A new `go doc` command is customised for command-line use.

	The release also contains one small language change involving map literals.
*/

/*
	CHANGES TO THE LANGUAGE
*/

/*
	1. Map literals
	Due to an oversight, the rule that allowed the element type to be elided from slice literals was not applied to map keys. This has been corrected in Go 1.5. An example will make this clear.

		m := map[Point]string{
			Point{29.935523, 52.891566}:   "Persepolis",
			Point{-25.352594, 131.034361}: "Uluru",
			Point{37.422455, -122.084306}: "Googleplex",
		}

	May be written as follows, without the `Point` type listed explicitly:
		m := map[Point]string{
			{29.935523, 52.891566}:   "Persepolis",
			{-25.352594, 131.034361}: "Uluru",
			{37.422455, -122.084306}: "Googleplex",
		}

*/

func mapTypeEliding() {

	type Point struct {
		x, y float64
	}

	m := map[Point]string{
		Point{29.935523, 52.891566}:   "Persepolis", // redundant type from array, slice, or map composite literalsimplifycompositelitdefault
		Point{-25.352594, 131.034361}: "Uluru",      // redundant type from array, slice, or map composite literalsimplifycompositelitdefault
		Point{37.422455, -122.084306}: "Googleplex", // redundant type from array, slice, or map composite literalsimplifycompositelitdefault

	}

	m = map[Point]string{
		{29.935523, 52.891566}:   "Persepolis",
		{-25.352594, 131.034361}: "Uluru",
		{37.422455, -122.084306}: "Googleplex",
	}

	_ = m
}

/*
	THE IMPLEMENTATION
*/

/*
	1. No more C
	The compiler and runtime are now implemented in Go and assembler, without C. The only C source left in the tree is related to testing or to `cgo`. There was a C compiler in the tree in 1.4 and earlier. It was used to build the runtime a custom compiler was necessary in part to guarantee the C code would work with the stack management and goroutines. Since the runtime is in Go now, there's no need for this C compiler and it's gone.

	The conversion from C was done with the help of custom tools created for the job. Most important, the compiler was actually moved by automatic translation of the C code into Go. It's in effect the same program in a different language. It's not a new implementation of the compiler so we expect the process won't have introduced new compiler bugs.
*/

/*
	2. Garbage Collector
	The GC has been re-engineered for 1.5 as part of the development outlined in the design document (https://docs.google.com/document/d/16Y4IsnNRCN43Mx0NZc5YXZLovrHvvLhK_h0KN8woTO4/edit?tab=t.0#heading=h.o8eay7ieosat). Expected latencies are much lower than with the collector in prior releases, through a combination of advanced algorithms, better scheduling of the collector, and running more of the collection in parallel with the user program. The `stop the world` phase of the collector will almost always under 10 milliseconds and usually much less.

	For systems that benefit from low latency, such as user-responsive web sites, the drop in expected latency with the new GC may be important.

	Details of the new collector were presented in a talk at GopherCon 2015 - https://go.dev/talks/2015/go-gc.pdf
*/

/*
	3. Runtime
	In Go 1.5, the order in which goroutines are scheduled has been changed. The properties of the scheduler were never defined by the language, but programs that depend on the scheduling order may be broken by this change. We have seen a few (erroneous) programs affected by this change. If you have programs that implicitly depend on the scheduling order, you will need to update them.

	Another potentially breaking change is that the runtime now sets the default number of threads to run sumiltaneously, defined GOMAXPROCS, to the number of cores available on the CPU. In prior releases the default was 1. Programs that don't expect to run with multiple cores may break inadvertently. They can be updated by removing the restriction or by setting GOMAXPROCS explicitly. For a more detailed discussion of this change, see the design document (https://go.googlesource.com/proposal/+/master/design/go15bootstrap.md)
*/

/*
	TOOLS
*/

/*
	1. Translating
	As part of the process to eliminate C from the tree, the compiler and linker were translared from C to Go. It was a genuine translation, so the new programs are essentially the old programs translated rather than new ones with new bugs. We are confident the translation process has introduced few if any new bugs, and in fact uncovered a number of previously unknown bugs, now fixed.

	The assembler is a new program, however; it's described below.
*/

/*
	2. Compiler
	The 1.5 compiler is mostly equivalent to the old, but some internal details have changed. One significant change is that evaluation of constants now uses the `math/big` package rather than a custom (and less well tested) implementation of high precision arithmetic. We don't expect this to affect the results.
*/

/*
	3. Assembler
	The new assembler is very nearly compatible with the previous ones, but there are a few changes that may affect some assembler source files.

	First, the expression evaluation used for constants is a little different. It now uses unsigned 64-bit arithmetic and the precedence of operations (+, -, <<, etc.) comes from Go, not C. We expect these changes to affect very few programs but manual verification may be required.
*/

/*
	4. Linker
	The linker in Go 1.5 is now one Go program, that replaces 6l, 8l, etc. Its operating system and instruction seet are specified by the environment variables GOOS and GOARCH.
 */

/*
	PERFORMANCE
	AS always the changes are so general and varied that precise statements about performance are difficult to make. The changes are even broader ranging than usual in this release, which includes a new GC and a conversion of the runtime of Go. Some programs may run faster, sone slower. On average the programs in the Go 1 benchmark suite run a few percent faster in Go 1.5 than they did in Go 1.4, while as mentioned above the GC's pauses are dramatically shorter, and almost always under 10 milliseconds.

	Builds in Go 1.5 will be slower by a factor of about two. The automatic translation of the compiler and linker from C to Go resulted in unidiomatic Go code that performs poorly compared to well-written Go. Analysis tools and refactoring helped to improve the code, but much remains to done. Further profiling and optimization will continue in Go 1.6 and future releases.
*/
