package chapter5

import (
	"fmt"
)

func minElementInArray() {
	array := []int{48, 96, 86, 68,
		57, 82, 63, 70,
		37, 34, 83, 27,
		19, 97, 9, 17,
	}

	min := array[0]
	for i := 1; i < len(array)-1; i++ {

		if array[i] > min {
			min = array[i+1]
		}
	}

	fmt.Println(min)
}
