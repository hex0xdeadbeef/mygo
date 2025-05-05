package main

import (
	"fmt"
	"reflect"
	"strings"
	"unsafe"
)

/*
	INTRO

	Go pointers cannot participate arithmetic ops, and for two arbitrary pointer types, it's very possible that their vals can't be converted to each other.

	Go supports type-unsafe pointers, which are pointers without restrictions made for safe pointers. Type-unsafe pointers are also called unsafe pointers in Go. Go unsafe pointers are much like C pointers, they are powerful, and also dangerous. For some cases, we can write more efficient code with the help of unsafe pointers. On the other hand, by using unsafe pointers, it's easy to write bad code which is too subtle to detect in time.

	Code depending on unsafe pointers works today could break since a later Go version.
*/

/*
	ABOUT THE `unsafe` STANDARD PACKAGE

	The unsafe.Pointer type is defined as:
		type Pointer *ArbitraryType

	ArbitraryType just hints that a unsafe.Pointer value can be converted
	to any safe pointer values in Go (and vice versa).

	Go unsafe pointers mean the types whose underlying types are unsafe.Pointer.

	The zero values of unsafe pointers are also represented with the predeclared
	identifier `nil`

	Before Go 1.17, the `unsafe` standard package has already provided three funcs:
		- func Alignof(variable ArbitraryType) uintptr, which is used to get the address alignment of a value. Please notes, the aligns for struct-field values and non-field values of the same type may be different, though for the standard Go compiler, they're always the same, For the gccgo compiler, they may be different.
		- func Offsetof(selector ArbitraryType) uintptr, which is used to get the address offset of a field in a struct value. The offset is relative to the address of the struct value. The return results should be always the same for the same corresponding field of values of the same struct type in the same program.
		- func Sizeof(variable ArtbitraryType) uintptr, which is used to get the size of a value (a.k.a the size of the type of the value). The return results should be always the same for all values of the same type in the same program.
	Note:
		- The types of the return results of the three funcs are all `uintptr`
		- Although the return results of calls of any of the three funcs are consistent in the same program, they might be different crossing operating systems, crossing architectures, crossing compilers, and crossing compiler versions.
		- !!! Calls to the three funcs are always evaluated at compile time. The evaluation results are typed constants with type `uintptr
		- The argument passed to a call to the unsafe.Offsetof` function must the struct field selector form value.field. The selector may denote an embedded field, but the field must be reachable without implicit pointer indirections.

*/

func OffsetofAlignofSizeofUsageA() {
	var x struct {
		a int64
		b bool
		c string
	}

	const M, N = unsafe.Sizeof(x.c), unsafe.Sizeof(x)
	fmt.Println(M, N) // 16 32

	fmt.Println(unsafe.Alignof(x.a)) // 8
	fmt.Println(unsafe.Alignof(x.b)) // 1
	fmt.Println(unsafe.Alignof(x.c)) // 8

	fmt.Println(unsafe.Offsetof(x.a)) // 0
	fmt.Println(unsafe.Offsetof(x.b)) // 8
	fmt.Println(unsafe.Offsetof(x.c)) // 16
}

func OffsetofAlignofSizeofUsageB() {
	type T struct {
		c string
	}

	type S struct {
		b bool
	}

	var x struct {
		a int64
		*S
		T
	}

	fmt.Println(unsafe.Offsetof(x.a)) // 0
	fmt.Println(unsafe.Offsetof(x.S)) // 8
	fmt.Println(unsafe.Offsetof(x.T)) // 16

	// This line compiles, for c can be reached
	// without implicit pointer indirections
	fmt.Println(unsafe.Offsetof(x.c)) // 16

	// This line doesn't compile, for b must be
	// reached with the implicit pointer field S.
	// Error: invalid argument: field b is embedded via
	// a pointer in struct{a int64; *S; T}
	// compilerInvalidOffsetof

	// fmt.Println(unsafe.Offsetof(x.b))

	// This line compiles. However, it prints the offset of field b in the value x.S
	fmt.Println(unsafe.Offsetof(x.S.b)) // 0
}

