package main

import (
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"strings"
	"syscall"
	"unsafe"
)

/*
	1. WHAT REALLY IS unsafe.Pointer?

	type Pointer *ArbitraryType

	Unlike a *int, which can only point to int, or *bool, which us restricted to bool values,
	unsafe.Pointer has the liberty to point to any arbitrary type.

	unsafe.Pointer is a special kind of pointer that turns off the usual safety rules. When we use
	it, we're telling the Go compiler: "I know what I'm doing, so trust me.".

	We should be careful, because we're going around the normal safety checks.

	So far, an important point: unsafe.Pointer can be converted to any type of pointer. For example,
	we can change an *int pointer to an unsafe.Pointer, and then change that unsafe.Pointer to a *float64.
*/

func FirstExample() {
	type person struct {
		name string
		age  int
	}

	mp := person{
		name: "John",
		age:  30,
	}

	p := (*int)(unsafe.Add(unsafe.Pointer(&mp), unsafe.Sizeof(mp.name)))
	*p = 10

	fmt.Println(mp.age)

}

func UnsafeTypeSwap() {
	var a int64 = 10
	aP := unsafe.Pointer(&a)

	b := (*float64)(aP)

	// What's the value of `b`? Is it 10 but as a floating-point number?
	//
	// No, it's not.
	// int64 and float64 both use 64 bits, but they store data differently. int64 stores numbers
	// as integers, and float64 stores numbers in a floating-point format.
	fmt.Println(*b) // 5e-323

}

/*
	WHAT IS uintptr?

	A uintptr is basically an integer designed to hold a memory address. It's an unsigned integer
	that can act like a pointer. What does this mean in simpler terms?

	!!! In the world of computers, each memory location is identified by a `number`, and that's what
	uintptr holds, the number that points to a specific memory location.

	NOW, HERE COMES THE TRAP
	If you're familiar with Go, we know it has a GC that clears out memory not being used anymore.
	Specifically, it scans your code to find values that are no logner being referenced. This makes
	life easier because because we don't have to manually manage memory like we would in languages
	such as C.

	var a int64 = 10
	var aPtr = &a

	As long as aPtr has the memory address of `a`, that section of memory will stay in use. This is
	because the Go garbage sees that this memory is still linked to your code, so it won't remove it.

	BUT!

	`uintptr` is a bit different, it's just a number that can hold a memory address. So unlike ohter pointers, it doesn't actually link to the memory in a way the garbage collector understands.

	var a int64 = 10
	uPtr := uintpr(unsafe.Pointer(&a)) // store the number, representing a memory location the
	// var a actually stored in.


	In this example, aPtr is a uintptr holding the memory address of a. But if all pointers that
	directly link to disappear, the GC gets confused. It won't realize that uPtr is supposed to be
	connected to a.

	THE RESULT?
	The garbage collector might just clear out the memory where `a` is stored. If that happens,
	uPtr becomes what's known a "dangling pointer". It would point to an area in memory that's
	been cleared and might be used for something later. If we try using uPtr after this, the behavior of your program becomes unpredictable.
*/

/*
	2. unsafe PACKAGE

	Beside unsafe.Pointer, this package provides us several functions providing low-level information, make us understand more about the internal structure of type.


	FUNC unsafe.Alignof

		The alignment is a concept which is brought from the low-level programming. It returns the
		required alignment of a type and the layout of memory can affect the performance of our
		application.

		For example, int32 has an alignment of 4 bytes. The alignment of a struct depends on the
		field with the largest alignment.

		Knowing these bytes offsets allows us to better understand the struct's memory organization
		and do pointer arithmetuic to access specific fields.
*/

func AlignmentExplanation() {
	type Person struct {
		Name string
		age  int
	}

	var i int32 = 1
	var s [3]int8 = [3]int8{1, 2, 3}
	var p Person = Person{Name: "Bob", age: 20}

	fmt.Println("alignof(int32) =", unsafe.Alignof(i))   // 4
	fmt.Println("alignof([3]int8) =", unsafe.Alignof(s)) // 1
	fmt.Println("alignof(Person) = ", unsafe.Alignof(p)) // 8
}

/*
	FUNC unsafe.Offsetof

		unsafe.Offsetof basically returns the distance in bytes between the start of a struct and
		the start of a specific field. This is useful for understanding struct memory layouts.
*/

