package chapter8_9

import (
	"fmt"
	mymath "goprs/pkg/chapter8_9/math"
	"math/rand"
)

// Demonstrates how packages interact with each other
func Package_Function_Using() {
	var randSlice []float64
	for i := 0; i < 10; i++ {
		randSlice = append(randSlice, rand.Float64())
	}

	fmt.Println(mymath.Average(randSlice))
}
