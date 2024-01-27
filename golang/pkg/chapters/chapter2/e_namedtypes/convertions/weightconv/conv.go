package weightconv

import "fmt"

func (kg Kilogram) String() string {
	return fmt.Sprintf("%g kg.", kg)
}

func (p Pound) String() string {
	return fmt.Sprintf("%g lb.", p)
}

// KGToP converts a Kilogram weight to Pound
func KGToP(kg Kilogram) Pound {
	return Pound(kg * Kilogram(kgToPoundCoefficient))
}

// PToKG converts a Pound weight to Kilogram
func PToKG(p Pound) Kilogram {
	return Kilogram(p / Pound(kgToPoundCoefficient))
}
