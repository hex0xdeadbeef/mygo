package chapter9

/*------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
9.1 RACE CONDITIONS

1. When we cannot confidently say that one event happens befor the ofther, then events "x" and "y" are concurrent

2. The function is concurrently-safe if it continues to work correctly even when called concurrently, that is, from two or more goroutines with no additional synchronization.

3. A type if concurrently-safe if all its accessible methods and operations are concurrently-safe.

4. We should access a variable concurrently only if the documentation for its type says that this is safe.

5. We avoid concurrent access to most variables either by:
	1) confining them to a single goroutine
	OR
	2) maintaining a higher-level invatiant of mutual-exclusion.

6. Exported package-level functions are generally expected to be concurrency-safe. Since package-level variables cannot be confined to a single goroutine, functions that modify them
must enforce mutual exclusion.

7. There are many reasons a function might not work when called concurrently, including "deadlock"/"livelock"/"resorce starvation"/"race condition"

8. CHECK THE CODE OF "Race()"

A race condition is a situation in whic the program doesn't give the correct result for some interleavings of the operations of multiple goroutines.
		Alice first				Bob first			Alice/Bob/Alice
		A1	200					B	100					A1	200
		A2	"=200"				A1	300					B	300
		B	"300"				A2	"= 300"				A2	"= 300"

There is the extra subtle case, when Bob's deposit executes in the middle of A1 operation. Say, A1 is two operations A1Read and A1Write and Bob's addition happens
in the middle of it. This operation falls through.
	Data Race:
A1Read	0	... = balance + amount
B		100
A1Write	200	balance = ...
A2		"= 200"

A data race occurs whenever two goroutines access the same variable concurrently and at least one of the accesses is a write.

9. CHECK THE CODE OF "SliceChanging()".
The value of "x" in the final statement is not defined:
	1) It could be nil
	OR
	2) A slice of length 10
	OR
	3) A slice of length 1,000,000.
But recall that there are three parts to a slice: the pointer, the length and the capacity. If the pointer comes from the first call to make and the length comes from the second, "x" would be a chimera, a slice whose nominal length is 1,000,000 but whose underlying array has only 10 elements. In this eventuality storing to
element 999,999 would clobber an arbitrary faraway memory location with consequences that are impossible to predict and hard to debug and localize. This semantic
minefield is called "undefined behavior"

10. A good rule of thumb: there is no such thing as a benign data race

11. Three ways to avoid a data race:
	1) NOT TO WRITE THE VARIABLE. CHECK THE CODE OF "Icon()".
	While the process a goroutine could keep up with another adding an element, so the second one adds this
	element as well. If instead we initialize the map with all necessary entries before creating additional goroutines and never modify it again, then any number of goroutines may safely call Icon() since each only reads the map. In the example, the icons variable is assigned during package initialization which happens before the program's main function starts running, Once initialized, icons is never modified. Data structures that are never modified or are immutable are inherently concurrency-safe and need no synchronization.

	2) AVOID ACCESSING THE VARIABLE FROM MULTIPLE GOROUTINES. CHECK THE OF "teller()"
	Since other goroutines cannot access the variable directly, they must use a channel to send the confining goroutine a request to query or update the variable. This is meant by the Go mantra: "Don't communicate by sharing memory; instead, share memory by communicating". A goroutine that brokers access to a confined variable using channel requests is called "monitor goroutine" for that variable.

	2!) SERIAL CONFINEMENT. CHECK THE CODE OF "Cake"
	Even when a variable cannot be connfined to a single goroutine for its entire lifetime, confinement may still be a solution to the problem of concurrent access. For example, it's common to share a variable between goroutines in a pipeline by passing its address from one stage to the next over a channel. If each stage of pipeline refrains from accessing the variable after sending it to the next stage, the all accesses to the variable are sequential. In effect, the variable is confined to one stage of the pipeline, then confined to the next and so on. This discipline is sometimes called "serial confinement".

	3) ALLOW MANY GOROUTINES TO ACCESS THE VARIABLE, BUT ONLY ONE AT A TIME. This approach is known as "mutual exclusion"
------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------*/