func OffsetofExample() {
	type StructA struct {
		A byte  // 1 byte alignment
		B int32 // 4 byte alignment
		C byte  // 1 byte alignment
		D int64 // 8 byte alignment
		E byte  // 1 byte alignment
	}

	a := StructA{}
	_ = a
	fmt.Println("sizeof(A byte) =", unsafe.Offsetof(a.A))
	fmt.Println("sizeof(B int32) =", unsafe.Offsetof(a.B))
	fmt.Println("sizeof(C byte) =", unsafe.Offsetof(a.C))
	fmt.Println("sizeof(D int64) =", unsafe.Offsetof(a.D))
	fmt.Println("sizeof(E byte) =", unsafe.Offsetof(a.E))

	//	sizeof(A byte) = 0
	// sizeof(B int32) = 4
	// sizeof(C byte) = 8
	// sizeof(D int64) = 16
	// sizeof(E byte) = 24
}

/*
	FUNC unsafe.Sizeof

		The unsafe.Sizeof function returns the size in bytes to store some value in memory. It's
		handy for seeing how much space a type takes up.

		"I've tried both unsafe.Sizeof and unsafe.Alignof, and they both return the same value"

		That's mostly true, but not the full picture. For basic types like `int` or `string`, these
		values often do match:

			int:     sizeof = 8   | alignof = 8
			int32:   sizeof = 4   | alignof = 4
			string:  sizeof = 16  | alignof = 8

	While unsafe.Sizeof tells you the total bytes a type occupies, unsafe,Alignof has a different
	role. It indicates the required spacing between elements of that type, depending on the CPU
	architecture.
*/

func AlignofUsageA() {
	var a complex128

	fmt.Println("complex128 sizeof =", unsafe.Sizeof(a), "; alignof = ", unsafe.Sizeof(a))
	// complex128 sizeof = 16 ; alignof =  16

	// In the instance outlined, although a complex128 utilizes 16 bytes of memory, CPU efficiency
	// is optimized when alignment occurs on 8-byte boundaries (given the 64-bit architecture of
	// the system)
}

/*
	USE unsafe.Pointer SAFELY WITH THESE 6 PATTERNS

	The goal of these six methods is to use 'unsafe' pointers in the safest way possible. By
	sticking to these guidlines, we're reducing the risks in an area that's inherently tricky.

	I noticed, that most of the methods advise against storing `uintptr` in a variable, for the
	reasons we've already discussed.

	Pattern 1: Converting a T1 Pointer to a T2 Pointer

		We might remember this from the example where we converted a pointer from *int64 to *float64
		Here's the code for a quick refresher:

		The important thing is to note is that we should only do this kind of conversion if we're
		100% sure that the memory structures of T1 and T2 are compatible.

		Wy should we use this pattern?
		A simple scenario where this comes handy invloves converting an array of one type to
		an array of another type.

	Pattern 2: Turning unsafe.Pointer to uintptr and back, BUT BE CAREFUL!

		We don't need to go into detail here because we already talked about the risks of using
		`uintptr`.

		So, if we want to play it safe, don't change a `uintptr` back to an unsafe.Pointer. We can
		do it, but only in very specific situations.

	Pattern 3: Switch from unsafe.Pointer to uintptr for Arithmetic

		If you remember the black magic part at the start of this article, this is the pattern we're
		talking about. Instead of using the unsafe.Add function (introduced in Go 1.17). We're using
		basic math to get the job done.

	Pattern 4: Using uintptr with Syscalls

		Having a function whose arguments are uintptr is generally something I consider bad
		practice. But Syscall us on of those special cases that require it.

	Pattern 5: `uintptr` returned from reflect.ValueOf(...).(Pointer/UnsafeAddress)

		If you want to get the memory address, we can use the reflect package like this:
			reflect.ValueOf(p).UnsafeAddr()
			reflect.ValueOf(p).Pointer()
		However, both of these methods give us a `uintptr`, which can be risky. The guideline here
		is not to store the output of these methods, but to use them directly to sidestep the
		pitfalls of `uintptr`.

		As of Go 1.19, we have a safer alternative:
			reflect.ValueOf(...).UnsafePointer()

	Pattern 6: Convert from string or slice to reflect.StringHeader and reflect.SliceHeader

		If you dive into how strings and slices work in Go, then you'll find something interesting.
		The pointers of strings or slices don't point directly to the memory of the array they
		represent but actually they point to a different structure known as a Header (StringHeader,
		SliceHeader).

		type StringHeader struct {
			Data uintptr
			Len int
		}

		type SliceHeader struct {
			Data uintptr
			Len int
			Cap int
		}

		It's actually the pointer that points to the string's bytes, while our string pointer just
		refers to this header. So, to get the real data, we have to go through that.

		To make the conversion `string` <-> `[]byte` without any allocation: we cast the *string as
		*T1 to an unsafe.Pointer and then to a *reflect.StringHeader as *T2.


*/

