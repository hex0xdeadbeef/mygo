package main

import (
	"log"
	"runtime"
	"sync/atomic"
)

// https://www.youtube.com/watch?v=UVqpl4PExkM
/*
	GO GARBAGE COLLECTOR
1. Arrays allocation:
	1) Arrays <= 10 MB are allocated in STACK
	2) Arrays >10 MB are allocated in HEAP
2. Slices allocation:
	1) Slices < 64KB are allocated in STACK
	2) Slices >= 64 KB are allocated in HEAP
3. Go's GC is based on the algorithm "Mark And Sweep"
4. Stages of "Mark And Sweep":
	1) Mark: GC marks all the reachable objects as "alive"
	2) Sweep: GC goes through all the tree of marked objects and sweeps all the unreachable objects from the heap
5. GC consumes:
	1) Physical memory
	2) CPU time
6. "Stop the World". When the GC is at the stages of:
	1) Freparing of Markering
	2) Finish of Markering
the overall application is stopped.
7. "Alive and Dead" heaps.
	1) "Alive" heap is the heap that includes all the objects that haven't been sweept by GC at the previous cycle
	2) "Dead" heap is the heap that was swept by GC at the previous cycle
When the overhead above the "Alive" heap reaches 100% of "Alive" heap the GC starts cleaning
8. We can manage the percentage of ratio of New Heap to "Alive" heap with using GOGC global variable.
9. There's the global varbiable GOMEMLIMIT that sets limit on max amount of memory can be used by Go runtime. It helps in the following situation:
	When the overall amount of memory reaches the limit set the GC is invoked.
10. "Death spiral " When there's the situation when we have infinite amount of leaked goroutines, the use of GOMEMLIMIT goes away because our amount of goroutines will grow and GC
will be invoked more often. It results in 100% time consuming by GC -> the overall time of execution will grom dramatically.
11. To beat the problem of "Death spiral" developers made the following:
	1) GOMEMLIMIT is soft manage of memory. GOMEMLIMIT doen't fully guarantee that after reaching GOMEMLIMIT value, GC will be called. The program is allowed to go beyond this limit.
	2)
*/

// https://www.youtube.com/watch?v=CX4GSErFenI&t=2757s
/*
	GO GARBAGE COLLECTOR
1. Reachability:
	1) Root objects
	2) Objects bound to root objects
	3) Reachable from bound to root objects
Root -> bound to root -> reachable from bound to root
2. "Mark and Sweep"
	1) Each object has a bit of reachability (0 or 1)
	2) From scratch all the elems except root objects aren't marked
The algorithm:
	1) Stop the execution
	2) Recursively mark all the objects setting on them "a bit of reachability"
	3) Delete all the unreachable objects
	4) Resume our execution
3. The lacks of direct z"Mark And Sweep":
	1) Pauses
	2) All the heap traversal
	3) Fragmentation
4. In go "Breath First Mark" is used. The algorithm is:
	1) Split all the objects into 3 groups.
	2) At each step mark all the reachable objects as "Scanned"
	3) In "Will be scanned" mark all that are referenced by a link
	4) After BFS finish, delete all that are not in 1st and 2nd group
5. "Breath Mark" Groups are":
	1) Black - investigated objects
	2) Gray - the objects that are waiting for investigation
	3) White - the objects haven't been investigated
6. If an object is added during markering, it'll be marked as gray.
*/

// https://habr.com/ru/amp/publications/753244/
/*
1. GC stages
	1) Sweep Termination.
	"Stop the world". It's needed so that all the physical kernels reach the point when it's safe to start GC. After this all the blocks marked as garbage are sent to
	Scavenger to "be eaten". After this, Scavenger will send it to the OS. It allows not to request OS to give more memory, than app queries.
	2) Mark.
	Write barrier is turned on. Tasks of markering root variables are queued. Root variables are global variables and all the contents of stack.

	"Start the world". All the markering work are executed in the special painter-workers started by scheduler. These workers go through all objects in stack/heap. New allocated objects get "Black"
	colour. So that the stack of goroutine will be marked, it's stopped, painted and started again.

	It resumes until the queue of gray objects empties. In this state GC consumes until 25% of CPU work.
	3) Mark termination.
	"Stop the world". Stop our workers. Clean Ms.
	4) Sweep
	Write barrier is off. "Start the world". New allocated objects will be marked as "White". Allocation can happen above blocks of memory marked as garbage. Moreover, the goroutine that will return
	the memory to OS will be started.
2. When GC started?
	1. The limit of GOGC exceeded.
	2. 2 minutes taken without GC invokation.
	3. Calling runtime.GC().

*/

/*
	WRITE BARRIER
1. Within the single thread the compiler can swap instructions in any order if it doesn't affect on the program.
2. The compiler can rewrite our code.
3. There are only two operations with memory:
	1) Load (reading)
	2) Write (writing)
And it results in four operations we don't want to mix:
	1) LoadLoad
	2) LoadStore
	3) StoreStore
	4) StoreLoad
4. Using "Memory Barrier" we restrict CPU/compiler to swap operations.
5. Full barrier says to CPU/compiler not to mix all the kinds of instructions (read/write) (StoreStore + LoadLoad)
6. Write barrier says to CPU/compiler not to mix writes. It means that writes cannot be mixed, but reads can.
7. Read barrier says to CPU/compiler not to mix reads. It means that reads cannot be mixed, but writes can.
8. Memory barriers guarantee only the order of instructions relatively to (DRAM, caches). It limits only instructions optimisations, not the other ones.
9. Atomic values forces CPU/Compiler to use "Full Barrier", it means that all the instruction above and below cannot go through this barrier.
*/

