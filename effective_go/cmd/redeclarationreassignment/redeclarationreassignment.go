package redeclarationreassignment

import "fmt"

func ShortRepetitions() {
	var (
		v int
	)

	{
		v := 10

		fmt.Println(v)
	}

	fmt.Println(v)

	fmt.Println()

	v, err := a()
	fmt.Println(v, err)

}

func a() (int, error) {
	return 10, nil
}
