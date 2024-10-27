package main

import (
	"fmt"
	"unsafe"
)

/*
	ARRAYS

	What's an array?
1. Arrays in Go are a lot like those in other programming languages. They've got fixed size and store elems of the same type in contiguous memory locations.

	This means Go can acces each element quickly since their addresses are calculated based on the starting address of the array and the element's index.

2. There are a couple of things to notice in the function ArrayA(...):
	1) The address of the array `arr` is the same as the address of the first element.
	2) The address of each element is 1 byte apart from each other because our element type is `byte`.

	Our stack is growing downwards from a higher to a lower address, right?. This pictute shows exactly how an arrays looks in the stack, from arr[4] to arr[0].

	So, does that mean we can access any element of an array by knowing the address of the first element (of the array) and the size of the element? Let's try this with an int array and `unsafe`
	package.

3. In the function ArrayB(...):
	We get the pointer to the first element and then calculate the pointers to the next elements by adding multiples of the size of an `int`, which is 8 bytes on a 64-bit architecture. Then we use these pointers to access and convert them back to the int values.

4. The example is just a play around with the `unsafe` package to access memory directly for educational purposes. Don't do this in production without understanding the consequences.

5. Now, an array of type T isn't a type by itself, byt an array with a specific size and type T, is considered a type. Even though both `a` and `b` are arrays of bytes, the Go compiler sees them
	as completely different types, the %T format makes this point clear.

6. Here's how the Go compiler sees it internally (src/cmd/compile/internal/types2/array.go)

	// An array represents an array type
	type Array struct {
		length int64
		eleme Type
	}

	// NewArray returns a new array type for the given element type and length. A negative length indicates an unknown length.
	func NewArray(elem Type, len int64) *Array {
		return &Array{len: len, elem: elem}
	}

	The length of the array is `encoded` in the type itself, so the compiler knows the length of the array from its type. Trying to assign an array of one size to another, or compare them, will
	result in a mismatched type error.

	ARRAY LITERALS
1. There are many ways to initialize an array in Go, and some of them might be rarely used in real projects. Check the function ArrayE(...)

	What we're doing in the function ArrayE(...) (except for the first one) is both defining and initializin their values, which is called a "composite literal". This term is also used for slices,
	maps, and structs.

2. Now here's an interesting thing: when we create an array with less than 4 elems, Go generates instructions to put the values into the array one by one.

	So when we do arr := [4]int{1,2,3,4}, what's actually happening is:
		arr := [4]int{}
		arr[0] = 1
		arr[1] = 2
		arr[2] = 3
		arr[3] = 4

	This strategy is called local-code initialization. This means that the initialization code is generated and executed within the scope of a specific function, rather than being part of the
	global or static initialization code.

	The compiler creates a static representation of the array in the binary which is known as `static initialization` strategy. This means the values of the array are stored in a read-only section
	of the binary. This static data is created at compile time, so the values are directly embedded into the binary. [5]int{1,2,3,4,5} looks like in Go assembly:
		main..stmp_1 SRODATA static size=40
			0x0000 01 00 00 00 00 00 00 00 02 00 00 00 00 00 00 00  ................
			0x0010 03 00 00 00 00 00 00 00 04 00 00 00 00 00 00 00  ................
			0x0020 05 00 00 00 00 00 00 00                          ........

	Our data is stored in stmp_1, which is read-only static data with a size of 40 bytes (8 bytes for each element), and the address of this data is harcoded in the binary.

	The compiler generates code to reference this static data. When our application runs, it can directly use this pre-initialized data without needing additional code to set up the array.
3. What about an array with 5 elems but only 3 of the initialized?
	Good question, this literal [5]int{1,2,3} falls into the first category, where Go puts the value into the array one by one.

	While talking about defining and initializing arrays, we should mention that not every array is allocated on the stack. If it's too big, it gets moved to the heap.

	As of Go 1.23, if the size of the variable, not just array, excedds a constant valye MaxStackVarSize, which is currently 10 MB, it'll be considered too large for stack allocation an will
	escape to the heap.

	ARRAY OPERATIONS
1. The length of the array is hardcoded in the type itself. Even though array don't have a cap property, we can still get it.

	The capacity equals the length, no doubt, but the most important thing is that we know this at compile time.

	So len(arr) doesn't make sense to the compiler it's not a runtime property, Go compiler knows the value at compile time. Go takes this key point and turns it into a constant behind the scenes.

	SLICING OVER AN ARRAY
1. Slicing is a way to get a slice from an array, and its full form is denoted by the syntax [start:end:capacity]. Usually, we'll see its variants:
	- [start:end]
	- [:end]
	- [:]

2. `start` is the index of the first element to include in the new slice (inclusively), `end` is the index of the last elemet to exclude from the new slice (exclusively), and `capacity` is an
	optional argument that specifies the capacity of the new slice.

3. The compiler evaluates the slicing indices (start, end, and capacity) to determine the bound of the new slice.

4. If any of the indices are missing, they default to:
	- `start` defaults to 0
	- `end` default to the length of original slice or array
	- `capacity` default to the capacity of the original slice or the length of the original array.

5. What about the new length and capacity of the slice?
	- The new length is determined by subtracting the start index from the end index
	- The new capacity is determined by subtracting the start index from either the capacity argument (if provided) or the original capacity

	When we write b := a[1:3], here's what's really going on:
		b.len = 3 - 1
		b.capacity = 5 - 1
6. Regarding the panic when we specify the end index out of bounds: because the end is exclusive, we can specify it to be equal to the length of the original array.

	We can also specify `start` at the length index like arr[len(arr):]. What we're creating is an empty slice, with no length and no capacity. So, the general rule for bound-checking the slicing
	is:
		`0 <= start <= end <= cap <= real capacity`

	ARRAY ARE VALUES, RAW VALUES
1. In some other languages, an array variable is basically a pointer to the first element on the array. When we pass an array to a functon, what's actually passed is a pointer, not the whole
	array. So, changing the array elements within the function will affect the original array.

	In contrast, Go treats arrays as value types. This means an array variable in Go represents the entire array, not just a reference to its first element, even though printing &a gives the same
	address as &a[0].

	When we pass an array to a function in Go, the entire array is copied.

	The output of ArrayK(...) is expected since we're modifying the copied array, not the original one.

2. Isn't it inefficient? It's always copied.
	It's only inefficient only when we pass a large array to a function and benchmarks or profilers show it's a bottleneck. Otherwise, we can just pass arrays as usual.

3. Here's a tricky thing when we loop over our array, especially using for-range loop.

	In the snippet of the func ArrayL(...) `a` and `b` are [3]int so we can assign them, but we assign `a = b` right in the loop. When we iterate to index 1, we immediately change the array, so
	the output should be `1 2 6` because v already evaluated to 2. Unexpectedly, the output is `1 2 3`, just like nothing happened. So, did a change or did our assignment have no effect?
		- The arr used inside the loop is a copy of the original arr.

	Go indeed makes a copy, but he copy is hidden from us, and only `v` can see that copied array. The array `a` we use in the loop is still our original a, and if we print it out after the loop,
	it'll be `[4, 5, 6]`, we could think of another scenario like the one shown in the picture

	That means it works like a pass-by-value case. If our array is much bigger than just several elements, making a copy like this will be inefficient and the Go team has optimized this for us by
	allowing for-range with a pointer to the array. In this case the output will be `1 2 6`. But the key takeaway here isn't encourage changing the array inside the loop, I'd not recommend that.
	Instead, it's to show that Go supports for-range with a pointer to an array, while it doesn't support pointers to slices.
*/

