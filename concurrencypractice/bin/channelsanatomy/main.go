package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

func main() {

	// creationNil()

	// creationMake()

	// greetUsage()

	// deadlockExample()

	// channelClosing()

	// squaresUsage()
	// squaresRangeUsage()

	// squaresPrintingWithBufferedChan()

	// capLenOfBufferedChannel()

	// senderUsage()
	// senderUsageModified()

	// readingFromClosedChannel()

	// workWithTwoGoroutines()

	// unidirectionalChannels()

	// channelTypeSwitching()

	// anonymousGoroutine()

	// channelWithChannelType()

	// selectUsageA()

	// nonBlockedChansInSelect()

	// defaultUsageA()

	// deadlockPreventing()

	// sendOnNilChan()

	// timeOuting()

	// emptySelectA()

	// wgUsageA()

	// workerPool()

	// raceCondA()
	// raceCondB()
}

func creationNil() {
	var c chan int
	fmt.Println(c)
}

func creationMake() {
	var c = make(chan int)

	fmt.Printf("type of channel c: %T\n", c)
	fmt.Printf("value of channel c: %#v\n", c)
}

func greet(done chan struct{}, c <-chan string) {
	fmt.Println(<-c)

	done <- struct{}{}
}

func greetUsage() {
	fmt.Println("main() started")

	done := make(chan struct{})
	c := make(chan string)

	go greet(done, c)

	c <- "Hello!"

	<-done

	fmt.Println("main() stopped")
}

func deadlockExample() {
	fmt.Println("main() started")

	c := make(chan string)
	c <- "Dmitriy"

	fmt.Println("main() stopped")

	// fatal error: all goroutines are asleep - deadlock!

	// goroutine 1 [chan send]:
	// main.deadlockExample()
	//
	//	/Users/dmitriymamykin/Desktop/mygo/concurrencypractice/bin/channelsanatomy/main.go:54 +0xa4
	//
	// main.main()
	//
	//	/Users/dmitriymamykin/Desktop/mygo/concurrencypractice/bin/channelsanatomy/main.go:13 +0x1c
}

func greetTwice(done chan struct{}, c <-chan string) {
	fmt.Println(<-c)
	fmt.Println(<-c)

	done <- struct{}{}
}

func channelClosing() {
	fmt.Println("main() started")
	c := make(chan string)
	done := make(chan struct{})

	go greetTwice(done, c)
	c <- "Dmitriy"

	close(c)

	c <- "Denis"
	// "send on closed channel"

	<-done

	fmt.Println("main() stopped")

}

func squares(c chan<- int) {
	for i := 0; i < 10; i++ {
		c <- i * i
	}

	close(c)
}

func squaresUsage() {
	fmt.Println("main() started")
	defer fmt.Println("main() stopped")

	c := make(chan int)

	go squares(c)

	for {
		val, ok := <-c
		if !ok {
			fmt.Println(val, ok, "<-- loop broken")

			return
		}

		fmt.Println(val, ok)

	}

}

func squaresRangeUsage() {
	fmt.Println("main() started")
	defer fmt.Println("main() stopped")

	c := make(chan int)

	go squares(c)

	for val := range c {
		fmt.Println(val)
	}

}

func squaresPrinting(c <-chan int) {
	for i := 0; i <= 3; i++ {
		num := <-c
		fmt.Println(num * num)
	}
}

func squaresPrintingWithBufferedChan() {
	fmt.Println("main() started")
	defer fmt.Println("main() stopped")

	c := make(chan int, 3)

	go squaresPrinting(c)
	c <- 1
	c <- 2
	c <- 3
	c <- 4

}

func capLenOfBufferedChannel() {
	c := make(chan int, 3)

	c <- 1
	c <- 2

	fmt.Println(len(c), cap(c))
}

func sender(c chan int) {
	c <- 1
	c <- 2
	c <- 3
	c <- 4

	close(c)
}

