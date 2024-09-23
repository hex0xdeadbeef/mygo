package main

import (
	"fmt"
	"io"
	"sync"
)

/*
	sync.Pool USAGE
1. Before the version 1.13 the sync.Pool is being cleaned at each iteration of the GC. It can negatively impact the performance of our programs that allocate much memory. Starting from 1.13
	more objects are kept alive after cleaning up by the GC.
2. Before backing an object up to the sync.Pool, the fields of the object should be cast to the default value relatively to the type. If we don't do it, we can get "dirty" object from the
	sync.Pool.
3. The case when we needn't null the object is when we use the memory, that is used to write into it.
*/

/*
NULLING OBJECTS BEFORE PUTTING TO THE POOL
*/
type AuthentificationResponse struct {
	Token  string
	UserID string
}

type Pool struct {
	sp  *sync.Pool
	rmu *sync.RWMutex
}

// Get() *AuthentificationResponse returns the *AuthentificationResponse from the p.sp *sync.Pool
func (p *Pool) Get() *AuthentificationResponse {
	return p.sp.Get().(*AuthentificationResponse)
}

// Put(ar *AuthentificationResponse) returns a pointer to a nulled AuthentificationResponse object
func (p *Pool) Put(ar *AuthentificationResponse) {
	ar.reset()
	p.sp.Put(ar)
}

// reset() resets all the fields of a
func (a *AuthentificationResponse) reset() {
	a.Token, a.UserID = "", ""
}

/*
	NOT NULLING AN OBJECT TO
*/

type CustomBuffer struct {
	bp *BufferPool

	w io.Writer
	r io.Reader
}

func (cb *CustomBuffer) Write() (n int, err error) {
	if cb.w == nil {
		return 0, fmt.Errorf("writer not specified")
	}

	if cb.bp == nil {
		return 0, fmt.Errorf("buffer pool not specified")
	}

	buf := cb.bp.GetBytes()
	n, err = cb.w.Write(buf)
	cb.bp.sp.Put(buf)

	return n, err
}

func (cb *CustomBuffer) Read() (n int, err error) {
	if cb.w == nil {
		return 0, fmt.Errorf("reader not specified")
	}

	if cb.bp == nil {
		return 0, fmt.Errorf("buffer pool not specified")
	}

	buf := cb.bp.GetBytes()
	n, err = cb.r.Read(cb.bp.sp.Get().([]byte))
	cb.bp.sp.Put(buf)

	return n, err
}

type BufferPool struct {
	sp         *sync.Pool
	maxBufSize int
}

func NewBufferPool(maxBufSize int) *BufferPool {
	return &BufferPool{
		sp: &sync.Pool{
			New: func() any { return make([]byte, 0, maxBufSize) },
		},
		maxBufSize: maxBufSize,
	}
}

func (bp *BufferPool) GetBytes() []byte {
	return bp.sp.Get().([]byte)
}

func (bp *BufferPool) PutBytes(buf []byte) (error, []byte) {
	if len(buf) > bp.maxBufSize {
		return fmt.Errorf("buffer is overfloved;"), buf
	}

	buf = buf[:0]
	return nil, nil
}
