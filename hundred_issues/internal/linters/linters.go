package linters

import "fmt"

/*
	LINTERS USING
1. vet
	1) vet allows us to find shadowed vars
2. staticcheck helps to:
	1) find redundant break statements
	2) unused program units
	3) shadowed vars
*/

func ShadowedVar() {
	var (
		i = -1
	)

	for j := 0; j < 1; j++ {
		i := 0
		fmt.Println(i)
	}

	fmt.Println(i)
}
