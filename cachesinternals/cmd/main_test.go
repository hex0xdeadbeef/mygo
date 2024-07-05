package main

import (
	"sync"
	"testing"
	"unsafe"
)

func BenchmarkMatrixCombination(b *testing.B) {
	const (
		size = 1 << 14
	)

	matrixA, matrixB := createMatrix(size), createMatrix(size)

	/*
		BEFORE OPTIMISATION
	*/

	// matrixA[i][j] and matrixB[i][j]
	/*
		BenchmarkMatrixCombination-12    	       1	15_054_902_375 ns/op	17181462072 B/op	   65571 allocs/op
	*/

	// matrixA[i][j] and matrixB[j][i]
	/*
		BenchmarkMatrixCombination-12    	       1	18_212_164_625 ns/op	17181458680 B/op	   65566 allocs/op
	*/

	/*
		AFTER OPTIMISATION
	*/
	// matrixA[i][j] and matrixB[i][j] before optimization
	/*
		BenchmarkMatrixCombination-12              27343             42752 ns/op              38 B/op          0 allocs/op
		BenchmarkMatrixCombination-12               6198            193675 ns/op             681 B/op          0 allocs/op
		BenchmarkMatrixCombination-12                  1        3_720_439_625 ns/op        4295764064 B/op    32794 allocs/op
	*/

	// matrixA[i][j] and matrixB[j][i]
	/*
		BenchmarkMatrixCombination-12              19419             57444 ns/op              54 B/op          0 allocs/op
		BenchmarkMatrixCombination-12               4894            261444 ns/op             862 B/op          0 allocs/op
		BenchmarkMatrixCombination-12                  1        4_556_983_625 ns/op        4295763960 B/op    32792 allocs/op

	*/

	// The difference with 256 size is 34%
	// The difference with 512 size is 34,9%
	// The difference with 16_384 is 22,5%

	for n := 0; n < b.N; n++ {
		MatrixCombination(matrixA, matrixB)
	}
}

func BenchmarkMatrixCombinationLoopBlocking(b *testing.B) {
	const (
		size = 1 << 14
	)

	matrixA, matrixB := createMatrix(size), createMatrix(size)

	// matrixA[i][j] and matrixB[i][j]
	/*
		32 elems BenchmarkMatrixCombination-12    	       1	4_453_703_417 ns/op	17181462072 B/op	   65571 allocs/op
		64 elems  BenchmarkMatrixCombinationLoopBlocking-12              1        4_387_627_625 ns/op        4295760952 B/op    32791 allocs/op
	*/

	// matrixA[i][j] and matrixB[j][i]
	/*
		32 elems  BenchmarkMatrixCombinationLoopBlocking-12              1        4_883_716_459 ns/op        4295764056 B/op    32793 allocs/op
		64 elems BenchmarkMatrixCombinationLoopBlocking-12              1        4_408_456_000 ns/op        4295760760 B/op    32789 allocs/op
	*/

	for n := 0; n < b.N; n++ {
		MatrixCombinationLoopBlocking(matrixA, matrixB)
	}
}

// This structures are arranged in the memory contiguously in the case of threads living on different cores (which is not necessarily in this case)
// Each goroutine gets its own L1 cache and the variables are marked "Shared" and then "Modified"
// Yet, as the structures are contiguous in memory, it'll be present in both cache lines.

/*
	Before memory padding
*/
//	BenchmarkFalseSharing-12            1170           1049223 ns/op              57 B/op          2 allocs/op
/*
	After memory padding
*/
// BenchmarkAlignment-12               1117            909322 ns/op              68 B/op          2 allocs/op

func BenchmarkFalseSharing(b *testing.B) {

	type (
		SimpleStruct struct {
			n int
		}
	)

	const (
		cyclesCount = 1e6
	)

	var (
		wg sync.WaitGroup
	)

	q, w := SimpleStruct{}, SimpleStruct{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wg.Add(2)
		go func() {
			for i := 0; i < cyclesCount; i++ {
				q.n += i
			}
			wg.Done()
		}()

		go func() {
			for i := 0; i < cyclesCount; i++ {
				w.n += i
			}
			wg.Done()
		}()

		wg.Wait()
	}

}

