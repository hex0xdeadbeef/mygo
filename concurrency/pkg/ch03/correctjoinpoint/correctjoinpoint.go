package correctjoinpoint

import (
	"fmt"
	"sync"
)

// Go uses the mechanism that's named "fork-join". This means that at the random point a child goroutine separates from the parent goroutine and subsequently joins back to it.

var wg sync.WaitGroup

func Using() {
	wg.Add(1)

	// Separate the goroutine from the parent so that it'll be scheduled and run.
	go func() {
		defer wg.Done()
		fmt.Println("Hello ")
	}()

	wg.Wait()
}
