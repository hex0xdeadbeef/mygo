package main

/*
	Understanding the Golang Goroutine Scheduler GPM Model
*/

/*
	1.1 The origin of the Golang "Scheduler"
	All software operates on top of an operating system, and the true computational engine behind these software applications is the CPU.
*/

/*
	Scheduler Requirements in the Era of Multi-Process/Multi-Thread

	A time slice, also referred to as a `quantum` or `processor slice`, is a microscopically allocated segment of CPU time that a time-sharing operating systems assigns to each running process (or preemptive kernels, the time from when a process begins running until it gets preempted). Modern OS like Windows, Linux, and macOS allow multiple processes to run simultaneously. For example, we can listen to music using a media player while browsing the internet and downloading files. In reality, even though a computer may have multiple CPUs, a single CPU cannot truly run multiple tasks simultaneously. In the context of a single CPU, these processes `appear to` run concurrently, but they actually execute in a interleaved fashion. This interleaving is made possible by the allocation of short time slices (typically ranging from 5 ms to 800 ms on Linux), which are too short for users to perceive.

	The OS system's scheduler allocates time slices to each process. (!!!) Initially, the kernel assigns an equal time slice to each process. Subsequently, processes take turns executing during their allotted time slices. When all processes have exhausted their time slices, the kernel recalculates and assigns time slices to each process and this process repeats.
*/

/*
	1.1.3 Coroutine Improve CPU Utilization

	Creating a thread for each task would result in an overwhelming number of threads running simultaneously, leading to not only high context switching overhead but also substantial memory consumption (process virtual memory can consume up to 4 GB on 32-bit OS, and threads typically require around 4 MB each). This profileferation of processes or threads introduced new problems:
		- High memory consumption
		- Significant CPU overhead in scheduling

	Engineers discovered that a thread could be divided into two states
	- Kernel mode
	- User mode
	User-mode threads essentially mimic kernel-mode threads in user space, with the aim of achieving greater lightweightness (lower memory consumption, reduced isolation, faster scheduling) and higher controllability (ability to manage the scheduler independently). In user mode, everything is visible to kernel mode, but from the kernel's perspective, user-mode threads are merely a set of memory data.

	A user-mode thread must be bound to a kernel-mode thread, but (!!!) the CPU is unaware of the existence of user-mode threads; it only knows that it's running a kernel-mode thread (Linux's PCB - Process Control Block).

	If we further classify threads, kernel threads remain as `threads`, while user threads are referred to as `coroutines` (Co-routine). Threads at the operating system level are what we commonly know as kernel-mode threads, while user-mode threads come in various forms. As long as they can execute multiple tasks on the same kernel-mode thread, they are considered user-mode threads. Examples include coroutines, Golang's Goroutines, C# Tasks, and more.

	Since one coroutine can be bound to one thread, is it possible for multiple coroutines to be bound to one or more threads? There are three mapping relationships between coroutines and threads: N:1, 1:1, and M:N.

	1. `N:1` relationship
	N coroutines are bound to 1 thread. The advantage of this approach is that coroutines perform context switches within user-mode threads without transitioning to kernel mode, making these switches extremely lightweight and fast. However, the drawback is evident: all coroutines of one process are bound to a single thread.

	Challenge faced in the N:1 relationship:
		- Inability to harness the hardware's multi-core acceleration capabilities for a particular program
		- When one coroutine is blocked, it leads to thread blocking, causing all other coroutines withing the same process to become non-executable, thus resulting in a lack of concurrency.

	2. 1:1 Relationship
	In a 1 coroutine to 1 thread relationship, it's the simplest to implement. Coroutine scheduling
	is entirely handled by the CPU. While it doesn't suffer the drawbacks of N:1 relationship, the cost and overhead of creating, deleting, and switching coroutines are all borne by the CPU, making it slightly more expensive.

	3. M:N Relationship
	In an M coroutines to N threads relationship, it combines aspects of both the N:1 and 1:1 models, overcoming the drawbacks of the two aforementioned approaches. However, it's the most complex to implement, as illustrated. Multiple coroutines are hosted on the same scheduler, while the scheduler's downstream resources comprise multiple CPU cores. It's important to note that coroutines and threads differ in their nature:
		- Thread scheduling by the CPU is preemptive
		- Coroutines scheduling in user mode is cooperative; once coroutine yields the CPU before one is executed.
	Consequently, the design of an intermediate-level scheduler in the M:N model becomes crucial, and improving the binding relationship and execution efficiency between threads and coroutines becomes a diverse set of languages' primary goals when designing schedulers.
*/

