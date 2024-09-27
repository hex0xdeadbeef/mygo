package semicolons

import "fmt"

func A() {
	// if 10 < B() //;  unexpected newline, expecting { after if clause
	// {

	// }

	if 10 < B() {
		fmt.Println()
	}
}

func B() int {
	return -1
}
