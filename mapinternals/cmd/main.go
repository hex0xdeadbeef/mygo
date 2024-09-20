package main

import (
	"fmt"
	"math/rand"
	"unsafe"
)

/*
	MAP INTERNALS
1. Applying hash-function to two different keys we can get the same index. This situation called "collision". They are solved by using two methods:
	1) Open addressing method
	2) Chains method
2. "Open addressing" method.
	Using "open addressing" method we put the key into a free cell that is:
		1) The next of the current
		2) With the specific striding through an array (Unit stride, Squared stride etc.)
	The search inside this array is performed in the same way. The element will be found using 1) or 2).
3. "Chains" method.
	Using "Chains" method we add to the element in the current cell a pointer to the new one. Hence, there's a connection between the old elements and the new ones.
		1) There might be a case when we put all the N elements into the same cell and the doubled linked list appears. The search in this structure takes O(N) time compelxity.
4. To meake search faster we can modify our table. Assume that we write our goods into the table and we distribute the elems using a particular property. For example it can be
fruits or vegetables. We've splitted our table to some parts called "buckets". For example there are 10 elems in fruits and 10 elems in the bucket of vegetables. So now we:
	1) define the property of an object
	2) loop over the bucket's keys and find the needed one
5. To control filling of the table we use the following term: "load factor". Using the "load factor" we can resctructure our table making distribution steady.
6. "load factor" - ratio of the amount of "key-value" pairs to the amount of the buckets. If they beat the specific value, the hash-table is being extended and indexes are
distributed.
7. The structure of the map.
	type hmap struct {
		// ... more code
		count      int 	 	 an amount of "key-val" pairs in the map
		B 		   uint8 	 an amount of the buckets expressed by log2(#). This expression extremely eases bit ops.
		noverflow  uint16    an approximate amount of overflowed buckets
		hash0	   uint32	 a number used to affect happenstance and overcome collisions
		buckets	   unsafe.Pointer a pointer to an array of buckets. In Go it's a cell of memory, where the buckets are stored.
		oldbuckets unsafe.Pointer a pointer to an array of buckets from which data is fetched during "data evacuation".
	}

	1) count - the amount of key-value pairs
	2) B - the amount of buckets expressed as log2. This expression extremely eases bit operations. The explanation is:
		log2(CountBuckets) = B
		2^B = CountBuckets
		log2(1) = 0
		log2(2) = 1
		log2(4) = 2
		log2(8) = 3
		...
	3) noverflowed - the approximate amount of overflowed buckets
	4) hash0 - a number used to affeect happenstance and overcome collisions. The collisions persistence is the key point in elems distribution because we could've gotten a si-
	tuation when all the elems had been put into the same buckets, so the value hash0 is relevant in the hmap structure.
	5) buckets - a pointer to an array of buckets. In Go it is a a cell where the buckets are stored.
	6) oldbuckets - a pointer to an array of buckets from which data is fetched during "data evacuation"
8. "data evacuation" - the situation when a map is overflowed and to make the pairs' distribution steadier data is relocated to other buckets. It happens during store/delete
operations. "data evacuation" happens when the load factor gets the 6.5 value, it means that a ratio of elems (key-value pairs) count to buckets count is 6.5.
8. There's another structure: MapType
	type MapType struct {
		Hasher func(unsafe.Pointer, uintptr) uintptr is a function that takes the arguments (key, hash0) and returns hash value
	}
9. There's another one structure: bmap
	bucketCnt = 1 << 3 - the size of each bucket in a map

	type bmap struct {
		// tophash generally contains the top byte of the hash value for each key in this bucket
		tophash [bucketCnt]uint8
	}
10. Insertion into map. Check the book.
11. Affter insertion of the 8th element (we've filled all the 272 bytes of bucket) we'll touch the overflow block. The "overflow" block is arranged at the end of the bucket.
The descrtibed structure eases the search in the bucket and we can easily store/read/delete elements.
12. Bucket's overflow. What does happen when overflow happen? There are the two approaches:
	1) Increasing the overall amount of buckets and pairs rearranging among these buckets.
	2) New bucket creation that is bound to the old one. If we need to wtire out an element into a bucket, we'll put the element to the new bucket.
13. Increasing of the amount of buckets and pairs distribution is called "data evacuation". The decision of need of "data evacuation" is based on "load factor"
	1) The value 6.5 (or 80% of a bucket) of "load factor" results in "data evacuation"'s start.
14. During "data evacuation":
	1) the amount of buckets is doubled
	2) the current pairs are rearranged throughout these new buckets
15. Buckets are stored sequentially in memory.
16 "overflow"'s 8 bytes take the rest of the bucket.
*/

