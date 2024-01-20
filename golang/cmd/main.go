package main

import (
	"fmt"
	s "golang/pkg/chapter4/b_slices"
)

func main() {
	sl := []int{1, 2, 3, 4, 5}
	fmt.Println(s.RotationRightSingle(sl, -3))
}
