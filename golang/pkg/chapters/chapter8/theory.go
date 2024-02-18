package chapter8

/*
INTRODUCTION---------------------------------------------------------------------------------------------------------------------------------------------------------------
Concurrent programming, the expression of a program as a composition of several autonomous activities, has never been more important
than it's today.

1) There's Communicating Sequential Processes (CSP) - the model of concurrency in which values are passed between independent activi
ties (goroutines) but variables are for the most part confined to a single activity.

2) There's another model is called Shared Memory Multithreading (SMM)

--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------


8.1 GOROUTINES---------------------------------------------------------------------------------------------------------------------------------------------------------------
1. Each concurrent executing activity is called a goroutine

2. When a program start, it's only goroutine is the one that calls the main function, so we call it the main goroutine.

3. New goroutines are created by the "go" statement. Syntactically, a "go" statement is an ordinary function/method call prefixed by
the keyword "go".

4. A "go" statement causes the function to be called in a newly created goroutine.

5. The "go" statement itself completes immediately.
	f() // call f(); wait for it to return
	go f() // create a new goroutine that calls f(); don't wait

6. Other than by returning from "main" or exiting the program, there's no programmatic way for one goroutine to stop another.

7. There are ways to communicate with a goroutine to request that it stop itself.

8. If the main goroutine ends with its children goroutines without waiting, the children goroutines hang after the execution of preceding
operations of the parent goroutine.
	func main() {
	fmt.Println("Hello!")
	go func() { fmt.Println(1) }()
	go func() { fmt.Println(2) }()
	go func() { fmt.Println(3) }()
	go func() { fmt.Println(4) }()
	go func() { fmt.Println(5) }()
	}

9. Gorotines are executed arbitrarily after running the parent goroutine.
	go func() { fmt.Println(1) }()
	go func() { fmt.Println(2) }()
	go func() { fmt.Println(3) }()
	go func() { fmt.Println(4) }()
	go func() { fmt.Println(5) }()

	for {
	}

	The output will be arbitrary: 31254
--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------


8.2 EXAMPLE: CONCURRENT CLOCK SERVER---------------------------------------------------------------------------------------------------------------------------------------------------------------
1. The time.Time.Format provides a way to format date and time information by example. Its argument is a template indicating how to
format a reference time, specifically Mon Jan 2 03:04:05PM 2006 UTC-0700. The reference time has 8 components.

2. The time package defines templates for many standart time formats, such as time.RFC1123. The same mechanism is used in reverse when
parsing a time using time.Parse.

3. The "killall" command is a Unix utility that kills all processes with the given name.
--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------


8.3 EXAMPLE: CONCURRENT ECHO SERVER---------------------------------------------------------------------------------------------------------------------------------------------------------------
1. Dealing with "os.Stdin" and "os.Stdout" we should be aware of the blocking operations. If we initially connect "os.Stdin" sender we'll be writing data
finitely because of the blocked consequent operations, so we should bind the "os.Stdout" recepient before binding the "os.Stdin" as a sender.
	! CHECK THE CODE !
--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------


8.4 CHANNELS---------------------------------------------------------------------------------------------------------------------------------------------------------------
1. Channels are the connections between goroutines. Channels let one goroutine send values to another goroutine.

2. Channel has its own "element type". The example of channel with type int: "var chanName chan int"

3. To create a channel we use the built-in "make" function. As with maps a channel is a reference to the data structure created by make, so when we copy a
channel or pass one as an argument to a fuction, we are copying its reference, so the caller and callee refer to the same data structure.
	ch := make(chan int)

4. The zero value of a channel is nil. Perhaps surprisingly, nil channels are sometimes useful. Because send and receive operations on a nil
channel block forever, a case in a "select" statement whose channel is nil IS NEVER SELECTED. This lets us use nil to enable or disable
cases that correspond to features like handling timeouts or cancellation responding to other input events, or emitting output.

5. Two channels of the same type may be compared using ==. The comparison is true if both are references to the same channel data structure. A channel may be
also be compared to nil.

6. A channel has two principal operations "send" and "receive", collectively known as communications.
	1) A "send" statement transmits a value from a goroutine, through the channel, to another goroutine executing a corresponding receive expression.

Both operations are written using "<-" operator.
	1) In a send statement, the "<-" separates the channel and value operands: "CHANNEL <- value" thus: "value is passed through the channel"

	2) In a receive expression, "<-" precedes the channel operand: "x = <- CHANNEL" thus: "the value goes from the channel"
	2!) A receive expression whose result is not used is a valid expression as well: "<- channel" is a receive statement and its result is discarded.

	3) Channels support a third operation "close", which sets a flag indicating that no more values will ever be sent on this channel.
		1) Consequent attempts to send will panic.

		2) Receive operations on a closed channel yield the values that have been sent until no more values are left. Any receive operations thereafter complete
		immediately and yield the zero value of the channel's element type.

	To close a channel, we call the built-in close function: close(channelName)
		1) There's no way to test directly whether a channel has been closed, but there's a variant of the receive operation that produces
		two results: channel element & boolean value, conventionally called "ok", which is true for a successful receive and false for a receive on a closed
		and drained channel.


7. There are the options while creating a channel with "make" operator:
	1) Omitting the channel's capacity we create "unbuffered" channel that's also called "synchronous" channel.

	2) Using the second non-zero argument we create a "buffered" channel with the capacity.
		ch = make(chan int)    // unbuffered channel
		ch = make(chan int, 0) // unbuffered channel
		ch = make(chan int, 3) // buffered channel with capacity 3
--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------


8.4.1 UNBUFFERED CHANNELS---------------------------------------------------------------------------------------------------------------------------------------------------------------
1.SEND.  A send operation on an "unbuffered" blocks the sending goroutine until another goroutine executes a corresponding receive op on the same channel, at
which point the value is transmitted and both goroutines may continue.

2. RECEIVE. Conversely, if the receive operation was attempted first, the receiving goroutine is blocked until another goroutine performs a send op on the same
channel.

3. Communication over an unbuffered channel causes the sending and receiving goroutines to synchronize. Because of this, unbuffered channels are sometimes called
"synchronous" channels. Sending and receiving the value on an unbuffered channel are synchronized and act on the principle "one-on-one".

4. In discussion of concurrency, when we say "x" happens before "y", we don't mean merely x occurs earlier in time than y. We mean that it's guaranteed to do so
and that all its prior effects, such as updates to variables, are complete and that you may rely on them.

5. When "x" neither happens before "y" nor after "y", we say that "x" is concurrent with "y". This doesn't mean that "x" and "y" necessarily simultaneous,
merely that we can not assume anything about their ordering.

6. Messages sent over the channels have two important aspects. Each message has a value, but sometimes the fact of communication
and the moment at which it occurs are just as important. We call messages event when we wish to stress this aspect. When the event
carries no additional information, that is, its sole purpose is synchronization, we'll emphasize this using a channel whose element
type is struct{}. Though it's common to use a channel of "bool" or "int" for the same purpose since done "<- 1" is shorter than done "<-
struct{}{}".

--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------


8.4.2 PIPELINES---------------------------------------------------------------------------------------------------------------------------------------------------------------
1. Channels can be used to connect goroutines together so that the output of one is the input of another. This is called a pipeline.

2. If we want say to the "hanging" goroutines that they needn't wait for others message, we close the goroutine using the built-in "close()" function.
	1) After a channel has been closed, any further send operations on it will panic.

	2) After the closed channel has been drained, that is, after the last sent element has been receive, all consequent receive operations
	will proceed without blocking but will yield a zero value of the channel's type.
We can check whether the channel is drained and closed in the following way:

	if value, ok := <- channelName; ok {
		...
	}

3. We can use a range loop to iterate over channels too. This is a more convenient syntax for receiving all the values sent on a channel and
terminating the loop after the last one.

	for value := range channelName {

	}

	This method is the similar to the following one:
	for {

	if value, ok := <- channelName; !ok {
		break
	}

	}

It's only necessary to close a channel when it's important to tell the receiving goroutines that all data have been sent.
	! ATTEMPTING TO CLOSE AN ALREADY-CLOSED CHANNEL CAUSES A PANIC. THE SAME LOGIC WILL BE EXECUTED AFTER AN ATTEMPT TO CLOSE A NIL CHANNEL !

4. A channel that the garbage collector determines to be unreachable will have its resources reclaimed whether or not it's closed.

--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------


8.4.3 UNIDIRECTIONAL CHANNEL TYPES---------------------------------------------------------------------------------------------------------------------------------------------------------------
1. Go type system provides unidirectional channel types that expose only one or the other of the send and receive operations.

2. SENDER. The type "channel <- typeName" is the send only channel of "typeName" type, allows sends but not receives.
	1) Since the close operation asserts that no more sends will occur on a channel, only the sending goroutine is in a position to call it,
	and for this reason it's a compile time error to attempt to close a receive-only channel.

3. RECEIVER. The type "<- chan typeName" , a receive-only channel of "typeName" type, allows receives but not sends.

4. After passing the channel the compiler implicitly converts the channel passed.

5. Convertions from bidirectional to unidirectional channel types are permitted in any assignment. There's no going back, however: once
you have a value of unidirectional type such as "chan <- int", there's no way to obtain from it a value of "chan int" that refers to the same
channel data structure.

--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------


8.4.4 BUFFERED CHANNELS---------------------------------------------------------------------------------------------------------------------------------------------------------------
1. A buffered channel has a queue of elements. The queue's maximum size is determined when it's created, by the capacity argument.
	channel := make(chan string, 10) creates a channel that's possible to contain 10 strings.

2. We can find out the channel's capacity using the function "cap(channelName)"

3. The built-in "len(channelName)" function returns the number of elements currently buffered. Since in a concurrent program this information is likely
to be state as soon as it's retrieved, its value is limited, but it could conceivably be useful during fault diagnosis or performance optimization.


2. Operations.
	1) SEND. A send operation on a buffered channel inserts an element at the back of the queue.
		1) If the channel is full, the send operation blocks its goroutine until space is made available by another goroutine's receive.
	2) RECEIVE. A receive operation removes an element from the front.
		1) If the channel is empty, a receive operation blocks its goroutine until a value is sent by another goroutine.

3. The channel's buffer decouples the sending and receiving goroutines.

4. We cannot use a buffered channel as a queue because channels deeply connected to goroutine scheduling, and without another goroutine receiving from the
channel, a sender - and perhaps the whole program - risks becoming blocked forever.

5. CHECK THE CODE
Had we used an unbuffered channel, the slower goroutines would have got stuck trying to send their responses on a channel from which no goroutine will ever
receive. This situation called a "goroutine leak", would be a bug. Unlike garbage variables, leaked goroutines are not automatically collected, so it's important
to make sure that goroutines terminate themselves when no longer needed.

6. Failure to allocate sufficient buffer capacity would cause the program to deadlock.

CHANNEL BUFFERING CONS EXPLANATION

UNBUFFERED processing:
Imagine 3 cooks in a cake shop, one baking, one icing and one inscribing each cake before passing it on to the next cook in the assebly line. In a kitchen with
little space, each cook that has finished a cake must wait for the next cook to become readt to accept it, This rendezvous is analogous to communication over an
unbuffered channel.

BUFFERED processing:
If there's a space for one cake between each cook, a cook may place a finished cake there and immediately start work on the next. This is analogous to a buffered
channel with capacity 1. So long as the cooks work at about the same rate on average, most of these handovers proceed quickly, smoothing out transient differencies
in their respective rates. More space between cooks - larger buffers - can smooth out bigger transient variations in their rates without stalling the assembly
line, such as happens when one cook takes a short break, then rushes to catch up.

THE DIFFERENCIES BETWEEN WORKERS:
1) If an earlier stage of the assembly line is consistently faster than the following stage, the buffer between them will spend most of its time full.
2) Conversely, if the later stage is faster the buffer will usually be empty.
A BUFFER PROVIDES NO BENEFIT IN THIS CASE.

GOROUTINE HELPER
If the second stage is more elaborate, a single cook may not be able to keep up with the supply from the first cook or meet the demand from the third. To solve
the problem, we could hire another cook to help the second, performing the same task, bit working independently.


--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------



8.5 LOOPING IN PARALLEL---------------------------------------------------------------------------------------------------------------------------------------------------------------
! CHECK THE CODE !

1.
	// Wait for all the goroutines
	func makeThumbnail(filenames []string) {
		done := make(chan bool)
		for _, f := range filenames {
			go func(f string) { // Defines the anonymous function that takes only one argument.
				ImageFile(f)
				done <- true
			}(f) 				// Catch the current value, so that each goroutine works with each own value.
		}

		for range filenames {
			<-done
		}
	}

2. If we don't know the count of the data entries we can use "sync.WaitGroup" to know when the last goroutine has finished. It blocks the goroutine
it were called in increments the number of working goroutines when a new goroutine has launched and decrement the counter when its work done.
	1) wg.Add(1) must be called before starting a new goroutine (worker). Otherwise we wouldn't be sure that the wg.Add() happens before
	the "closer" goroutine calls wait.

	2) wg.Done() must be called when the goroutine finished its work. We can use defer wg.Done(). It's implicit argument to be incremented
	with is -1.

	3) We should declare the goroutine where the "wg.Wait()" and "close(channelName)" will be placed.

		func MakeThumbnails6(filenames <-chan string) (size int64) {
			sizes := make(chan int64)
			var wg sync.WaitGroup 	 // The number of working goroutines at a moment

			for f := range filenames {
				wg.Add(1) 			 // Increment the current amount of workers
				// worker
				go func(f string) {
					defer wg.Done()  // Decrement the current amount of workers after all goroutine's activities done
					...
				}(f)

			// closer
			go func() {
				wg.Wait() 			 // Wait for the moment when all the goroutines finished its work
				close(sizes) 		 // close the channel
			}()

			...
			return
		}

--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------


8.6 CONCURRENT WEB CRAWLER---------------------------------------------------------------------------------------------------------------------------------------------------------------
! CHECK THE CODE !

1. The crawl goroutine takes "link" as an explicit parameter to avoid the problem of loop variable capture.

2. Initial send of the command-line arguments to the work list must run in its own goroutine to avoid "deadlock". An alternative
solution would be to use a buffered channel.

3. Unbounded parallelism is rarely a good idea since there's always a limiting factor in the system, such as:
	1) The number of CPU cores for computer-bound workloads
	2) The number of spindles and heads for local disk I/O operations
	3) The bandwidth of the network for streaming downloads
	4) The serving capacity of a web service.

The solution is to limit the number of parallel uses of the resource to match the level of parallelism that is available. A simple
way to do that in our example is to ensure that no more than N calls to Extract() are active at once, where N is comfortably less
than the file descriptor limit - 20, for example. This is analogous to the way a doorman at a crowded nightclub admits a guest only when
some other guest leaves.

We can limit parallelism using a buffered channel of capacity N to model a concurrency primitive called "counting semaphore". Conceptually,
each of the N vacant slots in the channel buffer represents a token entiting the holder to proceed. Sending a value into the channel acquires
a token, and receiving a value from the channel releases a token, creating a new vacant slot. This ensures that at most N sends can occur
without an intervening receive (Although it might be more intuitive to treat filled slots in the channel buffer as tokens, using vacant slots
avoids the need to fill the channel buffer after creating it). Since the channel element type isn't important, we'll use struct{}, which has
size zero.

4. The second problem is that the program never terminates, even when it has discovered all the links reachable from the initial URLs. We
can set up the counter of goroutines have been started by the main cycle and each time we're left the current iteration we decrement this
variable by one.
--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------


8.7 MULTIPLEXING WITH SELECT---------------------------------------------------------------------------------------------------------------------------------------------------------------
1. So that we can multiplex channels operations we need "select" statement:
	select {
	case <- ch1: 		// Discard the value of the channel "ch1"
		// ...
	case x := <-ch2: 	// Use the value of the channel "ch2"
		// ... use x ...
	case ch3 <- y: 		// Send the value "y" on the "ch3" channel
		// ...
	default:
		// ...
	}

2. A "select" waits until a communication is ready to proceed. It then performs that communication and executes the case's associated
statements. The other communications don't happen.
	1) If multiple cases are ready, select picks one at random, which ensures that every channel has an equal chance to be selected.

3. A select without any cases "select{}" waits forever.


4. Each case specifies communication (a send or receive operation on some channel) and an associated block of statements.
	RECEIVE:
	1) A receive expression can appear on its own
	2) A receive expression can apper within a short variable declaration. In this case we can refer to the received value.

5. Sometimes we wanna try to send/receive on a channel but avoid blocking if the channel is not ready. We need the "non-blocking communication".
A select statement can do that too using "default" statement inside "select" body, which specifies what to do when none of the other communicatons
can proceed immediately. The statements that are executed when there's no ability to send/receive are called "polling" because we wait a channel
to get us an ability to make the operation needed.

6. The zero value of a channel is nil. Perhaps surprisingly, nil channels are sometimes useful. Because send and receive operations on a nil
channel block forever, a case in a "select" statement whose channel is nil IS NEVER SELECTED. This lets us use nil to enable or disable
cases that correspond to features like handling timeouts or cancellation responding to other input events, or emitting output.

	select {
	case <- abort:
		fmt.Printf("Launch aborted!\n")
	default:
		time.Sleep(500*time.Millisecond)
	}

7. Deterministic output.

	func ForSelect() {
	ch := make(chan int, 1)

	for i := 0; i < 10; i++ {
		select {
		case value := <-ch:
			fmt.Println(value)
		case ch <- i:
		}
	}
	}

	EXPLANATION:
	1st iteration:
		1) Try to read the value from the channel. There's nothing on a channel.
		2) We put the value "0" on a channel

	2nd iteration:
		Read the "0" and print it.

	3rd iteration:
		1) Try to read the value from the channel. There's nothing on a channel.
		2) We put the value "2" on a channel

	AND SO ON

8. Nondeterministic output

	func NondeterministicForSelect() {
	ch := make(chan int, 2)

	for i := 0; i < 10; i++ {
		// Each case after the first iteration have the same chances to be executed
		select {
		case value := <-ch:
			fmt.Println(value)
		case ch <- i:
		}
	}
	}

	EXPLANATION:
	1st iteration:
		1) Try to read the value from the channel. There's nothing on a channel.
		2) We put the value "0" on a channel

	2nd iteration RANDOM OPERATION
		1) The operation will be chosen at random, so either reading or writing will be taken.

	3rd iteration
		1) The channel is full. The writing on a channel is impossible, so we print the element at the head of the queue.

	AND SO ON

9. The time.Tick() function behaves as if it creates a goroutine that calls time.Sleep() in a loop, sending an event each time it wakes up.
	1) When we've finished all the work with "ticker := time.NewTicker(1 * time.Second)" we should close it with "ticker.Stop()" operation.

*/
