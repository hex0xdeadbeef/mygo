package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

/*
	CONCURRENCY BASICS
1. Concurrenct is about a structure and organization of things.
2. Actions are arranged through the independent threads and they communicate with each other doing theirs duties.
3. We can add some parallelism in a concurrent process, but it's not required at all. Its boundwidth will be increased, but the process itself won't be changed in its structure.
Because of decreasing of concurrency between threads.
4. Concurrency provide a possibility of a parallelism.
5. Concurrency pushes us to create a structure to solve a problem decoupling it. We can make actions in the parts parallel.
6. A concurrent solution is not always faster than sequential one. The velocity of execution depends on:
	1) Effectiveness of a structure (concurrency)
	2) Fraction of parts that can be executed in parallel
	3) A level of concurrency between calculative blocks

	PLANNING IN GO
1. OS Thread (M) - is the minimal unit of execution that can be computed by OS. If a process wants to do some actions simultaneously, it starts some OS Threads (M).
2. OS Threads (M) can be:
	1) Concurrent: two or more threads can be started, executed, finished in overlapping periods of time.
	2) Parallel: the same task can be executed simultaneously
3. OS is responsible for optimal planning in order to:
	1) All the OS Threads (M) can use processor time waiting for a while
	2) Workload will be distributed between all the kernels of processor evenly
4. Kernel of CPU executes different OS Threads (M). When it switches between OS Threads (M), the operation such as context switching is executed:
	1) An active OS Thread (M) consuming CPU cycles and being in "executing state" switches to a "runnable state". It means that it is ready to be launched and waiting for a
	free CPU kernel.
Context switching is an expensive operation since OS must save the current state of an OS Thread (M) before switching (for example: saving the values of registers)
5. OS Threads (M) are managed by OS.
6. Go's Goroutine (G) is managed by Go executing environment.
7. Goroutine (G) requires less place than an OS Thread (M) (2KB instead of 2MB). The tinier size makes context switching faster.

	HOW GO PLANNER WORKS
1. Terminology:
	1) G - Goroutine
	2) M - OS Thread (M means machine)
	3) P - OS Kernel (P means processor)
2. Each M assigns to P by OS Planner. Then each G is launched on M. The variable GOMAXPROCS defines the limit of Ms are responsible for simultaneous code execution of user
level. If an M is blocked in system call (for example: I/O, going to network or else), the OS Planner can launch more Ms. Starting from Go 1.5 GOMAXPROCS is equal to an amount
of free CPU kerneks by default.
3. A Gorourine (G) has simpler lifecycle than OS Thread (M). It can be in a following state:
	1) Execution (Executing). An execution of G is assigned to M and its actions is performed on it.
	2) Runnable (Readiness to be executed). Goroutine (G) is waiting for state-to-state transition
	3) Waiting. Goroutine (G) is stopped and waiting for finish of something. For example: syscall, synchronization or something else.

	PLANNING REALIZATION
1. What's happening when a Goroutine (G) has been already created, but cannot be executed yet (For example: other Ms are executing Gs). What will Go environment do? It'll place this
G into the queue.
2. Go environment is executing 2 types of queues for Goroutines (Gs):
	1) Local queues for each OS Kernel (P)
	2) Global queue is responsible for execution on all the OS Kernels (P)
3. Properties of Go environment planning:
	1) Only 1/61 of time the Go environment checks the Global Queue of Goroutines (Gs) to have a free goroutine to be executed. If it's empty ->
	2) Go Environment checks the Local Queues. If it's also empty ->
	3) Go Environment tries to steal Goroutines (Gs) from the other OS Kernels (Ps). If nothing has been found ->
	4) Go Environment checks the Global Queue of Goroutines in repeat. If nothing has been found again ->
	5) Network poll
4. Starting from Go 1.14 the Planner is preemptible. It means that when a Goroutine is being executed more than 10 mulliseconds, it'll be marked as a repressed one and can be detached
and replaced by another one. It allows to use the processor in the preriod when a long task is executed and other tasks need to be executed.

	PARALLEL PROCESSING ISN'T ALWAYS GOOD
1. If workload we want to distribute is too tiny relatively to creation of goroutines, theirs planning, theirs execution, there won't be a profit of parallelism using.
2. We can set a treshold of using parallelism performing ops on tiny datasets.
3. Thresholds should be picked after a bunch of benchmarks.

	WHEN WE SHOULD USE MUTEX | CHANNEL
1. In general synchronization between goroutines that use mutual data must be achieved by mutexes.
2. Concurrent Goroutines' (Gs) actions must be coordinated by using channels. It's called orchestration. This coordination relates to communication and must be executed by using
channels.
3. If we must restrict an access to a resource on a specific stage, we must use channels to signal other goroutines to stop their actions (there's another way - sync.Cond).
4. Parallel Goroutines (Gs) need Mutexees, Concurrent Goroutines (Gs) need Channels.

	DATARACE
1. Datarace. Datarace is when two goroutines clobber the same area of memory and at least one of them is writing into it.
2. The ways to solve datarace problem:
	1) To make the operation atomic. To achieve this we can use sync/atomic pkg. An atomic operation cannot be preempted, it prevents the situation when the access is provided to 2
	places simultaneosly.
	2) To synchronize execution by using mutexes (mutual exclusions).
	3) To use channels to make a synchronous execution.

	RACE CONDITION & DETERMINISM OF GOROUTINE EXECUTION
1. Race Condition happens when behavior depends on a sequence|time of events completion we cannot control.
2. Specific goroutines' sequence execution provision - is a deal of coordination and orchestration. This orchestration can be achieved by channels using.
3. Coordination and Orchestration can guarantee that only one goroutine has an access to a specific part of our application.

	GO MEM-MODEL
1. There's no possibility within a single goroutine of:
	1) Asynchronous access to a variable
	2) The order of actions is guaranteed of a program
2. Receiving from an unbuffered channel happens before sending competion on this channel.
*/

