package fanoutin

import (
	"concurrency/pkg/ch04/pipeline"
	"fmt"
	"math"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

/*
One of the interesting properties of pipelines is the ability they give us to operate on the stream of data using a combination of separate, often reorderable stages. We can even reuse
stages of the pipeline multiple times.
*/

/*
Fan-out is a term to describe the process of starting multiple goroutines to handle input from the pipeline. Input from the pipeline -> many goroutines

Conditions so that a stage suits for utilizing fan-out:
	1) Order independence (The stage doesn't rely on values that the stage had calculated before). The property of order-independence is important because we have no guarantee in what
	order concurrent copies of our stage will run, nor in what order they will return.
	2) It takes a long time to run.

All we have to do is start multiple multiple versions of that stage.
*/

/*
Fan-in is a term to describe the process of combining multiple results into one channel. Many results -> one channel
*/

/*
A naive implementation of the FAN-IN, FAN-OUT algorithm only works if the order in which results arrive is unimportant.
*/

func FanInOutUsing() {

	const (
		topRandIntegerBoundary = 100_000_000_000
	)

	var (
		processorsNumber = runtime.GOMAXPROCS(0)

		randomNumber = func() interface{} {
			return rand.Intn(topRandIntegerBoundary)
		}

		toIntStream = func(done <-chan struct{}, dataProducer <-chan interface{}) <-chan int {
			ints := make(chan int)

			go func() {
				defer close(ints)

				for val := range dataProducer {
					select {
					case <-done:
						return
					case ints <- val.(int):
					}
				}
			}()

			return ints
		}

		getPrime = func(num int) int {
			dividersCounter := 0
			for i := int(math.Sqrt(float64(num))) + 1; i > 0; i-- {
				if num%i == 0 {
					dividersCounter++
				}

				if dividersCounter > 1 {
					return -1
				}
			}

			return num
		}

		primeFinder = func(done <-chan struct{}, intStream <-chan int) <-chan interface{} {
			primesStream := make(chan interface{})

			go func() {
				defer close(primesStream)

				for val := range intStream {
					num := getPrime(val)
					if num == -1 {
						continue
					}

					select {
					case <-done:
						return
					case primesStream <- num:
					}
				}

			}()

			return primesStream
		}

		/*
			FAN-IN EXAMPLE
			Many workers' resulting data -> the single channel
		*/
		fanIn = func(done <-chan struct{}, channels ...<-chan interface{}) <-chan interface{} {
			var (
				wg sync.WaitGroup

				multiplexedStream = make(chan interface{})
				multiplex         = func(c <-chan interface{}) {
					defer wg.Done()

					for val := range c {
						select {
						case <-done:
							return
						case multiplexedStream <- val:
						}
					}
				}
			)

			for _, c := range channels {
				wg.Add(1)
				go multiplex(c)
			}

			go func() {
				wg.Wait()
				close(multiplexedStream)
			}()

			return multiplexedStream
		}
	)

	// Canceling channel
	done := make(chan struct{})
	defer close(done)

	// The start point of primes finding
	start := time.Now()

	// Random numbers generator and multiple goroutines
	randIntStream := toIntStream(done, pipeline.RepeatFunction(done, randomNumber))

	/*
		FAN-OUT EXAMPLE.
		Input -> many workers
	*/
	primeFinders := make([]<-chan interface{}, processorsNumber)
	for i := 0; i < processorsNumber; i++ {
		primeFinders[i] = primeFinder(done, randIntStream)
	}

	// FAN-IN. Get the multiplexed channels into the single one.
	multiplexedFindersStream := fanIn(done, primeFinders...)

	fmt.Println("Primes and non-primes:")
	for val := range pipeline.LimitedJunction(done, multiplexedFindersStream, 10) {
		fmt.Println(val)
	}

	fmt.Printf("Time elapsed: %v\n", time.Since(start))

}
