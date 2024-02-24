package chapter9

import (
	"fmt"
	"sync"
)

var (
	wg   sync.WaitGroup
	x, y int
)

func DeterministicOutput() {
	wg.Add(1)
	go func() {
		defer wg.Done()
		x = 1                    // Sx
		fmt.Printf("y = %d ", y) // Oy
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		y = 1                    // Sy
		fmt.Printf("x = %d ", x) // Ox
	}()

	wg.Wait()
}
