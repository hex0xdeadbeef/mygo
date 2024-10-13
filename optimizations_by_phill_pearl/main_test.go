package main

import (
	"math/rand"
	"testing"
)

const (
	alphabetsCount           = 2
	alphabetSize, digitsSize = 26, 10

	alphabetPad                         = 1 << 5
	alphabetStartByte, numbersStartByte = 'a', '0'

	objectsCnt    = 1 << 10
	objectMaxSize = 1 << 8
)

var (
	symbols                   []byte
	itemTypes, clientIDs, IDs []string
)

func init() {
	symbols = make([]byte, 0, alphabetSize*alphabetsCount+digitsSize)

	for i := 0; i < 26; i++ {
		symbols = append(symbols, alphabetStartByte, alphabetStartByte-alphabetPad)
	}

	for i := 0; i < 10; i++ {
		symbols = append(symbols, numbersStartByte+byte(i))
	}

	genTestSets()
}

func genTestSets() {
	itemTypes, clientIDs, IDs = make([]string, objectsCnt), make([]string, objectsCnt), make([]string, objectsCnt)

	for i := 0; i < objectsCnt; i++ {
		itemTypes[i], clientIDs[i], IDs[i] = getObject(), getObject(), getObject()
	}
}

func getObject() string {
	var (
		size = rand.Intn(objectMaxSize)
		b    = make([]byte, size)
	)

	for i := 0; i < size; i++ {
		b[i] = symbols[rand.Intn(len(symbols))]
	}

	return string(b)
}

// BenchmarkNaiveConcat-12    	12800488	        87.87 ns/op	     401 B/op	       0 allocs/op
func BenchmarkNaiveConcat(b *testing.B) {
	for n := 0; n < b.N; n++ {
		NaiveConcat(itemTypes[rand.Intn(objectsCnt)], clientIDs[rand.Intn(objectsCnt)], IDs[rand.Intn(objectsCnt)])
	}
}

// BenchmarkFmtSprintfConcatt-12    	 6925022	       167.9 ns/op	     457 B/op	       3 allocs/op
func BenchmarkFmtSprintfConcatt(b *testing.B) {
	for n := 0; n < b.N; n++ {
		FmtSprintfConcat(itemTypes[rand.Intn(objectsCnt)], clientIDs[rand.Intn(objectsCnt)], IDs[rand.Intn(objectsCnt)])
	}
}

// BenchmarkBytesBufferConcat-12    	 4163598	       283.2 ns/op	    2047 B/op	       3 allocs/op
func BenchmarkBytesBufferConcat(b *testing.B) {
	for n := 0; n < b.N; n++ {
		BytesBufferConcat(itemTypes[rand.Intn(objectsCnt)], clientIDs[rand.Intn(objectsCnt)], IDs[rand.Intn(objectsCnt)])
	}
}