/*------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
9.2 MUTUAL EXCLUSION: sync.Mutex

1. We have used a buffered channel as a "counting semaphore" to ensure that no more N goroutines made simultaneous HTTP requests. With the same idea, we can use a channel of capacity 1
to ensure that at most one goroutine accesses a shared variable at a time. A semaphore that counts only to 1 1s called a "binary semaphore".

2. This pattern of mutual exclusion is so useful that is supported by the Mutex type from the "sync" package. Its "Lock()" method acquires the token (called a lock) and its "Unlock()" method releases it.
	1) Each time a goroutine accesses the variables of the bank, it must call the mutex's "Lock()" method to acquire an exclusive lock.
	2) If some other goroutine acquired the lock, this operation will be blocked until the other goroutine calls "Unlock()" and the block becomes available again.

3. The mutex guards the shared variables. By convention, the variables guarded by a mutex are declared immediately after the declaration of the mutex itself.

4. The region of code between "Lock()" and "Unlock()" in which a goroutine is free to read and modify the shared variables is called "critical section"

5. A set of exported functions encapsulates one or more variables so that the only way to access the variables is through these functions (or methods, for the variables of an object).Each function acquires a mutex lock at the beginning and releases it at the end, thereby ensuring that the shared variables are not accessed concurrently. This arrangement of functions, mutex lock, and variables is called a "monitor" (broker that ensures variables are accessed sequentially).

6. If the function body is very elaborate we can defer "Unlock()" op.
	1) A deferred "Unlock()" will run even if the critical section panics, which may be important in programs that make use of recover.

7. "atomic" function is a fucntion that executes some operations as a single thing rejecting any terminations and interventions by other goroutines.

8. MUTEXES DEADLOCK. If there's an attempt to lock the mutex again, for example in the internal function, there will be a "deadlock"
	func Withdraw(amount int) bool {
		mu.Lock()
		defer mu.Unlock()
		Deposit(-amount) // INSIDE THE "mu.Lock()" WILL BE CALLED AGAIN
		if Balance() < 0 {
			Deposit(amount)
			return false // insufficient funds
		}
		return true
	}

Deposit tries to acquire the mutex lock a second time by calling mu.Lock() but because mutex locks are not re-entrant - it's not possible to lock a mutex that's already locked. This leads to a deadlock where nothing can proceed, and Withdraw blocks forever.

In order to defeat this problem we should divide the fucntion Withdraw into two variants:
	1) An unexported function, deposit(), that assumes the lock is already held and does the real work
	2) An exported function Deposit() that acquires the lock before calling deposit()
We can then express Withdraw in terms of deposit like this:
	func Withdraw(amount int) bool {
		mu.Lock()
		defer mu.Unlock()
		deposit(-amount)
		if balance < 0 {
			deposit(amount)
			return false
		}
		return true
	}

	func Deposit(amount int) {
		mu.Lock()
		defer mu.Unlock()
		deposit(amount)
	}

	func Balance() int {
		mu.Lock()
		defer mu.Unlock()
		return balance
	}

	func deposit(amount int) {
		balance += amount
	}

9. For the same reason, encapsulation also helps us maintain concurrency invariants. When we use a mutex, we should make sure that both it and the variables it guards are not exported, whether they are package-level variables or the fields of a struct.
------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------*/

/*------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
9.3 READ/WRITE MUTEXES sync.RWMutex
1. Since the "Balance()" function only needs to read the state of the variable, it would in fact safe for multiple Balance calls to run concurrently, so long as no "Deposit()" or
"Withdraw()" call is running. In this scenario we need a special kind of lock that allows read-only multiple operations to proceed in parallel with each other, but write operations to have  fully exclusive access. This lock is called a "multiple readers", "single writer" lock.
	1) "Rlock()" allows multple reading, but the writing data will be allowed only when all the goroutines that read this resource call the "RUnlock()"

2. The "BalanceRWMutex()" function now calls the "RLock()" and "RUnclock()" methods to acquire and release a "readers" or "shared" lock. The Deposit function, which is unchanged, calls
the mu.Lock and mu.Unlock methods to acquire and release a "writer" or "exclusive" lock. With this change, Balance requests run in parallel with each other and finish more quickly.

3. A method that appears to be a simple accessor might also increment an internal usage counter, or update a cachee so that repeat calls are faster.

4. It's profitable to use RWMutex when most of the goroutines that acquire the lock are readers, and the lock is under contention, that is, goroutines routinely have to wait to acquire
it.
------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------*/

