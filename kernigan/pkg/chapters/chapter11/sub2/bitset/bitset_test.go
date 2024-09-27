package bitset

import "testing"

var addTests = []struct {
	name     string
	elems    []int
	expected string
}{
	{"Basic", []int{1, 2, 3}, "{ 1 2 3 }"},
	{"Same elements", []int{7, 7, 7}, "{ 7 }"},
	{"Empty list", []int{}, "{ }"},
	{"Large set", []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, "{ 1 2 3 4 5 6 7 8 9 10 }"},
	{"Zero value", []int{0}, "{ 0 }"},
	{"Repeated elements in different order", []int{3, 2, 1, 1, 2, 3}, "{ 1 2 3 }"},
}

func TestAdd(t *testing.T) {
	const (
		templ = "%s: Add(%#v) = %s, want %s"
	)

	var (
		set *IntSet = new(IntSet)
	)

	for _, testSuit := range addTests {

		// ops
		set.Add(testSuit.elems...)

		// check
		if got := set.String(); got != testSuit.expected {
			t.Errorf(templ, testSuit.name, testSuit.elems, got, testSuit.expected)
		}
		set.Clean()
	}
}

var removeTests = []struct {
	name string

	elems          []int
	removedElement int
	expected       string
}{
	{"Basic", []int{1, 2, 3}, 1, "{ 2 3 }"},
	{"Same elements", []int{7, 7, 7}, 7, "{ }"},
	{"Empty list", []int{}, 100, "{ }"},
}

func TestRemove(t *testing.T) {
	const (
		templ = "%s: Remove(%#v) = %s, want %s"
	)

	var (
		set *IntSet = new(IntSet)
	)

	for _, testSuit := range removeTests {

		// ops
		set.Add(testSuit.elems...)
		set.Remove(testSuit.removedElement)

		// check
		if got := set.String(); got != testSuit.expected {
			t.Errorf(templ, testSuit.name, testSuit.elems, got, testSuit.expected)
		}
		set.Clean()
	}
}

var elementsTests = []struct {
	name string

	elems    []int
	expected []int
}{
	{"Basic", []int{1, 2, 3}, []int{1, 2, 3}},
	{"Same elements", []int{7, 7, 7}, []int{7}},
	{"Empty list", []int{}, []int{}},
}

func areElementsSimilar(expected, got []int) bool {
	if len(expected) != len(got) {
		return false
	}

	for i := 0; i < len(expected); i++ {
		if expected[i] != got[i] {
			return false
		}
	}

	return true
}

func TestElements(t *testing.T) {
	const (
		templ = "%s: Elements(%#v) = %#v, want %#v"
	)

	var (
		set *IntSet = new(IntSet)
	)

	for _, testSuit := range elementsTests {

		// ops
		set.Add(testSuit.elems...)

		// check
		if got := set.Elements(); !areElementsSimilar(testSuit.expected, got) {
			t.Errorf(templ, testSuit.name, testSuit.elems, got, testSuit.expected)
		}
		set.Clean()
	}
}

var lenTests = []struct {
	name string

	elems    []int
	expected int
}{
	{"Basic", []int{1, 2, 3}, 3},
	{"Same elements", []int{7, 7, 7}, 1},
	{"Empty list", []int{}, 0},
}

func TestLen(t *testing.T) {
	const (
		templ = "%s: Elements(%#v) = %d, want %d"
	)

	var (
		set *IntSet = new(IntSet)
	)

	for _, testSuit := range lenTests {

		// ops
		set.Add(testSuit.elems...)

		// check
		if got := set.Len(); got != testSuit.expected {
			t.Errorf(templ, testSuit.name, testSuit.elems, got, testSuit.expected)
		}
		set.Clean()
	}
}

type setsPair struct {
	src *IntSet
	dst *IntSet
}

type copyTest struct {
	name  string
	elems []uint

	initFunc func(ct *copyTest)
	pair     *setsPair

	expected bool
}

type CopyTests []*copyTest

func (ct *CopyTests) init() {
	for _, test := range *ct {
		test.initFunc(test)
	}
}

func setCopiedPair(ct *copyTest) {
	if ct.elems == nil {
		ct.pair.src, ct.pair.dst = nil, nil
		return
	}
	ct.pair.src = &IntSet{words: ct.elems}
	ct.pair.dst = &IntSet{words: make([]uint, len(ct.pair.src.words))}

	copy(ct.pair.dst.words, ct.pair.src.words)
}

func isSetCopied(source, copy *IntSet) bool {
	if source == nil && copy == nil {
		return true
	}

	return &source.words != &copy.words
}

var copyTests = CopyTests{
	{name: "Usual", elems: []uint{1, 2, 3}, initFunc: setCopiedPair, pair: new(setsPair), expected: true},
	{name: "nil", elems: nil, initFunc: setCopiedPair, pair: new(setsPair), expected: true},
}

func TestCopy(t *testing.T) {
	const (
		templ = "%s: Copy(%#v)"
	)
	copyTests.init()

	for _, test := range copyTests {
		// ops
		if isSetCopied(test.pair.src, test.pair.src.Copy()) != test.expected {
			t.Errorf(templ, test.name, test.elems)
		}
	}
}

