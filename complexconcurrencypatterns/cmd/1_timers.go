package main

import (
	"fmt"
	"time"
)

/*
	TIMERS
As I started in the previous article timers are hard to use the right way
	so here are some handful tips.

	PREFACE
1. If you don't think that dealing with time and timers while also juggling
	goroutines is hard, here are some juicy bugs related to time.Time
2. If you're not satisfied, in the following code there are a deadlock and
	a race condition.

	func A() {
		tm := time.NewTimer(1)
		tm.Reset(100*time.Millisecond)

		<-tm.C
		if !tm.Stop() {
			<- tm.C
		}
	}
	THE time.Tick OBJECT
1. This is just something we should never use unless we plan to carry the returned chan around and
	keeping using it for the entire lifetime of the program.
2. As the documentation says: "the underlying Ticker cannot be recovered by the GC; it leaks"
	Please use this with caution and in case of doubt use Ticker isntead.

	THE time.After OBJECT
1. This is basically the same concept of Tick but instead of hiding a Ticker, hides Timer. This is
	slightly better because once the timer will fire, it'll be collected. Please note that timers
	use unbuffered channels.
2. As above, if we care about performance and we want to be able to cancel the call, we shouldn't
	use After.

	THE time.Timer OBJECT
1. I find this to be quiet an odd and unusual API for Go: NewTicker(duraation) returns a *Timer,
	which has an exported chan field called "C" that is the interesting value for the caller. This
	is odd on Ticker, but very odd in Timer.
2. Usually in Go exported fields mean that the user can both get and set them, while here setting
	"C" doesn't mean anything. It's quiet opposite: setting "C" and resetting the timer will still
	make the runtime send messages in the channel that was previously in "C". This is even
	worsened by the fact that the Timer returned by AfterFunc doesn't use C at all.
3. That said, Timer has a lot more interesting oddities. Here is an overview on the API.

	type Timer struct {
		C <- chan Time
	}

	func AfterFunc(d time.Duration, f func()) *Timer
	func NewTimer(d Duration) *Timer
	func (*Timer) Stop() bool
	func (*Timer) Reset(d Duration) bool
4. time.AfterFunc
	Officia documentation says:
		AfterFunc waits for the duration to elapse and then calls f in its own goroutine. It
		returns a Timer that can be used to cancel the call using its Stop method.

	While this is correct, be careful: when calling Stop, if false returned, it means that stopping
	failed and the function was already started. This doesn't say anything about that function
	having returned, for that we need to add some manual coordination.
5. In addition to that, the returned timer will not fire. It's only usable to call Stop.
	t := time.AfterFunc(1*time.Second, func { fmt.Println("Time has passed!") })
	<-t.C // will cause deadlock
6. Also, at the moment of writing, resetting the timer makes the runtime call f again after the
	passed duration, but it's not documented so it might change in the future.

	THE time.Timer OBJECT
1. Official documentation says: "NewTimer creates a new Timer that will send the current time on
	its channel after at least duration d"

	This means that there's no way to construct a valid Timer without starting it. If we need to
	construct one for future re-use, we either do it lazily or we have to create and stop it, which can be done with the code presented.

	We have to read from the channel. If the timer fired between the New and Stop calls the
	channels isn't drained there will be a value in C, so all future reads will be wrong.

	THE (*time.Timer) Stop METHOD
1. Stop prevents the Timer from firing. It returns true if the call stops the timer, false if the
	timer has already expired or been stopped.
2. The "or" in the sentence above is very important. All examples for Stop() in the doc show this
	snippet:
		if !t.Stop() {
			<-t.C
		}
	The problem is that "or" means that this pattern is valid 0 or 1 times. It's not valid if
	someone else already drained the channel, and it's not valid to do this multiple times without
	calling Reset() in between.

	This summarized means that Stop()+drain is safe to do if and only if no reads from the channel
	have been performed, including ones caused by other drains.

	Moreover, the pattern above is not thread safe, as the value returned by Stop() might already
	be stale when draining the channel is attempted, and this would cause a deadlock as two
	goroutines would try to drain C.

	THE TABLE
1. Here is a table with fields of runtimeTimer set by the time package:

	Constructor		when		period			f			arg
	NewTicker(d)	d			set to d		sendTime	C
	NewTimer(d)		d			not set			sendTime	C
	AfterFunc(d,f)	d			not set			goFunc

2. Runtime timers don't rely on goroutines and are bucketed together to fire in an efficient and
	precise manner.
*/

