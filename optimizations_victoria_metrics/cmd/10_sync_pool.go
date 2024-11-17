package main

import "sync"

/*
   GO sync.Pool AND THE MECHANICS BEHIND IT
*/

/*
 1. sync.Pool is great for how we handle temporary objects, especially byte buffers or slices. It's commonly used in the standart library. For instance in the encoding/json package:
      package json

      var encodeStatePool sync.Pool

      // An encodeState encodes JSON into a bytes.Buffer
      type encodeState struct {
         bytes.Buffer // accumulated output

         ptrLevel uint
         ptrSeen map[any]struct{}
      }

      In this case, sync.Pool is being used to reuse *encodeState objects, which handle the process of encoding JSON into a bytes.Buffer.

      Instead of just throwing these objects after each use, which would only give the GC more work, we stash them in a pool (sync.Pool). The next time we need something similar, we
      just grab it from the pool instead of making a new one from scratch.

      We'll also find multiple sync.Pool instances in the `net/http` package, that are used to optimize I/O ops:
         package http

         var (
            bufioReaderPool 	   sync.Pool
            bufioWriter2kPool 	sync.Pool
            bufioWriter4kPool 	sync.Pool
         )

      When the server reads request bodies or write responses, it can quickly pull a pre-allocated reader or writer from these pools, skipping extra allocations. Furthermore, the 2
      writer pools, *bufioWriter2kPool and *bufioWriter4kPool, are set up to handle different writing needs.

      func bufioWriterPool(size int) *sync.Pool {
         switch size {
         case 2 << 10:
            return &bufioWriter2kPool
         case 4 << 10:
            return &bufioWriter4kPool
         }

         return nil
      }
*/

/*
	WHAT IS sync.Pool?

1. To put it simply sync.Pool in Go is a place where we can keep temporary objects for later reuse.
2. We don't control how many objects stay in the pool, and anything we put there can be removed at any time, without any warning.
3. The pool is built to be thread-safe, so multiple goroutines can tap into it simultaneously.
*/

/*
		WHY DO WE BOTHER REUSING OBJECTS?

1. When we have a lot of goroutines running at once, they often need similar objects. Imagine running go f() multiple times concurrently. If each goroutine creates its own objects,
   memory usage can quickly increase and this puts a strain on the GC because it has to clean up all those objects once they're no longer needed.

   This situation creates a cycle where high concurrency leads to high memory usage, which then slows down the GC. sync.Pool is designed to help break this cycle.

2. To create a pool, we can provide a New() function that returns a new object when the pool is empty.
  - This function is optional, if we don't provide it, the pool just returns nil. If it's empty.

3. For instance, if the slice grows to 8192 bytes during use, we can reset its length to zero before putting it back in the pool. The underlying array still has a capacity of 8192, so
the next time we need it, those 8192 bytes are ready to be reused.

4. The flow is pretty simple:
We get an object from the pool, use it, reset it, and then put it back into the pool. Resetting the object can be done either before we put it back or right after we get it from the
pool, but it's not mandatory, it's a common practice.

5. If we're not fans of using type assertions pool.Get().(*Object), there are a couple of ways to avoid it:

   1. Use a dedicated function to get the object from the pool:
      func() getObjectFromPool() *Object {
         return pool.Get().(*Object)
      }

   2. Create our own generic version of sync.Pool
      type Pool[T any] struct {
         sync.Pool
      }

      func NewPool[T any](newF func() T) *Pool[T] {
         return &Pool[T]{
            New: func() interface{} {
               return newF()
            }
         }
      }

      func (p *Pool[T]) Get() T {
         return p.Pool.Get().(T)
      }

      func (p *Pool) Put(x T) {
         p.Pool.Put(x)
      }

      The generic wrapper gives us a more type-safe way to work with the pool, avoiding type assertions.

      Just note that, it adds a tiny bit of overhead due to the extra layer of indirection. In most cases, this overhead is minimal, but if we're in a highly CPU-sensitive environment,
      it's a good idea to run benchmarks to see if it's worth it.
*/

