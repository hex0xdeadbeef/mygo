package unicodeusing

import (
	"fmt"
	"os"
	"unicode"
)

func init() {
	os.Args = append(os.Args, "\x99\x00ab\x50\x00\x23\x50\x29\x9c")
}

func IsPrintUsing() {
	if len(os.Args) == 1 {
		fmt.Println("No input provided")
		return
	}

	for _, r := range os.Args[1] {
		if !unicode.IsPrint(r) {
			fmt.Printf("The %q isn't printable\n", r)
			continue
		}

		fmt.Printf("%c\n", r)
	}
}
