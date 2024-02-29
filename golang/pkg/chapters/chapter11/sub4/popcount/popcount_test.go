package popcount_test

import (
	"fmt"
	"golang/pkg/chapters/chapter2/popcount"
	"log"
	"math/rand"
	"testing"
	"time"
)

var (
	counts    []int
	functions []func(uint64) int

	seed int64
	rnd  *rand.Rand
)

func init() {
	counts = append(counts, []int{1e2, 1e5, 1e6, 1e7}...)

	functions = append(functions, []func(uint64) int{
		popcount.PopCount,
		popcount.PopCountLoop,
		popcount.PopCountWithOneBitShifting,
		popcount.PopCountFast,
	}...)

	seed = time.Now().UTC().UnixNano()
	rnd = rand.New(rand.NewSource(seed))

}

func benchmark(b *testing.B, f func(x uint64) int, numbers []uint64) {
	for i := 0; i < b.N; i++ {
		for _, number := range numbers {
			f(number)
		}
	}
}

func BenchmarkPopCounts(b *testing.B) {
	var (
		numbers []uint64
	)
	log.Printf("Seed: %d\n", seed)

	for _, count := range counts {
		numbers = popcount.GenerateNumbers(rnd, count)
		for _, popFunc := range functions {

			b.Run(fmt.Sprintf("Count: %d", count), func(b *testing.B) {
				benchmark(b, popFunc, numbers)
			})
		}

	}
}