/*
	1.1.4 Go Language's Goroutines
	Goroutines are derived from the concept of coroutines, allowing a set of reusable functions to run on a set of threads. Even if one Goroutine is blocked, the runtime scheduler can still move other Goroutines of the same thread to other runnable threads. Most importantly, programmers are shielded from these low-level details, which simplifies programming and offers easier concurrency.
*/

/*
	1.1.5 Deprecated Goroutine Scheduler
	All coroutines denoted as G, were placed in a GRQ. Outside the GRQ, where multiple M's shared resources, a lock was added for syncronization and mutual exclusion purposes.

	In order for an M to execute or return a G, it must access the GRQ. Since there are multiple Ms, which implies threads accessing the same resource, locking is required to ensure mutual exclusion and syncronization. Therefore, the GRQ is protected by a mutex lock.

	Several shortcomings of the old scheduler can be identified:
		- Creating, destroying, and scheduling G all require each M to acquire the lock, leading to intense lock contention
		- Transferring G between Ms results in delays and additional system overhead. For example, when G contains the creation of a new coroutine, and M creates G', in order to continue executing G, G' must be handed over to M2 (if allocated). This also leads to poor locality because G' and G are related, and it's preferable to execute them on the same M rather than a different M2.
		- Syscalls (CPU switches between M's) resulted in frequent thread blocking and ublocking ops, leading to increased system overhead.
*/

/*
	In Go, threads serve as the entities that execute Goroutines, and the scheduler's role is to allocate runnable Goroutines to worker threads.

	In the GPM model, there are several important concepts:

	1) GRQ: It stores waiting-to-run Gs. The GRQ is a shared resource accessible by any P, making it a global resource within the model. Consequently, mutual exclusion mechanisms are used for read and write ops on this queue.

	2) LRQ: Similar to the GRQ, it contains waiting-to-run Gs, but with a limited capacity, typically not exceeding 256 Goroutines. (!!!) When a new G' is created, it's initially added to the P's LRQ. (!!!) If the LRQ becomes full, half of its Goroutines are moved to the GRQ.

	3) P List: All (!!!) Ps are created at program startup and are stored in an array, with the maximum number of Ps typically set to GOMAXPROCS (a configurable value).

	4) M: Threads (M) request tasks by obtaining a P. They retrieve Gs from the P's LRQ. When the LRQ is empty, an M may attempt to fetch a batch of Gs from the GRQ and place them in its P's LRQ or steal a half from another P's LRQ. The M then executes a G, and after completion, it fetches the next G from its P, repeating this process.

	The Gs scheduler works in conjunction with the OS scheduler through the M, where each M represents a kernel thread. The OS scheduler allocates kernel threads to CPU cores for execution.

	1. Regarding the number of P's and M's
		1) The number of Ps is determined by the environment variable $GOMAXPROCS or the runtime method GOMAXPROCS(). This means that at any given moment during program execution, there are only $GOMAXPROCS Goroutines running simultaneosly.

		2) (!!!) The number of Ms is constrained by Go's own limits. When a Go program starts, it sets the maximum number of Ms, typically defaulting to 10,000. However, it's challenging for the kernel to support such a high number of threads, so this limit can often be disregarded. The `runtime/debug` package provides the SetMaxThreads() function to adjust the max number of Ms, (!!!) When an M becomes blocked, new Ms can be created.

		There's no strict relationship between the number of Ms and Ps. When an M becomes blocked, a P may either create a new M or switch to another available M. So, even of the
		default number of Ps is 1, it's still possible to create multiple Ms when needed.

	2. Regarding the Creation Timing of Ps and Ms
		1) Ps are created when the maximum number of Ps, denotes as `N`, is determined. The runtime system generates `N` Ps accordingly.
		2) Ms are created when there aren't enough Ms available to associate with Ps and execute the runnable Goroutines within them. For example, if all Ms are blocked, and there are still many ready tasks in the P, the system will search for idle Ms, and if none are available, it'll create new Ms.
*/

