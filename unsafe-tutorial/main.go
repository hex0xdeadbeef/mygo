package main

import (
	"fmt"
	"math/rand"
	"reflect"
	"unsafe"
)

func main() {
	// UintptrIntroduction()
	// StructFieldsManipulation()
	// IteratingOverSlice()
	// FastConversionFromBytesToString()

	// OffsetofAlignofSizeofUsageA()
	// OffsetofAlignofSizeofUsageB()

	// UsingUnsafePackageFuncs()
	// UnsafePointerOfUnsafePointer()

	// ConversionsBetweenSliceTypes()
	// Pattern2Usage()
	// Pattern3Usafe()
	// Pattern3UsafeWrong()

	// WorkingWithReflectRight()
	// WorkingWithReflectWrong()

	// ReflectHeadersWorkingA()
	// ReflectHeadersWorkingB()
	// ReflectHeadersWorkingC()
	ReflectHeadersWorkingDUsage()
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
