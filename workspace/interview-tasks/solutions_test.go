package tasks

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

const cnt = 5

// TASK 1
/* Bridged val: 785
Bridged val: 673
Bridged val: 433
Bridged val: 523
Bridged val: 720
*/
func TestBridge(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	from := getGenerator(cnt)
	to := make(chan int)

	Bridge(ctx, from, to)

	for v := range to {
		fmt.Println("Bridged val:", v)
	}

}

// TASK 2
/*
Send val 0 to 1
out1 val:  0
Send val 0 to 2
out2 val:  0
Send val 1 to 1
out1 val:  1
Send val 1 to 2
Send val 2 to 1
out1 val:  2
out2 val:  1
out2 val:  2
Send val 2 to 2
Send val 3 to 1
Send val 3 to 2
out1 val:  3
out2 val:  3
Send val 4 to 1
out1 val:  4
*/
func TestFanOut(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	in := make(chan int)
	out1, out2 := make(chan int), make(chan int)

	go FanOut(ctx, in, customCh{ch: out1, ID: 1}, customCh{ch: out2, ID: 2})

	go func() {
		for v := range out1 {
			fmt.Println("out1 val: ", v)
		}
	}()

	go func() {
		for v := range out2 {
			fmt.Println("out2 val: ", v)
		}
	}()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(in)

		for i := range cnt {
			in <- i
		}
	}()

	wg.Wait()
}

// TASK 3
//
// Filtered val: 2640
// Filtered val: 182040
// Filtered val: 120801
// Filtered val: 48
// Receiver chan for multiplier closed
// Receiver chan for divider closed
// Receiver chan for filtering closed
func TestPipeline(t *testing.T) {
	const (
		dividerCoef = 3
		filterCoef  = 7
	)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	gen := getGenerator(cnt)
	squarerCh := squarer(ctx, gen)
	dividerCh := divider(ctx, dividerCoef, squarerCh)
	filterCh := filterer(ctx, filterCoef, dividerCh)

	for v := range filterCh {
		fmt.Println("Filtered val:", v)
	}

}
