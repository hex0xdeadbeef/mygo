package main

import (
	"fmt"
	"math/rand"
	"reflect"
	"unsafe"
)

func main() {
	// StringUsageA()
	// StringUsageB()

	// SliceUsageA()
	// fmt.Printf("%b", Float64Bits(3.14134))

	// GetFieldWithOffset()
	// GetSecondFieldVal()
	// RangeOverSliceWithUinptr()

	// ConversionBetweenTypes()
	// AccessingStringHeader()

	// AdvancePointerFurther()
	// AligningFields()
	// PrivateFieldsAccessing()
	// StringConversionOptimization()

	// MisuseOfUinptr()

	// userManipulation()
	// IteratingOverASlice()

	FastConversions()
}

func StringUsageA() {
	s := "Hello, World!"
	sd := unsafe.StringData(s)
	s = unsafe.String(sd, 5)
	fmt.Println(s)
}

func StringUsageB() {
	s := ""
	sd := unsafe.StringData(s)
	fmt.Println(sd == nil)

	sd = new(byte)
	*sd = 'a'

	s = unsafe.String(sd, 1)

	fmt.Println(s)
}

func SliceUsageA() {
	s := make([]int, 5, 10)
	sd := unsafe.SliceData(s)

	s = unsafe.Slice(sd, 3)

	fmt.Println(s, len(s), cap(s))

	s = nil
	sd = unsafe.SliceData(s)
	fmt.Println(sd == nil)

	s = unsafe.Slice(sd, 0)
	fmt.Println(s)
	fmt.Println(s == nil)
}

func Float64Bits(f float64) uint64 {
	return *(*uint64)(unsafe.Pointer(&f))
}

type Numbers struct {
	a int
	b int
	c int
}

func GetFieldWithOffset() {
	ns := Numbers{a: 1, b: 2, c: 3}

	ff := unsafe.Pointer(uintptr(unsafe.Pointer(&ns)) + unsafe.Offsetof(ns.c))
	val := *(*int)(ff)

	fmt.Println(val)
}

func GetSecondFieldVal() {
	ns := Numbers{a: 1, b: 2, c: 3}

	uptr := (uintptr)(unsafe.Pointer(&ns))

	uptr += unsafe.Sizeof(int(0))

	sf := *(*int)(unsafe.Pointer(uptr))
	fmt.Println(sf)
}

func RangeOverSliceWithUinptr() {
	s := make([]int, 100)
	for i := range len(s) {
		s[i] = rand.Intn(101)
	}

	for i := range len(s) {
		e := unsafe.Pointer(uintptr(unsafe.Pointer(&s[0])) + uintptr(i)*unsafe.Sizeof(s[0]))

		fmt.Println(*(*int)(e))
	}
}

// Converting between different pointer types

func ConversionBetweenTypes() {
	var f float64 = 42.0

	p := (*uint64)(unsafe.Pointer(&f))

	fmt.Println(*p)
	fmt.Printf("Bits of the value %g: %b\n", f, *p)
}

// Accessing the internal representation of a string
func AccessingStringHeader() {
	type StringHeader struct {
		Data uintptr
		Len  int
	}

	s := "Hello, World!"

	stringHeader := *(*StringHeader)(unsafe.Pointer(&s))
	stringHeader.Data += 5 * unsafe.Sizeof(byte(0))
	stringHeader.Len -= 5
	fmt.Println(stringHeader)

	s = *(*string)(unsafe.Pointer(&stringHeader))

	fmt.Println(s)
}

// Pointer arithmetics
func AdvancePointerFurther() {
	val := 4
	uVal := uintptr(unsafe.Pointer(&val))

	uVal += unsafe.Sizeof(val)

	val = *(*int)(unsafe.Pointer(uVal))

	fmt.Println(val)
}

// Aligning of struct fields
func AligningFields() {
	type S struct {
		A int8
		B int64
	}

	s := S{}

	fmt.Println(unsafe.Alignof(s.A))
	fmt.Println(unsafe.Alignof(s.B))
}

func PrivateFieldsAccessing() {
	type SecretStruct struct {
		privateField string
	}

	secret := SecretStruct{privateField: "This is a secret!"}

	field := reflect.ValueOf(&secret).Elem().FieldByName("privateField")

	ptr := unsafe.Pointer(field.UnsafeAddr())
	realPtr := (*string)(ptr)

	*realPtr = "This is not a secret!"

	fmt.Println(secret.privateField)
}