var unionTests = []struct {
	name           string
	elems1, elems2 []int
	expected       string
}{
	{"Usual", []int{1, 2}, []int{3, 4}, "{ 1 2 3 4 }"},
	{"Either empty", []int{}, []int{3, 4}, "{ 3 4 }"},
	{"Both empty", []int{}, []int{}, "{ }"},
	{"Same elems", []int{1, 2, 3}, []int{1, 2, 3}, "{ 1 2 3 }"},
}

func TestUnion(t *testing.T) {
	const (
		templ = "%s: UnionWith(%#v, %#v) = %s, want %s"
	)

	var (
		set1, set2 = new(IntSet), new(IntSet)
	)

	for _, test := range unionTests {
		set1.Add(test.elems1...)
		set2.Add(test.elems2...)

		// ops
		set1.UnionWith(set2)
		if set1.String() != test.expected {
			t.Errorf(templ, test.name, test.elems1, test.elems2, set1.String(), test.expected)
		}

		set1.Clean()
		set2.Clean()
	}
}

var intersectionTests = []struct {
	name           string
	elems1, elems2 []int
	expected       string
}{
	{"No mutual elems", []int{1, 2}, []int{3, 4}, "{ }"},
	{"Single mutual elems", []int{1, 2}, []int{2, 3}, "{ 2 }"},
	{"Plenty mutual elems", []int{1, 2, 3, 4}, []int{1, 2, 3, 5}, "{ 1 2 3 }"},
	{"Either empty", []int{}, []int{1, 2, 3}, "{ }"},
	{"Both empty", []int{}, []int{}, "{ }"},
}

func TestIntersection(t *testing.T) {
	const (
		templ = "%s: IntersectWith(%#v, %#v) = %s, want %s"
	)

	var (
		set1, set2 = new(IntSet), new(IntSet)
	)

	for _, test := range intersectionTests {
		set1.Add(test.elems1...)
		set2.Add(test.elems2...)

		// ops
		set1.IntersectWith(set2)
		if set1.String() != test.expected {
			t.Errorf(templ, test.name, test.elems1, test.elems2, set1.String(), test.expected)
		}

		set1.Clean()
		set2.Clean()
	}
}

var differenceTests = []struct {
	name           string
	elems1, elems2 []int
	expected       string
}{
	{"No mutual elems", []int{1, 2}, []int{3, 4}, "{ 1 2 }"},
	{"Single mutual elems", []int{1, 2}, []int{2, 3}, "{ 1 }"},
	{"Plenty mutual elems", []int{1, 2, 3, 4}, []int{1, 2, 3, 5}, "{ 4 }"},
	{"Left empty", []int{}, []int{1, 2, 3}, "{ }"},
	{"Right empty", []int{1, 2, 3}, []int{}, "{ 1 2 3 }"},
	{"Either empty", []int{}, []int{1, 2, 3}, "{ }"},
	{"Both empty", []int{}, []int{}, "{ }"},
}

func TestDifference(t *testing.T) {
	const (
		templ = "%s: IntersectWith(%#v, %#v) = %s, want %s"
	)

	var (
		set1, set2 = new(IntSet), new(IntSet)
	)

	for _, test := range differenceTests {
		set1.Add(test.elems1...)
		set2.Add(test.elems2...)

		// ops
		set1.DifferenceWith(set2)
		if set1.String() != test.expected {
			t.Errorf(templ, test.name, test.elems1, test.elems2, set1.String(), test.expected)
		}

		set1.Clean()
		set2.Clean()
	}
}

var symmetricDifferenceTests = []struct {
	name           string
	elems1, elems2 []int
	expected       string
}{
	{"No mutual elems", []int{1, 2}, []int{3, 4}, "{ 1 2 3 4 }"},
	{"Single mutual elems", []int{1, 2}, []int{2, 3}, "{ 1 3 }"},
	{"Plenty mutual elems", []int{1, 2, 3, 4}, []int{1, 2, 3, 5}, "{ 4 5 }"},
	{"Left empty", []int{}, []int{1, 2, 3}, "{ 1 2 3 }"},
	{"Right empty", []int{1, 2, 3}, []int{}, "{ 1 2 3 }"},
	{"Both empty", []int{}, []int{}, "{ }"},
	{"All elemets the same", []int{1, 2, 3}, []int{1, 2, 3}, "{ }"},
}

func TestSymmetricDifference(t *testing.T) {
	const (
		templ = "%s: SymmetricDifference(%#v, %#v) = %s, want %s"
	)

	var (
		set1, set2 = new(IntSet), new(IntSet)
	)

	for _, test := range symmetricDifferenceTests {
		set1.Add(test.elems1...)
		set2.Add(test.elems2...)

		// ops

		if got := set1.SymmetricDifference(set2); got.String() != test.expected {
			t.Errorf(templ, test.name, test.elems1, test.elems2, got.String(), test.expected)
		}

		set1.Clean()
		set2.Clean()
	}
}
