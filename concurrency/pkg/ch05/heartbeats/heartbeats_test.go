package heartbeats

import (
	"testing"
	"time"
)

func TestDoWork(t *testing.T) {
	const (
		timeout = 2 * time.Second
	)

	var (
		intSlice = []int{1, 2, 3, 5, 7}

		done               = make(chan struct{})
		heartbeat, results = DoWork(done, timeout, intSlice...)

		i = 0
	)
	defer func() {
		close(done)
	}()

	// We're waiting for the finish of prework of DoWork
	<-heartbeat

	for {
		select {
		case r, ok := <-results:
			if !ok {
				return
			}

			if expected := intSlice[i]; r != expected {
				t.Errorf("index %v: expected %v, but received %v", i, expected, r)
			}

			i++

		case <-heartbeat:

		case <-time.After(0 * time.Second):
			t.Fatal("test time out")
		}
	}

}