// https://www.ardanlabs.com/blog/2018/12/garbage-collection-in-go-part1-semantics.html
/*
	GO GARBAGE COLLECTOR
1. Garbage Collectors have the responsibility of:
	1) Tracking heap memory allocations
	2) Freeing up allocations that are no longer needed
	3) Keeping allocations that are still in-use
2. As of version 1.12 the Go programming language uses a non-generational concurrent tri-color "Mark and Sweep" collector.
3. Collector Behavior.
	When a collection starts, the collector runs through three phases of work.
		1) Two of these phases create "Stop the World" (STW) latencies
		2) And the other one creates latencies that slow down the throughput of the app
	There are only three phases:
		1) Mark Setup - 		STW
		2) Marking - 			Concurrent
		3) Mark Termination - 	STW
4. Mark Setup (STW) phase.
	1) Goroutines stopping.
	When a collection starts, the first activity that must be performed is: turning on the "Write Barrier". The purpose of "Write Barrier" is to allow the collector to maintain data integrity
	on the heap during a collection sicne both the collector and application goroutines will be running concurrently.

	In order to turn on the Write Barrier, every app goroutine running must be stopped. This activity is usually very quick (10-30 microsecs) on average.
	That is, as long as the app goroutines are behaving properly.

	The only way to do this is for the collector to watch and wait for each goroutine to make a function call. These function calls guarantee the goroutines are at a safe point to be stopped.

	2) Waiting for a goroutine
	If there's a situation when the GC is watiting for a goroutine that is stalled, the other Ps cannot service any other goroutines, so it's important that goroutines make function calls in
	reasonable timeframes.
5. Marking (Concurrent)
	Once the Write Barrier is turned on, the collector commences with the Marking Phase.

	The first thing the collector does it take 25% of the available CPU capacity for itself.

	The collector uses goroutines to do the collection work and needs the same P's and M's the application Goroutines use.
	It means that for our 4 threaded Go program, one entire P will be dedicated to collection work.

	This work starts by inspecting the stacks for all existing goroutines to find root pointers to heap memory.
	Then the collector must traverse the heap memory graph from those root pointers found.
	While the Marking work is happening on P1, application work can continue on P2, P3, P4. This means the impact of the collector has been minimized to 25% of the current CPU capacity.

		1) The use of Mark Assist. What if it's identified during the collection that the Goroutine dedicated to P1 won't finish the Marking work before the heap memory in-use reaches its limit?
		What if only one of those 3 Goroutines performing app work is the reason the collector won't finish in time?

		In this case, new allocations have to be slowed down and specifically from that Goroutine. If the collector determines that it needs to slow down allocations, it'll recruit the app
		Goroutines to assist with the Marking work. This is called "Mark Assist". The amount of time any application Goroutine will be placed in "Mark Assist" is proportional to the amount of data
		it's adding to heap memory.

		One posisive effect of Mark Assist is that it helps to finish the collection faster.

		If any given collection ends up requiring a lot of Mark Assist, the collector can start the next garbage collection earlier. This is done in an attempt to reduce the amount of Mark Assist
		that will be necessary on the next collection.
6. Mark Termination - STW
	Once the Marking work is done, the next phase is Mark Termination. This is when the Write Barrier is turned off, various clean up tasks are performed, and the next collection goal is calculated.
	Goroutines that find themselves in a tight loop during the Marking Phase can also cause Mark Termination STW latencies to be extended.

	Mark Termination takes 60-90 microsecs on average.

	Once the collection finished, every P can be used by the application Goroutines again and the app is back to full throttle.
7. Sweeping
	There's another activity that happens after a collection is finished called Sweeping.

	Sweeping is when the memory associated with values in heap memory that were not marked as in-use are reclaimed. This activity occurs when app Goroutines attempt to allocate new vals in heap
	memory.
*/

func main() {
	// ArraysAllocation()

	// SlicesAllocation()

	// MemoryBarrierB()
}

// go build -gcflags="-m"
func ArraysAllocation() {
	a1 := [1_310_720]int{}
	a1[0] = 1

	a2 := [1_310_721]int{} // ./main.go:14:18: moved to heap: s
	a2[0] = 1
}

func SlicesAllocation() {
	s1 := make([]int, (64 * 1024 / 8))
	s1[0] = 1

	s2 := make([]int, (64*1024/8)+1) // ./main.go:32:12: make([]int, 8193) escapes to heap
	s2[0] = 1
}

func MemoryBarrierA() {
	first := 10
	second := 20
	third := 0

	first++  // independent
	second++ // independent
	third++  // independent

	third = first + second // dependent from the results above -> cannot be swapped by the compiler

}

var str string
var b atomic.Bool

func setup() {
	str = "hello, world"
	b.Store(true)

	if b.Load() {
		log.Println(len(str))
	}
}

func MemoryBarrierB() {
	go setup()

	for !b.Load() {
		runtime.Gosched()
	}

	log.Println(str)
}
