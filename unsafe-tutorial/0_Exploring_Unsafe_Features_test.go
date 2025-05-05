package main

import "testing"

// goos: darwin
// goarch: arm64
// pkg: unsafe-tutorial
// cpu: Apple M3 Pro
// BenchmarkStringConversionFast-12        16081910                64.33 ns/op           16 B/op          1 allocs/op
// BenchmarkStringConversionDefault-12     16904682                68.04 ns/op           16 B/op          1 allocs/op

func BenchmarkStringConversionFast(b *testing.B) {
	for b.Loop() {
		_ = StringConversionOptimization()
	}

}

func BenchmarkStringConversionDefault(b *testing.B) {
	for b.Loop() {
		_ = StringConversionDefault()
	}
}