/*
	1.2.2 Scheduler Design Strategies

	Strategy One: Thread Reuse
	The aim is to avoid frequent thread creation and destruction and instead focus on reusing threads.

		1) Work Stealing Mechanism
			(!!!) When the current thread has no runnable Goroutines to execute, it attempts to steal Gs from other threads bound to the same P, rather than terminating the thread.

			(!!!) It's important to note that the act of stealing Gs is initiated by a P, not an M, because the number of P's is fixed. If an M cannot acquire a P, it means that the M doesn't have a LRQ, let alone the capacity to steal from other Ps' queues.

		2) Hand Off Mechanism
			When the current thread blocks due to a syscall while executing a Goroutine, the thread releases the associated P and transfers it to another available thread for execution. If in the GPM combination of M1, G1 is currently being scheduled and becomes blocked, it triggers the design mechanism for handoff. In the GPM model, to maximize the utilization of M and Ps performance, a single blocked G1 shouldn't indefinitely delay the work. So, in such situations, the handoff mechanism is designed to immeadiately release the P.

			Since M1 is currently executing the active G1, the entire program stack space is contained within M1. Therefore, M1 should enter a blocked state along with G1. However, the P that has been released needs to be associated with another M. In this case, it selects M3 (or creates a new one or awakens a sleeping one if there's no M3 available at the time) to establish a new association. The new P then continues its work, either accepting new G's or implementing the stealing mechanism from other queues.

	Strategy Two: Leveraging Parallelism
		GOMAXPROCS sets the number of Ps, with maximum of GOMAXPROCS threads distributed across multiple CPUs simultaneously. GOMAXPROCS also limits the degree of concurrency; for example, if GOMAXPROCS == `coresNumber`/2, it utilizes at most half of the CPU cores for parallelism.

	Strategy Three: Preemption
		In Coroutine, one must wait for a coroutine to voluntarily yield the CPU before the next coroutine can execute. (!!!) In Go, a Goroutine can occupy the CPU for a maximum of 10 ms, preventing other Goroutines from being starved of resources. This is a distinguishing feature of Goroutines compared to regular Coroutines.

	Strategy Four: GRQ (Global Run Queue)
		In the new scheduler, the GRQ stil exists but with weakened functionality. When an M attempts to steal Goroutines but cannot find any in the queues of others' Ps, it can retrieve Goroutines from the GRQ.
*/

/*
	1.2.3 go func() Scheduling Process
	1. Creating a goroutine using `go func()`

	2. There are two queues for storing G:
		- One is the LRQ of the scheduler P
		- Other is the GRQ
	The newly created G is initially stored in the P's LRQ. (!!!) If the LRQ is full, it'll be stored in the GRQ.

	3. G can only run within an M, and each M must hold a P, establishing a 1:1 relationship between M and P. M pops an executable G from the LRQ of P for execution. If P's LRQ is empty, M retrieves a G from the GRQ. If the GRQ is also empty, M attempts to steal an executable G from other MP combinations.

	4. The process of scheduling G for execution in an M follows a looping mechanism.

	5. When an M is executing a G, in a syscall or other blocking op occurs, the M becomes blocked. If there are other Gs currently executing, the runtime detaches the thread M from P and creates a new OS thread (reusing an available idle thread if possible) to serve this P.

	6. (!!!) When the M completes the syscall, the G will attempt to acquire an idle P for execution and be placed in the LRQ of that P. If acquiring a P isn't successful, the thread enters a sleep state and joins the pool of idle threads. (!!!) Meanwhile, the G is placed in the GRQ.
*/

