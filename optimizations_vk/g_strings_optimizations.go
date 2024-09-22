package main

/*
	STRING OPTIMIZATIONS
1. The convertations: []byte <-> string and []rune <-> string in a for-range cycle doesn't lead to any allocations.

2. There are no allocations, when we compare byte slices casting them to strings from []byte.

3. There's no allocation when we make a convertation []bytes -> string to find an element in the map.

4. Relatively to https://github.com/golang/go/issues/61725#issuecomment-1788751754, sort.Strings() works 30% better than slices.Sort. It works in this way because sort.Strings makes the only one comparison instead of using 3. 
*/

func StringConvertationInForRange() {
	var (
		s = "abcdefg"
	)

	for _, b := range []byte(s) {
		_ = b
	}
}

func BytesComparison(a, b []byte) bool {
	return string(a) == string(b)
}

func IsElemInMap(m map[string]struct{}, key []byte) bool {
	_, ok := m[string(key)]
	return ok
}
