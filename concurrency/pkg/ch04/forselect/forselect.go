package forselect

import (
	"fmt"
	"time"
)

func LoopOverDataSequence() {
	var (
		stringData = []string{"fox", "a", "b", "c", "d", "f",
			"fox", "a", "b", "c", "d", "f",
			"fox", "a", "b", "c", "d", "f",
			"fox", "a", "b", "c", "d", "f",
			"fox", "a", "b", "c", "d", "f",
			"fox", "a", "b", "c", "d", "f",
			"fox", "a", "b", "c", "d", "f",
			"fox", "a", "b", "c", "d", "f",
			"fox", "a", "b", "c", "d", "f"}

		done         = make(chan struct{})
		stringStream = make(chan string)
	)

	time.AfterFunc(time.Microsecond*5, func() {
		close(done)
	})

	go func() {
		for s := range stringStream {
			fmt.Printf("%s ", s)
		}
	}()

loop:
	for _, s := range stringData {
		select {
		case <-done:
			break loop
		case stringStream <- s:
			/*
				Do work in another goroutine
			*/
		}
	}
	close(stringStream)

	fmt.Println("\nPrinting has done!")

}

func InfiniteLoopingWithWorkInLoop() {
	done := make(chan struct{})

	time.AfterFunc(time.Microsecond*5, func() {
		close(done)
	})

	i := 0
	for {
		select {
		// Stops the work after done is closed
		case <-done:
			return
		// Non-blocking continuation
		default:
		}
		/*
			Do non-preemptable work
		*/
		fmt.Printf("%c ", 'A'+i)
		i++
	}

}

func InfiniteLoopingWithWorkInDefault() {
	done := make(chan struct{})

	time.AfterFunc(time.Microsecond*5, func() {
		close(done)
	})

	i := 0
	for {
		select {
		case <-done:
			return
		default:
			fmt.Printf("%c ", 'A'+i)
			i++
		}

	}
}
