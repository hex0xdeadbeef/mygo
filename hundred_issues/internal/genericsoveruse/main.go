package main

/*
	USEFUL CASES OF GENERICS:
1. Data Structures: We can use generics in order to emphasize the type of an elem. For example: during implementation a binary tree|linked list|heap
2. Functions working with slices/maps/channels of any type. For example: a function unioning two chans will work with any type -> we can use type params in order to define
a channel type.
3. Factorizing behavior instead of types. For example: "sort" package includes the interface sort.Interface including three methods: Len, Less, Swap
*/

/*
	FORBIDDEN CASES OF USING GENERICS
1. An invocation a method with type args
2. When generics obfuscate the code. A bunch of abstractions complicates understanding of a program in general.
*/

/*
	CONCLUSION
We won't write a complicated code polluting it with unnecessary abstractions, until there are weighty needs. When we met a situation forcing us to write a formulaic code, we
can start to think about using of generics.	
*/

import (
	"fmt"
	"io"
	"sort"
	"strconv"
	"time"
)

func getKeys(m any) ([]any, error) {
	switch m := m.(type) {
	case map[string]int:
		var (
			keys []any
		)

		for k := range m {
			keys = append(keys, k)
		}
	case map[int]int:
		var (
			keys []any
		)

		for k := range m {
			keys = append(keys, k)
		}

		// Another cases will be repeated during adding new types
	}

	return nil, fmt.Errorf("incorrect type provided")
}

func getKeysGen[K comparable, T any](m map[K]T) []K {
	var (
		keys []K
	)

	for k := range m {
		keys = append(keys, k)
	}

	return keys
}

type customConstraintA interface {
	~int | ~string
}

func getKeysGenConstrainted[K customConstraintA, T any](m map[K]T) []K {
	var (
		keys []K
	)

	for k := range m {
		keys = append(keys, k)
	}

	return keys
}

func Usage() {
	m := map[string]int{
		"one":   1,
		"two":   2,
		"three": 3,
	}

	keys := getKeysGenConstrainted[string](m)

	for _, k := range keys {
		fmt.Println(k)
	}
}

type Node[T any] struct {
	Val  T
	Next *Node[T]
}

func (n *Node[T]) Add(next *Node[T]) *Node[T] {
	var (
		start          = n
		cur   *Node[T] = start
	)

	if n == nil {
		n = next
		return n
	}

	for cur.Next != nil {
		cur = cur.Next
	}

	cur.Next = next

	return n
}

type customConstraintB interface {
	~int            // Any type that is based on int including the int itself
	String() string // Any type implementing String() method
}

type customInt int

func (ci customInt) String() string {
	return strconv.Itoa(int(ci))
}

func GetDoubled[T customConstraintB](num T) (int, string) {
	return int(num) * 2, num.String()
}

type Foo struct{}

// func (f Foo) Bar[T any](t T) {} // method must have no type parameterssyntax

/*
	USEFUL CASES OF USING GENERICS
*/
// merge merges two chans into a new one
func merge[T any](chA, chB <-chan T) <-chan T {
	res := make(chan T)

	go func() {
		ticker := time.NewTicker(time.Second * 1)
		defer ticker.Stop()

		for {
			select {
			case v, ok := <-chA:
				if !ok {
					chA = nil
				}
				res <- v

			case v, ok := <-chB:
				if !ok {
					chB = nil
				}
				res <- v

			case <-ticker.C:
				return
			}
		}
	}()

	return res
}

type SliceFn[T any] struct {
	S       []T
	Compare func(a, b T) bool
}

func (sf SliceFn[T]) Len() int           { return len(sf.S) }
func (sf SliceFn[T]) Swap(i, j int)      { sf.S[i], sf.S[j] = sf.S[j], sf.S[i] }
func (sf SliceFn[T]) Less(i, j int) bool { return sf.Compare(sf.S[i], sf.S[j]) }

func factorizingUsage() {
	s := SliceFn[int]{
		S:       []int{1, 2, 3, 4, 5},
		Compare: func(a, b int) bool { return a < b },
	}

	sort.Sort(s)
	fmt.Println(s.S)
}

/*
	FORBIDDEN CASES OF USING GENERICS
*/

// It should be changed to an explicit parameter of type io.Writer
// func foo[T io.Writer](w T) {
// 	b := getBytes()
// 	_, _ = w.Write(b)
// }

func foo(w io.Writer) {
	b := getBytes()
	_, _ = w.Write(b)
}

func getBytes() []byte {
	return []byte{1, 2, 3, 4, 5}
}