func Pattern1UsageA() {
	var a int64 = 10
	aPtr := unsafe.Pointer(&a)

	fmt.Println(*(*float64)(aPtr)) // 5e-323

}

func Pattern1UsageB() {
	type MyInt int

	x := []MyInt{1, 2, 3}

	// PrintArr(x) // cannot use x (variable of type []MyInt) as []int value in argument to PrintArrcompilerIncompatibleAssign

	y := *(*[]int)(unsafe.Pointer(&x))
	PrintArr(y) // [1 2 3]

	// Remember: x and y are using the same memory, so changing one will affect the other.

}

func PrintArr(arr []int) {
	fmt.Println(arr)
}

func Patter2UsageA() {
	// GOOD: Changing to uintptr and immediately back to unsafe.Pointer
	newPtr := unsafe.Pointer(uintptr(unsafe.Pointer(new(int))))
	_ = newPtr

	// BAD: Saving the `uintptr` value for later usage.
	storedPtr := uintptr(unsafe.Pointer(new(int)))

	// ... some ops

	newPtr = unsafe.Pointer(storedPtr) // possible misuse of unsafe.Pointerunsafeptrdefault
	_ = newPtr
}

func Pattern3UsageA() {
	type Person struct {
		Name  string
		Age   int
		phone string
	}

	p := Person{
		Name:  "Tom",
		Age:   18,
		phone: "123456789",
	}

	agePtr := (*int)(
		// What we do here is briefly the unsafe.Pointer to uintptr, use some
		// quick math to find where the Age field is located, and then change it back
		// to an unsafe.Pointer
		unsafe.Pointer(
			uintptr(unsafe.Pointer(&p)) + unsafe.Sizeof(p.Name),
		),
	)

	*agePtr = 333

	fmt.Println(p.Age)

	// But here's what we mustn't do:
	pRawPtr := uintptr(unsafe.Pointer(&p))
	newAgePtr := unsafe.Pointer(pRawPtr + unsafe.Offsetof(p.Age))

	// Why?
	// Because storing a `uintptr` in a variable is risky business
	_ = newAgePtr

}

func Patter4UsageA() {
	// GOOD:
	fd := rand.Intn(math.MaxInt)
	p := rand.Intn(math.MaxInt)
	n := rand.Intn(math.MaxInt)

	syscall.Syscall(syscall.AF_CHAOS, uintptr(unsafe.Pointer(&fd)), uintptr(unsafe.Pointer(&p)), uintptr(unsafe.Pointer(&n)))

	// BAD:
	newArg := rand.Intn(math.MaxInt)
	u := uintptr(unsafe.Pointer(&newArg))

	syscall.Syscall(syscall.AF_CHAOS, uintptr(unsafe.Pointer(&fd)), uintptr(unsafe.Pointer(&p)), u)
}

func Patter5UsageA() {
	// GOOD:
	p := (*int)(unsafe.Pointer(reflect.ValueOf(new(int)).Pointer()))
	*p = 42

	fmt.Println(*p)

	// BAD
	u := reflect.ValueOf(new(int)).Pointer()
	newP := (*int)(unsafe.Pointer(u))

	*newP = 42

	fmt.Println(*newP)

	// Since Go 1.19 we have reflect.ValueOf(...).UnsafePointer()
	up := reflect.ValueOf(new(int)).UnsafePointer()
	*(*int)(up) = 42
	fmt.Println(*(*int)(up))
}

func Pattern6UsageA() {
	str := strings.Join([]string{"hel", "lo"}, "")
	var b []byte

	// Reinterpret the string's header as a slice header
	strHeader := (*reflect.StringHeader)(unsafe.Pointer(&str))
	sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&b))

	// Copy the data pointer and length from the string header to the slice header
	sliceHeader.Data = strHeader.Data
	sliceHeader.Len = strHeader.Len - 1
	sliceHeader.Cap = strHeader.Len

	// Now `b` points directly to the raw bytes of the original string without any new allocation

	fmt.Println(b)
	fmt.Println(string(b))
}

func Pattern6UsageB() {
	str := strings.Join([]string{"hel", "lo"}, "")

	b := unsafe.Add(unsafe.Pointer(unsafe.StringData(str)), 3)
	fmt.Println(unsafe.String((*byte)(b), 2))
}
