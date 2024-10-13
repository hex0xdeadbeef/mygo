package main

import (
	"runtime"
	"sync"
)

/*
	CONCURRENCY LIMITING
1. Limiting concurrency brings benefits for both CPU and memory usage:
	1) CPU - The machine(s) running app only have a finite number of cores. By limiting concurrency, we reduce contention between threads and context switching. Context switching between threads
		is expensive, and when spikes occur, we want our CPUs to work as efficiently as possible.

	2) Memory - each running thread adds to the memory overhead of our app. By reducing the number of concurrent threads, we reduce the memory usage.

2. That said, limiting concurrency is not without downsides. The biggest issue is that it introduces complexity to our app which could lead to accidental deadlocks or inefficient resource utilization.

	For example, a thread could be waiting on I/O from storage or the network. I/O is unpredictable and often comes with large delays. By limiting the number of threads that can run concurrently,
		we increase the risk of under-utilizing the CPU while waiting on I/O ops. This is why it's important to only limit concurremcy for CPU-intensive workloads that doesn't depend on I/O.

3. Worker pool implementation.
	The canonical way to limit the amount of concurrency in the app is to spawn a limited number of worker threads and share a channel between them for dispatching work.

	The code first finds out the number of CPUs the machine has using Go's runtime library. It then spawns a number of worker threads equal to the number of CPUs, and each runs processNumber
		function. It also creates a shared channel to distribute work among them, it's called prod.

	This approach works well, but comes with a few downsides:
		1) Complicated implementation - this approach requires implementing start and stop procedures for workers.
		2) Additional synchronization - if the units of work are small, then the overhead of distributing work through a channel is relatively large.

		This approach provides a good solution for limited concurrency, and might work well for our app. That said, it's worth comparing it to a token-base implementation as described below.

4. Token-based implementation (Semaphore)

	An alternative approach is to implement a gate that workers need to pass before they can run. One implementation of this approach is to create a channel with a limited buffer size and make
	workers attempt to put a dummy value into that channel, a `token`, before running.

	Workers then will then take a dummy value from the channel once they are finished. Since the channel has a limited buffer size, only a limited number of workers can run at once.

	The code below uses a channgel-based object to restrict parallelism by using put/release token ops.

	The internal goroutine does some CPU-heavy processing, and so before running, it requests a token from semaphore by placing dummy value in the channel.

	This approach guarantees that no more than len(semaphore.c) workers are running at the same time, since each worker has to secure a slot before it does any intensive work.

5. Summary
	Limiting concurrency can help to bound resource usage and protect our application from crashing during load spikes. It ensures that our app processes load without interruptions in an optimal
	manner instead of wasting resources on context switches.
*/

const numbersToProcess = 1 << 11

// startWorkerPool starts and waits workers until all the work
// is done.
func startWorkerPool() {
	var workingInParallelThreadsNum = runtime.NumCPU()

	var (
		wg   = &sync.WaitGroup{}
		prod = getWorkProducer()
	)

	for i := 0; i < workingInParallelThreadsNum; i++ {
		wg.Add(1)
		go processNumber(wg, prod)
	}

	wg.Wait()
}

// getWorkProducer returns a producer channel and passes a data into it.
// As a result, the channel will be closed.
func getWorkProducer() <-chan int {

	var prodCh = make(chan int, 1)

	go func() {
		defer close(prodCh)

		for i := 0; i < numbersToProcess; i++ {
			prodCh <- i
		}
	}()

	return prodCh

}

// processNumber is a worker that is doing all the work on the number
func processNumber(wg *sync.WaitGroup, prod <-chan int) {
	defer wg.Done()

	for num := range prod {
		_ = num * num
	}
}

type Semaphore struct {
	c chan struct{}
}

func New(maxParallelWorkers int) Semaphore {
	return Semaphore{
		c: make(chan struct{}, maxParallelWorkers),
	}
}

func (s Semaphore) BookToken() {
	s.c <- struct{}{}
}

func (s Semaphore) ReleaseToken() {
	<-s.c
}

// startSemaphorePool creates a new semaphore, starts all the goroutines
// and these goroutines alternatively do the work of squaring numbers
func startSemaphorePool() {
	var (
		wg = &sync.WaitGroup{}

		sem = New(runtime.NumCPU())
	)

	for i := 0; i < numbersToProcess; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup, sem Semaphore, num int) {
			defer wg.Done()

			_ = num * num
		}(wg, sem, i)
	}

	wg.Wait()
}
