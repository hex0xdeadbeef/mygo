package sub2

import (
	"fmt"
	"unsafe"
)

var x struct {
	a bool
	b int16
	c []int
}

func ChangeValueThroughOffset() {

	// Convert into (*int16) type (the type of b)
	// Convert the numeric address value into "unsafe.Pointer"
	// Take the address of x
	// Adds the shift from the start of the address
	pb := (*int16)(unsafe.Pointer(uintptr(unsafe.Pointer(&x)) + unsafe.Offsetof(x.b)))
	*pb = 100

	fmt.Println(x.b)
}

func SubtleVariableIntroducing() {
	// NOTE: subtly incorrect
	tmp := uintptr(unsafe.Pointer(&x)) +
		unsafe.Offsetof(x.b)
	pb := (*int16)(unsafe.Pointer(tmp))

	*pb = 42
}