/*
	sync.Pool AND ALLOCATION TRAP

1. If we've noticed, from many previous examples, including those in the standart library, what we store in the pool is typically not the object itself but a pointer to the object.

2. For example:
   var pool = sync.Pool{
      New: func() any {
         return []byte{}
      }
   }

   func main() {
      bytes := pool.Get().([]byte)

      // do something with bytes
      _ = bytes

      pool.Put(bytes)
   }

   We're using a pool of `[]byte`. Generally (though not always), when we pass a value to an interface, it may cause the value to be placed on the heap. This happens here too, not just
   with slices but with anything we pass to pool.Put() that isn't a pointer.

   Now, we don't say our variable `bytes` moves to the heap, we say `the value of bytes escapes to the heap`

3. To really get why this happens, we'd need to dig into how escape analysis work (which we might do in another article). However, if we pass a pointer to pool.Put(...), there's no extra
    allocation.

    var pool = sync.Pool{
      New: func() any {
         return new([]byte)
      }
   }

   func main() {
      bytes = *pool.Get().(*[]byte)

      // do something with bytes
      _ = bytes

      pool.Put(bytes)
   }

   If we run the escape analysis again, we won't see that the value escapes to the heap. (https://github.com/golang/go/blob/2580d0e08d5e9f979b943758d3c49877fb2324cb/src/sync/
   example_pool_test.go#L15)

4. Documentation says:
   The Pool's New() function should generally only return a pointer
   types, since a pointer can be put into the return interface value
   without an allocation.
*/

/*
	PMG NOTES

1. If we've got N logical processors (P), we can run up to N goroutines in parallel, as long as we've got at least N machine threads (M) available.

2. At any one time, only one goroutine (G) can run on a single processor (P). So, when a P1 is busy with a G, no other G can run on that P1 until the current G either gets blocked,
   finishes up, or something else happens to free it up.
*/

/*
	   sync.Pool INTERNALS

1. The thing is, a `sync.Pool` in Go isn't just one big pool, its actually made up of several `local` pools, with each one tied to a specific processor context, or P, that Go's runtime
   is managing at any given time.

2. The number of `local` pools depends on the number of runtime.GOMAXPROCS()

3. When a goroutine running on a processor (P) needs an object from the pool, it'll first check its own `P-local pool` before looking anywhere else.

4. This is a smart design choice because it means that each logical processor (P) has it's own set of objects to work with.
   - This reduces contention between goroutines since only one goroutine can access its `P-local pool` at a time.

5. So, the process is super fast because there's no chance of two goroutines trying to grab the same object from the same local pool.
*/

