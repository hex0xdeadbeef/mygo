package main

import (
	"fmt"
	"unsafe"
)

/*
	SLICES

	CHECKING MEMORY ADDRESSES OF SLICE
1. Using unsafe.SliceData(...)
	When we pass a slice to unsafe.SliceData(...), it does some checks to figure out what to return:
		- If the slice has a capacity greater than 0, the function returns a pointer to the first elemt of the slice.
		- If the slice is `nil`, the function just returns `nil`
		- If the slice isn't nil, but has zero capacity (an empty slice), the function gives us a pointer, but it's pointing to `Unspecified Memory Address`

	In Go, we can have types with zero size, like struct{} or [0]int. When the Go runtime allocates memory for these types, instead of giving each one a unique memory address, it just returns the
	address of a special variable called `zerobase`.

	FULL SLICE EXPRESSIONS
1. By default, if we don't specify the third parameter in the slicing op, the capacity is taken from the sliced slice or the length of the sliced array.

2. The capacity of a slice is the maximum number of elems it can hold before it needs to grow.

3. If we keep adding elems to a slice and it surpasses its current cap, Go will automatically create a larger array, copy the the elems and use that new array as the slice's underlying array.

4. This is also a common mistak new developers make when two slices that used to share the same data no longer do after an append() op. The most typical scenario is when a slice is passed as a
	function argument.

	THE CAPACITY OF SLICE
1. When Go creates a new array to adapt a growing slice, it usually doubles the capacity, but this changes once the slice reaches a certain size.

2. The table of slice growing: https://victoriametrics.com/blog/go-slice/index.html

3. But doubling capacity indefinitely would lead to huge memory allocations as the slice gets larger. To avoid that, Go adjusts the growth rate once thw slice reaches a certain size, typically
	around 256.

	At this point, the growth slows down, following this formula:
		`oldCap + (oldCap + 3 * 256) / 4`

	With some basic math, we can see that it equates to:
		`1.25 * oldCap + 192.`
	This keeps the slice growing efficiently without wasting too much
	memory. But keep in mind, this value is just a hint, an approximation.

4. The table above tells a different story depending on the type of slice, and that's because Go has to consider other factors like alignment, page size, size classes, and so on, all of which
	relate to how memory is allocated using certain predefined sizes.

5. If we need to increase the length of an existing slice in Go, we don't always have to use the append() function or create a new slice. We can simply extend the slice using a slicing op by
	specifying a length that's greater than the current length.

	This method works as long as the new length doesn't exceed the slice's capacity.

	HOW SLICE IS ALLOCATED?
1. We need to consider two things separately:
	1) The slice itself (which just the slice header)
	2) The underlying array
2. The first case: Both the slice and its underlying array are allocated on the stack.

	In this case the addresses are close to each other. It's easily observed by printing them.

	The slice is a simple struct with 3 fields, so it's usually allocated on the stack unless we do something that causes the slice itself to outlive the function.

	The underlying array, on other hand, is more likely to end up on the heap because it's created on the fly, either when we allocate a slice with a certaing number of elements or when the slice
	grows and needs more space.

3. Why is the underlying array in the example above allocated on the stack?
	The answer is: the size of the underlying array is known at compile time, so Go can optimize the allocation an place it on the stack.

4. The second case: The underlying array starts on the stack, then grows to the heap.

	In the second case the addres of the underlying array changes dramatically when the slice exceeds its capacity.

	At this point, it's no longer on our goroutine stack.

	This is why setting a predefined capacity is good idea to avoid unnecessary heap allocation, Even if we don't know the exact size at compile time, giving the slice an estimated capacity is
	better than leaving it at zero.

5. The third case: The underlying array is allocated on the heap.
	Even if we predefine a capacity that's known at compile time, there are still situations where the underlying array ends up on the heap.

	On these situations happens when using make() function, if the capacity exceeds 64 KB

	The underlying array of heapedSlice is allocated on the heap, but not the slice header itself.

6. Forcing the underlying array to be allocated on the heap.
	It's pretty easy to force the underlying array to be allocated on the heap, anything dynamic at compile time end up there.

	The stack size is determined at compile time, the underlying array of slice will definitely be on the heap as it size is n, which is determined at runtime.

7. What should we do to avoid heap allocations in these cases?
	- In most situations, it's enough to estimate the size of a slice at
		compile time, so we can't entirely avoid heap allocation the first
		time.
	- We can use make() with a capacity that's known at runtime to reduce
		the chances of additional heap allocations later on.
	- Using sync.Pool is a good option, as we can reuse the underlying array of the slice for the same purpose and expect that the task will have the same size. Before putting the slice back in
		the pool, set the length of the slice to 0. (slice = slice[:0]), so next time we can use append() naturally,
*/

