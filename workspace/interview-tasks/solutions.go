package tasks

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
)

/*
TASK 0 (Base for all others)
*/
func getGenerator(cnt int) <-chan int {
	const (
		genPad = 1
		maxNum = 1_000
	)

	ch := make(chan int)

	go func() {
		defer close(ch)

		for range cnt {
			ch <- genPad + rand.Intn(maxNum)
		}
	}()

	return ch
}

/*
TASK 1
*/
func Bridge(ctx context.Context, from <-chan int, to chan<- int) {
	go func() {
		defer close(to)

		for {
			select {
			case <-ctx.Done():
				return

			case v, ok := <-from:
				if !ok {
					return
				}

				select {
				case <-ctx.Done():
					return

				case to <- v:
				}
			}
		}
	}()
}

type customCh struct {
	ID int
	ch chan<- int
}

/*
TASK 2
*/
func FanOut(ctx context.Context, src <-chan int, dests ...customCh) {
	var wg sync.WaitGroup

	for val := range src {
		for _, dest := range dests {
			wg.Add(1)
			go func() {
				defer wg.Done()

				select {
				case <-ctx.Done():
					return

				case dest.ch <- val:
					fmt.Printf("Send val %d to %d\n", val, dest.ID)
				}
			}()
		}

		wg.Wait()
	}

	for _, out := range dests {
		close(out.ch)
	}
}

func FanOutLinear(ctx context.Context, src <-chan int, dests ...customCh) {
	// Decrease the likelyhood of being hung up on a busy destination
	m := map[customCh]struct{}{}
	for _, dest := range dests {
		m[dest] = struct{}{}
	}

	for val := range src {
		for out := range m {

			select {
			case <-ctx.Done():
				return

			case out.ch <- val:
				fmt.Printf("Send val %d to %d\n", val, out.ID)
			}
		}

	}

	for _, out := range dests {
		close(out.ch)
	}
}

/*
TASK 3
*/
func squarer(ctx context.Context, in <-chan int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)

		for {
			select {
			case <-ctx.Done():
				fmt.Println("Context done")
				return

			case v, ok := <-in:
				if !ok {
					fmt.Println("Receiver chan for multiplier closed")
					return
				}

				out <- v * v
			}
		}

	}()

	return out
}

func divider(ctx context.Context, coef int, in <-chan int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)

		for {
			select {
			case <-ctx.Done():
				fmt.Println("Context done")

			case v, ok := <-in:
				if !ok {
					fmt.Println("Receiver chan for divider closed")
					return
				}

				out <- v / coef
			}

		}

	}()

	return out
}

func filterer(ctx context.Context, coef int, in <-chan int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)

		for {
			select {
			case <-ctx.Done():
				fmt.Println("Context done")
				return

			case v, ok := <-in:
				if !ok {
					fmt.Println("Receiver chan for filtering closed")
					return
				}

				if v%coef == 0 {
					continue
				}

				out <- v
			}
		}

	}()

	return out
}
