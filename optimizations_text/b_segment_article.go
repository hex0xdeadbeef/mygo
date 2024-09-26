package main

import (
	"fmt"
	"hash/fnv"
	"runtime"
	"unsafe"
)

// https://segment.com/blog/allocation-efficiency-in-high-performance-go-services/
/*
	ALLOCATION EFFICIENCY IN HIGH-PERFORMANCE GO SERVICES

	ESCAPE ANALYSIS
1. Patterns which typically cause variables to escape to heap
	1) Sending pointers or values containing pointers to channels.
		At compile time there's no way to know which goroutine will receive the data on a channel.
		Therefore the compiler cannot determine when this data will no longer be referenced.
	2) Storing pointers or values containing pointers in a slice
		An example of this is a type like []*string or just []string. This always causes the
		contents of the slice to escape. Even though the backing array of the slice may still be on
		the stack, the refferenced data escaped to the heap.
	3) Backing arrays of slices that get reallocated because an append would exceed their capacity.
		In cases where the initial size of a slice is known at compile time. it'll begin its
		allocation on the stack. If this slice's underlying storage must be expanded based on data
		only known at runtime, it'll be allocated on the heap.
	4) Calling methods on an interface types.
		Method calls on interface types are a dynamic dispatch - the actual concrete implementation
		to use is only determinable at runtime. Consider a variable "r" with an interface type of
		io.Reader. A call to r.Read(b) will cause both the value "r" and the backing array of the
		byte slice "b" to escape and therefore be allocated on the heap.

	SOME POINTERS
1. The rule of thumb: pointers point to data allocated on the heap. Ergo, reducing the number of
	pointers in a program reduces the number of heap allocations.
	This is not an axiom, but we've found it to be the common case in real-world Go programs.
	In many cases copying a value is much less expensive than the overhead of using a pointer. Why?
		1) The compiler generates checks when dereferencing a pointer.
			The purpose is to avoid memory corruption by running panic() if the pointer is nil.
			This is extra code that must be executed at runtime. When data is passed by value, it
			cannot be nil.
		2) Pointers often have poor locality of reference.
			All of the values used within a function are collocated in memory of the stack.
			Locality of reference is an important aspect of efficient code. It dramatically
			increases the chance that a value is warm in CPU caches and reduces the risk of a miss
			penalty during prefetching.
		3) Copying objects within a cache line is the roughly equivalent to copying a single pointer
			CPUs move memory between caching layers and main memory on cache lines of constant size.
			On x86 this is 64 bytes. Further, Go uses a technique called Duff's device to make
			common memory operations like copies very efficient.
2. Pointers should primarily be used to reflect ownership semantics and mutability. In practice,
	the use of pointers to avoid copies should be infrequent.

	Don't fall into the trap of premature optimization. It's good to develop a habit of passing
	data by value, only failing back to passing pointers when necessary. An extra bonus is the
	increased safety of eliminating nil.
3. Reducing the number of pointers in a program can yield another helpful result as the garbage
	collector will skip regions of memory that it can prove will contain no pointers. For example:
	regions of the heap which back slices of type []byte aren't scanned at all. This also holds
	true for arrays of struct types that don't contain any fields with pointer types.
4. Not only does reducing pointers result in less work for the GC, it produces more cache-friendly
	code. Reading memory moves data from main memory into the CPU caches. Caches are finite, so
	some other piece of data must be evicted to make room. Evicted data may still be relevant to
	other portions of the program. The resulting cache thrashing can cause unexpected and sudden
	shifts the behavior of production services.

	DIGGING FOR POINTER GOLD
1. Reducing pointer usage often means digging into the source code of the types used to construct
	our programs.
2. For example we have user-structs for a service:
	type retryQueue struct {
		buckets 		[][]retryItem // each bucket represents a 1 second interval
		currentTime 	time.Time
		currentOffset 	int
	}

	type retryItem struct {
		id ksuid.KSUID 	// ID of the item to retry
		time time.Time 	// exact time at which the item has to be retried
	}

	The size of the outer array in buckets is constant, but the number of items in the contained []
	[]retryItem slice will vary at runtime. The more retries, the larger these slices will grow.

	Digging into the implementation details of each field of a retryItem, we learn that
		- KSUID is a type alias for [20]byte, which has no pointers, and therefore can be ruled
			out.
		- currentOffset is an int, which is a fixed-size primitive, and can also be ruled out.
		- Next, looking at the implementation of time.Time type.
3. As for Time struct:
	type Time struct {
		sec 	int64 *
		nsec 	int32 **
		loc *Location // pointer to the time zone structure
	}

	The time.Time struct contains an internal pointer for the "loc" field. Using it within the
	retryItem type causes the GC to "chase" the pointers on these structs each time it passes
	through this area of the heap.

	This is a typical case of cascading effects under unexpected circumstances.
	During normal operation failures are uncommon. Only a small amount of memory is used to stores
	retries. When failures suddenly spike, the number of items in the retry queue can increase by
	thousands per second, bringing with it a significally increased workload for the GC.
4. For this particular use case, the timezone information in time.TIme isn't necessary. These
	timestamps are kept in memory and are never serialized.

	Therefore these data structures can be refactored to avoid thus type entirely:
		type retryItem struct {
			id 		ksuid.KSUID
			sec 	int64 *
			nsec 	int32 **
		}

		func (item *retryItem) time() time.Time {
			return time.Unix(item.sec, int64(item.nsec))
		}

		func makeRetryItem(id ksuid.KSUID, time time.Time) retryItem {
			return retryItem {
				id: id,
				nsec: uint32(time.Nanosecond()),
				sec: time.Unix(),
			}
		}

	Now the retryItem doesn't contain any pointers. This dramatically reduces the load on the GC
	as the entire footprint of retryItem  is knowable at compile time.

	PASS ME A SLICE
1. Slices are fertile ground for inefficient allocation behavior in hot code paths. Unless the
	compiler knows the size of the slice at compile time, the backing arrays for slices (and
	maps!) are allocated on the heap.
2. The usage of time.Time in Go's is particularly expensive.
	The profiler shows a large percentage of the heap allocations in code that serializes a time.
	Time values so that it can be sent over the wire to somewhere.
3. The code calling the method Format() on the time.Time, which returns a string. Wait, aren't we
	talking about slices? Well, according to the official Go blog, a string is just a "read-only
	slices of bytes with a bit of extra syntactic support from the language". Most of the same
	rules around allocation apply.
4. The profiler tells us a massive 12.38% of the allocations are occuring when running this
	Format() method. What does Format() do?

	It turns out there's a much more efficient way to do the same thing that uses a common pattern
	across the standart library. While the Format() method is easy and convenient, code using
	AppendFormat() can be much easier on the allocator.

	Peering into the source code for the "time" package, we notice that all internal uses are
	AppendFormat() and not Format(). This is a pretty strong hint that AppendFormat() is going to
	yield more performant behavior.
5. func (t time) AppendFormat(b []byte, layout string) []byte is profitable.
	The method AppendFormat is like Format() but appends the textual representation to b and
	returns the extended buffer.

	In fact, the Format method just wraps the AppendFormat() method:
		func (t time) Format(layout string) string {
			const (
				bufSize = 1 << 6
			)

			var (
				b []byte
				max = len(layout) + 10
			)

			switch max < bufSize {
				case true:
					var buf [bufSize]byte
					b = buf[:0]
				default:
					b = make([]byte, 0, max)
			}
			b = AppendFormat(b, layout)

			return string(b)
		}
6. Most importantly, AppendFormat() gives the programmer far more control over allocation. It
	requires passing the slice to mutate rather than returning a string that it allocates
	internally like Format().

	Using AppendFormat() instead of Format() allows the same operation to use a fixed-size
	allocation and thus is eligible for stack placement.
7. Code changes explanation.
	- The first thing to notice is that "var a [1<<6]byte" is a fixed-size array. It's size is
	known at compile-time and its use is scoped entirely to this function, so we can deduce that
	this will be allocated on the stack. However, this type can't be passed to AppendFormat(),
	which accepts type []byte. Using the a[:0] notation converts the fixed-size array to a slice
	type represented by "b" that is backed by this array. This will pass the compiler's check and
	be allocated on the stack.
	- More critically, the memory that would otherwise be dynamically allocated is passed to
	AppendFormat(), a method which itself passes the compiler's stack allocation checks. In the
	previous version, Format() is used, which contains allocations of sizes that can't be
	determined at compile time and therefore don't qualify for stack allocation.

	The result of this relatively small change massively reduced allocationsin this code path!

	INTERFACE TYPES AND YOU
1. It's fairly common knowledge that method calls on interface types are more expensive than those
	on struct types. Method calls on interface types are executed via dynamic dispatch. This
	severely limits the ability for the compiler to determine the way that code will be executed
	at runtime. So far we've largely discussed shaping code so that the compiler can understand
	its behavior best at compile-time.

	Interface types throw all of this away!
2. Unfortunately interface types are a very useful abstraction - they let us write more flexible
	code. A common case of interfaces being used in the hot path of a program is the hashing
	functionality provided by a standart library's "hash" package. The hash package defines a set
	of generic interfaces and provides several concrete implementations.
3. The "go build -gcflags '-m' main.go" shows that passing an argument to fmt.Printf(...) causes
	"s" to escape to the heap, because passing an argument to a function as the interface value
	forces the value to escape.

	Moreover, interfaces use dynamic dispatch an we can see it in the next rows:
		./main.go:20:9: devirtualizing h.Write to *fnv.sum64
		./main.go:21:16: devirtualizing h.Sum64 to *fnv.sum64
	This causes some overheads at run time.
4. The way we need to use is described in "optimizations_vk"

	ONE WEIRD TRICK
1. Our final anecdote is a more amusing than practical. However, it's a useful example for
	understanding the mechanics of the compiler's escape analysis. When reviewing the standart
	library for the optimizations covered, we came across a rather curious piece of code.
2. The code is:
	// noescape hides a pointer from escape analysis.  noescape is
	// the identity function but escape analysis doesn't think the
	// output depends on the input.  noescape is inlined and currently
	// compiles down to zero instructions.
	// USE CAREFULLY!
	//go:nosplit
	func noescape(p unsafe.Pointer) unsafe.Pointer {
    	x := uintptr(p)
    	return unsafe.Pointer(x ^ 0)
	}

	The function will hide the passed pointer from the compiler's escape analysis functionality.
3. The escape analysis output from the compiler shows us that the FooTrick version doesn't escape.
4. This is the compiler recognizing that the NewFoo() function takes a reference to the string and
	stores it in the struct, causing it to escape. However, no such output appears for the
	NewFooTrick() function. If the call to noescape() is removed, the escape analysis moves the
	data referenced by the FooTrick() struct to the heap.
5. What's happening here?
	The noescape() function masks the dependency between the input argument and the return value.
	The compiler doesn't think that "p" escapes via "x" because the uintptr() produces a reference
	that is opaque to the compiler. The builtin uintptr() type's name may lead one to believe this
	is a bona fide pointer type, but from the compiler's perspective it's just an integer that
	just happens to be large enough to store a pointer. The final line of code constrcuts and
	returns an unsafe.Pointer value from a seemingly arbitrary integer value.

	Nothing to see here folks!

	noescape() is used in dozens of functions in the runtime package that use unsafe.Pointer().
	It's useful in cases where the author knows for certain that data referenced by an unsafe.
	Pointer doesn't escape, but the compiler naively thinks otherwise.


*/

