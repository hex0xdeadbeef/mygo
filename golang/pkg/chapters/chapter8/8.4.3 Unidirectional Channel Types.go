package chapter8

import "fmt"

func ThreeArityUnidirPipeline(boundary int) {
	naturals := make(chan int)
	squares := make(chan int)
	done := make(chan bool)

	// The compiler implicitly converts the channels passed
	go generator(naturals, boundary)
	go squarer(naturals, squares)
	go printer(squares, done)

	<-done
}

func generator(naturals chan<- int, boundary int) {
	defer close(naturals)

	for i := 1; i <= boundary; i++ {
		naturals <- i
	}
}

func squarer(naturals <-chan int, squares chan<- int) {
	defer close(squares)

	for naturalNumber := range naturals {
		squares <- naturalNumber * naturalNumber
	}
}

func printer(squares <-chan int, done chan<- bool) {
	for square := range squares {
		fmt.Printf("%d ", square)
	}

	done <- true
}
