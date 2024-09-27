package deepnumberequivalence

import (
	"fmt"
	"math"
	"reflect"
)

func IsNumbersEqual(a, b any) (bool, error) {

	const (
		conversionErrorTemplate = "converting %q's value to float64: %v"
	)

	var (
		reflectA = reflect.ValueOf(a)
		reflectB = reflect.ValueOf(b)
	)

	if aConverted, err := isNumeric(reflectA); err != nil {
		return false, fmt.Errorf(conversionErrorTemplate, "a", err)
	} else if bConverted, err := isNumeric(reflectB); err != nil {
		return false, fmt.Errorf(conversionErrorTemplate, "b", err)
	} else {
		return deepNumericConversion(aConverted, bConverted), nil
	}

}

func isNumeric(v reflect.Value) (float64, error) {
	switch v.Kind() {
	// Int
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(v.Int()), nil
		// Uint
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return float64(v.Uint()), nil
		// Float
	case reflect.Float32, reflect.Float64:
		return v.Float(), nil
		// Complex
	case reflect.Complex64, reflect.Complex128:
		return real(v.Complex()), nil
	default:
		return 0, fmt.Errorf("non numeric type provided")
	}
}

func deepNumericConversion(a, b float64) bool {
	return math.Abs(a-b) < 10e-9
}
