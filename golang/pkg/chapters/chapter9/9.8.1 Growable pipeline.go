package chapter9

import (
	"log"
	"sync"
)

func CreatePipeline() {
	var (
		head = make(chan int)
		prev = head
		wg   sync.WaitGroup
	)

	wg.Add(1)
	go func() {
		wg.Done()
		for i := 1; ; i++ {
			log.Print(i)
			head <- i
		}
	}()

	for {

		output := make(chan int)

		go func(prev, output chan int) {

			for prevValue := range prev {

				output <- prevValue
			}

		}(prev, output)

		prev = output
	}
}
