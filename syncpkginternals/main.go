package main

import (
	"bytes"
	"fmt"
	"sync"
	"time"
)

func main() {
	// condUsage()

	// onceUsage()

	poolUsage()
}

// Cond usage
func condUsage() {
	var (
		c   = sync.NewCond(&sync.Mutex{})
		msg = "hello"

		listenForSignal = func() {
			c.L.Lock()
			defer c.L.Unlock()

			for msg == "hello" {
				fmt.Println("cond is ok")

				fmt.Println("unlock the c.L and suspends")
				// Unlocks c.L and suspends this goroutine
				c.Wait()
				// After resuming the execution by Singal() | Broadcast() locks c.L
				fmt.Println("lock the c.L after Signal() | Broadcast")
			}
		}

		signal = func() {
			c.L.Lock()
			defer c.L.Unlock()

			time.Sleep(2 * time.Second)
			fmt.Println("Changes msg")
			msg = "hello world"

			c.Signal()
		}
	)

	go listenForSignal()

	signal()

}

// Once usage
func onceUsage() {
	type User struct {
		id   int
		name string

		mu *sync.Mutex
	}

	var (
		u  User
		wg = &sync.WaitGroup{}

		setupInvocationsCnt int
		setupUser           = func() {
			u = User{id: 333, name: "Denis ", mu: &sync.Mutex{}}

			setupInvocationsCnt++

		}

		once    = &sync.Once{}
		execute = func() {
			for i := 0; i < 100; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()

					once.Do(setupUser)

					u.mu.Lock()
					defer u.mu.Unlock()
					fmt.Println(u.name)

				}()

			}
		}
	)

	execute()
	wg.Wait()

	fmt.Println(setupInvocationsCnt)

}

func poolUsage() {
	type PoolCounter struct {
		mu *sync.Mutex

		pool                  *sync.Pool
		newFunc               func() any
		poolNewInvocationsCnt int
	}

	var (
		pc = &PoolCounter{pool: &sync.Pool{}, mu: &sync.Mutex{}}

		wg = new(sync.WaitGroup)
	)
	pc.newFunc = func() any {
		pc.mu.Lock()
		defer pc.mu.Unlock()

		pc.poolNewInvocationsCnt++

		b := bytes.NewBuffer(make([]byte, 8))
		return b
	}

	pc.pool.New = pc.newFunc

	for i := 0; i < 100_000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			b := pc.pool.Get().(*bytes.Buffer)
			b.Write([]byte("abcdefg"))
			b.WriteByte('a')
			b.WriteRune('П')
			b.WriteString("хуй")

			fmt.Println(b.String())
			b.Reset()

			pc.pool.Put(b)
		}()
	}

	wg.Wait()

	fmt.Println(pc.poolNewInvocationsCnt)
}