func main() {
	// DataraceA()

	// DataraceB()

	// DataraceC()

	// DataraceD()

	// NonDeterministicCalls()

	// ReceiveBeforeSendCompletionBuffered()

	ReceiveBeforeSendCompletionUnbuffered()

}

/*
1. Increment operation makes 3 things:
 1. Reading i
 2. Incrementing the i
 3. Writing the new val into i

2. There's no guarantee that the 1st goroutine will start first and finish before the 2nd.
3. There might be a case when both goroutines execute alternatively.
*/
func DataraceA() {
	var (
		wg = &sync.WaitGroup{}
		i  int
	)

	wg.Add(2)
	go func() {
		defer wg.Done()
		for j := 0; j < 1000; j++ {
			i++
		}
	}()

	go func() {
		defer wg.Done()
		for j := 0; j < 1000; j++ {
			i++
		}
	}()

	wg.Wait()

	fmt.Println(i)
}

func DataraceB() {
	var (
		wg = &sync.WaitGroup{}
		i  int64
	)

	wg.Add(2)
	go func() {
		defer wg.Done()
		for j := 0; j < 1000; j++ {
			atomic.AddInt64(&i, 1)
		}
	}()

	go func() {
		defer wg.Done()
		for j := 0; j < 1000; j++ {
			atomic.AddInt64(&i, 1)
		}
	}()

	wg.Wait()

	fmt.Println(i)
}

func DataraceC() {
	var (
		wg = &sync.WaitGroup{}
		mu = &sync.Mutex{}
		i  int
	)

	wg.Add(2)
	go func() {
		defer wg.Done()
		for j := 0; j < 1000; j++ {
			mu.Lock()
			i++
			mu.Unlock()
		}
	}()

	go func() {
		defer wg.Done()
		for j := 0; j < 1000; j++ {
			mu.Lock()
			i++
			mu.Unlock()
		}
	}()

	wg.Wait()

	fmt.Println(i)
}

func DataraceD() {
	var (
		i int

		incrementer = func() chan int {
			incr := make(chan int)

			go func() {
				defer close(incr)

				for j := 0; j < 1000; j++ {
					incr <- 1
				}

			}()

			return incr
		}

		ticker = time.NewTicker(time.Second * 1)
	)

	incrA, incrB := incrementer(), incrementer()
	defer ticker.Stop()

Cycle:
	for {
		select {
		case v, ok := <-incrA:
			if !ok {
				incrA = nil
				break
			}

			i += v
		case v, ok := <-incrB:
			if !ok {
				incrB = nil
				break
			}

			i += v
		case <-ticker.C:
			break Cycle
		}

	}

	fmt.Println(i)
}

/*
The calls of this function will print non deterministic vals of i.
It happens because we cannot guarantee a specific order of gorountines planning.

In this case the time of actions is the order of goroutines execution.
*/
func NonDeterministicCalls() {
	var (
		mu = &sync.Mutex{}
		wg = &sync.WaitGroup{}

		i int
	)

	wg.Add(2)
	go func() {
		defer wg.Done()

		mu.Lock()
		defer mu.Unlock()

		i = 1
	}()

	go func() {
		defer wg.Done()

		mu.Lock()
		defer mu.Unlock()

		i = 2
	}()

	wg.Wait()

	fmt.Println(i)
}

func ReceiveBeforeSendCompletionBuffered() {
	i := 0
	ch := make(chan struct{}, 1)

	go func() {
		i = 1 // It can happen before/after printing
		<-ch
	}()

	ch <- struct{}{}
	fmt.Println(i) // It can happen before/after replacing value
}


// Receiving from an unbuffered channel happens before the completion of sending on this channel.
func ReceiveBeforeSendCompletionUnbuffered() {
	i := 0
	ch := make(chan struct{})

	go func() {
		i = 1
		<-ch
	}()

	ch <- struct{}{}
	fmt.Println(i)
}
