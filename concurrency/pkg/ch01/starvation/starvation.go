package starvation

import (
	"fmt"
	"sync"
	"time"
)

/*
Starvation is any situation where a concurrent process cannot get all the resources it needs to perform work.

Starvation usually implies that there are one or more greedy concurrent processes that unfairly preventing one or more concurrent processes from accomplishing work as efficiently as
possible, or maybe at all.
*/

var (
	wg         sync.WaitGroup
	sharedLock sync.Mutex
)

const (
	timeOfRun = 1 * time.Second
)

func Starvation() {
	greedyWorker := func() {
		defer wg.Done()

		var count int
		for begin := time.Now(); time.Since(begin) <= timeOfRun; {
			sharedLock.Lock()
			time.Sleep(3 * time.Nanosecond)
			sharedLock.Unlock()
			count++
		}

		fmt.Printf("Greedy worker was able to execute %v work loops\n", count)
	}

	politeWorker := func() {
		defer wg.Done()

		var count int
		for begin := time.Now(); time.Since(begin) <= timeOfRun; {
			sharedLock.Lock()
			time.Sleep(1 * time.Nanosecond)
			sharedLock.Unlock()

			sharedLock.Lock()
			time.Sleep(1 * time.Nanosecond)
			sharedLock.Unlock()

			sharedLock.Lock()
			time.Sleep(1 * time.Nanosecond)
			sharedLock.Unlock()

			count++
		}

		fmt.Printf("Polite worker was able to execute %v work loops\n", count)
	}

	wg.Add(2)
	go greedyWorker()
	go politeWorker()

	wg.Wait()
}
