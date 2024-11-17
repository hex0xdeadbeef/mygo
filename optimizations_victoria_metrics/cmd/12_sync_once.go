package main

import (
	"fmt"
	"math/rand"
	"sync"
)

/*
	sync.Cond
*/

/*
	WHAT IS sync.Once?
1. The sync.Once is exactly what we reach for when we need a function to run just one time, now matter how many times it gets called or how many goroutines hit it simultaneously.

2. It's perfect for initializing a singleton resource, something that should only happen once in ours application lifecycle, such as setting up a database conn pool, initializing a logger or maybe configuring our
	metrics system.

	var once sync.Once
	var conf Config

	func GetConfig() {
		once.Do(func() {
			conf = fetchConfig()
		})

		return conf
	}

	If GetConfig() is called multiple times, fetchConfig() is executed only once.

	The real benefit of sync.Once is that it delays certain ops until they are first needed (lazy-loading), which can improve runtime performance and reduce initial memory usage. For instance, if a large lookup table is
	created only when accessed for the first time, then memory and processing time for creating the table are saved until that moment.

3. Most of the time, it's better for initializing (external) resources than using init()

4. Now, something important to keep in mind:
	Once we've used a sync.Once object to run a function, that's it - we can't reuse it. Once it's done, it's done.

5. There's no built-in way to reset a sync.Once either. Once it's done its job, it's retired for good.

6. Now, here's interesting twist. If the function we pass to sync.Once panics while running, sync.Once still treats that as "mission accomplished". That means future calls to Do(f) won't run the function again. This can
	be tricky, especially of we're trying to catch the panic and handle the error afterward, there's no retry.

7. If we need to handle errors that might come out of `f`, it can get a little awkward to write. The thing is, if fetchConfig() fails, only the first call will receive the error, Any subsequent calls - second, third and
	so on - won't return that error. To make the behaviour consistent across calls, we need to make `err` variable a package scoped variable.


	FUNCTIONS OnceFunc, OnceValue, and OnceValues FUNCTIONS
1. From Go 1.21 onward, we get OnceFunc, OnceValue, and OnceValues functions. These are basically handy wrappers around sync.Once that make things smoother, without sacrificing any performance.

2. OnceFunc is pretty straightforward, it takes our function `f` and wraps it in another function that we can call as many times as we want, but `f` itself will only run once.

	Even if we call wrapper() a bunch of times, printOnce() runs only on the first call.

	Now, if `f` panics during that first execution, every future call to wrapper() will also panic with the same error. It's like it locks in the failure state, so our app doesn't continue like nothing went wrong when
	something critical didn't initialize correctly.

	But this doesn't really solve the problem of catching errors in a nice way.

3. OnceValue and OnceValues. These are cool because they remember the result of `f` after the first execution and just return the cached result on future calls.

	After that first call to `getConfigOnce`, it just hands us the same result without returning -1 directly. Beside the lazy loading, we don't need to deal with closures here.

4. But what about errors? Fetching some usually involves error handling.

	This is where sync.OnceValues comes in. It works like sync.OnceValue, but lets us return multiple values, including errors. So, we can cache both the result and any error that comes up during the first run.

	Now, `getIntegerOnce` behaves just like a normal function - it's still concurrent-safe, caches the result, and only incurs the cost of running the function once. After that, every call is cheap.

	So, if an error happens, is that cached too?
		Unfortunately, yes. Whether it's a panic or an error, both the result and the failure state are cached. So, the calling code needs to be aware that it might be dealing with a cached error or failure. If we need
		to retry, we'd have to create a new instance from sync.OnceValues to re-run the initialization.

	Also, in the example, we return an error as the second value to match the typical function signature we're used to, but honestly, it could be anything depending on what we need.


	HOW IT WORKS?
1. It's really simple compared to something like sync.Mutex, which is one of the trickiest synchronization tool in Go. So let's take a step back and think about how sync.Once is actually implemented.

2. The structure of sync.Cond

	type Once struct {
		done atomic.Uint32
		m Mutex
	}

	The first thing we'll notice is `done` field, which uses an atomic op.

	And here's an interesting detail, `done` is placed at the very top of the struct for a reason. On many CPU architectures (like x86-64), accessing the first field in a struct is faster because it sits at the base
	address of the memory block. This little optimization allows the CPU to load the first field more directly, without calculating memory offsets.

	Also, putting it at the top helps with inlining optimization, which we'll get to in a bit.

3. So, what's the first thing that comes to mind when implementing sync.Once to make sure a function runs only once? Using a mutex, right?

	func (o *Once) Do(func()) {
		o.m.Lock()
		defer o.m.Unlock()

		if o.done.Load() == 0 {
			o.done.Store(1)
			f()
		}
	}

	It's simple and gets the job done. The idea is that the mutex allows only one goroutine to enter the critical section at a time. Then if `done` is still 0 (meaning the function hasn't run yet), it sets `done` to 1
	and runs the function f()

	This is actually the original version of sync.Once, written by Rob Pike back in 2010.

	Now, the version we just looked at works fine, but it's not the most performant one. Even after the first call, every time Do(f) is called, it still grabs a lock, which means goroutines are waiting on each other. We
	can definitely do better by adding a quick exit if the task is already done.

		func (o *Once) Do(func()) {
			if o.done.LoadUint32(&o.done) == 1 {
				return
			}

			// slow path
			o.m.Lock()
			defer o.m.Lock()

			if o.done.Load() == 0 {
				o.done.Store(1)
				f()
			}
		}

	This gives us a `fast path`, when the `done` flag is set, we skip the lock entirely and just return immediately. Nice and quick. But, if the flag isn't set, we fall back to the `slower path`, which locks the mutex,
	rechecks `done`, and then runs the function.

	Now, we have to re-check `done` after acquiring the lock because there's a small window between checking the flag and actually locking the mutex, where another goroutine might have already run f() and set the flag.
	We also set the `done` flag before calling f(). The idea is that even if f() panics, we still mark it as `success` to prevent it from running again.

	But the last action is a mistake. Imagine the last scenario, we set done to 1, but f() hasn't finished yet, maybe it's stuck on a long network call.

	Now a second goroutine comes along, checks the flag, sees that it's set, and mistakenly thinks, "Great, the resource is ready to go!" But in reality, it's still being fetched. So what happens? Nil dereference and
	panic! The resource isn't ready, and the system tries to use it early.

	We can fix the problem by using `defer` like this:

		func (o *Once) Do(f func()) {
			if o.done.Load() == 1 {
				return
			}

			// slow path
			o.m.Lock()
			defer o.m.Unlock()

			if done.Load() == 0 {
				defer o.done.Store(1)
				f()
			}
		}

	We might think, "Okay, this looks pretty solid now". But it's still not perfect.

	The idea is that Go supports something called inlining optimization.

4. If a function is simple enough, the Go compiler will "inline" it, meaning it'll take the function's code and paste it directly where the function is called, making it faster. Our Do() function is still too complex for
	inlining, though, because it has multiple branches, a defer, and function calls.

	To help the compiler make better decisions about inlining the code, we can move the slow path logic to another function:

		func (o *Once) Do(f func()) {
			if o.done.Load() == 0 {
				o.doSlow()
			}
		}

		func doSlow(o *Once) {
			o.m.Lock()
			defer o.m.Unlock()

			if o.done.Load() == 0 {
				defer o.done.Store(1)
				f()
			}
		}

	This makes the once.Do() function much simpler and can be inlined by the compiler.

	Even though from our point of view, it now has 2 function calls, it's not quite that way in practice. The o.done.Load() is an atomic op that Go's compiler handles in a special way (compiler intrinsic), so it doesn't
	count toward the function call complexity.

5. Why not just inline doSlow()?
	The reason is that after the first call to Do(f), the common scenario is the fast path - just checking if the function has already run.

	In real world apps, after f() runs once (which is the slow path), there are usually many more calls to once.Do(f) that just need to quickly check the `done` flag without locking or re-running the function.

	That's why we optimize for the fast path, where we just check if it's already done and immediately return. And remember when we talked about why the `done` field is placed first in the `Once` struct? That's because
	it makes the fast path quicker by being easier to access.

	Now we have the perfect version of `sync.Once`, but here's the final quiz. The Go team also mentioned an implementation version using a Compare-And-Swap (CAS) op, which makes our Do() function much simpler:

		func (o *Once) Do(f func()) {
			if o.done.CompareAndSwap(0, 1) {
				f()
			}
		}

	The idea is that whichever goroutine can successfully swap the value of `done` from 0 to 1 would "win" the race and run `f()`, while all other goroutines would just return. But why doesn't the Go team use this
	version? Can you guess why before reading the next section, as we already discussed this mistake?

	While the "winning" goroutine is still running `f()`, other goroutines might come and check the `done` flag, thinking `f()` is already finished, and proceed to use resources that aren't fully ready yet.

*/

