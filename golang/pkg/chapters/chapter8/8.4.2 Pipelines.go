package chapter8

import (
	"fmt"
)

func ThreeArityPipeline() {
	naturals := make(chan int)
	squares := make(chan int)

	/*
		There's no confidence that the Counter will start before Squarer and vice versa
	*/

	// Counter
	go func() {
		/*
			Closing the naturals channel above would cause the squarer's loop to spin as it
			receives a never-ending stream of zero values, and to send these zeros to the printer.
		*/
		for naturalNumber := 1; naturalNumber <= 5; naturalNumber++ {
			naturals <- naturalNumber
		}
		// The receiver will be given a stream of the zero value of type int finitely
		close(naturals)
	}()

	// Squarer
	go func() {
		for {
			naturalNumber, ok := <-naturals
			// If we are given a value from the drained steam, close the squares channel
			if !ok {
				close(squares)
				break
			}
			squares <- naturalNumber * naturalNumber
		}
	}()

	// Main goroutine is the Printer
	for {
		square, ok := <-squares
		if !ok {
			break
		}
		fmt.Println(square)
	}
}

func ThreeArityLoopPipeline() {
	naturals := make(chan int)
	squares := make(chan int)

	/*
		There's no confidence that the Counter will start before Squarer and vice versa
	*/

	// Counter
	go func() {
		// The receiver will be given a stream of the zero value of type int finitely
		defer close(naturals)
		/*
			Closing the naturals channel above would cause the squarer's loop to spin as it
			receives a never-ending stream of zero values, and to send these zeros to the printer.
		*/
		for naturalNumber := 1; naturalNumber <= 15; naturalNumber++ {
			naturals <- naturalNumber
		}
	}()

	// Squarer
	go func() {
		defer close(squares)
		// Channel will be closed after the last given element
		for natural := range naturals {
			squares <- natural * natural
		}

	}()

	// Main goroutine is the Printer
	for square := range squares {
		fmt.Printf("%d ", square)
	}

}
