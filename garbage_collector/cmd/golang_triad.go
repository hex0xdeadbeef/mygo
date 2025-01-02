package main

/*
	COMPREHENSIVE ANALYSIS OF GO'S HYBRID WRITE BARRIER GC
*/

/*
	THE DEFINITION OF GARBAGE COLLECTOR
	Garbage Collection (GC) is an automatic memory management mechanism provided in programming languages, capable of automatically releasing unnecessary memory objects and relinquishing storage resources without the need for manual intervention by programmers. The support for GC mechanisms is widespread in modern programming languages, and the performance of GC is often considered one of the comparative metrcis between different languages.

	Throughout its evolution in GC, Golang has undegone several significant changes. Up to Go 1.8, the modifications to Golang's GC have been substantial. The specific major changes are as follows:

		- Mark and Sweep method before Go V1.3
		- Three-color concurrent marking method in Go V1.5
		- "Strong-Weak" three color invariant, insertion barrier, and deletion barrier in Go V1.5
		- Hybrid write barrier mechanism in Go 1.8

	This chapter will introduce the algorithmic models of each GC method, their respective advantages and disadvantages, and how Golang's GC has evolved step by step from logical requirements to the hybrid barrier mode.
*/

/*
	1. MARK-SWEEP ALGORITHM IN GO UNTIL V1.3
	Next, let's take a look at the conventional mark-sweep algorithm primarily used in Golang before version 1.3. This algorithm consists of two main steps:
		- Mark phase
		- Sweep phase
	The process iq quite straightforward:
		identify memory data to be cleared and then sweep them all at once.
*/

/*
	1.1 Detailed process of Mark and Sweep Algorithm
	1) In the first step, pause the program's business logic, categorize objects into reachable and unreachable, and then mark them accordingly.
	2) In the second step, marking begins. The program identifies all the objects it can reach an mark them.
	3) In the third step, after marking is completed, the unmarked objects are then cleared.

	The operation is quiet simple, but there is one imporant point to note. (!!!) The Mark and Sweep algorithm requires the program to pause during execution, known as STW (Stop the World). (!!!) During STW, the CPU doesn't execute user code and is entirely dedicated to garbage collection. This process has a significant impact, making STW one of the biggest challenges and focal points for optimization in some garbage collection mechanisms. Therefore, during the third step, the program will temporarily halt ant work and wait for the garbage collection to complete.

	4) In the fourth step, the pause stops, allowing the program to resume running. This process is than repeated in a loop until the program's lifecycle ends. This describes the Mark and Sweep algoritm.
*/

/*
	1.2 Disadvantages of the Mark and Sweep Algorithm
	The Mark and Sweep algorithm is starightforward and clear in its process, but it has some significant issues.
	- The first STW (Stop the World), which causes the program to pause and may lead to noticable stuttering, and important concern.
	- The second issue is that marking requires scanning the entire heap.
	- The third issue us that clearing data can create heap fragmentation.

	Prior to version 1.3, Go implemented garbage collection using this approach. The basic GC process involved first initiating an STW pause, then performing marking, followed by data collection, and finally stopping the STW.

	(!!!) From the figure 3, it can been seen that the entire GC time is encapsulared within the STW period, which seems to cause the program pause to be too long, affecting the program's runtime performance. Therefore, in Go version 1.3, a simple optimization was made to advance the STW step, reducing the duration of the STW pause.

	The figure 4 mainly illustrates that the STW step was advanced because during the Sweep phase, an STW pause isn't necessary. These objects are already uncreachable, so there won't be any conflicts or issues during garbage collection. However, regardless of the optimiztion, Go version 1.3 still faces a significant issue: the Mark And Sweep algorithm pauses the entire program.
*/

/*
	2. TRICOLOR MARKING IN GO 1.5
	(!!!) In Golang, garbage collection mainly employs tricolor marking algorithm. The GC process can run concurrently with other user goroutines, but a certain amount of STW time is still required. (!!!) The tricolor marking algorithm determines which objects need to be cleared through three marking phases. This section will introduce the specific process.
*/

