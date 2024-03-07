package sub1

import (
	"fmt"
	"unsafe"
)

func SizeofUsing() {
	fmt.Println(unsafe.Sizeof(float64(0)))
}

type x struct {
	a bool
	b int16
	c []int
}
