package teechannel

import (
	pl "concurrency/pkg/ch04/pipeline"
	"concurrency/pkg/ch04/readerdonewrapper"
	"fmt"
	"math/rand"
	"sync"
)

/*
tee channel is used to split values coming in from a channel so that you can send them off into two separate areas of our codebase.

Imagine a channel of user-commands: You might want to take in a stream of user commands on a channel, send the to something that executes them, and also send them to something that
logs the commands for later auditing.
*/

func TeeChannel(done <-chan struct{}, dataProducer <-chan interface{}) (_, _ <-chan interface{}) {
	out1, out2 := make(chan interface{}), make(chan interface{})

	go func() {
		defer func() {
			close(out1)
			close(out2)
		}()

		for val := range readerdonewrapper.IsReaderDone(done, dataProducer) {
			// The copies of the outer variables. These are stored in heap.
			out1, out2 := out1, out2

			// The two iterations let us be sure that every channel has taken the same value.
			for i := 0; i < 2; i++ {
				// Select will be blocked until every channel has been given the same value
				select {
				case out1 <- val:
					// We need a copy to be nil value because every case statement is needed to be chosen once.
					out1 = nil
				case out2 <- val:
					// We need a copy to be nil value because every case statement is needed to be chosen once.
					out2 = nil
				}
			}

			// The iteration over dataProducer cannot be finished until both out1 and out2 are written in.
		}
	}()

	return out1, out2
}

func Using() {
	done := make(chan struct{})
	defer close(done)

	wg := sync.WaitGroup{}

	out1, out2 := TeeChannel(done, pl.LimitedJunction(done, pl.RepeatFunction(done, func() any { return rand.Intn(100) }), 10))

	wg.Add(1)
	go func() {
		defer wg.Done()

		for v := range out1 {
			fmt.Printf("1: %v\n", v)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		for v := range out2 {
			fmt.Printf("2: %v\n", v)
		}
	}()

	wg.Wait()
}
