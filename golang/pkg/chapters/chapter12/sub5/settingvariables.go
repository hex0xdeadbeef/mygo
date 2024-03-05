package sub5

import (
	"fmt"
	"os"
	"reflect"
)

func AddressableVariables() {
	var (
		x = 100
	)

	fmt.Println(x)

	a := reflect.ValueOf(2)
	b := reflect.ValueOf(x)
	c := reflect.ValueOf(&x)

	d := c.Elem()
	d.SetInt(1000)

	fmt.Println(x)

	fmt.Println(a.CanAddr())
	fmt.Println(b.CanAddr())
	fmt.Println(c.CanAddr())
	fmt.Println(d.CanAddr())

}

func VariableRecovering() {
	var (
		x = 100
	)

	// Creating an addressable Value
	addressableVal := reflect.ValueOf(&x).Elem()
	fmt.Println(addressableVal.Addr())

	addressableVal.Set(reflect.ValueOf(1))
	fmt.Println(x)

	// addressableVal.Set(reflect.ValueOf(2.0)) // "reflect.Set: value of type float64 is not assignable to type int"

	// unaddressableValue := reflect.ValueOf(x)
	// unaddressableValue.Set(reflect.ValueOf(3.0)) // "reflect: reflect.Value.Set using unaddressable value"

	addressableVal.SetInt(4)
	fmt.Println(x)

	// "Value" dereferencing
	xPtr := addressableVal.Addr().Interface().(*int)
	*xPtr = 5
	fmt.Println(x)
}

func TrickyOccasion() {
	var (
		x int
		y interface{}

		xPtr = reflect.ValueOf(&x).Elem()
		yPtr = reflect.ValueOf(&y).Elem()
	)

	xPtr.Set(reflect.ValueOf(100))
	xPtr.SetInt(200)
	// xPtr.Set(reflect.ValueOf("hello")) // "reflect.Set: value of type string is not assignable to type int"
	// xPtr.SetString("hello") // "reflect: call of reflect.Value.SetString on int Value"

	yPtr.Set(reflect.ValueOf(100))     // OK
	yPtr.Set(reflect.ValueOf("hello")) // OK
	// yPtr.SetInt(400) // "reflect: call of reflect.Value.SetInt on interface Value"
}

func TryToChangeUnexportedField() {
	stdoutAddr := reflect.ValueOf(os.Stdout).Elem()
	fmt.Println(stdoutAddr.Type())

	if !stdoutAddr.CanAddr() {
		fmt.Fprint(os.Stderr, "unaddressable variable")
		return
	}

	name := stdoutAddr.FieldByName("name")
	if name.IsZero() {
		fmt.Fprint(os.Stderr, "no fields matches")
		return
	}

	// read
	fmt.Println(name.String()) // OK
	// write
	// name.SetString("bebra") // "reflect: reflect.Value.SetInt using value obtained using unexported ...
	if !name.CanSet() {
		fmt.Fprintf(os.Stderr, "setting field %q: unsettabe field", name.String())
		return
	}
	name.SetString("bebra")

}
