package basics

import (
	"bytes"
	"fmt"
	"os"
	"sync"
	"time"
)

/*
Any goroutine that attempts to write to a channel that is full will wait until the channel has been emptied, and any goroutine that attempts to read from a channel
that is empty will wait until at least one item is placed on it.
*/

/*
The receiving bool value indicates whether the gotten value is genereted by sending from a channel or by the closed channel.
*/

/*
If we have N goroutines that are waiting on a single channel, instead of writing N times to the channel to unblock each goroutine, you can simply close the channel. Since
a closed channel can be read from an infinite number of times, it doesn't matter how many goroutines are waiting on it, and closing the channel is both cheaper and faster
than performing N writes.
*/

/*
Buffered channels are an in-memory FIFO queue for concurrent processes to communicate over.

If a buffered channel is empty and has a receiver, the buffer will be bypassed and the value will be passed directly from the sender to the receiver.
*/

/*
Buffered channels can be a sort of premature optimization and also hide deadlocks by making them more unlikely to happen.
*/

/*
If a goroutine making writes to a channel has knowledge of how many writes it will make, it can be useful to create a buffered channel whose capacity is the number of writes
to be made, and then make those writes as quickly as possible.
*/

/*
The default value of a channel is nil.
1) Reading from the nil channel results in a panic.
2) Writing to a nil channel cause a panic
3) Closing a nil channel results in the following error: close of nil channel.
*/

/*
Things to put channel in the right context:
1) To assign channel ownership. Channel owners have a write-access view into the channel (chan or chan <-). Channel utilizers only have a read-only view into the channel
(<- chan). The goroutine that owns a channel should:
	1) Instantiate a channel
	2) Perform writes, or pass ownership to another goroutine.
	3) Close the channel.
	4) Ecapsulate the 1), 2), 3) in this list and expose them via a reader channel.
*/

func UnbufferedUsing() {
	strStream := make(chan string)

	go func() {
		strStream <- "Hello!"
	}()

	fmt.Println(<-strStream)
}

func Deadlock() {
	strStream := make(chan string)

	go func() {
		if 1 != 0 {
			return
		}
		strStream <- "Hello!"
	}()

	<-strStream
}

/*
	fatal error: all goroutines are asleep - deadlock!

	goroutine 1 [chan receive]:
	concurrency/pkg/ch02/channels/basics.Deadlock()
		/Users/dmitriymamykin/Desktop/goprojects/concurrency/pkg/ch02/channels/basics/basics.go:30 +0x44
	main.main()
		/Users/dmitriymamykin/Desktop/goprojects/concurrency/bin/main.go:7 +0x1c
*/

func IsChannelOpenOk() {
	stringStream := make(chan string)

	go func() {
		stringStream <- "Hello!"
	}()

	if val, ok := <-stringStream; ok {
		fmt.Printf("(%t): %s", ok, val)
	}

}

func IsChannelClosedOk() {
	strStream := make(chan string)

	go func() {
		strStream <- "Hello!"
		close(strStream)
	}()

	time.Sleep(500 * time.Millisecond)

	for i := 0; i < 2; i++ {
		if val, ok := <-strStream; ok {
			fmt.Printf("(%t): %q\n", ok, val)
		} else {
			fmt.Printf("(%t): %q\n", ok, val)
		}
	}

}

func ImmediateClosing() {
	intStream := make(chan int)
	close(intStream)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		fmt.Println(<-intStream)
	}()

	wg.Wait()
}

func RangingOverChannel() {
	intStream := make(chan int)
	go func() {
		defer close(intStream)
		for i := 0; i < 5; i++ {
			intStream <- i
		}
	}()

	// it'll be closed after the channel is drained
	for val := range intStream {
		fmt.Println(val)
	}

}

func UnlockingNGoroutines() {
	begin := make(chan struct{})
	wg := sync.WaitGroup{}

	for i := 0; i < 33; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			<-begin
			fmt.Printf("%d goroutine has been woken up\n", i)

		}(i)
	}

	close(begin)
	wg.Wait()
}

func BufferedUsing() {
	var stdoutBuff bytes.Buffer
	defer stdoutBuff.WriteTo(os.Stdout)

	intStream := make(chan int, 4)
	go func() {
		defer close(intStream)
		defer fmt.Fprintf(&stdoutBuff, "Producer done.\n")

		for i := 0; i < 5; i++ {
			fmt.Fprintf(&stdoutBuff, "Sending: %d\n", i)
			intStream <- i
		}
	}()

	for integer := range intStream {
		fmt.Fprintf(&stdoutBuff, "Received %d.\n", integer)
	}
}

// The right lifcycle of the channel.
func SensibleChannelUsing() {
	// Create a channel and return a reader-channel
	chanOwner := func() <-chan int {
		// Create the channel
		resultStream := make(chan int, 5)

		go func() {
			// Close the channel after all the writes
			defer close(resultStream)

			// Write the values into the channel
			for i := 0; i <= 5; i++ {
				resultStream <- i
			}

		}()

		fmt.Println("Receiver channel sent.")
		// Return unidirectional reader-channel
		return resultStream
	}

	resultStream := chanOwner()
	for val := range resultStream {
		fmt.Printf("Received: %d\n", val)
	}

	fmt.Println("Receiver ended its work. All the values accepted.")
}
