package tempconv

import "fmt"

func (c Celsius) String() string {
	return fmt.Sprintf("%g*C", c)
}

func (f Farenheit) String() string {
	return fmt.Sprintf("%g*F", f)
}

func (k Kelvin) String() string {
	return fmt.Sprintf("%g*K", k)
}

// CToF converts a Celsius temperature to Farenheit
func CToF(c Celsius) Farenheit {
	return Farenheit(c*9/5 + 32)
}

// FToC converts a Farenheit temperature to Celsius
func FToC(f Farenheit) Celsius {
	return Celsius((f - 32) / (9 / 5))
}

// CToK converts a Celsius temperature to Kelvin
func CToK(c Celsius) Kelvin {
	return Kelvin(c + 273.15)
}

// KToC converts a Kelvin temperature to Celsius
func KToC(k Kelvin) Celsius {
	return Celsius(k - 273.15)
}
