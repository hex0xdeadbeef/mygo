package selectpkg

import (
	"fmt"
	"time"
)

/*
Unlike switch blocks, case statements in a select block aren't tested sequentially, and execution won't automatically fall through if non of the criteria are met.

Instead, all channel reads and writes considered simultaneously to see if any of them are ready: populated or closed channel in the case reads, and channels that are not at
capacity in the case of writes. If non of the channels are ready, the entire select statement blocks.

Then when one of the channels is ready, that operation will proceed, and its corresponding statements will execute.
*/

/*
When all channels are blocked an none of them are ready to give the event we can use: tme.After(time.Duration())
*/

func Example() {
	start := time.Now()

	c := make(chan interface{})

	go func() {
		time.Sleep(1500 * time.Millisecond)
		close(c)
	}()

	fmt.Println("Blocking on read...")
	select {
	case <-c:
		fmt.Printf("Unblocked: %v later.\n", time.Since(start))
	}
}

func RandomReadiness() {
	var (
		c1, c2             = make(chan interface{}), make(chan interface{})
		counter1, counter2 int
	)
	close(c1)
	close(c2)

	for i := 1000; i > 0; i-- {
		select {
		case <-c1:
			counter1++
		case <-c2:
			counter2++
		}
	}

	fmt.Printf("Counter1: %d\nCounter2: %d\n", counter1, counter2)

}

func BlockedChannelHandlingWithTimeout() {
	var c <-chan int

	select {
	case <-c:
	case <-time.After(time.Second * 1):
		fmt.Println("Timed out")
	}
}

func NoReadyChannelsHandling() {
	start := time.Now()
	var c1, c2 <-chan struct{}

	select {
	case <-c1:
	case <-c2:
	default:
		fmt.Printf("In default after: %v\n", time.Since(start))
	}
}

func WaitForAnotherGoroutine() {
	done := make(chan struct{})

	go func() {
		time.Sleep(5 * time.Second)
		close(done)
	}()

	workCounter := 0

loop:
	for {
		select {
		case <-done:
			break loop
		default:
		}
	}

	workCounter++
	time.Sleep(1 * time.Second)

	fmt.Printf("Achieved %v cycles of work before signalled to stop. \n", workCounter)
}

/*
fatal error: all goroutines are asleep - deadlock!

goroutine 1 [select (no cases)]:
concurrency/pkg/ch02/channels/selectpkg.InfiniteDeadlock()
	/Users/dmitriymamykin/Desktop/goprojects/concurrency/pkg/ch02/channels/selectpkg/selectpkg.go:107 +0x1c
main.main()
	/Users/dmitriymamykin/Desktop/goprojects/concurrency/bin/main.go:7 +0x1c
*/

func InfiniteDeadlock() {
	select {}
}
