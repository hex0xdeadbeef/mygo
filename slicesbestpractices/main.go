package main

import "fmt"

func main() {

	defaultValue()

	fmt.Println()

	makeCreation()

	fmt.Println()

	passingSliceToAFunc()

	fmt.Println()

	copyUsage()

	fmt.Println()
	appendResultUsing()
}

func defaultValue() {
	var (
		s []int
	)

	fmt.Println(s == nil)
	fmt.Println(len(s))

	s = []int{}
	fmt.Println(s == nil)
	fmt.Println(len(s))
}

func makeCreation() {
	var (
		double = func(s []int) []int {
			var (
				res = make([]int, 0, len(s))
			)

			for _, v := range s {
				res = append(res, 2*v)
			}

			return res
		}
	)
	list := []int{1, 2, 4, 8, 16}

	doubledList := double(list)

	fmt.Println(doubledList)
}

func passingSliceToAFunc() {
	var (
		handleA = func(s []int) {
			s[0] = 333
		}

		handleB = func(s []int) {
			s = append(s, 777)

			fmt.Println("append:", s)
		}

		list = []int{1, 2, 3, 4}
	)

	// handleA
	fmt.Println(list)
	handleA(list)
	fmt.Println(list)

	fmt.Println()
	// handleB.1
	fmt.Println("before append:", list)
	handleB(list)
	fmt.Println("after appendL", list)

	fmt.Println()
	// handleB.2
	fmt.Println("before append:", list)
	handleB(list[:1])
	fmt.Println("after appendL", list)

}

func copyUsage() {
	var (
		handle = func(s []int) []int {
			res := make([]int, len(s))

			copy(res, s)

			return res
		}

		list = []int{1, 2, 4, 8, 16}
	)

	newList := handle(list)
	newList[0] = 333

	fmt.Println(list)
	fmt.Println(newList)
}

func appendResultUsing() {
	var (
		list    = make([]int, 4, 5)
		newList []int
	)

	newList = append(list, 333)

	list[0] = 333
	newList[len(newList)-2] = 777

	fmt.Println(list)
	fmt.Println(newList)
}