func ArrayA() {
	arr := [5]byte{0, 1, 2, 3, 4}
	fmt.Println("arr:", arr)

	for i := range arr {
		fmt.Println(i, &arr[i])
	}

}

func ArrayB() {
	arr := [3]int{99, 100, 101}

	// Unsafe pointered pointer to the start of the array
	curElem := unsafe.Pointer(&arr[0])

	// The size of a stride to stride through the arr
	elemSizeStride := unsafe.Sizeof(arr[0])

	for i := 0; i < len(arr); i++ {
		fmt.Printf("%d: %d\n", i, *(*int)(curElem))

		// curElem is the unsafe pointer referencing to the next element
		// in the array.
		curElem = unsafe.Pointer(uintptr(curElem) + elemSizeStride)
	}
}

func ArrayC() {
	arrA := [4]byte{}
	arrB := [5]byte{}

	fmt.Printf("%T\n", arrA)
	fmt.Printf("%T\n", arrB)

	// cannot convert arrB (variable of type [5]byte) to type [4]bytecompilerInvalidConversion
	// fmt.Println(arrA == [4]byte(arrB))
}

func ArrayD() {
	type EmptyStruct struct{}

	esA, esB := EmptyStruct{}, EmptyStruct{}
	fmt.Println(esA == esB)

	type EmptyArray [0]int

	eaA, eaB := EmptyArray{}, EmptyArray{}
	fmt.Println(eaA == eaB)
}