func zerobaseIdentifying() {
	var emptyStruct struct{}

	fmt.Printf("zerobase: %p\n", &emptyStruct)

	var emptyArr [0]int
	fmt.Printf("zerobase: %p\n", &emptyArr)

	var emptiedSlice = make([]int, 0)
	fmt.Printf("zerobase: %p\n", emptiedSlice)

}

func sliceHeaderUsage() {
	arr := [6]byte{0, 1, 2, 3, 4, 5}

	slice := arr[1:3]
	fmt.Printf("slice: %v\n", slice)

	type sliceHeader struct {
		array unsafe.Pointer

		len int
		cap int
	}

	header := (*sliceHeader)(unsafe.Pointer(&slice))
	fmt.Printf("sliceHeader: %v %v %v\n", header.array, header.len, header.cap)
}

func fullSliceExpressions() {
	arr := [6]int{0, 1, 2, 3, 4, 5}
	slice := arr[1:3:4]

	extendedSlice := append(slice, 333)

	fmt.Println("Beforew appending:")
	fmt.Println(arr, "\n", slice, "\n", extendedSlice)

	slice = append(slice, 777)

	fmt.Println()
	fmt.Println("After appending into the `slice`:")
	fmt.Println(arr, "\n", slice, "\n", extendedSlice)

	extendedSlice = append(extendedSlice, 333)

	fmt.Println()
	fmt.Println("After appending into the `extendedSlice`:")
	fmt.Println(arr, "\n", slice, "\n", extendedSlice)
}

func tryToChangeSlice() {
	slice := []int{1, 2, 3}

	extendedSlice := make([]int, 3, 4)
	copy(extendedSlice, slice)

	changeSlice(extendedSlice)

	fmt.Println(*(*int)(unsafe.Pointer(uintptr(unsafe.Pointer(&extendedSlice[0])) + 3*unsafe.Sizeof(extendedSlice[0]))))

	fmt.Println(extendedSlice)
}

func changeSlice(s []int) {
	// Will be shown
	s[0] = -1

	// Will be stored in the same underlying array, but hidden from the caller's
	s = append(s, -333)

	// Won't be shown
	s = append(s, 333, 444, 555)
}

func lengthIncreasingTrick() {
	s := []int{1, 2, 3, 4, 5, 6}

	sPartCopied := s[1:3]
	sPartCopied = sPartCopied[:len(sPartCopied)+1]

	fmt.Println(sPartCopied)
}

func sliceAllocationA() {
	a := byte(1)
	fmt.Printf("a's address: %p\n", &a)

	// s and its underlying array don't escape to the heap.
	s := make([]byte, 1)
	fmt.Printf("slice's address: %p\n", &s)
	fmt.Printf("underlying's array's address: %p\n", s)

	/*
		// The addresses are close to each other
			slice's address: 0x140000b0000
			underlying's array's address: 0x140000a0005
	*/
}

func sliceAllocationB() {
	slice := make([]int, 0, 3)
	fmt.Printf("slice: %p - slice addr: %p\n", slice, &slice)

	slice = append(slice, 1, 2, 3)
	fmt.Printf("slice: %p - slice addr: %p\n", slice, &slice)

	slice = append(slice, 4)
	fmt.Printf("slice: %p - slice addr: %p\n", slice, &slice)
}

func sliceAllocationC() {
	stackedSlice := make([]byte, 64*1<<10)
	fmt.Printf("stackedSlice %p - stackedSlice array: %p\n", &stackedSlice, stackedSlice)

	heapedSlice := make([]byte, 64*1<<10+1)
	fmt.Printf("heapedSlice %p - heapedSlice array: %p\n", &heapedSlice, heapedSlice)
}

func heapAllocForcing(n int) {
	// The underlying array will be allocated on the heap because
	// n is runtime-defined value.
	slice := make([]int, n)

	_ = slice
}
