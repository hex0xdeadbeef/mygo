package main

import (
	"fmt"
	"unsafe"
)

/*
	INTRO
	It's important to note that Go us a garbage-collected language.
	Traditionally, `unsafe` provided a backdoor to Go's type safety,
	allowing for low-level memory manipulation. While
	powerful, it's often regarded as dangerous due to its potential risks.

	`
	Package unsafe contains ops that step around tge type safety of Go programns

	Packages that import unsafe may be non-portable and aren't protected by tge Go 1 compatibility
	guidelines.
	`

	If we visit scr/reflect/value.go, we'll notice that reflect.
	StringHeader is deprecated in favor of unsafe.String & unsafe.
	StringData and reflect.SliceHeader is deprecated in favor if
	unsafe.Slice and unsafe.SliceData.


	SOME NOTES ON unsafe.Pointer and uintptr.

	`
	Pointer represents a pointer ti an arbitrary type. There are four special operations available for type Pointer that aren't available for other types:

	- A usual pointer value of any type can be converted to unsafe.Pointer
	- An unsafe.Pointer can be converted to a pointer value of any type
	- A uintptr can be converted to an unsafe.Pointer
	- An unsafe.Pointer can be converted to a uintptr.
	`

	Another important note is that Go's standard type pointers don't have the ability to participate in pointer arithmetic, which is why we
	convert to and from a uintptr when we need to do those ops. A uinptr is
	simply an integer type that is large enough to contain the bit pattern of any pointer.

	`
	A uintptr is an integer, not a reference. Converting an unsafe.Pointer to a uintptr creates an integer value with no pointer semantics. Even if a uintptr holds the address of some object, the garbage collector will not update that uintptr's value if the object moves, nor will that uintptr keep the object from being reclaimed.
	`

	The main takeaway: it's generally invalid to store an address as a uintptr in a variable before using it. This is due to the fact that a uintptr is simple integer with no pointer semantics and objects can both change underneath you or be freed all together by the time we use them.
*/

func UintptrIntroduction() {
	// Allocate an integger, get pointer to it
	x := 10
	xP := &x

	// INVALID: Get a uintptr of the address x. Don't do this.
	xUinPtr := uintptr(unsafe.Pointer(xP))

	// At this point, `x` is unused and so could be garbage collected.
	// If that happens, we then have an uinptr (integer) that when casted back to an
	// unsafe.Pointer, points to some invalid piece of memory.

	fmt.Println(*(*int)(unsafe.Pointer(xUinPtr)))
}

// Manipulating Struct Fields

// p is a simple pretty print func
func p(v any) { fmt.Printf("%+v\n", v) }

type user struct {
	name    string
	age     int
	animals []string
}

// All this can be applied to private structs / struct members, however there are some
// implementation differences in terms of how to access and calculate things when we don't have
// direct access to the members.
func StructFieldsManipulation() {
	// Declare zero value 'user' struct and print its contents
	var u user

	p(u) // {name: age:0 animals:[]}

	// Retrieve an umsafe.Pointer to 'u', which points to the first
	// member of the struct - 'name' which is a string. Then we cast
	// the unsafe.Pointer to a string pointer. This is allows us to
	// manipulate the memeory pointed at as a string type.s
	uNamePtr := (*string)(unsafe.Pointer(&u))
	*uNamePtr = "Dmitriy"

	p(u) // {name:Dmitriy age:0 animals:[]}

	// Here we have a similar situation in that we want to get a pointer to a struct
	// member. This time it's the second member, so we need to calculate the address
	// within the struct by using offsets. The general idea is that we need to add the
	// size of the 'age' member. Finally, we gen an unsafe.Pointer from 'unsafe.Add' and cast
	// it to an '*int'.
	agePtr := (*int)(unsafe.Add(unsafe.Pointer(&u), unsafe.Offsetof(u.age)))
	*agePtr = 21
	p(u) // {name:Dmitriy age:21 animals:[]}

	// One other thing ti note is that 'unsafe.Add' was added to the unsafe package in Go 1.17
	// along with 'unsafe.Slice' prior to having 'Add' (which conviniently returns an 'unsafe.Pointer') we had to approach this slightly differently, convering and adding uintptrs and casting them back to unsafe.Pointers:
	anotherAgePtr := (*int)(unsafe.Pointer(uintptr(unsafe.Pointer(&u)) + unsafe.Offsetof(u.age)))
	*anotherAgePtr = 22
	p(u) // {name:Dmitriy age:22 animals:[]}

	// Here we're working with something a bit different. First we add a slice of animals to
	// the user struct we've been working with.
	u.animals = []string{"missy", "ellie", "toby"}

	// Now we want to get a pointer to the second slice element and make a change to it. We use
	// the new unsafe func here called 'SliceData'. This will return a pointer to the underlying
	// array of the argument slice. Now that we have a pointer to the array, we can add the size
	// of one string to the pointer to get the address of the second element. This means we could
	// say  2 * unsafe.Sizeof("") to get to the third element in this example if that is helpful
	// for visualizing.
	secondAnimal := (*string)(unsafe.Add(
		unsafe.Pointer(unsafe.SliceData(u.animals)),
		unsafe.Sizeof("")),
	)
	*secondAnimal = "Bobik"

	p(u) // {name:Dmitriy age:22 animals:[missy Bobik toby]}

}

