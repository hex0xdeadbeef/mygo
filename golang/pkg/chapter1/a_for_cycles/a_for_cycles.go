package a_for_cycles

import (
	"fmt"
	"os"
	"strings"
)

func Echo_SimpeFor() {
	var separator, resultString string
	separator = " "
	for i := 0; i < len(os.Args); i++ {
		resultString += os.Args[i] + separator
	}
	fmt.Println(resultString)
}

func For_While() {
	number := 0
	for number < 10 {
		number++
	}
	fmt.Println(number)
}

func For_Infinite() {
	number := 2
	for {
		if number < 10e10 {
			number *= 2
		} else {
			break
		}

	}
	fmt.Println(number)
}

func Echo_RangeFor() {
	str := ""
	for index, cmdArgument := range os.Args {
		fmt.Print(index, " ")
		str += cmdArgument + " "

	}
	fmt.Println(str)
}

func Echo_stingsJoin() {
	fmt.Println(strings.Join(os.Args[0:], " "))
}

func Echo_StraightforwardPrint() {
	fmt.Println(os.Args[0:])
}
