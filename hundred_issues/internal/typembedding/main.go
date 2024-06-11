package main

import (
	"fmt"
	"io"
	"sync"
)

func main() {
	UsageA()
}

type Foo struct {
	Bar // The anonymous field
}

type Bar struct {
	Baz int
}

func UsageA() {
	foo := Foo{Bar: Bar{Baz: 42}}

	fmt.Printf("%#v", foo)

	fmt.Println(foo.Baz)
	fmt.Println(foo.Bar.Baz)
}

// Mutex should be the unexported field so that no client can use its methods
type InMem struct {
	mu *sync.Mutex
	m  map[string]int
}

func NewInMem() *InMem {
	return &InMem{mu: &sync.Mutex{}, m: make(map[string]int)}
}

func (im *InMem) Get(k string) (int, bool) {
	im.mu.Lock()
	defer im.mu.Unlock()

	v, ok := im.m[k]
	return v, ok
}

// With embedding we can remove any routing from the code

// Code with re-routing
// type Logger struct {
// 	writeCloser io.WriteCloser
// }

// func (l Logger) Write(b []byte) (int, error) {
// 	return l.writeCloser.Write(b)
// }

// func (l Logger) Close() error {
// 	return l.writeCloser.Close()
// }

// Code without any routing
// Additionally the new Logger type implements io.WriterClose interface
type Logger struct {
	io.WriteCloser
}

func UsageB() error {
	l := Logger{}

	n, err := l.Write([]byte("foo"))
	if err != nil {
		return err
	}
	fmt.Println(n)

	if err = l.Close(); err != nil {
		return err
	}

	return nil
}
