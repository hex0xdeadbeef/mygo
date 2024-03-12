package runesstrings

import "fmt"

func Using() {
	aString := "Hello World ££¢∞•£¡£º"

	fmt.Printf("First character: %c\n", aString[0])

	for i, r := range aString {
		fmt.Printf("Byte start point: %10d | Character %10c | Unciode point: %10c\n	", i, r, r)
	}

}
