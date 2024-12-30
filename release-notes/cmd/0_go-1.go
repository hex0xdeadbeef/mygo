package main

import (
	"fmt"
	"log"
	"unicode"
)

/*
	GO 1 RELEASE NOTES
*/

/*
	CHANGES TO THE LANGUAGE
*/

/*
1. Append
	The append predeclared variadic function makes it easy to grow a slice by adding elements to the end. A common use is to add bytes to the end of a byte slice when generating output. However, append didn't provide a way to append a `string` to a `[]byte`, which is another common case.

		greeting := []byte{}
		greeting = append(greeting, []byte("hello")...)

	By analogy, with the similar property, of `copy`, Go 1 permits a string to be appended (byte-wise) directly to a byte slice, reduces the friction between strings and byte slices. The conversion is no longer necessary:

		greeting := []byte{}
		greeting = append(greeting, "hello"...)
*/

func directStringAppendToByteSlice() {
	var greeting []byte
	// We can append a string directly to a byte slice without any conversion.
	greeting = append(greeting, "hello"...)
	_ = greeting
}

/*
	2. Close
	The `close` predeclared function provides a mechanism for a sender to signal that no more values will be sent. It's important to the implementation of `for-range` loops over channels and is helpful in other situations. Partly by design and partly because of race conditions that can occur otherwise, it's intended for use only by the goroutine sending on the channel, not by the goroutine receiving data. However, before Go 1 there was no compile-time checking that `close` was being used correctly.

	To close this gap, at least in part, Go 1 disallows `close` on `receive-only` channels. Attempting to close such a channel is a `compile-time` error.

	UPD: Existing code that attempts to close a receive-only channel was erroneous even before Go 1 and should be fixed. The compiler won't reject such code.
*/

func tryToCloseReceiveOnlyChannel() {
	var c chan int

	// Channel-sender
	var cSend chan<- int = c
	// Channel-receiver
	var cRecv <-chan int = c

	_ = cRecv
	close(c)     // Legal
	close(cSend) // Legal, but panics because the channel has been closed earlier
	// close(cRecv) // invalid operation: cannot close receive-only channel cRecv (variable of type <-chan int)compilerInvalidClose
}

/*
	3. Composite literals
	In Go 1, a composite literal of array, slice, or map type can elide the type specification for the elements' initializers if they are of pointer type. All four of the initializations in this example are legal; the last one was illegal before Go 1.

	UPD: This change has no effect on existing code, but the commad `gofmt -s` applied to existing source will, among other thing, elide explicit types wherever permitted.
*/

type Date struct {
	month string
	day   int
}

// Struct values, fully qualified; always legal.
var holiday1 = []Date{
	Date{"Feb", 14},
	Date{"Nov", 11},
	Date{"Dec", 25},
}

// Struct values, type name elided; always legal
var holiday2 = []Date{
	{"Feb", 14},
	{"Nov", 11},
	{"Dec", 25},
}

// Pointers, fully qualified, always legal
var holiday3 = []*Date{
	&Date{"Feb", 14},
	&Date{"Nov", 11},
	&Date{"Dec", 25},
}

// Pointers, type name elided; legal in Go 1.
// Before Go 1, this initialization was illegal.
var holiday4 = []*Date{
	{"Feb", 14},
	{"Nov", 11},
	{"Dec", 25},
}

/*
	4. Goroutines during init function.
	The old language defined that `go` statements executed during initializations created goroutines but that they didn't begin to run until initialization of the entire program was complete. This introduced clumsiness in many places and, in effect, limited the utility of the `init` construct:

		If it was possible for another package to use the library during initialization, the library was forced to avoid goroutines. This design was done for reasons of simplicity and safety, but, as our confidence in the language grew, it seemed unnecessary. Running goroutines during initializations is no more complex or unsafe than running them during normal execution.

	In Go 1, code that uses goroutines can be called from `init` routines and global initialization exressions without introducing a deadlock.

	UPD: This is a new feature, so existing code needs no changes, although it's possible that code that depends on goroutines not starting before main will break. There was no such code in the `std`
*/

var PackageGlobal int

func init() {
	c := make(chan int)
	go initPackageGlobal(c)
}

func initPackageGlobal(c chan int) {
	c <- 42
}