// Large slice stack allocation
func LargeSliceStackAllocation() {
	printAlloc()
	var (
		s = []byte{1 << 28: 1}
	)
	printAlloc()

	s = append(s, s...)

	runtime.GC()
	printAlloc()

	runtime.KeepAlive(&s)
}

func printAlloc() {
	var (
		m runtime.MemStats
	)
	runtime.ReadMemStats(&m)

	fmt.Printf("%d MB\n", m.Alloc/(1<<20))
}

// hashing usage
func UsageOfHashIt() {
	s := "Hello!"
	fmt.Printf("The FNV64a hash of %q is \"%v\"\n", s, hashIt(s))
}

/*
# command-line-arguments
./main.go:8:6: can inline main
./main.go:15:12: inlining call to fmt.Printf
./main.go:19:16: inlining call to fnv.New64
./main.go:20:9: devirtualizing h.Write to *fnv.sum64
./main.go:20:9: inlining call to fnv.(*sum64).Write
./main.go:21:16: devirtualizing h.Sum64 to *fnv.sum64
./main.go:21:16: inlining call to fnv.(*sum64).Sum64
./main.go:18:13: in does not escape
./main.go:20:17: ([]byte)(in) does not escape
./main.go:20:17: zero-copy string->[]byte conversion
./main.go:15:12: ... argument does not escape
./main.go:15:50: s escapes to heap
./main.go:15:59: hashIt(s) escapes to heap
*/
func hashIt(in string) uint64 {
	h := fnv.New64()
	h.Write([]byte(in))
	out := h.Sum64()
	return out
}

