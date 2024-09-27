package cyclicdata

import (
	"reflect"
	"unsafe"
)

type empty struct{}

type pointerPair struct {
	ptr unsafe.Pointer
	t   reflect.Type
}

func IsCyclic(v any) bool {
	val := reflect.ValueOf(v)
	seen := make(map[pointerPair]empty)

	return isCyclic(val, seen)
}

func isCyclic(v reflect.Value, seen map[pointerPair]empty) bool {
	if !v.IsValid() {
		return false
	}

	if v.CanAddr() {
		pp := pointerPair{unsafe.Pointer(v.UnsafeAddr()), v.Type()}

		if _, ok := seen[pp]; ok {
			return true
		}
		seen[pp] = empty{}
	}

	switch v.Kind() {
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			if isCyclic(v.Field(i), seen) {
				return true
			}
		}
		return false
	case reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			if isCyclic(v.Index(i), seen) {
				return true
			}
		}
		return false

	case reflect.Pointer, reflect.Interface:
		return isCyclic(v.Elem(), seen)
	}
	return true
}
