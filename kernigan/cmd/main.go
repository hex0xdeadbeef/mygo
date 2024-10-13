package main

import "math/rand/v2"

func main() {
	// workWithProdA()
	// workWithProdB()
}

func getProducer() <-chan struct{} {
	res := make(chan struct{})

	go func() {
		var (
			randomSendingTimeInSecs = rand.IntN(10)

			time.Sleep()
		)
	}()


	return res


}
