package main

import "fmt"

/*
	MAP
1. Properties of the map type:
	- A for-range loop over a map won't return the keys in any specific order.
	- Maps aren't thread-safe, the Go runtime will cause a fatal error if we try to read/write/iterate.
	- We can check if a key is in a map doing a simple `ok` check: _, ok := m[key]
	- The key type for a map must be comparable
2. Interfaces can be both comparable and non-comparable.

	We can define a map with an `empty` interface as the key without any compile errors, but the internal strucures put to interface-key should be comparable.

	So, we shouldn't use the interafes as map keys.

	MAP ANATOMY
1. The strucure of map:
	type hmap struct {
	...
	buckets unsafe.Pointer
	...

	A map contains a pointer that points to the bucket array. This is why when we assign a map to a variable or pass it to a function, both the variable and the function's argument are sharing the same pointer.
	}

2. Maps are pointers to the `hmap` under the hood, but they aren't reference types, nor are they passed by reference like a ref argument in C#, if we change the map parameter setting it to nil, it won't reflect in the caller's functiom.

3. In Go everything is passed by value.

4. The hash function used for maps in Go is consistent across all maps with the same key type, the `seed` used by that hash function is different for each map instance. So, when we create a new map, Go generates a random seed just for that map.

5. A `map` in Go grows when one of two conditions is met:
	1) There are too many overflow buckets
	2) Map is overloaded, meaning the load factor is too high

	Because of these 2 conditions, there are also two kind of growth:
		1) One that doubles the size of the buckets (when overloaded)
		2) One that keeps the same size but redistributes entries (when there are too many overflow buckets)

	If there are too many overflow buckets, it's better to redistribute the entries rather than just adding more memory.

	When the load factor goes beyond this threshold, the map is considered overloaded. In this case, it'll grow by allocating a new array of buckets that's twice the size of the current one, and then rehashing the elements into these two buckets.

	The reason we need to grow the map even when a bucket is just almost full comes down to performance. We usually think that accessing and assigning values in a map is O(1), but it's not always that simple.

	The more slost in a bucket that are occupied, the slower things get.

	When we want to add another key-value pair, it's not just about checking if a bucket has space. It's about comparing the key with existing key in that bucket to decide if we're adding a new entry or updating an existing one. And it gets even worse when we have overflow buckets, because then we have to check each slot in those overflow buckets too. This slowdown affect access and delete ops as well.

	Go caches the tophash of `key` in the bucket as `uint8`, and uses that for a quick comparison with the tophash of any new key. This makes the initial check superfast. After comparing the `tophash`, if they match, it means the keys "might" be the same. Then, Go moves on to the slower process whether the keys are actually identical.

	WHY DOES CREATING A NEW MAP WITH make(map, int) NOT PROVIDE AN EXACT SIZE BUT JUST A HING?
1. The `hint` parameter in make(map, hint) tells Go the initial number of elements we expect the map to hold.

	This hint helps minimize the number of times the map needs to grow as we add elements.

	Since each growth operation involves allocating a new array of buckets and copying existing elements over, it's not the most efficient process. Starting with a larger initial capacity can help avoid some of these costly growth ops. 
*/

func interfaceMapKey() {
	m := map[any]int{
		1:   1,
		"a": 2,
	}

	// panic: runtime error: hash of unhashable type []int
	// panic: runtime error: hash of unhashable type func() {}
	// m[[]int{1, 2, 3}] = 3
	// m[func() {}] = 4

	_ = m
}

func changingMapInFunc() {
	var m = map[string]string{
		"a": "b",
		"c": "d",
	}

	fmt.Println(m)

	changeMap(m)
	fmt.Println("Changed")

	fmt.Println(m)
}

func changeMap(m map[string]string) {
	m["a"] = "a"
}

func nillingMapInFunc() {
	var m = map[string]string{
		"a": "b",
	}

	fmt.Println(m)
	nonNilledFunc(m)
	fmt.Println(m)
}

func nonNilledFunc(m map[string]string) {
	m = nil
}

func differentMapRandomSeeds() {
	var (
		mA = map[int]string{
			1: "a",
			2: "b",
			3: "c",
		}

		mB = map[int]string{
			1: "a",
			2: "b",
			3: "c",
		}
	)

	for k, v := range mA {
		fmt.Println(k, v)
	}

	for k, v := range mB {
		fmt.Println(k, v)
	}
}
