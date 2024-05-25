package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

func main() {
	// closedChannel()

	// rangingOverChannel()

	// selectUsingA()

	// selectUsingB()

	// selectUsingC()

	// nilChanUsingA()

	// nilChanUsingB()

	// waitManyUsing()

	// nilChanUsingC()

	// nilChanUsingD()

	// closedChanUsingA()

	// closedChanUsingB()

	// closedChanUsingС()

	// iceCreamMakers()

	// workStealingExample()

	// timeOutControlPatternUsage()

	// OrchestrationPatternUsage()

	workCancellationPatternUsage()

}

// Практика Go — Concurrency. - Хабр: https://habr.com/ru/articles/759584/
func closedChannel() {
	ch := make(chan bool, 2)
	ch <- true
	ch <- true

	close(ch)

	for i := 0; i < cap(ch)+1; i++ {
		v, ok := <-ch
		fmt.Println(v, ok)
	}

	// ch <- false // "send on closed channel"
}

func rangingOverChannel() {
	ch := make(chan bool, 2)
	ch <- true
	ch <- true

	close(ch)

	for v := range ch {
		fmt.Println(v)
	}
}

func selectUsingA() {
	var (
		finish = make(chan bool)
		wg     sync.WaitGroup
	)

	wg.Add(1)
	go func() {
		defer wg.Done()

		select {
		case <-time.After(time.Hour):
		case <-finish:
		}
	}()

	t0 := time.Now()
	finish <- true

	wg.Wait()

	fmt.Printf("Waited %v for goroutine to stop\n", time.Since(t0))

}

func selectUsingB() {
	const (
		n = 100
	)

	var (
		finish = make(chan bool)
		wg     sync.WaitGroup
	)

	for i := 0; i <= 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			select {
			case <-time.After(time.Hour):
			case <-finish:
			}
		}()
	}

	start := time.Now()
	close(finish)

	wg.Wait()

	fmt.Printf("Waited %v for %d goroutines to stop\n", time.Since(start), n)
}

func selectUsingC() {
	var (
		finish = make(chan struct{})
		wg     sync.WaitGroup
	)

	wg.Add(1)
	go func() {
		defer wg.Done()

		select {
		case <-time.After(time.Hour):
		case <-finish:
		}
	}()

	start := time.Now()
	close(finish)

	wg.Wait()

	fmt.Printf("Waited %v for goroutine to stop\n", time.Since(start))
}

func nilChanUsingA() {
	var (
		ch chan bool
	)

	// Blocks forever
	ch <- true
}

func nilChanUsingB() {
	var (
		ch chan bool
	)

	// Blocks forever
	<-ch
}

func waitManyChannelsToBeClosed(a, b <-chan bool) {
	for a != nil || b != nil {
		select {
		case <-a:
			a = nil

		case <-b:
			b = nil

		}
	}
}

func waitManyUsing() {
	var (
		aCh, bCh = make(chan bool), make(chan bool)
	)
	start := time.Now()

	go func() {
		time.Sleep(time.Millisecond * 250)
		close(aCh)
		close(bCh)
	}()

	waitManyChannelsToBeClosed(aCh, bCh)

	fmt.Printf("Waited %v for goroutines a,b to stop\n", time.Since(start))
}

func nilChanUsingC() {
	var (
		ch chan bool
	)

	// Blocks forever
	fmt.Println(<-ch)
}

func nilChanUsingD() {
	var (
		ch chan string
	)

	// Blocks forever
	ch <- "let's get started"
}

// func closedChanUsingA() {
// 	var (
// 		ch = make(chan int, 100)
// 	)

// 	for i := 0; i < 10; i++ {
// 		go func() {
// 			for j := 0; j < 10; j++ {
// 				ch <- j
// 			}

// 			close(ch)
// 		}()
// 	}

// 	for v := range ch {
// 		fmt.Println(v)
// 	}
// }

func closedChanUsingB() {
	c := make(chan int, 3)

	c <- 1
	c <- 2
	c <- 3

	close(c)

	for i := 0; i < cap(c)+1; i++ {
		fmt.Printf("%d ", <-c)
	}
	fmt.Println()

}

func closedChanUsingС() {
	c := make(chan int, 3)

	c <- 1
	c <- 2
	c <- 3

	close(c)

	for v := range c {
		fmt.Print(v, " ")
	}
	fmt.Println()

}

type (
	IceCreamMaker interface {
		Hello()
	}

	Ben struct {
		name string
	}

	Jerry struct {
		name string
	}
)

func (b *Ben) Hello() {
	fmt.Printf("Ben says, \"Hello my name is %s\"\n", b.name)
}

func (j *Jerry) Hello() {
	fmt.Printf("Jerry says, \"Hello my name is %s\"\n", j.name)
}

func iceCreamMakers() {
	var (
		ben   = &Ben{"Ben"}
		jerry = &Jerry{"Jerry"}

		maker IceCreamMaker = ben

		loop0, loop1 func()
	)

	loop0 = func() {
		maker = ben
		go loop1()
	}

	loop1 = func() {
		maker = jerry
		go loop0()
	}

	go loop0()

	for {
		maker.Hello()
	}

}

// Многопоточность и параллелизм в Go: Goroutines и каналы. Хабр - https://habr.com/ru/companies/mvideo/articles/778248/
func printNumbers(prefix string) {
	for i := 0; i < 5; i++ {
		fmt.Printf("%s: %d\n", prefix, i)
		time.Sleep(time.Millisecond * 1)
	}
}

func workStealingExample() {
	// Sets the max amount of kernels
	runtime.GOMAXPROCS(1)

	go printNumbers("Goroutine 1: ")
	go printNumbers("Goroutine 2: ")

	time.Sleep(time.Millisecond * 100)

}

func longOp(ch chan<- string) {
	// Long op imitation
	time.Sleep(time.Second * 2)
	ch <- "res"
}

func timeOutControlPatternUsage() {
	results := make(chan string)

	go longOp(results)

	select {
	case res := <-results:
		fmt.Println(res)
	case <-time.After(time.Second * 1):
		fmt.Println("Operation timeout.")
	}
}

func worker(id int, ch chan string) {
	for {
		time.Sleep(time.Millisecond)
		ch <- fmt.Sprintf("Worker %d: finised task", id)
	}
}

func OrchestrationPatternUsage() {
	ch1, ch2 := make(chan string), make(chan string)

	go worker(1, ch1)
	go worker(2, ch2)

	for i := 0; i < 5; i++ {
		select {
		case msg1 := <-ch1:
			fmt.Println(msg1)
		case msg2 := <-ch2:
			fmt.Println(msg2)
		case <-time.After(time.Second * 3):
			return
		}
	}
}

func workerWithCancelling(done chan struct{}) {
	for {
		select {
		case <-done:
			fmt.Println("Stop-signal gotten. Work has been finished.")
			return
		default:
			fmt.Println("Work in progress...")
			time.Sleep(time.Second * 1)
		}
	}
}

func workCancellationPatternUsage() {
	done := make(chan struct{})

	go workerWithCancelling(done)

	time.AfterFunc(time.Second*3, func() {
		close(done)
	})

	time.Sleep(time.Second * 5)
}