/*------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
9.4 MEMORY SYNCHRONIZATION
1.Two reasons we need a mutex in the "Balance()".
	1) It's equaly important that "Balance()" not execute in the middle of some other operation like "Withdraw()".
	2) Synchronization is about more than just order of execution of multiple goroutines; synctonization also affects memory.

2. In a modern computer there may be dozens of processors, each with its own local cache of the main memory. For efficiency, writes to memory are buffered within each processor and
flushed out to main memory only when necessary. They may even be commited to main memory in a different order than they were written by the writing goroutine. Synchronization primitives
like "channel communications" and "mutex operations" cause the processor to flush out and commit all its accumulated writes up the the effects of goroutine execution up to that point are
guaranteed to be visible to goroutines running on other processors.

3. Goroutines are "sequentially consistent", that means:
	1) ORDER. All the operations are executed in the order they were planned. This means if operation "A" were planned before "B", the "A" is always executed before "B"

	2) VISIBILITY. All the changes were made by a single operation get clear and visible to other operations. This means that if the operation "A" changes the value of variable and then
	the operation "B" reads this value, then this operation will see all the changes.

	3) ATOMICITY. Operations that were supposed to be processed as the single thing are executed in the atomic way. This means that if the operation "A" includes some steps within itself, all the steps will be executed as a single thing. And no operations can intervene inside it.

4. CHECK THE CODE OF "DeterministicOutput()"
	Operation sequence			Output
		Sx Oy Sy Ox			y = 0	x = 1
		Sy Ox Sx Oy			x = 0	y = 1
		Sx Sy Ox Oy 		x = 1	y = 1
		Sy Sx Oy Ox 		y = 1 	x = 1

Side effects
	1)
	Operation sequence			Output
		Sx Oy 					y = 0
		Sy Ox					x = 0

	2)
	Operation sequence			Output
		Sy Ox 					x = 0
		Sx Oy					y = 0
Although goroutine A must observe the effect of the write "x" before it reads the value of "y", it doesn't necessarily observe the write of "y" done by goroutine B, so A may print aa stale value
of "y"

5. Because the assignment and the Print refer to different variables, a compiler may conclude that the order of the two statements cannot affect the result and swap them. If the two
goroutines execute on different CPUs, each with its own cache, writes by one goroutine are not visible to the other goroutine's Print() until the caches are synchronized with main memory.
------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------*/

/*------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
9.5 Lazy Initialization: sync.Once
1. It's good practice to defer an expensive initialization step until the moment its needed. Initializing a variable  up front increases the start-up latency of a program and is
unnecessary if execution doesn't always reach the part of the program that uses that variable.

2. CHECK THE CODE OF "RawLazyIcon()".
	1) For a variable accessed by only a single goroutine, we can use the pattern of "loadIconsUnsafe()" but this pattern isn't safe if "RawLazyIcon()" is called concurrently.

	2) Because of a possibility of swap states of main memory between the two goroutines by compiler and CPU at the moment when the first goroutine reached the statements of filling
	the map and the second one reached only the statement of the check whether the map is nil and it really is:
		1) The first goroutine will get the nil map and the program will panic.

		2) The second goroutine will get non-nil map and the check will be passed.

		3) However, the cost of enforcing mutually exclusive access to icons is that two goroutines cannot access the variable concurrently, even once the variable has been safely
		initialized and will never be modified again. This suggests a multiple readers lock.

		4. There is no way to upgrade a "shared lock" to an "exclusive" one without first releasing the "shared lock" so we must recheck the icons variable in case another goroutine
		already initialized it in the interim. This pattern gives us pattern greater concurrency but is complex and error-prone.

3. For an one-time initialization the package "sync" provides a specialized solution: "sync.Once". Conceptually, a "Once()" consists of a mutex and a boolean variable that records
whether initialization has taken place. The mutex guards the boolean and the client's data structures.
	1) The sole method "Do()" accepts the initialization function as its argument.

	2) Each call to "Do()" locks the mutex and checks the boolean variable.
		1) In the first call, in which the variable is false, "Do()" calls "loadIconsUnsafe()" and sets the variable to "true".

		2) Subsequent calls do nothing, but the mutex synchronization ensures that the effects of "loadIconsUnsafe()" on memory (specifically icons) become visible to all goroutines.

		3) Using "sync.Once" in this way, we can avoid sharing variables with other goroutines until they have been properly constructed.

------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------*/

