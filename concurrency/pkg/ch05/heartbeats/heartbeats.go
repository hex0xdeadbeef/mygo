package heartbeats

import (
	"fmt"
	"math/rand"
	"time"
)

/*
Heartbeats are a way for concurrent processes to signal life to outside parties. They allow us insights into our system, and they can make testing the system deterministic when it might
otherwise not be.
*/

/*
There are two different types of heartbeats:
	1) Heartbeats that occur on a time interval. These ones are useful for concurrent code that might be waiting for something else to happen for it
	to process a unit of work. Because you don't know when that work might come in, your goroutine might be sitting around for a while waiting for
	something to happen. A heartbeat is a way to signal to its listeners that everything is well, and the silence is expected.
	2) Heartbeats that occur at the beginning of a unit of work.
*/

/*
Heartbeats help with the opposite case: they let us know that long-running goroutines remain up, but are just taking a while to produce a value to send on the values channel.
*/

/*
The unit hearbeats shines when it's used in writing tests. Interval-based heartbeats can be used in the same fashion, but if you only care that the goroutine has started doing its work, this style
of heartbeat is simple.
*/

/*
Heartbeats aren't strictly necessary when writing concurrent code, but this section demonstrates their utility.
*/

func TimeoutHeartbeat() {
	var (
		/*
			doWork returns two channels:
				1) heartbeat that is responsible for pulsation
				2) results that is responsible for sending results of the work
		*/
		doWork = func(done <-chan struct{}, pulseInterval time.Duration) (<-chan struct{}, <-chan time.Time) {
			heartbeat := make(chan struct{})
			results := make(chan time.Time)

			go func() {

				pulse := time.NewTicker(pulseInterval)
				workGen := time.NewTicker(pulseInterval * 2)

				/*
					Close the channels and stop the tickers
				*/
				defer func() {
					close(results)
					close(heartbeat)

					pulse.Stop()
					workGen.Stop()
				}()

				/*
					sendPulse is responsible for making pulsation
				*/
				sendPulse := func() {
					select {
					case heartbeat <- struct{}{}:
					default:
					}
				}

				/*
					sendResult is responsible for:
					1) Sending the results of work
					2) Making the pulsation alive
					3) Checking whether the full work of the process is done or not
				*/
				sendResult := func(r time.Time) {
					for {
						select {
						case <-done:
							return
						case <-pulse.C:
							sendPulse()
						case results <- r:
							return
						}
					}
				}

				for {
					// for i := 0; i < 2; i++ {
					select {
					case <-done:
						return
					case <-pulse.C:
						sendPulse()
					case r := <-workGen.C:
						sendResult(r)
					}
				}
			}()

			return heartbeat, results
		}
	)

	const (
		timeout = 2 * time.Second
	)

	done := make(chan struct{})
	time.AfterFunc(10*time.Second, func() { close(done) })

	heartbeat, results := doWork(done, timeout/2)
	for {
		select {
		case _, ok := <-heartbeat:
			if !ok {
				return
			}
			fmt.Println("pulse")

		case r, ok := <-results:
			if !ok {
				return
			}
			fmt.Printf("results %v\n", r.Second())

		case <-time.After(timeout):
			fmt.Println("worker goroutine is not healthy")
			return
		}
	}
}

func UnitHeartbeat() {
	var (
		doWork = func(done <-chan struct{}) (<-chan struct{}, <-chan int) {
			/*
				heartbeatStream accepts a heartbeat and is buffered so that the goroutine won't block while sending
			*/
			heartbeatStream := make(chan struct{}, 1)
			workStream := make(chan int)

			go func() {
				defer func() {
					close(workStream)
					close(heartbeatStream)
				}()

				for i := 0; i < 10; i++ {
					select {
					case heartbeatStream <- struct{}{}:
					default:
					}

					select {
					case <-done:
						return
					case workStream <- rand.Intn(10):
					}
				}
			}()

			return heartbeatStream, workStream
		}
	)

	done := make(chan struct{})
	defer close(done)

	heartbeat, results := doWork(done)

	for {
		select {
		case _, ok := <-heartbeat:
			if !ok {
				return
			}

			fmt.Println("pulse")

		case r, ok := <-results:
			if !ok {
				return
			}

			fmt.Printf("results: %d\n", r)
		}
	}
}

func DoWork(done <-chan struct{}, pulseInterval time.Duration, nums ...int) (<-chan struct{}, <-chan int) {

	var (
		heartbeat = make(chan struct{}, 1)
		intStream = make(chan int)

		pulse = time.NewTicker(pulseInterval)
	)

	// ... some work
	time.Sleep(2 * time.Second)

	go func() {
		defer func() {
			close(intStream)
			close(heartbeat)

			pulse.Stop()
		}()

	numLoop:
		for _, n := range nums {
			for {
				select {
				case <-done:
					return

					// This is made in order to be sure, that all the iterations take less time than the acceptable period of time.
				case <-pulse.C:
					select {
					case heartbeat <- struct{}{}:
					default:
					}
				case intStream <- n:
					continue numLoop
				}
			}
		}
	}()

	return heartbeat, intStream

}
