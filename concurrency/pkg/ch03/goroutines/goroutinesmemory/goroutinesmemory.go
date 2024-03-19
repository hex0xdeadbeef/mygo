package goroutinesmemory

import (
	"fmt"
	"runtime"
	"sync"
)

/*
If we have too many concurrent processes, we can spend all of our CPU time context switching beetween them an never get any real work done. At OS level with threads, this can be
quite costly. The OS thread must save things like register values, lookup tables, memory maps to successfully be able to switch back to the current thread when it's time. Then it has
to load the same information for the incoming thread.

*/

func GoroutineMeasuring() {
	const (
		goroutinesNumber = 1e6
	)

	memConsumed := func() uint64 {
		runtime.GC()
		var s runtime.MemStats
		runtime.ReadMemStats(&s)
		return s.Sys
	}

	var (
		c  <-chan interface{}
		wg sync.WaitGroup

		noop = func() {
			wg.Done()
			<-c
		}
	)

	wg.Add(goroutinesNumber)
	before := memConsumed()
	for i := goroutinesNumber; i > 0; i-- {
		go noop()

	}

	wg.Wait()
	after := memConsumed()
	fmt.Printf("%.3f KB\n", float64(after-before)/1024)
}

