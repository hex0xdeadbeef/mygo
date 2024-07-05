package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"golang.org/x/sync/errgroup"
)

/*
	CONCURRENCY PRACTICE

	PASSING INAPPROPRIATE CONTEXT
1. The context bound to a HTTP request can be canceled in the following conditions:
	1) The client conn is being closed
	2) In the case of a HTTP/2 request, when it's closed

	3) When the response has been sent to a client

	NOT TO STOP A GOROUTINE
0. A goroutine is the same resource as a file | conn and it must be closed respectively.
1. The max size of goroutine is 1GB for 64-bit system and 250MB for 32-bit system.
2. We should chase leaked goroutines so that we  won't get leaks of resourses as: HTTP-conns | DB conns | Opened files | Network Sockets and others.
3. The methods to stop a goroutine:
	1) Using Context pkg. In this approach we need to remember about closing all the resources allocated to an object.

	TREAT GOROUTINES AND CYCLE VARIABLES CARELESSLY
1. Referencing to a cycle within a closure that is used to be launched as a G is an error. The closures observe current value of a variable.
	1) We need to explicitly pass all the args to a closure.
	2) Create a local copy at each iteration

	EXPECT DETERMINISTIC BEHAVIOR FROM SELECT STATEMENT
1. select operator chooses a case randomly. It's made to prevent starvation of goroutines.
2. Channels guarantee the order of sending/receiving. The messages will be gotten in the order they've been sent.
3. In the case of some goroutines-producers writing on a channel, the order of msgs receiving isn't deterministic. It results in race condition. If we've gotten a
message that tells us to disconnect we can apply the following solution:
	1) Get msgs from msgsCh | disconnectCh
	2) If the disconnect-notification gotten: Read all the messages from msgsCh if any and Return

	NOT TO USE NOTIFICATION-CHANS
1. To communicate betweeen goroutines we should use notification channels.
2. The type of channel is preferably chosen as struct{}
3. struct{} weighs 0 bytes independently of an architecture

	NOT TO USE NIL-CHANS
4. Sending/Reading on/from nil channel blocks the calling goroutine
5. Closing nil channel panics
6. Select operator never chooses nil channel

	CHANNEL SIZE CHOOSING
1. Unbuffered channel guaranteees synchronization. There's the guarantee that 2 goroutines will be in the known state: one receives and another one sends.
2. Buffered channel provides light synchronization. The goroutine won't get a message before its receiving.
3. We should start from the default size of buffered channel that is 1.
4. How and when should we choose the size of buffered channel?
	1) Using Worker Pool pattern. The size can be value of the goroutines amount | GOMAXPROCS
	2) Using this chan in terms of limitation on execution speed (Semaphore). The size is choosen by demands.
5. The code mustn't have magical numbers.
6. Queues, as rule, are always close to fullness or emptiness because of difference between speed of work of consumers and producers. They're rarely work in
balance (paces are equal). "Martin Thompson"
7. An optimal size can be choosen based on benchmarks
8. Buffered chans lead to mutual locks.

	STRINGS FORMATTING SIDE EFFECTS
1. Mutual Locks. Using RWMutex we can force fmt.XxxF() to use String() method on a receiver with the "%s" directive. The situation is: in a different place we
apply Lock() on the mutex. And when we call RLock() on the same mutex, the calling goroutine will'be locked. The ways to solve it:
	1) Make all the checks before getting mutex.
	2) Use formatting methods using other directives.

	RACE CONDITIONS WTIH APPEND
1. If we have the filled slice (len == capacity) and goroutines creates slices using append, there won't be a datarace. We don't write into the source structure.
We create a new one. A new underlying array will be created and assigned to a new structure.
2. If len != cap and two goroutines append elems, there will a datarace because 2 goroutines can write to the same memory area simultaneously. The solution is:
	1) To create a deep copy firstly and append elems needed.

If the slice is filled fully (len == cap) append will not result in race presence. Otherwise (len < cap) there will be a datarace.
3. Dataraces' affection on slices/maps. When there are some Gs the following is right:
	1) Access to the same index by at least G updating the corresponding value - is the situation with a datarace. Gs clobber the same memory area.
	2) Access to different indices of slices don't result in a datarace. Different indices means an affection on different memory areas.
	3) Access to the same map (independently from whether the key is the same or different) by at least a single Goroutine updating this map results in datarace.
	The difference between slice and map: map is an array of segments and each segment is the pointer to an array of pairs (key -> value). A hash algorithm is used
	to define an index in the array of the segment. Since the hash algorithm includes a random part during tha map initialization, an execution can lead to the
	same index in array and another cannot. Race detector processes this case giving an alert independently from whether there's a factual race.

	INCORRECT WAY TO USE MUTEXES ON SLICES/MAPS
1. If we want to make processing over fixed state of a slice/map in separate Gs, we need to create a deep copy of elems inside it. The simple reassigning doesn't
make copy, we just copy the structure and reference to the same object. The solutions are:
	1) Make a deep copy of the structure needed. This method allocates more memory, but results in faster execution time. This method is useful when the operation
	over internal structure isn't fast.
	2) Restrict all the function with calculation with RWMutex. This method results in more execution time.

	INCORRECT USE OF sync.WaitGroup
1. If the internal counter gets negative value - a G will panic
2. We should remember that execution of Gs isn't deterministic without synchronization. Two goroutines can be set to different OS Threads (M) and we cannot
guarantee what thread will be performe first.
3. CPU must use memory (fence/barrier) in order to define the order. The memory barrier can be established by wg.WaitGroup, it established the relation between
wg.Add(...) and wg.Done().
4. We can use the wg.Add() in the following ways:
	1) Add an amount if it's known before the execution
	2) If we're in a cycle, we increment the internal counter before launching every goroutine

	USAGE OF sync.Cond
1. sync.Cond is used when we want to spread out the signal to multiple Gs many times.
2. If we try to imitate the sync.Cond with the looping over a condition, we creates busy waiting that consumes many CPU ticks and loads the CPU itself.
3. If we try to imitate the sync.Cond using channels we run into the following problem: the signal from a chan will be consumed by a single receiver, not all.
Multiple goroutines can only get channel closing message simultaneously, but the single msg cannot.
4. sync.Cond is used when we need to distribute a signal to multiple Gs. This method allows to avoid busy-waiting of Gs.
5. Mechanism of Wait:
	1) Wait() unlocks the Locker()
	2) Suspends the caller (G) knocking it from thread
	3) After Broadcast()/Signal() Wait() blocks the caller and allows the G to observe the condition
6. Observation / Updating the condition requires the Lock() of Locker
7. The lack of sync.Cond is:
	1) When we use chans and there are no active receivers, there's the guarantee that channel makes bufferization. Finally this signal will be received.
	2) If there's no Gs waiting for condition, the signal will be skipped.

	USING errgroup.Group
1. The benefit of errgroup is the mutual context of all the goroutines in group.

	COPYING OF sync OBJECTS
1. We cannot copy all the objects of sync package
2. The soluitons:
	1) Turn the type of a receiver to pointer
	2) Make a filed of pointer type. In this case we must also change the initialization of instance.
*/

