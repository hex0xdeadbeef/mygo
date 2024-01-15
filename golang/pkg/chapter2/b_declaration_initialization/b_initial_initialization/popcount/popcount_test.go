package popcount

import (
	"testing"
)

func BenchmarkPopCountSimple(b *testing.B) {
	for i := 0; i < b.N; i++ {
		PopCounSimple(32769)
	}
}

func BenchmarkPopCountLoop(b *testing.B) {
	for i := 0; i < b.N; i++ {
		PopCountLoop(32769)
	}
}

func BenchmarkPopCountWithOneBitShifting(b *testing.B) {
	for i := 0; i < b.N; i++ {
		PopCountWithOneBitShifting(32769)
	}
}

func BenchmarkPopCountWithSibractionPrecedingValue(b *testing.B) {
	for i := 0; i < b.N; i++ {
		PopCountWithSubractionPrecedingValue(32769)
	}
}
