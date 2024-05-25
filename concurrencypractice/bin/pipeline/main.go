package main

import (
	"fmt"
	"sync"
)

func main() {
	fmt.Println()
	pipelineUsingA()

	fmt.Println()
	fanOutInUsing()

	fmt.Println()
	cancellationUsageA()

}

func gen(done <-chan struct{}, nums ...int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)

		for _, n := range nums {
			select {
			case <-done:
				return

			case out <- n:
			}

		}
	}()

	return out
}

func squareNums(done <-chan struct{}, producer <-chan int) <-chan int {
	squared := make(chan int)

	go func() {
		defer close(squared)

		for n := range producer {
			select {
			case <-done:
				return

			case squared <- n * n:
			}

		}
	}()

	return squared
}

func pipelineUsingA() {
	for square := range squareNums(nil, gen(nil, []int{1, 2, 3, 4, 5, 6, 7, 8, 9}...)) {
		fmt.Print(square, " ")
	}
	fmt.Println()
}

func fanOutInUsing() {
	producer := gen(nil, []int{1, 2, 3, 4, 5, 6, 7, 8, 9}...)

	// We spread out the output to two wokers
	squarerA := squareNums(nil, producer)
	squarerB := squareNums(nil, producer)

	// We dump two outputs into one channel with merge function
	for square := range merge(nil, squarerA, squarerB) {
		fmt.Print(square, " ")
	}
	fmt.Println()

}

func merge(done <-chan struct{}, chs ...<-chan int) <-chan int {
	out := make(chan int)

	var (
		wg sync.WaitGroup

		output = func(c <-chan int) {
			for n := range c {
				select {
				case <-done:
					return

				case out <- n:
				}

			}
		}
	)

	for _, c := range chs {
		wg.Add(1)
		go func(c <-chan int) {
			defer wg.Done()

			output(c)
		}(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func cancellationUsageA() {
	done := make(chan struct{})
	defer close(done)

	producer := gen(done, []int{1, 2, 3, 4, 5, 6, 7, 8, 9}...)

	squarerA := squareNums(done, producer)
	squarerB := squareNums(done, producer)

	fmt.Println(<-merge(done, squarerA, squarerB))
}
