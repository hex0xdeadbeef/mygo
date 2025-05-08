package main

import (
	"fmt"
	"unsafe"
)

/*
	MEMORY LAYOUT

	Let's say we have 50 StructA instances. That'll eat up 1600 bytes (1.6KB) of memory. Not bad,
	right? But get this - those same 50 StructB instances would only use 800 bytes (0.8 KB)

	That difference can really start to impact our memory usage and performance in a hidden way.
*/

// StructA eats 32 bytes behind the scenes
type StructA struct {
	A byte  // 1-byte alignment
	B int32 // 4-byte alignment
	C byte  // 1-byte alignment
	D int64 // 8-byte alignment
	E byte  // 1-byte alignment
}

// StructB eats only 16 bytes behind the scenes
type StructB struct {
	B int32
	A byte
	C byte
	E byte
	D int64
}

/*
	WHAT IS ALIGNMENT?

	When experts mention alignmentm, they really mean the "required alignment" or "alignment guarantee" for a type.

	What does that mean? Well, say a type has an alignment value of N bytes. The compiler will enforce a rule that any variables of that types must live at a memory address that's a multiple of N.

	If it doesn't make sesne to you, let's look at the int32 type which has an alignment of 4 bytes. This means int32 variables must be stored at memory addresses that are multiples of 4. So valid addresses would be 0, 4, 8, 12 and so forth.

	Is that why our struct layouts seem so weird?
	Yep! Each type has its own alignment rules, so structs have to juggle those to keep their fields
	happy. That's why we see stuff like padding and extending alignment.

	But why do we even need alignment?

	Great question, modern CPUs are designed to access memory super efficiently when it's aligned to their word size (either 32 or 64 bits). Aligning data makes access predictable and smooth
	since it fits nicely into those present chunks.

	We may heard terms like "machine word", "native word", or just "word". These refer to the size of data a processor can handle in a single operation, based on if it's 32-bit or 64-bit architecture.

	To illustrate: say we're adding two 64-bit integers on a 32-bit processor. That processor would have to break it into multiple steps, handling the data in 32-bit chunks. It can't process the full 64-bits in one slot.


	HOW STRUCT ALIGNMENT LOOK LIKE?

	To understand how alignment works, nothing beats looking at real examples. Let's examine StructA from before example:

	Here we inspect StructA which has a largest alignment which is the D field with 8-bytes alignment, so this struct itself should have 8-byte alignment.

	The struct A will consume 8*4 = 32 bytes.

	From the initial section of this talk, we know that a byte has an alignment of 1 byte, int32 aligns to 4 bytes, and int64 aligns to 8 bytes, To explain the idea of StructA memory layout, it's easy to use alignment logic:

	- byte alignment begins at 0,1,2,3 ... (it can be anywhere)
	- int32 with a 4-byte alignment begins at addresses like 0, 4, 8, ...
	- int64 with an 8 byte alignment starts at addresses such as 0, 8, 16, ... In our example, int64 begins at address 16.

	Oh, I see the problem now, there are a lot of empty holes in the layout of StructA.
	We're right. Those gaps created by the alignment rules are called "padding".


	PADDING AND OPTIMIZE PADDING

	Padding's added to make sure each stuct field lines up properly in memory based on its needs, like we saw earlier.

	But while it enables efficient access, padding can also waste space if the fields ain't ordered well.

	The optimized guy takes 16 bytes total, with just 1 wasted byte of padding, way better now.

	We've looked at simple struct so far, but let's dig into more complex scenarios - structs containing other structs, pointers, arrays, maps and the like.


	HOW TYPES ALIGN

	Each type in Go has its own alignment, as we've seen before, and this depends on the architecture of the processor:

	- bool, uint8, int8: align to 1 byte.
	- int16, uint16: align to 2-byte boundaries
	- int32, uint32, float32: align to 4 bytes
	- int64, uint64, float64: align to 8 bytes
	- uint, int: these behave like uint32 and int32 on 32-bit architectures, and like uint64 and int64-bit architectures.
	- arrays: determined by the array's element type
	- structs: determined by the field with the largest alignment requirement
	- others: (string, channel, pointers, etc.) Align to size of machine word.

	What if we have a struct that contains an array with 5 elems of type uint32, a byte field, and a pointer to a byte like this:

	type StructA struct {
		A [5]int32
		B byte
		Bptr *byte
	}

	Let's inspect each field:
	- A: an int32 takes 4 bytes, so 5 int32s will take up 20 bytes (5*4), the array will take alignment of its type.
	- The byte field B takes 1 byte
	- The pointer BPtr takes up either 4 bytes (32-bit arch) or 8 bytes (64-bit arch) depending on the system arch (currently we are accounting for 64-bit)

	Given the fields, the total size of StructA will be either 28 bytes (32-bit) or 32 bytes (64-bit)

	But because struct field alignments work like the largest field, the entire struct will align to the 8-byte pointer.

	How about removing the pointer? Wouldn't the struct align to 4 bytes then?

	We got it, without that pointer, the array's alignment of 4 bytes for int32 become the largest alignment.


	DO I REALLY NEED TO CHANGE MY CODE JUST APPLY THIS OPTIMIZATION?

	No, we shouldn't overhaul working code just for this micro-optimization, unless we know memory issues are being caused by our struct specifically.
*/

func main() {
	ExampleA()
}

func ExampleA() {
	var s StructA

	// The compiler's making sure any int32 fields get slotted into address spots that work for
	// them.
	fmt.Println("alignof int32", unsafe.Alignof(s.B))      // 4
	fmt.Println("address offset 1:", unsafe.Offsetof(s.B)) // 4
	fmt.Println("address offset 2:", unsafe.Offsetof(s.E)) // 24
}

type OptimizedStructA struct {
	D int64 // 8-byte alignment
	B int32 // 4-byte alignment
	A byte  // 1-byte alignment
	C byte  // 1-byte alignment
	E byte  // 1-byte alignment
}