func main() {
	// CycleVarsAndG()

	// SelectErrorsA()

	// SelectErrorsB()

	// TwoProducersOnTheSameChan()

	// EmptyStructSizeInBytes()

	// NilChanA()

	// NilChanB()

	// NilChanC()

	// channelClosing()

	// MergeUsage()

	// MutualLock()

	// DataRaceAppendA()

	// DataRaceAppendB()

	// DataRaceAppendC()

	// DataRaceAppendD()

	// SliceMapDataraceA()

	// SliceMapDataraceB()

	// SliceMapDataraceC()

	// ConcurrentMutexedStructureA()

	// WaitGroupNegativeInternalCounter()

	// IncorrectWaitGroupUsageA()

	// NonDeterministicExecution()

	// SyncCondSimulationA()

	// SyncCondSimulationB()

	// SyncCondSimulationC()

	// ErrGroupContextUsage()
}

func InappropriateContextPassing() {
	var (
		doSomeTask = func(ctx context.Context, r *http.Request) (http.Response, error) {
			// some logic ...

			return http.Response{}, nil
		}

		publish = func(ctx context.Context, resp http.Response) error {
			// some logic ...

			return nil
		}

		writeResp = func(resp http.Response) {
			// some logic ...
		}

		// The context of request can be canceled in the following conditions:
		// 1. Client's conn is being closed
		// 2. In the case of a HTTP 2.0, when it's closed
		// 3. When an answer was directed to a client.
		handler = func(w http.ResponseWriter, r *http.Request) {
			// The task of execution of a task on the request
			response, err := doSomeTask(r.Context(), r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// The goroutine that publishes the response asynchronously

			// 1. If the client conn is being closed, we go into the publish(...) and return because the context is already closed.
			// 2. In the case of a HTTP 2.0 request when it's been canceled, it's the same to 1.

			// 3. When an aswer was sent to a client back, there will be a race condition.
			// A
			go func() {
				// err := publish(r.Context(), response)

				// To beat the datarace we should pass the empty context.
				if err := publish(context.TODO(), response); err != nil {
					log.Println(err)
				}

			}()
			// A

			// HTTP response writing
			// B
			writeResp(response)
			// B

			// There's no guarantee of A < B, so there can be the datarace
		}
	)

	handler(nil, nil)

}

type watcher struct{}

func NewWatcher() watcher {
	return watcher{}
}

func (w watcher) Close() error {
	// some logic ...
	return nil
}

func (w watcher) watch(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			fmt.Println("I'm watching")
		}
	}
}

