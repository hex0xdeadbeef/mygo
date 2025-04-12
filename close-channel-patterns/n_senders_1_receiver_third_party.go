package main

import (
	"log"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

func ManySendersOneReceiverThirdPartyClose() {

	log.SetFlags(0)

	const (
		Max             = 1000000
		NumReceivers    = 10
		NumSenders      = 1000
		NumThirdParties = 15
	)

	wgReceivers := sync.WaitGroup{}
	wgReceivers.Add(NumReceivers)

	dataCh := make(chan int)
	middleCh := make(chan int)
	closing := make(chan string)
	closed := make(chan struct{})

	var stoppedBy string

	stop := func(by string) {
		select {
		case closing <- by:
			<-closed
		case <-closed:
		}
	}

	// Middle layer
	go func() {
		exit := func(v int, needSend bool) {
			close(closed)
			if needSend {
				dataCh <- v
			}
			close(dataCh)
		}

		for {
			select {
			case stoppedBy = <-closing:
				exit(0, false)
				return
			case v := <-middleCh:
				select {
				case stoppedBy = <-closing:
					exit(v, true)
					return
				case dataCh <- v:
				}
			}
		}
	}()

	// Third Parties
	for i := range NumThirdParties {
		go func(id string) {
			r := 1 + rand.Intn(3)
			time.Sleep(time.Duration(r) * time.Second)
			stop("3rd-party#" + id)
		}(strconv.Itoa(i))
	}

	// Senders
	for i := range NumSenders {
		go func(id string) {
			for {
				value := rand.Intn(Max)
				if value == 0 {
					stop("sender#" + id)
					return
				}

				select {
				case <-closed:
					return
				default:
				}

				select {
				case <-closed:
					return
				case middleCh <- value:
				}
			}
		}(strconv.Itoa(i))
	}

	// Receivers
	for range NumReceivers {
		go func() {
			defer wgReceivers.Done()

			for value := range dataCh {
				log.Println(value)
			}
		}()
	}

	wgReceivers.Wait()
	log.Println("stopped by", stoppedBy)
}
