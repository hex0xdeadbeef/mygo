package main

import "testing"

func concatSubBench(b *testing.B, f func([]string) string, amount int, size int) {
	var (
		testSet = getTokenSet(amount, size)
		ret     string
	)
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		ret = f(testSet)
	}

	_ = ret
}

/*
BenchmarkConcatLinearlyString100x256-12       	    6844	    172342 ns/op	 1352456 B/op	      99 allocs/op
BenchmarkConcatLinearlyString1_000x256-12     	     100	  10769807 ns/op	131712187 B/op	    1006 allocs/op
BenchmarkConcatLinearlyString10_000x256-12    	       1	1097544208 ns/op	12840589632 B/op	   10061 allocs/op
BenchmarkConcatLinearlyBytes100x256-12        	   54603	     22480 ns/op	  140160 B/op	      14 allocs/op
BenchmarkConcatLinearlyBytes1_000x256-12      	    7574	    184129 ns/op	 1448198 B/op	      22 allocs/op
BenchmarkConcatLinearlyBytes10_000x256-12     	     751	   2037670 ns/op	15907082 B/op	      32 allocs/op
BenchmarkConcatBuffer100x256-12               	  306170	      3447 ns/op	   27264 B/op	       1 allocs/op
BenchmarkConcatBuffer1_000x256-12             	   39016	     25749 ns/op	  262145 B/op	       1 allocs/op
BenchmarkConcatBuffer10_000x256-12            	    5223	    228199 ns/op	 2564100 B/op	       1 allocs/op
*/
func BenchmarkConcatLinearlyString100x256(b *testing.B) {
	concatSubBench(b, ConcatLinearlyString, 100, 256)
}

func BenchmarkConcatLinearlyString1_000x256(b *testing.B) {
	concatSubBench(b, ConcatLinearlyString, 1_000, 256)
}

func BenchmarkConcatLinearlyString10_000x256(b *testing.B) {
	concatSubBench(b, ConcatLinearlyString, 10_000, 256)
}

func BenchmarkConcatLinearlyBytes100x256(b *testing.B) {
	concatSubBench(b, ConcatLinearlyBytes, 100, 256)
}

func BenchmarkConcatLinearlyBytes1_000x256(b *testing.B) {
	concatSubBench(b, ConcatLinearlyBytes, 1_000, 256)
}

func BenchmarkConcatLinearlyBytes10_000x256(b *testing.B) {
	concatSubBench(b, ConcatLinearlyBytes, 10_000, 256)
}

func BenchmarkConcatBuffer100x256(b *testing.B) {
	concatSubBench(b, ConcatBuffer, 100, 256)
}

func BenchmarkConcatBuffer1_000x256(b *testing.B) {
	concatSubBench(b, ConcatBuffer, 1_000, 256)
}

func BenchmarkConcatBuffer10_000x256(b *testing.B) {
	concatSubBench(b, ConcatBuffer, 10_000, 256)
}

// BenchmarkSmallObjectBuffer6x7-12    	 8323852	       144.8 ns/op	    1536 B/op	       1 allocs/op
func BenchmarkSmallObjectBuffer6x7(b *testing.B) {
	concatSubBench(b, ConcatBuffer, 6, 7)
}

// BenchmarkSmallObjectGetFullname-12    	41790975	        28.28 ns/op	      64 B/op	       1 allocs/op
func BenchmarkSmallObjectGetFullname(b *testing.B) {
	var (
		placeholder string
	)

	for n := 0; n < b.N; n++ {
		placeholder = GetFullname("Dmitriy", "Mamykin", "8-937-333-33-33")
	}

	_ = placeholder
}
