package main

import (
	"fmt"
	e "golang/pkg/chapters/chapter4/e_functions"
)

func main() {
	fmt.Println(e.Placeholder("My fullname is $firstname $lastname and my age is $age", e.Replacements))
}
