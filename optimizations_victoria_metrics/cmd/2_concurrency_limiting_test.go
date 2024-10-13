package main

import "testing"

// Benchmark_startWorkerPool-12    	    3157	    380391 ns/op	     434 B/op	      15 allocs/op
func Benchmark_startWorkerPool(b *testing.B) {
	for n := 0; n < b.N; n++ {
		startWorkerPool()
	}
}

// Benchmark_startSemaphorePool-12    	    3256	    368681 ns/op	   98437 B/op	    2050 allocs/op
func Benchmark_startSemaphorePool(b *testing.B) {
	for n := 0; n < b.N; n++ {
		startSemaphorePool()
	}
}