/*
	POOL LOCAL & FALSE SHARING SHARING PROBLEM

1. Earlier, we mentioned that `Only one goroutine can access the P-local pool at the same time`, but the reality is a bit more nuanced.

2. Each `P-local pool` actually has two main parts:
   - Private object `private`
   - The shared pool chain `shared`
   - Cache padding

   Here's the definition of local pool in Go source code:
      type poolLocal struct {
         poolLocalInternal
         pad [128 - unsafe.Sizeof(poolLocalInternal{})%128]byte
      }

      type poolLocalInternal struct {
         private any
         shared poolChain
      }

	The `private` field is where a single object is stored, and only the `P` that owns this `P-local pool` can access it, let's call it the` private object`.

	It's designed so a goroutine can quickly grab a reusable object (which is private object) without needing to mess with any mutex or synchronization tricks. In other words, only one
	goroutine can access its own private object, no other goroutine can compete with it.

	But if the private object isn't available, that's when the shared pool chain `shared` steps in.

3. Why wouldn't private object be available? I thought only one goroutine could get and put the object the private object back into the pool. So, who's the competition?
   While it's true that only one goroutine can access the private of a P at a time, there's a catch. If G `A` grabs the private object and then gets blocked or preempted, G `B`
   might start running on that same P. When that happens, G `B` won't be able to access the private object because G `A` still has it. Now, unlike the simple of private object, the
   `shared` pool chain is a bit more complex.

   The diagram 3 isn't entirely accurate since it doesn't account for the `victim pool`.

   If the `shared pool chain` is empty as well, sync.Pool will either create a new object (assuming we've provided a New() function) or just return `nil`.

   And, by the way, there's also a `victim mechanism` inside the shared pool.

4. Wait, I see the pad field in the `P-local pool`. What's that all about?
   One thing that jumps out when we look at the `P-local` pool structure is this pad attribute. It's something that VM' CTO adjusted in his commit.

   type poolLocal struct {
      poolLocalInternal
      pad [128 - unsafe.Sizeof(poolLocalInternal{})%128]byte
   }

   This pad might seem a bit odd since it doesn't add any direct functionality, but it's actually there to prevent a problem that can cron up on modern multi-core processors called
   `false sharing`

   To get why this matters, we need to dig into how CPUs handle memory and caching. So, brace yourself for a bit of a deep dive into the inner working of CPUs.

   Modern CPUs use a component called the CPU cache to speed up RAM access, this cache is divided into units called cache lines, which typically hold either 64 or 128 bytes of data. When
   the CPU needs to access memory, it doesn't just grab a single byte or word - it loads an entire cache line.

   This means that if two pieces of data are close together in memory, they might end up on the same cache line, even if they're logically separate.

   Now, in the context of Go's `sync.Pool`, each logical processor (P) has its own `poolLocal`, which is stored in an array, if the `poolLocal` structure is smaller than the size of a
   cache line, multiple `poolLocal` instances from different Ps can end up on the same cache line. This is where things can go sideways. If two Ps, running on different CPU cores, try to
   acces their own `poolLocal` at the same time, they could unintentionally step on each other's toes.

   Even though each P is only dealing with its own `poolLocal`, these structures might share the same cache line.

   When one processor modifies something in a chache line, `the entire cache line is invalidated in other processor' caches`, even if they're working with different variables within that
   same line. This can cause a serious performance hit of unnecessary cache invalidations and extra memory traffic.

   That's where `128 - unsafe.Sizeof(poolLocalInternal{}%128)` comes in. It calculates the number of bytes needed to pad the `P-local pool` so that its total size is a multiple of 128
   bytes. This padding helps each `poolLocal` gets its own cache line, preventing `false sharing` and keeping things running faster, free-conflict.
*/

