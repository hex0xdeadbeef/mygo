package forloop

import "fmt"

func ForString() {
	for pos, char := range "абвгд\x8012" {
		fmt.Println(pos, char)
	}
}

func MultipleArgsFor() {
	var (
		str = []byte("abcde")
	)

	for i, j := 0, len(str)-1; i < j; i, j = i+1, j-1 {
		str[i], str[j] = str[j], str[i]
	}
}
