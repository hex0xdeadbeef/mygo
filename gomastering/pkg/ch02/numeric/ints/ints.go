package ints

import "fmt"

func DivisionDifference() {
	var m, n float64
	var k, x int = 100, 200
	m = 1.223

	fmt.Println("m, n: ", m, n)

	y := 4 / 2.3
	fmt.Printf("y: %g", y)

	divFloat := float64(k) / float64(x)
	fmt.Println("divFloat:", divFloat)
	fmt.Printf("Type of divFloat: %T", divFloat)
}