/*
   POOL CHAIN & POOL DEQUEUE
1. The shared `pool chain` in sync.Cond is represented by a type called `poolChain`.
   From the name, we might guess it's a double-linked list. But here's the twist:
      - Each node in this list isn't just a reusable object.
      - Instead, it's another structure called a `pool dequeue` (`poolDequeue`)

   `shared` is:
   type poolChain struct {
      head *poolChainElt
      tail atomic.Pointer[poolChainElt]
   }

   type poolChainElt struct {
      poolDequeue
      next, prev atomic.Pointer[poolChainElt]
   }

   The design of the poolChain is pretty strategic, the diagram 4 show us something.

   When the current pool dequeue (the one at the head of the list) gets full, a new `pool dequeue` is created that's double the size of the previous one. This new, larger pool is then
   added to the chain.

   If we take a look at the poolChain we'll notice it has two fields:
         1) `head` *poolChainElt
         2) An atomic pointer `tail atomic.Pointer[poolChainElt]`
   These fields reveal how the mechanism works:
         1) The producer (the P that own the current `P-local pool`) only adds new items to the most recent pool dequeue, which we call the `head`. Since only the producer is touching
         the `head`, there's no need for locks or any fancy synchronization tricks, so it's really fast.
         2) Consumers (other Ps) take items from the pool dequeue at the `tail` of the list. Since multiple consumers might try to pop items at the same time, access to the tail is
         synchronized using atomic ops to keep things in order.
   But there's the key part:
      1) When the pool dequeue at the tail is completely emptied, it gets removed from the list, and the next pool dequeue in line becomes the new tail. But the situation at the head is
      a bit different.
      2) When the pool dequeue at the head runs out of items, it doesn't get removed. Instead, it stays in place, ready to be refilled when new items are added.

   Now let's take a look at how the pool dequeue is defined. As the name "dequeue" suggests, it's a double-ended queue. Unlike a regular queue where we can only add elements at the back
   and remove them from the front, a dequeue lets us insert and delete elements at both the front and the back. Its mechanism is actually quite similar to the pool chain. It's designed
   so that one producer can add or remove items from the head, while multiple consumers can take items from the tail.

   type poolDequeue struct {
      headTail atomic.Uint64
      vals     []eface
   }

   The producer (which is the current P) can add new items to the front of the queue or take items from it. Meanwhile, the consumers only take items from the tail of the queue. This
   queue is lock-free, which means it doesn't use locks to manage the coordination between the producer and consumers, only using atomic ops.

   We can think of this queue as a kind of ring buffer. A ring buffer, or circular buffer, is a data structure that uses a fixed-size array to store elems in a loop-like fashion. It's
   called a "ring" because, in a way, the end of the buffer wraps around to meet the beginning, making it look like a circle.

   In the context of the pool dequeue we're talking about the `headTail` field is a 64-bit integer that packs two 32-bit indexes into a single value. These indexes are the `head` and
   `tail` of the queue and help keep track of where data is being stored and accessed in the buffer:
      1) The `tail index` points to where the oldest item in the buffer is and when consumers (like other goroutines) read from the buffer, they start here and move forward.
      2) The `head index` is where the next piece of data will be written. As new data comes in, it's placed at this head index, and then the index moves to the next available slot.

   By packing the `head` and `tail` indices into a single 64-bit value, the code can update both indices in one go, making the operation atomic. This is especially useful when two
   consumers (or a consumer and a producer) try to pop an item from the queue at the same time. The CAS op, d.headTail.CompareAndSwap(ptrs, ptrs2), ensures that only one of them succeeds.
   The other one fails and retries, keeping things orderly without any complex locking.

   The actual data in the queue is stored in a circular buffer called `vals`, `which has to be a power of two in size`.

   This design choice makes it easier to handle the queue wrapping around when it reaches the end of the buffer. Each slot in this buffer is an `eface` value, which is how Go represents
   an `empty` interface (`interface{}`) under the hood.
      type eface struct {
         typ, val unsafe.Pointer
      }
   A slot in the buffer stays `in use` until two things happen:
      1) The tail index moves past the slot, meaning the data in that slot has been consumed by one of the consumers.
      2) The consumer who accessed that slot sets it to nil, signaling that the producer can now use that slot to store new data.
   In short, the pool chain combines a linked list and a ring buffer for each node. When one dequeue fills up, a new, larger one is created and linked to the head of the chain. This
   setup helps manage a high volume objects efficiently.
*/

