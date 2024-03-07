package chapter13

/*------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
13. LOW-LEVEL PROGRAMMING

1. There's no way to discover the memory layout of an aggregate type like a struct/the machine code for a function/the identity of the operation system thread on which the current goroutine
is running. Indeed, the Go scheduler freely moves goroutines from one thread to another. A pointer identifies a variable without revealing the variable's numeric address. Addresses may change as
the garbage collector moves variables; pointers are transparently updated.

Together, these features make Go programs, especially failing ones, more predictable and less mysterious than programs in C, the quintessential low-level language. By hiding the underlying details
they also make Go programs highly portable, since the language semantics are largely independent of any particular compiler/operating system/CPU architecture. (Not entirely indepednent: some
details leaks through, such as the word size of the processor/the order of evaluation of certain expressions/the set of implementation restrictions imposed by the compiler).

Occasionally, we may choose to forfeit some of these helpful guarantees to achieve the highest possible performance, to interoperate with libraries written in other languages, or to implement a
function that cannot be expressed in pure Go.

2. Use of "unsafe" also voids Go's warranty of compatibility with future releases, since, whether intended or inadvertent, it's easy to depend on unspecified implementation details that may change
unexpectedly.

3. "unsafe" provides access to a number of built-in language features that are not ordinarily available because they expose details of Go's memory layout.

------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------*/

/*------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
13.1 unsafe.Sizeof / unsafe.Alignof / unsafe.Offsetof

1. The "unsafe.Sizeof" function reports the size in bytes of the representation of its operand, which may be of any type. The expression isn't evaluated. A call to "unsafe.Sizeof" is a constant
expression of type uintptr, so the result may be used as the dimension of an array type, or to compute other constants.

2. For portability, we've given the sizes of reference types (or types containing references) in terms of words, where a word is 4 bytes on a 32-bit platform and 8 bytes on a 64-bit
platform.

3. Computers load and store values from memory most efficiently when those are properly aligned. For example:
	1) The address of a value of a two-byte type such as int16 should be an even number.
	2) The address of a four-byte value such as a rune should be a multiple of four
	3) The address of an eight-byte value such as int64, uint64, or 64-bit pointer should be a multiple of eight.
Alignment requirements of higher multiplies are unusual, even for larger data types such as complex128. For this reason, the size of an aggregate type (a struct/array) is at least the sum of the
sizes of its fields or elements but may be greater due to the presence of "holes" Holes are unused spaces added by the compiler to ensure that the following field or element is properly aligned
relative to the start of the struct or array.

	Type							Size
	bool							1 byte
	intN,..., complexN				N/8 bytes
	int, uint, uintptr				1 word (4/8 bytes)
	*T								1 word (4/8 bytes)
	string							2 words (data, len) (8/16 bytes)
	[]T								3 words	(len, cap, data) (12/24 bytes)
	map								1 word (4/8 bytes)
	func							1 word (4/8 bytes)
	chan 							1 word (4/8 bytes)
	interface						2 words (type, value) (8/16 bytes)

4. The language specification doesn't guarantee that the order in which fields are declared is the order in which they are laid out in memory, so in theory a compiler is free to rearrange them.

5. If the types of a struct's fields are of different sizes, it may be more space-efficient to declare the fields in an order that packs them as tightly as possible. The three structs below have the
same fields, but the first requires up to 50% more memory than the other two, so we should arrange the field in the ASC/DSC orders.
								64-bit			32-bit
struct {bool; floa64; int16} 	3 words			4 words
struct {bool, int16, float64}	2 words			3 words ASC
struct {float64, int16, bool}	2 words			3 words DESC

Efficient packing may make frequently allocated data structures more compact and therefore faster.

6. The "unsafe.Alignof()" function reports the required alignment of its argument's type. Like "unsafe.Sizeof", it may be applied to an expression of any type, and it yields a constant. Typically,
boolean and numeric types are aligned to their size (up to a maximum of 8 bytes) and all other types are word-aligned.

7. The "unsafe.Offsetof" function, whose operand must be a field selector "x.f" computes the offset of field "f" relative to the start of its enclosing struct "x", accounting for holes, if any.

------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------*/

/*------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
13.2 unsafe.Pointer
1. The unsafe.Pointer type is a special kind of pointer that can hold the address of any variable. Of course, we can't indirect through an "unsafe.Pointer" using *p because we don't know what type
that expression should have. Like ordinary pointers, "unsafe.Pointer"s are comparable and may be compared with nil, which is the zero of the type.

2. An ordinary pointer "*T" may be converted to an "unsafe.Pointer", and an "unsafe.Pointer" may be converted back to an ordinary pointer, not necessarily of the same type *T.

3. In general, "unsafe.Pointer" conversions let us write arbitrary values to memory and thus subvert the type system.

4. An "unsafe.Pointer" may also be converted to a "uintptr" that holds the pointer's numeric value, letting us perform arithmetic on addresses.

5. Many "unsafe.Pointer" values are thus intermediaries for converting ordinary pointers to raw numeric addresses and back again.

6. We mustn't introduce temporary variables of type "uintptr" to break the lines. The reason is very subtle.
	Some garbage collectors move variables around in memory to reduce fragmentation or bookkeeping. Garbage collectors of this kind are known as "moving GCs". When a variable is moved, all pointers that hold the address of the old location must be updated to point to the new one. From the perspective of the garbage collector, an "unsafe.Pointer" is a pointer and thus its value must change as the variables moves, but a uintptr is just a number so its value must not change. The incorrect code of "SubtleVariableIntroducing()" hides a pointer from the garbage collector in the non-pointer variable tmp. By the time the second statement executes, the variable x could have moved and the number in "tmp" would no longer be the address &x.b. The third statement clobbers an arbitrary memory location with the value 42.

7. After this statement has executed:
	pT := uintptr(unsafe.Pointer(new(T))), when T is an arbitrary type.

there are no pointers that refer to the variable created by new, so the GC is entitled to recycle when this statement completes, after which pT contains the address where the variable was but is no
longer.

8. No current Go implementation uses a MGC, but this is no reason for complacency: current version of Go do move some varaibles around in memory. When goroutine stack grows, all variables on the
old stack may be relocated to a new, larger stack, so we cannot rely on the numeric value of a variable's address remaining unchanged throughout its lifetime.

9. We MUST treat all "uintptr" values as if they contain the former address, and minimize the number of operations between converting an "unsafe.Pointer" to a "uintptr" and using that "uintptr".

10. When calling a library function that returns a "uintptr", the result MUST be immediately converted to an "unsafe.Pointer" to ensure that it continues to point to the same variable.

------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------*/

