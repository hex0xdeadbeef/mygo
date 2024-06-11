package foo

// It'll br initialized first

import "fmt"

func init() {
	fmt.Println("foo done")
}
