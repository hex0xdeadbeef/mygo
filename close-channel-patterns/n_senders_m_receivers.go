package main

import (
	"log"
	"math/rand/v2"
	"strconv"
	"sync"
)

func ManySendersManyReceivers() {
	const (
		maxNum = 1e4

		sNum     = 10_000
		sStopVal = 333

		rNum     = 1_000
		rStopVal = 777
	)

	dataCh := make(chan int)
	toStopCh := make(chan string, 1)
	stoppedCh := make(chan struct{})

	go func() {
		log.Println("Stopped by:", <-toStopCh)
		close(stoppedCh)
	}()

	var (
		wgS sync.WaitGroup
		wgR sync.WaitGroup
	)

	for i := range sNum {
		wgS.Add(1)
		go func(id string) {
			defer wgS.Done()

			for {
				v := rand.IntN(maxNum)
				if v == sStopVal {
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

	for i := range rNum {
		wgR.Add(1)
		go func(id string) {
			defer wgR.Done()

			select {
			case <-stoppedCh:
				return
			default:
			}

			select {
			case <-stoppedCh:
				return
			case v := <-dataCh:
				if v == rStopVal {
					select {
					case toStopCh <- "receiver-" + id:
					default:
					}

					return
				}

				log.Println(v)
			}

		}(strconv.Itoa(i))
	}

	wgR.Wait()
}
