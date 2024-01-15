package a_for_cycles

import "testing"

var functions []func() = []func(){
	Echo_SimpeFor,
	Echo_RangeFor,
	Echo_stingsJoin,
	Echo_StraightforwardPrint,
}

func BenchmarkEchos(b *testing.B) {
	for _, function := range functions {
		for i := 0; i < b.N; i++ {
			function()
		}
	}

}
