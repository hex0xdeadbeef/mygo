package main

import (
	"strings"
	"unsafe"
)

func main() {
	// commonUnderlyingBytesUsage()
	// getPointerByReflectUsage()
	takeAddressOfFirtstByteUsage()
}

func getBlock(n int) []byte {
	var b strings.Builder
	b.Grow(n)

	return unsafe.Slice(unsafe.StringData(b.String()), len(b.String()))
}
