package benchcontextswitch

import (
	"sync"
	"testing"
)

func Benchmark(b *testing.B) {
	var (
		wg    sync.WaitGroup
		begin = make(chan struct{})
		c     = make(chan struct{})

		token struct{}

		// Pass a tocken to another goroutine
		sender = func() {
			defer wg.Done()
			// When we close the channel the work will start
			<-begin
			for i := 0; i < b.N; i++ {
				c <- token
			}
		}

		// Get a token from the other goroutine
		receiver = func() {
			defer wg.Done()
			// When we close the channel the work will start
			<-begin
			for i := 0; i < b.N; i++ {
				<-c
			}
		}
	)

	wg.Add(2)
	go sender()
	go receiver()

	b.StartTimer()
	close(begin)

	wg.Wait()

}
