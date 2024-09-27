package main

import (
	"fmt"
	"unique"
	"unsafe"
)

/*
	UNIQUE PACKAGE
1. Interning is the process of storing only one copy of a value in memory and sharing a unique reference to it instead of allocating multiple copies and wasting memory. For instance: the Go
	compiler already performs interning of string constants at compile time, so instead of allocating multiple identical strings it'll allocate a single instance of the string and share
	reference to it.
2. Interning (or 'canonicalizing values') can now be performed using "unique" package via its "Handle" type that acts as a globally unique identity for any provided (comparable) value, meaning
	two values used to create the handles would have also compared equal.
		type Handle[T comparable] struct{}

		func (h Handle[T]) Value() T {}

		func Make[T comparable](value T) Handle[T] {}
3. Internally the "Handle" type is backed by a concurrent-safe map (well, an adaptive radix tree, to be precise due to its need to shrink and grow the tree more efficiently than a map) which
	acts as a read-through cache, storing unique values when not detected in the cache and returning a Handle[T] which is designed to be passed around as values/fields instead of the underlying
	value, saving you additional allocations and also resulting in cheaper comparisons.
4. Sure, the same interning behaviour could be achieved using a custom map to reduce duplicate allocations, however handling garbage collecting efficiently isn't possible. The "unique" package
	has a notion of "weak references" that means the garbage collector can clean up in a single cycle made possible by access to runtime internals, something a custom rolled solution wouldn't
	be able to perform.
5. The package "unique" uses weak pointers to work with interned values. It allows the GC to clean the objects using a single cycle instead of multiple of them if there are no strong links to these objects. At the moment when these strong links left, a weak link is not to obstruct the GC.

	TYPES OF LINKS TO OBJECTS
	STRONG LINKS
	Usual links - the standart mechanism of referencing to objects in the RAM. When an object exists, there is at least one strong link to it, the GC cannot remove this object since the program
	is still using it. More strong links to objects, more time an object exists in the RAM. But if the program is about not to use this object, but there are some links (for example: in a map
	or in a slice), this object won't be removed and it can lead to memory leaks.

	WEAK LINKS
1. In the terms of the default RAM managing model an object exists and is available until it's referenced by any links. These links are called "strong". A "weak" link is the link that doesn't
	hold the object in the RAM. If the object is referenced only by weak links, the GC can remove this object from the memory. After the deletion a weak link gets "nil" value or points to an
	invalid object.
2. Weak link - is the type of link, that doesn't obstruct the GC to remove an object. If this object is referenced only by weak links, the GC considers it as an unused and can free up the
	memory used by this object.
3. The main difference between strong links and weak links is that a weak link doesn't hold an object in memory. The object is removed by the GC when it's not referenced by any strong links.
	After the deletion this object all the weak links will get "nil" value, that is they don't reference to the object.
4. Caching data
	When we work with the systems that cache data, weak links allow to save objects, that
	can be removed automatically from the memory. These deletions will happen when there is a lack of resources, without the need of explicit lifecycle managing of them.
5. The behavior of weak links
	1. Creation of weak link.
		We create a weak link to an object that can be collected by the GC.
	2. The acces to the object.
		We can try to access the object using a weak link. If the object isn't collected,
		we'll get the access to it. If the object collected, we'll get the nil value or other indicator of object abscence.
	3. Object deletion
		When all the strong links left, the GC can remove the object, even if it's referenced by the weak links.
	4. Postprocessing after removing the object
		After removing the object the weak links become invalidated.

*/

func main() {
	// stringsInterningAtCompileTime()

	usingUniqueA()
}

func stringsInterningAtCompileTime() {
	const (
		greetA = "hello world"
		greetB = "hello world"
	)

	var (
		a = greetA
	)

	// greetA == greetB && greetB == a -> greetA == a
	fmt.Println(unsafe.StringData(greetA) == unsafe.StringData(greetB) && unsafe.StringData(greetB) == unsafe.StringData(a))
	fmt.Printf("%d\n%d\n%d\n", unsafe.StringData(greetA), unsafe.StringData(greetB), unsafe.StringData(a))

	fmt.Println("\na changed")
	a += "s"

	fmt.Printf("%d\n%d\n%d\n", unsafe.StringData(greetA), unsafe.StringData(greetB), unsafe.StringData(a))
}

type addrDetail struct {
	isV6   bool
	zoneV6 string
}

func usingUniqueA() {
	h1 := unique.Make(addrDetail{isV6: true, zoneV6: "2001"})
	// this addrDetail won't be allocated as it already exists in the underlying map (adaptive radix tree)
	h2 := unique.Make(addrDetail{isV6: true, zoneV6: "2001"})

	if h1 == h2 {
		fmt.Println("Addresses are equal")
	}

	// Value() returns a copy of the underlying value (i.e, different memory address)
	fmt.Println(h1.Value().isV6)
}