func senderUsage() {
	c := make(chan int, 3)

	go sender(c)

	fmt.Printf("Len: %d, Cap: %d\n", len(c), cap(c))

	for val := range c {
		fmt.Printf("Val: %d Len: %d, Cap: %d\n", val, len(c), cap(c))
	}
}

func squaresModified(c chan int) {
	for i := 0; i < 4; i++ {
		num := <-c
		fmt.Println(num * num)
	}
}

func senderUsageModified() {
	fmt.Println("main() started")
	defer fmt.Println("main() stopped")

	c := make(chan int, 3)
	go squaresModified(c)

	fmt.Println("Active goroutines: ", runtime.NumGoroutine())
	c <- 1
	c <- 2
	c <- 3
	c <- 4 // locks here
	fmt.Println("Active goroutines: ", runtime.NumGoroutine())

	go squaresModified(c)

	fmt.Println("Active goroutines: ", runtime.NumGoroutine())
	c <- 5
	c <- 6
	c <- 7
	c <- 8 // locks here
	fmt.Println("Active goroutines: ", runtime.NumGoroutine())

}

func readingFromClosedChannel() {
	c := make(chan int, 3)
	c <- 1
	c <- 2
	c <- 3
	close(c)

	for v := range c {
		fmt.Println(v)
	}
}

func workWithTwoGoroutines() {
	var (
		square = func(c chan int) {
			fmt.Println("[square] reading")
			num := <-c
			c <- num * num
		}

		cube = func(c chan int) {
			fmt.Println("[cube] reading")
			num := <-c
			c <- num * num * num
		}
	)

	fmt.Println("[main] main() started")
	defer fmt.Println("[main] main() stopped")

	squareCh, cubeCh := make(chan int), make(chan int)

	go square(squareCh)
	go cube(cubeCh)

	testNum := 3

	fmt.Println("[main] sent testNum to squareChan")
	squareCh <- testNum
	fmt.Println("[main] resuming")

	fmt.Println("[main] sent testNum to cubeChan")
	cubeCh <- testNum
	fmt.Println("[main] resuming")

	fmt.Println("[main] reading from channels")
	squaredVal, cubedVal := <-squareCh, <-cubeCh

	sum := squaredVal + cubedVal
	fmt.Println("[main] sum of squared and cubed nums:", sum)
}

func unidirectionalChannels() {
	roc := make(<-chan int)
	soc := make(chan<- int)

	fmt.Printf("roc data type is %T\n", roc)
	fmt.Printf("soc data type is %T\n", soc)
}

func channelTypeSwitching() {
	var (
		greet = func(done chan<- struct{}, roc <-chan string) {
			fmt.Println(<-roc)
			done <- struct{}{}
		}
	)

	fmt.Println("main() started")
	defer fmt.Println("main() stopped")

	c := make(chan string)
	done := make(chan struct{})

	go greet(done, c)

	c <- "Hello"
	<-done
}

func anonymousGoroutine() {
	var (
		c    = make(chan string)
		done = make(chan struct{})
	)

	fmt.Println("main() started")
	defer fmt.Println("main() stopped")

	go func(done chan<- struct{}, c <-chan string) {
		fmt.Println(<-c)
		done <- struct{}{}

	}(done, c)

	c <- "Hello"
	<-done
}

func channelWithChannelType() {
	var (
		greet = func(done chan struct{}, c chan string) {
			fmt.Println(<-c)
			done <- struct{}{}
		}

		greeter = func(cc chan chan string) {
			c := make(chan string)
			cc <- c
		}
	)

	fmt.Println("main() started")
	defer fmt.Println("main() stopped")

	cc := make(chan chan string)

	go greeter(cc)

	done := make(chan struct{})
	c := <-cc

	go greet(done, c)

	c <- "Hello!"
	<-done
}

