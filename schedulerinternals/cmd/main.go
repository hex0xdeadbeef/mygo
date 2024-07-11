package main

// https://www.ardanlabs.com/blog/2018/08/scheduling-in-go-part1.html
/*
	OS SCHEDULER

	QUICK OVERVIEW
1. A program is just a series of insructions that need to be executed one after the other sequentially. To make this happen the operating system uses the concept of a Thread (M).
	It's the job of a Thread to account for and sequentially execute the set of instructions it's assigned. Execution continues until there are no more instructions for the Thread to execute.
2. Every program we run creates a Process and each Process is given an initial Thread.
3. Threads have the ability to create more Threads. All these different Threads run independently of each other and scheduling decisions are made at the Thread level, not at the Process level.
4. Threads can run concurrently (each taking a turn on an individual core) or in parallel (each running at the same time on different cores)
5. Threads also maintain their own state to allow for the safe, local, and independent execution of their insructions.
6. The OS Scheduler is responsible for making sure that cores aren't idle if there are Threads that can be executing.
7. It must also create the illusion that all the Threads can execute are executing at the same time. In the process of creating this illusion, the scheduler needs to run Threads with a higher
priority over lower priority Threads. However, Threads with a lower priority can't be starved of execution time.
8. The OS Scheduler also needs to minimize scheduling latencies as much as possible by making quick and smart decsions.

	INSTRUCTION POINTER (IP) or Program Counter (PC)
1. The "Program Counter" (PC), which is sometimes called the "Instruction Pointer", is what allows a Thread to keep track of the next instruction to execute. In most processors, the PC points
	to the next instruction and not the current instruction.
2. The computer keeps track of the next line to be executed by keeping its address in a special register called the "Instruction Pointer" (IP) or "Program Counter" (PC). This register is
	relative to CS as segment register and points to the next instruction to be executed. The contents of this register is updated with every instruction executed. Thus a program is executed
	sequentially line by line.

	THREAD STATES
1. Thread state dictates the role the scheduler takes with the Thread. The Thread can be in one of three states: "Waiting", "Runnable" or "Executing"
2. "Waiting" status
	"Waiting" status means the Thread is stopped and waiting for something in order to continue. This could be for reasons like, waiting for hardware (disk, network), the OS (syscalls) or
	synchronization calls (atomic, mutexes). These types of latencies are a root cause for bad performance.
3. "Runnable" state
	"Runnable" state means the Thread wants time one a core so it can execute its assigned machine instructions. If we have a lot of Threads that want time, then Threads have to wait longer
	to get time. Also, the individual amount of time any given Thread gets is shortened, as more Threads compete for time. This type of scheduling latency can also be a cause of bad performance.
4. "Executing" state
	"Executing" state means that the Thread has been placed on a core and is executing its machine instructions. The work related to the app is getting done. This is everyone wants.

	TYPES OF WORK
1. There are two types of work a Thread can do. The first one is called "CPU-Bound" and the second is called "I/O-Bound"
2. "CPU-Bound" workload
	"CPU-Bound" workload is the work that never creates a situation where the Thread may be placed in "Waiting" states. This is work that is continuosly making calculations. A Thread calculating
	the Pi number to the Nth digit would be "CPU-Bound"
3. "I/O-Bound" workload
	"I/O-Bound" workload is the work that causes Threads to enter into "Waiting" states. This is work that consists in requesting access to a resource over the network or making system calls
	into the operating system. A Thread that need to an access a database would be "I/O-Bound". I would include synchronization (mutexes, atomic), that cause the Thread to wait as the part of
	this category.

	CONTEXT SWITCHING
1. If we running on Linux, Mac or Windows, we're running on an OS that has a "Preemptive" scheduler. It means that:
	1) The scheduler is unpredictable when it comes to what Threads will be to run at any given time. Threads' priorities together with events (like receiving data on the network) make it
	impossible to determine what the scheduler will choose to do and when.
	2) We must never write code based on some perceived behavior that we've been lucky to experience but it's not guaranteed to take place every time. We must control the synchronization and
	orchestration of Threads if we need determinism in our app.
2. "Context Switch" happens when the scheduler pulls an "Executing" Thread off a core and replace it with a "Runnable" Thread. The Thread what was selected from the run queue moves into an
	"Executing" state. The Thread that was pulled can move back into a "Runnable" state (if it still has the ability to run), or into a "Waiting" state (if it was replaced because of an "I/O-
	Bound" type of a request)
3. Context switches are considered to be expensive because it takes times to swap Threads on and off a core. The amount of latency incurrent during context switch depends on different factors
	but it's reasonable for it to take between ~1000 and ~1500 nanosecs.
4. Considering the hardware should be able to reasonably execute (on average) 12 instructions per nanosecond per core, a context switch can cost us from ~12,000 to ~18,000 instructions of latency
5. In essence, our program is losing the ability to execute a large number of instructions during context switch.
	1) If our program is focused on "I/O-Bound" workload, then context switches are going to be an advantage. Once a Thread moves into a "Waiting" state, another Thread in a "Runnable" state is
	there to take its place. This allows the core to always be doing work. This is one of the most important aspects of scheduling. We mustn't allow a core to go idle if there's work (Threads in
	a Runnable state) to be done.
	2) If our program is focused on "CPU-Bound" workload, then context switches are going to be a performance nightmare. Since the Thread always has work to do, the context switch is stopping
	that work from processing. This situation is in stark contrast with what happens with an "I/O-Bound" workload.

	LESS IS MORE
1. When there are more Threads to consider, and "I/O-Bound" work happening, there's more chaos and nondeterministic behavior. Things take longer to schedule and execute.
2. This is why the rule of the game is "Less is More".
	1) Less Threads in a "Runnable" state mean less scheduling overhead and more time each Thread gets over time.
	2) More Threads in a "Runnable" state mean less time each Thread gets over time. That means less of our work is getting done over time as well.

	FINDING THE BALANCE IN CORES and THREADS NUMBER
1. There's a balance we need to find between the number of Cores and Threads we need to get the best throughput for our app.
2. When it comes to managing this balance, "Thread Pool" were a great answer.

	CACHE LINES
1. Accessing data from main memory has such a high latency cost (~100 - 300 clock cycles) that processors and cores have local caches to keep data close to the hardware threads that need it.
2. Accessing data from caches have a much lower cost (~3 - 40 clock cycles) depening on the cache being accessed.
3. Data is exchanged between the processor and main memory using cache lines. A cache line is 64-byte chunk of memory that is exchanged between main memory and the caching system. Each core is
	given its own copy of any cache line it needs, which means the hardware uses value semantics. This is why mutations to memory in multithreaded apps can create performance nightmares.

	When multiple Threads running in paralell are accessing the same data value or even data values near one another, they will be accessing data on the same cache line. Any Thread running on
	any core will get its ow copy of that same cache line.
4. False sharing
	If one Thread on a given core makes a change to its copy of the cache line, then through the magic of hardware, all other copies of the same cache line have to be marked dirty. When a
	Thread attempts to read or write access to a dirty cache line, main memory access (~100 - 300 clock cycles) is required to get a new copy of the updated cache line.
*/

// https://www.youtube.com/watch?v=P2Tzdg8n9hw
/*
	GO SCHEDULER
1. Web-services work with I/O workload mostly and this work requires context switching.
2. Using threads is expensive.
	1) Thread has the weight of 2-8 MB
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
