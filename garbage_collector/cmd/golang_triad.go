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
		2) The reachable relationship between a gray object and a white object it references is broken (the gray object loses the reference to the white object).

		If both of the above conditions are met simultaneously, it results in object loss. (!!!) Furthermore, in the scenario shown above, if the white object 3 has many downstream objects, they will be cleared as well.

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

			Next, within the STW-protected are, the tricolor marking process continues until all reachable white objects are scanned and no gray nodes remain. The final state is achieved where object 1,2,3 and 9 are marked black. Since object 5 is still not scanned, it remains white.

			This time, when all memory objects are either white or black, the STW will be stopped, and the protection layer will be released.

			Finally, all remaining white nodes in the stack and heap space (objects 5 and 6) are scanned and collected. Typically, the STW duration for this stack space operation is around 10 to 100 milliseconds. After the final scan, all objects in memory will be black. This concludes the explanation of the tricolor marking and collection mechanism based on the insertion barrier.

			The purpose of the insertion barrier is to ensure that when a white object is inserted, there's a gray object to protect it, or the inserted object is turned gray. The insertion barrier essentially satisfies the strong tricolor invariant, preventing mistakenly deleted white objects.

	3.3 Deletion Barrier
	The specific operation of a `deletion barrier` is as follows: if a deleted object is `gray` or `white`, it's marked as gray

	The `deletion barrier` essentially satisfies the weak tricolor, aiming to ensure that the path from gray objects to white objects remains intact. The pseudocode for the deletion barrier is as follows:

		addDownstreamObject(currentDownstreamSlot, newDownstreamObjectPtr) {
			// Step 1
			if (currentDownstreamSlotIsGray || currentDownstreamSlotIsWhite) {
				// The slot being removed, mark it as gray
				markGray(currentDownstreamSlot)
			}

			// Step 2
			currentDownstreamSlot = newDownstreamObjectPtr
		}

	The pseudocode for the deletion barrier is different scenarios is as follows:

		// Object A removes the reference to object B. B is removed by A an marked gray (if B was previously white)
		addDownstreamObject(B, nil)

		// Object A replaces downstream B with C. B is removed by A and marked gray (if B was previously white)

	The deletion barrier ensures that when an object's reference is removed or replaced by an upstream object, the removed object is marked as gray. The purpose of marking it as gray is to prevent a white object from being collected if it's later referenced by another black object. This mechanism ensures that the deleted object, which was not intended to be collected, is properly protected. The following diagrams will simulate a detailed process to provide a clearer view of the overall workflow.

	If gray object 1 deletes the white object 5 without triggering the deletion barrier mechanism, white object 5, along with its downstream objects 2 and 3, will be disconnected from the main chain and eventually be collected.

	However, with the current tricolor marking algorithm using the deletion barrier mechanism, deleted objects are marked as gray. This is to protect object 5 and its downstream objects (consider why protection is necessary - what unexpected issues could arise if object 5 were not marked as gray? For example, if object 1 deletes object 5 but object 5 remains white, and if the process doesn't include STW protection, there's a high chance that, during the deletion process, an already black-marked object might reference object 5. Object 5 is still a valid memory object needed by the program, but it could be mistakenly collected by the GC as a white object. This happens because the downstream objects of black objects are not protected). Therefore, object 5 is marked as gray.

	Following the tricolor marking process, the next step is to traverse the gray marking table, which includes objects 1,4, and 5. During this traversal, their reachable objects are marked gray from white. Additionally, the gray objects being traversed are marked as black.

	This process of remarkering continues until there are no gray objects (two groups of black and white objects are filled).

	Finally, execute the collection and cleanup process, removing all white objects through GC.

	The above describes the tricolor marking process using the deletion barrier. The deletion barrier still allows for garbage collection in a parallel state, but this method has lower precision. Even if an object is deleted, the last pointer pointing to it can still `survive` this round of garbage collection and will only be cleared in the next GC cycle.

*/

