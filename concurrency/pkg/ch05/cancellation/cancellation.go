package cancellation

/*
We must do two things:
	1) Define the period whithin wich our concurrent process is preemptable.
	2) Ensure that any functionality that takes more time than this period is itself preemptable.

An easy way to do this is to break up the pieces.

We should aim for all nonpreemptable atomic operations to complete in less time than the period we've deemed acceptable.
*/


/*
There's the problem that's lurking here. If our goroutine happens to modify shared state - e.g. a database, a file, an in-memory data structure - what happens when the goroutine is canceled?
Does your goroutine try and roll back the intermediary work it's done? How long does it have to do this work?

If possible, we should build up intermediate results in-memory and then modify state as quickly as possible.
*/

/*
When designing your concurrent processes, be sure to take into account timeouts and cancellation. Like many other topics in software engineering, neglecting timeouts and cancellation from the
beginning and then attempting to put them in later is a bit like trying to add eggs to a cake after it has been baked.
*/