/*------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
13.3 DEEP EQUIVALENCE

1. The "DeepEquivalence" function from the "reflect" package reports whether two values are "deeply" equal.
	1) "DeepEqual()" compares basic values as if by the built-in "==" operator
	2) For composite values, it traverses them recursively, compairing corresponding elements.

Because it works for any pair of values, even once that are not comparable with "==", it finds widespread use in tests.

2. "DeepEqual()" doesn't consider a nil map equal to a non-nil empty map, nor a slice equal to a non-nil empty one.

3. WE CAN PUT THE panic("...") SO THAT SUPPRESS RETURN STATEMENT

4. We need the "reflect.Type" in the struct "comparison" because different variables can have the same address. For example: if x and y both are arrays, and x[0] and x have the same address,
as do y and y[0], and it's important to distinguish whether we have compared x and y or x[0] and y[0].
------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------*/

/*------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
13.4 CALLING C CODE with cgo

1. "cgo" is a tool that creates Go bindings fo C functions. Such tool are called "foreign-function interfaces" (FFIs).

2. "cgo" is not only fo Go programs. SWIG is another. It provides more complex features for integrating with C++ classes.

3. "compress/..." subtree of the standart library provides compressors and decompressors for popular compression algorithms.

4. From the "libbzip2" C package we need:
	1) "bz_stream" struct type, which holds the input and output buffers
	AND THREE C FUNCTIONS
	1) BZ2_bzCompressInit, which allocates the stream's buffers
	2) BZ2_bzCompress, which compresses data from the input buffer to the output buffer
	3) BZ2_bzCompressEnd, which releases the buffers

We'll call the BZ2_bzCompressInit and BZ2_bzCompressEnd C functions directly from Go, but for BZ2_bzCompress, we'll define a wrapper function in C.

5. The import "C" is the special import. There's no package C, but this import causes "go build" to preprocess the file using the "cgo" tool before the Go compiler sees it.
	During preprocessing "cgo" generates a temporary package that contains Go declarations corresponding to all the C functions and types used by the file, such as C.bz_stream and
	c.BZ2_bzCompressInit. The "cgo" tool discovers these types by invoking the C compiler in a special way on the contents of the comment that precedes the import declaration.

	The comment may also contain "#cgo" directives that specify extra options to the C toolchain. The "CFLAGS" and "LDFLAGS" contribute extra arguments to the compiler and linker commands
	so that they can locate the "bzlib.h" header file and the "libbz2.a" archive library.

6. We should observe that Go program may access C types like "bz_stream", "char", and "uint", C functions like bz2compress, and even object-like C preprocessor macross such "BZ_RUN", all
through the "C.x" notation. The "C.uint" type is distinct from Go's "uint" type, even if both have the same width.
------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------*/

/*------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
13.4 CALLING C CODE with cgo

1. "cgo" is a tool that creates Go bindings fo C functions. Such tool are called "foreign-function interfaces" (FFIs).

2. "cgo" is not only fo Go programs. SWIG is another. It provides more complex features for integrating with C++ classes.

3. "compress/..." subtree of the standart library provides compressors and decompressors for popular compression algorithms.

4. From the "libbzip2" C package we need:
	1) "bz_stream" struct type, which holds the input and output buffers
	AND THREE C FUNCTIONS
	1) BZ2_bzCompressInit, which allocates the stream's buffers
	2) BZ2_bzCompress, which compresses data from the input buffer to the output buffer
	3) BZ2_bzCompressEnd, which releases the buffers

We'll call the BZ2_bzCompressInit and BZ2_bzCompressEnd C functions directly from Go, but for BZ2_bzCompress, we'll define a wrapper function in C.

5. The import "C" is the special import. There's no package C, but this import causes "go build" to preprocess the file using the "cgo" tool before the Go compiler sees it.
	During preprocessing "cgo" generates a temporary package that contains Go declarations corresponding to all the C functions and types used by the file, such as C.bz_stream and
	c.BZ2_bzCompressInit. The "cgo" tool discovers these types by invoking the C compiler in a special way on the contents of the comment that precedes the import declaration.

	The comment may also contain "#cgo" directives that specify extra options to the C toolchain. The "CFLAGS" and "LDFLAGS" contribute extra arguments to the compiler and linker commands
	so that they can locate the "bzlib.h" header file and the "libbz2.a" archive library.

6. We should observe that Go program may access C types like "bz_stream", "char", and "uint", C functions like bz2compress, and even object-like C preprocessor macross such "BZ_RUN", all
through the "C.x" notation. The "C.uint" type is distinct from Go's "uint" type, even if both have the same width.
------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------*/
