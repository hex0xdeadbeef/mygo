package main

import (
	"sync"
	"time"
)

/*
	WAITGROUP
1. The structure of sync.WaitGroup is:
	type WaitGroup struct {
		noCopy noCopy

		state atomic.Uint64
		sema uint32
	}

	type noCopy struct{}

	func (*noCopy) Lock() {}
	func (*noCopy) Unlock() {}

2. In Go, it's easy to copy a struct by just assigning it to another variable. But some structs like WaitGroup, really shouldn't be copied.

	Copying a WaitGroup can mess things up because the internal state that tracks goroutines and their synchronization can get out of sync between the copies.

	THE FIELD noCopy
1. The noCopy struct is included in WaitGroup as a way to help prevent copying mistakes, not by throwing errors, but by serving as a warning.

2. The noCopy struct doesn't actually affect how our program runs. Instead, it acts as a marker tools like `go vet` can pick up on to detect when a struct has been copied in a way it
	shouldn't be.

3. The structure of noCopy:
	type noCopy struct{}

	func (*noCopy) Lock() {}
	func (*noCopy) Unlock() {}

	The methods are there just to work with the `-copylocks` in the `go vet` tool.

	When we run `go vet` on our code, it checks to see if any structs with noCopy field like WaitGroup, have been copied in a way that could cause issues.

	It'll throw an error to let us know there might be a problem. This gives us a heads-up to fix it before it turns into a bug.

	INTERNAL STATE
1. The state of WaitGroup is stored in an `atomic.Uint64` variable.
2. The `state` consists of two parts (the size of each part is 32 bits):
	1) Counter (high 32 bits):
		This part keeps track of the number of goroutines the WaitGroup is waiting for. When we call wg.Add(...) with a positive value, it bumps up this counter, and when we call wg.Done
		(), it decreases the counter by one.
	2) Waiter (low 32 bits):
		This tracks the number of goroutines currently waiting for that counter (the high 32 bits) to hit zero. Every time we call wg.Wait(), it increases this "waiter" count. Once the
		`Counter` part reaches zero, it releases all the goroutines that were waiting.

3. Then, there's the final field `sema uint32`, which is an internal semaphore managed by Go runtime. When a goroutine calls wg.Wait() and the counter part isn't zero, it increases the
	`Waiter` part count and then blocks by runtime_Semacquire(&wg.sema). This function call puts the goroutine to sleep until it gets woken up by a corresponding runtime_Semrelease(&wg.
	sema) call.

	ALIGNMENT PROBLEM.
1. The core of WaitGroup (the counter, the waiter, and semaphore) hasn't really changed across different Go versions. However, the way these elements are strucured has been modified many
	times.

	When we talk about alignment, we're referring to the need for data types to be stored at specific memory addresses to allow for efficient access.

	For example, on a 64-bit system, a 64-bit value like uint64 should be ideally be stored at a memory address that's a multiple of 8 bytes. The reason is, the CPU can grab aligned data
	in one step, but if the data isn't aligned, it might take multiple ops to access it.

2. On 32-bit architectures, the compiler doesn't guarantee that 64-bit values will be aligned on an 8-byte boundary. Instead, they might only be aligned on a 4-byte boundary.

	This becomes a problem when we use the atomic package to perform ops on the state variable. The atomic package specifically notes:
		On ARM, 386, and 32-bit MIPS, it is the caller’s responsibility to arrange for 64-bit alignment of 64-bit words accessed atomically via the primitive atomic functions
	What this means is if we don't align the `state` uint64 variable to an 8-byte boundary on these 32-bit architectures, it could cause the program to crash.

	THE SOLUTION OF GO'S 1.20 VERSION FOR WAITGROUP
1. In Go 1.19 a key optimization was introduced to ensure that uint64 values used in atomic ops are always aligned to 8-byte boundaries, even on 32-bit architectures.

	To achieve this, Russ Cox introduced a special struct atomic.Uint64, it's basically wrapper around uint64:
		type Uint64 struct {
			_ noCopy
			_ align64
			v uint64
		}

		// `align64` may be added to structs that must be 64-bit aligned.
		// This struct is recognized by a special case in the compiler
		// and won't be work if copied to any other package.
		type align64 struct{}

	So, what's the deal with align64?
		When we include align64 in a struct, it signals to the Go compiler that the entire struct needs to be aligned to an 8-byte boundary in memory.

	How exactly does align64 pull this off?
		align64 is just a plain, empty struct. It doesn't have any methods or special behaviors like the noCopy struct we talked about earlier.

	There's the magic happens.
		The align64 field itself doesn't take up any space, but it acts as a "marker" that tells Go compiler to handle the struct differently. The Go compiler recognizes align64 and
		automatically adjusts the struct's memory layout with the necessary padding, making sure that the entire struct starts at an 8-byte boundary.

	Thanks to atomic.Uint64, the state is guaranteed to be 8-byte aligned, so we don't have to worry about the alignment issues that could mess up atomic ops on 64-bit variables.

	HOW sync.WaitGroup INTERNALLY WORKS?
1. Why don't we just split the `Counter` and `Waiter` into two separate `uint32` variables?

	If we used a mutex to manage concurrency, we wouldn't have to worry about whether we're on a 32-bit or 64-bit system, no alignment issues to deal with.

	However, there's a trade-off. Using locks, like a mutex, adds overhead, especially when operations are happening frequently. Every time we lock and unlock, we're adding a bit of
	delay, which can stack up when dealing with high-frequency ops.

	On the other hand, we use atomic ops to modify this 64-bit variable safely, without the need to lock and unlock a mutex.
2. wg.Add(delta int)
	When we pass a value to the wg.Add(?) method, it adjusts the counter accordingly.

	If we provide a positive delta, it adds to the counter. Interestingly, we can also pass a negative delta, which will subtract from the counter.

	wg.Done() is just a shortcut for wg.Add(-1):
		func (wg *WaitGroup) Done() {
			wg.Add(-1)
		}

	However, if the negative delta causes the counter to drop below zero, the program will panic. The WaitGroup doesn't check the counter before updating it, it updates first and then
	checks. This means that if the counter goes negative after calling wg.Add(?), it stays negative until panic occurs.
3. The internals of wg.Add():
	func (wg *WaitGroup) Add(delta int) {
		// race stuffs

		// Add the delta to the counter
		state := wg.state.Add(uint(delta) << 32

		// Extract the counter and waiter count from the state
		v := int32(state >> 32)
		w := uint32(state)

		if v < 0 {
			panic("sync: negative WaitGroup counter")
		}

		if w != 0 && delta > 0 && v == int32(delta) {
			panic("sync: WaitGroup misuse: Add caller concurrently with Wait")
		)

		if v > 0 || w == 0 {
			return
		}

		if wg.state.Load() != state {
			panic("sync: WaitGroup misuse: Add called concurrently with Wait")
		}

		// Reset Waiters count to 0.
		wg.state.Store(0)
		for ; w != 0; w-- {
			runtime_Semrelease(&wg.sema, false, 0)
		}
	}

	Here's something that might not be widely known: when we add a positive delta to indicate the start of new tasks or goroutines, `this must happen before calling wg.Wait()`

	So, if we're reusing a WaitGroup, or if we want to wait for one batch of tasks, reset, and then wait for another batch, we need to make sure all calls to wg.Wait() are completed
	before starting new wg.Add(positive) calls for the next around of tasks.

	This avoids confusion about which tasks the WaitGroup is currently managing.

	On the other hand, we can call wg.Add(?) with a negative delta at any time, as long as it doesn't push the counter into negative territory.

4. wg.Wait()
	There isn't a whole lot to say about wg.Wait(). It basically loops and tries to increase the waiter count using CAS (Compare And Swap) atomic op.

	// Wait blocks until the [WaitGroup] counter is zero.
	func (wg *WaitGroup) Wait() {
		// race stuffs

		for {
			// Load the state of the WaitGroup
			state := wg.state.Load()
			v := int32(state >> 32)
			w := uint32(state)

			if v == 0 {
				// Counter is 0, no need to wait
				// ...
				return
			}

			// Increment waiters count
			if wg.state.CompareAndSwap(state, state + 1) {
				// ...

				runtime_Semacquire(&wg.sema)
				if wg.state.Load() != 0 {
					panic("sync: WaitGroup is reused before previous Wait has returned")
				}

				// ...

				return
			}
		}


	}

	If the CAS op fails, it means that another goroutine has modified the state, maybe the counter reached zero or it was incremented / decremented. In that case, wg.Wait() can't just
	assume everything is as it was, so it retries.

	When the CAS succeeds, wg.Wait() increments the waiter count and then puts goroutine to sleep using the semaphore.

	If we recall, at the end of the wg.Add(), it checks two conditions:
		1) if the counter is 0
		2) waiter count is greater than 0
	If both are true, it wakes up all the waiting goroutines and resets the state to 0.
*/

// This function represents the case, when `Waiter` part of the `state` field is used.
func multipleWaiters() {
	var ch1, ch2 = make(chan struct{}), make(chan struct{})
	var wg sync.WaitGroup

	go func() {
		wg.Wait()
		close(ch1)
	}()

	go func() {
		wg.Wait()
		close(ch2)
	}()

	for i := 0; i < 10_000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			time.Sleep(1 * time.Second)
		}()
	}

	<-ch1
	<-ch2
}
