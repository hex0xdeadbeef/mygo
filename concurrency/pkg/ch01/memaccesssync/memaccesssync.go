package memaccesssync

import (
	"fmt"
	"sync"
)

var m sync.Mutex

// The possible outcomes
func DataRace() {
	var data int

	go func() {
		m.Lock()
		data++
		m.Unlock()
	}()

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
