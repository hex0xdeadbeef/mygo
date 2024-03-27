package preventinggorleak

import (
	"fmt"
	"math/rand"
	"time"
)

/*
The goroutine has a few paths to termination:
	1) When it has completed its work
	2) When it cannot continue its work due to an unrecoverable error
	3) When it's told to stop working.
*/

/*
	If a goroutine is responsible for creating a goroutine, it's also responsible for ensuring it can stop the goroutine.
*/

func CancelationAbsence() {

	doWork := func(strings <-chan string) <-chan interface{} {
		completed := make(chan interface{})

		go func() {
			// It won't be reached because of the absence of closing the doWork channel
			defer fmt.Println("Child goroutine stopped.")
			// If we pass the nil channel this goroutine will be blocked.
			for s := range strings {
				fmt.Println(s)
			}
		}()

		return completed
	}

	doWork(nil)

	// some work here

	fmt.Println("Done.")

}

func SimpleReaderCancelation() {

	doWork := func(done <-chan struct{}, strings <-chan string) <-chan struct{} {
		terminated := make(chan struct{})
		go func() {
			defer func() {
				close(terminated)
				fmt.Println("doWork exited")
			}()

			for {
				select {
				case s := <-strings:
					fmt.Println(s)
				case <-done:
					return
				default:
					fmt.Println("Waiting for cancelation.")
					time.Sleep(750 * time.Millisecond)
				}
			}
		}()

		return terminated
	}

	done := make(chan struct{})
	terminated := doWork(done, nil)

	go func() {
		time.Sleep(1 * time.Second)
		fmt.Println("Canceling doWork goroutine...")
		close(done)
	}()

	<-terminated

	fmt.Println("Done.")
}

func SimpleWriterCancelation() {
	newRandStream := func(done <-chan struct{}) <-chan int {
		randStream := make(chan int)

		go func() {
			defer func() {
				close(randStream)
				fmt.Println("newRandStream closure exited.")
			}()

			for {
				select {
				case <-done:
					return
				case randStream <- rand.Int():
				}
			}
		}()

		return randStream
	}

	done := make(chan struct{})
	randStream := newRandStream(done)

	fmt.Println("3 random ints:")
	for i := 1; i <= 3; i++ {
		fmt.Printf("%d: %d\n", i, <-randStream)
	}
	close(done)

	time.Sleep(1 * time.Second)
}

// The case when we want union all the "done" channel into a single one.
func or(doneChannels ...<-chan struct{}) <-chan struct{} {
	switch len(doneChannels) {
	case 0:
		return nil
	case 1:
		return doneChannels[0]
	}

	// Build a signal tree so that when the first signal is caught, we will stop the execution.
	orDone := make(chan struct{})
	go func() {
		defer close(orDone)

		switch length := len(doneChannels); {
		case length == 2:
			select {
			case <-doneChannels[0]:
			case <-doneChannels[1]:
			}
		case length > 2:
			select {
			case <-doneChannels[0]:
			case <-doneChannels[1]:
			case <-doneChannels[2]:
			case <-or(append(doneChannels[3:], orDone)...):
			}
		}
	}()

	return orDone
}

func OrExample() {
	signal := func(after time.Duration) <-chan struct{} {
		c := make(chan struct{})

		go func() {
			defer close(c)
			time.Sleep(after)
		}()

		return c
	}

	start := time.Now()
	<-or(
		signal(2*time.Hour),
		signal(5*time.Minute),
		signal(1*time.Second),
		signal(1*time.Hour),
		signal(1*time.Minute),
	)

	fmt.Printf("done after: %v", time.Since(start))
}
