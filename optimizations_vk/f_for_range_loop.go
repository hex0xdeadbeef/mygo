package main

/*
	FOR RANGE/LOOP RANGING
1. The lack of the for-range cycle is that we superficially copy all the element while ranging over the slice or other structures.

2. It's better to use for-loop cycles to suppress the values copying.
*/

type CustomIndex struct {
	idx int
}

func SumRange(cs []CustomIndex) int {
	var (
		res int
	)

	for _, v := range cs {
		res += v.idx
	}

	return res
}

func SumLoop(cs []CustomIndex) int {
	var (
		res int
	)

	for i := 0; i < len(cs); i++ {
		res += cs[i].idx
	}

	return res
}
