package main

import (
	"log"
	"math/rand/v2"
	"sync"
	"time"
)

func OneSendersManyReceiverThirdPartyClose() {
	const (
		maxNum = 1e4

		rNum = 1_000

		thirdPartiesNum = 15
	)

	var wgR sync.WaitGroup

	dataCh := make(chan int)
	closing := make(chan struct{})
	closed := make(chan struct{})

	stop := func() {
		select {
		case closing <- struct{}{}:
			<-closed
		case <-closed:
		}
	}

	// Third Party Closers
	for range thirdPartiesNum {
		go func() {
			time.Sleep(time.Duration(1 + rand.IntN(3)))
			stop()
		}()
	}

	go func() {
		defer func() {
			close(closed)
			close(dataCh)
		}()

		for {

			select {
			case <-closing:
				return
			default:
			}

			select {
			case <-closing:
				return
			case dataCh <- rand.IntN(maxNum):
			}

		}
	}()

	for range rNum {
		wgR.Add(1)
		go func() {
			defer wgR.Done()

			for v := range dataCh {
				log.Println(v)
			}
		}()
	}

	wgR.Wait()
}
