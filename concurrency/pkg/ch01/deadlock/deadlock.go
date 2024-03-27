package deadlock

import (
	"fmt"
	"sync"
	"time"
)

// 1. a1 locks its mutex
// 2. G1 pr G2 waits a second
// 3. b2 locks its mutex
// 4. The second from the 2. passed and b1 locks its mutex, so there's the deadlock

// The result is: infinite waiting.

type value struct {
	mu    sync.Mutex
	value int
}

func Deadlock() {
	var (
		wg   sync.WaitGroup
		a, b value

		printSum = func(goroutineNumber int, v1, v2 *value) {
			fmt.Printf("First block by %d goroutine.\n", goroutineNumber)
			v1.mu.Lock()
			defer v1.mu.Unlock()

			time.Sleep(2 * time.Second)

			fmt.Printf("Second block by %d goroutine.\n", goroutineNumber)
			v2.mu.Lock()
			defer v2.mu.Unlock()

			fmt.Printf("sum=%v\n", v1.value+v2.value)
		}
	)

	wg.Add(2)
	go func() {
		defer wg.Done()
		printSum(1, &a, &b)
	}()

	go func() {
		defer wg.Done()
		printSum(2, &b, &a)
	}()

	wg.Wait()

}

/*
sync.runtime_Semacquire(0x14000114010?)
	/usr/local/go/src/runtime/sema.go:62 +0x2c
sync.(*WaitGroup).Wait(0x14000114010)
	/usr/local/go/src/sync/waitgroup.go:116 +0x98
concurrency/pkg/ch01/deadlock.Deadlock()
	/Users/dmitriymamykin/Desktop/goprojects/concurrency/pkg/ch01/deadlock/deadlock.go:49 +0x290
main.main()
	/Users/dmitriymamykin/Desktop/goprojects/concurrency/bin/main.go:7 +0x1c

goroutine 34 [sync.Mutex.Lock]:
sync.runtime_SemacquireMutex(0xbfae71ab?, 0x8?, 0x14000074d98?)
	/usr/local/go/src/runtime/sema.go:77 +0x28
sync.(*Mutex).lockSlow(0x14000114030)
	/usr/local/go/src/sync/mutex.go:171 +0x330
sync.(*Mutex).Lock(0x14000114030)
	/usr/local/go/src/sync/mutex.go:90 +0x94
concurrency/pkg/ch01/deadlock.Deadlock.func1(0x14000114020, 0x14000114030)
	/Users/dmitriymamykin/Desktop/goprojects/concurrency/pkg/ch01/deadlock/deadlock.go:32 +0xc4
concurrency/pkg/ch01/deadlock.Deadlock.func2()
	/Users/dmitriymamykin/Desktop/goprojects/concurrency/pkg/ch01/deadlock/deadlock.go:42 +0xc0
created by concurrency/pkg/ch01/deadlock.Deadlock in goroutine 1
	/Users/dmitriymamykin/Desktop/goprojects/concurrency/pkg/ch01/deadlock/deadlock.go:40 +0x184

goroutine 35 [sync.Mutex.Lock]:
sync.runtime_SemacquireMutex(0xbfae71ed?, 0x0?, 0x14000126d98?)
	/usr/local/go/src/runtime/sema.go:77 +0x28
sync.(*Mutex).lockSlow(0x14000114020)
	/usr/local/go/src/sync/mutex.go:171 +0x330
sync.(*Mutex).Lock(0x14000114020)
	/usr/local/go/src/sync/mutex.go:90 +0x94
concurrency/pkg/ch01/deadlock.Deadlock.func1(0x14000114030, 0x14000114020)
	/Users/dmitriymamykin/Desktop/goprojects/concurrency/pkg/ch01/deadlock/deadlock.go:32 +0xc4
concurrency/pkg/ch01/deadlock.Deadlock.func3()
	/Users/dmitriymamykin/Desktop/goprojects/concurrency/pkg/ch01/deadlock/deadlock.go:46 +0xc0
created by concurrency/pkg/ch01/deadlock.Deadlock in goroutine 1
	/Users/dmitriymamykin/Desktop/goprojects/concurrency/pkg/ch01/deadlock/deadlock.go:44 +0x288

*/
/*
Coffman conditions:
1) A concurrent process holds exclusive rights to a resource at any on time
2) A concurrent process must simultaneously hold a resource and be waiting for an additional resource
3) A resource that is held by a concurrent process can only be released by that process, so it fulfills this condition.
4) A concurrent process (P1) must be waiting on a chain of other concurrent processes (P2), which are in turn waiting on it (P1), so it fulfills this final condition too.
*/