// Escape analysis supressing
type Foo struct {
	S *string
}

func (f *Foo) String() string {
	return *f.S
}

func NewFoo(s string) Foo {
	return Foo{S: &s}
}

type FooTrick struct {
	S unsafe.Pointer
}

func (f *FooTrick) String() string {
	return *(*string)(f.S)
}

func NewFooTrick(s string) FooTrick {
	return FooTrick{S: noescape(unsafe.Pointer(&s))}
}

// noescape hides a pointer from escape analysis.  noescape is
// the identity function but escape analysis doesn't think the
// output depends on the input.  noescape is inlined and currently
// compiles down to zero instructions.
// USE CAREFULLY!
func noescape(p unsafe.Pointer) unsafe.Pointer {
	x := uintptr(p)
	return unsafe.Pointer(x ^ 0)
}

/*
# command-line-arguments
./main.go:18:6: can inline NewFoo with cost 5 as: func(string) Foo { return Foo{...} }
./main.go:39:6: can inline noescape with cost 9 as: func(unsafe.Pointer) unsafe.Pointer { up := uintptr(p); return unsafe.Pointer(up ^ uintptr(0)) }
./main.go:30:6: can inline NewFooTrick with cost 16 as: func(string) FooTrick { return FooTrick{...} }
./main.go:14:6: can inline (*Foo).String with cost 4 as: method(*Foo) func() string { return *f.S }
./main.go:26:6: can inline (*FooTrick).String with cost 4 as: method(*FooTrick) func() string { return *(*string)(f.S) }
./main.go:44:6: can inline UsageOfNoEscape with cost 68 as: func() { s := "hello"; f1 := NewFoo(s); f2 := NewFooTrick(s); s1 := (*Foo).String(f1); s2 := (*FooTrick).String(f2); _ = s1; _ = s2 }
./main.go:5:6: can inline main with cost 70 as: func() { UsageOfNoEscape() }
./main.go:6:17: inlining call to UsageOfNoEscape
./main.go:6:17: inlining call to NewFoo
./main.go:6:17: inlining call to NewFooTrick
./main.go:6:17: inlining call to (*Foo).String
./main.go:6:17: inlining call to (*FooTrick).String
./main.go:6:17: inlining call to noescape
./main.go:31:29: inlining call to noescape
./main.go:46:14: inlining call to NewFoo
./main.go:47:19: inlining call to NewFooTrick
./main.go:49:17: inlining call to (*Foo).String
./main.go:50:17: inlining call to (*FooTrick).String
./main.go:47:19: inlining call to noescape
./main.go:14:7: parameter f leaks to ~r0 with derefs=2:
./main.go:14:7:   flow: ~r0 = **f:
./main.go:14:7:     from f.S (dot of pointer) at ./main.go:15:11
./main.go:14:7:     from *f.S (indirection) at ./main.go:15:9
./main.go:14:7:     from return *f.S (return) at ./main.go:15:2
./main.go:14:7: leaking param: f to result ~r0 level=2
./main.go:18:13: s escapes to heap:
./main.go:18:13:   flow: ~r0 = &s:
./main.go:18:13:     from &s (address-of) at ./main.go:19:16
./main.go:18:13:     from Foo{...} (struct literal element) at ./main.go:19:12
./main.go:18:13:     from return Foo{...} (return) at ./main.go:19:2
./main.go:18:13: parameter s leaks to ~r0 with derefs=0:
./main.go:18:13:   flow: ~r0 = &s:
./main.go:18:13:     from &s (address-of) at ./main.go:19:16
./main.go:18:13:     from Foo{...} (struct literal element) at ./main.go:19:12
./main.go:18:13:     from return Foo{...} (return) at ./main.go:19:2
./main.go:18:13: moved to heap: s
./main.go:26:7: parameter f leaks to ~r0 with derefs=2:
./main.go:26:7:   flow: ~r0 = **f:
./main.go:26:7:     from f.S (dot of pointer) at ./main.go:27:21
./main.go:26:7:     from *(*string)(f.S) (indirection) at ./main.go:27:9
./main.go:26:7:     from return *(*string)(f.S) (return) at ./main.go:27:2
./main.go:26:7: leaking param: f to result ~r0 level=2
./main.go:30:18: s does not escape
./main.go:39:15: p does not escape
*/
func UsageOfNoEscape() {
	s := "hello"
	f1 := NewFoo(s)
	f2 := NewFooTrick(s)

	s1 := f1.String()
	s2 := f2.String()

	_ = s1
	_ = s2
}