func GoroutineLeak() {
	w := NewWatcher()
	/*
		Without closing the resources of watcher created we have a possibility to leak the resources allocated.
		We need explicitly lock the parent goroutine.
	*/
	defer w.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// The function will return after the outer function returning
	go w.watch(ctx)
}

func CycleVarsAndG() {
	var (
		s = []int{1, 2, 3}
	)

	for _, v := range s {
		go func(val int) {
			fmt.Println(val)
		}(v)
	}

	for _, v := range s {
		v := v

		go func() {
			fmt.Println(v)
		}()
	}
}

func SelectErrorsA() {
	var (
		IDProducer = func() <-chan string {

			ch := make(chan string)

			go func() {
				var (
					timer = time.NewTimer(time.Millisecond * 100)
				)
				defer timer.Stop()
				defer close(ch)

				for {
					select {
					case <-timer.C:
						return
					case ch <- genID():
					}
				}
			}()

			return ch
		}

		disconnectProd = func() chan struct{} {
			ch := make(chan struct{})

			time.AfterFunc(time.Millisecond*300, func() { close(ch) })

			return ch
		}
	)

	msgProd := IDProducer()
	disconnectCh := disconnectProd()
	storage := make([]string, 0, 1<<10)

	for {
		select {
		case msg := <-msgProd:

			storage = append(storage, msg)
		case <-disconnectCh:
			return
		}
	}

}

type Producer struct {
	msgs chan string
}

func NewProducer() Producer {
	return Producer{msgs: make(chan string)}
}

func (p Producer) Msgs() <-chan string {
	return p.msgs
}

func (p Producer) Close() error {
	// some logic ...
	return nil
}

func (p Producer) LaunchMessages(ctx context.Context) {
	go func() {
		var (
			timer = time.NewTimer(time.Millisecond * 100)
		)
		defer timer.Stop()
		defer close(p.msgs)

		for {
			select {
			case p.msgs <- genID():

			case <-timer.C:
				return

			case <-ctx.Done():
				return
			}
		}
	}()
}

func SelectErrorsB() {
	var (
		producer = NewProducer()
		storage  = make([]string, 0, 1<<10)
	)
	defer producer.Close()

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	producer.LaunchMessages(ctx)

Out:
	for {
		select {
		case msg, ok := <-producer.Msgs():
			if !ok {
				break Out
			}
			storage = append(storage, msg)
		}
	}

	for _, s := range storage {
		fmt.Println(s)
	}
}

func genID() string {
	const (
		alphabetSize = 26
		IDSize       = 30

		pad = 1
	)
	var (
		b = &strings.Builder{}
	)
	b.Grow(IDSize)

	for i := 0; i < IDSize; i++ {
		idx := rand.Intn(alphabetSize)
		b.WriteRune(rune('a' + idx))
	}

	return b.String()
}

func TwoProducersOnTheSameChan() {
	var (
		prod = func(ch chan string, timeoutInMillis int) {
			go func() {
				var (
					timer = time.NewTimer(time.Millisecond * time.Duration(timeoutInMillis))
				)
				defer timer.Stop()

				for {
					select {
					case ch <- genID():
					case <-timer.C:
						return
					}
				}
			}()

		}

		disconnectCh = func(timeoutInMillis int) chan struct{} {
			ch := make(chan struct{})

			time.AfterFunc(time.Millisecond*time.Duration(timeoutInMillis), func() { close(ch) })

			return ch
		}

		prodCh  = make(chan string)
		storage = make([]string, 0, 1<<10)
	)

	prod(prodCh, 200)
	prod(prodCh, 500)

	disconn := disconnectCh(300)
Out:
	for {
		select {
		case v := <-prodCh:
			storage = append(storage, v)

		case <-disconn:
			select {
			case v := <-prodCh:
				storage = append(storage, v)
				// We'll get into default when prodCh is drained
			default:
				fmt.Println("disconnect, return")
				break Out
			}
		}

	}

	for _, s := range storage {
		fmt.Println(s)
	}

}

