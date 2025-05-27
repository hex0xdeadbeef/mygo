package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

/*
	How to take addresses of the bytes in a string?
*/

/*
	1. Use unsafe.StringData

	Prior to Go 1.22 (and since 1.20) the only way to take the addresses of string bytes through the unsafe way.
*/

var h = "Hello "
var w = "world!"

func commonUnderlyingBytes(x, y string) bool {
	return &[]byte(x)[4] == &[]byte(y)[4]
}

func commonUnderlyingBytesUsage() {
	var s = h + w

	var p = unsafe.StringData(s[4:])

	fmt.Println(p, string(*p))
}

/*
	2. Use reflect.Value.UnsafePointer

	Since Go 1.23, the reflect.Value.UnsafePointer method has supported string values. So we can
	use it to get addresses of string bytes.
*/

func getPointerByReflectUsage() {
	var s = h + w
	var p = (*byte)(reflect.ValueOf(s[4:]).UnsafePointer())
	fmt.Println(p, string(*p))
}

/*
	3. Convert string to a byte slice, then take addresses of the byte slice

	Since Go toolchain 1.22 an optimization (https://github.com/golang/go/issues/2205) has been made so that memory allocation and bytes duplication are not needed anymore in sime []byte(aString) conversions (if the bytes in the conversion result are proven to be never modified)

	We can make use of this optimization to take addresses of the bytes of a string, by converting the string to a byte slice then indirectly take addresses of the bytes in the byte slice.
*/

func takeAddressOfFirtstByte(s string) *byte {
	// There's no new allocation because there won't be changes to the slice
	return &[]byte(s)[0]
}

func takeAddressOfFirtstByteUsage() {
	var s = h + w

	var p1 = &[]byte(s)[4]
	fmt.Println(p1, string(*p1))

	var p2 = &[]byte(s)[4]
	*p2 = 'X' // modify the byte
	fmt.Println(p2, string(*p2))

	fmt.Println(p1 == p2) // false

	// There won't be new allocation because there won't be a changes to the string
	var s2 = s
	fmt.Println(commonUnderlyingBytes(s, s2)) // true
}