/*
	NEWCOMINGS

	1. The new type `type IntegerType int`. The following is its definition. This type doesn't denote a specified type. It just represents any arbitrary integer type. We can view it as
	a generic type.

	2. func Add(ptr Pointer, len IntegerType) Pointer. This function adds an offset to the address
	represented by an unsafe pointer and return a new unsafe pointer wich represents the new
	address.

	3. func Slice(ptr *ArbitraryType, len IntegerType) []ArbitraryType. This function is used to
	derive a slice wtih the specified length from a safe pointer, where ArbitraryType is the element
	type of the result slice.

	4. func SliceData(slice []ArbitraryType) *ArbitraryType. This function is used to get the
	pointer to the underlying element sequence of a slice.

	5. func String(ptr *byte, len IntegerType) string. This function is used to derive a string with
	the specified length from a safe `byte` pointer.

	6. func StringData(str string) *byte. This function is used to get the pointer to the underlying
	byte sequence of a string. Please note, don't pass empty strings as arguments to this function.
*/

func UsingUnsafePackageFuncs() {
	a := [16]int{3: 3, 9: 9, 11: 11}
	fmt.Println(a) // [0 0 0 3 0 0 0 0 0 9 0 11 0 0 0 0]

	elemSize := int(unsafe.Sizeof(a[0]))

	p9 := &a[9]
	up9 := unsafe.Pointer(p9)

	// Shifts the memory back by 6 elems of the array
	p3 := (*int)(unsafe.Add(up9, -6*elemSize))
	fmt.Println(*p3) // 3

	// Get reslicing of the initial slice, starting from the 9th elem
	// and ending at the last elem of the initial slice.
	s := unsafe.Slice(p9, 7)
	fmt.Println(s) // [9 0 11 0 0 0 0]

	// Reslicing
	s = s[:3]
	fmt.Println(s, len(s), cap(s)) // [9 0 11] 3 7

	// `nil` slice creation
	t := unsafe.Slice((*int)(nil), 0)
	fmt.Println(t == nil, len(t), cap(t)) // true 0 0

	a[len(a)-1] = -1

	// The following two calls are dangerous.
	// They make the results reference
	// unknown memory blocks
	lastElem := (*int)(unsafe.Add(up9, 6*elemSize))
	fmt.Println(*lastElem) // -1

	unknownMemoryElem := (*int)(unsafe.Add(up9, 7*elemSize))
	fmt.Println(*unknownMemoryElem) // 0
}

