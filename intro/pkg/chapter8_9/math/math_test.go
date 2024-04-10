package math

import "testing"

type testpair struct {
	values  []float64
	average float64
}

var tests = []testpair{
	{[]float64{1, 2}, 1},
	{[]float64{1, 1, 1, 1, 1, 1, 1}, 1},
	{[]float64{-1, 1}, 0},
}

/*
func TestAverage(t *testing.T) {
	average := Average([]float64{1, 2})
	if average != 1.5 {
		t.Error("Expected 1.5, got ", average)
	}
} */

func TestAverage(test *testing.T) {
	for _, pair := range tests {
		actual := Average(pair.values)
		if actual != pair.average {
			test.Error("For:", pair, "expected:", pair.average, "got: ", actual)
		}
	}
}