/*
	1.2.4 Scheduler Lifecycle
	In the GPM model of the Golang scheduler, there are two distinct roles: M0 and G0

	M0 (zero-thread)
		- The main thread with the index 0 started after the program launches
		- (!!!) Located in the global command runtime.m0, requiring no allocation on the heap
		- Responsible for executing initialization ops and launching the first G
		- (!!!) After launching the first G, M0 functions like any other M

	G0 (zero-goroutine)
		- Each time an M is started, the first Goroutine created is G0
		- (!!!) G0 is solely dedicated to scheduling tasks
		- G0 doesn't point to any executable function
		- (!!!) Every M has its own G0
		- During scheduling or syscalls, M switches to G0, and scheduling is performed through G0
		- M0's G0 is placed in the global space.


	Let's make a breakdown of this code:

		package main

		import "fmt"

		func main() {
			fmt.Println("Hello world!")
		}

	The overall analysis is as follows:
		1) The runtime creates an initial thread M0 and Goroutine G0, associating the two.
		2) Scheduler initialization:
			- initializes M0
			- stacks
			- garbage collection
			- creates and initializes the P list consisting of GOMAXPROCS Ps.
		3) The main() func in the example code is `main.main`, and there's also a main() function in the runtime, `runtime.main`. After the code is compiled, `runtime.main` will call `main.main`. When the program starts, it'll create a Goroutine for `runtime.main` known as the Main Goroutine, and then add the Main Goroutine to the LRQ of P.
		4) Start m0, m0 has already bound P, it'll get G from P's LRQ, and get the Main Goroutine.
		5) G has a stack, M sets the running environment according to the stack information and scheduling information in G.
		6) M runs G
		7) G exits, and returns to M to get the runnable G, repeating this until main.main exits, runtime.main performs Defer and Panic processing, or calls `runtime.exit` to exit the program.

		The lifecycle of the scheduler almost fills the entire life of a Go program. All the preparations before the Goroutine of runtime.main are for the scheduler. The running of the Goroutine of runtime.main is the real start of the scheduler, which ends when runtime.main ends.
*/

