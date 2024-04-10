package chapter6

import (
	"fmt"
	"math/rand"
)

func main() {
	// Simple functions
	fmt.Println("Simple functions")

	var firstArray []float64 = []float64{98, 93, 77, 82, 83}
	fmt.Println(average(firstArray))
	fmt.Println()

	// Functions that return some values
	fmt.Println("Functions that return some values")

	x, y, err := functionWithMultipleReturnableArguments()
	fmt.Println(x, y, err)
	fmt.Println()

	// Functions that use multiple argument
	fmt.Println("Functions that use multiple argument")

	fmt.Println(functionWithModifiedFormOfMultipleArguments(0))
	fmt.Println(functionWithModifiedFormOfMultipleArguments(1, 2, 3, 4, 5))
	sliceOfIntegers := []int{1, 1, 5}
	anotherSliceOfIntegers := []int{}
	fmt.Println(functionWithModifiedFormOfMultipleArguments(sliceOfIntegers...))
	fmt.Println(functionWithModifiedFormOfMultipleArguments(anotherSliceOfIntegers...))
	fmt.Println()

}

func average(slice []float64) float64 {
	total := 0.0
	for _, element := range slice {
		total += element
	}

	return total / float64(len(slice))
}

func functionWithMultipleReturnableArguments() (num1, num2 float64, err error) {
	return float64(rand.Float32()), float64(rand.Float32()), nil
}

func functionWithModifiedFormOfMultipleArguments(slice ...int) int {
	sum := 0
	for _, element := range slice {
		sum += element
	}

	return sum
}
