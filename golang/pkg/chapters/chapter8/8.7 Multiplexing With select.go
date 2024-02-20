package chapter8

import (
	"fmt"
	"os"
	"time"
)

func RocketCountdown() {
	fmt.Println("Commencing countdown. Press return to abort.")

	abort := make(chan struct{})
	go func() {
		os.Stdin.Read(make([]byte, 1)) // read a single byte
		abort <- struct{}{}
	}()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for countdown := 3; countdown > 0; countdown-- {
		fmt.Println(countdown)
		select {
		case <-ticker.C:
		case <-abort:
			fmt.Println("Launch aborted!")
			return
		}
	}
	launch()

}

func launch() {
	fmt.Println("The rocket launched!")
}

func ForSelect() {
	ch := make(chan int, 1)

	for i := 0; i < 10; i++ {
		select {
		case value := <-ch:
			fmt.Println(value)
		case ch <- i:
		}
	}
}

func NondeterministicForSelect() {
	ch := make(chan int, 2)

	for i := 0; i < 10; i++ {
		// Each case after the first iteration have the same chances to be executed
		select {
		case value := <-ch:
			fmt.Println(value)
		case ch <- i:
		}
	}
}
