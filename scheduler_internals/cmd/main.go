package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

// Ardan Labs

// https://www.ardanlabs.com/blog/2018/08/scheduling-in-go-part1.html
/*
	OS SCHEDULER

	QUICK OVERVIEW
1. A program is just a series of insructions that need to be executed one after the other sequentially. To make this happen the operating system uses the concept of a Thread (M).
	It's the job of a Thread to account for and sequentially execute the set of instructions it's assigned. Execution continues until there are no more instructions for the Thread to execute.
2. Every program we run creates a Process and each Process is given an initial Thread (M).
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
	"Runnable" state means the Thread wants time one a core so that it can execute its assigned machine instructions. If we have a lot of Threads that want time, then Threads have to wait longer
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

	CONTEXT SWITCHES
1. If we running on Linux, Mac or Windows, we're running on an OS that has a "Preemptive" scheduler. It means that:
	1) The scheduler is unpredictable when it comes to what Threads will be to run at any given time. Threads' priorities together with events (like receiving data on the network) make it
	impossible to determine what the scheduler will choose to do and when.
	2) We must never write code based on some perceived behavior that we've been lucky to experience but it's not guaranteed to take place every time. We must control the synchronization and
	orchestration of Threads if we need determinism in our app.
2. "Context Switch" happens when the scheduler pulls an "Executing" Thread off a core and replace it with a "Runnable" Thread. The Thread what was selected from the run queue moves into an
	"Executing" state. The Thread that was pulled can move back into a "Runnable" state (if it still has the ability to run), or into a "Waiting" state (if it was replaced because of an "I/O-
	Bound" type of a request)
3. Context switches are considered to be expensive because it takes times to swap Threads on and off a core. The amount of latency incurrent during context switch depends on different factors
	but it's reasonable for it to take ~(1000 - 1500 nanosecs).
4. Considering the hardware should be able to reasonably execute (on average) 12 instructions per nanosecond per core, a context switch can cost us from ~12,000 to ~18,000 (1250 nanosecs)
instructions of latency
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
3. Data is exchanged between the processor and main memory using cache lines. A cache line is (64 | 128)-byte chunk of memory that is exchanged between main memory and the caching system. Each
	core is given its own copy of any cache line it needs, which means the hardware uses value semantics. This is why mutations to memory in multithreaded apps can create performance nightmares.

	When multiple Threads running in paralell are accessing the same data value or even data values near one another, they will be accessing data on the same cache line. Any Thread running on
	any core will get its ow copy of that same cache line.
4. False sharing
	If one Thread on a given core makes a change to its copy of the cache line, then through the magic of hardware, all other copies of the same cache line have to be marked dirty. When a
	Thread attempts to read or write access to a dirty cache line, main memory access (~100 - 300 clock cycles) is required to get a new copy of the updated cache line.
*/

