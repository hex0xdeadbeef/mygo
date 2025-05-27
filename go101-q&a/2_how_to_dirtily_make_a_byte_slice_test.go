package main

import (
	"strings"
	"testing"
)

// goos: darwin
// goarch: arm64
// pkg: go101qa
// cpu: Apple M3 Pro
// Benchmark_MakeByteSlice-12                896134              1293 ns/op
// Benchmark_make_append-12                  836109              1362 ns/op

func Benchmark_MakeDirtyByteSlice(tb *testing.B) {
	for tb.Loop() {
		b := MakeDirtyByteSlice(Len)
		_ = b
	}
}

func Benchmark_make(tb *testing.B) {
	for tb.Loop() {
		b := make([]byte, Len)
		_ = b
	}
}

func MakeByteSlice(initialValues ...[]byte) []byte {
	n := 0
	for i := range initialValues {
		n += len(initialValues[i])
	}

	n, r := 0, MakeDirtyByteSlice(n)
	for _, s := range initialValues {
		n += copy(r[n:], s)
	}
	return r
}

const N = Len / 4

var x = []byte(strings.Repeat("x", N))
var y = []byte(strings.Repeat("y", N))
var z = []byte(strings.Repeat("z", N))
var q = []byte(strings.Repeat("z", N))
var r []byte

func Benchmark_MakeByteSlice(tb *testing.B) {
	for tb.Loop() {
		r = MakeByteSlice(x, y, z, q)
	}
}

func Benchmark_make_append(tb *testing.B) {
	for tb.Loop() {
		r = make([]byte, 0, Len)
		r = append(r, x...)
		r = append(r, y...)
		r = append(r, z...)
		r = append(r, q...)
	}
}