// The following two functions may be used to do conversions between
// strings and byte slices, in type unsafe manners. Comparing with type safe manners,
// the type unsafe manners don't duplicate underlying byte sequences of strings and byte slices,
// so they're more performant
func String2ByteSlice(s string) []byte {
	if s == "" {
		return nil
	}

	// There's no copying
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

func ByteSliceToString(b []byte) string {
	if len(b) == 0 {
		return ""
	}

	// There's no copying
	return unsafe.String(unsafe.SliceData(b), len(b))
}

/*

	UNSAFE POINTERS RELATED CONVERSION RULES

	Currently (Go 1.24), Go compilers allow the following explicit conversions:
		- A safe pointer can be explicitly converted to an unsafe pointer, and vice versa.
		- A uintptr value can be explicitly converted to an unsafe pointer, and vice versa. !!! But
		please note, a nil unsafe.Pointer shouldn't be converted to uinptr and back with arithmetic.

	By using these conversions, we can convert a safe pointer value to an arbitrary safe pointer
	type.

	However, although these conversions are all legal at compile time, not all of them are valid
	(safe) at runtime. These conversions defeat the memory safety the whole Go type system (except
	the unsafe part) tries to maintain. We must follow the instructions listed in a later section
	below to write valid Go code with unsafe pointers.
*/

/*
	SOME FACTS IN GO WE SHOULD KNOW

	Fact 1:
		Unsafe pointers are pointers and `uintptr` are integers.

		Each of non-nil safe and unsafe pointers references another value. However, `uintptr` values
		don't reference any values, they're just plain integers, though often each of them stores an
		integer which can be used to represent a memory address.

		Go is a language supporting automatic garbage collection. When a Go program is running, Go
		runtime wull check which memory blocks aren't used by any value any more and collect the
		memory allocated for these unused blocks, from time to time. Pointers play an important
		role in the check process. If a memory block is unreachable from (references by) any values
		still in use, then Go runtime thinks it's an unused value and it can be safely garbage
		collected.


		!!! A `uintptr` values are integers, they can participate arithmetic ops


	Fact 2:
		At run time, the GC may run at an uncertain time, and each garbage collection process
		may last an uncertain duration. So when a memory block becomes unused, it may be collected
		at an uncertain time.


	Fact 3:
		The adresses of some values might change at run time

		Here, we should just know that when the size of the stack of a goroutine changes, the memory
		blocks allocated on the stack will be moved. In other words, the addresses of the values on
		these memory blocks will change.


	Fact 4:
		The life range of a value at run time may not be not as large as it looks in code


	Fact 5:
		*unsafe.Pointer is a general safe pointer type

		Yes, *unsafe.Pointer is a safe pointer type. Its base type is unsafe.Pointer. As it is a safe pointer, according the conversion rules listed above, it can be converted to unafe.Pointer type, and vice versa.
*/

func CaveatWithUinptrValA() {
	// Integers creation
	p0, y, z := createInt(), createInt(), createInt()
	p1 := unsafe.Pointer(y)

	// DANGEROUS because either:
	// - goroutine stack can grow and all vals will be put to other addresses or `y`
	// - `y` can be garbage collected
	p2 := uintptr(unsafe.Pointer(z))

	// At the time, even if the address of the int
	// value referenced by z is still stored in p2,
	// the int value has already become unused, so
	// garbage collector can collect memory allocated
	// for it now. On the other hand, the int vals referenced
	// by p0 and p1 are still in use

	// `uintptr cam participate arithmetic ops`
	p2 += 2
	p2--
	p2--

	*p0 = 1                         // ok
	*(*int)(p1) = 2                 // ok
	*(*int)(unsafe.Pointer(p2)) = 3 // dangerous!

	// In the above example, the fact that value p2 is still in use can't
	// guarantee that the memory block ever hosting the int val referenced
	// by z hasn't been garbage collected yet. In other words when
	// *(*int)(unsafe.Pointer(p2)) = 3 is executed, the memory block may be collected,
	// or not. It's dangerous to dereference the address stored in val `p2` to an int
	// val, for it's that the memory block has been already reallocated for another val
	// (even for another program)

}

// Assume createInt won't be inlined
//
//go:noinline
func createInt() *int {
	return new(int)
}

func LifespanOfObject() {
	type T struct {
		x int
		y *[1 << 23]byte
	}

	t := &T{y: new([1 << 23]byte)}
	// Dangerous because:
	// - the value of field `y` can be garbage collected unexpectedly
	// - the value of field `y` can be moved to other memory block
	p := uintptr(unsafe.Pointer(&t.y[0]))
	_ = p

	// ... // use T.x and T,.y

	// A smart compiler can detect that the value
	// t.y will never be used again and think the memory block
	// hosting t.y can be collected now

	// using *(*byte)(unsafe.Pointer(p)) is dangerous here.

	// Continue using value t, but only use its `x` field.
	fmt.Println(t.x)
}

func UnsafePointerOfUnsafePointer() {
	x := 123
	p := unsafe.Pointer(&x) // p -> x

	pp := &p       // pp -> p -> x
	pOriginal := p // Save the original val

	p = unsafe.Pointer(pp)    // p -> pp
	pp = (*unsafe.Pointer)(p) // pp -> p

	fmt.Println(*(*int)(pOriginal)) // Print out the original value
}

/*
	HOW TO USE UNSAFE POINTERS CORRECTLY?

	Pattern 1:
		Convert a *T1 value to unsafe Pointer, then convert the unsafe pointer value to *T2
		*T1 -> UP -> T2

		As mentioned above, by using the unsafe pointer conversion rules above, we can convert a
		value of *T1 to type *T2, where T1 and T2 are two arbitrary types. However, we should only
		do such conversions if the size of T1 is no smaller than T2, and only if the conversions
		are meaningful.

		Conditions:
			- unsafe.Sizeof(T2) <= unsafe.Sizeof(T1)
			- The types should be matching

		A practice by using the pattern is to convert a byte slice, which won't be used after the
		conversion, to a string, as the following code shows. In this conversion, a duplication of
		the underlying byte sequence is avoided.

	Pattern 2:
		Convert unsafe pointer to `uintptr`, then use the `uintptr` value

		This pattern isn't very useful. Usually, we print the result uintptr vals to check the
		memory address stored in them. However there are other both safe and less verbose ways to
		this job. So this pattern isn't much useful.


	Pattern 3:
		Convert unsafe pointer to uinptr, do arithmetic ops with the uintptr val, the convert it
		back.

		In this pattern, the result unsafe pointer must continue to point into the original
		allocated memory block.

		!!! Another detail which should be also noted is that, it's not recommended to store the end
		boundary of a memory block in a pointer (either safe or unsafe). Doing this will prevent
		another memory block which closely follos the former memory block from being garbage
		collected, or crash program if that boundary address is not valid for any allocated memory
		blocks (depending on compiler implementations).

	Pattern 4:
		Convert unsafe pointers to uintptr values as args of syscall.Syscall calls.

		// Assume this func won't be inlined
		func DoSomething(addr uintptr) {
			// read or write values at the passed address...
		}

		The reason why the above function is dangerous is that the function itself cannot guarantee
		the memory block at the passed argument address is not garbage collected yet. If the memory
		block is collected or is reallocated for other vals, then the ops made in the func body is
		dangerous.

	Pattern 5:
		Convert the uintptr result of reflect.Value.Pointer or reflect.Value.UnsafeAddr method call
		to unsafe.Pointer

		The methods Pointer and UnsafeAddr of the Value type in the reflect std package both return
		a result of type uintptr instead of unsafe.Pointer. This is a deliberate design, which is to
		avoid converting the results of calls (to the two methods) to any safe pointer types without
		importing the unsafe std package.

		The design requires the return result of a call to either of the two methods must be
		converted to an unsafe pointer immediately after making the call. Otherwise, there will be
		small time window in which the memory block allocated at the address stored in the result
		might lose all references and be garbage collected.

		Please note that, Go 1.19 introduces a new method reflect.Value.UnsafePointer(), which returns a unsafe.Pointer value and is preferred over the two just mentioned funcs. That means, the old deliberate design is thought as not good now.

	Patter 6:
		Convert a reflect.SliceHeader.Data or reflect.StringHeader.Data field to unsafe pointer, and the inverse.

		For the same reason mentioned for the last subsection, the Data fields of the struct type
		SliceHeader and StringHeader in the reflect std package are declared with type uintptr instead of unsafe.Pointer.

		We can convert a string pointer to a (*reflect.StringHeader) pointer value, so that we can manipulate the internal of the string. The same, we can convert a slice pointer to a *reflect.SliceHeader pointer value, so that we can manipulate the internal of the slice.

		In general, we should only get a *reflect.StringHeader pointer value from an actual (already existed) string, or get a *reflect.SliceHeader pointer value from an actual (already existed) slice. We shouldn't do the contrary, such as creating a string from a
		new allocated StringHeader, or creating a new allocated SliceHeader.
)
*/

func Float64ToBits(f float64) uint64 {
	return *(*uint64)(unsafe.Pointer(&f))
}

func Float64FromBits(b uint64) float64 {
	return *(*float64)(unsafe.Pointer(&b))
}

func ConversionsBetweenSliceTypes() {
	type MyString string

	ms := []MyString{"C", "C++", "Go"}
	fmt.Printf("%s\n", ms) // [C C++ Go]

	// ss := ([]string)(ms) // cannot convert ms (variable of type []MyString) to type []
	// stringcompilerInvalidConversion

	ss := *(*[]string)(unsafe.Pointer(&ms))
	ss[1] = "Zig"
	fmt.Printf("%s\n", ss) // [C Zig Go]

	// ms = []MyString(ss) // cannot convert ss (variable of type []string) to type []
	// MyStringcompilerInvalidConversion

	ms = *(*[]MyString)(unsafe.Pointer(&ss))
	fmt.Printf("%s\n", ms) // [C Zig Go]

	// Since Go 1.17, we may also use the unsafe.Slice function to do the conversions
	ss = unsafe.Slice((*string)(&ms[0]), len(ms))
	fmt.Printf("%s\n", ss) // [C Zig Go]

	ms = unsafe.Slice((*MyString)(&ss[0]), len(ss))
	fmt.Printf("%s\n", ms) // [C Zig Go]

}

// A practice by using the pattern is to convert a byte slice, which won't be used after the
// conversion, to a string, as the following code shows. In this conversion, a duplication of
// the underlying byte sequence is avoided.
//
// This is the implementation adopted by the `String` method of the Builder type supported since
// Go 1.10 in the `string` std package. The size of a byte slice is larger than a string, and their
// internal structures are similar, so the conversion is valid (for main stream Go compilers).
// However, despite the implementation may be safely used in standard packages now, it's not
// recommended to be used in general user code. Since Go 1.20, in general user code, we should try
// to use the implementation which uses the `unsafe.String` function.
//
// Note: when using the just introduced unsafe way to convert a byte slice to a string, please make
// sure not to modify the bytes in the byte slice of the result string still survives.

// The layout of `b` can be larger than string, so it's dangerous
func ByteSliceToStringNew(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// The converse, converting a string to a byte slice in the similar way, is invalid, for the size
// of a string is smaller than a byte slice.
func StringToByteSliceNew(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&s))
}

