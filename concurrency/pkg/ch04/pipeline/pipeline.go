package pipeline

import (
	"fmt"
	"time"
)

/*
A pipeline is just another tool we can use to form an abstraction in our system. In particular, it's a very powerful tool to use when our program need to process streams
or batches of data. A pipeling is nothing more that a series of things that take data in, perform an operation on it, and pass the data back out. We call each of these ops
a stage of the pipeline.
*/

/*
By using a pipeline we separate the concerns of each stage, which provides numerous of benefits. We can modify stages independent of one another, we can mix and match how
stages are combined independent of modifying the stages, we can process each stage concurrent or downstream stages, and we can fan-out, or rate-limit portions of our
pipeline.
*/

/*
Pipeline properties:
	1) A stage consumes and returns the same type
	2) A stage must be reified by the language so that it may be passed around. Functions in go are reified and fit this purpose nicely. "reified" means that we can pass functions
	around our program.
*/

/*
1. Butch processing means that stages operate on chunks of data all at once instead of one discrete value at a time.
2. Stream processing means that the stage receives and emits one element at a time.
*/

/*
Original data remains unaltered
*/

/*
Closing the done channel cascades through the pipiline. This is made possible by two things in each stage of the pipeline:
	1) Ranging over the incoming channel. When the incoming channel is closed, the range will exit.
	2) The send sharing a select statement with the done channel.

Regardless of what state the pipeline stage is in - waiting on the incoming channel, or waiting on the send closing the done channel will force the pipeline stage to terminate.
*/

/*
There's a recurrence realtion in play here. At the beginning of the pipeline, we've established that we must convert discrete values into a channel. There are two points in this process
that must be preemptable:
	1) Creation of the discrete value that is not nearly instantaneous.
	2) Sending of the discrete value on its channel.
*/

/*
If one of our stages is computationally expensive, this will certainly eclipse the performance overhead.
*/

func MultiplyBatchStage(values []int, multiplier int) []int {
	multipliedValues := make([]int, len(values))
	for i, v := range values {
		multipliedValues[i] = v * multiplier
	}

	return multipliedValues
}

func AddBatchStage(values []int, additive int) []int {
	addedValues := make([]int, len(values))
	for i, v := range values {
		addedValues[i] = v + additive
	}

	return addedValues
}

func AddMultiplyBatchCombined() {
	vals := []int{1, 2, 3, 4}
	result := AddBatchStage(MultiplyBatchStage(vals, 2), 10)
	for _, val := range result {
		fmt.Printf("%v ", val)
	}

	fmt.Println()
}

func MultiplyStreamStage(val int, multiplier int) int {
	return val * multiplier
}

func AddStreamStage(val int, additive int) int {
	return val + additive
}

func AddMultiplyStreamCombined() {
	vals := []int{1, 2, 3, 4}
	for _, v := range vals {
		fmt.Printf("%v ", AddStreamStage(MultiplyStreamStage(v, 2), 10))
	}
}

func ConcurrentPipeline() {
	generator := func(done <-chan struct{}, integers ...int) <-chan int {
		intStream := make(chan int)
		go func() {
			defer close(intStream)

			for _, val := range integers {
				select {
				case <-done:
					return
				case intStream <- val:
				}

			}

		}()

		return intStream
	}

	multiply := func(done <-chan struct{}, intStream <-chan int, multiplier int) <-chan int {
		multipliedStream := make(chan int)
		go func() {
			defer close(multipliedStream)
			for val := range intStream {
				select {
				case <-done:
					return
				case multipliedStream <- multiplier * val:
				}
			}
		}()

		return multipliedStream
	}

	add := func(done <-chan struct{}, intStream <-chan int, additive int) <-chan int {
		addedStream := make(chan int)
		go func() {
			defer close(addedStream)

			for val := range intStream {
				select {
				case <-done:
					return
				case addedStream <- additive * val:
				}
			}
		}()

		return addedStream
	}

	done := make(chan struct{})
	time.AfterFunc(time.Microsecond*38, func() { close(done) })

	intStream := generator(done, []int{1, 2, 3, 4, 5}...)

	pipelineMultiplyTail := multiply(done, add(done, intStream, 3), 10)

	for v := range pipelineMultiplyTail {
		fmt.Printf("%v ", v)
	}

	fmt.Println()

	fmt.Println("Done!")
}
