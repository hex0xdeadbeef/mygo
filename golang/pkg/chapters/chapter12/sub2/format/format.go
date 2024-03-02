package format

import (
	"fmt"
	"reflect"
	"strconv"
	"time"
)

func main() {
	var (
		x int64         = 1
		d time.Duration = 1 * time.Nanosecond
	)

	fmt.Println(Any(x))
	fmt.Println(Any(d))
}

func Any(value interface{}) string {
	return formatAtom(reflect.ValueOf(value))
}

func FormatAtom(v reflect.Value) string {
	return formatAtom(v)
}

func formatAtom(v reflect.Value) string {
	switch v.Kind() {
	// Invalid
	case reflect.Invalid:
		return "invalid"
		// Ints
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10)
		// Uints
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return strconv.FormatUint(v.Uint(), 10)
		// Bool
	case reflect.Bool:
		return strconv.FormatBool(v.Bool())
		// String
	case reflect.String:
		return strconv.Quote(v.String())
		// Reference types
	case reflect.Chan, reflect.Func, reflect.Ptr, reflect.Slice, reflect.Map:
		return v.Type().String() + " 0x" + strconv.FormatUint(uint64(v.Pointer()), 16)
		// Others
	default:
		return v.Type().String() + " value"
	}
}