func EmptyStructSizeInBytes() {
	var emptyStruct struct{}
	fmt.Println(unsafe.Sizeof(emptyStruct)) // 0 bytes

	var emptyInterface interface{}
	fmt.Println(unsafe.Sizeof(emptyInterface)) // 16 bytes
}

func NilChanA() {
	var ch chan int

	// Will block the calling goroutine
	ch <- 0
}

func NilChanB() {
	var ch chan int

	// Will block the calling goroutine
	<-ch
}

func NilChanC() {
	var ch chan int

	// Will panic
	close(ch)
}

func MergeUsage() {
	var (
		getIntChan = func(low, top int) <-chan int {
			ch := make(chan int)

			go func() {
				defer close(ch)

				for i := low; i <= top; i++ {
					ch <- i
				}
			}()

			return ch
		}
	)

	for v := range ChansMerge(getIntChan(10, 100), getIntChan(200, 1000)) {
		fmt.Print(v, " ")
	}
}

func ChansMerge[T any](chA, chB <-chan T) <-chan T {
	mergedCh := make(chan T)

	go func() {
		const (
			timeoutInSec = 3
		)

		var (
			timer = time.NewTimer(time.Second * timeoutInSec)
		)
		defer timer.Stop()
		defer close(mergedCh)

		for {
			select {
			case v, ok := <-chA:
				if !ok {
					chA = nil
					break
				}
				mergedCh <- v

			case v, ok := <-chB:
				if !ok {
					chB = nil
					break
				}
				mergedCh <- v

			case <-timer.C:
				if chA == nil && chB == nil {
					return
				}
			}
		}

	}()

	return mergedCh
}

func channelClosing() {
	var (
		ch = make(chan int)
	)

	close(ch)

	fmt.Println(<-ch, <-ch)
}

type Customer struct {
	mu  *sync.RWMutex
	id  string
	age int
}

// If we call the SetAge(...) in fmt.Errorf(...) the method c.String(...) will be called on the receiver
// RLock will be waiting until the Lock will be free.
func NewCustomer() *Customer {
	return &Customer{mu: &sync.RWMutex{}, id: "default", age: 18}
}

func (c *Customer) String() string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return fmt.Sprintf("id: %s, age: %d", c.id, c.age)
}

func (c *Customer) SetAgeA(age int) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if age < 0 {
		return fmt.Errorf("age is less than 0: %s", c)
	}

	c.age = age
	return nil
}

func (c *Customer) SetAgeB(age int) error {
	/*
		In this case we lock only when the age is correct
	*/
	if age < 0 {
		return fmt.Errorf("age is less than 0: %s", c)
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.age = age
	return nil
}

func (c *Customer) SetAgeC(age int) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	/*
		In this case we force the formatter no to use String() method on receiver
	*/
	if age < 0 {
		return fmt.Errorf("age is less than 0: %d", c.age)
	}

	c.age = age
	return nil
}

func MutualLock() {
	c := NewCustomer()

	if rand.Intn(2) == 1 {
		/*
			Will be locked
		*/
		c.SetAgeA(-1)
	}

	/*
		Won't be locked
	*/
	c.SetAgeB(-1)
	c.SetAgeC(-1)
}

func DataRaceAppendA() {
	var (
		s  = make([]int, 1)
		wg = &sync.WaitGroup{}
	)
	wg.Add(2)

	go func() {
		defer wg.Done()
		s1 := append(s, 4)

		fmt.Println(s1)
	}()

	go func() {
		defer wg.Done()
		s2 := append(s, 5)
		fmt.Println(s2)
	}()

	wg.Wait()
}

func DataRaceAppendB() {
	var (
		s  = make([]int, 0, 1)
		wg = &sync.WaitGroup{}
	)
	wg.Add(2)

	/*
		==================
		WARNING: DATA RACE
		Write at 0x00c0001b4128 by goroutine 7:
		  main.DataRaceAppendB.func1()
		      /Users/dmitriymamykin/Desktop/mygo/hundred_issues/internal/concurrencypractice/main.go:637 +0xcc

		Previous write at 0x00c0001b4128 by goroutine 8:
		  main.DataRaceAppendB.func2()
		      /Users/dmitriymamykin/Desktop/mygo/hundred_issues/internal/concurrencypractice/main.go:644 +0xcc

		Goroutine 7 (running) created at:
		  main.DataRaceAppendB()
		      /Users/dmitriymamykin/Desktop/mygo/hundred_issues/internal/concurrencypractice/main.go:635 +0x108
		  main.main()
		      /Users/dmitriymamykin/Desktop/mygo/hundred_issues/internal/concurrencypractice/main.go:105 +0x20

		Goroutine 8 (running) created at:
		  main.DataRaceAppendB()
		      /Users/dmitriymamykin/Desktop/mygo/hundred_issues/internal/concurrencypractice/main.go:642 +0x1b0
		  main.main()
		      /Users/dmitriymamykin/Desktop/mygo/hundred_issues/internal/concurrencypractice/main.go:105 +0x20
		==================
		[5]
		[4]
		Found 1 data race(s)
		exit status 66
	*/

	go func() {
		defer wg.Done()
		s1 := append(s, 4)

		fmt.Println(s1)
	}()

	go func() {
		defer wg.Done()
		s2 := append(s, 5)
		fmt.Println(s2)
	}()

	wg.Wait()
}

