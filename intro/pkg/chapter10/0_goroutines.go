package chapter10

import (
	"fmt"
	"math/rand"
	"time"
)

func createGoroutine(number int) {
	for i := 0; i < 10; i++ {
		delay := time.Duration(rand.Intn(100))
		time.Sleep(time.Millisecond * delay)
		fmt.Printf("Goroutine №%d: %d\n", number, i)
	}
}

func Firs_Using_Goroutines() {
	for i := 0; i < 10; i++ {
		go createGoroutine(i)
	}

	var input string
	fmt.Scanln(&input)
}
