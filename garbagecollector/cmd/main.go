package main

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
	GO SCHEDULER
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
	2. 2 minutes taken without GC invocation.
	3. Calling runtime.GC().

*/

func main() {
	// ArraysAllocation()

	// SlicesAllocation()
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
