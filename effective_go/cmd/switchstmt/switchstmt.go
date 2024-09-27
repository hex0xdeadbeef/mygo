package switchstmt

import "fmt"

func Unhex(c byte) byte {
	switch {
	case '0' <= c && c <= '9':
		return c - '0'
	case 'a' <= c && c <= 'f':
		return c - 'a' + 10
	case 'A' <= c && c <= 'F':
		return c - 'A' + 10
	}

	return 0
}

func FallthroughImitation(c byte) bool {
	switch c {
	case ' ', '?', '&', '=', '#', '+', '%':
		return true

	}

	return false
}

func BreakingLabel() {
Loop:
	for i := 10; i > 0; i++ {

		if i == 5 {
			break Loop
		}
	}

}

func TypeSwitching(v interface{}) {
	switch v := v.(type) {
	default:
		fmt.Printf("unexpected type %T\n", v)
	case bool:
		fmt.Println(v)
	case *bool:
		fmt.Println("bool pointer", v)
	}
}