// https://www.ardanlabs.com/blog/2018/08/scheduling-in-go-part2.html
/*
	GO SCHEDULER
1. When our program starts, it's given a Logical Processor (P) for every virtual core that is identified on the host machine.
	1) If we have a processor with multiple hardware threads per physical core, each hardware thread will be presented to our Go program as a virtual core.
	2) On MBP M3 Pro we have 12 physical cores, and under the hood there are 1 hardware threads per physical core. This will report to the Go program that 12 virtual cores are available for
	executing OS Threads in parallel.
	3) Any Go's program we run on MBP M3 Pro machine will be given 12 P's

	GMP MODEL
1. Every P is assigned an OS Thread (M). The "M" stands for machine. This thread is still managed by the OS and the OS is still responsible for placing the thread on a Core for execution.
	This all means that when we run a Go program on my machine, we have 12 threads available to execute all the work, each individually attached to a P.
2. Every Go program is also given an initial Goroutine ("G"), which is the path of execution for a Go program.
3. We can think of Goroutines as app-level threads and they are similar to OS Threads (Ms) in many ways. Just as OS Threads are context-switched on and off a core, Goroutines are
	context-switched on and off an M.
4. There are two different queues in the Go scheduler: Global Run Queue (GRQ) and the Local Run Queue (LRQ).
	1) Each P is given a LRQ that manages the Goroutines assigned to be executed within the context of a P. These Goroutines take turns being context-switched on and off the M assigned to that P
	2) The GRQ is for Goroutines that have not been assigned to a P yet. There's a process to move Goroutines from GRQ to a LRQ

	COOPERATIVE SCHEDULER
1. As we discussed in the first part, the OS Scheduler is preemptive. Essentially it means we cannot predict what the scheduler is going to do at any given time. The kernel is making decisions
and everything is non-deterministic. Apps that run on top of the OS have no control over what is happening inside the kernel with scheduling unless they leverage synchronization primitives.
2. The Go's scheduler is the part of Go's runtime and the Go's runtime is built into our app. This means Go's scheduler runs in user space, above the kernel.
3. The Go's scheduler is "implicitly" cooperative. Being cooperating scheduler means the scheduler needs well-defined user space events that happen at safe points in the code to make decisions.
4. It's relevant to think of the Go's scheduler as a preemptive scheduler and since the scheduler is non-determinisctic.

	GOROUTINES STATES
1. A Goroutine can be in one of 3 states similar to OS scheduler: "Waiting", "Runnable", "Executing"
2. "Waiting" state
	"Waiting" state means the Goroutine is stopped and waiting for something in order to continue. This could be for reasons like waiting for the operating system (system calls) or
	synchronization calls (atomic and mutexes operations). These types of latencies are a root cause for bad performance.
3. "Runnable" state
	"Runnable" state means the Goroutine wants time on an M so it can execute its assigned instructions. If we have a lot of Goroutines that want time, then Goroutines have to wait longer
	to get time. Also, the individual amount of time any given Goroutines gets is shortened as more Goroutines compete for time. This type of scheduling latency can also be a cause of bad
	performance.
4. "Executing" state
	"Executing" state means the Goroutine has been placed on an M and is executing its instructions. The work related to the app is getting done. This is what everyone wants

	CONTEXT SWITCHING
1. The Go's scheduler requires well-defined user-space events that occur at safe points in the code to context switch from. These safe points manifest themselves within function calls. Function
	calls are critical to the health of the Go's scheduler, if we run any tight loops that are not making function calls, we'll cause latencies within the scheduler and garbage collections. It's
	critically important that function calls happen within reasonable timeframes.
2. There are four classes of events that occur in a Go programs that allow the scheduler to make scheduling decisions. It doesn't mean it'll always happen on one of these events. It means the
	scheduler gets the opportunity:
		1) The use of keyword "go" FuncName(...) | func(...) {}()
		2) Garbage Collection
		3) System Calls
		4) Synchronization and Orchestration
3. "go" word invocation
	The keyword "go" is how we create Goroutines. Once a new Goroutine is created, it gives the scheduler an opportunity to make a scheduling decision.
4. Garbage collection
	Since the GC runs its own set of Goroutines on 25% of all the Threads (Ms), those Goroutines need time on an Ms to run. This causes the GC to create a lot of scheduling chaos. However, the
	scheduler is very smart about what a Goroutine is doing and it will leverage that intelligence to make smart decisions. One smart decision is context-switching a Goroutine that wants to
	touch the heap with those that don't touch the heap during GC. When GC is running a lot of scheduling decisions are being made.
5. System calls
	If a Goroutine makes a syscall that will cause the Goroutine to block the M, sometimes the scheduler is capable of context-switching Goroutine off the M and context-switch a new Goroutine
	onto the same M. However, sometimes a new M is required to keep executing Goroutines from the P (on that blocked M was beign executed) that are queued in the P.
6. Synchronization and Orchestration
	If an atomic, mutex, or channel operation (atomics + mutexes are used) call will cause the Goroutine to block, the scheduler can context-switch a new Goroutine to run. Once the Goroutine can
	run again, it can be re-queued and eventually context-switched back on a free M.

	ASYNCHRONOUS SYSTEM CALLS
1. When the OS has the ability to handle a syscall asychronously, something called the "network poller" can be used to process the syscall more efficiently. This is accomplished by using kqueue
(MacOS), epoll (Linux).
2. Networking-based system calls can be processed asychnronously by many of the OSs we use today. This is where network poller gets its name, since its primary use is handling networking ops.
3. By using the network poller for networking system calls, the scheduler can prevent Goroutines from blocking the M when those system calls are made. This helps to keep the M available to
	execute other Goroutines in the P's LRQ without the need to create new Ms. This helps to reduce scheduling load on the OS.
4. When the Goroutine's network call needs to happen, it's removed from M and moved to Network Poller to be executed asynchronously. After that a Goroutine in the P's LRQ of that M can be
	executed instead of the moved one to "Network Poller".
5. After finishing the network call of the moved Goroutine by "Network Poller", this goroutine will be moved back into the LRQ of P it was moved from.
6. The big win here is that, to execute network system calls, no extra Ms are needed.
7. The network poller has an OS Thread and it's handling an efficient event loop.

	SYNCHRONOUS SYSTEM CALLS
1. In the case when a Goroutine wants to make a syscall that can't be done asynchronously, the Goroutine making this syscall is going to block the M. It's unfortunate but there's no way to
	prevent this from happening.
2. One example of a syscall that can't be made asynchtonously is file-based syscalls ("I/O-Bound" workload)
3. When a Goroutine wants to make a synchronous syscall:
	1) The old M will be detached from the P with the blocking Goroutine
	2) To keep the busyness of the overall app the new M will be brought for that P:
		1) a new M created (if there's no free thread in the "Idle Thread Pool")
		2) a free M from the "Idle Thread Pool"
	3) After the execution of the blocking Goroutine, it'll be placed back into the LRQ of that P used

	WORK STEALING
1. Another aspect of the Go's scheduler is that it's a work-stealing scheduler. This helps in a few areas to keep scheduling efficient:
	1) The last thing we want is an M to move into the "Waiting" state because, once that happens, the OS will context-switch the M off the Core This means the P can't get any work done, even
	if there's a Goroutine in a runnable state, until M is context-switched back on a Core.
	2) The work-stealing also helps to balance the Goroutines across all the P's so the work is better distributed and getting dome more efficiently.
2. The rules of "Work Stealing":
	runtime.schedule() {
		// Only 1/61 of the time, check the GRQ for the Goroutines in Runnable state
		// If no Goroutines were found, check the LRQ
		// If no Goroutines were found, try to Steal from others Ps
		// If no Goroutines were found in other Ps, check the GRQ again
		// If not Goroutines were finally found, poll the network
	}
3. The good thing of "Work Stealing" is that it allows Ms to stay busy and not go idle. This work stealing is considered internally as spinning the M.
4. In order to minimize the hand-off, Go scheduler implements "spinning threads". Spinning threads consume a little extra CPU power but they minimize the preemption of the OS threads.
	A thread is spinning if:
		1) An M with a P assignment is looking for a runnable goroutine
		2) An M without a P assignment is looking for available Ps
		3) Scheduler also unparks an additional M and spins it when it's readying a Goroutine if there's an idle P and there are no other spinning threads

	There are at most GOMAXPROCS spinning Ms at any time.

	When a spinning thread finds work, it takes itself out of spinning state.
5. Idle threads with a P assignment don't block if there are idle Ms without a P assignment.

	When new Goroutines are created or an M is being blocked, scheduler ensures that there's at least one spinning M. This ensures that there are no runnable goroutines that can be otherwise
	running; and avoids excessive M blocking/ublocking

	CONCLUSION
1. Since all the context-switching is happening at the app levelm we don't lose the same ~12k instructions (on average) per context swtich that we were losing when using Threads.
2. In Go, those same context switches are costing us ~200 nanoseconds of ~2,4k instructions.
3. The scheduler also helps with gains on cache-line efficiencies
*/

