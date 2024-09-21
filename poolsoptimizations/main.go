package main

import "sync"

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
