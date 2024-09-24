package main

import (
	"fmt"
	"math"
	"math/rand"
	"runtime"
)

/*
	MEMORY LEAKS
1. Working with slices and maps we should remember about the situatuins when we want to prevent leaks of memory. We shouldn't just cut a slice or delete elements from a map.
2. If we're talking about slices, then after just shrinking we get the same underlying slice but the visible part is restricted to the index we've shrinked before.
3. If we're talking about maps, then after deleting elems, the buckets are just rearranged but still present in the internal structure. The way we need follow is to copy the elems into
	a brand new map using the "len" parameter and and to put all the elems from the old map to this one.
*/

const (
	sliceMaxSize     = 1 << 24
	idxToBeCutBefore = 1 << 8
)

func ShrinkSliceWrong() error {
	var (
		s   = getSlice(sliceMaxSize)
		err error
	)

	printMemAllocations()
	if s, err = cutSliceWrong(s, idxToBeCutBefore); err != nil {
		return err
	}

	printMemAllocations()
	runtime.GC()
	printMemAllocations()

	runtime.KeepAlive(&s)
	return nil
}

func ShrinkSliceRight() error {
	var (
		s   = getSlice(sliceMaxSize)
		err error
	)

	printMemAllocations()
	if s, err = cutSliceRight(s, idxToBeCutBefore); err != nil {
		return err
	}

	printMemAllocations()
	runtime.GC()
	printMemAllocations()

	runtime.KeepAlive(&s)
	return nil
}

func printMemAllocations() {
	const (
		ToKB = 1 << 10
	)
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	fmt.Printf("%d KB\n", m.Alloc/ToKB)
}

func cutSliceWrong(s []int, idxToBeCutTo int) ([]int, error) {
	const (
		idxPad = 1
	)

	if idxToBeCutTo >= len(s) {
		return nil, fmt.Errorf("index out of the range: %d ", idxToBeCutTo)
	}

	return s[:idxToBeCutTo+idxToBeCutTo], nil
}

func cutSliceRight(s []int, idxToBeCutBefore int) ([]int, error) {
	const (
		idxPad = 1
	)

	if idxToBeCutBefore >= len(s) {
		return nil, fmt.Errorf("index out of the range: %d ", idxToBeCutBefore)
	}

	var buf = make([]int, idxToBeCutBefore+idxToBeCutBefore)
	copy(buf, s[:idxToBeCutBefore+idxToBeCutBefore])

	return buf, nil
}

func getSlice(size int) []int {
	var (
		s = make([]int, size)
	)

	for i := 0; i < size; i++ {
		s[i] = rand.Intn(math.MaxUint32)
	}

	return s
}
