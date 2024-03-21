package queuing

import (
	"bufio"
	"concurrency/pkg/ch04/pipeline"
	"io"
	"log"
	"os"
	"testing"
)

// BenchmarkUnbufferedWrite-12    	  465237	      2566 ns/op	       1 B/op	       1 allocs/op
func BenchmarkUnbufferedWrite(b *testing.B) {
	performWrite(b, tmpFileOrFatal())
}

// BenchmarkBufferedWrite-12      	 2467894	       470.3 ns/op	       1 B/op	       1 allocs/op
func BenchmarkBufferedWrite(b *testing.B) {
	// In bufio.Writer the writes are queued internally into a buffer until a sufficient chunk has been accumulated, and then the chunk is written out.
	// This process is called chunking.
	bufferedFile := bufio.NewWriter(tmpFileOrFatal())
	performWrite(b, bufferedFile)
}

func tmpFileOrFatal() *os.File {
	file, err := os.CreateTemp("", "tmp")
	if err != nil {
		log.Fatalf("creating file: %v", err)
	}

	return file
}

func performWrite(b *testing.B, writer io.Writer) {
	done := make(chan struct{})
	defer close(done)

	b.ResetTimer()
	for bt := range pipeline.LimitedJunction(done, pipeline.Repeat(done, byte(0)), b.N) {
		writer.Write([]byte{bt.(byte)})
	}
}
