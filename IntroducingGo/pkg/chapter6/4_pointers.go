package chapter6

import (
	"fmt"
	"math/rand"
)

func firstPointersPresence() {
	number := 10                                     // Ссылается на number напрямую в рамках этого скоупа
	randomChangingParameterValueWithPointer(&number) // Даем адрес переменной, а не ее значение
	fmt.Println(number)
}

func randomChangingParameterValueWithPointer(numberPtr *int) { // Передаем указатель на переменную
	*numberPtr = rand.Intn(101) // изменяем значение переменной по указателю, также ссылается на number в том скоупе
}

func creatingAndUsingPointerWithNewOperator() {
	var number *int = new(int)
	randomChangingParameterValueWithPointer(number)
	fmt.Println(*number)
}
