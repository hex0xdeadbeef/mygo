package main

// https://www.youtube.com/watch?v=P2Tzdg8n9hw
/*
	GO SCHEDULER
1. Web-services work with I/O workload mostly and this work requires context switching.
2. Using threads is expensive.
	1) Thread has the wait of 2-8 MB
	2) Switching between contexts of threads takes place on kernel level
3. Thread is just a context to provide some code. To launch some code we need:
	0) Registers
	1) Instruction pointer
	2) Stack
4. To switch from thread to thread we need:
	1) Take of the registers of the old one and put them on stack.
	2) Get the registers of the new thread and upload them on the CPU
	3) Launch the code
In this situation we discard the use of CPU caches, so the cache is cleaned.
5. The Goroutine (G) is just the superstructure over OS Threads.
	1) It's created inside our application.
	2) It's just a lightweight thread
	3) The kernel of OS doesn't know about Goroutines.
	4) It's start from 2 KB, it dumps the workload of DRAM
6. The idea of Goroutines is
	1) A programmer doesn't work with OS threads directly
	2) We pretend that there are only Goroutines
	3) The complexity of goroutines' distribution is abstracted
7. The OS thread just performs the code of Goroutine.
8. The OS scheduler is preemptive.
9. The Go's scheduler is cooperative
	1) The compiler will add some points into the source code to provide
	context-switching.
10. Working with goroutines we use "N:M model"
	1) We perform N Goroutines on M OS threads
11. Local queues of Goroutines.
	1) Each Local Queue will be bound to a specific OS Thread (M)
	2) The OS thread will reference only to its own local queue
	3) This way eliminates work with synchronization primitives
	4) The max size of the local queue of Goroutines is 256 Goroutines
12. Processor (P) - is the abstraction that tells that we have the specific resources of execution:
	1) Rights of execution
	2) Resources of execution (queues of Goroutines)
	3) An OS Thread belongs to only one P (Processor)
13. GMP Model.
	1) G - Goroutine (What we execute)
	2) M - Machine (OS Thread) (Where we execute)
	3) P - Processor (Rights and resources of execution)
14. What the place we should put a new Goroutine created into?
	1) The goroutine created will be put into local queue
15. What we should do when the local queue is empty? (Work stealing)
	1) We should try to steal some goroutines from other Ps (Other local queues)
	2) In this case we need synchronization (lock-free synchronization (atomics))
	3) A local queue will be picked randomly
	4) We will try to steal 4 times
	5) A half of goroutines will be stealed from other P
16. Syscalls
Syscall is work with the kernel of OS. In this stage the OS thread is blocked fully. It turns into the "Waiting" stage. The OS scheduler removes this OS thread from the kernel and
puts a new one istead.
	1) When a syscall is performed, we unpin the OS thread (M) from its P (Local queue)
	2) Create a new OS thread and put it into the P
	3) We pass all the goroutines left of this M created
A bunch of these OS Threads is called "Thread Pool" and is used to keep the OS Threads that are in syscalls

The Go's scheduler will unpin the OS Thread (M) from the processor (P) when it understands that this OS Thread will be blocked by syscall for a long time. In other cases it'll allow this OS Thread to be blocked and not to be unpinned.

In the cases when the OS Thread (M) is in syscall and not unppinned but executes too long it'll be evicted from P by "sysmon".
	1) Sysmon checks periodically the blocked OS thread that is on a P and in syscall
	2) If this OS thread is in syscall too long, it'll be evicted by sysmon
	3) A new OS Thread will be pinned to this P
	4) Sysmon is just another OS thread
17. Where a goroutine should be put into after syscall?
	1) If the queue has enough place and this OS thread isn't performing another syscall, put this goroutine into the local queue of this P
	2) If the queue of this P is full (256 goroutines are already placed) or this OS thread is still performing syscall, we should search a new P (local queue)
	3) If there are no free queue, put this goroutine into the global queue
18. When we work with blocking calls, we'll use multiplexer (netpoller). It's used to catch signals of events (event loop). When there's a blocking call, we'll send this call to the netpoller.

19. Queues
	1. Global queue. What should we do wneh there's filled Global queue?
		runtime schedule() {
			// only 1/61 of the time. check the global runnable queue for a G
			// if not found, check the local queue
			// if not found, try to steal from other Ps
			// if not found, check the global runnable queue
			// if not found, poll the network
		}
21. Goroutine has three states
	1) Running (is performing)
	2) Runnable (is ready to be performed)
	3) Waiting (is blocked on mutex/syscall)
20. Synchronizations. When a goroutine is locked we just evict it from OS thread and put another one from its Ps queue. The evicted goroutine winds up ind "Waiting Queue".
	1) When a goroutine is in "Waiting Queue" and waiting more than 1 ms, it'll
	be picked as the next to use the mutex.
21. Long Cycles (Adding preemptive part)
	1) If there's a goroutine that is performing more than 10 ms, it'll be evicted from OS Thread. The goroutine will reach the safe pointer before it is evicted.
22. When goroutines switching happens:
	1) morestack() -> newstack()
	2) runtime.Gosched()
	3) locks
	4) networks I/O
	5) syscalls
*/

// https://www.youtube.com/watch?v=rloqQY9CT8I
/*
	GO SCHEDULER
1. Abstractions
	1. M - OS Thread (Machine)
		1) It's default to have the number of Ms equal to the number of Ps
		2) Only one goroutine can be active on a M
		3) Threads can move from P1 to P2
	2. P - Processor
		1) Go creates the amount of P that is equal to the number of physical kernels
		2) Each P has its own queue of Goroutines (Gs). It's called "Local queue" or "Runqueue"
	3. G - Goroutine
		1) The number of Goroutines is unlimited
2. Situations when goroutine can pass the time to another goroutine:
	1) input/output
	2) using channels
	3) syscalls
	4) runtime.Gosched()
3. schedule() call
	1) Check whether we need to collect garbage. If not, get the goroutine.
	2) If there's no goroutine we call steal_or_wait()
	3) execute(g)
*/

func main() {

}
