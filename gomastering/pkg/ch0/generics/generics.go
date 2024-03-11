package generics

import "fmt"

func print[T any](s []T) {
	for _, v := range s {
		fmt.Print(v, " ")
	}

	fmt.Println()
}

func Using() {
	Ints := []int{1, 2, 3}
	Strings := []string{"One", "Two", "Three"}

	print(Ints)
	print(Strings)
}
