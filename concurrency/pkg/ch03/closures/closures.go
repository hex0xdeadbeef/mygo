package closures

import (
	"fmt"
	"sync"
)

var wg sync.WaitGroup

func Single() {
	salutation := "Hello!"

	wg.Add(1)
	go func() {
		defer wg.Done()
		salutation = "welcome"
	}()

	wg.Wait()

	// There are two possible cow outcomes:
	// 1. The additional goroutine will start before the fmt.Println(salutation) operation. The output will be: "welcome"
	// 2. The additional goroutine will start after the fmt.Println(salutation) operation. The output will be: "Hello"
	fmt.Println(salutation)
}

func Loop() {
	for _, salutation := range []string{"hello", "greetings", "good day"} {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// There are many possible outcomes:
			// 1. The goroutines start befor for loop exits, so the output will be the slice contents in the random order. In this case every goroutine closure the local variable, so
			// the data doesn't escapes the main goroutine stack.
			// 2. The goroutines will run after the for loop exits. In this case for variable escapes the local domain (the main goroutine stack) and resides in the heap. In this
			// case the last change of the for variable will be printed by each goroutine.
			// 3. The part of all the goroutines run in the for range and the others goroutines will run after the for loop ends. In this case the data the data will be separated.
			fmt.Println(salutation)
		}()
	}

	wg.Wait()
}

func LoopFix() {
	for _, salutation := range []string{"hello", "greetings", "good day"} {
		wg.Add(1)
		// In this case we create the copy of every value of the for loop variable and points to in in the closure.
		go func(str string) {
			defer wg.Done()
			fmt.Println(str)
		}(salutation)
	}

	wg.Wait()
}