// Ardan Labs

// https://habr.com/ru/companies/first/articles/581158/
/*
	GO/OS SCHEDULER
1. Planning of execution - is the hardest process. From the point of an OS's view the tiniest unit of work is OS Thread. The OS cannit regulate an amount of OS Threads are being created by a
user. The user can create these OS Threads by themselves.
2. The OS has to distribute these OS Thread throughout the cores. Since an amount of cores is limited by the processor, the OS needs to limit the time of work of these OS Threads.
3. Whem an Operating System decides that a specific Thread played enough and it has to chill, we have a bunch of work:
	1) We must switch some instructions
	2) Rewrite the CPU caches
4. The fewer amount of caches, the faster the program.
5. The Goroutines in Go language are working over CPU's Threads.
*/

// https://habr.com/ru/companies/first/articles/582144/
/*
1. Terms of GMP model:
	1) G - Goroutine that we will launch
	2) M - Machine, a thread we will run onto the CPU and complete our tasks on
	3) P - Processor that will perform our work
2. When we run our program the Go's environment reads the value of GOMAXPROCS (the amount of Ps) in order to find out the number of virtual processors (cores). It happens in the function called
schedinit, that is responsible for launching of tasks scheduler. The comment for this function is:
	// The bootstrap sequence is:
	//	call osinit
	//	call schedinit
	//	make & queue new G
	//	call runtime mstart

	//	The new G calls runtime main

	During the call of Go's scheduler we also creates our program.
3. Before the start of the program we complete the sched.maxmcount = 10_000, limiting the number of Ms (OS threads).
4. Before rewriting the maximum number of Ms (OS Thread) with the referening to GOMAXPROCS, we refer to the OS to find out the available number of them.
5. If we rewrite the GOMAXPROCS environment variable with the higher number than our system can support, we could dramatically descrease the performance of the overall app.
6. After that we find the main Goroutine an attach it to the first OS Thread run by the OS.
7. M - is the stream of the execution. M - is that, the Operation Ssytem fights for as we saw before. The fewer amount of OS Threads, the easier the process of planning of tasks of the OS.
	At the too low amount of Ms we have many tasks undone.
8. If we don't have the work, we could park our Ms.

	PARKING AND LAUNCHING OF Ms
1. We should keep the balance between of launched and parked Ms (OS Threads) in order to use all the resources of the system. As we said before:
	1) Many Ms result in scheduler chaos.
	2) Tiny amount of Ms results in the more work undone.
2. The scheduler itself is distributed. The queues of work ("Runnable" Goroutines) are saved for each processor (core). So we cannot find out the sequence of Goroutines' execution.
3. The M (OS Thread) is spinning if:
	1) The M has finished all the Goroutines in the LRQ and is searching for a new one
	2) The GRQ is empty too
	3) The timers queue is empty too. (If we run time.Sleep(...) in our Goroutine, it goes for a nap during the time is going. After the time passes, this Goroutine winds up in the special queue
	of execution, where the same Goroutines are)
	4) There are no Goroutines in "Network Poller" Thread.

	The spinning thread checks these queues and if it finds something for an execution, it starts the work. If it has found nothing, it's parked and thrown into "Idle Thread Pool".

	All the new OS Threads created are in the "spinning" mode.
4. Creation of Ms (OS Threads)
	1) If there are new Goroutines, we run a new M (OS Thread) if there are idle Ps and no spinning Ms.
	2) If we have a spinning M, we just pass the work (Gs) to this M. This M also checks the presence of others spinning Ms. If it hasn't found any work, it launches a new M in the "spinning"
	state.
*/

