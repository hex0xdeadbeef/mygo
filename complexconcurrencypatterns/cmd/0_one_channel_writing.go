package main

import (
	"context"
	"reflect"
	"sync"
	"sync/atomic"
	"time"
)

/*
	TIMED CHANNELS OPS
1. Sometimes we want to time our channels ops: keep trying to do something, and if we can't do it in time just drop the ball.

2. To do so we can either use "context" or "time", both are fine. Context might be more idiomatic,
	time is a little bit more efficientm but they're almost identical.
3. Since performance isn't really relevant here (after all we're waiting) the only difference I
	found is that the solution using context performs more allocations (also because the one with
	the timer can be further optimized to recycle timers).

4. Beware that re-using timers is tricky, so keep in mind that it might not be worth the risk to
	hust save 10 allocs/ops.

	FIRST COME FIRST SERVED
1. Sometimes we want to write the same message to many channels, writing to whichever is available
	first, but never writing the same message twice on the same channel.
2. To do this there are two magic ways:
	- We can mask the channels with local variables, and disable select cases accordingly, or use
	goroutines and waits.
3. Please note that in this case performance might matter, and at the time of writing the solution
	that spawns goroutines takes 4 times slower than the one with select.
4. If the amount of of channels is not kwown at compile time, the first solution becomes trickier,
	but it still possible, while the second one stays basically unchanged.

	Note: it your program has many moving parts of unknown size, it might be worth revising your design, as it's very likely possible to simplify it.
5. If your code survives you review and still has unbound moving parts, here are the two solutions
	to support that.
6. Needless to say: the solution using reflection is almost two orders of magnitude slower than
	one with goroutines and undreadable, so please don't use it.

	PUT IT TOGETHER
1. In case you want to both try a several sends for a while and abort if it's taking too long here
	are two solutions:
		- One with time + select
		- One with context + go
	The first one is better if the amount of channels is known at compile time, while the other one
	should be used when it's not.
*/

/*
	TIMED CHANNELS OPS
*/

type Type struct{}

func ToChanContext(ctx context.Context, d time.Duration, message Type, c chan<- Type) (written bool) {
	ctx, cancel := context.WithTimeout(ctx, d)
	defer cancel()

	select {
	case c <- message:
		return true
	case <-ctx.Done():
		return false
	}
}

func ToChanTimer(d time.Duration, message Type, c chan<- Type) (written bool) {
	t := time.NewTimer(d)
	defer t.Stop()

	select {
	case c <- message:
		return true
	case <-t.C:
		return false
	}
}

/*
	FIRST COME FIRST SERVED
*/

func FirstComeFirsServedSelect(message Type, a, b chan<- Type) {
	for i := 0; i < 2; i++ {
		aCopied, bCopied := a, b

		select {
		case aCopied <- message:
			aCopied = nil
		case bCopied <- message:
			bCopied = nil
		}
	}
}

func FirstComeFirstGoroutines(message Type, a, b chan<- Type) {
	var (
		wg sync.WaitGroup
	)

	wg.Add(1)
	go func(msg Type, a chan<- Type) {
		defer wg.Done()

		a <- msg
	}(message, a)

	wg.Add(1)
	go func(msg Type, b chan<- Type) {
		defer wg.Done()

		b <- msg
	}(message, b)

	wg.Wait()
}

func FirstComeFirstServedGoroutinesVariadic(message Type, channels ...chan<- Type) {
	var (
		wg sync.WaitGroup
	)
	wg.Add(len(channels))

	for _, c := range channels {
		c := c

		go func(wg *sync.WaitGroup, msg Type, c chan<- Type) {
			defer wg.Done()

			c <- msg
		}(&wg, message, c)
	}

	wg.Wait()
}

func FirsComeFirstServedSelectVariadic(message Type, channels ...chan<- Type) {
	/*
		// A SelectCase describes a single case in a select operation.
		// The kind of case depends on Dir, the communication direction.
		//
		// If Dir is SelectDefault, the case represents a default case.
		// Chan and Send must be zero Values.
		//
		// If Dir is SelectSend, the case represents a send operation.
		// Normally Chan's underlying value must be a channel, and Send's underlying value must be
		// assignable to the channel's element type. As a special case, if Chan is a zero Value,
		// then the case is ignored, and the field Send will also be ignored and may be either zero
		// or non-zero.
		//
		// If Dir is [SelectRecv], the case represents a receive operation.
		// Normally Chan's underlying value must be a channel and Send must be a zero Value.
		// If Chan is a zero Value, then the case is ignored, but Send must still be a zero Value.
		// When a receive operation is selected, the received Value is returned by Select.
		type SelectCase struct {
			Dir  SelectDir // direction of case
			Chan Value     // channel to use (for send or receive)
			Send Value     // value to send (for send)
		}
	*/
	var (
		cases = make([]reflect.SelectCase, len(channels))
	)

	for i := 0; i < len(channels); i++ {
		cases[i] = reflect.SelectCase{
			Dir:  reflect.SelectSend,
			Chan: reflect.ValueOf(channels[i]),
			Send: reflect.ValueOf(message),
		}
	}

	/*
		Select executes a select operation described by the list of cases. Like the Go select statement, it blocks until at least one of the cases can proceed, makes a uniform pseudo-random
		choice, and then executes that case. It returns the index of the chosen case and, if that case was a receive operation, the value received and a boolean indicating whether the value
		corresponds to a send on the channel (as opposed to a zero value received because the channel is closed). Select supports a maximum of 65536 cases.
	*/
	for i := 0; i < len(channels); i++ {
		chosenChannelIdx, _, _ := reflect.Select(cases)
		cases[chosenChannelIdx].Chan = reflect.ValueOf(nil)
	}
}

/*
	PUT IT TOGETHER
*/

func ToChansTimedTimerSelect(d time.Duration, message Type, a, b chan Type) (written int) {
	t := time.NewTimer(d)
	for i := 0; i < 2; i++ {
		select {
		case a <- message:
			a = nil
		case b <- message:
			b = nil
		case <-t.C:
			return i
		}

	}

	t.Stop()
	return 2
}

func ToChansTimedContextGoroutines(ctx context.Context, d time.Duration, message Type, ch ...chan Type) (written int) {
	ctx, cancel := context.WithTimeout(ctx, d)
	defer cancel()

	var (
		writtenCnt atomic.Int32
		wg         sync.WaitGroup
	)

	for i := 0; i < len(ch); i++ {
		wg.Add(1)

		c := ch[i]

		go func(wg *sync.WaitGroup, msg Type, ch chan<- Type) {
			defer wg.Done()

			select {
			case c <- msg:
				_ = writtenCnt.Add(1)

			case <-ctx.Done():
				return
			}

		}(&wg, message, c)
	}

	wg.Wait()

	return int(writtenCnt.Load())

}
