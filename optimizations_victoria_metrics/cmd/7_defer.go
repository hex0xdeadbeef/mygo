package main

import (
	"fmt"
	"math/rand/v2"
	"sync"
)

/*
	DEFER
1. The defer statement actually has 3 types:
	- Open-coded defer
	- Heap-allocated defer
	- Stack-allocated defer

	Each one has different scenarios where they're best used, which is good to know if we want to optimize performance.

2. Defer is called when the function returns, even if a panic (runtime error) happens.

3. We mustn't put a wrapped recover handling into a closure. Instead, we must defer this processing function directly.

4. Each goroutine has its own defer linked list. So, we should process panics properly in these goroutines. We cannot handle panic happened in goroutineA in the goroutineB. We already know that defer chains belong to a specific goroutine.

	If we don't handle a panic in child goroutine, an application will crash.

	DEFER ARGUMENTS, INCLUDING RECEIVER ARE IMMEDIATELY EVALUATED.
1. To solve the problem of capture-the-value principle, we can use passing the pointer / closures, but closures are more idiomatic in Go, especially when dealing with simple variables captures.

	DEFER TYPES
1. Defer can be one of the following three types:
	1) Heap-allocated
	2) Stack-allocated
	3) Open-coded defer

2. When we call defer, we're creating a structure called a defer object `_defer`, which holds all the necessary information about the deferred call.

	This object gets pushed into the goroutine's defer chain, as we discussed earlier.

	Every time the function exits, whether normally or due to an error, the compiler ensures a call to runtime.deferreturn. This function is responsible for unwinding the chain of deferred calls, retrieving the stored information from the defer objects, and then executing the deferred functions in the correct order.

3. The difference between heap-allocated and stack-allocated defer types is where the defer object is allocated. Below Go 1.13, we only had heap-allocated defer.

4. Currently, in Go 1.22, if we use defer in a loop, it'll be heap-allocated.

	The heap allocation here is necessary because the number of defer objects can change at runtime. So, the heap ensures that the program can handle any number of defers, no matter how many or where they appear in the function, without bloating the stack.

	Now, don't panic, heap allocation is indeed considered bad for performance, but Go tries to optimize that by using a pool of defer objects.

	We have two pools:
		1) A local cache pool of the logical processor P to avoid lock contention
		2) A global cache pool shared and taken by all the goroutines, which then put defer objects into processor P's local pool.

	DEFER IN IF/FOR/FOR-RANGE STRUCTURES
1. Good catch, putting defer in the `if` statement can be unpredictable.

	Since Go 1.13, the defer can be stack-allocated and this means we craft the `_defer` object in the stack, then push it into the goroutine's defer chain.

	If the defer statement within the `if` block is invoked only once and not in a loop or another dynamic context, it benefits from the optimization in Go 1.13, meaning the defer object will be stack-allocated.

	With this optimization, according to the `Open-coded defers proposal` in the `cmd/go` binary, this optimization applies to 363/370 static defer sites. As a result, these sites see a 30% performance improvement compared to the previous approach where defer objects were heap-allocated.

2. Open-coded defer
	What if we put the defer at the end of the function?

	The performance of a direct call is much better than the other two. As of Go 1.13, most defer ops take about 35 nanoseconds (down from 50 ns in Go 1.12).

	In contrast, a direct call takes about 6 nanoseconds.

	If a function has at least one heap-allocated defer, any defer in the function won't be inlined or open-coded. That means, to optimize the function headAndStackAllocatedDefers(...), we should remove or move the heap allocated defer elsewhere.

	Another thing to keep in mind is that the product of the number of defers in the function and the number of return statements needs to be 15 or less to get into this category. This is because we put the defer before every return statement. Our binary code will get pretty bloated if we have too many exit paths like that.

	If the number of defer statements is more than 8, open-coded defer will not be applied.
*/

func erorrneousRecoverHandlingA() {
	defer func() {
		processRecover()
	}()

	panic("UNREACHABLE")
}

func processRecover() {
	if v := recover(); v != nil {
		fmt.Println("Recovered:", v)
	}
}

func rightPanicProcessing() {
	defer processRecover()

	panic("UNREACHABLE")
}

func erorrneousRecoverHandlingB() {
	defer processRecover()

	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		defer wg.Done()

		panic("UNREACHABLE")
	}()

	wg.Wait()
}

func deferInLoop() {
	for i := 0; i < 100; i++ {
		// Heap-allocated defer
		defer fmt.Println(i)
	}
}

func headAndStackAllocatedDefers() {
	if rand.IntN(100) > 50 {
		// Stack-allocated defer
		defer fmt.Println("Defer in if")
	}

	if rand.IntN(100) > 50 {
		// Stack-allocated defer
		defer fmt.Println("Defer in if")
	}

	for range rand.IntN(100) {
		// Heap-allocated defer
		defer fmt.Println("Defer in for")
	}
}
