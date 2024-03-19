package memaccesssync

import (
	"fmt"
	"sync"
)

var m sync.Mutex

// The possible outcomes
func DataRace() {
	var data int

	// We protect this code from the intervenings of the other goroutines that could possibly read the stale value of "data"
	// None of the goroutines that works simultaneously can read the variables in this critical section until it's done
	go func() {
		m.Lock()
		data++
		m.Unlock()
	}()

	// We protect this code from the intervenings of the other goroutines that could possibly write the new value into the "data"
	// None of the goroutines that works simultaneously can change the variables in this critical section
	m.Lock()
	if data == 0 {
		// 1. The data is 0
		// 2. The data after the check is updated
		fmt.Println("the value is 0")
	} else {
		// 3. The data isn't 0
		fmt.Printf("the value is %d\n", data)
	}
	m.Unlock()

}
