package main

import (
	"fmt"
	"runtime"
	"sync"
)

/*
	MUTEXES
1. Mutex, or MUtual EXclusion, in Go is basically a way to make sure that only one goroutine is messing with a shared resource at a time. This resource can be a piece of code, an integer, a map, a
	struct, a channel, or pretty much anything.

	WHY DO WE NEED sync.Mutex?
1. A datarace happens when multiple goroutines try to access and change shared data at the same time without proper synchronization.

2. In fact, The operation `counter++` consists of three ops:
	MOVD	main.counter(SB), R0
	ADD	$1 R0, R0
	MOVD	R0, main.counter(SB)

	The `counter++` is a read-modify operation and these steps above aren't atomic (in the case of NOT using mutex), meaning they're not executed as a single, uninterruptible action.

	For instance, goroutine G1 reads the value of counter, and `before` it writes the updated value, goroutine G2 reads the same value. Both then write their updated values back, but since they
	read the same original value, one increment is practically lost.

3. If we call Unlock() on an already unlocked mutex, it'll cause a fatal error:
	sync: unlock of unlocked mutex

	MUTEX STRUCTURE: THE ANATOMY
1. The mutex itself is:
		type Mutex struct {
			state int32
			sema uint32
		}

	At its core, a mutex in Go has two fields:
		1) state int32
		2) sema uint32
	They might look simple numbers, but there's more to them than meets the eye.

2. The `state` field is a 32-bit integer that shows the current state of the mutex. It's actually divided into multiple bits that encode various pieces of information about the mutex.

3. State bits (modes):
	1) Locked (junior 0th bit):
		Whether the mutex is currently locked. If it's set to 1, the mutex is locked and no other goroutine can grab it.

	2) Woken (junior 1st bit):
		Set to 1 if any goroutines has been woken up and is trying to acquire the mutex. Other goroutines shouldn't be woke up unnecessary.

	3) Starvation (junior 2nd bit):
		This bit shows whether the mutex is in starvation mode (set to 1). We'll dive into what this mode means in a bit.

	4) Waiter (3-31)
		The rest of the bits keep track of how many goroutines are waiting for the mutex.

4. The other field `sema` is a `uint32` that acts as a semaphore to manage and signal waiting goroutines. When the mutex is locked, one of the goroutines is woken up to acquire the lock.

	Unlike the state `sema` doesn't have a specific bit layout and relies on runtime internal code to handle the semaphore logic.

	MUTEX LOCK FLOW
1. In the mutex.Lock() function there are two paths:
	1) The fast path for the usual case
	2) The slow path for handling the unusual case

	func (m *Mutex) Lock() {
		// Fast path: grab unlocked mutex.
		if atomic.CompareAndSwap(&m.state, 0, mutexLocked) {
			if race.Enabled {
				race.Acquire(unsafe.Pointer(m))
			}
			return
		}

		// Slow path (outlined so that the past path can be inlined)
		m.lockSlow()
	}

	The fast path is designed to be really quick and is expected to handle most lock acquisitions where the mutex isn't already in use. This path is also inlined, meaning it's embedded directly
	into the calling function.

	FYI, this inlined fast path is neat trick that utilizes Go's inline optimization, and it's used a lot in Go's source code.

	When the CAS (Compare And Swap) operation in the fast path fails, it means the state field wasn't 0, so the mutex is currently locked.

	The real concern here is the slow path m.lockSlow(), which does most of the heavy lifting. In the slow path, the goroutine keeps actively spinning to try to acquire the lock, it doesn't just
	go straight to the waiting queue.

	WHAT DO YOU MEAN SPINNING?
1. Spinning means the goroutine enters a tight loop, repeatedly checking the state of the mutex without giving up the CPU.

	In this case, it's not a simple for loop but low-level assembly instructions to perform:
		TEXT runtime procyield(SB),NOSPLIT,$0-0
			MOVWU cycles+0(FP), R0
		again:
			YIELD
			SUBW	$1, R0
			CBNZ	R0, again
			RET
	YIELD returns an access to other threads in a system to prevent the process from busy waiting in a cycle.

	Firstly, a thread that is waiting for a mutex tries to acquire the mutex by using YIELD. YEILD is used as a `lightweigh pause`, and if the waiting exceeds the limit set, the thread can move to
	the blocking methods.

	CBNZ R0, again - the command "Compare and Branch if Not Zero" checks the value of the R0 registry. If it's not zero, the flow returs to `again`. It keeps the cycle working until the value is 0.

	The conclusion is: `YIELD` in mutexes' realization helps to optimize the concurrent access to resources. It provides the balance between active waiting and full blocking. This is relevant in
	multi-threaded environments, where high concurrency takes place.

	The Algorithm of Mutex working:
		1) In our case the number of cycles is 30 (runtime.procyield(30)). It means that the assembly code runs a tight loop for 30 cycles, repeatedly yeilding the CPU and decrementing the spin
		counter.

		2) After spinning, it tries to acquire the lock again. If it fails, it has three more chances to spin before giving up. So, in total, it tries for up to 120 cycles.

		3) If it still can't get the lock, it increases the waiter count, puts itself in the waiting queue, goes to sleep, waits for a signal to wake up and try again (1)

	WHY DO WE NEED SPINNING?
1. The idea behind spinning is to wait for a short while in hopes that the mutex will free up soon, letting the goroutine grab the mutex without the overhead of a sleep-wake cycle.

	If our computer doesn't have multiple cores, spinning isn't enabled because it would just waste CPU time.

	WHAT IF ANOTHER GOROUTINE IS ALREADY WAITING FOR A MUTEX? IT DOESN'T SEEM FAIR IF THIS GOROUTINE TAKES THE LOCK FIRST.
1. That's why our mutex has two modes:
	1) Normal
	2) Starvation. Spinning isn't working in starvation mode.

	In normal mode, goroutines waiting for the mutex are organized in a first-in, first-out (FIFO) queue. When a goroutine wakes up to try grab the mutex, it doesn't get control immediately.
	Instead, it has to compete with any new goroutines that also want the mutex at that time.

	This competition is tilted in favor of new goroutines because they're already running on the CPU and can quickly try to grab the mutex, while the queued goroutine is still waking up.

	As a result, the goroutine that just woke up might frequently lose the race to the new contenders and get put back at the front of the queue.

	WHAT IF THAT GOROUTINE IS UNLUCKY AND ALWAYS WAKES UP WHEN A NEW GOROUTINE ARRIVES?
1. Starvation mode kicks in if a goroutine fails acquire the lock for more than 1 millisecond. It's designed to make sure that waiting goroutines eventually get a fair chance at the mutex.

	In this mode, when a goroutine releases the mutex, it directly passes control to the goroutine at the front of the queue. This means no competition, no race, from new goroutines. They don't
	even try to acquire it and just join the end of the waiting queue.

	MUTEX UNLOCK FLOW
1. The unlock flow is simpler than the lock flow. We still have two paths:
		1) Inlined path. The fast path drops the locked bit in the state of the mutex. If dropping this bit makes the state zero, it means no other flags are set (like waiting goroutines), and our
		mutex is now completely free.

		2) Slow path, which handles unusual cases
2. The slow unlock.
	That's where the slow path comes in and it needs to know if our mutex is in normal mode or starvation mode.

	func (m *Mutex) unlockSlow(new int32) {
	// 1. Attempting to unlock an already unlocked mutex will cause a fatal error
		if (new + mutexLocked)&mutexLocked == 0 {
			fatal("sync: unlock of unlcoked mutex")
		}

		if new & mutexStarving == 0 {
			old := new

			for {
			// 2. If there are no waiters, or if the mutex is already locked, or woken, or in starvation mode, return.
				if old >> mutexWaiterShift == 0 || old&(mutexLocked|mutexWoken||mutexStarving) != 0 {
					return
				}
			}

			// Grab the right to someone
			new = (old - 1 << mutexWaiterShift) | mutexWoken
			if atomic.CompareAndSwap(&m.state, old, new) {
				runtime.Semrelease(&m.sema, false, 1)
				return
			}
			old = m.state
		} else {
			// 3. If the mutex is in starvation mode, hand off the ownership to the first waiting goroutine in the queue.
			runtime_Semrelease(&m.sema, true, 1)
		}
	}

	In normal mode, if there are waiters and no other goroutine has been woken or acquired the lock, the mutex tries to decrement the waiter count and turn on the mutexWoken flag atomically. If
	successful, it releases the semaphore to wake up one of the waiting goroutines to acquire the mutex.

	In starvation mode, it atomically increments the semaphore and hands off mutex ownership directly to the first waiting goroutine in the queue. The second argument of runtime_Semrelease
	determines if the handoff is true.


*/

func interruptibleIncrementing() {
	var (
		counter int
		wg      sync.WaitGroup
	)

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			counter++
		}()
	}

	wg.Wait()
}

func mutexedIncrementing() {
	var (
		counter int
		wg      sync.WaitGroup
		mu      sync.Mutex
	)

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			mu.Lock()
			counter++
			mu.Unlock()
		}()
	}

	wg.Wait()
}

func singleProc() {
	var (
		counter int
		wg      sync.WaitGroup
	)
	runtime.GOMAXPROCS(1)

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter++
		}()
	}

	wg.Wait()

	fmt.Println(counter)
}
