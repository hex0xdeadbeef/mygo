package main

import "fmt"

/*
THE WAYS WHEN INTERFACES SHOULD BE USED
1. Common behavior
	1) Sorting some elements can be torn apart into 3 activities: getting elems count, instruction about order of two elems, swap of two elems
2. Coupling descending
3. Behavior restriction
*/

func main() {

}

/*
3. Behavior restriction
*/

type intConfigGetter interface {
	Get() int
}

type Foo struct {
	threshold intConfigGetter
}

func NewFoo(threshold intConfigGetter) Foo {
	return Foo{threshold: threshold}
}

func (f Foo) Bar() {
	threshold := f.threshold.Get()
	fmt.Println(threshold)
}
