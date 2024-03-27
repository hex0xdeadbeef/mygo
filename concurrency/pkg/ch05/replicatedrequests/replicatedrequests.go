package replicatedrequests

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

/*
For some applications, receiving a response as quickly as possible is the top priority. For example, maybe the application is servicing user's HTTP request, or retrieving a replicated blob of
data. In these cases we can make a trade-off: we can replicate the request to multiple handlers (whether those be goroutines, processes, or servers), and one of them will return faster than the
others ones and then we can immediately return the result.

The downside of this approach is we'll have to utilize resources to keep multiple copies of the handlers running.
*/

/*
We should only replicate out requests like this to handlers that have different runtime conditions: different processes, machines, paths to a data strore, or access to different data stores
altogether.
*/

/*
Although this is can be expensive to set up and maintain, if speed is our goal, this is a valuable technique. In addition this naturally provides fault tolerance and scalability.
*/

func RequestsReplicationUsing() {
	const (
		timePaddingInSec = 1
	)

	var (
		doWork = func(done <-chan struct{}, id int, wg *sync.WaitGroup, result chan<- int) {
			defer func() {
				wg.Done()
			}()

			var (
				started = time.Now()
				// Simulate random load
				simulatedLoadTime = time.Duration(timePaddingInSec+rand.Intn(5)) * time.Second
			)

			// Produce the goroutine lock and unlock after simulatedLoadTime
			select {
			case <-done:
				return
			case <-time.After(simulatedLoadTime):
			}

			// Send the id to the caller
			select {
			case <-done:
				return
			case result <- id:
			}

			// The time that was spent
			took := time.Since(started)

			// Display how long handlers would have taken
			if took < simulatedLoadTime {
				took = simulatedLoadTime
			}

			fmt.Printf("%v took %v\n", id, took)
		}

		wg sync.WaitGroup
	)

	done := make(chan struct{})
	result := make(chan int)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			doWork(done, i, &wg, result)
		}(i)
	}

	firstReturned := <-result
	close(done)
	wg.Wait()

	fmt.Printf("Received an answer form %#v\n", firstReturned)
}
