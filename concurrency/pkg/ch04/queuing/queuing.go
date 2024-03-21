package queuing

/*
Sometimes it's useful to begin accepting work for your pipeline even though the pipeline isn't yet ready for more. This process is called "queuing".

All this means is that once our stage has completed some work, it stores the result in a temporary location in memory so that other stages can retrieve it later, and this stage
doesn't need to hold a reference to it.
*/

/*
Adding queuing prematurely can hide synchronization issues such as deadlocks and livelocks, and further, as your program converges toward correctness, you may find that you need more
or less queuing.
*/

/*
Queuing will almost never speed up the total runtime of your program. It will only allow the program to behave differently.
*/

/*
The situations in which queuing can increase the overall performance on the nature of your pipeline:
	1) If batching requests in a stage that saves time. For example: a stage that buffers input in something faster (e.g. memory) than it's designed to send to (e.g. disk). This of
	course the entire purpose of Go's "bufio" package.

	We can encounter chunking anytime performing an operation that requires an overhead, chunking may increase system performance. Some examples of this are: opening database transactions calculating message checksums, and allocating contiguous space.

	Aside from chunking, queuing can also help if your algorithm can be optimized by supporting lookbehinds or ordering.

	2) If delays in a stage produce a feedback loop into the system (For example: if one of the stages of data processing in the system process data slowly, this results in others
	stages' waiting on the ending of the first stage. As a result the overall time will be increased).

	If the efficiency of the pipeline drops below a certain critical threshold, the systems upstream from the pipeline begin increasing their inputs into the pipeline, which causes
	the pipeline to lose more efficiency, and the death-spiral begins. Without some sort of fail-safe, the system utilizing the pipeline will never recover.

	By introducing a queue at the entrance to the pipeline, you can break the feedback loop at the cost of creating lag for requests. From the perspective of the caller into the
	pipeline, the request appears to be processing, but taking a very long time. As long as the caller doesn't time out, your pipeline will remain stable. If the caller does time out,
	we need to be sure that we support some kind of check for readiness when dequeuing. If we don't, we can inadvertently create a feedback loop by processing dead requests thereby
	decreasing the efficiency of out pipepine.
*/

/*
Queuing should be implemented either:
	1) At the entrance to our pipeline.
	2) In stages where batching will lead to higher efficiency.
*/

/*
How large our queues should be? We will follow Little's Law.
L = λW, where:

- L is the average number of units in the system
- λ is the average arrival rate of units
- W is the average time a unit spends in the system

1) Decreasing W by a factor of N (W/W) -> We have 2 options:
	1) L*N
	2) λ*N
2) Increasing L (Adding queues to our stages) (L*N):
	1) L*N
	2) λ*N
So, using Little's law we've proven that queuing won't help decrease the amount of time spent in a system.

This equation only applies to "stable" systems. In a pipeline a stable system is one in which the rate that work enters the pipelinee, or ingress, is equal to the rate in which it
exits the system, or egress.
	1) If the rate of ingress exceeds the rate of egress, your system is unstable and has entered a death spiral.
	2) If the rate of ingress is less than the rate of egress, we still have an unstable system, but all that's happening is that your resources aren't being utilized completely.
*/