/*
	5. The `rune` type
	The language spec allows the int type to be int32 or 64 bits wide, but current implementations set `int` to 32 bits even on 64-bit platforms. It would be preferable to have `int` be 64 bits on 64-bit platforms. (There are important consequences for indexing large slices.) However, this change would waste space when processing Unicode characters with the old language because the `int` type was also used to hold Unicode code points: each code point would waste an extra 32 bits of storage if `int` grew from 32 bits to 64.

	To make changing to 64-bit `int` feasible, Go 1 introduces a new basic type, `rune`, to represent individual Unicode code points. It's an alias for `int32`, analogous to byte as an alias for `uint8`.

	Character literals such as 'a', '語', '\u0345', now have default type `rune`, analogous to 1.0 having default type `float64`, A variable initialized to a character constant will therefore have type `rune` unless otherwise specified.

	Libraries have been updated to use `rune` rather than `int` when appropriate. For instance, the functions `unicode.ToLower` and relatives now take and return a `rune`

	Updating: Most source code will be unaffected by this because the type inference from `:=` initializers introduces the new type silently, and it propagates from there. Some code may get type errors that a trivial conversions will resolve.
*/

func runeDefaultType() {
	delta := 'δ' // delta has type `rune`
	var DELTA = unicode.ToUpper(delta)
	_ = DELTA
	epsilon := unicode.ToLower(DELTA + 1)
	if epsilon != 'δ'+1 {
		log.Fatal("inconsistent casing for Geek")
	}
}

/*
	6. The `error` type.
	Go 1 introduces a new built-in type, `error`, which has the following definition:

		type error interface {
			Error() string
		}

	Since the consequences of this type are all in the package library, it's discussed below.
*/

/*
	7. Deleting from maps.
	In the old language, to delete the entry with key `k` from map `m`, one wrote the statement:

		m[k] = value, false

	This syntax was a peculiar special case, the only two-to-one assignment. It required passing a value (usually ignored) `that is evaluated but discarded`, plus a boolean that was nearly always the constant `false`. It did the job, but was odd and a point of contention.

	In Go 1, that syntax has gone; instead there's a new built-in function `delete`. The call:

		delete(m, k)

	will delete the map entry retrieved by the expression m[k]. There's no return value. Deleting a non-existent entry is no-op.

	UPD: Running `go fix` will convert expressions `m[k] = value, false` into delete(m, k) when it's clear that the ignored value can be safely discarded from the program and false refers to the predefined boolean constant. The fix tool will flag other uses of the syntax for inspection by the programmer.

*/

/*
	8. Iterating in maps
	The old language spec didn't define the order of iteration for maps, and in practice it differed across hardware platforms. This caused tests that iterated over maps to be fragile and non-portable, with the unpleasant property that a test might always pass on one machine but break on another.

	In Go 1, the order in which elems are visited when iterating over a map using `for-range` statement is defined to be unpredictable, even if the same loop is run multiple times with the same map. Code shouldn't assume that the elems are visited in any particular order.

	This change means that code that depends on iteration order is very likely to break early and be fixed long before it becomes a problem. Just as important, it allows the map implementation to ensure better map balancing even when programs are using range loops to select an element from a map.

	UPD: This is one change where tools cannot help. Most existing code will be unaffected, but some programs may break or misbehave; we recommend manual checking of all range statements over maps to verify they don't depend on iteration order. There were a few such examples in the `std` repository; they have been fixed. Note that it was already incorrect to depend on the iteration order, which was unspecified. This change codifies the unpredictability.
*/

func iteratingOverMaps() {
	m := map[string]int{
		"Sunday": 0,
		"Monday": 1,
	}

	// This loop shouldn't assume Sunday will be visited first.
	for name, value := range m {
		fmt.Println(name, value)
	}
}

/*
	9. Multiple assignment
	The language spec. has long guaranteed that in assignments the right-hand-side expressions are all evaluated before any left-hand-side are assigned. To guarantee predictable behavior, Go 1 refines the spec. further.

	If the left-hand side of the assignment contains expressions that require evaluation, such as function calls or array indexing ops, these will all be done using the usual left-to-right rule before any variables are assigned their value. Once `everything` is evaluated, the actual assignments proceed in left-to-right order.

	UPD: This is one change where tools cannot help, but breakage is unlikely. No code in the std
	repository was broken by this change, and code that dependeds on the previous unspecified behavior was already incorrect.
*/

func multipleAssignments() {
	sa := []int{1, 2, 3}
	i := 0
	i, sa[i] = 1, 2 // sets i = 1, sa[0] = 2 (references to the old value of i)
	fmt.Println(sa)

	sb := []int{1, 2, 3}
	j := 0
	sb[j], j = 2, 1 // sb[0] = 2 (references to the old value of j) , j = 1
	fmt.Println(sb)

	sc := []int{1, 2, 3}
	sc[0], sc[0] = 1, 2 // sets sc[0] = 1, then sc[0] = 2 (so sc[0] = 2 at the end)
	fmt.Println(sc)
}

