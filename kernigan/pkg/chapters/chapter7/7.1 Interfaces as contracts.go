package chapter7

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type ByteCounter int

func (bc *ByteCounter) Write(p []byte) (n int, err error) {
	*bc += ByteCounter(len(p))

	return len(p), nil
}

func ByteCounterUsing() {
	var bc ByteCounter

	bc.Write([]byte("Hello, sun!"))
	fmt.Println(bc)

	bc = 0
	fmt.Fprintf(&bc, "Hello, %s", "$name")
	fmt.Println(bc)
}

type WordCounter struct {
	count int
	words []string
}

func (wc *WordCounter) Write(p []byte) (n int, err error) {
	if len(wc.words) == 0 {
		wc.words = make([]string, 0, 1)
	}
	wordScanner := bufio.NewScanner(strings.NewReader(string(p)))
	wordScanner.Split(bufio.ScanWords)

	for wordScanner.Scan() {
		wc.count++
		wc.words = append(wc.words, string(wordScanner.Text()))
	}

	wc.words = wc.words[:wc.count]

	return len(p), nil
}

func (wc *WordCounter) PrintWords() {
	for _, word := range wc.words {
		fmt.Println(word)
	}
}

type LinesCounter struct {
	count int
	lines []string
}

func (wc *LinesCounter) Write(p []byte) (n int, err error) {
	if len(wc.lines) == 0 {
		wc.lines = make([]string, 0, 1)
	}
	lineScanner := bufio.NewScanner(strings.NewReader(string(p)))
	lineScanner.Split(bufio.ScanLines)

	for lineScanner.Scan() {
		wc.count++
		wc.lines = append(wc.lines, string(lineScanner.Text()))
	}

	wc.lines = wc.lines[:wc.count]

	return len(p), nil
}

func (wc *LinesCounter) PrintLines() {
	for _, line := range wc.lines {
		fmt.Println(line)
	}
}

type CountWriter struct {
	writer       io.Writer
	bytesWritten int64
}

func (cw *CountWriter) Write(p []byte) (n int, err error) {
	n, err = cw.writer.Write(p)

	cw.bytesWritten += int64(n)

	return
}

func NewCountingWriter(w io.Writer) (io.Writer, *int64) {
	cr := &CountWriter{writer: w, bytesWritten: 0}

	return cr, &cr.bytesWritten
}
