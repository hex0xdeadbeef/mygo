package arraystmts

import "fmt"

func SumOfArray(a *[3]float64) (sum float64) {
	for i := 0; i < len(a); i++ {
		sum += a[i]
	}

	return
}

func CLikeArrayUsing() {
	arrP := &[...]float64{1, 2, 3}
	fmt.Println(SumOfArray(arrP))
}
