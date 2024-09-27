package lengthconv

import "fmt"

func (m Meter) String() string {
	return fmt.Sprintf("%g m.", m)
}

func (foot Foot) String() string {
	return fmt.Sprintf("%g ft.", foot)
}

// FtToM converts a Meter length to Foot
func MToFt(m Meter) Foot {
	return Foot(m * Meter(meterToFootCoefficient))
}

// FtToM converts a Foot length to Meter
func FtToM(ft Foot) Meter {
	return Meter(ft / Foot(meterToFootCoefficient))
}
