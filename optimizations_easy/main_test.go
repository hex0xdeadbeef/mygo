package main

import (
	"math/rand"
	"testing"
)

// BenchmarkStrconvFmt100-12          	12075698	        83.19 ns/op	     136 B/op	       3 allocs/op
// BenchmarkStrconvFmt100000-12       	   96361	     12383 ns/op	  213149 B/op	       4 allocs/op
// BenchmarkStrconvFmt100000000-12    	     270	   4342680 ns/op	200017686 B/op	       6 allocs/op
// BenchmarkFmtFmt100-12              	28344726	        42.37 ns/op	     120 B/op	       2 allocs/op
// BenchmarkFmtFmt100000-12           	  223352	      5502 ns/op	  106505 B/op	       2 allocs/op
// BenchmarkFmtFmt100000000-12        	     510	   2459934 ns/op	100007958 B/op	       2 allocs/op

func generateString(size int) string {
	var (
		buf = make([]byte, size)
	)

	for i := 0; i < size; i++ {
		buf[i] = tokenSymbolsSet[rand.Intn(len(tokenSymbolsSet))]
	}

	return string(buf)
}

const (
	number = 1e7
)

func fmtFmtSubBench(size int, b *testing.B) {
	var (
		s = generateString(size)
	)
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		FmtFmt(s, number)
	}
}

func strconvFmtSubBench(size int, b *testing.B) {
	var (
		s = generateString(size)
	)
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		StrconvFmt(s, number)
	}
}

func BenchmarkStrconvFmt100(b *testing.B) {
	fmtFmtSubBench(100, b)
}

func BenchmarkStrconvFmt100000(b *testing.B) {
	fmtFmtSubBench(100_000, b)
}

func BenchmarkStrconvFmt100000000(b *testing.B) {
	fmtFmtSubBench(100_000_000, b)
}

func BenchmarkFmtFmt100(b *testing.B) {
	strconvFmtSubBench(100, b)
}

func BenchmarkFmtFmt100000(b *testing.B) {
	strconvFmtSubBench(100_000, b)
}

func BenchmarkFmtFmt100000000(b *testing.B) {
	strconvFmtSubBench(100_000_000, b)
}