func BenchmarkAlignment(b *testing.B) {

	const (
		sizeOfL1DCache = 1 << 4
	)

	type (
		CacheLinePad struct {
			_ [sizeOfL1DCache - 1]int64
		}

		SimpleStruct struct {
			n int

			CacheLinePad
		}
	)

	const (
		cyclesCount = 1e6
	)

	var (
		wg sync.WaitGroup
	)

	q, w := SimpleStruct{}, SimpleStruct{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wg.Add(2)
		go func() {
			for i := 0; i < cyclesCount; i++ {
				q.n += i
			}
			wg.Done()
		}()

		go func() {
			for i := 0; i < cyclesCount; i++ {
				w.n += i
			}
			wg.Done()
		}()

		wg.Wait()
	}

}

/*
BenchmarkSliceLooping-12            4048            251123 ns/op               0 B/op          0 allocs/op
BenchmarkSliceLooping-12            4688            249577 ns/op               0 B/op          0 allocs/op
BenchmarkSliceLooping-12            4058            251698 ns/op               0 B/op          0 allocs/op
*/
func BenchmarkSliceLooping(b *testing.B) {
	var (
		s = createArray()
	)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, v := range s {
			_ = v
		}
	}
}

// BenchmarkLinkedListLooping-12               1114           1036809 ns/op               0 B/op          0 allocs/op
// BenchmarkLinkedListLooping-12               1143           1039353 ns/op               0 B/op          0 allocs/op
// BenchmarkLinkedListLooping-12               1155           1053870 ns/op               0 B/op          0 allocs/op

func BenchmarkLinkedListLooping(b *testing.B) {
	var (
		cur = createLinkedList()
	)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		root := cur
		for cur != nil {
			cur = cur.next
		}

		cur = root
	}

}

const (
	sliceLength       = 512 * (1 << 20)
	elemsInCacheLine  = 128 / 8                        // 128 (cache line size in Bytes) / 8 (1 int64 size in Bytes)
	cacheLinesInCache = sliceLength / elemsInCacheLine // cachelines needed to fill all cache
)

/*
BenchmarkUnitStrideIteration-12              130           8_697_000 ns/op               0 B/op          0 allocs/op
*/
func BenchmarkUnitStrideIteration(b *testing.B) {

	var (
		s = make([]int64, sliceLength)
	)

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		for i := 0; i < cacheLinesInCache; i++ {
			s[i]++
		}
	}

}

// The step is the same to cacheline size.
/*
BenchmarkConstantStrideIterationA-12                   9         119_888_032 ns/op               0 B/op          0 allocs/op
*/
func BenchmarkConstantStrideIterationA(b *testing.B) {
	var (
		s = make([]int64, sliceLength)
	)

	b.ResetTimer()

	for n := 0; n < b.N; n++ {

		for i := 0; i < sliceLength; i += elemsInCacheLine {
			s[i]++
		}

	}

}

// The step is the same to cacheline size * 8
/*
BenchmarkConstantStrideIterationB-12                   5         239_941_883 ns/op               0 B/op          0 allocs/op
*/
func BenchmarkConstantStrideIterationB(b *testing.B) {
	var (
		s = make([]int64, sliceLength)
	)

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		for i := 0; i < sliceLength; i += elemsInCacheLine * 2 {
			for j := i; j < i+elemsInCacheLine*2; j++ {
				s[j]++
			}
		}
	}

}

// BenchmarkFalseSharingA-12            931           1243620 ns/op             343 B/op          8 allocs/op
func BenchmarkFalseSharingA(b *testing.B) {
	var (
		wg = &sync.WaitGroup{}

		s = []int{0: 100, 1: 200, 2: 300, 3: 400, 48: -1}
	)
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		for i := 0; i < 4; i++ {

			wg.Add(1)

			go func(incrementIndex int) {
				defer wg.Done()

				for range int(1e6) {
					s[incrementIndex]++
				}

			}(i)

		}

		wg.Wait()
	}

}

