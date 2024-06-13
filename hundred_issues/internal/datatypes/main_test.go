package main

import (
	"math/rand"
	"testing"
)

/*
BenchmarkConvertSlow10-12           	12414242	        85.05 ns/op	     248 B/op	       5 allocs/op
BenchmarkConvertSlow1_000-12        	  601699	      2024 ns/op	   25208 B/op	      12 allocs/op
BenchmarkConvertSlow1_000_000-12    	     472	   2537682 ns/op	41695079 B/op	      38 allocs/op
*/
func benchmarkConvertSlow(n int, b *testing.B) []int {
	var (
		res = make([]int, 0, n)
	)

	for range n {
		res = append(res, rand.Intn(100))
	}

	for i := 0; i < b.N; i++ {
		ConvertSlow(res)
	}

	return res
}

func BenchmarkConvertSlow10(b *testing.B) {
	benchmarkConvertSlow(10, b)
}

func BenchmarkConvertSlow1_000(b *testing.B) {
	benchmarkConvertSlow(1_000, b)
}

func BenchmarkConvertSlow1_000_000(b *testing.B) {
	benchmarkConvertSlow(1_000_000, b)
}

/*
BenchmarkConverWithReuseLengthAppend0-12            	67679161	        17.35 ns/op	      80 B/op	       1 allocs/op
BenchmarkConverWithReuseLengthAppend1_000-12        	 1481790	       827.5 ns/op	    8192 B/op	       1 allocs/op
BenchmarkConverWithReuseLengthAppend1_000_000-12    	     735	   1622110 ns/op	 8014482 B/op	       1 allocs/op
*/
func benchmarkConverWithReuseLengthAppend(n int, b *testing.B) []int {
	var (
		res = make([]int, 0, n)
	)

	for range n {
		res = append(res, rand.Intn(100))
	}

	for i := 0; i < b.N; i++ {
		ConverWithReuseLengthAppend(res)
	}

	return res
}

func BenchmarkConverWithReuseLengthAppend10(b *testing.B) {
	benchmarkConverWithReuseLengthAppend(10, b)
}

func BenchmarkConverWithReuseLengthAppend1_000(b *testing.B) {
	benchmarkConverWithReuseLengthAppend(1_000, b)
}

func BenchmarkConverWithReuseLengthAppend1_000_000(b *testing.B) {
	benchmarkConverWithReuseLengthAppend(1_000_000, b)
}

/*
BenchmarkConverWithReuseLengthInit10-12           	66557008	        17.73 ns/op	      80 B/op	       1 allocs/op
BenchmarkConverWithReuseLengthAppend1_000-12      	 1000000	      1096 ns/op	    8192 B/op	       1 allocs/op
BenchmarkConverWithReuseLengthInit1_000_000-12    	     734	   1603327 ns/op	 8014495 B/op	       1 allocs/op
*/
func benchmarkConverWithReuseLengthInit(n int, b *testing.B) []int {
	var (
		res = make([]int, 0, n)
	)

	for range n {
		res = append(res, rand.Intn(100))
	}

	for i := 0; i < b.N; i++ {
		ConverWithReuseLengthInit(res)
	}

	return res
}

func BenchmarkConverWithReuseLengthInit10(b *testing.B) {
	benchmarkConverWithReuseLengthInit(10, b)
}

func BenchmarkConverWithReuseLengthInit1_000(b *testing.B) {
	benchmarkConverWithReuseLengthInit(1_000, b)
}

func BenchmarkConverWithReuseLengthInit1_000_000(b *testing.B) {
	benchmarkConverWithReuseLengthInit(1_000_000, b)
}
