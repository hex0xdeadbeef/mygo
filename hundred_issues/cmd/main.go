package main

import (
	"fmt"
	"log"
	"math/rand/v2"
	"strconv"
	"sync"
	"time"
)

func main() {
	// OneSenderManyReceivers()

	// OneSendersManyReceiverThirdPartyClose()

	// ManySendersOneReceiver()

	// ManySendersManyReceivers()
}

func OneSenderManyReceivers() {
	log.SetFlags(1)
	const (
		maxNum  = 1e4
		stopNum = 333

		rNum = 10000
	)

	var (
		wgR sync.WaitGroup
	)

	dataCh := make(chan int)

	go func() {
		defer close(dataCh)

		for {
			v := rand.IntN(maxNum)
			if v == stopNum {
				return
			}

			dataCh <- v
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

func ManySendersOneReceiver() {
	log.SetFlags(1)
	const (
		maxNum = 1e4

		sNum     = 10000
		sStopNum = 333
	)

	var (
		wgS sync.WaitGroup
		wgR sync.WaitGroup
	)

	dataCh := make(chan int)
	toStopCh := make(chan string, 1)
	stopCh := make(chan struct{})

	go func() {
		log.Printf("Stopped by: %s", <-toStopCh)
		close(stopCh)
	}()

	for i := range sNum {
		wgS.Add(1)
		go func() {
			defer wgS.Done()

			for {
				v := rand.IntN(maxNum)
				if v == sStopNum {
					select {
					case toStopCh <- "receiver" + strconv.Itoa(i):
					default:
					}

					return
				}

				select {
				case <-stopCh:
					return
				default:
				}

				select {
				case <-stopCh:
					return
				case dataCh <- v:
				}
			}

		}()
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

func ManySendersManyReceivers() {
	log.SetFlags(1)
	const (
		maxNum = 1e4

		rNum     = 1000
		rStopNum = 333

		sNum     = 10000
		sStopNum = 999
	)

	var (
		wgR sync.WaitGroup
		wgS sync.WaitGroup
	)

	dataCh := make(chan int)
	toStopCh := make(chan string, 1)
	stopCh := make(chan struct{})

	done := make(chan struct{})

	// Moderator
	go func() {
		log.Printf("Stopped by: %s", <-toStopCh)
		close(stopCh)
	}()

	// Senders
	for i := range sNum {
		wgS.Add(1)
		go func(i int) {
			wgS.Done()

			for {
				v := rand.IntN(maxNum)
				if v == sStopNum {
					select {
					case toStopCh <- "sender" + strconv.Itoa(i):
					default:
					}
					return
				}

				select {
				case <-stopCh:
					return
				default:
				}

				select {
				case <-stopCh:
					return
				case dataCh <- v:
				}
			}
		}(i)

	}

	// Receivers
	for i := range rNum {
		wgR.Add(1)
		go func(i int) {
			defer wgR.Done()

			for {
				select {
				case <-stopCh:
					return
				default:
				}

				select {
				case <-stopCh:
					return
				case v := <-dataCh:
					if v == rStopNum {
						select {
						case toStopCh <- "receiver" + strconv.Itoa(i):
						default:
						}

						return
					}

					log.Println(v)
				}
			}

		}(i)
	}

	go func() {
		wgS.Wait()
		wgR.Wait()
		close(done)
	}()

	<-done

}

// ???
func ManySendersOneReceiverThirdPartyClose() {
	log.SetFlags(0)

	const (
		maxNum = 1e4

		rNum            = 10
		sNum            = 1000
		thirdPartiesNum = 15
	)

	var wgR sync.WaitGroup

	var stoppedBy string

	dataCh := make(chan int)
	middleCh := make(chan int)
	closing := make(chan string)
	closed := make(chan struct{})

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
			case v := <-middleCh:
				select {
				case stoppedBy = <-closing:
					exit(v, true)
				case dataCh <- v:
				}
			}
		}
	}()

	// Third Parties
	for i := range thirdPartiesNum {
		go func(id string) {
			time.Sleep(time.Duration(1+rand.IntN(3)) * time.Second)
			stop("3rd-party" + id)
		}(strconv.Itoa(i))
	}

	// Senders
	for i := range sNum {
		go func(id string) {
			for {
				v := rand.IntN(maxNum)
				if v == 0 {
					stop("sender-" + id)
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
				case middleCh <- v:
				}
			}
		}(strconv.Itoa(i))
	}

	// Receiver
	for range rNum {
		wgR.Add(1)
		go func() {
			wgR.Done()

			for v := range dataCh {
				log.Println(v)
			}

		}()
	}

	wgR.Wait()
	fmt.Println("stopped by", stoppedBy)

}