func selectUsageA() {
	var (
		t time.Time

		initTime = func(t *time.Time) {
			*t = time.Now()
		}

		serviceA = func(c chan string) {
			// time.Sleep(time.Second * 3)
			c <- "Hello A!"
		}

		serviceB = func(c chan string) {
			// time.Sleep(time.Second * 5)
			c <- "Hello B!"
		}
	)

	fmt.Println("main() started")
	defer fmt.Println("main() stopped")

	chA, chB := make(chan string), make(chan string)
	initTime(&t)

	go serviceA(chA)
	go serviceB(chB)

	select {
	case res := <-chA:
		fmt.Println("Response from chA", res, time.Since(t))
	case res := <-chB:
		fmt.Println("Response from chB", res, time.Since(t))
	}
}

func nonBlockedChansInSelect() {
	var (
		t time.Time

		initTime = func(t *time.Time) {
			*t = time.Now()
		}
	)

	fmt.Println("main() started")
	defer fmt.Println("main() stopped")

	chA, chB := make(chan string, 2), make(chan string, 2)
	initTime(&t)

	chA <- "Val1"
	chA <- "Val2"

	chB <- "Val1"
	chB <- "Val2"

	select {
	case res := <-chA:
		fmt.Println("Response from chA", res, time.Since(t))
	case res := <-chB:
		fmt.Println("Response from chA", res, time.Since(t))
	}
}

func defaultUsageA() {
	var (
		t time.Time

		initTime = func(t *time.Time) {
			*t = time.Now()
		}

		serviceA = func(c chan string) {
			time.Sleep(time.Second * 3)
			c <- "Hello A!"
		}

		serviceB = func(c chan string) {
			time.Sleep(time.Second * 5)
			c <- "Hello B!"
		}
	)

	fmt.Println("main() started")
	defer fmt.Printf("main() stopped\n")

	chA, chB := make(chan string), make(chan string)
	initTime(&t)

	go serviceA(chA)
	go serviceB(chB)

Out:
	for {
		select {
		case res := <-chA:
			fmt.Println("Response from chA", res, time.Since(t))
			break Out
		case res := <-chB:
			fmt.Println("Response from chB", res, time.Since(t))
			break Out
		default:
			fmt.Println("No response received!", "Sleep for 500 ms...")
			time.Sleep(time.Millisecond * 500)
		}
	}

}

func deadlockPreventing() {
	fmt.Println("main() started")
	defer fmt.Printf("main() stopped\n")

	chA, chB := make(chan string), make(chan string)
	select {
	case res := <-chA:
		fmt.Println("Response from chA", res)
	case res := <-chB:
		fmt.Println("Response from chA", res)
	default:
		fmt.Println("All the channels are empty. Returning!")
	}

}

func sendOnNilChan() {
	var (
		service = func(c chan string) {
			// goroutine 20 [chan send (nil chan)]:
			// main.sendOnNilChan.func1(0x0)
			// 	/Users/dmitriymamykin/Desktop/mygo/concurrencypractice/bin/channelsanatomy/main.go:517 +0x28
			// created by main.sendOnNilChan in goroutine 1
			// 	/Users/dmitriymamykin/Desktop/mygo/concurrencypractice/bin/channelsanatomy/main.go:526 +0x160
			c <- "resp"
		}
	)

	fmt.Println("main() started")
	defer fmt.Printf("main() stopped\n")

	var c chan string

	go service(c)

	select {
	// fatal error: all goroutines are asleep - deadlock!

	// goroutine 1 [chan receive (nil chan)]:
	// main.sendOnNilChan()
	// 	/Users/dmitriymamykin/Desktop/mygo/concurrencypractice/bin/channelsanatomy/main.go:529 +0x178
	// main.main()
	// 	/Users/dmitriymamykin/Desktop/mygo/concurrencypractice/bin/channelsanatomy/main.go:51 +0x1c
	case res := <-c:
		fmt.Println("Response received:", res)
	default:
		fmt.Println("No resp!")
	}

}

