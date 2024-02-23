package chapter9

import "fmt"

var x, y int

func DeterministicOutput() {
	go func() {
		x = 1                   // Sx
		fmt.Printf("y = %d", y) // Oy
	}()

	go func() {
		y = 1                   // Sy
		fmt.Printf("x = %d", x) // Ox
	}()
}
