package conditions

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// An "event" is any arbitrary signal between two or more goroutines that carries mo information other than the fact it's occured.
// It's a way for a goroutine to efficiently sleep until it was signaled to wake and check its condition. This is exactly what the Cond type does for us.

// Wait() does the following:
// 1. Unlocks before entering
// 2. Blocks and suspends the goroutine it's called from
// 3. After a Signal()/Broadcast() and before exiting it calls Lock()

// Internally, the runtime maintains a FIFO list of goroutines waiting to be signalled.
// 1.Signal() is one of the two methods that the Cond type provides for notifying goroutines blocked on a wait call that the condition has been
// triggered. Signal() finds the goroutine that's been waiting the longest and notifies that.
// 2. Broadcast() is the next method of two methods that the Cond type provides for notifying goroutines blocked on a wait call that the condition
// has been triggered. Broadcast() sends a signal to all goroutines that are waiting.

// The using of Cond type is more performant than utilizing channels.

func Example() {

	c := sync.NewCond(&sync.Mutex{}) // The NewCond takes in a type that satisfies the "sync.Locker" interface. This is what facilitate coordination with other goroutines

	// in a concurrent-safe way.
	c.L.Lock() // Call to Wait() automatically calls Unlock() on the Locker when entered
	for lessThanFifty() == false {
		// Wait to be notified that the condition has occured. This is a blocking call and the goroutine will be suspended.
		c.Wait()
	}
	c.L.Unlock() // When the call to Wait() exits, it calls Lock on the Locker for the condition.

}

func lessThanFifty() bool {
	return rand.Intn(100) < 50
}

func Using() {
	c := sync.NewCond(&sync.Mutex{})

	queue := make([]interface{}, 0, 10)
	removeFromQueue := func(delay time.Duration) {
		// Sleep for a moment
		time.Sleep(delay)

		// Lock others removers from this critical section
		c.L.Lock()

		// Remove the first element
		queue = queue[1:]
		fmt.Println("Removed from the queue!")

		// Unlock others removers from the deleting
		c.L.Unlock()

		// Wake up the only one goroutine to now that an element has been deleted.
		c.Signal()
	}

	for i := 0; i < 10; i++ {
		c.L.Lock()

		// Every time when the length of the queue is 2 we block and suspend the main goroutine
		for len(queue) == 2 {
			c.Wait()
		}

		// At the start we'll add 2 elements
		fmt.Println("Adding to queue")
		queue = append(queue, struct{}{})

		// Start some removers
		go removeFromQueue(1 * time.Second)

		c.L.Unlock()
	}
}

type Button struct {
	Clicked *sync.Cond
}

func ClickBroadcasting() {
	var (
		button = Button{Clicked: sync.NewCond(&sync.Mutex{})}

		subscribe = func(c *sync.Cond, fn func()) {
			var goroutineRunning sync.WaitGroup

			goroutineRunning.Add(1)
			go func() {
				goroutineRunning.Done()
				// Lock the goroutine
				c.L.Lock()
				// Unlock the goroutine after return call
				defer c.L.Unlock()

				c.Wait()
				fn()
			}()

			goroutineRunning.Wait()
		}

		clickRegistered sync.WaitGroup
	)

	clickRegistered.Add(3)

	subscribe(button.Clicked, func() {
		defer clickRegistered.Done()
		fmt.Println("Maximizing window.")
	})
	fmt.Println(1)

	time.Sleep(750 * time.Millisecond)

	subscribe(button.Clicked, func() {
		defer clickRegistered.Done()
		fmt.Println("Displaying annoying dialog box!")
	})
	fmt.Println(2)
	time.Sleep(750 * time.Millisecond)

	subscribe(button.Clicked, func() {
		defer clickRegistered.Done()
		fmt.Println("Mouse clicked.")
	})
	fmt.Println(3)
	time.Sleep(750 * time.Millisecond)

	button.Clicked.Broadcast()

	clickRegistered.Wait()

}