func deferA() (res int) {
	res = 1
	defer func() {
		res++
	}()

	return res
}

func deferB() int {
	var res int

	res = 1
	defer func() {
		res++
	}()

	return res
}

func mapTrick() {
	for range 100 {
		var (
			m = map[int]bool{10: true}
		)

		for k, _ := range m {
			if k%10 == 0 {
				m[k+10] = true
			}
		}

		fmt.Println(m)

	}
}

func main() {
	fmt.Println(deferA())
	fmt.Println(deferB())
	mapTrick()
}

type hmap struct {
	/*
		MORE CODE HERE...
	*/

	// An amount of elements (key-val pairs)
	count int
	// An amount of buckets expressed by log2(#). This value extremely eases bit ops
	B uint8
	// An approximate value of overflowed buckets
	nooverflow uint16

	// The value that affect happenstance and overcome collisions. It's used in Hasher(unsafe.Pointer, uintptr) uintptr that returns a hash value.
	hash0 uint32

	// A pointer to an array of buckets. In Go this is a cell of a memory
	buckets unsafe.Pointer
	// A pointer to an array of buckets from which data is fetched during "data evacuation"
	oldbuckets unsafe.Pointer
}

type MapType struct {
	/*
		MORE CODE HERE
	*/

	// Hasher is a function that uses key and hash0 and returns hash value.
	Hasher func(unsafe.Pointer, uintptr) uintptr

	// for example: hash := t.Hasher(noescape(unsafe.Pointer(&s)), uintptr(h.hash0))
	// s - key
	// h.hash0 - hash0's current value
}

const (
	// This is the size of each bucket
	bucketCnt = 1 << 3
)

type bmap struct {
	// tophash generally contains the top byte of the hash value for each value in this bucket
	tophash [bucketCnt]uint8
}

func MapDefining() {
	var (
		_ map[string]string // this is the nil value

		/*
			In this case we initialize the hmap structure too.

			During this type of map creation the function makemap(...) is used.
		*/
		_ = make(map[string]string) // this is the nil value
	)

}

type maptype struct{}

// makemap creates map when we create a map using make(...) function
func makemap(t *maptype, hint int, h *hmap) *hmap {
	/*
		MORE CODE HERE...
	*/

	// 1. Initialize Hmap.
	/*
		We check the presence of h *hmap and if it's not present, we initialize it. Then we get hash0 value.
	*/
	if h == nil {
		h = new(hmap)
	}
	h.hash0 = fastrand()

	// 2.
	/*
		We calculate an amount of buckets B in the format of: log2 based on hint.
		hint - the second argument in the function make(map[DataTypeA]DataTypeB, hint)
	*/
	B := uint8(0)
	for overLoadFactor(hint, B) {
		B++
	}

	// 3.
	/*
		If the amount of buckets isn't equal to 1 (log2(1) == 0 == B), the memory is allocated for our buckets keeping.
	*/
	if h.B != 0 {
		// var nextOverflom *bmap
		// h.buckets, nextOverflom = makeBucketArray(t, h.B, nil)
		/*
			MORE CODE HERE...
		*/
	}

	/*
		MORE CODE HERE...
	*/

	return nil
}

func fastrand() uint32 {
	return rand.Uint32()
}

func overLoadFactor(hint int, B uint8) bool {
	return false
}
