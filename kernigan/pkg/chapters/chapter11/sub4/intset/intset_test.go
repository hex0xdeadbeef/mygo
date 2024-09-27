package intset_test

import (
	"fmt"
	"golang/pkg/chapters/chapter11/sub4/intset/mapbasedintset"
	"golang/pkg/chapters/chapter6"
	"log"
	"math"
	"math/rand"
	"testing"
	"time"
)

const (
	bitsetSize = 1e6

	powIncrement = 8
	limitsNumber = 4
)

var (
	seed int64
	rnd  *rand.Rand

	wordSizeLimits []int
)

func init() {
	seed = time.Now().UTC().UnixNano()
	rnd = rand.New(rand.NewSource(seed))

	wordSizeLimits = make([]int, 0, limitsNumber)
	for i := 1; i <= limitsNumber; i++ {
		wordSizeLimits = append(wordSizeLimits, int(math.Pow(2, float64(i*powIncrement))))
	}
}

type filler interface {
	Add(...int)
	Clear() error
}

func addbenchmarkWords(b *testing.B, flr filler, wordSizeLimit int) {
	for i := 0; i < b.N; i++ {
		for i := 0; i < bitsetSize; i++ {
			flr.Add(rnd.Intn(wordSizeLimit))
		}
		flr.Clear()
	}
}

func BenchmarkAdd(b *testing.B) {

	var (
		fillers = []filler{
			chapter6.New(),
			mapbasedintset.New(),
		}
	)

	log.Printf("SEED: %d", seed)

	for _, limit := range wordSizeLimits {
		for _, filler := range fillers {
			b.Run(fmt.Sprintf("SIZE LIMIT: %d", limit),
				func(b *testing.B) {
					limit := limit
					addbenchmarkWords(b, filler, limit)
				})

		}
	}

}
