package main

/*
	CHANNELS

	ONCE
1. Once is a fairly simple primitive: the first call to Do(func()) will cause all other concurrent calls to block until the argument of Do() returns. After this happens all blocked calls and successive ones will do nothing and return immediately.
2. This is useful for lazy initialization and singleton instances.

	MUTEX
1. Two little digressions will be made for mutex:
	- A Semaphore of size N allows at most N goroutines to hold its lock at
	any given time. Mutexes are special cases of Semaphores of size 1.
	- Mutexes might benefit from a TryLock method.
2. The "sync" package doesn't provide semaphores or triable locks but since
	we're trying to prove expressive power of channels and select let's
	implement both

	RWMUTEX
1. RWMutex is a slightly more complicated primitive: it allows an arbitrary
	amount of concurrent read locks but only a single write lock at any
	given time. It also guaranteed that if someone is holding a write lock
	no one should be able to hace or acquire a read lock
2. The standart library implementation also grants that if a write lock is
	attempted, further read locks will queue up and wait to avoid starving
	write lock. We're going to relax this condition for brevity.
3. There's a way to do this by using a 1-sized channels mutex to protect
	the RWLock state, but it's a boring example to show. This
	implementation will starve writers if there is always at least one
	reader.
4. RWMutex has three states: free, with a writer and readers. This means we
	need two channels:
		- when the mutex is free, we'll have both channels empty
		- when there's a writer we'll have a channel with an empty struct
		in it;
		- when there are readers we'll have a value in both channels, one
		of which will be the count of readers.
5.  TryLock and TryRLock methods could be implemented as I did in the
	previous example by adding default cases to channel ops

	POOL
1. Pool is used to relieve strees on the GC and re-use objects that are
	frequently allocated and destroyed
2. The standart library uses many techniques for this:
	- thread local storage
	- removing oversized objects and generations to make objects to survive
	in the pool only if they are used between two GC collections.
3. We're not going to provide all of these functionalities: we don't have
	access to thread local storage nor garbage collection informations, nor
	we care.

	Our pool will be fixed size as we just want to express the
	semantic of the type, not the implementation details
4. One utility that we're going to provide that Pool doesn't have is a
	cleanup function.
5. It's very common whem using Pool to not clear up the objects that come
	out of it, thus causing re-use of non-zeroed memory, which can lead to
	nasty bugs and vulnerabilities. Our implementation is going to
	guarantee that a cleaner function is called if and only ifthe returned
	object is recycled.
6. If a cleanup isn't specified we won't call it

	WAITGROUP
1. WaitGroups allow for a variety of uses but the must common one is to
	create a group, Add a count to it, spawn as many goroutines as that
	count and then wait for all of them to complete.
2. Every time a goroutine is done running it'll call Done on the group to
	signal it has finished working.
3. WaitGroup can reach "0" count by calling Done or calling Add with a
	negative count (even greater than -1).
4. When a group reaches 0 all the waiters currently blocked on a Wait call
	will resume execution. Future calls to Wait won't be blocking.
5. One little known feature is that wait groups can be re-used: Add can
	still be called after the counter reached 0 and it'll put the wait
	group back in the blocking state.
6. This means we have a sort of "generation" for each given wait group:
	- a generation begins when the counter moves from 0 to a positive number
	- a generation ends when the counter reaches 0
	- when a generation ends all waiters of that generation are unblocked
*/

/*
ONCE
*/
type Once chan struct{}

func NewOnce() Once {
	o := make(Once, 1)
	// Filling buffer of Once so that it'll be full
	o <- struct{}{}

	return o
}

func (o Once) Do(f func()) {
	// Read from a closed chan always succeeds
	// This only blocks during initialization.
	if _, ok := <-o; !ok {
		return
	}

	// Call f.
	// Only one goroutine will get here as there's only one value in the
	// channel

	// Run the function of the initializing instance
	f()

	// This unreleases all waiting goroutines and future ones
	close(o)
}

/*
MUTEX
*/
type Semaphore chan struct{}

func NewSemaphore(size int) Semaphore {
	return make(Semaphore, size)
}

// Lock acquires a free place in the buffered channel if any
func (s Semaphore) Lock() {
	// Writes will only succeed if there's a room in s
	s <- struct{}{}
}

// TryLock tries to make a write into a buffered channel if it's possible
// and returns true if the a write has succeeded, returns false otherwise
func (s Semaphore) TryLock() bool {
	// Select with default case: if no cases are ready just fall in the
	// default block
	select {
	case s <- struct{}{}:
		return true
	default:
		return false
	}
}

// Unlock frees a room of the buffered channel if any
func (s Semaphore) Unlock() {
	// Make room for other users of the semaphore
	<-s
}

// The resulting type "Mutex":
type Mutex Semaphore

func NewMutex() Mutex {
	return Mutex(NewSemaphore(1))
}

/*
	RWMUTEX
*/

type RWMutex struct {
	write   chan struct{}
	readers chan int
}

func NewRWMutex() RWMutex {
	return RWMutex{
		// This is used as a normal Mutex
		write: make(chan struct{}, 1),
		/*
			This is used to protect the readers count.
			By receiving the value it's guaranteed that no
			other goroutine is changing it at the same time
		*/
		readers: make(chan int, 1),
	}
}

