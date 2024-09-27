package chapter9

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

var (
	wgMemo   sync.WaitGroup
	muMemo   sync.Mutex
	RWMuMemo sync.RWMutex
)

func incomingURLs() []string {
	return os.Args[1:]
}

func httpGetBody(url string) (interface{}, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	return io.ReadAll(response.Body)
}

// A Memo caches the results of calling a Func.
type Memo struct {
	f     Func
	cache map[string]result
}

// Func is the type of the function to memoize
type Func func(key string) (interface{}, error)

type result struct {
	value interface{}
	err   error
}

func New(f Func) *Memo {
	return &Memo{f: f, cache: make(map[string]result)}
}

/*
1) Writing to the cache by two or more goroutines at the same time for the same key
*/
func (memo *Memo) GetUnsafe(key string) (interface{}, error) {
	res, ok := memo.cache[key]
	if !ok {
		res.value, res.err = memo.f(key)
		memo.cache[key] = res
	}
	return res.value, res.err
}

func MemoUnsafeUsing() {
	m := New(httpGetBody)

	for _, url := range incomingURLs() {
		wgMemo.Add(1)
		go func(url string) {
			defer wgMemo.Done()
			start := time.Now()

			value, err := m.GetUnsafe(url)
			if err != nil {
				log.Print(err)
			}
			fmt.Printf("%s %s, %d bytes\n", url, time.Since(start), len(value.([]byte)))

		}(url)

	}

	wgMemo.Wait()

}

/*
Waiting for all the operations of other goroutines to get into the critical section
*/
func (memo *Memo) GetMutex(key string) (interface{}, error) {
	muMemo.Lock()
	defer muMemo.Unlock()
	res, ok := memo.cache[key]
	if !ok {
		res.value, res.err = memo.f(key)
		memo.cache[key] = res
	}
	return res.value, res.err
}

func MemoMutexUsing() {
	m := New(httpGetBody)

	for _, url := range incomingURLs() {
		wgMemo.Add(1)
		go func(url string) {
			defer wgMemo.Done()
			start := time.Now()

			value, err := m.GetMutex(url)
			if err != nil {
				log.Print(err)
			}
			fmt.Printf("%s %s, %d bytes\n", url, time.Since(start), len(value.([]byte)))

		}(url)

	}

	wgMemo.Wait()

}

/*
1) Waiting for the all read operations
2) A possibility to call the "memo.f(key)" at the same time by two or more goroutines, the result of the last will overwrite the previous
*/
func (memo *Memo) GetTwoMutexes(key string) (interface{}, error) {
	// 1st Critical section
	RWMuMemo.RLock()
	res, ok := memo.cache[key]
	RWMuMemo.RUnlock()

	// Here several goroutines may race to compute f(key) and update the map.

	if !ok {
		res.value, res.err = memo.f(key)
		muMemo.Lock()
		memo.cache[key] = res
		muMemo.Unlock()
	}

	return res.value, res.err
}

func MemoTwoMutexesUsing() {
	m := New(httpGetBody)

	for _, url := range incomingURLs() {
		wgMemo.Add(1)
		go func(url string) {
			defer wgMemo.Done()
			start := time.Now()

			value, err := m.GetTwoMutexes(url)
			if err != nil {
				log.Print(err)
			}
			fmt.Printf("%s %s, %d bytes\n", url, time.Since(start), len(value.([]byte)))

		}(url)

	}

	wgMemo.Wait()

}

type FuncAndCancel struct {
	f    Func
	done chan struct{}
}

func (fc *FuncAndCancel) isCancelled() bool {
	select {
	case <-fc.done:
		return true
	default:
		return false
	}
}

type entry struct {
	res   result
	ready chan struct{}
}

type MemoEntry struct {
	fc    FuncAndCancel
	mu    sync.Mutex
	cache map[string]*entry
}

func NewMemoEntry(f Func, c ...chan struct{}) *MemoEntry {
	fc := FuncAndCancel{f: f}
	if len(c) == 1 && c != nil {
		fc.done = c[0]

		go func() {
			os.Stdin.Read(make([]byte, 1))
			close(fc.done)
		}()

	}
	return &MemoEntry{fc: fc, cache: make(map[string]*entry)}
}

func (memoEntry *MemoEntry) GetDuplicateSuppressing(key string) (interface{}, error) {
	memoEntry.mu.Lock()

	e := memoEntry.cache[key]
	if e == nil {
		/*
			1) This is the first request for this key.
			2) This goroutine becomes responsible for computing the value and broadcasting the ready condition.
		*/
		e = &entry{ready: make(chan struct{})}
		memoEntry.cache[key] = e
		memoEntry.mu.Unlock()
		e.res.value, e.res.err = memoEntry.fc.f(key)

		close(e.ready) // After closing starts to pass the zero value of struct{}
	} else {
		// This is repeat request for this key.
		memoEntry.mu.Unlock()

		<-e.ready
	}

	if memoEntry.fc.isCancelled() {
		memoEntry.mu.Lock()
		delete(memoEntry.cache, key)
		memoEntry.mu.Unlock()
		return []byte{}, errors.New("operation canceled")
	}

	return e.res.value, e.res.err
}

func MemoEntryUsing() {
	m := NewMemoEntry(httpGetBody, make(chan struct{}))

	for _, url := range incomingURLs() {
		wgMemo.Add(1)
		go func(url string) {
			defer wgMemo.Done()
			start := time.Now()

			value, err := m.GetDuplicateSuppressing(url)
			if err != nil {
				log.Print(err)
			}
			fmt.Printf("%s %s, %d bytes\n", url, time.Since(start), len(value.([]byte)))

		}(url)

	}

	wgMemo.Wait()
}

type MemoChan struct {
	requests chan request
}

type request struct {
	key      string
	response chan<- result
}

func NewMemoChan(f Func) *MemoChan {
	memo := &MemoChan{requests: make(chan request)}
	go memo.server(f)
	return memo
}

func (memoChan *MemoChan) server(f Func) {
	cache := make(map[string]*entry)

	for request := range memoChan.requests {
		e := cache[request.key]

		if e == nil {
			// This is the first request for this key
			e = &entry{ready: make(chan struct{})}
			cache[request.key] = e
			go e.call(f, request.key)
		}
		go e.deliver(request.response)
	}
}

func (e *entry) call(f Func, key string) {
	// Evaluate the function
	e.res.value, e.res.err = f(key)
	// Broadcast the ready condition
	close(e.ready)
}

func (e *entry) deliver(response chan<- result) {
	// Wait for ready condition
	<-e.ready
	// Send the result to the client
	response <- e.res
}
