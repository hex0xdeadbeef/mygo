package main_test

import (
	"strings"
	"testing"
	"unique"
)

/*
	Running tool: /opt/homebrew/bin/go test -benchmem -run=^$ -coverprofile=/var/folders/ky/43rc3l3j2ml222448pyxcnd40000gn/T/vscode-goE8hl4y/go-code-cover -bench . uniquepkg/src

	goos: darwin
	goarch: arm64
	pkg: uniquepkg/src
	cpu: Apple M3 Pro
	BenchmarkStringCompareSmall-12             	347887723	         3.165 ns/op	       0 B/op	       0 allocs/op
	BenchmarkStringCompareMedium-12            	143908264	         8.586 ns/op	       0 B/op	       0 allocs/op
	BenchmarkStringCompareLarge-12             	  146972	      8501 ns/op	       0 B/op	       0 allocs/op
	BenchmarkStringCanonicalisingSmall-12      	1000000000	         0.2499 ns/op	       0 B/op	       0 allocs/op
	BenchmarkStringCanonicalisingmMedium-12    	1000000000	         0.2493 ns/op	       0 B/op	       0 allocs/op
	BenchmarkStringCanonicalisingmLarge-12     	1000000000	         0.2500 ns/op	       0 B/op	       0 allocs/op
	PASS
	coverage: 0.0% of statements
	ok  	uniquepkg/src	5.875s
*/

func BenchmarkStringCompareSmall(b *testing.B)  { benchStringComparison(b, 10) }
func BenchmarkStringCompareMedium(b *testing.B) { benchStringComparison(b, 100) }
func BenchmarkStringCompareLarge(b *testing.B)  { benchStringComparison(b, 100_000) }

func BenchmarkStringCanonicalisingSmall(b *testing.B)   { benchCanonicalising(b, 10) }
func BenchmarkStringCanonicalisingmMedium(b *testing.B) { benchCanonicalising(b, 100) }
func BenchmarkStringCanonicalisingmLarge(b *testing.B)  { benchCanonicalising(b, 100_000) }

func benchStringComparison(b *testing.B, count int) {
	s1 := strings.Repeat("2024,", count)
	s2 := strings.Repeat("2024,", count)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		if s1 != s2 {
			b.Fatal()
		}
	}
	b.ReportAllocs()
}

func benchCanonicalising(b *testing.B, count int) {
	s1 := unique.Make(strings.Repeat("2024,", count))
	s2 := unique.Make(strings.Repeat("2024,", count))

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		if s1 != s2 {
			b.Fatal()
		}
	}
	b.ReportAllocs()
}
