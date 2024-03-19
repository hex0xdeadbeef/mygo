package livelock

import (
	"bytes"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

/*
Livelock are programs that are actively performing concurrent operations, but these operations do nothing to move the state of the program forward

Have you ever been in a hallway walking toward another person? She moves to one side to let you pass, but you've just done the same. So you move to the other side, but she's also
done the same. Imagine this going on forever, and you understand livelocks.
*/

func Livelock() {
	cadence := sync.NewCond(&sync.Mutex{})

	go func() {
		for range time.Tick(1 * time.Millisecond) {
			// Wakes up all the goroutines that wait on this condition
			cadence.Broadcast()
		}
	}()

	takeStep := func() {
		cadence.L.Lock()
		cadence.Wait()
		cadence.L.Unlock()
	}

	tryDirection := func(directionName string, dir *int32, out *bytes.Buffer) bool {
		fmt.Fprintf(out, " %v ", directionName)
		atomic.AddInt32(dir, 1)
		takeStep()
		if atomic.LoadInt32(dir) == 1 {
			fmt.Fprintf(out, ". Success!")
			return true
		}
		takeStep()
		atomic.AddInt32(dir, -1)
		return false
	}

	var left, right int32
	tryLeft := func(out *bytes.Buffer) bool {
		return tryDirection("left", &left, out)
	}
	tryRight := func(out *bytes.Buffer) bool {
		return tryDirection("right", &right, out)
	}

	walk := func(walking *sync.WaitGroup, name string) {
		var out bytes.Buffer

		// Print the result of walking
		defer func() {
			fmt.Println(out.String())
		}()

		// The first func to be invoked at the finishing of execution
		defer walking.Done()

		fmt.Fprintf(&out, "%v is trying to scoot", name)

		for i := 0; i < 5; i++ {
			if tryLeft(&out) || tryRight(&out) {
				return
			}
		}

		fmt.Fprintf(&out, "\n%v tosses her hands up in exasperation!", name)
	}

	var peopleInHallway sync.WaitGroup
	peopleInHallway.Add(2)
	go walk(&peopleInHallway, "Alice")
	go walk(&peopleInHallway, "Barbara")
	peopleInHallway.Wait()
}