// Lock locks the mutex by sending to the buffered channel the value,
// so every goroutine trying to lock the write mutex will have to wait of
// freeing of this buffered channel
func (rwm RWMutex) Lock() { rwm.write <- struct{}{} }

// Unlock unlocks the mutex by freeing the buffered write channel up
func (rwm RWMutex) Unlock() { <-rwm.write }

// RLock waits for freedom of the write mutex and locks the mutex
// for some reads or just increments the number of readers and sends
// the incremented count to the channel of readers count
func (rwm RWMutex) RLock() {
	var (
		// Count current readers. Default to 0.
		rs int
	)

	/*
		Select on the channels without default.
		One and only one case will be selected and this
		will block until one case becomes available
	*/
	select {

	case rwm.write <- struct{}{}:
		/*
			One sending case for write.
			If the write lock is available we have no readers.
			We grab the write lock to prevent concurrent read-writers
		*/
	case rs = <-rwm.readers:
		/*
			One receiving case for read.
			There are already are readers, let's grab and update the
			readers count
		*/
	}

	// If we grabbed a write lock this is 0.
	rs++

	// Updated the readers count. If there are none
	// just adds an item to the empty readers channel
	rwm.readers <- rs
}

func (rwm RWMutex) RUnlock() {
	// Take the value of readers and decrement it
	rs := <-rwm.readers
	rs--

	// If zero, make the write lock available again and return
	if rs == 0 {
		<-rwm.write
		return
	}

	/*
		If not zero, just update the readers count.
		0 will never be written to the readers channel,
		at most one of the two channels will have a value
		at any given time
	*/
	rwm.readers <- rs
}

func (rwm RWMutex) TryLock() bool {
	select {
	case rwm.write <- struct{}{}:
		return true
	default:
		return false
	}
}

func (rwm RWMutex) TryRLock() bool {
	var (
		rs int
	)

	select {
	/*
		We'll be satisfied either
			- the write will pass (in this case the rs will be 0 and it'll
				be incremented and passed to rwm.readers)
			- the rs will be given by rwm.readers, incremented by the
				instructions lower and passed to the readers channel

		In the other cases we'll return false
	*/
	case rwm.write <- struct{}{}:
	case rs = <-rwm.readers:
	default:
		return false
	}

	// In any case of achieving this code we increment the count of readers
	// and pass it to the readers channel
	rs++
	rwm.readers <- rs

	return true
}

/*
	POOL
*/

type Item = interface{}

type Pool struct {
	buf   chan Item   // object buffer
	alloc func() Item // function to allocate new objects when it's
	// needed
	clean func(Item) Item // function that cleans returned objects
}

// NewPool returns a new pool created by using size, alloc, clean parameters
func NewPool(size int, alloc func() Item, clean func(Item) Item) *Pool {
	return &Pool{
		buf:   make(chan interface{}, size),
		alloc: alloc,
		clean: clean,
	}
}

// Get returns [cleaned] object from the pool's buffer
func (p *Pool) Get() Item {
	select {
	case i := <-p.buf:
		// If our pool has the clean function,
		// a cleaned object will be returned
		if p.clean != nil {
			return p.clean(i)
		}

		// Otherwise a used before item will be returned
		return i
	default:
		return p.alloc()
	}
}

// Put puts a given object to the pool's buffer or if it's full sends the
// item to be collected
func (p *Pool) Put(x Item) {
	select {
	case p.buf <- x:
	default:
		// If our buffer is full, the object will be just garbage collected
	}
}

/*
	WAITGROUP
*/

type generation struct {
	// A barrier for waiters to wait on.
	// This will never be used for sending, only receive and close
	wait chan struct{}

	// The counter for remaining jobs to wait for
	n int
}

// newGeneration returns a new generation created with newly created
// channel and waiters count set to 0
func newGeneration() generation {
	return generation{wait: make(chan struct{}), n: 0}
}

// end unlocks a WaitGroup by closing the channel
func (g generation) end() {
	// The end of a generation is signalled by closing its channel
	close(g.wait)
}

// Here we use a channel to protect the current generation.
// This is basically a mutex for the state of the WaitGroup
type WaitGroup chan generation

func NewWaitGroup() WaitGroup {
	wg := make(WaitGroup, 1)
	g := newGeneration()

	// On a new waitgroup waits should just return,
	// so it behaves exactly as after a terminated generation.
	g.end()
	wg <- g
	return wg
}

func (wg WaitGroup) Add(delta int) {
	// Acquire the current generation
	g := <-wg

	if g.n == 0 {
		// We were at 0, create the next generation
		g = newGeneration()
	}

	// Adding up the delta value
	g.n += delta

	// The negative n of the generation's n field
	if g.n < 0 {
		// This is the same behavior of stdlib
		panic("negative WaitGroup count")
	}

	if g.n == 0 {
		// We reached zero, signal waiters to return from Wait
		g.end()
	}

	// If we just add the positive value and the g.n is >0,
	// put the current generation back
	wg <- g
}

func (wg WaitGroup) Done() {
	wg.Add(-1)
}

func (wg WaitGroup) Wait() {
	// Acquire the current generation
	g := <-wg

	// Save a reference to the current waiting chan
	wait := g.wait

	// Release the current generation
	wg <- g

	// Wait for the chan to be closed
	// We'll go through this line after the g.end() will be called
	<-wait
}
