package main

import "sync"

/*
	STRING INTERNING
1. Strings give a lot of flexibility, as they can represent anything in metadata and introduce new labels and values whenever they like. In practice, however, metadata strings don't change often, creating a lot of repetition during collecttion,

2. Take, for example, the go_info metric. Its metadata had a Go version label-value pair. There are only so many potential Go versions, and it's unlikely that the version of Go being used
	changes very often. But each time we collect this metric from our apps, we need to parse its metadata and allocate it in the memory until it's garbage-coolected. Taking into account this
	metric could be exposed by not one, but by thousands of apps, the metrics collector will have to parse and allocate in memory the same strings over and over.

3. To avoid storing the same strings lots of times, what if we stored each unique string once and referred to it when we needed to? This is called string interning (https://en.wikipedia.org/
	wiki/String_interning) and it can save a significant chunk of memory.

4. In the illustration 1_string_interning_usage.png, vmagent sees the same metadata string across multiple scrapes. In the left image, three copies of the same string are stored in memory,
	but in the right image only a single copy of the string is stored in memory. This allows saving memory by 3x.

5. The StringInterningNaive(...) wotks perfectly for a single-threaded app, but vmagent has many threads that work across many targets concurrently. It's possible to add a lock to the
	StringInterningNaive function, but that doesn't scale nicely on multi-core systems as there's likely to be a lot of contention when accessing this map.

6. The solution of beating lock-contention is to use sync.Map, a thread-safe implementation built into the Go standart library.

	The best part is that sync.Map simplifies our original code! It comes with a LoadOrStore method that means that we no longer need to check whether the string is already present in the
	map ourselves.

	sync.Map is optimized for two use cases:
		1) When a given key is only ever written once, but used many times, i.e. the cache has high hit ratio.
		2) When multiple goroutines read, write and overwrite entries for disjoint sets of keys, i.e., each goroutines uses a different set of keys.

	Whenever either of these two cases applies, sync.Map reduses lock-contention and improves the performance of our app compared to if we'd used a regular Go map paired with a Mutex or
	RWMutex.

7. There are a couple of "gotchas" that we should be aware of when using sync.Map.
	1) Unconstrained memory growth is dangerous. To prevent the map from growing infinitely, we need to either occasionaly rotate the map by deleting it and recreating it, or implement some
		form of time to live logic for our keys.

	2) Sanity check the arguments to intern
		The "intern" function performs really well when passing regular string, but unless we lock down the interface, someone will eventually try to pass something that might break our
		function.

		For example, in Go byte slices are mutable, this means they could change at any time and aren't suitable for use as a key in our map. This is usually a case of using "unsafe
		conversion to string", a common optimization in Go. And, in the same time, the most common source of bugs.
*/

func StringInterningNaive(m map[string]string, strToFind string) string {
	if v, ok := m[strToFind]; ok {
		return v
	}

	m[strToFind] = strToFind
	return strToFind
}

func StringInterningMultiThreaded(m *sync.Map, strToFind string) string {
	if v, ok := m.LoadOrStore(strToFind, strToFind); ok {
		return v.(string)
	}

	return strToFind
}
