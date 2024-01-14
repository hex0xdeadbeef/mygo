package chapter8_9

import (
	"flag"
	"fmt"
	"math/rand"
)

func ParsingCommandLineArgs() {
	maxp := flag.Int("max", 10, "the max value")
	flag.Parse()

	fmt.Println(rand.Intn(*maxp))
}
