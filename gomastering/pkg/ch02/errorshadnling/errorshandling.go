package errorshadnling

import (
	"errors"
	"fmt"
	"os"
)

const (
	invalidZero = 0
)

func checkErrors(a, b int) error {

	if a == invalidZero && b == invalidZero {
		return errors.New("This is a custom error message")
	}

	return nil
}

func checkFmt(a, b int) error {
	if a == invalidZero && b == invalidZero {
		return fmt.Errorf("a %d and b %d. UserID: %d", a, b, os.Geteuid())
	}

	return nil
}

func Using() {
	if err := checkErrors(0, 0); err != nil {
		fmt.Println(err)
	}

	if err := checkFmt(0, 0); err != nil {
		fmt.Println(err)
	}

}