// https://habr.com/ru/articles/804145/
/*
	BASICS
1. Go doesn't give us the access to M (OS Threads)
2. Concurrency is about of execution of tasks simultaneously, but not necessarily at the same moment.
3. Parallelism - Some tasks are executed at the same moment and likely using more than one CPU's kernels.

	GMP MODEL
1. Goroutine (G)
	Goroutine - is the tiniest unit of Go's execution. The Goroutine is like a lightweight OS Thread.
	The Goroutine is also the "g" structure. Once a Goroutine is created it's passed into the LRQ of the P (logical Processor) and this P itself passes this Goroutine into M (OS Thread) to
	execute.
2. States of Goroutine
	1) "Waiting" state
		In "Waiting" state the Goroutine is idle. For example: G is paused for an operation of channels/mutexes/atomics or it can be stopped by syscalls
	2) "Runnable"
		In "Runnable" state the Goroutine is ready to be executed, but it's not being executed yet. It's just waiting for execution on an M (OS Thread)
	3) "Running"
		In "Running" state the Goroutine is being executed on a Thread (M). It will resuwe until the work is done or until it's stopped by the scheduler or smth else.
3. Logical Processor (P)
	P - is just an abstraction that is also called logical processor. As default the amount of Ps is set to the number of available kernels on the host.
	Each P includes its own list of "Runnable" Goroutines, that is called Local Run Queue (LRQ). The limit of Goroutines in LRQ is 256 (1<<8) Goroutines.
		1) If the LRQ is filled fully (It's current amount of Goroutines is 256), the Goroutine created is put into GRQ (Global Run Queue)

	The amount of Ps calls the maximum amount of Gorotines that can be run simultaneosly.
4. OS Thread (M)
	The Go's maximum number of Ms is 10_000
	If all the threads are locked, the Go's scheduler creates a new OS Thread for a new Goroutine.
	The OS Threads are reusable instances because of the weight of delete/create ops over OS Threads.

	HOW GMP WORKS?
1. If the Goroutine performs syscall that takes too long (for example a read from file on the HDD), the M will be blocked and unpinned from P. The P will be pinned to:
	1) a free M in "Thread Pool" if any
	2) a new M created
2. Work Stealing
	When the M has finished its work, it will search the Goroutines in others' LRQs.
		0) LRQs picked randomly
		1) It'll try to do it 4 times
		2) A half of Goroutines will be stolen
	If no Gorotines have been found,
*/

