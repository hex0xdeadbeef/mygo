package mutexes

import (
	"fmt"
	"math"
	"os"
	"sync"
	"text/tabwriter"
	"time"
)

// A Mutex shares memory by creating a convention developers must follow to synchromize access to the memory.

/*
	Critical sections are expensive to enter and exit a critical section, and so generally people attempt to minimize the time spent in critical sections. One strategy for doing so
	is to reduce the cross-section. There may be memory that needs to be shared between multiple concurrent processes, but perhaps no all of these processes will read and write to
	this memory. In this cases we use "sync.RWMutex"
*/

func Incrementing() {
	const (
		goroutinesNumber = 5
	)

	var (
		count int

		m  sync.Mutex
		wg sync.WaitGroup

		increment = func() {
			m.Lock()
			defer m.Unlock() // Works in the situations when the function is panicing. If m.Unlock() isn't deferred, it could be a deadlock.
			count++
			fmt.Printf("Incrementing: %d\n", count)
		}

		decrement = func() {
			m.Lock()
			defer m.Unlock() // Works in the situations when the function is panicing. If m.Unlock() isn't deferred, it could be a deadlock.
			count--
			fmt.Printf("Decrementing: %d\n", count)
		}
	)

	// Increment
	for i := 0; i < goroutinesNumber; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			increment()
		}()
	}

	// Decrement
	for i := 0; i < goroutinesNumber; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			decrement()
		}()
	}

	wg.Wait()
	fmt.Println("Arithmetic complete.")
}

func RWUsing() {
	var (
		// producer always uses the simple mutex
		producer = func(wg *sync.WaitGroup, l sync.Locker) {
			defer wg.Done()

			for i := 5; i > 0; i-- {
				l.Lock()
				l.Unlock()

				time.Sleep(1)
			}
		}

		// observer uses the different mutexes in the tests
		// 1. In the first test it uses the usual mutex that locks the critical section from all the intervenings
		// 2. In the second test it uses only the RWMutex so that all other goroutines can get the resource
		observer = func(wg *sync.WaitGroup, l sync.Locker) {
			defer wg.Done()
			l.Lock()
			defer l.Unlock()
		}

		test = func(count int, mutex, rwMutex sync.Locker) time.Duration {
			var wg sync.WaitGroup

			beginTestTime := time.Now()

			// Start the single locking with the usual mutex
			wg.Add(1)
			go producer(&wg, mutex)

			// Start multiple locking with the RWMutex/Mutex with the current count. The amout of goroutines will depend of the current pow of the 2.
			for i := count; i > 0; i-- {
				wg.Add(1)
				go observer(&wg, rwMutex)
			}

			wg.Wait()
			return time.Since(beginTestTime)
		}

		m  sync.RWMutex
		tw = tabwriter.NewWriter(os.Stdout, 0, 1, 2, ' ', 0)
	)
	defer tw.Flush()

	fmt.Fprintf(tw, "Readers\tRWMutex\tMutex\n")

	for i := 0; i < 20; i++ {
		count := int(math.Pow(2, float64(i)))

		fmt.Fprintf(
			tw,
			"%d\t%v\t%v\n",
			count,
			test(count, &m, m.RLocker()),
			test(count, &m, &m),
		)

	}
}