// Conversion optimization
// NOTE: b mustn't be changed after comparison
func StringConversionOptimization() {
	b := []byte("Hello!")

	s := *(*string)(unsafe.Pointer(&b))

	b[0] = 'h'

	_ = s
}

func StringConversionDefault() {
	b := []byte("Hello!")
	s := string(b)

	_ = s
}

// uintptr doesn't have the semantics of a usual pointer
// we cannot count on correctness with pointer arithmetic,
// because the GC can move pointer away
func MisuseOfUinptr() {
	// Allocate integer, get a pointer to it
	x := 10
	xPtr := &x
	y := 42
	_ = y

	// Get a uintptr of the address of x. Do NOT do this.
	xUinptr := uintptr(unsafe.Pointer(xPtr))

	// At this point, `x` is unused and so could be garbage collected
	// If that happen, we have an uinptr (integer) that when
	// casted back to an unsafe.Pointer, points to some invalid
	// piece of memory
	fmt.Println(*(*int)(unsafe.Pointer(xUinptr)))
}

// Manipulating Struct Fields
func p(a any) { fmt.Printf("%+v\n", a) }

type user struct {
	name    string
	age     int
	animals []string
}

func userManipulation() {
	var u user
	p(u)

	// Retrieve an unsafe.Pointer to 'u', which points to the first
	// member of the struct - 'name' - which is a string. Then we
	// cast the unsafe.Pointer to a string pointer. This allows us
	// to manipulate the memory pointed at as a string type.
	uNamePtr := (*string)(unsafe.Pointer(&u))
	*uNamePtr = "Dmitriy"

	p(u)

	// Here we have a similar situation in that we want to get a pointer
	// to a struct member. This time it's the second member, so we need
	// to calculate the address within the struct by using offsets. The
	// general idea is that we need to add the size of 'name' to the
	// address of the struct to get the start of the 'age' member.
	// Finally we get an unsafe.Pointer from 'unsafe.Add' and cast it to
	// an *int
	age := (*int)(unsafe.Add(unsafe.Pointer(&u), unsafe.Offsetof(u.age)))

	*age = 21

	p(u)

	u.animals = []string{"missy", "ellie", "toby"}
	// Now we want to get a pointer to the second slice element and make
	// a change to it. We use the new unsafe func here called 'SliceData'.
	// This will return a pointer to the underlying array of the argument
	// slice. Now that we have a pointer to the array, we can add the size
	// of one string to the pointer to get the address of the second element.
	// This means we could say 2*unsafe.Sizeof("") to get the third element
	// in this example if that is helpful for visualizing.

	secondAnimal := (*string)(unsafe.Add(
		unsafe.Pointer(unsafe.SliceData(u.animals)),
		unsafe.Sizeof("")),
	)

	p(u)

	*secondAnimal = "Belly"

	p(u)

}

// All this can be applied to private structs / struct members, however there are some
// implementation differences in terms of how to access and calculate things when we don't have
// direct access to the members.
func IteratingOverASlice() {
	fruits := []string{"apples", "oranges", "bananas", "kansas"}

	start := unsafe.Pointer(unsafe.SliceData(fruits))

	size := unsafe.Sizeof(fruits[0])

	for i := range len(fruits) {
		p(*(*string)(unsafe.Add(start, uintptr(i)*size)))
	}
}

// Strings and Bytes
func FastConversions() {
	myStr := "neato burrito"
	byteSlice := unsafe.Slice(unsafe.StringData(myStr), len(myStr))

	p(byteSlice)

	// ByteSliceToString
	myBytes := []byte{
		115, 111, 32, 109, 97, 110, 121, 32, 110,
		101, 97, 116, 32, 98, 121, 116, 101, 115,
	}

	str := unsafe.String(unsafe.SliceData(myBytes), len(myBytes))
	fmt.Println(str)

	// While unsafe provides a high-performance alternative for string-byte slice conversions
	// it should be used judiciously. Note that we should never modify the underlying bytes of
	// string after these conversions to ensure data integrity and avoid unexpected behavior. The
	// benefits of using unsafe for these conversions must be weighed against the increased
	// complexity and potential risks.
}