/*
	10. Returns and shadowed variables
	A common mistake is to use return (without arguments) after an assignment to a variable that has the same name as a result variale but is not the same variable. This situation is called `shadowing`: the result variable has been shadowed by another variable with the same name in an inner scope.

	In funcs with named return vals, the Go 1 compilers disallow return statements `without arguments` if any of the named return values is shadowed at the point of the return statement. (It's not part of the spec., because this is one area we're still exploring; the situation is analogous to the compilers rejecting funcs that don't end with an explicit return statement.)

	This func implicitly returns a shadowed return value and will be rejected by the compiler:

		func Bug() (i, j, k int) {
			for i = 0; i < 5; i++ {
				for j := 0; j < 5; j++ { // Redeclares j.
					k += i * j

					if k > 100 {
						return // Rejected j is shadowed here
					}
				}
			}

			return // OK; j isn't shadowed here.
		}

	UPD: Code that shadows return values in this way will be rejected by the compiler and will need to be fixed by hand. The few cases that arose in the `std` were mostly bugs.
*/

/*
	11. Copying structs with unexported fields.
	The old lang didn't allow a package to make a copy of a struct value containing unexported fields belonging to a different package. There was, however, a required exception for a method receiver; also, the implementations of `copy(dst, src)` and `append(dst, elems...)` have never honored this restriction.

	Go 1 will allow packages to copy struct values containing unexported fields from other packages. Besides resolving the inconsistency, this change admits a new kind of API; a package can return an opaque value without resorting to a pointer or interface. The new implementations of time.Time and reflect.Value are examples of types taking advantages of this new property.

	As an example, if package `p` includes the definitions,

		package p

		type Struct struct {
			Public int
			secret int
		}

		func NewStruct(a int) Struct { // Note: not a pointer.
			return Struct{a, f(a)}
		}

		func (s Struct) String() string {
			return fmt.Sprintf("{%d (secret %d)}", s.Public, s.secret)
		}

	a package that imports `p` can assign and copy vals of type p.Struct at will. Behind the scenes the unexported fields will be assigned and copied just as if they were unexported, but the client code will never be aware of them. The code:

		package my

		import "p"

		myStruct := p.NewStruct(23)
		copyOfMyStruct := myStruct
		fmt.Println(myStruct, copyOfMyStruct)

	will show that the secret field of the struct has been copied to the new value.

	UPD: This is a new feature, so existing code needs no changes.
*/

/*
	12. Equality
	Before Go 1, the lang didn't define equality on `struct` and `array` vals. This meant, among other things, that structs and arrays couldn't be used as map keys. On the other hand, Go did define equality on `function` and `map` vals. Function equality was problematic in the presence of closures (when are two closures equal?) while map equality compared pointers, not the maps' content, which was usually not what the user would want.

	Go 1 addressed these issues. First, structs and arrays can be compared for equality and inequality (== and !=), and therefore be used as map keys, provided they're composed from elements for which equality is also defined, using element-wise comparison.

		type Day struct {
			long string
			short string
		}

		Christmas := Day{"Christmas", "Xmas"}
		Thanksgiving := Day{"Thanksgiving", "Turkey"}

		holiday := map[Day]bool {
			Christmas: true,
			Thanksgiving: true,
		}

		fmt.Printf("Christmas is a holiday: %t\n", holiday[Christmas])

	Second, Go 1 removes the definition of equality for `func` values, except for comparison with `nil`. Finally map equality is gone too, also except for comparison with `nil`.

	Note that equality is still be undefined for slices, for which the calculation is in general infeasible. Also, note that the ordered comparison operators (< <= > >=) are still undefined for structs and arrays.

	UPD: Struct and array equality is a new feature, so existing code needs no changes. Existing code that depends on function or map equality will be rejected by the compiler and will need to be fixed by hand. Few programs will be affected, but the fix may require some redesign.
*/

/*
	MAJOR CHANGES TO THE LIBRARY
*/

/*
	1. The `error` type and `errors` package.
	The placement of `os.Error` in package `os` is mostly historical: errors came up when implementing package `os`, and they seemed system-related at the time. Since then it has become clear that erros are more fundamental than the operating system. For example, it'd be nice to use `Errors` in packages that `os` depends on like `syscall`. Also, having `Error` in `os` introduces many dependencies on `os` that would otherwise not exist.

	Go 1 solves these problems by introducing a built-in error interface type and a separate `errors` package (analogous to `bytes` and `strings`) that contains utility funcs. It replaces `os.NewError` with `errors.New`, giving errors a more central place in the environment.

	So the widely-used String method doesn't cause accidental satisfaction of the error interface, the `error` interface uses instead the name `Error` for that method:

		type error interface {
			Error() string
		}

	The `fmt` library automatically invokes Error, as it already does for String, for easy printing of error values.

		type SyntaxError struct {
			File string
			Line int
			Message string
		}

		func (se *SyntaxError) Error() string {
			return fmt.Sprintf("%s:%d: %s", se.File, se.Line, se.Message)
		}

	All standart packages have been updated to use the new interface; the old os.Error is gone.

	A new package, `errors`, contains the func:

		func New(text string) error

	to turn a string into an error. It replaces the old os.NewError.

		var ErrSyntax = errors.New("syntax error")

	UPD: Running `go fix` will update almost all code affected by the change. Code that defines error types with a String method will need to be updated by hand to rename the methods to Error.
*/

