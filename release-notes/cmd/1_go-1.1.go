package main

/*
	GO 1.1 RELEASE NOTES
*/

/*
	INTRO
*/

/*
	The focus of this version is to provide more efficient lang. Benchmarking is an inexact science at best, but we see significant, sometimes dramatic speedups for many of our test programs.
*/

/*
	CHANGES TO THE LANG
*/

/*
	1. Integer division by zero
	In Go 1, integer division by a constant zero produced a run-time panic.

		func f(x int) int {
			const zeroVal = 0
			return x / zeroVal // There was a panic in Go 1.0
		}

	In Go 1.1, an integer division by constant zero is not a legal program, so it's compile-time error.
*/

/*
	2. Surrogates in Unicode literals
	The definition of string and rune literals has been refined to exclude surrogate halves from the set of valid Unicode code points.
*/

/*
	3. Method values
	Go 1.1 now implements method values, which are funcs that have been bound to a specific receiver value. For instance, given a Writer value `w`, the expression w.Write, a method value, is a func that will always write to `w`; that is equivalent to a function literal closing over `w`.

		func (p []byte) (n int, err error) {
			return w.Write(p)
		}

	Method values are distinct from methods expresisons, which generate functions from methods of a given type; the method expression (*bufio.Writer).Write is equivalent to a function with an extra first argument, a receiver of type (*bufio.Writer)

		func (w *bufio.Writer, p []byte) (n int, err error) {
			return w.Write(p)
		}

	UPD: No existing code is affected the change is strictly backward-compatible.
*/

/*
	4. Return requirements
	Before 1.1 Go func that returned a value needed an explicit "return" or call to panic at the end of the func; this was a simple way to make the programmer be explicit about the meaning of the function. But there are many cases where a final "return" is clearly unnecessary, such as a func with only an infinite `for` loop.

	In Go 1.1 the rule about final `return` statements is more permissive. It introduces the concept of a terminating statement, a statement that is guaranteed to be the last one a function executes. Examples include `for` loops with no condition an `if-else` statements in which each half ends in a `return`. If the final statement of a func can be shown syntactically to be a terminating statement, no final return is needed.

	Note that the rule is purely syntactic: it pays no attention to the values in the code and therefore requires no complex analysis.

	UPD: The change is backward-compatible, but existing code with superfluous `return` statements and calls to panic may be simplified manually. Such code can be identified by `go vet`
*/

/*
	CHANGES TO THE IMPLEMENTATIONS AND TOOLS
*/

/*
	1. Status of `gccgo`
	The GCC release schedule doesn't coincide with the Go release schedule, so some skew is inevitable in gccgo's releases.
*/

/*
	2. Command-line flag parsing
	In the `gc` toolchain, the compilers and linkers now use the same command-line flag parsing rules as the Go flag package, a departure from the traditional Unix flag parsing. This may affect scripts that invoke the tool directly. For example `go tool 6c -Fw -Dfoo` must now be written `go tool 6c -F -w -D foo`
*/

/*
	3. Size of `int` on 64-bit platforms
	The lang allows the implementation to choose whether the `int` and `uint` types are 32 or 64 bits. Previous Go implementation made `int` and `uint` 32 bits on all systems. Both the `gc` and `gccgo` implementations now make `int` and `uint` 64 bits on 64-bit platforms such as AMD64/x86-64. Among other things, this enables the allocation of slices with more than 2 billion elems on 64-bit platforms.

	UPD: Most programs will be unaffected by this change. Because Go doesn't allow impicit conversions between distinct numeric types, no programs will stop compiling due to this change. However programs that contain implicit assumptions that `int` is only 32 bits may change this behavior. For example, this code prints a positive number on 64-bit and a negative on 32-bit systems:

		x := ^uint32(0) // x is 0xffffffff
		i := int(x) // i is -1 on 32-bit systems, 0xffffffff on 64 bit
		fmt.Println(i)

	Portable code intending 32-bit sign extension (yelding -1 on all systems) would instead say:

		i := int(int32(x))
*/

/*
	4. Heap size on 64-bit atchitectures
	On 64-bit architectures, the maximum heap size has been enlarged substantially, from a few gigabytes to several tens of gigabytes. (The exact details depend on the system and may change)

	On 32-bit architectures, the heap size hasn't been changed.

	UPD: This change should have no effect on existing programs beyond allowing them to run with larger heaps.
*/


/*
	5. Unicode
	To make it possible to represent code points greater than 65535 in UTF-16, Unicode defines surrogate halves, a range of code points to be used only in the assembly of large values, and only in UTF-16. The code points in that surrogate range are illegal for any other purpose. In Go 1.1, this constraint is honored by the compiler, libraries, and run-time: a surrogate half is illegal as a rune value, when encoded as UTF-8, or when encoded in isolation as UTF-16. When encountered, for example in converting from a rune to UTF-8, it's treated as an encoding error and will yield the replacement rune, utf8.RuneError, u+FFFD.

		import "fmt"

		func main() {
			fmt.Print("%+q\n", string(0xD800))
		}

	prints "\ud800" in Go 1.0, but prints "\ufffd" in Go 1.1.

	Surrogate-half Unicode values are now illegal in rune and string constants, so constants such as '\ud800' and "\ud800" are now rejected by the compilers. When written expicitly as UTF-8 encoded bytes, such strings can still be create, as in "\xed\xa0\x80". However when such a string is decoded as a sequence of runes, as in a range loop. it'll yield only utf8.RuneError vals.

	The Unicode byte order mark (BOM) U+FEFF, encoded in UTF-8, is now permitted as the first character of a Go source file. Even though its appearance in the byte-order-free UTF-8 encoding is clearly unnecessary, some editors add the mark as a kind of `magic number` identifying a UTF-8 encoded file.

	UPD: Most programs will be unaffected byt he surrogate change. Programs that depend on the old behavior should be modified to avoid the issue. The BOM change is strictly backward-compatible.
 */

 /*
	6. Race detector
	A major addition to the tools is a `race detector`, a way to find bugs in program caused by concurrent access of the same variable, where at least one of the accesses is a write. This new facility is built into the `go tool`. To enable it, set the `-race` flag when building or testing our program.
 */

 /*
	PERFORMANCE
 */

 /*
	There are too many small performance-driven tweaks through the tools and libraries to list them all here, but the following major changes are worth noting:
		1) The `gc` compilers generate better code in many cases, most noticeably for floating point on the 32-bit Intel architecture.
		2) The `gc` compilers do more in-lining, including for some ops in the run-time such as `append` and `interface conversions.`
		3) There's a new implementation of Go maps with significant reduction in memory footprint and CPU time.
		4) The GC has been made more parallel, which can reduce latencies for programs running on multiple CPUs
		5) The GC is also more precise, which costs a small amount of CPU time but can reduce the size of the heap significantly, especially on 32-bit architectures.
		6) Due to tighter coupling of the run-time and network libraries, fewer context switches are required on network ops.
 */