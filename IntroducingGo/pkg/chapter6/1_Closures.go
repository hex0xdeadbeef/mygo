package chapter6

import (
	"math/rand"
)

func firstFunctionWithClosure(arguments ...int) float64 {
	averageValue := func(closureMultipleArguments ...int) float64 {
		var total float64 = 0.0
		for _, element := range closureMultipleArguments {
			total += float64(element)
		}
		return total / float64(len(closureMultipleArguments))
	}
	return averageValue(arguments...)

}

func generateRandomSlice(amountOfNecessaryElements int) []int {

	var randomElements []int = make([]int, 0, amountOfNecessaryElements)

	for i := 0; i < amountOfNecessaryElements; i++ {
		randomInt := rand.Intn(101)
		randomElements = append(randomElements, randomInt)
	}

	return randomElements
}

func MakeEvenGenerator() func() uint {
	i := uint(0)               // На эту переменную будет ссылаться переменная хранящая внешнюю функцию
	return func() (ret uint) { // Эта функция имеет право изменять локальный внешние переменные
		ret = i
		i += 2
		return ret
	}
}