/*------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
9.6 THE RACE DETECTOR
1. Go has the tool is named "Race Detector". This utility checks whether there's a race in our code.

2. To use this tool we just add to the "go filename.go build/run/test" the flag "-race". For example: "go run -race main.go" This causes the compiler to build a modified version of our
application or test with additional instrumentation that effectively records all accesses to shared variables that occured during execution, along with identity of the goroutine that
read or wrote the variable.

3. In addition, the modified program records all synchronization events, such as "go" statements, channel operations, and calls to "(*sync.Mutex).Lock()", "(*sync.WaitGroup).Wait()" and
so on.

4. The "race detector" studies this stream of events, looking for cases in which one goroutine reads or writes a shared variable that was most recently written by a different goroutine
without an intervening synchronization operation. This indicates a concurrent access to the shared variable, and thus a data race. The tool prints a report that includes:
	1)  An identity of the variable

	2) 	Stacks of active function calls in the reading goroutine and the writing goroutine.

This is usually sufficient to pinpoint the problem.

5. Report example:

➜  cmd git:(main) ✗ go run -race main.go
x = 0 ==================
WARNING: DATA RACE
Write at 0x000102b0e098 by goroutine 7:
  golang/pkg/chapters/chapter9.DeterministicOutput.func1()
      /Users/dmitriymamykin/Desktop/goprojects/golang/pkg/chapters/chapter9/9.4 Memory Synchronization.go:17 +0x70

Previous read at 0x000102b0e098 by goroutine 8:
  golang/pkg/chapters/chapter9.DeterministicOutput.func2()
      /Users/dmitriymamykin/Desktop/goprojects/golang/pkg/chapters/chapter9/9.4 Memory Synchronization.go:25 +0x8c

Goroutine 7 (running) created at:
  golang/pkg/chapters/chapter9.DeterministicOutput()
      /Users/dmitriymamykin/Desktop/goprojects/golang/pkg/chapters/chapter9/9.4 Memory Synchronization.go:15 +0x38
  main.main()
      /Users/dmitriymamykin/Desktop/goprojects/golang/cmd/main.go:6 +0x20

Goroutine 8 (running) created at:
  golang/pkg/chapters/chapter9.DeterministicOutput()
      /Users/dmitriymamykin/Desktop/goprojects/golang/pkg/chapters/chapter9/9.4 Memory Synchronization.go:22 +0x54
  main.main()
      /Users/dmitriymamykin/Desktop/goprojects/golang/cmd/main.go:6 +0x20
==================
==================
WARNING: DATA RACE
Read at 0x000102b0e0a0 by goroutine 7:
  golang/pkg/chapters/chapter9.DeterministicOutput.func1()
      /Users/dmitriymamykin/Desktop/goprojects/golang/pkg/chapters/chapter9/9.4 Memory Synchronization.go:18 +0x8c

Previous write at 0x000102b0e0a0 by goroutine 8:
  golang/pkg/chapters/chapter9.DeterministicOutput.func2()
      /Users/dmitriymamykin/Desktop/goprojects/golang/pkg/chapters/chapter9/9.4 Memory Synchronization.go:24 +0x70

Goroutine 7 (running) created at:
  golang/pkg/chapters/chapter9.DeterministicOutput()
      /Users/dmitriymamykin/Desktop/goprojects/golang/pkg/chapters/chapter9/9.4 Memory Synchronization.go:15 +0x38
  main.main()
      /Users/dmitriymamykin/Desktop/goprojects/golang/cmd/main.go:6 +0x20

Goroutine 8 (running) created at:
  golang/pkg/chapters/chapter9.DeterministicOutput()
      /Users/dmitriymamykin/Desktop/goprojects/golang/pkg/chapters/chapter9/9.4 Memory Synchronization.go:22 +0x54
  main.main()
      /Users/dmitriymamykin/Desktop/goprojects/golang/cmd/main.go:6 +0x20
==================
y = 1 Found 2 data race(s)
exit status 66

6. The race detector reports all data races that were acually executed. However, it can only detect race conditions that occur during a run.
	1) it cannot prove that none will ever occur
	2) For best results, make sure that your tests exercise your packages using concurrency.

7. Due to extra bookkeeping a progrsm built with race detection needs more time and memory to run, but the overhead is tolerable. For infrequently occuring race conditions, letting the
race detector do its job can save hours or days of debugging.
------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------*/