func Pattern2Usage() {
	type T struct{ a int }
	var t T

	fmt.Printf("%p\n", &t)                          //             0x140000a0020
	println(&t)                                     // 0x140000a0020
	fmt.Printf("%x\n", uintptr(unsafe.Pointer(&t))) // 140000a0020

}

func Pattern3Usafe() {
	type T struct {
		x bool
		y [3]int16
	}

	const (
		M = unsafe.Offsetof(T{}.y)
		N = unsafe.Sizeof(T{}.y[0])
	)

	t := T{x: true, y: [3]int16{123, 456, 789}}

	// uintptr(unsafe.Pointer(&t)) + M + 2*N
	ty2 := (*int16)(unsafe.Add(unsafe.Pointer(&t), M+2*N))
	*ty2 = -333

	fmt.Printf("%+v\n", t) // {x:true y:[123 456 -333]}
}

func Pattern3UsafeWrong() {
	type T struct {
		x bool
		y [3]int16
	}

	const (
		M = unsafe.Offsetof(T{}.y)
		N = unsafe.Sizeof(T{}.y[0])
	)

	t := T{x: true, y: [3]int16{123, 456, 789}}
	p := unsafe.Pointer(&t)

	addr := uintptr(p) + M + 2*N

	// Some other ops...

	// Now the `t` value becomes unused, its memory may be
	// garbage collected at this time. So the following use
	// of the address of t.y[2] may become invalid and dangerous!
	//
	// Another potential danger is, if some ops make the stack grow
	// or shrink here, then the address of `t` might change, so that
	// the address saved in addr will become invalid.
	ty2 := (*int16)(unsafe.Pointer(addr))
	fmt.Println(*ty2)

	// Such bugs are very subtle and hard to detect, which is why the uses of
	// unsafe pointers are dangerous.

}

