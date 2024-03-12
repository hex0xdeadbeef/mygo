package inputstring

import (
	"fmt"
	"os"
	"strconv"
)

func Using() {
	if len(os.Args) != 2 {
		fmt.Println("Print provide an integer.")
		return
	}

	n, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}

	input := strconv.Itoa(n)
	fmt.Println(input)

	input = strconv.FormatInt(int64(n), 10)
	fmt.Println(input)

	input = string(rune(n))
	fmt.Println(input)

}