/*
	2.1 THE PROCESS OF TRICOLOR MARKING
	1) First, every newly created object is initially marked as `white`.
	2) In the second step, at the start of each GC collection, the process begins by traversing all objects from the root nodes. Objects encountered during this traversal are moved from the white set to the `gray` set. (!!!) it's important to note that this traversal is a single-level, non-recursive process. It traverses one level of objects reachable from the program. If the currently reachable objects are object 1 and object 4, then by the end of this round of traversal, objects 1 and 4 will be marked as `gray`, and the `gray` marking table will include these two objects.
	3) In the third step, the `gray` set is traversed. Object referenced by `gray` objects are moved from `white` set to the `gray` set. Afterward, the original `gray` objects are moved to the `black` set. This traversal only scans `gray` objects. The first layer of objects reachable from the `gray` objects is moved from the `white` set to the `gray` set, such as objects 2 and 7. Meanwhile, the previously `gray` objects 1 and 4 are marked as `black` and moved from the `gray` marking table to the `black` marking table.
	4) In the fourth step, repeat the third step until there are no objects left in the `gray` set.

	Once all reachable objects have been traversed, the `gray` marking table will no longer contain any `gray` objects. At this point, all memory data will only have two colors: `black` and `white`. `Black` objects are the ones that are logically reachable (needed) by the program. These are the valid and useful data that support the normal operation of the program and cannot be deleted. `White` objects, on the other hand, are uncreachable by the current program logic and are considered garbage data in memory, which needs to be cleared.

	5) In the fifth step, all objects in the `white` marking table are reclaimed, which means collecting the garbage. All `white` objects are deleted and reclaimed, leaving only the `black` objects, which are the necessary dependencies.

	This is the tricolor concurrent marking algorithm. It clearly demonstrates the characteristics of the tricolor marking. However, many concurrent processes might be scanned, and the memory of these concurrent processes might be interdependent. (!!!) To ensure data safety during the GC process an STW is added before starting the tricolor marking. After scanning and determining the black and white objects, the STW if lifted. However, it's evident that the performance of such GC scan is quite low.
*/

/*
	2.2 TRICOLOR MARKING WITHOUTH STW
	If there's no STW, performance issues would be eliminated. But what would happen if the tricolor marking didn't include SWT?

	Based on the aforementioned tricolor concurrent marking algorithm, it's indeed dependent on STW. Without pausing the program, any changes in object reference relationships during the marking phase could affect the correctness of the marking results. Consider a scenario where the tricolor marking process doesn't use STW and what might occur.

	Assume the initial state has already undergone the first round of scanning. Currently, objects 1 and 4 are black, object 2 and 7 are gray, and the rest are white. Object 2 is pointing to object 3 through a pointer `p`

	(!!!) Now, if the tricolor marking process doesn't initiate STW, any object may undergo read or write ops during the GC scan. As shown in figure 1, before object 2 has been scanned, the alredy black-marked object 4 creates a pointer `q` and points to the white object 3.

	(!!!) Meanwhile, the gray object 2 removes a pointer `p`. As a result, the white object 3 is now effectively attached to the already scanned black object 4, as shown in figure 2.

	Then, following the normal logic of the Tricolor Marking Algorithm, all gray objects are marked as black. Consequently, objects 2 and 7 are marked as black, as shown in figure 3.

	Next, the final step of the tricolor marking process is executed, where all white objects are treated as garbage and collected, as shown in Figure 4.

	(!!!) However, the final result is that object 3, which was legitimately referenced by object 4, is mistakenly collected by the GC. This is not the desired outcome for Golang's garbage collection or for developers.
*/

/*
	2.3 NECESSARY CONDITIONS FOR TRIGGERING UNSAFE TRICOLOR MARKING
	It can be seen that there are two situations in the Tricolor Marking Algorithm that are undesirable:

		1) A white object is referrenced by a black object (the white object is attached under the black object)

		2) The reachable relationship between a gray object and a white object it references is broken (the gray object loses the reference to the white object). If both of the above conditions are met simultaneously, it results in object loss. (!!!) Furthermore, in the scenario shown above, if the white object 3 has many downstream objects, they will be cleared as well.

	To prevent this phenomenon, the simplest method is STW, which directly prohibits other user programs from interfering with object reference relationships. However, the STW process leads to significant resource waste and greatly impacts all user programs. Is it possible to maximize GC efficiency and minimize STW time while ensuring no objects are lost? The answer is yes. By using a mechanism that attempts to break the above two necessary conditions, this can be achieved.
*/

