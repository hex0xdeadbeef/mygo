package d_booleans

import "fmt"

/* 1. %t verb allows to shows the boolean type in a formatted output */

func shortCircuitForms() {
	str := ""
	fmt.Println(str != "" && str[0] == 'x') // The right side won't be evaluated
}

func combinationOrAnd() {
	// && has higher precedence than || so parentheses isn't required
	if 'a' <= 'c' && 'c' <= 'z' ||
		'A' <= 'c' && 'c' <= 'Z' {
		fmt.Println("Hello!")
	}
}

func explicitConversionFromIntToBoolean(value int) bool {
	if value > 50 {
		return true
	} else {
		return false
	}
}