func WorkingWithReflectRight() {
	// This is the safe case
	p := unsafe.Pointer(reflect.ValueOf(new(int)).Pointer())
	usualP := (*int)(p)
	*usualP = 100

	fmt.Println(*usualP)
}

func WorkingWithReflectWrong() {
	u := reflect.ValueOf(new(int)).Pointer()
	// At this moment, the memory block at the address
	// stored in u might have been collected already.
	p := (*int)(unsafe.Pointer(u))

	*p = 100
	fmt.Println(*p)
}

func ReflectHeadersWorkingA() {
	a := [...]byte{'G', 'o', 'l', 'a', 'n', 'g'}
	s := "Java"

	hdr := (*reflect.StringHeader)(unsafe.Pointer(&s))
	hdr.Data = uintptr(unsafe.Pointer(&a))
	hdr.Len = len(a)

	fmt.Println(s) // Golang

	// New `s` and `a` share the same byte sequence, which makes
	// the bytes in the string `s` become mutable.
	a[2], a[3], a[4], a[5] = 'o', 'g', 'l', 'e'

	fmt.Println(s) // Google
}

func ReflectHeadersWorkingB() {
	a := [6]byte{'G', 'o', '1', '0', '1'}
	bs := []byte("Golang")

	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&bs))
	hdr.Data = uintptr(unsafe.Pointer(&a))
	hdr.Len = 2
	hdr.Cap = len(a)

	fmt.Printf("%s\n", bs) // Go

	bs = bs[:cap(bs)]

	fmt.Printf("%s\n", bs) // Go101

}