func DataRaceAppendC() {
	var (
		s  = make([]int, 0, 1)
		wg = &sync.WaitGroup{}
	)
	wg.Add(2)

	go func() {
		defer wg.Done()
		s1 := append([]int(nil), s...)
		s1 = append(s1, 4)
		fmt.Println(s1)
	}()

	go func() {
		defer wg.Done()
		s2 := append([]int(nil), s...)
		s2 = append(s2, 5)
		fmt.Println(s2)
	}()

	wg.Wait()
}

func DataRaceAppendD() {
	var (
		s  = make([]int, 0, 1)
		wg = &sync.WaitGroup{}
	)
	wg.Add(2)

	go func() {
		defer wg.Done()
		s1 := append(make([]int, len(s), cap(s)), s...)
		s1 = append(s1, 4)
		fmt.Println(s1)
	}()

	go func() {
		defer wg.Done()
		s2 := append(make([]int, len(s), cap(s)), s...)
		s2 = append(s2, 5)
		fmt.Println(s2)
	}()

	wg.Wait()
}

/*
	Dataraces in slices/maps
*/

func SliceMapDataraceA() {
	var (
		wg = &sync.WaitGroup{}
		s  = make([]int, 0, 1)
	)

	wg.Add(2)
	go func() {
		wg.Done()
		s = append(s, 1)
	}()

	go func() {
		wg.Done()
		s = append(s, 2)
	}()

	wg.Wait()

	fmt.Println(s)

	/*
		==================
		WARNING: DATA RACE
		Read at 0x00c0001b0108 by goroutine 7:
		  main.SliceMapDataraceA.func1()
		      /Users/dmitriymamykin/Desktop/mygo/hundred_issues/internal/concurrencypractice/main.go:754 +0x3c

		Previous write at 0x00c0001b0108 by goroutine 8:
		  main.SliceMapDataraceA.func2()
		      /Users/dmitriymamykin/Desktop/mygo/hundred_issues/internal/concurrencypractice/main.go:759 +0xb0

		Goroutine 7 (running) created at:
		  main.SliceMapDataraceA()
		      /Users/dmitriymamykin/Desktop/mygo/hundred_issues/internal/concurrencypractice/main.go:752 +0x158
		  main.main()
		      /Users/dmitriymamykin/Desktop/mygo/hundred_issues/internal/concurrencypractice/main.go:118 +0x20

		Goroutine 8 (finished) created at:
		  main.SliceMapDataraceA()
		      /Users/dmitriymamykin/Desktop/mygo/hundred_issues/internal/concurrencypractice/main.go:757 +0x1f8
		  main.main()
		      /Users/dmitriymamykin/Desktop/mygo/hundred_issues/internal/concurrencypractice/main.go:118 +0x20
		==================
		==================
		WARNING: DATA RACE
		Read at 0x00c0001b4138 by main goroutine:
		  reflect.maplen()
		      /opt/homebrew/Cellar/go/1.22.2/libexec/src/runtime/map.go:1406 +0x7c
		  reflect.packEface()
		      /opt/homebrew/Cellar/go/1.22.2/libexec/src/reflect/value.go:135 +0xa4
		  reflect.valueInterface()
		      /opt/homebrew/Cellar/go/1.22.2/libexec/src/reflect/value.go:1526 +0x148
		  reflect.Value.Interface()
		      /opt/homebrew/Cellar/go/1.22.2/libexec/src/reflect/value.go:1496 +0x74
		  fmt.(*pp).printValue()
		      /opt/homebrew/Cellar/go/1.22.2/libexec/src/fmt/print.go:769 +0x80
		  fmt.(*pp).printValue()
		      /opt/homebrew/Cellar/go/1.22.2/libexec/src/fmt/print.go:910 +0x1cd4
		  fmt.(*pp).printArg()
		      /opt/homebrew/Cellar/go/1.22.2/libexec/src/fmt/print.go:759 +0x8a0
		  fmt.(*pp).doPrintln()
		      /opt/homebrew/Cellar/go/1.22.2/libexec/src/fmt/print.go:1221 +0x40
		  fmt.Fprintln()
		      /opt/homebrew/Cellar/go/1.22.2/libexec/src/fmt/print.go:304 +0x48
		  fmt.Println()
		      /opt/homebrew/Cellar/go/1.22.2/libexec/src/fmt/print.go:314 +0x278
		  main.SliceMapDataraceA()
		      /Users/dmitriymamykin/Desktop/mygo/hundred_issues/internal/concurrencypractice/main.go:764 +0x204
		  main.main()
		      /Users/dmitriymamykin/Desktop/mygo/hundred_issues/internal/concurrencypractice/main.go:118 +0x20

		Previous write at 0x00c0001b4138 by goroutine 8:
		  main.SliceMapDataraceA.func2()
		      /Users/dmitriymamykin/Desktop/mygo/hundred_issues/internal/concurrencypractice/main.go:759 +0x98

		Goroutine 8 (finished) created at:
		  main.SliceMapDataraceA()
		      /Users/dmitriymamykin/Desktop/mygo/hundred_issues/internal/concurrencypractice/main.go:757 +0x1f8
		  main.main()
		      /Users/dmitriymamykin/Desktop/mygo/hundred_issues/internal/concurrencypractice/main.go:118 +0x20
		==================
		==================
		WARNING: DATA RACE
		Write at 0x00c0001b0108 by goroutine 7:
		  main.SliceMapDataraceA.func1()
		      /Users/dmitriymamykin/Desktop/mygo/hundred_issues/internal/concurrencypractice/main.go:754 +0xb0

		Previous read at 0x00c0001b0108 by main goroutine:
		  main.SliceMapDataraceA()
		      /Users/dmitriymamykin/Desktop/mygo/hundred_issues/internal/concurrencypractice/main.go:764 +0x20c
		  main.main()
		      /Users/dmitriymamykin/Desktop/mygo/hundred_issues/internal/concurrencypractice/main.go:118 +0x20

		Goroutine 7 (running) created at:
		  main.SliceMapDataraceA()
		      /Users/dmitriymamykin/Desktop/mygo/hundred_issues/internal/concurrencypractice/main.go:752 +0x158
		  main.main()
		      /Users/dmitriymamykin/Desktop/mygo/hundred_issues/internal/concurrencypractice/main.go:118 +0x20
		==================
		[2]
		Found 3 data race(s)
		exit status 66
	*/
}

