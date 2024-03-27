package confinement

import (
	"bytes"
	"fmt"
	"sync"
)

/*
There are two kinds of the confinement:
	1) ad hoc. ad hoc confinement is when you achieve confinement through a convention - whether it be set by the languages community, the group you work within, or the
	codebase you work within.
	2) lexical. lexical confinement involves using lexical scope to expose only the correct data and concurrency primitives for multiple concurrent processes to use. It
	makes impossible to do the wrong thing.
*/

// Ad hoc
func AdHocExample() {
	data := make([]int, 4)

	loopData := func(handleData chan<- int) {
		defer close(handleData)

		for i := range data {
			handleData <- data[i]
		}

	}

	handleData := make(chan int)
	go loopData(handleData)

	for num := range handleData {
		fmt.Println(num)
	}
}

func LexicalExample() {
	// Create the owner of a channel that will provide the channel with elements
	chanOwner := func() <-chan int {
		results := make(chan int, 5)
		go func() {
			defer close(results)
			for i := 0; i < 100; i++ {
				results <- i
			}

		}()

		// Return the receiver channel
		return results
	}

	// The logic for a receiver
	consumer := func(workerId int, results <-chan int) {
		for result := range results {
			fmt.Printf("Received %d by worker №%d\n", result, workerId)
		}
	}

	// Create an owner starting giving data
	results := chanOwner()

	wg := sync.WaitGroup{}

	// Create 10 workers that will work with the reader-channel
	for i := 10; i > 0; i-- {
		wg.Add(1)
		go func() {
			defer wg.Done()
			consumer(i, results)
		}()
	}

	wg.Wait()
	fmt.Println("Receiving has done!")
}

func LexicalWithStruct() {
	printData := func(wg *sync.WaitGroup, data []byte) {
		defer wg.Done()

		var buf bytes.Buffer
		for _, b := range data {
			fmt.Fprintf(&buf, "%c", b)
		}

		fmt.Println(buf.String())
	}

	var wg sync.WaitGroup
	wg.Add(2)

	data := []byte("golang")
	go printData(&wg, data[:3])
	go printData(&wg, data[3:])

	wg.Wait()

}
