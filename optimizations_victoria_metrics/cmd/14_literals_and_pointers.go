package main

import (
	"fmt"
)

/*
	LITERALS AND POINTERS.


	PASSING VALUES IS FAST AND RARELY SLOWER THAN PASSING POINTERS
1. Copying a small amount of data is very efficient and can often be quicker than the indirection required when using pointers.

2. It reduces the workload of the GC, when values are passed directly, the GC has fewer pointer references to track.

3. Data that is passed by values tends to be stored closer together in memory, allowung the CPU to access the data more quickly.

4. Sure, pointers can speed things up for large or growing (unbounded) struct, bu we've got to prove it's worth it.


	POINTER RECEIVERS AND VALUE RECEIVERS NUANCES.
1. Both (addressable) Values T and Pointers *T can invoke methods with Value and Pointer Receivers.
	This is the most basic principle:
		- If a type has a value receiver method, it can be called on both the value itself and a pointer to that value.
		- Similarly, for a pointer receiver, it can be called on both its pointer and the value itself.

2. So, there are two notes here
	1) Go automatically takes the value's address to call a pointer receiver.
	2) It automatically dereferences a pointer to call a value receiver.


	NUANCES WITH `nil`
1. No example panics in the function NilCarUsage() function, because we don't use any field of a structure.

2. Attempting to call a method with a value receiver on a nil pointer (var car *Car; car.Move()) leads to panic, because Go can't dereference nil to pass it as a value (at runtime)


	THE UNADDRESSABLE PROBLEM
1. Go automatically takes a value's address to call a pointer receiver, but what happens when the value is `unaddressable`? Like values returned from a function or temporary values (a struct literal or value from a map) are typically not addressable directly.

2. So what's the outcome when we try calling a method that expects a pointer receiver on such values?

	Neither m[0] nor GetDefaultCar() can call the Transform method, we face a compile error.

	This is because Go needs to extract an address to call the method with a pointer receiver, but it can't do this with temporary, unaddressable values.

3. Of course, calling a method with a value receiver on these instances works just fine, as expected.


	INTERFACES AND METHOD SET NUANCES
1. What are method sets?
	For a type T, its method set includes all methods with T as the receiver.

2. But for a pointer to a type (*T), the method set expands, it includes methods with both T and *T as the receiver.

	Both *AnotherCar and AnotherCar can have the Move() method in their method sets, allowing them to satisfy the Mover interface.

3. Even though T and *T might seem to interchangeably use each others methods, their method sets diverge when interfacing with ..., well, interfaces.


	PREFER USING A POINTER RECEIVER WHEN DEFINING METHODS
1. The rule isn't black and white: "Use pointer receiver to modify, and value receivers otherwise".

2. When to opt for a pointer receiver?
	-  To modify the receiver's state.
	- For structs that are considered "large", which can be a bit subjective as I've mentioned in previous section.
	- When a struct includes synchronizing fields like sync.Mutex, sync.WaitGroup, sync.Once, sync.Cond, sync.Pool, sync.RWMutex, sync.Map, opting for a pointer helps avoid copying the lock.
	- If unsure, learning towards a pointer receiver might be wise.

3. When is a value receiver suitable?
	- For types that are small and won't change (immutab;e)
	- If our type is a map, func, channel, or involves slices that aren't altered in size or capacity (though elements may change)

	Why slices that aren't altered in size or capacity?
		Copying the slice an modifying its size will break the connection.

		While we can modify the elements of the slice or the content of the underlying array through a value receiver (affecting the original slice), resizing the slice (e.g., using append that increases the capacity) will not affect the original slice outside the method.

4. Consistency is key.
	Avoid mixing receiver types for a given struct to maintain consistency.

	If any method requires a pointer receiver due to mutations, it's generally better to use pointer receivers for all methods of that struct, even those that don't cause mutations.

	The main reasons are:
		- Using both can lead to confusion and inconsistency in how instance of that struct are manipulated, especially when changing the receiver type.
		- To keep object interactions with interfaces consistent and straightforward (Remember the `method set` method we discussed in Take 3)
*/

type Car struct {
	Model string
}

func (c Car) Move() {
	fmt.Println(c.Model, "car moves.")
}

func (c *Car) Transform(model string) {
	c.Model = model
}

func CarUsage() {
	m := Car{Model: "Tesla"}
	m.Move()
	m.Transform("Toyota")

	mP := &Car{"Tesla"}
	mP.Move()
	mP.Transform("Toyota")
}

type NewCar struct {
	Model string
}

func (nc *NewCar) Transform() {
	fmt.Println("Transforming...")
}

func (nc NewCar) Move() {
	fmt.Println(nc.Model, "moving...")
}

func GetDefaultCar() NewCar {
	return NewCar{Model: "Tesla"}
}

func NilCarUsage() {
	var nc *NewCar

	nc.Transform() // doesn't panic

	(*NewCar)(nil).Transform() // doesn't panic

	nc.Move() // panic: runtime error: invalid memory address or nil pointer dereference

	(*NewCar)(nil).Move() // panic: runtime error: invalid memory address or nil pointer dereference

}

func UnaddressableUsage() {
	m := map[int]NewCar{}
	_ = m

	// m[0].Transform() // cannot call pointer method Transform on NewCarcompilerInvalidMethodExpr

	// GetDefaultCar().Transform() // cannot call pointer method Transform on NewCarcompilerInvalidMethodExpr

}

type Mover interface {
	Move()
}

type AnotherCar struct {
	Model string
}

func (ac AnotherCar) Move() {
	fmt.Println("Moving...")
}

func MoverUsageA() {
	var m Mover = AnotherCar{Model: "Tesla"}
	m.Move()

	var m2 Mover = &Car{Model: "Tesla"}
	m2.Move()
}

type AnotherNewCar struct {
	Model string
}

func (anc *AnotherNewCar) Move() {
	fmt.Println("Moving...")
}

func NoverUsageB() {
	var m Mover = &AnotherNewCar{Model: "Tesla"}
	m.Move()

	// var m2 Mover = AnotherNewCar{Model: "Tesla"} // cannot use AnotherNewCar{…} (value of type AnotherNewCar) as Mover value in variable declaration: AnotherNewCar does not implement Mover (method Move has pointer receiver)compilerInvalidIfaceAssign

}
