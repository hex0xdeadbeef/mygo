package main

/*
	THE LACK OF INTERFACES
1. Consider an interface haning only two methods of summarizing objects. The first method accepts
	values by the value and the second one by the pointer.
2. As a result we'll see that the method that accepts an object as a pointer works slower and
	allocates something in the heap. The reason consists in the "escaping" because it's guided by the following rule:
		- If a variable is accessed by more than a goroutine, the variable should be
		passed to the heap. If the access to the variable is given to more than one goroutine, the goroutine must be moved to the heap or if it's not possible to prove that the only one
		goroutine can access this variable.
3. The rule of forcing the runtime to move the variable to the heap:
	Any casting the complex type to the interface - everything allocation except the following
	cases:
		1) Pointers to interfaces
		2) int, int32, int64, uint, uint32, uint64
		3) Simple types: empty struct{} or the interface that consists in the field that doesn't cross the border of only one machine word.
4. The caveat of using interfaces:
	We should avoid using any interfaces in our code.
*/

type Summer interface {
	SumValue(Object) int
	SumPointer(*Object) int
}

type Object struct {
	A, B int
}

type Sum struct{}

func (v Sum) SumValue(obj Object) int {
	return obj.A + obj.B
}

func (v Sum) SumPointer(obj *Object) int {
	return obj.A + obj.B
}

func (v Sum) SumUsualValue(obj Object) int {
	return obj.A + obj.B
}

func (v Sum) SumUsualPointer(obj *Object) int {
	return obj.A + obj.B
}
