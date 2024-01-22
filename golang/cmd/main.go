package main

import (
	"fmt"
	a "golang/pkg/chapter4/b_slices"
)

func main() {
	fmt.Print(string(a.ReverseUTF8([]byte("異х異ф"))))
}