// Iterating Over a Slice
func IteratingOverSlice() {
	fruits := []string{"apples", "oranges", "bananas", "kansas"}

	// Get an unsafe.Pointer to the slice data
	start := unsafe.SliceData(fruits)

	// Get the size of an item in the slice, Another way to do this
	// could also be written as 'size := unsafe.Sizeof("")' here as
	// we know the items are strings.
	size := unsafe.Sizeof(fruits[0])

	// Get the size of an item in the slice and print data in each time.
	// Arrays in Go are stored contiguously and sequentially in memory,
	// so we are able to directly access each item through indexing:
	//
	// base_address + (index * size_of_elem)
	//
	// In each iteration, we take the pointer to the array data ('start')
	// and add the (index * size_of_an_elem) to get the address of each item along
	// the block of memory. Finally, we cast the item to a '*string' to print it.

	for i := range len(fruits) {
		fmt.Println(*(*string)(unsafe.Add(unsafe.Pointer(start), uintptr(i)*size)))
	}

	// apples
	// oranges
	// bananas
	// kansas
}

// Strings and Bytes
func FastConversionFromBytesToString() {
	// When converting between strings and byte slices in Go, the standard library's string()
	// and []byte(str) are commonly used for their safety and simplicity. These methods create a
	// new copy of the data, ensuring that the original data remains immutable and that the type
	// safety is maintained. However, this also means that every conversion involves memory
	// allocation and copying, which can be a performance concern in certain high-efficiency
	// scenarios.

	// StringToByteSlice
	myStr := "Hello, World!"
	bytesOfStr := unsafe.Slice(unsafe.StringData(myStr), len(myStr))
	p(bytesOfStr) // [72 101 108 108 111 44 32 87 111 114 108 100 33]

	// ByteSliceToString
	myBytes := []byte{
		115, 111, 32, 109, 97, 110, 121, 32, 110,
		101, 97, 116, 32, 98, 121, 116, 101, 115,
	}

	myStr = unsafe.String(unsafe.SliceData(myBytes), len(myBytes))
	p(myStr) // so many neat bytes

	// While unsafe provides a high performance alternative for string-byte slice conversions,
	// it should be used judiciously. Note that we should never modify the underlying bytes of
	// string after these conversions to ensure data intergrity and avoid unexpected behavior.
	// The benefits of using unsafe for these conversions must be weighed against the increased ?
	// complexity and potential risks.
}

// Conversion optimization
// NOTE: b mustn't be changed after comparison
func StringConversionOptimization() string {
	b := []byte("Hello!")
	return unsafe.String(unsafe.SliceData(b), len(b))
}

func StringConversionDefault() string {
	b := []byte("Hello!")
	return string(b)
}
