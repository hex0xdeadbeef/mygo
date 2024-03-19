package once

import (
	"fmt"
	"sync"
)

// The "sync.Once" instances are used to initialize the varible once.
// We should provide the method Do() with the function, that will initialize the variable. The "sync.Once" counts only the number of calls the Do() method instead of
// counting the count of the unique function that were called within it.

// The check how many times the go itself uses the sync.Once: grep -ir sync.Once $(go env GOROOT)/src | wc -l
// The result is: 234 (18.03.2024)

// We should formalize the coupling by wrapping any usage of "sync.Once" in a small lexical block:
// 1. A small function
// 2. Wrapping both in a type

func Using() {
	var (
		count     int
		increment = func() {
			count++
		}

		once sync.Once

		wg sync.WaitGroup
	)

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			wg.Done()
			once.Do(increment)
		}()
	}

	wg.Wait()

	fmt.Println(count)
}

func RewritingDoFunction() {
	var (
		count int
		once  sync.Once

		increment = func() {
			count++
		}
		decrement = func() {
			count--
		}
	)

	once.Do(increment)
	once.Do(decrement)

	fmt.Println(count)

}

/*
fatal error: all goroutines are asleep - deadlock!

goroutine 1 [sync.Mutex.Lock]:
sync.runtime_SemacquireMutex(0x0?, 0x0?, 0x0?)
	/usr/local/go/src/runtime/sema.go:77 +0x28
sync.(*Mutex).lockSlow(0x1400000e0c4)
	/usr/local/go/src/sync/mutex.go:171 +0x330
sync.(*Mutex).Lock(0x1400000e0c4)
	/usr/local/go/src/sync/mutex.go:90 +0x94
sync.(*Once).doSlow(0x1400000e0c0, 0x14000078f18)
	/usr/local/go/src/sync/once.go:70 +0x34
sync.(*Once).Do(0x1400000e0c0, 0x14000078f18)
	/usr/local/go/src/sync/once.go:65 +0x44
concurrency/pkg/ch02/syncpkg/once.OnceDeadlock.func2()
	/Users/dmitriymamykin/Desktop/goprojects/concurrency/pkg/ch02/syncpkg/once/once.go:69 +0x30
sync.(*Once).doSlow(0x1400000e0d0, 0x14000078ef8)
	/usr/local/go/src/sync/once.go:74 +0x148
sync.(*Once).Do(0x1400000e0d0, 0x14000078ef8)
	/usr/local/go/src/sync/once.go:65 +0x44
concurrency/pkg/ch02/syncpkg/once.OnceDeadlock.func1()
	/Users/dmitriymamykin/Desktop/goprojects/concurrency/pkg/ch02/syncpkg/once/once.go:68 +0x34
sync.(*Once).doSlow(0x1400000e0c0, 0x14000078f18)
	/usr/local/go/src/sync/once.go:74 +0x148
sync.(*Once).Do(0x1400000e0c0, 0x14000078f18)
	/usr/local/go/src/sync/once.go:65 +0x44
concurrency/pkg/ch02/syncpkg/once.OnceDeadlock()
	/Users/dmitriymamykin/Desktop/goprojects/concurrency/pkg/ch02/syncpkg/once/once.go:71 +0x180
main.main()
	/Users/dmitriymamykin/Desktop/goprojects/concurrency/bin/main.go:7 +0x1c
Process 26536 has exited with status 2
Detaching
dlv dap (26484) exited with code: 0
*/

func OnceDeadlock() {
	var onceA, onceB sync.Once
	var initB func()

	initA := func() { onceB.Do(initB) } // |1| A is locked, go to B |2| Try to lock A, but A has already locked. DEADLOCK
	initB = func() { onceA.Do(initA) }  // |1| B is locked, go to B

	onceA.Do(initA)
}