func ArrayE() {
	var arr1 [10]int // [0, 0, 0, 0, 0, 0, 0, 0, 0, 0]

	// With value, infer-length
	arr2 := [...]int{1, 2, 3, 4, 5} // [1, 2, 3, 4, 5]

	// With index, infer length
	arr3 := [...]int{11: 3} // [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 3]

	// Combined index and value
	arr4 := [5]int{1, 4: 5} // [0, 0, 0, 0, 5]

	// Combined index and value
	arr5 := [5]int{2: 3, 4, 4: 5} // [0, 0, 3, 4, 5]

	_, _, _, _, _ = arr1, arr2, arr3, arr4, arr5
}

func ArrayF() {
	// 10 MB -> stack
	a := [10 * 1 << 20]byte{}

	_ = a

	// 10 MB + 1 byte -> heap
	b := [10*1<<20 + 1]byte{}

	_ = b
}

func ArrayG() {
	var arr = [5]byte{1, 2, 3, 4, 5}
	fmt.Println(len(arr)) // 5
	fmt.Println(cap(arr)) // 5
}

func AraryH() {
	var arr = [5]byte{1, 2, 3, 4, 5}

	fmt.Println(len(arr)) // The Go's compiler just references to the previous row behind the scenes
	fmt.Println(cap(arr)) // The Go's compiler just references to the previous row behind the scenes
}

func ArrayI() {
	arr := [5]int{0, 1, 2, 3, 4}

	// new slice from a[1] to arr[3-1]
	b := arr[1:3]
	_ = b

	// new slice from a[0] to arr[3-1]
	c := arr[:3]
	_ = c

	// new slice from a[1] to arr5-1]
	d := arr[1:]
	_ = d
}

func ArrayJ() {
	arr := [5]int{0, 1, 2, 3, 4}

	// len(arr) will not be considered
	b := arr[:5]
	fmt.Println(b)

	c := arr[5:]
	fmt.Println(c)
}

func ArrayK() {
	arr := [5]int{0, 1, 2, 3, 4}
	fmt.Println(arr)

	tryToMutateArr(arr)
	fmt.Println(arr)

}

func tryToMutateArr(arr [5]int) {
	arr[0] = -1
	_ = arr
}

// 1 2 3
func ArrayL() {
	a := [3]int{1, 2, 3}
	b := [3]int{4, 5, 6}

	for i, v := range a {
		if i == 1 {
			a = b
		}

		fmt.Print(v, " ")
	}
}

// 1 2 6
func ArrayM() {
	a := [3]int{1, 2, 3}
	b := [3]int{4, 5, 6}

	for i, v := range &a {
		if i == 1 {
			a = b
		}

		fmt.Print(v, " ")
	}
}