// There's no datarace because we clobber different memory areas.
func SliceMapDataraceB() {
	var (
		wg = &sync.WaitGroup{}
		s  = make([]int, 1<<3)
	)

	wg.Add(2)
	go func() {
		defer wg.Done()
		s1 := s[:4]
		// s1 = s1[:0]

		s1 = append(s1, []int{2, 2, 2, 2}...)
	}()

	go func() {
		defer wg.Done()
		s2 := s[4:]
		// s2 = s2[:0]

		s2 = append(s2, []int{3, 3, 3, 3}...)
	}()

	wg.Wait()

	fmt.Println(s)
}

func SliceMapDataraceC() {
	var (
		wg = &sync.WaitGroup{}
		ch = make(chan int, 1)

		m = make(map[int]string, 1<<10)
	)

	wg.Add(2)

	// There will be the datarace because of there's the goroutine updating data concurrently
	go func() {
		defer wg.Done()
		defer close(ch)

		for i := 0; i < 1<<10; i++ {
			m[i] = genID()
			ch <- i
		}

	}()

	go func() {
		defer wg.Done()

		for v := range ch {
			m[v] = genID()
		}
	}()

	wg.Wait()
}

const (
	defaultCacheSize = 1 << 10
	multiplierBase   = 1e4
)

type Cache struct {
	mu       *sync.RWMutex
	balances map[string]float64
}

func NewCache() *Cache {
	return &Cache{mu: &sync.RWMutex{}, balances: make(map[string]float64, defaultCacheSize)}
}

func (c *Cache) Fill() {

	for i := 0; i < defaultCacheSize; i++ {
		c.balances[genID()] = rand.Float64() * float64(multiplierBase+rand.Intn(multiplierBase))
	}
}

func (c *Cache) AddBalance(ID string, newBalance float64) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.balances[ID] = newBalance
}

func (c *Cache) AvgBalanceA() float64 {
	c.mu.RLock()
	/*
		This method is silly because we don't make a deep copy
	*/
	copiedBalances := c.balances

	for k, v := range c.balances {
		copiedBalances[k] = v
	}
	c.mu.RUnlock()

	var (
		balancesSum float64
	)

	for _, clientBalance := range copiedBalances {
		balancesSum += clientBalance
	}

	return balancesSum / float64(len(copiedBalances))
}

