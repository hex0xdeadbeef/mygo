package e_namedtypes

import (
	"bufio"
	"fmt"
	"golang/pkg/chapter2/e_namedtypes/convertions/lengthconv"
	"golang/pkg/chapter2/e_namedtypes/convertions/tempconv"
	"golang/pkg/chapter2/e_namedtypes/convertions/weightconv"
	"os"
	"strconv"
	"strings"
)

type Celsius float64 // This is the "Named Type"
type Farenheit float64

const (
	AbsoluteZeroC Celsius = -273.15
	FreezingC     Celsius = 0
	BoilingC      Celsius = 100
)

func (c Celsius) String() string {
	return fmt.Sprintf("%g*C", c)
}

func (f Farenheit) String() string {
	return fmt.Sprintf("%g*F", f)
}

func CToF(c Celsius) Farenheit {
	return Farenheit(c*9/5 + 32)
}

func FToC(f Farenheit) Celsius {
	return Celsius((f - 32) * 5 / 9)
}

// Named types behave as their underlying types, but operations between two types that have the same underlying
// types are forbidden (We should converse one type to another)
func NamedTypesStatements() {
	fmt.Printf("%g\n", BoilingC-FreezingC) // This operation is equal to float64-float64 (underlying types are the same)
	boilingF := CToF(BoilingC)
	fmt.Printf("%g\n", boilingF-CToF(FreezingC))
	// fmt.Printf("%g\n" boilboilingF - BoilingC) compile error: types mismatch

	var c Celsius
	var f Farenheit
	fmt.Println(c == 0)
	fmt.Println(f >= 0)
	// fmt.Println(c == f) compile error: types mismatch
	fmt.Println(c == Celsius(f)) // true because of conversion f from Farenheit to Celsius

	c = FToC(212.0)

	// When we declare String() function for a named type, there is no need to call it explicitly
	fmt.Println(c.String())
	fmt.Printf("%v\n", c) // no need to call String() explicitly
	fmt.Printf("%s\n", c)
	fmt.Println(c)
	fmt.Println("%g\n", c)
	fmt.Println(float64(c))
}

func Package_Using() {
	fmt.Println("Brrrr! %v\n", tempconv.AbsoluteZeroC)
	fmt.Println(tempconv.CToF(tempconv.BoilingC))
}

func ParseCmdArgConvertToCorF() {
	for _, arg := range os.Args[1:] {
		temperature, err := strconv.ParseFloat(arg, 64)
		if err != nil {
			fmt.Fprint(os.Stderr, "cf %v\n", err)
			os.Exit(1)
		}
		f := tempconv.Farenheit(temperature)
		c := tempconv.Celsius(temperature)

		fmt.Printf("%s = %s, %s = %s\n", f, tempconv.FToC(f), c, tempconv.CToF(c))
	}
}

func GeneralPurposeUnitConversion() {
	if len(os.Args) > 1 {
		for _, arg := range os.Args[1:] {
			unit, err := strconv.ParseFloat(arg, 64)
			if err != nil {
				fmt.Fprintf(os.Stderr, "The error is: %v\n", err)
				os.Exit(1)
			}
			unitsConvertation(unit)
		}
	} else {
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()

		userInput := scanner.Text()
		userInputArgs := strings.Split(userInput, " ")

		for _, arg := range userInputArgs {
			unit, err := strconv.ParseFloat(arg, 64)
			if err != nil {
				errorPrintAndExit(err)
			}
			unitsConvertation(unit)
		}

	}
}

func errorPrintAndExit(err error) {
	fmt.Fprintf(os.Stderr, "The error is: %v\n", err)
	os.Exit(1)
}

func unitsConvertation(unit float64) {
	farenheit := tempconv.Farenheit(unit)
	celsius := tempconv.Celsius(unit)

	kilogram := weightconv.Kilogram(unit)
	pound := weightconv.Pound(unit)

	meter := lengthconv.Meter(unit)
	foot := lengthconv.Foot(unit)

	fmt.Printf("%s = %s, %s = %s\n", farenheit, tempconv.FToC(farenheit), celsius, tempconv.CToF(celsius))
	fmt.Printf("%s = %s, %s = %s\n", kilogram, weightconv.KGToP(kilogram), pound, weightconv.PToKG(pound))
	fmt.Printf("%s = %s, %s = %s\n", meter, lengthconv.MToFt(meter), foot, lengthconv.FtToM(foot))
	fmt.Println()
}
