package pipeline

import (
	"bytes"
	"fmt"
	"math/rand"
	"time"
)

/*
A generator for a pipeline is any function that converts a set of discrete values into a stream of values on a channel.
*/

// Repeat will send to a channel repeatable discrete elements until we close the done channel.
func Repeat(done <-chan struct{}, values ...any) <-chan any {
	valueStream := make(chan any)

	go func() {
		defer close(valueStream)

		for {
			for _, val := range values {
				select {
				case <-done:
					return
				case valueStream <- val:
				}
			}
		}

	}()

	return valueStream
}

func LimitedJunction(done <-chan struct{}, valueStream <-chan interface{}, num int) <-chan interface{} {
	takeStream := make(chan interface{})

	go func() {
		defer close(takeStream)

		for i := 0; i < num; i++ {
			select {
			case <-done:
				return
			// The valueStream gives the value and takeStream accepts it.
			case takeStream <- <-valueStream:
			}
		}
	}()

	return takeStream
}

func LimitedJunctionUsing() {
	done := make(chan struct{})
	time.AfterFunc(time.Microsecond*10, func() {
		close(done)
	})

	for v := range LimitedJunction(done, Repeat(done, struct {
		name    string
		surname string
		age     uint
	}{"Denis", "Nemchenko", 0}),
		10) {
		fmt.Printf("%v ", v)
	}

	fmt.Println()
}

func RepeatFunction(done <-chan struct{}, fn func() any) <-chan any {
	valueStream := make(chan any)
	go func() {
		defer close(valueStream)
		for {
			select {
			case <-done:
				return
			case valueStream <- fn():
			}
		}
	}()

	return valueStream
}

func RepeatFunctionUsing() {
	done := make(chan struct{})
	defer close(done)

	rand := func() any { return rand.Int() }

	for v := range LimitedJunction(done, RepeatFunction(done, rand), 5) {
		fmt.Println(v)
	}
}

func StringTaker(done <-chan struct{}, valueStream <-chan interface{}) <-chan string {
	stringStream := make(chan string)

	go func() {
		defer close(stringStream)

		for v := range valueStream {
			select {
			case <-done:
				return
			case stringStream <- v.(string):
			}
		}
	}()

	return stringStream
}

func StringTakerUsing() {
	done := make(chan struct{})
	defer close(done)

	var buff bytes.Buffer
	for v := range StringTaker(done, LimitedJunction(done, Repeat(done, `I`, `am.`), 10)) {
		buff.WriteString(v)
	}

	fmt.Println(buff.String())
}