func ReflectHeadersWorkingC() {
	var hdr reflect.StringHeader
	hdr.Data = uintptr(unsafe.Pointer(new([5]byte)))

	// Now the just allocated byte array has lose all
	// references and it can be garbage collected now.

	hdr.Len = 5

	s := *(*string)(unsafe.Pointer(&hdr)) // dangerous!
	for i := range len(s) {
		fmt.Printf("%d - %q\n", i, s[i])
	}

	// 0 - '\x00'
	// 1 - '\x00'
	// 2 - '\x00'
	// 3 - '\x00'
	// 4 - '\x00'
}

func ReflectHeadersWorkingD(s string) (bs []byte) {
	strHdr := (*reflect.StringHeader)(unsafe.Pointer(&s))
	sliceHdr := (*reflect.SliceHeader)(unsafe.Pointer(&bs))

	sliceHdr.Data = strHdr.Data
	sliceHdr.Len = strHdr.Len
	sliceHdr.Cap = strHdr.Len

	// Note, when using the just introduced unsafe way to convert a string to a byte slice, please
	// make sure not to modify the bytes in the result byte slice if the string still survives (for
	// demonstration purpose, the above example violates this principle).
	return
}

func ReflectHeadersWorkingDUsage() {
	// str := "Golang"
	// For the official standard compiler, the above
	// line will make the bytes in str allocated on
	// an immutable memory zone.

	// So we use the following line instead.
	str := strings.Join([]string{"Go", "land"}, "")
	s := ReflectHeadersWorkingD(str)

	fmt.Printf("%s\n", s)
	s[5] = 'g'
	fmt.Println(str)

}

func ReflectHeadersWorkingE(bs []byte) (s string) {
	sliceHdr := (*reflect.SliceHeader)(unsafe.Pointer(&bs))
	strHdr := (*reflect.StringHeader)(unsafe.Pointer(&s))

	strHdr.Data = sliceHdr.Data
	strHdr.Len = sliceHdr.Len

	// Note, when using the just introduced unsafe way to convert a string to a byte slice, please make sure not to modify the bytes in the result byte slice if the string still survives (for demonstration purpose, the above example violates this principle).
	return
}

// The Go core development team also realized that the two types are inconvenient and error-prone,
// so the two types have been not recommended any more since Go 1.20 and they have been deprecated
// since Go 1.21. Instead, we should try to use the:
//   - unsafe.String
//   - unsafe.StringData
//   - unsafe.Slice
//   - unsafe.SliceData
//
// functions described earlier in this article.
func AdvanceBytesA() *byte {
	str := "godoc"
	hdr := (*reflect.StringHeader)(unsafe.Pointer(&str))
	return (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(&hdr.Data)) + 2))
}

func AdvanceBytesB() *byte {
	str := "godoc"
	hdr := (*reflect.StringHeader)(unsafe.Pointer(&str))
	return (*byte)(unsafe.Add(unsafe.Pointer(&hdr.Data), 2))
}
