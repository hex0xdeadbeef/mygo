package main

import (
	"log"
	"math/rand/v2"
	"strconv"
	"sync"
)

func ManySendersOneReceiver() {
	const (
		maxNum = 10_000

		sNum = 10_000

		stopVal = 333
	)

	dataCh := make(chan int)
	toStopCh := make(chan string, 1)
	stoppedCh := make(chan struct{})

	var wgS sync.WaitGroup
	var wgR sync.WaitGroup

	go func() {
		log.Println("Stopped by:", <-toStopCh)
		close(stoppedCh)
	}()

	for i := range sNum {
		wgS.Add(1)

		go func(id string) {
			defer wgS.Done()
			for {
				v := rand.IntN(maxNum)
				if v == stopVal {
					select {
					case toStopCh <- "sender-" + id:
					default:
					}

					return
				}

				select {
				case <-stoppedCh:
					return
				default:
				}

				select {
				case <-stoppedCh:
					return
				case dataCh <- v:
				}
			}
		}(strconv.Itoa(i))
	}

	go func() {
		wgS.Wait()
		close(dataCh)
	}()

	wgR.Add(1)
	go func() {
		defer wgR.Done()

		for v := range dataCh {
			log.Println(v)
		}
	}()

	wgR.Wait()
}
