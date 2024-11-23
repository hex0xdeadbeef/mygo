package main

import (
	"fmt"
	"unique"
)

/*
	UNIQUE PACKAGE
1. The Go team rolled out this "new" unique package. It's all about dealing with duplicates in a "smart" way.

2. The idea is pretty simple.
	When we've got several identical values in our code, we only store one copy. Instead of having several copies of the same thing, they all just point to this one version,
	which is a lot more efficient. It's a process often called "interning" in programming circles.

3. The big win is memory savings. We're not wasting space on multiple copies of the same value floating around.

4. We've got just one true, official version, and everything else points back to it. And the beauty of the "unique" package is that it manages all of this for us behind the
	scenes. We don't really have to think about it's too much and Go takes care of the heavy lifting.

5. Let's reference back to the String Interning technique.
	When we pass a string "s" to the intern function, it checks if that string is already hanging out in internStringMap. If it is, we skip storing a duplicate and just return
	the version that's already there. If it's not in the map, we go ahead and store it, and next time the same string pops up, we just reuse the stored version.

6. Strings in Go behave a lot like slices and when we assign a string to a new variable, say from
	"a" to "b", Go doesn't copy the actual value, like the string "VictoriaMetrics", but instead
	just copies what's called the "string header", similar to a slice header.

	So both "a" and "b" end up pointing to the same underlying byte array.

	type stringStruct struct {
		str unsafe.Pointer // []byte
		len int
	}

	Go enforces that strings are immutable (at least in normal situations), so even though they share the same memory, neither can change the underlying value. The key point here
	is, how do we take advantage of this shared underlying byte array to save memory when we've got identical string values?

	After interning we'll get the string that has the same string header and it's the same to having a unique string that doesn't reside in the RAM repeatedly.

	But the simple solution of interning isn't thread-safe because the "map" in Go isn't safe in situations of writing multiple values into it concurrently.

	Because of thread-unsafety of the simple solution we've moved to the side of the "sync.Map" type. Alright, it solves the problem, which is great.

	But we've still got another issue as the more unique vals we add, the bigger that map gets, and eventually, we're staring down the barrel of unbounded memory growth. Without
	some kind of cleanup mechanism, this thing will just keep growing forever, which isn't ideal.

	We could this with a simple time-to-live (TTL) mechanism for each key or the entire key-value pair. This increases complexity, but we can reference how VM handles it.

7. The Go's unique package takes a different strategy, similar to sync.Pool. There's another perk, or a benefit to string interning that Michael Knyszek pointed out:
	- faster equality checks

	Normally, comparing two strings works like this:
		1) First, we check if they start at the same memory address (basically comparing pointers) and whether their lengths match.
		2) If neither of those is true, we're stuck comparing the strings byte by byte.

	Sure, Go might optimize it so it's not literally checking byte by byte (chunk by chunk), but still, comparing two long strings is expensive if they're identical since we're
	cheking all the way to the last byte.

	If we've interned all of our strings, meaning there's only one canonical copy of each distinct string in memory, we can just compare the memory addresses (the pointers). If
	the pointers are the same, the strings are definitely equal.

	No need to go through the bytes, no need to worry about the lengths - just a quick pointer comparison, and we're done.

8. Unique package.
	The unique package only exposes two main public pieces:
		1) The type Handle[T]
		2) The function Make[T]()

	Here's what it looks like:
		// Handle is a globally unique identity for some value of type T
		type Handle[T comparable] struct {
			value *T
		}

		// Make returns a globally unique handle for a value of type T. Handles are equal if and only if the values used to produce them are equal.
		func Make[T comparable](value T) Handle[T] {...}

	The Make function is designed to give us a canonical, globally unique handle for any comparable value we pass in.

9. Instead of just returning the value itself, it hands us a Handle[T], which works as a kind of reference to that value. This handle lets us efficiently compare values and
	manage memory without worrying about having multiple copies of the same thing floating around in memory.

	In other words, when we call Make with the same value over and over, we'll always get the saame Handle[T]

	So, if two handles are equal, we know for sure that their underlying values are also equal. Comparing handles is way more efficient than comparing the values themselves,
	especially if those values are large or complex. Instead of doing a deep dive into the actual content of the value, we just compare the handles, since they're really just
	pointers, it's fast.

10. In the output of the function UniqueUsageA() we'll notice that h1 and h2 both point to the same memory address (0x14000090270), while w1 has a different address. The equality
	check between h1 and h2 (h1==h2) is based on comparing these pointers, not the actual content of the strings. This makes the comparison sidesteps the need to check each
	character of the strings themselves.

11. And what makes unique even more powerful is, it's not limited to just strings, we can intern any comparable type.

	We can indeed see how this is used in practice by taking a look at the source code of the net/http package, the Go team decided to use the unique package to intern the IPV6
	zone name within a struct called addrDetail:
		package netip

		type addrDetail struct {
			isV6 bool // IPv4 is false, IPv6 is true
			zoneV6 string // != "" only if IsV6 is true
		}

		var (
			z0 unique.Handle[addrDetail]
			z4 = unique.Make[addrDetail{}]
			z6noz = unique.Make(addrDetail{isV6: true})
		)

		func (ip Addr) WithZone(zone string) Addr {
			if !ip.Is6() {
				return ip
			}

			if zone == "" {
				ip.z = z6noz
				return ip
			}

			ip.z = unique.Make(addrDetail{isV6: true, zoneV6: zone})

			return ip
		}

	If we have an IPv6 address with a zone, it gets stored in a unique handle.

	Before switching to the unique package, this code is also used a custom interning solution from their internal package, internal/intern, which was itself a port of the
	interning logic.

11. Another hidden gem about interning strings is, if we intern a substring of a large string, even though both the substring and the original string share the same underlying
	array, the large string can stil be garbage collected if it's no longer in use anywhere else in the program.

12. What is "Weak Pointer"?
	Normally, when we've got a pointer to something, the garbage collector (GC) treats that object as "in use", meaning it won't be freed.

	A weak pointer (or weak reference) works differently, as it's a reference to a memory object that doesn't stop the GC from freeing that memory. But, we can convert it back
	into a regular (strong) pointer if we need to hold onto it.

	In other words, with a weak pointer, the memory can still be collected by the GC if no strong references (normal pointers) are holding onto it.

13. The real use case for weak pointers.
	In the Go proposal by Michael Knyszek, weak pointers are designed for cases where you need a reference to an object but don't want to control `how long that object lives`. If
	later, we decide, "Actually, I do need this to stick around", we can convert that weak pointer into a `strong` pointer.

14. The beauty of weak pointers.
	The beauty of weak pointers is that they're "non-intrusive", they don't manage the memory themselves.

15. How do we know when the memory is no longer valid?
	When the GC decides it's time to free up the memory for an object, it just updates the weak pointer to nil.
		- The catch is that weak pointers can become nil without we realizing it.

	So every time we want to use a weak pointer we first have to check whether it's still valid (i.e., not nil). If it's nil, the memory is gone and we're out of luck.

	Interestingly, a weak pointer doesn't actually point directly to the object in memory. Instead, it points to what's called an "indirection object," which is like a middleman
	of type atomic.Uintptr (correct me). This indirection object is just a small 8-byte block of memory that contains a memory address to the actual memory being referenced.

	The reason for this indirection layer is efficiency.

	When the GC frees the memory for an object, it only needs to update this one pointer inside the indirection object. This way, all weak pointers referencing that memory get set
	to nil instantly, without any extra work.

	It speeds up the whole process since the GC doesn't have to through and manually clear each weak pointer individually.
*/

func UniqueUsageA() {
	h1 := unique.Make("Hello")
	h2 := unique.Make("Hello")
	w1 := unique.Make("World")

	fmt.Println("h1:", h1)
	fmt.Println("h2:", h2)
	fmt.Println("w1", w1)

	fmt.Println("h1 == h2", h1 == h2)
	fmt.Println("h1 == w1", h1 == w1)
}