/*
	2. System call errors
	The old `syscall` package, which predated `os.Error` (and just about everything else), returned as `int` vals. In turn, the `os` package forwarded many of these errors, such as `EINVAL`, but using a different set of errors on each platform. This behavior was unpleasant and unportable.

	In Go 1, `syscall` package instead returns an `error` for system call errors. On Unix, the implementation is done by a `syscall.Errno` type that satisfies `error` and replaces the old `os.Errno`

	The changes affecting `os.EINVAL` and relatives are described elsewhere (https://tip.golang.org/doc/go1#os)

	UPD: Running `go fix` will update almost all code affected by the change. Regardless, most code should use the `os` package rather than `syscall` and so will be unaffected.
*/

/*
	3. Time
	The old Go `time` package had int64 units, no real type safety, and no distinction between absolute times and durations.

	One of the most sweeping changes in the Go 1 library is therefore a complete redesign of the `time` package. Instead of an integer number of nanoseconds as an `int64`, and a separate `*time.Time` type to deal with human units such as hours and years, there are now two fundamental types:
		- time.Time (a value, so the * is gone), which represents a moment in time
		- time.Duration, which represents an interval.
	Both have nanosecond resolution. A Time can represent any time into the ancient past and remote future, while a Duration can span plus or minus only about 290 years. There are methods on these types, plus a number of helpful predefined constant durations such as `time.Second`.

	Among the new methods are things like `Time.Add`, which adds a Duration to a Time, and Time.Sub, which subtracts two Times to yield a Duration.

	The most important semantic change is that the Unix epoch (Jan 1, 1970) is now relevant only for those funcs and methods that mention Unix: `time.Unix` and the `Unix` and `UnixNano` methods of the Time type. In particular, `time.Now` returns a `time.Time` value rather than, in the old API, an integer nanosecond count since the Unix epoch.

	// sleepUntil sleeps until the specified time. It returns immediately if it's too late.
	func sleepUntil(wakeup time.Time) {
		now := time.Now() // A Time
		if !wakeup.After(now) {
			return
		}

		delta := wakeup.Sub(now) // A Duration
		fmt.Printf("Sleeping fof %.3fs\n", delta.Seconds())
		time.Sleep(delta)
	}

	The new types, methods, and constants have been propagated through all the std packages that use time, such as `os` and its representation of file time stamps.

	UPD: The `go fix` tool will update many uses of the old `time` package to use the new types and methods, although it doesn't replace values such as `1e9` representing nanoseconds per second. Also, because of type changes in some of the values that arise, some of the expressions rewritten by the fix tool may require further hand editing; in such cases the rewrite will include the correct func or method for the old functionality, but may have the wrong type or require further analysis.
*/

/*
	MINOR CHANGES TO THE LIBRARY
	// TODO - https://clck.ru/3FFJmd
*/

/*
	THE `go` COMMAND
	Go 1 introduces the `go` command, a tool for fetching, building, and installing Go packages and commands. The `go` command does away with makefiles, instead using Go source code to find dependencies and determine build conditions. Most existing Go programs will no longer require makefiles to be built.

	UPD: Projects that depend on the Go project's old makefile-based build infrastructure (Make.pkg, Make.cmd and so on) should switch to using the `go` command for building `Go` code and, if necessary, rewrite their makefiles to perform any auxiliary build tasks.
*/

/*
	The cgo command
	In Go 1, the `cgo` command uses a different `_cgo_export.h` file, which is generated for packages containing `// export` lines. The `_cgo_export.h` file now begins with the `C` preamble comment, so that exported function definitions can use types defined there. This has the effect of compiling the preamble multiple times, so a package using `//export` mustn't put function definitions or variable initializations in the `C` preamble.
*/

/*
	Packaged releases
	One of the most significant changes associated with Go 1 is the availability of prepackaged, downloadable distributions. They are available for many combinations or architecture and OS (including Windows) and the list will grow. Installation details are described on the `Getting Started` page, while the distributions themselves are listed on the `download page`.
*/
