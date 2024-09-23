package main

import (
	"math/rand"
	"testing"
)

/*
SMALL SIZE OBJECTS
*/

// BenchmarkSumA-12    	1000000000	         0.5294 ns/op	       0 B/op	       0 allocs/op
func BenchmarkSumA(b *testing.B) {
	var (
		ret int
	)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		obj := A{a: n, d: n}
		ret = SumA(obj)
	}

	_ = ret
}

// BenchmarkSumB-12    	1000000000	         0.5147 ns/op	       0 B/op	       0 allocs/op
func BenchmarkSumB(b *testing.B) {
	var (
		res int
	)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		obj := B{a: n, d: n, e: n}
		res = SumB(obj)
	}

	_ = res
}

/*
	THE LACK OF INTERFACES
*/

// BenchmarkSumPointer-12    	100000000	        10.23 ns/op	      16 B/op	       1 allocs/op
func BenchmarkSumPointer(b *testing.B) {
	var (
		res int
		p   Summer
	)

	p = Sum{}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		obj := Object{A: n, B: n}
		res = p.SumPointer(&obj)
	}

	_ = res

}

// BenchmarkSumValue-12      	894017342	         1.309 ns/op	       0 B/op	       0 allocs/op
func BenchmarkSumValue(b *testing.B) {
	var (
		res int
		p   Summer
	)

	p = Sum{}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		obj := Object{A: n, B: n}
		res = p.SumValue(obj)
	}

	_ = res
}

// BenchmarkSumUsualValue-12    	1000000000	         0.2576 ns/op	       0 B/op	       0 allocs/op
func BenchmarkSumUsualValue(b *testing.B) {
	var (
		res int
		p   = Sum{}
	)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		obj := Object{A: n, B: n}
		res = p.SumUsualValue(obj)
	}

	_ = res
}

// BenchmarkSumUsualPointer-12    	1000000000	         0.2531 ns/op	       0 B/op	       0 allocs/op
func BenchmarkSumUsualPointer(b *testing.B) {
	var (
		res int
		p   = Sum{}
	)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		obj := Object{A: n, B: n}
		res = p.SumUsualPointer(&obj)
	}

	_ = res
}

/*
	SEARCH THE REALIZATION OF INTERFACE'S FUNCTION IN "itab"
*/

// BenchmarkMakeDumpsA-12    	303237055	         3.740 ns/op	       0 B/op	       0 allocs/op
func BenchmarkMakeDumpsA(b *testing.B) {
	var (
		dumpers = []Dumper{fileDumper{}, dbDumper{}}
	)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		MakeDumpsA(dumpers...)
	}

}

// BenchmarkMakeDumpsB-12    	504708154	         2.275 ns/op	       0 B/op	       0 allocs/op
func BenchmarkMakeDumpsB(b *testing.B) {
	var (
		dumpKinds = []DumperCustom{
			{kind: 1},
			{kind: 2},
		}
	)
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		MakeDumpsB(dumpKinds...)
	}

}

/*
	STACK TRICKS
*/

// BenchmarkSliceAllocationA-12    	1000000000	         0.2608 ns/op	       0 B/op	       0 allocs/op
func BenchmarkSliceAllocationA(b *testing.B) {
	for n := 0; n < b.N; n++ {
		SliceAllocationA()
	}
}

// BenchmarkSliceAllocationB-12    	     142	   7775394 ns/op	1073741877 B/op	       1 allocs/op
func BenchmarkSliceAllocationB(b *testing.B) {
	for n := 0; n < b.N; n++ {
		SliceAllocationB()
	}
}

// BenchmarkSliceAllocationC-12    	1000000000	         0.2566 ns/op	       0 B/op	       0 allocs/op
func BenchmarkSliceAllocationC(b *testing.B) {
	for n := 0; n < b.N; n++ {
		SliceAllocationC()
	}
}

/*
	BCE CALLING
*/

// BenchmarkToUint64A-12    	1000000000	         0.2596 ns/op	       0 B/op	       0 allocs/op
func BenchmarkToUint64A(b *testing.B) {
	var (
		buf = make([]byte, 8)
		res uint64
	)
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		res = ToUint64A(buf)
	}

	_ = res
}

// BenchmarkToUint64B-12    	1000000000	         0.2540 ns/op	       0 B/op	       0 allocs/op
func BenchmarkToUint64B(b *testing.B) {
	var (
		buf = make([]byte, 8)
		res uint64
	)
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		res = ToUint64B(buf)
	}

	_ = res
}

/*
	FOR RANGE/LOOP RANGING
*/

const (
	sliceSize = 1 << 30
	topNumber = 1 << 8
)

func generateSlice() []CustomIndex {
	var (
		res = make([]CustomIndex, sliceSize)
	)

	for i := 0; i < sliceSize; i++ {
		res[i] = CustomIndex{idx: rand.Intn(topNumber)}
	}

	return res
}

// BenchmarkSumRange-12    	       4	 265_833_500 ns/op	       0 B/op	       0 allocs/op
func BenchmarkSumRange(b *testing.B) {
	var (
		s   = generateSlice()
		res int
	)
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		res = SumRange(s)
	}

	_ = res
}

// BenchmarkSumLoop-12    	       4	 266_073_792 ns/op	       0 B/op	       0 allocs/op
func BenchmarkSumLoop(b *testing.B) {
	var (
		s   = generateSlice()
		res int
	)
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		res = SumLoop(s)
	}

	_ = res
}

/*
	STRING OPTIMIZATIONS
*/

// BenchmarkStringConvertationInForRange-12    	321686134	         3.216 ns/op	       0 B/op	       0 allocs/op
func BenchmarkStringConvertationInForRange(b *testing.B) {
	for n := 0; n < b.N; n++ {
		StringConvertationInForRange()
	}
}

func generateBytes() []byte {
	const (
		stringSize    = 1 << 20
		topByteNumber = 256
	)

	var (
		strBytes = make([]byte, stringSize)
	)

	for i := 0; i < stringSize; i++ {
		strBytes[i] = byte(rand.Intn(topByteNumber))
	}

	return strBytes
}

// BenchmarkStringsComparison-12    	650185662	         1.770 ns/op	       0 B/op	       0 allocs/op
func BenchmarkStringsComparison(b *testing.B) {
	var (
		sA, sB = generateBytes(), generateBytes()
	)
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		BytesComparison(sA, sB)
	}
}

func genMapAndKeys() (m map[string]struct{}, keys [][]byte) {
	const (
		mapSize = 1 << 10
	)

	m, keys = make(map[string]struct{}, mapSize), make([][]byte, 0, mapSize)

	for range mapSize {
		key := generateBytes()
		m[string(key)], keys = struct{}{}, append(keys, key)
	}

	return m, keys
}

// BenchmarkIsElemInMap-12    	      34	  34587028 ns/op	       0 B/op	       0 allocs/op
func BenchmarkIsElemInMap(b *testing.B) {
	var (
		m, keys = genMapAndKeys()
	)
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		for i := 0; i < len(keys); i++ {
			IsElemInMap(m, keys[i])
		}
	}
}
