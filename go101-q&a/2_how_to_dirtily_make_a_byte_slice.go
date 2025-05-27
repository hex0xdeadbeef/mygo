package main

import (
	"fmt"
	"strings"
	"unsafe"
)

/*
	How to dirtily make a byte slice (without zeroing its byte elements)?

	With the 1.21+ Go implementation, now we get a chance to create byte slice without initializing
	their elements to zero (through achieving this functionality requires the use of unsafe funcs)
*/

const Len = 20_000

func init() {
	fmt.Println("=============== Len =", Len)
}

func MakeDirtyByteSlice(n int) []byte {
	var b strings.Builder
	b.Grow(n)

	return unsafe.Slice(unsafe.StringData(b.String()), n)
}
