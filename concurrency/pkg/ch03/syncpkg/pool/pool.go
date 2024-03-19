package pool

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

// Pool is concurrent-safe implementation of the object pool pattern.

/*
	A pool pattern is a way to create and make available a fixed number, or pool, of things for use. It's commonly used to constrain the creation of things that are
	expensive (e.g. database connections) so that only a fixed number of them are ever created, but an indeterminate number of operations can still request access to these
	things.

	Pool's primary interface is its "Get()" method. When called, "Get()" will first check whether there are any available instances within the pool to return to the caller
	and if not, call its "New()" member variable to create a new one. When finished, callers call Put to place the instance they were working with back in the pool for use
	by other processes.
*/

// In the case of Go's "sync.Pool", this data type can be safely used by multiple goroutines.

// The sync.Pool is useful in the following situations:
// 1. We want to store the certain quantity of free/buzy objects to be used in the execution and not to be garbage collected.
// 2. Warming a cache of preallocated objects for operations that must run as quickly as possible. In this case, instead of trying to guard the host machine's memory by\
// constraining the number of objects created, we're trying to guard consumers' time by front-loading the time it takes to get a reference to another object. This is very
// common when writing high-throughput network servers that attempt to respond as quickly as possible.

/*
If the code that utilizes sync.Pool requires things that are roughly homogenous, we may spend more time converting what you've retrivied from the pool that it would have
taken to just instatiate in the first place. For instance: if our program requires slices of random and variable length, a Pool isn't going to help us much. The probability
that you'll receive a slice the length tou require is low.
*/

/*
The rules of working with "sync.Pool"
1. When instantiating sync.Pool, give it a New() member variable that is thread-safe when called.
2. When you retrieve an instance from Get(), make no assumptions regarding the state of the object you receive back.
3. Make sure to call Put when you're finished with the object you pulled out of the pool. Otherwise, the Pool is useless. Usually this is done with defer.
4. Objects in the pool must be roughly uniform in makeup.
*/

func Example() {
	myPool := &sync.Pool{
		New: func() any {
			fmt.Println("Creating new instance.")
			return struct{}{}
		},
	}

	myPool.Get() // This call will invoke the "New()" function defined on the pool since instances haven't yet been instantiated.
	instance := myPool.Get()
	myPool.Put(instance) // We put an instance previously retrieved back in the pool. This increases the available number of instances to one.
	myPool.Get()         // When this call is executed we will reuse the instance previously allocated and put it back in the pool. The "New()" function won't be invoked.
}

func Using() {
	const (
		size          = 1024
		workersNumber = 1024 * 1024
	)

	var (
		calcsCreatedCounter int

		calcPool = &sync.Pool{
			New: func() any {
				calcsCreatedCounter++
				mem := make([]byte, size)
				// NOTE: We return the pointer to a slice.
				return &mem
			},
		}

		wg sync.WaitGroup
	)

	// Create and put the slices of size 1KB into the pool.
	calcPool.Put(calcPool.New())
	calcPool.Put(calcPool.New())
	calcPool.Put(calcPool.New())
	calcPool.Put(calcPool.New())

	for i := workersNumber; i > 0; i-- {
		wg.Add(1)
		go func() {
			defer wg.Done()

			// If there's a free calculator in the pool, we'll get it. Otherwise the New() method will be called.
			mem := calcPool.Get().(*[]byte)
			defer calcPool.Put(mem)

		}()
	}

	wg.Wait()
	// In average we create 30 calculators.
	fmt.Printf("%d calculators were created\n", calcsCreatedCounter)

}

// connectToService simulates the creation of a connection to a service.
func connectToService() interface{} {
	time.Sleep(1 * time.Second)
	return struct{}{}
}

func warmServiceConnCache() *sync.Pool {
	p := &sync.Pool{
		New: connectToService,
	}
	for i := 0; i < 10; i++ {
		p.Put(p.New())
	}
	return p
}

func startNetworkDaemon() *sync.WaitGroup {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		connPool := warmServiceConnCache()
		server, err := net.Listen("tcp", "localhost:8080")
		if err != nil {
			log.Fatalf("cannot listen: %v", err)
		}
		defer server.Close()
		wg.Done()
		for {
			conn, err := server.Accept()
			if err != nil {
				log.Printf("cannot accept connection: %v", err)
				continue
			}
			svcConn := connPool.Get()
			fmt.Fprintln(conn, "")
			connPool.Put(svcConn)
			conn.Close()
		}
	}()
	return &wg
}
