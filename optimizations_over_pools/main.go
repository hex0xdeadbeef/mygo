package main

import (
	"fmt"
	"sync"
)

// https://habr.com/ru/companies/yadro/articles/842314/
/*
	POOLS OPTIMIZATIONS
1. The moments when GC is called:
	- Regularly, for example: each 2 minutes.
	- After the dead heap is overflowed
	- When a developer wants
2. There's the problem when we don't controll the memory. It consists in the following subtlety:
	- It might happen that it'll require the memory size that is close to GOMEMLIMIT value during a period between GC calling.

	In the picture 1 we see that the application get memory, but the GC hasn't come yet to free the memory up. The memory will be freed at the next GC call isn't available yet. We cannot use it until the GC call. We have to wait until the GC has started and returns the memory areas into the list of the available memory areas.

	It's a-okay for programs that isn't demanding. But if we need to work with memory, this type of optimization doesn't work. The only one advantage is that a programmer doensn't need to bother about memory, all the managing is duty of the GC.
3. The main idea of memory optimization is to allocate a specific memory volume and then free them up manually. Doing this way we make the GC life much easier.

	THE WAYS WE CAN MAKE OPTIMIZATIONS WITH
1. Buffers channel using.
	We will use a buffered channel to store the distributed fragments of memory. This method is called: "Channel Pool".

	With this realization GetBytes() will wait until a free buffer from the channel, and PutBytes() will return the buffer used back.

	Since we cannot change the size of buffered channel after the initialization, we need to predict the number of buffers that will fit over needs. If we cannot do it, then the situation when the program starves of buffers can happen.

	We cannot usualy predict the number of buffers we need precisely.
2. sync.Pool channel pool usage
	If we make an operator of sync.Pool to return a slice of other sructure of large size, we will use this structure and put it back to the pool. If there are no free strucures the New() function will create and return a new instance of it.

	However, if the memory that has been isn't being used at the moment, the GC will be called somewhere and it'll collect this memory. We must not consider this data to be persistent. We cannot store data in objects that will be put in sync.Pool
3. arena
	We can do work with memory better by using arena's.

	Memory arena - is the area of memory that is not controlled by the GC. When we were working with the heap using sync.Pool's we had the problem: the GC came and dusted the memory areas. With sync.Pool and the base of memory behalf of memory arena we coped with these cleans.

	There's a corner case. If we haven't guessed the size of the block of a buffer, we cannot change it at all. We need to choose the size of a buffer that we'll require from the pool.

	THE OVERFLOW OF BUFFERS
1. If we store the memory in slices and at the time of execution we cross the border of the 	possible size limitation, we should remove this buffer from the pool whie using the PutBytes() method.
2. The non-returned area of the memory will be cleaned by the GC and put back to the areas we can
grab and use.

	HOW TO CHOOSE THE BEST METHOD?
1. Without memory management
	If we don't use memory actively, we needn't care about memory at all. The maximum we should do
	is to discover our needs in terms of memory consumption and try to minimize using of it. We can make the discovery of escaping objects to the heap and prevent our code from it.

	Objects are moved to the heap when:
		- The size of object cross the border of class' size.
		- The object is returned from a function as a pointer
		- The object is passed to the function as an interface variable
2. Channel Pool
	For example we write an editor of video snapshots that is used for the smooth animation. For this program the Channel Pool fits the best because we know the size of the vault in advance.

	The Channel Pool also fits the programs that is performing a fragmentational workload on memory. The gaps beetween peaks of the workload in such programs take place, but they're not often.
3. sync.Pool Pool
	In the cases when we don't know the number of the buffers we will want in the future we'd better to use the Pool based on the sync.Pool because we allocate memory dynamically.
*/

type ChanPool struct {
	c chan *[]byte

	contentCap int
}

func NewChanPool(chanSize int, contentCap int) *ChanPool {
	return &ChanPool{c: make(chan *[]byte, chanSize), contentCap: contentCap}
}

func (cp *ChanPool) GetBytes() *[]byte {
	select {
	case buf, ok := <-cp.c:
		if ok {
			return buf
		}
	// This block allow to make a non-blocking call, so if we don't get a buf, we'll just move to the code lower
	default:
	}

	b := make([]byte, 0, cp.contentCap)
	return &b
}

func (cp *ChanPool) PutBytes(b *[]byte) {
	// The GC will free this memory up later
	if cap(*b) > cp.contentCap {
		return
	}

	*b = (*b)[:0]
	cp.c <- b
}

type SyncPool struct {
	*sync.Pool
	contentCapTopBoundary int

	mu *sync.RWMutex
}

func getBufferConstructor(size int) func() any {
	return func() any {
		buf := make([]byte, 0, size)
		return &buf
	}
}

func NewSyncPool(initialBufferSize int, contentCap int) *SyncPool {
	return &SyncPool{Pool: &sync.Pool{
		New: getBufferConstructor(initialBufferSize)},
		contentCapTopBoundary: contentCap,
		mu:                    &sync.RWMutex{},
	}
}

func (sp *SyncPool) SetNewBuffer(newBufferSize int) {
	newBufferProducer := getBufferConstructor(newBufferSize)

	// LOOK HERE
	sp.mu.Lock()
	sp.New = newBufferProducer
	sp.mu.Unlock()
}

func (sp *SyncPool) GetBytes() *[]byte {
	// LOOK HERE
	sp.mu.RLock()
	buf := sp.Get().(*[]byte)
	sp.mu.RUnlock()

	return buf
}

func (sp *SyncPool) PutBytes(buf *[]byte) {
	// The GC will utilize this unused memory
	if cap(*buf) > sp.contentCapTopBoundary {
		return
	}

	*buf = (*buf)[:0]

	// LOOK HERE
	sp.mu.RLock()
	sp.Put(buf)
	sp.mu.RUnlock()
}

func main() {
	workWithProdA()
}

func workWithProdA() {
	const (
		topBound = 50
	)

	var (
		prod = getProducerA()
	)

	for v := range prod {
		fmt.Println(v)

		if v == topBound {
			prod = make(<-chan int)
			continue
		}
	}
}

func getProducerA() <-chan int {
	var (
		prod = make(chan int)
	)

	go func() {
		defer close(prod)

		for i := 1; i <= 100; i++ {

			prod <- i
		}

	}()

	return prod
}

func workWithProdB() {
	defer fmt.Println("return from func")

	const (
		topBound = 50
	)

	var (
		prod, done = getProducerB()
	)
	<-done

	for {
		select {
		case v, ok := <-prod:
			if !ok {
				return
			}

			fmt.Println(v)
			if v == topBound {
				prod = make(<-chan int)
				fmt.Println("New producer provided")
			}

		default:
			return
		}
	}

}

func getProducerB() (<-chan int, <-chan struct{}) {
	var (
		prod = make(chan int, 100)
		done = make(chan struct{})
	)

	go func() {
		defer close(prod)

		for i := 1; i <= 100; i++ {
			prod <- i
		}

		close(done)
	}()

	return prod, done
}