/*
	3. BARRIER MECHANISM IN Go V1.5
	The GC ensures that objects aren't lost when one of the following two conditions is met: the `strong tricolor invariant` and the `weak tricolor invariant`

	3.1 `Strong-Weak` Tricolor Invariants
		1) Strong Tricolor Invariant: No `black` object should reference a `white` object.

		The `Strong Tricolor Invariant` essentially enforces that `black` objects aren't allowed to reference `white` objects, preventing the accidental deletion of white objects.

		2) Weak Tricolor Invariant: All white objects referenced by black objects are in a `gray` protected state, as shown in figure 1.

		The `Weak Tricolor Invariant` emphasizes that while black objects can reference `white` objects, those `white` objects must have references from other `gray` objects or have a chain of `gray` objects upstream in their reachable path. This means that although a white object referenced by a black object is in a potentially dangerous state of being deleted, the references from upstream `gray` objects protect it and ensure its safety.

		To adhere to these two invariants, the GC algorithm has evolved to include two types of barriers:
			- Insertion barriers
			- Deletion barriers

	3.2 Insertion Barrier
	The specific operation of an insertion is as follows:

		When object `A` references object `B`, object `B` is marked as gray (object `B` must be marked as gray when it's attached downstream of object `A`)

	The insertion barrier effectively satisfies the `Strong Tricolor Invariant` (preventing the situation where a black object references a white object because the white object will be forcibly marked as gray)

	The pseudocode for the insertion barrier is as follows:
		addDownstreamObject(currentDownstreamSlot, newDownstreamObjectPtr) {
			// Step 1
			MarkGray(newDownstreamObjectPtr)

			// Step 2
			currentDownstreamSlot = newDownstreamObjectPtr
		}

		The pseudocode for inserting a write barrier in different scenarios is as follows:

			// A previously had no downstream, now add a downstream object B, and B is marked gray
			addDownstreamObject(nil, B)

			// A replaces the downstream object C with B, and B is marked gray
			addDownstreamObject(C, B)

		This pseudocode logic represents the write barrier. Memory slots for black objects can be located in two places:
			- Stack
			- Heap
		The stack is characterized by its small capacity but requires fast access due to frequent function calls. Therefore, the `insertion barrier` mechanism isn't used for objects in stack space but is applied in heap space.

		Check the screenshots until the 5 step, After the 5th screenshot:
			At this point, the insertion barrier will not immediately trigger GC but will instead perform an additional process.

			However, if the stack isn't scanned, there may still be white objects referenced in the stack after the tricolor marking is complete (such as object 9 in the previous figure). Therefore, the stack needs to be re-scanned with the tricolor marking to ensure no objects are lost. This re-scan requires initiating an STW (Stop-The-World) pause until the tricolor marking of the stack space is finished.

			All memory objects in the stack are re-marked as white. During this process, STW is activated to protect these white objects, preventing any read or write ops on these protected objects by intercepting and blocking them, thus avoiding external interference (e.g., new white objects being added by black objects)

			Meanwhile, other heap space objects will not trigger STW, ensuring the performance of GC collection in the heap space.

			Next, within the STW-protected are, the tricolor marking process continues until all reachable white objects are scanned and no gray nodes remain. The final state is achieved where object 1,2,3 and 9 are marked black. Since object 5 is still not scanne, it remains white.

			This time, when all memory objects are either white or black, the STW will be stopped, and the protection layer will be released.

			Finally, all remaining white nodes in the stack and heap space (objects 5 and 6) are scanned and collected. Typically, the STW duration for this stack space operation is around 10 to 100 milliseconds. After the final scan, all objects in memory will be black. This concludes the explanation of the tricolor marking and collection mechanism based on the insertion barrier.

			The purpose of the insertion barrier is to ensure that when a white object is inserted, there's a gray object to protect it, or the inserted object is turned gray. The insertion barrier essentially satisfies the strong tricolor invariant, preventing mistakenly deleted white objects.
*/