// GopherCon

// https://www.youtube.com/watch?v=YHRO5WQGh0k
/*
1. When does Go's scheduler do its work?
	1) Goroutine creation
	2) Goroutine blocking
	3) Blocking syscall. The Thread itself blocks too!
2. Long Running Goroutines.
	There's the "sysmon" Goroutine that has its own OS Thread (M). It's used to preempt long running Goroutines. It work in this way:
		1) It runs a background Thread called "sysmon" to detect long-running Goroutines
		2) If a Goroutines has been running >10 ms, it'll be unscheduled at the safe point
		3) The unscheduled goroutine will be put into GRQ (Global Runb Queue)
3. Thread Spinning
	1) Threads without work "spins" looking for work before parking
	2) They check GRQ, poll the network, attempt to run GC tasks, and work-steal.

	This burns CPU cycles, but maximally leverages available parallelism
4. Ps and Runqueues
	1) The per-core runqueues are stored in heap-allocated "p" struct
	2) "p" structure stores other resources a thread needs to run Goroutines too, like a memory cache.
	3) A thread claims "p" to run Goroutines and the entire "p" is haded-off when this M is blocked. This hand-off is taken care of by the "sysmon" too.
*/

// GopherCon

// Oleg Kozyrev
/*
1. The internals of Goroutines LRQ
	The LRQ is LIFO + FIFO
		1) If there's a Goroutine in LIFO, it'll be put first and executed. After it the elems will be picked and run with FIFO principle.
		2) If there are elements only in FIFO queue, they'll be picked with FIFO and executed respectively

	The LRQ has the limit on number of elems. This limit is 256 elems.
*/

// Oleg Kozyrev