func impossibilityOnceReusage() {
	var once sync.Once

	once.Do(func() {
		fmt.Println("This will be printed once")
	})

	once.Do(func() {
		fmt.Println("This will be printed once")
	})
}

func errorHandling() {
	var (
		once sync.Once
		err  error
	)

	once.Do(func() {
		if rand.Intn(2) == 1 {
			err = fmt.Errorf("An error occured")
			return
		}
	})

	if err != nil {
		fmt.Println(err)
	}
}

func printOnce() {
	fmt.Println("This will be printed once")
}

func OnceFuncUsage() {
	var wrapper = sync.OnceFunc(printOnce)

	wrapper()
	wrapper()
	wrapper()
}

func OnceValueUsage() {
	var getConfigOnce = sync.OnceValue(func() int {
		fmt.Println("Returning the value -1") // This will be printed Once
		return -1                             // The result of the first call will be cached.
	})

	valA := getConfigOnce() // Returning the value -1
	valB := getConfigOnce() // No print, the cached result will be just returned

	fmt.Println(valA == valB)
}

func OnceValuesUsage() {
	var getIntegerOnce = sync.OnceValues(func() (*int, error) {
		fmt.Println("Generating number!")

		if rand.Intn(2) == 1 {
			return nil, fmt.Errorf("generating a random number an error occureed")
		}

		res := 0
		return &res, nil
	})

	var (
		ptrA *int
		err  error
	)

	ptrA, err = getIntegerOnce() // Generating number!
	ptrA, err = getIntegerOnce() // Just getting the cached results

	fmt.Println(ptrA != nil, err == nil)

}
