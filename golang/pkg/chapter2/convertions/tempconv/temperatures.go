package tempconv // It defines the package name

type Celsius float64   // This type is visible to all files in tempconv package
type Farenheit float64 // This type is visible to all files in tempconv package
type Kelvin float64    // This type is visible to all files in tempconv package

const (
	AbsoluteZeroC Celsius = -273.15 // This type is visible to all files in tempconv package
	FreezingC     Celsius = 0       // This type is visible to all files in tempconv package
	BoilingC      Celsius = 0       // This type is visible to all files in tempconv package

	ZeroKerlvin Kelvin = -273.15
)
