package main

import (
	"log"
	"math/rand/v2"
	"sync"
)

func OneSenderManyReceivers() {

	const (
		maxNum = 10_000
		rNum = 1_000
		stopVal = 333
	)

	dataCh := make(chan int)

	var wgR sync.WaitGroup
	for range rNum {
		wgR.Add(1)

		go func() {
			defer wgR.Done()

			for v := range dataCh {
				log.Println(v)
			}

		}()
	}

	go func() {
		defer close(dataCh)

		for {
			v := rand.IntN(maxNum)
			if v == stopVal {
				return
			}

			dataCh <- v
		}
	}()

	wgR.Wait()
}