func (c *Cache) AvgBalanceB() float64 {
	c.mu.RLock()
	// In this case we make the deep copy.
	copiedBalances := make(map[string]float64, len(c.balances))
	for k, v := range c.balances {
		copiedBalances[k] = v
	}
	defer c.mu.RUnlock()

	var (
		balancesSum float64
	)

	for _, clientBalance := range copiedBalances {
		balancesSum += clientBalance
	}

	return balancesSum / float64(len(copiedBalances))
}

func (c *Cache) AvgBalanceC() float64 {
	// In this case we don't unlock the RWMutex until the calculation is done
	c.mu.RLock()
	defer c.mu.RUnlock()

	copiedBalances := make(map[string]float64, len(c.balances))
	for k, v := range c.balances {
		copiedBalances[k] = v
	}

	var (
		balancesSum float64
	)

	for _, clientBalance := range copiedBalances {
		balancesSum += clientBalance
	}

	return balancesSum / float64(len(copiedBalances))
}

func ConcurrentMutexedStructureA() {
	var (
		wg = &sync.WaitGroup{}
		c  = NewCache()
	)
	c.Fill()

	wg.Add(1)
	go func() {
		wg.Done()

		for i := 0; i < 10_000; i++ {
			for k := range c.balances {
				c.AddBalance(k, multiplierBase*rand.Float64())
			}
		}

	}()

	wg.Add(1)
	go func() {
		wg.Done()
		for i := 0; i < 10_000; i++ {
			c.AvgBalanceC()
		}
	}()

	wg.Wait()
}

func WaitGroupNegativeInternalCounter() {
	var (
		wg = &sync.WaitGroup{}
	)

	for i := 0; i < 10_000; i++ {
		wg.Done()
	}

	/*
		panic: sync: negative WaitGroup counter

		goroutine 1 [running]:
		sync.(*WaitGroup).Add(0x14000010180, 0xffffffffffffffff)
			/opt/homebrew/Cellar/go/1.22.2/libexec/src/sync/waitgroup.go:62 +0x1bc
		sync.(*WaitGroup).Done(0x14000010180)
			/opt/homebrew/Cellar/go/1.22.2/libexec/src/sync/waitgroup.go:87 +0x24
		main.WaitGroupNegativeInternalCounter()
			/Users/dmitriymamykin/Desktop/mygo/hundred_issues/internal/concurrencypractice/main.go:1062 +0x5c
		main.main()
			/Users/dmitriymamykin/Desktop/mygo/hundred_issues/internal/concurrencypractice/main.go:139 +0x1c
	*/
}

func IncorrectWaitGroupUsageA() {
	var (
		wg = &sync.WaitGroup{}
		v  uint64
	)

	for i := 0; i < 10_000_000; i++ {
		/*
			The right placement of wg.Add(1)
			In this case we can guarantee that the internal counter of wg will be incremented until wg.Wait()
		*/
		// wg.Add(1)
		go func() {
			defer wg.Done()
			/*
				The incorrect placement of wg.Add(1)
				In this case we cannot guarantee that this section
				will be executed until wg.Wait()
			*/
			wg.Add(1)
			atomic.AddUint64(&v, 1)
		}()
	}

	wg.Wait()

	// When the placement of wg.Add(1) is inside the goroutine, we'll get a random value of v
	fmt.Println(v)
}

func NonDeterministicExecution() {
	var (
		wg = &sync.WaitGroup{}
	)
	wg.Add(3)

	go func() {
		wg.Done()
		fmt.Print("a")
	}()

	go func() {
		wg.Done()

		fmt.Print("b")
	}()

	go func() {
		wg.Done()

		fmt.Print("c")
	}()

	wg.Wait()
}

/*
The main lack of this code is busy-loop. Every goroutine-listener resumes the cycle until the goal isn't reached.
It consumes many cycles of CPU and loads it.
*/
func SyncCondSimulationA() {
	type (
		Donation struct {
			mu      *sync.RWMutex
			balance int
		}
	)

	var (
		wg       = &sync.WaitGroup{}
		donation = &Donation{mu: &sync.RWMutex{}, balance: 0}
	)

	var (
		f = func(wg *sync.WaitGroup, goal int) {
			defer wg.Done()

			donation.mu.RLock()

			/*
				There's busy waiting.
				It spends many CPU cycles and loads it.
			*/
			for donation.balance < goal {
				donation.mu.RUnlock()
				donation.mu.RLock()
			}

			fmt.Printf("$%d goal reached\n", donation.balance)
			donation.mu.RUnlock()
		}
	)

	wg.Add(2)
	go f(wg, 10)
	go f(wg, 15)

	go func() {
		for {
			time.Sleep(time.Millisecond * 100)
			donation.mu.Lock()
			donation.balance++
			donation.mu.Unlock()
		}
	}()

	wg.Wait()

}

