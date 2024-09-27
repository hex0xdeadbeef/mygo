package equal

import (
	"reflect"
	"unsafe"
)

type empty struct{}

type comparison struct {
	x, y unsafe.Pointer
	t    reflect.Type
}

func equal(x, y reflect.Value, seen map[comparison]empty) bool {
	if !x.IsValid() || !y.IsValid() {
		return x.IsValid() == y.IsValid()
	}

	if x.Type() != y.Type() {
		return false
	}

	// cycle check
	if x.CanAddr() && y.CanAddr() {
		xptr := unsafe.Pointer(x.UnsafeAddr())
		yptr := unsafe.Pointer(y.UnsafeAddr())
		if xptr == yptr {
			return true
		} // identrical references

		c := comparison{xptr, yptr, x.Type()}
		if _, ok := seen[c]; ok {
			return true
		}

		seen[c] = empty{}
	}

	switch x.Kind() {
	case reflect.Bool:
		return x.Bool() == y.Bool()
	case reflect.String:
		return x.String() == y.String()

	// numeric cases omitted for brevity
	case reflect.Chan, reflect.UnsafePointer, reflect.Func:
		return x.Pointer() == y.Pointer()

	case reflect.Pointer, reflect.Interface:
		return equal(x.Elem(), y.Elem(), seen)

	case reflect.Array, reflect.Slice:
		if x.Len() != y.Len() {
			return false
		}
		for i := 0; i < x.Len(); i++ {
			if !equal(x.Index(i), y.Index(i), seen) {
				return false
			}
		}
		return true
		// struct and map cases omitted for brevity
	}

	panic("unreachable")
}

// Equal reports whether x and y are deeply equal
func Equal(x, y any) bool {
	seen := make(map[comparison]empty)
	return equal(reflect.ValueOf(x), reflect.ValueOf(y), seen)
}