/*
	4.HYBRID WRITE BARRIER IN GO V1.8

	While both `insertion` and `deletion` write barriers can address the issue of non-parallel processing caused by STW to some extent, each has its own limitations:
		- Insertion Write Barrier: (!!!) Requires STW at the end to re-scan the stack and mark the live white objects referrenced on the stack.
		- Deletion Write Barrier: Has lower accuracy in garbage collection. At the start of GC, STW scans the heap to record an initial snapshot, which protects all live objects at that moment.

	Go V1.8 introduces the Hybrid Write Barrier mechanism, which avoids the need for stack re-scanning, significantly reducing STW time. This mechanism combines the advantages of both insertion and deletion barriers.


	4.1 Hybrid Write Barrier Rules
	The specific ops of the hybrid write barrier generally need to follow the following conditions:

		1) At the start of GC, all objects on the stack are scanned and marked as black (there's no need for a second scan or STW)
		2) Any new objects created on the stack during GC are marked as black

		3) Deleted objects are marked as gray
		4) Added objects are marked as gray

	The hybrid write barrier effectively satisfies a variant of the weak tricolor invariant. The pseudocode is as follows:

		addDownstreamObject(currentDownstreamSlot, newDownstreamPtr) {
			// 1
			markGray(currentDownstreamSlot) // mark the current downstream

			// 2
			markGray(newDownstreamPtr)

			// 3
			currentDownstreamSlot = newDownstreamPtr
		}

	Note: (!!!) the barrier technique is not applied to stack objects to ensure stack ops efficiency. (!!!) The hybrid write barrier is a GC barrier mechanism and is triggered only during GC execution.

	The flow of Execution:

	Put all objects to the `white` group.

	As GC begins, following the steps of the hybrid write barrier, the first step is to scan the stack area and mark all reachable objects as black. Therefore, when the stack scan is complete, reachable objects of the stack are marked as black and added to the black marking table.

	Next, we'll analyze several scenarios involving the hybrid write barrier. This section will outline four scenarios, all starting from the point shown in figure `1_Configure_Black_Group_From_Stack_Objects`, where the stack space has been scanned and reachable objects are marked as black. The four scenarios are:
		- Heap Deletion of Reference, becoming a Stack Downstream
		- Stack Deletion of Reference, becoming a Stack Downstream
		- Heap Deletion of Reference, becoming a Heap Downstream
		- Stack Deletion of Reference, becoming a Heap Downstream

	4.2 Scenario 1: Heap Deletion of Reference, Becoming a Stack Downsteam
	Scenario 1 mainly describes the case where an object is deleted from a heap object's references and becomes a downstream object of a stack object. The pseudocode is as follows:

		StackObj1 -> Obj7 = HeapObj7 // Attach heap object 7 to stack object 1 as a downstream
		HeapObj4 -> Obj7 = null; // HeapObj4 deletes reference to Obj7

	Now, execute the first line of code, from the above scenario: StackObj1 -> Obj7 = HeapObj7, which adds the white object 7 to the downstream of the black stack object 1. It's important to note that because the stack doesn't trigger the write barrier mechanism, the white object 7 is directly attached under the black object 1, and object 7 remains white. Now, scan object 4, marking object 4 as gray.

	Then, execute the second line of code from the avove scenario: HeapObj4 -> Object7 = null. The gray object 4 deletes the white object 7 (deletion means reassigning to null). Since object 4 is in the heap space, the write barrier is triggerred, and the deleted object 7 is marked as gray.

	Thus, from the Scene One, object 7 end up being attached to the downstream of object 1. Since object 7 is gray, it won't be collected as garbage, thus it's protected. In the hybrid write barrier of Scene One, STW will not be triggered again to re-scan the stack space objects. The subsequent process continues to follow the hybrid write barrier's tricolor marking logic.

	4.3 Scnenario Two: Stack Deletes Reference, Becomes Stack Downstream

	Scenario Two primarily describes a situation where an object is deleted by one stack object and becomes the downstream of another stack object. The pseudocode is as follows:

		new stack object9
		object9 -> object3 = object3
		object2 -> object3 = null

	Now, execute the first line of code in the above scenario: `new stack object9`. According to the hybrid write's barrier's rule, any new memory object created in the stack space will be marked as black. Here, we inherit the memory layot scenario describet in the prev. figure, so object 9 is currently a black object referenced by the program's root set.

	Next, execute the second line of code in the above scenario: `object9 -> object3 = object3`. Object 9 adds downstream object3. Since object9 is in the stack space, this addition doesn't trigger the write barrier. Object 3 is directly attached to the downstream of object 9.

	Finally, execute the third line of code in the above scenario: `object2 -> object3 = null`. Object 2 will remove downstream object 3. Since object 2 is within the stack space, the write barrier mechanism isn't triggered. Object 2 directly detaches object 3 from its downstream.

	From the above process, we can see hat in the hybrid write barrier mechanism, an object transferring from the downstream of one stack object to another stack object's dowmnstream ensures the object's safety without the need to trigger the write barrier or the STW mechanism, since stack objects are all black. This highlights the ingenious design of the hybrid write barrier.

	4.3 Scenario Three: Heap Reference Deletion, Becoming a Heap Downstream

	Scenario three mainly describes a situation where an object is removed from the reference of one heap object and becomes the downstream of another heap object. The pseudocode is as follows:
		heapObject10 -> object7 = heapObject7
		heapObject4 -> object7 = null

	Now let's restore the memory layout for scenario three. In the heap space, there's a black object10 (whether it's black or not doesn't matter, because object10 is a reachable heap object and it'll eventually be marked as black). Here, we don't consider the possibility of object10 being any other color. Since black is a special case, if object10 is white, its downstream will eventually be scanned and thus safe. Similarly, if object10 is gray, its downstream is also safe. Only when object 10 is black, there's a risk of its downstream memory being unsafe. So, the current colour of the object10 is black.

	Executing the 1st line of code in the above scenario: `heapObject10 -> object7 = heapObject7`, heapObject10 adds a dwonstream reference to the white heapObject7. Since object10 is within the heap space, this write op will trigger the barrier mechanism. According to the conditions of the hybrid write barrier, the added object will be marked as gray. Thus, the white object7 will be marked as gray, indirectly protecting the write object6.

	Next, executing the second line of code in the above scenario: `heapObject4 -> object7 = null`, the gray heapObject4 is within the heap space, this op triggers the barrier mechanism. According to the conditions of the hybrid write barrier, the removed object will be marked as gray. Therefore, object7 will be marked as gray (even though it's already gray).

	Through the processes described, the originally white object7 has successfully been transferred from beneath heapObject4 to beneath black heapObject10. Both object 7 and its downstream object 6 are protected, and throughout this process, no STW was used, avoiding any delay in the program's performance.

	4.4 Scenario Four: Stack Deletes Reference, Becomes Heap Downstream
	Scenario Four primarily describes the case where an object is deleted from a stack and becomes a downstream object of another heap object. The pseudocode is as follows:

		stackObject1 -> Object2 = null
		heapObject4 -> object7 = StackObject2

	Continuing from the starting memory layout, executing the first line of code, `stackObject1 -> Object2 = null`, removes stackObject1's reference to stackObject2. Since stackObject1 is in the stack space, this op doesn't trigger the write barrier mechanism. Consequently, stackObject1 directly deletes object2 and all downstream objects associated with object2.

	Next, executing the second line of code in the scenario: `heapObject4 -> object7 = StackObject2`, actually performs two actions:
		- It removes the prev. downstream white object, object7, from heapObject4
		- It adds stackObject2 as the new downstream object for heapObject4
	When heapObject4 performs the two actions described, it triggers the write barrier mechanism since heapObject4 is in the heap space. According to the rules of the hybrid write barrier, deleted objects are marked as gray, and newly added objects are also marked as gray. Consequently, object7 is marked as gray, thereby protecting its downstream object6 (object6). Object2, being newly added, will also go through the gray marking process. Since object2 is already black (and thus safe), it will remain black.

	Ultimately, a question arises: Why were object7 and object6, which are already disconnected from the program's Root Set, not reclaimed? This is due to the delay issue with the hybrid write barrier. To avoid triggering STW, some memory may be delayed in being reclaimed for one cycle. In the next round of GC, if object7 and object6 are not referrenced by external objects, they will eventually be reclaimed as white garbage memory.
*/

/*
	5. SUMMARY
	STW is necessary for memory safety, but the hybrid write barrier mechanism significantly reduces the need for STW, allowing for concurrent garbage collection without pausing the program.

	We can note the key milestones in the evolution of Golang's GC:
	1) Go 1.3: The traditional mark-and-sweep approach required STW, resulting in low efficiency
	2) Go 1.5: Introduced tri-color marking combined with insertion and deletion write barriers, where write barriers were applied to heap space but not stack space. This method required an additional stack scan after initial scan, leading to moderate efficiency.
	3) Go 1.8: Implemented tri-color marking with a hybrid write barrier mechanism, where the stack space didn't require barriers while the heap space did. This approach minimized the need for STW and achieved higher efficiency.
*/