// https://www.youtube.com/watch?v=P2Tzdg8n9hw
/*
	GO SCHEDULER
1. Web-services work with I/O workload mostly and this work requires context switching.
2. Using threads is expensive.
	1) The OS Thread has the weight of 2-8 MB
	2) Switching between contexts of threads takes place on OS kernel level
3. Thread is just a context to provide some code. To launch some code we need:
	0) Registers
	1) Instruction pointer (IP) (Program Counter)
	2) Stack
4. To switch from thread to thread we need:
	1) Take off the registers of the old one and put them on stack.
	2) Get the registers of the new thread and upload them on the CPU
	3) Launch the code
In this situation we discard the use of CPU caches, so the caches are invalidated.
5. The Goroutine (G) is just the superstructure over OS Threads.
	1) It's created inside our application (userspace).
	2) It's just a lightweight thread
	3) The kernel of OS doesn't know about Goroutines.
	4) It's start from 2 KB, it dumps the workload of DRAM
6. The idea of Goroutines is
	1) A programmer doesn't work with OS threads directly
	2) We pretend that there are only Goroutines
	3) The complexity of goroutines' distribution is abstracted
7. The OS Thread (M) just performs the code of Goroutine.
8. The OS scheduler is preemptive.
9. The Go's scheduler is implicitly cooperative
	1) The compiler will add some points into the source code to provide context-switching at runtime.
10. Working with goroutines we use "N:M model"
	1) We perform N Goroutines on M OS threads
11. Local queues of Goroutines.
	1) Each Local Ru Queue (LRQ) will be bound to a specific OS Thread (M)
	2) The OS thread will reference only to its own local run queue
	3) This way eliminates work with synchronization primitives
	4) The max size of the local queue of Goroutines is 256 (1 << 8) Goroutines
12. Processor (P) - is the abstraction that tells that we have the specific resources of execution (Goroutines):
	1) Rights of execution
	2) Resources of execution (queues of Goroutines)
	3) An OS Thread belongs to only one P (Processor)
13. GMP Model.
	1) G - Goroutine (What we execute)
	2) M - Machine (OS Thread) (Where we execute)
	3) P - Processor (Rights and resources of execution)
14. What the place we should put a new Goroutine created into?
	1) The goroutine created will be put into local queue
15. Work stealing.
	What we should do when the local queue is empty?
		1) We should try to steal some goroutines from other Ps (Other local queues)
		2) In this case we need synchronization (lock-free synchronization (atomics))
		3) A local queue will be picked randomly
		4) We will try to steal 4 times
		5) A half of goroutines will be stealed from other P
16. Syscalls
	Syscalls work with the kernel of OS. In this stage the OS thread is blocked fully. It turns into the "Waiting" stage. The OS scheduler removes this OS Thread from the kernel and
	puts a new one istead.
		1) When a syscall is performed, we unpin the OS Thread (M) from its P (Local queue)
		2) Create a new OS hread and put it into the P
		3) We pass all the goroutines left of this M created
	A bunch of these OS Threads is called "Thread Pool" and is used to keep the OS Threads that are in syscalls

	The Go's scheduler will unpin the OS Thread (M) from the processor (P) when it understands that this OS Thread will be blocked by syscall for a long time. In other cases it'll allow this OS Thread to be blocked and not to be unpinned.

	In the cases when the OS Thread (M) is in syscall and not unppinned but can execute too long it'll be evicted from P by "sysmon".

	The work of "sysmon" thread:
		1) "sysmon" checks periodically the blocked OS thread that is on a P and in syscall
		2) If this OS thread is in syscall too long, it'll be evicted by sysmon
		3) A new OS Thread will be pinned to this P
		4) Sysmon is just another OS thread
17. Where a goroutine should be put into after syscall?
	1) If the queue has enough place and this OS thread isn't performing another syscall, put this goroutine into the local queue of this P
	2) If the queue of this P is full (256 goroutines are already placed) or this OS thread is still performing syscall, we should search a new LRQ (local run queue)
	3) If there are no free queue, put this goroutine into the GTQ (global run queue)
18. When we work with blocking calls, we'll use multiplexer (netpoller). It's used to catch signals of events (event loop). When there's a blocking call, we'll send this call to the netpoller.

19. Queues
	1. Global queue. What should we do wneh there's filled Global queue?
		runtime schedule() {
			// only 1/61 of the time. check the global run queue for a G
			// if not found, check the local queue
			// if not found, try to steal from other Ps (4 retries, picked randomly)
			// if not found, check the GRQ
			// if not found, poll the network
		}
21. Goroutine has three states
	1) Executing (is performing)
	2) Runnable (is ready to be performed)
	3) Waiting (is blocked on mutex/syscall)
20. Synchronizations. When a goroutine is locked we just evict it from OS thread and put another one from its Ps queue. The evicted goroutine winds up ind "Waiting Queue".
	1) When a goroutine is in "Waiting Queue" and waiting more than 1 ms, it'll be picked as the next to use the mutex.
21. Tight Cycles (Adding preemptive part)
	1) If there's a goroutine that is performing more than 10 ms, it'll be evicted from OS Thread. The goroutine will reach the safe pointer before it is evicted.
22. When goroutines switching happens:
	1) morestack() -> newstack()
	2) runtime.Gosched()
	3) locks
	4) networks I/O
	5) syscalls
	6) Garbage Collection
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

// Google Dev (Dmitriy Vyukov)
/*
	THE LACKS OF THREADS
1. Memory consumption at least 32 KB
2. Performance (syscalls)
3. No infinite stacks

	M:N THREADING
1. OS doesn't know anything about Goroutines, it works only with OS Threads

	BLOCKED GOROUTINES
1. The channel object has its own "Wait Queue"
2. When Goroutine is blocked on a channel it's listed in this channel.
3. When this blocked Goroutine is unblocked, it's returned to GRQ
4. The same mechanism is applied to:
	1) Mutexes
	2) Timers
	3) Network I/O

	SYSTEM CALLS
1. When the execution goes into the kernel, we completely lose visibility.
2. We only see when our execution enters/exits the syscall.
3. When a Thread enters a syscall, it creates/wakes up another thread to resume work
4. #Threads is > than #Cores

	M:P:N THREADING
1. P (Processor) - is the resources required to run Go code.
2. The number of Ps (Processors) is equal to the number of cores.

	SYSTEM CALLS HANDLING (HANDOFF)
1. When a Thread goes into system call it wakes another Thread and then it passes the Processor object to this Thread woken.
2. This thread woken starts to take Goroutines from the LRQ of this Processor and perform it.

	FAIRNESS
1. Fairness is a property of the Scheduler that simply says:
	If a goroutine in the "Runnable" state, it'll run eventually
2. If we don't have fairness we get:
	1) bad tail latencies
	2) livelocks
	3) pathological behaviors
3. FIFO is bad because we run the latest thing in Goroutine Run Queue.
4. LIFO is good because we run the newest thing, so we use the "hot"-caches.
5. "Infinite Goroutines"
	When a Gorotine in an infinite loop, the other Goroutines are starved by it.
	The solution is to preempt this Goroutine in the "tight-loop" after 10 milliseconds.

	This 10 milliseconds is soft limit of time.
	The Goroutine preempted goes to the tail of Global Run Queue
6. Local Run Queue (LRQ)
	The Local Run Queue is built with FIFO queue (Double Linked List) + 1-element LIFO buffer

	The 1-element buffer is neeeded to provide better locality of caches

	The Goroutine in LIFO slot cannot be stolen. This restriction lasts for ~3 microseconds.

7. LRQ Starvation
	1-element LIFO buffer is enough to get starvation

	The situation is: two goroutines response each other, so each of them wind up in the 1-elem LIFO buffer.
		The Gs in FIFO are never executed, so they're starved.

	To resolve this problem we have Time Slice Inheritance:
		The goroutine has the time slice (10 milliseconds). When we get a Goroutine from the 1-elem LIFO buffer this Goroutine inherits the time slice.

	Inherit Time Slice -> Looks like infinite loop -> Preemption (~10 milliseconds)
8. GRQ Starvation & The choice of the number 61:
	1) Not too small. In other case we'll poll the GRQ too frequently and get lock-contention problem.
	2) Not too large. In other case we'll not resolve the starvation problem
	3) Prime
9. Network Poller Starvation.
	The Solution: background Thread polls network occasionally if it wasn't polled by main-worker Threads

	THREAD STACK
1. Function frame
	1) Local variables
	2) Return address
	3) Previous frame pointer
2. Thread Stack is the place when we store Function Frames.
	1) There's the Stack Pointer that points where the stack top is.
3. Stack implementations
	1) Stack is implemented as the large range of virtual address space.
	2) The size of this virtual range is about 1-8 MB
	3) It consists of physical pages (4KB). These pages are allocated lazily.
4. Types of pages:
	1) Preallocated stack pages
	2) Empty stack pages
	3) Protected stack pages (are used to define stack overwflows)

	GOROUTINE STACK
0. Some properties:
	1) The stack size is limited to 1 GB (it's "infinite enough")
	2) We do the check of Stack overflow at runtime working with Goroutines
	3) If we don't have enough space, we call morestack()

	4) If function returns and it happens in the Nth frame of stack, the frame will just be removed (it's applied to Split Stacks)
5. The benefits of "Split Stack":
	1) We can run 10^6 Goroutines
	2) Works on 32-bits systems
	3) Better granularity
	4) Cheap "page-out"
	5) Huge pages
	The costs of this type of stack:
		1) A normal function call: ~2 nanosecs
		2) Stack split:	~60 nanosecs
7. Growable stacks.
	Go has switched to "Growable stack"
	The idea of "Growable stack" is simple: When we reach the limit of current stack:
		1) We allocate the stack of size x2 of Old Stack
		2) Copy all the stack contents to the new stack
		3) Deallocate the old stack
8. The differences between "Split" and "Growable" stacks
	Split stack:
		1) O(1) cost per function call
		2) Repeated
	Worst case: stack split in hot loop

	Growable stack
		1) O(N) cost per function call
		2) Amortized
	Worst case: Grow stack for short Goroutine

	PREEMPTION
1. What is Preemption?
	Preemption - is asynchronously asking a goroutine to yield.
2. We need Preemption for fairness, multiplexing, GC calls and crashes
*/

