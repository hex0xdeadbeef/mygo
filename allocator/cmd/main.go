package main

import (
	"fmt"
)

/*
	GO ALLOCATOR
1. Pointers can be arranged in stack too.
2. Golang always tries to allocate memory on stack except some situations.
3. The mechanism used to define a place (Stack/Heap) where a variable will reside is called "Escape Analysis"
	1) "Escape Analysis" works during compile time
4. The variable will be relocated to heap by "Escape Analysis" if:
	1) The result is returned by using pointer
	2) Argument is put into interface{} parameter
	3) The size of variable exceeds stack's size
5. The algorith of "Escape Analysis":
	1) The weights graph is built
	2) Graph traversal to find variables we need to put into heap
	3) Put escape-flag for each variable that will be put into stack
6. Stages of Batch(...) (...) that is used during escape analysis:
	1) Graph building
	2) Variables size check (each type has its own limit)
	3) Graph traversal
	4) Variables markering
7. Heap's idea: memory consists of some layers or different blocks of memory for each thread
8. Go allocator
	1) The memory is given by blocks. This block is called "arena". When arena is full, a new "arena" is given.
	2) The size of "arena" is bound to the specific OS. The arm's size of "arena" is 4MB
	3) Arena is splitted to 8KB pages.
	4) The "arena"s and their's pages are managed by hmap structure (heap arena)
	5) Arena is also splitted into "Objects pool". The pool of objects is called in Go's sources as "mspan"
	6) Golang has 67 classes of objects
	7) When an object is put into arena: its class is defined based on its size
	8) When the pool is full the new pool is created and bound to the previous as next (doubled linked list)
	9)
	type mspan struct {
		startAddr uintptr 	says where the mspan resides in memory
		npages  			says how many pages are used to support this pool
		freeIndex 			points to a free block of memory in this pool to allocate objects faster
		nelems 				says how many elems are arranged in this pool
		allocBits			shows the free cells in this pool and is needed to be used in search and for GC using
	}
	10) A list of mspans is managed by "mcentral"

	type mcentral struct {
	...
	spanclass
	nonempty 	 	mspan <-> mspan <-> mspan <-> mspan <-> mspan <-> mspan ...
	empty 			mspan <-> mspan <-> mspan <-> mspan <-> mspan <-> mspan ...
	}

	We have 67 "mcentral"s
	11) Arena -> "mcentral"s -> list of "mspan"s
	12) Each thread has its own Cache that is represented as "mcache"
*/

func main() {}

// On the Stack
/*
 1. func A()
    x = 2
    y = 0 // we don't know what will be returned from function
 2. func pow()
    a = 2
 3. func A()
    x = 2
    y = 4
 4. fmt.Println(y)
    a = 4
*/
func A() {
	x := 2
	y := pow(x)

	fmt.Println(y) // the variable y is put into interface{} so it'll be relocated to heap
}

func pow(a int) int {
	return a * a
}

// On the Stack
/*
1. func B()
	x = 2 at the address 0xA
2. func inc()
	x = 0xA at the address 0xB
3. func B()
	x = 3 at the address 0xA
4. fmt.Println(x)
	a = 3
*/
func B() {
	x := 2
	inc(&x)

	fmt.Println(x)
}

func inc(x *int) {
	*x++
}

// On the heap
/*
	1. func C()
		x = nil at the address 0xA
	2. func getVal() Stack pointer UP
		x = 4 at the address 0xB
	3. func C() Stack pointer DOWN
		x = 0xB at the address 0xA
	4. func fmt.Println(*x/2) Stack pointer UP and shades the stack of function C where x = 4 at the address 0xB resides.
We got "Dagling pointer"
*/
func C() {
	x := getVal()
	fmt.Println(*x / 2)
}

func getVal() *int {
	x := 4 // x escapes to heap (the pointer to it is returned)
	return &x
}

// Escape Analysis graph building
type T struct {
	Value int
}

/*
Simple assignment 	= 0
Reference 			= -1
Dereference 		= 1
*/
func tryToEscape(a T) (*T, bool) {
	l1 := a      // The simple assignment (0)
	l1.Value = 1 // The simple assignment (0)

	l2 := &l1 // The pointer assignment (-1)
	l3 := *l2 // Dereferencing (1)

	return &l3, a.Value == 1
}
