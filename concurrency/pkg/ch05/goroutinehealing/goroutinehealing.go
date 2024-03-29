package goroutinehealing

import (
	"concurrency/pkg/ch04/bridgechannel"
	"concurrency/pkg/ch04/pipeline"
	pgl "concurrency/pkg/ch04/preventinggorleak"
	"fmt"
	"log"
	"os"
	"time"
)

/*
To heal goroutines, we'll use our heartbeat pattern to check up on the liveness of the goroutine we're monitoring.

The type of heartbeat will be determined by what we're trying to monitor, but if our goroutine can become livelocked, make sure that the heartbeat contains some kind of information
indicating that the goroutine is not only up, but doing useful work.
*/

/*
1. The logic that monitors a goroutine's health is a "steward". Stewards will also be responsible for restarting a ward's goroutine should it become unhealthy. To do so, it will need a
reference to a function that can start the goroutine.
2. The goroutine is monitored is called "ward".
*/

// The goroutine that can be monitored and restarted.
type startGorFn func(done <-chan struct{}, pulseInterval time.Duration) (heartbeat <-chan struct{})

/*
timeout is the timeout of the function will being motitored
startGor is used to restart the unhealthy goroutine
*/
func NewSteward(timeout time.Duration, startGor startGorFn) startGorFn {

	return func(done <-chan struct{}, pulseInterval time.Duration) <-chan struct{} {
		// Steward's timeout to give an upward goroutine the info that it's still alive
		heartbeat := make(chan struct{})

		go func() {
			defer close(heartbeat)

			// The variables to control the liveness of the ward.
			var (
				wardDone      chan struct{}
				wardHeartbeat <-chan struct{}
				startWard     = func() {
					wardDone = make(chan struct{})
					wardHeartbeat = startGor(pgl.Or(done, wardDone), timeout/2)
				}

				pulse = time.NewTicker(pulseInterval)
			)
			defer pulse.Stop()

			startWard()

		monitorLoop:
			for {

				// Update the timeoutStream
				timeoutStream := time.After(timeout)
				for {

					select {

					// Signal to the upward goroutine that steward is alive
					case <-pulse.C:
						select {
						case heartbeat <- struct{}{}:
						default:
						}

						// Check whether the ward is alive
					case <-wardHeartbeat:
						continue monitorLoop

						// Check whether the timeout is reached
					case <-timeoutStream:
						log.Println("steward: ward unhealthy; restarting")
						close(wardDone)
						startWard()
						continue monitorLoop

						// Check whether the all work has done or not.
					case <-done:
						return

					}

				}

			}
		}()

		return heartbeat
	}

}

func UsingFirst() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime | log.LUTC)

	doWork := func(done <-chan struct{}, _ time.Duration) <-chan struct{} {
		log.Println("ward: Hi, I'm irresponsible!")

		go func() {
			<-done
			log.Println("ward: i'm halting")
		}()

		return nil
	}

	doWorkWithSteward := NewSteward(4*time.Second, doWork)

	done := make(chan struct{})
	time.AfterFunc(9*time.Second, func() {
		log.Println("main: halting steward and ward.")
		close(done)
	})

	for range doWorkWithSteward(done, 4*time.Second) {
	}

	log.Println("Done.")
}

func UsingSecond() {
	doWorkFn := func(done <-chan struct{}, intList ...int) (startGorFn, <-chan interface{}) {
		var (
			intChanStream = make(chan (<-chan interface{}))
			intStream     = bridgechannel.Bridge(done, intChanStream)

			doWork = func(done <-chan struct{}, pulseInterval time.Duration) <-chan struct{} {
				var (
					intStream = make(chan interface{})
					heartbeat = make(chan struct{})
				)

				go func() {
					defer close(intStream)

					select {
					case intChanStream <- intStream:
					case <-done:
						return
					}

					pulse := time.NewTicker(pulseInterval)
					defer pulse.Stop()

					for {

					valueLoop:
						for _, intVal := range intList {

							if intVal < 0 {
								log.Printf("negative value: %v\n", intVal)
								return
							}

							for {
								select {
								case <-pulse.C:
									select {
									case heartbeat <- struct{}{}:
									default:
									}
								case intStream <- intVal:
									continue valueLoop
								case <-done:
									return
								}
							}

						}
					}
				}()

				return heartbeat
			}
		)
		return doWork, intStream
	}

	log.SetFlags(log.Ltime | log.LUTC)
	log.SetOutput(os.Stdout)

	done := make(chan struct{})
	defer close(done)

	doWork, intStream := doWorkFn(done, 1, 2, -1, 3, 4, 5)
	doWorkWithSteward := NewSteward(1*time.Millisecond, doWork)
	doWorkWithSteward(done, 1*time.Hour)

	for intVal := range pipeline.LimitedJunction(done, intStream, 6) {
		fmt.Printf("Received: %v\n", intVal)
	}
}
