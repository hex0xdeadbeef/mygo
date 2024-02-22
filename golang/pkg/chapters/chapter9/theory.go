package chapter9

/*------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
9.1 RACE CONDITIONS

1. When we cannot confidently say that one event happens befor the ofther, then events "x" and "y" are concurrent

2. The function is concurrently-safe if it continues to work correctly even when called concurrently, that is, from two or more goroutines with no additional synchronization.

3. A type is concurrently-safe if all its accessible methods and operations are concurrently-safe.

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

------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------*/
