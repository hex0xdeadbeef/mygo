package chapter5

import (
	"fmt"
)

func creatingASlice() {
	var x []float64
	fmt.Println(x)

	y := make([]float64, 5)
	fmt.Println(y)

	z := make([]float64, 5, 10)
	fmt.Println(z)

	firstArray := [5]float64{1, 2, 3, 4, 5}
	a := firstArray[0:2]
	fmt.Println(a)
	b := firstArray[0:]
	fmt.Println(b)
	c := firstArray[:len(firstArray)-1]
	fmt.Println(c)
	d := firstArray[:]
	fmt.Println(d)

}

func appendingElementsIntoASlice() {
	firstSlice := []int{1, 2}
	fmt.Println(firstSlice, len(firstSlice), cap(firstSlice))

	secondSlice := append(firstSlice, 3)
	fmt.Println(secondSlice, len(secondSlice), cap(secondSlice))
}

func copyingSlices() {
	firstSlice := []int{1, 2, 3}
	fmt.Println(len(firstSlice), cap(firstSlice))
	secondSlice := make([]int, 2)

	copy(secondSlice, firstSlice)

	fmt.Println(secondSlice)
}
