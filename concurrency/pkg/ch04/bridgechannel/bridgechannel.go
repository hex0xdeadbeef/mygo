package bridgechannel

import (
	"concurrency/pkg/ch04/readerdonewrapper"
	"fmt"
)

/*
The technique that defines destructuring a chanel of channels into a simple channel is called "bridging"

chanStream <- <- chan interface{} is a channel that is being given some channels during runtime.
*/

/*
Thanks to "Bridge()", we can use the channel of channels from within a single range statement and focus on out loop's logic. Destructing the channel of channels is left to code that
is specific to this concern.
*/

func Bridge(done <-chan struct{}, chanStream <-chan <-chan interface{}) <-chan interface{} {
	// This channel will be given values from the subchannels that chanStream has.
	dataProducer := make(chan interface{})

	go func() {
		defer close(dataProducer)

		for {
			var subProducer <-chan interface{}

			select {
			// Try to pull the channel
			case underlyingStream, ok := <-chanStream:
				if !ok {
					return
				}

				subProducer = underlyingStream

			case <-done:
				return

			}

			for val := range readerdonewrapper.IsReaderDone(done, subProducer) {
				select {
				// Send the value to the
				case dataProducer <- val:
				case <-done:
				}

			}

		}
	}()

	return dataProducer
}

func Using() {
	// The parentheses is required.
	genVals := func() <-chan (<-chan interface{}) {

		// The parentheses is required.
		chanStream := make(chan (<-chan interface{}))

		go func() {
			defer close(chanStream)

			for i := 0; i < 10; i++ {
				// The channel is buffered so that this goroutine won't be blocked.
				stream := make(chan interface{}, 1)

				stream <- i
				close(stream)

				chanStream <- stream
			}
		}()

		return chanStream
	}

	// Channels will be invoked with LIFO principle
	for v := range Bridge(nil, genVals()) {
		fmt.Printf(" %v", v)
	}
}
