package methods

import (
	"fmt"
	"reflect"
	"strings"
)

// Print prints the method set of the value x
func Print(x interface{}) {
	var (
		v = reflect.ValueOf(x)
		t = v.Type()

		methodType reflect.Type
	)

	fmt.Printf("type %s\n", t)
	for i := 0; i < v.NumMethod(); i++ {
		methodType = v.Method(i).Type()
		fmt.Printf("func(%s) %s%s\n", t, t.Method(i).Name, strings.TrimPrefix(methodType.String(), "func"))
	}
}
