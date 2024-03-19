package waitgroup

import (
	"fmt"
	"sync"
	"time"
)

// WaitGroup is used when we need to wait for a set of concurrent operations to complete when
// 1) We don't care about the result.
// OR
// 2) We have other means of collecting their results

// The operation Wait() blocks the goroutine it's called from.

var wg sync.WaitGroup

func Using() {
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("1st goroutine is sleeping...")
		time.Sleep(1 * time.Second)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("2nd goroutine is sleeping...")
		time.Sleep(1 * time.Second)
	}()

	wg.Wait()
	fmt.Println("All the goroutines are complete")
}

func FuckUp() {
	go func() {
		// wg.Add(1) // It's the race condition, because the multiple goroutines try to write to the WaitGroup
		defer wg.Done()
		fmt.Println("1st goroutine is sleeping...")
		time.Sleep(1 * time.Second)
	}()
	go func() {
		// wg.Add(1) // It's the race condition, because the multiple goroutines try to write to the WaitGroup
		defer wg.Done()
		fmt.Println("2nd goroutine is sleeping...")
		time.Sleep(1 * time.Second)
	}()

	// We could reach this statement before any goroutine executes wg.Add(1)
	wg.Wait()
	fmt.Println("All the goroutines are complete")
}

func MultipleAdd() {
	const (
		numGreeters = 5
	)

	var (
		hello = func(wg *sync.WaitGroup, id int) {
			defer wg.Done()
			fmt.Printf("Hello from %v!\n", id)
		}

		wg sync.WaitGroup
	)

	wg.Add(numGreeters)
	for i := 0; i < numGreeters; i++ {
		go hello(&wg, i)
	}

	wg.Wait()
}
