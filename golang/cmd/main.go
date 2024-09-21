package main

import (
	"fmt"
	"sync"
)

func main() {
	// workWithProdA()
	// workWithProdB()
}

func Func() {
	var (
		wg sync.WaitGroup
	)

	for i := 0; i < 5; i++ {

		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			fmt.Println(i)
		}(i)
	}

	wg.Wait()
}