/*------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
9.7 CONCURRENT NON-BLOCKING CACHE
1. GetMutex(), GetTwoMutexes() provides us with the safe operations but have side effects. We want to have the cache that will be written once and accessible to read ftom after write
operation. So that we can create it we need implement "duplicate suppression"

2. In the fourth version of Memo:
	1) Each map element is a pointer to an "entry" struct.

	2) Each "entry" contains the memoized result of a call to the function "f()", as before, but it additionally contains a channel called "ready". Just after the entry's "result" has
	been set, this channel will be closed.

	3) We should notice that the variables "e.res.value" and "e.res.err" in the entry are shared among multiple goroutines. The goroutine that creates the entry sets their values, and
	other goroutines read their values once the "ready" condition has been broadcast. Despite being accessed by multiple goroutines, no mutex lock is necessary. The closing of the
	"ready" channel happens before any other goroutine receives the broadcast, so the write to those variables in the first goroutine happens before they are read by subsequent
	goroutines. There's no data race.
------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------*/

/*------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
9.8 GOROUTINES AND THREADS


9.8.1 GROWABLE STACKS
THREADS:
1) Each OS thread has a fixed-size block of memory (often as large as 2MB) for its "stack". Stack is the work area where it saves the local variables of function calls that are in
progress or temporarily suspended while another function is called. This fixed-size stack is simultaneously too much and too little.

2) Changing the fixed size can improve space efficiency and allow more threads to be created, or it can more deeply recursive functions, but it cannot do both.

GOROUTINES
1) In contrast, a goroutine starts life with a small stack, typically 2KB. A goroutine's stack, like the stack of an OS thread, holds the local variables of active and suspended
function calls, but unlike an OS thread, a goroutine's stack is not fixed. It grows and shrinks as needed. The size limit for a goroutine stack may be as much as 1GB, orders of
magnitude larger than a typical fixed-size thread stack, though of course few goroutines use that much.


9.8.2 GOROUTINE SCHEDULING
THREADS
1) OS Threads are scheduled by the OS kernel. Every few milliseconds, a hardware timer interrupts the processor, which causes a kernel function called the "scheduler" to be
invoked. This function suspends the currently executing thread and saves its registers in memory, looks over the list of threads and decides which one should run next, restores
that thread's registers from memory, then resumes the execution of that thread. Because OS threads are scheduled by the kernel, passing control from one thread to another requires
a full "context switch", that is, saving the state of one user thread to memory, restoring the state of another, and updating the scheduler's data structures. This operation is slow, due
to its poor locality and the number of memory accesses required, and has historically only gotten worse as the number of CPU cycles required to access memory has increased.

GOROUTINES
1) The Go runtime contains its own scheduler that uses a technique known as "M:N scheduling", because it multiplexes (or schedules) M goroutines on N OS threads. The job of the Go scheduler
is analogous to that of the kernel scheduler, but it's concerned only with the goroutines of a single Go program.

2) Unlike the operating system's thread scheduler, the Go scheduler is not invoked periodically by a hardware timer, but implicitly by certain Go language constructs. For example, when
a goroutine calls "time.Sleep()" or blocks in a channel or mutex operation, the scheduler puts it to sleep and runs another goroutine until it's time to wake the first one up. Because it
doesn't need a switch to kernel context, rescheduling a goroutine is much cheaper than rescheduling a thread.


9.8.3 GOMAXPROCS
1) The Go scheduler uses a parameter called "GOMAXPROCS" to determine how many OS threads may be actively executing Go code simultaneously. Its default value is the number of CPUs on the
machine, so on a machine with 8 CPUs, the scheduler will schedule Go Code on up to 8 OS threads at once (GOMAXPROCS is the N in "M:N scheduling")

2) Operations that need or don't an OS thread
	1) Goroutines that are sleeping or blocked in a communication don't need a thread at all.
	2) Goroutines that are blocked in I/O or other system calls or are calling non-Go functions, do need an OS thread.
But "GOMAXPROCS" needn't account for them.

3) We can explicitly control this parameter using the "GOMAXPROCS" environment variable or the "runtime.GOMAXPROCS" function.


9.8.3 GOROUTINES HAVE NO IDENTITY
THREADS
1) In most operating systems and programming languages that support multithreading, the current thread has a distinct "identity" that can be easily obtained as an ordinary value, typically an
integer or pointer. This make it easy to buold an abstraction called "thread-local storage", which is essentially a global map keyed by thread identity, so that each thread can store and retrieve
values independent of other threads.

GOROUTINES
1) Goroutines have no notion of "identity" that is accessible to the programmer.

2) Go encourages a simpler style of programming in which parameters that affect the behavior of a function are explicit. Not only does this make programs easier to read, but it
lets us freely assign subtasks of a given function to many different goroutines without worrying about their identity.

------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------*/