/*
	Pool.Put()

1. Let's start with the Put(...) flow because it's a bit more straightforward than Get(), plus it ties into another process:

  - pinning a goroutine to a P.

    When a goroutine calls Put() on a sync.Pool, the first thing it tries to do is to store the object in the `private` spot of the `P-local pool` for the current P. If that private spot is already occupied, that object gets pushed to the `head` of the pool chain, wbich is the shared part.

    func (p *Pool) Put(x interface{}) {
    // if the object is nil, just return
    if x == nil {
    return
    }

    // Pin the current P's P-local pool
    l, _ := p.pin()

    // If the private pool is not there, create it and set the object to it
    if l.private == nil {
    l.private = x
    x = nil
    }

    // If the private object is there, push it to the head of the shared chain
    if x != nil {
    l.shared.pushHead(x)
    }

    // Unpin the current P
    runtime_procUnpin()
    }

    We haven't talked about the pin() or runtime_procUnpin() functions yet, but they're important for both Get() and Put() ops because they ensure the goroutines stays `pinned` to the current `P`. It means:
    Starting with Go 1.14, Go introduced preemptive scheduling, which means the runtime can pause a goroutine if it's been running on a processor P for too long, usually around 10ms, to give other goroutines a chance to run.

    This is generally good for keeping things fair and responsive,
    but it can cause issues when dealing with sync.Pool.
    Operations like Put() and Get() in sync.Pool assume that the goroutine stays on the same processor (say, P1) throughout the entire op. If the goroutine is preempted in the middle of these ops and then resumed on a different processor (P2), the local data it was working with could end up being from the wrong processor.

    So, what doest the pin() function do?
    // pin pins the current G to P, disables preemption and
    // returns poolLocal pool for the P and the P's id.
    // Caller must call runtime_procUnpin() when done with the pool.
    func (p *Pool) pin() (*poolLocal, int) {...}
    Basically, pin() temporarily disables the scheduler's ability to preempt the goroutine while it's putting an object into the pool.

    Even though it says `pins the current goroutine to P`, what's actually happening is that the current thread M is locked to the processor P, which prevents it from being preempted. As a result, the G running on that thread will also not be preempted.

    As a side effect, pin() also updates the number of processors (Ps) if we happen to change GOMAXPROCS(n) (which controls the number of Ps) at runtime.

    How about the shared pool chain?
    When we need to add an item to the chain, the op first checks the head of the chain, Remember the *poolChainElt pointer? That's the most recent pool dequeue in the list.
    Depending on the situation, here's what can happen:
    1) If the head buffer of the chain is `nil`, meaning there's no pool dequeue in the chain yet, a new pool dequeue is created with an initial buffer size of `8`. The item is then placed into this brand-new pool dequeue.
    2) If the head buffer of the chain isn't `nil` and that buffer isn't full, the item is simply added to the buffer at the `head` position.
    3) If the head buffer of the chain isn't `nil`, but that buffer is full, meaning that the head index has wrapped around and caught up with the tail index, then a new pool dequeue is created. This new pool has a buffer size that's double of the current head. The item is placed into this new pool dequeue, and the head of the pool chain is updated to point to this new pool.

    And that's pretty much it for the Put() flow. It's a relatively simple process because it doesn't involve interacting with the local pool other processors (Ps); everything happens within the current head of the pool chain.

    sync.Pool().Get()

1. At first glance, the Get() function seems pretty similar to Put().

2. It starts by pinning the current goroutine to its P to prevent preemption, then checks and grabs the private object from its `P-local pool` without needing any synchronization. If the private object isn't there, it checks the shared pool chain and pops the head of the chain.

		Only the goroutine running on the current `P-local pool` can access the head of the chain, which is why we use popHead():
		   func (p *Pool) Get() interface{} {
		      // Pin the current P's P-local pool
		      l, pid := p.pin()

		      // Get the private obj from the current P-local pool
		      x := l.private
		      l.private = nil

		      // If the private obj isn't there, pop the head of the shared pool chain
		      if x == nil {
		         x, _ = l.shared.popHead()

		         // Steal from other P's cache
		         if x == nil {
		            x = p.getSlow(pid)
		         }
		      }
		      runtime_procUnpin()

		      // If the object is still not there, create a new object by using the factory function
		      if x == nil && p.New != nil {
		         x = p.New()
		      }

		      return x
		   }

		   Unlike in p.pin() for Put(), here we also get the `pid`, which is the ID of the P that the
		   current goroutine is running on. We need this for the stealing process, which comes into play
		   if the fast path fails.

		   The fast path is when the object is available in the current P's cache. But if that doesn't work out, meaning the private object and the head of the shared chain are both empty, the slow path (getSlow) takes over.

		   In the slow path, we try to steal objects from the cache pools of other processors (Ps).

		   The idea behind stealing is to reuse objects that might be sitting idle in the caches of other processors, instead of creating new objects from scratch. If another P has extra objects in its cache pool, the current P can grab those objects and put them to use.

		   The stealing process basically loops through all the Ps, except the current one (pid) and tries to grab an object from each P's shared pool chain.
		      for i := 0; i < int(size); i++ {
		         l := indexLocal(locals, (pid+i+1)%int(size))
		         if x, _ := l.shared.popTail(); x != nil {
		            return x
		         }
		      }
		   As we've talked about before, in a poolChain, the provider (the current P) pushes and pops at the head, while multiple consumers (other Ps) pop from the tail. So popTail looks at the last pool dequeue in the linked list and tries to grab data from the end of that dequeue:
		      1) If it finds data, the steal is successful, and the data is returned.
		      2) If it doesn't find any data in that pool dequeue, the tail index increases, and that pool dequeue gets removed from the chain.
		   This process continues until it either successfully steals some data or runs out of options in all the pool chains.

		   If, after all the stealing attempts, it still can't find any data, the function tries to get data from what's called the `victim`. This is a new concept related to how sync.Pool cleans up objects.

		   We're trying to grab an object in every possible way, and if nothing is found, it finally creates a new object using New(). But if New() is nil, it just returns `nil`. Simple as that.

		   Now, after the attempt with the victim pool, it's atomically marked as empty (though concurrent accesses may still retrieve from it). Subsequent Get() ops will skip checking the victim cache until it's filled up again.

		   VICTIM POOL

	 1. Even though `sync.Pool` is built for better manage resources, it doesn't give us, developers, direct tools to clean up or manage object lifecycles. Instead, `sync.Pool` handles
	    cleanup behind the scenes to avoid uncheked growth, which could lead to memory leaks.

	    The primary way this cleanup happens is through Go's GC.

	    Remember when we talked about pin()? It turns out pin() has another side effect.
	    Every time a sync.Pool calls pin() for the first time (or after the number of Ps has changed via GOMAXPROCS), it gets added to a global slice called allPools in the sync package.

	    package sync

	    var (
	    allPoolsMu Mutex

	    // allPools is the set of pools that have non-empty primary caches. Protected by either 1) allPoolsMu and pinning or 2) STW
	    allPools []*Pool

	    // oldPools is the set of pools that may have non-empty victim caches. Protected by STW.
	    oldPools []*Pool
	    )
	    This `allPools []*Pool` slice keeps track of all the active sync.Pool instances in our application.

	    Before each GC cycle starts, Go's runtime triggers a cleanup process that clears out the `allPools` slice. Here's how it works:

	 1. Before the GC kicks in, it calls `clearPool`, which transfers all the objects, including the private objects and shared pool chains, over to what's called the victim area.

	 2. These objects aren't immediately thrown away, they're held in this victim area for now.

	 3. Meanwhile, the objects that were already in the `victim area` from the last GC cycle get fully cleared out during the current GC cycle.

	    func poolCleanup() {
	    // Drop victim caches from all pools
	    for _, p := range oldPools {
	    p.victim = nil
	    p.victimSize = 0
	    }

	    // Move primary cache to victim cache
	    for _, p := range allPools {
	    p.victim = p.local
	    p.victimSize = p.localSize
	    p.local = nil
	    p.localSize = 0
	    }

	    // The pools with non-empty primary caches now have non-empty victim caches and no pools have primary caches.oldPools, allPools = allPools, nil
	    oldPools, allPools = allPools, nil
	    }

	    But why do we need this victim mechanism? Why does it take up to 2 GC cycles to clear all the objects in the pool?
	    The reason for using the victim mechanism in sync.Pool is to avoid suddenly and completely emptying the pool right after a GC cycle. If the pool were emptied all at once, it
	    could lead to performance issues, as any new requests for objects would require them to be recreated from scratch. So we move objects to the victim area first, sync.Pool ensures
	    there's a buffer period where objects can still be reused before they're fully discarded.
	    To sum up, an object in sync.Pool takes `at least 2 GC cycles` to be fully removed.

	    This could be a problem for programs with a low GOGC value, which controls how frequently the GC runs to clean up unused objects. If GOGC is set tool low, the cleanup process might
	    remove unused objects too quickly, leading to more cache misses.

	 2. Final words: Even with sync.Pool, if we're dealing with extremely high concurrency and slow GC, we might experience more overhead. In this case a good solution could be to implement
	    rate limiting on sync.Pool usage.
*/
const defaultCap = 1 << 10

type Object struct {
	Data []byte
}

func (o *Object) Reset() {
	const capTopBoundary = 1 << 13
	if cap(o.Data) >= capTopBoundary {
		o.Data = make([]byte, 0, defaultCap)
		return
	}

	o.Data = o.Data[:0]
}

func objectPoolCreation() {

	var objectPool = &sync.Pool{
		New: func() any {
			return &Object{Data: make([]byte, 0, defaultCap)}
		},
	}

	_ = objectPool
}

func ResetAndPut() {
	var objectPool = &sync.Pool{
		New: func() any {
			return &Object{Data: make([]byte, 0, defaultCap)}
		},
	}

	obj := objectPool.Get().(*Object)
	// some work
	obj.Reset()
	objectPool.Put(obj)
}
