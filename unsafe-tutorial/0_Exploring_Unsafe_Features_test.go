package main

import "testing"

// goos: darwin
// goarch: arm64
// pkg: unsafe-tutorial
// cpu: Apple M3 Pro
// BenchmarkStringConversionFast-12        1000000000               1.000 ns/op
// BenchmarkStringConversionDefault-12     436401061                2.812 ns/op

func BenchmarkStringConversionFast(b *testing.B) {
	for b.Loop() {
		StringConversionOptimization()
	}

}

func BenchmarkStringConversionDefault(b *testing.B) {
	for b.Loop() {
		StringConversionDefault()
	}
}
