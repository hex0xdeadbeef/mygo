package a_introduction

import "fmt"

const boilingF = 212.0
func BoilingCelsiusFarenheit() {
	var farenheit = boilingF
	var celsius = (farenheit - 32) * 5 / 9
	fmt.Printf("boiling point = %.2f F* or %.2f C*\n", farenheit, celsius)
}

func BoilingCelsiusFarenheitReturnable(farenheit float64) float64 {
	return (farenheit - 32) * 5 / 9
}