func timeOuting() {
	var (
		t time.Time

		initTime = func(t *time.Time) {
			*t = time.Now()
		}

		serviceA = func(c chan string) {
			time.Sleep(time.Second * 3)
			c <- "Hello A!"
		}

		serviceB = func(c chan string) {
			time.Sleep(time.Second * 5)
			c <- "Hello B!"
		}
	)

	fmt.Println("main() started")
	defer fmt.Printf("main() stopped\n")

	chA, chB := make(chan string), make(chan string)
	initTime(&t)

	go serviceA(chA)
	go serviceB(chB)

	select {
	case res := <-chA:
		fmt.Println("Response from chA", res, time.Since(t))
	case res := <-chB:
		fmt.Println("Response from chB", res, time.Since(t))
		// case <-time.After(time.Second * 2):
	// fmt.Println("Timeout. No response received!")
	case <-time.After(time.Second * 10):
		fmt.Println("Timeout. No response received!")
	}
}

func emptySelectA() {
	fmt.Println("main() started")
	defer fmt.Printf("main() stopped\n")

	// 	fatal error: all goroutines are asleep - deadlock!

	// goroutine 1 [select (no cases)]:
	// main.emptySelectA()
	// 	/Users/dmitriymamykin/Desktop/mygo/concurrencypractice/bin/channelsanatomy/main.go:597 +0xac
	// main.main()
	// 	/Users/dmitriymamykin/Desktop/mygo/concurrencypractice/bin/channelsanatomy/main.go:55 +0x1c
	select {}
}

func wgUsageA() {
	var (
		resource int

		service = func(id int, mu *sync.Mutex, wg *sync.WaitGroup) {
			defer wg.Done()

			mu.Lock()
			resource++
			fmt.Println("ID:", id, "Resource val:", resource)
			mu.Unlock()
		}

		wg = &sync.WaitGroup{}
		mu = &sync.Mutex{}
	)

	fmt.Println("main() started")
	defer fmt.Printf("main() stopped\n")

	for i := 0; i < 3; i++ {
		i := i

		wg.Add(1)
		go service(i, mu, wg)
	}

	wg.Wait()
}

func workerPool() {
	type (
		Result struct {
			sq int
			id int
		}
	)

	var (
		sqrWorker = func(wg *sync.WaitGroup, id int, tasks <-chan int, results chan<- Result) {
			var (
				sendRes = func(id int, results chan<- Result, resNum int) {
					time.Sleep(time.Millisecond * 1)
					fmt.Printf("Worker [%d] sending result\n", id)
					results <- Result{sq: resNum, id: id}
				}
			)
			defer wg.Done()

			for num := range tasks {
				sendRes(id, results, num*num)
			}

		}

		wg = &sync.WaitGroup{}
	)

	fmt.Println("main() started")
	defer fmt.Printf("main() stopped\n")

	tasks, results := make(chan int), make(chan Result, 10)

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(id int) {
			sqrWorker(wg, id, tasks, results)
		}(i)
	}

	for i := 0; i < 5; i++ {
		tasks <- i * 2
	}
	close(tasks)

	wg.Wait()

	for i := 0; i < 5; i++ {
		fmt.Printf("R %+v\n", <-results)
	}

	fmt.Println("[main] wrote 5 tasks")

}

func raceCondA() {
	var (
		i int

		worker = func(wg *sync.WaitGroup) {
			i++
			wg.Done()
		}

		wg = &sync.WaitGroup{}
	)

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go worker(wg)
	}

	wg.Wait()

	fmt.Println(i)

}

func raceCondB() {
	var (
		i int

		worker = func(mu *sync.Mutex, wg *sync.WaitGroup) {
			mu.Lock()
			defer mu.Unlock()

			i++
			wg.Done()
		}

		wg = &sync.WaitGroup{}
		mu = &sync.Mutex{}
	)

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go worker(mu, wg)
	}

	wg.Wait()

	fmt.Println(i)

}