// BenchmarkFalseSharingB-12           1224           1056492 ns/op             326 B/op          8 allocs/op
func BenchmarkFalseSharingB(b *testing.B) {
	var (
		wg = &sync.WaitGroup{}

		s = []int{0: 100, 16: 200, 32: 300, 48: 400}
	)
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		wg.Add(4)
		for i := 0; i < 4; i++ {
			go func(incrementIndex int) {
				defer wg.Done()

				for range int(1e6) {
					s[incrementIndex]++
				}

			}(i)

		}

		wg.Wait()
	}

}

type Small struct {
	str string // 16 bytes
	i   int    // 8 bytes
	b   bool   // 1 byte + 7 padding bytes
}

// 4797	    248396 ns/op	       0 B/op	       0 allocs/op
func BenchmarkValue_L1(b *testing.B) {
	s := Small{"str", 1, true} // 32 bytes, fits L1 cache
	println(unsafe.Sizeof(s), unsafe.Offsetof(s.str), unsafe.Offsetof(s.i), unsafe.Offsetof(s.b))

	// prepare data
	items := make([]Small, 0, 100)
	for j := 0; j < 100; j++ {
		items = append(items, s)
	}

	// bench
	for i := 0; i < b.N; i++ {
		processSmallByVal(items)
	}
}

// BenchmarkPointer_L1-12              2354            429332 ns/op               0 B/op          0 allocs/op
func BenchmarkPointer_L1(b *testing.B) {
	s := &Small{"str", 1, true} // 32 bytes, fits L1 cache

	// prepare data
	items := make([]*Small, 0, 100)
	for j := 0; j < 100; j++ {
		items = append(items, s)
	}

	// bench
	for i := 0; i < b.N; i++ {
		processSmallByPointer(items)
	}
}

/* data manipulation */

func processSmallByVal(items []Small) {
	for _, item := range items { // get item value, item var reused to store values from slice
		for i := 0; i < 10000; i++ {
			item.i = i                                             // access item's field from cache, mutate field, store in cache line directly
			item.str = "123"                                       // same
			someint, somestr, somebool := item.i, item.str, item.b // read item's fields from cache
			_, _, _ = someint, somestr, somebool                   // stub
		}
	}
}

func processSmallByPointer(items []*Small) {
	for _, item := range items { // get pointer, store in 'item'
		for i := 0; i < 10000; i++ {
			item.i = i                                             // go to the main mem
			item.str = "123"                                       // go to the main mem again
			someint, somestr, somebool := item.i, item.str, item.b // go to the main mem again
			_, _, _ = someint, somestr, somebool                   // stub
		}
	}
}

type Big struct {
	arr [10_000]int // 78.125KB
}

// BenchmarkValue_L2-12              294324              4021 ns/op               0 B/op          0 allocs/op
func BenchmarkValue_L2(b *testing.B) {
	a := [10000]int{} // 40kb, fits L2 cache
	for i := 0; i < 10000; i++ {
		a[i] = i // fill array with values
	}
	s := Big{arr: a}

	b.ResetTimer()
	// bench
	for i := 0; i < b.N; i++ {
		processBigByVal(s)
	}
}

// BenchmarkPointer_L2-12            473072              2562 ns/op               0 B/op          0 allocs/op
func BenchmarkPointer_L2(b *testing.B) {
	a := [10000]int{} // 78.125KB
	for i := 0; i < 10000; i++ {
		a[i] = i
	}
	s := &Big{arr: a}

	b.ResetTimer()
	// bench
	for i := 0; i < b.N; i++ {
		processBigByPointer(s)
	}
}

/* data manipulation */
func processBigByVal(item Big) {
	for i := 0; i < 10000; i++ {
		item.arr[i] = i
		somearr := item.arr
		_ = somearr
	}
}

func processBigByPointer(item *Big) {
	for i := 0; i < 10000; i++ {
		item.arr[i] = i
		somearr := item.arr
		_ = somearr
	}
}