// Google Dev (Dmitriy Vyukov)

// Vicki Niu
/*
	INTERNALS OF GOROUTINE
1. What is a Goroutine?
	type g struct {
	stack 				stack 	// stack offsets
	m					*m 		// current m (OS Thread)
	sched				gobuf	// saves context of g for recovering after context switch
	atomicstatus		uint32 	// current status of g (waiting, runnable, running)

								// activeStackChans indicates that there are unlocked channels
								// pointing at this Goroutine's stack. If true,
								// stack copying needs to acquire channel locks to protect these
								// areas of the stack.
	activateStackChans 	bool

								// parkingOnChan indicates that the goroutine is about to
								// park on a channel or chanrecv. Used to signal an unsafe point
								// for stack shrinking. It's a boolean value, but it's updated atomically
	parkingOnChan 		uint8
	}
2. Goroutine stack
	1) All Gs are initially allocated with 2 KB stack
	2) Each go function has a small preamble, which calls morestack if it runs out of memory. Then the runtime allocates a new memory segment with the doubled size and copies the old segnment
	and restarts the execution
	3) Effectively makes goroutines infinitely growable with efficient shrinking.
*/

// Vicki Niu

func main() {

	A()

	// LFIFO()
}

func A() {
	runtime.GOMAXPROCS(1)

	go B()

	time.Sleep(1 * time.Second)
}

func B() {
	for {
	}
}

func NumberOfVirtualCPUs() {
	fmt.Println(runtime.NumCPU())
}

func LFIFO() {
	runtime.GOMAXPROCS(1)
	var (
		wg sync.WaitGroup
	)

	for i := 0; i < 5; i++ {

		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			fmt.Println(i)
		}(i)
	}

	wg.Wait()
}

// Function Frame is:
/*
	1) Local variables
	2) Return address
	3) Previous frame pointer
*/
func Foo() {}
