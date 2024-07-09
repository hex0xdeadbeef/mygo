package main

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

/*
	COMPLEX CONCURRENCY PATTERNS

SNEAKY RACE CONDITIONS AND GRANULAR LOCKS
1. Race Detector.
	We should always use flag "-race" to build/run an app to check the presence of data races.

TIME CHANNELS OPERATIONS
1. Timing of channels can be done with time pkg (ticker, timer or else tools) or context.Context both are fine.
	1) Context is more idionatic but makes more allocations
	2) time is a little bit efficient

	FIRST COME, FIRST SERVED
1. There are some methods to spread a message to some chans.
2. If the program has many moving parts of unknown size, it might be worth revising the design of it, as it's very likely possible to simplify it.
*/

func main() {
	si := SharedInt{}
	wg := sync.WaitGroup{}

	for i := 0; i < 1e4; i++ {
		wg.Add(2)
		go func() {
			wg.Done()
			Answer(&si)
		}()

		go func() {
			wg.Done()
			si.Increment()
		}()

	}

	wg.Wait()

}

// THE NAKED MUTEX
type SharedInt struct {
	mu  sync.Mutex
	val int
}

// THE ENCAPSULATION
func (si *SharedInt) SetVal(newVal int) {
	si.mu.Lock()
	defer si.mu.Unlock()

	si.val = newVal
}

func (si *SharedInt) Val() int {
	si.mu.Lock()
	defer si.mu.Unlock()

	return si.val
}

func (si *SharedInt) Increment() {
	si.mu.Lock()
	defer si.mu.Unlock()

	si.val++
}

func Answer(si *SharedInt) {
	si.mu.Lock()
	defer si.mu.Unlock()

	si.val = 42
}

// THE MISUSE
/*
	There will be the data race because some calls of si.Increment() can happen
	between si.Val() and si.SetVal(...), so we loose any updates happen between these calls.
	The stale value will be set to the si.val

	So we put into the si.SetVal(...) a stale value

	This call is wrong logically, but doesn't include any data races
*/
func Increment(si *SharedInt) {
	si.SetVal(si.Val() + 1)
}

type Type interface{}

// TIMED CHANNELS MAKING
func ToTimedChannelContext(ctx context.Context, d time.Duration, message Type, c chan<- Type) (written bool) {
	ctx, cancel := context.WithTimeout(ctx, d)
	defer cancel()

	select {
	case c <- message:
		return true
	case <-ctx.Done():
		return false
	}
}

func ToTimedChannelTime(d time.Duration, message Type, c chan<- Type) (written bool) {
	var (
		timer = time.NewTicker(d)
	)
	defer timer.Stop()

	select {
	case c <- message:
		return true
	case <-timer.C:
		return false
	}
}

// SPREADING A MESSAGE TO TWO CHANS
/*
	After nilling a channel the select statement will never choose
	a case with nilled recerence of a or b
*/
func SpreadMessageToTwoChans(msg Type, a, b chan<- Type) {
	for i := 0; i < 2; i++ {
		select {
		case a <- msg:
			a = nil
		case b <- msg:
			b = nil
		}
	}
}

func SpreadMessageToMultipleChans(msg Type, channels ...chan<- Type) {
	var (
		wg sync.WaitGroup
	)

	for _, c := range channels {
		wg.Add(1)

		c := c

		go func() {
			wg.Done()

			c <- msg
		}()
	}

	wg.Wait()
}

// UNIFIED APPROACHES Timed chans + Spreading a message on multiple chans with contex/time pkgs
func SpreadMessageToTwoChansTimed(d time.Duration, msg Type, a, b chan<- Type) (written int) {
	var (
		timer = time.NewTimer(d)
	)

	defer timer.Stop()

	for i := 0; i < 2; i++ {
		select {
		case a <- msg:
			a = nil

		case b <- msg:
			b = nil

		case <-timer.C:
			return -1
		}

	}

	return 2
}

func SpreadMessageToMultipleChansTimed(ctx context.Context, d time.Duration, msg Type, channels ...chan<- Type) (written int) {
	var (
		wg sync.WaitGroup
		wr atomic.Int64
	)

	ctx, cancel := context.WithTimeout(ctx, d)
	defer cancel()

	for _, c := range channels {
		wg.Add(1)

		go func(c chan<- Type) {
			wg.Done()

			select {
			case c <- msg:
				wr.Add(1)
			case <-ctx.Done():
			}

			c <- msg
		}(c)
	}

	wg.Wait()

	return int(wr.Load())
}
