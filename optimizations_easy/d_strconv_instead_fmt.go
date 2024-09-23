package main

import (
	"fmt"
	"strconv"
)

/*
	USING "strconv" INSTEAD OF "fmt" package
1. "fmt" is the most used package in Go. But when it comes to types convertation, it's better to use "strconv".
2. The "strconv" shows better performance relatively to "fmt"
3. Because the "fmt" accepts interface{} parameters mostly we have some cons:
	- We lose type safety
	- We can rise an amount of memory allocations. Passing values instead of pointers to values leads to allocations in heap
4. The results of benchmarks:
	BenchmarkStrconvFmt100-12          	12075698	        83.19 ns/op	     136 B/op	       3 allocs/op
	BenchmarkStrconvFmt100000-12       	   96361	     12383 ns/op	  213149 B/op	       4 allocs/op
	BenchmarkStrconvFmt100000000-12    	     270	   4342680 ns/op	200017686 B/op	       6 allocs/op
	BenchmarkFmtFmt100-12              	28344726	        42.37 ns/op	     120 B/op	       2 allocs/op
	BenchmarkFmtFmt100000-12           	  223352	      5502 ns/op	  106505 B/op	       2 allocs/op
	BenchmarkFmtFmt100000000-12        	     510	   2459934 ns/op	100007958 B/op	       2 allocs/op
*/

func StrconvFmt(s string, n int) string {
	return s + ":" + strconv.Itoa(n)
}

func FmtFmt(s string, n int) string {
	return fmt.Sprintf("%s:%d", s, n)
}
