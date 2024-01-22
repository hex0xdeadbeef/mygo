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

const (
	AbsoluteZeroC Celsius = -273.15
	FreezingC     Celsius = 0
	BoilingC      Celsius = 100
)

type Celsius float64 // This is the "Named Type"
type Farenheit float64

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
	return Celsius((f - 32) * 5 / 9) // (f - 32) is placed on the left-hand side of the expression for casting the result
	// up to float64
}

// Named types behave as their underlying types, but operations between two types that have the same underlying
// types are forbidden (We should converse one type to another)
func NamedTypesStatements() {

	var (
		c Celsius
		f Farenheit
	)

	fmt.Printf("%g\n", BoilingC-FreezingC) // This operation is equal to float64-float64 (underlying types are the same)
	boilingF := CToF(BoilingC)
	fmt.Printf("%g\n", boilingF-CToF(FreezingC))
	// fmt.Printf("%g\n" boilboilingF - BoilingC) compile error: types mismatch

	// Comparison
	fmt.Println(c == 0)
	fmt.Println(f >= 0)
	// fmt.Println(c == f) compile error: types mismatch
	fmt.Println(c == Celsius(f)) // true because of conversion f from Farenheit to Celsius

	c = FToC(212.0)

	// When we declare String() function for a named type, there is no need to call it explicitly
	fmt.Println(c.String()) // The explicit call of the String() function

	fmt.Printf("%s\n", c) // The implicit call of the String() function
	fmt.Println(c)        // Another implicit call

	fmt.Printf("%g\n", c)   // Printing the value as float64 number with "g" verb implicitly
	fmt.Println(float64(c)) // Printing the value as float64 number with "g" verb explicitly
}

func Package_Using() {
	fmt.Printf("Brrrr! %v\n", tempconv.AbsoluteZeroC)
	fmt.Println(tempconv.CToF(tempconv.BoilingC))
}

func ParseCmdArgConvertToCorF() {
	for _, arg := range os.Args[1:] { // Indexes are skipped
		temperature, err := strconv.ParseFloat(arg, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "cf %v\n", err)
			os.Exit(1)
		}
		f := tempconv.Farenheit(temperature)
		c := tempconv.Celsius(temperature)

		fmt.Printf("%s = %s, %s = %s\n", f, tempconv.FToC(f), c, tempconv.CToF(c))
	}
}

func GeneralPurposeUnitConversion() {
	if len(os.Args) > 1 {
		for _, arg := range os.Args[1:] { // Read cmd args
			unit, err := strconv.ParseFloat(arg, 64) // Parse float value
			if err != nil {
				fmt.Fprintf(os.Stderr, "The error is: %v\n", err)
				os.Exit(1)
			}
			unitsConvertation(unit) // Print converted values
		}
	} else {
		scanner := bufio.NewScanner(os.Stdin) // Create scanner
		scanner.Scan()                        // Scan a cmd input

		userInput := scanner.Text()                    // Get a text from the input
		userInputArgs := strings.Split(userInput, " ") // Split and put into a string array the args we got as a string

		for _, arg := range userInputArgs {
			unit, err := strconv.ParseFloat(arg, 64) // Try to parse the next float value
			if err != nil {
				errorPrintExit(err)
			}
			unitsConvertation(unit) // // Print converted values
		}

	}
}

func errorPrintExit(err error) {
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
