package main

import (
	"fmt"
	"math/rand/v2"
	"sync"
	"time"
)

/*
	sync.Cond INTERNALS
1. When a goroutine needs to wait for something to happen, like some shared data changing, it can `block`, meaning it just pauses its work until it gets the go-ahead to continue. The
	most basic way to do this is with a loop, maybe even adding a time.Sleep to prevent the CPU from going crazy with `busy-waiting`.

	// wait until condition is true
	for !condition {
	}

	// or
	for !condition {
		time.Sleep(100*time.Millisecond)
	}

	Not this isn't really efficient as that loop is still running in the background, burning through CPU cycles, even when nothing's changed.

2. Technically, sync.Cond is a `condition variable` if we're running from a more academic background.
	- When a goroutine is waiting for something to happen (waiting for a certain condition to become true), it can call Wait()
	- Another goroutine, once it knows that the condition might be met, can call Signal() or Broadcast() to wake up the waiting goroutine(s) and let them know it's time to move on.

3. The interface that sync.Cond provides:
	// Wait suspends the calling goroutine until the condition is met
	func (c *Cond) Wait() {...}

	// Signal wakes up a single goroutine, if there's one
	func (c *Cond) Signal() {...}

	// Broadcast wakes up all waiting goroutines
	func (c *Cond) Broadcast() {...}

4. In the function CondUsage:
	- One goroutine is waiting for `Pikachu` to show up, while another one (the producer) randomly selects a Pokemon from the list and signals the consumer when a new one appears.
	- When a producer sends the signal, the consumer wakes up and checks if the right Pokemon has appeared:
		1) If it has, we catch the Pokemon
		2) If not, the consumer goes back to sleep and waits for the next one.

	The problem is, there's a gap between the producer sending the signal and the consumer actually waking up. In the meantime, the Pokemon could change, because the consumer goroutine
	might wake up later than 1 ms (rarely) or other goroutine modifies the shared pokemon. So sync.Cond is basically saying: `Hey! something changed! Wake up and check it out, but if
	you're too late, it might change again.`

5. If the consumer wakes up late, the Pokemon might run away, and the goroutine will go back to sleep.

6. Huh, I could use a channel to send the pokemon name or a signal to the other goroutine.
	Absolutely. In fact, channels are generally preferred over `sync.Cond` in Go because they're simpler, more idiomatic, and familiar to most developers.

	In our case we could easily send the Pokemon name through a channel, or just use an empty struct{} value to signal without sending any data. But our issue isn't just about passing messages
	through channels, it's about dealing with a shared state.

7. Our example is pretty simple, but if multiple goroutines are accessing the shared pokemon variable, let's look at what happens if we use a channel:
	- If we use a channel to send the Pokemon name, we'd still need a mutex to protect the shared pokemon variable.
	- If we use a channel just to signal, a mutex is still necessary to manage access to the shared state.
	- If we check for Pikachu in the producer and then send it trough the channel, we'd also need a mutex. On top of that, we'd violate the separation of concerns principle, where the
		producer is taking on the logic that really belongs to the consumer.

	That said, when multiple goroutines are modifying shared data, a mutex is still necessary to protect it. We'll often see a combination of channels and mutexes in these cases to ensure
	proper synchronization and data safety.

8. Okay, but what about broadcasting signals?
	Good question! We can indeed mimic a broadcast signal to all waiting goroutines using a channel by simply closing it `close(ch)`. When we close a channel, all goroutines receiving
	from that channel get notified. But keep in mind, a closed channel can't be reused once it's closed, it stays closed.

9. What's sync.Cond good for, then?
	Well there are certain scenarios where sync.Cond can be more appropriate than channels:
		- With a channel, we can either send a signal to one goroutine by sending a value or notify goroutines by closing the channel, but we cannot do both, sync.Cond gives us more
			fine-grained control. We can call Signal() to wake up a single goroutine or Broadcast() to wake up all of them.
		- We can call Broadcast() as many times as we need, which channels can't do once they're closed (closing a closed channel will trigger a panic)
		- Channels don't provide a built-in way to protect shared data - we'd need to manage that separately with a mutex. sync.Cond, on the other hand, gives us a more integrated
			approach by combining locking and signaling in one package (and better performance).

10. Why is the lock embedded in sync.Cond?
	In theory, a condition variable like sync.Cond doesn't have to be tied to a lock for its signaling to work.

	We could have the users manage their own locks outside of the condition variable, which might sound like it gives more flexibility. It's not really a technical limitation but more
	about human error.

	Managing it manually can easily lead to mistakes because the pattern isn't really intuitive, we have to unlock the mutex before calling Wait(), then lock it again when the goroutine
	wakes up. This process can feel awkward and is pretty prone to errors, like forgetting to lock or unlock at the right time.

	Typically, goroutines that call cond.Wait() need to check some shared state in a loop, like this:
		for !checkSomeSharedState() {
			cond.Wait()
		}
	The lock embedded in sync.Cond() helps handle the lock/unlock process for us, making the code cleaner and less error prone, we'll discuss the pattern in detail soon.


	HOW TO USE sync.Cond?
1. We always lock the mutex before waiting (.Wait()) on the condition, and we unlock it after the condition is met.
	func f(c *sync.Cond) {
		c.L.Lock()

		for !condition() {
			c.Wait()
		}

		// ...

		c.L.Unlock()

	}

2. Plus, we wrap the waiting condition inside a loop, here's a refresher:
	// Consumer
	go func() {
		cond.L.Lock()
		defer cond.L.Unlock()

		// waits until Pickachu appears
		for curPokemon != "Pickachu" {
			cond.Wait()
		}

		fmt.Println("Caught:", curPokemon)
	}

3. When we call Wait() on a sync.Cond, we're telling the current goroutine to hang tight until some condition is met. Here's what's happening behind the scenes:
	- The goroutine gets added to a list of other goroutines that are also waiting on this same condition. All these goroutines are blocked, meaning they can't continue until they're
		"woken up" by either a Signal() or Broadcast() call.
	- The key part here is that the mutex must be locked before calling Wait() because Wait() does something important, it automatically releases the lock (calls Unlock()) before putting
		the goroutine to sleep. This allows other goroutines to grab the lock and do their work while the original goroutine is waiting.
	- When the waiting goroutine gets woken up (by Signal() or Broadcast()), it doesn't immediately resume work. First, it has to re-acquire the lock (Lock()).

4. Here's a look at how Wait() works under the good:

	func (c *Cond) Wait() {
		// Check if Cond has been copied and panic() if it's been copied
		c.checker.check()

		// Get the ticket number
		t := runtime_notifyListAdd(&c.notify)

		// Unlock the mutex
		c.L.Unlock()

		// Suspend the goroutine until being woken up
		runtime_notifyListWait(&c.notify, t)

		// Re-lock the mutex
		c.L.Lock()
	}

	Even though, it's simple, we can take 4 main points:
		- There's a checker to prevent copying the `Cond` instance, it would be panic if we do so.
		- Calling cond.Wait() immediately unlocks the mutex, so the mutex must be locked before calling cond.Wait(), otherwise, it'll panic.
		- After being woken up, cond.Wait() re-locks the mutex, which means we'll need to unlock it again after we're done with the shared data.
		- Most of sync.Cond's functionality is implemented in the Go runtime with an internal data structure called `notifyList`, which uses a ticket-based system for notifications.

5. Because of this lock/unlock behavior, there's a typical pattern we'll follow when using sync.Cond.Wait() to avoid common mistakes:
	// Consumer function
	func f(cond *sync.Cond) {
		c.L.Lock()
		defer c.L.Unlock()

		for !condition {
			c.Wait()
		}

		// ... make use of condition
	}

6. Why not just use c.Wait() directly without a loop?
	When Wait() returns, we can't just assume that the condition we're waiting for is immediately true. While our goroutine is waking up, other goroutines could've messed with the shared
	state and the condition might not be true anymore. So, to handle this properly, we always want to use Wait() inside a loop.

	We've also mentioned this delay issue in the Pokemon example.

	The loop keeps things in check by continuously testing the condition, and only when that condition is true does our goroutine move forward.


	FUNCTIONS Cons.Signal() & Cond.Broadcast()
1. The Signal() method is used to wake up a single goroutine that's currently waiting on a condition variable.
	- If there are no goroutines currently waiting, Signal() doesn't do anything, it's basically no-op in this case.
	- If there are goroutines waiting, Signal() wakes up `the first` one in the queue. So if we've fired off a bunch of goroutines, like from 0 to N, the 0th goroutine will be the first
		one woken up by the Signal() call.

2. The idea of SignalGoroutinePriority() is that Signal() is used to wake up one goroutine and tell it that the condition might be satisfied. Here's what the Signal() implementation
	looks like:
		func (c *Cond) Singal() {
			c.checker.check()

			runtime_notifyListNotifyOne(&c.notify)
		}

	We don't have to lock the mutex before calling Signal(), `but it's generally a good idea to do so`, especially if we're modifying shared data and it's being accessed concurrently.

3. The contents of Broadcast is:
	func (c *Cond) Broadcast() {
		c.checker.check()

		runtime_notifyListNotifyAll(&c.notify)
	}

	When we call Broadcast(), it wakes up all the waiting goroutines and removes them from the queue. The internal logic here is still pretty simple, hidden behind the
	runtime_notifyListNotifyAll() function.

	In the case of Broadcast(), all the goroutines are woken up within the 100 milliseconds, but `there's no specific order how they're woken up`.

	When Broadcast() is called, it marks all the waiting goroutines as ready to run, but they don't run immediately, they're picked based on the Go scheduler's underlying algorithm,
	which can be a bit unpredictable.


	HOW IT WORKS INTERNALLY
1. Copy checker
	The copyChecker (copyChecker) in the sync.Package is designed to catch if a `Cond` object has been copied after it's been used for the first time. The `first time` could be any of
	the public methods like Wait(), Signal(), Broadcast().

	If the `Cond` gets copied after that first use, the program will panic with the error `sync.Cond is copied`

	We might have seen something similar in sync.WaitGroup or sync.Pool, where they use a `noCopy` field to prevent copying, but in those cases, it just avoids the issue without causing
	a panic.

	Now, this `copyChecker` is actually just a `uintptr`, which is basically an `integer that holds a memory address`, here's how it works:
		- After the first time we use sync.Cond, the `copyChecker` stores the memory address of itself, basically pointing to the `cond.copyChecker` object.
		- If the object gets copied, the memory address of the copy checker (&cond.copyChecker) changes (since a new copy lives in a different location in memory), but the `uintptr` that
			the copy checker holds doesn't change.

	The check is simple:
		Compare the memory addresses. If they're different, boom, there's a panic.

	Even though this logic is simple, the implementation might seem a bit tricky if we're not familiar with Go's atomic ops and the unsafe package.
		// copyChecker holds back pointer to itself to detect object copying
		type copyChecker uintptr

		func (c *copyChecker) check() {
			if 	uintptr(*c) != uintptr(unsafe.Pointer(c)) &&
				!atomic.CompareAndSwapUintptr((*uintptr)(c), 0, uintptr(unsafe.Pointer(c))) &&
				uintptr(*c) != uintptr(unsafe.Pointer(c)) {
					panic("sync.Cond is copied")
				}
		}

	Let's break this down into two main checks, since the first and last checks are doing pretty much the same thing.
		1) The first check uintptr(*c) != uintptr(unsafe.Pointer(c)), looks to see if the memory address has changed. If it has, the object's been copied. But, there's a catch, if this
			is the first time `copyChecker` is being used, it'll be 0 since it's not initialized yet.

		2) The second check !atomic.CompareAndSwapUintptr((*uintptr)(c), 0, uintptr(unsafe.Pointer(c))), is where we use a Compare-And-Swap (CAS) op to handle both initialization and
			checking.
				- If CAS succeeds, it means the `copyChecker` was just initialized, so the object hasn't been copied yet, and we're good to go.
				- If CAS fails, it means the `copyChecker` was already initialized, and we need to do that final check (uintptr(*c)) != uintptr(unsafe.Pointer(c)) to make sure the object
					hasn't been copied.
	The final check uintptr(*c) != uintptr(unsafe.Pointer(c)) (it's the same as the first check) makes sure that the object hasn't been copied after all that.

2. Why the extra check at the end? Isn't two checks enough to panic?
	The reason for the third check is that the first and second checks aren't atomic.

	There's a race condition during the initialization.

	if this is the first time the copyChecker is being used, it hasn't been initialized yet, and its value will be zero. In this case, the check will pass incorrectly, even though the
	object hasn't been copied but just hasn't been initialized.


	notifyList - TICK-BASED NOTIFICATION LIST
1. Beyond the locking and copy-checking mechanisms, one of the other important parts of the sync.Cond is the `notifyList`.

2. The internal structure of Cond type is:

	type Cond struct {
		noCopy 	noCopy
		L 		Locker

		`notify 	notifyList`

		checker copyChecker
	}

	The version in the `sync` package:
		type notifyList struct {
			wait 	uint32
			notify 	uint32

			lock 	uintptr

			head 	unsafe.Pointer
			tail 	unsafe.Pointer
		}

	Now, the `notifyList` in the `sync` package and the one in the `runtime` package are different but share the same name and memory layout (in purpose). To really understand how it
	works, we'll need to look at the version in the `runtime` package:

		package runtime

		type notifyList struct {
			wait atomic.Uint32
			notify uint32

			lock mutex

			head *sudog // the first goroutine in the linked list
			tail *sudog // the last goroutine in the linked list
		}

	If we look at the `head` and `tail`, we probably guess this is some kind of linked list, and we'd be right. It's a linked list of `sudog` (short of `pseudo-goroutine`), which
	represents a goroutine waiting on synchronization events, like waiting to receive or send data on a channel or waiting on a condition variable.

	The `head` and `tail` are pointers to the first and last goroutine in this list.

	Meanwhile, the `wait` and `notify` fields act as `ticket` numbers that are continuously increasing, each representing a position in the queue of waiting goroutines.
		- `wait`: This number represents the `next` ticket that's going to be issued to a waiting goroutine.
		- `notify`: This tracks the next ticket number that's supposed to be notified, or woken up.
	And that's the core idea behind `notifyList`, let's put them together to see how it works.


	Wait()'s notifyListAdd() FUNCTION
1. When a goroutine is about to wait for a notification, it calls notifyListAdd() to get its
	"ticket" first.

	func (c *Cond) Wait() {
		// Check if the `c` has been copied
		c.checker.check()

		// Get the ticket number
		t := runtime_notifyListAdd(*c.notify)

		c.L.Unlock()
		// Add the goroutine to the list and suspend it
		runtime_notifyListWait(&c.notify, t)
		c.L.Lock()
	}

	func notifyListAdd(l *notifyList) uint32 {
		return l.wait.Add(1) - 1
	}

	The ticket assignment is handled by an atomic counter. So when a goroutine calls
	notifyListAdd(), that counter ticks up, and the goroutine is handed the next available ticket
	number.

	Every goroutine gets its own unique ticket number, and this process happens `without
	any lock`. This means that multiple goroutines can request tickets at the same time without
	waiting for each other.

	For example, if the current ticket number at 5, the goroutine that calls notifyListAdd() next
	will get ticket number 5, and the wait counter will then bump up to 6, ready for the next one
	in line. The wait `field` always points to the next ticket number that'll be issued.

	But here's where things get a little tricky

	Since many goroutines can grab a ticket at the same time, there's a small gap between when they call notifyListAdd() and when they actually enter notifyListWait(). Their order isn't
	necessarily guaranteed, even though the ticket numbers are issued sequentially, the order in which the goroutines get added to the linked list might not be `1,2,3`. Instead, it
	could end up being `3,2,1`, or `2,1,3`, it all depends on the timing.

*/