/*
The main lack of this is that receiving on a channel is performed by a single goroutine,
but we want multiple goroutines to get the message
*/
func SyncCondSimulationB() {
	type (
		Donation struct {
			ch      chan int
			balance int
		}
	)

	var (
		wg       = &sync.WaitGroup{}
		donation = &Donation{ch: make(chan int), balance: 0}
	)

	var (
		f = func(wg *sync.WaitGroup, goal int) {
			defer wg.Done()

			for v := range donation.ch {
				if v >= goal {
					fmt.Printf("$%d goal reached!\n", goal)
					return
				}

			}

		}
	)

	for i := 10; i <= 500_000; i += 10 {
		wg.Add(1)
		go f(wg, i)
	}

	go func() {
		for {
			time.Sleep(time.Nanosecond * 1)
			donation.balance++
			donation.ch <- donation.balance

		}
	}()

	wg.Wait()

}

func SyncCondSimulationC() {
	type (
		Donation struct {
			cond    *sync.Cond
			balance int
		}
	)

	var (
		wg       = &sync.WaitGroup{}
		donation = &Donation{cond: &sync.Cond{L: &sync.RWMutex{}}, balance: 0}
	)

	var (
		f = func(wg *sync.WaitGroup, goal int) {
			defer wg.Done()

			donation.cond.L.Lock() // <- this part is required by documentation

			// When we reach this section the Wait() locks Locker
			for donation.balance < goal {
				// Firstly Wait() unlocks the Locker
				donation.cond.Wait() // Goroutine is suspended (arranged to a local/global queue of Gs)
				// Secondly Wait() locks the locker
			}

			// When we reach this section Wait() has locked the Locker
			fmt.Printf("$%d goal reached!\n", goal)
			donation.cond.L.Unlock() // <- this part is required
		}
	)

	wg.Add(2)
	go f(wg, 10)
	go f(wg, 15)

	go func() {
		for {
			time.Sleep(time.Millisecond * 100)

			donation.cond.L.Lock() // <- this part is required by documentation

			donation.balance++
			donation.cond.L.Unlock() // <- this part is required

			donation.cond.Broadcast()
		}
	}()

	wg.Wait()

}

type (
	Point struct {
		x, y int
	}
	Circle struct {
		Point
		radius int
	}

	Result string
)

func handler(ctx context.Context, circles []Circle) ([]Result, error) {
	const (
		goroutinesDefaultLimit = 1 << 4
	)
	var (
		results = make([]Result, len(circles))
	)
	// errgroup.WithContext(ctx context.Context) returns a group an associated with given context
	// the derived context will be canceled if:
	// 		1) A goroutine returns a non-nil error
	// 		2) Wait() returns
	// Whichever happens first
	g, ctx := errgroup.WithContext(ctx)
	// g.SetLimit(n int) sets limit on active goroutines at each moment working on an overall task together
	g.SetLimit(goroutinesDefaultLimit)

	for i, circle := range circles {
		i := i
		circle := circle

		// Starts a func() error given in a new goroutine.
		// It suspends the execution of this function if the amount of active goroutines at this moment
		// is equal to limit set
		g.Go(func() error {
			res, err := foo(ctx, circle)
			if err != nil {
				return err
			}

			results[i] = res
			return nil
		})
	}

	// g.Wait() returns the first non-nil error (if any) or nil otherwise
	// Wait() is waiting for all returns from goroutines created happen
	if err := g.Wait(); err != nil {
		return nil, err
	}

	return results, nil
}

func foo(_ context.Context, _ Circle) (Result, error) {
	var (
		res = genID()
		err error
	)

	if rand.Intn(2) == 1 {
		err = errors.New("an error")
	}

	return Result(res), err
}

func ErrGroupContextUsage() {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*1)
	defer cancel()

	g, groupCtx := errgroup.WithContext(ctx)
	for i := 1; i <= 3; i++ {
		i := i
		g.Go(func() error {

			time.Sleep(time.Second * 1)
			if err := groupCtx.Err(); err != nil {
				fmt.Printf("%d stopped\n", i)
				return nil
			}

			return fmt.Errorf("err: %d", i)

		})
	}

	if err := g.Wait(); err != nil {
		fmt.Println(err)
	}
}
