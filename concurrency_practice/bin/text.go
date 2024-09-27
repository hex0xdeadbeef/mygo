package main

// Практика Go — Concurrency. - Хабр: https://habr.com/ru/articles/759584/
/*
CLOSED CHANNEL ---------------------------------------------------------------------------------------------------------------------------------
1. If a channel closed we cannot send data on it.
2. If a channel closed we cannot the elements sent before closing.


NIL CHANNEL ---------------------------------------------------------------------------------------------------------------------------------
1. nil channel always blocks
2. "select" operator never chooses a block with nil channel

CHANNEL AXIOMS ---------------------------------------------------------------------------------------------------------------------------------
1. Sending on a nil channel always blocks
2. Getting from a nil channel always blocks

3. Sending on a closed channel always panics
4. Getting info from a closed channel is instant. We get zero values and ok val is got the "false" vals.

5. The default value for a channel is "nil"
*/

// Многопоточность и параллелизм в Go: Goroutines и каналы. Хабр - https://habr.com/ru/companies/mvideo/articles/778248/
/*
GOROUTINE INTERNALS ---------------------------------------------------------------------------------------------------------------------------------
1. Goroutine consists of 70 lines of code.

type g struct {
	stack 		stack
	stackguard0 uintptr
	stackguard1 uintptr

	_panic		*_panic
	_defer		*_defer

	m			*m

	sched		gobuf

	syscallsp	uintptr
	syscallpc	uintptr
	...
	param		unsafe.Pointer

	atomicstatus	uint32

	stackLock		uint32

	goid			int64
	...
	preempt			bool
	preemptStop		bool
	preemptShrink	bool
	...
	waiting			*sudog
	...
	gcAssistBytes	int64
}

1. stack - the stucture describing a stack of that goroutine. It includes the pointers to low/top boundaries of stack
2. stackguard0 and stackguard1 - are used to check whether the stack is exceeded

3. _panic - the pointer to a structure of panic. If this goroutine is in panic.
4. _defer - the pointer to a structure of deferred call used to realize "defer" mechanism

5. m - the pointer to a machine structure (M) this goroutine bound to
	1) In Go the machine is OS thread

6. sched (gobuf) - the structure of scheduler's context including the pointers to a stack and the registers.

7. syscallsp and syscallpc (uintptr) - are used while performing some system calls.

8. param (unsafe.Pointer) - an arbitrary parameter used to transmit data from one goroutine to another one

9. atomicstatus (uint32) - a status of goroutine used to manage its state.

10. stackLock (uint32) - a lock used to manage an access to goroutine's stack

11. goid (int64) - a global unique goroutine's identifier

12. preempt, preemptStop, preemptShrink (bool) - are flags used to manage interruptions of goroutines process.

13. gcAssistBytes (int64) - an amount of bytes, that goroutine must help to arrange to garbage collector (GC).


GO SCHEDULER ---------------------------------------------------------------------------------------------------------------------------------
1. Go scheduler multiplexes all the active goroutines on the limited amount of OS threads (M). It allows to use multiprocessing efficiently without creation an excessive amount of OS
threads (M).

2. Go uses cooperative locks. It means that context switches happens in the specific points of a program process, for example: while I/O operations, explicit calls of functions of
locking / unlocking.

3. Go scheduler uses the algorith of "work stealing" to balance the workload between the OS threads (M). If an OS thread finishes performing of its goroutines, it migh steal the
goroutines from the queue of another OS thread to supply the uniform work distribution.


GO STACK ---------------------------------------------------------------------------------------------------------------------------------
1. The variables that must be saved after goroutine finishes are tranmitted to the heap automatically. The compiler does this work.

2. The goroutine blocking happens while performing is stopped until the specific condition or action is performed.

3. Goroutine can be blocked trying to catch the mutex already gotten by another goroutine. Similarly a goroutine can be blocked calling Wait() on sycn.WaitGroup instance until others
ones call Done().

4. Network queries or I/O operations on files can also cause the lock of goroutine until the actione is finally finished.

5. "Nonlock" means that goroutine's execution continues without a stop, even other operations are finished. Incorrect managing of locks can result in performance descent and even
deadlocks (some goroutines lock each other).


CHANNELS USING PATTERNS ---------------------------------------------------------------------------------------------------------------------------------
1. Timeout Control

2. Channel Orchestration

3. Work cancellation
*/
