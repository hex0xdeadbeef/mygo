package main

import (
	"bytes"
	"sync"
)

type Pool struct {
	pool *sync.Pool

	maxBufCap int
}

func NewPool(maxBufCap int) *Pool {
	const (
		initialBufSize = 1 << 12
		indexPad       = 1
	)

	return &Pool{
		pool: &sync.Pool{
			New: func() any {
				s := []byte{initialBufSize - indexPad: 0}
				return bytes.NewBuffer(s)
			},
		},
	}
}

func (p *Pool) Get() *bytes.Buffer {
	b := p.pool.Get().(*bytes.Buffer)
	return b
}

func (p *Pool) Put(b *bytes.Buffer) {
	if b.Cap() >= p.maxBufCap {
		return
	}

	b.Reset()
	p.Put(b)
}