func CondUsage() {
	var (
		pokemons = []string{
			"Pickachu",
			"Charmander",
			"Squirtle",
			"Bulbasaur",
			"Jigglypuff",
		}
		curPokemon string

		cond = sync.NewCond(&sync.Mutex{})

		wg sync.WaitGroup
	)

	wg.Add(1)
	go func() {
		defer wg.Done()

		cond.L.Lock()
		defer cond.L.Unlock()

		for curPokemon != "Pickachu" {
			cond.Wait()
		}

		fmt.Println("caught:", curPokemon)
		curPokemon = ""
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		for range 100 {
			time.Sleep(time.Millisecond)

			// It's the place when other goroutines can change
			// the shared variable even if the consumer has been
			// woken up.
			cond.L.Lock()
			curPokemon = pokemons[rand.IntN(len(pokemons))]
			cond.L.Unlock()

			cond.Signal()
		}
	}()

	wg.Wait()
}

func SignalGoroutinePriority() {
	var (
		cond = sync.Cond{L: &sync.Mutex{}}
	)

	for i := range 10 {
		go func(i int) {
			cond.L.Lock()
			defer cond.L.Unlock()

			cond.Wait()

			fmt.Println(i)
		}(i)
	}

	time.Sleep(100 * time.Millisecond)
	cond.Signal()
	time.Sleep(100 * time.Millisecond)

}

func BroadcastUsage() {
	var (
		cond = sync.Cond{L: &sync.Mutex{}}
	)

	for i := range 10 {
		go func(i int) {
			cond.L.Lock()
			defer cond.L.Unlock()

			cond.Wait()

			fmt.Println(i)
		}(i)
	}

	time.Sleep(100 * time.Millisecond)
	cond.Broadcast()
	time.Sleep(100 * time.Millisecond)
}

func SyncCondCopyingPanic() {
	var cond = sync.Cond{L: &sync.Mutex{}}

	cond.Signal()

	// panic: sync.Cond is copied
	// func passes lock by value: sync.Cond contains sync.noCopycopylocksdefault
	func(cond sync.Cond) {
		cond.Signal()
	}(cond)
}