func TimerDeadlockA() {
	// NewTimer creates a new Timer that will send the current time on its
	// channel after at least duration d.
	//
	// Before Go 1.23, the GC didn't recover timers that had not yet expired or been stopped,
	// so code often immediately deferred t.Stop after calling NewTimer, to make the timer recoverable
	// when it was no longer needed.
	//
	// As of Go 1.23, the GC can recover unreferenced timers, even if they haven't
	// expired or been stopped. The Stop method is no longer necessary to help the garbage collector.
	// (Code may of course still want to call Stop to stop the timer for other reasons.)
	//
	// Before Go 1.23, the channel associated with a Timer was asynchronous (buffered, capacity 1),
	// which meant that stale time values could be received even after [Timer.Stop] or [Timer.Reset]
	// returned. As of Go 1.23, the channel is synchronous (unbuffered, capacity 0), eliminating the
	// possibility of those stale values.
	//
	// The GODEBUG setting asynctimerchan=1 restores both pre-Go 1.23 behaviors: when set,
	// unexpired timers won't be garbage collected, and channels will have buffered capacity.
	// This setting may be removed in Go 1.27 or later.

	// This is the timer with the existence of 1 nanosecond
	// The timer starts to count the time and waits the expiration
	tm := time.NewTimer(1)

	// Reset changes the timer to expire after duration d.
	//
	// It returns true if the timer had been active, false if the timer had
	// expired or been stopped.
	//
	// For a func-based timer created with [AfterFunc](d, f), Reset either
	// reschedules when f will run, in which case Reset returns true, or
	// schedules f to run again, in which case it returns false.
	//
	// When Reset returns false, Reset neither waits for the prior f to
	// complete before returning nor does it guarantee that the subsequent
	// goroutine running f does not run concurrently with the prior
	// one. If the caller needs to know whether the prior execution of
	// f is completed, it must coordinate with f explicitly.
	//
	// For a chan-based timer created with NewTimer, as of Go 1.23,
	// any receive from t.C after Reset has returned is guaranteed not
	// to receive a time value corresponding to the previous timer settings;
	// if the program has not received from t.C already and the timer is
	// running, Reset is guaranteed to return true.
	// Before Go 1.23, the only safe way to use Reset was to [Stop] and
	// explicitly drain the timer first.
	// See the [NewTimer] documentation for more details.

	// The Reset() method changes the duration to 100 milliseconds and resets the channel tm.C.
	// That is, all the written values are cleaned.
	tm.Reset(100 * time.Millisecond)

	// The program locks trying to read the values from the channel tm.C and waits until the timer
	// has exiped. It'll happen in 100 milliseconds and program will be unlocked getting the time
	// from the channel
	<-tm.C

	// Stop prevents the [Timer] from firing.
	// It returns true if the call stops the timer, false if the timer has already
	// expired or been stopped.
	//
	// For a func-based timer created with [AfterFunc](d, f),
	// if t.Stop returns false, then the timer has already expired
	// and the function f has been started in its own goroutine;
	// Stop does not wait for f to complete before returning.
	// If the caller needs to know whether f is completed,
	// it must coordinate with f explicitly.
	//
	// For a chan-based timer created with NewTimer(d), as of Go 1.23,
	// any receive from t.C after Stop has returned is guaranteed to block
	// rather than receive a stale time value from before the Stop;
	// if the program has not received from t.C already and the timer is
	// running, Stop is guaranteed to return true.
	//
	// Before Go 1.23, the only safe way to use Stop was insert an extra
	// <-t.C if Stop returned false to drain a potential stale value.
	// See the [NewTimer] documentation for more details.

	// Stop() trues to stop the timer if it's active so far.
	// The method returns:
	// - true if the timer has been successfully stopped until its expirations
	// - false if the timer has expired or hasn't been active at the moment of Stop()'s call.

	// If tm.Stop() returns false, it means the timer has expired and the time value has been
	// already written into the channel tm.C. After that the program tries to read from the tm.C
	// channel.

	// However, tm.C is the unbuffered channel and there are no values in it so far since the time
	// value corresponding to the previous timer has been read at the first iteration of <-tm.C.
	// That is, the program will blocked forever, and it leads to the deadlock.
	if !tm.Stop() {
		<-tm.C
	}
}

func TimerDeadlockB(t *time.Timer, ch chan int) {
	// - If the timer was active before the calling of t.Reset(), the Reset
	// () returns true.
	// - If the timer has expired, Reset() returns false and the timer will
	// be started from scratch
	t.Reset(1 * time.Second)

	defer func() {
		// - If t.Stop() returns false, it means that the timer has already expired and the time value has been sent on the channel t.C
		// In this case the <-t.C is used to read from the channel,
		// otherwise this value will hang out in this channel and it'll lead
		// to the resource leaking.
		if !t.Stop() {
			<-t.C
		}
	}()

	select {
	case ch <- 42:
		// Read the value from the timer.
		// - If the timer has expired, the value will be thrown from the
		// channel
		/*
			- If the ch is locked, the operation ch <- 42 will be blocked
			- Simultaneously, if the t.C hasn't got any values, the t.C will be also blocked.
		*/
	case <-t.C:
	}

}

func AfterFuncCancelationController() {
	var (
		// done is used to signal that the function f of the timer has returned.
		// When the function f ends its work, it'll close this channel and it'll signal to a parent
		// goroutine that the function has finished.
		done = make(chan struct{})

		f = func() {
			defer close(done)
			fmt.Println("Do work of f()...")

		}

		// - After a millisecond has passed, the timer will asynchronously call the function f
		t = time.AfterFunc(1*time.Second, f)
	)

	time.Sleep(2 * time.Second)

	// - t.Stop() returns true if the timer has been successfully stopped and the function f won't
	// be called
	// - t.Stop() returns trie if the timer has expired and the function f has started the
	// execution of f
	if !t.Stop() {
		select {
		case <-done:
			fmt.Println("f() has finished its execution")
		default:
			fmt.Println("f() has been prevented from an execution.")
		}

	}

}

func ReusableTimerCreation() {
	var (
		t = time.NewTimer(0)
	)
	// We have to read from the channel. If the timer fired between the New and Stop calls the
	// channels isn't drained there will be a value in C, so all future reads will be wrong.
	if !t.Stop() {
		<-t.C
	}
}

// This ensures the timer is ready to be re-used after toChanReturns.
func toChanTimed(t *time.Timer, ch chan int) {
	t.Reset(1 * time.Second)
	// No defer, as we don't know which case will be selected

	select {
	case ch <- 42:
	case <-t.C:
		// C is drained early, return
		return
	}

	// We still need to check the return value
	// of Stop, because t could have fired
	// between the send on ch and this line
	if !t.Stop() {
		<-t.C
	}
}
