package deepequivalence

import (
	"fmt"
	"reflect"
)

func DeepSliceMap() {
	var a, aEmpty []int = []int{}, nil
	fmt.Println(reflect.DeepEqual(a, aEmpty))

	var b, bEmpty map[int]int = make(map[int]int), nil
	fmt.Println(reflect.DeepEqual(b, bEmpty))
}
