package slices

import "fmt"

func SliceCapacityDoubling() {
	s := make([]int, 4)
	s = append(s, -1)

	s = append(s, 1, 2, 3, 4)
	fmt.Println(len(s), cap(s))

	s = append(s, []int{1, 2, 3, 4}...)

	fmt.Println(len(s), cap(s))
}
