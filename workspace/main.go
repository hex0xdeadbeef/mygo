package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

func main() {
	fmt.Println("Func result:", NamedReturn())
}

// Defer with closure close phuong-secrets.txt: file already closed
// Func result: close phuong-secrets.txt: file already closed
func NamedReturn() (err error) {
	const (
		filePath = "phuong-secrets.txt"
	)

	f, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("opening the file %q: %w", filePath, err)
	}
	defer func() {
		err = errors.Join(err, f.Close())
		fmt.Println("Defer with closure", err)
	}()

	sr := bufio.NewScanner(f)
	for sr.Scan() {
		_ = sr.Text()
	}

	// It's made deliberately to show the changes in defer statement
	f.Close()

	return
}