/*
	1.3 Go Scheduler: Comprehensive Analysis of Scheduling Scenarios.

	1.3.1 Scenario 1: G1 creates G2
	This scenario primarily demonstrates the locality of GPM scheduling. P owns G1, and M1 starts running G1 after acquiring P.

	(!!!) When G1 uses go func() to create G2, for the sake of locality, G2 is given priority and is added to the local queue of P1.

	By default, if a G1 creates a G2, G2 is prioritized and placed in the LRQ of P1. This is because G1 and G2 share similar memory and stack information, making the cost of switching between them relatively small for M1 and P. This is a characteristic that locality aims to preserve.

	1.3.2 Scenario 2: G1 execution Completed
	Currently, M1 is bound to G1, and the LRQ of P1 contains G2. Meanwhile, M2 is bound to G4, and the LRQ of P2 contains G5.

	(!!!) After the completion of the execution of G1, the Goroutine running on M1 is switched to G0.

	(!!!) G0 is responsible for coordinating the context switch of Goroutines during scheduling. It retrieves G2 from the LRQ of P, switches from G0 to G2, and starts executing G2, achieving the reuse of thread M1.

	Scenario 2 primarily demonstrates the reusability of M in the GPM scheduling model. When a Goroutine (G) has completed its execution, regardless of the availability of others Ps, M will prioritize fetching a pending G from its bound P's local queue. In this case, the pending G is G2.

	1.3.3 Scenario 3: G2 Creates Excessive Gs
	Suppose each P's local queue can only store 4 Gs, and the LRQ of P1 is currently empty. (!!!) If G2 needs to create 6 Gs, then G2 can only the first 4 Gs (G3, G4, G5, G6) and place them in its current queue. The excess Gs (G7, G8), will not be added to the LRQ of P1. The detailed explanation of where the excess Gs are placed will be covered in Scenario 4.

	Scenario 3 illustrates that if a G creates an excessive number of Gs, the LRQ may become full. The excess Gs need to be arranged according to the logic explained in the Scneario 4.

	1.3.4 Scenario 4: G2 Locally Full, Creating G7
	(!!!) When G2 attempts to create G7 and discovers that P1's local queue is full, it triggers a load balancing algorithm. This algorithm moves the first half of the Gs from P1's LRQ, along with the newly created G7, to the GRQ.

	When these Gs are moved to the GRQ, their order becomes disrupted. Therefore, G3, G4, and G7 are added to the GRQ in a non-sequential manner.

	1.3.5 Scenario 5: P2 LRQ Not Full, G2 creating G8

	When G2 creates G8 and the LRQ of P1 isn't full, G8 will be added to the LRQ of P1.

	(!!!) Newly created G8 is given top priority and is placed into the LRQ first, adhering to the locality property. Since the LRQ already contains other Gs at its head, G8 will sequentially enter from the queue's tail.

	1.3.6 Scenario 6: Waking Up a Sleeping M
	In the GPM model, (!!!) when creating a new G, the running G will attempt to wake up other idle P and M combinations for execution. This implies that there might be surplus M instances from previous ops, and these M instances won't be immediately reclaimed by the OS. Instead, they're placed in a sleeping thread queue. (!!!) The condition triggering the awakening of an M from the sleeping queue is attempting to create a new G.

	The situation is: only P1 and M1 are bound, and P1 is executing G2. (!!!) When G2 attempts to create a new G, it triggers an attempt to fetch an M from the sleeping thread queue and attempts to bind it to a new P along with the local queue of P.

	G2 is waking up M2 from the sleeping thread queue.

	If M2 discovers available P resources, it'll be activated and bound to P2.

	At this point, if M2 is bound to P2, it needs to find other Gs. Therefore, if there are no normal Gs available for M2 and P2, G0 will be scheduled by P. If P's LRQ is empty and P is in G0, the combination of M2, P2 and G0 is referred to as a `spinning thread`. The `spinning thread` is constantly looking for Gs.

	One characteristic of a spinning thread is the continuous search for Gs. The process of finding Gs will lead to the flow of scenarios 7 and 8 to be introduced next.

	1.3.7 Scenario 7: M2 Awoken Fetches a Batch of Gs from the GRQ.
	M2 attempts to fetch a batch of Gs from the GRQ and places them in P2's LRQ. The number of Gs that M2 fetches from the GRQ follows the formula:

		N = min(len(GRQ)/GOMAXPROCS + 1, cap(LRQ)/2)

	At least 1 G is fetched from the GRQ each time, but it avoids moving too many Gs from the GRQ to the P queue at once, leaving some for other Ps. This aims to achieve load balancing from the GRQ to the P LRQ.

	If M2 is a spinning thread state at this moment, the GRQ has a count of 3, and the LRQ capacity of P2 is 4. Accoring to the load balancing formula, the number of Gs fetched from the GRQ at once is 1. M2 will then fetch G3 from the head of the GRQ and add it to P2's LRQ.

	Once G3 is added to the LRQ, it needs to be scheduled by G0. G0 is then replaced by G3, and at the same time, G7 and G4 from the GRQ will move to the head of the GRQ.

	Once G3 is scheduled, M2 is no longer a spinning thread.

	1.3.8 Scenario 8: M2 Steals from M1
	Assuming G2 has been running on M1 all along, after two rounds, M2 has succesfully acquired G7 and G4 from the GRQ and added them to the LRQ of P2, completing their execution. Both the GRQ and the LRQ of P2 will be empty.

	M2 will once again be in a state of scheduling G0. In this state, a spinning MP combination continuously looks for Gs that can be scheduled. Otherwise, waiting in idle would mean wasting thread resources. Since the GRQ is empty, M2 will perform stealing: it steals half of the Gs from another P that has Gs and adds them to its own P's LRQ. P2 takes half of the Gs from the tail of P1's LRQ.

	In this case, half means only the G8. Therefore, M2 will attempt to steal G8, placing it in P2's LRQ.

	The subsequent process is similar to the previous scenario. G0 will schedule G8, and as a result, M2 will no longer be in a spinning thread state.

	1.3.9 Scene 9: Maximum of Spinning Threads
	M1's LRQ with G5 and G6 has been taken by M2 and executed, and currently, M1 and M2 are running G2 and G8, respectively. In the GPM model, (!!!) the GOMAXPROCS variables determines the number of Ps. The number of Ps is fixed in the GPM model and cannot be dynamically added or removed.

	Assuming GOMAXPROCS=4, the number of Ps is also 4. At this point, P3 and P4 are bound to M3 and M4.

	M3 and M4 don't have any Gs to run, so they are currently bound to their respective G0, and both M3 and M4 are in a spinning thread state, continually searching for Gs.

	Why allow M3 and M4 to spin? Spinning is essentially running, and when a thread is running without executing a G, it becomes a waste of CPU. Why not destroy the context to save CPU resources? Because creating and destroying CPU also incurs a time cost, and the philosophy of the Golang scheduler GPM model is to have an M ready to run a new G immediately upons its creation. If we were to destroy and recreate, it would inroduce delays and reduce efficiency. Of course, the system also considers that having too many spinning threads wastes CPU, so there's a maximum limit of GOMAXPROCS spinning threads in the system. Any excess idle threads will be put to sleep. If there's an M5 currently running but no available P to bind with, M5 will be added to the sleeping thread queue.

	1.3.10 Scenario 10: Blocking Syscall in G
	Assuming that, in addition to M3 and M4 being spinning threads there are also idle threads M5 and M6 (not bound to any P, noting that there can be a maximum of 4 P's, so the number of P's should always be M >= P). G8 creates G9.

	If G8 makes a blocking syscall at this point, M2 and P2 are immediately unbound. P2 performs the following checks: If P2's LRQ has a G or the GRQ has a G that needs execution, and there's an idle M available, P2 will promptly wake up one M to bing with (if there are no sleeping Ms, a new one will be created for binding). Otherwise, P2 will be added to the list of idle Ps, waiting for an M to acquire an available P.

	In this scenario, P2's LRQ has G9, and it can be bound to another idle thread, M5. During the blocking period of G8, M2's resources are temporary occupied by G8.

	Scenario 11: G undergoes a non-blocking syscall.
	Continuing from Scenario 10, G8 creates G9. Assuming that the previous syscall made by G8 has ended, it's currently in a non-blocking state.

	In Scenario 10, even though M2 and P2 become unbound, M2 remembers P2. Then when G8 and M2 exit the syscall, M2 attempts to reacquire the previously bound P2. If P2 is available, M2 will rebind with P2. If it cannot be obtained, M2 will handle it in other ways.

	In the current example, at this moment, P2 is bound to M5, so M2 naturally cannot acquire P2. M2 will attempt to find an available P from the idle P queue.

	If there are still no idle Ps in the `Idle Ps queue`, M2 effectively cannot find an available P to bind to. G8 will be marked as runnable and added to the GRQ. M2 without a bound P, enters a sleeping state (long-term sleep may lead to GC recycling and destruction).
*/
