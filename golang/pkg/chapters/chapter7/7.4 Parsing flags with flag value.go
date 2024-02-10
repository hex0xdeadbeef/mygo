package chapter7

import (
	"flag"
	"fmt"
	. "golang/pkg/chapters/chapter2/e_namedtypes"
	tc "golang/pkg/chapters/chapter2/e_namedtypes/convertions/tempconv"

	"time"
)

func Period() {
	var period = flag.Duration("period", 1*time.Second, "sleep period")

	flag.Parse()
	fmt.Printf("Sleeping for %v...", *period)
	time.Sleep(*period)
	fmt.Println()
}

// Celsius does already have the String method, so we need just to declare Set(string)
type celsiusFlag struct{ Celsius }

func (f *celsiusFlag) Set(s string) error {
	var value float64
	var unit string

	fmt.Sscanf(s, "%f%s", &value, &unit)

	switch unit {
	case "C", "*C":
		f.Celsius = Celsius(value)
		return nil

	case "F", "*F":
		f.Celsius = FToC(Farenheit(value))
		return nil
	case "K", "*K":
		f.Celsius = Celsius(tc.KToC(tc.Kelvin(value)))

		return nil
	default:
		return fmt.Errorf("invalid temperature %q", s)
	}
}

func CelsiusFlag(name string, value Celsius, usage string) *Celsius {
	f := celsiusFlag{Celsius: value}
	flag.CommandLine.Var(&f, name, usage)
	return &f.Celsius
}

func FlagUsing() {
	var temp = CelsiusFlag("temp", 20.0, "the temperature")

	flag.Parse()
	fmt.Println(temp)
}
