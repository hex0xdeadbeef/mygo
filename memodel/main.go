package main

import (
	"fmt"
	"sync"
)

func main() {
	// MutexUsage()

	// OnceUsage()

	// wrongSyncA()

	// wrongSyncB()

	wrongSyncC()
}

func MutexUsage() {
	var (
		mu sync.Mutex
		a  string

		f = func() {
			// It's the runtime error if the mutex hasn't been locked until there's a try to unlock it
			defer mu.Unlock()

			a = "hello, world!"
		}

		mainFunc = func() {
			// The first lock is guaranteed to be completed
			mu.Lock()
			go f()

			// The second lock blocks the caller if there's already the lock allocated
			mu.Lock()
			fmt.Println(a)
		}
	)

	mainFunc()
}

func OnceUsage() {
	var (
		a    string
		once sync.Once
		wg   sync.WaitGroup

		setupA = func() {
			a = "hello, world!"
		}

		printA = func() {
			once.Do(setupA)
			fmt.Println(a)
		}

		multiplePrint = func() {
			for i := 0; i < 100; i++ {
				wg.Add(1)

				go func() {
					defer wg.Done()
					printA()
				}()

			}
		}
	)

	multiplePrint()

	wg.Wait()
}

// Incorrect synchronization
/*
!!! All these examples must be fixed by the explicit synchronization !!!
*/

func wrongSyncA() {
	var (
		a, b int

		f = func() {
			a = 1
			b = 2
		}

		g = func() {
			fmt.Println(b)
			fmt.Println(a)
		}

		execute = func() {
			go f()
			g()
		}
	)

	execute()
}

func wrongSyncB() {
	var (
		wg   sync.WaitGroup
		once sync.Once

		a    string
		done bool

		setup = func() {
			// done = true // can be arranged here after swap
			a = "hello world!"
			done = true
		}

		doPrint = func() {
			defer wg.Done()

			if !done {
				once.Do(setup)
			}

			fmt.Println(a)
		}

		execute = func() {
			wg.Add(2)

			go doPrint()
			go doPrint()
		}
	)

	execute()
	wg.Wait()
}

func wrongSyncC() {
	var (
		a    string
		done bool

		setup = func() {
			// done = true // can be arranged here after swap
			a = "hello world!"
			done = true
		}

		execute = func() {
			// runtime.GOMAXPROCS(1)

			go setup()

			for !done {
				// runtime.Gosched()
			}

			fmt.Println(a)
		}
	)

	execute()

}

func wrongSyncD() {
	type (
		T struct {
			msg string
		}
	)

	var (
		g *T

		setup = func() {
			t := new(T)
			// g = t // can be arranged here after a swap by a compiler
			t.msg = "hello world"

			g = t
		}

		execute = func() {
			go setup()

			for g == nil {
			}

			fmt.Println(g.msg)
		}
	)

	execute()
}
